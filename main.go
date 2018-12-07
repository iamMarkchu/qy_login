package main

import (
	"net/http"
	"log"
	"fmt"
	"io/ioutil"
	"github.com/skip2/go-qrcode"
	"encoding/json"
	"net/url"
)

const (
	GET_KEY_URL = "https://work.weixin.qq.com/wework_admin/wwqrlogin/get_key?r=0.5914192326208461&login_type=login_admin&callback=wwqrloginCallback_1544150579268&redirect_uri=https%3A%2F%2Fwork.weixin.qq.com%2Fwework_admin%2Floginpage_wx%3Fpagekey%3D1544150579268956"
	QRCODE_URL = "https://work.weixin.qq.com/wwqrlogin/scan?login_type=login_admin&qrcode_key="
	CHECK_STATUS_URL = "https://work.weixin.qq.com/wework_admin/wwqrlogin/check"
	LOGIN_PAGE_URL = "https://work.weixin.qq.com/wework_admin/loginpage_wx"
)

// 模拟企业微信二维码登录，获取cookie, 用于请求相关接口
func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/getkey", handleGetKey)
	http.HandleFunc("/getqrcode", handleGetQrcode)
	http.HandleFunc("/checkstatus", handleCheckStatus)
	http.HandleFunc("/loginpagewx", handleLoginPageWx)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
func handleLoginPageWx(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.Form.Get("code")
	qrcodeKey := r.Form.Get("qrcode_key")
	pagekey := "1544178373478405"
	wwqrlogin := "1"

	u, err := url.Parse(LOGIN_PAGE_URL)
	if err != nil {
		fmt.Println(err)
	}
	q := u.Query()

	q.Set("code", code)
	q.Set("qrcode_key", qrcodeKey)
	q.Set("pagekey", pagekey)
	q.Set("wwqrlogin", wwqrlogin)
	u.RawQuery = q.Encode()
	fmt.Println(u.String())
	// w.Write([]byte(u.String()))
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		fmt.Println(err)
	}
	client := &http.Client{}
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.80 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Cookies:", res.Cookies())
	fmt.Println("Header:", res.Header)
	content, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	//w.Write(content)
	fmt.Println(string(content))
}
func handleCheckStatus(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	status := r.Form.Get("status")
	qrcodeKey := r.Form.Get("qrcode_key")
	randNum := "0.4500575181163635"
	u, err := url.Parse(CHECK_STATUS_URL)
	if err != nil {
		fmt.Println(err)
	}
	q := u.Query()

	q.Set("status", status)
	q.Set("qrcode_key", qrcodeKey)
	q.Set("randNum", randNum)
	u.RawQuery = q.Encode()
	fmt.Println(u.String())

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		fmt.Println(err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	fmt.Println("status code:", res.StatusCode)
	if err != nil {
		fmt.Println(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	w.Write(content)
}
func handleGetQrcode(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	if key == "" {
		w.Header().Set("status", "404")
		return
	}

	png, err := qrcode.Encode(QRCODE_URL + key, qrcode.Medium, 256)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(png)
}
func handleGetKey(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", GET_KEY_URL, nil)
	if err != nil {
		fmt.Println(err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	var dat map[string]interface{}
	json.Unmarshal(content, &dat)
	key := dat["data"].(map[string]interface{})["qrcode_key"].(string)

	w.Write([]byte(key))
}

func handleHome(w http.ResponseWriter, r *http.Request)  {
	http.ServeFile(w, r, "./public/index.html")
}