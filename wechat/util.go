package wechat

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/labstack/echo"
)

const (
	IPAY_WECHAT_PREPAY       = "IPAY_WECHAT_PREPAY"
	IPAY_WECHAT_PREPAY_INNER = "IPAY_WECHAT_PREPAY_INNER"
	IPAY_WECHAT_PREPAY_ERROR = "IPAY_WECHAT_PREPAY_ERROR"
)

func SetCookie(key, value string, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = key
	value = url.QueryEscape(value)
	cookie.Value = value
	cookie.Domain = "p2shop.cn"
	cookie.Path = "/"
	//cookie.Expires = time.Now().Add(1 * time.Hour)
	c.SetCookie(cookie)
}

func SetCookieObj(key string, value interface{}, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = key
	b, _ := json.Marshal(value)
	cookie.Value = url.QueryEscape(string(b))
	cookie.Domain = "p2shop.cn"
	cookie.Path = "/"
	//cookie.Expires = time.Now().Add(1 * time.Hour)
	c.SetCookie(cookie)
}
