// 商户模块表单验证规则

import { reactive } from 'vue'
import { validateBusinessLicense } from './dataProcessor'

// 商户名称验证规则
export const merchantNameRules = [
  { required: true, message: '请输入商户名称', trigger: ['input', 'blur'] },
  { whitespace: true, message: '不能只输入空格', trigger: ['input', 'blur'] },
  { min: 2, max: 100, message: '商户名称长度应在2-100个字符之间', trigger: ['input', 'blur'] }
]

// 商户类型验证规则
export const merchantTypeRules = [
  { required: true, message: '请选择商户类型', trigger: ['change'] }
]

// 开关状态验证规则
export const isEnabledRules = [
  { required: true, message: '请选择商户开关状态', trigger: ['change'] }
]

// 商户等级验证规则
export const merchantLevelRules = [
  { required: true, message: '请选择商户等级', trigger: ['change'] }
]

// 营业执照验证规则
export const businessLicenseRules = [
  { required: true, validator: validateBusinessLicense, trigger: ['input', 'blur'] }
]

// 法人姓名验证规则
export const legalPersonRules = [
  { required: true, message: '请输入法人姓名', trigger: ['input', 'blur'] },
  { min: 2, max: 50, message: '法人姓名长度应在2-50个字符之间', trigger: ['input', 'blur'] }
]

// 注册地址验证规则
export const registeredAddressRules = [
  { required: true, message: '请输入注册地址', trigger: ['input', 'blur'] },
  { min: 5, max: 255, message: '注册地址长度应在5-255个字符之间', trigger: ['input', 'blur'] }
]

// 经营范围验证规则
export const businessScopeRules = [
  { required: true, message: '请输入经营范围', trigger: ['input', 'blur'] },
  { min: 5, max: 255, message: '经营范围长度应在5-255个字符之间', trigger: ['input', 'blur'] }
]

// 父商户ID验证规则
export const parentIDRules = [
  { type: 'number', message: '父商户ID必须是数字', trigger: ['input', 'blur'] }
]

// 有效期开始时间验证规则
export const validStartTimeRules = (formData) => [
  { type: 'date', message: '请选择有效的开始时间', trigger: ['change'] },
  {
    validator: (rule, value, callback) => {
      if (value && formData.value.validEndTime && value > formData.value.validEndTime) {
        callback(new Error('开始时间不能晚于结束时间'))
      } else {
        callback()
      }
    },
    trigger: ['change']
  }
]

// 有效期结束时间验证规则
export const validEndTimeRules = (formData) => [
  { type: 'date', message: '请选择有效的结束时间', trigger: ['change'] },
  {
    validator: (rule, value, callback) => {
      if (value && formData.value.validStartTime && value < formData.value.validStartTime) {
        callback(new Error('结束时间不能早于开始时间'))
      } else {
        callback()
      }
    },
    trigger: ['change']
  }
]

// 统一导出所有验证规则
export default {
  merchantNameRules,
  merchantTypeRules,
  isEnabledRules,
  merchantLevelRules,
  businessLicenseRules,
  legalPersonRules,
  registeredAddressRules,
  businessScopeRules,
  parentIDRules,
  validStartTimeRules,
  validEndTimeRules
}

/**
 * 创建商户表单验证规则
 * @param {object} formData - 表单数据对象（用于有效期时间比较）
 * @returns {object} 验证规则对象
 */
export const createMerchantValidationRules = (formData) => {
  return reactive({
    merchantName: [
      { required: true, message: '请输入商户名称', trigger: ['input', 'blur'] },
      { whitespace: true, message: '不能只输入空格', trigger: ['input', 'blur'] },
      { min: 2, max: 100, message: '商户名称长度应在2-100个字符之间', trigger: ['input', 'blur'] }
    ],
    merchantType: [
      { required: true, message: '请选择商户类型', trigger: ['change'] }
    ],
    isEnabled: [
      { required: true, message: '请选择商户开关状态', trigger: ['change'] }
    ],
    merchantLevel: [
      { required: true, message: '请选择商户等级', trigger: ['change'] }
    ],
    businessLicense: [
      { required: true, message: '请输入营业执照编号', trigger: ['input', 'blur'] },
      { pattern: /^[A-Z0-9]{15,20}$/, message: '营业执照编号格式不正确', trigger: ['input', 'blur'] }
    ],
    legalPerson: [
      { required: true, message: '请输入法人姓名', trigger: ['input', 'blur'] },
      { min: 2, max: 50, message: '法人姓名长度应在2-50个字符之间', trigger: ['input', 'blur'] }
    ],
    registeredAddress: [
      { required: true, message: '请输入注册地址', trigger: ['input', 'blur'] },
      { min: 5, max: 255, message: '注册地址长度应在5-255个字符之间', trigger: ['input', 'blur'] }
    ],
    businessScope: [
      { required: true, message: '请输入经营范围', trigger: ['input', 'blur'] },
      { min: 5, max: 255, message: '经营范围长度应在5-255个字符之间', trigger: ['input', 'blur'] }
    ],
    parentID: [
      { type: 'number', message: '父商户ID必须是数字', trigger: ['input', 'blur'] }
    ],
    validStartTime: [
      { type: 'date', message: '请选择有效的开始时间', trigger: ['change'] },
      {
        validator: (rule, value, callback) => {
          if (value && formData.value.validEndTime && value > formData.value.validEndTime) {
            callback(new Error('开始时间不能晚于结束时间'))
          } else {
            callback()
          }
        },
        trigger: ['change']
      }
    ],
    validEndTime: [
      { type: 'date', message: '请选择有效的结束时间', trigger: ['change'] },
      {
        validator: (rule, value, callback) => {
          if (value && formData.value.validStartTime && value < formData.value.validStartTime) {
            callback(new Error('结束时间不能早于开始时间'))
          } else {
            callback()
          }
        },
        trigger: ['change']
      }
    ]
  })
}

/**
 * 创建搜索表单验证规则
 * @returns {object} 搜索表单验证规则对象
 */
export const createSearchValidationRules = () => {
  return reactive({
    merchantName: [
      { max: 100, message: '商户名称不能超过100个字符', trigger: ['input', 'blur'] }
    ],
    merchantType: [
      { type: 'number', message: '商户类型必须是数字', trigger: ['change'] }
    ],
    isEnabled: [
      { type: 'boolean', message: '开关状态必须是布尔值', trigger: ['change'] }
    ],
    merchantLevel: [
      { type: 'number', message: '商户等级必须是数字', trigger: ['change'] }
    ],
    parentID: [
      { type: 'number', message: '父商户ID必须是数字', trigger: ['input', 'blur'] }
    ]
  })
}

/**
 * 验证日期范围
 * @param {Date} startTime - 开始时间
 * @param {Date} endTime - 结束时间
 * @returns {boolean} 验证结果
 */
export const validateDateRange = (startTime, endTime) => {
  if (!startTime || !endTime) return true
  return new Date(startTime) <= new Date(endTime)
}

/**
 * 验证数字范围
 * @param {number} value - 要验证的数字
 * @param {number} min - 最小值（可选）
 * @param {number} max - 最大值（可选）
 * @returns {boolean} 验证结果
 */
export const validateNumberRange = (value, min = undefined, max = undefined) => {
  const num = Number(value)
  if (isNaN(num)) return false
  if (min !== undefined && num < min) return false
  if (max !== undefined && num > max) return false
  return true
}

/**
 * 验证字符串长度
 * @param {string} value - 要验证的字符串
 * @param {number} min - 最小长度（可选）
 * @param {number} max - 最大长度（可选）
 * @returns {boolean} 验证结果
 */
export const validateStringLength = (value, min = undefined, max = undefined) => {
  if (typeof value !== 'string') return false
  if (min !== undefined && value.length < min) return false
  if (max !== undefined && value.length > max) return false
  return true
}