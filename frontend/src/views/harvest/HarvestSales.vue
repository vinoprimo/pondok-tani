<script setup lang="ts">
import { computed, ref } from 'vue';
import { ShoppingCart, Package, Check, X } from 'lucide-vue-next';

type Grade = { grade: string; quantity: number; price: number };

type HarvestDry = {
  id: string;
  batchId: string;
  investor: string;
  type: 'Dry';
  grades: Grade[];
  totalQuantity: number;
  province: string;
  harvestDate: string;
};

type HarvestWet = {
  id: string;
  batchId: string;
  investor: string;
  type: 'Wet';
  quantity: number;
  price: number;
  province: string;
  harvestDate: string;
};

type Harvest = HarvestDry | HarvestWet;

const selectedHarvest = ref<string | null>(null);
const showConfirmModal = ref(false);

const harvestsReadyForSale: Harvest[] = [
  {
    id: 'HRV-001',
    batchId: 'BATCH-045',
    investor: 'Sarah Johnson',
    type: 'Dry',
    grades: [
      { grade: 'A', quantity: 25, price: 5000000 },
      { grade: 'B', quantity: 15, price: 4250000 },
    ],
    totalQuantity: 40,
    province: 'West Java',
    harvestDate: '2024-01-15',
  },
  {
    id: 'HRV-002',
    batchId: 'BATCH-038',
    investor: 'Michael Chen',
    type: 'Wet',
    quantity: 120,
    price: 750000,
    province: 'West Java',
    harvestDate: '2024-01-18',
  },
  {
    id: 'HRV-003',
    batchId: 'BATCH-052',
    investor: 'Emma Davis',
    type: 'Dry',
    grades: [
      { grade: 'A', quantity: 30, price: 5000000 },
      { grade: 'B', quantity: 20, price: 4250000 },
      { grade: 'C', quantity: 10, price: 3500000 },
    ],
    totalQuantity: 60,
    province: 'West Java',
    harvestDate: '2024-01-20',
  },
];

function calculateTotalValue(harvest: Harvest) {
  if (harvest.type === 'Wet') {
    return harvest.quantity * harvest.price;
  }
  return harvest.grades.reduce(
    (sum, grade) => sum + grade.quantity * grade.price,
    0
  );
}

function handleProcessSale(harvestId: string) {
  selectedHarvest.value = harvestId;
  showConfirmModal.value = true;
}

function confirmSale() {
  window.alert('Sale processed successfully!');
  showConfirmModal.value = false;
  selectedHarvest.value = null;
}

const selectedHarvestData = computed(() =>
  harvestsReadyForSale.find((h) => h.id === selectedHarvest.value)
);

function typeBadgeClass(type: string) {
  return type === 'Dry'
    ? 'bg-orange-100 text-orange-800'
    : 'bg-blue-100 text-blue-800';
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-semibold text-gray-900">Harvest Sales</h2>
      <p class="text-gray-600 mt-1">Manage and process harvest sales</p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Ready for Sale</p>
          <Package class="w-5 h-5 text-green-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">12</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Total Quantity</p>
          <Package class="w-5 h-5 text-blue-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">420 kg</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Est. Total Value</p>
          <ShoppingCart class="w-5 h-5 text-green-600" />
        </div>
        <p class="text-2xl font-semibold text-green-600">$1.85M</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Sales This Month</p>
          <ShoppingCart class="w-5 h-5 text-orange-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">8</p>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Harvest ID</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Investor</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Type</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Grade Breakdown</th>
              <th class="text-right px-6 py-3 text-sm font-medium text-gray-900">Quantity</th>
              <th class="text-right px-6 py-3 text-sm font-medium text-gray-900">Est. Value</th>
              <th class="text-center px-6 py-3 text-sm font-medium text-gray-900">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="harvest in harvestsReadyForSale" :key="harvest.id" class="hover:bg-gray-50">
              <td class="px-6 py-4">
                <div>
                  <p class="font-medium text-gray-900">{{ harvest.id }}</p>
                  <p class="text-sm text-gray-500">{{ harvest.batchId }}</p>
                </div>
              </td>
              <td class="px-6 py-4 text-gray-900">{{ harvest.investor }}</td>
              <td class="px-6 py-4">
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                  :class="typeBadgeClass(harvest.type)"
                >
                  {{ harvest.type }}
                </span>
              </td>
              <td class="px-6 py-4">
                <div v-if="harvest.type === 'Dry'" class="space-y-1">
                  <div v-for="(g, idx) in harvest.grades" :key="idx" class="text-sm">
                    <span class="font-medium">Grade {{ g.grade }}:</span> {{ g.quantity }} kg
                  </div>
                </div>
                <span v-else class="text-sm text-gray-500">-</span>
              </td>
              <td class="px-6 py-4 text-right font-medium text-gray-900">
                {{
                  harvest.type === 'Dry' ? harvest.totalQuantity : harvest.quantity
                }}
                kg
              </td>
              <td class="px-6 py-4 text-right font-semibold text-green-600">
                ${{ calculateTotalValue(harvest).toLocaleString() }}
              </td>
              <td class="px-6 py-4 text-center">
                <button
                  type="button"
                  class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors text-sm font-medium"
                  @click="handleProcessSale(harvest.id)"
                >
                  Process Sale
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div
      v-if="showConfirmModal && selectedHarvestData"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
    >
      <div class="bg-white rounded-2xl max-w-lg w-full">
        <div class="border-b border-gray-200 px-6 py-4">
          <h3 class="text-lg font-semibold text-gray-900">Confirm Sale</h3>
          <p class="text-sm text-gray-600 mt-1">Review sale details before processing</p>
        </div>

        <div class="p-6 space-y-4">
          <div class="bg-gray-50 rounded-lg p-4 space-y-3">
            <div class="flex justify-between">
              <span class="text-sm text-gray-600">Harvest ID</span>
              <span class="font-medium text-gray-900">{{ selectedHarvestData.id }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-sm text-gray-600">Investor</span>
              <span class="font-medium text-gray-900">{{ selectedHarvestData.investor }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-sm text-gray-600">Type</span>
              <span class="font-medium text-gray-900">{{ selectedHarvestData.type }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-sm text-gray-600">Total Quantity</span>
              <span class="font-medium text-gray-900">
                {{
                  selectedHarvestData.type === 'Dry'
                    ? selectedHarvestData.totalQuantity
                    : selectedHarvestData.quantity
                }}
                kg
              </span>
            </div>
            <div class="flex justify-between pt-2 border-t border-gray-200">
              <span class="text-sm font-medium text-gray-900">Total Sale Value</span>
              <span class="text-xl font-bold text-green-600">
                ${{ calculateTotalValue(selectedHarvestData).toLocaleString() }}
              </span>
            </div>
          </div>

          <div class="flex gap-3">
            <button
              type="button"
              class="flex-1 px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors font-medium text-gray-700 inline-flex items-center justify-center gap-2"
              @click="showConfirmModal = false"
            >
              <X class="w-4 h-4" />
              Cancel
            </button>
            <button
              type="button"
              class="flex-1 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-medium inline-flex items-center justify-center gap-2"
              @click="confirmSale"
            >
              <Check class="w-4 h-4" />
              Confirm Sale
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
