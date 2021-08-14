//go:generate fyne bundle -o bundle.go -append myapp.png
//go:generate fyne bundle -o bundle.go -append AlibabaPuHuiTi-2-55-Regular.ttf
//
//fyne package -os linux -icon myapp.png
//fyne package -os windows -icon myapp.png
//fyne package -os android -appID cn.corehub.goyoung -icon myapp.png
package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"GoYoung/model"
	"GoYoung/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/qifengzhang007/goCurl"
)

const (
	baiduURL         = "baidu.com:443"
	feiYoungRedirect = "http://www.msftconnecttest.com/redirect"
)

var (
	version    = "Version 1.0.0"
	httpClient = goCurl.CreateHttpClient()
	myApp      MyApp
	message    = ""
	user       model.User
)

type MyApp struct {
	app      fyne.App
	win      fyne.Window
	lab      *widget.Label
	msg      *widget.Entry
	prefix   *widget.Select
	username *widget.Entry
	password *widget.Entry
}

func main() {
	myApp = MyApp{app: app.NewWithID("GoYoung")}
	myApp.app.SetIcon(resourceMyappPng)
	myApp.app.Settings().SetTheme(&MyTheme{})
	myApp.win = myApp.app.NewWindow("LinkGoYoung")
	myApp.win.Resize(fyne.Size{Width: 320, Height: 480})

	myApp.lab = widget.NewLabel(version)
	myApp.msg = widget.NewMultiLineEntry()
	myApp.prefix = widget.NewSelect([]string{"!^Adcm0", "!^Iqnd0"}, func(value string) {})
	myApp.prefix.SetSelectedIndex(0)
	myApp.username = widget.NewEntry()
	myApp.username.SetPlaceHolder("账号")
	myApp.password = widget.NewPasswordEntry()
	myApp.password.SetPlaceHolder("密码")
	user.ReadUserInfoJson()
	myApp.prefix.SetSelected(user.UserHard)
	myApp.username.SetText(user.UserAccount)
	myApp.password.SetText(user.PassWord)
	message += CheckServer(baiduURL) + "\n"
	myApp.msg.SetText(message)

	form := &widget.Form{
		BaseWidget: widget.BaseWidget{},
		Items: []*widget.FormItem{
			{Text: "账号", Widget: myApp.username},
			{Text: "密码", Widget: myApp.password},
		},
		OnSubmit:   submitHandle,
		OnCancel:   cancelHandle,
		SubmitText: "登录",
		CancelText: "下线",
	}

	content := container.New(layout.NewGridLayoutWithColumns(1), myApp.msg,
		container.NewVBox(myApp.lab, myApp.prefix, form))
	myApp.win.SetContent(content)
	myApp.win.ShowAndRun()
}

func submitHandle() {
	loginURL, userIp, nasIp, userMac := nextUrl(myApp.msg, message, feiYoungRedirect)
	fmt.Printf("校园网关IP：%s\n", nasIp)
	if nasIp == "" {
		message += "重定向失败，请检查网络环境\n"
		myApp.msg.SetText(message)
		return
	}
	message = "当前网络信息："
	message += "\nnasip:" + nasIp
	message += "\nuserip:" + userIp
	message += "\nusermac:" + userMac

	resp, err := httpClient.Get(loginURL, goCurl.Options{
		Headers: map[string]interface{}{
			"User-Agent": "CDMA+WLAN(Mios)",
		},
		SetResCharset: "utf-8",
	})
	if err != nil {
		return
	}
	body, err := resp.GetContents()
	if err != nil {
		return
	}

	loginURL = utils.ParseXML(body, "WISPAccessGatewayParam", "Redirect", "LoginURL")
	fmt.Println("loginURL:" + loginURL)
	tag := login(&myApp, message, loginURL, myApp.prefix.Selected, strings.TrimSpace(myApp.username.Text), strings.TrimSpace(myApp.password.Text))
	dialog.ShowInformation("登陆……", tag, myApp.win)
	if tag == "50：认证成功" {
		user.UserHard = myApp.prefix.Selected
		user.UserAccount = myApp.username.Text
		user.PassWord = myApp.password.Text
		user.LastLoginURL = loginURL
		user.SaveUserInfoJson()
	}
}

func cancelHandle() {
	if user.LastLoginURL != "" {
		dialog.ShowInformation("下线……", logout(&myApp, message, user), myApp.win)
		return
	}
	dialog.ShowInformation("Error", "请先登录", myApp.win)
}

func login(app *MyApp, message string, url, userHard, userName, passWord string) string {
	token := "UserName=" + userHard + userName + "&Password=" + passWord + "&AidcAuthAttr1=" + time.Now().Format("20060102150405") +
		"&AidcAuthAttr3=keuyGQlK&AidcAuthAttr4=zrDgXllCChyJHjwkcRwhygP0&AidcAuthAttr5=kfe1GQhXdGqOFDteego5zwP9IsNoxX7djTWspPrYm1A%3D%3D&" +
		"AidcAuthAttr6=5Ia4cQhDfXSFbTtUDGY1yx8%3D&AidcAuthAttr7=6ZWiVlwdNiHMXCpOagQv2w2MQs0ohTWJnTu8qK5OibhCydTpTxkI88wadKPWby%2F2PKCVaZ" +
		"UxglbBs96%2FtmLE89M8AJ6y28o7qolpFep%2FcYFFRLd7H4MAMrDUMRO0F%2B93jh14fiAZYmtk9hdp%2BZ5w%2BjMQUoV4TCtM9VJ07XQwxlMVg%2F0YKrS1s3hXA" +
		"stdQ1fvdSn3nAVGgdxc%2BJQDrQ%3D%3D&AidcAuthAttr8=jPSyBQxVaXWTQWUaakluj06scJ98nyqCyX7y%2FLUk1OkXiNjkXhVGvJhyTuLDaCPhK%2FOFJttlxxi" +
		"VqNKupnDXkp9%2BR9D9j8p2j5h8FOxoatMaGu0oRdk%3D&createAuthorFlag=0"
	resp, err := httpClient.Post(url, goCurl.Options{
		Headers: map[string]interface{}{
			"User-Agent":   "CDMA+WLAN(Mios)",
			"Content-Type": "application/x-www-form-urlencoded",
		},
		XML:           token,
		SetResCharset: "utf-8",
		Timeout:       2,
	})
	if err != nil {
		message += fmt.Sprintf("\nLogin请求出错：%s", err.Error())
		app.msg.SetText(message)
		return "Login Timeout"
	}
	body, err := resp.GetContents()
	if err != nil {
		message += fmt.Sprintf("\nLogin请求失败,错误明细：%s", err.Error())
		app.msg.SetText(message)
		return "Login请求出错"
	}
	body = utils.ParseXML(body, "WISPAccessGatewayParam", "AuthenticationReply", "ReplyMessage")
	message += fmt.Sprintf("\n请求结果：%s", body)
	app.msg.SetText(message)
	return body
}

func logout(app *MyApp, message string, user model.User) string {
	u := "http://58.53.199.144:8001/wispr_logout.jsp?" + strings.Split(user.LastLoginURL, "?")[1]
	var httpClient = goCurl.CreateHttpClient()
	resp, err := httpClient.Get(u, goCurl.Options{
		Headers: map[string]interface{}{
			"User-Agent": "CDMA+WLAN(Mios)",
		},
		SetResCharset: "utf-8",
		Timeout:       2,
	})
	if err != nil {
		return "下线请求超时"
	}
	body, err := resp.GetContents()
	if err != nil {
		return "下线信息未知"
	}
	body = utils.ParseXML(body, "WISPAccessGatewayParam", "LogoffReply", "ResponseCode")
	if body == "150" {
		message += "150:下线成功\n"
		app.msg.SetText(message)
		return "150:下线成功"
	} else {
		message += "255:下线失败\n"
		app.msg.SetText(message)
		return "255:下线失败"
	}
}

func nextUrl(msg *widget.Entry, message, urlIn string) (urlOut, userIp, nasIp, userMac string) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 30 * time.Second,
	}
	res, err := client.Get(urlIn)
	if err != nil {
		message += fmt.Sprintf("\n重定向请求出错：%s", err.Error())
		msg.SetText(message)
		return
	}
	if res.StatusCode != http.StatusFound {
		message += fmt.Sprintf("\n[Error]StatusCode:%v", res.StatusCode)
		msg.SetText(message)
		return
	}
	u, _ := url.Parse(res.Header.Get("Location"))
	query := u.Query()
	userIp, nasIp, userMac = query.Get("userip"), query.Get("nasip"), query.Get("usermac")
	urlOut = fmt.Sprintf("http://58.53.199.144:8001/?userip=%s&wlanacname=&nasip=%s&usermac=%s&aidcauthtype=0", userIp, nasIp, userMac)
	return urlOut, userIp, nasIp, userMac
}

func CheckServer(url string) string {
	timeout := 5 * time.Second
	//t1 := time.Now()
	_, err := net.DialTimeout("tcp", url, timeout)
	//massage += "\n网络测试时长 :" + time.Now().Sub(t1).String()

	if err == nil {
		return "已接入互联网，只能进行下线操作"
	} else {
		return "未接入互联网"
	}
}
