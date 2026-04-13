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
  { month: 'Feb', revenue: 9200, costs: 7500 },
  { month: 'Mar', revenue: 10800, costs: 7800 },
  { month: 'Apr', revenue: 12500, costs: 8100 },
  { month: 'May', revenue: 11800, costs: 8300 },
  { month: 'Jun', revenue: 12200, costs: 8500 },
]

const costBreakdown = [
  { category: 'Labor', amount: 18000, percentage: 38 },
  { category: 'Materials & Supplies', amount: 12000, percentage: 25 },
  { category: 'Equipment', amount: 8000, percentage: 17 },
  { category: 'Maintenance', amount: 6000, percentage: 13 },
  { category: 'Other', amount: 3500, percentage: 7 },
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
      <h2 class="text-2xl font-semibold text-gray-900">Financial Projections</h2>
      <p class="text-gray-600 mt-1">Long-term revenue forecasts and cost analysis</p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-10 h-10 bg-green-50 rounded-lg flex items-center justify-center">
            <TrendingUp class="w-5 h-5 text-green-600" />
          </div>
          <span class="text-sm text-gray-600">Projected Revenue (2026)</span>
        </div>
        <p class="text-2xl font-semibold text-gray-900">$65,000</p>
        <p class="text-sm text-green-600 mt-1">+261% from 2025</p>
      </div>
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-10 h-10 bg-blue-50 rounded-lg flex items-center justify-center">
            <DollarSign class="w-5 h-5 text-blue-600" />
          </div>
          <span class="text-sm text-gray-600">Net Profit (2026)</span>
        </div>
        <p class="text-2xl font-semibold text-gray-900">$17,000</p>
        <p class="text-sm text-green-600 mt-1">Break-even achieved</p>
      </div>
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-10 h-10 bg-purple-50 rounded-lg flex items-center justify-center">
            <TrendingDown class="w-5 h-5 text-purple-600" />
          </div>
          <span class="text-sm text-gray-600">Operating Costs</span>
        </div>
        <p class="text-2xl font-semibold text-gray-900">$48,000</p>
        <p class="text-sm text-gray-600 mt-1">Annual (2026)</p>
      </div>
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-10 h-10 bg-yellow-50 rounded-lg flex items-center justify-center">
            <Calendar class="w-5 h-5 text-yellow-600" />
          </div>
          <span class="text-sm text-gray-600">Break-even Point</span>
        </div>
        <p class="text-2xl font-semibold text-gray-900">Q2 2026</p>
        <p class="text-sm text-gray-600 mt-1">On track</p>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-6">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">5-Year Financial Projection</h3>
      <p class="text-xs text-gray-500 mb-3">Multi-line chart placeholder — table + scaled bars</p>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead class="bg-gray-50 text-left">
            <tr>
              <th class="px-3 py-2 font-medium text-gray-600">Year</th>
              <th class="px-3 py-2 font-medium text-gray-600 text-right">Revenue</th>
              <th class="px-3 py-2 font-medium text-gray-600 text-right">Costs</th>
              <th class="px-3 py-2 font-medium text-gray-600 text-right">Net profit</th>
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
                <p class="text-[10px] text-gray-400 mt-0.5">green revenue · red costs · blue profit</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">2026 Monthly Performance</h3>
        <p class="text-xs text-gray-500 mb-3">Grouped bar chart placeholder</p>
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead class="bg-gray-50 text-left">
              <tr>
                <th class="px-3 py-2 font-medium text-gray-600">Month</th>
                <th class="px-3 py-2 font-medium text-gray-600 text-right">Revenue</th>
                <th class="px-3 py-2 font-medium text-gray-600 text-right">Costs</th>
                <th class="px-3 py-2 font-medium text-gray-600">Bars</th>
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
                      title="Revenue"
                    />
                    <div
                      class="w-4 rounded-t bg-amber-500"
                      :style="{ height: `${(row.costs / maxMonthlyBar) * 100}%` }"
                      title="Costs"
                    />
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Annual Cost Breakdown (2026)</h3>
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
            <span class="font-semibold text-gray-900">Total Annual Costs</span>
            <span class="text-xl font-semibold text-gray-900">$47,500</span>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-6">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">ROI Projection by Investment Tier</h3>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="border border-gray-200 rounded-lg p-6">
          <div class="mb-4">
            <h4 class="font-semibold text-gray-900">Basic Tier</h4>
            <p class="text-sm text-gray-600">$5,000 - $15,000</p>
          </div>
          <div class="space-y-3">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Year 1 ROI</span>
              <span class="font-medium">-15%</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Year 3 ROI</span>
              <span class="font-medium text-green-600">+18%</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Year 5 ROI</span>
              <span class="font-medium text-green-600">+42%</span>
            </div>
          </div>
        </div>
        <div class="border-2 border-green-600 rounded-lg p-6 bg-green-50">
          <div class="mb-4">
            <div class="flex items-center gap-2">
              <h4 class="font-semibold text-gray-900">Premium Tier</h4>
              <span class="px-2 py-0.5 bg-green-600 text-white text-xs rounded-full">Popular</span>
            </div>
            <p class="text-sm text-gray-600">$15,000 - $30,000</p>
          </div>
          <div class="space-y-3">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Year 1 ROI</span>
              <span class="font-medium">-12%</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Year 3 ROI</span>
              <span class="font-medium text-green-600">+22%</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Year 5 ROI</span>
              <span class="font-medium text-green-600">+52%</span>
            </div>
          </div>
        </div>
        <div class="border border-gray-200 rounded-lg p-6">
          <div class="mb-4">
            <h4 class="font-semibold text-gray-900">Elite Tier</h4>
            <p class="text-sm text-gray-600">$30,000+</p>
          </div>
          <div class="space-y-3">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Year 1 ROI</span>
              <span class="font-medium">-10%</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Year 3 ROI</span>
              <span class="font-medium text-green-600">+25%</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Year 5 ROI</span>
              <span class="font-medium text-green-600">+58%</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-blue-50 border border-blue-200 rounded-xl p-6">
      <h4 class="font-semibold text-gray-900 mb-3">Projection Assumptions</h4>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm text-gray-700">
        <div class="flex items-start gap-2">
          <div class="w-1.5 h-1.5 bg-blue-600 rounded-full mt-1.5 shrink-0" />
          <span>Average vanilla market price: $450/kg</span>
        </div>
        <div class="flex items-start gap-2">
          <div class="w-1.5 h-1.5 bg-blue-600 rounded-full mt-1.5 shrink-0" />
          <span>Plant maturity: 18-24 months</span>
        </div>
        <div class="flex items-start gap-2">
          <div class="w-1.5 h-1.5 bg-blue-600 rounded-full mt-1.5 shrink-0" />
          <span>Average yield: 1.5kg per plant annually</span>
        </div>
        <div class="flex items-start gap-2">
          <div class="w-1.5 h-1.5 bg-blue-600 rounded-full mt-1.5 shrink-0" />
          <span>Annual cost inflation: 3-5%</span>
        </div>
      </div>
    </div>
  </div>
</template>
