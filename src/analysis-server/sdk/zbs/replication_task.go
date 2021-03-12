package zbs

import (
	"errors"

	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
)

type Rt struct {
}

func (s *Rt) Rescheduler(opts *options.ReschedulerOptions) (*model.ReschedulerView, error) {
	action := "Reschedule"
	switch {
	case opts.PoolId == "":
		return nil, errors.New("PoolId is required")
	}
	params := model.ReschedulerParams{
		PoolId: &opts.PoolId,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.ReschedulerView{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (s *Rt) TransferLeader(opts *options.TransferLeaderOptions) (*model.TransferLeaderView, error) {
	action := "TransferLeader"
	switch {
	case opts.RgId == "":
		return nil, errors.New("RgId is required")
	case opts.ReplicaId == "":
		return nil, errors.New("ReplicaId is required")
	}
	params := model.TransferLeaderParams{
		RgId:      &opts.RgId,
		ReplicaId: &opts.ReplicaId,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.TransferLeaderView{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (s *Rt) MoveReplica(opts *options.MoveReplicaOptions) (*model.MoveReplicaView, error) {
	action := "MoveReplica"
	switch {
	case opts.ReplicaId == "":
		return nil, errors.New("ReplicaId is required")
	case opts.TargetDiskId == "":
		return nil, errors.New("TargetDiskId is required")
	}
	params := model.MoveReplicaParams{
		ReplicaId:    &opts.ReplicaId,
		TargetDiskId: &opts.TargetDiskId,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.MoveReplicaView{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (s *Rt) DeleteReplica(opts *options.DeleteReplicaOptions) (
	*model.DeleteReplicaView, error) {

	action := "DeleteReplica"
	switch {
	case opts.ReplicaId == "":
		return nil, errors.New("ReplicaId is required")
	}
	params := model.DeleteReplicaParams{
		ReplicaId: &opts.ReplicaId,
		Force:     &opts.Force,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.DeleteReplicaView{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (s *Rt) AddReplica(opts *options.AddReplicaOptions) (
	*model.AddReplicaView, error) {

	action := "AddReplica"
	switch {
	case opts.RgID == "":
		return nil, errors.New("RgID is required")
	case opts.DiskID == "":
		return nil, errors.New("DiskID is required")
	}
	params := model.AddReplicaParams{
		RgID:   &opts.RgID,
		DiskID: &opts.DiskID,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.AddReplicaView{}
	util.FormatView(result.Data, &view)
	return view, nil
}
