package zbs

import (
	"errors"

	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
)

type Snapshot struct {
}

func (s *Snapshot) CreateSnapshot(opts *options.CreateSnapshotOptions) (*model.SnapshotView, error) {
	if opts.VolumeId == "" || opts.TenantId == "" || opts.SnapshotName == "" {
		return nil, errors.New("VolumeId, TenantId, SnapshotName are required")
	}
	action := "CreateSnapshot"
	para := &model.CreateSnapshotParams{
		SnapshotName: &opts.SnapshotName,
		TenantId:     &opts.TenantId,
		VolumeId:     &opts.VolumeId,
		Description:  &opts.Description,
	}
	result, err := util.DoRequestwithToken(opts.ClientToken, action, para)
	if err != nil {
		return nil, err
	}
	snapshotView := new(model.SnapshotView)
	util.FormatView(result.Data, snapshotView)
	return snapshotView, nil

}

func (s *Snapshot) ConvertImageToSnapshot(opts *options.ConvertImageToSnapshotOptions) (*model.SnapshotView, error) {
	if opts.ImageId == "" || opts.TenantId == "" || opts.SnapshotName == "" ||
		opts.ImageFormat == "" || opts.ImageLocation == "" || opts.ImageType == "" ||
		opts.ImageHash == "" {
		return nil, errors.New("ImageId, TenantId, SnapshotName, ImageFormat ,ImageLocation,ImageType ,ImageHash are required")
	}
	action := "ConvertImageToSnapshot"
	para := &model.ConvertImageToSnapshotParams{
		SnapshotName:  &opts.SnapshotName,
		TenantId:      &opts.TenantId,
		ImageId:       &opts.ImageId,
		ImageFormat:   &opts.ImageFormat,
		ImageLocation: &opts.ImageLocation,
		ImageType:     &opts.ImageType,
		PublicImage:   &opts.PublicImage,
		ImageHash:     &opts.ImageHash,
		AzName:        opts.AzName,
	}
	result, err := util.DoRequest(action, para)
	if err != nil {
		return nil, err
	}
	snapshotView := new(model.SnapshotView)
	util.FormatView(result.Data, snapshotView)
	return snapshotView, nil

}
func (s *Snapshot) RestoreVolume(opts *options.RestoreSnapshotOptions) error {
	if opts.SnapshotId == "" || opts.TenantId == "" {
		return errors.New("SnapshotId, TenantId  are required")
	}
	action := "RestoreVolume"
	para := &model.RestoreSnapshotParams{
		TenantId:   &opts.TenantId,
		SnapshotId: &opts.SnapshotId,
	}
	_, err := util.DoRequest(action, para)
	return err

}
func (s *Snapshot) DeleteSnapshot(opts *options.DeleteSnapshotOptions) error {
	if opts.TenantId == "" || opts.SnapshotId == "" {
		return errors.New("SnapshotId, TenantId are required")
	}
	action := "DeleteSnapshot"
	para := &model.DeleteSnapshotParams{
		TenantId:   &opts.TenantId,
		SnapshotId: &opts.SnapshotId,
	}
	_, err := util.DoRequest(action, para)
	return err
}

func (s *Snapshot) DescribeSnapshot(opts *options.DescribeSnapshotOptions) (*model.SnapshotView, error) {
	action := "DescribeSnapshot"
	para := &model.DescribeSnapshotParams{
		TenantId:   &opts.TenantId,
		SnapshotId: &opts.SnapshotId,
	}
	result, err := util.DoRequest(action, para)
	if err != nil {
		return nil, err
	}
	snapshotView := new(model.SnapshotView)
	util.FormatView(result.Data, snapshotView)
	return snapshotView, nil
}

//return list
func (s *Snapshot) ListSnapshots(opts *options.ListOptions) (int64, []*model.SnapshotView, error) {
	action := "ListSnapshots"
	var snapshotViewSlice []*model.SnapshotView
	desc, err := ListTenantResources(action, opts)
	if err != nil {
		return -1, nil, err
	}
	if err := util.FormatView(desc.Elements, &snapshotViewSlice); err != nil {
		return -1, nil, err
	}
	return desc.Tc, snapshotViewSlice, nil
}

func (s *Snapshot) ChangeSnapshot(opts *options.ChangeSnapshotOptions) error {
	if opts.TenantId == "" || opts.SnapshotId == "" {
		return errors.New("SnapshotId, TenantId are required")
	}
	action := "ChangeSnapshot"
	para := &model.ChangeSnapshotParams{
		TenantId:     &opts.TenantId,
		SnapshotId:   &opts.SnapshotId,
		SnapshotName: &opts.SnapshotName,
		Description:  &opts.Description,
		Share:        &opts.Share,
	}
	_, err := util.DoRequest(action, para)
	if err != nil {
		return err
	}
	return nil
}

func (s *Snapshot) CopySnapshotList(opts *options.CopySnapshotListOptions) (
	int64, []*model.CopySnapshotView, error) {
	var (
		copySnpViews []*model.CopySnapshotView
		action       = "CopySnapshotList"
		descData     model.DescData
	)

	if opts.Region == "" || opts.TenantId == "" || opts.DstTenantId == "" ||
		len(opts.SnapshotIds) == 0 {
		msg := "Region, TenantId, SnapshotIds, DstTenantid are required"
		return -1, nil, errors.New(msg)
	}

	param := &model.CopySnapshotListParams{
		DstTenantId: opts.DstTenantId,
		SnapshotIds: opts.SnapshotIds,
		TenantId:    opts.TenantId,
		Region:      opts.Region,
	}
	result, err := util.DoRequest(action, param)
	if err != nil {
		return -1, nil, err
	}

	if err = util.FormatView(result.Data, &descData); err != nil {
		return -1, nil, err
	}

	if err = util.FormatView(descData.Elements, &copySnpViews); err != nil {
		return -1, nil, err
	}

	return descData.Tc, copySnpViews, nil
}
