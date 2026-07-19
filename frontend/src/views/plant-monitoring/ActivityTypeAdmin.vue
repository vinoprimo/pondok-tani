<script setup lang="ts">
import { onMounted, ref, computed } from 'vue';
import { Plus, Edit2, Trash2 } from 'lucide-vue-next';
import SearchBar from '../../components/common/SearchBar.vue';
import Pagination from '../../components/common/Pagination.vue';
import {
  ActivityType,
  getActivityTypes,
  createActivityType,
  updateActivityType,
  deleteActivityType
} from '../../services/maintenance/activityType';

const items = ref<ActivityType[]>([]);
const isLoading = ref(false);
const currentPage = ref(1);
const limit = ref(10);
const totalPages = ref(1);
const totalRows = ref(0);
const searchTerm = ref('');

const showModal = ref(false);
const isEditing = ref(false);
const isSaving = ref(false);
const form = ref({
  id: 0,
  name: '',
  description: ''
});

async function loadData() {
  isLoading.value = true;
  try {
    const res = await getActivityTypes({
      page: currentPage.value,
      limit: limit.value,
      search: searchTerm.value
    });
    
    if (res.data && res.data.meta) {
      items.value = res.data.data;
      totalPages.value = res.data.meta.total_pages;
      totalRows.value = res.data.meta.total_rows;
      currentPage.value = res.data.meta.page;
    } else {
      items.value = res.data;
    }
  } catch (error) {
    console.error('Failed to load activity types:', error);
    items.value = [];
  } finally {
    isLoading.value = false;
  }
}

function handleSearch(val: string) {
  searchTerm.value = val;
  currentPage.value = 1;
  loadData();
}

function handlePageChange(val: number) {
  currentPage.value = val;
  loadData();
}

function openAddModal() {
  isEditing.value = false;
  form.value = { id: 0, name: '', description: '' };
  showModal.value = true;
}

function openEditModal(item: ActivityType) {
  isEditing.value = true;
  form.value = { ...item };
  showModal.value = true;
}

function closeModal() {
  showModal.value = false;
}

async function submitForm() {
  if (!form.value.name.trim()) {
    alert("Nama aktivitas wajib diisi");
    return;
  }
  
  isSaving.value = true;
  try {
    if (isEditing.value) {
      await updateActivityType(form.value.id, {
        name: form.value.name,
        description: form.value.description
      });
    } else {
      await createActivityType({
        name: form.value.name,
        description: form.value.description
      });
    }
    closeModal();
    loadData();
  } catch (error: any) {
    alert(error.response?.data?.error || "Terjadi kesalahan saat menyimpan data");
  } finally {
    isSaving.value = false;
  }
}

async function deleteItem(id: number) {
  if (!confirm('Apakah Anda yakin ingin menghapus jenis aktivitas ini?')) return;
  
  try {
    await deleteActivityType(id);
    loadData();
  } catch (error: any) {
    alert(error.response?.data?.error || "Gagal menghapus jenis aktivitas");
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
        <h2 class="text-2xl font-semibold text-gray-900">Jenis Aktivitas</h2>
        <p class="text-gray-600 mt-1">Kelola master data jenis aktivitas perawatan tanaman</p>
      </div>
      <button
        type="button"
        class="inline-flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700 transition-colors"
        @click="openAddModal"
      >
        <Plus class="h-4 w-4" />
        Tambah Jenis Aktivitas
      </button>
    </div>

    <div class="rounded-xl border border-gray-200 bg-white overflow-hidden">
      <div class="flex items-center justify-between border-b border-gray-200 px-6 py-4">
        <h3 class="text-lg font-semibold text-gray-900">Daftar Jenis Aktivitas</h3>
        <div class="w-full max-w-xs">
          <SearchBar v-model="searchTerm" placeholder="Cari aktivitas..." @search="handleSearch" />
        </div>
      </div>

      <div v-if="isLoading" class="px-6 py-8 text-center text-sm text-gray-500">Memuat data...</div>
      <div v-else-if="items.length === 0" class="px-6 py-8 text-center text-sm text-gray-500">
        Belum ada jenis aktivitas.
      </div>
      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-sm text-gray-600">
          <thead class="bg-gray-50 border-b border-gray-200 text-gray-700 font-medium">
            <tr>
              <th class="px-6 py-3 w-1/3">Nama Aktivitas</th>
              <th class="px-6 py-3">Deskripsi</th>
              <th class="px-6 py-3 w-32 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="item in items" :key="item.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 font-medium text-gray-900">{{ item.name }}</td>
              <td class="px-6 py-4">{{ item.description || '-' }}</td>
              <td class="px-6 py-4 text-right">
                <div class="flex items-center justify-end gap-2">
                  <button
                    type="button"
                    class="p-1.5 text-blue-600 hover:bg-blue-50 rounded-md transition-colors"
                    title="Edit"
                    @click="openEditModal(item)"
                  >
                    <Edit2 class="w-4 h-4" />
                  </button>
                  <button
                    type="button"
                    class="p-1.5 text-red-600 hover:bg-red-50 rounded-md transition-colors"
                    title="Hapus"
                    @click="deleteItem(item.id)"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <Pagination
        v-if="items.length > 0"
        :current-page="currentPage"
        :total-pages="totalPages"
        :total-rows="totalRows"
        :limit="limit"
        @update:page="handlePageChange"
      />
    </div>

    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div class="w-full max-w-md rounded-xl bg-white shadow-xl">
        <div class="border-b border-gray-200 px-6 py-4">
          <h3 class="text-lg font-semibold text-gray-900">
            {{ isEditing ? 'Edit Jenis Aktivitas' : 'Tambah Jenis Aktivitas' }}
          </h3>
        </div>
        <form @submit.prevent="submitForm" class="p-6">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Nama Aktivitas *</label>
              <input
                v-model="form.name"
                type="text"
                class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-green-500 focus:outline-none focus:ring-1 focus:ring-green-500"
                placeholder="Contoh: Pemupukan"
                required
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Deskripsi</label>
              <textarea
                v-model="form.description"
                rows="3"
                class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-green-500 focus:outline-none focus:ring-1 focus:ring-green-500"
                placeholder="Penjelasan opsional"
              ></textarea>
            </div>
          </div>
          <div class="mt-6 flex items-center justify-end gap-3">
            <button
              type="button"
              class="rounded-lg px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100"
              @click="closeModal"
            >
              Batal
            </button>
            <button
              type="submit"
              class="rounded-lg bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700 disabled:opacity-50"
              :disabled="isSaving"
            >
              {{ isSaving ? 'Menyimpan...' : 'Simpan' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
