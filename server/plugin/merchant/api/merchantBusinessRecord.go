package api

import (

	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/plugin/merchant/model"
    "github.com/flipped-aurora/gin-vue-admin/server/plugin/merchant/model/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "github.com/flipped-aurora/gin-vue-admin/server/utils"
)

var MerchantBusinessRecord = new(merchant_record)

type merchant_record struct {}

// CreateMerchantBusinessRecord 创建商户业务记录
// @Tags MerchantBusinessRecord
// @Summary 创建商户业务记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.MerchantBusinessRecord true "创建商户业务记录"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /merchant_record/createMerchantBusinessRecord [post]
func (a *merchant_record) CreateMerchantBusinessRecord(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var info model.MerchantBusinessRecord
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    info.CreatedBy = utils.GetUserID(c)
	err = serviceMerchantBusinessRecord.CreateMerchantBusinessRecord(ctx,&info)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteMerchantBusinessRecord 删除商户业务记录
// @Tags MerchantBusinessRecord
// @Summary 删除商户业务记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.MerchantBusinessRecord true "删除商户业务记录"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /merchant_record/deleteMerchantBusinessRecord [delete]
func (a *merchant_record) DeleteMerchantBusinessRecord(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
    userID := utils.GetUserID(c)
	err := serviceMerchantBusinessRecord.DeleteMerchantBusinessRecord(ctx,ID,userID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("删除成功", c)
}

// DeleteMerchantBusinessRecordByIds 批量删除商户业务记录
// @Tags MerchantBusinessRecord
// @Summary 批量删除商户业务记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /merchant_record/deleteMerchantBusinessRecordByIds [delete]
func (a *merchant_record) DeleteMerchantBusinessRecordByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
    userID := utils.GetUserID(c)
	err := serviceMerchantBusinessRecord.DeleteMerchantBusinessRecordByIds(ctx,IDs,userID)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("批量删除成功", c)
}

// UpdateMerchantBusinessRecord 更新商户业务记录
// @Tags MerchantBusinessRecord
// @Summary 更新商户业务记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.MerchantBusinessRecord true "更新商户业务记录"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /merchant_record/updateMerchantBusinessRecord [put]
func (a *merchant_record) UpdateMerchantBusinessRecord(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var info model.MerchantBusinessRecord
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    info.UpdatedBy = utils.GetUserID(c)
	err = serviceMerchantBusinessRecord.UpdateMerchantBusinessRecord(ctx,info)
    if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("更新成功", c)
}

// FindMerchantBusinessRecord 用id查询商户业务记录
// @Tags MerchantBusinessRecord
// @Summary 用id查询商户业务记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询商户业务记录"
// @Success 200 {object} response.Response{data=model.MerchantBusinessRecord,msg=string} "查询成功"
// @Router /merchant_record/findMerchantBusinessRecord [get]
func (a *merchant_record) FindMerchantBusinessRecord(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	remerchant_record, err := serviceMerchantBusinessRecord.GetMerchantBusinessRecord(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
    response.OkWithData(remerchant_record, c)
}
// GetMerchantBusinessRecordList 分页获取商户业务记录列表
// @Tags MerchantBusinessRecord
// @Summary 分页获取商户业务记录列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.MerchantBusinessRecordSearch true "分页获取商户业务记录列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /merchant_record/getMerchantBusinessRecordList [get]
func (a *merchant_record) GetMerchantBusinessRecordList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo request.MerchantBusinessRecordSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceMerchantBusinessRecord.GetMerchantBusinessRecordInfoList(ctx,pageInfo)
	if err != nil {
	    global.GVA_LOG.Error("获取失败!", zap.Error(err))
        response.FailWithMessage("获取失败:" + err.Error(), c)
        return
    }
    response.OkWithDetailed(response.PageResult{
        List:     list,
        Total:    total,
        Page:     pageInfo.Page,
        PageSize: pageInfo.PageSize,
    }, "获取成功", c)
}
// GetMerchantBusinessRecordDataSource 获取MerchantBusinessRecord的数据源
// @Tags MerchantBusinessRecord
// @Summary 获取MerchantBusinessRecord的数据源
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /merchant_record/getMerchantBusinessRecordDataSource [get]
func (a *merchant_record) GetMerchantBusinessRecordDataSource(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口为获取数据源定义的数据
   dataSource, err := serviceMerchantBusinessRecord.GetMerchantBusinessRecordDataSource(ctx)
   if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
        response.FailWithMessage("查询失败:" + err.Error(), c)
		return
   }
    response.OkWithData(dataSource, c)
}
// GetMerchantBusinessRecordPublic 不需要鉴权的商户业务记录接口
// @Tags MerchantBusinessRecord
// @Summary 不需要鉴权的商户业务记录接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /merchant_record/getMerchantBusinessRecordPublic [get]
func (a *merchant_record) GetMerchantBusinessRecordPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权 示例为返回了一个固定的消息接口,一般本接口用于C端服务,需要自己实现业务逻辑
    serviceMerchantBusinessRecord.GetMerchantBusinessRecordPublic(ctx)
    response.OkWithDetailed(gin.H{"info": "不需要鉴权的商户业务记录接口信息"}, "获取成功", c)
}
