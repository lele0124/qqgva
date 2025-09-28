package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
	SortField string `json:"sortField" form:"sortField"` // 排序字段
	SortOrder string `json:"sortOrder" form:"sortOrder"` // 排序顺序: asc 或 desc
}
