package zbs

import (
	"errors"
	"common/message"
	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
)

type ZbsProxy struct {
}

func (c *ZbsProxy) AddZbsProxy(opts *options.ZbsProxyOptions) (*message.ZbsProxyList, error) {
	action := "AddZbsProxy"
	switch {
	case opts.Addr == "":
		return nil, errors.New("addr is required")
	}
	params := model.ZbsProxyParams{
		Addr: &opts.Addr,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &message.ZbsProxyList{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (c *ZbsProxy) DeleteZbsProxy(opts *options.ZbsProxyOptions) (*message.ZbsProxyList, error) {
	action := "DeleteZbsProxy"
	switch {
	case opts.Addr == "":
		return nil, errors.New("addr is required")
	}
	params := model.ZbsProxyParams{
		Addr: &opts.Addr,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &message.ZbsProxyList{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (c *ZbsProxy) ListZbsProxy(opts *options.ZbsProxyOptions) (*message.ZbsProxyList, error) {
	action := "ListZbsProxy"
	params := model.ZbsProxyParams{
		Addr: &opts.Addr,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &message.ZbsProxyList{}
	util.FormatView(result.Data, &view)
	return view, nil
}
