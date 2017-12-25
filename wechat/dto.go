package wechat

import (
	wxpay "github.com/relax-space/lemon-wxpay-sdk"
)

type ReqPayDto struct {
	*wxpay.ReqPayDto
	EId int64 `json:"e_id"`
}
type ReqQueryDto struct {
	*wxpay.ReqQueryDto
	EId int64 `json:"e_id"`
}
type ReqRefundDto struct {
	*wxpay.ReqRefundDto
	EId int64 `json:"e_id"`
}
type ReqReverseDto struct {
	*wxpay.ReqReverseDto
	EId int64 `json:"e_id"`
}
type ReqRefundQueryDto struct {
	*wxpay.ReqRefundQueryDto
	EId int64 `json:"e_id"`
}
type ReqPrepayDto struct {
	*wxpay.ReqPrepayDto
	EId int64 `json:"e_id"`
}

type ReqPrepayEasyDto struct {
	*wxpay.ReqPrepayDto
	EId   int64  `json:"e_id" query:"e_id"`
	AppId string `json:"app_id" query:"app_id"` //required
	Scope string `json:"scope" query:"scope"`   //option
	State string `json:"state" query:"state"`   //option

	//Secret      string `json:"secret"`
	RedirectUrl string `json:"redirect_url" query:"redirect_url"`
	PageUrl     string `json:"page_url" query:"page_url"` //option
	//ReqPrepayDto *ReqPrepayDto `json:"prepay_param" query:"prepay_param"`
}
