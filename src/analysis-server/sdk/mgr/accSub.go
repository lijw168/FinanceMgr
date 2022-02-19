package mgr

import (
	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
	"encoding/json"
	"errors"
	"fmt"
)

type AccSub struct {
}

func (as *AccSub) CreateAccSub(opts *options.CreateSubjectOptions) (*model.AccSubjectView, error) {
	action := "CreateAccSub"
	switch {
	case opts.SubjectName == "":
		return nil, errors.New("SubjectName is required")
	case opts.CommonID == "":
		return nil, errors.New("CommonID is required")
	case opts.SubjectLevel <= 0:
		return nil, errors.New("SubjectLevel is required")
	case opts.CompanyID <= 0:
		return nil, errors.New("CompanyID is required")
	case opts.SubjectDirection <= 0:
		return nil, errors.New("SubjectDirection is required")
	case opts.SubjectType <= 0:
		return nil, errors.New("SubjectType is required")
	case opts.SubjectStyle == "":
		return nil, errors.New("SubjectStyle is required")
	}
	params := model.CreateSubjectParams{
		SubjectName:      &opts.SubjectName,
		CommonID:         &opts.CommonID,
		SubjectLevel:     &opts.SubjectLevel,
		CompanyID:        &opts.CompanyID,
		SubjectType:      &opts.SubjectType,
		SubjectDirection: &opts.SubjectDirection,
		SubjectStyle:     &opts.SubjectStyle,
	}

	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.AccSubjectView{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (as *AccSub) CreateAccSub_json(params []byte) (*model.AccSubjectView, error) {
	action := "CreateAccSub"
	result, err := util.DoRequest_json(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.AccSubjectView{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (as *AccSub) DeleteAccSub(opts *options.BaseOptions) error {
	action := "DeleteAccSub"
	err := DeleteOpsResource(action, opts)
	if err != nil {
		return err
	}
	fmt.Printf("DeleteAccSub succeed")
	return nil
}

func (as *AccSub) GetAccSub(opts *options.BaseOptions) (*model.AccSubjectView, error) {
	action := "GetAccSub"
	dr, err := DescribeOpsResource(action, opts)
	if err != nil {
		return nil, err
	}
	view := &model.AccSubjectView{}
	err = util.FormatView(dr, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}

func (as *AccSub) ListAccSub_json(params []byte) ([]byte, error) {
	action := "ListAccSub"
	return ListOpsResources_json(action, params)
}

func (as *AccSub) ListAccSub(opts *options.ListOptions) (int64, []*model.AccSubjectView, error) {
	action := "ListAccSub"
	var ret []*model.AccSubjectView
	desc, err := ListOpsResources(action, opts)
	if err != nil {
		return -1, nil, err
	}
	if err := util.FormatView(desc.Elements, &ret); err != nil {
		return -1, nil, err
	}
	return desc.Tc, ret, nil
}

func (as *AccSub) UpdateAccSub(opts *options.ModifySubjectOptions) error {
	action := "UpdateAccSub"
	switch {
	case opts.SubjectID <= 0:
		return errors.New("SubjectID are required")
	}
	param := &model.ModifySubjectParams{}
	if opts.SubjectID != 0 {
		param.SubjectID = &opts.SubjectID
	}
	if opts.SubjectLevel != 0 {
		param.SubjectLevel = &opts.SubjectLevel
	}
	if opts.SubjectName != "" {
		param.SubjectName = &opts.SubjectName
	}
	if opts.CommonID != "" {
		param.CommonID = &opts.CommonID
	}
	if opts.CompanyID > 0 {
		param.CompanyID = &opts.CompanyID
	}
	if opts.SubjectDirection > 0 {
		param.SubjectDirection = &opts.SubjectDirection
	}
	if opts.SubjectType > 0 {
		param.SubjectType = &opts.SubjectType
	}
	_, err := util.DoRequest(action, param)
	return err
}

func (as *AccSub) UpdateAccSub_json(params []byte) error {
	action := "UpdateAccSub"
	_, err := util.DoRequest_json(action, params)
	return err
}

func (as *AccSub) QueryAccSubReference(opts *options.BaseOptions) (int64, error) {
	action := "QueryAccSubReference"
	switch {
	case opts.ID <= 0:
		return 0, errors.New("ID is required")
	}
	params := &model.BaseParams{
		ID: &opts.ID,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return 0, err
	}
	var iRefCount int64
	if err = util.FormatView(result.Data, &iRefCount); err != nil {
		return 0, err
	}
	return iRefCount, nil
	// if iRefCount, ok := result.Data.(int64); ok {
	// 	return iRefCount, nil
	// }
	//return 0, errors.New("the type of result.Data is wrong")
}

// func (as *AccSub) GetYearBalanceById_json(params []byte) ([]byte, error) {
// 	action := "GetYearBalance"
// 	result, err := util.DoRequest_json(action, params)
// 	if err != nil {
// 		return nil, err
// 	}
// 	balValue, err := json.Marshal(result.Data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return balValue, nil
// }

// func (as *AccSub) UpdateYearBalance_json(params []byte) error {
// 	action := "UpdateYearBalance"
// 	_, err := util.DoRequest_json(action, params)
// 	return err
// }

// func (as *AccSub) ListYearBalance_json(params []byte) ([]byte, error) {
// 	action := "ListYearBalance"
// 	return ListOpsResources_json(action, params)
// }

func (as *AccSub) CopyAccSubTemplate(opts *options.BaseOptions) ([]byte, error) {
	action := "CopyAccSubTemplate"
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
	return json.Marshal(result.Data)
}

func (as *AccSub) GenerateAccSubTemplate(opts *options.BaseOptions) error {
	action := "GenerateAccSubTemplate"
	switch {
	case opts.ID <= 0:
		return errors.New("ID is required")
	}
	params := &model.BaseParams{
		ID: &opts.ID,
	}
	_, err := util.DoRequest(action, params)
	if err != nil {
		return err
	}
	fmt.Printf("GenerateAccSubTemplate succeed")
	return nil
}
