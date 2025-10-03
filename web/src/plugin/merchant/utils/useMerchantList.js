import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useMerchantStore } from '../store/merchant'
import { deleteMerchantByIds } from '../api/merchant'
import { createSearchValidationRules } from './validationRules'

/**
 * 管理商户列表数据的钩子
 * @returns {Object} 包含列表数据、分页信息、搜索条件和相关方法
 */
export default function useMerchantList() {
  const merchantStore = useMerchantStore()
  
  // 列表数据（从store中获取）
  const merchantList = computed(() => merchantStore.merchantList)
  const total = computed(() => merchantStore.total)
  const loading = computed(() => merchantStore.loading)
  
  // 分页信息
  const pagination = reactive({
    currentPage: 1,
    pageSize: 10,
    total: 0
  })
  
  // 搜索条件
  const searchParams = reactive({
    merchantName: '',
    merchantType: '',
    isEnabled: '',
    startTime: '',
    endTime: ''
  })
  
  // 搜索验证规则
  const searchRules = computed(() => createSearchValidationRules())
  
  // 选中的行
  const selectedRows = ref([])
  const selectionChange = (rows) => {
    selectedRows.value = rows
  }
  
  // 获取列表数据
  const getMerchantList = (resetPage = false) => {
    if (resetPage) {
      pagination.currentPage = 1
    }
    
    // 准备搜索参数
    const params = {
      page: pagination.currentPage,
      pageSize: pagination.pageSize,
      ...searchParams
    }
    
    // 执行搜索
    merchantStore.fetchMerchantList(params)
  }
  
  // 搜索
  const handleSearch = async (formRef) => {
    try {
      if (formRef) {
        await formRef.validate()
      }
      getMerchantList(true)
    } catch (error) {
      ElMessage.error('搜索条件验证失败')
    }
  }
  
  // 重置搜索条件
  const resetSearch = () => {
    Object.keys(searchParams).forEach(key => {
      searchParams[key] = ''
    })
    getMerchantList(true)
  }
  
  // 分页变化
  const handleSizeChange = (size) => {
    pagination.pageSize = size
    getMerchantList()
  }
  
  const handleCurrentChange = (current) => {
    pagination.currentPage = current
    getMerchantList()
  }
  
  // 批量删除
  const batchDelete = async () => {
    if (selectedRows.value.length === 0) {
      ElMessage.warning('请选择要删除的商户')
      return
    }
    
    try {
      const ids = selectedRows.value.map(item => item.id)
      await deleteMerchantByIds({ ids })
      ElMessage.success('删除成功')
      // 刷新列表
      getMerchantList()
      // 清空选中状态
      selectedRows.value = []
    } catch (error) {
      ElMessage.error(error?.message || '删除失败')
    }
  }
  
  // 排序处理
  const handleSortChange = (sort) => {
    if (sort.prop && sort.order) {
      searchParams.sortField = sort.prop
      searchParams.sortOrder = sort.order === 'ascending' ? 'asc' : 'desc'
    } else {
      delete searchParams.sortField
      delete searchParams.sortOrder
    }
    getMerchantList(true)
  }
  
  // 启用/禁用商户
  const toggleMerchantStatus = async (id, isEnabled) => {
    try {
      await merchantStore.updateMerchantStatus(id, isEnabled ? 1 : 0)
      ElMessage.success(`${isEnabled ? '启用' : '禁用'}成功`)
      getMerchantList()
    } catch (error) {
      ElMessage.error(`${isEnabled ? '启用' : '禁用'}失败`)
    }
  }
  
  // 初始化时加载数据
  const init = () => {
    getMerchantList()
  }
  
  return {
    merchantList,
    total,
    loading,
    pagination,
    searchParams,
    searchRules,
    selectedRows,
    getMerchantList,
    handleSearch,
    resetSearch,
    handleSizeChange,
    handleCurrentChange,
    selectionChange,
    batchDelete,
    handleSortChange,
    toggleMerchantStatus,
    init
  }
}