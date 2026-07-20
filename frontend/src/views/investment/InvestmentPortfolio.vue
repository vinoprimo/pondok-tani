<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { TrendingUp, DollarSign, Calendar, Download, Eye } from 'lucide-vue-next'
import { getCurrentUser } from '../../services/user/user'
import { listOperationalCosts } from '../../services/financial/operational-cost'
import { listActualRevenues } from '../../services/financial/revenue'
import { listPlantMonitorings } from '../../services/plant-monitoring/monitoring'

type ApiInvestment = {
  id: number
  package_id: number
  amount: number | string
  status: string
  investment_date: string
  package?: {
    package_name?: string
  }
}

type PortfolioRow = {
  id: string
  name: string
  invested: number
  status: string
  startDate: string
}

type MonitoringRow = {
  plant_batch_id: number
  batch_code: string
  package_name: string
}

type OperationalCostApiRow = {
  plant_batch_id: number
  batch_code: string
  package_name: string
  amount: number | string
  cost_date: string
}

type RevenueApiRow = {
  plant_batch_id: number
  batch_code: string
  package_name: string
  amount: number | string
  revenue_date: string
}

type BatchPerformanceRow = {
  key: string
  batchCode: string
  packageName: string
  invested: number
  expense: number
  revenue: number
  netBalance: number
  totalReturn: number
  roiPercent: number
  latestDate: string
}

const loading = ref(false)
const errorMessage = ref('')
const investments = ref<PortfolioRow[]>([])
const batchPerformanceRows = ref<BatchPerformanceRow[]>([])
const monthlyNetSeries = ref<Array<{ key: string; month: string; net: number }>>([])

const formatRupiah = (value: number) =>
  new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(value)

const formatSignedRupiah = (value: number) => {
  if (value < 0) {
    return `-${formatRupiah(Math.abs(value))}`
  }
  return formatRupiah(value)
}

const formatSignedPercent = (value: number) => {
  if (!Number.isFinite(value)) return '0.0%'
  const sign = value > 0 ? '+' : ''
  return `${sign}${value.toFixed(1)}%`
}

const formatDate = (value: string) => {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  })
}

const statusMeta = (status: string) => {
  const normalized = String(status || '').toLowerCase()
  if (normalized === 'active' || normalized === 'on_process') {
    return {
      label: normalized === 'active' ? 'Aktif' : 'On process',
      badgeClass: 'bg-emerald-50 text-emerald-700',
    }
  }
  if (normalized === 'pending') {
    return {
      label: 'Pending',
      badgeClass: 'bg-amber-50 text-amber-700',
    }
  }
  return {
    label: normalized || 'Unknown',
    badgeClass: 'bg-gray-100 text-gray-700',
  }
}

const totalInvested = computed(() => investments.value.reduce((acc, item) => acc + item.invested, 0))
const totalOperationalCost = computed(() => batchPerformanceRows.value.reduce((acc, item) => acc + item.expense, 0))
const totalRevenue = computed(() => batchPerformanceRows.value.reduce((acc, item) => acc + item.revenue, 0))
const netBalance = computed(() => totalRevenue.value - totalOperationalCost.value)
const totalReturn = computed(() => netBalance.value - totalInvested.value)
const portfolioRoiPercent = computed(() => {
  if (totalInvested.value <= 0) return 0
  return (totalReturn.value / totalInvested.value) * 100
})

const breakEvenStatus = computed(() => {
  if (totalInvested.value <= 0) return 'Belum ada modal investasi tercatat'
  if (totalReturn.value >= 0) return 'Break-even tercapai'
  return 'Belum break-even'
})

const breakEvenClass = computed(() => {
  return totalReturn.value >= 0 ? 'text-emerald-100' : 'text-rose-100'
})

const ytdInvested = computed(() => {
  const now = new Date()
  const startOfYear = new Date(now.getFullYear(), 0, 1)
  return investments.value.reduce((acc, item) => {
    const d = new Date(item.startDate)
    if (Number.isNaN(d.getTime())) return acc
    return d >= startOfYear ? acc + item.invested : acc
  }, 0)
})

const currentMonthInvested = computed(() => {
  const now = new Date()
  return monthlyNetSeries.value.reduce((acc, item) => {
    const d = new Date(`${item.key}-01`)
    if (Number.isNaN(d.getTime())) return acc
    const isCurrentMonth = d.getFullYear() === now.getFullYear() && d.getMonth() === now.getMonth()
    return isCurrentMonth ? acc + item.net : acc
  }, 0)
})

const latestInvestment = computed(() => {
  if (!batchPerformanceRows.value.length) return null
  const sorted = [...batchPerformanceRows.value].sort((a, b) => b.latestDate.localeCompare(a.latestDate))
  return sorted[0]
})

const portfolioData = computed(() => {
  if (monthlyNetSeries.value.length > 0) {
    let running = 0
    return monthlyNetSeries.value.map((item) => {
      running += item.net
      return {
        month: item.month,
        value: running,
      }
    })
  }

  const now = new Date()
  const monthStart = new Date(now.getFullYear(), now.getMonth(), 1)
  const monthKeys: string[] = []
  for (let i = 5; i >= 0; i -= 1) {
    const d = new Date(monthStart.getFullYear(), monthStart.getMonth() - i, 1)
    monthKeys.push(`${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`)
  }

  const byMonth = new Map<string, number>()
  monthKeys.forEach((key) => byMonth.set(key, 0))

  investments.value.forEach((item) => {
    const d = new Date(item.startDate)
    if (Number.isNaN(d.getTime())) return
    const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`
    if (!byMonth.has(key)) return
    byMonth.set(key, (byMonth.get(key) || 0) + item.invested)
  })

  let running = 0
  return monthKeys.map((key) => {
    running += byMonth.get(key) || 0
    const [year, month] = key.split('-').map(Number)
    const d = new Date(year, month - 1, 1)
    return {
      month: d.toLocaleDateString('id-ID', { month: 'short' }),
      value: running,
    }
  })
})

const balanceChart = {
  width: 820,
  height: 360,
  marginTop: 18,
  marginRight: 28,
  marginBottom: 32,
  marginLeft: 54,
}

const balanceChartInnerWidth = computed(
  () => balanceChart.width - balanceChart.marginLeft - balanceChart.marginRight
)

const balanceChartInnerHeight = computed(
  () => balanceChart.height - balanceChart.marginTop - balanceChart.marginBottom
)

const balanceChartMax = computed(() => {
  const maxValue = Math.max(...portfolioData.value.map((d) => Math.max(d.value, 0)), 0)
  const step = 15000
  return Math.max(60000, Math.ceil(maxValue / step) * step || step)
})

const balanceChartTicks = computed(() => {
  const step = balanceChartMax.value / 4
  return Array.from({ length: 5 }, (_, index) => index * step)
})

function getBalanceChartX(index: number) {
  if (portfolioData.value.length <= 1) {
    return balanceChart.marginLeft + balanceChartInnerWidth.value / 2
  }

  return balanceChart.marginLeft + (index / (portfolioData.value.length - 1)) * balanceChartInnerWidth.value
}

function getBalanceChartY(value: number) {
  const normalized = Math.max(0, Math.min(value, balanceChartMax.value))
  const ratio = normalized / balanceChartMax.value
  return balanceChart.marginTop + (1 - ratio) * balanceChartInnerHeight.value
}

const balanceChartPoints = computed(() =>
  portfolioData.value.map((item, index) => ({
    ...item,
    x: getBalanceChartX(index),
    y: getBalanceChartY(item.value),
  }))
)

const balanceChartLinePath = computed(() => {
  if (!balanceChartPoints.value.length) return ''
  return balanceChartPoints.value
    .map((point, index) => `${index === 0 ? 'M' : 'L'} ${point.x.toFixed(2)} ${point.y.toFixed(2)}`)
    .join(' ')
})

const balanceChartAreaPath = computed(() => {
  if (!balanceChartPoints.value.length) return ''
  const baselineY = balanceChart.height - balanceChart.marginBottom
  const first = balanceChartPoints.value[0]
  const last = balanceChartPoints.value[balanceChartPoints.value.length - 1]
  return [
    `M ${first.x.toFixed(2)} ${baselineY}`,
    balanceChartLinePath.value.replace(/^M/, 'L'),
    `L ${last.x.toFixed(2)} ${baselineY}`,
    'Z',
  ].join(' ')
})

async function loadPortfolio() {
  loading.value = true
  errorMessage.value = ''
  try {
    const [userRes, costsRes, revenuesRes, monitoringRes] = await Promise.all([
      getCurrentUser(),
      listOperationalCosts(),
      listActualRevenues(),
      listPlantMonitorings(),
    ])

    const rows = Array.isArray(userRes?.data?.investments) ? (userRes.data.investments as ApiInvestment[]) : []
    const costs = Array.isArray(costsRes?.data) ? (costsRes.data as OperationalCostApiRow[]) : []
    const revenues = Array.isArray(revenuesRes?.data) ? (revenuesRes.data as RevenueApiRow[]) : []
    const monitorings = Array.isArray(monitoringRes?.data) ? (monitoringRes.data as MonitoringRow[]) : []

    investments.value = rows
      .map((item) => {
        const invested = Number(item.amount) || 0
        return {
          id: `INV-${item.id}`,
          name: item.package?.package_name || `Paket #${item.package_id}`,
          invested,
          status: item.status || 'pending',
          startDate: item.investment_date,
        }
      })
      .sort((a, b) => b.startDate.localeCompare(a.startDate))

    const packageInvestmentMap = new Map<string, number>()
    investments.value.forEach((item) => {
      const key = item.name.trim().toLowerCase()
      packageInvestmentMap.set(key, (packageInvestmentMap.get(key) || 0) + item.invested)
    })

    const batchMetaMap = new Map<string, { batchCode: string; packageName: string }>()
    monitorings.forEach((item) => {
      if (!item.batch_code) return
      const key = item.batch_code.trim().toLowerCase()
      if (!batchMetaMap.has(key)) {
        batchMetaMap.set(key, {
          batchCode: item.batch_code,
          packageName: item.package_name || 'Paket tanpa nama',
        })
      }
    })

    costs.forEach((item) => {
      if (!item.batch_code) return
      const key = item.batch_code.trim().toLowerCase()
      if (!batchMetaMap.has(key)) {
        batchMetaMap.set(key, {
          batchCode: item.batch_code,
          packageName: item.package_name || 'Paket tanpa nama',
        })
      }
    })

    revenues.forEach((item) => {
      if (!item.batch_code) return
      const key = item.batch_code.trim().toLowerCase()
      if (!batchMetaMap.has(key)) {
        batchMetaMap.set(key, {
          batchCode: item.batch_code,
          packageName: item.package_name || 'Paket tanpa nama',
        })
      }
    })

    const packageBatchCount = new Map<string, number>()
    Array.from(batchMetaMap.values()).forEach((meta) => {
      const key = meta.packageName.trim().toLowerCase()
      packageBatchCount.set(key, (packageBatchCount.get(key) || 0) + 1)
    })

    const batchRows = new Map<string, BatchPerformanceRow>()
    batchMetaMap.forEach((meta, key) => {
      const packageKey = meta.packageName.trim().toLowerCase()
      const packageInvested = packageInvestmentMap.get(packageKey) || 0
      const batchCount = packageBatchCount.get(packageKey) || 1
      const investedAllocated = packageInvested > 0 ? packageInvested / batchCount : 0
      batchRows.set(key, {
        key,
        batchCode: meta.batchCode,
        packageName: meta.packageName,
        invested: investedAllocated,
        expense: 0,
        revenue: 0,
        netBalance: 0,
        totalReturn: 0,
        roiPercent: 0,
        latestDate: '',
      })
    })

    const monthMap = new Map<string, number>()
    const pushMonthDelta = (rawDate: string, delta: number) => {
      const date = new Date(rawDate)
      if (Number.isNaN(date.getTime())) return
      const key = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
      monthMap.set(key, (monthMap.get(key) || 0) + delta)
    }

    costs.forEach((item) => {
      const key = String(item.batch_code || '').trim().toLowerCase()
      if (!key || !batchRows.has(key)) return
      const row = batchRows.get(key)!
      row.expense += Number(item.amount) || 0
      if (!row.latestDate || item.cost_date > row.latestDate) {
        row.latestDate = item.cost_date
      }
      pushMonthDelta(item.cost_date, -(Number(item.amount) || 0))
    })

    revenues.forEach((item) => {
      const key = String(item.batch_code || '').trim().toLowerCase()
      if (!key || !batchRows.has(key)) return
      const row = batchRows.get(key)!
      row.revenue += Number(item.amount) || 0
      if (!row.latestDate || item.revenue_date > row.latestDate) {
        row.latestDate = item.revenue_date
      }
      pushMonthDelta(item.revenue_date, Number(item.amount) || 0)
    })

    batchPerformanceRows.value = Array.from(batchRows.values())
      .map((row) => {
        const netBalanceValue = row.revenue - row.expense
        const totalReturnValue = netBalanceValue - row.invested
        const roiPercentValue = row.invested > 0 ? (totalReturnValue / row.invested) * 100 : 0

        return {
          ...row,
          netBalance: netBalanceValue,
          totalReturn: totalReturnValue,
          roiPercent: roiPercentValue,
          latestDate: row.latestDate || new Date().toISOString().slice(0, 10),
        }
      })
      .sort((a, b) => b.latestDate.localeCompare(a.latestDate))

    const now = new Date()
    const monthStart = new Date(now.getFullYear(), now.getMonth(), 1)
    const keys: string[] = []
    for (let i = 5; i >= 0; i -= 1) {
      const d = new Date(monthStart.getFullYear(), monthStart.getMonth() - i, 1)
      keys.push(`${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`)
    }

    monthlyNetSeries.value = keys.map((key) => {
      const [year, month] = key.split('-').map(Number)
      const date = new Date(year, month - 1, 1)
      return {
        key,
        month: date.toLocaleDateString('id-ID', { month: 'short' }),
        net: monthMap.get(key) || 0,
      }
    })
  } catch (error) {
    console.error('Failed to load portfolio data', error)
    errorMessage.value = 'Gagal memuat data portofolio dari server.'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadPortfolio()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-semibold text-gray-900">Portofolio investasi</h2>
        <p class="text-gray-600 mt-1">Pantau investasi perkebunan vanili dan imbal hasil Anda</p>
      </div>
      <button
        type="button"
        class="flex items-center gap-2 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
      >
        <Download class="w-4 h-4" />
        Ekspor laporan
      </button>
    </div>

    <div v-if="errorMessage" class="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
      {{ errorMessage }}
    </div>

    <div class="bg-gradient-to-br from-green-600 to-green-700 rounded-xl p-6 text-white">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div>
          <p class="text-green-100 text-sm mb-1">Total investasi</p>
          <p class="text-3xl font-semibold">{{ formatRupiah(totalInvested) }}</p>
        </div>
        <div>
          <p class="text-green-100 text-sm mb-1">Saldo bersih</p>
          <p class="text-3xl font-semibold">{{ formatSignedRupiah(netBalance) }}</p>
        </div>
        <div>
          <p class="text-green-100 text-sm mb-1">Total imbal hasil</p>
          <p class="text-3xl font-semibold">{{ formatSignedRupiah(totalReturn) }}</p>
        </div>
        <div>
          <p class="text-green-100 text-sm mb-1">ROI portofolio</p>
          <p class="text-3xl font-semibold">{{ formatSignedPercent(portfolioRoiPercent) }}</p>
        </div>
      </div>
      <p class="mt-4 text-xs" :class="breakEvenClass">{{ breakEvenStatus }}</p>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-6">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">Pertumbuhan saldo bersih</h3>
      <p class="text-xs text-gray-500 mb-3">Akumulasi arus kas bersih 6 bulan terakhir (pendapatan - biaya operasional)</p>
      <div class="mb-4 overflow-x-auto">
        <svg
          class="h-[340px] min-w-[720px] w-full"
          :viewBox="`0 0 ${balanceChart.width} ${balanceChart.height}`"
          role="img"
          aria-label="Grafik pertumbuhan saldo bersih"
        >
          <defs>
            <linearGradient id="net-balance-area" x1="0" x2="0" y1="0" y2="1">
              <stop offset="0%" stop-color="#34d399" stop-opacity="0.38" />
              <stop offset="62%" stop-color="#a7f3d0" stop-opacity="0.18" />
              <stop offset="100%" stop-color="#ffffff" stop-opacity="0.78" />
            </linearGradient>
          </defs>

          <g>
            <line
              :x1="balanceChart.marginLeft"
              :x2="balanceChart.width - balanceChart.marginRight"
              :y1="balanceChart.height - balanceChart.marginBottom"
              :y2="balanceChart.height - balanceChart.marginBottom"
              stroke="#9ca3af"
              stroke-width="1"
            />
            <line
              :x1="balanceChart.marginLeft"
              :x2="balanceChart.marginLeft"
              :y1="balanceChart.marginTop"
              :y2="balanceChart.height - balanceChart.marginBottom"
              stroke="#9ca3af"
              stroke-width="1"
            />

            <g v-for="tick in balanceChartTicks" :key="`y-${tick}`">
              <line
                :x1="balanceChart.marginLeft"
                :x2="balanceChart.width - balanceChart.marginRight"
                :y1="getBalanceChartY(tick)"
                :y2="getBalanceChartY(tick)"
                stroke="#dbe3ec"
                stroke-dasharray="3 4"
                stroke-width="1"
              />
              <text
                :x="balanceChart.marginLeft - 8"
                :y="getBalanceChartY(tick) + 5"
                text-anchor="end"
                class="fill-slate-400 text-[14px]"
              >
                {{ Math.round(tick) }}
              </text>
            </g>

            <g v-for="(point, index) in balanceChartPoints" :key="`x-${point.month}-${index}`">
              <line
                :x1="point.x"
                :x2="point.x"
                :y1="balanceChart.marginTop"
                :y2="balanceChart.height - balanceChart.marginBottom"
                stroke="#dbe3ec"
                stroke-dasharray="3 4"
                stroke-width="1"
              />
              <text
                :x="point.x"
                :y="balanceChart.height - 14"
                text-anchor="middle"
                class="fill-slate-400 text-[14px]"
              >
                {{ point.month }}
              </text>
            </g>

            <path :d="balanceChartAreaPath" fill="url(#net-balance-area)" />
            <path
              :d="balanceChartLinePath"
              fill="none"
              stroke="#10b981"
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
            />
          </g>
        </svg>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead class="bg-gray-50 text-left text-gray-600">
            <tr>
              <th class="px-3 py-2 font-medium">Bulan</th>
              <th class="px-3 py-2 font-medium text-right">Saldo bersih akumulatif</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="row in portfolioData" :key="row.month">
              <td class="px-3 py-2 font-medium text-gray-900">{{ row.month }}</td>
              <td class="px-3 py-2 text-right" :class="row.value >= 0 ? 'text-gray-900' : 'text-rose-700'">{{ formatSignedRupiah(row.value) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden">
      <div class="p-6 border-b border-gray-200">
        <h3 class="text-lg font-semibold text-gray-900">Kinerja per batch</h3>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Batch</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Investasi dialokasikan
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Saldo bersih</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Total imbal hasil</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ROI %</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-if="loading">
              <td colspan="6" class="px-6 py-4 text-sm text-gray-500">Memuat data portofolio...</td>
            </tr>
            <tr v-else-if="!batchPerformanceRows.length">
              <td colspan="6" class="px-6 py-4 text-sm text-gray-500">Belum ada data batch untuk dihitung.</td>
            </tr>
            <tr v-for="row in batchPerformanceRows" :key="row.key" class="hover:bg-gray-50">
              <td class="px-6 py-4">
                <div>
                  <p class="font-medium text-gray-900">{{ row.batchCode }}</p>
                  <p class="text-sm text-gray-500">{{ row.packageName }}</p>
                </div>
              </td>
              <td class="px-6 py-4 text-gray-900">{{ formatRupiah(row.invested) }}</td>
              <td class="px-6 py-4" :class="row.netBalance >= 0 ? 'text-gray-900' : 'text-rose-700'">{{ formatSignedRupiah(row.netBalance) }}</td>
              <td class="px-6 py-4" :class="row.totalReturn >= 0 ? 'text-gray-900' : 'text-rose-700'">{{ formatSignedRupiah(row.totalReturn) }}</td>
              <td class="px-6 py-4">
                <span
                  class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-sm font-medium"
                  :class="row.roiPercent >= 0 ? 'bg-green-50 text-green-700' : 'bg-rose-50 text-rose-700'"
                >
                  <TrendingUp class="w-3 h-3" />
                  {{ formatSignedPercent(row.roiPercent) }}
                </span>
              </td>
              <td class="px-6 py-4">
                <span
                  class="inline-flex rounded-full px-2.5 py-1 text-xs font-semibold"
                  :class="row.totalReturn >= 0 ? 'bg-emerald-50 text-emerald-700' : 'bg-amber-50 text-amber-700'"
                >
                  {{ row.totalReturn >= 0 ? 'Break-even' : 'Belum break-even' }}
                </span>
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
          <h3 class="font-semibold text-gray-900">Bulan lalu</h3>
        </div>
        <p class="text-2xl font-semibold mb-1" :class="currentMonthInvested >= 0 ? 'text-gray-900' : 'text-rose-700'">{{ formatSignedRupiah(currentMonthInvested) }}</p>
        <p class="text-sm text-gray-600">Arus kas bersih bulan berjalan</p>
      </div>
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-4">
          <div class="w-10 h-10 bg-green-50 rounded-lg flex items-center justify-center">
            <Calendar class="w-5 h-5 text-green-600" />
          </div>
          <h3 class="font-semibold text-gray-900">Batch terbaru</h3>
        </div>
        <p class="text-2xl font-semibold text-gray-900 mb-1">{{ latestInvestment ? latestInvestment.batchCode : '-' }}</p>
        <p class="text-sm text-gray-600">{{ latestInvestment ? formatDate(latestInvestment.latestDate) : 'Belum ada transaksi batch' }}</p>
      </div>
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-4">
          <div class="w-10 h-10 bg-purple-50 rounded-lg flex items-center justify-center">
            <TrendingUp class="w-5 h-5 text-purple-600" />
          </div>
          <h3 class="font-semibold text-gray-900">ROI YTD</h3>
        </div>
        <p class="text-2xl font-semibold mb-1" :class="portfolioRoiPercent >= 0 ? 'text-gray-900' : 'text-rose-700'">{{ formatSignedPercent(portfolioRoiPercent) }}</p>
        <p class="text-sm text-gray-600">Persentase ROI terhadap total investasi</p>
      </div>
    </div>
  </div>
</template>
