<script setup lang="ts">
import { computed, ref } from 'vue';
import { Package, TrendingDown, TrendingUp, ArrowRight } from 'lucide-vue-next';
import SearchBar from '../../components/common/SearchBar.vue';
import Pagination from '../../components/common/Pagination.vue';

const stockByGrade = [
  { grade: 'A', quantity: 145, value: 725000000, percentage: 85 },
  { grade: 'B', quantity: 98, value: 416500000, percentage: 60 },
  { grade: 'C', quantity: 52, value: 182000000, percentage: 35 },
  { grade: 'Split', quantity: 28, value: 70000000, percentage: 20 },
  { grade: 'Asalan', quantity: 15, value: 30000000, percentage: 10 },
];

const stockMovements = [
  {
    id: 'MOV-001',
    source: 'Panen - BATCH-045',
    type: 'In' as const,
    grade: 'A',
    quantityIn: 25,
    quantityOut: 0,
    date: '2024-01-20',
    batchRef: 'HRV-001',
  },
  {
    id: 'MOV-002',
    source: 'Penjualan - Ekspor Premium',
    type: 'Out' as const,
    grade: 'A',
    quantityIn: 0,
    quantityOut: 40,
    date: '2024-01-19',
    batchRef: 'SALE-089',
  },
  {
    id: 'MOV-003',
    source: 'Panen - BATCH-038',
    type: 'In' as const,
    grade: 'B',
    quantityIn: 35,
    quantityOut: 0,
    date: '2024-01-18',
    batchRef: 'HRV-002',
  },
  {
    id: 'MOV-004',
    source: 'Penjualan - Distributor Lokal',
    type: 'Out' as const,
    grade: 'B',
    quantityIn: 0,
    quantityOut: 25,
    date: '2024-01-17',
    batchRef: 'SALE-088',
  },
  {
    id: 'MOV-005',
    source: 'Panen - BATCH-052',
    type: 'In' as const,
    grade: 'A',
    quantityIn: 30,
    quantityOut: 0,
    date: '2024-01-16',
    batchRef: 'HRV-003',
  },
];

const monthlyMovement = [
  { month: 'Aug', stockIn: 120, stockOut: 85 },
  { month: 'Sep', stockIn: 145, stockOut: 95 },
  { month: 'Oct', stockIn: 160, stockOut: 110 },
  { month: 'Nov', stockIn: 135, stockOut: 100 },
  { month: 'Dec', stockIn: 155, stockOut: 120 },
  { month: 'Jan', stockIn: 142, stockOut: 105 },
];

const chartMax = computed(() =>
  Math.max(
    ...monthlyMovement.flatMap((m) => [m.stockIn, m.stockOut]),
    1
  )
);

function barPct(value: number) {
  return `${(value / chartMax.value) * 100}%`;
}

const searchQuery = ref('');
const currentPage = ref(1);
const limit = ref(5);

const filteredMovements = computed(() => {
  if (!searchQuery.value) return stockMovements;
  const q = searchQuery.value.toLowerCase();
  return stockMovements.filter(m => 
    m.id.toLowerCase().includes(q) || 
    m.source.toLowerCase().includes(q) || 
    m.grade.toLowerCase().includes(q) ||
    m.batchRef.toLowerCase().includes(q)
  );
});

const totalPages = computed(() => Math.ceil(filteredMovements.value.length / limit.value));

const paginatedMovements = computed(() => {
  const start = (currentPage.value - 1) * limit.value;
  return filteredMovements.value.slice(start, start + limit.value);
});

function handleSearch(val: string) {
  searchQuery.value = val;
  currentPage.value = 1;
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-semibold text-gray-900">Stok Gudang</h2>
      <p class="text-gray-600 mt-1">Pantau level stok vanili dan pergerakannya</p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Total stok</p>
          <Package class="w-5 h-5 text-blue-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">338 kg</p>
        <p class="text-xs text-gray-500 mt-1">Semua mutu digabung</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Nilai total</p>
          <Package class="w-5 h-5 text-green-600" />
        </div>
        <p class="text-2xl font-semibold text-green-600">$1.42M</p>
        <p class="text-xs text-gray-500 mt-1">Perkiraan nilai pasar</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Masuk gudang (MTD)</p>
          <TrendingUp class="w-5 h-5 text-green-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">142 kg</p>
        <p class="text-xs text-green-600 mt-1">+12% vs bulan lalu</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Keluar gudang (MTD)</p>
          <TrendingDown class="w-5 h-5 text-orange-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">105 kg</p>
        <p class="text-xs text-orange-600 mt-1">-5% vs bulan lalu</p>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-6">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">Stok per mutu</h3>
      <div class="space-y-4">
        <div v-for="item in stockByGrade" :key="item.grade" class="space-y-2">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <span class="font-medium text-gray-900 w-16">Mutu {{ item.grade }}</span>
              <div
                class="relative flex-1 h-8 bg-gray-100 rounded-lg overflow-hidden"
                style="width: 300px"
              >
                <div
                  class="absolute inset-y-0 left-0 bg-gradient-to-r from-green-500 to-green-600 flex items-center justify-end px-2"
                  :style="{ width: `${item.percentage}%` }"
                >
                  <span class="text-xs font-medium text-white">{{ item.percentage }}%</span>
                </div>
              </div>
            </div>
            <div class="text-right">
              <p class="font-semibold text-gray-900">{{ item.quantity }} kg</p>
              <p class="text-xs text-gray-500">
                {{ '$' + (item.value / 1000000).toFixed(2) }}M
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-6">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">Tren pergerakan stok</h3>
      <p class="text-xs text-gray-500 mb-2 flex gap-4">
        <span class="inline-flex items-center gap-1">
          <span class="inline-block w-3 h-3 rounded-sm bg-[#10b981]" /> Masuk (kg)
        </span>
        <span class="inline-flex items-center gap-1">
          <span class="inline-block w-3 h-3 rounded-sm bg-[#f59e0b]" /> Keluar (kg)
        </span>
      </p>
      <div class="w-full h-[300px] flex flex-col">
        <div class="flex-1 flex items-end gap-2 px-2 min-h-0">
          <div
            v-for="m in monthlyMovement"
            :key="m.month"
            class="flex-1 flex flex-col items-center justify-end h-full min-w-0"
          >
            <div
              class="flex gap-1 items-end justify-center w-full"
              style="height: 220px"
              :title="`${m.month}: Masuk ${m.stockIn} kg, Keluar ${m.stockOut} kg`"
            >
              <div
                class="w-[28%] max-w-5 rounded-t-md bg-[#10b981] transition-all"
                :style="{ height: barPct(m.stockIn) }"
              />
              <div
                class="w-[28%] max-w-5 rounded-t-md bg-[#f59e0b] transition-all"
                :style="{ height: barPct(m.stockOut) }"
              />
            </div>
            <span class="text-xs text-[#6b7280] mt-2">{{ m.month }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden">
      <div class="flex flex-wrap items-center justify-between gap-3 border-b border-gray-200 px-6 py-4">
        <div>
          <h3 class="text-lg font-semibold text-gray-900">Riwayat pergerakan stok</h3>
          <p class="text-sm text-gray-600 mt-1">Transaksi masuk dan keluar gudang terbaru</p>
        </div>
        <div class="w-full max-w-sm">
          <SearchBar
            v-model="searchQuery"
            placeholder="Cari ID, sumber, batch..."
            @search="handleSearch"
          />
        </div>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">ID pergerakan</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Sumber</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Mutu</th>
              <th class="text-center px-6 py-3 text-sm font-medium text-gray-900">Jumlah masuk</th>
              <th class="text-center px-6 py-3 text-sm font-medium text-gray-900">Jumlah keluar</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Tanggal</th>
              <th class="text-center px-6 py-3 text-sm font-medium text-gray-900">Referensi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="movement in paginatedMovements" :key="movement.id" class="hover:bg-gray-50">
              <td class="px-6 py-4">
                <div class="flex items-center gap-2">
                  <TrendingUp v-if="movement.type === 'In'" class="w-4 h-4 text-green-600" />
                  <TrendingDown v-else class="w-4 h-4 text-orange-600" />
                  <span class="font-medium text-gray-900">{{ movement.id }}</span>
                </div>
              </td>
              <td class="px-6 py-4 text-gray-900">{{ movement.source }}</td>
              <td class="px-6 py-4">
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
                >
                  Mutu {{ movement.grade }}
                </span>
              </td>
              <td class="px-6 py-4 text-center">
                <span v-if="movement.quantityIn > 0" class="font-semibold text-green-600">
                  +{{ movement.quantityIn }} kg
                </span>
                <span v-else class="text-gray-400">-</span>
              </td>
              <td class="px-6 py-4 text-center">
                <span v-if="movement.quantityOut > 0" class="font-semibold text-orange-600">
                  -{{ movement.quantityOut }} kg
                </span>
                <span v-else class="text-gray-400">-</span>
              </td>
              <td class="px-6 py-4 text-gray-900">{{ movement.date }}</td>
              <td class="px-6 py-4 text-center">
                <button
                  type="button"
                  class="text-blue-600 hover:text-blue-700 font-medium text-sm flex items-center gap-1 mx-auto"
                >
                  {{ movement.batchRef }}
                  <ArrowRight class="w-3 h-3" />
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
  </div>
</template>
