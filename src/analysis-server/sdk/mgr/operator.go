package mgr

import (
	"errors"
	"financeMgr/src/analysis-server/model"
	"financeMgr/src/analysis-server/sdk/options"
	"financeMgr/src/analysis-server/sdk/util"
	"fmt"
)

type Operator struct {
}

func (or *Operator) CreateOperator(opts *options.CreateOptInfoOptions) (*model.OperatorInfoView, error) {
	action := "CreateOperator"
	switch {
	case opts.CompanyID <= 0:
		return nil, errors.New("CompanyID is required")
	case opts.Name == "":
		return nil, errors.New("Name is required")
	case opts.Password == "":
		return nil, errors.New("Password is required")
	}
	params := model.CreateOptInfoParams{
		CompanyID:  &opts.CompanyID,
		Name:       &opts.Name,
		Password:   &opts.Password,
		Department: &opts.Department,
		Job:        &opts.Job,
		Role:       &opts.Role,
		Status:     &opts.Status,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.OperatorInfoView{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (or *Operator) CreateOperator_json(params []byte) (*model.OperatorInfoView, error) {
	action := "CreateOperator"
	result, err := util.DoRequest_json(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.OperatorInfoView{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (or *Operator) DeleteOperator(opts *options.BaseOptions) error {
	action := "DeleteOperator"
	err := DeleteOpsResource(action, opts)
	if err != nil {
		return err
	}
	fmt.Printf("DeleteOperator succeed")
	return nil
}

func (or *Operator) GetOperatorInfo(opts *options.BaseOptions) (*model.OperatorInfoView, error) {
	action := "GetOperatorInfo"
	dr, err := DescribeOpsResource(action, opts)
	if err != nil {
		return nil, err
	}
	view := &model.OperatorInfoView{}
	err = util.FormatView(dr, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}

func (or *Operator) ListOperatorInfo(opts *options.ListOptions) (int64, []*model.OperatorInfoView, error) {
	action := "ListOperatorInfo"
	var ret []*model.OperatorInfoView
	desc, err := ListOpsResources(action, opts)
	if err != nil {
		return -1, nil, err
	}
	if err := util.FormatView(desc.Elements, &ret); err != nil {
		return -1, nil, err
	}
	return desc.Tc, ret, nil
}

func (or *Operator) ListOperatorInfo_json(params []byte) ([]byte, error) {
	action := "ListOperatorInfo"
	return ListOpsResources_json(action, params)
}

func (or *Operator) UpdateOperator(opts *options.ModifyOptInfoOptions) error {
	action := "UpdateOperator"

	if opts.OperatorID <= 0 {
		return errors.New("OperatorID is required")
	}
	param := model.ModifyOptInfoParams{}
	if opts.Name != "" {
		param.Name = &opts.Name
	}
	if opts.Password != "" {
		param.Password = &opts.Password
	}
	if opts.Department != "" {
		param.Department = &opts.Department
	}
	if opts.Job != "" {
		param.Job = &opts.Job
	}
	if opts.Role != 0 {
		param.Role = &opts.Role
	}
	if opts.Status != 0 {
		param.Status = &opts.Status
	}
	_, err := util.DoRequest(action, param)
	return err
}

func (or *Operator) UpdateOperator_json(param []byte) error {
	action := "UpdateOperator"
	_, err := util.DoRequest_json(action, param)
	return err
}
