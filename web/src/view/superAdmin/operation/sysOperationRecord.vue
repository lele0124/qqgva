<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="请求方法">
          <el-input v-model="searchInfo.method" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="请求路径">
          <el-input v-model="searchInfo.path" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item label="结果状态码">
          <el-input v-model="searchInfo.status" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit"
            >查询</el-button
          >
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button
          icon="delete"
          :disabled="!multipleSelection.length"
          @click="onDelete"
          >删除</el-button
        >
      </div>
      <el-table
        ref="multipleTable"
        :data="tableData"
        style="width: 100%"
        tooltip-effect="dark"
        row-key="ID"
        @selection-change="handleSelectionChange"
        @sort-change="handleSortChange"
        :default-sort="{prop: 'ID', order: 'descending'}"
      >
        <el-table-column align="left" type="selection" width="40" />
        <el-table-column align="left" label="操作ID" prop="id" width="100" sortable />
        <el-table-column align="left" label="创建时间" prop="createdAt" width="180" sortable>
          <template #default="scope">{{ formatDate(scope.row.createdAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="更新时间" prop="updatedAt" width="180" sortable>
          <template #default="scope">{{ formatDate(scope.row.updatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="请求IP" prop="ip" width="120" sortable />
        <el-table-column align="left" label="请求方法" prop="method" width="120" sortable />
        <el-table-column align="left" label="请求路径" prop="path" min-width="200" sortable />
        <el-table-column align="left" label="状态码" prop="status" width="100" sortable>
          <template #default="scope">
            <div>
              <el-tag :type="getStatusCodeTagType(scope.row.status)">{{ scope.row.status }}</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="请求耗时" prop="latency" width="120" sortable>
          <template #default="scope">{{ formatLatency(scope.row.latency) }}</template>
        </el-table-column>
        <el-table-column align="left" label="用户ID" prop="user_id" width="100" sortable />
        <el-table-column align="left" label="姓名" prop="user_name" width="100" />
        <el-table-column align="left" label="用户终端" width="180">
          <template #default="scope">
            <el-popover
              placement="top-start"
              :width="350"
              trigger="hover"
            >
              <div>{{ scope.row.agent }}</div>
              <template #reference>
                <span class="truncate-text">{{ scope.row.agent }}</span>
              </template>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column align="left" label="错误信息" width="120">
          <template #default="scope">
            <el-popover
              v-if="scope.row.error_message"
              placement="left-start"
              :width="350"
            >
              <div class="error-message-box">{{ scope.row.error_message }}</div>
              <template #reference>
                <el-tag type="danger" size="small">错误</el-tag>
              </template>
            </el-popover>
            <span v-else>无</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="请求体" width="110">
          <template #default="scope">
            <div>
              <el-popover
                v-if="scope.row.body"
                placement="left-start"
                :width="350"
              >
                <div class="popover-box">
                  <pre>{{ fmtBody(scope.row.body) }}</pre>
                </div>
                <template #reference>
                  <el-icon style="cursor: pointer"><Warning /></el-icon>
                </template>
              </el-popover>
              <span v-else>无</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="响应体" width="110">
          <template #default="scope">
            <div>
              <el-popover
                v-if="scope.row.resp"
                placement="left-start"
                :width="350"
              >
                <div class="popover-box">
                  <pre>{{ fmtBody(scope.row.resp) }}</pre>
                </div>
                <template #reference>
                  <el-icon style="cursor: pointer"><Warning /></el-icon>
                </template>
              </el-popover>
              <span v-else>无</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="用户信息" width="80">
          <template #default="scope">
            <div>
              <el-popover
                v-if="scope.row.user"
                placement="left-start"
                :width="350"
              >
                <div class="user-info-box">
                  <p><strong>用户名:</strong> {{ scope.row.user.userName }}</p>
                  <p><strong>昵称:</strong> {{ scope.row.user.nickName }}</p>
                  <p><strong>权限ID:</strong> {{ scope.row.user.authorityId }}</p>
                </div>
                <template #reference>
                  <el-button type="primary" size="small" link>详情</el-button>
                </template>
              </el-popover>
              <span v-else>无</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" width="60">
          <template #default="scope">
            <el-button
              icon="delete"
              type="primary"
              link
              @click="deleteSysOperationRecordFunc(scope.row)"
              >删除</el-button
            >
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
  import {
    deleteSysOperationRecord,
    getSysOperationRecordList,
    deleteSysOperationRecordByIds
  } from '@/api/sysOperationRecord' // 此处请自行替换地址
  import { formatDate } from '@/utils/format'
  import { ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { Warning } from '@element-plus/icons-vue'

  defineOptions({
    name: 'SysOperationRecord'
  })

  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])
  const searchInfo = ref({})
  // 排序相关状态
  const sortField = ref('')
  const sortOrder = ref('')
  
  const onReset = () => {
    searchInfo.value = {}
    sortField.value = ''
    sortOrder.value = ''
  }
  // 条件搜索前端看此方法
  const onSubmit = () => {
    page.value = 1
    if (searchInfo.value.status === '') {
      searchInfo.value.status = null
    }
    getTableData()
  }

  // 分页
  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  // 排序处理函数
  const handleSortChange = ({ prop, order }) => {
    sortField.value = prop
    sortOrder.value = order === 'ascending' ? 'asc' : order === 'descending' ? 'desc' : ''
    page.value = 1
    getTableData()
  }

  // 查询
  const getTableData = async () => {
    const table = await getSysOperationRecordList({
      page: page.value,
      pageSize: pageSize.value,
      SortField: sortField.value,
      SortOrder: sortOrder.value,
      ...searchInfo.value
    })
    if (table.code === 0) {
      tableData.value = table.data.list
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  }

  getTableData()

  const multipleSelection = ref([])
  const handleSelectionChange = (val) => {
    multipleSelection.value = val
  }
  const onDelete = async () => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      const ids = []
      multipleSelection.value &&
        multipleSelection.value.forEach((item) => {
          ids.push(item.ID)
        })
      const res = await deleteSysOperationRecordByIds({ ids })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === ids.length && page.value > 1) {
          page.value--
        }
        getTableData()
      }
    })
  }
  const deleteSysOperationRecordFunc = async (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      const res = await deleteSysOperationRecord({ ID: row.ID })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === 1 && page.value > 1) {
          page.value--
        }
        getTableData()
      }
    })
  }
  const fmtBody = (value) => {
    try {
      return JSON.parse(value)
    } catch (_) {
      return value
    }
  }
  
  // 根据HTTP状态码返回不同的标签类型
  const getStatusCodeTagType = (status) => {
    if (status >= 200 && status < 300) {
      return 'success'
    } else if (status >= 400 && status < 500) {
      return 'warning'
    } else if (status >= 500) {
      return 'danger'
    }
    return 'info'
  }
  
  // 格式化延迟时间，假设单位是纳秒
  const formatLatency = (latency) => {
    if (!latency) return '0'
    // 转换为毫秒
    return (latency / 1000000).toFixed(2)
  }
</script>

<style lang="scss">
  .table-expand {
    padding-left: 60px;
    font-size: 0;
    label {
      width: 90px;
      color: #99a9bf;
      .el-form-item {
        margin-right: 0;
        margin-bottom: 0;
        width: 50%;
      }
    }
  }
  .popover-box {
    background: #112435;
    color: #f08047;
    height: 600px;
    width: 420px;
    overflow: auto;
  }
  .popover-box::-webkit-scrollbar {
    display: none; /* Chrome Safari */
  }
  .truncate-text {
    display: inline-block;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .error-message-box {
    background: #fef0f0;
    color: #f56c6c;
    padding: 10px;
    max-height: 300px;
    overflow-y: auto;
  }
  .user-info-box {
    background: #f0f9ff;
    color: #2c3e50;
    padding: 10px;
  }
  /* 修复排序按钮点击问题的样式 */
  .el-table th {
    position: relative;
  }
  .el-table th > .cell {
    position: relative;
    display: flex;
    align-items: center;
  }
  .el-table th > .cell > .caret-wrapper {
    position: absolute;
    right: 8px;
    display: inline-flex;
    flex-direction: column;
    justify-content: center;
    height: 100%;
    width: 16px;
  }
  .el-table th > .cell > .caret-wrapper .sort-caret {
    cursor: pointer;
    opacity: 0.4;
    transition: opacity 0.2s, color 0.2s;
  }
  .el-table th > .cell > .caret-wrapper .sort-caret:hover {
    opacity: 1;
  }
</style>
