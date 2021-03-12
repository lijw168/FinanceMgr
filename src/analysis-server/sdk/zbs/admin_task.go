package zbs

import (
	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
)

type At struct {
}

func (r *At) DescribeAt(opts *options.BaseOptions) (*model.AdminTaskView, error) {
	action := "DescribeAdminTask"
	bt, err := DescribeOpsResource(action, opts)
	if err != nil {
		return nil, err
	}
	view := &model.AdminTaskView{}
	err = util.FormatView(bt, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}
