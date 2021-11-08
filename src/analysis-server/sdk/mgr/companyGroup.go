package mgr

import (
	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
	"errors"
	"fmt"
)

type CompanyGroup struct {
}

func (c *CompanyGroup) CreateCompanyGroup(opts *options.CreateCompanyGroupOptions) (*model.CompanyGroupView, error) {
	action := "CreateCompanyGroup"
	switch {
	case opts.GroupName == "":
		return nil, errors.New("GroupName is required")
	case opts.GroupStatus < 0:
		return nil, errors.New("GroupStatus is required")
	}
	params := model.CreateCompanyGroupParams{
		GroupName:   &opts.GroupName,
		GroupStatus: &opts.GroupStatus,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.CompanyGroupView{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (c *CompanyGroup) DeleteCompanyGroup(opts *options.BaseOptions) error {
	action := "DeleteCompanyGroup"
	err := DeleteOpsResource(action, opts)
	if err != nil {
		return err
	}
	fmt.Printf("DeleteCompany succeed")
	return nil
}

func (c *CompanyGroup) GetCompanyGroup(opts *options.BaseOptions) (*model.CompanyGroupView, error) {
	action := "GetCompanyGroup"
	dr, err := DescribeOpsResource(action, opts)
	if err != nil {
		return nil, err
	}
	view := &model.CompanyGroupView{}
	err = util.FormatView(dr, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}

func (c *CompanyGroup) ListCompanyGroup(opts *options.ListOptions) (int64, []*model.CompanyGroupView, error) {
	action := "ListCompanyGroup"
	var ret []*model.CompanyGroupView
	desc, err := ListOpsResources(action, opts)
	if err != nil {
		return -1, nil, err
	}
	if err := util.FormatView(desc.Elements, &ret); err != nil {
		return -1, nil, err
	}
	return desc.Tc, ret, nil
}

func (c *CompanyGroup) UpdateCompanyGroup(opts *options.ModifyCompanyGroupOptions) error {
	action := "UpdateCompanyGroup"
	switch {
	case opts.CompanyGroupID <= 0:
		return errors.New("CompanyGroupID are required")
	}
	param := &model.ModifyCompanyGroupParams{}
	if opts.CompanyGroupID != 0 {
		param.CompanyGroupID = &opts.CompanyGroupID
	}
	if opts.GroupName != "" {
		param.GroupName = &opts.GroupName
	}
	if opts.GroupStatus >= 0 {
		param.GroupStatus = &opts.GroupStatus
	}
	_, err := util.DoRequest(action, param)
	return err
}
