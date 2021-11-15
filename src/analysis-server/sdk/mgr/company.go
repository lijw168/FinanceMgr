package mgr

import (
	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
	"errors"
	"fmt"
)

type Company struct {
}

func (c *Company) CreateCompany(opts *options.CreateCompanyOptions) (*model.CompanyView, error) {
	action := "CreateCompany"
	switch {
	case opts.CompanyName == "":
		return nil, errors.New("CompanyName is required")
	case opts.Phone == "":
		return nil, errors.New("Phone is required")
	case opts.StartAccountPeriod <= 0:
		return nil, errors.New("StartAccountPeriod is required")
	}
	params := model.CreateCompanyParams{
		CompanyName:        &opts.CompanyName,
		Phone:              &opts.Phone,
		AbbrevName:         &opts.AbbrevName,
		Backup:             &opts.Backup,
		CompanyAddr:        &opts.CompanyAddr,
		Corporator:         &opts.Corporator,
		StartAccountPeriod: &opts.StartAccountPeriod,
		Email:              &opts.Email,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.CompanyView{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (c *Company) CreateCompany_json(params []byte) (*model.CompanyView, error) {
	action := "CreateCompany"
	result, err := util.DoRequest_json(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.CompanyView{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (c *Company) DeleteCompany(opts *options.BaseOptions) error {
	action := "DeleteCompany"
	err := DeleteOpsResource(action, opts)
	if err != nil {
		return err
	}
	fmt.Printf("DeleteCompany succeed")
	return nil
}

func (c *Company) GetCompany(opts *options.BaseOptions) (*model.CompanyView, error) {
	action := "GetCompany"
	dr, err := DescribeOpsResource(action, opts)
	if err != nil {
		return nil, err
	}
	view := &model.CompanyView{}
	err = util.FormatView(dr, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}

func (c *Company) ListCompany(opts *options.ListOptions) (int64, []*model.CompanyView, error) {
	action := "ListCompany"
	var ret []*model.CompanyView
	desc, err := ListOpsResources(action, opts)
	if err != nil {
		return -1, nil, err
	}
	if err := util.FormatView(desc.Elements, &ret); err != nil {
		return -1, nil, err
	}
	return desc.Tc, ret, nil
}

func (c *Company) ListCompany_json(params []byte) ([]byte, error) {
	action := "ListCompany"
	return ListOpsResources_json(action, params)
}

func (c *Company) UpdateCompany(opts *options.ModifyCompanyOptions) error {
	action := "UpdateCompany"
	switch {
	case opts.CompanyID <= 0:
		return errors.New("CompanyID are required")
	}
	param := &model.ModifyCompanyParams{}
	if opts.CompanyID != 0 {
		param.CompanyID = &opts.CompanyID
	}
	if opts.CompanyName != "" {
		param.CompanyName = &opts.CompanyName
	}
	if opts.Phone != "" {
		param.Phone = &opts.Phone
	}
	if opts.AbbrevName != "" {
		param.AbbrevName = &opts.AbbrevName
	}
	if opts.Backup != "" {
		param.Backup = &opts.Backup
	}

	if opts.CompanyAddr != "" {
		param.CompanyAddr = &opts.CompanyAddr
	}
	if opts.Corporator != "" {
		param.Corporator = &opts.Corporator
	}
	if opts.Email != "" {
		param.Email = &opts.Email
	}
	if opts.LatestAccountYear > 0 {
		param.LatestAccountYear = &opts.LatestAccountYear
	}
	_, err := util.DoRequest(action, param)
	return err
}

func (c *Company) UpdateCompany_json(param []byte) error {
	action := "UpdateCompany"
	_, err := util.DoRequest_json(action, param)
	return err
}

func (c *Company) AssociatedCompanyGroup(opts *options.AssociatedCompanyGroupOptions) error {
	action := "AssociatedCompanyGroup"
	switch {
	case opts.CompanyGroupID <= 0:
		return errors.New("CompanyGroupID is required")
	case opts.CompanyID <= 0:
		return errors.New("CompanyID is required")
	}
	params := model.AssociatedCompanyGroupParams{
		CompanyGroupID: &opts.CompanyGroupID,
		CompanyID:      &opts.CompanyID,
		IsAttach:       &opts.IsAttach,
	}
	_, err := util.DoRequest(action, params)
	return err
}

func (c *Company) InitResourceInfo(opts *options.BaseOptions) (*model.DescData, error) {
	action := "InitResourceInfo"
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
	desc := &(model.DescData{})
	if err := util.FormatView(result.Data, desc); err != nil {
		return nil, err
	}
	return desc, nil
}
