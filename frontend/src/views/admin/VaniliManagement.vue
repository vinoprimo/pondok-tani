<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import {
  BadgeDollarSign,
  CalendarDays,
  CheckCircle2,
  Edit3,
  Layers3,
  Leaf,
  MapPinned,
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
  price_per_kg: number;
  effective_date: string;
  harvest_type: string;
  grade_id?: number | null;
  grade?: GradeRef | null;
};

const activeTab = ref<"grades" | "prices">("grades");
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
const editingGradeId = ref<number | null>(null);
const editingPriceId = ref<number | null>(null);
const searchQuery = ref("");

const gradeForm = reactive({
  grade_name: "",
  description: "",
});

const priceForm = reactive({
  province: "",
  price_per_kg: 0,
  effective_date: new Date().toISOString().slice(0, 10),
  harvest_type: "basah",
  grade_id: "",
});

const harvestTypeOptions = [
  { value: "basah", label: "Basah" },
  { value: "kering", label: "Kering" },
] as const;

const isWetHarvestType = computed(() => priceForm.harvest_type === "basah");

const formatRupiah = (value: number) =>
  new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(value);

const filteredPrices = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
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

function resetGradeForm() {
  gradeForm.grade_name = "";
  gradeForm.description = "";
  editingGradeId.value = null;
}

function resetPriceForm() {
  priceForm.province = "";
  priceForm.price_per_kg = 0;
  priceForm.effective_date = new Date().toISOString().slice(0, 10);
  priceForm.harvest_type = "basah";
  priceForm.grade_id = "";
  editingPriceId.value = null;
}

async function loadGrades() {
  loadingGrades.value = true;
  gradeError.value = "";
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
  priceError.value = "";
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

function startEditGrade(item: GradeItem) {
  editingGradeId.value = item.id;
  gradeForm.grade_name = item.grade_name;
  gradeForm.description = item.description || "";
  activeTab.value = "grades";
}

function startEditPrice(item: PriceItem) {
  editingPriceId.value = item.national_price_id;
  priceForm.province = item.province;
  priceForm.price_per_kg = Number(item.price_per_kg) || 0;
  priceForm.effective_date = item.effective_date ? String(item.effective_date).slice(0, 10) : new Date().toISOString().slice(0, 10);
  const normalizedHarvestType = String(item.harvest_type || "").trim().toLowerCase();
  priceForm.harvest_type = normalizedHarvestType === "kering" ? "kering" : "basah";
  priceForm.grade_id = priceForm.harvest_type === "basah" ? "" : item.grade_id ? String(item.grade_id) : "";
  activeTab.value = "prices";
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

    resetGradeForm();
    await loadGrades();
    await loadPrices();
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
  if (priceForm.price_per_kg <= 0) {
    priceError.value = "Harga per kg harus lebih dari 0.";
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
      price_per_kg: Number(priceForm.price_per_kg),
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

    resetPriceForm();
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
  if (!window.confirm("Hapus grade vanili ini?")) return;

  try {
    await deleteGrade(id);
    if (editingGradeId.value === id) resetGradeForm();
    gradeSuccess.value = "Grade vanili berhasil dihapus.";
    await loadGrades();
    await loadPrices();
  } catch (error) {
    console.error("Failed to delete grade", error);
    gradeError.value = "Gagal menghapus grade vanili.";
  }
}

async function removePrice(id: number) {
  priceError.value = "";
  priceSuccess.value = "";
  if (!window.confirm("Hapus data harga vanili ini?")) return;

  try {
    await deleteNationalPrice(id);
    if (editingPriceId.value === id) resetPriceForm();
    priceSuccess.value = "Harga vanili berhasil dihapus.";
    await loadPrices();
  } catch (error) {
    console.error("Failed to delete national price", error);
    priceError.value = "Gagal menghapus harga vanili.";
  }
}

function cancelGradeEdit() {
  resetGradeForm();
}

function cancelPriceEdit() {
  resetPriceForm();
}

onMounted(async () => {
  await Promise.all([loadGrades(), loadPrices()]);
});
</script>

<template>
  <div class="space-y-6">
    <section class="overflow-hidden rounded-3xl border border-emerald-100 bg-gradient-to-br from-emerald-50 via-white to-amber-50 shadow-sm">
      <div class="grid gap-6 px-6 py-6 lg:grid-cols-[1.3fr_0.7fr] lg:px-8 lg:py-8">
        <div class="space-y-4">
          <div class="inline-flex items-center gap-2 rounded-full bg-white px-3 py-1 text-xs font-semibold uppercase tracking-[0.2em] text-emerald-700 shadow-sm ring-1 ring-emerald-100">
            <Leaf class="h-4 w-4" />
            Panel Admin Vanili
          </div>
          <div class="space-y-3">
            <h2 class="text-3xl font-semibold tracking-tight text-gray-900 lg:text-4xl">
              Manajemen harga vanili dan grade mutu
            </h2>
            <p class="max-w-2xl text-sm leading-6 text-gray-600 lg:text-base">
              Kelola referensi mutu vanili, lalu simpan data harga nasional berdasarkan provinsi,
              jenis panen, dan grade yang berlaku.
            </p>
          </div>
          <div class="flex flex-wrap gap-3">
            <button
              type="button"
              class="inline-flex items-center gap-2 rounded-xl bg-emerald-600 px-4 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-emerald-700"
              @click="activeTab = 'prices'"
            >
              <BadgeDollarSign class="h-4 w-4" />
              Kelola Harga
            </button>
            <button
              type="button"
              class="inline-flex items-center gap-2 rounded-xl border border-emerald-200 bg-white px-4 py-2.5 text-sm font-semibold text-emerald-700 transition-colors hover:bg-emerald-50"
              @click="activeTab = 'grades'"
            >
              <Layers3 class="h-4 w-4" />
              Kelola Grade
            </button>
          </div>
        </div>

        <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-1 xl:grid-cols-2">
          <div class="rounded-2xl border border-white/70 bg-white/90 p-4 shadow-sm backdrop-blur">
            <p class="text-xs font-semibold uppercase tracking-[0.2em] text-gray-500">Total Grade</p>
            <div class="mt-3 flex items-end justify-between gap-3">
              <div>
                <p class="text-3xl font-semibold text-gray-900">{{ grades.length }}</p>
                <p class="mt-1 text-sm text-gray-500">referensi mutu aktif</p>
              </div>
              <div class="rounded-2xl bg-emerald-50 p-3 text-emerald-700">
                <Layers3 class="h-6 w-6" />
              </div>
            </div>
          </div>
          <div class="rounded-2xl border border-white/70 bg-white/90 p-4 shadow-sm backdrop-blur">
            <p class="text-xs font-semibold uppercase tracking-[0.2em] text-gray-500">Total Harga</p>
            <div class="mt-3 flex items-end justify-between gap-3">
              <div>
                <p class="text-3xl font-semibold text-gray-900">{{ prices.length }}</p>
                <p class="mt-1 text-sm text-gray-500">catatan harga nasional</p>
              </div>
              <div class="rounded-2xl bg-amber-50 p-3 text-amber-700">
                <BadgeDollarSign class="h-6 w-6" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <div class="grid gap-6 xl:grid-cols-[0.92fr_1.08fr]">
      <section class="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm">
        <div class="flex items-start justify-between gap-4">
          <div>
            <p class="text-xs font-semibold uppercase tracking-[0.2em] text-emerald-600">CRUD Grade</p>
            <h3 class="mt-2 text-xl font-semibold text-gray-900">Grade vanili</h3>
            <p class="mt-2 text-sm leading-6 text-gray-600">Tambahkan, ubah, atau hapus referensi grade yang dipakai di data harga.</p>
          </div>
          <button
            type="button"
            class="inline-flex items-center gap-2 rounded-xl border border-gray-200 px-3 py-2 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-50"
            @click="cancelGradeEdit"
          >
            <RefreshCcw class="h-4 w-4" />
            Reset
          </button>
        </div>

        <form class="mt-6 space-y-4" @submit.prevent="submitGradeForm">
          <div class="grid gap-4">
            <label class="space-y-2">
              <span class="text-sm font-medium text-gray-700">Nama grade</span>
              <input
                v-model="gradeForm.grade_name"
                type="text"
                class="w-full rounded-xl border border-gray-300 bg-white px-4 py-3 text-sm outline-none transition focus:border-emerald-500 focus:ring-4 focus:ring-emerald-100"
                placeholder="Contoh: Grade A"
              />
            </label>
            <label class="space-y-2">
              <span class="text-sm font-medium text-gray-700">Deskripsi</span>
              <textarea
                v-model="gradeForm.description"
                rows="4"
                class="w-full rounded-xl border border-gray-300 bg-white px-4 py-3 text-sm outline-none transition focus:border-emerald-500 focus:ring-4 focus:ring-emerald-100"
                placeholder="Karakteristik mutu grade ini"
              />
            </label>
          </div>

          <div v-if="gradeError" class="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {{ gradeError }}
          </div>
          <div v-if="gradeSuccess" class="rounded-xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
            {{ gradeSuccess }}
          </div>

          <button
            type="submit"
            :disabled="savingGrade"
            class="inline-flex w-full items-center justify-center gap-2 rounded-xl bg-emerald-600 px-4 py-3 text-sm font-semibold text-white transition-colors hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-70"
          >
            <Save class="h-4 w-4" />
            {{ savingGrade ? 'Menyimpan...' : editingGradeId ? 'Perbarui Grade' : 'Simpan Grade' }}
          </button>
        </form>

        <div class="mt-6 space-y-3">
          <div class="flex items-center justify-between">
            <h4 class="text-sm font-semibold uppercase tracking-[0.2em] text-gray-500">Daftar Grade</h4>
            <span class="text-xs text-gray-400">{{ grades.length }} item</span>
          </div>

          <div v-if="loadingGrades" class="rounded-2xl border border-dashed border-gray-200 p-6 text-center text-sm text-gray-500">
            Memuat grade...
          </div>

          <div v-else-if="!grades.length" class="rounded-2xl border border-dashed border-gray-200 p-6 text-center text-sm text-gray-500">
            Belum ada grade vanili.
          </div>

          <div v-else class="space-y-3">
            <article
              v-for="item in grades"
              :key="item.id"
              class="rounded-2xl border border-gray-200 p-4 transition hover:border-emerald-200 hover:bg-emerald-50/40"
            >
              <div class="flex items-start justify-between gap-4">
                <div>
                  <div class="flex items-center gap-2">
                    <span class="rounded-full bg-emerald-100 px-3 py-1 text-xs font-semibold text-emerald-700">{{ item.grade_name }}</span>
                    <CheckCircle2 class="h-4 w-4 text-emerald-500" />
                  </div>
                  <p class="mt-2 text-sm leading-6 text-gray-600">{{ item.description || 'Tidak ada deskripsi.' }}</p>
                </div>
                <div class="flex shrink-0 items-center gap-2">
                  <button
                    type="button"
                    class="rounded-xl border border-gray-200 p-2 text-gray-600 transition hover:bg-gray-100"
                    @click="startEditGrade(item)"
                  >
                    <Edit3 class="h-4 w-4" />
                  </button>
                  <button
                    type="button"
                    class="rounded-xl border border-red-200 p-2 text-red-600 transition hover:bg-red-50"
                    @click="removeGrade(item.id)"
                  >
                    <Trash2 class="h-4 w-4" />
                  </button>
                </div>
              </div>
            </article>
          </div>
        </div>
      </section>

      <section class="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm">
        <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
          <div>
            <p class="text-xs font-semibold uppercase tracking-[0.2em] text-amber-600">CRUD Harga</p>
            <h3 class="mt-2 text-xl font-semibold text-gray-900">Harga vanili nasional</h3>
            <p class="mt-2 text-sm leading-6 text-gray-600">Catat harga per kg berdasarkan provinsi, grade, jenis panen, dan tanggal berlaku.</p>
          </div>
          <div class="w-full max-w-sm">
            <label class="space-y-2">
              <span class="text-sm font-medium text-gray-700">Cari harga</span>
              <div class="relative">
                <Search class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
                <input
                  v-model="searchQuery"
                  type="text"
                  class="w-full rounded-xl border border-gray-300 bg-white py-3 pl-10 pr-4 text-sm outline-none transition focus:border-emerald-500 focus:ring-4 focus:ring-emerald-100"
                  placeholder="Provinsi, grade, atau jenis panen"
                />
              </div>
            </label>
          </div>
        </div>

        <form class="mt-6 grid gap-4 md:grid-cols-2" @submit.prevent="submitPriceForm">
          <label class="space-y-2 md:col-span-1">
            <span class="text-sm font-medium text-gray-700">Provinsi</span>
            <input
              v-model="priceForm.province"
              type="text"
              class="w-full rounded-xl border border-gray-300 bg-white px-4 py-3 text-sm outline-none transition focus:border-emerald-500 focus:ring-4 focus:ring-emerald-100"
              placeholder="Contoh: Jawa Tengah"
            />
          </label>
          <label class="space-y-2 md:col-span-1">
            <span class="text-sm font-medium text-gray-700">Grade</span>
            <select
              v-model="priceForm.grade_id"
              :disabled="isWetHarvestType"
              class="w-full rounded-xl border border-gray-300 bg-white px-4 py-3 text-sm outline-none transition focus:border-emerald-500 focus:ring-4 focus:ring-emerald-100"
              :class="isWetHarvestType ? 'cursor-not-allowed bg-gray-100 text-gray-400' : ''"
            >
              <option value="">Tanpa grade</option>
              <option v-for="grade in grades" :key="grade.id" :value="String(grade.id)">
                {{ grade.grade_name }}
              </option>
            </select>
            <p v-if="isWetHarvestType" class="text-xs text-gray-500">Panen basah tidak menggunakan grade.</p>
          </label>
          <label class="space-y-2 md:col-span-1">
            <span class="text-sm font-medium text-gray-700">Jenis panen</span>
            <select
              v-model="priceForm.harvest_type"
              class="w-full rounded-xl border border-gray-300 bg-white px-4 py-3 text-sm outline-none transition focus:border-emerald-500 focus:ring-4 focus:ring-emerald-100"
            >
              <option v-for="option in harvestTypeOptions" :key="option.value" :value="option.value">
                {{ option.label }}
              </option>
            </select>
          </label>
          <label class="space-y-2 md:col-span-1">
            <span class="text-sm font-medium text-gray-700">Tanggal berlaku</span>
            <input
              v-model="priceForm.effective_date"
              type="date"
              class="w-full rounded-xl border border-gray-300 bg-white px-4 py-3 text-sm outline-none transition focus:border-emerald-500 focus:ring-4 focus:ring-emerald-100"
            />
          </label>
          <label class="space-y-2 md:col-span-1">
            <span class="text-sm font-medium text-gray-700">Harga per kg</span>
            <input
              v-model.number="priceForm.price_per_kg"
              type="number"
              min="0"
              step="1000"
              class="w-full rounded-xl border border-gray-300 bg-white px-4 py-3 text-sm outline-none transition focus:border-emerald-500 focus:ring-4 focus:ring-emerald-100"
              placeholder="Contoh: 850000"
            />
          </label>
          <div class="flex items-end gap-3 md:col-span-1">
            <button
              type="submit"
              :disabled="savingPrice"
              class="inline-flex flex-1 items-center justify-center gap-2 rounded-xl bg-amber-500 px-4 py-3 text-sm font-semibold text-white transition-colors hover:bg-amber-600 disabled:cursor-not-allowed disabled:opacity-70"
            >
              <Save class="h-4 w-4" />
              {{ savingPrice ? 'Menyimpan...' : editingPriceId ? 'Perbarui Harga' : 'Simpan Harga' }}
            </button>
            <button
              type="button"
              class="rounded-xl border border-gray-200 px-4 py-3 text-sm font-medium text-gray-600 transition hover:bg-gray-50"
              @click="cancelPriceEdit"
            >
              <X class="h-4 w-4" />
            </button>
          </div>
        </form>

        <div v-if="priceError" class="mt-4 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
          {{ priceError }}
        </div>
        <div v-if="priceSuccess" class="mt-4 rounded-xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
          {{ priceSuccess }}
        </div>

        <div class="mt-6 overflow-hidden rounded-2xl border border-gray-200">
          <div class="flex items-center justify-between border-b border-gray-200 bg-gray-50 px-4 py-3">
            <h4 class="text-sm font-semibold text-gray-700">Daftar Harga</h4>
            <span class="text-xs text-gray-500">{{ filteredPrices.length }} item</span>
          </div>

          <div v-if="loadingPrices" class="p-6 text-center text-sm text-gray-500">
            Memuat harga vanili...
          </div>

          <div v-else-if="!filteredPrices.length" class="p-6 text-center text-sm text-gray-500">
            Belum ada data harga vanili.
          </div>

          <div v-else class="divide-y divide-gray-200">
            <article
              v-for="item in filteredPrices"
              :key="item.national_price_id"
              class="grid gap-4 px-4 py-4 lg:grid-cols-[1fr_auto] lg:items-center"
            >
              <div class="space-y-2">
                <div class="flex flex-wrap items-center gap-2">
                  <span class="inline-flex items-center gap-1 rounded-full bg-amber-100 px-3 py-1 text-xs font-semibold text-amber-700">
                    <MapPinned class="h-3.5 w-3.5" />
                    {{ item.province }}
                  </span>
                  <span class="inline-flex items-center gap-1 rounded-full bg-emerald-100 px-3 py-1 text-xs font-semibold text-emerald-700">
                    <Layers3 class="h-3.5 w-3.5" />
                    {{ item.grade?.grade_name || 'Tanpa grade' }}
                  </span>
                  <span class="inline-flex items-center gap-1 rounded-full bg-gray-100 px-3 py-1 text-xs font-semibold text-gray-600">
                    <CalendarDays class="h-3.5 w-3.5" />
                    {{ item.effective_date.slice(0, 10) }}
                  </span>
                </div>
                <div class="space-y-1">
                  <p class="text-base font-semibold text-gray-900">
                    {{ formatRupiah(Number(item.price_per_kg)) }} / kg
                  </p>
                  <p class="text-sm text-gray-600">Jenis panen: {{ item.harvest_type }}</p>
                </div>
              </div>
              <div class="flex items-center gap-2 lg:justify-end">
                <button
                  type="button"
                  class="rounded-xl border border-gray-200 p-2 text-gray-600 transition hover:bg-gray-100"
                  @click="startEditPrice(item)"
                >
                  <Edit3 class="h-4 w-4" />
                </button>
                <button
                  type="button"
                  class="rounded-xl border border-red-200 p-2 text-red-600 transition hover:bg-red-50"
                  @click="removePrice(item.national_price_id)"
                >
                  <Trash2 class="h-4 w-4" />
                </button>
              </div>
            </article>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>
