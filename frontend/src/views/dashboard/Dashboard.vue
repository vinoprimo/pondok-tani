<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import {
  TrendingUp,
  Sprout,
  Banknote,
  Users,
  ArrowUpRight,
  ArrowDownRight,
  AlertCircle,
  X,
  Calculator,
} from "lucide-vue-next";
import { getAdminDashboard, getUserDashboard } from "../../services/dashboard/dashboard";
import { Bar, Doughnut, Line } from "vue-chartjs";
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale,
  ArcElement,
  PointElement,
  LineElement
} from "chart.js";

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale, ArcElement, PointElement, LineElement);

type UserRole = "investor" | "mitra" | "admin";

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

const userRole = ref<UserRole>(
  (localStorage.getItem("userRole") as UserRole) || "investor"
);
const isAdmin = computed(() => userRole.value === "admin");

const showProjectionModal = ref(false);

const stats = ref<any[]>([]);
const monthlyData = ref<any[]>([]);
const plantStatusData = ref<any[]>([]);
const revenueComparison = ref<any[]>([]);
const salesHistory = ref<any[]>([]);
const recentActivity = ref<any[]>([]);
const upcomingMaintenance = ref<any[]>([]);
const harvestStatus = ref<any[]>([]);
const roiProgress = ref<any>(null);
const loading = ref(true);

const loadDashboardData = async () => {
  loading.value = true;
  try {
    let data;
    if (isAdmin.value) {
      data = await getAdminDashboard();
    } else {
      data = await getUserDashboard();
    }
    stats.value = data.stats || [];
    monthlyData.value = data.monthly_data || [];
    plantStatusData.value = data.plant_status_data || [];
    revenueComparison.value = data.revenue_comparison || [];
    salesHistory.value = data.sales_history || [];
    recentActivity.value = data.recent_activity || [];
    if (data.upcoming_maintenance) upcomingMaintenance.value = data.upcoming_maintenance;
    if (data.harvest_status) harvestStatus.value = data.harvest_status;
    if (data.roi_progress) roiProgress.value = data.roi_progress;
  } catch (err) {
    console.error("Gagal memuat dashboard", err);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  const stored = localStorage.getItem("userRole") as UserRole | null;
  if (stored === "admin" || stored === "investor" || stored === "mitra") {
    userRole.value = stored;
  }

  loadDashboardData();
});

const monthlyChartData = computed(() => ({
  labels: monthlyData.value.map(d => d.month),
  datasets: [
    {
      label: 'Investasi',
      backgroundColor: '#10b981',
      data: monthlyData.value.map(d => d.revenue)
    },
    {
      label: 'Biaya',
      backgroundColor: '#3b82f6',
      data: monthlyData.value.map(d => d.costs)
    }
  ]
}));

const monthlyChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    y: {
      beginAtZero: true,
      suggestedMax: 1000000,
      ticks: {
        callback: function(value: any) {
          return new Intl.NumberFormat('id-ID', {
            style: 'currency',
            currency: 'IDR',
            minimumFractionDigits: 0,
            maximumFractionDigits: 0
          }).format(value);
        }
      }
    }
  },
  plugins: {
    tooltip: {
      callbacks: {
        label: function(context: any) {
          let label = context.dataset.label || '';
          if (label) {
            label += ': ';
          }
          if (context.parsed.y !== null) {
            label += new Intl.NumberFormat('id-ID', {
              style: 'currency',
              currency: 'IDR',
              minimumFractionDigits: 0,
              maximumFractionDigits: 0
            }).format(context.parsed.y);
          }
          return label;
        }
      }
    }
  }
};

const plantStatusChartData = computed(() => ({
  labels: plantStatusData.value.map(d => d.name),
  datasets: [
    {
      backgroundColor: plantStatusData.value.map(d => d.color),
      data: plantStatusData.value.map(d => d.value)
    }
  ]
}));

const plantStatusChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'bottom' as const
    }
  }
};

// Removed revenueComparisonChartData and revenueComparisonChartOptions

const iconMap: Record<string, any> = {
  Banknote,
  TrendingUp,
  Sprout,
  Users
};



function activityDotClass(type: string) {
  if (type === "success") return "bg-green-500";
  if (type === "warning") return "bg-yellow-500";
  if (type === "error") return "bg-red-500";
  return "bg-blue-500";
}

// projection modal removed
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
            <component :is="iconMap[stat.icon]" class="w-6 h-6 text-green-600" />
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

    <div v-if="loading" class="flex justify-center items-center py-20">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-green-600"></div>
    </div>
    
    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="lg:col-span-2 bg-white rounded-xl border border-gray-200 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Kinerja investasi</h3>
        <div class="h-80">
          <Bar :data="monthlyChartData" :options="monthlyChartOptions" />
        </div>
      </div>



      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Status tanaman</h3>
        <div class="h-64 mt-4">
          <Doughnut :data="plantStatusChartData" :options="plantStatusChartOptions" />
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
            @click="$router.push('/dashboard/maintenance')"
          >
            Lihat semua
          </button>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div
            v-for="(task, index) in upcomingMaintenance"
            :key="index"
            :class="[
              'bg-white rounded-lg p-4 border-2',
              task.is_overdue ? 'border-red-300' : 'border-orange-200'
            ]"
          >
            <div class="flex items-start justify-between mb-2">
              <div>
                <p class="font-semibold text-gray-900">{{ task.task_name }}</p>
                <p class="text-sm text-gray-600 mt-1">Jatuh tempo: {{ task.due_date }}</p>
              </div>
              <span
                :class="[
                  'px-3 py-1 rounded-full text-xs font-bold',
                  task.is_overdue ? 'bg-red-100 text-red-700' : 'bg-orange-100 text-orange-700'
                ]"
              >
                {{ task.is_overdue ? `Terlambat ${Math.abs(task.days_diff)} hari` : `Tersisa ${task.days_diff} hari` }}
              </span>
            </div>
            <p
              :class="[
                'text-sm mt-2',
                task.is_overdue ? 'text-red-600' : 'text-gray-600'
              ]"
            >
              {{ task.is_overdue ? 'Perlu tindakan segera' : 'Segera datang' }}
            </p>
          </div>
          <div v-if="upcomingMaintenance.length === 0" class="col-span-full py-4 text-center text-gray-500">
            Tidak ada tugas perawatan dalam waktu dekat.
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div class="bg-white rounded-xl border border-gray-200 p-6 lg:col-span-2">
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

      <!-- Removed Perkiraan vs pendapatan aktual and Kemajuan ROI -->

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

    <!-- Modal removed -->
  </div>
</template>
