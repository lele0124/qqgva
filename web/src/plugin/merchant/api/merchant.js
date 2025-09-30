import service from '@/utils/request'
// @Tags Merchant
// @Summary 创建商户信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Merchant true "创建商户信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /merchant/createMerchant [post]
export const createMerchant = (data) => {
  return service({
    url: '/merchant/createMerchant',
    method: 'post',
    data
  })
}

// @Tags Merchant
// @Summary 删除商户信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Merchant true "删除商户信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /merchant/deleteMerchant [delete]
export const deleteMerchant = (params) => {
  return service({
    url: '/merchant/deleteMerchant',
    method: 'delete',
    params
  })
}

// @Tags Merchant
// @Summary 批量删除商户信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除商户信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /merchant/deleteMerchant [delete]
export const deleteMerchantByIds = (params) => {
  return service({
    url: '/merchant/deleteMerchantByIds',
    method: 'delete',
    params
  })
}

// @Tags Merchant
// @Summary 更新商户信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Merchant true "更新商户信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /merchant/updateMerchant [put]
export const updateMerchant = (data) => {
  return service({
    url: '/merchant/updateMerchant',
    method: 'put',
    data
  })
}

// @Tags Merchant
// @Summary 用id查询商户信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.Merchant true "用id查询商户信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /merchant/findMerchant [get]
export const findMerchant = (params) => {
  return service({
    url: '/merchant/findMerchant',
    method: 'get',
    params
  })
}

// @Tags Merchant
// @Summary 分页获取商户信息列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "分页获取商户信息列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /merchant/getMerchantList [post]
export const getMerchantList = (data) => {
  return service({
    url: '/merchant/getMerchantList',
    method: 'post',
    data
  })
}
// @Tags Merchant
// @Summary 不需要鉴权的商户信息接口
// @Accept application/json
// @Produce application/json
// @Param data query request.MerchantSearch true "分页获取商户信息列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /merchant/getMerchantPublic [get]
export const getMerchantPublic = () => {
  return service({
    url: '/merchant/getMerchantPublic',
    method: 'get',
  })
}
