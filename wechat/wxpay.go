package wechat

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"lemon-ipay-api/core"
	"lemon-ipay-api/model"
	"net/http"
	"time"

	"github.com/relax-space/go-kitt/random"

	"github.com/relax-space/lemon-wxmp-sdk/mpAuth"

	"github.com/relax-space/go-kit/base"
	"github.com/relax-space/go-kit/data"
	"github.com/relax-space/go-kit/log"
	"github.com/relax-space/go-kit/sign"

	wxpay "github.com/relax-space/lemon-wxpay-sdk"

	"github.com/labstack/echo"
	kmodel "github.com/relax-space/go-kit/model"
)

func Pay(c echo.Context) error {
	b, err := ioutil.ReadAll(c.Request().Body)
	defer c.Request().Body.Close()

	log.Info(fmt.Sprintf("Bind request param:%+v--err:%+v", string(b), err))
	var reqDto ReqPayDto
	if err := json.Unmarshal(b, &reqDto); err != nil {
		log.Error(fmt.Sprintf("err:%+v", err))
		return c.JSON(http.StatusBadRequest, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}

	// reqDto := ReqPayDto{}
	// if err := c.Bind(&reqDto); err != nil {
	// 	return c.JSON(http.StatusBadRequest, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	// }

	account, err := model.WxAccount{}.Get(reqDto.EId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}
	reqDto.ReqBaseDto = &wxpay.ReqBaseDto{
		AppId:    account.AppId,
		SubAppId: account.SubAppId,
		MchId:    account.MchId,
		SubMchId: account.SubMchId,
	}
	customDto := wxpay.ReqCustomerDto{
		Key: account.Key,
	}

	_, _, result, err := wxpay.Pay(reqDto.ReqPayDto, &customDto)
	if err != nil {
		if err.Error() == wxpay.MESSAGE_PAYING {
			outTradeNo := result["out_trade_no"].(string)
			queryDto := wxpay.ReqQueryDto{
				ReqBaseDto: reqDto.ReqBaseDto,
				OutTradeNo: outTradeNo,
			}
			_, _, result, err = wxpay.LoopQuery(&queryDto, &customDto, 40, 2)
			if err == nil {
				return c.JSON(http.StatusOK, kmodel.Result{Success: true, Result: result})
			} else {
				reverseDto := wxpay.ReqReverseDto{
					ReqBaseDto: reqDto.ReqBaseDto,
					OutTradeNo: outTradeNo,
				}
				if _, _, _, err = wxpay.Reverse(&reverseDto, &customDto, 10, 10); err != nil {
					return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
				} else {
					return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: "reverse sucess"}})
				}

			}
		} else {
			return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
		}
	}
	return c.JSON(http.StatusOK, kmodel.Result{Success: true, Result: result})
}

func Query(c echo.Context) error {
	reqDto := ReqQueryDto{}
	if err := c.Bind(&reqDto); err != nil {
		return c.JSON(http.StatusBadRequest, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}

	account, err := model.WxAccount{}.Get(reqDto.EId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}
	reqDto.ReqBaseDto = &wxpay.ReqBaseDto{
		AppId:    account.AppId,
		SubAppId: account.SubAppId,
		MchId:    account.MchId,
		SubMchId: account.SubMchId,
	}
	customDto := wxpay.ReqCustomerDto{
		Key: account.Key,
	}
	_, _, result, err := wxpay.Query(reqDto.ReqQueryDto, &customDto)
	if err != nil {
		return c.JSON(http.StatusOK, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}
	return c.JSON(http.StatusOK, kmodel.Result{Success: true, Result: result})
}
func Refund(c echo.Context) error {
	reqDto := ReqRefundDto{}
	if err := c.Bind(&reqDto); err != nil {
		return c.JSON(http.StatusBadRequest, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}
	account, err := model.WxAccount{}.Get(reqDto.EId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}
	reqDto.ReqBaseDto = &wxpay.ReqBaseDto{
		AppId:    account.AppId,
		SubAppId: account.SubAppId,
		MchId:    account.MchId,
		SubMchId: account.SubMchId,
	}
	custDto := wxpay.ReqCustomerDto{
		Key:          account.Key,
		CertPathName: account.CertName,
		CertPathKey:  account.CertKey,
		RootCa:       account.RootCa,
	}
	_, _, result, err := wxpay.Refund(reqDto.ReqRefundDto, &custDto)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})

	}
	return c.JSON(http.StatusOK, kmodel.Result{Success: true, Result: result})

}
func Reverse(c echo.Context) error {
	reqDto := ReqReverseDto{}
	if err := c.Bind(&reqDto); err != nil {
		return c.JSON(http.StatusBadRequest, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}
	account, err := model.WxAccount{}.Get(reqDto.EId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}
	reqDto.ReqBaseDto = &wxpay.ReqBaseDto{
		AppId:    account.AppId,
		SubAppId: account.SubAppId,
		MchId:    account.MchId,
		SubMchId: account.SubMchId,
	}
	custDto := wxpay.ReqCustomerDto{
		Key:          account.Key,
		CertPathName: account.CertName,
		CertPathKey:  account.CertKey,
		RootCa:       account.RootCa,
	}
	_, _, result, err := wxpay.Reverse(reqDto.ReqReverseDto, &custDto, 10, 10)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})

	}
	return c.JSON(http.StatusOK, kmodel.Result{Success: true, Result: result})
}

func RefundQuery(c echo.Context) error {
	reqDto := ReqRefundQueryDto{}
	if err := c.Bind(&reqDto); err != nil {
		return c.JSON(http.StatusBadRequest, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}

	account, err := model.WxAccount{}.Get(reqDto.EId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}
	reqDto.ReqBaseDto = &wxpay.ReqBaseDto{
		AppId:    account.AppId,
		SubAppId: account.SubAppId,
		MchId:    account.MchId,
		SubMchId: account.SubMchId,
	}
	customDto := wxpay.ReqCustomerDto{
		Key: account.Key,
	}
	_, _, result, err := wxpay.RefundQuery(reqDto.ReqRefundQueryDto, &customDto)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}
	return c.JSON(http.StatusOK, kmodel.Result{Success: true, Result: result})
}

func Prepay(c echo.Context) error {
	reqDto := ReqPrepayDto{}
	if err := c.Bind(&reqDto); err != nil {
		return c.JSON(http.StatusBadRequest, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}

	account, err := model.WxAccount{}.Get(reqDto.EId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}
	reqDto.ReqBaseDto = &wxpay.ReqBaseDto{
		AppId:    account.AppId,
		SubAppId: account.SubAppId,
		MchId:    account.MchId,
		SubMchId: account.SubMchId,
	}
	reqDto.TimeStart = base.GetDateFormat(ChinaDatetime(), 121)
	reqDto.TimeExpire = base.GetDateFormat(ChinaDatetime().Add(10*time.Minute), 121)
	customDto := wxpay.ReqCustomerDto{
		Key: account.Key,
	}
	/*
		1.set customer notify_url into attach
		2.set ipay united notify_url to reqDto
	*/
	reqDto.ReqPrepayDto.Attach = SetNotifyAttach(reqDto.NotifyUrl, reqDto.Attach, reqDto.EId)
	reqDto.ReqPrepayDto.NotifyUrl = fmt.Sprintf("%v/wx/%v", core.Env.HostUrl, "notify")

	reqDto.ReqPrepayDto.TimeStart = ChinaDatetime().Format("20060102150405")
	reqDto.ReqPrepayDto.TimeExpire = ChinaDatetime().Add(10 * time.Minute).Format("20060102150405")
	_, _, result, err := wxpay.Prepay(reqDto.ReqPrepayDto, &customDto)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}

	prePayParam := make(map[string]interface{}, 0)
	prePayParam["package"] = "prepay_id=" + base.ToString(result["prepay_id"])
	prePayParam["timeStamp"] = base.ToString(ChinaDatetime().Unix())
	prePayParam["nonceStr"] = result["nonce_str"]
	prePayParam["signType"] = "MD5"
	prePayParam["appId"] = result["appid"]
	prePayParam["paySign"] = sign.MakeMd5Sign(base.JoinMapObject(prePayParam), account.Key)

	return c.JSON(http.StatusOK, kmodel.Result{Success: true, Result: prePayParam})
}

func Notify(c echo.Context) error {
	fmt.Printf("\n%v-%v", time.Now(), "wx notify received.")

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return NotifyError(c, err.Error())
	}
	xmlBody := string(body)
	if len(xmlBody) == 0 {
		return NotifyError(c, "xml is empty")
	}
	fmt.Println(string(xmlBody))
	//1.get dto data
	var notifyDto model.NotifyWechat
	err = xml.Unmarshal([]byte(xmlBody), &notifyDto)
	if err != nil {
		return NotifyError(c, err.Error())
	}
	//1.1 get mapData
	wxData := data.New()
	err = wxData.FromXml(xmlBody)
	if err != nil {
		return NotifyError(c, err.Error())
	}
	//2.valid
	if err = NotifyValid(notifyDto.Attach, notifyDto.Sign, notifyDto.OutTradeNo, notifyDto.TotalFee, wxData); err != nil {
		return NotifyError(c, err.Error())

	}

	//3.save into data base
	err = model.NotifyWechat{}.InsertOne(&notifyDto)
	if err != nil {
		return NotifyError(c, err.Error())
	}

	successResult := struct {
		XMLName    xml.Name `xml:"xml"`
		ReturnCode string   `xml:"return_code"`
		ReturnMsg  string   `xml:"return_msg"`
	}{xml.Name{}, "SUCCESS", "OK"}
	return c.XML(http.StatusOK, successResult)
}

/*
redirect to "https:/xxxx/v3/prepayopenid" for get openid
*/
func PrepayEasy(c echo.Context) error {

	prepay_param := c.QueryParam("prepay_param")

	reqDto := ReqPrepayEasyDto{}
	err := json.Unmarshal([]byte(prepay_param), &reqDto)
	if err != nil {
		return prepayError(c, "request param format is not right")
	}
	urlStr, err := prepayPageUrl(reqDto.PageUrl)
	if err != nil {
		return prepayError(c, err.Error())
	}
	reqDto.PageUrl = fmt.Sprintf(urlStr, random.Uuid("")) + "?" + random.Uuid("")
	account, err := model.WxAccount{}.Get(reqDto.EId)
	if err != nil {
		return prepayErrorDirect(c, reqDto.PageUrl, err.Error())
	}

	openIdUrlParam := &mpAuth.ReqDto{
		AppId:       account.AppId,
		State:       "state",
		RedirectUrl: fmt.Sprintf("%v/wx/%v", core.Env.HostUrl, "prepayopenid"),
		PageUrl:     reqDto.PageUrl,
	}
	SetCookie(IPAY_WECHAT_PREPAY_INNER, prepay_param, c)
	SetCookie(IPAY_WECHAT_PREPAY, "", c)
	SetCookie(IPAY_WECHAT_PREPAY_ERROR, "", c)
	return c.Redirect(http.StatusFound, mpAuth.GetUrlForAccessToken(openIdUrlParam))
}
