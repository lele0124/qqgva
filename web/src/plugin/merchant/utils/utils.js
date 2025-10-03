// 商户模块工具函数

/**
 * 格式化日期时间
 * @param {string|number} date - 日期时间值
 * @param {string} format - 格式化模板
 * @returns {string} 格式化后的日期时间字符串
 */
export const formatDateTime = (date, format = 'YYYY-MM-DD HH:mm:ss') => {
  if (!date) return ''
  
  const d = new Date(date)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')
  const seconds = String(d.getSeconds()).padStart(2, '0')
  
  return format
    .replace('YYYY', year)
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds)
}

/**
 * 格式化日期
 * @param {string|number} date - 日期值
 * @returns {string} 格式化后的日期字符串
 */
export const formatDate = (date) => {
  return formatDateTime(date, 'YYYY-MM-DD')
}

/**
 * 格式化金额
 * @param {number|string} amount - 金额值
 * @param {number} decimal - 小数位数
 * @returns {string} 格式化后的金额字符串
 */
export const formatAmount = (amount, decimal = 2) => {
  if (amount === undefined || amount === null || amount === '') return '0.00'
  
  const num = Number(amount)
  if (isNaN(num)) return '0.00'
  
  return num.toFixed(decimal).replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

/**
 * 转换商户类型
 * @param {number} type - 商户类型值
 * @returns {string} 商户类型名称
 */
export const formatMerchantType = (type) => {
  const typeMap = {
    1: '线上商城',
    2: '实体店铺',
    3: '服务提供商',
    4: '其他'
  }
  return typeMap[type] || '未知类型'
}

/**
 * 转换商户等级
 * @param {number} level - 商户等级值
 * @returns {string} 商户等级名称
 */
export const formatMerchantLevel = (level) => {
  const levelMap = {
    1: '钻石',
    2: '金牌',
    3: '银牌',
    4: '铜牌',
    5: '普通'
  }
  return levelMap[level] || '未知等级'
}

/**
 * 转换启用状态
 * @param {boolean|number} status - 启用状态值
 * @returns {string} 状态名称
 */
export const formatStatus = (status) => {
  return status ? '启用' : '禁用'
}

/**
 * 过滤空值对象
 * @param {object} obj - 待过滤的对象
 * @returns {object} 过滤后的对象
 */
export const filterEmptyValues = (obj) => {
  if (!obj || typeof obj !== 'object') return obj
  
  const result = {}
  for (const key in obj) {
    if (obj.hasOwnProperty(key)) {
      const value = obj[key]
      if (value !== undefined && value !== null && value !== '' && !(Array.isArray(value) && value.length === 0)) {
        result[key] = value
      }
    }
  }
  return result
}

/**
 * 深拷贝对象
 * @param {any} source - 源对象
 * @returns {any} 拷贝后的对象
 */
export const deepClone = (source) => {
  if (source === null || typeof source !== 'object') return source
  
  if (source instanceof Date) return new Date(source.getTime())
  if (source instanceof Array) return source.map(item => deepClone(item))
  
  const target = {}
  for (const key in source) {
    if (source.hasOwnProperty(key)) {
      target[key] = deepClone(source[key])
    }
  }
  return target
}

/**
 * 生成唯一ID
 * @returns {string} 唯一ID字符串
 */
export const generateUniqueId = () => {
  return Date.now().toString(36) + Math.random().toString(36).substr(2)
}

/**
 * 防抖函数
 * @param {Function} func - 要防抖的函数
 * @param {number} wait - 等待时间（毫秒）
 * @returns {Function} 防抖后的函数
 */
export const debounce = (func, wait = 300) => {
  let timeout = null
  return function(...args) {
    const context = this
    if (timeout) clearTimeout(timeout)
    timeout = setTimeout(() => {
      func.apply(context, args)
    }, wait)
  }
}

/**
 * 节流函数
 * @param {Function} func - 要节流的函数
 * @param {number} wait - 等待时间（毫秒）
 * @returns {Function} 节流后的函数
 */
export const throttle = (func, wait = 300) => {
  let previous = 0
  return function(...args) {
    const now = Date.now()
    if (now - previous > wait) {
      previous = now
      func.apply(this, args)
    }
  }
}