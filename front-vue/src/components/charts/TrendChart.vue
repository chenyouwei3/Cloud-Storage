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
    grid: {
      top: 0,
      right: 0,
      bottom: 0,
      left: 0
    },
    xAxis: {
      type: 'category',
      show: false,
      boundaryGap: false
    },
    yAxis: {
      type: 'value',
      show: false
    },
    series: [
      {
        type: 'line',
        data: props.data,
        symbol: 'none',
        lineStyle: {
          color: '#1890ff',
          width: 2
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            {
              offset: 0,
              color: 'rgba(24,144,255,0.3)'
            },
            {
              offset: 1,
              color: 'rgba(24,144,255,0.1)'
            }
          ])
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
  height: 50px;
}
</style> 