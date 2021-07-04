package mgr

import (
	"errors"

	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
	"encoding/json"
)

// func DescribeTenantResource(action string, opts *options.BaseOptions) (interface{}, error) {
// 	switch {
// 	case opts.Id == "":
// 		return nil, errors.New("ID is required")
// 	}
// 	params := &model.TenantBaseParams{
// 		TenantId: &util.Tenant,
// 	}
// 	params.Id = &opts.Id
// 	result, err := util.DoRequest(action, params)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result.Data, nil
// }

// func ListTenantResources(action string, opts *options.ListOptions) (*model.DescData, error) {
// 	para := &model.ListTenantResourceParams{
// 		TenantId: &util.Tenant,
// 	}
// 	para.DescLimit = &opts.Limit
// 	para.DescOffset = &opts.Offset
// 	for k, v := range opts.Filter {
// 		kv := k
// 		filterP := &model.FilterItem{Field: &kv, Value: v}
// 		para.Filter = append(para.Filter, filterP)
// 	}
// 	for k, v := range opts.Orders {
// 		orderKey := k
// 		orderV := v
// 		orderP := &model.OrderItem{Field: &orderKey, Direction: &orderV}
// 		para.Order = append(para.Order, orderP)
// 	}

// 	result, err := util.DoRequest(action, para)
// 	if err != nil {
// 		return nil, err
// 	}
// 	desc := &(model.DescData{})
// 	if err := util.FormatView(result.Data, desc); err != nil {
// 		return nil, err
// 	}
// 	return desc, nil
// }

func ListOpsResources(action string, opts *options.ListOptions) (*model.DescData, error) {
	para := &model.ListParams{
		DescLimit:  &opts.Limit,
		DescOffset: &opts.Offset,
	}
	for k, v := range opts.Filter {
		kv := k
		filterP := &model.FilterItem{Field: &kv, Value: v}
		para.Filter = append(para.Filter, filterP)
	}
	for k, v := range opts.Orders {
		orderKey := k
		orderV := v
		orderP := &model.OrderItem{Field: &orderKey, Direction: &orderV}
		para.Order = append(para.Order, orderP)
	}

	result, err := util.DoRequest(action, para)
	if err != nil {
		return nil, err
	}
	desc := &(model.DescData{})
	if err := util.FormatView(result.Data, desc); err != nil {
		return nil, err
	}
	return desc, nil
}

func ListOpsResources_json(action string, params []byte) ([]byte, error) {
	result, err := util.DoRequest_json(action, params)
	if err != nil {
		return nil, err
	}
	return json.Marshal(result.Data)

	// desc := &(model.DescData{})
	// if err := util.FormatView(result.Data, desc); err != nil {
	// 	return nil, err
	// }
	//return b, nil
}

func DeleteOpsResource(action string, opts *options.BaseOptions) error {
	switch {
	case opts.ID <= 0:
		return errors.New("ID is required")
	}
	params := &model.BaseParams{
		ID: &opts.ID,
	}
	_, err := util.DoRequest(action, params)
	return err
}

func DescribeOpsResource(action string, opts *options.BaseOptions) (interface{}, error) {
	switch {
	case opts.ID <= 0:
		return nil, errors.New("ID is required")
	}
	params := &model.BaseParams{
		ID: &opts.ID,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}
