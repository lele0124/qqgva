/**
 * 统一数据处理工具函数
 * 用于处理表单数据的类型转换
 */

// 处理商户表单数据
export const processMerchantFormData = (data) => {
  const processedData = { ...data }

  // 处理数字类型字段
  const numericFields = ['parentID', 'merchantType', 'merchantLevel']
  numericFields.forEach(field => {
    if (processedData[field] !== undefined && processedData[field] !== null && processedData[field] !== '') {
      processedData[field] = parseInt(processedData[field])
    } else {
      // 删除空值字段
      delete processedData[field]
    }
  })

  // 处理布尔类型字段
  if (typeof processedData.isEnabled === 'string') {
    processedData.isEnabled = processedData.isEnabled === 'true' || processedData.isEnabled === '1'
  }

  return processedData
}

// 处理搜索条件数据
export const processSearchData = (data) => {
  const processedData = { ...data }

  // 处理数字类型字段
  const numericFields = ['parentID', 'merchantType', 'merchantLevel']
  numericFields.forEach(field => {
    if (processedData[field] !== undefined && processedData[field] !== null && processedData[field] !== '') {
      processedData[field] = parseInt(processedData[field])
    } else {
      // 删除空值字段而不是设置为undefined
      delete processedData[field]
    }
  })

  // 处理布尔类型字段
  if (processedData.isEnabled !== undefined && processedData.isEnabled !== null && processedData.isEnabled !== '') {
    processedData.isEnabled = processedData.isEnabled === '1' || processedData.isEnabled === true
  } else {
    // 删除空值字段而不是设置为undefined
    delete processedData.isEnabled
  }
  
  // 处理字符串类型字段，删除空字符串字段
  const stringFields = ['merchantName', 'merchantIcon', 'businessLicense', 'legalPerson', 
                       'registeredAddress', 'businessScope', 'validStartTime', 'validEndTime', 'address', 'status']
  stringFields.forEach(field => {
    if (processedData[field] !== undefined && processedData[field] !== null && processedData[field] === '') {
      delete processedData[field]
    }
  })

  return processedData
}

// 处理日期时间字段
export const processDateFields = (data, dateFields) => {
  const processedData = { ...data }
  
  dateFields.forEach(field => {
    if (processedData[field] && typeof processedData[field] === 'string') {
      processedData[field] = new Date(processedData[field])
    }
  })
  
  return processedData
}

// 验证营业执照编号格式
export const validateBusinessLicense = (rule, value, callback) => {
  if (!value) {
    callback(new Error('请输入营业执照编号'))
  } else if (!/^[A-Z0-9]{10,30}$/.test(value)) {
    callback(new Error('营业执照编号格式不正确，应为10-30位大写字母或数字'))
  } else {
    callback()
  }
}

// 清理空值字段
export const cleanEmptyFields = (data) => {
  const cleanedData = { ...data }
  Object.keys(cleanedData).forEach(key => {
    if (cleanedData[key] === '' || cleanedData[key] === null || cleanedData[key] === undefined) {
      delete cleanedData[key]
    }
  })
  return cleanedData
}