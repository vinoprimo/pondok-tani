<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import {
  ArrowRight,
  BarChart3,
  Check,
  Droplets,
  Loader2,
  Package,
  Plus,
  TrendingDown,
  TrendingUp,
  Wind,
  X,
} from 'lucide-vue-next'
import { getWarehouseSummary, listStockMovements, listWarehouseStocks } from '../../services/warehouse/warehouse'
import { listDryingProcesses, completeDryingProcess } from '../../services/postharvest/drying'
import { listGradingBatches, createGradingBatch } from '../../services/postharvest/grading'
import { getGrades } from '../../services/price/vanili'

const summary = ref<any>({
  total_quantity: 0,
  total_batches: 0,
  total_grades: 0,
  grade_breakdown: [],
  monthly_movements: [],
})
const stocks = ref<any[]>([])
const movements = ref<any[]>([])
const dryingProcesses = ref<any[]>([])
const gradingBatches = ref<any[]>([])
const grades = ref<any[]>([])

const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const showCompleteModal = ref(false)
const selectedDrying = ref<any | null>(null)
const completing = ref(false)
const dryingForm = ref({
  endDate: '',
  finalQuantity: 0,
  notes: '',
})

const showGradingModal = ref(false)
const creatingGrading = ref(false)
const gradingForm = ref({
  dryingProcessId: 0,
  gradingDate: '',
  details: [{ gradeId: '', quantity: 0 }],
})

const showOriginModal = ref(false)
const selectedGradeName = ref('')
const originPage = ref(1)
const ORIGIN_PAGE_SIZE = 6

const ongoingDrying = computed(() =>
  dryingProcesses.value.filter((item) => item.status === 'ongoing')
)
const readyForGrading = computed(() =>
  dryingProcesses.value.filter((item) => item.status === 'completed' && !item.has_grading)
)

const gradeBreakdown = computed(() => summary.value.grade_breakdown || [])
const monthlyMovements = computed(() => summary.value.monthly_movements || [])
const currentMonthKey = computed(() => new Date().toLocaleDateString('en-US', { month: 'short' }))
const currentMonthStats = computed(() => {
  const match = monthlyMovements.value.find((item: any) => item.month === currentMonthKey.value)
  return {
    stockIn: Number(match?.stock_in || 0),
    stockOut: Number(match?.stock_out || 0),
  }
})

const originStocks = computed(() => {
  if (!selectedGradeName.value) return []
  return stocks.value.filter((item) => String(item.grade_name) === selectedGradeName.value)
})

const originTotalPages = computed(() =>
  Math.max(1, Math.ceil(originStocks.value.length / ORIGIN_PAGE_SIZE))
)

const paginatedOrigins = computed(() => {
  const start = (originPage.value - 1) * ORIGIN_PAGE_SIZE
  return originStocks.value.slice(start, start + ORIGIN_PAGE_SIZE)
})

const chartMax = computed(() => {
  if (!monthlyMovements.value.length) return 1
  return Math.max(
    ...monthlyMovements.value.flatMap((m: any) => [m.stock_in || 0, m.stock_out || 0]),
    1
  )
})

function barPct(value: number) {
  return `${(value / chartMax.value) * 100}%`
}

function formatDate(value: string | Date) {
  const date = value instanceof Date ? value : new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  })
}

async function fetchAll() {
  loading.value = true
  errorMessage.value = ''
  try {
    const [summaryRes, stocksRes, movementRes, dryingRes, gradingRes, gradesRes] = await Promise.all([
      getWarehouseSummary(),
      listWarehouseStocks(),
      listStockMovements(),
      listDryingProcesses(),
      listGradingBatches(),
      getGrades(),
    ])
    summary.value = summaryRes?.data || summary.value
    stocks.value = Array.isArray(stocksRes.data) ? stocksRes.data : []
    movements.value = Array.isArray(movementRes.data) ? movementRes.data : []
    dryingProcesses.value = Array.isArray(dryingRes.data) ? dryingRes.data : []
    gradingBatches.value = Array.isArray(gradingRes.data) ? gradingRes.data : []
    grades.value = Array.isArray(gradesRes.data) ? gradesRes.data : []
  } catch (error) {
    console.error('Failed to load warehouse data:', error)
    errorMessage.value = 'Gagal memuat data stok gudang.'
  } finally {
    loading.value = false
  }
}

function openCompleteModal(item: any) {
  selectedDrying.value = item
  dryingForm.value.endDate = new Date().toISOString().slice(0, 10)
  dryingForm.value.finalQuantity = item.initial_quantity || 0
  dryingForm.value.notes = ''
  showCompleteModal.value = true
}

async function confirmCompleteDrying() {
  if (!selectedDrying.value) return
  if (dryingForm.value.finalQuantity <= 0) {
    errorMessage.value = 'Kuantitas akhir wajib diisi.'
    return
  }

  completing.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    await completeDryingProcess(selectedDrying.value.id, {
      end_date: dryingForm.value.endDate,
      final_quantity: Number(dryingForm.value.finalQuantity),
      notes: dryingForm.value.notes ? dryingForm.value.notes.trim() : null,
    })
    successMessage.value = 'Proses pengeringan berhasil diselesaikan.'
    showCompleteModal.value = false
    await fetchAll()
  } catch (error: any) {
    console.error('Failed to complete drying process:', error)
    errorMessage.value = error?.response?.data?.error || 'Gagal menyelesaikan pengeringan.'
  } finally {
    completing.value = false
  }
}

function openGradingModal(item: any) {
  gradingForm.value.dryingProcessId = item.id
  gradingForm.value.gradingDate = new Date().toISOString().slice(0, 10)
  gradingForm.value.details = [{ gradeId: '', quantity: 0 }]
  showGradingModal.value = true
}

function addGradingRow() {
  gradingForm.value.details.push({ gradeId: '', quantity: 0 })
}

function removeGradingRow(index: number) {
  gradingForm.value.details.splice(index, 1)
}

function openOriginModal(gradeName: string) {
  selectedGradeName.value = gradeName
  originPage.value = 1
  showOriginModal.value = true
}

async function confirmGrading() {
  const detailPayload = gradingForm.value.details
    .filter((item) => item.gradeId && Number(item.quantity) > 0)
    .map((item) => ({
      grade_id: Number(item.gradeId),
      quantity: Number(item.quantity),
    }))

  if (!gradingForm.value.dryingProcessId || detailPayload.length === 0) {
    errorMessage.value = 'Detail grading wajib diisi.'
    return
  }

  creatingGrading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    await createGradingBatch({
      drying_process_id: gradingForm.value.dryingProcessId,
      grading_date: gradingForm.value.gradingDate,
      details: detailPayload,
    })
    successMessage.value = 'Grading berhasil disimpan.'
    showGradingModal.value = false
    await fetchAll()
  } catch (error: any) {
    console.error('Failed to create grading batch:', error)
    errorMessage.value = error?.response?.data?.error || 'Gagal menyimpan grading.'
  } finally {
    creatingGrading.value = false
  }
}

onMounted(() => {
  fetchAll()
})
</script>

<template>
  <div class="warehouse-lab space-y-6">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <div>
        <h2 class="text-2xl font-semibold text-gray-900">Stok Gudang</h2>
        <p class="mt-1 text-gray-600">Pantau level stok vanili, pengeringan, dan grading.</p>
      </div>
      <div class="rounded-xl border border-emerald-100 bg-emerald-50 px-4 py-2 text-xs text-emerald-700">
        Update terakhir: {{ new Date().toLocaleDateString('id-ID') }}
      </div>
    </div>

    <div v-if="errorMessage" class="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
      {{ errorMessage }}
    </div>
    <div v-if="successMessage" class="rounded-xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
      {{ successMessage }}
    </div>

    <div class="grid grid-cols-1 gap-4 md:grid-cols-4">
      <div class="rounded-2xl border border-emerald-100 bg-gradient-to-br from-emerald-50 to-white p-6">
        <div class="flex items-center justify-between">
          <p class="text-sm text-gray-600">Total stok</p>
          <Package class="h-5 w-5 text-emerald-600" />
        </div>
        <p class="mt-2 text-2xl font-semibold text-gray-900">{{ Number(summary.total_quantity || 0).toFixed(1) }} kg</p>
        <p class="text-xs text-gray-500">Semua mutu digabung</p>
      </div>
      <div class="rounded-2xl border border-indigo-100 bg-gradient-to-br from-indigo-50 to-white p-6">
        <div class="flex items-center justify-between">
          <p class="text-sm text-gray-600">Stok masuk (bulan ini)</p>
          <TrendingUp class="h-5 w-5 text-indigo-600" />
        </div>
        <p class="mt-2 text-2xl font-semibold text-gray-900">{{ currentMonthStats.stockIn.toFixed(1) }} kg</p>
        <p class="text-xs text-gray-500">Total masuk bulan berjalan</p>
      </div>
      <div class="rounded-2xl border border-sky-100 bg-gradient-to-br from-sky-50 to-white p-6">
        <div class="flex items-center justify-between">
          <p class="text-sm text-gray-600">Stok keluar (bulan ini)</p>
          <TrendingDown class="h-5 w-5 text-sky-600" />
        </div>
        <p class="mt-2 text-2xl font-semibold text-gray-900">{{ currentMonthStats.stockOut.toFixed(1) }} kg</p>
        <p class="text-xs text-gray-500">Total keluar bulan berjalan</p>
      </div>
      <div class="rounded-2xl border border-amber-100 bg-gradient-to-br from-amber-50 to-white p-6">
        <div class="flex items-center justify-between">
          <p class="text-sm text-gray-600">Pengeringan aktif</p>
          <Wind class="h-5 w-5 text-amber-600" />
        </div>
        <p class="mt-2 text-2xl font-semibold text-gray-900">{{ ongoingDrying.length }}</p>
        <p class="text-xs text-gray-500">Batch kering berjalan</p>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
      <div class="rounded-2xl border border-gray-200 bg-white p-6">
        <h3 class="text-lg font-semibold text-gray-900">Stok per mutu</h3>
        <div v-if="loading" class="mt-4 flex items-center gap-2 text-sm text-gray-500">
          <Loader2 class="h-4 w-4 animate-spin" />
          Memuat stok...
        </div>
        <div v-else class="mt-4 space-y-3">
          <button
            v-for="item in gradeBreakdown"
            :key="item.grade_id || item.grade_name"
            type="button"
            class="w-full rounded-xl border border-emerald-100 bg-white px-4 py-3 text-left transition hover:border-emerald-300 hover:bg-emerald-50"
            @click="openOriginModal(item.grade_name)"
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-3">
                <span class="font-medium text-gray-900">{{ item.grade_name }}</span>
                <div class="relative h-2 w-48 overflow-hidden rounded-full bg-gray-100">
                  <div
                    class="absolute inset-y-0 left-0 rounded-full bg-gradient-to-r from-emerald-400 to-emerald-600"
                    :style="{ width: `${Math.min((Number(item.quantity || 0) / Math.max(summary.total_quantity || 1, 1)) * 100, 100)}%` }"
                  />
                </div>
              </div>
              <div class="text-right">
                <p class="font-semibold text-gray-900">{{ Number(item.quantity || 0).toFixed(1) }} kg</p>
                <p class="text-xs text-gray-500">{{ ((Number(item.quantity || 0) / Math.max(summary.total_quantity || 1, 1)) * 100).toFixed(0) }}%</p>
              </div>
            </div>
          </button>
        </div>
      </div>

      <div class="rounded-2xl border border-gray-200 bg-white p-6">
        <h3 class="text-lg font-semibold text-gray-900">Tren pergerakan stok</h3>
        <p class="mt-1 text-xs text-gray-500">Masuk (hijau) vs keluar (amber) per bulan.</p>
        <div v-if="loading" class="mt-4 flex items-center gap-2 text-sm text-gray-500">
          <Loader2 class="h-4 w-4 animate-spin" />
          Memuat grafik...
        </div>
        <div v-else class="mt-4 h-[220px]">
          <div class="flex h-full items-end gap-2">
            <div v-for="item in monthlyMovements" :key="item.month" class="flex flex-1 flex-col items-center">
              <div class="flex h-40 items-end gap-1">
                <div
                  class="w-4 rounded-t-lg bg-emerald-500"
                  :style="{ height: barPct(Number(item.stock_in || 0)) }"
                />
                <div
                  class="w-4 rounded-t-lg bg-amber-400"
                  :style="{ height: barPct(Number(item.stock_out || 0)) }"
                />
              </div>
              <span class="mt-2 text-xs text-gray-500">{{ item.month }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
      <div class="rounded-2xl border border-gray-200 bg-white p-6">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Pengeringan aktif</h3>
            <p class="text-sm text-gray-500">Batch kering yang sedang diproses.</p>
          </div>
          <Wind class="h-5 w-5 text-amber-500" />
        </div>
        <div class="mt-4 overflow-x-auto">
          <table class="w-full text-sm">
            <thead class="bg-gray-50 text-gray-700">
              <tr>
                <th class="px-4 py-2 text-left font-medium">Batch</th>
                <th class="px-4 py-2 text-left font-medium">Mulai</th>
                <th class="px-4 py-2 text-right font-medium">Qty awal</th>
                <th class="px-4 py-2 text-center font-medium">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-if="!ongoingDrying.length">
                <td colspan="4" class="px-4 py-4 text-center text-gray-500">Tidak ada proses pengeringan aktif.</td>
              </tr>
              <tr v-for="item in ongoingDrying" :key="item.id">
                <td class="px-4 py-3">
                  <p class="font-medium text-gray-900">{{ item.batch_code }}</p>
                  <p class="text-xs text-gray-500">{{ item.package_name }}</p>
                </td>
                <td class="px-4 py-3 text-gray-600">{{ formatDate(item.start_date) }}</td>
                <td class="px-4 py-3 text-right text-gray-700">{{ Number(item.initial_quantity || 0).toFixed(1) }} kg</td>
                <td class="px-4 py-3 text-center">
                  <button
                    type="button"
                    class="inline-flex items-center gap-2 rounded-lg bg-amber-500 px-3 py-2 text-xs font-semibold text-white hover:bg-amber-600"
                    @click="openCompleteModal(item)"
                  >
                    <Check class="h-4 w-4" />
                    Selesaikan
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="rounded-2xl border border-gray-200 bg-white p-6">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Siap grading</h3>
            <p class="text-sm text-gray-500">Batch yang sudah selesai dikeringkan.</p>
          </div>
          <Droplets class="h-5 w-5 text-sky-500" />
        </div>
        <div class="mt-4 overflow-x-auto">
          <table class="w-full text-sm">
            <thead class="bg-gray-50 text-gray-700">
              <tr>
                <th class="px-4 py-2 text-left font-medium">Batch</th>
                <th class="px-4 py-2 text-left font-medium">Selesai</th>
                <th class="px-4 py-2 text-right font-medium">Qty akhir</th>
                <th class="px-4 py-2 text-center font-medium">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-if="!readyForGrading.length">
                <td colspan="4" class="px-4 py-4 text-center text-gray-500">Belum ada batch siap grading.</td>
              </tr>
              <tr v-for="item in readyForGrading" :key="item.id">
                <td class="px-4 py-3">
                  <p class="font-medium text-gray-900">{{ item.batch_code }}</p>
                  <p class="text-xs text-gray-500">{{ item.package_name }}</p>
                </td>
                <td class="px-4 py-3 text-gray-600">{{ item.end_date ? formatDate(item.end_date) : '-' }}</td>
                <td class="px-4 py-3 text-right text-gray-700">{{ Number(item.final_quantity || 0).toFixed(1) }} kg</td>
                <td class="px-4 py-3 text-center">
                  <button
                    type="button"
                    class="inline-flex items-center gap-2 rounded-lg bg-emerald-600 px-3 py-2 text-xs font-semibold text-white hover:bg-emerald-700"
                    @click="openGradingModal(item)"
                  >
                    <Plus class="h-4 w-4" />
                    Input grading
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <div class="rounded-2xl border border-gray-200 bg-white p-6">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-lg font-semibold text-gray-900">Riwayat grading</h3>
          <p class="text-sm text-gray-500">Hasil grading yang sudah masuk gudang.</p>
        </div>
      </div>
      <div class="mt-4 overflow-x-auto">
        <table class="w-full text-sm">
          <thead class="bg-gray-50 text-gray-700">
            <tr>
              <th class="px-4 py-2 text-left font-medium">Batch</th>
              <th class="px-4 py-2 text-left font-medium">Tanggal grading</th>
              <th class="px-4 py-2 text-left font-medium">Total</th>
              <th class="px-4 py-2 text-left font-medium">Rincian grade</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-if="!gradingBatches.length">
              <td colspan="4" class="px-4 py-4 text-center text-gray-500">Belum ada data grading.</td>
            </tr>
            <tr v-for="item in gradingBatches" :key="item.id">
              <td class="px-4 py-3">
                <p class="font-medium text-gray-900">{{ item.batch_code_plant }}</p>
                <p class="text-xs text-gray-500">{{ item.package_name }}</p>
              </td>
              <td class="px-4 py-3 text-gray-600">{{ formatDate(item.grading_date) }}</td>
              <td class="px-4 py-3 text-gray-700">{{ Number(item.total_quantity || 0).toFixed(1) }} kg</td>
              <td class="px-4 py-3">
                <div class="flex flex-wrap gap-2">
                  <span
                    v-for="detail in item.details"
                    :key="detail.grade_id"
                    class="rounded-full bg-emerald-50 px-2 py-0.5 text-xs text-emerald-700"
                  >
                    {{ detail.grade_name }}: {{ Number(detail.quantity || 0).toFixed(1) }} kg
                  </span>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="rounded-2xl border border-gray-200 bg-white p-6">
      <div class="border-b border-gray-200 pb-3">
        <h3 class="text-lg font-semibold text-gray-900">Riwayat pergerakan stok</h3>
        <p class="text-sm text-gray-500">Transaksi masuk dan keluar gudang terbaru.</p>
      </div>
      <div class="mt-4 overflow-x-auto">
        <table class="w-full text-sm">
          <thead class="bg-gray-50 text-gray-700">
            <tr>
              <th class="px-4 py-2 text-left font-medium">ID pergerakan</th>
              <th class="px-4 py-2 text-left font-medium">Sumber</th>
              <th class="px-4 py-2 text-left font-medium">Mutu</th>
              <th class="px-4 py-2 text-center font-medium">Masuk</th>
              <th class="px-4 py-2 text-center font-medium">Keluar</th>
              <th class="px-4 py-2 text-left font-medium">Tanggal</th>
              <th class="px-4 py-2 text-center font-medium">Referensi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-if="!movements.length">
              <td colspan="7" class="px-4 py-4 text-center text-gray-500">Belum ada pergerakan stok.</td>
            </tr>
            <tr v-for="movement in movements" :key="movement.id" class="hover:bg-gray-50">
              <td class="px-4 py-3">
                <div class="flex items-center gap-2">
                  <TrendingUp v-if="movement.movement_type === 'in'" class="h-4 w-4 text-emerald-600" />
                  <TrendingDown v-else class="h-4 w-4 text-amber-600" />
                  <span class="font-medium text-gray-900">{{ movement.id }}</span>
                </div>
              </td>
              <td class="px-4 py-3 text-gray-700">{{ movement.reference_type }}</td>
              <td class="px-4 py-3">
                <span class="rounded-full bg-blue-100 px-2 py-0.5 text-xs font-medium text-blue-700">
                  {{ movement.grade_name }}
                </span>
              </td>
              <td class="px-4 py-3 text-center">
                <span v-if="movement.movement_type === 'in'" class="font-semibold text-emerald-600">
                  +{{ Number(movement.quantity || 0).toFixed(1) }} kg
                </span>
                <span v-else class="text-gray-400">-</span>
              </td>
              <td class="px-4 py-3 text-center">
                <span v-if="movement.movement_type === 'out'" class="font-semibold text-amber-600">
                  -{{ Number(movement.quantity || 0).toFixed(1) }} kg
                </span>
                <span v-else class="text-gray-400">-</span>
              </td>
              <td class="px-4 py-3 text-gray-600">{{ formatDate(movement.movement_date) }}</td>
              <td class="px-4 py-3 text-center">
                <button type="button" class="text-emerald-600 hover:text-emerald-700 text-xs font-semibold">
                  {{ movement.reference_id }}
                  <ArrowRight class="ml-1 inline h-3 w-3" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <Pagination
        v-if="filteredMovements.length > 0"
        :current-page="currentPage"
        :total-pages="totalPages"
        :total-rows="filteredMovements.length"
        :limit="limit"
        @update:page="currentPage = $event"
      />
    </div>

    <div
      v-if="showOriginModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
    >
      <div class="w-full max-w-4xl rounded-2xl bg-white">
        <div class="border-b border-gray-200 px-6 py-4">
          <h3 class="text-lg font-semibold text-gray-900">Detail stok mutu {{ selectedGradeName }}</h3>
          <p class="mt-1 text-sm text-gray-600">Asal stok berdasarkan harvest batch dan pemiliknya.</p>
        </div>
        <div class="p-6">
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead class="bg-gray-50 text-gray-700">
                <tr>
                  <th class="px-4 py-2 text-left font-medium">Harvest batch</th>
                  <th class="px-4 py-2 text-left font-medium">Plant batch</th>
                  <th class="px-4 py-2 text-left font-medium">User</th>
                  <th class="px-4 py-2 text-left font-medium">Paket</th>
                  <th class="px-4 py-2 text-right font-medium">Qty (kg)</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                <tr v-if="!originStocks.length">
                  <td colspan="5" class="px-4 py-4 text-center text-gray-500">Belum ada stok untuk mutu ini.</td>
                </tr>
                <tr v-for="item in paginatedOrigins" :key="item.id">
                  <td class="px-4 py-3">
                    <p class="font-medium text-gray-900">{{ item.batch_code }}</p>
                  </td>
                  <td class="px-4 py-3 text-gray-700">#{{ item.plant_batch_id }}</td>
                  <td class="px-4 py-3 text-gray-700">{{ item.user_name || item.user_id || '-' }}</td>
                  <td class="px-4 py-3 text-gray-600">{{ item.package_name }}</td>
                  <td class="px-4 py-3 text-right font-semibold text-gray-900">{{ Number(item.total_quantity || 0).toFixed(1) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="mt-4 flex items-center justify-between">
            <p class="text-xs text-gray-500">Menampilkan {{ paginatedOrigins.length }} dari {{ originStocks.length }} data</p>
            <div class="flex items-center gap-2">
              <button
                type="button"
                class="rounded-lg border border-gray-300 px-3 py-1.5 text-xs text-gray-700 disabled:opacity-50"
                :disabled="originPage <= 1"
                @click="originPage -= 1"
              >
                Sebelumnya
              </button>
              <span class="text-xs text-gray-600">Halaman {{ originPage }} / {{ originTotalPages }}</span>
              <button
                type="button"
                class="rounded-lg border border-gray-300 px-3 py-1.5 text-xs text-gray-700 disabled:opacity-50"
                :disabled="originPage >= originTotalPages"
                @click="originPage += 1"
              >
                Berikutnya
              </button>
            </div>
          </div>
        </div>
        <div class="flex justify-end border-t border-gray-200 px-6 py-4">
          <button
            type="button"
            class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50"
            @click="showOriginModal = false"
          >
            <X class="mr-1 inline h-4 w-4" />
            Tutup
          </button>
        </div>
      </div>
    </div>

    <div
      v-if="showCompleteModal && selectedDrying"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
    >
      <div class="w-full max-w-lg rounded-2xl bg-white">
        <div class="border-b border-gray-200 px-6 py-4">
          <h3 class="text-lg font-semibold text-gray-900">Selesaikan pengeringan</h3>
          <p class="mt-1 text-sm text-gray-600">Input hasil akhir pengeringan.</p>
        </div>
        <div class="space-y-4 p-6 text-sm text-gray-700">
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">Tanggal selesai</label>
            <input
              v-model="dryingForm.endDate"
              type="date"
              class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-amber-500 focus:outline-none focus:ring-2 focus:ring-amber-100"
            />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">Kuantitas akhir (kg)</label>
            <input
              v-model.number="dryingForm.finalQuantity"
              type="number"
              min="0"
              step="0.1"
              class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-amber-500 focus:outline-none focus:ring-2 focus:ring-amber-100"
            />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">Catatan</label>
            <textarea
              v-model="dryingForm.notes"
              rows="3"
              class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-amber-500 focus:outline-none focus:ring-2 focus:ring-amber-100"
            />
          </div>
        </div>
        <div class="flex gap-3 border-t border-gray-200 px-6 py-4">
          <button
            type="button"
            class="flex-1 rounded-lg border border-gray-300 px-4 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50"
            @click="showCompleteModal = false"
          >
            <X class="mr-1 inline h-4 w-4" />
            Batal
          </button>
          <button
            type="button"
            class="flex-1 rounded-lg bg-amber-600 px-4 py-2 text-sm font-semibold text-white hover:bg-amber-700 disabled:opacity-60"
            :disabled="completing"
            @click="confirmCompleteDrying"
          >
            <Loader2 v-if="completing" class="mr-1 inline h-4 w-4 animate-spin" />
            <Check v-else class="mr-1 inline h-4 w-4" />
            Simpan
          </button>
        </div>
      </div>
    </div>

    <div
      v-if="showGradingModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
    >
      <div class="w-full max-w-2xl rounded-2xl bg-white">
        <div class="border-b border-gray-200 px-6 py-4">
          <h3 class="text-lg font-semibold text-gray-900">Input grading</h3>
          <p class="mt-1 text-sm text-gray-600">Bagi hasil vanili kering berdasarkan grade.</p>
        </div>
        <div class="space-y-4 p-6 text-sm text-gray-700">
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">Tanggal grading</label>
            <input
              v-model="gradingForm.gradingDate"
              type="date"
              class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-100"
            />
          </div>
          <div class="space-y-3">
            <div
              v-for="(detail, index) in gradingForm.details"
              :key="index"
              class="grid grid-cols-1 gap-3 rounded-xl border border-gray-200 p-3 md:grid-cols-[1.2fr_1fr_auto]"
            >
              <select
                v-model="detail.gradeId"
                class="rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-100"
              >
                <option value="">Pilih grade</option>
                <option v-for="grade in grades" :key="grade.id" :value="String(grade.id)">
                  {{ grade.grade_name }}
                </option>
              </select>
              <input
                v-model.number="detail.quantity"
                type="number"
                min="0"
                step="0.1"
                placeholder="Qty (kg)"
                class="rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-100"
              />
              <button
                type="button"
                class="rounded-lg border border-gray-300 px-3 py-2 text-xs font-semibold text-gray-600 hover:bg-gray-50"
                @click="removeGradingRow(index)"
              >
                Hapus
              </button>
            </div>
            <button
              type="button"
              class="inline-flex items-center gap-2 rounded-lg border border-emerald-500 px-3 py-2 text-xs font-semibold text-emerald-700 hover:bg-emerald-50"
              @click="addGradingRow"
            >
              <Plus class="h-4 w-4" />
              Tambah grade
            </button>
          </div>
        </div>
        <div class="flex gap-3 border-t border-gray-200 px-6 py-4">
          <button
            type="button"
            class="flex-1 rounded-lg border border-gray-300 px-4 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50"
            @click="showGradingModal = false"
          >
            <X class="mr-1 inline h-4 w-4" />
            Batal
          </button>
          <button
            type="button"
            class="flex-1 rounded-lg bg-emerald-600 px-4 py-2 text-sm font-semibold text-white hover:bg-emerald-700 disabled:opacity-60"
            :disabled="creatingGrading"
            @click="confirmGrading"
          >
            <Loader2 v-if="creatingGrading" class="mr-1 inline h-4 w-4 animate-spin" />
            <Check v-else class="mr-1 inline h-4 w-4" />
            Simpan
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

