<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rules" label-width="120px">
        <el-form-item label="商户名称:" prop="merchantName">
          <el-input v-model="formData.merchantName" :clearable="true" placeholder="请输入商户名称" />
        </el-form-item>
        <el-form-item label="商户图标:" prop="merchantIcon">
          <el-input v-model="formData.merchantIcon" :clearable="true" placeholder="请输入商户图标URL" />
        </el-form-item>
        <el-form-item label="商户类型:" prop="merchantType">
          <el-select v-model="formData.merchantType" placeholder="请选择商户类型" style="width: 100%">
            <el-option label="企业" :value="1" />
            <el-option label="个体" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="父商户ID:" prop="parentID">
          <el-input v-model="formData.parentID" :clearable="true" placeholder="请输入父商户ID（可选）" />
        </el-form-item>
        <el-form-item label="营业执照:" prop="businessLicense">
          <el-input v-model="formData.businessLicense" :clearable="true" placeholder="请输入营业执照编号" />
        </el-form-item>
        <el-form-item label="法人姓名:" prop="legalPerson">
          <el-input v-model="formData.legalPerson" :clearable="true" placeholder="请输入法人姓名" />
        </el-form-item>
        <el-form-item label="注册地址:" prop="registeredAddress">
          <el-input v-model="formData.registeredAddress" :clearable="true" placeholder="请输入注册地址" />
        </el-form-item>
        <el-form-item label="经营范围:" prop="businessScope">
          <el-input v-model="formData.businessScope" :clearable="true" placeholder="请输入经营范围" type="textarea" />
        </el-form-item>
        <el-form-item label="开关状态:" prop="isEnabled">
          <el-select v-model="formData.isEnabled" placeholder="请选择开关状态" style="width: 100%">
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item label="商户等级:" prop="merchantLevel">
          <el-select v-model="formData.merchantLevel" placeholder="请选择商户等级" style="width: 100%">
            <el-option label="普通商户" :value="1" />
            <el-option label="高级商户" :value="2" />
            <el-option label="VIP商户" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="有效期开始:" prop="validStartTime">
          <el-date-picker v-model="formData.validStartTime" type="datetime" placeholder="选择有效期开始时间" style="width: 100%" />
        </el-form-item>
        <el-form-item label="有效期结束:" prop="validEndTime">
          <el-date-picker v-model="formData.validEndTime" type="datetime" placeholder="选择有效期结束时间" style="width: 100%" />
        </el-form-item>
        <el-form-item>
          <el-button :loading="btnLoading" type="primary" @click="save">{{ type === 'create' ? '创建' : '更新' }}</el-button>
          <el-button type="default" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { createMerchant, updateMerchant, findMerchant } from '@/plugin/merchant/api/merchant'
import { useRoute, useRouter } from "vue-router"
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, computed, watch } from 'vue'
import { useMerchantStore } from '@/plugin/merchant/store/merchant'
import { processMerchantFormData, processDateFields } from '@/plugin/merchant/utils/dataProcessor'
import validationRules from '@/plugin/merchant/utils/validationRules'

// 定义组件属性
const props = defineProps({
  type: {
    type: String,
    default: 'create'
  },
  data: {
    type: Object,
    default: () => ({})
  }
})

// 定义事件
const emit = defineEmits(['submit', 'cancel'])

// 表单数据
const formData = ref({
  ID: '',
  merchantName: '',
  merchantIcon: '',
  merchantType: null,
  parentID: null,
  businessLicense: '',
  legalPerson: '',
  registeredAddress: '',
  businessScope: '',
  isEnabled: true,
  validStartTime: null,
  validEndTime: null,
  merchantLevel: null,
})

// 验证规则
const rules = computed(() => ({
  merchantName: validationRules.merchantNameRules,
  merchantType: validationRules.merchantTypeRules,
  isEnabled: validationRules.isEnabledRules,
  merchantLevel: validationRules.merchantLevelRules,
  businessLicense: validationRules.businessLicenseRules,
  legalPerson: validationRules.legalPersonRules,
  registeredAddress: validationRules.registeredAddressRules,
  businessScope: validationRules.businessScopeRules,
  parentID: validationRules.parentIDRules,
  validStartTime: validationRules.validStartTimeRules(formData),
  validEndTime: validationRules.validEndTimeRules(formData)
}))

const route = useRoute()
const router = useRouter()
const btnLoading = ref(false)
const elFormRef = ref()
const merchantStore = useMerchantStore()

// 计算属性
const isUpdate = computed(() => props.type === 'edit')

// 监听数据变化
watch(() => props.data, (newVal) => {
  if (newVal && Object.keys(newVal).length > 0) {
    formData.value = { ...newVal }
  }
}, { immediate: true })

// 初始化方法
const init = async () => {
  if (route.query.id) {
    try {
      const res = await findMerchant({ ID: route.query.id })
      if (res.code === 0) {
        // 使用统一的日期处理函数
        const processedData = processDateFields(res.data, ['validStartTime', 'validEndTime'])
        formData.value = processedData
      }
    } catch (error) {
      ElMessage.error('获取商户数据失败，请稍后重试')
      console.error('Failed to get merchant data:', error)
    }
  }
}

// 如果是通过路由访问而不是作为组件使用，则执行初始化
if (!props.data || Object.keys(props.data).length === 0) {
  init()
}

// 保存按钮
const save = async() => {
  btnLoading.value = true
  elFormRef.value?.validate(async (valid) => {
    if (!valid) {
      btnLoading.value = false
      return
    }
    
    try {
      // 使用统一的数据处理函数
      const submitData = processMerchantFormData({ ...formData.value })
      
      let res
      switch (props.type) {
        case 'create':
          res = await createMerchant(submitData)
          break
        case 'edit':
          res = await updateMerchant(submitData)
          break
        default:
          res = await createMerchant(submitData)
          break
      }
      
      btnLoading.value = false
      
      if (res.code === 0) {
        // 更新成功后，通知状态管理刷新列表数据
        merchantStore.fetchMerchantList()
        
        ElMessage({ type: 'success', message: props.type === 'create' ? '创建成功' : '更新成功' })
        
        // 触发提交事件
        emit('submit')
      } else {
        ElMessage.error(res.msg || (props.type === 'create' ? '创建失败' : '更新失败'))
      }
    } catch (error) {
      btnLoading.value = false
      ElMessage.error('操作失败：' + error.message || '未知错误')
      console.error('Save operation failed:', error)
    }
  })
}

// 返回按钮
const back = () => {
  emit('cancel')
  router.go(-1)
}
</script>

<style scoped>
.gva-form-box {
  padding: 20px;
  background-color: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
</style>