<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { AlertCircle, Calendar, CheckCircle, Clock, RefreshCcw } from "lucide-vue-next";
import { getAdminUsers } from "../../services/user/user";
import {
  getAdminMaintenanceScheduleSummary,
  getAdminMaintenanceSchedules,
} from "../../services/maintenance/schedule";

type UserItem = {
  id: string;
  name: string;
  email: string;
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

type UserScheduleSummary = {
  id: string;
  name: string;
  email: string;
  total: number;
  nearestActivity: string | null;
  nearestDueDate: string | null;
  status: "pending" | "overdue" | "aman";
};

const router = useRouter();
const users = ref<UserItem[]>([]);
const schedules = ref<ScheduleItem[]>([]);
const summary = ref({
  pending_count: 0,
  due_today_count: 0,
  overdue_count: 0,
  total_count: 0,
});

const isLoading = ref(false);

const todayDateOnly = computed(() => {
  const now = new Date();
  now.setHours(0, 0, 0, 0);
  return now;
});

const pendingCount = computed(() => summary.value.pending_count);
const dueTodayCount = computed(() => summary.value.due_today_count);
const overdueCount = computed(() => summary.value.overdue_count);
const totalScheduleCount = computed(() => summary.value.total_count);

const userSummaries = computed<UserScheduleSummary[]>(() => {
  return users.value.map((user) => {
    const userSchedules = schedules.value.filter((s) => s.user_id === user.id);
    const total = userSchedules.length;

    let nearestActivity: string | null = null;
    let nearestDueDate: string | null = null;
    if (total > 0) {
      const nearestSchedule = [...userSchedules].sort(
        (a, b) => new Date(a.next_due_date).getTime() - new Date(b.next_due_date).getTime()
      )[0];
      nearestActivity = nearestSchedule.activity_type;
      nearestDueDate = nearestSchedule.next_due_date;
    }

    const hasOverdue = userSchedules.some((s) => {
      const due = new Date(s.next_due_date);
      due.setHours(0, 0, 0, 0);
      return s.status === "overdue" || (s.status === "pending" && due.getTime() < todayDateOnly.value.getTime());
    });

    const hasPending = userSchedules.some((s) => s.status === "pending");

    let status: "pending" | "overdue" | "aman" = "aman";
    if (hasOverdue) status = "overdue";
    else if (hasPending) status = "pending";

    return {
      id: user.id,
      name: user.name,
      email: user.email,
      total,
      nearestActivity,
      nearestDueDate,
      status,
    };
  });
});

async function loadUsers() {
  const res = await getAdminUsers();
  users.value = Array.isArray(res.data)
    ? res.data.map((item: any) => ({
        id: item.id,
        name: item.name,
        email: item.email,
      }))
    : [];
}

async function loadSchedules() {
  isLoading.value = true;
  try {
    const [summaryRes, schedulesRes] = await Promise.all([
      getAdminMaintenanceScheduleSummary(),
      getAdminMaintenanceSchedules(),
    ]);
    summary.value = summaryRes.data || summary.value;
    schedules.value = Array.isArray(schedulesRes.data) ? schedulesRes.data : [];
  } finally {
    isLoading.value = false;
  }
}

function goToUserDetail(userId: string) {
  router.push(`/dashboard/monitoring-schedule/${userId}`);
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
  if (status === "pending") return "Pending";
  return "Aman";
}

function statusClass(status: string) {
  if (status === "overdue") return "bg-red-50 text-red-700";
  if (status === "pending") return "bg-yellow-50 text-yellow-700";
  return "bg-green-50 text-green-700";
}

onMounted(async () => {
  await loadUsers();
  await loadSchedules();
});
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-semibold text-gray-900">Jadwal Monitoring</h2>
        <p class="text-gray-600 mt-1">Daftar ringkasan jadwal per user dan kelola task monitoring.</p>
      </div>
      <button
        type="button"
        class="inline-flex items-center gap-2 rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
        @click="loadSchedules"
      >
        <RefreshCcw class="h-4 w-4" />
        Muat Ulang
      </button>
    </div>

    <div class="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4">
      <div class="rounded-xl border border-gray-200 bg-white p-6">
        <div class="flex items-start justify-between">
          <p class="text-sm text-gray-600">Menunggu Tindakan</p>
          <Clock class="h-5 w-5 text-yellow-600" />
        </div>
        <p class="mt-3 text-4xl font-semibold text-gray-900">{{ pendingCount }}</p>
        <p class="mt-2 text-sm text-gray-500">Jadwal pending yang belum lewat jatuh tempo</p>
      </div>

      <div class="rounded-xl border border-gray-200 bg-white p-6">
        <div class="flex items-start justify-between">
          <p class="text-sm text-gray-600">Jatuh Tempo Hari Ini</p>
          <Calendar class="h-5 w-5 text-blue-600" />
        </div>
        <p class="mt-3 text-4xl font-semibold text-gray-900">{{ dueTodayCount }}</p>
        <p class="mt-2 text-sm text-gray-500">Perlu diprioritaskan hari ini</p>
      </div>

      <div class="rounded-xl border border-gray-200 bg-white p-6">
        <div class="flex items-start justify-between">
          <p class="text-sm text-gray-600">Overdue</p>
          <AlertCircle class="h-5 w-5 text-red-600" />
        </div>
        <p class="mt-3 text-4xl font-semibold text-gray-900">{{ overdueCount }}</p>
        <p class="mt-2 text-sm text-red-600">Butuh tindak lanjut admin</p>
      </div>

      <div class="rounded-xl border border-gray-200 bg-white p-6">
        <div class="flex items-start justify-between">
          <p class="text-sm text-gray-600">Total Jadwal Monitoring</p>
          <CheckCircle class="h-5 w-5 text-green-600" />
        </div>
        <p class="mt-3 text-4xl font-semibold text-gray-900">{{ totalScheduleCount }}</p>
        <p class="mt-2 text-sm text-gray-500">Semua task yang tercatat</p>
      </div>
    </div>

    <div class="rounded-xl border border-gray-200 bg-white overflow-hidden">
      <div class="border-b border-gray-200 px-6 py-4">
        <h3 class="text-lg font-semibold text-gray-900">Daftar Maintenance Schedule</h3>
      </div>

      <div v-if="isLoading" class="px-6 py-8 text-sm text-gray-600">Memuat data...</div>
      <div v-else-if="userSummaries.length === 0" class="px-6 py-8 text-sm text-gray-600">
        Belum ada user.
      </div>
      <div v-else class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 text-left text-sm text-gray-700">
            <tr>
              <th class="px-6 py-3">User</th>
              <th class="px-6 py-3">Email</th>
              <th class="px-6 py-3">Total Jadwal</th>
              <th class="px-6 py-3">Jadwal Terdekat</th>
              <th class="px-6 py-3">Jatuh Tempo</th>
              <th class="px-6 py-3">Status</th>
              <th class="px-6 py-3">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 text-sm">
            <tr v-for="summary in userSummaries" :key="summary.id" class="hover:bg-gray-50">
              <td class="px-6 py-3 font-medium text-gray-900">{{ summary.name }}</td>
              <td class="px-6 py-3 text-gray-700">{{ summary.email || "-" }}</td>
              <td class="px-6 py-3 text-gray-800">{{ summary.total }}</td>
              <td class="px-6 py-3 text-gray-800">{{ summary.nearestActivity || "-" }}</td>
              <td class="px-6 py-3 text-gray-800">{{ formatDate(summary.nearestDueDate) }}</td>
              <td class="px-6 py-3">
                <span :class="['inline-flex rounded-full px-2.5 py-1 text-xs font-medium', statusClass(summary.status)]">
                  {{ statusLabel(summary.status) }}
                </span>
              </td>
              <td class="px-6 py-3">
                <div class="flex flex-wrap gap-2">
                  <button
                    type="button"
                    class="rounded-lg border border-gray-300 px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50"
                    @click="goToUserDetail(summary.id)"
                  >
                    Lihat Detail
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
