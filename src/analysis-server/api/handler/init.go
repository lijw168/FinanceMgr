package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"analysis-server/api/service"
	"analysis-server/model"
	"common/log"
)

var GAccessTokenH = NewAccessTokenHandler()

type DescData struct {
	Tc       int64       `json:"total_count"`
	Elements interface{} `json:"elements"`
}

type CCHandler struct{}

func (h *CCHandler) Response(ctx context.Context, logger log.ILog, w http.ResponseWriter, ce service.CcError, data interface{}) {
	rsp := model.CommResp{Code: 0, Data: data}
	if ce != nil {
		rsp.Code = ce.GetCode()
		rsp.Message = ce.Error()
		rsp.Detail = ce.Detail()
	}

	jsonRsp, _ := json.Marshal(&rsp)
	if ce != nil {
		logger.InfoContext(ctx, "[init/handler/Response] [response: %s, error: %s]", string(jsonRsp), ce.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRsp)
}

func (h *CCHandler) GetTraceId(r *http.Request) string {
	if traceid := r.Header.Get("Client-Token"); len(traceid) > 0 {
		return traceid
	} else if traceid = r.Header.Get("Request-Id"); len(traceid) > 0 {
		return traceid
	} else {
		traceid = r.Header.Get("Trace-Id")
		return traceid
	}
}

func (h *CCHandler) HttpRequestParse(r *http.Request, param interface{}) error {
	jsonReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonReq, &param)
	if err != nil {
		return err
	}
	return nil
}

func (h *CCHandler) Success(ctx context.Context, logger *log.Logger, w http.ResponseWriter, data interface{}) {
	h.Response(ctx, logger, w, nil, data)
}

func (h *CCHandler) IsUnSet(v reflect.Value, exclude map[string]bool) (bool, string) {
	switch v.Kind() {
	case reflect.Func:
		return v.IsNil(), ""
	case reflect.Map:
		if v.IsNil() {
			return true, ""
		}
		for _, k := range v.MapKeys() {
			result, field := h.IsUnSet(v.MapIndex(k), exclude)
			if result {
				return result, field
			}
		}
		return false, ""
	case reflect.Array, reflect.Slice:
		if v.Len() == 0 {
			return true, ""
		}
		for i := 0; i < v.Len(); i++ {
			result, field := h.IsUnSet(v.Index(i), exclude)
			if result {
				return true, field
			}
		}
		return false, ""
	case reflect.Struct:
		structFields := v.Type()
		for i := 0; i < v.NumField(); i++ {
			fieldName := structFields.Field(i).Name
			if _, ok := exclude[fieldName]; !ok {
				if v.Field(i).CanSet() {
					result, field := h.IsUnSet(v.Field(i), exclude)
					if result {
						if field != "" {
							return true, field
						}

						return true, structFields.Field(i).Tag.Get("json")
					}
				}
			}
		}
		return false, ""
	case reflect.Ptr:
		if v.IsNil() {
			return true, ""
		}
		return h.IsUnSet(reflect.Indirect(v), exclude)
	case reflect.String:
		return v.Len() == 0, ""
	}

	return false, ""
}

func (h *CCHandler) Check(p interface{}, exclude map[string]bool) error {
	if reflect.ValueOf(p).Kind() != reflect.Ptr {
		return errors.New("must be ptr")
	}

	object := reflect.ValueOf(p).Elem()
	result, fieldName := h.IsUnSet(object, exclude)
	if result {
		if fieldName == "" {
			fieldName = "param"
		}
		return errors.New(strings.ToLower(fieldName))
	}
	return nil
}
