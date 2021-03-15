package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"common/message"
)

type SimpleHttpServer struct {
}

func (s *SimpleHttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var reader io.ReadCloser
	if reader = r.Body; reader == nil {
		return
	}

	r.FormValue("Action")

	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var data interface{}
	if err := json.Unmarshal(buf, &data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var rsp message.ResponseParam
	rsp.Code = 0
	rsp.Detail = "Detail"
	rsp.Message = "OK"
	rsp.Data = data

	if buf, err = json.Marshal(rsp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func (s *SimpleHttpServer) Start() {
	http.ListenAndServe(":8080", s)
}

func TestSendRequestWithObject(t *testing.T) {
	var s SimpleHttpServer
	go s.Start()
	a := assert.New(t)
	h := NewHttpWrapper(time.Second)
	url := "http://127.0.0.1:8080"
	var ss, ss2 []string
	ss = []string{"12", "22"}
	rsp, err := h.SendRequest(url, ss, &ss2)
	a.NoError(err, "SendRequest error")
	a.EqualValues(0, rsp.Code)
	a.EqualValues("Detail", rsp.Detail)
	a.EqualValues("OK", rsp.Message)
	a.EqualValues(ss, ss2, "value should equal")
}

func TestSendRequestWithNil(t *testing.T) {
	var s SimpleHttpServer
	go s.Start()
	a := assert.New(t)
	h := NewHttpWrapper(time.Second)
	url := "http://127.0.0.1:8080"
	var ss []string
	ss = []string{"12", "22"}
	rsp, err := h.SendRequest(url, ss, nil)
	a.NoError(err, "SendRequest error")
	a.EqualValues(0, rsp.Code)
	a.EqualValues("Detail", rsp.Detail)
	a.EqualValues("OK", rsp.Message)
	a.EqualValues([]interface{}{"12", "22"}, rsp.Data, "value should equal")
}
