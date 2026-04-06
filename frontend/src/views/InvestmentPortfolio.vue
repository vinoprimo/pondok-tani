<script setup lang="ts">
import { computed } from 'vue'
import { TrendingUp, DollarSign, Calendar, Download, Eye } from 'lucide-vue-next'

const portfolioData = [
  { month: 'Jan', value: 45000 },
  { month: 'Feb', value: 46800 },
  { month: 'Mar', value: 48200 },
  { month: 'Apr', value: 50100 },
  { month: 'May', value: 52500 },
  { month: 'Jun', value: 55980 },
]

const investments = [
  {
    id: 'INV-001',
    name: 'Premium Vanilla Batch A',
    invested: 15000,
    currentValue: 18650,
    roi: 24.3,
    plants: 30,
    status: 'active',
    startDate: '2025-01-15',
  },
  {
    id: 'INV-002',
    name: 'Organic Vanilla Batch B',
    invested: 20000,
    currentValue: 24800,
    roi: 24.0,
    plants: 40,
    status: 'active',
    startDate: '2025-03-20',
  },
  {
    id: 'INV-003',
    name: 'Standard Vanilla Batch C',
    invested: 10000,
    currentValue: 12530,
    roi: 25.3,
    plants: 15,
    status: 'active',
    startDate: '2025-06-10',
  },
]

const maxPortfolioValue = computed(() => Math.max(...portfolioData.map((d) => d.value), 1))
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-semibold text-gray-900">Investment Portfolio</h2>
        <p class="text-gray-600 mt-1">Track your vanilla plantation investments and returns</p>
      </div>
      <button
        type="button"
        class="flex items-center gap-2 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
      >
        <Download class="w-4 h-4" />
        Export Report
      </button>
    </div>

    <div class="bg-gradient-to-br from-green-600 to-green-700 rounded-xl p-6 text-white">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div>
          <p class="text-green-100 text-sm mb-1">Total Invested</p>
          <p class="text-3xl font-semibold">$45,000</p>
        </div>
        <div>
          <p class="text-green-100 text-sm mb-1">Current Value</p>
          <p class="text-3xl font-semibold">$55,980</p>
        </div>
        <div>
          <p class="text-green-100 text-sm mb-1">Total Returns</p>
          <p class="text-3xl font-semibold">$10,980</p>
        </div>
        <div>
          <p class="text-green-100 text-sm mb-1">Average ROI</p>
          <p class="text-3xl font-semibold">24.5%</p>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-6">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">Portfolio Value Growth</h3>
      <p class="text-xs text-gray-500 mb-3">Area chart placeholder — bar strip + table</p>
      <div class="flex items-end gap-1 h-32 mb-4 px-1 border border-gray-100 rounded-lg bg-gray-50/50 p-2">
        <div
          v-for="d in portfolioData"
          :key="d.month"
          class="flex-1 flex flex-col items-center justify-end gap-1 min-w-0"
        >
          <div
            class="w-full max-w-[48px] mx-auto rounded-t bg-emerald-500 min-h-[4px]"
            :style="{ height: `${(d.value / maxPortfolioValue) * 100}%` }"
            :title="`${d.month}: $${d.value.toLocaleString()}`"
          />
          <span class="text-[10px] text-gray-500 truncate w-full text-center">{{ d.month }}</span>
        </div>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead class="bg-gray-50 text-left text-gray-600">
            <tr>
              <th class="px-3 py-2 font-medium">Month</th>
              <th class="px-3 py-2 font-medium text-right">Portfolio value</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="row in portfolioData" :key="row.month">
              <td class="px-3 py-2 font-medium text-gray-900">{{ row.month }}</td>
              <td class="px-3 py-2 text-right text-gray-900">${{ row.value.toLocaleString() }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden">
      <div class="p-6 border-b border-gray-200">
        <h3 class="text-lg font-semibold text-gray-900">Active Investments</h3>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Investment
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Invested
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Current Value
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ROI</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Plants</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Start Date
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="investment in investments" :key="investment.id" class="hover:bg-gray-50">
              <td class="px-6 py-4">
                <div>
                  <p class="font-medium text-gray-900">{{ investment.name }}</p>
                  <p class="text-sm text-gray-500">{{ investment.id }}</p>
                </div>
              </td>
              <td class="px-6 py-4 text-gray-900">${{ investment.invested.toLocaleString() }}</td>
              <td class="px-6 py-4 text-gray-900">${{ investment.currentValue.toLocaleString() }}</td>
              <td class="px-6 py-4">
                <span
                  class="inline-flex items-center gap-1 px-2.5 py-1 bg-green-50 text-green-700 rounded-full text-sm font-medium"
                >
                  <TrendingUp class="w-3 h-3" />
                  {{ investment.roi }}%
                </span>
              </td>
              <td class="px-6 py-4 text-gray-900">{{ investment.plants }}</td>
              <td class="px-6 py-4 text-gray-600 text-sm">{{ investment.startDate }}</td>
              <td class="px-6 py-4">
                <button
                  type="button"
                  class="text-green-600 hover:text-green-700 font-medium text-sm flex items-center gap-1"
                >
                  <Eye class="w-4 h-4" />
                  View Details
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-4">
          <div class="w-10 h-10 bg-blue-50 rounded-lg flex items-center justify-center">
            <DollarSign class="w-5 h-5 text-blue-600" />
          </div>
          <h3 class="font-semibold text-gray-900">Last Month</h3>
        </div>
        <p class="text-2xl font-semibold text-gray-900 mb-1">$918</p>
        <p class="text-sm text-gray-600">June 2026 Returns</p>
      </div>
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-4">
          <div class="w-10 h-10 bg-green-50 rounded-lg flex items-center justify-center">
            <Calendar class="w-5 h-5 text-green-600" />
          </div>
          <h3 class="font-semibold text-gray-900">Next Payout</h3>
        </div>
        <p class="text-2xl font-semibold text-gray-900 mb-1">$950</p>
        <p class="text-sm text-gray-600">Expected on Aug 1, 2026</p>
      </div>
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-4">
          <div class="w-10 h-10 bg-purple-50 rounded-lg flex items-center justify-center">
            <TrendingUp class="w-5 h-5 text-purple-600" />
          </div>
          <h3 class="font-semibold text-gray-900">YTD Returns</h3>
        </div>
        <p class="text-2xl font-semibold text-gray-900 mb-1">$5,340</p>
        <p class="text-sm text-gray-600">Year to Date</p>
      </div>
    </div>
  </div>
</template>
