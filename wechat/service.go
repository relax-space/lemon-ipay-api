package wechat

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"lemon-ipay-api/core"
	"lemon-ipay-api/model"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/relax-space/go-kit/base"
	"github.com/relax-space/go-kit/httpreq"
	"github.com/relax-space/go-kit/log"
	"github.com/relax-space/go-kit/sign"
	"github.com/relax-space/go-kitt/auth"
	"github.com/relax-space/lemon-wxmp-sdk/mpAuth"
	paysdk "github.com/relax-space/lemon-wxpay-sdk"
	wxpay "github.com/relax-space/lemon-wxpay-sdk"
)

func NotifyQuery(account *model.WxAccount, outTradeNo string) (result map[string]interface{}, err error) {
	var reqDto paysdk.ReqQueryDto
	reqDto.ReqBaseDto = &wxpay.ReqBaseDto{
		AppId:    account.AppId,
		SubAppId: account.SubAppId,
		MchId:    account.MchId,
		SubMchId: account.SubMchId,
	}
	customDto := &wxpay.ReqCustomerDto{
		Key: account.Key,
	}
	reqDto.OutTradeNo = outTradeNo
	result, err = paysdk.Query(&reqDto, customDto)
	return
}

func NotifyBodyParse(body string) (bodyMap map[string]interface{}, eId int64, err error) {
	bodyMap = base.ParseMapObject(body, core.NOTIFY_BODY_SEP1, core.NOTIFY_BODY_SEP2)
	err = errors.New("e_id is not existed in body or format is not correct")
	eIdObj, ok := bodyMap["e_id"]
	if !ok {
		return
	}
	if eId, err = strconv.ParseInt(eIdObj.(string), 10, 64); err != nil {
		return
	}
	err = nil
	return
}

func NotifyValid(body, signParam, outTradeNo string, totalAmount int64, mapParam map[string]interface{}) (err error) {
	bodyMap, eId, err := NotifyBodyParse(body)
	if err != nil {
		return
	}
	account, err := model.WxAccount{}.Get(eId)
	if err != nil {
		return
	}

	//1.valid sign
	signStr := signParam
	delete(mapParam, "sign")
	fmt.Println(base.JoinMapObject(mapParam))

	if !sign.CheckMd5Sign(base.JoinMapObject(mapParam), account.Key, signStr) {
		err = errors.New("sign valid failure")
		return
	}
	//2.valid data
	queryMap, err := NotifyQuery(&account, outTradeNo)
	if err != nil {
		return
	}
	if !(queryMap["total_fee"].(string) == base.ToString(totalAmount)) {
		err = errors.New("amount is exception")
		return
	}
	//3.send data to sub_mch
	if subNotifyUrl, ok := bodyMap["sub_notify_url"]; ok {
		go func(signParam string) {
			mapParam["sign"] = signParam
			_, err = httpreq.POST("", subNotifyUrl.(string), mapParam, nil)
		}(signStr)
	}
	return
}

func NotifyError(c echo.Context, errMsg string) error {
	errResult := struct {
		XMLName    xml.Name `xml:"xml"`
		ReturnCode string   `xml:"return_code"`
		ReturnMsg  string   `xml:"return_msg"`
	}{xml.Name{}, "FAIL", ""}
	errResult.ReturnMsg = errMsg
	log.Error(errMsg)
	return c.XML(http.StatusBadRequest, errResult)

}
func prepayError(c echo.Context, errMsg string) error {
	prepayErrCookie(c, errMsg)
	return c.String(http.StatusBadRequest, errMsg)
}
func prepayErrorDirect(c echo.Context, reqUrl, errMsg string) error {
	prepayErrCookie(c, errMsg)
	return c.Redirect(http.StatusFound, reqUrl)
}
func prepayErrCookie(c echo.Context, errMsg string) {
	SetCookie(IPAY_WECHAT_PREPAY_ERROR, errMsg, c)
	SetCookie(IPAY_WECHAT_PREPAY_INNER, "", c)
	SetCookie(IPAY_WECHAT_PREPAY, "", c)
}
func prepayPageUrl(pageUrl string) (result string, err error) {
	result, err = url.QueryUnescape(pageUrl)
	if err != nil {
		return
	}
	if len(result) == 0 {
		err = errors.New("page_url miss")
		return
	}
	indexTag := strings.Index(result, "#")
	result = result[0:indexTag] + "%v?" + result[indexTag:]
	return
}

/*
1.get openid
2.get prepay param
*/
func PrepayOpenId(c echo.Context) error {
	code := c.QueryParam("code")
	reqUrl := c.QueryParam("reurl")
	reqDto, err := PrepayReqParam(c)
	if err != nil {
		return prepayErrorDirect(c, reqUrl, err.Error())
	}
	//1.get account
	account, err := model.WxAccount{}.Get(reqDto.EId)
	if err != nil {
		return prepayErrorDirect(c, reqUrl, err.Error())
	}
	//2.get openId
	respDto, err := mpAuth.GetAccessTokenAndOpenId(code, account.AppId, account.Secret)
	if err != nil {
		return prepayErrorDirect(c, reqUrl, err.Error())
	}
	reqDto.OpenId = respDto.OpenId
	//3.get prepay param
	prePayParam, err := PrepayRespParam(reqDto, &account)
	if err != nil {
		return prepayErrorDirect(c, reqUrl, err.Error())
	}
	SetCookieObj(IPAY_WECHAT_PREPAY, prePayParam, c)
	SetCookie(IPAY_WECHAT_PREPAY_ERROR, "", c)
	SetCookie(IPAY_WECHAT_PREPAY_INNER, "", c)
	return c.Redirect(http.StatusFound, reqUrl)

}

func PrepayReqParam(c echo.Context) (reqDto *ReqPrepayEasyDto, err error) {
	cookie, err := c.Cookie(IPAY_WECHAT_PREPAY_INNER)
	if err != nil {
		return
	}
	param, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return
	}
	// param := "%7B%0A%09%09%22page_url%22%3A%22https%3A%2F%2Fipay.p2shop.cn%2F%23%2Fpay%22%2C%0A%09%09%22attach%22%3A%22e_id%7C%7C%7C%7C10001%22%2C%0A%09%09%22e_id%22%3A10001%2C%0A%09%09%22body%22%3A%22xiaomiao+test%22%2C%0A%09%09%22total_fee%22%3A1%2C%0A%09%09%22trade_type%22%3A%22JSAPI%22%2C%0A%09%09%22notify_url%22%3A%22http%3A%2F%2Fxiao.xinmiao.com%22%0A%09%7D"
	// param, _ = url.QueryUnescape(param)
	reqDto = &ReqPrepayEasyDto{}
	err = json.Unmarshal([]byte(param), reqDto)
	if err != nil {
		return
	}
	return
}

func PrepayRespParam(reqDto *ReqPrepayEasyDto, account *model.WxAccount) (prePayParam map[string]interface{}, err error) {
	reqDto.ReqBaseDto = &wxpay.ReqBaseDto{
		AppId:    account.AppId,
		SubAppId: account.SubAppId,
		MchId:    account.MchId,
		SubMchId: account.SubMchId,
	}
	customDto := wxpay.ReqCustomerDto{
		Key: account.Key,
	}
	result, err := wxpay.Prepay(reqDto.ReqPrepayDto, &customDto)
	if err != nil {
		return
	}

	prePayParam = make(map[string]interface{}, 0)
	prePayParam["package"] = "prepay_id=" + base.ToString(result["prepay_id"])
	prePayParam["timeStamp"] = base.ToString(time.Now().Unix())
	prePayParam["nonceStr"] = result["nonce_str"]
	prePayParam["signType"] = "MD5"
	prePayParam["appId"] = result["appid"]
	prePayParam["pay_sign"] = sign.MakeMd5Sign(base.JoinMapObject(prePayParam), account.Key)
	prePayParam["jwtToken"], _ = auth.NewToken(map[string]interface{}{"type": "ticket"})
	return
}
