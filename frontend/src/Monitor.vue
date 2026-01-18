<template>
  <div v-if="monitor_addr >= 0" style="height: 300px; margin-bottom: 50px">
    <button style="float: right" @click="monitor_addr = -1">close</button>
    <h4>Monitor #{{ monitor_addr }}</h4>
    <LineChart :data="chartData" :options="chartOptions" updateMode="resize" />
  </div>
</template>

<script>

import { SessionDataStore } from './stores/session';
import { mapState, mapWritableState } from 'pinia';
import { Line } from "vue-chartjs";
import { Chart as ChartJS, Title, Tooltip, Legend, LineElement, PointElement, CategoryScale, TimeScale, LinearScale } from 'chart.js';

// Register necessary components with Chart.js
ChartJS.register(Title, Tooltip, Legend, LineElement, PointElement, CategoryScale, TimeScale, LinearScale);



export default {
  name: "GraphComponent",
  components: {
    LineChart: Line
  },
  setup() {
  },
  data(){
    return {
      chartOptions:{
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          x: {
            type: 'linear',
            title: {
              display: false,
              text: "Time passed"
            },
            ticks: {
              callback: value => {
                // Format time dynamically
                const seconds = value
                if (seconds > -60) {
                  return `${seconds}s`
                } else if (seconds > -3600) {
                  return `${(seconds / 60).toFixed(1)}m`
                } else {
                  return `${(seconds / 3600).toFixed(2)}h`
                }
              }
            }
          },
          y: {
            title: {
              display: true,
              text: "Value"
            },
            ticks: {
              stepSize: 1
            }
          }
        },
        plugins: {
          legend: {
            display: false
          }
        }
      },
      loading: true,
      timeseriesDataValues:[],
      timeseriesDataLabels:[]
    }
  },
  computed: {
    ...mapState(SessionDataStore, ["monitor_buffer"]),
    ...mapWritableState(SessionDataStore, ["monitor_addr"]),
    chartData() {
      const ts = new Date()
      if (this.monitor_buffer.length > 0) {
        this.timeseriesDataValues = this.monitor_buffer.map(item => item.data.value),
        this.timeseriesDataLabels = this.monitor_buffer.map(item => (item.ts-ts)/1000)
      } else {
        this.timeseriesDataValues = []
        this.timeseriesDataLabels = []
      }
      return {
        labels: this.timeseriesDataLabels,
        datasets: [
        {
            data: this.timeseriesDataValues, 
            fill: false, 
            borderColor: 'rgba(75, 192, 192, 1)', 
            tension: 0.1,
            stepped: true,
          }
        ],
      }
    },
  }
};
</script>

<style scoped>
/* Optional: You can add additional custom styles if needed */
</style>