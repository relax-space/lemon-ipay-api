package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/relax-space/go-kit/test"

	"lemon-ipay-api/alipay"

	"github.com/labstack/echo"
	"github.com/relax-space/go-kit/model"
)

func Test_AliPay(t *testing.T) {
	bodyStr := `
	{
		"e_id":10001,
		"auth_code":"283209675485586567",
		"subject":"xiaomiao test apilay",
		"total_amount":0.01
	}`
	req, err := http.NewRequest(echo.POST, "/v3/al/pay", strings.NewReader(bodyStr))
	test.Ok(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	test.Ok(t, alipay.Pay(c))
	v := model.Result{}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	fmt.Printf("\n%+v", v)
	test.Equals(t, http.StatusOK, rec.Code)

}

func Test_AliRefund(t *testing.T) {
	bodyStr := `
	{
		"e_id":10001,
		"out_trade_no":"1117112912739763007486053235",
		"refund_amount":0.01
	}`
	req, err := http.NewRequest(echo.POST, "/v3/wx/refund", strings.NewReader(bodyStr))
	test.Ok(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	test.Ok(t, alipay.Refund(c))
	v := model.Result{}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	fmt.Printf("\n%+v", v)
	test.Equals(t, http.StatusOK, rec.Code)

}

func Test_AliQuery(t *testing.T) {
	bodyStr := `
	{
		"e_id":10001,
		"out_trade_no":"1117112912739763007486053235"
	}`
	req, err := http.NewRequest(echo.POST, "/v3/wx/query", strings.NewReader(bodyStr))
	test.Ok(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	test.Ok(t, alipay.Query(c))
	v := model.Result{}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	fmt.Printf("\n%+v", v)
	test.Equals(t, http.StatusOK, rec.Code)

}

func Test_AliReverse(t *testing.T) {
	bodyStr := `
	{
		"e_id":10001,
		"out_trade_no":"1117112912739763007486053235"
	}`
	req, err := http.NewRequest(echo.POST, "/v3/wx/reverse", strings.NewReader(bodyStr))
	test.Ok(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	test.Ok(t, alipay.Reverse(c))
	v := model.Result{}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	fmt.Printf("\n%+v", v)
	test.Equals(t, http.StatusOK, rec.Code)

}

func Test_AliPrepay(t *testing.T) {
	bodyStr := `
	{
		"e_id":10001,
		"subject":"xiaomiao test ali",
		"total_amount":0.01
	}`
	req, err := http.NewRequest(echo.POST, "/v3/wx/prepay", strings.NewReader(bodyStr))
	test.Ok(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	test.Ok(t, alipay.Prepay(c))
	v := model.Result{}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	fmt.Printf("\n%+v", v)
	test.Equals(t, http.StatusOK, rec.Code)

}

func Test_AliNotify(t *testing.T) {
	bodyStr := `gmt_create=2017-12-07+11%3A15%3A39&amp;charset=UTF-8&amp;seller_email=eland_pay%40elandsystems.cn&amp;subject=xiaomiao+test+ali&amp;sign=***&amp;buyer_id=2088702305824122&amp;invoice_amount=0.01&amp;notify_id=50b1bbc78907f7d891e14f3209ecde0gxe&amp;fund_bill_list=%5B%7B%22amount%22%3A%220.01%22%2C%22fundChannel%22%3A%22PCREDIT%22%7D%5D&amp;notify_type=trade_status_sync&amp;trade_status=TRADE_SUCCESS&amp;receipt_amount=0.01&amp;buyer_pay_amount=0.01&amp;app_id=2015081700219294&amp;sign_type=RSA&amp;seller_id=2088312582701209&amp;gmt_payment=2017-12-07+11%3A15%3A43&amp;notify_time=2017-12-07+11%3A15%3A52&amp;version=1.0&amp;out_trade_no=131712072074192779392994717&amp;total_amount=0.01&amp;trade_no=2017120721001004120213770419&amp;auth_app_id=2015081700219294&amp;buyer_logon_id=xia***%40163.com&amp;point_amount=0.00`
	req, err := http.NewRequest(echo.POST, "/v3/wx/notify", strings.NewReader(bodyStr))
	test.Ok(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	test.Ok(t, alipay.Notify(c))
	fmt.Printf("\n%+v", string(rec.Body.Bytes()))
	test.Equals(t, http.StatusOK, rec.Code)

}
