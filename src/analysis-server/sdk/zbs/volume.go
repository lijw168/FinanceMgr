package zbs

import (
	"errors"

	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
)

type Volume struct {
}

func (v *Volume) CreateVolume(opts *options.CreateVolumeOptions) (*model.VolumeView, error) {
	if opts.VolumeTypeName == "" || opts.Size <= 0 || opts.TenantId == "" || opts.VolumeName == "" || opts.AzName == "" {
		return nil, errors.New("VolumeId, TenantId, Size ,volumeTypeName ,AzName  are required")
	}
	action := "CreateVolume"
	para := &model.CreateVolumeParams{
		VolumeName:     &opts.VolumeName,
		TenantId:       &opts.TenantId,
		Size:           &opts.Size,
		VolumeTypeName: &opts.VolumeTypeName,
		Description:    &opts.Description,
		SnapshotId:     &opts.SnapshotId,
		AzName:         &opts.AzName,
		Bootable:       &opts.Bootable,
	}
	result, err := util.DoRequestwithToken(opts.ClientToken, action, para)
	if err != nil {
		return nil, err
	}
	volumeView := new(model.VolumeView)
	util.FormatView(result.Data, volumeView)
	return volumeView, nil
}
func (v *Volume) ChangeVolume(opts *options.ChangeVolumeOptions) error {
	if opts.TenantId == "" || opts.VolumeId == "" {
		return errors.New("VolumeId, TenantId are required")
	}
	action := "ChangeVolume"
	para := &model.ChangeVolumeParams{
		TenantId:    &opts.TenantId,
		VolumeId:    &opts.VolumeId,
		VolumeName:  opts.VolumeName,
		Description: opts.Description,
	}
	_, err := util.DoRequest(action, para)
	if err != nil {
		return err
	}
	return nil
}

func (v *Volume) ResizeVolume(opts *options.ResizeVolumeOptions) (*model.VolumeView, error) {
	if opts.VolumeId == "" || opts.Size <= 0 || opts.TenantId == "" {
		return nil, errors.New("VolumeId, TenantId, Size  are required")
	}
	action := "ResizeVolume"
	para := &model.ResizeVolumeParams{
		TenantId: &opts.TenantId,
		Size:     &opts.Size,
		VolumeId: &opts.VolumeId,
	}
	result, err := util.DoRequest(action, para)
	if err != nil {
		return nil, err
	}
	volumeView := new(model.VolumeView)
	util.FormatView(result.Data, volumeView)
	return volumeView, nil

}
func (v *Volume) DeleteVolume(opts *options.DeleteVolumeOptions) error {
	if opts.TenantId == "" || opts.VolumeId == "" {
		return errors.New("VolumeId, TenantId are required")
	}
	action := "DeleteVolume"
	para := &model.DeleteVolumeParams{
		TenantId: &opts.TenantId,
		Id:       &opts.VolumeId,
	}
	_, err := util.DoRequest(action, para)
	if err != nil {
		return err
	}
	return nil
}

func (v *Volume) DescribeVolume(opts *options.DescribeVolumeOptions) (*model.VolumeView, error) {
	action := "DescribeVolume"
	para := &model.DescribeVolumeParams{
		TenantId: &opts.TenantId,
		Id:       &opts.VolumeId,
	}
	result, err := util.DoRequest(action, para)
	if err != nil {
		return nil, err
	}
	volumeView := new(model.VolumeView)
	util.FormatView(result.Data, volumeView)
	return volumeView, nil
}

//return list
func (v *Volume) ListVolumes(opts *options.ListOptions) (int64, []*model.VolumeView, error) {
	action := "ListVolumes"
	var volumeViewSlice []*model.VolumeView
	desc, err := ListTenantResources(action, opts)
	if err != nil {
		return -1, nil, err
	}
	if err := util.FormatView(desc.Elements, &volumeViewSlice); err != nil {
		return -1, nil, err
	}
	return desc.Tc, volumeViewSlice, nil
}

func (v *Volume) AttachVolume(opts *options.AttachVolumeOptions) (*model.AttachResultView, error) {
	action := "AttachVolume"
	para := &model.AttachVolumeParams{
		VolumeId:     &opts.VolumeId,
		HostIp:       &opts.HostIp,
		InstanceUuid: &opts.InstanceUuid,
		InstanceType: &opts.InstanceType,
		Multiple:     &opts.MultiAttachment,
		TenantId:     &opts.TenantId,
		AttachMode:   &opts.AttachMode,
	}
	result, err := util.DoRequest(action, para)
	if err != nil {
		return nil, err
	}
	attachResultView := new(model.AttachResultView)
	util.FormatView(result.Data, attachResultView)
	return attachResultView, nil
}

func (v *Volume) DetachVolume(opts *options.DetachVolumeOptions) error {
	action := "DetachVolume"
	para := &model.DetachVolumeParams{
		VolumeId:     &opts.VolumeId,
		TenantId:     &opts.TenantId,
		AttachmentId: &opts.AttachmentId,
	}
	_, err := util.DoRequest(action, para)
	return err
}
func (v *Volume) DescribeAttachment(opts *options.DescribeAttachmentOptions) (*model.AttachmentView, error) {
	action := "DescribeAttachment"
	para := &model.DescribeAttachmentParams{
		TenantId:     &opts.TenantId,
		AttachmentId: &opts.AttachmentId,
	}
	result, err := util.DoRequest(action, para)
	if err != nil {
		return nil, err
	}
	volumeAttView := new(model.AttachmentView)
	util.FormatView(result.Data, volumeAttView)
	return volumeAttView, nil
}

//return list
func (v *Volume) ListAttachments(opts *options.ListOptions) (int64, []*model.AttachmentView, error) {
	action := "ListAttachments"
	var volumeAttViewSlice []*model.AttachmentView
	desc, err := ListTenantResources(action, opts)
	if err != nil {
		return -1, nil, err
	}
	if err := util.FormatView(desc.Elements, &volumeAttViewSlice); err != nil {
		return -1, nil, err
	}
	return desc.Tc, volumeAttViewSlice, nil
}

func (v *Volume) ResumeVolume(opts *options.AdminVolumeOptions) error {
	action := "ResumeVolume"
	para := &model.AdminVolumeParams{
		VolumeId: opts.VolumeId,
		TenantId: opts.TenantId,
	}
	_, err := util.DoRequest(action, para)
	return err
}

func (v *Volume) StopVolume(opts *options.AdminVolumeOptions) error {
	action := "StopVolume"
	para := &model.AdminVolumeParams{
		VolumeId: opts.VolumeId,
		TenantId: opts.TenantId,
	}
	_, err := util.DoRequest(action, para)
	return err
}
