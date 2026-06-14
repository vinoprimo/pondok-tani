<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import {
  createInvestmentPackage,
  deleteInvestmentPackage,
  getInvestmentPackages,
  updateInvestmentPackage,
} from "../../services/investment/package";
import SearchBar from "../../components/common/SearchBar.vue";
import Pagination from "../../components/common/Pagination.vue";

type InvestmentPackageItem = {
  id: number;
  package_name: string;
  description?: string | null;
  min_quantity: number;
  price: number;
  status: string;
  roi?: string;
  duration?: string;
  benefits?: string[];
  is_popular?: boolean;
};

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
  roi: "",
  duration: "",
  benefits: "",
  is_popular: false,
});

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
  packageForm.roi = "";
  packageForm.duration = "";
  packageForm.benefits = "";
  packageForm.is_popular = false;
  editingPackageId.value = null;
}

const currentPage = ref(1);
const totalPages = ref(1);
const totalRows = ref(0);
const searchQuery = ref("");
const limit = ref(10);

async function loadAdminPackages() {
  packageLoading.value = true;
  packageError.value = "";
  try {
    const res = await getInvestmentPackages({ 
      page: currentPage.value, 
      limit: limit.value,
      search: searchQuery.value
    });
    if (res.data && res.data.meta) {
      adminPackages.value = Array.isArray(res.data.data) ? res.data.data : [];
      totalPages.value = res.data.meta.total_pages;
      totalRows.value = res.data.meta.total_rows;
      currentPage.value = res.data.meta.page;
    } else {
      adminPackages.value = Array.isArray(res.data) ? res.data : [];
      totalPages.value = 1;
      totalRows.value = adminPackages.value.length;
    }
  } catch (err) {
    console.error("Failed to load investment packages", err);
    packageError.value = "Gagal memuat data paket investasi.";
  } finally {
    packageLoading.value = false;
  }
}

function handleSearch(val: string) {
  searchQuery.value = val;
  currentPage.value = 1;
  loadAdminPackages();
}

function handlePageChange(val: number) {
  currentPage.value = val;
  loadAdminPackages();
}

function startEditPackage(item: InvestmentPackageItem) {
  editingPackageId.value = item.id;
  packageForm.package_name = item.package_name;
  packageForm.description = item.description || "";
  packageForm.min_quantity = item.min_quantity;
  packageForm.price = item.price;
  packageForm.status = item.status || "active";
  packageForm.roi = item.roi || "";
  packageForm.duration = item.duration || "";
  packageForm.benefits = Array.isArray(item.benefits) ? item.benefits.join("\n") : "";
  packageForm.is_popular = item.is_popular || false;
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
    roi: packageForm.roi.trim(),
    duration: packageForm.duration.trim(),
    benefits: packageForm.benefits.split("\n").map(b => b.trim()).filter(b => b !== ""),
    is_popular: packageForm.is_popular,
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
  loadAdminPackages();
});
</script>

<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-semibold text-gray-900">
        CRUD Paket Investasi
      </h2>
      <p class="text-gray-600 mt-1">
        Kelola paket investasi yang akan tampil di landing page.
      </p>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-6 space-y-5">
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-semibold text-gray-900">Form Paket Investasi</h3>
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

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Estimasi ROI</label>
          <input
            v-model="packageForm.roi"
            type="text"
            class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-green-500"
            placeholder="Contoh: 30-40%"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Durasi</label>
          <input
            v-model="packageForm.duration"
            type="text"
            class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-green-500"
            placeholder="Contoh: 24-36 bulan"
          />
        </div>

        <div class="md:col-span-2">
          <label class="block text-sm font-medium text-gray-700 mb-1">Benefit (Pisahkan dengan baris baru)</label>
          <textarea
            v-model="packageForm.benefits"
            rows="4"
            class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-green-500"
            placeholder="Cocok untuk pemula&#10;Laporan bulanan terperinci"
          />
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

        <div class="md:col-span-2 flex items-center mb-4">
          <input
            v-model="packageForm.is_popular"
            type="checkbox"
            id="is_popular"
            class="h-4 w-4 text-green-600 focus:ring-green-500 border-gray-300 rounded"
          />
          <label for="is_popular" class="ml-2 block text-sm text-gray-900">
            Tandai sebagai Paket Paling Populer
          </label>
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

      <div class="flex items-center justify-between mb-4 mt-8">
        <h4 class="text-md font-medium text-gray-900">Daftar Paket Investasi</h4>
        <SearchBar v-model="searchQuery" @search="handleSearch" placeholder="Cari nama atau deskripsi..." />
      </div>

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
      
      <Pagination 
        v-if="adminPackages.length > 0"
        :current-page="currentPage" 
        :total-pages="totalPages" 
        :total-rows="totalRows" 
        :limit="limit"
        @update:page="handlePageChange"
      />
    </div>
  </div>
</template>
