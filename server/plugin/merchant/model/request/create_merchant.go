package request

import (
	"strconv"
	"time"
)

// CreateMerchantRequest 创建商户的请求模型
// 该模型使用string类型的时间字段和类型字段，以避免类型转换错误

type CreateMerchantRequest struct {
	MerchantName      string  `json:"merchantName" form:"merchantName" binding:"required"`
	MerchantIcon      string  `json:"merchantIcon" form:"merchantIcon"`
	ParentID          *uint   `json:"parentID" form:"parentID"`
	MerchantType      string  `json:"merchantType" form:"merchantType" binding:"required"` // 使用string类型接收
	BusinessLicense   string  `json:"businessLicense" form:"businessLicense"`
	LegalPerson       string  `json:"legalPerson" form:"legalPerson"`
	RegisteredAddress string  `json:"registeredAddress" form:"registeredAddress"`
	BusinessScope     string  `json:"businessScope" form:"businessScope"`
	IsEnabled         string  `json:"isEnabled" form:"isEnabled"` // 使用string类型接收
	ValidStartTime    string  `json:"validStartTime" form:"validStartTime"` // 使用string类型接收时间
	ValidEndTime      string  `json:"validEndTime" form:"validEndTime"`     // 使用string类型接收时间
	MerchantLevel     string  `json:"merchantLevel" form:"merchantLevel" binding:"required"` // 使用string类型接收
}

// ToMerchantModel 将请求模型转换为数据模型
func (req *CreateMerchantRequest) ToMerchantModel() (model interface{}, err error) {
	// 将MerchantType字符串转换为uint
	merchantTypeUint, err := strconv.ParseUint(req.MerchantType, 10, 32)
	if err != nil {
		return nil, err
	}
	merchantType := uint(merchantTypeUint)

	// 将MerchantLevel字符串转换为uint
	merchantLevelUint, err := strconv.ParseUint(req.MerchantLevel, 10, 32)
	if err != nil {
		return nil, err
	}
	merchantLevel := uint(merchantLevelUint)

	// 这里我们返回map，让service层处理具体的转换逻辑
	result := map[string]interface{}{
		"MerchantName":      req.MerchantName,
		"MerchantIcon":      req.MerchantIcon,
		"ParentID":          req.ParentID,
		"MerchantType":      merchantType,
		"BusinessLicense":   req.BusinessLicense,
		"LegalPerson":       req.LegalPerson,
		"RegisteredAddress": req.RegisteredAddress,
		"BusinessScope":     req.BusinessScope,
		"IsEnabled":         req.IsEnabled,
		"MerchantLevel":     merchantLevel,
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