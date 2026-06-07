<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { Calendar, CheckCircle, Eye, User, XCircle } from "lucide-vue-next";
import {
  getAdminMaintenanceActivities,
  getAdminMaintenanceActivitiesSummary,
  reviewAdminMaintenanceActivity,
} from "../../services/maintenance/schedule";
import SearchBar from "../../components/common/SearchBar.vue";
import Pagination from "../../components/common/Pagination.vue";

type Summary = {
  waiting_count: number;
  approved_today: number;
  rejected_today: number;
  reviewed_month: number;
};

type ActivityItem = {
  id: number;
  activity_code: string;
  user_id: string;
  user_name: string;
  batch_code: string;
  activity_type: string;
  activity_date: string;
  submission_date: string;
  description?: string;
  photo_url?: string;
  status: string;
};

const isLoading = ref(false);
const isReviewing = ref(false);
const summary = ref<Summary>({
  waiting_count: 0,
  approved_today: 0,
  rejected_today: 0,
  reviewed_month: 0,
});
const submittedActivities = ref<ActivityItem[]>([]);
const currentPage = ref(1);
const totalPages = ref(1);
const totalRows = ref(0);
const limit = ref(10);
const searchTerm = ref("");

function handleSearch(val: string) {
  searchTerm.value = val;
  currentPage.value = 1;
  loadActivities();
}

function handlePageChange(val: number) {
  currentPage.value = val;
  loadActivities();
}

const selectedActivityId = ref<number | null>(null);
const showReviewModal = ref(false);
const showRejectReasonModal = ref(false);
const rejectionReason = ref("");

const selectedActivityData = computed(() =>
  submittedActivities.value.find((a) => a.id === selectedActivityId.value)
);

function formatDate(dateStr?: string | null) {
  if (!dateStr) return "-";
  const date = new Date(dateStr);
  if (Number.isNaN(date.getTime())) return "-";
  return date.toLocaleDateString("id-ID", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  });
}

function photoUrl(url?: string) {
  if (!url) return "";
  if (url.startsWith("http")) return url;
  return `http://localhost:8000${url}`;
}

async function loadSummary() {
  try {
    const summaryRes = await getAdminMaintenanceActivitiesSummary();
    summary.value = summaryRes.data || summary.value;
  } catch (err) {
    console.error("Failed to load summary", err);
  }
}

async function loadActivities() {
  isLoading.value = true;
  try {
    const activitiesRes = await getAdminMaintenanceActivities({
      status: "pending",
      page: currentPage.value,
      limit: limit.value,
      search: searchTerm.value,
    });
    
    if (activitiesRes.data && activitiesRes.data.meta) {
      submittedActivities.value = Array.isArray(activitiesRes.data.data) ? activitiesRes.data.data : [];
      totalPages.value = activitiesRes.data.meta.total_pages;
      totalRows.value = activitiesRes.data.meta.total_rows;
      currentPage.value = activitiesRes.data.meta.page;
    } else {
      submittedActivities.value = Array.isArray(activitiesRes.data) ? activitiesRes.data : [];
      totalPages.value = 1;
      totalRows.value = submittedActivities.value.length;
    }
  } catch (err) {
    console.error("Failed to load activities", err);
  } finally {
    isLoading.value = false;
  }
}

async function loadData() {
  await Promise.all([loadSummary(), loadActivities()]);
}

function handleReview(activityId: number) {
  selectedActivityId.value = activityId;
  showReviewModal.value = true;
  rejectionReason.value = "";
}

async function handleApprove() {
  if (!selectedActivityId.value) return;
  isReviewing.value = true;
  try {
    await reviewAdminMaintenanceActivity(selectedActivityId.value, { action: "approve" });
    showReviewModal.value = false;
    selectedActivityId.value = null;
    await loadData();
  } catch (error: any) {
    alert(error?.response?.data?.error || "Gagal menyetujui aktivitas");
  } finally {
    isReviewing.value = false;
  }
}

async function handleReject() {
  if (!selectedActivityId.value) return;
  if (!rejectionReason.value.trim()) {
    alert("Mohon isi alasan penolakan");
    return;
  }

  isReviewing.value = true;
  try {
    await reviewAdminMaintenanceActivity(selectedActivityId.value, {
      action: "reject",
      notes: rejectionReason.value.trim(),
    });
    showRejectReasonModal.value = false;
    showReviewModal.value = false;
    selectedActivityId.value = null;
    rejectionReason.value = "";
    await loadData();
  } catch (error: any) {
    alert(error?.response?.data?.error || "Gagal menolak aktivitas");
  } finally {
    isReviewing.value = false;
  }
}

function openRejectReasonModal() {
  rejectionReason.value = "";
  showRejectReasonModal.value = true;
}

function closeRejectReasonModal() {
  showRejectReasonModal.value = false;
  rejectionReason.value = "";
}

function closeReviewModal() {
  showReviewModal.value = false;
  showRejectReasonModal.value = false;
  selectedActivityId.value = null;
  rejectionReason.value = "";
}

onMounted(() => {
  loadData();
});
</script>

<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-semibold text-gray-900">Validasi perawatan</h2>
      <p class="text-gray-600 mt-1">Tinjau dan validasi aktivitas perawatan yang dikirim investor</p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Menunggu tinjauan</p>
          <Calendar class="w-5 h-5 text-yellow-600" />
        </div>
        <p class="text-2xl font-semibold text-yellow-600">{{ summary.waiting_count }}</p>
        <p class="text-xs text-gray-500 mt-1">Perlu perhatian</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Disetujui hari ini</p>
          <CheckCircle class="w-5 h-5 text-green-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">{{ summary.approved_today }}</p>
        <p class="text-xs text-green-600 mt-1">Aktivitas lolos validasi</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Ditolak hari ini</p>
          <XCircle class="w-5 h-5 text-red-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">{{ summary.rejected_today }}</p>
        <p class="text-xs text-red-600 mt-1">Masalah kualitas</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Total bulan ini</p>
          <Calendar class="w-5 h-5 text-blue-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">{{ summary.reviewed_month }}</p>
        <p class="text-xs text-gray-500 mt-1">Aktivitas ditinjau</p>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden">
      <div class="flex flex-wrap items-center justify-between gap-3 border-b border-gray-200 px-6 py-4">
        <div>
          <h3 class="text-lg font-semibold text-gray-900">Aktivitas menunggu</h3>
          <p class="text-sm text-gray-600 mt-1">Tinjau laporan perawatan yang dikirim</p>
        </div>
        <div class="w-full max-w-sm">
          <SearchBar
            v-model="searchTerm"
            placeholder="Cari pengguna, batch, aktivitas..."
            @search="handleSearch"
          />
        </div>
      </div>
      <div v-if="isLoading" class="px-6 py-8 text-sm text-gray-600">Memuat data...</div>
      <div v-else class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Pengguna</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Batch</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Jenis aktivitas</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Tanggal aktivitas</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Tanggal pengiriman</th>
              <th class="text-center px-6 py-3 text-sm font-medium text-gray-900">Foto</th>
              <th class="text-center px-6 py-3 text-sm font-medium text-gray-900">Status</th>
              <th class="text-center px-6 py-3 text-sm font-medium text-gray-900">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="activity in submittedActivities" :key="activity.id" class="hover:bg-gray-50">
              <td class="px-6 py-4">
                <div class="flex items-center gap-2">
                  <div class="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
                    <User class="w-4 h-4 text-green-600" />
                  </div>
                  <div>
                    <p class="font-medium text-gray-900">{{ activity.user_name }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 text-gray-900">{{ activity.batch_code }}</td>
              <td class="px-6 py-4">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                  {{ activity.activity_type }}
                </span>
              </td>
              <td class="px-6 py-4 text-gray-900">{{ formatDate(activity.activity_date) }}</td>
              <td class="px-6 py-4 text-gray-900">{{ formatDate(activity.submission_date) }}</td>
              <td class="px-6 py-4 text-center">
                <div class="flex justify-center">
                  <img
                    v-if="activity.photo_url"
                    :src="photoUrl(activity.photo_url)"
                    alt="Pratinjau aktivitas"
                    class="w-12 h-12 rounded-lg object-cover cursor-pointer hover:opacity-75 transition-opacity"
                    @click="handleReview(activity.id)"
                  />
                  <span v-else class="text-xs text-gray-400">-</span>
                </div>
              </td>
              <td class="px-6 py-4 text-center">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
                  Menunggu
                </span>
              </td>
              <td class="px-6 py-4 text-center">
                <button
                  type="button"
                  class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm font-medium flex items-center gap-2 mx-auto"
                  @click="handleReview(activity.id)"
                >
                  <Eye class="w-4 h-4" />
                  Tinjau
                </button>
              </td>
            </tr>
            <tr v-if="!submittedActivities.length">
              <td colspan="8" class="px-6 py-10 text-center text-sm text-gray-500">Belum ada aktivitas menunggu validasi.</td>
            </tr>
          </tbody>
        </table>
      </div>
      <Pagination
        v-if="submittedActivities.length > 0"
        :current-page="currentPage"
        :total-pages="totalPages"
        :total-rows="totalRows"
        :limit="limit"
        @update:page="handlePageChange"
      />
    </div>

    <div
      v-if="showReviewModal && selectedActivityData"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
    >
      <div class="bg-white rounded-2xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
        <div class="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
          <div>
            <h3 class="text-xl font-semibold text-gray-900">Tinjau aktivitas - {{ selectedActivityData.activity_code }}</h3>
            <p class="text-sm text-gray-600 mt-1">Dikirim oleh {{ selectedActivityData.user_name }}</p>
          </div>
          <button
            type="button"
            class="w-8 h-8 rounded-lg hover:bg-gray-100 flex items-center justify-center transition-colors"
            @click="closeReviewModal"
          >
            <XCircle class="w-5 h-5 text-gray-500" />
          </button>
        </div>

        <div class="p-6 space-y-6">
          <div>
            <h4 class="text-sm font-medium text-gray-700 mb-3">Bukti foto</h4>
            <div class="border border-gray-200 rounded-xl overflow-hidden">
              <img
                v-if="selectedActivityData.photo_url"
                :src="photoUrl(selectedActivityData.photo_url)"
                alt="Bukti aktivitas"
                class="w-full h-96 object-contain bg-gray-50"
              />
              <div v-else class="h-40 flex items-center justify-center text-sm text-gray-500">Tidak ada foto</div>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="bg-gray-50 rounded-lg p-4">
              <p class="text-sm text-gray-600 mb-1">Jenis aktivitas</p>
              <p class="font-semibold text-gray-900">{{ selectedActivityData.activity_type }}</p>
            </div>
            <div class="bg-gray-50 rounded-lg p-4">
              <p class="text-sm text-gray-600 mb-1">ID batch</p>
              <p class="font-semibold text-gray-900">{{ selectedActivityData.batch_code }}</p>
            </div>
            <div class="bg-gray-50 rounded-lg p-4">
              <p class="text-sm text-gray-600 mb-1">Tanggal aktivitas</p>
              <p class="font-semibold text-gray-900">{{ formatDate(selectedActivityData.activity_date) }}</p>
            </div>
            <div class="bg-gray-50 rounded-lg p-4">
              <p class="text-sm text-gray-600 mb-1">Tanggal pengiriman</p>
              <p class="font-semibold text-gray-900">{{ formatDate(selectedActivityData.submission_date) }}</p>
            </div>
          </div>

          <div>
            <h4 class="text-sm font-medium text-gray-700 mb-2">Deskripsi</h4>
            <div class="bg-gray-50 rounded-lg p-4">
              <p class="text-gray-900">{{ selectedActivityData.description || '-' }}</p>
            </div>
          </div>

          <div class="flex gap-3 pt-4">
            <button
              type="button"
              class="flex-1 px-6 py-3 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors font-medium flex items-center justify-center gap-2 disabled:opacity-60"
              @click="openRejectReasonModal"
              :disabled="isReviewing"
            >
              <XCircle class="w-5 h-5" />
              Tolak aktivitas
            </button>
            <button
              type="button"
              class="flex-1 px-6 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-medium flex items-center justify-center gap-2 disabled:opacity-60"
              @click="handleApprove"
              :disabled="isReviewing"
            >
              <CheckCircle class="w-5 h-5" />
              {{ isReviewing ? 'Memproses...' : 'Setujui aktivitas' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <div
      v-if="showRejectReasonModal"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-[60] p-4"
    >
      <div class="bg-white rounded-xl w-full max-w-lg border border-gray-200">
        <div class="px-5 py-4 border-b border-gray-200">
          <h4 class="text-lg font-semibold text-gray-900">Alasan penolakan</h4>
          <p class="text-sm text-gray-600 mt-1">Wajib diisi sebelum menolak aktivitas.</p>
        </div>
        <div class="p-5 space-y-4">
          <textarea
            v-model="rejectionReason"
            rows="4"
            placeholder="Tuliskan alasan penolakan yang jelas..."
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent"
          />
          <div class="flex gap-3">
            <button
              type="button"
              class="flex-1 px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 text-gray-700"
              :disabled="isReviewing"
              @click="closeRejectReasonModal"
            >
              Batal
            </button>
            <button
              type="button"
              class="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-60"
              :disabled="isReviewing"
              @click="handleReject"
            >
              {{ isReviewing ? 'Memproses...' : 'Konfirmasi Tolak' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
