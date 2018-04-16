package alipay

import (
	"errors"
	"lemon-ipay-api/core"
	"lemon-ipay-api/model"
	"strconv"

	"github.com/relax-space/go-kit/base"
	"github.com/relax-space/go-kit/sign"
	paysdk "github.com/relax-space/lemon-alipay-sdk"
)

func NotifyQuery(account *model.AlAccount, outTradeNo string) (result *paysdk.RespQueryDto, err error) {
	var reqDto paysdk.ReqQueryDto
	reqDto.ReqBaseDto = &paysdk.ReqBaseDto{
		AppId:        account.AppId,
		AppAuthToken: account.AuthToken,
	}

	customDto := &paysdk.ReqCustomerDto{
		PriKey: account.PriKey,
		PubKey: account.PubKey,
	}
	reqDto.OutTradeNo = outTradeNo
	_, _, result, err = paysdk.Query(&reqDto, customDto)
	return
}

func NotifyValidN(body, signParam, outTradeNo, totalAmount string, mapParam map[string]interface{}) (err error) {

	//0.get account info
	bodyMap := base.ParseMapObject(body, core.NOTIFY_BODY_SEP1, core.NOTIFY_BODY_SEP2)
	var eId int64
	var flag bool
	if eIdObj, ok := bodyMap["e_id"]; ok {
		if eId, err = strconv.ParseInt(eIdObj.(string), 10, 64); err == nil {
			flag = true
		}
	}
	if !flag {
		err = errors.New("e_id(int64) is not existed in param(param name:body) or format is not correct")
		return
	}

	account, err := model.AlAccount{}.Get(eId)
	if err != nil {
		return
	}

	//1.valid sign
	signStr := signParam
	delete(mapParam, "sign")
	delete(mapParam, "sign_type")

	if !sign.CheckSha1Sign(base.JoinMapObject(mapParam), signStr, account.PubKey) {
		err = errors.New("sign valid failure")
		return
	}

	//2.valid data
	queryDto, err := NotifyQuery(account, outTradeNo)
	if err != nil {
		return
	}
	if !(queryDto.TotalAmount == totalAmount) {
		err = errors.New("amount is exception")
		return
	}
	return
}
