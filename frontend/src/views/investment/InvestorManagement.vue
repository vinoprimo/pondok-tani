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
  MoreVertical,
  Package,
  Loader2,
} from 'lucide-vue-next';
import { activateUserPackage, getAdminUsers } from '../../services/user/user';

const users = ref<any[]>([]);
const loading = ref(false);
const activatingUserId = ref<string | null>(null);
const searchTerm = ref('');
const statusFilter = ref('all');

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

const totalSelected = computed(() => users.value.filter((user) => user.selected_package_id).length);
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
  if (status === 'pending') return 'bg-yellow-100 text-yellow-700';
  return 'bg-gray-100 text-gray-700';
}

function investorStatusLabel(status: string) {
  if (status === 'active') return 'Aktif';
  if (status === 'pending') return 'Menunggu';
  return 'Belum memilih';
}

async function loadUsers() {
  loading.value = true;
  try {
    const response = await getAdminUsers();
    users.value = Array.isArray(response.data) ? response.data : [];
  } catch (error) {
    console.error('Failed to load users:', error);
    users.value = [];
  } finally {
    loading.value = false;
  }
}

async function handleActivate(userId: string) {
  activatingUserId.value = userId;
  try {
    await activateUserPackage(userId);
    await loadUsers();
  } catch (error) {
    console.error('Failed to activate package:', error);
  } finally {
    activatingUserId.value = null;
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
        <h2 class="text-2xl font-semibold text-gray-900">Manajemen investor</h2>
        <p class="text-gray-600 mt-1">Kelola pilihan paket dan aktivasi akun investor</p>
      </div>
      <button
        type="button"
        class="flex items-center gap-2 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
        @click="loadUsers"
      >
        <UserCheck class="w-4 h-4" />
        Refresh data
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
          <div class="relative">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              v-model="searchTerm"
              type="text"
              placeholder="Cari nama, email, atau ID user..."
              class="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
            />
          </div>
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
                    <p class="text-sm text-gray-500">{{ user.id }}</p>
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
                <div v-if="user.selected_package" class="space-y-1">
                  <p class="font-medium text-gray-900">{{ user.selected_package.package_name }}</p>
                  <p class="text-sm text-gray-500">{{ user.selected_package.description }}</p>
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
                <div class="flex items-center gap-2">
                  <button
                    v-if="user.selected_package_id && user.package_status !== 'active'"
                    type="button"
                    class="inline-flex items-center gap-2 rounded-lg bg-green-600 px-3 py-2 text-sm font-medium text-white hover:bg-green-700 disabled:opacity-60"
                    :disabled="activatingUserId === user.id"
                    @click="handleActivate(user.id)"
                  >
                    <Loader2 v-if="activatingUserId === user.id" class="h-4 w-4 animate-spin" />
                    <span>{{ activatingUserId === user.id ? 'Mengaktifkan...' : 'Aktifkan paket' }}</span>
                  </button>
                  <button type="button" class="text-gray-400 hover:text-gray-600">
                    <MoreVertical class="w-5 h-5" />
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
