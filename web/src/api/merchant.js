import service from '@/utils/request'

// @Tags Merchant
// @Summary 创建商户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body {MerchantName:"string", MerchantType:1, BusinessLicense:"string", LegalPerson:"string", RegisteredAddress:"string", BusinessScope:"string", IsEnabled:true, ValidStartTime:"string", ValidEndTime:"string", MerchantLevel:1} true "创建商户信息"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /merchant/createMerchant [post]
export const createMerchant = (data) => {
  return service({
    url: '/merchant/createMerchant',
    method: 'post',
    data: data
  })
}

// @Tags Merchant
// @Summary 删除商户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body {ID:"string"} true "删除商户信息"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /merchant/deleteMerchant [delete]
export const deleteMerchant = (data) => {
  return service({
    url: '/merchant/deleteMerchant',
    method: 'delete',
    data: data
  })
}

// @Tags Merchant
// @Summary 批量删除商户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body {ids:["string"]} true "批量删除商户信息"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /merchant/deleteMerchantByIds [delete]
export const deleteMerchantByIds = (data) => {
  return service({
    url: '/merchant/deleteMerchantByIds',
    method: 'delete',
    data: data
  })
}

// @Tags Merchant
// @Summary 更新商户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body {ID:"string", MerchantName:"string", MerchantType:1, BusinessLicense:"string", LegalPerson:"string", RegisteredAddress:"string", BusinessScope:"string", IsEnabled:true, ValidStartTime:"string", ValidEndTime:"string", MerchantLevel:1} true "更新商户信息"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /merchant/updateMerchant [put]
export const updateMerchant = (data) => {
  return service({
    url: '/merchant/updateMerchant',
    method: 'put',
    data: data
  })
}

// @Tags Merchant
// @Summary 用id查询商户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query string true "用id查询商户信息"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /merchant/findMerchant [get]
export const findMerchant = (params) => {
  return service({
    url: '/merchant/findMerchant',
    method: 'get',
    params: params
  })
}

// @Tags Merchant
// @Summary 分页获取商户信息列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body {page:1, pageSize:10, merchantName:"string", merchantType:1, isEnabled:true, merchantLevel:1, createdAtRange:["string", "string"]} true "分页获取商户信息列表"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /merchant/getMerchantList [post]
export const getMerchantList = (data) => {
  return service({
    url: '/merchant/getMerchantList',
    method: 'post',
    data: data
  })
}

// @Tags Merchant
// @Summary 不需要鉴权的商户信息接口
// @accept application/json
// @Produce application/json
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /merchant/getMerchantPublic [get]
export const getMerchantPublic = () => {
  return service({
    url: '/merchant/getMerchantPublic',
    method: 'get'
  })
}

// 导出所有API函数作为一个对象，方便使用
export const merchantApi = {
  createMerchant,
  deleteMerchant,
  deleteMerchantByIds,
  updateMerchant,
  findMerchant,
  getMerchantList,
  getMerchantPublic
}

export default merchantApi