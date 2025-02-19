package mgr

import (
	"encoding/json"
	"errors"
	"financeMgr/src/analysis-server/model"
	"financeMgr/src/analysis-server/sdk/options"
	"financeMgr/src/analysis-server/sdk/util"
	"fmt"
)

type YearBalance struct {
}

func (yb *YearBalance) CreateYearBalance(opts *options.YearBalanceOption) error {
	action := "CreateYearBalance"
	switch {
	case opts.CompanyID <= 0:
		return errors.New("CompanyID is required")
	case opts.SubjectID <= 0:
		return errors.New("SubjectID is required")
	case opts.Year <= 0:
		return errors.New("year is required")
	case opts.Balance < 0:
		return errors.New("balance is required")
	}
	params := model.OptYearBalanceParams{
		CompanyID: &opts.CompanyID,
		SubjectID: &opts.SubjectID,
		Year:      &opts.Year,
		Balance:   &opts.Balance}
	//Status:    &opts.Status}
	_, err := util.DoRequest(action, params)
	if err != nil {
		return err
	}
	fmt.Printf("CreateYearBalance succeed")
	return nil
}

func (yb *YearBalance) CreateYearBalance_json(params []byte) error {
	action := "CreateYearBalance"
	_, err := util.DoRequest_json(action, params)
	return err
}

func (yb *YearBalance) BatchCreateYearBalance_json(params []byte) error {
	action := "BatchCreateYearBalance"
	_, err := util.DoRequest_json(action, params)
	return err
}

func (yb *YearBalance) BatchDeleteYearBalance_json(params []byte) error {
	action := "BatchDeleteYearBalance"
	_, err := util.DoRequest_json(action, params)
	return err
}

func (yb *YearBalance) DeleteYearBalance(opts *options.BasicYearBalance) error {
	action := "DeleteYearBalance"
	switch {
	case opts.CompanyID <= 0:
		return errors.New("CompanyID is required")
	case opts.SubjectID <= 0:
		return errors.New("SubjectID is required")
	// case opts.Balance == "":
	// 	return errors.New("Summary is required")
	case opts.Year <= 0:
		return errors.New("year is required")
	}
	params := model.BasicYearBalanceParams{
		CompanyID: &opts.CompanyID,
		SubjectID: &opts.SubjectID,
		Year:      &opts.Year}
	_, err := util.DoRequest(action, params)
	if err != nil {
		return err
	}
	fmt.Printf("DeleteYearBalanceByID succeed")
	return nil
}

func (yb *YearBalance) DeleteYearBalance_json(params []byte) error {
	action := "DeleteYearBalance"
	_, err := util.DoRequest_json(action, params)
	if err != nil {
		return err
	}
	fmt.Printf("DeleteYearBalanceByID succeed")
	return nil
}

func (yb *YearBalance) GetAccSubYearBalValue(opts *options.BasicYearBalance) (float64, error) {
	action := "GetAccSubYearBalValue"
	switch {
	case opts.CompanyID <= 0:
		return 0, errors.New("CompanyID is required")
	case opts.SubjectID <= 0:
		return 0, errors.New("SubjectID is required")
	case opts.Year <= 0:
		return 0, errors.New("year is required")
	}
	params := model.BasicYearBalanceParams{
		CompanyID: &opts.CompanyID,
		SubjectID: &opts.SubjectID,
		Year:      &opts.Year}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return 0, err
	}
	//var yearBal float64 = result.Data.(float64)
	var yearBal float64
	err = util.FormatView(result.Data, &yearBal)
	if err != nil {
		return 0, err
	}
	return yearBal, nil
}

func (yb *YearBalance) GetAccSubYearBalValue_json(params []byte) ([]byte, error) {
	action := "GetAccSubYearBalValue"
	result, err := util.DoRequest_json(action, params)
	if err != nil {
		return nil, err
	}
	return json.Marshal(result.Data)
}

func (yb *YearBalance) GetYearBalance(opts *options.BasicYearBalance) (*model.YearBalanceView, error) {
	action := "GetYearBalance"
	switch {
	case opts.CompanyID <= 0:
		return nil, errors.New("CompanyID is required")
	case opts.SubjectID <= 0:
		return nil, errors.New("SubjectID is required")
	case opts.Year <= 0:
		return nil, errors.New("year is required")
	}
	params := model.BasicYearBalanceParams{
		CompanyID: &opts.CompanyID,
		SubjectID: &opts.SubjectID,
		Year:      &opts.Year}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	yearBalView := model.YearBalanceView{}
	err = util.FormatView(result.Data, &yearBalView)
	if err != nil {
		return nil, err
	}
	return &yearBalView, nil
}

func (yb *YearBalance) GetYearBalance_json(params []byte) ([]byte, error) {
	action := "GetYearBalance"
	result, err := util.DoRequest_json(action, params)
	if err != nil {
		return nil, err
	}
	return json.Marshal(result.Data)
}

func (yb *YearBalance) UpdateYearBalance(opts *options.YearBalanceOption) error {
	action := "UpdateYearBalance"
	switch {
	case opts.CompanyID <= 0:
		return errors.New("companyId is required")
	case opts.Year <= 0:
		return errors.New("year is required")
	case opts.SubjectID <= 0:
		return errors.New("SubjectID is required")
	}
	params := model.OptYearBalanceParams{
		SubjectID: &opts.SubjectID,
		CompanyID: &opts.CompanyID,
		Year:      &opts.Year,
		Balance:   &opts.Balance,
		Status:    &opts.Status}
	_, err := util.DoRequest(action, params)
	if err != nil {
		return err
	}
	return nil
}

func (yb *YearBalance) UpdateYearBalance_json(param []byte) error {
	action := "UpdateYearBalance"
	_, err := util.DoRequest_json(action, param)
	return err
}

// 就不增加该接口的cli命令了。
func (yb *YearBalance) BatchUpdateBals_json(param []byte) error {
	action := "BatchUpdateBals"
	_, err := util.DoRequest_json(action, param)
	return err
}

func (yb *YearBalance) ListYearBalance_json(params []byte) ([]byte, error) {
	action := "ListYearBalance"
	return ListOpsResources_json(action, params)
}

func (yb *YearBalance) ListYearBalance(opts *options.ListOptions) (int64, []*model.YearBalanceView, error) {
	action := "ListYearBalance"
	var ret []*model.YearBalanceView
	desc, err := ListOpsResources(action, opts)
	if err != nil {
		return -1, nil, err
	}
	if err := util.FormatView(desc.Elements, &ret); err != nil {
		return -1, nil, err
	}
	return desc.Tc, ret, nil
}

// 暂不为年度结算和取消年度结算添加cli命令，因为暂时没有用到。
func (yb *YearBalance) AnnualClosing_json(params []byte) error {
	action := "AnnualClosing"
	_, err := util.DoRequest_json(action, params)
	return err
}

func (yb *YearBalance) CancelAnnualClosing_json(params []byte) error {
	action := "CancelAnnualClosing"
	_, err := util.DoRequest_json(action, params)
	return err
}

func (yb *YearBalance) GetAnnualClosingStatus_json(params []byte) ([]byte, error) {
	action := "GetAnnualClosingStatus"
	result, err := util.DoRequest_json(action, params)
	if err != nil {
		return nil, err
	}
	return json.Marshal(result.Data)
}
