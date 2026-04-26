<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import {
  BadgeDollarSign,
  Edit3,
  Layers3,
  RefreshCcw,
  Save,
  Search,
  Trash2,
  X,
} from "lucide-vue-next";
import {
  createGrade,
  createNationalPrice,
  deleteGrade,
  deleteNationalPrice,
  getGrades,
  getNationalPrices,
  updateGrade,
  updateNationalPrice,
} from "../../services/price/vanili";

type GradeItem = {
  id: number;
  grade_name: string;
  description?: string | null;
};

type GradeRef = {
  id: number;
  grade_name: string;
};

type PriceItem = {
  national_price_id: number;
  province: string;
  price_min_per_kg?: number;
  price_max_per_kg?: number;
  price_per_kg: number;
  effective_date: string;
  harvest_type: string;
  grade_id?: number | null;
  grade?: GradeRef | null;
};

const grades = ref<GradeItem[]>([]);
const prices = ref<PriceItem[]>([]);
const loadingGrades = ref(false);
const loadingPrices = ref(false);
const savingGrade = ref(false);
const savingPrice = ref(false);

const gradeError = ref("");
const priceError = ref("");
const gradeSuccess = ref("");
const priceSuccess = ref("");

const gradeModalOpen = ref(false);
const priceModalOpen = ref(false);
const editingGradeId = ref<number | null>(null);
const editingPriceId = ref<number | null>(null);

type DeleteTarget = {
  type: "grade" | "price";
  id: number;
  label: string;
};

const deleteModalOpen = ref(false);
const deleteSubmitting = ref(false);
const deleteTarget = ref<DeleteTarget | null>(null);

const gradeSearchQuery = ref("");
const priceSearchQuery = ref("");
const activeTableView = ref<"grades" | "prices">("grades");
const gradePage = ref(1);
const pricePage = ref(1);
const gradePageSize = 8;
const pricePageSize = 8;

const gradeForm = reactive({
  grade_name: "",
  description: "",
});

const priceForm = reactive({
  province: "",
  price_min_per_kg: 0,
  price_max_per_kg: 0,
  effective_date: new Date().toISOString().slice(0, 10),
  harvest_type: "basah",
  grade_id: "",
});

const harvestTypeOptions = [
  { value: "basah", label: "Basah" },
  { value: "kering", label: "Kering" },
] as const;

const isWetHarvestType = computed(() => priceForm.harvest_type === "basah");

const filteredGrades = computed(() => {
  const query = gradeSearchQuery.value.trim().toLowerCase();
  if (!query) return grades.value;

  return grades.value.filter((item) => {
    return (
      item.grade_name.toLowerCase().includes(query) ||
      String(item.description || "").toLowerCase().includes(query)
    );
  });
});

const filteredPrices = computed(() => {
  const query = priceSearchQuery.value.trim().toLowerCase();
  if (!query) return prices.value;

  return prices.value.filter((item) => {
    const gradeName = item.grade?.grade_name?.toLowerCase() || "";
    return (
      item.province.toLowerCase().includes(query) ||
      item.harvest_type.toLowerCase().includes(query) ||
      gradeName.includes(query)
    );
  });
});

const gradeTotalPages = computed(() => Math.max(1, Math.ceil(filteredGrades.value.length / gradePageSize)));
const priceTotalPages = computed(() => Math.max(1, Math.ceil(filteredPrices.value.length / pricePageSize)));

const paginatedGrades = computed(() => {
  const start = (gradePage.value - 1) * gradePageSize;
  return filteredGrades.value.slice(start, start + gradePageSize);
});

const paginatedPrices = computed(() => {
  const start = (pricePage.value - 1) * pricePageSize;
  return filteredPrices.value.slice(start, start + pricePageSize);
});

watch(gradeSearchQuery, () => {
  gradePage.value = 1;
});

watch(priceSearchQuery, () => {
  pricePage.value = 1;
});

watch(filteredGrades, () => {
  if (gradePage.value > gradeTotalPages.value) {
    gradePage.value = gradeTotalPages.value;
  }
});

watch(filteredPrices, () => {
  if (pricePage.value > priceTotalPages.value) {
    pricePage.value = priceTotalPages.value;
  }
});

const formatRupiah = (value: number) =>
  new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(value);

function getPriceRange(item: PriceItem) {
  const legacyPrice = Number(item.price_per_kg) || 0;
  const min = Number(item.price_min_per_kg) || legacyPrice;
  const max = Number(item.price_max_per_kg) || legacyPrice;
  return { min, max };
}

function resetGradeForm() {
  gradeForm.grade_name = "";
  gradeForm.description = "";
  editingGradeId.value = null;
}

function resetPriceForm() {
  priceForm.province = "";
  priceForm.price_min_per_kg = 0;
  priceForm.price_max_per_kg = 0;
  priceForm.effective_date = new Date().toISOString().slice(0, 10);
  priceForm.harvest_type = "basah";
  priceForm.grade_id = "";
  editingPriceId.value = null;
}

function openCreateGradeModal() {
  gradeError.value = "";
  resetGradeForm();
  gradeModalOpen.value = true;
}

function openEditGradeModal(item: GradeItem) {
  gradeError.value = "";
  editingGradeId.value = item.id;
  gradeForm.grade_name = item.grade_name;
  gradeForm.description = item.description || "";
  gradeModalOpen.value = true;
}

function closeGradeModal() {
  gradeModalOpen.value = false;
  resetGradeForm();
}

function openCreatePriceModal() {
  priceError.value = "";
  resetPriceForm();
  priceModalOpen.value = true;
}

function openEditPriceModal(item: PriceItem) {
  priceError.value = "";
  editingPriceId.value = item.national_price_id;
  priceForm.province = item.province;

  const legacyPrice = Number(item.price_per_kg) || 0;
  priceForm.price_min_per_kg = Number(item.price_min_per_kg) || legacyPrice;
  priceForm.price_max_per_kg = Number(item.price_max_per_kg) || legacyPrice;
  priceForm.effective_date = item.effective_date
    ? String(item.effective_date).slice(0, 10)
    : new Date().toISOString().slice(0, 10);

  const normalizedHarvestType = String(item.harvest_type || "").trim().toLowerCase();
  priceForm.harvest_type = normalizedHarvestType === "kering" ? "kering" : "basah";
  priceForm.grade_id = priceForm.harvest_type === "basah" ? "" : item.grade_id ? String(item.grade_id) : "";

  priceModalOpen.value = true;
}

function closePriceModal() {
  priceModalOpen.value = false;
  resetPriceForm();
}

async function loadGrades() {
  loadingGrades.value = true;
  try {
    const response = await getGrades();
    grades.value = Array.isArray(response.data) ? response.data : [];
  } catch (error) {
    console.error("Failed to load grades", error);
    gradeError.value = "Gagal memuat data grade vanili.";
  } finally {
    loadingGrades.value = false;
  }
}

async function loadPrices() {
  loadingPrices.value = true;
  try {
    const response = await getNationalPrices();
    prices.value = Array.isArray(response.data) ? response.data : [];
  } catch (error) {
    console.error("Failed to load national prices", error);
    priceError.value = "Gagal memuat data harga vanili.";
  } finally {
    loadingPrices.value = false;
  }
}

async function refreshData() {
  await Promise.all([loadGrades(), loadPrices()]);
}

async function submitGradeForm() {
  gradeError.value = "";
  gradeSuccess.value = "";

  if (!gradeForm.grade_name.trim()) {
    gradeError.value = "Nama grade wajib diisi.";
    return;
  }

  savingGrade.value = true;
  try {
    const payload = {
      grade_name: gradeForm.grade_name.trim(),
      description: gradeForm.description.trim() || null,
    };

    if (editingGradeId.value) {
      await updateGrade(editingGradeId.value, payload);
      gradeSuccess.value = "Grade vanili berhasil diperbarui.";
    } else {
      await createGrade(payload);
      gradeSuccess.value = "Grade vanili berhasil dibuat.";
    }

    closeGradeModal();
    await refreshData();
  } catch (error) {
    console.error("Failed to save grade", error);
    gradeError.value = "Gagal menyimpan grade vanili.";
  } finally {
    savingGrade.value = false;
  }
}

async function submitPriceForm() {
  priceError.value = "";
  priceSuccess.value = "";

  if (!priceForm.province.trim()) {
    priceError.value = "Provinsi wajib diisi.";
    return;
  }
  if (!harvestTypeOptions.some((option) => option.value === priceForm.harvest_type)) {
    priceError.value = "Jenis panen wajib diisi.";
    return;
  }
  if (priceForm.price_min_per_kg <= 0 || priceForm.price_max_per_kg <= 0) {
    priceError.value = "Harga minimum dan maksimum per kg harus lebih dari 0.";
    return;
  }
  if (priceForm.price_min_per_kg > priceForm.price_max_per_kg) {
    priceError.value = "Harga minimum tidak boleh lebih besar dari harga maksimum.";
    return;
  }
  if (!priceForm.effective_date) {
    priceError.value = "Tanggal berlaku wajib diisi.";
    return;
  }

  savingPrice.value = true;
  try {
    if (isWetHarvestType.value) {
      priceForm.grade_id = "";
    }

    const payload = {
      province: priceForm.province.trim(),
      price_min_per_kg: Number(priceForm.price_min_per_kg),
      price_max_per_kg: Number(priceForm.price_max_per_kg),
      effective_date: priceForm.effective_date,
      harvest_type: priceForm.harvest_type,
      grade_id: isWetHarvestType.value ? null : priceForm.grade_id ? Number(priceForm.grade_id) : null,
    };

    if (editingPriceId.value) {
      await updateNationalPrice(editingPriceId.value, payload);
      priceSuccess.value = "Harga vanili berhasil diperbarui.";
    } else {
      await createNationalPrice(payload);
      priceSuccess.value = "Harga vanili berhasil dibuat.";
    }

    closePriceModal();
    await loadPrices();
  } catch (error) {
    console.error("Failed to save national price", error);
    priceError.value = "Gagal menyimpan harga vanili.";
  } finally {
    savingPrice.value = false;
  }
}

async function removeGrade(id: number) {
  gradeError.value = "";
  gradeSuccess.value = "";

  try {
    await deleteGrade(id);
    if (editingGradeId.value === id) closeGradeModal();
    gradeSuccess.value = "Grade vanili berhasil dihapus.";
    await refreshData();
  } catch (error) {
    console.error("Failed to delete grade", error);
    gradeError.value = "Gagal menghapus grade vanili.";
  }
}

async function removePrice(id: number) {
  priceError.value = "";
  priceSuccess.value = "";

  try {
    await deleteNationalPrice(id);
    if (editingPriceId.value === id) closePriceModal();
    priceSuccess.value = "Harga vanili berhasil dihapus.";
    await loadPrices();
  } catch (error) {
    console.error("Failed to delete national price", error);
    priceError.value = "Gagal menghapus harga vanili.";
  }
}

function askDeleteGrade(item: GradeItem) {
  deleteTarget.value = {
    type: "grade",
    id: item.id,
    label: item.grade_name,
  };
  deleteModalOpen.value = true;
}

function askDeletePrice(item: PriceItem) {
  deleteTarget.value = {
    type: "price",
    id: item.national_price_id,
    label: `${item.province} (${item.harvest_type})`,
  };
  deleteModalOpen.value = true;
}

function closeDeleteModal() {
  if (deleteSubmitting.value) return;
  deleteModalOpen.value = false;
  deleteTarget.value = null;
}

async function confirmDelete() {
  if (!deleteTarget.value) return;

  deleteSubmitting.value = true;
  try {
    if (deleteTarget.value.type === "grade") {
      await removeGrade(deleteTarget.value.id);
    } else {
      await removePrice(deleteTarget.value.id);
    }
    closeDeleteModal();
  } finally {
    deleteSubmitting.value = false;
  }
}

onMounted(refreshData);
</script>

<template>
  <div class="space-y-6">
    <div class="rounded-xl border border-gray-200 bg-white p-6">
      <div class="flex flex-wrap items-start justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Harga & Grade Vanili</h2>
          <p class="mt-1 text-gray-600">Kelola referensi grade dan rentang harga vanili nasional.</p>
        </div>

        <div class="flex flex-wrap items-center gap-2">
          <button
            type="button"
            class="inline-flex items-center gap-2 rounded-lg border border-gray-300 px-3 py-2 text-sm font-medium text-gray-700 transition hover:bg-gray-50"
            @click="refreshData"
          >
            <RefreshCcw class="h-4 w-4" />
            Refresh data
          </button>
          <button
            type="button"
            class="inline-flex items-center gap-2 rounded-lg bg-emerald-600 px-3 py-2 text-sm font-medium text-white transition hover:bg-emerald-700"
            @click="openCreateGradeModal"
          >
            <Layers3 class="h-4 w-4" />
            Kelola Grade
          </button>
          <button
            type="button"
            class="inline-flex items-center gap-2 rounded-lg bg-amber-500 px-3 py-2 text-sm font-medium text-white transition hover:bg-amber-600"
            @click="openCreatePriceModal"
          >
            <BadgeDollarSign class="h-4 w-4" />
            Kelola Harga
          </button>
        </div>
      </div>
    </div>

    <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <div class="rounded-xl border border-gray-200 bg-white p-4">
        <p class="text-xs font-semibold uppercase tracking-wide text-gray-500">Total Grade</p>
        <p class="mt-2 text-2xl font-semibold text-gray-900">{{ grades.length }}</p>
        <p class="mt-1 text-xs text-gray-500">Referensi mutu aktif</p>
      </div>
      <div class="rounded-xl border border-gray-200 bg-white p-4">
        <p class="text-xs font-semibold uppercase tracking-wide text-gray-500">Total Harga</p>
        <p class="mt-2 text-2xl font-semibold text-gray-900">{{ prices.length }}</p>
        <p class="mt-1 text-xs text-gray-500">Data harga nasional</p>
      </div>
      <div class="rounded-xl border border-gray-200 bg-white p-4">
        <p class="text-xs font-semibold uppercase tracking-wide text-gray-500">Panen Basah</p>
        <p class="mt-2 text-2xl font-semibold text-gray-900">{{ prices.filter((item) => item.harvest_type === 'basah').length }}</p>
        <p class="mt-1 text-xs text-gray-500">Tanpa grade</p>
      </div>
      <div class="rounded-xl border border-gray-200 bg-white p-4">
        <p class="text-xs font-semibold uppercase tracking-wide text-gray-500">Panen Kering</p>
        <p class="mt-2 text-2xl font-semibold text-gray-900">{{ prices.filter((item) => item.harvest_type === 'kering').length }}</p>
        <p class="mt-1 text-xs text-gray-500">Dengan grade</p>
      </div>
    </div>

    <div v-if="gradeSuccess" class="rounded-lg border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
      {{ gradeSuccess }}
    </div>
    <div v-if="priceSuccess" class="rounded-lg border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
      {{ priceSuccess }}
    </div>
    <div v-if="gradeError" class="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
      {{ gradeError }}
    </div>
    <div v-if="priceError" class="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
      {{ priceError }}
    </div>

    <section class="overflow-hidden rounded-xl border border-gray-200 bg-white">
      <div class="border-b border-gray-200 px-4 py-3">
        <div class="inline-flex rounded-lg border border-gray-300 bg-gray-50 p-1">
          <button
            type="button"
            class="rounded-md px-3 py-1.5 text-sm font-medium transition"
            :class="activeTableView === 'grades' ? 'bg-white text-emerald-700 shadow-sm' : 'text-gray-600 hover:text-gray-900'"
            @click="activeTableView = 'grades'"
          >
            Tabel Grade
          </button>
          <button
            type="button"
            class="rounded-md px-3 py-1.5 text-sm font-medium transition"
            :class="activeTableView === 'prices' ? 'bg-white text-amber-700 shadow-sm' : 'text-gray-600 hover:text-gray-900'"
            @click="activeTableView = 'prices'"
          >
            Tabel Harga
          </button>
        </div>
      </div>

      <div v-if="activeTableView === 'grades'">
        <div class="flex flex-wrap items-center justify-between gap-3 border-b border-gray-200 px-4 py-3">
          <h3 class="text-lg font-semibold text-gray-900">Data Grade Vanili</h3>
          <div class="relative w-full max-w-sm">
            <Search class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
            <input
              v-model="gradeSearchQuery"
              type="text"
              class="w-full rounded-lg border border-gray-300 bg-white py-2.5 pl-10 pr-3 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
              placeholder="Cari grade"
            />
          </div>
        </div>

        <div class="max-h-[440px] overflow-auto">
          <table class="w-full min-w-[640px]">
            <thead class="bg-gray-50 text-left">
              <tr>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Nama Grade</th>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Deskripsi</th>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-if="loadingGrades">
                <td colspan="3" class="px-4 py-6 text-center text-sm text-gray-500">Memuat grade...</td>
              </tr>
              <tr v-else-if="!filteredGrades.length">
                <td colspan="3" class="px-4 py-6 text-center text-sm text-gray-500">Tidak ada data grade.</td>
              </tr>
              <tr v-for="item in paginatedGrades" :key="item.id" class="hover:bg-gray-50">
                <td class="px-4 py-3 text-sm font-medium text-gray-900">{{ item.grade_name }}</td>
                <td class="px-4 py-3 text-sm text-gray-600">{{ item.description || '-' }}</td>
                <td class="px-4 py-3">
                  <div class="flex items-center justify-end gap-2">
                    <button
                      type="button"
                      class="rounded-lg border border-gray-200 p-2 text-gray-600 transition hover:bg-gray-100"
                      @click="openEditGradeModal(item)"
                    >
                      <Edit3 class="h-4 w-4" />
                    </button>
                    <button
                      type="button"
                      class="rounded-lg border border-red-200 p-2 text-red-600 transition hover:bg-red-50"
                      @click="askDeleteGrade(item)"
                    >
                      <Trash2 class="h-4 w-4" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="flex items-center justify-between border-t border-gray-200 px-4 py-3">
          <p class="text-xs text-gray-500">
            Menampilkan {{ paginatedGrades.length }} dari {{ filteredGrades.length }} data grade
          </p>
          <div class="flex items-center gap-2">
            <button
              type="button"
              class="rounded-lg border border-gray-300 px-3 py-1.5 text-sm text-gray-700 disabled:cursor-not-allowed disabled:opacity-50"
              :disabled="gradePage <= 1"
              @click="gradePage -= 1"
            >
              Sebelumnya
            </button>
            <span class="text-sm text-gray-600">Halaman {{ gradePage }} / {{ gradeTotalPages }}</span>
            <button
              type="button"
              class="rounded-lg border border-gray-300 px-3 py-1.5 text-sm text-gray-700 disabled:cursor-not-allowed disabled:opacity-50"
              :disabled="gradePage >= gradeTotalPages"
              @click="gradePage += 1"
            >
              Berikutnya
            </button>
          </div>
        </div>
      </div>

      <div v-else>
        <div class="flex flex-wrap items-center justify-between gap-3 border-b border-gray-200 px-4 py-3">
          <h3 class="text-lg font-semibold text-gray-900">Data Harga Vanili Nasional</h3>
          <div class="relative w-full max-w-sm">
            <Search class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
            <input
              v-model="priceSearchQuery"
              type="text"
              class="w-full rounded-lg border border-gray-300 bg-white py-2.5 pl-10 pr-3 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
              placeholder="Cari provinsi, grade, jenis panen"
            />
          </div>
        </div>

        <div class="max-h-[440px] overflow-auto">
          <table class="w-full min-w-[980px]">
            <thead class="bg-gray-50 text-left">
              <tr>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Provinsi</th>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Grade</th>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Jenis Panen</th>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Harga / Kg</th>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Tanggal Berlaku</th>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-if="loadingPrices">
                <td colspan="6" class="px-4 py-6 text-center text-sm text-gray-500">Memuat harga vanili...</td>
              </tr>
              <tr v-else-if="!filteredPrices.length">
                <td colspan="6" class="px-4 py-6 text-center text-sm text-gray-500">Tidak ada data harga vanili.</td>
              </tr>
              <tr v-for="item in paginatedPrices" :key="item.national_price_id" class="hover:bg-gray-50">
                <td class="px-4 py-3 text-sm text-gray-900">{{ item.province }}</td>
                <td class="px-4 py-3 text-sm text-gray-700">{{ item.grade?.grade_name || 'Tanpa grade' }}</td>
                <td class="px-4 py-3 text-sm text-gray-700">{{ item.harvest_type }}</td>
                <td class="px-4 py-3 text-sm font-medium text-gray-900">
                  {{ formatRupiah(getPriceRange(item).min) }} - {{ formatRupiah(getPriceRange(item).max) }}
                </td>
                <td class="px-4 py-3 text-sm text-gray-700">{{ item.effective_date.slice(0, 10) }}</td>
                <td class="px-4 py-3">
                  <div class="flex items-center justify-end gap-2">
                    <button
                      type="button"
                      class="rounded-lg border border-gray-200 p-2 text-gray-600 transition hover:bg-gray-100"
                      @click="openEditPriceModal(item)"
                    >
                      <Edit3 class="h-4 w-4" />
                    </button>
                    <button
                      type="button"
                      class="rounded-lg border border-red-200 p-2 text-red-600 transition hover:bg-red-50"
                      @click="askDeletePrice(item)"
                    >
                      <Trash2 class="h-4 w-4" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="flex items-center justify-between border-t border-gray-200 px-4 py-3">
          <p class="text-xs text-gray-500">
            Menampilkan {{ paginatedPrices.length }} dari {{ filteredPrices.length }} data harga
          </p>
          <div class="flex items-center gap-2">
            <button
              type="button"
              class="rounded-lg border border-gray-300 px-3 py-1.5 text-sm text-gray-700 disabled:cursor-not-allowed disabled:opacity-50"
              :disabled="pricePage <= 1"
              @click="pricePage -= 1"
            >
              Sebelumnya
            </button>
            <span class="text-sm text-gray-600">Halaman {{ pricePage }} / {{ priceTotalPages }}</span>
            <button
              type="button"
              class="rounded-lg border border-gray-300 px-3 py-1.5 text-sm text-gray-700 disabled:cursor-not-allowed disabled:opacity-50"
              :disabled="pricePage >= priceTotalPages"
              @click="pricePage += 1"
            >
              Berikutnya
            </button>
          </div>
        </div>
      </div>
    </section>

    <div v-if="deleteModalOpen" class="fixed inset-0 z-50 bg-gray-900/45 p-4">
      <div class="mx-auto mt-20 w-full max-w-md rounded-xl border border-gray-200 bg-white p-6 shadow-xl">
        <h3 class="text-lg font-semibold text-gray-900">Konfirmasi Hapus</h3>
        <p class="mt-2 text-sm text-gray-600">
          Yakin ingin menghapus
          <span class="font-semibold text-gray-900">{{ deleteTarget?.label }}</span>
          ?
        </p>
        <p class="mt-1 text-xs text-gray-500">Tindakan ini tidak dapat dibatalkan.</p>

        <div class="mt-5 flex justify-end gap-2">
          <button
            type="button"
            class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
            :disabled="deleteSubmitting"
            @click="closeDeleteModal"
          >
            Batal
          </button>
          <button
            type="button"
            class="rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700 disabled:opacity-60"
            :disabled="deleteSubmitting"
            @click="confirmDelete"
          >
            {{ deleteSubmitting ? "Menghapus..." : "Ya, Hapus" }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="gradeModalOpen" class="fixed inset-0 z-50 bg-gray-900/45 p-4">
      <div class="mx-auto mt-10 w-full max-w-lg rounded-xl border border-gray-200 bg-white p-6 shadow-xl">
        <div class="flex items-start justify-between gap-3">
          <div>
            <h3 class="text-xl font-semibold text-gray-900">{{ editingGradeId ? 'Edit Grade Vanili' : 'Tambah Grade Vanili' }}</h3>
            <p class="mt-1 text-sm text-gray-600">Isi data grade yang akan dipakai pada harga vanili.</p>
          </div>
          <button
            type="button"
            class="rounded-lg border border-gray-200 p-2 text-gray-600 transition hover:bg-gray-100"
            @click="closeGradeModal"
          >
            <X class="h-4 w-4" />
          </button>
        </div>

        <form class="mt-5 space-y-4" @submit.prevent="submitGradeForm">
          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Nama grade</span>
            <input
              v-model="gradeForm.grade_name"
              type="text"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
              placeholder="Contoh: Grade A"
            />
          </label>

          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Deskripsi</span>
            <textarea
              v-model="gradeForm.description"
              rows="4"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
              placeholder="Karakteristik mutu grade ini"
            />
          </label>

          <div v-if="gradeError" class="rounded-lg border border-red-200 bg-red-50 px-3 py-2.5 text-sm text-red-700">
            {{ gradeError }}
          </div>

          <div class="flex justify-end gap-2">
            <button
              type="button"
              class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
              @click="closeGradeModal"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="savingGrade"
              class="inline-flex items-center gap-2 rounded-lg bg-emerald-600 px-4 py-2 text-sm font-medium text-white hover:bg-emerald-700 disabled:opacity-70"
            >
              <Save class="h-4 w-4" />
              {{ savingGrade ? 'Menyimpan...' : editingGradeId ? 'Perbarui Grade' : 'Simpan Grade' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <div v-if="priceModalOpen" class="fixed inset-0 z-50 bg-gray-900/45 p-4">
      <div class="mx-auto mt-10 w-full max-w-2xl rounded-xl border border-gray-200 bg-white p-6 shadow-xl">
        <div class="flex items-start justify-between gap-3">
          <div>
            <h3 class="text-xl font-semibold text-gray-900">{{ editingPriceId ? 'Edit Harga Vanili' : 'Tambah Harga Vanili' }}</h3>
            <p class="mt-1 text-sm text-gray-600">Simpan rentang harga per kilogram berdasarkan provinsi dan jenis panen.</p>
          </div>
          <button
            type="button"
            class="rounded-lg border border-gray-200 p-2 text-gray-600 transition hover:bg-gray-100"
            @click="closePriceModal"
          >
            <X class="h-4 w-4" />
          </button>
        </div>

        <form class="mt-5 grid gap-4 md:grid-cols-2" @submit.prevent="submitPriceForm">
          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Provinsi</span>
            <input
              v-model="priceForm.province"
              type="text"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
              placeholder="Contoh: Jawa Tengah"
            />
          </label>

          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Grade</span>
            <select
              v-model="priceForm.grade_id"
              :disabled="isWetHarvestType"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
              :class="isWetHarvestType ? 'cursor-not-allowed bg-gray-100 text-gray-400' : ''"
            >
              <option value="">Tanpa grade</option>
              <option v-for="grade in grades" :key="grade.id" :value="String(grade.id)">
                {{ grade.grade_name }}
              </option>
            </select>
          </label>

          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Jenis panen</span>
            <select
              v-model="priceForm.harvest_type"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
            >
              <option v-for="option in harvestTypeOptions" :key="option.value" :value="option.value">
                {{ option.label }}
              </option>
            </select>
          </label>

          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Tanggal berlaku</span>
            <input
              v-model="priceForm.effective_date"
              type="date"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
            />
          </label>

          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Harga minimum per kg</span>
            <input
              v-model.number="priceForm.price_min_per_kg"
              type="number"
              min="0"
              step="1000"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
              placeholder="Contoh: 850000"
            />
          </label>

          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Harga maksimum per kg</span>
            <input
              v-model.number="priceForm.price_max_per_kg"
              type="number"
              min="0"
              step="1000"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
              placeholder="Contoh: 920000"
            />
          </label>

          <div v-if="isWetHarvestType" class="md:col-span-2 text-xs text-gray-500">Panen basah tidak menggunakan grade.</div>

          <div v-if="priceError" class="md:col-span-2 rounded-lg border border-red-200 bg-red-50 px-3 py-2.5 text-sm text-red-700">
            {{ priceError }}
          </div>

          <div class="md:col-span-2 flex justify-end gap-2">
            <button
              type="button"
              class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
              @click="closePriceModal"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="savingPrice"
              class="inline-flex items-center gap-2 rounded-lg bg-amber-500 px-4 py-2 text-sm font-medium text-white hover:bg-amber-600 disabled:opacity-70"
            >
              <Save class="h-4 w-4" />
              {{ savingPrice ? 'Menyimpan...' : editingPriceId ? 'Perbarui Harga' : 'Simpan Harga' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
