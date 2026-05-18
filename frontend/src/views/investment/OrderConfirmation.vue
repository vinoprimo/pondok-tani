<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ArrowLeft, CheckCircle, MessageCircle } from 'lucide-vue-next';
import LandingLayout from '../../layouts/LandingLayout.vue';
import PaketCard from '../../components/landing-paket/PaketCard.vue';
import { getInvestmentPackages } from '../../services/investment/package';
import { getCurrentUser } from '../../services/user/user';

type PaketItem = {
  id: number;
  name: string;
  range: string;
  min_quantity: number;
  roi: string;
  duration: string;
  benefits: string[];
  is_popular: boolean;
};

const router = useRouter();
const route = useRoute();
const loadingPackages = ref(false);
const paketItems = ref<PaketItem[]>([]);
const selectedPackageIds = ref<number[]>([]);
const userProfile = ref({
  name: 'User',
  email: '-',
});

const formatRupiah = (value: number) =>
  new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(value);

function parseSelectedIds() {
  return String(route.query.package_ids || '')
    .split(',')
    .map((value) => Number(value.trim()))
    .filter((value) => Number.isFinite(value) && value > 0);
}

async function loadPackages() {
  loadingPackages.value = true;
  try {
    const response = await getInvestmentPackages({ status: 'active' });
    const list = Array.isArray(response.data) ? response.data : [];
    paketItems.value = list.map((item: any) => ({
      id: item.id,
      name: item.package_name,
      range: formatRupiah(item.price),
      min_quantity: item.min_quantity ?? 1,
      roi: item.roi || 'Estimasi 15% - 25% / tahun',
      duration: item.duration || '12 - 24 bulan',
      benefits: Array.isArray(item.benefits)
        ? item.benefits
        : item.benefits?.split('\n').filter(Boolean) || [],
      is_popular: item.is_popular ?? item.isPopular ?? false,
    }));

    const selectedIds = parseSelectedIds();
    const activePackageIds = new Set(paketItems.value.map((item) => item.id));
    selectedPackageIds.value = selectedIds.filter((id) => activePackageIds.has(id));
  } catch (error) {
    console.error('Failed to load investment packages:', error);
    paketItems.value = [];
  } finally {
    loadingPackages.value = false;
  }
}

async function loadUserProfile() {
  try {
    const response = await getCurrentUser();
    userProfile.value = {
      name: response.data?.name || 'User',
      email: response.data?.email || '-',
    };
  } catch (error) {
    console.error('Failed to load user profile:', error);
    userProfile.value = {
      name: 'User',
      email: '-',
    };
  } finally {
  }
}

const selectedPackages = computed(() =>
  paketItems.value.filter((item) => selectedPackageIds.value.includes(item.id))
);

const whatsappMessage = computed(() => {
  const userName = userProfile.value.name || 'User';
  const userEmail = userProfile.value.email || '-';

  if (!selectedPackages.value.length) {
    return `halo admin, saya ingin melanjutkan konfirmasi pesanan paket investasi.\n\nNama: ${userName}\nEmail: ${userEmail}`;
  }

  return `halo admin, saya ingin konfirmasi paket investasi berikut:\n\nNama: ${userName}\nEmail: ${userEmail}\n\nPaket dipilih:\n- ${selectedPackages.value.map((item) => item.name).join('\n- ')}`;
});

function goBackToSelection() {
  router.push({
    path: '/pilih-paket',
    query: selectedPackageIds.value.length
      ? { package_ids: selectedPackageIds.value.join(',') }
      : {},
  });
}

function continueToWhatsApp() {
  const message = encodeURIComponent(whatsappMessage.value);
  const whatsappUrl = `https://wa.me/6281328164003?text=${message}`;
  window.open(whatsappUrl, '_blank', 'noopener,noreferrer');
}

onMounted(() => {
  loadPackages();
  loadUserProfile();
});
</script>

<template>
  <LandingLayout>
    <section class="mx-auto max-w-7xl px-6 py-10 lg:py-14">
      <div class="mb-8 flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div class="max-w-2xl space-y-4">
          <span class="inline-flex items-center gap-2 rounded-full border border-green-200 bg-green-50 px-3 py-1 text-xs font-semibold uppercase tracking-[0.2em] text-green-700">
            <CheckCircle class="h-3.5 w-3.5" />
            Konfirmasi pesanan
          </span>
          <div class="space-y-3">
            <h1 class="text-3xl font-bold text-gray-900 lg:text-5xl">Tinjau paket yang Anda pilih</h1>
            <p class="text-base leading-7 text-gray-600 lg:text-lg">
              Periksa kembali detail paket sebelum melanjutkan konfirmasi ke WhatsApp admin.
            </p>
          </div>
        </div>
        <div class="rounded-2xl border border-green-200 bg-green-50 px-5 py-4 text-sm text-green-800 shadow-sm">
          <p class="font-semibold">Langkah berikutnya</p>
          <p class="mt-1 max-w-sm">Pastikan paket sudah sesuai, lalu lanjutkan pesan konfirmasi.</p>
        </div>
      </div>

      <p v-if="loadingPackages" class="mb-4 text-sm text-gray-500">Memuat detail paket...</p>

      <div
        v-if="selectedPackages.length"
        :class="[
          'grid gap-6',
          selectedPackages.length === 1
            ? 'grid-cols-1 justify-items-center'
            : 'grid-cols-1 lg:grid-cols-3'
        ]"
      >
        <div
          v-for="item in selectedPackages"
          :key="item.id"
          class="w-full max-w-sm"
        >
          <PaketCard
            :name="item.name"
            :range="item.range"
            :description="''"
            :roi="item.roi"
            :duration="item.duration"
            :benefits="item.benefits"
            :min-quantity="item.min_quantity"
            :highlight="false"
            :popular="item.is_popular"
            :selected="true"
            :disabled="true"
            button-label="Paket Dipilih"
          />
        </div>
      </div>

      <div v-else class="rounded-3xl border border-dashed border-gray-300 bg-white p-8 text-center text-gray-600">
        Belum ada paket yang dipilih. Silakan kembali ke halaman pemilihan paket.
      </div>

      <div class="mt-8 rounded-3xl border border-gray-200 bg-white p-6 shadow-sm lg:p-8">
        <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          <div>
            <h3 class="text-2xl font-bold text-gray-900">Siap lanjut konfirmasi?</h3>
            <p class="mt-1 text-sm text-gray-600">
              Anda bisa kembali memilih paket atau langsung mengirim pesan ke WhatsApp admin.
            </p>
          </div>

          <div class="flex flex-col gap-3 lg:min-w-[280px]">
            <button
              type="button"
              class="inline-flex items-center justify-center gap-2 rounded-xl border border-gray-200 px-5 py-3 font-semibold text-gray-700 transition-colors hover:bg-gray-50"
              @click="goBackToSelection"
            >
              <ArrowLeft class="h-4 w-4" />
              Kembali memilih paket
            </button>
            <button
              type="button"
              class="inline-flex items-center justify-center gap-2 rounded-xl bg-green-600 px-5 py-3 font-semibold text-white transition-colors hover:bg-green-700 disabled:cursor-not-allowed disabled:bg-gray-300"
              :disabled="!selectedPackages.length"
              @click="continueToWhatsApp"
            >
              <MessageCircle class="h-4 w-4" />
              Konfirmasi via WhatsApp
            </button>
          </div>
        </div>
      </div>
    </section>
  </LandingLayout>
</template>