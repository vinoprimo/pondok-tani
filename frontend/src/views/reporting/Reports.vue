<script setup lang="ts">
import { computed, ref } from 'vue';
import { Download, FileText, Calendar, Filter, TrendingUp } from 'lucide-vue-next';
import SearchBar from '../../components/common/SearchBar.vue';
import Pagination from '../../components/common/Pagination.vue';

type ReportCategory = 'financial' | 'plant' | 'investment';

const selectedPeriod = ref('monthly');
const selectedReport = ref('all');

const reportTypes = [
  { id: 'all', label: 'Semua laporan' },
  { id: 'financial', label: 'Keuangan' },
  { id: 'plant', label: 'Kesehatan tanaman' },
  { id: 'investment', label: 'Investasi' },
];

const reports = [
  {
    id: 1,
    name: 'Laporan rekening investasi bulanan',
    category: 'investment' as ReportCategory,
    date: '2026-01-01',
    period: 'Januari 2026',
    size: '245 KB',
    format: 'PDF',
  },
  {
    id: 2,
    name: 'Ringkasan keuangan Q2',
    category: 'financial' as ReportCategory,
    date: '2025-07-01',
    period: 'Kuartal II 2025',
    size: '512 KB',
    format: 'PDF',
  },
  {
    id: 3,
    name: 'Laporan kesehatan & pertumbuhan tanaman',
    category: 'plant' as ReportCategory,
    date: '2026-01-15',
    period: 'Januari 2026',
    size: '1.2 MB',
    format: 'PDF',
  },
  {
    id: 4,
    name: 'Laporan analisis ROI',
    category: 'investment' as ReportCategory,
    date: '2025-12-31',
    period: 'Tahunan 2025',
    size: '380 KB',
    format: 'PDF',
  },
  {
    id: 5,
    name: 'Analisis rincian biaya',
    category: 'financial' as ReportCategory,
    date: '2025-12-01',
    period: 'Desember 2025',
    size: '298 KB',
    format: 'PDF',
  },
  {
    id: 6,
    name: 'Laporan hasil panen',
    category: 'plant' as ReportCategory,
    date: '2025-11-20',
    period: 'November 2025',
    size: '445 KB',
    format: 'PDF',
  },
];

const categoryLabels: Record<ReportCategory, string> = {
  financial: 'Keuangan',
  plant: 'Kesehatan tanaman',
  investment: 'Investasi',
};

const searchQuery = ref('');
const currentPage = ref(1);
const limit = ref(5);

const filteredReports = computed(() => {
  let res = reports;
  if (selectedReport.value !== 'all') {
    res = res.filter((report) => report.category === selectedReport.value);
  }
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase();
    res = res.filter((report) => 
      report.name.toLowerCase().includes(q) || 
      report.period.toLowerCase().includes(q) || 
      report.date.toLowerCase().includes(q)
    );
  }
  return res;
});

const totalPages = computed(() => Math.ceil(filteredReports.value.length / limit.value));

const paginatedReports = computed(() => {
  const start = (currentPage.value - 1) * limit.value;
  return filteredReports.value.slice(start, start + limit.value);
});

function handleSearch(val: string) {
  searchQuery.value = val;
  currentPage.value = 1;
}

function reportTypeBadgeClass(category: ReportCategory) {
  if (category === 'financial') return 'bg-blue-100 text-blue-700';
  if (category === 'investment') return 'bg-green-100 text-green-700';
  return 'bg-purple-100 text-purple-700';
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-semibold text-gray-900">Laporan & dokumen</h2>
        <p class="text-gray-600 mt-1">Unduh dan ekspor laporan investasi Anda</p>
      </div>
      <button
        type="button"
        class="flex items-center gap-2 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
      >
        <Download class="w-4 h-4" />
        Buat laporan kustom
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
        <h3 class="font-semibold text-gray-900 mb-1">Rekening portofolio</h3>
        <p class="text-sm text-gray-600">Kepemilikan saat ini & kinerja</p>
      </button>

      <button
        type="button"
        class="bg-white border border-gray-200 rounded-xl p-6 hover:border-green-600 hover:shadow-md transition-all text-left"
      >
        <div class="w-12 h-12 bg-blue-50 rounded-lg flex items-center justify-center mb-4">
          <TrendingUp class="w-6 h-6 text-blue-600" />
        </div>
        <h3 class="font-semibold text-gray-900 mb-1">Analisis ROI</h3>
        <p class="text-sm text-gray-600">Imbal hasil & proyeksi</p>
      </button>

      <button
        type="button"
        class="bg-white border border-gray-200 rounded-xl p-6 hover:border-green-600 hover:shadow-md transition-all text-left"
      >
        <div class="w-12 h-12 bg-purple-50 rounded-lg flex items-center justify-center mb-4">
          <Calendar class="w-6 h-6 text-purple-600" />
        </div>
        <h3 class="font-semibold text-gray-900 mb-1">Dokumen pajak</h3>
        <p class="text-sm text-gray-600">Surat keterangan tahunan</p>
      </button>

      <button
        type="button"
        class="bg-white border border-gray-200 rounded-xl p-6 hover:border-green-600 hover:shadow-md transition-all text-left"
      >
        <div class="w-12 h-12 bg-yellow-50 rounded-lg flex items-center justify-center mb-4">
          <FileText class="w-6 h-6 text-yellow-600" />
        </div>
        <h3 class="font-semibold text-gray-900 mb-1">Laporan tanaman</h3>
        <p class="text-sm text-gray-600">Data pertumbuhan & kesehatan</p>
      </button>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-4">
      <div class="flex flex-wrap items-center gap-4">
        <div class="flex items-center gap-2">
          <Filter class="w-4 h-4 text-gray-600" />
          <span class="text-sm font-medium text-gray-700">Filter:</span>
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
            <option value="all">Sepanjang waktu</option>
            <option value="monthly">Bulan ini</option>
            <option value="quarterly">Kuartal ini</option>
            <option value="yearly">Tahun ini</option>
          </select>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden">
      <div class="flex flex-wrap items-center justify-between gap-3 border-b border-gray-200 px-6 py-4">
        <div>
          <h3 class="text-lg font-semibold text-gray-900">Daftar laporan</h3>
        </div>
        <div class="w-full max-w-sm">
          <SearchBar
            v-model="searchQuery"
            placeholder="Cari nama, periode..."
            @search="handleSearch"
          />
        </div>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Nama laporan
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Jenis
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Periode
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Tanggal
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Ukuran
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Aksi
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="report in paginatedReports" :key="report.id" class="hover:bg-gray-50">
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
                  class="flex items-center gap-2 px-3 py-1.5 text-green-600 hover:bg-green-50 rounded-lg font-medium text-sm transition-colors"
                >
                  <Download class="w-4 h-4" />
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

    <div class="bg-gray-50 border border-gray-200 rounded-xl p-6">
      <h3 class="font-semibold text-gray-900 mb-4">Opsi ekspor kustom</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="bg-white border border-gray-200 rounded-lg p-4">
          <p class="text-sm font-medium text-gray-900 mb-2">Rentang tanggal</p>
          <input
            type="date"
            class="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm"
            value="2025-01-01"
          />
        </div>
        <div class="bg-white border border-gray-200 rounded-lg p-4">
          <p class="text-sm font-medium text-gray-900 mb-2">Jenis laporan</p>
          <select class="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm">
            <option>Semua jenis</option>
            <option>Hanya keuangan</option>
            <option>Hanya investasi</option>
            <option>Hanya kesehatan tanaman</option>
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
            Buat laporan
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
