import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useMerchantStore } from '../store/merchant'
import { createMerchant, updateMerchant } from '../api/merchant'
import { formatDate } from './utils'
import { createMerchantValidationRules } from './validationRules'

/**
 * 管理商户表单弹窗的钩子
 * @returns {Object} 包含弹窗状态、表单数据、验证规则和相关方法
 */
export default function useMerchantDialog() {
  const merchantStore = useMerchantStore()
  const dialogVisible = ref(false)
  const loading = ref(false)
  const isEdit = ref(false)
  const merchantId = ref('')
  
  // 表单数据
  const formData = reactive({
    merchantName: '',
    merchantType: 1,
    parentID: '',
    merchantLevel: 1,
    businessLicense: '',
    legalPerson: '',
    licenseStartDate: '',
    licenseEndDate: '',
    contactPerson: '',
    contactPhone: '',
    address: '',
    isEnabled: 1
  })
  
  // 验证规则
  const rules = computed(() => createMerchantValidationRules())
  
  // 重置表单
  const resetForm = () => {
    Object.keys(formData).forEach(key => {
      formData[key] = ''
    })
    formData.merchantType = 1
    formData.merchantLevel = 1
    formData.isEnabled = 1
    merchantId.value = ''
    isEdit.value = false
  }
  
  // 打开新增弹窗
  const openAddDialog = () => {
    resetForm()
    dialogVisible.value = true
  }
  
  // 打开编辑弹窗
  const openEditDialog = (merchant) => {
    resetForm()
    // 填充表单数据
    Object.keys(formData).forEach(key => {
      if (merchant.hasOwnProperty(key)) {
        // 日期类型特殊处理
        if (key === 'licenseStartDate' || key === 'licenseEndDate') {
          if (merchant[key]) {
            formData[key] = new Date(merchant[key])
          }
        } else {
          formData[key] = merchant[key]
        }
      }
    })
    merchantId.value = merchant.id
    isEdit.value = true
    dialogVisible.value = true
  }
  
  // 关闭弹窗
  const closeDialog = () => {
    dialogVisible.value = false
    setTimeout(() => {
      resetForm()
    }, 300)
  }
  
  // 提交表单
  const submitForm = async (formRef) => {
    try {
      await formRef.validate()
      loading.value = true
      
      // 准备提交数据
      const submitData = { ...formData }
      // 格式化日期
      if (submitData.licenseStartDate) {
        submitData.licenseStartDate = formatDate(submitData.licenseStartDate)
      }
      if (submitData.licenseEndDate) {
        submitData.licenseEndDate = formatDate(submitData.licenseEndDate)
      }
      
      if (isEdit.value) {
        // 更新商户
        await updateMerchant({ ...submitData, id: merchantId.value })
        ElMessage.success('更新商户成功')
      } else {
        // 创建商户
        await createMerchant(submitData)
        ElMessage.success('创建商户成功')
      }
      
      // 关闭弹窗并刷新列表
      closeDialog()
      merchantStore.fetchMerchantList()
    } catch (error) {
      ElMessage.error(error?.message || '操作失败')
    } finally {
      loading.value = false
    }
  }
  
  return {
    dialogVisible,
    loading,
    formData,
    rules,
    isEdit,
    openAddDialog,
    openEditDialog,
    closeDialog,
    submitForm
  }
}