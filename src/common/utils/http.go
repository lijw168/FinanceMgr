package utils

import (
	// "bytes"
	// "encoding/json"
	// "fmt"
	// "io"
	// "io/ioutil"
	// "net/http"
	// "time"

	//"common/message"
)

// type IHttpWrapper interface {
// 	SendRequest(url string,
// 		req interface{},
// 		rspdata interface{}) (*message.ResponseParam, error)
// 	SendZbsRequest(url string,
// 		req interface{},
// 		rspdata interface{}) error
// 	SendRequestWithByte(url string,
// 		req []byte,
// 		rspdata interface{}) (*message.ResponseParam, error)
// }

// func NewHttpWrapper(timeout time.Duration) IHttpWrapper {
// 	return &HttpWrapper{timeout: timeout}
// }

// type HttpWrapper struct {
// 	timeout time.Duration
// }

// // req request []byte
// func (h HttpWrapper) SendRequestWithByte(url string,
// 	req []byte,
// 	rspdata interface{}) (*message.ResponseParam, error) {
// 	code, buf, err := h.sendRequest("POST", url, req)
// 	if err != nil {
// 		return nil, err
// 	} else if code != 200 {
// 		err = fmt.Errorf("return code not 200 %s", buf)
// 		return nil, err
// 	} else if buf == nil {
// 		err = fmt.Errorf("request bug nil")
// 		return nil, err
// 	}
// 	rsp := &message.ResponseParam{Data: rspdata}

// 	if err = json.Unmarshal(buf, rsp); err != nil {
// 		err = fmt.Errorf("unmarshal body err %s", string(buf))
// 	}
// 	return rsp, err
// }

// // req request objec
// func (h HttpWrapper) SendRequest(url string,
// 	req interface{},
// 	rspdata interface{}) (*message.ResponseParam, error) {
// 	reqBuf, err := json.Marshal(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return h.SendRequestWithByte(url, reqBuf, rspdata)
// }

// func (h HttpWrapper) SendZbsRequest(url string,
// 	req interface{},
// 	rspdata interface{}) error {
// 	rsp, err := h.SendRequest(url, req, rspdata)
// 	if err != nil {
// 		return err
// 	}

// 	if rsp.Code != 0 {
// 		err := fmt.Errorf("Remote Server returned error code %d",
// 			rsp.Code)
// 		return err
// 	}

// 	return nil
// }

// func (h HttpWrapper) sendRequest(method, url string,
// 	reqBody []byte) (int, []byte, error) {
// 	var (
// 		code   = -1
// 		body   []byte
// 		err    error
// 		reader io.Reader
// 	)
// 	if reqBody != nil {
// 		reader = bytes.NewReader(reqBody)
// 	}
// 	// FIXME 此处代码会导致HTTP请求完成后不关闭
// 	//client := http.Client{
// 	//	Transport: &http.Transport{
// 	//		Dial: (&net.Dialer{
// 	//			Timeout:   30 * time.Second,
// 	//			KeepAlive: 30 * time.Second,
// 	//		}).Dial,
// 	//		TLSHandshakeTimeout:   30 * time.Second,
// 	//		ResponseHeaderTimeout: 30 * time.Second,
// 	//		ExpectContinueTimeout: 30 * time.Second,
// 	//	}}

// 	client := http.Client{
// 		Timeout: h.timeout,
// 	}
// 	req, err := http.NewRequest(method, url, reader)
// 	if err != nil {
// 		return code, body, err
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	rsp, err := client.Do(req)
// 	if err != nil {
// 		return code, body, err
// 	}
// 	defer rsp.Body.Close()
// 	code = rsp.StatusCode
// 	body, err = ioutil.ReadAll(rsp.Body)
// 	return code, body, err
// }
