<template>
  <div>
    <el-card>
      <template #header>
        <div style="display:flex;justify-content:space-between">
          <span>组织管理</span>
          <el-button type="primary" @click="openDialog()">新增根组织</el-button>
        </div>
      </template>

      <el-tree
        :data="treeData"
        node-key="id"
        :props="{ label: 'org_name', children: 'children' }"
        default-expand-all
        :expand-on-click-node="false"
      >
        <template #default="{ node, data }">
          <span style="display:flex;align-items:center;width:100%">
            <span style="flex:1">{{ data.org_name }}</span>
            <span style="color:#999;margin-right:12px">{{ data.org_code }}</span>
            <el-button size="small" @click.stop="openDialog(data, true)">新增子</el-button>
            <el-button size="small" @click.stop="openDialog(data)">编辑</el-button>
            <el-button size="small" type="danger" @click.stop="remove(data)">删除</el-button>
          </span>
        </template>
      </el-tree>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑组织' : (isChild ? '新增子组织' : '新增根组织')" width="450px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="组织编码">
          <el-input v-model="form.org_code" />
        </el-form-item>
        <el-form-item label="组织名称">
          <el-input v-model="form.org_name" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort_order" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="save">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { orgApi } from '../api/mdm'

const treeData = ref<any[]>([])
const dialogVisible = ref(false)
const isChild = ref(false)
const form = reactive<any>({})

const buildTree = (list: any[], parentId = 0): any[] => {
  return list
    .filter((n: any) => n.parent_id === parentId)
    .map((n: any) => ({ ...n, children: buildTree(list, n.id) }))
}

const loadTree = async () => {
  const res: any = await orgApi.list(0)
  const all: any[] = res.data || []
  treeData.value = buildTree(all)
}

const openDialog = (row?: any, child = false) => {
  isChild.value = child
  if (child && row) {
    Object.assign(form, { parent_id: row.id, org_code: '', org_name: '', sort_order: 0 })
  } else {
    Object.assign(form, row || { parent_id: 0, org_code: '', org_name: '', sort_order: 0 })
  }
  dialogVisible.value = true
}

const save = async () => {
  if (form.id) {
    await orgApi.update(form.id, form)
    ElMessage.success('更新成功')
  } else {
    await orgApi.create(form)
    ElMessage.success('创建成功')
  }
  dialogVisible.value = false
  loadTree()
}

const remove = async (row: any) => {
  await ElMessageBox.confirm('确认删除该组织？子组织将一并删除！', '提示', { type: 'warning' })
  await orgApi.remove(row.id)
  ElMessage.success('删除成功')
  loadTree()
}

onMounted(loadTree)
</script>
