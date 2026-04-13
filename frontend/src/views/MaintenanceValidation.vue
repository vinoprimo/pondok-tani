<script setup lang="ts">
import { ref, computed } from 'vue';
import { CheckCircle, XCircle, Eye, User, Calendar } from 'lucide-vue-next';

const selectedActivity = ref<string | null>(null);
const showReviewModal = ref(false);
const rejectionReason = ref('');

const submittedActivities = [
  {
    id: 'MA-003',
    userName: 'Sarah Johnson',
    userId: 'INV-001',
    batchId: 'BATCH-045',
    activityType: 'Pemangkasan',
    submissionDate: '2024-01-23',
    activityDate: '2024-01-23',
    description:
      'Pemangkasan rutin sulur vanili, menghapus daun mati dan tunas berlebih untuk merangsang berbunga.',
    quantity: '15 sulur dipangkas',
    notes: 'Semua alat disterilisasi sebelum digunakan',
    photoUrl: 'https://images.unsplash.com/photo-1464226184884-fa280b87c399?w=800',
    status: 'Pending',
  },
  {
    id: 'MA-006',
    userName: 'Michael Chen',
    userId: 'INV-002',
    batchId: 'BATCH-038',
    activityType: 'Pemupukan',
    submissionDate: '2024-01-23',
    activityDate: '2024-01-22',
    description: 'Pupuk organik diterapkan pada seluruh tanaman. Formula NPK 15-15-15.',
    quantity: '25 kg pupuk',
    notes: 'Cuaca mendukung, tidak ada hujan dalam 24 jam',
    photoUrl: 'https://images.unsplash.com/photo-1416879595882-3373a0480b5b?w=800',
    status: 'Pending',
  },
  {
    id: 'MA-007',
    userName: 'Emma Davis',
    userId: 'INV-003',
    batchId: 'BATCH-052',
    activityType: 'Penyiraman',
    submissionDate: '2024-01-22',
    activityDate: '2024-01-22',
    description: 'Penyiraman mendalam pagi hari untuk seluruh tanaman.',
    quantity: '500 liter',
    notes: 'Kelembaban tanah dicek sebelum menyiram',
    photoUrl: 'https://images.unsplash.com/photo-1523348837708-15d4a09cfac2?w=800',
    status: 'Pending',
  },
];

const selectedActivityData = computed(() =>
  submittedActivities.find((a) => a.id === selectedActivity.value)
);

function handleReview(activityId: string) {
  selectedActivity.value = activityId;
  showReviewModal.value = true;
}

function handleApprove() {
  alert(`Aktivitas ${selectedActivity.value} berhasil disetujui!`);
  showReviewModal.value = false;
  selectedActivity.value = null;
}

function handleReject() {
  if (!rejectionReason.value.trim()) {
    alert('Mohon isi alasan penolakan');
    return;
  }
  alert(`Aktivitas ${selectedActivity.value} ditolak. Alasan: ${rejectionReason.value}`);
  showReviewModal.value = false;
  selectedActivity.value = null;
  rejectionReason.value = '';
}

function closeReviewModal() {
  showReviewModal.value = false;
  rejectionReason.value = '';
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-semibold text-gray-900">Validasi perawatan</h2>
      <p class="text-gray-600 mt-1">
        Tinjau dan validasi aktivitas perawatan yang dikirim investor
      </p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Menunggu tinjauan</p>
          <Calendar class="w-5 h-5 text-yellow-600" />
        </div>
        <p class="text-2xl font-semibold text-yellow-600">3</p>
        <p class="text-xs text-gray-500 mt-1">Perlu perhatian</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Disetujui hari ini</p>
          <CheckCircle class="w-5 h-5 text-green-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">12</p>
        <p class="text-xs text-green-600 mt-1">+8 dari kemarin</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Ditolak hari ini</p>
          <XCircle class="w-5 h-5 text-red-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">2</p>
        <p class="text-xs text-red-600 mt-1">Masalah kualitas</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Total bulan ini</p>
          <Calendar class="w-5 h-5 text-blue-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">87</p>
        <p class="text-xs text-gray-500 mt-1">Aktivitas ditinjau</p>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-200">
        <h3 class="text-lg font-semibold text-gray-900">Aktivitas menunggu</h3>
        <p class="text-sm text-gray-600 mt-1">Tinjau laporan perawatan yang dikirim</p>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">ID aktivitas</th>
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
              <td class="px-6 py-4 font-medium text-gray-900">{{ activity.id }}</td>
              <td class="px-6 py-4">
                <div class="flex items-center gap-2">
                  <div class="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
                    <User class="w-4 h-4 text-green-600" />
                  </div>
                  <div>
                    <p class="font-medium text-gray-900">{{ activity.userName }}</p>
                    <p class="text-xs text-gray-500">{{ activity.userId }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 text-gray-900">{{ activity.batchId }}</td>
              <td class="px-6 py-4">
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
                >
                  {{ activity.activityType }}
                </span>
              </td>
              <td class="px-6 py-4 text-gray-900">{{ activity.activityDate }}</td>
              <td class="px-6 py-4 text-gray-900">{{ activity.submissionDate }}</td>
              <td class="px-6 py-4 text-center">
                <div class="flex justify-center">
                  <img
                    :src="activity.photoUrl"
                    alt="Pratinjau aktivitas"
                    class="w-12 h-12 rounded-lg object-cover cursor-pointer hover:opacity-75 transition-opacity"
                    @click="handleReview(activity.id)"
                  />
                </div>
              </td>
              <td class="px-6 py-4 text-center">
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800"
                >
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
          </tbody>
        </table>
      </div>
    </div>

    <div
      v-if="showReviewModal && selectedActivityData"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
    >
      <div class="bg-white rounded-2xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
        <div class="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
          <div>
            <h3 class="text-xl font-semibold text-gray-900">
              Tinjau aktivitas — {{ selectedActivityData.id }}
            </h3>
            <p class="text-sm text-gray-600 mt-1">
              Dikirim oleh {{ selectedActivityData.userName }}
            </p>
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
                :src="selectedActivityData.photoUrl"
                alt="Bukti aktivitas"
                class="w-full h-96 object-cover"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="bg-gray-50 rounded-lg p-4">
              <p class="text-sm text-gray-600 mb-1">Jenis aktivitas</p>
              <p class="font-semibold text-gray-900">{{ selectedActivityData.activityType }}</p>
            </div>
            <div class="bg-gray-50 rounded-lg p-4">
              <p class="text-sm text-gray-600 mb-1">ID batch</p>
              <p class="font-semibold text-gray-900">{{ selectedActivityData.batchId }}</p>
            </div>
            <div class="bg-gray-50 rounded-lg p-4">
              <p class="text-sm text-gray-600 mb-1">Tanggal aktivitas</p>
              <p class="font-semibold text-gray-900">{{ selectedActivityData.activityDate }}</p>
            </div>
            <div class="bg-gray-50 rounded-lg p-4">
              <p class="text-sm text-gray-600 mb-1">Tanggal pengiriman</p>
              <p class="font-semibold text-gray-900">{{ selectedActivityData.submissionDate }}</p>
            </div>
          </div>

          <div>
            <h4 class="text-sm font-medium text-gray-700 mb-2">Deskripsi</h4>
            <div class="bg-gray-50 rounded-lg p-4">
              <p class="text-gray-900">{{ selectedActivityData.description }}</p>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <h4 class="text-sm font-medium text-gray-700 mb-2">Kuantitas / material</h4>
              <div class="bg-gray-50 rounded-lg p-4">
                <p class="text-gray-900">{{ selectedActivityData.quantity }}</p>
              </div>
            </div>
            <div>
              <h4 class="text-sm font-medium text-gray-700 mb-2">Catatan tambahan</h4>
              <div class="bg-gray-50 rounded-lg p-4">
                <p class="text-gray-900">{{ selectedActivityData.notes }}</p>
              </div>
            </div>
          </div>

          <div>
            <h4 class="text-sm font-medium text-gray-700 mb-2">Alasan penolakan (opsional)</h4>
            <textarea
              v-model="rejectionReason"
              rows="3"
              placeholder="Jika menolak, berikan alasan yang jelas..."
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          <div class="flex gap-3 pt-4">
            <button
              type="button"
              class="flex-1 px-6 py-3 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors font-medium flex items-center justify-center gap-2"
              @click="handleReject"
            >
              <XCircle class="w-5 h-5" />
              Tolak aktivitas
            </button>
            <button
              type="button"
              class="flex-1 px-6 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-medium flex items-center justify-center gap-2"
              @click="handleApprove"
            >
              <CheckCircle class="w-5 h-5" />
              Setujui aktivitas
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
