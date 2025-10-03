package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/merchant/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/merchant/service"
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

	// 初始化商户模型
	var req request.CreateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 转换为数据模型
	info, err := req.ToMerchantModel()
	if err != nil {
		response.FailWithMessage("数据转换失败: "+err.Error(), c)
		return
	}

	// 设置操作人信息
	info.OperatorId = utils.GetUserID(c)
	info.OperatorName = utils.GetUserName(c)

	// 调用服务层创建商户
	if err := service.Service.Merchant.CreateMerchant(ctx, &info); err != nil {
		global.GVA_LOG.Error("创建失败！", zap.Error(err))
		response.FailWithMessage("创建失败: "+err.Error(), c)
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
	if err := service.Service.Merchant.DeleteMerchant(ctx, ID); err != nil {
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
	if err := service.Service.Merchant.DeleteMerchantByIds(ctx, IDs); err != nil {
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

	// 转换为数据模型
	info, err := req.ToMerchantModel()
	if err != nil {
		response.FailWithMessage("数据转换失败: "+err.Error(), c)
		return
	}

	// 设置操作人信息
	info.OperatorId = utils.GetUserID(c)
	info.OperatorName = utils.GetUserName(c)

	if err := service.Service.Merchant.UpdateMerchant(ctx, info); err != nil {
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
	remerchant, err := service.Service.Merchant.GetMerchant(ctx, ID)
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

	var pageInfo request.MerchantSearch
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 调用服务层获取数据
	list, total, err := service.Service.Merchant.GetMerchantInfoList(ctx, pageInfo)
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
	service.Service.Merchant.GetMerchantPublic(ctx)
	response.OkWithDetailed(gin.H{"info": "不需要鉴权的商户信息接口信息"}, "获取成功", c)
}
