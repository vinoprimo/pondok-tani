<script setup>
import {
  Chart as ChartJS,
  LineElement,
  CategoryScale,
  LinearScale,
  PointElement,
} from "chart.js";

import { Line } from "vue-chartjs";

ChartJS.register(LineElement, CategoryScale, LinearScale, PointElement);

const chartData = {
  labels: ["Jan", "Feb", "Mar", "Apr", "Mei", "Jun"],
  datasets: [
    {
      label: "Profit",
      data: [5, 10, 8, 15, 20, 25],
      borderColor: "#16a34a",
      tension: 0.4,
    },
  ],
};
</script>

<template>
  <div class="bg-white p-6 rounded-xl shadow-sm">
    <p class="font-semibold mb-4">Performa Investasi</p>

    <Line :data="chartData" />
  </div>
</template>