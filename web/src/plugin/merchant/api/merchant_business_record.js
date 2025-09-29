import service from '@/utils/request'
// @Tags MerchantBusinessRecord
// @Summary 创建商户业务记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.MerchantBusinessRecord true "创建商户业务记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /merchant_record/createMerchantBusinessRecord [post]
export const createMerchantBusinessRecord = (data) => {
  return service({
    url: '/merchant_record/createMerchantBusinessRecord',
    method: 'post',
    data
  })
}

// @Tags MerchantBusinessRecord
// @Summary 删除商户业务记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.MerchantBusinessRecord true "删除商户业务记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /merchant_record/deleteMerchantBusinessRecord [delete]
export const deleteMerchantBusinessRecord = (params) => {
  return service({
    url: '/merchant_record/deleteMerchantBusinessRecord',
    method: 'delete',
    params
  })
}

// @Tags MerchantBusinessRecord
// @Summary 批量删除商户业务记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除商户业务记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /merchant_record/deleteMerchantBusinessRecord [delete]
export const deleteMerchantBusinessRecordByIds = (params) => {
  return service({
    url: '/merchant_record/deleteMerchantBusinessRecordByIds',
    method: 'delete',
    params
  })
}

// @Tags MerchantBusinessRecord
// @Summary 更新商户业务记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.MerchantBusinessRecord true "更新商户业务记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /merchant_record/updateMerchantBusinessRecord [put]
export const updateMerchantBusinessRecord = (data) => {
  return service({
    url: '/merchant_record/updateMerchantBusinessRecord',
    method: 'put',
    data
  })
}

// @Tags MerchantBusinessRecord
// @Summary 用id查询商户业务记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.MerchantBusinessRecord true "用id查询商户业务记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /merchant_record/findMerchantBusinessRecord [get]
export const findMerchantBusinessRecord = (params) => {
  return service({
    url: '/merchant_record/findMerchantBusinessRecord',
    method: 'get',
    params
  })
}

// @Tags MerchantBusinessRecord
// @Summary 分页获取商户业务记录列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取商户业务记录列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /merchant_record/getMerchantBusinessRecordList [get]
export const getMerchantBusinessRecordList = (params) => {
  return service({
    url: '/merchant_record/getMerchantBusinessRecordList',
    method: 'get',
    params
  })
}
// @Tags MerchantBusinessRecord
// @Summary 获取数据源
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /merchant_record/findMerchantBusinessRecordDataSource [get]
export const getMerchantBusinessRecordDataSource = () => {
  return service({
    url: '/merchant_record/getMerchantBusinessRecordDataSource',
    method: 'get',
  })
}
// @Tags MerchantBusinessRecord
// @Summary 不需要鉴权的商户业务记录接口
// @Accept application/json
// @Produce application/json
// @Param data query request.MerchantBusinessRecordSearch true "分页获取商户业务记录列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /merchant_record/getMerchantBusinessRecordPublic [get]
export const getMerchantBusinessRecordPublic = () => {
  return service({
    url: '/merchant_record/getMerchantBusinessRecordPublic',
    method: 'get',
  })
}
