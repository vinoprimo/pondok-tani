<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import {
  Calendar,
  Download,
  FileText,
  Filter,
  Loader2,
  Package,
  Sprout,
  TrendingUp,
  Wallet,
} from 'lucide-vue-next'
import SearchBar from '../../components/common/SearchBar.vue'
import Pagination from '../../components/common/Pagination.vue'
import { listHarvestRequests } from '../../services/harvest/harvest'
import { listOperationalCosts } from '../../services/financial/operational-cost'
import { listActualRevenues } from '../../services/financial/revenue'
import { getCurrentUser } from '../../services/user/user'
import { listStockMovements, listWarehouseStocks } from '../../services/warehouse/warehouse'

type ReportCategory = 'financial' | 'plant' | 'investment' | 'warehouse'

type ReportRow = Record<string, string | number>

type ReportItem = {
  id: string
  name: string
  category: ReportCategory
  date: string
  period: string
  size: string
  format: string
  rows: ReportRow[]
}

const selectedPeriod = ref('all')
const selectedReport = ref('all')
const searchQuery = ref('')
const currentPage = ref(1)
const limit = ref(5)
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const harvestRequests = ref<any[]>([])
const operationalCosts = ref<any[]>([])
const actualRevenues = ref<any[]>([])
const investments = ref<any[]>([])
const warehouseStocks = ref<any[]>([])
const stockMovements = ref<any[]>([])

const customForm = reactive({
  startDate: new Date(new Date().getFullYear(), 0, 1).toISOString().slice(0, 10),
  endDate: new Date().toISOString().slice(0, 10),
  reportType: 'all',
})

const reportTypes = [
  { id: 'all', label: 'Semua laporan' },
  { id: 'financial', label: 'Keuangan' },
  { id: 'plant', label: 'Panen & tanaman' },
  { id: 'investment', label: 'Investasi' },
  { id: 'warehouse', label: 'Stok gudang' },
]

const categoryLabels: Record<ReportCategory, string> = {
  financial: 'Keuangan',
  plant: 'Panen & tanaman',
  investment: 'Investasi',
  warehouse: 'Stok gudang',
}

const reportCards = [
  {
    id: 'investment',
    title: 'Rekening portofolio',
    description: 'Kepemilikan dan status investasi',
    icon: FileText,
    color: 'bg-green-50 text-green-600',
  },
  {
    id: 'financial',
    title: 'Ringkasan keuangan',
    description: 'Pendapatan, biaya, dan estimasi bersih',
    icon: TrendingUp,
    color: 'bg-blue-50 text-blue-600',
  },
  {
    id: 'warehouse',
    title: 'Laporan stok gudang',
    description: 'Stok akhir dan mutasi gudang',
    icon: Package,
    color: 'bg-amber-50 text-amber-600',
  },
  {
    id: 'plant',
    title: 'Laporan panen',
    description: 'Ajuan, validasi, dan hasil panen',
    icon: Sprout,
    color: 'bg-purple-50 text-purple-600',
  },
]

const periodLabel = computed(() => {
  if (selectedPeriod.value === 'monthly') return 'Bulan ini'
  if (selectedPeriod.value === 'quarterly') return 'Kuartal ini'
  if (selectedPeriod.value === 'yearly') return 'Tahun ini'
  return 'Sepanjang waktu'
})

const dateRange = computed(() => getPeriodRange(selectedPeriod.value))

const totalRevenue = computed(() =>
  actualRevenues.value.reduce((sum, item) => sum + Number(item.amount || 0), 0)
)
const totalCost = computed(() =>
  operationalCosts.value.reduce((sum, item) => sum + Number(item.amount || 0), 0)
)
const totalHarvest = computed(() =>
  harvestRequests.value.reduce((sum, item) => sum + Number(item.total_quantity || 0), 0)
)
const totalWarehouseStock = computed(() =>
  warehouseStocks.value.reduce((sum, item) => sum + Number(item.total_quantity || 0), 0)
)

const reports = computed<ReportItem[]>(() => [
  buildInvestmentReport(),
  buildFinancialReport(),
  buildWarehouseReport(),
  buildHarvestReport(),
])

const filteredReports = computed(() => {
  let result = reports.value
  if (selectedReport.value !== 'all') {
    result = result.filter((report) => report.category === selectedReport.value)
  }
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter((report) =>
      report.name.toLowerCase().includes(query) ||
      report.period.toLowerCase().includes(query) ||
      categoryLabels[report.category].toLowerCase().includes(query)
    )
  }
  return result
})

const totalPages = computed(() =>
  Math.max(1, Math.ceil(filteredReports.value.length / limit.value))
)

const paginatedReports = computed(() => {
  const start = (currentPage.value - 1) * limit.value
  return filteredReports.value.slice(start, start + limit.value)
})

watch([selectedReport, selectedPeriod], () => {
  currentPage.value = 1
})

function getPeriodRange(period: string) {
  const now = new Date()
  if (period === 'monthly') {
    return {
      start: new Date(now.getFullYear(), now.getMonth(), 1),
      end: now,
    }
  }
  if (period === 'quarterly') {
    const quarterStartMonth = Math.floor(now.getMonth() / 3) * 3
    return {
      start: new Date(now.getFullYear(), quarterStartMonth, 1),
      end: now,
    }
  }
  if (period === 'yearly') {
    return {
      start: new Date(now.getFullYear(), 0, 1),
      end: now,
    }
  }
  return null
}

function isInsideSelectedPeriod(value: string | Date | undefined) {
  const range = dateRange.value
  if (!range || !value) return true
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return false
  return date >= range.start && date <= range.end
}

function filterRowsByDate<T extends Record<string, any>>(rows: T[], keys: string[]) {
  return rows.filter((row) => {
    const dateValue = keys.map((key) => row[key]).find(Boolean)
    return isInsideSelectedPeriod(dateValue)
  })
}

function formatDate(value: string | Date | undefined) {
  if (!value) return '-'
  const date = value instanceof Date ? value : new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  })
}

function formatRupiah(value: number) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(value || 0)
}

function formatNumber(value: number) {
  return new Intl.NumberFormat('id-ID', {
    maximumFractionDigits: 2,
  }).format(value || 0)
}

function makeSize(rows: ReportRow[]) {
  return `${Math.max(1, Math.ceil(JSON.stringify(rows).length / 1024))} KB`
}

function getToday() {
  return new Date().toISOString().slice(0, 10)
}

function buildInvestmentReport(): ReportItem {
  const rows = investments.value.map((item) => ({
    'ID Investasi': item.id || '-',
    Paket: item.package?.package_name || item.package_name || '-',
    Status: item.status || item.package_status || '-',
    'Jumlah pohon': item.plant_count || item.total_plants || '-',
    'Tanggal dibuat': formatDate(item.created_at),
  }))
  return {
    id: 'investment-portfolio',
    name: 'Laporan rekening portofolio',
    category: 'investment',
    date: getToday(),
    period: periodLabel.value,
    size: makeSize(rows),
    format: 'CSV',
    rows,
  }
}

function buildFinancialReport(): ReportItem {
  const revenueRows = filterRowsByDate(actualRevenues.value, ['revenue_date']).map((item) => ({
    Tipe: 'Pendapatan',
    Batch: item.batch_code || '-',
    Kategori: item.source || '-',
    Nominal: Number(item.amount || 0),
    Tanggal: formatDate(item.revenue_date),
    Status: item.status || '-',
  }))
  const costRows = filterRowsByDate(operationalCosts.value, ['cost_date']).map((item) => ({
    Tipe: 'Biaya',
    Batch: item.batch_code || '-',
    Kategori: item.category || '-',
    Nominal: Number(item.amount || 0),
    Tanggal: formatDate(item.cost_date),
    Status: '-',
  }))
  const rows = [...revenueRows, ...costRows]
  return {
    id: 'financial-summary',
    name: 'Ringkasan keuangan',
    category: 'financial',
    date: getToday(),
    period: periodLabel.value,
    size: makeSize(rows),
    format: 'CSV',
    rows,
  }
}

function buildWarehouseReport(): ReportItem {
  const stockRows = warehouseStocks.value.map((item) => ({
    Tipe: 'Stok akhir',
    Batch: item.batch_code || '-',
    Paket: item.package_name || '-',
    Mutu: item.grade_name || 'Basah',
    Kuantitas: Number(item.total_quantity || 0),
    Tanggal: formatDate(item.updated_at),
  }))
  const movementRows = filterRowsByDate(stockMovements.value, ['movement_date']).map((item) => ({
    Tipe: item.movement_type === 'out' ? 'Keluar' : 'Masuk',
    Batch: item.batch_code || '-',
    Paket: item.package_name || '-',
    Mutu: item.grade_name || 'Basah',
    Kuantitas: Number(item.quantity || 0),
    Tanggal: formatDate(item.movement_date),
  }))
  const rows = [...stockRows, ...movementRows]
  return {
    id: 'warehouse-stock',
    name: 'Laporan stok gudang',
    category: 'warehouse',
    date: getToday(),
    period: periodLabel.value,
    size: makeSize(rows),
    format: 'CSV',
    rows,
  }
}

function buildHarvestReport(): ReportItem {
  const rows = filterRowsByDate(harvestRequests.value, ['harvest_date', 'created_at']).map((item) => ({
    Batch: item.batch_code || '-',
    Paket: item.package_name || '-',
    Status: item.status || '-',
    'Total panen': Number(item.total_quantity || 0),
    Basah: Number(item.wet_quantity || 0),
    Kering: Number(item.dry_quantity || 0),
    Tanggal: formatDate(item.harvest_date || item.created_at),
  }))
  return {
    id: 'harvest-report',
    name: 'Laporan hasil panen',
    category: 'plant',
    date: getToday(),
    period: periodLabel.value,
    size: makeSize(rows),
    format: 'CSV',
    rows,
  }
}

function reportTypeBadgeClass(category: ReportCategory) {
  if (category === 'financial') return 'bg-blue-100 text-blue-700'
  if (category === 'investment') return 'bg-green-100 text-green-700'
  if (category === 'warehouse') return 'bg-amber-100 text-amber-700'
  return 'bg-purple-100 text-purple-700'
}

function handleSearch(val: string) {
  searchQuery.value = val
  currentPage.value = 1
}

function escapeCsv(value: string | number) {
  const text = String(value ?? '')
  if (/[",\n]/.test(text)) {
    return `"${text.replace(/"/g, '""')}"`
  }
  return text
}

function downloadCsv(name: string, rows: ReportRow[]) {
  if (!rows.length) {
    successMessage.value = ''
    errorMessage.value = 'Belum ada data untuk laporan ini.'
    return
  }
  const headers = Array.from(rows.reduce((set, row) => {
    Object.keys(row).forEach((key) => set.add(key))
    return set
  }, new Set<string>()))
  const csv = [
    headers.map(escapeCsv).join(','),
    ...rows.map((row) => headers.map((header) => escapeCsv(row[header] ?? '')).join(',')),
  ].join('\n')
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `${name.toLowerCase().replace(/[^a-z0-9]+/g, '-')}-${getToday()}.csv`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
  errorMessage.value = ''
  successMessage.value = 'Laporan berhasil dibuat.'
}

function downloadReport(report: ReportItem) {
  downloadCsv(report.name, report.rows)
}

function makeCustomReportRows() {
  const start = new Date(customForm.startDate)
  const end = new Date(customForm.endDate)
  end.setHours(23, 59, 59, 999)

  const inRange = (value: string | Date | undefined) => {
    if (!value) return false
    const date = new Date(value)
    if (Number.isNaN(date.getTime())) return false
    return date >= start && date <= end
  }

  const selected = customForm.reportType
  const rows: ReportRow[] = []

  if (selected === 'all' || selected === 'financial') {
    actualRevenues.value.filter((item) => inRange(item.revenue_date)).forEach((item) => {
      rows.push({
        Laporan: 'Keuangan',
        Tipe: 'Pendapatan',
        Batch: item.batch_code || '-',
        Deskripsi: item.source || '-',
        Nilai: Number(item.amount || 0),
        Tanggal: formatDate(item.revenue_date),
      })
    })
    operationalCosts.value.filter((item) => inRange(item.cost_date)).forEach((item) => {
      rows.push({
        Laporan: 'Keuangan',
        Tipe: 'Biaya',
        Batch: item.batch_code || '-',
        Deskripsi: item.category || '-',
        Nilai: Number(item.amount || 0),
        Tanggal: formatDate(item.cost_date),
      })
    })
  }

  if (selected === 'all' || selected === 'plant') {
    harvestRequests.value
      .filter((item) => inRange(item.harvest_date || item.created_at))
      .forEach((item) => {
        rows.push({
          Laporan: 'Panen',
          Tipe: item.status || '-',
          Batch: item.batch_code || '-',
          Deskripsi: item.package_name || '-',
          Nilai: Number(item.total_quantity || 0),
          Tanggal: formatDate(item.harvest_date || item.created_at),
        })
      })
  }

  if (selected === 'all' || selected === 'warehouse') {
    stockMovements.value.filter((item) => inRange(item.movement_date)).forEach((item) => {
      rows.push({
        Laporan: 'Stok gudang',
        Tipe: item.movement_type === 'out' ? 'Keluar' : 'Masuk',
        Batch: item.batch_code || '-',
        Deskripsi: item.grade_name || 'Basah',
        Nilai: Number(item.quantity || 0),
        Tanggal: formatDate(item.movement_date),
      })
    })
  }

  if (selected === 'all' || selected === 'investment') {
    investments.value.forEach((item) => {
      rows.push({
        Laporan: 'Investasi',
        Tipe: item.status || item.package_status || '-',
        Batch: item.batch_code || '-',
        Deskripsi: item.package?.package_name || item.package_name || '-',
        Nilai: item.plant_count || item.total_plants || 0,
        Tanggal: formatDate(item.created_at),
      })
    })
  }

  return rows
}

function createCustomReport() {
  if (!customForm.startDate || !customForm.endDate) {
    errorMessage.value = 'Rentang tanggal wajib diisi.'
    return
  }
  if (new Date(customForm.startDate) > new Date(customForm.endDate)) {
    errorMessage.value = 'Tanggal mulai tidak boleh melewati tanggal selesai.'
    return
  }
  downloadCsv('laporan-kustom', makeCustomReportRows())
}

async function fetchReportsData() {
  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''

  const loaders = [
    {
      name: 'data panen',
      load: async () => {
        const response = await listHarvestRequests()
        harvestRequests.value = Array.isArray(response.data) ? response.data : []
      },
    },
    {
      name: 'biaya operasional',
      load: async () => {
        const response = await listOperationalCosts()
        operationalCosts.value = Array.isArray(response.data) ? response.data : []
      },
    },
    {
      name: 'pendapatan',
      load: async () => {
        const response = await listActualRevenues()
        actualRevenues.value = Array.isArray(response.data) ? response.data : []
      },
    },
    {
      name: 'investasi',
      load: async () => {
        const response = await getCurrentUser()
        investments.value = Array.isArray(response.data?.investments) ? response.data.investments : []
      },
    },
    {
      name: 'stok gudang',
      load: async () => {
        const response = await listWarehouseStocks()
        warehouseStocks.value = Array.isArray(response.data) ? response.data : []
      },
    },
    {
      name: 'pergerakan stok',
      load: async () => {
        const response = await listStockMovements()
        stockMovements.value = Array.isArray(response.data) ? response.data : []
      },
    },
  ]

  try {
    const results = await Promise.allSettled(loaders.map((loader) => loader.load()))
    const failedLoaders = results
      .map((result, index) => (result.status === 'rejected' ? loaders[index].name : null))
      .filter(Boolean)

    if (failedLoaders.length) {
      errorMessage.value = `Sebagian data laporan gagal dimuat: ${failedLoaders.join(', ')}.`
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchReportsData()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <div>
        <h2 class="text-2xl font-semibold text-gray-900">Laporan & dokumen</h2>
        <p class="mt-1 text-gray-600">Buat dan unduh laporan dari data operasional yang sudah tercatat.</p>
      </div>
      <button
        type="button"
        class="inline-flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-green-700 disabled:opacity-60"
        :disabled="loading"
        @click="fetchReportsData"
      >
        <Loader2 v-if="loading" class="h-4 w-4 animate-spin" />
        <Download v-else class="h-4 w-4" />
        Refresh data
      </button>
    </div>

    <div v-if="errorMessage" class="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
      {{ errorMessage }}
    </div>
    <div v-if="successMessage" class="rounded-xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
      {{ successMessage }}
    </div>

    <div class="grid grid-cols-1 gap-4 md:grid-cols-4">
      <div class="rounded-xl border border-emerald-100 bg-white p-5">
        <div class="flex items-center justify-between">
          <p class="text-sm text-gray-600">Total stok gudang</p>
          <Package class="h-5 w-5 text-emerald-600" />
        </div>
        <p class="mt-2 text-2xl font-semibold text-gray-900">{{ formatNumber(totalWarehouseStock) }} kg</p>
      </div>
      <div class="rounded-xl border border-purple-100 bg-white p-5">
        <div class="flex items-center justify-between">
          <p class="text-sm text-gray-600">Total panen</p>
          <Sprout class="h-5 w-5 text-purple-600" />
        </div>
        <p class="mt-2 text-2xl font-semibold text-gray-900">{{ formatNumber(totalHarvest) }} kg</p>
      </div>
      <div class="rounded-xl border border-blue-100 bg-white p-5">
        <div class="flex items-center justify-between">
          <p class="text-sm text-gray-600">Pendapatan</p>
          <TrendingUp class="h-5 w-5 text-blue-600" />
        </div>
        <p class="mt-2 text-2xl font-semibold text-gray-900">{{ formatRupiah(totalRevenue) }}</p>
      </div>
      <div class="rounded-xl border border-amber-100 bg-white p-5">
        <div class="flex items-center justify-between">
          <p class="text-sm text-gray-600">Biaya operasional</p>
          <Wallet class="h-5 w-5 text-amber-600" />
        </div>
        <p class="mt-2 text-2xl font-semibold text-gray-900">{{ formatRupiah(totalCost) }}</p>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-4 md:grid-cols-4">
      <button
        v-for="card in reportCards"
        :key="card.id"
        type="button"
        class="rounded-xl border border-gray-200 bg-white p-6 text-left transition-all hover:border-green-600 hover:shadow-md"
        @click="downloadReport(reports.find((report) => report.category === card.id)!)"
      >
        <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg" :class="card.color">
          <component :is="card.icon" class="h-6 w-6" />
        </div>
        <h3 class="mb-1 font-semibold text-gray-900">{{ card.title }}</h3>
        <p class="text-sm text-gray-600">{{ card.description }}</p>
      </button>
    </div>

    <div class="rounded-xl border border-gray-200 bg-white p-4">
      <div class="flex flex-wrap items-center gap-4">
        <div class="flex items-center gap-2">
          <Filter class="h-4 w-4 text-gray-600" />
          <span class="text-sm font-medium text-gray-700">Filter:</span>
        </div>

        <div class="flex flex-wrap gap-2">
          <button
            v-for="type in reportTypes"
            :key="type.id"
            type="button"
            :class="[
              'rounded-lg px-4 py-2 text-sm font-medium transition-colors',
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
            class="rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-green-500"
          >
            <option value="all">Sepanjang waktu</option>
            <option value="monthly">Bulan ini</option>
            <option value="quarterly">Kuartal ini</option>
            <option value="yearly">Tahun ini</option>
          </select>
        </div>
      </div>
    </div>

    <div class="overflow-hidden rounded-xl border border-gray-200 bg-white">
      <div class="flex flex-wrap items-center justify-between gap-3 border-b border-gray-200 px-6 py-4">
        <h3 class="text-lg font-semibold text-gray-900">Daftar laporan</h3>
        <div class="w-full max-w-sm">
          <SearchBar
            v-model="searchQuery"
            placeholder="Cari nama, periode..."
            @search="handleSearch"
          />
        </div>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full min-w-[840px]">
          <thead class="border-b border-gray-200 bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Nama laporan</th>
              <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Jenis</th>
              <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Periode</th>
              <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Tanggal</th>
              <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Ukuran</th>
              <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-if="!paginatedReports.length">
              <td colspan="6" class="px-6 py-6 text-center text-sm text-gray-500">Tidak ada laporan untuk filter ini.</td>
            </tr>
            <tr v-for="report in paginatedReports" :key="report.id" class="hover:bg-gray-50">
              <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                  <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-red-50">
                    <FileText class="h-5 w-5 text-red-600" />
                  </div>
                  <div>
                    <p class="font-medium text-gray-900">{{ report.name }}</p>
                    <p class="text-sm text-gray-500">{{ report.format }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4">
                <span
                  class="inline-flex rounded-full px-2.5 py-1 text-xs font-medium"
                  :class="reportTypeBadgeClass(report.category)"
                >
                  {{ categoryLabels[report.category] }}
                </span>
              </td>
              <td class="px-6 py-4 text-gray-900">{{ report.period }}</td>
              <td class="px-6 py-4 text-gray-600">{{ report.date }}</td>
              <td class="px-6 py-4 text-gray-600">{{ report.size }}</td>
              <td class="px-6 py-4">
                <button
                  type="button"
                  class="flex items-center gap-2 rounded-lg px-3 py-1.5 text-sm font-medium text-green-600 transition-colors hover:bg-green-50"
                  @click="downloadReport(report)"
                >
                  <Download class="h-4 w-4" />
                  Unduh
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <Pagination
        v-if="filteredReports.length > 0"
        :current-page="currentPage"
        :total-pages="totalPages"
        :total-rows="filteredReports.length"
        :limit="limit"
        @update:page="currentPage = $event"
      />
    </div>

    <div class="rounded-xl border border-gray-200 bg-gray-50 p-6">
      <h3 class="mb-4 font-semibold text-gray-900">Buat laporan kustom</h3>
      <div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
        <label class="rounded-lg border border-gray-200 bg-white p-4">
          <span class="mb-2 block text-sm font-medium text-gray-900">Tanggal mulai</span>
          <input
            v-model="customForm.startDate"
            type="date"
            class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm"
          />
        </label>
        <label class="rounded-lg border border-gray-200 bg-white p-4">
          <span class="mb-2 block text-sm font-medium text-gray-900">Tanggal selesai</span>
          <input
            v-model="customForm.endDate"
            type="date"
            class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm"
          />
        </label>
        <label class="rounded-lg border border-gray-200 bg-white p-4">
          <span class="mb-2 block text-sm font-medium text-gray-900">Jenis laporan</span>
          <select v-model="customForm.reportType" class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm">
            <option value="all">Semua jenis</option>
            <option value="financial">Keuangan</option>
            <option value="investment">Investasi</option>
            <option value="plant">Panen & tanaman</option>
            <option value="warehouse">Stok gudang</option>
          </select>
        </label>
        <div class="flex items-end">
          <button
            type="button"
            class="w-full rounded-lg bg-green-600 px-4 py-2 font-medium text-white transition-colors hover:bg-green-700"
            @click="createCustomReport"
          >
            Buat laporan
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
