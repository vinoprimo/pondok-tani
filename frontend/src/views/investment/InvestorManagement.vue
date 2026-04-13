<script setup lang="ts">
import { computed } from 'vue';
import {
  Users,
  DollarSign,
  TrendingUp,
  UserPlus,
  Search,
  Filter,
  Mail,
  Phone,
  MoreVertical,
} from 'lucide-vue-next';

const investors = [
  {
    id: 'INV-001',
    name: 'Sarah Johnson',
    email: 'sarah.j@email.com',
    phone: '+1 (555) 123-4567',
    investment: 45000,
    roi: 24.5,
    plants: 85,
    status: 'active',
    joinDate: '2024-03-15',
  },
  {
    id: 'INV-002',
    name: 'Michael Chen',
    email: 'michael.c@email.com',
    phone: '+1 (555) 234-5678',
    investment: 75000,
    roi: 26.2,
    plants: 140,
    status: 'active',
    joinDate: '2024-01-10',
  },
  {
    id: 'INV-003',
    name: 'Emma Davis',
    email: 'emma.d@email.com',
    phone: '+1 (555) 345-6789',
    investment: 32000,
    roi: 22.8,
    plants: 60,
    status: 'active',
    joinDate: '2024-06-20',
  },
  {
    id: 'INV-004',
    name: 'James Wilson',
    email: 'james.w@email.com',
    phone: '+1 (555) 456-7890',
    investment: 55000,
    roi: 25.1,
    plants: 105,
    status: 'active',
    joinDate: '2024-02-28',
  },
  {
    id: 'INV-005',
    name: 'Linda Martinez',
    email: 'linda.m@email.com',
    phone: '+1 (555) 567-8901',
    investment: 28000,
    roi: 21.5,
    plants: 52,
    status: 'pending',
    joinDate: '2026-01-15',
  },
];

const totalInvestment = computed(() =>
  investors.reduce((sum, inv) => sum + inv.investment, 0)
);

const activeInvestors = computed(() =>
  investors.filter((inv) => inv.status === 'active').length
);

const avgROI = computed(() =>
  (investors.reduce((sum, inv) => sum + inv.roi, 0) / investors.length).toFixed(1)
);

function investorInitials(name: string) {
  return name
    .split(' ')
    .map((n) => n[0])
    .join('');
}

function statusBadgeClass(status: string) {
  return status === 'active'
    ? 'bg-green-100 text-green-700'
    : 'bg-yellow-100 text-yellow-700';
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-semibold text-gray-900">Investor Management</h2>
        <p class="text-gray-600 mt-1">Manage and track all investor accounts</p>
      </div>
      <button
        type="button"
        class="flex items-center gap-2 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
      >
        <UserPlus class="w-4 h-4" />
        Add New Investor
      </button>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-12 h-12 bg-blue-50 rounded-lg flex items-center justify-center">
            <Users class="w-6 h-6 text-blue-600" />
          </div>
          <span class="text-sm text-gray-600">Total Investors</span>
        </div>
        <p class="text-3xl font-semibold text-gray-900">{{ investors.length }}</p>
        <p class="text-sm text-green-600 mt-1">{{ activeInvestors }} active</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-12 h-12 bg-green-50 rounded-lg flex items-center justify-center">
            <DollarSign class="w-6 h-6 text-green-600" />
          </div>
          <span class="text-sm text-gray-600">Total Investment</span>
        </div>
        <p class="text-3xl font-semibold text-gray-900">
          ${{ totalInvestment.toLocaleString() }}
        </p>
        <p class="text-sm text-gray-600 mt-1">Cumulative</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-12 h-12 bg-purple-50 rounded-lg flex items-center justify-center">
            <TrendingUp class="w-6 h-6 text-purple-600" />
          </div>
          <span class="text-sm text-gray-600">Average ROI</span>
        </div>
        <p class="text-3xl font-semibold text-gray-900">{{ avgROI }}%</p>
        <p class="text-sm text-gray-600 mt-1">Across all investors</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-12 h-12 bg-yellow-50 rounded-lg flex items-center justify-center">
            <UserPlus class="w-6 h-6 text-yellow-600" />
          </div>
          <span class="text-sm text-gray-600">New This Month</span>
        </div>
        <p class="text-3xl font-semibold text-gray-900">12</p>
        <p class="text-sm text-green-600 mt-1">+18% vs last month</p>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-4">
      <div class="flex flex-wrap items-center gap-4">
        <div class="flex-1 min-w-[300px]">
          <div class="relative">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              type="text"
              placeholder="Search investors by name, email, or ID..."
              class="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
            />
          </div>
        </div>
        <button
          type="button"
          class="flex items-center gap-2 px-4 py-2 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
        >
          <Filter class="w-4 h-4" />
          Filter
        </button>
        <select
          class="px-4 py-2 border border-gray-200 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-green-500"
        >
          <option>All Status</option>
          <option>Active</option>
          <option>Pending</option>
          <option>Inactive</option>
        </select>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Investor
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Contact
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Investment
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                ROI
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Plants
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Join Date
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="investor in investors" :key="investor.id" class="hover:bg-gray-50">
              <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                  <div
                    class="w-10 h-10 bg-gradient-to-br from-green-400 to-green-600 rounded-full flex items-center justify-center text-white font-semibold"
                  >
                    {{ investorInitials(investor.name) }}
                  </div>
                  <div>
                    <p class="font-medium text-gray-900">{{ investor.name }}</p>
                    <p class="text-sm text-gray-500">{{ investor.id }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4">
                <div class="space-y-1">
                  <div class="flex items-center gap-2 text-sm text-gray-600">
                    <Mail class="w-4 h-4" />
                    <span>{{ investor.email }}</span>
                  </div>
                  <div class="flex items-center gap-2 text-sm text-gray-600">
                    <Phone class="w-4 h-4" />
                    <span>{{ investor.phone }}</span>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 font-medium text-gray-900">
                ${{ investor.investment.toLocaleString() }}
              </td>
              <td class="px-6 py-4">
                <span
                  class="inline-flex items-center gap-1 px-2.5 py-1 bg-green-50 text-green-700 rounded-full text-sm font-medium"
                >
                  <TrendingUp class="w-3 h-3" />
                  {{ investor.roi }}%
                </span>
              </td>
              <td class="px-6 py-4 text-gray-900">{{ investor.plants }}</td>
              <td class="px-6 py-4">
                <span
                  class="inline-flex px-2.5 py-1 rounded-full text-xs font-medium"
                  :class="statusBadgeClass(investor.status)"
                >
                  {{ investor.status }}
                </span>
              </td>
              <td class="px-6 py-4 text-gray-600 text-sm">{{ investor.joinDate }}</td>
              <td class="px-6 py-4">
                <button type="button" class="text-gray-400 hover:text-gray-600">
                  <MoreVertical class="w-5 h-5" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="flex items-center justify-between bg-white rounded-xl border border-gray-200 p-4">
      <p class="text-sm text-gray-600">
        Showing 1 to {{ investors.length }} of 234 investors
      </p>
      <div class="flex gap-2">
        <button
          type="button"
          class="px-3 py-1 border border-gray-200 rounded-lg hover:bg-gray-50 text-sm"
        >
          Previous
        </button>
        <button type="button" class="px-3 py-1 bg-green-600 text-white rounded-lg text-sm">
          1
        </button>
        <button
          type="button"
          class="px-3 py-1 border border-gray-200 rounded-lg hover:bg-gray-50 text-sm"
        >
          2
        </button>
        <button
          type="button"
          class="px-3 py-1 border border-gray-200 rounded-lg hover:bg-gray-50 text-sm"
        >
          3
        </button>
        <button
          type="button"
          class="px-3 py-1 border border-gray-200 rounded-lg hover:bg-gray-50 text-sm"
        >
          Next
        </button>
      </div>
    </div>
  </div>
</template>
