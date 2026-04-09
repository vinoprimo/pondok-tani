<script setup>
import DashboardLayout from "../layouts/DashboardLayout.vue";
import StatCard from "../components/StatCard.vue";
import InvestmentChart from "../components/InvestmentChart.vue";
</script>

<template>
  <DashboardLayout>

    <h1 class="text-2xl font-semibold mb-6">Dashboard</h1>

    <!-- STAT -->
    <div class="grid grid-cols-4 gap-6 mb-6">
      <StatCard title="Total Investasi" value="Rp 25.000.000" change="+12%" />
      <StatCard title="ROI" value="32.5%" change="+2.1%" />
      <StatCard title="Tanaman" value="120" change="+10" />
      <StatCard title="Keuntungan" value="Rp 8.200.000" change="+15%" />
    </div>

    <!-- CHART -->
    <div class="grid grid-cols-3 gap-6">

      <div class="col-span-2">
        <InvestmentChart />
      </div>

      <div class="bg-white p-6 rounded-xl shadow-sm">
        <p class="font-semibold mb-4">Distribusi Tanaman</p>

        <div class="h-60 flex items-center justify-center text-gray-400">
          Pie Chart (next)
        </div>
      </div>

    </div>

  </DashboardLayout>
</template>