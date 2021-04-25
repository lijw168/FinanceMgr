package mgr

import (
	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
	"errors"
	"fmt"
)

type Voucher struct {
}

func (vr *Voucher) CreateVoucher(opts *options.VoucherOptions) (*model.DescData, error) {
	action := "CreateVoucher"
	switch {
	case opts.InfoOptions.CompanyID <= 0:
		return nil, errors.New("CompanyID is required")
	case opts.InfoOptions.VoucherMonth <= 0:
		return nil, errors.New("VoucherMonth is required")
	case len(opts.RecordsOptions) == 0:
		return nil, errors.New("VoucherRecords is required")
	}
	params := model.VoucherParams{}
	vouInfoParam := model.VoucherInfoParams{CompanyID: &(opts.InfoOptions.CompanyID), VoucherMonth: &(opts.InfoOptions.VoucherMonth)}
	params.InfoParams = &vouInfoParam
	var recordParamSlice []*model.CreateVoucherRecordParams
	for _, val := range opts.RecordsOptions {
		recordItem := val
		recordParam := model.CreateVoucherRecordParams{
			SubjectName: &recordItem.SubjectName,
			Summary:     &recordItem.Summary,
			BillCount:   &recordItem.BillCount,
			CreditMoney: &recordItem.CreditMoney,
			DebitMoney:  &recordItem.DebitMoney,
			SubID1:      &recordItem.SubID1,
			SubID2:      &recordItem.SubID2,
			SubID3:      &recordItem.SubID3,
			SubID4:      &recordItem.SubID4,
		}
		recordParamSlice = append(recordParamSlice, &recordParam)
	}
	params.RecordsParams = append(params.RecordsParams, recordParamSlice...)
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

func (vr *Voucher) DeleteVoucher(opts *options.BaseOptions) error {
	action := "DeleteVoucher"
	err := DeleteOpsResource(action, opts)
	if err != nil {
		return err
	}
	//fmt.Printf("DeleteVoucher succeed\n")
	return nil
}

func (vr *Voucher) GetVoucher(opts *options.BaseOptions) (*model.VoucherView, error) {
	action := "GetVoucher"
	dr, err := DescribeOpsResource(action, opts)
	if err != nil {
		return nil, err
	}
	view := &model.VoucherView{}
	err = util.FormatView(dr, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}

func (vr *Voucher) CreateVoucherRecords(opts []options.CreateVoucherRecordOptions) (*model.DescData, error) {
	action := "CreateVoucherRecords"
	switch {
	case len(opts) == 0:
		return nil, errors.New("voucher records is required")
	}
	var recordParamSlice []*model.CreateVoucherRecordParams
	for _, val := range opts {
		recordItem := val
		recordParam := model.CreateVoucherRecordParams{
			VoucherID:   &recordItem.VoucherID,
			SubjectName: &recordItem.SubjectName,
			Summary:     &recordItem.Summary,
			BillCount:   &recordItem.BillCount,
			CreditMoney: &recordItem.CreditMoney,
			DebitMoney:  &recordItem.DebitMoney,
			SubID1:      &recordItem.SubID1,
			SubID2:      &recordItem.SubID2,
			SubID3:      &recordItem.SubID3,
			SubID4:      &recordItem.SubID4,
		}
		recordParamSlice = append(recordParamSlice, &recordParam)
	}
	result, err := util.DoRequest(action, recordParamSlice)
	if err != nil {
		return nil, err
	}
	desc := &(model.DescData{})
	if err := util.FormatView(result.Data, desc); err != nil {
		return nil, err
	}
	return desc, nil
}

func (vr *Voucher) DeleteVoucherRecord(opts *options.BaseOptions) error {
	action := "DeleteVoucherRecord"
	err := DeleteOpsResource(action, opts)
	if err != nil {
		return err
	}
	fmt.Printf("DeleteVoucherRecord succeed")
	return nil
}

func (vr *Voucher) GetVoucherInfo(opts *options.BaseOptions) (*model.VoucherInfoView, error) {
	action := "GetVoucherInfo"
	dr, err := DescribeOpsResource(action, opts)
	if err != nil {
		return nil, err
	}
	view := &model.VoucherInfoView{}
	err = util.FormatView(dr, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}

func (vr *Voucher) ListVoucherInfo(opts *options.ListOptions) (int64, []*model.VoucherInfoView, error) {
	action := "ListVoucherInfo"
	var ret []*model.VoucherInfoView
	desc, err := ListOpsResources(action, opts)
	if err != nil {
		return -1, nil, err
	}
	if err := util.FormatView(desc.Elements, &ret); err != nil {
		return -1, nil, err
	}
	return desc.Tc, ret, nil
}

func (vr *Voucher) ListVoucherRecords(opts *options.ListOptions) (int64, []*model.VoucherRecordView, error) {
	action := "ListVoucherRecords"
	var ret []*model.VoucherRecordView
	desc, err := ListOpsResources(action, opts)
	if err != nil {
		return -1, nil, err
	}
	if err := util.FormatView(desc.Elements, &ret); err != nil {
		return -1, nil, err
	}
	return desc.Tc, ret, nil
}

func (vr *Voucher) UpdateVoucherRecord(opts *options.ModifyVoucherRecordOptions) error {
	action := "UpdateVoucherRecord"
	switch {
	case opts.VouRecordID <= 0:
		return errors.New("VouRecordID is required")
	}
	param := model.ModifyVoucherRecordParams{}
	if opts.VouRecordID != 0 {
		param.VouRecordID = &opts.VouRecordID
	}
	if opts.Summary != "" {
		param.Summary = &opts.Summary
	}
	if opts.SubjectName != "" {
		param.SubjectName = &opts.SubjectName
	}
	if opts.BillCount != -1 {
		param.BillCount = &opts.BillCount
	}
	if opts.CreditMoney != -1 {
		param.CreditMoney = &opts.CreditMoney
	}
	if opts.DebitMoney != -1 {
		param.DebitMoney = &opts.DebitMoney
	}
	if opts.SubID1 != 0 {
		param.SubID1 = &opts.SubID1
	}
	if opts.SubID2 != 0 {
		param.SubID2 = &opts.SubID2
	}
	if opts.SubID3 != 0 {
		param.SubID3 = &opts.SubID3
	}
	if opts.SubID4 != 0 {
		param.SubID4 = &opts.SubID4
	}
	_, err := util.DoRequest(action, param)
	return err
}
