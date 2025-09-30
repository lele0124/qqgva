package request

import (
	"time"
)



// UpdateMerchantRequest 更新商户的请求模型
// 该模型使用string类型的时间字段和类型字段，以避免类型转换错误

type UpdateMerchantRequest struct {
	ID                uint        `json:"id" form:"id" binding:"required"`
	MerchantName      string      `json:"merchantName" form:"merchantName" binding:"required"`
	MerchantIcon      string      `json:"merchantIcon" form:"merchantIcon"`
	ParentID          *uint       `json:"parentID" form:"parentID"`
	MerchantType      uint        `json:"merchantType" form:"merchantType" binding:"required"` // 改为uint类型，与数据模型保持一致
	BusinessLicense   string      `json:"businessLicense" form:"businessLicense"`
	LegalPerson       string      `json:"legalPerson" form:"legalPerson"`
	RegisteredAddress string      `json:"registeredAddress" form:"registeredAddress"`
	BusinessScope     string      `json:"businessScope" form:"businessScope"`
	IsEnabled         bool        `json:"isEnabled" form:"isEnabled"` // 改为bool类型，与数据模型保持一致
	ValidStartTime    string      `json:"validStartTime" form:"validStartTime"` // 使用string类型接收时间
	ValidEndTime      string      `json:"validEndTime" form:"validEndTime"`     // 使用string类型接收时间
	MerchantLevel     uint        `json:"merchantLevel" form:"merchantLevel" binding:"required"` // 改为uint类型，与数据模型保持一致
}

// ToMerchantModel 将请求模型转换为数据模型
func (req *UpdateMerchantRequest) ToMerchantModel() (interface{}, error) {
	// 创建返回的map
	result := map[string]interface{}{
		"ID":                req.ID,
		"MerchantName":      req.MerchantName,
		"MerchantIcon":      req.MerchantIcon,
		"ParentID":          req.ParentID,
		"MerchantType":      req.MerchantType,
		"BusinessLicense":   req.BusinessLicense,
		"LegalPerson":       req.LegalPerson,
		"RegisteredAddress": req.RegisteredAddress,
		"BusinessScope":     req.BusinessScope,
		"IsEnabled":         req.IsEnabled,
		"MerchantLevel":     req.MerchantLevel,
	}

	// 处理时间字段，只有非空时才尝试解析
	if req.ValidStartTime != "" {
		startTime, timeErr := time.Parse(time.RFC3339, req.ValidStartTime)
		if timeErr != nil {
			// 尝试其他常见的时间格式
			startTime, timeErr = time.Parse("2006-01-02 15:04:05", req.ValidStartTime)
			if timeErr != nil {
				return nil, timeErr
			}
		}
		result["ValidStartTime"] = &startTime
	}

	if req.ValidEndTime != "" {
		endTime, timeErr := time.Parse(time.RFC3339, req.ValidEndTime)
		if timeErr != nil {
			// 尝试其他常见的时间格式
			endTime, timeErr = time.Parse("2006-01-02 15:04:05", req.ValidEndTime)
			if timeErr != nil {
				return nil, timeErr
			}
		}
		result["ValidEndTime"] = &endTime
	}

	return result, nil
}