<template>
  <div ref="chartRef" class="chart-container"></div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import * as echarts from 'echarts'

const props = defineProps({
  data: {
    type: Array,
    required: true
  }
})

const chartRef = ref(null)
let chart = null

const initChart = () => {
  if (chart) {
    chart.dispose()
  }
  chart = echarts.init(chartRef.value)
  updateChart()
}

const updateChart = () => {
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    legend: {
      data: ['上传', '下载']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: props.data.map(item => item.date)
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '上传',
        type: 'bar',
        data: props.data.map(item => item.uploads),
        itemStyle: {
          color: '#1890ff'
        }
      },
      {
        name: '下载',
        type: 'bar',
        data: props.data.map(item => item.downloads),
        itemStyle: {
          color: '#52c41a'
        }
      }
    ]
  }
  chart.setOption(option)
}

watch(() => props.data, () => {
  updateChart()
}, { deep: true })

onMounted(() => {
  initChart()
  window.addEventListener('resize', () => {
    chart?.resize()
  })
})
</script>

<style scoped>
.chart-container {
  width: 100%;
  height: 300px;
}
</style> 