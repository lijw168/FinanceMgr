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
	case opts.Summary == "":
		return errors.New("Summary is required")
	case opts.SubjectDirection <= 0:
		return errors.New("SubjectDirection is required")
	}
	params := model.YearBalanceParams{
		SubjectID:        &opts.SubjectID,
		Summary:          &opts.Summary,
		SubjectDirection: &opts.SubjectDirection,
		Balance:          &opts.Balance}
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

func (yb *YearBalance) DeleteYearBalanceByID(opts *options.BaseOptions) error {
	action := "DeleteYearBalanceByID"
	err := DeleteOpsResource(action, opts)
	if err != nil {
		return err
	}
	fmt.Printf("DeleteYearBalanceByID succeed")
	return nil
}

func (yb *YearBalance) DeleteYearBalanceByID_json(params []byte) error {
	action := "DeleteYearBalanceByID"
	_, err := util.DoRequest_json(action, params)
	if err != nil {
		return err
	}
	fmt.Printf("DeleteYearBalanceByID succeed")
	return nil
}

func (yb *YearBalance) GetYearBalanceById(opts *options.BaseOptions) (*model.YearBalanceView, error) {
	action := "GetYearBalanceById"
	dr, err := DescribeOpsResource(action, opts)
	if err != nil {
		return nil, err
	}
	view := &model.YearBalanceView{}
	err = util.FormatView(dr, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}

func (yb *YearBalance) GetYearBalanceById_json(params []byte) ([]byte, error) {
	action := "GetYearBalanceById"
	result, err := util.DoRequest_json(action, params)
	if err != nil {
		return nil, err
	}
	return json.Marshal(result.Data)
}

// func (yb *YearBalance) UpdateYearBalanceById(opts *options.YearBalanceOption) error {
// 	action := "UpdateYearBalanceById"
// 	if opts.SubjectID <= 0 {
// 		return errors.New("SubjectID is required")
// 	}
// 	param := model.YearBalanceParams{}
// 	if opts.Summary != "" {
// 		param.Summary = &opts.Summary
// 	}
// 	if opts.SubjectDirection != 0 {
// 		param.SubjectDirection = &opts.SubjectDirection
// 	}
// 	// if math.Abs(opts.Balance) >= 0.001 {
// 	// 	param.Balance = &opts.Balance
// 	// }
// 	_, err := util.DoRequest(action, param)
// 	return err
// }

func (yb *YearBalance) UpdateYearBalanceById_json(param []byte) error {
	action := "UpdateYearBalanceById"
	_, err := util.DoRequest_json(action, param)
	return err
}
