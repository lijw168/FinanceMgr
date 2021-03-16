package utils

import (
	"bytes"
	"encoding/json"
	//"errors"
	"fmt"
	//"io/ioutil"
	"net/http"
	"time"

	"common/log"
	//"common/message"
)

type HttpClient struct {
	logger  *log.Logger
	verbose bool
	traceID string
	client  *http.Client
}

func NewHttpClient(verbose bool, traceID string, timeout time.Duration, logger *log.Logger) *HttpClient {
	httpClient := new(HttpClient)
	httpClient.verbose = verbose
	httpClient.traceID = traceID
	httpClient.logger = logger
	httpClient.client = new(http.Client)
	httpClient.client.Timeout = timeout
	return httpClient
}
func (hc *HttpClient) SetLogger(logger *log.Logger) {
	hc.logger = logger
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

// func (hc *HttpClient) SendRequest(domain, reqMethod, action string, dataContent interface{}) (*message.ResponseParam, error) {
// 	input, _ := json.Marshal(dataContent)
// 	url := domain + "?Action=" + action
// 	req, err := http.NewRequest(reqMethod, url, bytes.NewReader(input))
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	req.Header.Set("User-Agent", "ZbsClient")
// 	if hc.traceID == "" {
// 		hc.traceID = Uuid()
// 	}
// 	req.Header.Set("Trace-Id", hc.traceID)
// 	if hc.verbose {
// 		if hc.logger != nil {
// 			hc.logger.DebugContext(req.Context(), "Request: %s]", genCurl(req, input))
// 		}
// 	}

// 	resp, err := hc.client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	body, _ := ioutil.ReadAll(resp.Body)
// 	resp.Body.Close()

// 	if hc.verbose {
// 		if hc.logger != nil {
// 			hc.logger.DebugContext(req.Context(), "Response:")
// 			hc.logger.DebugContext(req.Context(), "Trace-Id:"+resp.Header.Get("Trace-Id"))
// 			hc.logger.DebugContext(req.Context(), "X-Requst-Trace:"+resp.Header.Get("X-Requst-Trace"))
// 			hc.logger.DebugContext(req.Context(), fmt.Sprintf("%s", body))
// 		}
// 	}
// 	if resp.StatusCode != 200 {
// 		return nil, fmt.Errorf("HTTP Error Code: %d", resp.StatusCode)
// 	}
// 	result := &message.ResponseParam{}
// 	decoder := json.NewDecoder(bytes.NewReader(body))
// 	decoder.UseNumber()
// 	decoder.Decode(result)

// 	if result.Code != 0 {
// 		msg := result.Message
// 		if result.Detail != "" {
// 			msg += ":" + result.Detail
// 		}
// 		return nil, errors.New(msg)
// 	}
// 	return result, nil
// }

func (hc *HttpClient) FormatView(data interface{}, view interface{}) error {
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
