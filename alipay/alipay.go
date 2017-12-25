package alipay

import (
	"fmt"
	"io/ioutil"
	"lemon-ipay-api/core"
	"lemon-ipay-api/model"
	"net/http"
	"time"

	"github.com/relax-space/go-kit/log"

	paysdk "github.com/relax-space/lemon-alipay-sdk"

	"github.com/labstack/echo"
	"github.com/relax-space/go-kit/base"
	kmodel "github.com/relax-space/go-kit/model"
)

func Pay(c echo.Context) error {
	reqDto := ReqPayDto{}
	if err := c.Bind(&reqDto); err != nil {
		return c.JSON(http.StatusBadRequest, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}

	account, err := model.AlAccount{}.Get(reqDto.EId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}
	reqDto.ReqBaseDto = &paysdk.ReqBaseDto{
		AppId:        account.AppId,
		AppAuthToken: account.AuthToken,
	}
	if len(account.SysServiceProviderId) != 0 {
		reqDto.ExtendParams = &paysdk.ExtendParams{
			SysServiceProviderId: account.SysServiceProviderId,
		}
	}
	customDto := &paysdk.ReqCustomerDto{
		PriKey: account.PriKey,
		PubKey: account.PubKey,
	}

	result, err := paysdk.Pay(reqDto.ReqPayDto, customDto)
	if err != nil {
		if err.Error() == paysdk.MESSAGE_PAYING {
			outTradeNo := result.OutTradeNo
			queryDto := paysdk.ReqQueryDto{
				ReqBaseDto: reqDto.ReqBaseDto,
				OutTradeNo: outTradeNo,
			}
			result, err = paysdk.LoopQuery(&queryDto, customDto, 40, 2)
			if err == nil {
				return c.JSON(http.StatusOK, kmodel.Result{Success: true, Result: result})
			} else {
				reverseDto := paysdk.ReqReverseDto{
					ReqBaseDto: reqDto.ReqBaseDto,
					OutTradeNo: outTradeNo,
				}
				if _, err = paysdk.Reverse(&reverseDto, customDto, 10, 10); err != nil {
					return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
				} else {
					return c.JSON(http.StatusOK, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: "reverse success"}})
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

	account, err := model.AlAccount{}.Get(reqDto.EId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}
	reqDto.ReqBaseDto = &paysdk.ReqBaseDto{
		AppId:        account.AppId,
		AppAuthToken: account.AuthToken,
	}

	customDto := &paysdk.ReqCustomerDto{
		PriKey: account.PriKey,
		PubKey: account.PubKey,
	}
	result, err := paysdk.Query(reqDto.ReqQueryDto, customDto)
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
	account, err := model.AlAccount{}.Get(reqDto.EId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}
	reqDto.ReqBaseDto = &paysdk.ReqBaseDto{
		AppId:        account.AppId,
		AppAuthToken: account.AuthToken,
	}

	customDto := &paysdk.ReqCustomerDto{
		PriKey: account.PriKey,
		PubKey: account.PubKey,
	}
	result, err := paysdk.Refund(reqDto.ReqRefundDto, customDto)
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
	account, err := model.AlAccount{}.Get(reqDto.EId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}
	reqDto.ReqBaseDto = &paysdk.ReqBaseDto{
		AppId:        account.AppId,
		AppAuthToken: account.AuthToken,
	}

	customDto := &paysdk.ReqCustomerDto{
		PriKey: account.PriKey,
		PubKey: account.PubKey,
	}
	result, err := paysdk.Reverse(reqDto.ReqReverseDto, customDto, 10, 10)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})

	}
	return c.JSON(http.StatusOK, kmodel.Result{Success: true, Result: result})
}

// func RefundQuery(c echo.Context) error {
// 	reqDto := ReqRefundQueryDto{}
// 	if err := c.Bind(&reqDto); err != nil {
// 		return c.JSON(http.StatusBadRequest, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
// 	}

// 	account, err := model.AlAccount{}.Get(reqDto.EId)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
// 	}
// 	reqDto.ReqBaseDto = &paysdk.ReqBaseDto{
// 		AppId:        account.AppId,
// 		AppAuthToken: account.AuthToken,
// 	}

// 	customDto := &paysdk.ReqCustomerDto{
// 		PriKey: account.PriKey,
// 		PubKey: account.PubKey,
// 	}
// 	result, err := paysdk.RefundQuery(reqDto.ReqRefundQueryDto, customDto)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
// 	}
// 	return c.JSON(http.StatusOK, kmodel.Result{Success: true, Result: result})
// }

func Prepay(c echo.Context) error {
	reqDto := ReqPrepayDto{}
	if err := c.Bind(&reqDto); err != nil {
		return c.JSON(http.StatusBadRequest, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}

	account, err := model.AlAccount{}.Get(reqDto.EId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}
	reqDto.ReqBaseDto = &paysdk.ReqBaseDto{
		AppId:        account.AppId,
		AppAuthToken: account.AuthToken,
	}
	if len(account.SysServiceProviderId) != 0 {
		reqDto.ExtendParams = &paysdk.ExtendParams{
			SysServiceProviderId: account.SysServiceProviderId,
		}
	}
	customDto := &paysdk.ReqCustomerDto{
		PriKey: account.PriKey,
		PubKey: account.PubKey,
	}
	result, err := paysdk.Prepay(reqDto.ReqPrepayDto, customDto)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, kmodel.Result{Success: false, Error: kmodel.Error{Code: 10004, Message: err.Error()}})
	}
	return c.JSON(http.StatusOK, kmodel.Result{Success: true, Result: result})
}

func Notify(c echo.Context) error {
	fmt.Printf("\n%v-%v", time.Now(), "al notify received.")
	sbody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusBadRequest, "failure")
	}
	formParam := string(sbody)
	if len(formParam) == 0 {
		log.Error("param is empty")
		return c.String(http.StatusBadRequest, "failure")
	}
	var reqDto model.NotifyAlipay
	mapParam := base.ParseMapObjectEncode(formParam, "&", "=")
	err = core.Decode(mapParam, &reqDto)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusBadRequest, "failure")
	}

	//1.validate
	if err = NotifyValidN(reqDto.Body, reqDto.Sign, reqDto.OutTradeNo, reqDto.TotalAmount, mapParam); err != nil {
		log.Error(err)
		return c.String(http.StatusBadRequest, "failure")
	}

	//2.save notify info
	err = model.NotifyAlipay{}.InsertOne(&reqDto)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "failure")
	}
	return c.String(http.StatusOK, "success")
}
