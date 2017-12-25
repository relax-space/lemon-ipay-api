package model

import (
	"errors"
	"time"
)

type WxAccount struct {
	EId      int64  `json:"e_id" xorm:"int64 notnull pk 'e_id'"`
	AppId    string `json:"app_id" xorm:"varchar(25)"`
	SubAppId string `json:"sub_app_id" xorm:"varchar(25)"`
	Key      string `json:"key" xorm:"varchar(50)"`
	MchId    string `json:"mch_id" xorm:"varchar(25)"`

	SubMchId string `json:"sub_mch_id" xorm:"varchar(25)"`
	CertName string `json:"cert_name" xorm:"varchar(50)"`
	CertKey  string `json:"cert_key" xorm:"varchar(50)"`
	RootCa   string `json:"root_ca" xorm:"varchar(50)"`
	Secret   string `json:"secret" xorm:"varchar(50)"`

	Description string    `json:"description" xorm:"nvarchar(100)"`
	CreatedAt   time.Time `json:"created_at" xorm:"created"`
	UpdatedAt   time.Time `json:"updated_at" xorm:"updated"`
	InUserId    int       `json:"in_user_id" xorm:"INT"`
	ModiUserId  int       `json:"modi_user_id" xorm:"INT"`
}

func (WxAccount) TableName() string {
	return "wechat"
}

func (WxAccount) Get(eId int64) (account WxAccount, err error) {
	has, err := Db.Where("e_id =?", eId).Get(&account)
	if err != nil {
		return
	} else if !has {
		err = errors.New("no data has found.")
		return
	}
	return
}
