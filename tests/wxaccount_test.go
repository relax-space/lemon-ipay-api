package tests

import (
	"fmt"
	"kit/test"
	"testing"

	"lemon-ipay-api/model"
)

func Test_WxAccount_Get(t *testing.T) {
	account, err := model.WxAccount{}.Get(10001)
	test.Ok(t, err)
	fmt.Println(account)
}
