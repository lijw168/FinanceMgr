package mgr

import (
	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
	"errors"
	//"fmt"
)

type Authen struct {
}

func (au *Authen) Login(opts *options.AuthenInfoOptions) (*model.LoginInfoView, error) {
	action := "Login"
	switch {
	case opts.CompanyID <= 0:
		return nil, errors.New("CompanyID is required")
	case opts.Name == "":
		return nil, errors.New("Name is required")
	case opts.Password == "":
		return nil, errors.New("Password is required")
	}
	params := model.AuthenInfoParams{
		CompanyID: &opts.CompanyID,
		Name:      &opts.Name,
		Password:  &opts.Password,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.LoginInfoView{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (au *Authen) Login_json(params []byte) (*model.LoginInfoView, error) {
	action := "Login"
	result, err := util.DoRequest_json(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.LoginInfoView{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (au *Authen) Logout(opts *options.BaseOptions) error {
	action := "Logout"
	if opts.ID <= 0 {
		return errors.New("operatorID is required")
	}
	params := &model.BaseParams{
		ID: &opts.ID,
	}
	_, err := util.DoRequest(action, params)
	return err
}

func (au *Authen) StatusCheckout(opts *options.BaseOptions) (*model.StatusCheckoutView, error) {
	action := "StatusCheckout"
	if opts.ID <= 0 {
		return nil, errors.New("operatorID is required")
	}
	params := &model.BaseParams{
		ID: &opts.ID,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.StatusCheckoutView{}
	err = util.FormatView(result.Data, &view)
	if err != nil {
		return nil, err
	}
	return view, nil
}

func (au *Authen) ListLoginInfo(opts *options.ListOptions) (int64, []*model.LoginInfoView, error) {
	action := "ListLoginInfo"
	var ret []*model.LoginInfoView
	desc, err := ListOpsResources(action, opts)
	if err != nil {
		return -1, nil, err
	}
	if err := util.FormatView(desc.Elements, &ret); err != nil {
		return -1, nil, err
	}
	return desc.Tc, ret, nil
}

func (au *Authen) ListLoginInfo_json(params []byte) ([]byte, error) {
	action := "ListLoginInfo"
	return ListOpsResources_json(action, params)
}
