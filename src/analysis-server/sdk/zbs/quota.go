package zbs

import (
	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
)

type Quota struct {
}

//return list
func (q *Quota) ListQuotas(opts *options.ListOptions) (int64, []*model.QuotaView, error) {
	action := "ListQuota"
	var quotaViewSlice []*model.QuotaView
	desc, err := ListTenantResources(action, opts)
	if err != nil {
		return -1, nil, err
	}
	if err := util.FormatView(desc.Elements, &quotaViewSlice); err != nil {
		return -1, nil, err
	}
	return desc.Tc, quotaViewSlice, nil
}
