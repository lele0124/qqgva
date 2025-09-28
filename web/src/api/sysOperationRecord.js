import service from '@/utils/request'

/**
 * @Tags SysOperationRecord
 * @Summary 删除单条操作日志
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body model.SysOperationRecord true "删除操作日志"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
 * @Router /sysOperationRecord/deleteSysOperationRecord [delete]
 */
export const deleteSysOperationRecord = (data) => {
  return service({
    url: '/sysOperationRecord/deleteSysOperationRecord',
    method: 'delete',
    data
  })
}

/**
 * @Tags SysOperationRecord
 * @Summary 批量删除操作日志
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body request.IdsReq true "批量删除操作日志"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
 * @Router /sysOperationRecord/deleteSysOperationRecordByIds [delete]
 */
export const deleteSysOperationRecordByIds = (data) => {
  return service({
    url: '/sysOperationRecord/deleteSysOperationRecordByIds',
    method: 'delete',
    data
  })
}

/**
 * @Tags SysOperationRecord
 * @Summary 分页获取操作日志列表
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data query request.PageInfo true "分页获取操作日志列表"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
 * @Router /sysOperationRecord/getSysOperationRecordList [get]
 */
export const getSysOperationRecordList = (params) => {
  return service({
    url: '/sysOperationRecord/getSysOperationRecordList',
    method: 'get',
    params
  })
}

/**
 * @Tags SysOperationRecord
 * @Summary 根据ID查询操作日志
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data query request.GetById true "根据ID查询操作日志"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
 * @Router /sysOperationRecord/findSysOperationRecord [get]
 */
export const findSysOperationRecord = (params) => {
  return service({
    url: '/sysOperationRecord/findSysOperationRecord',
    method: 'get',
    params
  })
}
