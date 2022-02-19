package mgr

import (
	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
	"encoding/json"
	"errors"
	"fmt"
)

type YearBalance struct {
}

func (yb *YearBalance) CreateYearBalance(opts *options.YearBalanceOption) error {
	action := "CreateYearBalance"
	switch {
	case opts.SubjectID <= 0:
		return errors.New("SubjectID is required")
	// case opts.Balance == "":
	// 	return errors.New("Summary is required")
	case opts.Year <= 0:
		return errors.New("Year is required")
	}
	params := model.OptYearBalanceParams{
		SubjectID: &opts.SubjectID,
		Year:      &opts.Year,
		Balance:   &opts.Balance}
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

func (yb *YearBalance) DeleteYearBalance(opts *options.BasicYearBalance) error {
	action := "DeleteYearBalance"
	switch {
	case opts.SubjectID <= 0:
		return errors.New("SubjectID is required")
	// case opts.Balance == "":
	// 	return errors.New("Summary is required")
	case opts.Year <= 0:
		return errors.New("Year is required")
	}
	params := model.BasicYearBalanceParams{
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

func (yb *YearBalance) GetYearBalance(opts *options.BasicYearBalance) (float64, error) {
	action := "GetYearBalance"
	switch {
	case opts.SubjectID <= 0:
		return 0, errors.New("SubjectID is required")
	// case opts.Balance == "":
	// 	return errors.New("Summary is required")
	case opts.Year <= 0:
		return 0, errors.New("Year is required")
	}
	params := model.BasicYearBalanceParams{
		SubjectID: &opts.SubjectID,
		Year:      &opts.Year}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return 0, err
	}
	var yearBal float64 = result.Data.(float64)
	// err = util.FormatView(result.Data, &yearBal)
	// if err != nil {
	// 	return nil, err
	// }
	return yearBal, nil
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
	case opts.SubjectID <= 0:
		return errors.New("SubjectID is required")
	// case opts.Balance == "":
	// 	return errors.New("Summary is required")
	case opts.Year <= 0:
		return errors.New("Year is required")
	}
	params := model.OptYearBalanceParams{
		SubjectID: &opts.SubjectID,
		Year:      &opts.Year,
		Balance:   &opts.Balance}
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
