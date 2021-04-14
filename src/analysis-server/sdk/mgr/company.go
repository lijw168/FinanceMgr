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
	}
	params := model.CreateCompanyParams{
		CompanyName: &opts.CompanyName,
		Phone:       &opts.Phone,
		AbbrevName:  &opts.AbbrevName,
		Backup:      &opts.Backup,
		CompanyAddr: &opts.CompanyAddr,
		Corporator:  &opts.Corporator,
		Email:       &opts.Email,
		Summary:     &opts.Summary,
	}
	result, err := util.DoRequest(action, params)
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
	action := "Listcompany"
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

func (c *Company) UpdateCompany(opts *options.ModifyCompanyOptions) error {
	action := "UpdateCompany"
	switch {
	case opts.CompanyID <= 0:
		return errors.New("CompanyID are required")
	}
	param := &model.ModifyCompanyParams{
		CompanyID:   &opts.CompanyID,
		CompanyName: &opts.CompanyName,
		Phone:       &opts.Phone,
		AbbrevName:  &opts.AbbrevName,
		Backup:      &opts.Backup,
		CompanyAddr: &opts.CompanyAddr,
		Corporator:  &opts.Corporator,
		Email:       &opts.Email,
		Summary:     &opts.Summary,
	}
	_, err := util.DoRequest(action, param)
	return err
}