<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import {
  ArrowRight,
  Loader2,
  Package,
  TrendingDown,
  TrendingUp,
  Warehouse,
} from 'lucide-vue-next'
import { getWarehouseSummary, listStockMovements, listWarehouseStocks } from '../../services/warehouse/warehouse'
import Pagination from '../../components/common/Pagination.vue'

const summary = ref<any>({
  total_quantity: 0,
  total_batches: 0,
  total_grades: 0,
  grade_breakdown: [],
  monthly_movements: [],
})
const stocks = ref<any[]>([])
const movements = ref<any[]>([])
const loading = ref(false)
const errorMessage = ref('')

const currentPage = ref(1)
const limit = ref(8)

const availableStocks = computed(() =>
  stocks.value.filter((item) => Number(item.total_quantity || 0) > 0)
)

const gradeBreakdown = computed(() => {
  const rows = Array.isArray(summary.value.grade_breakdown) ? summary.value.grade_breakdown : []
  if (rows.length > 0) return rows

  const grouped = new Map<string, { grade_id: number | null; grade_name: string; quantity: number }>()
  availableStocks.value.forEach((item) => {
    const gradeName = item.grade_name || 'Basah'
    const key = `${item.grade_id || 0}-${gradeName}`
    const existing = grouped.get(key) || {
      grade_id: item.grade_id || null,
      grade_name: gradeName,
      quantity: 0,
    }
    existing.quantity += Number(item.total_quantity || 0)
    grouped.set(key, existing)
  })
  return Array.from(grouped.values())
})
const monthlyMovements = computed(() => summary.value.monthly_movements || [])
const currentMonthKey = computed(() => new Date().toLocaleDateString('en-US', { month: 'short' }))
const currentMonthStats = computed(() => {
  const match = monthlyMovements.value.find((item: any) => item.month === currentMonthKey.value)
  return {
    stockIn: Number(match?.stock_in || 0),
    stockOut: Number(match?.stock_out || 0),
  }
})

const totalPages = computed(() =>
  Math.max(1, Math.ceil(movements.value.length / limit.value))
)

const paginatedMovements = computed(() => {
  const start = (currentPage.value - 1) * limit.value
  return movements.value.slice(start, start + limit.value)
})

const chartMax = computed(() => {
  if (!monthlyMovements.value.length) return 1
  return Math.max(
    ...monthlyMovements.value.flatMap((item: any) => [item.stock_in || 0, item.stock_out || 0]),
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

function formatReferenceType(value: string) {
  const labels: Record<string, string> = {
    harvest_output: 'Panen basah',
    grading_detail: 'Hasil grading',
    sales_detail: 'Penjualan',
  }
  const normalized = String(value || '').trim().toLowerCase()
  if (!normalized) return '-'
  return labels[normalized] || normalized
    .split('_')
    .filter(Boolean)
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ')
}

async function fetchWarehouseData() {
  loading.value = true
  errorMessage.value = ''
  try {
    const [summaryRes, stocksRes, movementsRes] = await Promise.all([
      getWarehouseSummary(),
      listWarehouseStocks(),
      listStockMovements(),
    ])
    summary.value = summaryRes?.data || summary.value
    stocks.value = Array.isArray(stocksRes.data) ? stocksRes.data : []
    movements.value = Array.isArray(movementsRes.data) ? movementsRes.data : []
  } catch (error) {
    console.error('Failed to load user warehouse data:', error)
    errorMessage.value = 'Gagal memuat data stok gudang Anda.'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchWarehouseData()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <div>
        <h2 class="text-2xl font-semibold text-gray-900">Stok Gudang Saya</h2>
        <p class="mt-1 text-gray-600">Pantau hasil panen Anda yang sudah tercatat di gudang.</p>
      </div>
      <button
        type="button"
        class="inline-flex items-center gap-2 rounded-lg border border-gray-300 px-3 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50"
        :disabled="loading"
        @click="fetchWarehouseData"
      >
        <Loader2 v-if="loading" class="h-4 w-4 animate-spin" />
        <Warehouse v-else class="h-4 w-4" />
        Refresh
      </button>
    </div>

    <div v-if="errorMessage" class="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
      {{ errorMessage }}
    </div>

    <div class="grid grid-cols-1 gap-4 md:grid-cols-4">
      <div class="rounded-2xl border border-emerald-100 bg-gradient-to-br from-emerald-50 to-white p-6">
        <div class="flex items-center justify-between">
          <p class="text-sm text-gray-600">Total stok</p>
          <Package class="h-5 w-5 text-emerald-600" />
        </div>
        <p class="mt-2 text-2xl font-semibold text-gray-900">{{ Number(summary.total_quantity || 0).toFixed(1) }} kg</p>
        <p class="text-xs text-gray-500">Hasil panen milik Anda</p>
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
          <p class="text-sm text-gray-600">Batch di gudang</p>
          <Warehouse class="h-5 w-5 text-amber-600" />
        </div>
        <p class="mt-2 text-2xl font-semibold text-gray-900">{{ summary.total_batches || 0 }}</p>
        <p class="text-xs text-gray-500">Batch hasil panen Anda</p>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
      <section class="rounded-2xl border border-gray-200 bg-white p-6">
        <h3 class="text-lg font-semibold text-gray-900">Stok per mutu</h3>
        <div v-if="loading" class="mt-4 flex items-center gap-2 text-sm text-gray-500">
          <Loader2 class="h-4 w-4 animate-spin" />
          Memuat stok...
        </div>
        <div v-else-if="!gradeBreakdown.length" class="mt-4 rounded-lg bg-gray-50 px-4 py-5 text-sm text-gray-500">
          Belum ada hasil panen Anda di gudang.
        </div>
        <div v-else class="mt-4 space-y-3">
          <button
            v-for="item in gradeBreakdown"
            :key="item.grade_id || item.grade_name"
            type="button"
            class="w-full rounded-xl border border-emerald-100 bg-white px-4 py-3 text-left"
          >
            <div class="flex items-center justify-between gap-4">
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
      </section>

      <section class="rounded-2xl border border-gray-200 bg-white p-6">
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
      </section>
    </div>

    <section class="rounded-xl border border-gray-200 bg-white p-6">
        <h3 class="text-lg font-semibold text-gray-900">Detail stok saya</h3>
        <div class="mt-4 overflow-x-auto">
          <table class="w-full min-w-[720px] text-sm">
            <thead class="bg-gray-50 text-gray-700">
              <tr>
                <th class="px-4 py-2 text-left font-medium">Batch</th>
                <th class="px-4 py-2 text-left font-medium">Paket</th>
                <th class="px-4 py-2 text-left font-medium">Mutu</th>
                <th class="px-4 py-2 text-right font-medium">Qty</th>
                <th class="px-4 py-2 text-left font-medium">Update</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-if="!availableStocks.length">
                <td colspan="5" class="px-4 py-5 text-center text-gray-500">Belum ada stok Anda di gudang.</td>
              </tr>
              <tr v-for="item in availableStocks" :key="item.id">
                <td class="px-4 py-3 font-medium text-gray-900">{{ item.batch_code }}</td>
                <td class="px-4 py-3 text-gray-700">{{ item.package_name || '-' }}</td>
                <td class="px-4 py-3">
                  <span class="rounded-full bg-emerald-50 px-2 py-0.5 text-xs font-medium text-emerald-700">
                    {{ item.grade_name || 'Basah' }}
                  </span>
                </td>
                <td class="px-4 py-3 text-right font-semibold text-gray-900">
                  {{ Number(item.total_quantity || 0).toFixed(1) }} {{ item.unit || 'kg' }}
                </td>
                <td class="px-4 py-3 text-gray-600">{{ formatDate(item.updated_at) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
    </section>

    <section class="rounded-xl border border-gray-200 bg-white p-6">
      <div class="border-b border-gray-200 pb-3">
        <h3 class="text-lg font-semibold text-gray-900">Riwayat pergerakan stok saya</h3>
        <p class="text-sm text-gray-500">Catatan masuk dan keluar dari hasil panen Anda.</p>
      </div>
      <div class="mt-4 overflow-x-auto">
        <table class="w-full min-w-[760px] text-sm">
          <thead class="bg-gray-50 text-gray-700">
            <tr>
              <th class="px-4 py-2 text-left font-medium">ID</th>
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
              <td colspan="7" class="px-4 py-5 text-center text-gray-500">Belum ada pergerakan stok Anda.</td>
            </tr>
            <tr v-for="movement in paginatedMovements" :key="movement.id" class="hover:bg-gray-50">
              <td class="px-4 py-3">
                <div class="flex items-center gap-2">
                  <TrendingUp v-if="movement.movement_type === 'in'" class="h-4 w-4 text-emerald-600" />
                  <TrendingDown v-else class="h-4 w-4 text-amber-600" />
                  <span class="font-medium text-gray-900">{{ movement.id }}</span>
                </div>
              </td>
              <td class="px-4 py-3 text-gray-700">{{ formatReferenceType(movement.reference_type) }}</td>
              <td class="px-4 py-3 text-gray-700">{{ movement.grade_name || 'Basah' }}</td>
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
                <span class="inline-flex items-center gap-1 text-xs font-semibold text-emerald-600">
                  {{ movement.reference_id }}
                  <ArrowRight class="h-3 w-3" />
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <Pagination
        v-if="movements.length > 0"
        :current-page="currentPage"
        :total-pages="totalPages"
        :total-rows="movements.length"
        :limit="limit"
        @update:page="currentPage = $event"
      />
    </section>
  </div>
</template>
