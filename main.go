package main

import (
	"net/http"
	"log"
	"fmt"
	"io/ioutil"
	"github.com/skip2/go-qrcode"
	"encoding/json"
	"strings"
)

var (
	cookieDict = make(map[string]string)
)

const (
	GET_KEY_URL = "https://work.weixin.qq.com/wework_admin/wwqrlogin/get_key?r=0.5914192326208461&login_type=login_admin&callback=wwqrloginCallback_1544150579268&redirect_uri=https%3A%2F%2Fwork.weixin.qq.com%2Fwework_admin%2Floginpage_wx%3Fpagekey%3D1544150579268956"
	QRCODE_URL = "https://work.weixin.qq.com/wwqrlogin/scan?login_type=login_admin&qrcode_key="
	CHECK_STATUS_URL = "https://work.weixin.qq.com/wework_admin/wwqrlogin/check"
	LOGIN_PAGE_URL = "https://work.weixin.qq.com/wework_admin/loginpage_wx"
	CORP_APP_URL = "https://work.weixin.qq.com/wework_admin/getCorpApplication?lang=zh_CN&f=json&ajax=1&timeZoneInfo%5Bzone_offset%5D=-8&random=0.29567277370227707"
)

// 模拟企业微信二维码登录，获取cookie, 用于请求相关接口
func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/getkey", handleGetKey)
	http.HandleFunc("/getqrcode", handleGetQrcode)
	http.HandleFunc("/checkstatus", handleCheckStatus)
	http.HandleFunc("/loginpagewx", handleLoginPageWx)
	http.HandleFunc("/getcorpapp", handleGetCorpApp)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
func handleGetCorpApp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	headers := map[string]string{
		"Cookie": cookieDict[r.Form.Get("qrcode_key")],
		"Origin": "https://work.weixin.qq.com",
		"Referer": "https://work.weixin.qq.com/wework_admin/frame",
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36",
		"x-requested-with": "XMLHttpRequest",
		"content-type": "application/x-www-form-urlencoded",
	}
	res, err := Get(CORP_APP_URL, nil, headers)
	CheckError(err)
	content, err := ioutil.ReadAll(res.Body)
	CheckError(err)
	defer res.Body.Close()
	w.Write(content)
}

// 扫码后的登录请求 code换cookie
func handleLoginPageWx(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	params := map[string]string{
		"code": r.Form.Get("code"),
		"qrcode_key": r.Form.Get("qrcode_key"),
		"pagekey" : "1544178373478405",
		"wwqrlogin": "1",
	}
	headers := map[string]string{
		"upgrade-insecure-requests": "1",
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.80 Safari/537.36",
	}
	res, err := Get(LOGIN_PAGE_URL, params, headers)
	CheckError(err)
	var strSli = make([]string, len(res.Cookies()))
	for _,c := range res.Cookies() {
		if s := c.String(); s != "" {
			strSli = append(strSli, c.String())
		}
	}
	cookieDict[params["qrcode_key"]] = strings.Join(strSli, ";")
	fmt.Println(cookieDict)
	w.Write([]byte("ok"))
}
// 监听扫码动作
func handleCheckStatus(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	params := map[string]string{
		"status": r.Form.Get("status"),
		"qrcode_key": r.Form.Get("qrcode_key"),
		"randNum": "0.4500575181163635",
	}
	res, err := Get(CHECK_STATUS_URL, params, nil)
	fmt.Println("status code:", res.StatusCode)
	CheckError(err)
	content, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	CheckError(err)
	w.Write(content)
}

// 获取登录二维码
func handleGetQrcode(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	if key == "" {
		w.Header().Set("status", "404")
		return
	}

	png, err := qrcode.Encode(QRCODE_URL + key, qrcode.Medium, 256)
	CheckError(err)
	w.Write(png)
}

// 获取二维码的key
func handleGetKey(w http.ResponseWriter, r *http.Request) {
	res, err := Get(GET_KEY_URL, nil, nil)
	CheckError(err)

	content, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	CheckError(err)

	var dat map[string]interface{}
	json.Unmarshal(content, &dat)
	key := dat["data"].(map[string]interface{})["qrcode_key"].(string)

	w.Write([]byte(key))
}

// 登录页面
func handleHome(w http.ResponseWriter, r *http.Request)  {
	http.ServeFile(w, r, "./public/index.html")
}