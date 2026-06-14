<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import {
  Users,
  DollarSign,
  TrendingUp,
  UserCheck,
  Search,
  Filter,
  Mail,
  Phone,
  Package,
  Loader2,
  Eye,
  MapPin,
  Sprout,
  Plus,
  Pencil,
  Trash2,
  Image,
} from 'lucide-vue-next';
import {
  activateUserPackage,
  updateInvestmentStatus,
  getAdminUsers,
  getAdminUserPlantBatches,
  createAdminUser,
  updateAdminUser,
  deleteAdminUser,
  validateMitra,
} from '../../services/user/user';
import SearchBar from '../../components/common/SearchBar.vue';
import Pagination from '../../components/common/Pagination.vue';

const users = ref<any[]>([]);
const loading = ref(false);
const activatingUserId = ref<string | null>(null);
const updatingUserId = ref<string | null>(null);
const validatingMitraId = ref<string | null>(null);
const activationFormUser = ref<any | null>(null);
const activationFormError = ref('');
const activationForms = ref<any[]>([]);
const reviewingUser = ref<any | null>(null);
const reviewLoading = ref(false);
const reviewPlantBatches = ref<any[]>([]);
const userFormMode = ref<'create' | 'edit'>('create');
const userFormOpen = ref(false);
const userFormLoading = ref(false);
const userFormError = ref('');
const editingUser = ref<any | null>(null);
const userForm = ref({
  name: '',
  email: '',
  password: '',
  role: 'investor',
  phone_number: '',
  address: '',
});
const searchTerm = ref('');
const statusFilter = ref('all');

function getUserInvestmentItems(user: any) {
  if (Array.isArray(user.investments) && user.investments.length > 0) {
    return user.investments;
  }

  if (user.selected_package) {
    return [
      {
        id: `legacy-${user.id}`,
        package_id: user.selected_package_id,
        package: user.selected_package,
        status: user.package_status,
      },
    ];
  }

  return [];
}

const visibleUsers = computed(() => {
  const keyword = searchTerm.value.trim().toLowerCase();
  return users.value.filter((user) => {
    const packageStatus = user.package_status || 'none';
    const matchesKeyword =
      !keyword ||
      user.name?.toLowerCase().includes(keyword) ||
      user.email?.toLowerCase().includes(keyword) ||
      user.id?.toLowerCase().includes(keyword);
    const matchesStatus = statusFilter.value === 'all' || packageStatus === statusFilter.value;
    return matchesKeyword && matchesStatus;
  });
});

const totalSelected = computed(() => users.value.filter((user) => getUserInvestmentItems(user).length > 0).length);
const activeUsers = computed(() => users.value.filter((user) => user.package_status === 'active').length);
const pendingUsers = computed(() => users.value.filter((user) => user.package_status === 'pending').length);
const activeRatio = computed(() => {
  if (!users.value.length) return '0';
  return ((activeUsers.value / users.value.length) * 100).toFixed(0);
});

function investorInitials(name: string) {
  return (
    name
      ?.split(' ')
      .filter(Boolean)
      .map((part) => part[0])
      .join('')
      .slice(0, 2) || '?'
  );
}

function statusBadgeClass(status: string) {
  if (status === 'active') return 'bg-green-100 text-green-700';
  if (status === 'on_process') return 'bg-blue-100 text-blue-700';
  if (status === 'pending_validation') return 'bg-amber-100 text-amber-700';
  if (status === 'pending') return 'bg-yellow-100 text-yellow-700';
  return 'bg-gray-100 text-gray-700';
}

function investorStatusLabel(status: string) {
  if (status === 'active') return 'Aktif';
  if (status === 'on_process') return 'Proses';
  if (status === 'pending_validation') return 'Menunggu validasi';
  if (status === 'pending') return 'Menunggu';
  return 'Belum memilih';
}

function getActionState(user: any) {
  if (user.role === 'mitra' && user.package_status === 'pending_validation') return 'validate-mitra';
  if (getUserInvestmentItems(user).length && user.package_status === 'pending') return 'mark-paid';
  if (getUserInvestmentItems(user).length && user.package_status === 'on_process') return 'activate';
  if (getUserInvestmentItems(user).length && user.package_status === 'active') return 'review';
  return 'none';
}

function formatDisplayDate(value: string) {
  if (!value) return '-';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleDateString('id-ID');
}

const currentPage = ref(1);
const totalPages = ref(1);
const totalRows = ref(0);
const limit = ref(10);

async function loadUsers() {
  loading.value = true;
  try {
    const response = await getAdminUsers({
      page: currentPage.value,
      limit: limit.value,
      search: searchTerm.value,
    });
    if (response.data && response.data.meta) {
      users.value = Array.isArray(response.data.data) ? response.data.data : [];
      totalPages.value = response.data.meta.total_pages;
      totalRows.value = response.data.meta.total_rows;
      currentPage.value = response.data.meta.page;
    } else {
      users.value = Array.isArray(response.data) ? response.data : [];
      totalPages.value = 1;
      totalRows.value = users.value.length;
    }
  } catch (error) {
    console.error('Failed to load users:', error);
    users.value = [];
  } finally {
    loading.value = false;
  }
}

function handleSearch(val: string) {
  searchTerm.value = val;
  currentPage.value = 1;
  loadUsers();
}

function handlePageChange(val: number) {
  currentPage.value = val;
  loadUsers();
}

async function handleMarkAsPaid(userId: string) {
  updatingUserId.value = userId;
  try {
    await updateInvestmentStatus(userId, 'on_process');
    await loadUsers();
  } catch (error) {
    console.error('Failed to mark as paid:', error);
  } finally {
    updatingUserId.value = null;
  }
}

async function handleValidateMitra(userId: string) {
  validatingMitraId.value = userId;
  try {
    await validateMitra(userId);
    await loadUsers();
  } catch (error) {
    console.error('Failed to validate mitra:', error);
  } finally {
    validatingMitraId.value = null;
  }
}

async function handleActivate(userId: string) {
  if (!activationForms.value.length) {
    activationFormError.value = 'Data paket on process tidak ditemukan';
    return;
  }

  for (const item of activationForms.value) {
    if (!String(item.batchCode || '').trim()) {
      activationFormError.value = `Batch code wajib diisi untuk paket ${item.packageName}`;
      return;
    }
    if (!item.plantingDate) {
      activationFormError.value = `Tanggal tanam wajib diisi untuk paket ${item.packageName}`;
      return;
    }
    if (!String(item.location || '').trim()) {
      activationFormError.value = `Lokasi wajib diisi untuk paket ${item.packageName}`;
      return;
    }

    const seedCount = Number(item.seedCount);
    if (!Number.isFinite(seedCount) || seedCount <= 0) {
      activationFormError.value = `Jumlah bibit harus lebih dari 0 untuk paket ${item.packageName}`;
      return;
    }

    const landArea = Number(item.landArea);
    if (!Number.isFinite(landArea) || landArea <= 0) {
      activationFormError.value = `Luas lahan harus lebih dari 0 untuk paket ${item.packageName}`;
      return;
    }

    if (!(item.imageFile instanceof File)) {
      activationFormError.value = `Foto batch wajib diunggah untuk paket ${item.packageName}`;
      return;
    }
  }

  activatingUserId.value = userId;
  activationFormError.value = '';
  try {
    const plantBatchesPayload = activationForms.value.map((item) => ({
        package_id: item.packageId,
        batch_code: String(item.batchCode || '').trim(),
        planting_date: item.plantingDate,
        location: String(item.location || '').trim(),
        seed_count: Math.floor(Number(item.seedCount)),
        land_area: Number(item.landArea),
    }));

    const formData = new FormData();
    formData.append('plant_batches', JSON.stringify(plantBatchesPayload));
    activationForms.value.forEach((item) => {
      if (item.imageFile instanceof File) {
        formData.append(`batch_image_${item.packageId}`, item.imageFile);
      }
    });

    await activateUserPackage(userId, formData);
    await loadUsers();
    closeActivationForm();
  } catch (error) {
    console.error('Failed to activate package:', error);
    activationFormError.value = 'Aktivasi gagal. Periksa data plant batch dan coba lagi.';
  } finally {
    activatingUserId.value = null;
  }
}

function getOnProcessInvestmentItems(user: any) {
  return getUserInvestmentItems(user).filter((investment: any) => investment.status === 'on_process');
}

function openActivationForm(user: any) {
  const onProcessItems = getOnProcessInvestmentItems(user);
  if (!onProcessItems.length) {
    activationFormError.value = 'User ini belum memiliki paket on process';
    return;
  }

  activationFormUser.value = user;
  activationFormError.value = '';
  activationForms.value = onProcessItems.map((investment: any, index: number) => ({
    packageId: investment.package_id,
    packageName: investment.package?.package_name || `Paket #${investment.package_id}`,
    batchCode: `PB-${String(user.id).slice(0, 6).toUpperCase()}-${investment.package_id}-${index + 1}`,
    plantingDate: new Date().toISOString().slice(0, 10),
    location: '',
    seedCount: '',
    landArea: '',
    imageFile: null,
    imagePreview: '',
  }));
}

function handleBatchImageChange(formItem: any, event: Event) {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0] ?? null;

  if (formItem.imagePreview) {
    URL.revokeObjectURL(formItem.imagePreview);
  }

  formItem.imageFile = file;
  formItem.imagePreview = file ? URL.createObjectURL(file) : '';
}

function closeActivationForm() {
  activationForms.value.forEach((item) => {
    if (item.imagePreview) {
      URL.revokeObjectURL(item.imagePreview);
    }
  });
  activationFormUser.value = null;
  activationFormError.value = '';
  activationForms.value = [];
}

async function openReviewPlantBatch(user: any) {
  reviewingUser.value = user;
  reviewLoading.value = true;
  reviewPlantBatches.value = [];
  try {
    const response = await getAdminUserPlantBatches(user.id);
    reviewPlantBatches.value = Array.isArray(response.data) ? response.data : [];
  } catch (error) {
    console.error('Failed to load plant batch data:', error);
    reviewPlantBatches.value = [];
  } finally {
    reviewLoading.value = false;
  }
}

function closeReviewPlantBatch() {
  reviewingUser.value = null;
  reviewPlantBatches.value = [];
}

function openCreateUserForm() {
  userFormMode.value = 'create';
  editingUser.value = null;
  userFormError.value = '';
  userForm.value = {
    name: '',
    email: '',
    password: '',
    role: 'investor',
    phone_number: '',
    address: '',
  };
  userFormOpen.value = true;
}

function openEditUserForm(user: any) {
  userFormMode.value = 'edit';
  editingUser.value = user;
  userFormError.value = '';
  userForm.value = {
    name: user.name || '',
    email: user.email || '',
    password: '',
    role: user.role || 'investor',
    phone_number: user.phone_number || '',
    address: user.address || '',
  };
  userFormOpen.value = true;
}

function closeUserForm() {
  userFormOpen.value = false;
  userFormError.value = '';
}

async function submitUserForm() {
  if (!userForm.value.name.trim() || !userForm.value.email.trim()) {
    userFormError.value = 'Nama dan email wajib diisi';
    return;
  }
  if (userFormMode.value === 'create' && userForm.value.password.trim().length < 6) {
    userFormError.value = 'Password minimal 6 karakter';
    return;
  }

  userFormLoading.value = true;
  userFormError.value = '';
  try {
    const payload: any = {
      name: userForm.value.name.trim(),
      email: userForm.value.email.trim(),
      role: userForm.value.role,
      phone_number: userForm.value.phone_number.trim() || null,
      address: userForm.value.address.trim() || null,
    };

    if (userForm.value.password.trim()) {
      payload.password = userForm.value.password.trim();
    }

    if (userFormMode.value === 'create') {
      await createAdminUser(payload);
    } else if (editingUser.value) {
      await updateAdminUser(editingUser.value.id, payload);
    }

    await loadUsers();
    closeUserForm();
  } catch (error) {
    console.error('Failed to save user:', error);
    userFormError.value = 'Gagal menyimpan user. Cek data lalu coba lagi.';
  } finally {
    userFormLoading.value = false;
  }
}

async function handleDeleteUser(user: any) {
  const confirmed = window.confirm(`Hapus user ${user.name}? Tindakan ini tidak bisa dibatalkan.`);
  if (!confirmed) return;

  try {
    await deleteAdminUser(user.id);
    await loadUsers();
  } catch (error) {
    console.error('Failed to delete user:', error);
  }
}

onMounted(() => {
  loadUsers();
});
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-semibold text-gray-900">Manajemen user</h2>
        <p class="text-gray-600 mt-1">Kelola investor, mitra, dan validasi akun</p>
      </div>

      <button
        type="button"
        class="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        @click="openCreateUserForm"
      >
        <Plus class="w-4 h-4" />
        Tambah User
      </button>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-12 h-12 bg-blue-50 rounded-lg flex items-center justify-center">
            <Users class="w-6 h-6 text-blue-600" />
          </div>
          <span class="text-sm text-gray-600">Total user</span>
        </div>
        <p class="text-3xl font-semibold text-gray-900">{{ users.length }}</p>
        <p class="text-sm text-green-600 mt-1">{{ activeUsers }} aktif</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-12 h-12 bg-yellow-50 rounded-lg flex items-center justify-center">
            <Package class="w-6 h-6 text-yellow-600" />
          </div>
          <span class="text-sm text-gray-600">Menunggu aktivasi</span>
        </div>
        <p class="text-3xl font-semibold text-gray-900">{{ pendingUsers }}</p>
        <p class="text-sm text-gray-600 mt-1">Paket dipilih namun belum aktif</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-12 h-12 bg-green-50 rounded-lg flex items-center justify-center">
            <TrendingUp class="w-6 h-6 text-green-600" />
          </div>
          <span class="text-sm text-gray-600">Sudah memilih paket</span>
        </div>
        <p class="text-3xl font-semibold text-gray-900">{{ totalSelected }}</p>
        <p class="text-sm text-gray-600 mt-1">User dengan paket tersimpan</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-12 h-12 bg-purple-50 rounded-lg flex items-center justify-center">
            <DollarSign class="w-6 h-6 text-purple-600" />
          </div>
          <span class="text-sm text-gray-600">Rasio aktif</span>
        </div>
        <p class="text-3xl font-semibold text-gray-900">{{ activeRatio }}%</p>
        <p class="text-sm text-gray-600 mt-1">Dari seluruh user</p>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-4">
      <div class="flex flex-wrap items-center gap-4">
        <div class="flex-1 min-w-[300px]">
          <SearchBar 
            v-model="searchTerm" 
            @search="handleSearch" 
            placeholder="Cari nama, email, atau ID user..." 
          />
        </div>
        <button
          type="button"
          class="flex items-center gap-2 px-4 py-2 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
        >
          <Filter class="w-4 h-4" />
          Saring
        </button>
        <select
          v-model="statusFilter"
          class="px-4 py-2 border border-gray-200 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-green-500"
        >
          <option value="all">Semua status</option>
          <option value="active">Aktif</option>
          <option value="on_process">Proses</option>
          <option value="pending_validation">Menunggu validasi</option>
          <option value="pending">Menunggu</option>
          <option value="none">Belum memilih</option>
        </select>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden">
      <div v-if="loading" class="px-6 py-10 text-center text-gray-500 flex items-center justify-center gap-2">
        <Loader2 class="h-4 w-4 animate-spin" />
        Memuat data investor...
      </div>
      <div v-else class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                User
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Kontak
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Paket Dipilih
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Bergabung
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Aksi
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="user in visibleUsers" :key="user.id" class="hover:bg-gray-50">
              <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 bg-gradient-to-br from-green-400 to-green-600 rounded-full flex items-center justify-center text-white font-semibold">
                    {{ investorInitials(user.name) }}
                  </div>
                  <div>
                    <p class="font-medium text-gray-900">{{ user.name }}</p>

                  </div>
                </div>
              </td>
              <td class="px-6 py-4">
                <div class="space-y-1">
                  <div class="flex items-center gap-2 text-sm text-gray-600">
                    <Mail class="w-4 h-4" />
                    <span>{{ user.email }}</span>
                  </div>
                  <div v-if="user.phone_number" class="flex items-center gap-2 text-sm text-gray-600">
                    <Phone class="w-4 h-4" />
                    <span>{{ user.phone_number }}</span>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4">
                <div v-if="getUserInvestmentItems(user).length" class="space-y-3">
                  <div v-for="investment in getUserInvestmentItems(user)" :key="investment.id" class="rounded-2xl border border-gray-200 bg-gray-50 p-4">
                    <div class="flex items-start justify-between gap-3">
                      <div>
                        <p class="font-semibold text-gray-900">
                          {{ investment.package?.package_name || `Paket #${investment.package_id}` }}
                        </p>
                        <p class="mt-1 text-sm text-gray-500">
                          {{ investment.package?.description || 'Tidak ada deskripsi' }}
                        </p>
                      </div>
                      <span class="inline-flex items-center rounded-full bg-slate-100 px-3 py-1 text-xs font-medium uppercase tracking-wide text-slate-700">
                        {{ investorStatusLabel(investment.status) }}
                      </span>
                    </div>
                  </div>
                </div>
                <p v-else class="text-sm text-gray-500">Belum memilih paket</p>
              </td>
              <td class="px-6 py-4">
                <span class="inline-flex px-2.5 py-1 rounded-full text-xs font-medium" :class="statusBadgeClass(user.package_status)">
                  {{ investorStatusLabel(user.package_status) }}
                </span>
              </td>
              <td class="px-6 py-4 text-gray-600 text-sm">
                {{ user.created_at ? new Date(user.created_at).toLocaleDateString('id-ID') : '-' }}
              </td>
              <td class="px-6 py-4">
                <div class="space-y-2">
                  <div class="flex items-center">
                  <button
                    v-if="getActionState(user) === 'validate-mitra'"
                    type="button"
                    class="inline-flex min-w-[210px] items-center justify-center gap-2 rounded-lg bg-amber-600 px-4 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-amber-700 disabled:cursor-not-allowed disabled:opacity-60"
                    :disabled="validatingMitraId === user.id"
                    @click="handleValidateMitra(user.id)"
                  >
                    <Loader2 v-if="validatingMitraId === user.id" class="h-4 w-4 animate-spin" />
                    <span>{{ validatingMitraId === user.id ? 'Memvalidasi...' : 'Validasi Mitra' }}</span>
                  </button>

                  <button
                    v-if="getActionState(user) === 'mark-paid'"
                    type="button"
                    class="inline-flex min-w-[210px] items-center justify-center gap-2 rounded-lg bg-blue-600 px-4 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-60"
                    :disabled="updatingUserId === user.id"
                    @click="handleMarkAsPaid(user.id)"
                  >
                    <Loader2 v-if="updatingUserId === user.id" class="h-4 w-4 animate-spin" />
                    <span>{{ updatingUserId === user.id ? 'Memproses...' : 'Tandai Telah Bayar' }}</span>
                  </button>

                  <button
                    v-else-if="getActionState(user) === 'activate'"
                    type="button"
                    class="inline-flex min-w-[210px] items-center justify-center gap-2 rounded-lg bg-green-600 px-4 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-green-700 disabled:cursor-not-allowed disabled:opacity-60"
                    :disabled="activatingUserId === user.id"
                    @click="openActivationForm(user)"
                  >
                    <Loader2 v-if="activatingUserId === user.id" class="h-4 w-4 animate-spin" />
                    <span>{{ activatingUserId === user.id ? 'Membuka Form...' : 'Aktifkan Penanaman Awal' }}</span>
                  </button>

                  <button
                    v-else-if="getActionState(user) === 'review'"
                    type="button"
                    class="inline-flex min-w-[210px] items-center justify-center gap-2 rounded-lg border border-green-200 bg-green-50 px-4 py-2.5 text-sm font-semibold text-green-700 transition-colors hover:bg-green-100"
                    @click="openReviewPlantBatch(user)"
                  >
                    <Eye class="h-4 w-4" />
                    <span>Tinjau Plant Batch</span>
                  </button>

                  <span
                    v-else
                    class="inline-flex min-w-[210px] items-center justify-center rounded-lg border border-gray-200 bg-gray-50 px-4 py-2.5 text-sm font-medium text-gray-500"
                  >
                    Tidak ada aksi
                  </span>
                  </div>

                  <div class="flex items-center gap-2">
                    <button
                      type="button"
                      class="inline-flex items-center gap-1 rounded-md border border-gray-200 px-2.5 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-50"
                      @click="openEditUserForm(user)"
                    >
                      <Pencil class="h-3.5 w-3.5" />
                      Edit
                    </button>
                    <button
                      type="button"
                      class="inline-flex items-center gap-1 rounded-md border border-red-200 px-2.5 py-1.5 text-xs font-medium text-red-600 hover:bg-red-50"
                      @click="handleDeleteUser(user)"
                    >
                      <Trash2 class="h-3.5 w-3.5" />
                      Hapus
                    </button>
                  </div>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <Pagination 
        v-if="users.length > 0"
        :current-page="currentPage" 
        :total-pages="totalPages" 
        :total-rows="totalRows" 
        :limit="limit"
        @update:page="handlePageChange"
      />
    </div>

    <div v-if="activationFormUser" class="fixed inset-0 z-50 overflow-y-auto bg-black/40 p-4">
      <div class="mx-auto my-6 w-full max-w-lg rounded-2xl bg-white p-6 shadow-xl max-h-[88vh] overflow-y-auto">
        <h3 class="text-xl font-semibold text-gray-900">Validasi Penanaman Awal</h3>
        <p class="mt-1 text-sm text-gray-600">
          Isi data plant batch untuk aktivasi paket user <span class="font-medium text-gray-900">{{ activationFormUser.name }}</span>.
        </p>

        <div class="mt-5 space-y-4">
          <div
            v-for="(formItem, index) in activationForms"
            :key="`${formItem.packageId}-${index}`"
            class="rounded-xl border border-gray-200 p-4"
          >
            <p class="mb-3 text-sm font-semibold text-gray-800">{{ formItem.packageName }}</p>

            <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <label class="mb-1 block text-sm font-medium text-gray-700">Batch Code</label>
                <input
                  v-model="formItem.batchCode"
                  type="text"
                  class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-100"
                  placeholder="Contoh: PB-INVESTOR01"
                />
              </div>

              <div>
                <label class="mb-1 block text-sm font-medium text-gray-700">Tanggal Tanam</label>
                <input
                  v-model="formItem.plantingDate"
                  type="date"
                  class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-100"
                />
              </div>

              <div>
                <label class="mb-1 block text-sm font-medium text-gray-700">Lokasi Penanaman</label>
                <input
                  v-model="formItem.location"
                  type="text"
                  class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-100"
                  placeholder="Contoh: Blok A - Kebun Timur"
                />
              </div>

              <div>
                <label class="mb-1 block text-sm font-medium text-gray-700">Jumlah Bibit</label>
                <input
                  v-model="formItem.seedCount"
                  type="number"
                  min="1"
                  step="1"
                  class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-100"
                  placeholder="Contoh: 120"
                />
              </div>

              <div class="md:col-span-2">
                <label class="mb-1 block text-sm font-medium text-gray-700">Luas Lahan (m²)</label>
                <input
                  v-model="formItem.landArea"
                  type="number"
                  min="0"
                  step="0.01"
                  class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-100"
                  placeholder="Contoh: 250"
                />
              </div>

              <div class="md:col-span-2">
                <label class="mb-1 block text-sm font-medium text-gray-700">Foto Batch</label>
                <input
                  type="file"
                  accept="image/png,image/jpeg,image/jpg,image/webp"
                  class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm file:mr-3 file:rounded-md file:border-0 file:bg-green-50 file:px-3 file:py-1.5 file:text-green-700 hover:file:bg-green-100"
                  @change="handleBatchImageChange(formItem, $event)"
                />
                <p class="mt-1 text-xs text-gray-500">Format: JPG, PNG, WEBP</p>

                <div v-if="formItem.imagePreview" class="mt-3 overflow-hidden rounded-lg border border-gray-200">
                  <img :src="formItem.imagePreview" alt="Preview foto batch" class="h-40 w-full object-cover" />
                </div>
              </div>
            </div>
          </div>

          <p v-if="activationFormError" class="text-sm text-red-600">{{ activationFormError }}</p>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <button
            type="button"
            class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
            @click="closeActivationForm"
          >
            Batal
          </button>
          <button
            type="button"
            class="inline-flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700 disabled:opacity-60"
            :disabled="activatingUserId === activationFormUser.id"
            @click="handleActivate(activationFormUser.id)"
          >
            <Loader2 v-if="activatingUserId === activationFormUser.id" class="h-4 w-4 animate-spin" />
            <span>{{ activatingUserId === activationFormUser.id ? 'Memproses...' : 'Aktifkan Paket' }}</span>
          </button>
        </div>
      </div>
    </div>

    <div v-if="reviewingUser" class="fixed inset-0 z-50 overflow-y-auto bg-black/40 p-4">
      <div class="mx-auto my-6 w-full max-w-3xl rounded-2xl bg-white p-6 shadow-xl max-h-[88vh] overflow-y-auto">
        <div class="flex items-start justify-between gap-4">
          <div>
            <h3 class="text-xl font-semibold text-gray-900">Tinjau Data Plant Batch</h3>
            <p class="mt-1 text-sm text-gray-600">
              Data penanaman awal untuk user <span class="font-medium text-gray-900">{{ reviewingUser.name }}</span>.
            </p>
          </div>
          <button
            type="button"
            class="rounded-lg border border-gray-300 px-3 py-1.5 text-sm font-medium text-gray-700 hover:bg-gray-50"
            @click="closeReviewPlantBatch"
          >
            Tutup
          </button>
        </div>

        <div v-if="reviewLoading" class="mt-6 flex items-center justify-center gap-2 text-sm text-gray-500">
          <Loader2 class="h-4 w-4 animate-spin" />
          Memuat data plant batch...
        </div>

        <div v-else-if="!reviewPlantBatches.length" class="mt-6 rounded-xl border border-dashed border-gray-300 p-6 text-center text-sm text-gray-500">
          Belum ada data plant batch untuk user ini.
        </div>

        <div v-else class="mt-6 space-y-4">
          <article
            v-for="item in reviewPlantBatches"
            :key="item.id"
            class="rounded-xl border border-gray-200 bg-white p-4"
          >
            <div class="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
              <div>
                <p class="text-sm font-semibold text-gray-900">{{ item.package_name || `Paket #${item.package_id}` }}</p>
                <p class="mt-1 text-xs uppercase tracking-wide text-gray-500">Batch Code: {{ item.batch_code }}</p>
              </div>
              <span class="inline-flex rounded-full bg-green-100 px-2.5 py-1 text-xs font-medium text-green-700">{{ item.status }}</span>
            </div>

            <div class="mt-4 grid grid-cols-1 gap-3 text-sm text-gray-700 md:grid-cols-2">
              <p class="inline-flex items-center gap-2"><MapPin class="h-4 w-4 text-gray-500" /> {{ item.location || '-' }}</p>
              <p class="inline-flex items-center gap-2"><Sprout class="h-4 w-4 text-gray-500" /> {{ item.seed_count }} bibit</p>
              <p>Tanggal Tanam: {{ formatDisplayDate(item.planting_date) }}</p>
              <p>Luas Lahan: {{ Number(item.land_area || 0).toLocaleString('id-ID') }} m²</p>
            </div>

            <div v-if="item.photo_url" class="mt-4 overflow-hidden rounded-lg border border-gray-200 bg-gray-50">
              <img :src="`http://localhost:8000${item.photo_url}`" alt="Foto batch" class="h-48 w-full object-cover" />
            </div>
            <div v-else class="mt-4 inline-flex items-center gap-2 text-xs text-gray-500">
              <Image class="h-3.5 w-3.5" />
              Foto batch belum tersedia.
            </div>
          </article>
        </div>
      </div>
    </div>

    <div v-if="userFormOpen" class="fixed inset-0 z-50 overflow-y-auto bg-black/40 p-4">
      <div class="mx-auto my-6 w-full max-w-xl rounded-2xl bg-white p-6 shadow-xl">
        <h3 class="text-xl font-semibold text-gray-900">
          {{ userFormMode === 'create' ? 'Tambah User Baru' : 'Edit User' }}
        </h3>
        <p class="mt-1 text-sm text-gray-600">Isi data user yang akan dikelola oleh admin.</p>

        <div class="mt-5 grid grid-cols-1 gap-4 md:grid-cols-2">
          <div class="md:col-span-2">
            <label class="mb-1 block text-sm font-medium text-gray-700">Nama</label>
            <input v-model="userForm.name" type="text" class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100" />
          </div>
          <div class="md:col-span-2">
            <label class="mb-1 block text-sm font-medium text-gray-700">Email</label>
            <input v-model="userForm.email" type="email" class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">Role</label>
            <select v-model="userForm.role" class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100">
              <option value="investor">Investor</option>
              <option value="mitra">Mitra</option>
              <option value="admin">Admin</option>
            </select>
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">No. Telepon</label>
            <input v-model="userForm.phone_number" type="text" class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100" />
          </div>
          <div class="md:col-span-2">
            <label class="mb-1 block text-sm font-medium text-gray-700">Alamat</label>
            <textarea v-model="userForm.address" rows="2" class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100" />
          </div>
          <div class="md:col-span-2">
            <label class="mb-1 block text-sm font-medium text-gray-700">Password {{ userFormMode === 'edit' ? '(opsional)' : '' }}</label>
            <input v-model="userForm.password" type="password" class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100" />
          </div>
        </div>

        <p v-if="userFormError" class="mt-3 text-sm text-red-600">{{ userFormError }}</p>

        <div class="mt-6 flex justify-end gap-3">
          <button type="button" class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50" @click="closeUserForm">
            Batal
          </button>
          <button
            type="button"
            class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-60"
            :disabled="userFormLoading"
            @click="submitUserForm"
          >
            <Loader2 v-if="userFormLoading" class="h-4 w-4 animate-spin" />
            <span>{{ userFormLoading ? 'Menyimpan...' : 'Simpan' }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
