<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { BarChart3, Check, ClipboardCheck, Eye, Loader2, MessageCircle, Search, X } from 'lucide-vue-next'
import {
  createHarvestRequest,
  listHarvestRequests,
  listHarvests,
  getHarvestSummary,
  validateHarvestRequest,
} from '../../services/harvest/harvest'
import { listPlantMonitorings } from '../../services/plant-monitoring/monitoring'
import Pagination from '../../components/common/Pagination.vue'

const storedRole = localStorage.getItem('userRole') as 'investor' | 'mitra' | 'admin' | null
const userRole = ref<'investor' | 'mitra' | 'admin'>(storedRole || 'investor')

const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const searchQuery = ref('')
const summaryLoading = ref(false)
const summaryError = ref('')
const harvestSummary = ref<any[]>([])
const batchOptions = ref<any[]>([])
const selectedBatchFilter = ref('')
const dateFilters = ref({
  start: '',
  end: '',
})

const harvestRequests = ref<any[]>([])
const harvests = ref<any[]>([])
const panenMonitorings = ref<any[]>([])

const showRequestModal = ref(false)
const selectedMonitoringId = ref<number | null>(null)
const requesting = ref(false)

const showValidateModal = ref(false)
const selectedRequest = ref<any | null>(null)
const validating = ref(false)
const validationForm = ref({
  harvestDate: '',
  notes: '',
  totalQuantity: 0,
  wetQuantity: 0,
  dryQuantity: 0,
})
const lastEditedHarvestType = ref<'wet' | 'dry' | null>(null)

const showDetailModal = ref(false)
const selectedDetail = ref<any | null>(null)

const isAdmin = computed(() => userRole.value === 'admin')

function numericValue(value: unknown) {
  const number = Number(value)
  return Number.isFinite(number) ? number : 0
}

function roundQuantity(value: number) {
  return Math.round(value * 10) / 10
}

function completeCounterpartQuantity(changedField: 'total' | 'wet' | 'dry') {
  const totalQty = numericValue(validationForm.value.totalQuantity)
  const wetQty = numericValue(validationForm.value.wetQuantity)
  const dryQty = numericValue(validationForm.value.dryQuantity)

  if (changedField === 'wet') {
    lastEditedHarvestType.value = 'wet'
    validationForm.value.dryQuantity = roundQuantity(Math.max(totalQty - wetQty, 0))
    return
  }

  if (changedField === 'dry') {
    lastEditedHarvestType.value = 'dry'
    validationForm.value.wetQuantity = roundQuantity(Math.max(totalQty - dryQty, 0))
    return
  }

  if (lastEditedHarvestType.value === 'wet') {
    validationForm.value.dryQuantity = roundQuantity(Math.max(totalQty - wetQty, 0))
  } else if (lastEditedHarvestType.value === 'dry') {
    validationForm.value.wetQuantity = roundQuantity(Math.max(totalQty - dryQty, 0))
  }
}

const isValidationFormComplete = computed(() => {
  const totalQty = numericValue(validationForm.value.totalQuantity)
  const wetQty = numericValue(validationForm.value.wetQuantity)
  const dryQty = numericValue(validationForm.value.dryQuantity)

  return Boolean(
    selectedRequest.value &&
      validationForm.value.harvestDate &&
      totalQty > 0 &&
      wetQty >= 0 &&
      dryQty >= 0 &&
      wetQty + dryQty > 0 &&
      Math.abs(totalQty - (wetQty + dryQty)) <= 0.01
  )
})

const pendingRequests = computed(() => harvestRequests.value.filter((item) => item.status === 'pending'))
const validatedRequests = computed(() => harvestRequests.value.filter((item) => item.status === 'validated'))

const latestHarvestDate = computed(() => {
  if (!harvests.value.length) return '-'
  const latest = harvests.value[0]
  return formatDate(latest.harvest_date || latest.harvestDate)
})

const summaryCards = computed(() => [
  { label: 'Total ajuan', value: harvestRequests.value.length },
  { label: 'Menunggu validasi', value: pendingRequests.value.length },
  { label: 'Tervalidasi', value: validatedRequests.value.length },
  { label: 'Total panen', value: harvests.value.length },
])

const phaseLabelMap: Record<string, string> = {
  penanaman: 'Penanaman',
  pertumbuhan_awal: 'Pertumbuhan Awal',
  vegetatif: 'Vegetatif',
  'pra-berbunga': 'Pra-berbunga',
  berbunga: 'Berbunga',
  panen: 'Panen',
}

const healthLabelMap: Record<string, string> = {
  sehat: 'Sehat',
  sebagian_terdampak: 'Sebagian terdampak',
  mati: 'Mati',
}

const diseaseLabelMap: Record<string, string> = {
  '': 'Tidak ada',
  batang_busuk: 'Batang busuk',
  lainnya: 'Lainnya',
}

function formatDate(value: string | Date) {
  const date = value instanceof Date ? value : new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'long',
    year: 'numeric',
  })
}

function formatShortDate(value: string | Date) {
  const date = value instanceof Date ? value : new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'short',
  })
}

function healthLabel(value: string) {
  return healthLabelMap[value] || value || '-'
}

function diseaseLabel(value: string) {
  return diseaseLabelMap[value] || value || 'Tidak ada'
}

const summaryMax = computed(() => {
  if (!harvestSummary.value.length) return 1
  return Math.max(...harvestSummary.value.map((item: any) => Number(item.total_quantity || 0)), 1)
})

const summaryTotal = computed(() =>
  harvestSummary.value.reduce((acc: number, item: any) => acc + Number(item.total_quantity || 0), 0)
)

const summaryAverage = computed(() => {
  if (!harvestSummary.value.length) return 0
  return summaryTotal.value / harvestSummary.value.length
})

function summaryBarHeight(value: number) {
  return `${(value / summaryMax.value) * 100}%`
}

const harvestMap = computed(() => {
  const map = new Map<number, any>()
  harvests.value.forEach((item) => {
    if (item?.id) {
      map.set(item.id, item)
    }
  })
  return map
})

const tableRows = computed(() => {
  return harvestRequests.value.map((item) => {
    const harvest = item.harvest_id ? harvestMap.value.get(item.harvest_id) : null
    return {
      id: item.id,
      status: item.status,
      batchCode: item.batch_code,
      packageName: item.package_name,
      requester: item.requested_by_name || item.requested_by,
      monitoringDate: item.monitoring_date,
      monitoringNote: item.note,
      totalQuantity: item.total_quantity || 0,
      wetQuantity: item.wet_quantity || 0,
      dryQuantity: item.dry_quantity || 0,
      harvestId: item.harvest_id,
      harvestDate: item.harvest_date || harvest?.harvest_date || harvest?.harvestDate,
      harvestNotes: harvest?.notes || null,
    }
  })
})

const filteredRows = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  if (!query) return tableRows.value

  return tableRows.value.filter((row) => {
    const tokens = [
      row.batchCode,
      row.packageName,
      row.requester,
      row.status,
      row.monitoringNote,
      row.harvestNotes,
      String(row.totalQuantity || ''),
      String(row.wetQuantity || ''),
      String(row.dryQuantity || ''),
      row.harvestDate ? formatDate(row.harvestDate) : '',
      row.monitoringDate ? formatDate(row.monitoringDate) : '',
    ]
      .filter(Boolean)
      .map((value) => String(value).toLowerCase())

    return tokens.some((value) => value.includes(query))
  })
})

// Pagination state & logic for harvest requests
const currentPage = ref(1)
const limit = ref(8)

const totalPages = computed(() =>
  Math.max(1, Math.ceil(filteredRows.value.length / limit.value))
)

const paginatedRows = computed(() => {
  const start = (currentPage.value - 1) * limit.value
  return filteredRows.value.slice(start, start + limit.value)
})

watch(filteredRows, () => {
  if (currentPage.value > totalPages.value) {
    currentPage.value = totalPages.value
  }
})

watch(searchQuery, () => {
  currentPage.value = 1
})

const selectedMonitoring = computed(() =>
  panenMonitorings.value.find((item: any) => item.id === selectedMonitoringId.value) || null
)

function buildHarvestMessage(monitoring: any) {
  const batchLabel = monitoring?.batch_code || monitoring?.batchCode || '-'
  const packageLabel = monitoring?.package_name || monitoring?.packageName || '-'
  const monitoringDate = formatDate(monitoring?.created_at || monitoring?.monitoring_date || new Date())
  const requestDate = formatDate(new Date())

  const lines = [
    'Halo Admin Pondok Tani,',
    '',
    'Saya ingin mengajukan panen dengan detail monitoring berikut:',
    `Tanggal pengajuan: ${requestDate}`,
    `Batch: ${batchLabel}`,
    `Paket: ${packageLabel}`,
    `Fase: ${phaseLabelMap[monitoring?.phase] || monitoring?.phase || '-'}`,
    `Status kesehatan: ${healthLabel(monitoring?.health_status || monitoring?.healthStatus || '-')}`,
    `Disease: ${diseaseLabel(monitoring?.disease || '')}`,
    `Tanaman terdampak: ${monitoring?.affected_count ?? monitoring?.affectedCount ?? 0} / ${monitoring?.total_plants ?? monitoring?.totalPlants ?? 0}`,
  ]

  if (monitoring?.disease_note || monitoring?.diseaseNote) {
    lines.push(`Disease note: ${monitoring?.disease_note || monitoring?.diseaseNote}`)
  }

  lines.push(`Catatan: ${monitoring?.note || '-'}`)
  lines.push(`Tanggal monitoring: ${monitoringDate}`)
  lines.push('', 'Terima kasih.')

  return lines.join('\n')
}

function openWhatsApp(message: string) {
  const phoneNumber = '6281328164003'
  const url = `https://wa.me/${phoneNumber}?text=${encodeURIComponent(message)}`
  window.open(url, '_blank')
}

async function loadHarvestRequests() {
  const response = await listHarvestRequests()
  harvestRequests.value = Array.isArray(response.data) ? response.data : []
}

async function loadHarvests() {
  const response = await listHarvests()
  harvests.value = Array.isArray(response.data) ? response.data : []
}

async function loadHarvestSummary() {
  summaryLoading.value = true
  summaryError.value = ''
  try {
    const params: Record<string, string> = {}
    if (selectedBatchFilter.value) {
      params.plant_batch_id = selectedBatchFilter.value
    }
    if (dateFilters.value.start) {
      params.start_date = dateFilters.value.start
    }
    if (dateFilters.value.end) {
      params.end_date = dateFilters.value.end
    }
    const response = await getHarvestSummary(params)
    harvestSummary.value = Array.isArray(response.data) ? response.data : []
  } catch (error) {
    console.error('Failed to load harvest summary:', error)
    summaryError.value = 'Gagal memuat ringkasan panen.'
  } finally {
    summaryLoading.value = false
  }
}

async function loadPanenMonitorings() {
  const response = await listPlantMonitorings()
  const rows = Array.isArray(response.data) ? response.data : []
  const map = new Map<number, any>()
  rows.forEach((item: any) => {
    if (!map.has(item.plant_batch_id)) {
      map.set(item.plant_batch_id, {
        id: item.plant_batch_id,
        batchCode: item.batch_code,
        packageName: item.package_name,
      })
    }
  })
  batchOptions.value = Array.from(map.values())
  panenMonitorings.value = isAdmin.value ? [] : rows.filter((item: any) => item.phase === 'panen')
}

async function refreshData() {
  loading.value = true
  errorMessage.value = ''
  try {
    await Promise.all([loadHarvestRequests(), loadHarvests(), loadPanenMonitorings()])
    await loadHarvestSummary()
  } catch (error) {
    console.error('Failed to load harvest data:', error)
    errorMessage.value = 'Gagal memuat data panen.'
  } finally {
    loading.value = false
  }
}

function openRequestModal() {
  selectedMonitoringId.value = null
  showRequestModal.value = true
}

async function confirmRequest() {
  if (!selectedMonitoring.value) return
  requesting.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    await createHarvestRequest({ plant_monitoring_id: selectedMonitoring.value.id })
    successMessage.value = 'Ajuan panen berhasil dibuat.'
    showRequestModal.value = false
    await loadHarvestRequests()
    openWhatsApp(buildHarvestMessage(selectedMonitoring.value))
  } catch (error: any) {
    console.error('Failed to create harvest request:', error)
    errorMessage.value = error?.response?.data?.error || 'Gagal mengajukan panen.'
  } finally {
    requesting.value = false
  }
}

function openValidateModal(item: any) {
  if (!item) return
  selectedRequest.value = item
  validationForm.value.harvestDate = new Date().toISOString().slice(0, 10)
  validationForm.value.notes = item.note || ''
  validationForm.value.totalQuantity = Number(item.total_quantity || 0)
  validationForm.value.wetQuantity = Number(item.wet_quantity || 0)
  validationForm.value.dryQuantity = Number(item.dry_quantity || 0)
  lastEditedHarvestType.value = null
  showValidateModal.value = true
}

async function confirmValidation() {
  if (!selectedRequest.value) return
  if (!isValidationFormComplete.value) return
  if (!validationForm.value.harvestDate) {
    errorMessage.value = 'Tanggal panen wajib diisi.'
    return
  }

  const totalQty = numericValue(validationForm.value.totalQuantity)
  const wetQty = numericValue(validationForm.value.wetQuantity)
  const dryQty = numericValue(validationForm.value.dryQuantity)
  if (totalQty <= 0 || wetQty + dryQty <= 0) {
    errorMessage.value = 'Kuantitas panen wajib diisi.'
    return
  }
  if (Math.abs(totalQty - (wetQty + dryQty)) > 0.01) {
    errorMessage.value = 'Total kuantitas harus sama dengan basah + kering.'
    return
  }

  validating.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    await validateHarvestRequest(selectedRequest.value.id, {
      harvest_date: validationForm.value.harvestDate,
      notes: validationForm.value.notes ? validationForm.value.notes.trim() : null,
      total_quantity: totalQty,
      wet_quantity: wetQty,
      dry_quantity: dryQty,
    })
    successMessage.value = 'Ajuan panen berhasil divalidasi.'
    showValidateModal.value = false
    await refreshData()
  } catch (error: any) {
    console.error('Failed to validate harvest request:', error)
    errorMessage.value = error?.response?.data?.error || 'Gagal memvalidasi ajuan panen.'
  } finally {
    validating.value = false
  }
}

function openDetailModal(item: any) {
  selectedDetail.value = item
  showDetailModal.value = true
}

onMounted(() => {
  refreshData()
})

watch(
  [selectedBatchFilter, () => dateFilters.value.start, () => dateFilters.value.end],
  () => {
    loadHarvestSummary()
  }
)
</script>

<template>
  <div class="harvest-hub space-y-6">
    <div>
      <h2 class="text-2xl font-semibold text-gray-900">Panen</h2>
      <p class="mt-1 text-gray-600">Kelola ajuan panen dan data panen terbaru.</p>
    </div>

    <div v-if="errorMessage" class="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
      {{ errorMessage }}
    </div>
    <div v-if="successMessage" class="rounded-xl border border-green-200 bg-green-50 px-4 py-3 text-sm text-green-700">
      {{ successMessage }}
    </div>

    <div class="grid grid-cols-1 gap-4 md:grid-cols-4">
      <div v-for="card in summaryCards" :key="card.label" class="rounded-xl border border-gray-200 bg-white p-5">
        <p class="text-sm text-gray-600">{{ card.label }}</p>
        <p class="mt-2 text-2xl font-semibold text-gray-900">{{ card.value }}</p>
      </div>
    </div>

    <section class="rounded-2xl border border-emerald-100 bg-gradient-to-br from-emerald-50/70 via-white to-white p-6">
      <div class="flex flex-wrap items-center justify-between gap-4">
        <div class="flex items-center gap-3">
          <div class="flex h-11 w-11 items-center justify-center rounded-xl bg-emerald-600 text-white shadow-md shadow-emerald-200">
            <BarChart3 class="h-5 w-5" />
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Tren kuantitas panen</h3>
            <p class="text-sm text-gray-600">Pantau total panen dari waktu ke waktu.</p>
          </div>
        </div>
        <div class="flex flex-wrap items-center gap-3">
          <select
            v-model="selectedBatchFilter"
            class="rounded-lg border border-emerald-200 bg-white px-3 py-2 text-sm text-gray-700 shadow-sm focus:border-emerald-400 focus:outline-none focus:ring-2 focus:ring-emerald-100"
          >
            <option value="">Semua batch</option>
            <option v-for="option in batchOptions" :key="option.id" :value="String(option.id)">
              {{ option.batchCode }}
            </option>
          </select>
          <input
            v-model="dateFilters.start"
            type="date"
            class="rounded-lg border border-emerald-200 bg-white px-3 py-2 text-sm text-gray-700 shadow-sm focus:border-emerald-400 focus:outline-none focus:ring-2 focus:ring-emerald-100"
          />
          <input
            v-model="dateFilters.end"
            type="date"
            class="rounded-lg border border-emerald-200 bg-white px-3 py-2 text-sm text-gray-700 shadow-sm focus:border-emerald-400 focus:outline-none focus:ring-2 focus:ring-emerald-100"
          />
        </div>
      </div>

      <div v-if="summaryError" class="mt-4 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
        {{ summaryError }}
      </div>

      <div class="mt-5 grid grid-cols-1 gap-4 lg:grid-cols-3">
        <div class="rounded-xl border border-emerald-100 bg-white/80 p-4">
          <p class="text-xs uppercase tracking-wide text-emerald-600">Total panen</p>
          <p class="mt-2 text-2xl font-semibold text-gray-900">{{ summaryTotal.toFixed(1) }} kg</p>
        </div>
        <div class="rounded-xl border border-emerald-100 bg-white/80 p-4">
          <p class="text-xs uppercase tracking-wide text-emerald-600">Rata-rata</p>
          <p class="mt-2 text-2xl font-semibold text-gray-900">{{ summaryAverage.toFixed(1) }} kg</p>
        </div>
        <div class="rounded-xl border border-emerald-100 bg-white/80 p-4">
          <p class="text-xs uppercase tracking-wide text-emerald-600">Terakhir panen</p>
          <p class="mt-2 text-2xl font-semibold text-gray-900">{{ latestHarvestDate }}</p>
        </div>
      </div>

      <div class="mt-6 rounded-2xl border border-emerald-100 bg-white p-5">
        <div v-if="summaryLoading" class="flex items-center gap-2 text-sm text-gray-500">
          <Loader2 class="h-4 w-4 animate-spin" />
          Memuat grafik panen...
        </div>
        <div v-else-if="!harvestSummary.length" class="text-sm text-gray-500">
          Belum ada data panen untuk rentang waktu ini.
        </div>
        <div v-else class="flex items-end gap-3 overflow-x-auto pb-2">
          <div
            v-for="(item, index) in harvestSummary"
            :key="`${item.date}-${index}`"
            class="flex min-w-[52px] flex-col items-center justify-end"
          >
            <div class="flex h-44 w-full items-end justify-center">
              <div
                class="w-6 rounded-full bg-gradient-to-t from-emerald-500 via-emerald-400 to-emerald-200 shadow-md shadow-emerald-200/80"
                :style="{ height: summaryBarHeight(Number(item.total_quantity || 0)) }"
              />
            </div>
            <p class="mt-2 text-xs font-medium text-gray-600">{{ formatShortDate(item.date) }}</p>
            <p class="text-[11px] text-gray-500">{{ Number(item.total_quantity || 0).toFixed(1) }} kg</p>
          </div>
        </div>
      </div>
    </section>

    <section class="rounded-xl border border-gray-200 bg-white p-5">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold text-gray-900">Data ajuan panen</h3>
          <p class="mt-1 text-sm text-gray-600">Terakhir panen: {{ latestHarvestDate }}</p>
        </div>
        <div class="flex flex-wrap items-center gap-3">
          <div class="relative w-full max-w-sm">
            <Search class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
            <input
              v-model="searchQuery"
              type="text"
              class="w-full rounded-lg border border-gray-300 bg-white py-2.5 pl-10 pr-3 text-sm outline-none transition focus:border-green-500 focus:ring-2 focus:ring-green-100"
              placeholder="Cari batch, paket, status, catatan"
            />
          </div>
          <button
            v-if="!isAdmin"
            type="button"
            class="inline-flex items-center gap-2 rounded-lg border border-green-600 px-4 py-2 text-sm font-semibold text-green-700 hover:bg-green-50"
            @click="openRequestModal"
          >
            <MessageCircle class="h-4 w-4" />
            Ajukan Panen
          </button>
        </div>
      </div>

      <div v-if="loading" class="mt-4 flex items-center gap-2 text-sm text-gray-500">
        <Loader2 class="h-4 w-4 animate-spin" />
        Memuat data panen...
      </div>

      <div v-else class="mt-4 overflow-x-auto">
        <table class="min-w-full text-sm">
          <thead class="bg-gray-50 text-gray-700">
            <tr>
              <th class="px-4 py-2 text-left font-medium">Batch</th>
              <th class="px-4 py-2 text-left font-medium">Pengaju</th>
              <th class="px-4 py-2 text-left font-medium">Tanggal monitoring</th>
              <th class="px-4 py-2 text-left font-medium">Status</th>
              <th class="px-4 py-2 text-right font-medium">Total (kg)</th>
              <th class="px-4 py-2 text-left font-medium">Tanggal panen</th>
              <th class="px-4 py-2 text-left font-medium">Catatan panen</th>
              <th class="px-4 py-2 text-center font-medium">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-if="!filteredRows.length">
              <td colspan="8" class="px-4 py-4 text-center text-gray-500">Belum ada data panen.</td>
            </tr>
            <tr v-for="row in paginatedRows" :key="row.id">
              <td class="px-4 py-3">
                <p class="font-medium text-gray-900">{{ row.batchCode }}</p>
                <p class="text-xs text-gray-500">{{ row.packageName }}</p>
              </td>
              <td class="px-4 py-3 text-gray-700">{{ row.requester }}</td>
              <td class="px-4 py-3 text-gray-700">{{ formatDate(row.monitoringDate) }}</td>
              <td class="px-4 py-3">
                <span
                  class="rounded-full px-2 py-0.5 text-xs font-medium"
                  :class="row.status === 'validated' ? 'bg-green-100 text-green-700' : 'bg-amber-100 text-amber-700'"
                >
                  {{ row.status === 'validated' ? 'Tervalidasi' : 'Pending' }}
                </span>
              </td>
              <td class="px-4 py-3 text-right text-gray-700">{{ row.totalQuantity.toFixed(1) }}</td>
              <td class="px-4 py-3 text-gray-700">{{ row.harvestDate ? formatDate(row.harvestDate) : '-' }}</td>
              <td class="px-4 py-3 text-gray-700">{{ row.harvestNotes || '-' }}</td>
              <td class="px-4 py-3 text-center">
                <div class="flex items-center justify-center gap-2">
                  <button
                    v-if="isAdmin && row.status === 'pending'"
                    type="button"
                    class="inline-flex items-center gap-2 rounded-lg bg-green-600 px-3 py-2 text-xs font-semibold text-white hover:bg-green-700"
                    @click="openValidateModal(harvestRequests.find((item) => item.id === row.id))"
                  >
                    <ClipboardCheck class="h-4 w-4" />
                    Validasi
                  </button>
                  <button
                    v-else-if="row.status === 'validated'"
                    type="button"
                    class="inline-flex items-center gap-2 rounded-lg border border-gray-300 px-3 py-2 text-xs font-semibold text-gray-700 hover:bg-gray-50"
                    @click="openDetailModal(row)"
                  >
                    <Eye class="h-4 w-4" />
                    Tinjau
                  </button>
                  <span v-else class="text-xs text-gray-500">-</span>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <Pagination
        v-if="filteredRows.length > 0"
        :current-page="currentPage"
        :total-pages="totalPages"
        :total-rows="filteredRows.length"
        :limit="limit"
        @update:page="currentPage = $event"
      />
    </section>

    <div
      v-if="showRequestModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
    >
      <div class="w-full max-w-2xl rounded-2xl bg-white">
        <div class="border-b border-gray-200 px-6 py-4">
          <h3 class="text-lg font-semibold text-gray-900">Ajukan panen</h3>
          <p class="mt-1 text-sm text-gray-600">Pilih monitoring fase panen yang ingin diajukan.</p>
        </div>
        <div class="space-y-4 p-6 text-sm text-gray-700">
          <div v-if="!panenMonitorings.length" class="rounded-lg border border-dashed border-gray-300 p-4 text-sm text-gray-500">
            Belum ada monitoring dengan fase panen.
          </div>
          <div v-else class="max-h-[320px] space-y-3 overflow-auto">
            <label
              v-for="item in panenMonitorings"
              :key="item.id"
              class="flex cursor-pointer items-start gap-3 rounded-xl border border-gray-200 p-4 transition hover:border-green-400"
            >
              <input
                v-model="selectedMonitoringId"
                type="radio"
                name="panen-monitoring"
                class="mt-1"
                :value="item.id"
              />
              <div class="flex-1">
                <div class="flex items-center justify-between">
                  <div>
                    <p class="text-sm text-gray-500">{{ formatDate(item.created_at) }}</p>
                    <p class="font-semibold text-gray-900">{{ item.batch_code }}</p>
                    <p class="text-sm text-gray-600">{{ item.package_name || '-' }}</p>
                  </div>
                  <span class="rounded-full bg-green-100 px-2 py-0.5 text-xs font-medium text-green-700">Panen</span>
                </div>
                <div class="mt-2 grid grid-cols-2 gap-2 text-sm text-gray-700">
                  <p>Status: <span class="font-medium">{{ healthLabel(item.health_status) }}</span></p>
                  <p>Disease: <span class="font-medium">{{ diseaseLabel(item.disease) }}</span></p>
                  <p class="col-span-2">Terdampak: <span class="font-medium">{{ item.affected_count }} / {{ item.total_plants }}</span></p>
                </div>
              </div>
            </label>
          </div>
        </div>
        <div class="flex gap-3 border-t border-gray-200 px-6 py-4">
          <button
            type="button"
            class="flex-1 rounded-lg border border-gray-300 px-4 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50"
            @click="showRequestModal = false"
          >
            <X class="h-4 w-4 inline-block mr-1" />
            Batal
          </button>
          <button
            type="button"
            class="flex-1 rounded-lg bg-green-600 px-4 py-2 text-sm font-semibold text-white hover:bg-green-700 disabled:opacity-60"
            :disabled="requesting || !selectedMonitoringId"
            @click="confirmRequest"
          >
            <Loader2 v-if="requesting" class="h-4 w-4 inline-block mr-1 animate-spin" />
            <Check v-else class="h-4 w-4 inline-block mr-1" />
            Ajukan
          </button>
        </div>
      </div>
    </div>

    <div
      v-if="showValidateModal && selectedRequest"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
    >
      <div class="w-full max-w-lg rounded-2xl bg-white">
        <div class="border-b border-gray-200 px-6 py-4">
          <h3 class="text-lg font-semibold text-gray-900">Validasi dan input data panen</h3>
          <p class="mt-1 text-sm text-gray-600">Masukkan tanggal dan catatan panen.</p>
        </div>
        <div class="space-y-4 p-6 text-sm text-gray-700">
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">Harvest ID</label>
            <input
              type="text"
              :value="selectedRequest?.harvest_id || 'Akan dibuat setelah disimpan'"
              class="w-full rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 text-sm text-gray-500"
              disabled
            />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">Tanggal panen</label>
            <input
              v-model="validationForm.harvestDate"
              type="date"
              class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-100"
            />
          </div>
          <div class="grid grid-cols-1 gap-3 md:grid-cols-3">
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700">Total panen (kg)</label>
              <input
                v-model.number="validationForm.totalQuantity"
                type="number"
                min="0"
                step="0.1"
                class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-100"
                @input="completeCounterpartQuantity('total')"
              />
            </div>
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700">Basah (kg)</label>
              <input
                v-model.number="validationForm.wetQuantity"
                type="number"
                min="0"
                step="0.1"
                class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-100"
                @input="completeCounterpartQuantity('wet')"
              />
            </div>
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700">Kering (kg)</label>
              <input
                v-model.number="validationForm.dryQuantity"
                type="number"
                min="0"
                step="0.1"
                class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-100"
                @input="completeCounterpartQuantity('dry')"
              />
            </div>
          </div>
          <p class="text-xs text-gray-500">Pastikan total panen = basah + kering.</p>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">Catatan panen</label>
            <textarea
              v-model="validationForm.notes"
              rows="3"
              class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-100"
              placeholder="Catatan tambahan"
            />
          </div>
        </div>
        <div class="flex gap-3 border-t border-gray-200 px-6 py-4">
          <button
            type="button"
            class="flex-1 rounded-lg border border-gray-300 px-4 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50"
            @click="showValidateModal = false"
          >
            <X class="h-4 w-4 inline-block mr-1" />
            Batal
          </button>
          <button
            type="button"
            class="flex-1 rounded-lg bg-green-600 px-4 py-2 text-sm font-semibold text-white hover:bg-green-700 disabled:cursor-not-allowed disabled:bg-gray-300 disabled:text-gray-500 disabled:hover:bg-gray-300"
            :disabled="validating || !isValidationFormComplete"
            @click="confirmValidation"
          >
            <Loader2 v-if="validating" class="h-4 w-4 inline-block mr-1 animate-spin" />
            <Check v-else class="h-4 w-4 inline-block mr-1" />
            Simpan
          </button>
        </div>
      </div>
    </div>

    <div
      v-if="showDetailModal && selectedDetail"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
    >
      <div class="w-full max-w-lg rounded-2xl bg-white">
        <div class="border-b border-gray-200 px-6 py-4">
          <h3 class="text-lg font-semibold text-gray-900">Tinjau data panen</h3>
          <p class="mt-1 text-sm text-gray-600">Detail panen yang sudah divalidasi.</p>
        </div>
        <div class="space-y-4 p-6 text-sm text-gray-700">
          <div class="rounded-lg bg-gray-50 p-4">
            <p class="text-xs text-gray-500">Batch</p>
            <p class="font-semibold text-gray-900">{{ selectedDetail.batchCode }}</p>
            <p class="text-sm text-gray-600">{{ selectedDetail.packageName || '-' }}</p>
          </div>
          <div class="space-y-2">
            <p>Pengaju: <span class="font-medium">{{ selectedDetail.requester }}</span></p>
            <p>Tanggal monitoring: <span class="font-medium">{{ formatDate(selectedDetail.monitoringDate) }}</span></p>
            <p>Harvest ID: <span class="font-medium">{{ selectedDetail.harvestId || '-' }}</span></p>
            <p>Tanggal panen: <span class="font-medium">{{ selectedDetail.harvestDate ? formatDate(selectedDetail.harvestDate) : '-' }}</span></p>
            <p>Total panen: <span class="font-medium">{{ selectedDetail.totalQuantity.toFixed(1) }} kg</span></p>
            <p>Basah: <span class="font-medium">{{ selectedDetail.wetQuantity.toFixed(1) }} kg</span></p>
            <p>Kering: <span class="font-medium">{{ selectedDetail.dryQuantity.toFixed(1) }} kg</span></p>
            <p>Catatan panen: <span class="font-medium">{{ selectedDetail.harvestNotes || '-' }}</span></p>
          </div>
        </div>
        <div class="flex justify-end border-t border-gray-200 px-6 py-4">
          <button
            type="button"
            class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50"
            @click="showDetailModal = false"
          >
            <X class="h-4 w-4 inline-block mr-1" />
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
