<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { 
  CalendarDays, Clock, Wallet, CheckCircle, TrendingUp, DollarSign, 
  Shield, BarChart3, Award, ChevronRight, FileText, Users 
} from 'lucide-vue-next';
import LandingLayout from '../../layouts/LandingLayout.vue';
import { getInvestmentPackages } from '../../services/investment/package';
import { getCurrentUser } from '../../services/user/user';

type PaketItem = {
  id: number;
  name: string;
  range: string;
  min_quantity: number;
  description: string;
  roi: string;
  duration: string;
  benefits: string[];
  is_popular: boolean;
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
    paketItems.value = list.map((item: any) => ({
      id: item.id,
      name: item.package_name,
      range: `${formatRupiah(item.price)}`,
      min_quantity: item.min_quantity,
      description: item.description || 'Deskripsi paket belum tersedia.',
      roi: item.roi || 'Estimasi 15% - 25%',
      duration: item.duration || '12 - 24 bulan',
      benefits: Array.isArray(item.benefits) ? item.benefits : [],
      is_popular: item.is_popular ?? item.isPopular ?? false,
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

function selectPackage(id: number) {
  selectedPackageId.value = id;
}

function choosePackage(id?: number) {
  const targetId = id || (paketItems.value.length > 0 ? paketItems.value[0].id : null);
  if (!targetId) return;
  
  selectedPackageId.value = targetId;
  if (localStorage.getItem('token')) {
    router.push({
      path: '/pilih-paket',
      query: { package_id: String(targetId) },
    });
    return;
  }

  router.push({
    path: '/login',
    query: {
      mode: 'register',
      package_id: String(targetId),
    },
  });
}

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
    <!-- Hero Section -->
    <section class="bg-gradient-to-r from-green-600 to-green-600 text-white py-20 relative overflow-hidden">
      <div class="absolute inset-0 opacity-10 bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] from-white via-transparent to-transparent"></div>
      
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        <div class="max-w-3xl">
          <h1 class="text-4xl md:text-5xl font-bold mb-6 leading-tight drop-shadow-sm">
            Temukan Paket Investasi Sesuai Target Anda
          </h1>
          <p class="text-xl text-green-100 leading-relaxed">
            Berikut simulasi paket investasi untuk memudahkan skema perencanaan di
            Pondok Tani Land. Pilih paket yang sesuai dengan kemampuan dan target investasi Anda.
          </p>
        </div>
      </div>
    </section>

    <!-- Investment Packages -->
    <section class="py-20 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        <!-- On Process state from the old view -->
        <div v-if="hasOnProcessInvestment" class="mb-16 rounded-3xl border border-yellow-200 bg-gradient-to-br from-yellow-50 via-amber-50 to-orange-50 p-6 lg:p-8">
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

        <div v-else class="grid md:grid-cols-3 gap-8 mb-16 items-stretch">
          <div 
            v-for="item in paketItems" 
            :key="item.id"
            @click="selectPackage(item.id)"
            :class="[
              'rounded-2xl p-8 flex flex-col h-full transform transition-all duration-300 cursor-pointer',
              selectedPackageId === item.id
                ? 'bg-gradient-to-br from-green-600 to-green-600 text-white relative overflow-hidden md:scale-105 shadow-2xl z-10 border-4 border-green-400/50 hover:scale-[1.07] hover:shadow-2xl' 
                : 'bg-white border-2 border-gray-200 hover:border-green-500 hover:shadow-xl hover:-translate-y-2'
            ]"
          >
            <!-- Badge for popular -->
            <div v-if="item.is_popular" class="absolute top-0 right-0 bg-yellow-400 text-gray-900 px-4 py-1.5 text-sm font-bold rounded-bl-lg shadow-sm">
              PALING POPULER
            </div>

            <div class="flex items-center justify-between mb-6" :class="{ 'mt-2': item.is_popular }">
              <h3 class="text-2xl font-bold" :class="selectedPackageId === item.id ? 'text-white' : 'text-gray-900'">{{ item.name }}</h3>
              <div class="w-12 h-12 rounded-full flex items-center justify-center bg-green-100">
                <CheckCircle class="w-6 h-6 text-green-600" />
              </div>
            </div>

            <div class="mb-6">
              <div class="text-sm mb-2" :class="selectedPackageId === item.id ? 'text-white' : 'text-gray-600'">Mulai dari</div>
              <div class="text-4xl font-bold mb-1" :class="selectedPackageId === item.id ? 'text-white' : 'text-gray-900'">{{ item.range }}</div>
              <div class="text-sm" :class="selectedPackageId === item.id ? 'text-white' : 'text-gray-600'">per {{ item.min_quantity }} unit</div>
            </div>

            <div :class="[
              'p-4 rounded-lg mb-6 space-y-3',
              selectedPackageId === item.id ? 'bg-white/10 backdrop-blur-sm border border-white/20' : 'bg-green-50'
            ]">
              <div class="flex items-center gap-2" :class="selectedPackageId === item.id ? 'text-white' : 'text-green-600'">
                <Clock class="w-5 h-5" :class="selectedPackageId === item.id ? 'text-white' : 'text-green-600'" />
                <span class="font-semibold">Durasi: {{ item.duration }}</span>
              </div>
              <div class="flex items-center gap-2" :class="selectedPackageId === item.id ? 'text-white' : 'text-green-600'">
                <TrendingUp class="w-5 h-5" :class="selectedPackageId === item.id ? 'text-white' : 'text-green-600'" />
                <span class="font-semibold">Estimasi ROI: {{ item.roi }}</span>
              </div>
              <div class="flex items-center gap-2" :class="selectedPackageId === item.id ? 'text-white' : 'text-green-600'">
                <BarChart3 class="w-5 h-5" :class="selectedPackageId === item.id ? 'text-white' : 'text-green-600'" />
                <span class="font-semibold">Min. Investasi: {{ item.min_quantity }} unit</span>
              </div>
            </div>

            <div class="border-t pt-6 mb-6 flex-grow" :class="selectedPackageId === item.id ? 'border-white/20' : 'border-gray-200'">
              <h4 class="font-bold mb-4" :class="selectedPackageId === item.id ? 'text-white' : 'text-gray-900'">Benefit yang Didapatkan:</h4>
              <ul class="space-y-3">
                <li v-for="(benefit, i) in item.benefits" :key="i" class="flex items-start gap-3">
                  <CheckCircle class="w-5 h-5 flex-shrink-0 mt-0.5" :class="selectedPackageId === item.id ? 'text-white' : 'text-green-600'" />
                  <span :class="selectedPackageId === item.id ? 'text-white' : 'text-gray-700'">{{ benefit }}</span>
                </li>
              </ul>
            </div>

            <button 
              @click.stop="choosePackage(item.id)"
              :class="[
                'w-full py-3.5 rounded-lg font-bold transition-colors mt-auto',
                selectedPackageId === item.id
                  ? 'bg-white text-green-600 hover:bg-green-50 shadow-lg' 
                  : 'bg-green-600 text-white hover:bg-green-700'
              ]"
            >
              Pilih Paket Ini
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- How It Works -->
    <section class="py-20 bg-slate-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-16">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
            Cara Kerja Investasi
          </h2>
          <p class="text-xl text-gray-600">
            Proses investasi yang mudah, terstruktur, dan transparan
          </p>
        </div>
        <div class="grid md:grid-cols-4 gap-8">
          <div class="bg-white p-8 rounded-xl shadow-sm text-center border border-slate-100 hover:-translate-y-1 transition-transform">
            <div class="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6">
              <FileText class="w-8 h-8 text-green-600" />
            </div>
            <div class="w-10 h-10 bg-green-600 text-white rounded-full flex items-center justify-center mx-auto mb-4 font-bold text-lg shadow-md">
              1
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">Pilih Paket</h3>
            <p class="text-gray-600 leading-relaxed">
              Pilih paket investasi yang sesuai dengan kemampuan dan target Anda
            </p>
          </div>
          <div class="bg-white p-8 rounded-xl shadow-sm text-center border border-slate-100 hover:-translate-y-1 transition-transform">
            <div class="w-16 h-16 bg-emerald-100 rounded-full flex items-center justify-center mx-auto mb-6">
              <DollarSign class="w-8 h-8 text-emerald-600" />
            </div>
            <div class="w-10 h-10 bg-emerald-600 text-white rounded-full flex items-center justify-center mx-auto mb-4 font-bold text-lg shadow-md">
              2
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">Lakukan Pembayaran</h3>
            <p class="text-gray-600 leading-relaxed">
              Transfer ke rekening resmi dan upload bukti pembayaran
            </p>
          </div>
          <div class="bg-white p-8 rounded-xl shadow-sm text-center border border-slate-100 hover:-translate-y-1 transition-transform">
            <div class="w-16 h-16 bg-teal-100 rounded-full flex items-center justify-center mx-auto mb-6">
              <Award class="w-8 h-8 text-teal-600" />
            </div>
            <div class="w-10 h-10 bg-teal-600 text-white rounded-full flex items-center justify-center mx-auto mb-4 font-bold text-lg shadow-md">
              3
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">Verifikasi & Aktivasi</h3>
            <p class="text-gray-600 leading-relaxed">
              Tim kami akan memverifikasi dan mengaktifkan akun investasi Anda
            </p>
          </div>
          <div class="bg-white p-8 rounded-xl shadow-sm text-center border border-slate-100 hover:-translate-y-1 transition-transform">
            <div class="w-16 h-16 bg-green-600/10 rounded-full flex items-center justify-center mx-auto mb-6">
              <BarChart3 class="w-8 h-8 text-white" />
            </div>
            <div class="w-10 h-10 bg-green-600 text-white rounded-full flex items-center justify-center mx-auto mb-4 font-bold text-lg shadow-md">
              4
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">Monitor Progres</h3>
            <p class="text-gray-600 leading-relaxed">
              Pantau perkembangan investasi melalui dashboard real-time
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- Benefits -->
    <section class="py-20 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-16">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
            Keuntungan Berinvestasi
          </h2>
          <p class="text-xl text-gray-600">
            Berbagai benefit strategis yang akan Anda dapatkan
          </p>
        </div>
        <div class="grid md:grid-cols-3 gap-8">
          <div class="bg-gradient-to-br from-green-50 to-emerald-100 p-8 rounded-xl border border-green-100 shadow-sm">
            <div class="w-12 h-12 bg-green-600 rounded-lg flex items-center justify-center mb-6 shadow-md">
              <TrendingUp class="w-6 h-6 text-white" />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">ROI Kompetitif</h3>
            <p class="text-gray-700 leading-relaxed">
              Estimasi return 30-50% dari nilai investasi awal dalam periode 24-36 bulan dengan perhitungan yang sangat transparan.
            </p>
          </div>
          <div class="bg-gradient-to-br from-emerald-50 to-teal-100 p-8 rounded-xl border border-emerald-100 shadow-sm">
            <div class="w-12 h-12 bg-emerald-600 rounded-lg flex items-center justify-center mb-6 shadow-md">
              <Shield class="w-6 h-6 text-white" />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">Investasi Aman</h3>
            <p class="text-gray-700 leading-relaxed">
              Dilindungi dengan perjanjian legal yang jelas dan sistem monitoring ketat untuk menjamin keamanan dana investasi Anda.
            </p>
          </div>
          <div class="bg-gradient-to-br from-teal-50 to-green-100 p-8 rounded-xl border border-teal-100 shadow-sm">
            <div class="w-12 h-12 bg-green-600 rounded-lg flex items-center justify-center mb-6 shadow-md">
              <Users class="w-6 h-6 text-white" />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">Pendampingan Profesional</h3>
            <p class="text-gray-700 leading-relaxed">
              Tim ahli agrikultur siap memberikan update dan konsultasi lapangan secara reguler selama periode investasi berjalan.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- FAQ Preview -->
    <section class="py-20 bg-slate-50">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <h2 class="text-3xl font-bold text-gray-900 mb-12 text-center">
          Pertanyaan Umum
        </h2>
        <div class="space-y-6">
          <div class="bg-white p-6 rounded-xl shadow-sm border border-slate-100">
            <h3 class="text-lg font-bold text-gray-900 mb-2">
              Berapa minimal investasi yang harus saya keluarkan?
            </h3>
            <p class="text-gray-600 leading-relaxed">
              Minimal investasi dimulai dari Rp 10.000.000 untuk Paket C. Anda dapat memilih paket sesuai kemampuan finansial.
            </p>
          </div>
          <div class="bg-white p-6 rounded-xl shadow-sm border border-slate-100">
            <h3 class="text-lg font-bold text-gray-900 mb-2">
              Berapa lama periode investasi?
            </h3>
            <p class="text-gray-600 leading-relaxed">
              Periode investasi berkisar antara 24-36 bulan tergantung siklus budidaya vanili dan kondisi tren pasar.
            </p>
          </div>
          <div class="bg-white p-6 rounded-xl shadow-sm border border-slate-100">
            <h3 class="text-lg font-bold text-gray-900 mb-2">
              Bagaimana saya memantau perkembangan investasi?
            </h3>
            <p class="text-gray-600 leading-relaxed">
              Anda akan mendapat akses ke dashboard online yang menampilkan progres real-time, laporan bulanan/mingguan, dan dokumentasi foto.
            </p>
          </div>
          <div class="bg-white p-6 rounded-xl shadow-sm border border-slate-100">
            <h3 class="text-lg font-bold text-gray-900 mb-2">
              Apakah investasi ini dijamin aman?
            </h3>
            <p class="text-gray-600 leading-relaxed">
              Ya, investasi dilindungi dengan perjanjian legal yang jelas, monitoring ketat, dan pengelolaan oleh tim profesional berpengalaman.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="py-24 bg-green-600 text-white relative overflow-hidden">
      <div class="absolute inset-0 opacity-20 bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] from-white via-transparent to-transparent"></div>
      
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center relative z-10">
        <h2 class="text-3xl md:text-4xl font-bold mb-6 drop-shadow-sm">
          Siap Memulai Investasi Vanili Anda?
        </h2>
        <p class="text-xl text-green-100 mb-10 leading-relaxed max-w-3xl mx-auto">
          Bergabunglah dengan ratusan investor yang telah mempercayakan investasi mereka
          kepada Pondok Tani Land. Mulai wujudkan pertumbuhan finansial berkelanjutan.
        </p>
        <div class="flex flex-wrap justify-center gap-4">
          <button 
            @click="choosePackage(paketItems[1]?.id)"
            class="bg-white text-green-600 px-10 py-4 rounded-xl font-bold text-lg hover:bg-green-50 transition-all flex items-center gap-2 shadow-xl hover:-translate-y-0.5"
          >
            Pilih Paket & Investasi Sekarang
            <ChevronRight class="w-5 h-5" />
          </button>
          <button class="border-2 border-green-400/50 text-white px-10 py-4 rounded-xl font-bold text-lg hover:bg-white/10 hover:border-white transition-all">
            Konsultasi Gratis
          </button>
        </div>
      </div>
    </section>
  </LandingLayout>
</template>
