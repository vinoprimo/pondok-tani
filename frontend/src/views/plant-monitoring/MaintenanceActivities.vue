<script setup lang="ts">
import { ref, reactive } from 'vue';
import {
  Calendar,
  CheckCircle,
  XCircle,
  Clock,
  AlertCircle,
  Plus,
  Image as ImageIcon,
} from 'lucide-vue-next';

const showSubmitForm = ref(false);
const selectedFile = ref<string | null>(null);
const formData = reactive({
  activityType: '',
  activityDate: '',
  description: '',
  quantity: '',
  notes: '',
});

const upcomingSchedule = [
  {
    id: 'MA-001',
    activityType: 'Penyiraman',
    dueDate: '2024-01-25',
    status: 'Upcoming',
    batchId: 'BATCH-045',
    daysLeft: 3,
  },
  {
    id: 'MA-002',
    activityType: 'Pemupukan',
    dueDate: '2024-01-26',
    status: 'Upcoming',
    batchId: 'BATCH-038',
    daysLeft: 4,
  },
  {
    id: 'MA-003',
    activityType: 'Pemangkasan',
    dueDate: '2024-01-23',
    status: 'Submitted',
    batchId: 'BATCH-045',
    daysLeft: 1,
  },
  {
    id: 'MA-004',
    activityType: 'Pengendalian Hama',
    dueDate: '2024-01-22',
    status: 'Verified',
    batchId: 'BATCH-052',
    daysLeft: 0,
  },
  {
    id: 'MA-005',
    activityType: 'Penyiraman',
    dueDate: '2024-01-20',
    status: 'Rejected',
    batchId: 'BATCH-038',
    daysLeft: -2,
    rejectionReason:
      'Kualitas foto kurang jelas - mohon kirim ulang dengan foto yang lebih jelas',
  },
];

const activityHistory = [
  {
    id: 'AH-001',
    activityType: 'Penyiraman',
    date: '2024-01-18',
    status: 'Verified',
    verifiedBy: 'Admin',
    verifiedDate: '2024-01-19',
  },
  {
    id: 'AH-002',
    activityType: 'Pemupukan',
    date: '2024-01-15',
    status: 'Verified',
    verifiedBy: 'Admin',
    verifiedDate: '2024-01-16',
  },
  {
    id: 'AH-003',
    activityType: 'Pemangkasan',
    date: '2024-01-12',
    status: 'Verified',
    verifiedBy: 'Admin',
    verifiedDate: '2024-01-13',
  },
];

function isScheduleOverdue(item: (typeof upcomingSchedule)[number]) {
  return item.daysLeft < 0 && item.status === 'Upcoming';
}

function handleFileUpload(e: Event) {
  const input = e.target as HTMLInputElement;
  const file = input.files?.[0];
  if (file) {
    const reader = new FileReader();
    reader.onloadend = () => {
      selectedFile.value = reader.result as string;
    };
    reader.readAsDataURL(file);
  }
}

function handleSubmit() {
  alert('Laporan aktivitas berhasil dikirim!');
  showSubmitForm.value = false;
  formData.activityType = '';
  formData.activityDate = '';
  formData.description = '';
  formData.quantity = '';
  formData.notes = '';
  selectedFile.value = null;
}

function closeModal() {
  showSubmitForm.value = false;
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-semibold text-gray-900">Aktivitas Perawatan</h2>
        <p class="text-gray-600 mt-1">
          Pantau dan kirimkan laporan aktivitas perawatan perkebunan Anda
        </p>
      </div>
      <button
        type="button"
        class="px-6 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-medium flex items-center gap-2"
        @click="showSubmitForm = true"
      >
        <Plus class="w-5 h-5" />
        Kirim Laporan Aktivitas
      </button>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Tugas Mendatang</p>
          <Clock class="w-5 h-5 text-blue-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">2</p>
        <p class="text-xs text-blue-600 mt-1">Berikutnya dalam 3 hari</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Menunggu Verifikasi</p>
          <Clock class="w-5 h-5 text-yellow-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">1</p>
        <p class="text-xs text-gray-500 mt-1">Menunggu tinjauan admin</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Terverifikasi Bulan Ini</p>
          <CheckCircle class="w-5 h-5 text-green-600" />
        </div>
        <p class="text-2xl font-semibold text-green-600">8</p>
        <p class="text-xs text-green-600 mt-1">+2 dari bulan lalu</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Terlambat/Ditolak</p>
          <AlertCircle class="w-5 h-5 text-red-600" />
        </div>
        <p class="text-2xl font-semibold text-red-600">1</p>
        <p class="text-xs text-red-600 mt-1">Perlu perhatian</p>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-200">
        <h3 class="text-lg font-semibold text-gray-900">Jadwal perawatan</h3>
        <p class="text-sm text-gray-600 mt-1">Tugas mendatang dan aktivitas terkini Anda</p>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">ID aktivitas</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Jenis aktivitas</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Batch</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Jatuh tempo</th>
              <th class="text-center px-6 py-3 text-sm font-medium text-gray-900">Status</th>
              <th class="text-center px-6 py-3 text-sm font-medium text-gray-900">Sisa hari</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr
              v-for="item in upcomingSchedule"
              :key="item.id"
              class="hover:bg-gray-50"
            >
              <td class="px-6 py-4 font-medium text-gray-900">{{ item.id }}</td>
              <td class="px-6 py-4">
                <div class="flex items-center gap-2">
                  <Calendar class="w-4 h-4 text-gray-400" />
                  <span class="text-gray-900">{{ item.activityType }}</span>
                </div>
              </td>
              <td class="px-6 py-4 text-gray-900">{{ item.batchId }}</td>
              <td class="px-6 py-4 text-gray-900">{{ item.dueDate }}</td>
              <td class="px-6 py-4 text-center">
                <span
                  v-if="isScheduleOverdue(item)"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800"
                >
                  <AlertCircle class="w-3 h-3 mr-1" />
                  Terlambat
                </span>
                <span
                  v-else-if="item.status === 'Upcoming'"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
                >
                  <Clock class="w-3 h-3 mr-1" />
                  Mendatang
                </span>
                <span
                  v-else-if="item.status === 'Submitted'"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800"
                >
                  <Clock class="w-3 h-3 mr-1" />
                  Menunggu verifikasi
                </span>
                <span
                  v-else-if="item.status === 'Verified'"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800"
                >
                  <CheckCircle class="w-3 h-3 mr-1" />
                  Terverifikasi
                </span>
                <span
                  v-else-if="item.status === 'Rejected'"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800"
                >
                  <XCircle class="w-3 h-3 mr-1" />
                  Ditolak
                </span>
              </td>
              <td class="px-6 py-4 text-center">
                <template v-if="item.status === 'Upcoming' || item.status === 'Submitted'">
                  <span
                    :class="[
                      'font-semibold',
                      isScheduleOverdue(item)
                        ? 'text-red-600'
                        : item.daysLeft <= 2
                          ? 'text-orange-600'
                          : 'text-gray-900',
                    ]"
                  >
                    {{
                      isScheduleOverdue(item)
                        ? `Terlambat ${Math.abs(item.daysLeft)} hari`
                        : `${item.daysLeft} hari`
                    }}
                  </span>
                </template>
                <span v-else class="text-gray-400">-</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-6">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">Riwayat aktivitas</h3>
      <div class="space-y-4">
        <div
          v-for="(activity, index) in activityHistory"
          :key="activity.id"
          class="flex gap-4"
        >
          <div class="flex flex-col items-center">
            <div class="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center">
              <CheckCircle class="w-5 h-5 text-green-600" />
            </div>
            <div
              v-if="index < activityHistory.length - 1"
              class="w-0.5 h-full bg-gray-200 mt-2"
            />
          </div>
          <div class="flex-1 pb-8">
            <div class="bg-gray-50 rounded-lg p-4 border border-gray-200">
              <div class="flex items-center justify-between mb-2">
                <h4 class="font-semibold text-gray-900">{{ activity.activityType }}</h4>
                <span class="text-xs text-gray-500">{{ activity.date }}</span>
              </div>
              <p class="text-sm text-gray-600 mb-2">ID aktivitas: {{ activity.id }}</p>
              <div class="flex items-center justify-between text-xs">
                <span class="text-green-600 font-medium">
                  Diverifikasi oleh {{ activity.verifiedBy }} pada {{ activity.verifiedDate }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div
      v-if="showSubmitForm"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
    >
      <div class="bg-white rounded-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div class="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
          <div>
            <h3 class="text-xl font-semibold text-gray-900">Kirim laporan aktivitas</h3>
            <p class="text-sm text-gray-600 mt-1">
              Isi detail dan bukti perawatan yang telah dilakukan
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
              Jenis aktivitas <span class="text-red-500">*</span>
            </label>
            <select
              v-model="formData.activityType"
              required
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent"
            >
              <option value="">Pilih jenis aktivitas</option>
              <option value="watering">Penyiraman</option>
              <option value="fertilizing">Pemupukan</option>
              <option value="pruning">Pemangkasan</option>
              <option value="pest-control">Pengendalian hama</option>
            </select>
          </div>

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
              Bukti foto <span class="text-red-500">*</span>
            </label>
            <div
              class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-green-500 transition-colors"
            >
              <div v-if="selectedFile" class="space-y-2">
                <img :src="selectedFile" alt="Pratinjau" class="max-h-48 mx-auto rounded-lg" />
                <button
                  type="button"
                  class="text-sm text-red-600 hover:text-red-700"
                  @click="selectedFile = null"
                >
                  Hapus gambar
                </button>
              </div>
              <template v-else>
                <ImageIcon class="w-12 h-12 text-gray-400 mx-auto mb-2" />
                <p class="text-sm text-gray-600 mb-2">Klik untuk mengunggah atau seret file</p>
                <input
                  id="file-upload"
                  type="file"
                  accept="image/*"
                  required
                  class="hidden"
                  @change="handleFileUpload"
                />
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
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Kuantitas material (opsional)
            </label>
            <input
              v-model="formData.quantity"
              type="text"
              placeholder="mis. 5 kg pupuk"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Catatan tambahan (opsional)
            </label>
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
              class="flex-1 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-medium"
            >
              Kirim laporan
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
