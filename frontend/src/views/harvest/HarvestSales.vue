<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { BadgeDollarSign, Calendar, Check, Package, RefreshCcw, ShoppingCart, X } from "lucide-vue-next";
import { listWarehouseStocks } from "../../services/warehouse/warehouse";
import { createSalesOrder, listSalesOrders } from "../../services/sales/sales";
import SearchBar from "../../components/common/SearchBar.vue";

type WarehouseStockRow = {
  id: number;
  plant_batch_id: number;
  batch_code: string;
  package_name: string;
  user_id: string;
  user_name: string;
  grade_id: number | null;
  grade_name: string;
  total_quantity: number;
  unit: string;
  updated_at: string;
};

type SalesDetailApiRow = {
  id: number;
  warehouse_stock_id: number;
  plant_batch_id: number;
  batch_code: string;
  package_name: string;
  grade_name: string;
  quantity: number;
  unit_price: number;
  total_price: number;
};

type SalesOrderApiRow = {
  id: number;
  order_number: string;
  sales_date: string;
  buyer: string;
  total_amount: number;
  status: string;
  created_by: string;
  created_at: string;
  details: SalesDetailApiRow[];
};

type SalesDetailRow = {
  id: number;
  batchCode: string;
  packageName: string;
  gradeName: string;
  quantity: number;
  unitPrice: number;
  totalPrice: number;
};

type SalesOrderRow = {
  id: number;
  orderNumber: string;
  salesDate: string;
  buyer: string;
  totalAmount: number;
  status: string;
  createdAt: string;
  details: SalesDetailRow[];
};

const stocks = ref<WarehouseStockRow[]>([]);
const orders = ref<SalesOrderRow[]>([]);
const loading = ref(false);
const errorMessage = ref("");
const infoMessage = ref("");
const stockSearch = ref("");
const stockGradeFilter = ref("all");

const saleModalOpen = ref(false);
const submittingSale = ref(false);
const selectedStock = ref<WarehouseStockRow | null>(null);

const today = new Date().toISOString().slice(0, 10);

const saleForm = reactive({
  buyer: "",
  salesDate: today,
  quantity: 0,
  unitPrice: 0,
});

const availableStocks = computed(() => stocks.value.filter((item) => Number(item.total_quantity) > 0));

const availableStockGrades = computed(() => {
  const grades = new Set<string>();
  availableStocks.value.forEach((item) => {
    grades.add(item.grade_name || "Basah");
  });
  return Array.from(grades).sort((a, b) => a.localeCompare(b));
});

const filteredAvailableStocks = computed(() => {
  const keyword = stockSearch.value.trim().toLowerCase();

  return availableStocks.value.filter((item) => {
    const gradeName = item.grade_name || "Basah";
    const matchesGrade = stockGradeFilter.value === "all" || gradeName === stockGradeFilter.value;
    const searchableText = [
      item.batch_code,
      item.package_name,
      item.user_name,
      item.user_id,
      gradeName,
      item.unit,
    ]
      .filter(Boolean)
      .join(" ")
      .toLowerCase();

    const matchesSearch = !keyword || searchableText.includes(keyword);
    return matchesGrade && matchesSearch;
  });
});

const totalStockCount = computed(() => availableStocks.value.length);
const totalStockQuantity = computed(() =>
  availableStocks.value.reduce((acc, item) => acc + (Number(item.total_quantity) || 0), 0)
);
const totalRevenue = computed(() => orders.value.reduce((acc, item) => acc + (Number(item.totalAmount) || 0), 0));
const salesThisMonth = computed(() => {
  const now = new Date();
  return orders.value.filter((item) => {
    const date = new Date(item.salesDate);
    if (Number.isNaN(date.getTime())) return false;
    return date.getFullYear() === now.getFullYear() && date.getMonth() === now.getMonth();
  }).length;
});

const formatDate = (value: string) => {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "-";
  return date.toLocaleDateString("id-ID", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });
};

const formatRupiah = (value: number) =>
  new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(value || 0);

const formatQuantity = (value: number) =>
  new Intl.NumberFormat("id-ID", {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(value || 0);

function resetSaleForm(stock: WarehouseStockRow) {
  saleForm.buyer = "";
  saleForm.salesDate = today;
  saleForm.quantity = Number(stock.total_quantity) || 0;
  saleForm.unitPrice = 0;
}

function openSaleModal(stock: WarehouseStockRow) {
  selectedStock.value = stock;
  resetSaleForm(stock);
  errorMessage.value = "";
  infoMessage.value = "";
  saleModalOpen.value = true;
}

function closeSaleModal() {
  if (submittingSale.value) return;
  saleModalOpen.value = false;
  selectedStock.value = null;
}

async function fetchStocks() {
  const response = await listWarehouseStocks();
  const rows = Array.isArray(response.data) ? (response.data as WarehouseStockRow[]) : [];
  stocks.value = rows;
}

async function fetchOrders() {
  const response = await listSalesOrders();
  const rows = Array.isArray(response.data) ? (response.data as SalesOrderApiRow[]) : [];
  orders.value = rows.map((item) => ({
    id: item.id,
    orderNumber: item.order_number,
    salesDate: item.sales_date,
    buyer: item.buyer,
    totalAmount: Number(item.total_amount) || 0,
    status: item.status,
    createdAt: item.created_at,
    details: Array.isArray(item.details)
      ? item.details.map((detail) => ({
          id: detail.id,
          batchCode: detail.batch_code,
          packageName: detail.package_name,
          gradeName: detail.grade_name,
          quantity: Number(detail.quantity) || 0,
          unitPrice: Number(detail.unit_price) || 0,
          totalPrice: Number(detail.total_price) || 0,
        }))
      : [],
  }));
}

async function refreshData() {
  loading.value = true;
  errorMessage.value = "";
  infoMessage.value = "";
  try {
    await Promise.all([fetchStocks(), fetchOrders()]);
  } catch (error) {
    console.error("Failed to load sales data", error);
    errorMessage.value = "Gagal memuat data penjualan dari server.";
  } finally {
    loading.value = false;
  }
}

async function submitSale() {
  if (!selectedStock.value) return;

  errorMessage.value = "";
  infoMessage.value = "";

  if (!saleForm.buyer.trim()) {
    errorMessage.value = "Nama pembeli wajib diisi.";
    return;
  }

  if (!saleForm.salesDate) {
    errorMessage.value = "Tanggal penjualan wajib diisi.";
    return;
  }

  if (!Number.isFinite(saleForm.quantity) || saleForm.quantity <= 0) {
    errorMessage.value = "Kuantitas penjualan harus lebih dari 0.";
    return;
  }

  if (saleForm.quantity > Number(selectedStock.value.total_quantity || 0)) {
    errorMessage.value = "Kuantitas melebihi stok yang tersedia.";
    return;
  }

  if (!Number.isFinite(saleForm.unitPrice) || saleForm.unitPrice <= 0) {
    errorMessage.value = "Harga per kg harus lebih dari 0.";
    return;
  }

  submittingSale.value = true;
  try {
    await createSalesOrder({
      sales_date: saleForm.salesDate,
      buyer: saleForm.buyer.trim(),
      details: [
        {
          warehouse_stock_id: selectedStock.value.id,
          quantity: Number(saleForm.quantity),
          unit_price: Number(saleForm.unitPrice),
        },
      ],
    });

    await refreshData();

    infoMessage.value = "Penjualan berhasil diproses.";
    closeSaleModal();
  } catch (error) {
    console.error("Failed to create sales order", error);
    errorMessage.value = "Gagal memproses penjualan.";
  } finally {
    submittingSale.value = false;
  }
}

const orderTotalForForm = computed(() => (Number(saleForm.quantity) || 0) * (Number(saleForm.unitPrice) || 0));

onMounted(async () => {
  await refreshData();
});
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-wrap items-start justify-between gap-4">
      <div>
        <h2 class="text-2xl font-semibold text-gray-900">Penjualan panen</h2>
        <p class="text-gray-600 mt-1">Kelola penjualan stok gudang dan pantau transaksi yang sudah terjadi.</p>
      </div>
      <button
        type="button"
        class="inline-flex items-center gap-2 rounded-lg border border-gray-300 px-3 py-2 text-sm font-medium text-gray-700 transition hover:bg-gray-50"
        @click="refreshData"
      >
        <RefreshCcw class="h-4 w-4" />
        Refresh data
      </button>
    </div>

    <div v-if="errorMessage" class="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
      {{ errorMessage }}
    </div>
    <div v-if="infoMessage" class="rounded-lg border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
      {{ infoMessage }}
    </div>

    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Stok siap dijual</p>
          <Package class="w-5 h-5 text-emerald-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">{{ totalStockCount }}</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Total kuantitas</p>
          <Package class="w-5 h-5 text-blue-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">{{ formatQuantity(totalStockQuantity) }} kg</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Total pendapatan</p>
          <BadgeDollarSign class="w-5 h-5 text-emerald-600" />
        </div>
        <p class="text-2xl font-semibold text-emerald-600">{{ formatRupiah(totalRevenue) }}</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Order bulan ini</p>
          <Calendar class="w-5 h-5 text-orange-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">{{ salesThisMonth }}</p>
      </div>
    </div>

    <section class="bg-white rounded-xl border border-gray-200 overflow-hidden">
      <div class="flex flex-wrap items-center justify-between gap-3 border-b border-gray-200 px-6 py-4">
        <div>
          <h3 class="text-lg font-semibold text-gray-900">Stok gudang siap dijual</h3>
          <p class="mt-1 text-sm text-gray-600">
            Menampilkan {{ filteredAvailableStocks.length }} dari {{ availableStocks.length }} stok siap jual.
          </p>
        </div>
        <span v-if="loading" class="text-sm text-gray-500">Memuat data...</span>
      </div>
      <div class="flex flex-wrap items-center gap-3 border-b border-gray-100 px-6 py-4">
        <div class="min-w-[260px] flex-1">
          <SearchBar
            v-model="stockSearch"
            placeholder="Cari batch, investor, paket, atau mutu..."
          />
        </div>
        <select
          v-model="stockGradeFilter"
          class="rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-700 outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
        >
          <option value="all">Semua mutu</option>
          <option v-for="grade in availableStockGrades" :key="grade" :value="grade">
            {{ grade }}
          </option>
        </select>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full min-w-[900px]">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="text-left px-6 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Batch</th>
              <th class="text-left px-6 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Investor</th>
              <th class="text-left px-6 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Mutu</th>
              <th class="text-right px-6 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Stok (kg)</th>
              <th class="text-left px-6 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Update</th>
              <th class="text-center px-6 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-if="!filteredAvailableStocks.length">
              <td colspan="6" class="px-6 py-6 text-center text-sm text-gray-500">
                {{ availableStocks.length ? "Tidak ada stok yang cocok dengan pencarian atau filter." : "Belum ada stok gudang yang siap dijual." }}
              </td>
            </tr>
            <tr v-for="stock in filteredAvailableStocks" :key="stock.id" class="hover:bg-gray-50">
              <td class="px-6 py-4">
                <p class="font-medium text-gray-900">{{ stock.batch_code }}</p>
                <p class="text-xs text-gray-500">{{ stock.package_name || "Paket tanpa nama" }}</p>
              </td>
              <td class="px-6 py-4 text-sm text-gray-700">
                {{ stock.user_name || "-" }}
              </td>
              <td class="px-6 py-4 text-sm text-gray-700">
                {{ stock.grade_name || "Basah" }}
              </td>
              <td class="px-6 py-4 text-right font-medium text-gray-900">
                {{ formatQuantity(stock.total_quantity) }} {{ stock.unit || "kg" }}
              </td>
              <td class="px-6 py-4 text-sm text-gray-600">
                {{ formatDate(stock.updated_at) }}
              </td>
              <td class="px-6 py-4 text-center">
                <button
                  type="button"
                  class="inline-flex items-center gap-2 rounded-lg bg-emerald-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-emerald-700"
                  @click="openSaleModal(stock)"
                >
                  <ShoppingCart class="h-4 w-4" />
                  Proses penjualan
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <section class="bg-white rounded-xl border border-gray-200 overflow-hidden">
      <div class="border-b border-gray-200 px-6 py-4">
        <h3 class="text-lg font-semibold text-gray-900">Riwayat penjualan</h3>
        <p class="text-sm text-gray-600 mt-1">Rekap transaksi yang sudah tercatat di sistem.</p>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full min-w-[980px]">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="text-left px-6 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Order</th>
              <th class="text-left px-6 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Pembeli</th>
              <th class="text-left px-6 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Tanggal</th>
              <th class="text-left px-6 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Detail</th>
              <th class="text-right px-6 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Total</th>
              <th class="text-center px-6 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-if="!orders.length">
              <td colspan="6" class="px-6 py-6 text-center text-sm text-gray-500">
                Belum ada transaksi penjualan yang tercatat.
              </td>
            </tr>
            <tr v-for="order in orders" :key="order.id" class="hover:bg-gray-50">
              <td class="px-6 py-4">
                <p class="font-medium text-gray-900">{{ order.orderNumber }}</p>
                <p class="text-xs text-gray-500">{{ formatDate(order.createdAt) }}</p>
              </td>
              <td class="px-6 py-4 text-sm text-gray-700">{{ order.buyer }}</td>
              <td class="px-6 py-4 text-sm text-gray-600">{{ formatDate(order.salesDate) }}</td>
              <td class="px-6 py-4 text-sm text-gray-700">
                <div v-if="order.details.length" class="space-y-1">
                  <div v-for="detail in order.details" :key="detail.id" class="text-sm">
                    <span class="font-medium text-gray-900">{{ detail.batchCode }}</span>
                    <span class="text-xs text-gray-500">({{ detail.packageName }})</span>
                    - {{ detail.gradeName }}: {{ formatQuantity(detail.quantity) }} kg x {{ formatRupiah(detail.unitPrice) }}
                  </div>
                </div>
                <span v-else class="text-xs text-gray-500">-</span>
              </td>
              <td class="px-6 py-4 text-right font-semibold text-emerald-600">{{ formatRupiah(order.totalAmount) }}</td>
              <td class="px-6 py-4 text-center">
                <span class="inline-flex rounded-full bg-emerald-100 px-2.5 py-1 text-xs font-semibold text-emerald-700">
                  {{ order.status || "completed" }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <div v-if="saleModalOpen && selectedStock" class="fixed inset-0 z-50 bg-black/50 p-4">
      <div class="mx-auto mt-10 w-full max-w-2xl rounded-2xl bg-white p-6 shadow-xl">
        <div class="flex items-start justify-between gap-4">
          <div>
            <h3 class="text-xl font-semibold text-gray-900">Proses penjualan</h3>
            <p class="mt-1 text-sm text-gray-600">Masukkan detail transaksi sebelum konfirmasi penjualan.</p>
          </div>
          <button
            type="button"
            class="rounded-lg border border-gray-200 p-2 text-gray-500 transition hover:bg-gray-100"
            @click="closeSaleModal"
          >
            <X class="h-4 w-4" />
          </button>
        </div>

        <div class="mt-5 rounded-lg border border-gray-200 bg-gray-50 p-4 text-sm text-gray-700">
          <div class="flex flex-wrap items-center justify-between gap-3">
            <div>
              <p class="font-medium text-gray-900">{{ selectedStock.batch_code }}</p>
              <p class="text-xs text-gray-500">{{ selectedStock.package_name }}</p>
            </div>
            <div class="text-right">
              <p class="text-xs text-gray-500">Stok tersedia</p>
              <p class="font-semibold text-gray-900">
                {{ formatQuantity(selectedStock.total_quantity) }} {{ selectedStock.unit || "kg" }}
              </p>
            </div>
          </div>
          <div class="mt-3 flex flex-wrap gap-3 text-xs text-gray-600">
            <span class="rounded-full bg-white px-2.5 py-1">Investor: {{ selectedStock.user_name || "-" }}</span>
            <span class="rounded-full bg-white px-2.5 py-1">Mutu: {{ selectedStock.grade_name || "Basah" }}</span>
          </div>
        </div>

        <form class="mt-5 grid gap-4 md:grid-cols-2" @submit.prevent="submitSale">
          <label class="space-y-2 md:col-span-2">
            <span class="text-sm font-medium text-gray-700">Nama pembeli</span>
            <input
              v-model="saleForm.buyer"
              type="text"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
              placeholder="Contoh: PT Mitra Vanili"
            />
          </label>

          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Tanggal penjualan</span>
            <input
              v-model="saleForm.salesDate"
              type="date"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
            />
          </label>

          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Kuantitas (kg)</span>
            <input
              v-model.number="saleForm.quantity"
              type="number"
              min="0"
              step="0.1"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
            />
          </label>

          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Harga per kg</span>
            <input
              v-model.number="saleForm.unitPrice"
              type="number"
              min="0"
              step="1000"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
              placeholder="Contoh: 750000"
            />
          </label>

          <div class="md:col-span-2 rounded-lg border border-emerald-100 bg-emerald-50 px-4 py-3">
            <p class="text-xs font-semibold uppercase tracking-wide text-emerald-600">Total nilai penjualan</p>
            <p class="text-lg font-semibold text-emerald-700">{{ formatRupiah(orderTotalForForm) }}</p>
          </div>

          <div class="md:col-span-2 flex justify-end gap-2">
            <button
              type="button"
              class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
              @click="closeSaleModal"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="submittingSale"
              class="inline-flex items-center gap-2 rounded-lg bg-emerald-600 px-4 py-2 text-sm font-medium text-white hover:bg-emerald-700 disabled:opacity-70"
            >
              <Check class="h-4 w-4" />
              {{ submittingSale ? "Menyimpan..." : "Konfirmasi penjualan" }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
