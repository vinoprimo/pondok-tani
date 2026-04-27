<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import {
  AlertCircle,
  Calendar,
  CheckCircle,
  Clock,
  Image as ImageIcon,
  Send,
  XCircle,
} from "lucide-vue-next";
import {
  getMyMaintenanceActivities,
  getMyMaintenanceSchedules,
  submitMyMaintenanceActivity,
} from "../../services/maintenance/schedule";

type ScheduleItem = {
  id: number;
  schedule_code: string;
  plant_batch_id: number;
  batch_code: string;
  activity_type: string;
  frequency_days: number;
  next_due_date: string;
  status: string;
  remark?: string;
  has_submitted_proof: boolean;
};

type ActivityHistoryItem = {
  id: number;
  schedule_id: number;
  schedule_code: string;
  batch_code: string;
  activity_type: string;
  description?: string;
  activity_date: string;
  photo_url?: string;
  validation_status: string;
  validation_notes?: string;
  remark?: string;
  validated_at?: string;
};

const isLoading = ref(false);
const isSubmitting = ref(false);
const schedules = ref<ScheduleItem[]>([]);
const activityHistory = ref<ActivityHistoryItem[]>([]);
const showSubmitForm = ref(false);
const selectedFile = ref<string | null>(null);
const selectedPhotoFile = ref<File | null>(null);
const selectedSchedule = ref<ScheduleItem | null>(null);

const formData = reactive({
  activityDate: "",
  description: "",
  notes: "",
});

const upcomingSchedule = computed(() =>
  [...schedules.value].sort(
    (a, b) => new Date(a.next_due_date).getTime() - new Date(b.next_due_date).getTime()
  )
);

const statUpcoming = computed(() => upcomingSchedule.value.filter((item) => item.status === "mendatang").length);
const statWaitingValidation = computed(() =>
  upcomingSchedule.value.filter((item) => item.status === "menunggu_verifikasi").length
);
const statVerified = computed(() =>
  upcomingSchedule.value.filter((item) => item.status === "terverifikasi").length
);
const statOverdueOrRejected = computed(() =>
  upcomingSchedule.value.filter((item) => item.status === "overdue" || item.status === "ditolak").length
);

function startOfToday() {
  const date = new Date();
  date.setHours(0, 0, 0, 0);
  return date;
}

function computeDaysLeft(dueDate: string) {
  const due = new Date(dueDate);
  due.setHours(0, 0, 0, 0);
  const diffMs = due.getTime() - startOfToday().getTime();
  return Math.ceil(diffMs / (1000 * 60 * 60 * 24));
}

function formatDate(dateStr?: string | null) {
  if (!dateStr) return "-";
  const date = new Date(dateStr);
  if (Number.isNaN(date.getTime())) return "-";
  return date.toLocaleDateString("id-ID", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  });
}

function activityPhotoUrl(url?: string | null) {
  if (!url) return "";
  if (url.startsWith("http")) return url;
  return `http://localhost:8000${url}`;
}

function scheduleStatusLabel(status: string) {
  if (status === "overdue") return "Overdue";
  if (status === "menunggu_verifikasi") return "Menunggu verifikasi";
  if (status === "terverifikasi") return "Terverifikasi";
  if (status === "ditolak") return "Ditolak";
  return "Mendatang";
}

function scheduleStatusClass(status: string) {
  if (status === "overdue") return "bg-red-100 text-red-700";
  if (status === "menunggu_verifikasi") return "bg-yellow-100 text-yellow-700";
  if (status === "terverifikasi") return "bg-green-100 text-green-700";
  if (status === "ditolak") return "bg-red-100 text-red-700";
  return "bg-blue-100 text-blue-700";
}

function validationLabel(status: string) {
  if (status === "approved" || status === "verified") return "Diverifikasi";
  if (status === "rejected") return "Ditolak";
  return "Menunggu verifikasi";
}

function validationClass(status: string) {
  if (status === "approved" || status === "verified") return "text-green-600";
  if (status === "rejected") return "text-red-600";
  return "text-yellow-600";
}

function canSubmitReport(item: ScheduleItem) {
  return item.status === "mendatang" || item.status === "overdue" || item.status === "ditolak";
}

function openSubmitModal(item: ScheduleItem) {
  selectedSchedule.value = item;
  showSubmitForm.value = true;
  formData.activityDate = new Date().toISOString().slice(0, 10);
  formData.description = "";
  formData.notes = "";
  selectedFile.value = null;
  selectedPhotoFile.value = null;
}

function closeModal() {
  showSubmitForm.value = false;
  selectedSchedule.value = null;
  selectedPhotoFile.value = null;
}

function handleFileUpload(e: Event) {
  const input = e.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;

  selectedPhotoFile.value = file;

  const reader = new FileReader();
  reader.onloadend = () => {
    selectedFile.value = reader.result as string;
  };
  reader.readAsDataURL(file);
}

async function loadData() {
  isLoading.value = true;
  try {
    const [scheduleRes, historyRes] = await Promise.all([
      getMyMaintenanceSchedules(),
      getMyMaintenanceActivities(),
    ]);
    schedules.value = Array.isArray(scheduleRes.data) ? scheduleRes.data : [];
    activityHistory.value = Array.isArray(historyRes.data) ? historyRes.data : [];
  } finally {
    isLoading.value = false;
  }
}

async function handleSubmit() {
  if (!selectedSchedule.value || !formData.activityDate || !selectedPhotoFile.value) {
    alert("Data laporan belum lengkap.");
    return;
  }

  isSubmitting.value = true;
  try {
    const description = [formData.description, formData.notes].filter(Boolean).join("\n").trim();
    const payload = new FormData();
    payload.append("schedule_id", String(selectedSchedule.value.id));
    payload.append("activity_date", formData.activityDate);
    if (description) payload.append("description", description);
    payload.append("photo", selectedPhotoFile.value);

    await submitMyMaintenanceActivity(payload);

    await loadData();
    closeModal();
  } catch (error: any) {
    const message = error?.response?.data?.error || "Gagal mengirim laporan aktivitas";
    alert(message);
  } finally {
    isSubmitting.value = false;
  }
}

onMounted(() => {
  loadData();
});
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-semibold text-gray-900">Aktivitas Perawatan</h2>
        <p class="text-gray-600 mt-1">Pantau dan kirimkan laporan aktivitas perawatan perkebunan Anda</p>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Tugas Mendatang</p>
          <Clock class="w-5 h-5 text-blue-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">{{ statUpcoming }}</p>
        <p class="text-xs text-blue-600 mt-1">Status mendatang</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Menunggu Verifikasi</p>
          <Clock class="w-5 h-5 text-yellow-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">{{ statWaitingValidation }}</p>
        <p class="text-xs text-gray-500 mt-1">Menunggu tinjauan admin</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Terverifikasi</p>
          <CheckCircle class="w-5 h-5 text-green-600" />
        </div>
        <p class="text-2xl font-semibold text-green-600">{{ statVerified }}</p>
        <p class="text-xs text-green-600 mt-1">Aktivitas disetujui admin</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Terlambat/Ditolak</p>
          <AlertCircle class="w-5 h-5 text-red-600" />
        </div>
        <p class="text-2xl font-semibold text-red-600">{{ statOverdueOrRejected }}</p>
        <p class="text-xs text-red-600 mt-1">Perlu perhatian</p>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-200">
        <h3 class="text-lg font-semibold text-gray-900">Jadwal perawatan</h3>
        <p class="text-sm text-gray-600 mt-1">Tugas mendatang dan aktivitas terkini Anda</p>
      </div>
      <div v-if="isLoading" class="px-6 py-8 text-sm text-gray-600">Memuat data...</div>
      <div v-else class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Jenis aktivitas</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Batch</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Jatuh tempo</th>
              <th class="text-center px-6 py-3 text-sm font-medium text-gray-900">Status</th>
              <th class="text-center px-6 py-3 text-sm font-medium text-gray-900">Sisa hari</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Remark</th>
              <th class="text-center px-6 py-3 text-sm font-medium text-gray-900">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="item in upcomingSchedule" :key="item.id" class="hover:bg-gray-50">
              <td class="px-6 py-4">
                <div class="flex items-center gap-2">
                  <Calendar class="w-4 h-4 text-gray-400" />
                  <span class="text-gray-900">{{ item.activity_type }}</span>
                </div>
              </td>
              <td class="px-6 py-4 text-gray-900">{{ item.batch_code || `BATCH-${item.plant_batch_id}` }}</td>
              <td class="px-6 py-4 text-gray-900">{{ formatDate(item.next_due_date) }}</td>
              <td class="px-6 py-4 text-center">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium" :class="scheduleStatusClass(item.status)">
                  {{ scheduleStatusLabel(item.status) }}
                </span>
              </td>
              <td class="px-6 py-4 text-center">
                <span class="font-semibold" :class="computeDaysLeft(item.next_due_date) < 0 ? 'text-red-600' : computeDaysLeft(item.next_due_date) <= 2 ? 'text-orange-600' : 'text-gray-900'">
                  {{ computeDaysLeft(item.next_due_date) < 0 ? `Terlambat ${Math.abs(computeDaysLeft(item.next_due_date))} hari` : `${computeDaysLeft(item.next_due_date)} hari` }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-700">{{ item.remark || '-' }}</td>
              <td class="px-6 py-4 text-center">
                <button
                  v-if="canSubmitReport(item)"
                  type="button"
                  class="inline-flex items-center gap-2 rounded-lg bg-green-600 px-3 py-2 text-xs font-medium text-white hover:bg-green-700"
                  @click="openSubmitModal(item)"
                >
                  <Send class="w-3.5 h-3.5" />
                  Kirim Laporan
                </button>
                <span v-else class="text-xs text-gray-400">-</span>
              </td>
            </tr>
            <tr v-if="!upcomingSchedule.length">
              <td colspan="7" class="px-6 py-8 text-center text-sm text-gray-500">Belum ada jadwal perawatan.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-6">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">Riwayat aktivitas</h3>
      <div class="space-y-4">
        <div v-for="(activity, index) in activityHistory" :key="activity.id" class="flex gap-4">
          <div class="flex flex-col items-center">
            <div class="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center">
              <CheckCircle class="w-5 h-5 text-green-600" />
            </div>
            <div v-if="index < activityHistory.length - 1" class="w-0.5 h-full bg-gray-200 mt-2" />
          </div>
          <div class="flex-1 pb-8">
            <div class="bg-gray-50 rounded-lg p-4 border border-gray-200">
              <div class="flex items-center justify-between mb-2">
                <h4 class="font-semibold text-gray-900">{{ activity.activity_type }}</h4>
                <span class="text-xs text-gray-500">{{ formatDate(activity.activity_date) }}</span>
              </div>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-3 mb-3">
                <div>
                  <p class="text-xs text-gray-500">Batch</p>
                  <p class="text-sm text-gray-700">{{ activity.batch_code || '-' }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Remark</p>
                  <p class="text-sm text-gray-700">{{ activity.remark || '-' }}</p>
                </div>
              </div>
              <div class="mb-3">
                <p class="text-xs text-gray-500">Deskripsi</p>
                <p class="text-sm text-gray-700">{{ activity.description || '-' }}</p>
              </div>
              <div class="mb-3">
                <p class="text-xs text-gray-500 mb-1">Gambar</p>
                <img
                  v-if="activity.photo_url"
                  :src="activityPhotoUrl(activity.photo_url)"
                  alt="Bukti aktivitas"
                  class="h-24 w-24 rounded-lg object-cover border border-gray-200"
                />
                <p v-else class="text-sm text-gray-500">-</p>
              </div>
              <div class="flex items-center justify-between text-xs">
                <span class="font-medium" :class="validationClass(activity.validation_status)">
                  {{ validationLabel(activity.validation_status) }}
                </span>
                <span v-if="activity.validated_at" class="text-gray-500">{{ formatDate(activity.validated_at) }}</span>
              </div>
            </div>
          </div>
        </div>
        <p v-if="!activityHistory.length" class="text-sm text-gray-500">Belum ada riwayat aktivitas.</p>
      </div>
    </div>

    <div v-if="showSubmitForm" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div class="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
          <div>
            <h3 class="text-xl font-semibold text-gray-900">Kirim laporan aktivitas</h3>
            <p class="text-sm text-gray-600 mt-1">
              {{ selectedSchedule ? `${selectedSchedule.activity_type} - ${selectedSchedule.batch_code}` : '' }}
            </p>
          </div>
          <button
            type="button"
            class="w-8 h-8 rounded-lg hover:bg-gray-100 flex items-center justify-center transition-colors"
            @click="closeModal"
          >
            <XCircle class="w-5 h-5 text-gray-500" />
          </button>
        </div>

        <form class="p-6 space-y-4" @submit.prevent="handleSubmit">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Tanggal aktivitas <span class="text-red-500">*</span>
            </label>
            <input
              v-model="formData.activityDate"
              type="date"
              required
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Deskripsi <span class="text-red-500">*</span>
            </label>
            <textarea
              v-model="formData.description"
              required
              rows="3"
              placeholder="Jelaskan aktivitas perawatan yang dilakukan..."
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Bukti foto
            </label>
            <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-green-500 transition-colors">
              <div v-if="selectedFile" class="space-y-2">
                <img :src="selectedFile" alt="Pratinjau" class="max-h-48 mx-auto rounded-lg" />
                <button
                  type="button"
                  class="text-sm text-red-600 hover:text-red-700"
                  @click="selectedFile = null; selectedPhotoFile = null"
                >
                  Hapus gambar
                </button>
              </div>
              <template v-else>
                <ImageIcon class="w-12 h-12 text-gray-400 mx-auto mb-2" />
                <p class="text-sm text-gray-600 mb-2">Klik untuk mengunggah bukti foto</p>
                <input id="file-upload" type="file" accept="image/*" class="hidden" @change="handleFileUpload" />
                <label
                  for="file-upload"
                  class="inline-block px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors cursor-pointer"
                >
                  Pilih file
                </label>
              </template>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Catatan tambahan (opsional)</label>
            <textarea
              v-model="formData.notes"
              rows="2"
              placeholder="Informasi tambahan..."
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent"
            />
          </div>

          <div class="flex gap-3 pt-4">
            <button
              type="button"
              class="flex-1 px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors font-medium text-gray-700"
              @click="closeModal"
            >
              Batal
            </button>
            <button
              type="submit"
              class="flex-1 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-medium disabled:opacity-60"
              :disabled="isSubmitting"
            >
              {{ isSubmitting ? "Mengirim..." : "Kirim laporan" }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
