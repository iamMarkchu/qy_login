package main

import (
	"net/http"
	"fmt"
	"net/url"
)

func Get(resourceUrl string, params map[string]string, headers map[string]string) (http.Response, error) {
	u, err := url.Parse(resourceUrl)
	CheckError(err)
	q := u.Query()
	if params != nil {
		for key,param := range params {
			q.Set(key, param)
		}
	}
	u.RawQuery = q.Encode()
	fmt.Println("GET请求:", u.String())
	req, err := http.NewRequest("GET", u.String(), nil)
	CheckError(err)

	// 初始化http client
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// 设置请求头信息
	if headers != nil {
		for key,header := range headers {
			req.Header.Set(key, header)
		}
	}
	res, err := client.Do(req)
	CheckError(err)
	return *res, err
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
