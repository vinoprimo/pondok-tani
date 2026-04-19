<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { CalendarDays, Clock, Wallet } from 'lucide-vue-next';
import LandingLayout from '../../layouts/LandingLayout.vue';
import PaketHero from '../../components/landing-paket/PaketHero.vue';
import PaketCard from '../../components/landing-paket/PaketCard.vue';
import PaketBenefitList from '../../components/landing-paket/PaketBenefitList.vue';
import { getInvestmentPackages } from '../../services/investment/package';
import { getCurrentUser } from '../../services/user/user';

type PaketItem = {
  id: number;
  name: string;
  range: string;
  roi: string;
  duration: string;
  highlight?: boolean;
};

const paketItems = ref<PaketItem[]>([]);
const loadingPackages = ref(false);
const loadingUser = ref(false);
const router = useRouter();
const selectedPackageId = ref<number | null>(null);
const userInvestments = ref<any[]>([]);

const formatRupiah = (value: number) =>
  new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(value);

const formatDate = (value: string | Date) => {
  if (!value) return '-';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '-';
  return new Intl.DateTimeFormat('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  }).format(date);
};

async function loadPackages() {
  loadingPackages.value = true;
  try {
    const res = await getInvestmentPackages({ status: 'active' });
    const list = Array.isArray(res.data) ? res.data : [];
    paketItems.value = list.map((item: any, idx: number) => ({
      id: item.id,
      name: item.package_name,
      range: `${formatRupiah(item.price)} (min ${item.min_quantity} unit)`,
      roi: 'Estimasi placeholder 15% - 25% / tahun',
      duration: 'Placeholder 12 - 24 bulan',
      highlight: idx === 1,
    }));
  } catch (err) {
    console.error('Failed to load investment packages:', err);
    paketItems.value = [];
  } finally {
    loadingPackages.value = false;
  }
}

async function loadUserStatus() {
  if (!localStorage.getItem('token')) {
    userInvestments.value = [];
    return;
  }

  loadingUser.value = true;
  try {
    const response = await getCurrentUser();
    userInvestments.value = Array.isArray(response.data?.investments) ? response.data.investments : [];
  } catch (error) {
    console.error('Failed to load user status:', error);
    userInvestments.value = [];
  } finally {
    loadingUser.value = false;
  }
}

function choosePackage(packageItem: PaketItem) {
  selectedPackageId.value = packageItem.id;
  if (localStorage.getItem('token')) {
    router.push({
      path: '/pilih-paket',
      query: { package_id: String(packageItem.id) },
    });
    return;
  }

  router.push({
    path: '/login',
    query: {
      mode: 'register',
      package_id: String(packageItem.id),
    },
  });
}

const benefits = [
  'Akses dashboard perkembangan kebun secara real-time.',
  'Laporan bulanan performa kebun dan estimasi hasil panen.',
  'Notifikasi milestone penting pada fase budidaya.',
  'Pendampingan tim agronomi selama periode investasi.',
  'Dokumentasi aktivitas lapangan dengan foto dan catatan.',
  'Ringkasan finansial dan proyeksi ROI berkala.',
];

const onProcessPackages = computed(() =>
  userInvestments.value.filter((investment) => investment.status === 'on_process')
);

const hasOnProcessInvestment = computed(() => onProcessPackages.value.length > 0);

onMounted(() => {
  loadPackages();
  loadUserStatus();
});
</script>

<template>
  <LandingLayout>

    <PaketHero
      title="Temukan Paket Investasi Sesuai Target Anda"
      description="Berikut simulasi paket investasi placeholder untuk memudahkan eksplorasi skema pendanaan di Omah Vanili. Konten ini nantinya akan diganti dengan data resmi."
    />

    <section class="max-w-7xl mx-auto px-6 pb-10 lg:pb-14">
      <p v-if="loadingPackages" class="mb-4 text-sm text-gray-500">Memuat paket investasi terbaru...</p>
      <p v-else-if="paketItems.length === 0" class="mb-4 text-sm text-gray-500">
        Belum ada paket aktif. Konten akan diperbarui setelah data tersedia.
      </p>

      <div v-if="loadingUser" class="mb-4 text-sm text-gray-500">
        Memuat status investasi Anda...
      </div>

      <div v-if="hasOnProcessInvestment" class="rounded-3xl border border-yellow-200 bg-gradient-to-br from-yellow-50 via-amber-50 to-orange-50 p-6 lg:p-8">
        <div class="flex flex-col items-center gap-6 text-center">
          <div class="flex max-w-3xl flex-col items-center gap-3">
            <div class="inline-flex h-11 w-11 items-center justify-center rounded-full bg-yellow-100">
              <Clock class="h-6 w-6 text-yellow-700" />
            </div>
            <h3 class="text-lg font-bold text-gray-900">Status paket: On Process</h3>
            <p class="text-sm text-gray-700">
              Pembayaran paket investasi Anda sudah diterima. Saat ini kami menunggu tahap selanjutnya yaitu penanaman awal.
            </p>
          </div>

          <div class="w-full">
            <p class="mb-3 text-sm font-semibold text-gray-800">Paket yang sudah terbayar</p>
            <div class="grid w-full grid-cols-1 gap-3 md:grid-cols-2">
              <article
                v-for="investment in onProcessPackages"
                :key="investment.id"
                class="h-full w-full rounded-2xl border border-yellow-200/80 bg-white/90 p-4 shadow-sm"
              >
                <div class="flex flex-col items-center gap-2 text-center">
                  <h4 class="text-sm font-semibold text-gray-900">
                    {{ investment.package?.package_name || `Paket #${investment.package_id}` }}
                  </h4>
                  <span class="rounded-full bg-yellow-100 px-2.5 py-1 text-xs font-medium text-yellow-800">
                    On Process
                  </span>
                </div>

                <div class="mt-4 space-y-2">
                  <p class="flex items-center justify-center gap-2 text-sm text-gray-700">
                    <Wallet class="h-4 w-4 text-yellow-700" />
                    <span>{{ formatRupiah(Number(investment.amount || 0)) }}</span>
                  </p>
                  <p class="flex items-center justify-center gap-2 text-sm text-gray-700">
                    <CalendarDays class="h-4 w-4 text-yellow-700" />
                    <span>{{ formatDate(investment.investment_date) }}</span>
                  </p>
                </div>
              </article>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <PaketCard
          v-for="item in paketItems"
          :key="item.id"
          :name="item.name"
          :range="item.range"
          :roi="item.roi"
          :duration="item.duration"
          :highlight="item.highlight"
          :selected="selectedPackageId === item.id"
          button-label="Pilih Paket Ini"
          @choose="choosePackage(item)"
        />
      </div>
    </section>

    <PaketBenefitList :items="benefits" />
  </LandingLayout>
</template>
