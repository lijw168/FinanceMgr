package zbs

import (
	"errors"
	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
)

type Disk struct {
}

func (z *Disk) CreateDisk(opts *options.CreateDiskOptions) (*model.DiskView, error) {
	action := "CreateDisk"
	switch {
	case opts.DeviceId == "":
		return nil, errors.New("DeviceId is required")
	}
	params := model.CreateDiskParams{
		DeviceId:    &opts.DeviceId,
		HostId:      &opts.HostId,
		ManageAddr:  &opts.ManageAddr,
		StorageAddr: &opts.StorageAddr,
		ClientAddr:  &opts.ClientAddr,
		TraceAddr:   &opts.TraceAddr,
		VolTypeName: &opts.VolumeTypeName,
		Capacity:    &opts.Capacity,
		AdminStatus: &opts.AdminStatus,
	}

	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.DiskView{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (s *Disk) DeleteDisk(opts *options.DeleteOptions) (*model.AtDeleteView, error) {
	if opts.Id == "" {
		return nil, errors.New("Id are required")
	}
	action := "DeleteDisk"
	para := &model.DeletePoolParams{
		Id: &opts.Id,
	}
	dv, err := util.DoRequest(action, para)
	if err != nil {
		return nil, err
	}

	view := &model.AtDeleteView{}
	err = util.FormatView(dv.Data, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}

func (s *Disk) DescribeDisk(opts *options.BaseOptions) (*model.DiskView, error) {
	action := "DescribeDisk"
	bt, err := DescribeOpsResource(action, opts)
	if err != nil {
		return nil, err
	}
	view := &model.DiskView{}
	err = util.FormatView(bt, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}

func (z *Disk) DescribeDisks(opts *options.ListOptions) (int64, []*model.DiskView, error) {
	action := "DescribeDisks"
	var ret []*model.DiskView
	desc, err := ListOpsResources(action, opts)
	if err != nil {
		return -1, nil, err
	}
	if err := util.FormatView(desc.Elements, &ret); err != nil {
		return -1, nil, err
	}
	return desc.Tc, ret, nil
}

func (s *Disk) MarkDownDisk(opts *options.DeleteOptions) (*model.AtDeleteView, error) {
	if opts.Id == "" {
		return nil, errors.New("Id are required")
	}
	action := "MarkDownDisk"
	para := &model.DeletePoolParams{
		Id: &opts.Id,
	}
	dv, err := util.DoRequest(action, para)
	if err != nil {
		return nil, err
	}

	view := &model.AtDeleteView{}
	err = util.FormatView(dv.Data, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}

func (s *Disk) MarkUpDisk(opts *options.DeleteOptions) (*model.AtDeleteView, error) {
	if opts.Id == "" {
		return nil, errors.New("Id are required")
	}
	action := "MarkUpDisk"
	para := &model.DeletePoolParams{
		Id: &opts.Id,
	}
	dv, err := util.DoRequest(action, para)
	if err != nil {
		return nil, err
	}

	view := &model.AtDeleteView{}
	err = util.FormatView(dv.Data, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}
