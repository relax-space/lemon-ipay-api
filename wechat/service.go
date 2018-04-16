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
	"github.com/relax-space/go-kit/data"
	"github.com/relax-space/go-kit/log"
	"github.com/relax-space/go-kit/sign"
	"github.com/relax-space/go-kitt/auth"
	"github.com/relax-space/lemon-wxmp-sdk/mpAuth"
	paysdk "github.com/relax-space/lemon-wxpay-sdk"
)

func NotifyQuery(account *model.WxAccount, outTradeNo string) (result map[string]interface{}, err error) {
	var reqDto paysdk.ReqQueryDto
	reqDto.ReqBaseDto = &paysdk.ReqBaseDto{
		AppId:    account.AppId,
		SubAppId: account.SubAppId,
		MchId:    account.MchId,
		SubMchId: account.SubMchId,
	}
	customDto := &paysdk.ReqCustomerDto{
		Key: account.Key,
	}
	reqDto.OutTradeNo = outTradeNo
	_, _, result, err = paysdk.Query(&reqDto, customDto)
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

func NotifyValid(body, signParam, outTradeNo string, totalAmount int64, dataParam *data.Data) (err error) {
	subNotifyUrl, rawAttach, eId, err := GetNotifyAttach(body)
	if err != nil {
		return
	}
	account, err := model.WxAccount{}.Get(eId)
	if err != nil {
		return
	}

	//1.valid sign
	signStr := signParam
	mapParam := dataParam.DataAttr
	delete(mapParam, "sign")
	if !sign.CheckMd5Sign(base.JoinMapObject(mapParam), account.Key, signStr) {
		err = errors.New("sign valid failure")
		return
	}
	mapParam["attach"] = rawAttach
	//2.valid data
	queryMap, err := NotifyQuery(account, outTradeNo)
	if err != nil {
		return
	}
	if !(queryMap["total_fee"].(string) == base.ToString(totalAmount)) {
		err = errors.New("amount is exception")
		return
	}
	mapParam["sign"] = signParam
	//3.send data to sub_mch
	if len(subNotifyUrl) != 0 {
		go func(signParam, subNotifyUrl string, dataParam *data.Data) {
			SubNotify(subNotifyUrl, dataParam.ToXml())
		}(signStr, subNotifyUrl, dataParam)
	}
	return
}

type SuccessResult struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
}

func SubNotify(subNotifyUrl, xmlParam string) (successResult SuccessResult) {
	token, err := getToken()
	if err != nil {
		return
	}
	resp, err := POSTXml(token, subNotifyUrl, xmlParam, &successResult)
	if err == nil && resp != nil &&
		resp.StatusCode == http.StatusOK && successResult.ReturnCode == "SUCCESS" {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "notify success", xmlParam)
		fmt.Println(xmlParam)
		return
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
	prePayParam, err := PrepayRespParam(reqDto, account)
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
	reqDto.ReqBaseDto = &paysdk.ReqBaseDto{
		AppId:    account.AppId,
		SubAppId: account.SubAppId,
		MchId:    account.MchId,
		SubMchId: account.SubMchId,
	}
	customDto := paysdk.ReqCustomerDto{
		Key: account.Key,
	}
	_, _, result, err := paysdk.Prepay(reqDto.ReqPrepayDto, &customDto)
	if err != nil {
		return
	}

	prePayParam = make(map[string]interface{}, 0)
	prePayParam["package"] = "prepay_id=" + base.ToString(result["prepay_id"])
	prePayParam["timeStamp"] = base.ToString(ChinaDatetime().Unix())
	prePayParam["nonceStr"] = result["nonce_str"]
	prePayParam["signType"] = "MD5"
	prePayParam["appId"] = result["appid"]
	prePayParam["paySign"] = sign.MakeMd5Sign(base.JoinMapObject(prePayParam), account.Key)
	prePayParam["jwtToken"], _ = auth.NewToken(map[string]interface{}{"type": "ticket"})
	return
}

/*
e_id,sub_notify_url,attach
*/
func SetNotifyAttach(subNotifyUrl, attach string, eId int64) (newAttach string) {
	newAttach = strconv.FormatInt(eId, 10) + core.NOTIFY_BODY_SEP2 +
		url.QueryEscape(subNotifyUrl) + core.NOTIFY_BODY_SEP2 +
		url.QueryEscape(attach)
	return
}

/*
e_id,sub_notify_url,attach
*/
func GetNotifyAttach(attach string) (subNotifyUrl, rawAttach string, eId int64, err error) {
	vs := strings.Split(attach, core.NOTIFY_BODY_SEP2)
	if len(vs) != 3 {
		err = errors.New("attach param is missing[e_id,sub_notify_url,attach]")
		return
	}
	eId, err = strconv.ParseInt(vs[0], 10, 64)
	if err != nil {
		return
	}
	subNotifyUrl, err = url.QueryUnescape(vs[1])
	rawAttach, err = url.QueryUnescape(vs[2])
	return
}
