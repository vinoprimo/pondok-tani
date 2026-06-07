<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  AlertCircle,
  Calendar,
  CalendarDays,
  CheckCircle,
  ChevronLeft,
  ChevronRight,
  Clock,
  Plus,
} from "lucide-vue-next";
import { getAdminUsers, getAdminUser, getAdminUserPlantBatches } from "../../services/user/user";
import {
  createAdminMaintenanceSchedule,
  getAdminMaintenanceSchedules,
} from "../../services/maintenance/schedule";

type UserItem = {
  id: string;
  name: string;
  email: string;
};

type PlantBatchItem = {
  id: number;
  batch_code: string;
  location: string;
};

type ScheduleItem = {
  id: number;
  user_id: string;
  user_name: string;
  user_email: string;
  plant_batch_id: number;
  batch_code: string;
  activity_type: string;
  frequency_days: number;
  next_due_date: string;
  status: string;
};

const route = useRoute();
const router = useRouter();
const userId = computed(() => String(route.params.userId || ""));

const user = ref<UserItem | null>(null);
const schedules = ref<ScheduleItem[]>([]);
const plantBatches = ref<PlantBatchItem[]>([]);

const isLoading = ref(false);
const isSaving = ref(false);
const isAddModalOpen = ref(false);
const currentMonth = ref(startOfMonth(new Date()));
const selectedDateKey = ref(formatDateKey(new Date()));

const form = ref({
  plant_batch_id: "",
  activity_type: "",
  frequency_days: 7,
  next_due_date: "",
});

const sortedSchedules = computed(() => {
  return [...schedules.value].sort(
    (a, b) => new Date(a.next_due_date).getTime() - new Date(b.next_due_date).getTime()
  );
});

const scheduleMap = computed(() => {
  const map = new Map<string, ScheduleItem[]>();

  for (const item of schedules.value) {
    const key = formatDateKey(new Date(item.next_due_date));
    const existing = map.get(key) || [];
    existing.push(item);
    map.set(key, existing);
  }

  return map;
});

const calendarCells = computed(() => {
  const firstDayOfMonth = startOfMonth(currentMonth.value);
  const startDate = new Date(firstDayOfMonth);
  startDate.setDate(startDate.getDate() - startDate.getDay());

  const cells: Array<{
    key: string;
    date: Date;
    inMonth: boolean;
    isToday: boolean;
    isSelected: boolean;
    items: ScheduleItem[];
  }> = [];

  for (let index = 0; index < 42; index += 1) {
    const date = new Date(startDate);
    date.setDate(startDate.getDate() + index);
    const key = formatDateKey(date);
    const items = scheduleMap.value.get(key) || [];

    cells.push({
      key,
      date,
      inMonth: date.getMonth() === currentMonth.value.getMonth(),
      isToday: key === formatDateKey(new Date()),
      isSelected: key === selectedDateKey.value,
      items,
    });
  }

  return cells;
});

const selectedDateSchedules = computed(() => scheduleMap.value.get(selectedDateKey.value) || []);

const selectedDateLabel = computed(() => {
  const date = parseDateKey(selectedDateKey.value);
  return date.toLocaleDateString("id-ID", {
    day: "2-digit",
    month: "long",
    year: "numeric",
  });
});

const monthLabel = computed(() =>
  currentMonth.value.toLocaleDateString("id-ID", {
    month: "long",
    year: "numeric",
  })
);

function startOfMonth(date: Date) {
  return new Date(date.getFullYear(), date.getMonth(), 1);
}

function formatDateKey(date: Date) {
  const year = date.getFullYear();
  const month = `${date.getMonth() + 1}`.padStart(2, "0");
  const day = `${date.getDate()}`.padStart(2, "0");
  return `${year}-${month}-${day}`;
}

function parseDateKey(dateKey: string) {
  const [year, month, day] = dateKey.split("-").map(Number);
  return new Date(year, month - 1, day);
}

function formatDate(dateStr?: string | null) {
  if (!dateStr) return "-";
  const date = new Date(dateStr);
  if (Number.isNaN(date.getTime())) return "-";
  return date.toLocaleDateString("id-ID", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });
}

function statusLabel(status: string) {
  if (status === "overdue") return "Overdue";
  if (status === "mendatang") return "Mendatang";
  if (status === "menunggu_verifikasi") return "Menunggu verifikasi";
  if (status === "terverifikasi") return "Terverifikasi";
  if (status === "ditolak") return "Ditolak";
  return "Aman";
}

function statusClass(status: string) {
  if (status === "overdue") return "bg-red-50 text-red-700";
  if (status === "mendatang") return "bg-blue-50 text-blue-700";
  if (status === "menunggu_verifikasi") return "bg-yellow-50 text-yellow-700";
  if (status === "terverifikasi") return "bg-green-50 text-green-700";
  if (status === "ditolak") return "bg-red-50 text-red-700";
  return "bg-green-50 text-green-700";
}

function activityColorClass(status: string) {
  if (status === "overdue") return "bg-red-500";
  if (status === "mendatang") return "bg-blue-500";
  if (status === "menunggu_verifikasi") return "bg-yellow-500";
  if (status === "terverifikasi") return "bg-green-500";
  if (status === "ditolak") return "bg-red-500";
  return "bg-green-500";
}

function dayHasDanger(items: ScheduleItem[]) {
  return items.some((item) => item.status === "overdue");
}

function dayHasWarning(items: ScheduleItem[]) {
  return !dayHasDanger(items) && items.some((item) => item.status === "mendatang" || item.status === "menunggu_verifikasi");
}

function selectDay(key: string) {
  selectedDateKey.value = key;
  form.value.next_due_date = key;
}

function prevMonth() {
  currentMonth.value = new Date(currentMonth.value.getFullYear(), currentMonth.value.getMonth() - 1, 1);
  selectedDateKey.value = formatDateKey(currentMonth.value);
}

function nextMonth() {
  currentMonth.value = new Date(currentMonth.value.getFullYear(), currentMonth.value.getMonth() + 1, 1);
  selectedDateKey.value = formatDateKey(currentMonth.value);
}

function goToToday() {
  const today = new Date();
  currentMonth.value = startOfMonth(today);
  selectedDateKey.value = formatDateKey(today);
}

async function loadUserData() {
  try {
    const userRes = await getAdminUser(userId.value);
    const found = userRes.data;

    if (!found) {
      throw new Error("User not found");
    }

    user.value = {
      id: found.id,
      name: found.name,
      email: found.email,
    };
  } catch (error: any) {
    alert("User tidak ditemukan");
    router.push("/dashboard/monitoring-schedule");
  }
}

async function loadSchedules() {
  isLoading.value = true;
  try {
    const res = await getAdminMaintenanceSchedules({ user_id: userId.value });
    schedules.value = Array.isArray(res.data) ? res.data : [];

    if (schedules.value.length > 0) {
      const firstScheduleDate = new Date(schedules.value[0].next_due_date);
      currentMonth.value = startOfMonth(firstScheduleDate);
      selectedDateKey.value = formatDateKey(firstScheduleDate);
    } else {
      goToToday();
    }
  } finally {
    isLoading.value = false;
  }
}

async function openAddModal() {
  form.value.plant_batch_id = "";
  form.value.activity_type = "";
  form.value.frequency_days = 7;
  form.value.next_due_date = selectedDateKey.value;

  const res = await getAdminUserPlantBatches(userId.value);
  plantBatches.value = Array.isArray(res.data)
    ? res.data.map((item: any) => ({
        id: item.id,
        batch_code: item.batch_code,
        location: item.location,
      }))
    : [];

  isAddModalOpen.value = true;
}

function closeAddModal() {
  isAddModalOpen.value = false;
}

async function submitSchedule() {
  if (!form.value.plant_batch_id || !form.value.activity_type || !form.value.next_due_date) {
    alert("Lengkapi semua field wajib sebelum menyimpan.");
    return;
  }

  isSaving.value = true;
  try {
    await createAdminMaintenanceSchedule({
      user_id: userId.value,
      plant_batch_id: Number(form.value.plant_batch_id),
      activity_type: form.value.activity_type,
      frequency_days: Number(form.value.frequency_days || 0),
      next_due_date: form.value.next_due_date,
    });

    await loadSchedules();
    closeAddModal();
  } catch (error: any) {
    const message = error?.response?.data?.error || "Gagal menambahkan jadwal monitoring";
    alert(message);
  } finally {
    isSaving.value = false;
  }
}

onMounted(async () => {
  if (!userId.value) {
    router.push("/dashboard/monitoring-schedule");
    return;
  }

  await loadUserData();
  await loadSchedules();
});
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-semibold text-gray-900">Detail Jadwal Monitoring</h2>
        <p class="text-gray-600 mt-1">
          {{ user ? `${user.name} (${user.email})` : "Memuat user..." }}
        </p>
      </div>
      <div class="flex items-center gap-2">
        <button
          type="button"
          class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
          @click="router.push('/dashboard/monitoring-schedule')"
        >
          Kembali
        </button>
        <button
          type="button"
          class="inline-flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700"
          @click="openAddModal"
        >
          <Plus class="h-4 w-4" />
          Tambah Jadwal
        </button>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-4 md:grid-cols-3 xl:grid-cols-4">
      <div class="rounded-xl border border-gray-200 bg-white p-5">
        <div class="flex items-start justify-between">
          <p class="text-sm text-gray-600">Total Jadwal</p>
          <CalendarDays class="h-5 w-5 text-green-600" />
        </div>
        <p class="mt-3 text-3xl font-semibold text-gray-900">{{ schedules.length }}</p>
        <p class="mt-2 text-sm text-gray-500">Semua task untuk user ini</p>
      </div>
      <div class="rounded-xl border border-gray-200 bg-white p-5">
        <div class="flex items-start justify-between">
          <p class="text-sm text-gray-600">Menunggu</p>
          <Clock class="h-5 w-5 text-yellow-600" />
        </div>
        <p class="mt-3 text-3xl font-semibold text-gray-900">
          {{ schedules.filter((item) => item.status === 'pending').length }}
        </p>
        <p class="mt-2 text-sm text-gray-500">Belum selesai dikerjakan</p>
      </div>
      <div class="rounded-xl border border-gray-200 bg-white p-5">
        <div class="flex items-start justify-between">
          <p class="text-sm text-gray-600">Overdue</p>
          <AlertCircle class="h-5 w-5 text-red-600" />
        </div>
        <p class="mt-3 text-3xl font-semibold text-gray-900">
          {{ schedules.filter((item) => item.status === 'overdue').length }}
        </p>
        <p class="mt-2 text-sm text-red-600">Perlu tindak lanjut</p>
      </div>
      <div class="rounded-xl border border-gray-200 bg-white p-5">
        <div class="flex items-start justify-between">
          <p class="text-sm text-gray-600">Jadwal Bulan Ini</p>
          <CheckCircle class="h-5 w-5 text-blue-600" />
        </div>
        <p class="mt-3 text-3xl font-semibold text-gray-900">
          {{ schedules.filter((item) => new Date(item.next_due_date).getMonth() === currentMonth.getMonth()).length }}
        </p>
        <p class="mt-2 text-sm text-gray-500">Di bulan {{ monthLabel }}</p>
      </div>
    </div>

    <div class="rounded-2xl border border-gray-200 bg-white p-5 md:p-6">
      <div class="flex flex-wrap items-center justify-between gap-3 border-b border-gray-200 pb-4">
        <div>
          <h3 class="text-lg font-semibold text-gray-900">Kalender Monitoring</h3>
          <p class="text-sm text-gray-600 mt-1">Tanda indikator muncul pada tanggal yang memiliki aktivitas.</p>
        </div>
        <div class="flex items-center gap-2">
          <button
            type="button"
            class="inline-flex items-center justify-center rounded-lg border border-gray-300 bg-white px-3 py-2 text-gray-700 hover:bg-gray-50"
            @click="prevMonth"
          >
            <ChevronLeft class="h-4 w-4" />
          </button>
          <div class="min-w-36 text-center text-sm font-semibold text-gray-900">{{ monthLabel }}</div>
          <button
            type="button"
            class="inline-flex items-center justify-center rounded-lg border border-gray-300 bg-white px-3 py-2 text-gray-700 hover:bg-gray-50"
            @click="nextMonth"
          >
            <ChevronRight class="h-4 w-4" />
          </button>
          <button
            type="button"
            class="rounded-lg border border-green-200 bg-green-50 px-3 py-2 text-sm font-medium text-green-700 hover:bg-green-100"
            @click="goToToday"
          >
            Hari Ini
          </button>
        </div>
      </div>

      <div class="mt-4 grid gap-4 lg:grid-cols-[minmax(0,1.6fr)_minmax(320px,0.9fr)]">
        <div>
          <div class="grid grid-cols-7 gap-2 text-center text-sm font-medium text-gray-500 mb-2">
            <div>Min</div>
            <div>Sen</div>
            <div>Sel</div>
            <div>Rab</div>
            <div>Kam</div>
            <div>Jum</div>
            <div>Sab</div>
          </div>

          <div class="grid grid-cols-7 gap-2">
            <button
              v-for="cell in calendarCells"
              :key="cell.key"
              type="button"
              class="group min-h-28 rounded-xl border p-3 text-left transition-all"
              :class="[
                cell.inMonth ? 'bg-white' : 'bg-gray-50 text-gray-300',
                cell.isSelected ? 'border-green-500 ring-2 ring-green-100' : 'border-gray-200 hover:border-green-300',
                cell.isToday && !cell.isSelected ? 'border-blue-300 bg-blue-50/50' : '',
              ]"
              @click="selectDay(cell.key)"
            >
              <div class="flex items-start justify-between gap-2">
                <span class="text-sm font-semibold" :class="cell.inMonth ? 'text-gray-900' : 'text-gray-300'">
                  {{ cell.date.getDate() }}
                </span>
                <span
                  v-if="cell.items.length > 0"
                  class="inline-flex items-center rounded-full px-2 py-0.5 text-[11px] font-semibold"
                  :class="dayHasDanger(cell.items)
                    ? 'bg-red-100 text-red-700'
                    : dayHasWarning(cell.items)
                      ? 'bg-yellow-100 text-yellow-700'
                      : 'bg-green-100 text-green-700'"
                >
                  {{ cell.items.length }}
                </span>
              </div>

              <div v-if="cell.items.length > 0" class="mt-3 space-y-2">
                <div class="flex flex-wrap gap-1.5">
                  <span
                    v-for="(item, index) in cell.items.slice(0, 3)"
                    :key="`${cell.key}-${item.id}-${index}`"
                    class="h-2.5 w-2.5 rounded-full"
                    :class="activityColorClass(item.status)"
                  />
                </div>
                <p class="text-xs text-gray-500">
                  {{ cell.items[0].activity_type }}
                </p>
              </div>
            </button>
          </div>
        </div>

        <div class="rounded-2xl border border-gray-200 bg-gray-50 p-5">
          <h3 class="text-lg font-semibold text-gray-900">Aktivitas - {{ selectedDateLabel }}</h3>
          <p class="mt-1 text-sm text-gray-600">Daftar task pada tanggal terpilih.</p>

          <div v-if="selectedDateSchedules.length === 0" class="mt-8 flex flex-col items-center justify-center py-10 text-center text-gray-500">
            <Calendar class="h-14 w-14 text-gray-300" />
            <p class="mt-4 text-sm">Tidak ada aktivitas terjadwal</p>
          </div>

          <div v-else class="mt-5 space-y-3 max-h-[520px] overflow-y-auto pr-1">
            <div
              v-for="item in selectedDateSchedules"
              :key="item.id"
              class="rounded-xl border border-gray-200 bg-white p-4 shadow-sm"
            >
              <div class="flex items-start justify-between gap-3">
                <div>
                  <p class="font-semibold text-gray-900">{{ item.activity_type }}</p>
                  <p class="text-sm text-gray-600 mt-1">Batch {{ item.batch_code }}</p>
                </div>
                <span :class="['inline-flex rounded-full px-2.5 py-1 text-xs font-medium', statusClass(item.status)]">
                  {{ statusLabel(item.status) }}
                </span>
              </div>
              <div class="mt-3 grid grid-cols-2 gap-3 text-sm text-gray-600">
                <div>
                  <p class="text-xs uppercase tracking-wide text-gray-400">Frekuensi</p>
                  <p class="mt-1 font-medium text-gray-900">{{ item.frequency_days }} hari</p>
                </div>
                <div>
                  <p class="text-xs uppercase tracking-wide text-gray-400">Jatuh Tempo</p>
                  <p class="mt-1 font-medium text-gray-900">{{ formatDate(item.next_due_date) }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="rounded-xl border border-gray-200 bg-white overflow-hidden">
      <div class="border-b border-gray-200 px-6 py-4">
        <h3 class="text-lg font-semibold text-gray-900">Daftar Jadwal User</h3>
      </div>

      <div v-if="isLoading" class="px-6 py-8 text-sm text-gray-600">Memuat data...</div>
      <div v-else-if="sortedSchedules.length === 0" class="px-6 py-8 text-sm text-gray-600">
        User ini belum memiliki jadwal.
      </div>
      <div v-else class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 text-left text-sm text-gray-700">
            <tr>
              <th class="px-6 py-3">Batch</th>
              <th class="px-6 py-3">Aktivitas</th>
              <th class="px-6 py-3">Frekuensi</th>
              <th class="px-6 py-3">Jatuh Tempo</th>
              <th class="px-6 py-3">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 text-sm">
            <tr v-for="item in sortedSchedules" :key="item.id" class="hover:bg-gray-50">
              <td class="px-6 py-3 text-gray-800">{{ item.batch_code || `#${item.plant_batch_id}` }}</td>
              <td class="px-6 py-3 text-gray-800">{{ item.activity_type }}</td>
              <td class="px-6 py-3 text-gray-800">{{ item.frequency_days }} hari</td>
              <td class="px-6 py-3 text-gray-800">{{ formatDate(item.next_due_date) }}</td>
              <td class="px-6 py-3">
                <span :class="['inline-flex rounded-full px-2.5 py-1 text-xs font-medium', statusClass(item.status)]">
                  {{ statusLabel(item.status) }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="isAddModalOpen" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
      <div class="w-full max-w-2xl rounded-2xl bg-white p-6">
        <h3 class="text-xl font-semibold text-gray-900">Tambah Jadwal - {{ user?.name || '-' }}</h3>
        <p class="mt-1 text-sm text-gray-600">User dipilih otomatis, langsung tambah jadwal baru.</p>

        <form class="mt-5 grid grid-cols-1 md:grid-cols-2 gap-4" @submit.prevent="submitSchedule">
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">Batch</label>
            <select
              v-model="form.plant_batch_id"
              class="w-full rounded-lg border border-gray-300 px-3 py-2 focus:border-green-500 focus:outline-none"
              required
            >
              <option value="">Pilih batch</option>
              <option v-for="batch in plantBatches" :key="batch.id" :value="String(batch.id)">
                {{ batch.batch_code }} - {{ batch.location }}
              </option>
            </select>
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">Tanggal Jatuh Tempo</label>
            <input
              v-model="form.next_due_date"
              type="date"
              class="w-full rounded-lg border border-gray-300 px-3 py-2 focus:border-green-500 focus:outline-none"
              disabled
              required
            />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">Jenis Aktivitas</label>
            <input
              v-model="form.activity_type"
              type="text"
              class="w-full rounded-lg border border-gray-300 px-3 py-2 focus:border-green-500 focus:outline-none"
              placeholder="Contoh: Pemupukan"
              required
            />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">Frekuensi (hari)</label>
            <input
              v-model.number="form.frequency_days"
              type="number"
              min="0"
              class="w-full rounded-lg border border-gray-300 px-3 py-2 focus:border-green-500 focus:outline-none"
              required
            />
          </div>

          <div class="md:col-span-2 flex gap-3 pt-2">
            <button
              type="button"
              class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
              @click="closeAddModal"
            >
              Batal
            </button>
            <button
              type="submit"
              class="inline-flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700 disabled:opacity-60"
              :disabled="isSaving"
            >
              <CalendarDays class="h-4 w-4" />
              {{ isSaving ? "Menyimpan..." : "Simpan Jadwal" }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
