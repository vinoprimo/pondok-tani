<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import {
  TrendingUp,
  Sprout,
  DollarSign,
  Users,
  ArrowUpRight,
  ArrowDownRight,
  AlertCircle,
  X,
  Calculator,
} from "lucide-vue-next";
import {
  createInvestmentPackage,
  deleteInvestmentPackage,
  getInvestmentPackages,
  updateInvestmentPackage,
} from "../../services/investment/package";

type UserRole = "investor" | "mitra" | "admin";

type InvestmentPackageItem = {
  id: number;
  package_name: string;
  description?: string | null;
  min_quantity: number;
  price: number;
  status: string;
};

const userRole = ref<UserRole>(
  (localStorage.getItem("userRole") as UserRole) || "investor"
);

const showProjectionModal = ref(false);

const adminPackages = ref<InvestmentPackageItem[]>([]);
const packageLoading = ref(false);
const packageSaving = ref(false);
const packageError = ref("");
const packageSuccess = ref("");
const editingPackageId = ref<number | null>(null);

const packageForm = reactive({
  package_name: "",
  description: "",
  min_quantity: 1,
  price: 0,
  status: "active",
});

const isAdmin = computed(() => userRole.value === "admin");

const formatRupiah = (value: number) =>
  new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(value);

function resetPackageForm() {
  packageForm.package_name = "";
  packageForm.description = "";
  packageForm.min_quantity = 1;
  packageForm.price = 0;
  packageForm.status = "active";
  editingPackageId.value = null;
}

async function loadAdminPackages() {
  if (!isAdmin.value) return;

  packageLoading.value = true;
  packageError.value = "";
  try {
    const res = await getInvestmentPackages();
    adminPackages.value = Array.isArray(res.data) ? res.data : [];
  } catch (err) {
    console.error("Failed to load investment packages", err);
    packageError.value = "Gagal memuat data paket investasi.";
  } finally {
    packageLoading.value = false;
  }
}

function startEditPackage(item: InvestmentPackageItem) {
  editingPackageId.value = item.id;
  packageForm.package_name = item.package_name;
  packageForm.description = item.description || "";
  packageForm.min_quantity = item.min_quantity;
  packageForm.price = item.price;
  packageForm.status = item.status || "active";
}

async function submitPackageForm() {
  packageError.value = "";
  packageSuccess.value = "";

  if (!packageForm.package_name.trim()) {
    packageError.value = "Nama paket wajib diisi.";
    return;
  }
  if (packageForm.min_quantity <= 0) {
    packageError.value = "Minimum kuantitas harus lebih dari 0.";
    return;
  }
  if (packageForm.price <= 0) {
    packageError.value = "Harga harus lebih dari 0.";
    return;
  }

  const payload = {
    package_name: packageForm.package_name.trim(),
    description: packageForm.description.trim() || null,
    min_quantity: Number(packageForm.min_quantity),
    price: Number(packageForm.price),
    status: packageForm.status,
  };

  packageSaving.value = true;
  try {
    if (editingPackageId.value) {
      await updateInvestmentPackage(editingPackageId.value, payload);
      packageSuccess.value = "Paket investasi berhasil diperbarui.";
    } else {
      await createInvestmentPackage(payload);
      packageSuccess.value = "Paket investasi berhasil dibuat.";
    }

    resetPackageForm();
    await loadAdminPackages();
  } catch (err) {
    console.error("Failed to save investment package", err);
    packageError.value = "Gagal menyimpan paket investasi.";
  } finally {
    packageSaving.value = false;
  }
}

async function removePackage(id: number) {
  packageError.value = "";
  packageSuccess.value = "";
  const confirmed = window.confirm("Hapus paket investasi ini?");
  if (!confirmed) return;

  try {
    await deleteInvestmentPackage(id);
    packageSuccess.value = "Paket investasi berhasil dihapus.";
    if (editingPackageId.value === id) {
      resetPackageForm();
    }
    await loadAdminPackages();
  } catch (err) {
    console.error("Failed to delete investment package", err);
    packageError.value = "Gagal menghapus paket investasi.";
  }
}

onMounted(() => {
  const stored = localStorage.getItem("userRole") as UserRole | null;
  if (stored === "admin" || stored === "investor" || stored === "mitra") {
    userRole.value = stored;
  }

  loadAdminPackages();
});

const monthlyData = [
  { month: "Jan", revenue: 12000, costs: 8000 },
  { month: "Peb", revenue: 15000, costs: 8500 },
  { month: "Mar", revenue: 18000, costs: 9000 },
  { month: "Apr", revenue: 22000, costs: 9500 },
  { month: "Mei", revenue: 25000, costs: 10000 },
  { month: "Jun", revenue: 28000, costs: 10500 },
];

const plantStatusData = [
  { name: "Berbunga", value: 35, color: "#10b981" },
  { name: "Tumbuh", value: 45, color: "#3b82f6" },
  { name: "Panen", value: 15, color: "#f59e0b" },
  { name: "Baru", value: 5, color: "#6366f1" },
];

const investorStats = [
  {
    label: "Total investasi",
    value: "$45,000",
    change: "+12%",
    isPositive: true,
    icon: DollarSign,
  },
  {
    label: "ROI saat ini",
    value: "24.5%",
    change: "+3.2%",
    isPositive: true,
    icon: TrendingUp,
  },
  {
    label: "Tanaman aktif",
    value: "85",
    change: "+5",
    isPositive: true,
    icon: Sprout,
  },
  {
    label: "Perkiraan imbal tahunan",
    value: "$11,025",
    change: "+15%",
    isPositive: true,
    icon: DollarSign,
  },
];

const adminStats = [
  {
    label: "Total pendapatan",
    value: "$128,450",
    change: "+18%",
    isPositive: true,
    icon: DollarSign,
  },
  {
    label: "Investor aktif",
    value: "234",
    change: "+12",
    isPositive: true,
    icon: Users,
  },
  {
    label: "Total tanaman",
    value: "1,850",
    change: "+45",
    isPositive: true,
    icon: Sprout,
  },
  {
    label: "ROI rata-rata",
    value: "22.8%",
    change: "-0.8%",
    isPositive: false,
    icon: TrendingUp,
  },
];

const harvestStatus = [
  { status: "Siap panen", count: 35, color: "#10b981" },
  { status: "Pengeringan", count: 25, color: "#f59e0b" },
  { status: "Terjual", count: 15, color: "#3b82f6" },
];

const revenueComparison = [
  { month: "Jan", estimated: 8000, actual: 7500 },
  { month: "Peb", estimated: 9000, actual: 9200 },
  { month: "Mar", estimated: 10000, actual: 9800 },
  { month: "Apr", estimated: 11000, actual: 11500 },
  { month: "Mei", estimated: 12000, actual: 11800 },
  { month: "Jun", estimated: 13000, actual: 13200 },
];

const salesHistory = [
  {
    id: "SALE-089",
    date: "2024-01-19",
    grade: "A",
    quantity: 40,
    revenue: "$200,000",
    buyer: "PT Ekspor Premium",
  },
  {
    id: "SALE-088",
    date: "2024-01-17",
    grade: "B",
    quantity: 25,
    revenue: "$106,250",
    buyer: "Distributor Lokal",
  },
  {
    id: "SALE-087",
    date: "2024-01-15",
    grade: "A",
    quantity: 35,
    revenue: "$175,000",
    buyer: "PT Impor Global",
  },
];

const recentActivity = [
  { action: "Tanaman mencapai fase berbunga", time: "2 jam lalu", type: "success" as const },
  { action: "Imbal bulanan $918 telah dicairkan", time: "5 jam lalu", type: "info" as const },
  { action: "Inspeksi kesehatan tanaman selesai", time: "1 hari lalu", type: "success" as const },
  { action: "Siklus panen batch #34 dimulai", time: "2 hari lalu", type: "warning" as const },
];

const stats = computed(() => (userRole.value === "admin" ? adminStats : investorStats));

const maxMonthly = computed(() =>
  Math.max(...monthlyData.map((d) => Math.max(d.revenue, d.costs)), 1)
);

const maxRevenueCompare = computed(() =>
  Math.max(...revenueComparison.map((d) => Math.max(d.estimated, d.actual)), 1)
);

function activityDotClass(type: string) {
  if (type === "success") return "bg-green-500";
  if (type === "warning") return "bg-yellow-500";
  return "bg-blue-500";
}

const harvestQty = ref("");
const harvestType = ref<"wet" | "dry">("wet");
const grade = ref("");
const province = ref("");
const shrinkage = ref("75");
const operationalCost = ref("");

const priceData: Record<string, { wet: number; dry: Record<string, number> }> = {
  "West Java": {
    wet: 750000,
    dry: { A: 5000000, B: 4250000, C: 3500000, Split: 2500000, Asalan: 2000000 },
  },
  "Central Java": {
    wet: 725000,
    dry: { A: 4800000, B: 4080000, C: 3360000, Split: 2400000, Asalan: 1920000 },
  },
  "East Java": {
    wet: 780000,
    dry: { A: 5200000, B: 4420000, C: 3640000, Split: 2600000, Asalan: 2080000 },
  },
  Bali: {
    wet: 800000,
    dry: { A: 5400000, B: 4590000, C: 3780000, Split: 2700000, Asalan: 2160000 },
  },
};

const projection = computed(() => {
  if (!harvestQty.value || !province.value) return null;
  const qty = parseFloat(harvestQty.value);
  const provinceData = priceData[province.value];
  if (!provinceData || Number.isNaN(qty)) return null;

  let revenue = 0;
  let effectiveQty = qty;

  if (harvestType.value === "wet") {
    revenue = qty * provinceData.wet;
  } else {
    if (!grade.value) return null;
    const shrinkagePercent = parseFloat(shrinkage.value) / 100;
    if (Number.isNaN(shrinkagePercent)) return null;
    effectiveQty = qty * (1 - shrinkagePercent);
    const dryPrice = provinceData.dry[grade.value];
    if (dryPrice === undefined) return null;
    revenue = effectiveQty * dryPrice;
  }

  const costs = operationalCost.value
    ? parseFloat(operationalCost.value)
    : revenue * 0.15;
  const netProfit = revenue - costs;
  const roi = costs !== 0 ? (netProfit / costs) * 100 : 0;

  return { revenue, costs, netProfit, roi, effectiveQty };
});

function closeModal() {
  showProjectionModal.value = false;
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-semibold text-gray-900">
        {{ userRole === "admin" ? "Dasbor admin" : userRole === "mitra" ? "Dasbor mitra" : "Ringkasan investasi" }}
      </h2>
      <p class="text-gray-600 mt-1">
        {{
          userRole === "admin"
            ? "Kelola operasional perkebunan vanili Anda"
            : userRole === "mitra"
              ? "Pantau aktivitas operasional kebun vanili Anda"
              : "Pantau investasi perkebunan vanili Anda"
        }}
      </p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div
        v-for="(stat, index) in stats"
        :key="index"
        class="bg-white rounded-xl border border-gray-200 p-6"
      >
        <div class="flex items-center justify-between mb-4">
          <div class="w-12 h-12 bg-green-50 rounded-lg flex items-center justify-center">
            <component :is="stat.icon" class="w-6 h-6 text-green-600" />
          </div>
          <div
            :class="[
              'flex items-center gap-1 text-sm',
              stat.isPositive ? 'text-green-600' : 'text-red-600',
            ]"
          >
            <ArrowUpRight v-if="stat.isPositive" class="w-4 h-4" />
            <ArrowDownRight v-else class="w-4 h-4" />
            <span>{{ stat.change }}</span>
          </div>
        </div>
        <h3 class="text-2xl font-semibold text-gray-900 mb-1">{{ stat.value }}</h3>
        <p class="text-sm text-gray-600">{{ stat.label }}</p>
      </div>
    </div>

    <template v-if="isAdmin">
      <div class="bg-white rounded-xl border border-gray-200 p-6 space-y-5">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-lg font-semibold text-gray-900">CRUD Paket Investasi</h3>
            <p class="text-sm text-gray-600 mt-1">
              Kelola paket investasi yang akan tampil di landing page.
            </p>
          </div>
          <button
            type="button"
            class="px-4 py-2 border border-gray-300 rounded-lg text-sm font-medium hover:bg-gray-50"
            @click="loadAdminPackages"
          >
            Refresh Data
          </button>
        </div>

        <p v-if="packageError" class="text-sm text-red-600">{{ packageError }}</p>
        <p v-if="packageSuccess" class="text-sm text-emerald-600">{{ packageSuccess }}</p>

        <form class="grid grid-cols-1 md:grid-cols-2 gap-4" @submit.prevent="submitPackageForm">
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-1">Nama Paket</label>
            <input
              v-model="packageForm.package_name"
              type="text"
              class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-green-500"
              placeholder="Contoh: Growth Vanili"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Harga (IDR)</label>
            <input
              v-model.number="packageForm.price"
              type="number"
              min="1"
              class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-green-500"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Min Quantity</label>
            <input
              v-model.number="packageForm.min_quantity"
              type="number"
              min="1"
              class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-green-500"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
            <select
              v-model="packageForm.status"
              class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-green-500"
            >
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
            </select>
          </div>

          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-1">Deskripsi</label>
            <textarea
              v-model="packageForm.description"
              rows="3"
              class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-green-500"
              placeholder="Deskripsi singkat paket investasi"
            />
          </div>

          <div class="md:col-span-2 flex items-center gap-3">
            <button
              type="submit"
              class="px-5 py-2.5 rounded-lg bg-green-600 text-white font-medium hover:bg-green-700 disabled:opacity-50"
              :disabled="packageSaving"
            >
              {{ editingPackageId ? "Update Paket" : "Tambah Paket" }}
            </button>
            <button
              v-if="editingPackageId"
              type="button"
              class="px-5 py-2.5 rounded-lg border border-gray-300 font-medium hover:bg-gray-50"
              @click="resetPackageForm"
            >
              Batal Edit
            </button>
          </div>
        </form>

        <div class="overflow-x-auto rounded-lg border border-gray-100">
          <table class="w-full text-sm">
            <thead class="bg-gray-50 text-left text-gray-600">
              <tr>
                <th class="px-3 py-2 font-medium">Nama Paket</th>
                <th class="px-3 py-2 font-medium">Harga</th>
                <th class="px-3 py-2 font-medium">Min Qty</th>
                <th class="px-3 py-2 font-medium">Status</th>
                <th class="px-3 py-2 font-medium">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-if="packageLoading">
                <td colspan="5" class="px-3 py-4 text-center text-gray-500">Memuat data paket...</td>
              </tr>
              <tr v-else-if="adminPackages.length === 0">
                <td colspan="5" class="px-3 py-4 text-center text-gray-500">Belum ada paket investasi.</td>
              </tr>
              <tr v-for="item in adminPackages" :key="item.id">
                <td class="px-3 py-2 font-medium text-gray-900">{{ item.package_name }}</td>
                <td class="px-3 py-2">{{ formatRupiah(item.price) }}</td>
                <td class="px-3 py-2">{{ item.min_quantity }}</td>
                <td class="px-3 py-2">
                  <span
                    class="px-2 py-1 rounded-full text-xs font-medium"
                    :class="item.status === 'active' ? 'bg-emerald-100 text-emerald-700' : 'bg-gray-100 text-gray-600'"
                  >
                    {{ item.status }}
                  </span>
                </td>
                <td class="px-3 py-2">
                  <div class="flex items-center gap-2">
                    <button
                      type="button"
                      class="px-3 py-1.5 rounded-lg border border-gray-300 text-xs font-medium hover:bg-gray-50"
                      @click="startEditPackage(item)"
                    >
                      Edit
                    </button>
                    <button
                      type="button"
                      class="px-3 py-1.5 rounded-lg border border-red-200 text-red-600 text-xs font-medium hover:bg-red-50"
                      @click="removePackage(item.id)"
                    >
                      Hapus
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="lg:col-span-2 bg-white rounded-xl border border-gray-200 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Kinerja investasi</h3>
        <p class="text-xs text-gray-500 mb-3">Grafik batang (placeholder) — data di bawah</p>
        <div class="overflow-x-auto rounded-lg border border-gray-100">
          <table class="w-full text-sm">
            <thead class="bg-gray-50 text-left text-gray-600">
              <tr>
                <th class="px-3 py-2 font-medium">Bulan</th>
                <th class="px-3 py-2 font-medium text-right">Pendapatan</th>
                <th class="px-3 py-2 font-medium text-right">Biaya</th>
                <th class="px-3 py-2 font-medium">Batang pendapatan</th>
                <th class="px-3 py-2 font-medium">Batang biaya</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-for="row in monthlyData" :key="row.month">
                <td class="px-3 py-2 font-medium text-gray-900">{{ row.month }}</td>
                <td class="px-3 py-2 text-right text-gray-900">${{ row.revenue.toLocaleString() }}</td>
                <td class="px-3 py-2 text-right text-gray-900">${{ row.costs.toLocaleString() }}</td>
                <td class="px-3 py-2 w-32">
                  <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
                    <div
                      class="h-full rounded-full bg-emerald-500"
                      :style="{ width: `${(row.revenue / maxMonthly) * 100}%` }"
                    />
                  </div>
                </td>
                <td class="px-3 py-2 w-32">
                  <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
                    <div
                      class="h-full rounded-full bg-blue-500"
                      :style="{ width: `${(row.costs / maxMonthly) * 100}%` }"
                    />
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Status tanaman</h3>
        <p class="text-xs text-gray-500 mb-3">Grafik pie (placeholder) — distribusi di bawah</p>
        <div class="space-y-3 mb-4">
          <div v-for="(item, index) in plantStatusData" :key="index">
            <div class="flex items-center justify-between text-sm mb-1">
              <span class="text-gray-600">{{ item.name }}</span>
              <span class="font-medium text-gray-900">{{ item.value }}%</span>
            </div>
            <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
              <div
                class="h-full rounded-full"
                :style="{ width: `${item.value}%`, backgroundColor: item.color }"
              />
            </div>
          </div>
        </div>
        <div class="mt-4 space-y-2">
          <div
            v-for="(item, index) in plantStatusData"
            :key="`leg-${index}`"
            class="flex items-center justify-between text-sm"
          >
            <div class="flex items-center gap-2">
              <div class="w-3 h-3 rounded-full" :style="{ backgroundColor: item.color }" />
              <span class="text-gray-600">{{ item.name }}</span>
            </div>
            <span class="font-medium text-gray-900">{{ item.value }}%</span>
          </div>
        </div>
      </div>
    </div>

    <template v-if="userRole === 'investor'">
      <div class="bg-gradient-to-r from-orange-50 to-yellow-50 rounded-xl border-2 border-orange-200 p-6">
        <div class="flex items-center justify-between mb-4">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-orange-500 rounded-full flex items-center justify-center">
              <AlertCircle class="w-6 h-6 text-white" />
            </div>
            <div>
              <h3 class="text-lg font-semibold text-gray-900">Tugas perawatan mendatang</h3>
              <p class="text-sm text-gray-600">Jangan lewatkan jadwal aktivitas Anda</p>
            </div>
          </div>
          <button
            type="button"
            class="px-4 py-2 bg-white text-orange-600 border border-orange-300 rounded-lg hover:bg-orange-50 transition-colors font-medium text-sm"
          >
            Lihat semua
          </button>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="bg-white rounded-lg p-4 border-2 border-red-300">
            <div class="flex items-start justify-between mb-2">
              <div>
                <p class="font-semibold text-gray-900">Penyiraman - BATCH-038</p>
                <p class="text-sm text-gray-600 mt-1">Jatuh tempo: 2024-01-20</p>
              </div>
              <span class="px-3 py-1 bg-red-100 text-red-700 rounded-full text-xs font-bold">
                Terlambat 2 hari
              </span>
            </div>
            <p class="text-sm text-red-600 mt-2">Perlu tindakan segera</p>
          </div>
          <div class="bg-white rounded-lg p-4 border-2 border-orange-200">
            <div class="flex items-start justify-between mb-2">
              <div>
                <p class="font-semibold text-gray-900">Pemupukan - BATCH-045</p>
                <p class="text-sm text-gray-600 mt-1">Jatuh tempo: 2024-01-26</p>
              </div>
              <span class="px-3 py-1 bg-orange-100 text-orange-700 rounded-full text-xs font-bold">
                Tersisa 3 hari
              </span>
            </div>
            <p class="text-sm text-gray-600 mt-2">Segera datang</p>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div class="bg-gradient-to-br from-blue-50 to-blue-100 rounded-xl border-2 border-blue-200 p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-2">Proyeksi pendapatan</h3>
          <p class="text-sm text-gray-600 mb-4">
            Perkirakan pendapatan dan ROI untuk panen Anda
          </p>
          <button
            type="button"
            class="w-full px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium flex items-center justify-center gap-2"
            @click="showProjectionModal = true"
          >
            <DollarSign class="w-5 h-5" />
            Buat proyeksi pendapatan
          </button>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Status panen</h3>
          <div class="space-y-3">
            <div
              v-for="(item, index) in harvestStatus"
              :key="index"
              class="flex items-center justify-between p-3 rounded-lg border border-gray-200"
            >
              <div class="flex items-center gap-3">
                <div class="w-4 h-4 rounded-full" :style="{ backgroundColor: item.color }" />
                <span class="text-sm font-medium text-gray-900">{{ item.status }}</span>
              </div>
              <span class="text-lg font-semibold" :style="{ color: item.color }">{{
                item.count
              }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-gray-900">Perkiraan vs pendapatan aktual</h3>
          <div class="flex items-center gap-4 text-sm">
            <div class="flex items-center gap-2">
              <div class="w-3 h-3 rounded-full bg-blue-500" />
              <span class="text-gray-600">Perkiraan</span>
            </div>
            <div class="flex items-center gap-2">
              <div class="w-3 h-3 rounded-full bg-green-500" />
              <span class="text-gray-600">Aktual</span>
            </div>
          </div>
        </div>
        <p class="text-xs text-gray-500 mb-3">Grafik garis (placeholder) — seri di bawah</p>
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead class="bg-gray-50 text-left">
              <tr>
                <th class="px-3 py-2 font-medium text-gray-600">Bulan</th>
                <th class="px-3 py-2 font-medium text-gray-600 text-right">Perkiraan</th>
                <th class="px-3 py-2 font-medium text-gray-600 text-right">Aktual</th>
                <th class="px-3 py-2 font-medium text-gray-600">Bandingkan</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-for="row in revenueComparison" :key="row.month">
                <td class="px-3 py-2 font-medium">{{ row.month }}</td>
                <td class="px-3 py-2 text-right text-blue-600">${{ row.estimated.toLocaleString() }}</td>
                <td class="px-3 py-2 text-right text-emerald-600">${{ row.actual.toLocaleString() }}</td>
                <td class="px-3 py-2">
                  <div class="flex gap-1 h-8 items-end">
                    <div
                      class="w-3 bg-blue-500 rounded-t"
                      :style="{ height: `${(row.estimated / maxRevenueCompare) * 100}%` }"
                    />
                    <div
                      class="w-3 bg-emerald-500 rounded-t"
                      :style="{ height: `${(row.actual / maxRevenueCompare) * 100}%` }"
                    />
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Kemajuan ROI</h3>
        <div class="space-y-4">
          <div>
            <div class="flex items-center justify-between mb-2">
              <span class="text-sm text-gray-600">ROI saat ini</span>
              <span class="text-lg font-semibold text-green-600">24.5%</span>
            </div>
            <div class="relative h-4 bg-gray-100 rounded-full overflow-hidden">
              <div
                class="absolute inset-y-0 left-0 bg-gradient-to-r from-green-500 to-green-600 rounded-full"
                style="width: 61.25%"
              />
            </div>
            <div class="flex justify-between mt-1 text-xs text-gray-500">
              <span>0%</span>
              <span>Target: 40%</span>
            </div>
          </div>
          <div class="flex items-center justify-between p-4 bg-green-50 rounded-lg border border-green-200">
            <span class="text-sm text-gray-700">Perkiraan capai target dalam:</span>
            <span class="text-lg font-semibold text-green-600">8 bulan</span>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Riwayat penjualan</h3>
            <p class="text-sm text-gray-600 mt-1">Penjualan panen dan pendapatan terbaru</p>
          </div>
          <button
            type="button"
            class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors text-sm font-medium flex items-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
              />
            </svg>
            Unduh laporan
          </button>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-gray-50 border-b border-gray-200">
              <tr>
                <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">ID penjualan</th>
                <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Tanggal</th>
                <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Mutu</th>
                <th class="text-right px-6 py-3 text-sm font-medium text-gray-900">Kuantitas</th>
                <th class="text-right px-6 py-3 text-sm font-medium text-gray-900">Pendapatan</th>
                <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Pembeli</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-for="sale in salesHistory" :key="sale.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 font-medium text-gray-900">{{ sale.id }}</td>
                <td class="px-6 py-4 text-gray-900">{{ sale.date }}</td>
                <td class="px-6 py-4">
                  <span
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800"
                  >
                    Mutu {{ sale.grade }}
                  </span>
                </td>
                <td class="px-6 py-4 text-right font-medium text-gray-900">{{ sale.quantity }} kg</td>
                <td class="px-6 py-4 text-right font-semibold text-green-600">{{ sale.revenue }}</td>
                <td class="px-6 py-4 text-gray-900">{{ sale.buyer }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <div class="bg-white rounded-xl border border-gray-200 p-6">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">Aktivitas terbaru</h3>
      <div class="space-y-4">
        <div
          v-for="(activity, index) in recentActivity"
          :key="index"
          class="flex items-start gap-4 pb-4 border-b border-gray-100 last:border-0"
        >
          <div :class="['w-2 h-2 rounded-full mt-2', activityDotClass(activity.type)]" />
          <div class="flex-1">
            <p class="text-gray-900">{{ activity.action }}</p>
            <p class="text-sm text-gray-500 mt-1">{{ activity.time }}</p>
          </div>
        </div>
      </div>
    </div>

    <Teleport to="body">
      <div
        v-if="showProjectionModal"
        class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
        role="dialog"
        aria-modal="true"
        @click.self="closeModal"
      >
        <div class="bg-white rounded-2xl max-w-3xl w-full max-h-[90vh] overflow-y-auto">
          <div class="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-blue-50 rounded-lg flex items-center justify-center">
                <Calculator class="w-6 h-6 text-blue-600" />
              </div>
              <div>
                <h2 class="text-xl font-semibold text-gray-900">Kalkulator proyeksi pendapatan</h2>
                <p class="text-sm text-gray-600">Perkirakan pendapatan panen dan ROI Anda</p>
              </div>
            </div>
            <button
              type="button"
              class="w-8 h-8 rounded-lg hover:bg-gray-100 flex items-center justify-center transition-colors"
              @click="closeModal"
            >
              <X class="w-5 h-5 text-gray-500" />
            </button>
          </div>
          <div class="p-6 space-y-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Perkiraan jumlah panen (kg)
                </label>
                <input
                  v-model="harvestQty"
                  type="number"
                  placeholder="mis. 100"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Jenis panen</label>
                <div class="grid grid-cols-2 gap-2">
                  <button
                    type="button"
                    :class="[
                      'px-4 py-2 rounded-lg font-medium transition-colors',
                      harvestType === 'wet'
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200',
                    ]"
                    @click="harvestType = 'wet'"
                  >
                    Basah
                  </button>
                  <button
                    type="button"
                    :class="[
                      'px-4 py-2 rounded-lg font-medium transition-colors',
                      harvestType === 'dry'
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200',
                    ]"
                    @click="harvestType = 'dry'"
                  >
                    Kering
                  </button>
                </div>
              </div>
              <div v-if="harvestType === 'dry'">
                <label class="block text-sm font-medium text-gray-700 mb-2">Pilihan mutu</label>
                <select
                  v-model="grade"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  <option value="">Pilih mutu</option>
                  <option value="A">Mutu A — Premium</option>
                  <option value="B">Mutu B — Standar</option>
                  <option value="C">Mutu C — Ekonomi</option>
                  <option value="Split">Split</option>
                  <option value="Asalan">Asalan</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Provinsi (acuan harga)
                </label>
                <select
                  v-model="province"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  <option value="">Pilih provinsi</option>
                  <option value="West Java">Jawa Barat</option>
                  <option value="Central Java">Jawa Tengah</option>
                  <option value="East Java">Jawa Timur</option>
                  <option value="Bali">Bali</option>
                </select>
              </div>
              <div v-if="harvestType === 'dry'">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Perkiraan susut pengeringan (%)
                </label>
                <input
                  v-model="shrinkage"
                  type="number"
                  placeholder="75"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Perkiraan biaya operasional (opsional)
                </label>
                <input
                  v-model="operationalCost"
                  type="number"
                  placeholder="Otomatis: 15% dari pendapatan"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
            </div>
            <div
              v-if="projection"
              class="bg-gradient-to-br from-green-50 to-blue-50 rounded-xl p-6 border-2 border-green-200"
            >
              <div class="flex items-center gap-2 mb-4">
                <TrendingUp class="w-5 h-5 text-green-600" />
                <h3 class="text-lg font-semibold text-gray-900">Hasil proyeksi</h3>
              </div>
              <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div class="bg-white rounded-lg p-4">
                  <p class="text-sm text-gray-600 mb-1">Perkiraan pendapatan</p>
                  <p class="text-2xl font-bold text-green-600">
                    ${{ projection.revenue.toLocaleString() }}
                  </p>
                </div>
                <div class="bg-white rounded-lg p-4">
                  <p class="text-sm text-gray-600 mb-1">Biaya operasional</p>
                  <p class="text-2xl font-bold text-orange-600">
                    ${{ projection.costs.toLocaleString() }}
                  </p>
                </div>
                <div class="bg-white rounded-lg p-4">
                  <p class="text-sm text-gray-600 mb-1">Laba bersih</p>
                  <p class="text-2xl font-bold text-green-600">
                    ${{ projection.netProfit.toLocaleString() }}
                  </p>
                </div>
                <div class="bg-white rounded-lg p-4">
                  <p class="text-sm text-gray-600 mb-1">Proyeksi ROI</p>
                  <p class="text-2xl font-bold text-blue-600">{{ projection.roi.toFixed(1) }}%</p>
                </div>
              </div>
              <div
                v-if="harvestType === 'dry'"
                class="mt-4 p-3 bg-white/50 rounded-lg border border-green-200"
              >
                <p class="text-sm text-gray-700">
                  <span class="font-medium">Catatan:</span> Setelah susut {{ shrinkage }}%, jumlah efektif:
                  <span class="font-semibold">{{ projection.effectiveQty.toFixed(2) }} kg</span>
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
