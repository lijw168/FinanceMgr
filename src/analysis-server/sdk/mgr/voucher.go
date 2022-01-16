package mgr

import (
	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
	"encoding/json"
	"errors"
	"fmt"
	//"math"
)

type Voucher struct {
}

func (vr *Voucher) CreateVoucher(opts *options.VoucherOptions) (*model.DescData, error) {
	action := "CreateVoucher"
	switch {
	case opts.InfoOptions.CompanyID <= 0:
		return nil, errors.New("CompanyID is required")
	// case opts.InfoOptions.VoucherMonth <= 0:
	// 	return nil, errors.New("VoucherMonth is required")
	case opts.InfoOptions.VoucherFiller == "":
		return nil, errors.New("VoucherFiller is required")
	}
	params := model.CreateVoucherParams{}
	vouInfoParam := model.VoucherInfoParams{CompanyID: &(opts.InfoOptions.CompanyID),
		VoucherFiller: &(opts.InfoOptions.VoucherFiller)}
	params.InfoParams = &vouInfoParam
	var recordParamSlice []*model.CreateVoucherRecordParams
	for _, val := range opts.RecordsOptions {
		recordItem := val
		recordParam := model.CreateVoucherRecordParams{
			SubjectName: &recordItem.SubjectName,
			Summary:     &recordItem.Summary,
			CreditMoney: &recordItem.CreditMoney,
			DebitMoney:  &recordItem.DebitMoney,
			SubID1:      &recordItem.SubID1,
			// SubID2:      &recordItem.SubID2,
			// SubID3:      &recordItem.SubID3,
			// SubID4:      &recordItem.SubID4,
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

func (vr *Voucher) DeleteVoucher(opts *options.DeleteYearAndIDOptions) error {
	action := "DeleteVoucher"
	switch {
	case opts.ID <= 0:
		return errors.New("ID is required")
	case opts.VoucherYear <= 0:
		return errors.New("VoucherYear is required")
	}
	params := &model.DeleteYearAndIDParams{
		ID:          &opts.ID,
		VoucherYear: &opts.VoucherYear,
	}
	_, err := util.DoRequest(action, params)
	if err != nil {
		return err
	}
	fmt.Printf("DeleteVoucher succeed\n")
	return nil
}

func (vr *Voucher) GetVoucher(opts *options.DescribeYearAndIDOptions) (*model.VoucherView, error) {
	action := "GetVoucher"
	switch {
	case opts.ID <= 0:
		return nil, errors.New("ID is required")
	case opts.VoucherYear <= 0:
		return nil, errors.New("VoucherYear is required")
	}
	params := &model.DescribeYearAndIDParams{
		ID:          &opts.ID,
		VoucherYear: &opts.VoucherYear,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.VoucherView{}
	err = util.FormatView(result.Data, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}

func (vr *Voucher) ArrangeVoucher(opts *options.VoucherArrangeOptions) error {
	action := "ArrangeVoucher"
	switch {
	case opts.CompanyID <= 0:
		return errors.New("CompanyID is required")
	case opts.VoucherMonth <= 0:
		return errors.New("VoucherMonth is required")
	case opts.VoucherYear <= 0:
		return errors.New("VoucherYear is required")
	}
	param := model.VoucherArrangeParams{}
	param.CompanyID = &opts.CompanyID
	param.VoucherYear = &opts.VoucherYear
	param.VoucherMonth = &opts.VoucherMonth
	if opts.ArrangeVoucherNum {
		param.ArrangeVoucherNum = &opts.ArrangeVoucherNum
	}
	_, err := util.DoRequest(action, param)
	return err
}

//该参数直接就是相应的json格式的数据。所以不需要转换了。
// func (vr *Voucher) CreateVoucherRecords_json(params []byte) (*model.DescData, error) {
// 	action := "CreateVoucherRecords"
// 	result, err := util.DoRequest_json(action, params)
// 	if err != nil {
// 		return nil, err
// 	}
// 	desc := &(model.DescData{})
// 	if err := util.FormatView(result.Data, desc); err != nil {
// 		return nil, err
// 	}
// 	return desc, nil
// }

// func (vr *Voucher) CreateVoucherRecords(opts []options.CreateVoucherRecordOptions) (*model.DescData, error) {
// 	action := "CreateVoucherRecords"
// 	switch {
// 	case len(opts) == 0:
// 		return nil, errors.New("voucher records is required")
// 	}
// 	var recordParamSlice []*model.CreateVoucherRecordParams
// 	for _, val := range opts {
// 		recordItem := val
// 		recordParam := model.CreateVoucherRecordParams{
// 			VoucherID:   &recordItem.VoucherID,
// 			SubjectName: &recordItem.SubjectName,
// 			Summary:     &recordItem.Summary,
// 			CreditMoney: &recordItem.CreditMoney,
// 			DebitMoney:  &recordItem.DebitMoney,
// 			SubID1:      &recordItem.SubID1,
// 			// SubID2:      &recordItem.SubID2,
// 			// SubID3:      &recordItem.SubID3,
// 			// SubID4:      &recordItem.SubID4,
// 		}
// 		recordParamSlice = append(recordParamSlice, &recordParam)
// 	}
// 	result, err := util.DoRequest(action, recordParamSlice)
// 	if err != nil {
// 		return nil, err
// 	}
// 	desc := &(model.DescData{})
// 	if err := util.FormatView(result.Data, desc); err != nil {
// 		return nil, err
// 	}
// 	return desc, nil
// }

// func (vr *Voucher) DeleteVoucherRecord(opts *options.BaseOptions) error {
// 	action := "DeleteVoucherRecord"
// 	err := DeleteOpsResource(action, opts)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Printf("DeleteVoucherRecord succeed")
// 	return nil
// }

// //该参数直接就是相应的json格式的数据。所以不需要转换了。
// func (vr *Voucher) DeleteVoucherRecords_json(params []byte) error {
// 	action := "DeleteVoucherRecords"
// 	_, err := util.DoRequest_json(action, params)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (vr *Voucher) GetVoucherInfo(opts *options.DescribeYearAndIDOptions) (*model.VoucherInfoView, error) {
	action := "GetVoucherInfo"
	switch {
	case opts.ID <= 0:
		return nil, errors.New("ID is required")
	case opts.VoucherYear <= 0:
		return nil, errors.New("VoucherYear is required")
	}
	params := &model.DescribeYearAndIDParams{
		ID:          &opts.ID,
		VoucherYear: &opts.VoucherYear,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.VoucherInfoView{}
	err = util.FormatView(result.Data, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}

func (vr *Voucher) GetLatestVoucherInfo(opts *options.DescribeYearAndIDOptions) (int64, []*model.VoucherInfoView, error) {
	action := "GetLatestVoucherInfo"
	switch {
	case opts.ID <= 0:
		return -1, nil, errors.New("ID is required")
	case opts.VoucherYear <= 0:
		return -1, nil, errors.New("VoucherYear is required")
	}
	params := &model.DescribeYearAndIDParams{
		ID:          &opts.ID,
		VoucherYear: &opts.VoucherYear,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return -1, nil, err
	}
	desc := &(model.DescData{})
	if err := util.FormatView(result.Data, desc); err != nil {
		return -1, nil, err
	}

	var ret []*model.VoucherInfoView
	if err := util.FormatView(desc.Elements, &ret); err != nil {
		return -1, nil, err
	}
	return desc.Tc, ret, nil
}

func (vr *Voucher) GetMaxNumOfMonth(opts *options.QueryMaxNumOfMonthOption) (int64, error) {
	action := "GetMaxNumOfMonth"
	if opts.CompanyID <= 0 {
		return 0, errors.New("CompanyID is required")
	}
	if opts.VoucherMonth <= 0 {
		return -1, errors.New("VoucherMonth is required")
	}
	if opts.VoucherYear <= 0 {
		return -1, errors.New("VoucherYear is required")
	}
	params := &model.QueryMaxNumOfMonthParams{
		CompanyID:    &opts.CompanyID,
		VoucherYear:  &opts.VoucherYear,
		VoucherMonth: &opts.VoucherMonth,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return -1, err
	}
	var iCount int64
	if err := util.FormatView(result.Data, &iCount); err != nil {
		return 0, err
	}
	return iCount, nil
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

// func (vr *Voucher) UpdateVoucherRecordByID(opts *options.ModifyVoucherRecordOptions) error {
// 	action := "UpdateVoucherRecordByID"
// 	switch {
// 	case opts.VouRecordID <= 0:
// 		return errors.New("VouRecordID is required")
// 	}
// 	param := model.ModifyVoucherRecordParams{}
// 	if opts.VouRecordID != 0 {
// 		param.VouRecordID = &opts.VouRecordID
// 	}
// 	if opts.Summary != "" {
// 		param.Summary = &opts.Summary
// 	}
// 	if opts.SubjectName != "" {
// 		param.SubjectName = &opts.SubjectName
// 	}
// 	if math.Abs(opts.CreditMoney) >= 0.001 {
// 		param.CreditMoney = &opts.CreditMoney
// 	}
// 	if math.Abs(opts.DebitMoney) >= 0.001 {
// 		param.DebitMoney = &opts.DebitMoney
// 	}
// 	if opts.SubID1 >= 0 {
// 		param.SubID1 = &opts.SubID1
// 	}
// 	// if opts.SubID2 != 0 {
// 	// 	param.SubID2 = &opts.SubID2
// 	// }
// 	// if opts.SubID3 != 0 {
// 	// 	param.SubID3 = &opts.SubID3
// 	// }
// 	// if opts.SubID4 != 0 {
// 	// 	param.SubID4 = &opts.SubID4
// 	// }
// 	_, err := util.DoRequest(action, param)
// 	return err
// }

// func (vr *Voucher) UpdateVoucherRecord_json(params []byte) error {
// 	action := "UpdateVoucherRecordByID"
// 	_, err := util.DoRequest_json(action, params)
// 	return err
// }

func (vr *Voucher) UpdateVoucherInfo(opts *options.ModifyVoucherInfoOptions) error {
	action := "UpdateVoucherInfo"
	switch {
	case opts.VoucherID <= 0:
		return errors.New("VouRecordID is required")
	case opts.VoucherYear <= 0:
		return errors.New("VoucherYear is required")
	}
	param := model.ModifyVoucherInfoParams{}
	param.VoucherID = &opts.VoucherID
	param.VoucherYear = &opts.VoucherYear
	if opts.VoucherDate > 0 {
		param.VoucherDate = &opts.VoucherDate
	}
	if opts.VoucherFiller != "" {
		param.VoucherFiller = &opts.VoucherFiller
	}
	if opts.VoucherAuditor != "" {
		param.VoucherAuditor = &opts.VoucherAuditor
	}
	if opts.BillCount > 0 {
		param.BillCount = &opts.BillCount
	}
	if opts.Status > 0 {
		param.Status = &opts.Status
	}
	_, err := util.DoRequest(action, param)
	return err
}

//该参数直接就是相应的json格式的数据。所以不需要转换了。
func (vr *Voucher) CreateVoucher_json(params []byte) ([]byte, error) {
	action := "CreateVoucher"
	result, err := util.DoRequest_json(action, params)
	if err != nil {
		return nil, err
	}
	return json.Marshal(result.Data)
}

//该参数直接就是相应的json格式的数据。所以不需要转换了。
func (vr *Voucher) UpdateVoucher_json(params []byte) ([]byte, error) {
	action := "UpdateVoucher"
	result, err := util.DoRequest_json(action, params)
	if err != nil {
		return nil, err
	}
	return json.Marshal(result.Data)
}

func (vr *Voucher) DeleteVoucher_json(params []byte) error {
	action := "DeleteVoucher"
	_, err := util.DoRequest_json(action, params)
	if err != nil {
		return err
	}
	fmt.Printf("DeleteVoucher succeed\n")
	return nil
}

func (vr *Voucher) GetVoucher_json(params []byte) ([]byte, error) {
	action := "GetVoucher"
	result, err := util.DoRequest_json(action, params)
	if err != nil {
		return nil, err
	}
	return json.Marshal(result.Data)
}

func (vr *Voucher) ArrangeVoucher_json(params []byte) error {
	action := "ArrangeVoucher"
	_, err := util.DoRequest_json(action, params)
	return err
}

func (vr *Voucher) GetVoucherInfo_json(params []byte) ([]byte, error) {
	action := "GetVoucherInfo"
	result, err := util.DoRequest_json(action, params)
	if err != nil {
		return nil, err
	}
	return json.Marshal(result.Data)
}

func (vr *Voucher) GetLatestVoucherInfo_json(params []byte) ([]byte, error) {
	action := "GetLatestVoucherInfo"
	result, err := util.DoRequest_json(action, params)
	if err != nil {
		return nil, err
	}
	return json.Marshal(result.Data)
}

func (vr *Voucher) GetMaxNumOfMonth_json(params []byte) (int64, error) {
	action := "GetMaxNumOfMonth"
	result, err := util.DoRequest_json(action, params)
	if err != nil {
		return 0, err
	}
	var iCount int64
	if err := util.FormatView(result.Data, &iCount); err != nil {
		return 0, err
	}
	return iCount, nil
}

func (vr *Voucher) ListVoucherInfo_json(params []byte) ([]byte, error) {
	action := "ListVoucherInfo"
	return ListOpsResources_json(action, params)
}

func (vr *Voucher) ListVoucherInfoByMulCondition_json(params []byte) ([]byte, error) {
	action := "ListVoucherInfoByMulCondition"
	return ListOpsResources_json(action, params)
}

func (vr *Voucher) ListVoucherRecords_json(params []byte) ([]byte, error) {
	action := "ListVoucherRecords"
	return ListOpsResources_json(action, params)
}

func (vr *Voucher) UpdateVoucherInfo_json(params []byte) error {
	action := "UpdateVoucherInfo"
	_, err := util.DoRequest_json(action, params)
	return err
}

func (vr *Voucher) BatchAuditVouchers_json(params []byte) error {
	action := "BatchAuditVouchers"
	_, err := util.DoRequest_json(action, params)
	return err
}

//voucher template begin
func (vr *Voucher) CreateVoucherTemplate_json(params []byte) (int, error) {
	action := "CreateVoucherTemplate"
	result, err := util.DoRequest_json(action, params)
	if err != nil {
		return 0, err
	}
	var iSerialNum int
	if err := util.FormatView(result.Data, &iSerialNum); err != nil {
		return 0, err
	}
	return iSerialNum, nil
}

//该参数直接就是相应的json格式的数据。所以不需要转换了。
func (vr *Voucher) ListVoucherTemplate_json(params []byte) ([]byte, error) {
	action := "ListVoucherTemplate"
	return ListOpsResources_json(action, params)
}

func (vr *Voucher) DeleteVoucherTemplate(opts *options.BaseOptions) error {
	action := "DeleteVoucherTemplate"
	err := DeleteOpsResource(action, opts)
	if err != nil {
		return err
	}
	fmt.Printf("DeleteVoucher succeed\n")
	return nil
}

func (vr *Voucher) GetVoucherTemplate(opts *options.BaseOptions) ([]byte, error) {
	action := "GetVoucherTemplate"
	resData, err := DescribeOpsResource(action, opts)
	if err != nil {
		return nil, err
	}
	return json.Marshal(resData)
}

//voucher template end

//voucher generate report forms,begin
func (vr *Voucher) CalculateAccumulativeMoney_json(params []byte) ([]byte, error) {
	action := "CalculateAccumulativeMoney"
	resData, err := util.DoRequest_json(action, params)
	if err != nil {
		return nil, err
	}
	return json.Marshal(resData)
}

//voucher generate report forms,end
