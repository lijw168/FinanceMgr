package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"common/url"
	"jd.com/zbs/zbs-openapi/model"
)

type OpenHandler struct{}

func (h *OpenHandler) Response(ctx context.Context, w http.ResponseWriter, data interface{}) {
	var rsp interface{}
	requestId := ""
	if ctx.Value("trace_id") != nil {
		requestId = ctx.Value("trace_id").(string)
	}
	if data == nil {
		rsp = model.NullRespData{RequestId: requestId}
	} else {
		rsp = data
	}
	jsonRsp, err := json.Marshal(&rsp)
	if err != nil {
		//log.InfoContext(ctx, "[init/handler/Response] [response: %s, error: %s]", string(jsonRsp), err.Error())
		fmt.Printf("Marshal failed: the error is :%s", err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRsp)
}

type SnapshotHandler struct {
	OpenHandler
}

func (sh *SnapshotHandler) Register(rootRouter *url.RestUrlRouter) {
	router1 := rootRouter.AddOjbect("/snapshots")
	router1.Path("/").Methods("get").HandlerFunc(sh.DescribeSnapshots)
	router1.Path("/").Methods("post").HandlerFunc(sh.CreateSnapshot)
	router2 := rootRouter.AddOjbect("/snapshots/{snapshotId}")
	router2.Path("/").Methods("get").HandlerFunc(sh.DescribeSnapshot)
	router2.Path("/").Methods("delete").HandlerFunc(sh.DeleteSnapshot)
}

func (sh *SnapshotHandler) CreateSnapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		fmt.Println("the Vars is empty")
	} else {
		fmt.Println("the Vars is:", vars)
	}
	fmt.Println(vars["regionId"])
	fmt.Println(vars["snapshotId"])
	fmt.Println("SnapshotHandler.CreateSnapshot")
	sh.Response(r.Context(), w, model.RespData{RequestId: r.Context().Value("trace_id").(string), Result: "CreateSnapshot"})
}

func (sh *SnapshotHandler) DescribeSnapshot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SnapshotHandler.DescribeSnapshot")
	sh.Response(r.Context(), w, model.RespData{RequestId: r.Context().Value("trace_id").(string), Result: "DescribeSnapshot"})
}
func (sh *SnapshotHandler) DeleteSnapshot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SnapshotHandler.DeleteSnapshot")
	sh.Response(r.Context(), w, model.RespData{RequestId: r.Context().Value("trace_id").(string), Result: "DeleteSnapshot"})
}

func (sh *SnapshotHandler) DescribeSnapshots(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SnapshotHandler.DescribeSnapshots")
	sh.Response(r.Context(), w, model.RespData{RequestId: r.Context().Value("trace_id").(string), Result: "DescribeSnapshots"})
}
