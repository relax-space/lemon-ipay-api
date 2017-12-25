package tests

import (
	"lemon-ipay-api/core"
	"lemon-ipay-api/model"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func init() {
	initTest()
}

func initTest() {
	envParam := &core.EnvParamDto{
		AppEnv:      "",
		ConnEnv:     os.Getenv("IPAY_CONN"),
		BmappingUrl: "",
		HostUrl:     os.Getenv("IPAY_HOST"),
	}
	core.InitEnv(envParam)
	model.Db = InitDB("mysql", envParam.ConnEnv)
	//model.Db.Sync(new(model.WxAccount))
}

func InitDB(dialect, conn string) (newDb *xorm.Engine) {
	newDb, err := xorm.NewEngine(dialect, conn)
	if err != nil {
		panic(err)
	}
	return
}
