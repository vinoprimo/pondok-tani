<script setup lang="ts">
import { ref, computed } from 'vue';
import { X, Calculator, TrendingUp } from 'lucide-vue-next';

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  close: [];
}>();

const harvestQty = ref('');
const harvestType = ref<'wet' | 'dry'>('wet');
const grade = ref('');
const province = ref('');
const shrinkage = ref('75');
const operationalCost = ref('');

const priceData: Record<string, { wet: number; dry: Record<string, number> }> = {
  'West Java': {
    wet: 750000,
    dry: { A: 5000000, B: 4250000, C: 3500000, Split: 2500000, Asalan: 2000000 },
  },
  'Central Java': {
    wet: 725000,
    dry: { A: 4800000, B: 4080000, C: 3360000, Split: 2400000, Asalan: 1920000 },
  },
  'East Java': {
    wet: 780000,
    dry: { A: 5200000, B: 4420000, C: 3640000, Split: 2600000, Asalan: 2080000 },
  },
  Bali: {
    wet: 800000,
    dry: { A: 5400000, B: 4590000, C: 3780000, Split: 2700000, Asalan: 2160000 },
  },
};

const projection = computed(() => {
  if (!harvestQty.value || !province.value) return null;

  const qty = parseFloat(harvestQty.value);
  const provinceData = priceData[province.value];
  if (!provinceData) return null;

  let revenue = 0;
  let effectiveQty = qty;

  if (harvestType.value === 'wet') {
    revenue = qty * provinceData.wet;
  } else {
    if (!grade.value) return null;
    const shrinkagePercent = parseFloat(shrinkage.value) / 100;
    effectiveQty = qty * (1 - shrinkagePercent);
    const dryPrice = provinceData.dry[grade.value];
    if (dryPrice === undefined) return null;
    revenue = effectiveQty * dryPrice;
  }

  const costs = operationalCost.value ? parseFloat(operationalCost.value) : revenue * 0.15;
  const netProfit = revenue - costs;
  const roi = (netProfit / costs) * 100;

  return {
    revenue,
    costs,
    netProfit,
    roi,
    effectiveQty,
  };
});
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
  >
    <div class="bg-white rounded-2xl max-w-3xl w-full max-h-[90vh] overflow-y-auto">
      <div class="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 bg-blue-50 rounded-lg flex items-center justify-center">
            <Calculator class="w-6 h-6 text-blue-600" />
          </div>
          <div>
            <h2 class="text-xl font-semibold text-gray-900">Kalkulator proyeksi pendapatan</h2>
            <p class="text-sm text-gray-600">Perkirakan pendapatan panen dan ROI Anda</p>
          </div>
        </div>
        <button
          type="button"
          class="w-8 h-8 rounded-lg hover:bg-gray-100 flex items-center justify-center transition-colors"
          @click="emit('close')"
        >
          <X class="w-5 h-5 text-gray-500" />
        </button>
      </div>

      <div class="p-6 space-y-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Perkiraan jumlah panen (kg)
            </label>
            <input
              v-model="harvestQty"
              type="number"
              placeholder="mis. 100"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Jenis panen</label>
            <div class="grid grid-cols-2 gap-2">
              <button
                type="button"
                class="px-4 py-2 rounded-lg font-medium transition-colors"
                :class="
                  harvestType === 'wet'
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                "
                @click="harvestType = 'wet'"
              >
                Basah
              </button>
              <button
                type="button"
                class="px-4 py-2 rounded-lg font-medium transition-colors"
                :class="
                  harvestType === 'dry'
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                "
                @click="harvestType = 'dry'"
              >
                Kering
              </button>
            </div>
          </div>

          <div v-if="harvestType === 'dry'">
            <label class="block text-sm font-medium text-gray-700 mb-2">Pilihan mutu</label>
            <select
              v-model="grade"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="">Pilih mutu</option>
              <option value="A">Mutu A — Premium</option>
              <option value="B">Mutu B — Standar</option>
              <option value="C">Mutu C — Ekonomi</option>
              <option value="Split">Split</option>
              <option value="Asalan">Asalan</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Provinsi (acuan harga)</label>
            <select
              v-model="province"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="">Pilih provinsi</option>
              <option value="West Java">Jawa Barat</option>
              <option value="Central Java">Jawa Tengah</option>
              <option value="East Java">Jawa Timur</option>
              <option value="Bali">Bali</option>
            </select>
          </div>

          <div v-if="harvestType === 'dry'">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Perkiraan susut pengeringan (%)
            </label>
            <input
              v-model="shrinkage"
              type="number"
              placeholder="75"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Perkiraan biaya operasional (opsional)
            </label>
            <input
              v-model="operationalCost"
              type="number"
              placeholder="Otomatis: 15% dari pendapatan"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
        </div>

        <div
          v-if="projection"
          class="bg-gradient-to-br from-green-50 to-blue-50 rounded-xl p-6 border-2 border-green-200"
        >
          <div class="flex items-center gap-2 mb-4">
            <TrendingUp class="w-5 h-5 text-green-600" />
            <h3 class="text-lg font-semibold text-gray-900">Hasil proyeksi</h3>
          </div>

          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div class="bg-white rounded-lg p-4">
              <p class="text-sm text-gray-600 mb-1">Perkiraan pendapatan</p>
              <p class="text-2xl font-bold text-green-600">
                ${{ projection.revenue.toLocaleString() }}
              </p>
            </div>

            <div class="bg-white rounded-lg p-4">
              <p class="text-sm text-gray-600 mb-1">Biaya operasional</p>
              <p class="text-2xl font-bold text-orange-600">
                ${{ projection.costs.toLocaleString() }}
              </p>
            </div>

            <div class="bg-white rounded-lg p-4">
              <p class="text-sm text-gray-600 mb-1">Laba bersih</p>
              <p class="text-2xl font-bold text-green-600">
                ${{ projection.netProfit.toLocaleString() }}
              </p>
            </div>

            <div class="bg-white rounded-lg p-4">
              <p class="text-sm text-gray-600 mb-1">Proyeksi ROI</p>
              <p class="text-2xl font-bold text-blue-600">{{ projection.roi.toFixed(1) }}%</p>
            </div>
          </div>

          <div v-if="harvestType === 'dry'" class="mt-4 p-3 bg-white/50 rounded-lg border border-green-200">
            <p class="text-sm text-gray-700">
              <span class="font-medium">Catatan:</span> Setelah susut {{ shrinkage }}%, jumlah efektif:
              <span class="font-semibold">{{ projection.effectiveQty.toFixed(2) }} kg</span>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
