package main

import (
	"flag"
	"fmt"
	"lemon-ipay-api/alipay"
	"lemon-ipay-api/core"
	"lemon-ipay-api/ipay"
	"lemon-ipay-api/model"
	"lemon-ipay-api/wechat"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	appEnv      = flag.String("APP_ENV", os.Getenv("APP_ENV"), "APP_ENV")
	connEnv     = flag.String("IPAY_CONN", os.Getenv("IPAY_CONN"), "IPAY_CONN")
	bmappingUrl = flag.String("BMAPPING_URL", os.Getenv("BMAPPING_URL"), "BMAPPING_URL")
	hostUrl     = flag.String("IPAY_HOST", os.Getenv("IPAY_HOST"), "IPAY_HOST")
)

func init() {
	flag.Parse()
	envParam := &core.EnvParamDto{
		AppEnv:      *appEnv,
		ConnEnv:     *connEnv,
		BmappingUrl: *bmappingUrl,
		HostUrl:     *hostUrl,
	}
	core.InitEnv(envParam)
	model.Db = InitDB("mysql", envParam.ConnEnv)
	model.Db.Sync(new(model.WxAccount), new(model.NotifyWechat), new(model.AlAccount), new(model.NotifyAlipay))
}

func InitDB(dialect, conn string) (newDb *xorm.Engine) {
	newDb, err := xorm.NewEngine(dialect, conn)
	if err != nil {
		panic(err)
	}
	return
}

func main() {

	e := echo.New()
	e.Use(middleware.CORS())
	RegisterApi(e)
	e.Start(":5000")
}

func RegisterApi(e *echo.Echo) {

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "lemon epay")
	})
	e.GET("/ping", func(c echo.Context) error {
		fmt.Println("pong")
		return c.String(http.StatusOK, "pong")
	})
	track := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//util.Lang = c.Request().Header["Accept-Language"]
			return next(c)
		}
	}

	v3 := e.Group("/v3", track)
	v3.POST("/pay", ipay.Pay)
	v3.POST("/query", ipay.Query)
	v3.POST("/reverse", ipay.Reverse)
	v3.POST("/refund", ipay.Refund)
	v3.POST("/Prepay", ipay.Prepay)

	v3.GET("/record/:Id", ipay.GetRecord)

	wx := v3.Group("/wx")
	wx.POST("/pay", wechat.Pay)
	wx.POST("/query", wechat.Query)
	wx.POST("/reverse", wechat.Reverse)
	wx.POST("/refund", wechat.Refund)
	wx.POST("/prepay", wechat.Prepay)
	wx.POST("/notify", wechat.Notify)
	wx.GET("/prepayeasy", wechat.PrepayEasy)
	wx.GET("/prepayopenid", wechat.PrepayOpenId)

	al := v3.Group("/al")
	al.POST("/pay", alipay.Pay)
	al.POST("/query", alipay.Query)
	al.POST("/reverse", alipay.Reverse)
	al.POST("/refund", alipay.Refund)
	al.POST("/prepay", alipay.Prepay)
	al.POST("/notify", alipay.Notify)

}
