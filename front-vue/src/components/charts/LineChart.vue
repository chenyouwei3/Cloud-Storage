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
      trigger: 'axis'
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: props.data.map(item => item.date)
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '存储空间',
        type: 'line',
        smooth: true,
        data: props.data.map(item => item.value),
        areaStyle: {
          opacity: 0.3
        },
        lineStyle: {
          width: 3
        },
        itemStyle: {
          borderWidth: 2
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