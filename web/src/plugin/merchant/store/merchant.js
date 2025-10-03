import { defineStore } from 'pinia'
import { getMerchantList, findMerchant, createMerchant, updateMerchant, deleteMerchant, deleteMerchantByIds } from '../api/merchant'
import { ElMessage } from 'element-plus'

// 定义商户管理的Store
export const useMerchantStore = defineStore('merchant', {
  state: () => ({
    // 表格数据
    tableData: [],
    // 分页数据
    page: 1,
    pageSize: 10,
    total: 0,
    // 搜索条件
    searchInfo: {
      merchantName: '',
      merchantIcon: '',
      merchantType: '',
      parentID: '',
      businessLicense: '',
      legalPerson: '',
      registeredAddress: '',
      businessScope: '',
      isEnabled: '',
      validStartTime: '',
      validEndTime: '',
      merchantLevel: '',
      address: '',
      status: '',
      updatedAtRange: []
    },
    // 选中的数据
    selectedData: [],
    // 表单数据
    formData: {
      ID: '',
      merchantName: '',
      merchantIcon: '',
      merchantType: '',
      parentID: '',
      businessLicense: '',
      legalPerson: '',
      registeredAddress: '',
      businessScope: '',
      isEnabled: '1',
      validStartTime: '',
      validEndTime: '',
      merchantLevel: '1',
      address: ''
    },
    // 详情数据
    detailForm: {},
    // 加载状态
    loading: {
      table: false,
      submit: false,
      detail: false
    },
    // 弹窗状态
    dialogVisible: {
      form: false,
      detail: false
    },
    // 表单类型 (create/edit)
    formType: 'create'
  }),

  getters: {
    // 格式化的选中ID列表
    selectedIds: (state) => state.selectedData.map(item => item.ID)
  },

  actions: {
    // 重置搜索条件
    resetSearchInfo() {
      this.searchInfo = {
        merchantName: '',
        merchantIcon: '',
        merchantType: '',
        parentID: '',
        businessLicense: '',
        legalPerson: '',
        registeredAddress: '',
        businessScope: '',
        isEnabled: '',
        validStartTime: '',
        validEndTime: '',
        merchantLevel: '',
        address: '',
        status: '',
        updatedAtRange: []
      }
      this.page = 1
    },

    // 重置表单数据
    resetFormData() {
      this.formData = {
        ID: '',
        merchantName: '',
        merchantIcon: '',
        merchantType: '',
        parentID: '',
        businessLicense: '',
        legalPerson: '',
        registeredAddress: '',
        businessScope: '',
        isEnabled: '1',
        validStartTime: '',
        validEndTime: '',
        merchantLevel: '1',
        address: ''
      }
    },

    // 获取商户列表
    async fetchMerchantList() {
      try {
        this.loading.table = true
        
        // 创建提交数据的副本，进行类型转换
        const submitData = {
          page: this.page,
          pageSize: this.pageSize,
          ...this.searchInfo
        }
        
        // 对merchantType进行类型转换
        if (submitData.merchantType !== '') {
          submitData.merchantType = parseInt(submitData.merchantType)
        } else {
          // 如果为空字符串，设置为undefined，避免传递空字符串
          delete submitData.merchantType
        }
        
        // 对parentID进行类型转换（如果有值）
        if (submitData.parentID !== '') {
          submitData.parentID = parseInt(submitData.parentID)
        }
        
        // 对merchantLevel进行类型转换（如果有值）
        if (submitData.merchantLevel !== '') {
          submitData.merchantLevel = parseInt(submitData.merchantLevel)
        }
        
        // 对isEnabled进行类型转换（如果有值）
        if (submitData.isEnabled !== '') {
          submitData.isEnabled = submitData.isEnabled === '1' ? true : false
        }
        
        const res = await getMerchantList(submitData)
        if (res.code === 0) {
          this.tableData = res.data.list || []
          this.total = res.data.total || 0
        } else {
          ElMessage.error(res.msg || '获取数据失败')
        }
      } catch (error) {
        ElMessage.error('获取数据失败')
      } finally {
        this.loading.table = false
      }
    },

    // 获取商户详情
    async fetchMerchantDetail(id) {
      try {
        this.loading.detail = true
        const res = await findMerchant({ ID: id })
        if (res.code === 0) {
          this.detailForm = res.data
          return res.data
        } else {
          ElMessage.error(res.msg || '获取详情失败')
          return null
        }
      } catch (error) {
        ElMessage.error('获取详情失败')
        return null
      } finally {
        this.loading.detail = false
      }
    },

    // 创建商户
    async createMerchant(data) {
      try {
        this.loading.submit = true
        const res = await createMerchant(data)
        if (res.code === 0) {
          ElMessage.success('创建成功')
          return true
        } else {
          ElMessage.error(res.msg || '创建失败')
          return false
        }
      } catch (error) {
        ElMessage.error('创建失败')
        return false
      } finally {
        this.loading.submit = false
      }
    },

    // 更新商户
    async updateMerchant(data) {
      try {
        this.loading.submit = true
        const res = await updateMerchant(data)
        if (res.code === 0) {
          ElMessage.success('更新成功')
          return true
        } else {
          ElMessage.error(res.msg || '更新失败')
          return false
        }
      } catch (error) {
        ElMessage.error('更新失败')
        return false
      } finally {
        this.loading.submit = false
      }
    },

    // 删除单个商户
    async deleteMerchant(id) {
      try {
        this.loading.submit = true
        const res = await deleteMerchant({ ID: id })
        if (res.code === 0) {
          ElMessage.success('删除成功')
          return true
        } else {
          ElMessage.error(res.msg || '删除失败')
          return false
        }
      } catch (error) {
        ElMessage.error('删除失败')
        return false
      } finally {
        this.loading.submit = false
      }
    },

    // 批量删除商户
    async deleteMerchantBatch(ids) {
      try {
        this.loading.submit = true
        const res = await deleteMerchantByIds({ "IDs[]": ids })
        if (res.code === 0) {
          ElMessage.success('删除成功')
          return true
        } else {
          ElMessage.error(res.msg || '删除失败')
          return false
        }
      } catch (error) {
        ElMessage.error('删除失败')
        return false
      } finally {
        this.loading.submit = false
      }
    },

    // 打开表单弹窗
    openFormDialog(type, data = null) {
      this.formType = type
      this.resetFormData()
      
      if (type === 'edit' && data) {
        // 赋值
        Object.keys(this.formData).forEach(key => {
          this.formData[key] = data[key] !== undefined ? data[key] : ''
        })
      }
      
      this.dialogVisible.form = true
    },

    // 关闭表单弹窗
    closeFormDialog() {
      this.dialogVisible.form = false
      this.resetFormData()
    },

    // 打开详情弹窗
    openDetailDialog(data) {
      this.detailForm = data
      this.dialogVisible.detail = true
    },

    // 关闭详情弹窗
    closeDetailDialog() {
      this.dialogVisible.detail = false
      this.detailForm = {}
    },

    // 分页处理
    handlePageChange(val) {
      this.page = val
      this.fetchMerchantList()
    },

    // 页码大小变化
    handleSizeChange(val) {
      this.pageSize = val
      this.page = 1
      this.fetchMerchantList()
    },

    // 排序处理
    handleSortChange(obj) {
      if (obj.order === 'ascending') {
        this.searchInfo.sort = 'asc'
        this.searchInfo.field = obj.prop
      } else if (obj.order === 'descending') {
        this.searchInfo.sort = 'desc'
        this.searchInfo.field = obj.prop
      } else {
        this.searchInfo.sort = ''
        this.searchInfo.field = ''
      }
      this.fetchMerchantList()
    },

    // 选择数据变化
    handleSelectionChange(val) {
      this.selectedData = val
    },

    // 提交表单
    async submitForm() {
      // 创建提交数据的副本，进行类型转换
      const submitData = { ...this.formData }
      
      // 对merchantType进行严格的类型转换
      if (submitData.merchantType !== undefined && submitData.merchantType !== null && submitData.merchantType !== '') {
        submitData.merchantType = parseInt(submitData.merchantType)
      } else {
        // 确保有一个有效的整数值
        submitData.merchantType = 0
      }
      
      // 对parentID进行类型转换（如果有值）
      if (submitData.parentID !== undefined && submitData.parentID !== null && submitData.parentID !== '') {
        submitData.parentID = parseInt(submitData.parentID)
      }
      
      // 对merchantLevel进行严格的类型转换
      if (submitData.merchantLevel !== undefined && submitData.merchantLevel !== null && submitData.merchantLevel !== '') {
        submitData.merchantLevel = parseInt(submitData.merchantLevel)
      } else {
        // 确保有一个有效的整数值
        submitData.merchantLevel = 0
      }
      
      // 对isEnabled进行类型转换
      if (typeof submitData.isEnabled === 'string') {
        submitData.isEnabled = submitData.isEnabled === '1' ? true : false
      }
      
      let success
      if (this.formType === 'create') {
        success = await this.createMerchant(submitData)
      } else {
        success = await this.updateMerchant(submitData)
      }
      
      if (success) {
        this.closeFormDialog()
        this.fetchMerchantList()
      }
      
      return success
    }
  }
})