package model

import (
	"errors"
	"time"
)

type NotifyAlipay struct {
	NotifyTime string `json:"notify_time,omitempty" mapstruct:"notify_time"`
	NotifyType string `json:"notify_type,omitempty" mapstruct:"notify_type"`
	NotifyId   string `json:"notify_id,omitempty" mapstruct:"notify_id"`
	SignType   string `json:"sign_type,omitempty" mapstruct:"sign_type"`
	Sign       string `json:"sign,omitempty"  mapstruct:"sign" xorm:"varchar(256)"`

	TradeNo    string `json:"trade_no,omitempty" mapstruct:"trade_no"`
	AppId      string `json:"app_id,omitempty" mapstruct:"app_id" xorm:"appid"`
	OutTradeNo string `json:"out_trade_no,omitempty" mapstruct:"out_trade_no"`
	OutBizNo   string `json:"out_biz_no,omitempty" mapstruct:"out_biz_no"`
	BuyerId    string `json:"buyer_id,omitempty" mapstruct:"buyer_id"`

	BuyerLogonId string `json:"buyer_logon_id,omitempty" mapstruct:"buyer_logon_id"`
	SellerId     string `json:"seller_id,omitempty" mapstruct:"seller_id"`
	SellerEmail  string `json:"seller_email,omitempty" mapstruct:"seller_email"`
	TradeStatus  string `json:"trade_status,omitempty" mapstruct:"trade_status"`
	TotalAmount  string `json:"total_amount,omitempty" mapstruct:"total_amount"` //float64

	ReceiptAmount  string `json:"receipt_amount,omitempty" mapstruct:"receipt_amount"`     //float64
	InvoiceAmount  string `json:"invoice_amount,omitempty" mapstruct:"invoice_amount"`     //float64
	BuyerPayAmount string `json:"buyer_pay_amount,omitempty" mapstruct:"buyer_pay_amount"` //float64
	PointAmount    string `json:"point_amount,omitempty" mapstruct:"point_amount"`         //float64
	RefundFee      string `json:"refund_fee,omitempty" mapstruct:"refund_fee"`             //float64

	SendBackFee string `json:"send_back_fee,omitempty" mapstruct:"send_back_fee"` //float64
	Subject     string `json:"subject,omitempty" mapstruct:"subject" xorm:"varchar(256)"`
	Body        string `json:"body,omitempty" mapstruct:"body" xorm:"varchar(400)"`
	GmtCreate   string `json:"gmt_create,omitempty" mapstruct:"gmt_create"`
	GmtPayment  string `json:"gmt_payment,omitempty" mapstruct:"gmt_payment"`

	GmtRefund    string    `json:"gmt_refund,omitempty" mapstruct:"gmt_refund"`
	GmtClose     string    `json:"gmt_close,omitempty" mapstruct:"gmt_close"`
	FundBillList string    `json:"fund_bill_list,omitempty" mapstruct:"fund_bill_list" xorm:"varchar(512)"`
	CreatedAt    time.Time `json:"created_at" mapstruct:"created_at" xorm:"created"`
}

func (NotifyAlipay) Get(appId, outTradeNo string) (notify NotifyAlipay, err error) {
	has, err := Db.Where("appId =?", appId).And("out_trade_no=?", outTradeNo).Get(&notify)
	if err != nil {
		return
	} else if !has {
		err = errors.New("no data has found.")
		return
	}
	return
}

func (NotifyAlipay) InsertOne(notify *NotifyAlipay) (err error) {
	has, err := Db.Where("appId =?", notify.AppId).And("out_trade_no=?", notify.OutTradeNo).Get(&NotifyAlipay{})
	if err != nil {
		return
	} else if has { //success,when data exsits
		//err = errors.New("insert failure, because data is exist")
		return
	}
	r, err := Db.InsertOne(notify)
	if err != nil {
		return
	} else if r == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}
