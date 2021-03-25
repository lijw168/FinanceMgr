package zbs

import (
	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
	"errors"
	"fmt"
)

type Operator struct {
}

func (or *Operator) CreateOperator(opts *options.OperatorInfoOptions) (*model.OperatorInfoView, error) {
	action := "CreateOperator"
	switch {
	case opts.CompanyID <= 0:
		return nil, errors.New("CompanyID is required")
	case opts.Name == "":
		return nil, errors.New("Name is required")
	case opts.Password == "":
		return nil, errors.New("Password is required")
	}
	params := model.OperatorInfoParams{
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

func (or *Operator) DeleteOperator(opts *options.NameOptions) error {
	action := "DeleteOperator"
	switch {
	case opts.Name == "":
		return errors.New("Name is required")
	}
	params := model.DeleteOperatorParams{Name: &opts.Name}
	_, err := util.DoRequest(action, params)
	if err != nil {
		return err
	}
	fmt.Printf("DeleteOperator succeed")
	return nil
}

func (or *Operator) GetOperatorInfo(opts *options.NameOptions) (*model.OperatorInfoView, error) {
	action := "GetOperatorInfo"
	switch {
	case opts.Name == "":
		return nil, errors.New("Name is required")
	}
	params := model.DescribeNameParams{Name: &opts.Name}
	opt, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.OperatorInfoView{}
	err = util.FormatView(opt, view)
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

func (or *Operator) UpdateOperator(opts *options.OperatorInfoOptions) error {
	action := "UpdateOperator"
	switch {
	case opts.CompanyID <= 0:
		return errors.New("CompanyID is required")
	case opts.Name == "":
		return errors.New("Name is required")
	}
	params := model.OperatorInfoParams{
		CompanyID:  &opts.CompanyID,
		Name:       &opts.Name,
		Password:   &opts.Password,
		Department: &opts.Department,
		Job:        &opts.Job,
		Role:       &opts.Role,
		Status:     &opts.Status,
	}
	_, err := util.DoRequest(action, params)
	return err
}
