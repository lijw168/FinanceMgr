package zbs

import (
	"errors"
	"common/constant"
	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
	"regexp"
	"strconv"
)

type Host struct {
}

type IsSuccess struct {
	Success *int
}

func (h *Host) CreateHost(opts *options.CreateHostOptions) (*model.HostView, error) {
	switch {
	case opts.Name == "":
		fallthrough
	case opts.Rack == "":
		return nil, errors.New("Name, IpAddress, Rack are required")
	}
	action := "CreateHost"
	para := &model.CreateHostParams{
		Name:        &opts.Name,
		MgmtIp:      &opts.MgmtIp,
		DataIp:      &opts.DataIp,
		ClientIp:    &opts.ClientIp,
		TraceIp:     &opts.MgmtIp,
		Rack:        &opts.Rack,
		PoolId:      &opts.PoolId,
		AdminStatus: &opts.AdminStatus,
	}
	result, err := util.DoRequest(action, para)
	if err != nil {
		return nil, err
	}
	var hv *model.HostView = new(model.HostView)
	util.FormatView(result.Data, hv)
	return hv, nil

}

//return list
func (h *Host) DescribeHosts(opts *options.ListOptions) (int64, []*model.HostView, error) {
	action := "DescribeHosts"
	var ret []*model.HostView
	desc, err := ListOpsResources(action, opts)
	if err != nil {
		return -1, nil, err
	}
	if err := util.FormatView(desc.Elements, &ret); err != nil {
		return -1, nil, err
	}
	return desc.Tc, ret, nil
}

//return true of false
func (h *Host) ModifyAdminStaus(opts *options.ModifyAdminStatusOptions) (*model.HostView, error) {
	action := "ModifyAdminStaus"
	switch {
	case opts.Name == "":
		fallthrough
	case opts.AdminStatus == "":
		return nil, errors.New("Name, AdminStatus are required")
	}
	switch opts.AdminStatus {
	case "DOWN":
		opts.AdminStatus = "0"
	case "UP":
		opts.AdminStatus = "1"
	case "GRAY":
		opts.AdminStatus = "2"
	default:
		return nil, errors.New("host AdminiStatus Value Error.")
	}
	//may be redundance
	ok, err := CheckAdminstatus(opts.AdminStatus)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("host adminStatus value error.")
	}

	asValue, _ := strconv.Atoi(opts.AdminStatus)
	param := &model.ModifyHostAdminStatusParam{
		Name:        &opts.Name,
		AdminStatus: &asValue,
	}
	result, err := util.DoRequest(action, param)
	if err != nil {
		return nil, err
	}
	hv := &model.HostView{}
	util.FormatView(result.Data, hv)
	return hv, nil
}

func (h *Host) ModifyStatus(opts *options.ModifyStatusOptions) error {
	action := "ModifyHostStatus"
	switch {
	case opts.Name == "":
		return errors.New("Name are required")
	}
	switch opts.Status {
	case constant.HOST_STATUS_UP:
	case constant.HOST_STATUS_DOWN:
	default:
		return errors.New("Host Status Value Error.")
	}
	param := &model.ModifyHostStatusParams{
		Name:   &opts.Name,
		Status: &opts.Status,
		Force:  &opts.Force,
	}
	_, err := util.DoRequest(action, param)
	return err
}

//return hostview
func (h *Host) DescribeHost(opts *options.DescribeHostOptions) (*model.HostView, error) {
	action := "DescribeHost"
	switch {
	case opts.Id == "":
		return nil, errors.New("host Id value required.")
	}
	param := &model.DescribeHostParams{
		Id: &opts.Id,
	}
	result, err := util.DoRequest(action, param)
	if err != nil {
		return nil, err
	}
	var hostView *model.HostView = new(model.HostView)
	util.FormatView(result.Data, &hostView)
	return hostView, nil
}

func (h *Host) DeleteHost(opts *options.DeleteHostOptions) (*model.AtDeleteView, error) {
	switch {
	case opts.Id == "":
		return nil, errors.New("host Id value required.")
	}
	action := "DeleteHost"
	// use DescribeHostParams replace DeleteHostParams
	para := &model.DeleteHostParams{
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

func CheckType(t string) (bool, error) {
	tv, err := strconv.Atoi(t)
	if err != nil {
		return false, errors.New("Host Type must be a int type(0[vs_host],1[vr_host],2[vpngw],3[snatgw] available.)")
	}
	if tv != 0 && tv != 1 && tv != 2 && tv != 3 {
		return false, errors.New("Host Type's Value must be 0[vs_host] or 1[vr_host] or 2[vpngw] or 3[snatgw]")
	}
	return true, nil
}

func CheckAdminstatus(adminStatus string) (bool, error) {
	as, err := strconv.Atoi(adminStatus)
	if err != nil {
		return false, errors.New("Host adminStatus must be a int type(0[host_status_down],1[host_status_up] available.)")
	}
	if as != constant.HOST_ADMINSTATUS_UP && as != constant.HOST_ADMINSTATUS_DOWN && as != constant.HOST_ADMINSTATUS_GRAY {
		return false, errors.New("Host adminStatus's Value must be 0[host_status_down],1[host_status_up],2[host_status_test].")
	}
	return true, nil
}

func CheckIpAddr(ipaddr string) (bool, error) {
	ipaddr_regexp_pattern := "^(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9])\\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[0-9])$"
	is_match, err := regexp.Match(ipaddr_regexp_pattern, []byte(ipaddr))
	if err != nil {
		return false, err
	}
	return is_match, nil
}
