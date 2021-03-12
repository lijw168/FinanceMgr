package utils

import "analysis-server/model"

func FilterKeyRename(filters []*model.FilterItem, keyMap map[string]string) {
	for _, item := range filters {
		if value, ok := keyMap[*item.Field]; ok {
			item.Field = &value
		}
	}
}
