package zbs

import (
	"errors"
	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
)

type Proxy struct {
}

func (s *Proxy) CreateProxy(opts *options.CreateProxyOptions) (*model.ProxyView, error) {
	action := "CreateProxy"
	switch {
	case opts.Addr == "":
		return nil, errors.New("addr is required")
	}
	params := model.CreateProxyParams{
		Addr: &opts.Addr,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.ProxyView{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (s *Proxy) DeleteProxy(opts *options.BaseOptions) error {
	action := "DeleteProxy"
	return DeleteOpsResource(action, opts)
}

func (s *Proxy) DescribeProxy(opts *options.BaseOptions) (*model.ProxyView, error) {
	action := "DescribeProxy"
	bt, err := DescribeOpsResource(action, opts)
	if err != nil {
		return nil, err
	}
	view := &model.ProxyView{}
	err = util.FormatView(bt, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}

func (z *Proxy) DescribeProxys(opts *options.ListOptions) (int64, []*model.ProxyView, error) {
	action := "DescribeProxys"
	var ret []*model.ProxyView
	desc, err := ListOpsResources(action, opts)
	if err != nil {
		return -1, nil, err
	}
	if err := util.FormatView(desc.Elements, &ret); err != nil {
		return -1, nil, err
	}
	return desc.Tc, ret, nil
}
