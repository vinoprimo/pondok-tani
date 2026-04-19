<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { CheckCircle2, MessageCircle, Package, ArrowRight } from 'lucide-vue-next';
import LandingLayout from '../../layouts/LandingLayout.vue';
import PaketCard from '../../components/landing-paket/PaketCard.vue';
import { getInvestmentPackages } from '../../services/investment/package';
import { saveSelectedPackage } from '../../services/user/user';

type PaketItem = {
  id: number;
  name: string;
  range: string;
  roi: string;
  duration: string;
  highlight?: boolean;
};

const router = useRouter();
const route = useRoute();
const loadingPackages = ref(false);
const savingSelection = ref(false);
const paketItems = ref<PaketItem[]>([]);
const selectedPackageIds = ref<number[]>([]);

const formatRupiah = (value: number) =>
  new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(value);

async function loadPackages() {
  loadingPackages.value = true;
  try {
    const response = await getInvestmentPackages({ status: 'active' });
    const list = Array.isArray(response.data) ? response.data : [];
    paketItems.value = list.map((item: any, index: number) => ({
      id: item.id,
      name: item.package_name,
      range: `${formatRupiah(item.price)} (min ${item.min_quantity} unit)`,
      roi: 'Estimasi placeholder 15% - 25% / tahun',
      duration: 'Placeholder 12 - 24 bulan',
      highlight: index === 0,
    }));

    const fromPackageIds = String(route.query.package_ids || '')
      .split(',')
      .map((value) => Number(value.trim()))
      .filter((value) => Number.isFinite(value) && value > 0);

    const fallbackSingleId = Number(route.query.package_id || 0);
    const preselectedIds = fromPackageIds.length
      ? fromPackageIds
      : Number.isFinite(fallbackSingleId) && fallbackSingleId > 0
        ? [fallbackSingleId]
        : [];

    const activePackageIds = new Set(paketItems.value.map((item) => item.id));
    selectedPackageIds.value = preselectedIds.filter((id) => activePackageIds.has(id));
  } catch (error) {
    console.error('Failed to load investment packages:', error);
    paketItems.value = [];
  } finally {
    loadingPackages.value = false;
  }
}

const selectedPackages = computed(() =>
  paketItems.value.filter((item) => selectedPackageIds.value.includes(item.id))
);

const selectedPackageSummary = computed(() => {
  if (!selectedPackages.value.length) {
    return 'Belum ada paket yang dipilih';
  }

  if (selectedPackages.value.length === 1) {
    return selectedPackages.value[0].name;
  }

  return `${selectedPackages.value.length} paket dipilih`;
});

function choosePackage(packageItem: PaketItem) {
  const currentIds = new Set(selectedPackageIds.value);
  if (currentIds.has(packageItem.id)) {
    currentIds.delete(packageItem.id);
  } else {
    currentIds.add(packageItem.id);
  }

  selectedPackageIds.value = Array.from(currentIds);

  const query: Record<string, string> = {};
  if (selectedPackageIds.value.length === 1) {
    query.package_id = String(selectedPackageIds.value[0]);
  }
  if (selectedPackageIds.value.length > 1) {
    query.package_ids = selectedPackageIds.value.join(',');
  }

  router.replace({
    path: '/pilih-paket',
    query,
  });
}

async function confirmSelection() {
  if (!selectedPackageIds.value.length) return;

  savingSelection.value = true;
  try {
    await saveSelectedPackage(selectedPackageIds.value);

    const packageNames = selectedPackages.value.map((item) => item.name);
    const messageContent = packageNames.length
      ? `saya tertarik dengan paket berikut:\n- ${packageNames.join('\n- ')}`
      : 'saya tertarik dengan paket investasi';

    const message = encodeURIComponent(messageContent);
    const whatsappUrl = `https://wa.me/6281328164003?text=${message}`;
    window.open(whatsappUrl, '_blank', 'noopener,noreferrer');
  } catch (error) {
    console.error('Failed to save package selection:', error);
  } finally {
    savingSelection.value = false;
  }
}

onMounted(() => {
  loadPackages();
});
</script>

<template>
  <LandingLayout>
    <section class="mx-auto max-w-7xl px-6 py-10 lg:py-14">
      <div class="mb-8 flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div class="max-w-2xl space-y-4">
          <span class="inline-flex items-center gap-2 rounded-full border border-green-200 bg-green-50 px-3 py-1 text-xs font-semibold uppercase tracking-[0.2em] text-green-700">
            <Package class="h-3.5 w-3.5" />
            Langkah berikutnya
          </span>
          <div class="space-y-3">
            <h1 class="text-3xl font-bold text-gray-900 lg:text-5xl">Pilih paket yang ingin Anda aktifkan</h1>
            <p class="text-base leading-7 text-gray-600 lg:text-lg">
              Paket yang Anda pilih dari landing page sudah dibawa ke halaman ini. Tinjau ulang detailnya, lalu lanjutkan ke WhatsApp admin untuk konfirmasi.
            </p>
          </div>
        </div>
        <div class="rounded-2xl border border-green-200 bg-green-50 px-5 py-4 text-sm text-green-800 shadow-sm">
          <p class="font-semibold">Proses singkat</p>
          <p class="mt-1 max-w-sm">Pilih paket, simpan pilihan, lalu kirim pesan konfirmasi ke admin.</p>
        </div>
      </div>

      <p v-if="loadingPackages" class="mb-4 text-sm text-gray-500">Memuat paket investasi terbaru...</p>
      <p v-else-if="paketItems.length === 0" class="mb-4 text-sm text-gray-500">Belum ada paket aktif.</p>

      <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <PaketCard
          v-for="item in paketItems"
          :key="item.id"
          :name="item.name"
          :range="item.range"
          :roi="item.roi"
          :duration="item.duration"
          :highlight="item.highlight"
          :selected="selectedPackageIds.includes(item.id)"
          button-label="Pilih / Batalkan"
          @choose="choosePackage(item)"
        />
      </div>

      <div class="mt-8 rounded-3xl border border-gray-200 bg-white p-6 shadow-sm lg:p-8">
        <div class="flex flex-col gap-6 lg:flex-row lg:items-center lg:justify-between">
          <div class="space-y-3">
            <p class="inline-flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.2em] text-green-700">
              <CheckCircle2 class="h-4 w-4" />
              Paket terpilih
            </p>
            <div>
              <h2 class="text-2xl font-bold text-gray-900">{{ selectedPackageSummary }}</h2>
              <p class="mt-1 text-sm text-gray-600">
                {{ selectedPackages.length ? 'Pilihan ini akan dikirim ke admin lewat WhatsApp setelah Anda menekan tombol konfirmasi.' : 'Klik satu atau beberapa paket di atas untuk melanjutkan.' }}
              </p>
              <ul v-if="selectedPackages.length" class="mt-3 list-disc pl-5 text-sm text-gray-700 space-y-1">
                <li v-for="item in selectedPackages" :key="item.id">{{ item.name }}</li>
              </ul>
            </div>
          </div>

          <div class="flex flex-col gap-3 lg:min-w-[260px]">
            <button
              type="button"
              class="inline-flex items-center justify-center gap-2 rounded-xl bg-green-600 px-5 py-3 font-semibold text-white transition-colors hover:bg-green-700 disabled:cursor-not-allowed disabled:bg-gray-300"
              :disabled="!selectedPackageIds.length || savingSelection"
              @click="confirmSelection"
            >
              <MessageCircle class="h-4 w-4" />
              {{ savingSelection ? 'Memproses...' : 'Konfirmasi via WhatsApp' }}
            </button>
            <button
              type="button"
              class="inline-flex items-center justify-center gap-2 rounded-xl border border-gray-200 px-5 py-3 font-semibold text-gray-700 transition-colors hover:bg-gray-50"
              @click="router.push('/paket-investasi')"
            >
              <ArrowRight class="h-4 w-4" />
              Kembali ke paket
            </button>
          </div>
        </div>
      </div>
    </section>
  </LandingLayout>
</template>
