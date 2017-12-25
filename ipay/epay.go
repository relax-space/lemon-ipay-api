package ipay

import (
	"net/http"

	"github.com/relax-space/go-kit/model"

	"github.com/labstack/echo"
)

func Pay(c echo.Context) error {
	reqDto := new(PayReqDto)
	if err := c.Bind(reqDto); err != nil {
		return c.JSON(http.StatusBadRequest, model.Result{Success: false, Error: model.Error{Message: err.Error()}})
	}
	//1.get payType
	//2.get eId by serviceType
	//3.get account by eId
	//4.

	return nil
}

func Query(c echo.Context) error {
	return nil
}

func Reverse(c echo.Context) error {
	return nil
}

func Refund(c echo.Context) error {
	return nil
}

func Prepay(c echo.Context) error {
	return nil
}

func Notify(c echo.Context) error {
	return nil
}

func GetRecord(c echo.Context) error {
	return nil
}
