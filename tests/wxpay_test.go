package tests

import (
	"encoding/json"
	"fmt"
	"kit/test"
	"lemon-ipay-api/wechat"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/relax-space/go-kit/model"
)

func Test_WxPay(t *testing.T) {
	bodyStr := `
	{
		"e_id":10001,
		"auth_code":"135298324463700425",
		"body":"xiaoxinmiao test",
		"total_fee":1
	}`
	req, err := http.NewRequest(echo.POST, "/v3/wx/pay", strings.NewReader(bodyStr))
	test.Ok(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	test.Ok(t, wechat.Pay(c))
	v := model.Result{}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	fmt.Printf("\n%+v", v)
	test.Equals(t, http.StatusOK, rec.Code)

}

func Test_WxRefund(t *testing.T) {
	bodyStr := `
	{
		"e_id":10001,
		"out_trade_no":"147688874645492354650",
		"refund_fee":1
	}`
	req, err := http.NewRequest(echo.POST, "/v3/wx/refund", strings.NewReader(bodyStr))
	test.Ok(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	test.Ok(t, wechat.Refund(c))
	v := model.Result{}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	fmt.Printf("\n%+v", v)
	test.Equals(t, http.StatusOK, rec.Code)

}

func Test_WxReverse(t *testing.T) {
	bodyStr := `
	{
		"e_id":10001,
		"out_trade_no":"143420620288156126697"
	}`
	req, err := http.NewRequest(echo.POST, "/v3/wx/reverse", strings.NewReader(bodyStr))
	test.Ok(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	test.Ok(t, wechat.Reverse(c))
	v := model.Result{}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	fmt.Printf("\n%+v", v)
	test.Equals(t, http.StatusOK, rec.Code)

}

func Test_WxQuery(t *testing.T) {
	bodyStr := `
	{
		"e_id":10001,
		"out_trade_no":"14201711085205823413229775520"
	}`
	req, err := http.NewRequest(echo.POST, "/v3/wx/query", strings.NewReader(bodyStr))
	test.Ok(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	test.Ok(t, wechat.Query(c))
	v := model.Result{}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	fmt.Printf("\n%+v", v)
	test.Equals(t, http.StatusOK, rec.Code)

}

func Test_WxRefundQuery(t *testing.T) {
	bodyStr := `
	{
		"e_id":10001,
		"out_trade_no":"144650782494807835413"
	}`
	req, err := http.NewRequest(echo.POST, "/v3/wx/refundquery", strings.NewReader(bodyStr))
	test.Ok(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	test.Ok(t, wechat.RefundQuery(c))
	v := model.Result{}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	fmt.Printf("\n%+v", v)
	test.Equals(t, http.StatusOK, rec.Code)

}

func Test_WxPrepay(t *testing.T) {
	bodyStr := `
	{
		"e_id":10001,
		"body":"xiaomiao test",
		"total_fee":1,
		"trade_type":"JSAPI",
		"notify_url":"http://xiao.xinmiao.com",
		"openid":"os2u9uPKLkCKL08FwCM6hQAQ_LtI"
	}`
	req, err := http.NewRequest(echo.POST, "/v3/wx/Prepay", strings.NewReader(bodyStr))
	test.Ok(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	test.Ok(t, wechat.Prepay(c))
	v := model.Result{}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	fmt.Printf("\n%+v", v)
	test.Equals(t, http.StatusOK, rec.Code)

}

func Test_WxNotify(t *testing.T) {
	xmlBody := `<xml>
	<appid><![CDATA[wx2421b1c4370ec43b]]></appid>
	<attach><![CDATA[{"sub_notify_url":"https://baidu.com","e_id":10001}]]></attach>
	<bank_type><![CDATA[CFT]]></bank_type>
	<fee_type><![CDATA[CNY]]></fee_type>
	<is_subscribe><![CDATA[Y]]></is_subscribe>
	<mch_id><![CDATA[10000100]]></mch_id>
	<nonce_str><![CDATA[5d2b6c2a8db53831f7eda20af46e531c]]></nonce_str>
	<openid><![CDATA[oUpF8uMEb4qRXf22hE3X68TekukE]]></openid>
	<out_trade_no><![CDATA[1409811653]]></out_trade_no>
	<result_code><![CDATA[SUCCESS]]></result_code>
	<return_code><![CDATA[SUCCESS]]></return_code>
	<sign><![CDATA[7D24E7B803ED7574785872A50105046D]]></sign>
	<sub_mch_id><![CDATA[10000100]]></sub_mch_id>
	<time_end><![CDATA[20140903131540]]></time_end>
	<total_fee>1</total_fee>
	<trade_type><![CDATA[JSAPI]]></trade_type>
	<transaction_id><![CDATA[B2AE05C99B9C81A640472406AA3C2710]]></transaction_id>
 </xml>`
	req, err := http.NewRequest(echo.POST, "/v3/wx/notify", strings.NewReader(xmlBody))
	test.Ok(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	test.Ok(t, wechat.Notify(c))
	fmt.Printf("\n%+v", string(rec.Body.Bytes()))
	test.Equals(t, http.StatusOK, rec.Code)

}

func Test_WxNotify2(t *testing.T) {
	xmlBody := `<xml><appid><![CDATA[wx856df5e42a345096]]></appid>
	<attach><![CDATA[e_id||||10001&sub_notify_url||||https://baidu.com]]></attach>
	<bank_type><![CDATA[CMB_CREDIT]]></bank_type>
	<cash_fee><![CDATA[1]]></cash_fee>
	<fee_type><![CDATA[CNY]]></fee_type>
	<is_subscribe><![CDATA[Y]]></is_subscribe>
	<mch_id><![CDATA[1294997801]]></mch_id>
	<nonce_str><![CDATA[1240648768328515708]]></nonce_str>
	<openid><![CDATA[os2u9uBHeJRPtCkisjVf-kWZWjKQ]]></openid>
	<out_trade_no><![CDATA[169126120915612414200792892307]]></out_trade_no>
	<result_code><![CDATA[SUCCESS]]></result_code>
	<return_code><![CDATA[SUCCESS]]></return_code>
	<sign><![CDATA[DACC65EAD9461590F693C08EAB2F0A10]]></sign>
	<sub_appid><![CDATA[wx38db2bfbb79a3cea]]></sub_appid>
	<sub_is_subscribe><![CDATA[Y]]></sub_is_subscribe>
	<sub_mch_id><![CDATA[1464381802]]></sub_mch_id>
	<sub_openid><![CDATA[o2-sBj3ozQQ6gxiyYKI2JzJFcUhY]]></sub_openid>
	<time_end><![CDATA[20171209235456]]></time_end>
	<total_fee>1</total_fee>
	<trade_type><![CDATA[JSAPI]]></trade_type>
	<transaction_id><![CDATA[4200000031201712091054287297]]></transaction_id>
	</xml>`
	req, err := http.NewRequest(echo.POST, "/v3/wx/notify", strings.NewReader(xmlBody))
	test.Ok(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	test.Ok(t, wechat.Notify(c))
	fmt.Printf("\n%+v", string(rec.Body.Bytes()))
	test.Equals(t, http.StatusOK, rec.Code)

}

func TestPing(t *testing.T) {
	req, err := http.NewRequest(echo.GET, "/ping", nil)
	test.Ok(t, err)
	rec := httptest.NewRecorder()
	echo.New().NewContext(req, rec)
	test.Equals(t, http.StatusOK, rec.Code)
}

func Test_WxPrepayEasy(t *testing.T) {
	result := "https://ipay.p2shop.cn/#/pay"
	indexTag := strings.Index(result, "#")
	result = result[0:indexTag] + "%v" + result[indexTag:]
	fmt.Println(result)

	/*
		localhost:5000/v3/wx/prepayeasy?app_id=&page_url=ttps%3A%2F%2Fgateway.p2shop.cn%2Fipay%2Fping&prepay_param={"e_id":10001,"body":"xiaomiao test","total_fee":1,"trade_type":"JSAPI","notify_url":"http://xiao.xinmiao.com"}
	*/
	q := make(url.Values)
	q.Set("prepay_param", `{
		"page_url":"https://ipay.p2shop.cn/#/pay",
		"attach":"e_id||||10001",
		"body":"xiaomiao test",
		"total_fee":1,
		"trade_type":"JSAPI",
		"notify_url":"https://xiao.xinmiao.com"
	}`)

	req, err := http.NewRequest(echo.GET, "/v3/wx/prepayeasy?"+q.Encode(), nil)
	test.Ok(t, err)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	test.Ok(t, wechat.PrepayEasy(c))
	fmt.Printf("\n%+v", rec.HeaderMap)
	test.Equals(t, http.StatusFound, rec.Code)
	//cookie:IPAY_WECHAT_PREPAY=%7B%22body%22%3A%22xiaomiao+test%22%2C%22total_fee%22%3A1%2C%22notify_url%22%3A%22http%3A%2F%2Fxiao.xinmiao.com%22%2C%22trade_type%22%3A%22JSAPI%22%2C%22scene_info%22%3A%7B%7D%2C%22e_id%22%3A10001%7D; Path=/
}

func Test_WxPrepayOpenId(t *testing.T) {
	path := fmt.Sprintf("/v3/wx/prepayopenid?code=%v&reurl=%v", "081YsUt32WVM9M0y1Nv32BvXt32YsUt-", "https%3A%2F%2Fgateway.p2shop.cn%2Fipay%2Fping")
	req, err := http.NewRequest(echo.GET, path, nil)
	test.Ok(t, err)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	cookie := new(http.Cookie)
	cookie.Name = wechat.IPAY_WECHAT_PREPAY_INNER
	value := url.QueryEscape(`{
		"page_url":"https://ipay.p2shop.cn/#/pay",
		"attach":"e_id||||10001",
		"e_id":10001,
		"body":"xiaomiao test",
		"total_fee":1,
		"trade_type":"JSAPI",
		"notify_url":"http://xiao.xinmiao.com"
	}`)
	cookie.Value = value
	cookie.Path = "/"
	c.SetCookie(cookie)
	test.Ok(t, wechat.PrepayOpenId(c))
	fmt.Printf("\n%+v", rec.HeaderMap)
	test.Equals(t, http.StatusFound, rec.Code)

}

func Test_NotifyBodyParse(t *testing.T) {
	body := "e_id||||100&sub_notify_url||||https://www.baidu.com"
	bodyMap, eId, err := wechat.NotifyBodyParse(body)
	fmt.Printf("\n%+v", bodyMap)
	fmt.Printf("\n%T,%v", eId, eId)
	test.Ok(t, err)
}
