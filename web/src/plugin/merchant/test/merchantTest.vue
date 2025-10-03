<template>
  <div class="merchant-test-container p-4">
    <h1 class="text-2xl font-bold mb-6">商户模块测试页面</h1>
    
    <!-- API测试区域 -->
    <el-card class="mb-6">
      <template #header>
        <div class="font-bold">API接口测试</div>
      </template>
      
      <el-tabs v-model="activeApiTab" type="border-card">
        <!-- 获取商户列表 -->
        <el-tab-pane label="获取商户列表" name="getMerchantList">
          <el-form label-width="100px" label-position="left" size="small">
            <el-form-item label="页码">
              <el-input-number v-model="listParams.page" :min="1" :precision="0" />
            </el-form-item>
            <el-form-item label="每页数量">
              <el-input-number v-model="listParams.pageSize" :min="1" :max="100" :precision="0" />
            </el-form-item>
            <el-form-item label="搜索关键词">
              <el-input v-model="listParams.keyword" placeholder="请输入商户名称" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="testGetMerchantList">执行测试</el-button>
            </el-form-item>
          </el-form>
          
          <div v-if="apiResponse.getMerchantList" class="mt-4">
            <h4 class="font-semibold mb-2">测试结果:</h4>
            <pre class="bg-gray-100 p-3 rounded whitespace-pre-wrap break-words max-h-80 overflow-y-auto">{{ JSON.stringify(apiResponse.getMerchantList, null, 2) }}</pre>
          </div>
        </el-tab-pane>
        
        <!-- 获取单个商户 -->
        <el-tab-pane label="获取单个商户" name="findMerchant">
          <el-form label-width="100px" label-position="left" size="small">
            <el-form-item label="商户ID">
              <el-input v-model="singleParams.merchantId" placeholder="请输入商户ID" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="testFindMerchant">执行测试</el-button>
            </el-form-item>
          </el-form>
          
          <div v-if="apiResponse.findMerchant" class="mt-4">
            <h4 class="font-semibold mb-2">测试结果:</h4>
            <pre class="bg-gray-100 p-3 rounded whitespace-pre-wrap break-words max-h-80 overflow-y-auto">{{ JSON.stringify(apiResponse.findMerchant, null, 2) }}</pre>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>
    
    <!-- 验证规则测试 -->
    <el-card class="mb-6">
      <template #header>
        <div class="font-bold">验证规则测试</div>
      </template>
      
      <el-form :model="validationForm" ref="validationFormRef" label-width="150px" label-position="left" size="small">
        <el-form-item label="商户名称" prop="merchantName" :rules="validationRules.merchantName">
          <el-input v-model="validationForm.merchantName" placeholder="请输入商户名称" />
        </el-form-item>
        
        <el-form-item label="营业执照号" prop="licenseNumber" :rules="validationRules.licenseNumber">
          <el-input v-model="validationForm.licenseNumber" placeholder="请输入营业执照号" />
        </el-form-item>
        
        <el-form-item label="法人姓名" prop="legalPerson" :rules="validationRules.legalPerson">
          <el-input v-model="validationForm.legalPerson" placeholder="请输入法人姓名" />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="testValidationRules">执行验证</el-button>
        </el-form-item>
      </el-form>
      
      <div v-if="validationResult" class="mt-4">
        <h4 class="font-semibold mb-2">验证结果:</h4>
        <pre class="bg-gray-100 p-3 rounded whitespace-pre-wrap break-words">{{ JSON.stringify(validationResult, null, 2) }}</pre>
      </div>
    </el-card>
    
    <!-- 工具函数测试 -->
    <el-card class="mb-6">
      <template #header>
        <div class="font-bold">工具函数测试</div>
      </template>
      
      <el-tabs v-model="activeToolTab" type="border-card">
        <!-- 日期格式化 -->
        <el-tab-pane label="日期格式化" name="formatDate">
          <el-form label-width="120px" label-position="left" size="small">
            <el-form-item label="日期值">
              <el-date-picker v-model="toolParams.dateValue" type="datetime" placeholder="选择日期时间" />
            </el-form-item>
            <el-form-item label="格式化模板">
              <el-input v-model="toolParams.dateFormat" placeholder="如: yyyy-MM-dd HH:mm:ss" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="testFormatDate">执行测试</el-button>
            </el-form-item>
          </el-form>
          
          <div v-if="toolResults.formatDate" class="mt-4">
            <h4 class="font-semibold mb-2">测试结果:</h4>
            <p>{{ toolResults.formatDate }}</p>
          </div>
        </el-tab-pane>
        
        <!-- 金额格式化 -->
        <el-tab-pane label="金额格式化" name="formatAmount">
          <el-form label-width="120px" label-position="left" size="small">
            <el-form-item label="金额值">
              <el-input v-model="toolParams.amountValue" placeholder="请输入金额" type="number" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="testFormatAmount">执行测试</el-button>
            </el-form-item>
          </el-form>
          
          <div v-if="toolResults.formatAmount" class="mt-4">
            <h4 class="font-semibold mb-2">测试结果:</h4>
            <p>{{ toolResults.formatAmount }}</p>
          </div>
        </el-tab-pane>
        
        <!-- 商户状态转换 -->
        <el-tab-pane label="商户状态转换" name="formatMerchantStatus">
          <el-form label-width="120px" label-position="left" size="small">
            <el-form-item label="状态值">
              <el-select v-model="toolParams.statusValue" placeholder="请选择状态">
                <el-option label="正常" :value="1" />
                <el-option label="禁用" :value="2" />
                <el-option label="审核中" :value="3" />
                <el-option label="审核失败" :value="4" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="testFormatMerchantStatus">执行测试</el-button>
            </el-form-item>
          </el-form>
          
          <div v-if="toolResults.formatMerchantStatus" class="mt-4">
            <h4 class="font-semibold mb-2">测试结果:</h4>
            <p>{{ toolResults.formatMerchantStatus }}</p>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>
    
    <!-- 自定义钩子测试 -->
    <el-card>
      <template #header>
        <div class="font-bold">自定义钩子测试</div>
      </template>
      
      <el-button type="primary" @click="testMerchantDialogHook" class="mb-4">测试表单弹窗钩子</el-button>
      
      <div v-if="hookResults.useMerchantDialog" class="mt-4">
        <h4 class="font-semibold mb-2">钩子测试结果:</h4>
        <pre class="bg-gray-100 p-3 rounded whitespace-pre-wrap break-words">{{ JSON.stringify(hookResults.useMerchantDialog, null, 2) }}</pre>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { getMerchantList, findMerchant } from '@/plugin/merchant/api/merchant.js'
import { formatDate, formatAmount, formatMerchantStatus } from '@/plugin/merchant/utils/utils.js'
import { createMerchantValidationRules } from '@/plugin/merchant/utils/validationRules.js'
import { useMerchantDialog } from '@/plugin/merchant/utils/useMerchantDialog.js'

// API测试相关
const activeApiTab = ref('getMerchantList')
const listParams = reactive({
  page: 1,
  pageSize: 10,
  keyword: ''
})
const singleParams = reactive({
  merchantId: ''
})
const apiResponse = reactive({})

// 验证规则测试相关
const validationFormRef = ref(null)
const validationForm = reactive({
  merchantName: '',
  licenseNumber: '',
  legalPerson: ''
})
const validationRules = createMerchantValidationRules()
const validationResult = ref(null)

// 工具函数测试相关
const activeToolTab = ref('formatDate')
const toolParams = reactive({
  dateValue: new Date(),
  dateFormat: 'yyyy-MM-dd HH:mm:ss',
  amountValue: 1234567.89,
  statusValue: 1
})
const toolResults = reactive({})

// 自定义钩子测试相关
const hookResults = reactive({})

// API测试方法
const testGetMerchantList = async () => {
  try {
    const res = await getMerchantList(listParams)
    apiResponse.getMerchantList = res
    ElMessage.success('获取商户列表测试成功')
  } catch (error) {
    ElMessage.error('获取商户列表测试失败: ' + error.message)
    apiResponse.getMerchantList = { error: error.message }
  }
}

const testFindMerchant = async () => {
  if (!singleParams.merchantId) {
    ElMessage.warning('请输入商户ID')
    return
  }
  
  try {
    const res = await findMerchant({ id: singleParams.merchantId })
    apiResponse.findMerchant = res
    ElMessage.success('获取单个商户测试成功')
  } catch (error) {
    ElMessage.error('获取单个商户测试失败: ' + error.message)
    apiResponse.findMerchant = { error: error.message }
  }
}

// 验证规则测试方法
const testValidationRules = () => {
  if (!validationFormRef.value) return
  
  validationFormRef.value.validate((valid, errors) => {
    validationResult.value = {
      isValid: valid,
      errors: errors
    }
    
    if (valid) {
      ElMessage.success('表单验证通过')
    } else {
      ElMessage.warning('表单验证未通过')
    }
  })
}

// 工具函数测试方法
const testFormatDate = () => {
  try {
    const result = formatDate(toolParams.dateValue, toolParams.dateFormat)
    toolResults.formatDate = result
    ElMessage.success('日期格式化测试成功')
  } catch (error) {
    ElMessage.error('日期格式化测试失败: ' + error.message)
    toolResults.formatDate = { error: error.message }
  }
}

const testFormatAmount = () => {
  try {
    const result = formatAmount(toolParams.amountValue)
    toolResults.formatAmount = result
    ElMessage.success('金额格式化测试成功')
  } catch (error) {
    ElMessage.error('金额格式化测试失败: ' + error.message)
    toolResults.formatAmount = { error: error.message }
  }
}

const testFormatMerchantStatus = () => {
  try {
    const result = formatMerchantStatus(toolParams.statusValue)
    toolResults.formatMerchantStatus = result
    ElMessage.success('商户状态转换测试成功')
  } catch (error) {
    ElMessage.error('商户状态转换测试失败: ' + error.message)
    toolResults.formatMerchantStatus = { error: error.message }
  }
}

// 自定义钩子测试方法
const testMerchantDialogHook = () => {
  try {
    const hook = useMerchantDialog()
    hookResults.useMerchantDialog = {
      hasDialogVisible: typeof hook.dialogVisible === 'object',
      hasFormData: typeof hook.formData === 'object',
      hasResetForm: typeof hook.resetForm === 'function',
      hasOpenDialog: typeof hook.openDialog === 'function',
      hasCloseDialog: typeof hook.closeDialog === 'function',
      hasSubmitForm: typeof hook.submitForm === 'function'
    }
    ElMessage.success('自定义钩子测试成功')
  } catch (error) {
    ElMessage.error('自定义钩子测试失败: ' + error.message)
    hookResults.useMerchantDialog = { error: error.message }
  }
}
</script>

<style scoped>
.merchant-test-container {
  max-width: 1200px;
  margin: 0 auto;
}
</style>