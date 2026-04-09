<script setup lang="ts">
import { computed, ref } from 'vue';
import { Download, FileText, Calendar, Filter, TrendingUp } from 'lucide-vue-next';

const selectedPeriod = ref('monthly');
const selectedReport = ref('all');

const reportTypes = [
  { id: 'all', label: 'All Reports' },
  { id: 'financial', label: 'Financial' },
  { id: 'plant', label: 'Plant Health' },
  { id: 'investment', label: 'Investment' },
];

const reports = [
  {
    id: 1,
    name: 'Monthly Investment Statement',
    type: 'Investment',
    date: '2026-01-01',
    period: 'January 2026',
    size: '245 KB',
    format: 'PDF',
  },
  {
    id: 2,
    name: 'Q2 Financial Summary',
    type: 'Financial',
    date: '2025-07-01',
    period: 'Q2 2025',
    size: '512 KB',
    format: 'PDF',
  },
  {
    id: 3,
    name: 'Plant Health & Growth Report',
    type: 'Plant Health',
    date: '2026-01-15',
    period: 'January 2026',
    size: '1.2 MB',
    format: 'PDF',
  },
  {
    id: 4,
    name: 'ROI Analysis Report',
    type: 'Investment',
    date: '2025-12-31',
    period: '2025 Annual',
    size: '380 KB',
    format: 'PDF',
  },
  {
    id: 5,
    name: 'Cost Breakdown Analysis',
    type: 'Financial',
    date: '2025-12-01',
    period: 'December 2025',
    size: '298 KB',
    format: 'PDF',
  },
  {
    id: 6,
    name: 'Harvest Yield Report',
    type: 'Plant Health',
    date: '2025-11-20',
    period: 'November 2025',
    size: '445 KB',
    format: 'PDF',
  },
];

const filteredReports = computed(() =>
  reports.filter(
    (report) =>
      selectedReport.value === 'all' ||
      report.type.toLowerCase() === selectedReport.value
  )
);

function reportTypeBadgeClass(type: string) {
  if (type === 'Financial') return 'bg-blue-100 text-blue-700';
  if (type === 'Investment') return 'bg-green-100 text-green-700';
  return 'bg-purple-100 text-purple-700';
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-semibold text-gray-900">Reports & Documents</h2>
        <p class="text-gray-600 mt-1">Download and export your investment reports</p>
      </div>
      <button
        type="button"
        class="flex items-center gap-2 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
      >
        <Download class="w-4 h-4" />
        Generate Custom Report
      </button>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <button
        type="button"
        class="bg-white border border-gray-200 rounded-xl p-6 hover:border-green-600 hover:shadow-md transition-all text-left"
      >
        <div class="w-12 h-12 bg-green-50 rounded-lg flex items-center justify-center mb-4">
          <FileText class="w-6 h-6 text-green-600" />
        </div>
        <h3 class="font-semibold text-gray-900 mb-1">Portfolio Statement</h3>
        <p class="text-sm text-gray-600">Current holdings & performance</p>
      </button>

      <button
        type="button"
        class="bg-white border border-gray-200 rounded-xl p-6 hover:border-green-600 hover:shadow-md transition-all text-left"
      >
        <div class="w-12 h-12 bg-blue-50 rounded-lg flex items-center justify-center mb-4">
          <TrendingUp class="w-6 h-6 text-blue-600" />
        </div>
        <h3 class="font-semibold text-gray-900 mb-1">ROI Analysis</h3>
        <p class="text-sm text-gray-600">Returns & projections</p>
      </button>

      <button
        type="button"
        class="bg-white border border-gray-200 rounded-xl p-6 hover:border-green-600 hover:shadow-md transition-all text-left"
      >
        <div class="w-12 h-12 bg-purple-50 rounded-lg flex items-center justify-center mb-4">
          <Calendar class="w-6 h-6 text-purple-600" />
        </div>
        <h3 class="font-semibold text-gray-900 mb-1">Tax Documents</h3>
        <p class="text-sm text-gray-600">Annual tax statements</p>
      </button>

      <button
        type="button"
        class="bg-white border border-gray-200 rounded-xl p-6 hover:border-green-600 hover:shadow-md transition-all text-left"
      >
        <div class="w-12 h-12 bg-yellow-50 rounded-lg flex items-center justify-center mb-4">
          <FileText class="w-6 h-6 text-yellow-600" />
        </div>
        <h3 class="font-semibold text-gray-900 mb-1">Plant Reports</h3>
        <p class="text-sm text-gray-600">Growth & health data</p>
      </button>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-4">
      <div class="flex flex-wrap items-center gap-4">
        <div class="flex items-center gap-2">
          <Filter class="w-4 h-4 text-gray-600" />
          <span class="text-sm font-medium text-gray-700">Filter by:</span>
        </div>

        <div class="flex gap-2">
          <button
            v-for="type in reportTypes"
            :key="type.id"
            type="button"
            :class="[
              'px-4 py-2 rounded-lg text-sm font-medium transition-colors',
              selectedReport === type.id
                ? 'bg-green-600 text-white'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200',
            ]"
            @click="selectedReport = type.id"
          >
            {{ type.label }}
          </button>
        </div>

        <div class="ml-auto">
          <select
            v-model="selectedPeriod"
            class="px-4 py-2 border border-gray-200 rounded-lg text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-green-500"
          >
            <option value="all">All Time</option>
            <option value="monthly">This Month</option>
            <option value="quarterly">This Quarter</option>
            <option value="yearly">This Year</option>
          </select>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Report Name
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Type
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Period
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Date
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Size
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="report in filteredReports" :key="report.id" class="hover:bg-gray-50">
              <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 bg-red-50 rounded-lg flex items-center justify-center flex-shrink-0">
                    <FileText class="w-5 h-5 text-red-600" />
                  </div>
                  <div>
                    <p class="font-medium text-gray-900">{{ report.name }}</p>
                    <p class="text-sm text-gray-500">{{ report.format }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4">
                <span
                  class="inline-flex px-2.5 py-1 rounded-full text-xs font-medium"
                  :class="reportTypeBadgeClass(report.type)"
                >
                  {{ report.type }}
                </span>
              </td>
              <td class="px-6 py-4 text-gray-900">{{ report.period }}</td>
              <td class="px-6 py-4 text-gray-600">{{ report.date }}</td>
              <td class="px-6 py-4 text-gray-600">{{ report.size }}</td>
              <td class="px-6 py-4">
                <button
                  type="button"
                  class="flex items-center gap-2 px-3 py-1.5 text-green-600 hover:bg-green-50 rounded-lg font-medium text-sm transition-colors"
                >
                  <Download class="w-4 h-4" />
                  Download
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="bg-gray-50 border border-gray-200 rounded-xl p-6">
      <h3 class="font-semibold text-gray-900 mb-4">Custom Export Options</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="bg-white border border-gray-200 rounded-lg p-4">
          <p class="text-sm font-medium text-gray-900 mb-2">Date Range</p>
          <input
            type="date"
            class="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm"
            value="2025-01-01"
          />
        </div>
        <div class="bg-white border border-gray-200 rounded-lg p-4">
          <p class="text-sm font-medium text-gray-900 mb-2">Report Type</p>
          <select class="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm">
            <option>All Types</option>
            <option>Financial Only</option>
            <option>Investment Only</option>
            <option>Plant Health Only</option>
          </select>
        </div>
        <div class="bg-white border border-gray-200 rounded-lg p-4">
          <p class="text-sm font-medium text-gray-900 mb-2">Format</p>
          <select class="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm">
            <option>PDF</option>
            <option>Excel (XLSX)</option>
            <option>CSV</option>
          </select>
        </div>
        <div class="flex items-end">
          <button
            type="button"
            class="w-full px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-medium"
          >
            Generate Report
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
