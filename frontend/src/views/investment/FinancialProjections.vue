<script setup lang="ts">
import { computed } from 'vue'
import { TrendingUp, TrendingDown, DollarSign, Calendar } from 'lucide-vue-next'

const projectionData = [
  { year: '2024', revenue: 0, costs: 45000, profit: -45000 },
  { year: '2025', revenue: 18000, costs: 52000, profit: -34000 },
  { year: '2026', revenue: 65000, costs: 48000, profit: 17000 },
  { year: '2027', revenue: 95000, costs: 52000, profit: 43000 },
  { year: '2028', revenue: 125000, costs: 55000, profit: 70000 },
  { year: '2029', revenue: 148000, costs: 58000, profit: 90000 },
]

const monthlyBreakdown = [
  { month: 'Jan', revenue: 8500, costs: 7200 },
  { month: 'Peb', revenue: 9200, costs: 7500 },
  { month: 'Mar', revenue: 10800, costs: 7800 },
  { month: 'Apr', revenue: 12500, costs: 8100 },
  { month: 'Mei', revenue: 11800, costs: 8300 },
  { month: 'Jun', revenue: 12200, costs: 8500 },
]

const costBreakdown = [
  { category: 'Tenaga kerja', amount: 18000, percentage: 38 },
  { category: 'Bahan & perlengkapan', amount: 12000, percentage: 25 },
  { category: 'Peralatan', amount: 8000, percentage: 17 },
  { category: 'Pemeliharaan', amount: 6000, percentage: 13 },
  { category: 'Lainnya', amount: 3500, percentage: 7 },
]

const maxProjectionAbs = computed(() => {
  let m = 1
  for (const row of projectionData) {
    m = Math.max(m, Math.abs(row.revenue), Math.abs(row.costs), Math.abs(row.profit))
  }
  return m
})

const maxMonthlyBar = computed(() =>
  Math.max(...monthlyBreakdown.map((d) => Math.max(d.revenue, d.costs)), 1)
)
</script>

<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-semibold text-gray-900">Proyeksi keuangan</h2>
      <p class="text-gray-600 mt-1">Prakiraan pendapatan jangka panjang dan analisis biaya</p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-10 h-10 bg-green-50 rounded-lg flex items-center justify-center">
            <TrendingUp class="w-5 h-5 text-green-600" />
          </div>
          <span class="text-sm text-gray-600">Pendapatan proyeksi (2026)</span>
        </div>
        <p class="text-2xl font-semibold text-gray-900">$65,000</p>
        <p class="text-sm text-green-600 mt-1">+261% dari 2025</p>
      </div>
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-10 h-10 bg-blue-50 rounded-lg flex items-center justify-center">
            <DollarSign class="w-5 h-5 text-blue-600" />
          </div>
          <span class="text-sm text-gray-600">Laba bersih (2026)</span>
        </div>
        <p class="text-2xl font-semibold text-gray-900">$17,000</p>
        <p class="text-sm text-green-600 mt-1">Titik impas tercapai</p>
      </div>
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-10 h-10 bg-purple-50 rounded-lg flex items-center justify-center">
            <TrendingDown class="w-5 h-5 text-purple-600" />
          </div>
          <span class="text-sm text-gray-600">Biaya operasional</span>
        </div>
        <p class="text-2xl font-semibold text-gray-900">$48,000</p>
        <p class="text-sm text-gray-600 mt-1">Per tahun (2026)</p>
      </div>
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-10 h-10 bg-yellow-50 rounded-lg flex items-center justify-center">
            <Calendar class="w-5 h-5 text-yellow-600" />
          </div>
          <span class="text-sm text-gray-600">Titik impas</span>
        </div>
        <p class="text-2xl font-semibold text-gray-900">Kuartal II 2026</p>
        <p class="text-sm text-gray-600 mt-1">Sesuai rencana</p>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-6">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">Proyeksi keuangan 5 tahun</h3>
      <p class="text-xs text-gray-500 mb-3">Grafik multi-garis (placeholder) — tabel + batang skala</p>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead class="bg-gray-50 text-left">
            <tr>
              <th class="px-3 py-2 font-medium text-gray-600">Tahun</th>
              <th class="px-3 py-2 font-medium text-gray-600 text-right">Pendapatan</th>
              <th class="px-3 py-2 font-medium text-gray-600 text-right">Biaya</th>
              <th class="px-3 py-2 font-medium text-gray-600 text-right">Laba bersih</th>
              <th class="px-3 py-2 font-medium text-gray-600 min-w-[180px]">Visual</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="row in projectionData" :key="row.year">
              <td class="px-3 py-2 font-medium text-gray-900">{{ row.year }}</td>
              <td class="px-3 py-2 text-right text-emerald-600">${{ row.revenue.toLocaleString() }}</td>
              <td class="px-3 py-2 text-right text-red-500">${{ row.costs.toLocaleString() }}</td>
              <td
                class="px-3 py-2 text-right font-medium"
                :class="row.profit >= 0 ? 'text-blue-600' : 'text-gray-700'"
              >
                ${{ row.profit.toLocaleString() }}
              </td>
              <td class="px-3 py-2">
                <div class="flex flex-col gap-1">
                  <div class="h-1.5 bg-gray-100 rounded-full overflow-hidden">
                    <div
                      class="h-full bg-emerald-500 rounded-full"
                      :style="{ width: `${(Math.abs(row.revenue) / maxProjectionAbs) * 100}%` }"
                    />
                  </div>
                  <div class="h-1.5 bg-gray-100 rounded-full overflow-hidden">
                    <div
                      class="h-full bg-red-500 rounded-full"
                      :style="{ width: `${(Math.abs(row.costs) / maxProjectionAbs) * 100}%` }"
                    />
                  </div>
                  <div class="h-1.5 bg-gray-100 rounded-full overflow-hidden">
                    <div
                      class="h-full rounded-full bg-blue-500"
                      :style="{ width: `${(Math.abs(row.profit) / maxProjectionAbs) * 100}%` }"
                    />
                  </div>
                </div>
                <p class="text-[10px] text-gray-400 mt-0.5">hijau pendapatan · merah biaya · biru laba</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Kinerja bulanan 2026</h3>
        <p class="text-xs text-gray-500 mb-3">Grafik batang grup (placeholder)</p>
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead class="bg-gray-50 text-left">
              <tr>
                <th class="px-3 py-2 font-medium text-gray-600">Bulan</th>
                <th class="px-3 py-2 font-medium text-gray-600 text-right">Pendapatan</th>
                <th class="px-3 py-2 font-medium text-gray-600 text-right">Biaya</th>
                <th class="px-3 py-2 font-medium text-gray-600">Batang</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-for="row in monthlyBreakdown" :key="row.month">
                <td class="px-3 py-2 font-medium">{{ row.month }}</td>
                <td class="px-3 py-2 text-right text-emerald-600">${{ row.revenue.toLocaleString() }}</td>
                <td class="px-3 py-2 text-right text-amber-600">${{ row.costs.toLocaleString() }}</td>
                <td class="px-3 py-2">
                  <div class="flex items-end gap-1 h-10">
                    <div
                      class="w-4 rounded-t bg-emerald-500"
                      :style="{ height: `${(row.revenue / maxMonthlyBar) * 100}%` }"
                      title="Pendapatan"
                    />
                    <div
                      class="w-4 rounded-t bg-amber-500"
                      :style="{ height: `${(row.costs / maxMonthlyBar) * 100}%` }"
                      title="Biaya"
                    />
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Rincian biaya tahunan (2026)</h3>
        <div class="space-y-4">
          <div v-for="(item, index) in costBreakdown" :key="index">
            <div class="flex items-center justify-between mb-2">
              <span class="text-sm font-medium text-gray-900">{{ item.category }}</span>
              <span class="text-sm text-gray-600"
                >${{ item.amount.toLocaleString() }} ({{ item.percentage }}%)</span
              >
            </div>
            <div class="w-full bg-gray-100 rounded-full h-3">
              <div
                class="bg-green-600 h-3 rounded-full transition-all"
                :style="{ width: `${item.percentage}%` }"
              />
            </div>
          </div>
        </div>
        <div class="mt-6 pt-6 border-t border-gray-200">
          <div class="flex items-center justify-between">
            <span class="font-semibold text-gray-900">Total biaya tahunan</span>
            <span class="text-xl font-semibold text-gray-900">$47,500</span>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-6">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">Proyeksi ROI menurut tier investasi</h3>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="border border-gray-200 rounded-lg p-6">
          <div class="mb-4">
            <h4 class="font-semibold text-gray-900">Tier dasar</h4>
            <p class="text-sm text-gray-600">$5.000 - $15.000</p>
          </div>
          <div class="space-y-3">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">ROI tahun ke-1</span>
              <span class="font-medium">-15%</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">ROI tahun ke-3</span>
              <span class="font-medium text-green-600">+18%</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">ROI tahun ke-5</span>
              <span class="font-medium text-green-600">+42%</span>
            </div>
          </div>
        </div>
        <div class="border-2 border-green-600 rounded-lg p-6 bg-green-50">
          <div class="mb-4">
            <div class="flex items-center gap-2">
              <h4 class="font-semibold text-gray-900">Tier premium</h4>
              <span class="px-2 py-0.5 bg-green-600 text-white text-xs rounded-full">Populer</span>
            </div>
            <p class="text-sm text-gray-600">$15.000 - $30.000</p>
          </div>
          <div class="space-y-3">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">ROI tahun ke-1</span>
              <span class="font-medium">-12%</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">ROI tahun ke-3</span>
              <span class="font-medium text-green-600">+22%</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">ROI tahun ke-5</span>
              <span class="font-medium text-green-600">+52%</span>
            </div>
          </div>
        </div>
        <div class="border border-gray-200 rounded-lg p-6">
          <div class="mb-4">
            <h4 class="font-semibold text-gray-900">Tier elite</h4>
            <p class="text-sm text-gray-600">$30.000+</p>
          </div>
          <div class="space-y-3">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">ROI tahun ke-1</span>
              <span class="font-medium">-10%</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">ROI tahun ke-3</span>
              <span class="font-medium text-green-600">+25%</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">ROI tahun ke-5</span>
              <span class="font-medium text-green-600">+58%</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-blue-50 border border-blue-200 rounded-xl p-6">
      <h4 class="font-semibold text-gray-900 mb-3">Asumsi proyeksi</h4>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm text-gray-700">
        <div class="flex items-start gap-2">
          <div class="w-1.5 h-1.5 bg-blue-600 rounded-full mt-1.5 shrink-0" />
          <span>Harga pasar vanili rata-rata: $450/kg</span>
        </div>
        <div class="flex items-start gap-2">
          <div class="w-1.5 h-1.5 bg-blue-600 rounded-full mt-1.5 shrink-0" />
          <span>Kematangan tanaman: 18–24 bulan</span>
        </div>
        <div class="flex items-start gap-2">
          <div class="w-1.5 h-1.5 bg-blue-600 rounded-full mt-1.5 shrink-0" />
          <span>Hasil rata-rata: 1,5 kg per tanaman per tahun</span>
        </div>
        <div class="flex items-start gap-2">
          <div class="w-1.5 h-1.5 bg-blue-600 rounded-full mt-1.5 shrink-0" />
          <span>Inflasi biaya tahunan: 3–5%</span>
        </div>
      </div>
    </div>
  </div>
</template>
