package zbs

import (
	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
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
	case opts.SubjectLevel <= 0:
		return nil, errors.New("SubjectLevel is required")
	}
	params := model.CreateSubjectParams{
		SubjectName:  &opts.SubjectName,
		SubjectLevel: &opts.SubjectLevel,
	}

	result, err := util.DoRequest(action, params)
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
	param := &model.ModifySubjectParams{
		SubjectID:    &opts.SubjectID,
		SubjectLevel: &opts.SubjectLevel,
		SubjectName:  &opts.SubjectName,
	}
	_, err := util.DoRequest(action, param)
	return err
}
