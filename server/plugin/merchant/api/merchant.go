package api

import (
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/merchant/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/merchant/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Merchant = new(merchant)

type merchant struct{}

// CreateMerchant 创建商户信息
// @Tags Merchant
// @Summary 创建商户信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.CreateMerchantRequest true "创建商户信息"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /merchant/createMerchant [post]
func (a *merchant) CreateMerchant(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 先使用map接收JSON数据，处理可能的类型转换问题
	var rawData map[string]interface{}
	err := c.ShouldBindJSON(&rawData)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 初始化请求模型
	var req request.CreateMerchantRequest

	// 手动提取并转换字段，确保类型正确
	if merchantName, ok := rawData["merchantName"].(string); ok {
		req.MerchantName = merchantName
	}
	if merchantIcon, ok := rawData["merchantIcon"].(string); ok {
		req.MerchantIcon = merchantIcon
	}
	if businessLicense, ok := rawData["businessLicense"].(string); ok {
		req.BusinessLicense = businessLicense
	}
	if legalPerson, ok := rawData["legalPerson"].(string); ok {
		req.LegalPerson = legalPerson
	}
	if registeredAddress, ok := rawData["registeredAddress"].(string); ok {
		req.RegisteredAddress = registeredAddress
	}
	if businessScope, ok := rawData["businessScope"].(string); ok {
		req.BusinessScope = businessScope
	}
	if validStartTime, ok := rawData["validStartTime"].(string); ok {
		req.ValidStartTime = validStartTime
	}
	if validEndTime, ok := rawData["validEndTime"].(string); ok {
		req.ValidEndTime = validEndTime
	}

	// 处理merchantType字段，进行类型转换
	if merchantTypeValue, ok := rawData["merchantType"]; ok && merchantTypeValue != nil {
		switch v := merchantTypeValue.(type) {
		case string:
			// 如果是字符串，尝试转换为uint
			if temp, err := strconv.ParseUint(v, 10, 64); err == nil {
				req.MerchantType = uint(temp)
			}
		case float64:
			// 如果是数字类型(float64是JSON数字的默认类型)
			req.MerchantType = uint(v)
		case uint:
			req.MerchantType = v
		}
	}

	// 处理merchantLevel字段，进行类型转换
	if merchantLevelValue, ok := rawData["merchantLevel"]; ok && merchantLevelValue != nil {
		switch v := merchantLevelValue.(type) {
		case string:
			// 如果是字符串，尝试转换为uint
			if temp, err := strconv.ParseUint(v, 10, 64); err == nil {
				req.MerchantLevel = uint(temp)
			}
		case float64:
			// 如果是数字类型(float64是JSON数字的默认类型)
			req.MerchantLevel = uint(v)
		case uint:
			req.MerchantLevel = v
		}
	}

	// 处理parentID字段，进行类型转换
	if parentIDValue, ok := rawData["parentID"]; ok && parentIDValue != nil && parentIDValue != "" {
		var parentID uint
		switch v := parentIDValue.(type) {
		case string:
			// 如果是字符串，尝试转换为uint
			if temp, err := strconv.ParseUint(v, 10, 64); err == nil {
				parentID = uint(temp)
				req.ParentID = &parentID
			}
		case float64:
			// 如果是数字类型(float64是JSON数字的默认类型)
			parentID = uint(v)
			req.ParentID = &parentID
		case uint:
			parentID = v
			req.ParentID = &parentID
		}
	}

	// 处理isEnabled字段，进行类型转换
	if isEnabledValue, ok := rawData["isEnabled"]; ok && isEnabledValue != nil {
		switch v := isEnabledValue.(type) {
		case string:
			// 如果是字符串，尝试转换为bool
			if v == "1" || v == "true" {
				req.IsEnabled = true
			} else {
				req.IsEnabled = false
			}
		case bool:
			req.IsEnabled = v
		case float64:
			// 如果是数字类型，非0为true
			req.IsEnabled = v != 0
		}
	}

	// 创建商户模型
	info := model.Merchant{}

	// 复制基本字段，不需要类型转换
	info.MerchantName = &req.MerchantName
	info.MerchantIcon = &req.MerchantIcon
	info.ParentID = req.ParentID
	info.MerchantType = &req.MerchantType
	info.BusinessLicense = &req.BusinessLicense
	info.LegalPerson = &req.LegalPerson
	info.RegisteredAddress = &req.RegisteredAddress
	info.BusinessScope = &req.BusinessScope
	info.IsEnabled = req.IsEnabled
	info.MerchantLevel = &req.MerchantLevel

	// 处理时间字段，只有非空时才尝试解析
	if req.ValidStartTime != "" {
		startTime, timeErr := time.Parse(time.RFC3339, req.ValidStartTime)
		if timeErr != nil {
			// 尝试其他常见的时间格式
			startTime, timeErr = time.Parse("2006-01-02 15:04:05", req.ValidStartTime)
			if timeErr != nil {
				response.FailWithMessage("开始时间格式错误:"+timeErr.Error(), c)
				return
			}
		}
		info.ValidStartTime = &startTime
	}

	if req.ValidEndTime != "" {
		endTime, timeErr := time.Parse(time.RFC3339, req.ValidEndTime)
		if timeErr != nil {
			// 尝试其他常见的时间格式
			endTime, timeErr = time.Parse("2006-01-02 15:04:05", req.ValidEndTime)
			if timeErr != nil {
				response.FailWithMessage("结束时间格式错误:"+timeErr.Error(), c)
				return
			}
		}
		info.ValidEndTime = &endTime
	}

	// 设置操作人信息
	info.OperatorId = utils.GetUserID(c)
	info.OperatorName = utils.GetUserName(c)

	err = serviceMerchant.CreateMerchant(ctx, &info)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteMerchant 删除商户信息
// @Tags Merchant
// @Summary 删除商户信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Merchant true "删除商户信息"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /merchant/deleteMerchant [delete]
func (a *merchant) DeleteMerchant(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := serviceMerchant.DeleteMerchant(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteMerchantByIds 批量删除商户信息
// @Tags Merchant
// @Summary 批量删除商户信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /merchant/deleteMerchantByIds [delete]
func (a *merchant) DeleteMerchantByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := serviceMerchant.DeleteMerchantByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateMerchant 更新商户信息
// @Tags Merchant
// @Summary 更新商户信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.UpdateMerchantRequest true "更新商户信息"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /merchant/updateMerchant [put]
func (a *merchant) UpdateMerchant(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var req request.UpdateMerchantRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 将请求模型转换为数据模型
	modelData, err := req.ToMerchantModel()
	if err != nil {
		response.FailWithMessage("数据转换错误:"+err.Error(), c)
		return
	}

	// 获取原始商户数据
	originalMerchant, err := serviceMerchant.GetMerchant(ctx, strconv.Itoa(int(req.ID)))
	if err != nil {
		response.FailWithMessage("获取商户信息失败:"+err.Error(), c)
		return
	}

	// 创建更新的商户模型
	info := originalMerchant

	// 复制请求中的字段到模型
	if dataMap, ok := modelData.(map[string]interface{}); ok {
		if merchantType, ok := dataMap["MerchantType"].(uint); ok {
			info.MerchantType = &merchantType
		}
		if merchantName, ok := dataMap["MerchantName"].(string); ok {
			info.MerchantName = &merchantName
		}
		if merchantIcon, ok := dataMap["MerchantIcon"].(string); ok {
			info.MerchantIcon = &merchantIcon
		}
		if parentID, ok := dataMap["ParentID"].(*uint); ok {
			info.ParentID = parentID
		}
		if businessLicense, ok := dataMap["BusinessLicense"].(string); ok {
			info.BusinessLicense = &businessLicense
		}
		if legalPerson, ok := dataMap["LegalPerson"].(string); ok {
			info.LegalPerson = &legalPerson
		}
		if registeredAddress, ok := dataMap["RegisteredAddress"].(string); ok {
			info.RegisteredAddress = &registeredAddress
		}
		if businessScope, ok := dataMap["BusinessScope"].(string); ok {
			info.BusinessScope = &businessScope
		}
		if isEnabled, ok := dataMap["IsEnabled"].(bool); ok {
			info.IsEnabled = isEnabled
		}
		if merchantLevel, ok := dataMap["MerchantLevel"].(uint); ok {
			info.MerchantLevel = &merchantLevel
		}
		if validStartTime, ok := dataMap["ValidStartTime"].(*time.Time); ok {
			info.ValidStartTime = validStartTime
		}
		if validEndTime, ok := dataMap["ValidEndTime"].(*time.Time); ok {
			info.ValidEndTime = validEndTime
		}
	}

	// 设置操作人信息
	info.OperatorId = utils.GetUserID(c)
	info.OperatorName = utils.GetUserName(c)

	err = serviceMerchant.UpdateMerchant(ctx, info)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindMerchant 用id查询商户信息
// @Tags Merchant
// @Summary 用id查询商户信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询商户信息"
// @Success 200 {object} response.Response{data=model.Merchant,msg=string} "查询成功"
// @Router /merchant/findMerchant [get]
func (a *merchant) FindMerchant(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	remerchant, err := serviceMerchant.GetMerchant(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(remerchant, c)
}

// GetMerchantList 分页获取商户信息列表
// @Tags Merchant
// @Summary 分页获取商户信息列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.MerchantSearch true "分页获取商户信息列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /merchant/getMerchantList [post]
func (a *merchant) GetMerchantList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 先使用map接收JSON数据，处理可能的类型转换问题
	var rawData map[string]interface{}
	err := c.ShouldBindJSON(&rawData)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 初始化pageInfo
	var pageInfo request.MerchantSearch

	// 手动提取并转换字段，确保类型正确
	if page, ok := rawData["page"].(float64); ok {
		pageInfo.Page = int(page)
	}
	if pageSize, ok := rawData["pageSize"].(float64); ok {
		pageInfo.PageSize = int(pageSize)
	}
	if sort, ok := rawData["sort"].(string); ok {
		pageInfo.Sort = sort
	}
	if order, ok := rawData["order"].(string); ok {
		pageInfo.Order = order
	}

	// 处理merchantType字段，进行类型转换
	if merchantTypeValue, ok := rawData["merchantType"]; ok && merchantTypeValue != nil {
		var merchantType uint
		switch v := merchantTypeValue.(type) {
		case string:
			// 如果是字符串，尝试转换为uint
			var temp int64
			if temp, err = strconv.ParseInt(v, 10, 64); err == nil {
				merchantType = uint(temp)
				pageInfo.MerchantType = &merchantType
			}
		case float64:
			// 如果是数字类型(float64是JSON数字的默认类型)
			merchantType = uint(v)
			pageInfo.MerchantType = &merchantType
		}
	}

	// 处理其他可选的指针类型字段
	if merchantName, ok := rawData["merchantName"].(string); ok && merchantName != "" {
		pageInfo.MerchantName = &merchantName
	}
	if address, ok := rawData["address"].(string); ok && address != "" {
		pageInfo.Address = &address
	}
	if businessScope, ok := rawData["businessScope"].(string); ok && businessScope != "" {
		pageInfo.BusinessScope = &businessScope
	}
	if isEnabled, ok := rawData["isEnabled"].(bool); ok {
		pageInfo.IsEnabled = &isEnabled
	}
	if merchantLevel, ok := rawData["merchantLevel"].(float64); ok {
		level := uint(merchantLevel)
		pageInfo.MerchantLevel = &level
	}

	// 处理时间范围
	if createdAtRange, ok := rawData["createdAtRange"].([]interface{}); ok && len(createdAtRange) == 2 {
		pageInfo.CreatedAtRange = make([]time.Time, 2)
		// 解析开始时间
		if startTimeStr, ok := createdAtRange[0].(string); ok && startTimeStr != "" {
			if startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr); err == nil {
				pageInfo.CreatedAtRange[0] = startTime
			}
		}
		// 解析结束时间
		if endTimeStr, ok := createdAtRange[1].(string); ok && endTimeStr != "" {
			if endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr); err == nil {
				pageInfo.CreatedAtRange[1] = endTime
			}
		}
	}

	// 调用服务层获取数据
	list, total, err := serviceMerchant.GetMerchantInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetMerchantPublic 不需要鉴权的商户信息接口
// @Tags Merchant
// @Summary 不需要鉴权的商户信息接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /merchant/getMerchantPublic [get]
func (a *merchant) GetMerchantPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权 示例为返回了一个固定的消息接口,一般本接口用于C端服务,需要自己实现业务逻辑
	serviceMerchant.GetMerchantPublic(ctx)
	response.OkWithDetailed(gin.H{"info": "不需要鉴权的商户信息接口信息"}, "获取成功", c)
}
