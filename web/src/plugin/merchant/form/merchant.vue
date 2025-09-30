
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="120px">
        <el-form-item label="商户名称:" prop="merchantName">
          <el-input v-model="formData.merchantName" :clearable="true" placeholder="请输入商户名称" />
        </el-form-item>
        <el-form-item label="商户图标:" prop="merchantIcon">
          <el-input v-model="formData.merchantIcon" :clearable="true" placeholder="请输入商户图标URL" />
        </el-form-item>
        <el-form-item label="商户类型:" prop="merchantType">
          <el-select v-model="formData.merchantType" placeholder="请选择商户类型" style="width: 100%">
            <el-option label="线上电商" value="线上电商" />
            <el-option label="线下实体店" value="线下实体店" />
            <el-option label="平台服务商" value="平台服务商" />
            <el-option label="其他" value="其他" />
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
            <el-option label="启用" value="1" />
            <el-option label="禁用" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="商户等级:" prop="merchantLevel">
          <el-select v-model="formData.merchantLevel" placeholder="请选择商户等级" style="width: 100%">
            <el-option label="一级商户" value="1" />
            <el-option label="二级商户" value="2" />
            <el-option label="三级商户" value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="有效期开始:" prop="validStartTime">
          <el-date-picker v-model="formData.validStartTime" type="datetime" placeholder="选择有效期开始时间" style="width: 100%" />
        </el-form-item>
        <el-form-item label="有效期结束:" prop="validEndTime">
          <el-date-picker v-model="formData.validEndTime" type="datetime" placeholder="选择有效期结束时间" style="width: 100%" />
        </el-form-item>
        <el-form-item>
          <el-button :loading="btnLoading" type="primary" @click="save">保存</el-button>
          <el-button type="default" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { createMerchant, updateMerchant, findMerchant } from '@/plugin/merchant/api/merchant'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'

// 表单数据
const formData = ref({
  ID: '',
  merchantName: '',
  merchantIcon: '',
  merchantType: '',
  parentID: null,
  businessLicense: '',
  legalPerson: '',
  registeredAddress: '',
  businessScope: '',
  isEnabled: 1,
  validStartTime: null,
  validEndTime: null,
  merchantLevel: '',
})

// 验证规则
const rule = reactive({
  merchantName: [
    { required: true, message: '请输入商户名称', trigger: ['input', 'blur'] },
    { whitespace: true, message: '不能只输入空格', trigger: ['input', 'blur'] }
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
    { required: true, message: '请输入营业执照编号', trigger: ['input', 'blur'] }
  ],
  legalPerson: [
    { required: true, message: '请输入法人姓名', trigger: ['input', 'blur'] }
  ],
  registeredAddress: [
    { required: true, message: '请输入注册地址', trigger: ['input', 'blur'] }
  ],
  businessScope: [
    { required: true, message: '请输入经营范围', trigger: ['input', 'blur'] }
  ]
})

const route = useRoute()
const router = useRouter()
const btnLoading = ref(false)
const type = ref('')
const elFormRef = ref()

// 初始化方法
const init = async () => {
  if (route.query.id) {
    const res = await findMerchant({ ID: route.query.id })
    if (res.code === 0) {
      formData.value = res.data
      type.value = 'update'
    }
  } else {
    type.value = 'create'
  }
}

init()

// 保存按钮
const save = async() => {
  btnLoading.value = true
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return btnLoading.value = false
    let res
    switch (type.value) {
      case 'create':
        res = await createMerchant(formData.value)
        break
      case 'update':
        res = await updateMerchant(formData.value)
        break
      default:
        res = await createMerchant(formData.value)
        break
    }
    btnLoading.value = false
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '创建/更改成功' })
    }
  })
}

// 返回按钮
const back = () => {
  router.go(-1)
}
</script>

<style>
</style>
