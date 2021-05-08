package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"analysis-server/model"
	"common/utils"
	"sync"
)

var (
	Domain      string
	Tenant      string
	Verbose     bool
	Admin       bool
	Client      *http.Client
	TraceId     string
	AccessToken string
	TokenMutex  sync.RWMutex
)

type DescData struct {
	Tc       int64
	Elements []interface{}
}

func genCurl(req *http.Request, body []byte) string {
	msg := "curl"
	msg += " -X " + req.Method
	msg += " " + req.URL.String()
	for k, vs := range req.Header {
		for _, v := range vs {
			msg += fmt.Sprintf(" -H '%s: %s'", k, v)
		}
	}
	if len(body) > 0 && string(body) != "null" {
		msg += " -d '" + string(body) + "'"
	}
	return msg
}

func DoRequest(action string, params interface{}) (*model.CommResp, error) {
	TokenMutex.Lock()
	defer TokenMutex.Unlock()
	return DoRequestwithToken(AccessToken, action, params)
}

func addCookie(request *http.Request, cookieName, cookieVal string) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    cookieVal,
		HttpOnly: false,
		MaxAge:   0,
	}
	request.AddCookie(cookie)
}

func DoRequestwithToken(accessToken, action string,
	params interface{}) (*model.CommResp, error) {
	input, _ := json.Marshal(params)
	url := Domain + "?Action=" + action
	req, err := http.NewRequest("POST", url, bytes.NewReader(input))
	if err != nil {
		return nil, &RespErr{Code: -1, Err: err}
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", "MgrClient")
	if TraceId == "" {
		TraceId = utils.Uuid()
	}
	req.Header.Set("Trace-Id", TraceId)
	if Admin {
		req.Header.Set("Secret-Id", "MgrClientSecretId")
	}
	if accessToken != "" {
		addCookie(req, "access_token", accessToken)
	}
	if Verbose {
		fmt.Println("Request:")
		fmt.Println(genCurl(req, input))
	}

	resp, err := Client.Do(req)
	if err != nil {
		return nil, &RespErr{Code: -1, Err: err}
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if Verbose {
		fmt.Println("Response:")
		fmt.Println("Trace-Id:" + resp.Header.Get("Trace-Id"))
		fmt.Println("X-Requst-Trace:" + resp.Header.Get("X-Requst-Trace"))
		fmt.Println(fmt.Sprintf("%s", body))
	}

	if resp.StatusCode != 200 {
		return nil, &RespErr{Code: -1,
			Err: errors.New(fmt.Sprintf("HTTP Error Code: %d", resp.StatusCode))}

	}
	result := &model.CommResp{}
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	decoder.Decode(result)

	if result.Code != 0 {
		msg := result.Message
		if result.Detail != "" {
			msg += ":" + result.Detail
		}
		return nil, &RespErr{Code: result.Code, Err: errors.New(msg)}
	}
	return result, nil
}

func FormatView(data interface{}, view interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	err = decoder.Decode(view)
	if err != nil {
		return err
	}
	return nil
}
