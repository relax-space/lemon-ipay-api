package model

import (
	"errors"
	"time"
)

type AlAccount struct {
	EId      int64  `json:"e_id" xorm:"int64 notnull pk 'e_id'"`
	AppId    string `json:"app_id" xorm:"varchar(25)"`
	SubAppId string `json:"sub_app_id" xorm:"varchar(25)"`
	PriKey   string `json:"pri_key" xorm:"varchar(1024)"`
	PubKey   string `json:"pub_key" xorm:"varchar(1024)"`

	AuthToken            string    `json:"auth_token" xorm:"varchar(40)"`
	SysServiceProviderId string    `json:"sys_service_provider_id" xorm:"varchar(64)"`
	Description          string    `json:"description" xorm:"varchar(100)"`
	CreatedAt            time.Time `json:"created_at" xorm:"created"`
	UpdatedAt            time.Time `json:"updated_at" xorm:"updated"`

	InUserId   int `json:"in_user_id" xorm:"INT"`
	ModiUserId int `json:"modi_user_id" xorm:"INT"`
}

func (AlAccount) TableName() string {
	return "alipay"
}

func (AlAccount) Get(eId int64) (account *AlAccount, err error) {
	account = &AlAccount{}
	has, err := Db.Where("e_id =?", eId).Get(account)
	if err != nil {
		return
	} else if !has {
		err = errors.New("no data has found.")
		return
	}
	return
}

func (AlAccount) GetByAppId(appId string) (account AlAccount, err error) {
	has, err := Db.Where("app_id =?", appId).Get(&account)
	if err != nil {
		return
	} else if !has {
		err = errors.New("no data has found.")
		return
	}
	return
}
