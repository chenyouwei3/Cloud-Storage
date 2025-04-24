<!-- 文件云盘主界面模板 -->
<template>
  <!-- 主体容器 -->
  <div class="body">
    <!-- 顶部操作栏 -->
    <div id="header">
      <a-row justify="space-between" type="flex">
        <!-- 左侧操作按钮组 -->
        <a-col>
          <a-button-group>
            <!-- 刷新按钮 -->
            <a-button class="action-btn" @click="refreshList">
              <template #icon><sync-outlined /></template>
              刷新
            </a-button>
            <!-- 创建文件夹按钮 -->
            <a-button class="action-btn" @click="handleCreateFolder">
              <template #icon><folder-add-outlined /></template>
              新建文件夹
            </a-button>
            <!-- 上传文件按钮 -->
            <a-upload
              :show-upload-list="false"
              :before-upload="handleUpload"
              :customRequest="customUpload"
            >
              <a-button class="action-btn">
                <template #icon><upload-outlined /></template>
                上传文件
              </a-button>
            </a-upload>
            <!-- 批量删除按钮 -->
            <a-button class="action-btn" @click="handleBatchDelete">
              <template #icon><delete-outlined /></template>
              批量删除
            </a-button>
            <!-- 可视化按钮 -->
            <a-button class="action-btn" @click="goToVisualization">
              <template #icon><bar-chart-outlined /></template>
              数据可视化
            </a-button>
          </a-button-group>
        </a-col>
        <!-- 右侧磁盘使用情况显示 -->
        <a-col style="margin-right:20px">
          <div class="disk-info">
            <div class="disk-text">
              {{ resData.totalDist || '0 B' }} / 50.0 MB
            </div>
            <div class="disk-progress">
              <a-progress 
                :percent="diskUsagePercent" 
                :stroke-color="progressColor(diskUsagePercent)" 
                :show-info="false"
              />
            </div>
          </div>
        </a-col>
        <!-- 用户信息和退出按钮 -->
        <a-col>
          <div class="user-info">
            <span class="username">{{ userInfo.name }}</span>
            <a-button type="link" @click="handleLogout">
              <template #icon><logout-outlined /></template>
              退出
            </a-button>
          </div>
        </a-col>
      </a-row>
    </div>

    <!-- 路径导航栏 -->
    <div class="path">
      <!-- 根目录显示 -->
      <template v-if="currentPath === ''">
        <div class="path-item current">全部文件</div>
      </template>
      <!-- 子目录显示 -->
      <template v-else>
        <div class="path-item" @click="navigateToParent">
          <left-outlined /> 返回上一级
        </div>
        <div class="path-item" @click="navigateToRoot">全部文件</div>
        <!-- 显示当前路径的各个层级 -->
        <template v-for="(part, index) in pathParts" :key="index">
          <div class="path-separator">></div>
          <div 
            :class="['path-item', index === pathParts.length - 1 ? 'current' : '']" 
            @click="navigateToPath(index)"
          >
            {{ part }}
          </div>
        </template>
      </template>
    </div>

    <!-- 文件列表表格 -->
    <a-table 
      :columns="columns" 
      :data-source="currentItems"
      :pagination="false"
      :rowKey="record => record.filename"
      :rowSelection="{
        selectedRowKeys: selectedRowKeys,
        onChange: onSelectChange,
        preserveSelectedRowKeys: true
      }"
      class="file-table"
    >
      <!-- 自定义表格单元格内容 -->
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'name'">
          <span>
            <!-- 根据类型显示文件或文件夹图标 -->
            <component :is="record.isDir ? 'folder-outlined' : 'file-outlined'" />
            &nbsp;&nbsp;
            <!-- 文件夹可点击进入 -->
            <template v-if="record.isDir">
              <a href="javascript:;" @click="enterFolder(record.filename)">{{ record.filename }}</a>
            </template>
            <!-- 文件只显示名称 -->
            <template v-else>
              <span>{{ record.filename }}</span>
            </template>
          </span>
        </template>
        <!-- -------------------------------------------- -->
        <template v-if="column.key === 'action'">
          <a-space>
            <!-- 只有文件才显示下载按钮 -->
            <template v-if="!record.isDir">
              <a-tooltip title="下载">
                <a-button type="text" @click="handleDownload(record)">
                  <template #icon><download-outlined /></template>
                </a-button>
              </a-tooltip>
            </template>
            <a-tooltip title="删除">
              <a-button type="text" @click="handleDelete(record)">
                <template #icon><delete-outlined /></template>
              </a-button>
            </a-tooltip>
            <a-tooltip title="重命名">
              <a-button type="text" @click="openRenameModal(record)">
                <template #icon><edit-outlined /></template>
              </a-button>
            </a-tooltip>
            <a-tooltip title="复制">
              <a-button type="text" @click="handleCopy(record)">
                <template #icon><copy-outlined /></template>
              </a-button>
            </a-tooltip>
            <a-tooltip title="移动">
              <a-button type="text" @click="handleMove(record)">
                <template #icon><swap-outlined /></template>
              </a-button>
            </a-tooltip>
          </a-space>
        </template>
      </template>
    </a-table>
  </div>

  <!-- 重命名对话框 -->
  <a-modal
  v-model:open="renameModalVisible"
  title="重命名"
  ok-text="确认重命名"
  cancel-text="放弃修改"
  :ok-button-props="{ type: 'primary', style: { backgroundColor: '#00AEEC', borderColor: '#00AEEC' } }"
  :cancel-button-props="{ type: 'default', style: { color: '#999' } }"
  @ok="handleRename"
  wrap-class-name="bili-modal-style"
>
  <a-form
    :model="renameForm"
    :label-col="{ span: 6 }"
    :wrapper-col="{ span: 18 }"
  >
    <a-form-item label="旧名称">
      <a-input v-model:value="renameForm.oldName" disabled class="bili-input" />
    </a-form-item>
    <a-form-item label="新名称">
      <a-input v-model:value="renameForm.newName" class="bili-input" />
    </a-form-item>
  </a-form>
  </a-modal>

  <!-- 新建文件夹对话框 -->
  <a-modal
    v-model:open="createFolderModalVisible"
    title="新建文件夹"
    ok-text="确认创建"
    cancel-text="放弃创建"
    :ok-button-props="{ type: 'primary', style: { backgroundColor: '#00AEEC', borderColor: '#00AEEC' } }"
    :cancel-button-props="{ type: 'default', style: { color: '#999' } }"
    @ok="confirmCreateFolder"
    wrap-class-name="bili-modal-style"
  >
    <a-form
      :model="createFolderForm"
      :label-col="{ span: 6 }"
      :wrapper-col="{ span: 18 }"
    >
      <a-form-item label="文件夹名称">
        <a-input v-model:value="createFolderForm.folderName" class="bili-input" />
      </a-form-item>
    </a-form>
  </a-modal>

  <!-- 移动文件对话框 -->
  <a-modal
    v-model:open="moveModalVisible"
    title="移动文件"
    ok-text="确认移动"
    cancel-text="放弃移动"
    :ok-button-props="{ type: 'primary', style: { backgroundColor: '#00AEEC', borderColor: '#00AEEC' } }"
    :cancel-button-props="{ type: 'default', style: { color: '#999' } }"
    @ok="confirmMove"
    wrap-class-name="bili-modal-style"
  >
    <a-form
      :model="moveForm"
      :label-col="{ span: 6 }"
      :wrapper-col="{ span: 18 }"
    >
      <a-form-item label="目标文件夹">
        <a-select
          v-model:value="moveForm.targetPath"
          :options="folderOptions"
          @change="onFolderChange"
        />
      </a-form-item>
    </a-form>
  </a-modal>

</template>

<script setup>
// 导入Vue相关功能
import { ref, onMounted, computed } from 'vue'
// 导入图标组件
import { 
  FolderOutlined, 
  FileOutlined, 
  SyncOutlined,
  LeftOutlined,
  DeleteOutlined,
  EditOutlined,
  CopyOutlined,
  SwapOutlined,
  FolderAddOutlined,
  UploadOutlined,
  DownloadOutlined,
  LogoutOutlined,
  BarChartOutlined
} from '@ant-design/icons-vue'
// 导入API函数
import { distList,distMkdir,distRename,distRemove,distCopy,distMove, DropdownMenu,distDownload,distUpload} from '@/tools/api'
// 导入消息提示组件
import { message } from 'ant-design-vue'
import { useRouter } from 'vue-router'
import storage from 'store'

/*----------------------------------------定义结构体----------------------------------------*/
// 定义表格列配置
const columns = [
  { 
    title: '文件名', 
    dataIndex: 'filename', 
    key: 'name',
    width: '55%'
  },
  { 
    title: '大小', 
    dataIndex: 'size', 
    key: 'size',
    width: '20%'
  },
  { 
    title: '修改时间', 
    dataIndex: 'date',
    key: 'date'
  },
  {
    title: '操作',
    key: 'action',
    width: '15%'
  }
]

// 组件挂载时初始化
onMounted(() => {
  getList()
})

const currentPath = ref('') // 当前路径
const selectedRowKeys = ref([]) // 选中的行
const renameModalVisible = ref(false) // 重命名对话框可见性
const createFolderModalVisible = ref(false) // 新建文件夹对话框可见性
const moveModalVisible = ref(false) // 移动文件对话框可见性
const folderOptions = ref([]) // 文件夹选项列表
const renameForm = ref({ // 重命名表单数据
  oldName: '',
  newName: '',
  record: null
})
const createFolderForm = ref({ // 新建文件夹表单数据
  folderName: ''
})
const moveForm = ref({ // 移动文件表单数据
  targetPath: '',
  record: null
})
const resData = ref({ // 文件列表数据
  defaultPath: "", // 默认路径
  totalDist: "0 B", // 已用空间
  items: {}, // 文件列表
  tree: null // 目录树结构
})



// 根据使用比例返回进度条颜色
const progressColor = (percent) => {
  if (percent >= 80) {
    return 'red'
  } else if (percent >= 50) {
    return '#EAC100'
  }
  return '#1890ff'
}
const globalPath="../cloud"

/* ---------------------------------------------------------------------------------------*/
// 处理行选择变化
const onSelectChange = (selectedKeys) => {
  selectedRowKeys.value = selectedKeys
}

/*-----------------------------------退出----------------------------------------------*/

// 获取用户信息
const userInfo = ref({
  name: storage.get('User-Info')?.name || '未登录',
  account: storage.get('User-Info')?.account || ''
})

// 处理退出登录
const router = useRouter()
const handleLogout = () => {
  storage.remove('Access-Token')
  storage.remove('User-Info')
  message.success('退出成功')
  router.push('/cloud_storage/login')
}

/*-----------------------------------刷新/数据可视化按钮----------------------------------------------*/

// 刷新列表
const refreshList = () => {
  getList(currentPath.value)
}
// 跳转到可视化页面
const goToVisualization = () => {
  router.push('/visualization')
}

/*--------------------------------- 下载文件 ---------------------------------*/
// 下载文件
const handleDownload = async (record) => {
  if (record.isDir) {
    message.warning('文件夹暂不支持下载')
    return
  }
  
  try {
    const filePath = currentPath.value ? `${globalPath}/${currentPath.value}/${record.filename}` : `${globalPath}/${record.filename}`
    
    const res = await distDownload({ path: filePath })
    
    if (res?.data?.code === 200) {
      // 创建一个临时的 a 标签来下载文件
      const link = document.createElement('a')
      link.href = res.data.data
      link.download = record.filename
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      message.success('下载成功')
    } else {
      message.error(res?.data?.message?.['zh-CN'] || '下载失败')
    }
  } catch (error) {
    console.error('下载失败:', error)
    message.error('下载失败，请稍后重试')
  }
}

/*--------------------------------- 上传文件 ---------------------------------*/
// 处理文件上传前的验证
const handleUpload = (file) => {
  // 这里可以添加文件大小、类型等验证
  return true
}

// 自定义上传请求
const customUpload = async ({ file, onSuccess, onError }) => {
  try {
 
    const path = currentPath.value ? `${globalPath}/${currentPath.value}` : globalPath
    console.log("test",path)
    const res = await distUpload(file, path)
    
    if (res?.data?.code === 200) {
      message.success('上传成功')
      onSuccess()
      getList(currentPath.value) // 刷新文件列表
    } else {
      message.error(res?.data?.message?.['zh-CN'] || '上传失败')
      onError()
    }
  } catch (error) {
    console.error('上传失败:', error)
    message.error('上传失败，请稍后重试')
    onError()
  }
}

/* ------------------------------------移动文件---------------------------------------------------*/
const handleMove = async (record) => {
  // 获取文件夹列表
  try {
    const res = await DropdownMenu({  path:globalPath,query:"move"})
    if (res?.data?.code === 200 && Array.isArray(res.data.data)) {
      // 格式化选项数据
      folderOptions.value = res.data.data.map(folder => ({
        value: folder,
        label: folder
      }))
    } else {
      folderOptions.value = []
      message.warning('获取文件夹列表失败')
    }
  } catch (error) {
    console.error('获取文件夹列表失败:', error)
    folderOptions.value = []  
  } 
  
  // 打开移动文件对话框
  moveForm.value = {
    targetPath: currentPath.value || '',
    record: record
  }
  moveModalVisible.value = true
}

// 处理文件夹选择变化
const onFolderChange = (value) => {
  moveForm.value.targetPath = value
}

// 确认移动文件
const confirmMove = async () => {
  const { targetPath, record } = moveForm.value
  if (!targetPath) {
    message.warning('请选择目标文件夹')
    return
  }
  
  try {
    // 构造源文件路径
    const sourcePath = currentPath.value ? `${globalPath}/${currentPath.value}/${record.filename}`: `${globalPath}/${record.filename}`
    // 构造目标路径
    const newPath = targetPath ? `${globalPath}/${targetPath}/${record.filename}`: `${globalPath}/${record.filename}`
    
    // 调用移动文件接口
    const res = await distMove({
      oldPath: sourcePath,
      newPath: newPath
    })
    
    if (res?.data?.code === 200) {
      message.success('移动成功')
      getList(currentPath.value) // 刷新文件列表
    } else {
      message.error(res?.data?.message?.['zh-CN'] || '移动失败')
    }
  } catch (error) {
    console.error('移动失败:', error)
    message.error('移动失败，请稍后重试')
  } finally {
    moveModalVisible.value = false
  }
}

/* ------------------------------------复制文件---------------------------------------------------*/
// 复制文件
const handleCopy = async (record) => {
  try {
    // 这里实现移动或复制逻辑
    const filePath = currentPath.value ? `${globalPath}/${currentPath.value}/${record.filename}`: `${globalPath}/${record.filename}`
    const res = await distCopy({  path:filePath} )   // 调用复制文件接口
    if (res?.data?.code==200){
      message.success('复制成功')
      getList(currentPath.value) // 刷新文件列表
    }else{
      message.error(res?.data?.message?.['zh-CN'] || '删除失败')
    }
  } catch (error) {
    console.error('复制失败:', error)
    message.error('复制失败，请稍后重试')
  }
}

/*--------------------------------- 删除/批量删除 ---------------------------------*/
// 删除文件或文件夹
const handleDelete = async (record) => {
  try {
    // 构造要删除的文件路径
    const filePath = currentPath.value ? `${globalPath}/${currentPath.value}/${record.filename}`: `${globalPath}/${record.filename}`
    const res = await distRemove({  distsPath: [filePath] })   // 调用删除接口
    if (res?.data?.code === 200) {
      message.success('删除成功')
      getList(currentPath.value) // 刷新文件列表
    } else {
      message.error(res?.data?.message?.['zh-CN'] || '删除失败')
    }
  } catch (error) {
    console.error('删除失败:', error)
    message.error('删除失败，请稍后重试')
  }
}

// 批量删除按钮
const handleBatchDelete = async () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请先选择要删除的文件')
    return
  }
  try {
    // 构造要删除的文件路径数组
    const pathsToDelete = selectedRowKeys.value.map(filename => {
      return currentPath.value ? `${globalPath}/${currentPath.value}/${filename}`: `${globalPath}/${filename}`
    })

    // 调用删除接口
    const res = await distRemove({distsPath: pathsToDelete})

    if (res?.data?.code === 200) {
      message.success('批量删除成功')
      selectedRowKeys.value = [] // 清空选中状态
      getList(currentPath.value) // 刷新文件列表
    } else {
      message.error(res?.data?.message?.['zh-CN'] || '批量删除失败')
    }
  } catch (error) {
    console.error('批量删除失败:', error)
    message.error('批量删除失败，请稍后重试')
  }
}

/*--------------------------------- 重命名 ---------------------------------*/
const handleRename = async () => {
  const { oldName, newName, record } = renameForm.value

  // 简单校验
  if (!newName || newName === oldName) {
    renameModalVisible.value = false
    return
  }

  try {
    // 构造旧路径
    const oldPath = currentPath.value?`${globalPath}/${currentPath.value}/${record.filename}`:`${globalPath}/${record.filename}`
    // 调用重命名接口
    const res = await distRename({oldPath: oldPath,newPath: newName})
    if (res?.data?.code === 200) {
      message.success('重命名成功')
      getList(currentPath.value) // 刷新文件列表
    } else {
      message.error(res?.data?.message?.['zh-CN'] || '重命名失败')
    }
  } catch (err) {
    console.error('重命名失败:', err)
    message.error('重命名失败，请稍后重试')
  } finally {
    renameModalVisible.value = false
  }
}

// 打开重命名对话框
const openRenameModal = (record) => {
  renameForm.value = {
    oldName: record.filename,
    newName: record.filename,
    record: record
  }
  renameModalVisible.value = true
}

/*--------------------------------- 新创建文件夹 ---------------------------------*/
// 创建文件夹
const handleCreateFolder = () => {
  createFolderForm.value.folderName = ''
  createFolderModalVisible.value = true
}

// 确认创建文件夹
const confirmCreateFolder = async () => {
  const { folderName } = createFolderForm.value
  if (!folderName) {
    message.warning('请输入文件夹名称')
    return
  }
  
  try {
    const fullPath = currentPath.value ? `${globalPath}/${currentPath.value}/${folderName}`: `${globalPath}/${folderName}`
    const res = await distMkdir({ path: fullPath})
    if (res?.data?.code === 200) {
      message.success('创建成功')
      getList(currentPath.value)
    } else {
      message.error(res?.data?.message?.['zh-CN'] || '创建失败')
    }
  } catch (error) {
    console.error('创建失败:', error)
    message.error('创建失败，请稍后重试')
  } finally {
    createFolderModalVisible.value = false
  }
}

/*--------------------------------- 数据渲染 ---------------------------------*/
// 获取文件列表
const getList = async () => {
  try {
    const res = await distList({ path: globalPath })  
    if (res && res.data && res.data.code === 200 && res.data.data) {
      resData.value = {
        defaultPath: res.data.data.defaultPath || "",
        totalDist: res.data.data.totalDist || "0 B",
        items: res.data.data.item || {},
        tree: res.data.data.tree || { name: 'root', children: [] }
      }
    } else {
      console.error('API 返回数据格式错误:', res)
      resData.value = {
        defaultPath: "",
        totalDist: "0 B",
        items: {},
        tree: { name: 'root', children: [] }
      }
    }
  } catch (error) {
    console.error('获取文件列表失败:', error)
    resData.value = {
      defaultPath: "",
      totalDist: "0 B",
      items: {},
      tree: { name: 'root', children: [] }
    }
  }
}

// 计算当前路径的各个部分
const pathParts = computed(() => {
  return currentPath.value ? currentPath.value.split('/') : []
})

// 计算磁盘使用百分比
const diskUsagePercent = computed(() => {
  const totalDist = resData.value.totalDist || '0 B'
  const match = totalDist.match(/(\d+(\.\d+)?)\s*(B|KB|MB|GB)/i)
  if (!match) return 0
  
  const value = parseFloat(match[1])
  const unit = match[3].toUpperCase()
  
  let bytes = 0
  if (unit === 'B') bytes = value
  else if (unit === 'KB') bytes = value * 1024
  else if (unit === 'MB') bytes = value * 1024 * 1024
  else if (unit === 'GB') bytes = value * 1024 * 1024 * 1024
  
  // 假设总空间为50MB
  const totalBytes = 50 * 1024 * 1024
  return Math.min(100, Math.round((bytes / totalBytes) * 100))
})


// 进入文件夹
const enterFolder = (folderName) => {
  const newPath = currentPath.value ? `${currentPath.value}/${folderName}` : folderName
  currentPath.value = newPath
  getList(newPath)
}

// 返回上一级
const navigateToParent = () => {
  if (currentPath.value) {
    const parts = currentPath.value.split('/')
    parts.pop()
    currentPath.value = parts.join('/')
    getList(currentPath.value)
  }
}

// 返回根目录
const navigateToRoot = () => {
  currentPath.value = ''
  getList()
}

// 跳转到指定路径
const navigateToPath = (index) => {
  const parts = pathParts.value.slice(0, index + 1)
  currentPath.value = parts.join('/')
  getList(currentPath.value)
}

// 计算当前目录下的文件和文件夹
const currentItems = computed(() => {
  if (!resData.value.tree) return []  //如果树状结构还没加载出来（例如数据请求还没回来），就直接返回空数组。 
  let currentNode = resData.value.tree // 从根节点开始
  if (!currentPath.value) {// 如果是根目录，直接返回  的子节点
    if (!currentNode.children) return []//如果当前路径为空（说明用户在根目录），就执行下面的逻辑。
    const cloudNode = currentNode.children.find(node => node.name === 'cloud')
    if (!cloudNode || !cloudNode.children) return []
    const innerCloudNode = cloudNode.children.find(node => node.name === 'cloud')  // 找到 cloud 的子节点中的 cloud 节点
    if (!innerCloudNode || !innerCloudNode.children) return []
    
    // 返回 cloud 的子节点
    return (innerCloudNode.children || []).map(item => {
      const isDir = !!item.children || resData.value.items[`cloud/${item.name}`]?.isDir
      const itemKey = `cloud/${item.name}`
      const itemInfo = resData.value.items[itemKey] || {}
      
      return {
        filename: item.name,
        isDir,
        size: itemInfo.size || '0B',
        date: itemInfo.date || '',
        displayName: item.name
      }
    }).sort((a, b) => {
      if (a.isDir && !b.isDir) return -1
      if (!a.isDir && b.isDir) return 1
      return a.displayName.localeCompare(b.displayName)
    })
  }
  
  // 如果不是根目录，找到当前目录节点
  const pathParts = currentPath.value.split('/')
  const fullPathParts = ['cloud', 'cloud', ...pathParts]
  let currentPathStr = ''
  
  for (const part of fullPathParts) {
    if (!currentNode.children) return []
    
    currentPathStr = currentPathStr ? `${currentPathStr}/${part}` : part
    const child = currentNode.children.find(c => c.name === part)
    if (!child) return []
    currentNode = child
  }
  
  // 返回当前目录的子节点，即使为空也返回空数组
  return (currentNode.children || []).map(item => {
    const isDir = !!item.children || resData.value.items[`${currentPathStr}/${item.name}`]?.isDir
    const itemKey = `${currentPathStr}/${item.name}`
    const itemInfo = resData.value.items[itemKey] || {}
    
    // 调试输出
    console.log('文件信息:', {
      name: item.name,
      isDir,
      itemKey,
      itemInfo
    })
    
    return {
      filename: item.name,
      isDir,
      size: itemInfo.size || '0B',
      date: itemInfo.date || '',
      displayName: item.name
    }
  }).sort((a, b) => {
    if (a.isDir && !b.isDir) return -1
    if (!a.isDir && b.isDir) return 1
    return a.displayName.localeCompare(b.displayName)
  })
})


</script>

<style scoped>
/* 主体容器样式 */
.body {
  padding: 20px;
  background-color: #f4f5f6;
  min-height: 100vh;
}

/* 顶部操作栏样式 */
#header {
  background-color: #fff;
  padding: 15px 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  margin-bottom: 20px;
}

.action-btn {
  background-color: #f4f5f6;
  border: none;
  color: #18191c;
}

.action-btn:hover {
  background-color: #e3e5e7;
  color: #18191c;
}

.disk-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.disk-text {
  color: #18191c;
  font-size: 14px;
}

.disk-progress {
  width: 200px;
}

/* 路径导航栏样式 */
.path {
  background-color: #fff;
  padding: 15px 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.path-item {
  color: #61666d;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.3s;
}

.path-item:hover {
  background-color: #f4f5f6;
  color: #18191c;
}

.path-item.current {
  color: #18191c;
  font-weight: 500;
}

.path-separator {
  color: #9499a0;
}

/* 文件表格样式 */
.file-table {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.file-table :deep(.ant-table-thead > tr > th) {
  background-color: #fafafa;
  color: #18191c;
  font-weight: 500;
}

.file-table :deep(.ant-table-tbody > tr:hover > td) {
  background-color: #f4f5f6;
}

.file-table :deep(.ant-table-tbody > tr > td) {
  border-bottom: 1px solid #f0f0f0;
}

/* 自定义滚动条样式 */
::-webkit-scrollbar {
  width: 7px;
  height: 10px;
}

::-webkit-scrollbar-thumb {
  background-color: #a1a3a9;
  border-radius: 3px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.username {
  color: #18191c;
  font-size: 14px;
  font-weight: 500;
}
</style> 