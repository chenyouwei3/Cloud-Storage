<template>
  <div class="visualization-page">
    <div class="header">
      <a-row justify="space-between" align="middle">
        <a-col>
          <h1 class="page-title">数据可视化</h1>
        </a-col>
        <a-col>
          <a-button type="primary" @click="goBack">
            <template #icon><arrow-left-outlined /></template>
            返回文件云盘
          </a-button>
        </a-col>
      </a-row>
    </div>

    <div class="content">
      <!-- 数据概览卡片 -->
      <a-row :gutter="16" class="data-overview">
        <a-col :span="6">
          <a-card>
            <template #title>总文件数</template>
            <div class="card-content">
              <h2>{{ mockData.totalFiles }}</h2>
              <trend-chart :data="mockData.fileTrend" />
            </div>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <template #title>总存储空间</template>
            <div class="card-content">
              <h2>{{ mockData.totalStorage }}</h2>
              <trend-chart :data="mockData.storageTrend" />
            </div>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <template #title>本月上传</template>
            <div class="card-content">
              <h2>{{ mockData.monthlyUploads }}</h2>
              <trend-chart :data="mockData.uploadTrend" />
            </div>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <template #title>本月下载</template>
            <div class="card-content">
              <h2>{{ mockData.monthlyDownloads }}</h2>
              <trend-chart :data="mockData.downloadTrend" />
            </div>
          </a-card>
        </a-col>
      </a-row>

      <!-- 图表区域 -->
      <a-row :gutter="16" class="charts-section">
        <a-col :span="12">
          <a-card title="文件类型分布">
            <pie-chart :data="mockData.fileTypeDistribution" />
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="存储空间使用趋势">
            <line-chart :data="mockData.storageUsageTrend" />
          </a-card>
        </a-col>
      </a-row>

      <a-row :gutter="16" class="charts-section">
        <a-col :span="24">
          <a-card title="每日文件操作统计">
            <bar-chart :data="mockData.dailyOperations" />
          </a-card>
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeftOutlined } from '@ant-design/icons-vue'
import PieChart from '@/components/charts/PieChart.vue'
import LineChart from '@/components/charts/LineChart.vue'
import BarChart from '@/components/charts/BarChart.vue'
import TrendChart from '@/components/charts/TrendChart.vue'

const router = useRouter()

// 模拟数据
const mockData = ref({
  totalFiles: '1,234',
  totalStorage: '45.6 GB',
  monthlyUploads: '328',
  monthlyDownloads: '156',
  fileTrend: [30, 40, 35, 50, 49, 60, 70],
  storageTrend: [20, 30, 40, 35, 45, 50, 60],
  uploadTrend: [10, 15, 20, 25, 30, 35, 40],
  downloadTrend: [5, 10, 15, 20, 25, 30, 35],
  fileTypeDistribution: [
    { type: '图片', value: 40 },
    { type: '文档', value: 30 },
    { type: '视频', value: 20 },
    { type: '其他', value: 10 }
  ],
  storageUsageTrend: [
    { date: '1月', value: 30 },
    { date: '2月', value: 35 },
    { date: '3月', value: 40 },
    { date: '4月', value: 45 },
    { date: '5月', value: 50 },
    { date: '6月', value: 55 }
  ],
  dailyOperations: [
    { date: '周一', uploads: 20, downloads: 10 },
    { date: '周二', uploads: 25, downloads: 15 },
    { date: '周三', uploads: 30, downloads: 20 },
    { date: '周四', uploads: 35, downloads: 25 },
    { date: '周五', uploads: 40, downloads: 30 },
    { date: '周六', uploads: 45, downloads: 35 },
    { date: '周日', uploads: 50, downloads: 40 }
  ]
})

// 返回文件云盘
const goBack = () => {
  router.push('/filecloud')
}
</script>

<style scoped>
.visualization-page {
  padding: 24px;
  background-color: #f0f2f5;
  min-height: 100vh;
}

.header {
  margin-bottom: 24px;
}

.page-title {
  margin: 0;
  font-size: 24px;
  font-weight: 500;
  color: #1f1f1f;
}

.content {
  background-color: #fff;
  padding: 24px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.data-overview {
  margin-bottom: 24px;
}

.card-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.card-content h2 {
  margin: 0;
  font-size: 28px;
  color: #1f1f1f;
}

.charts-section {
  margin-bottom: 24px;
}

:deep(.ant-card) {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

:deep(.ant-card-head) {
  border-bottom: 1px solid #f0f0f0;
  padding: 16px 24px;
}

:deep(.ant-card-head-title) {
  font-size: 16px;
  font-weight: 500;
  color: #1f1f1f;
}

:deep(.ant-card-body) {
  padding: 24px;
}
</style> 