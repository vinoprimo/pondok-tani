<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { ArrowRight, Sprout, TrendingUp, Shield, BarChart3, Users, Handshake, BookOpenText, CalendarDays, Network, Globe } from 'lucide-vue-next';
import { useRouter } from "vue-router";
import LandingLayout from '../../layouts/LandingLayout.vue';
import heroImage from '../../assets/pondok-tani-hero.jpg';
import { getPublishedArticles } from '../../services/content/articles';

const router = useRouter();
const articles = ref<any[]>([]);

async function loadArticles() {
  try {
    const res = await getPublishedArticles({ status: 'publish' });
    articles.value = Array.isArray(res.data) ? res.data.slice(0, 3) : [];
  } catch (error) {
    console.error('Failed to load articles:', error);
    articles.value = [];
  }
}

onMounted(() => {
  loadArticles();
});
</script>

<template>
  <LandingLayout>

    <!-- Hero Section -->
    <section
      class="relative overflow-hidden bg-cover bg-center py-20"
      :style="{ backgroundImage: `url(${heroImage})` }"
    >
      <div class="absolute inset-0 bg-black/70"></div>
      <div class="relative max-w-7xl mx-auto px-6">
        <div class="flex flex-col items-center justify-center text-center">
          <div class="rounded-[2rem] p-10 shadow-2xl max-w-4xl">
            <h2 class="text-5xl font-bold mb-6 leading-tight text-white">
              Mulai Investasi Vanili sebagai
              <span class="text-green-600">  Aset Emas Hijau</span> Anda di Masa Depan
            </h2>
            <p class="text-xl mb-8 leading-relaxed text-white/85">
              Pondok Tani membantu Anda berinvestasi pada vanili sebagai aset emas hijau melalui pengelolaan yang profesional dan terstruktur, sehingga memberikan potensi return jangka panjang yang stabil dan berkelanjutan.
            </p>
            <div class="flex flex-col gap-4 sm:flex-row justify-center">
              <button
                type="button"
                class="px-8 py-4 bg-green-600 text-white rounded-xl hover:bg-green-600 transition-colors font-semibold flex items-center justify-center gap-2 text-lg"
                @click="router.push('/login')"
              >
                Mulai Investasi
                <ArrowRight class="w-5 h-5" />
              </button>
              <button
                type="button"
                class="px-8 py-4 border-2 border-white/90 text-white rounded-xl hover:bg-white/15 transition-colors font-semibold text-lg"
                @click="router.push('/paket-investasi')"
              >
                Pelajari Lebih Lanjut
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Stats Section -->
    <section class="bg-green-50/40 py-16 border-b border-green-100">
      <div class="max-w-7xl mx-auto px-6 overflow-hidden relative">
        <div class="w-full relative z-10 flex items-center">
          <!-- Marquee Track -->
          <div class="flex animate-marquee whitespace-nowrap">
            <!-- Marquee Content -->
            <div class="flex shrink-0 gap-20 px-10 items-center justify-around min-w-max">
              
              <div class="flex flex-col items-center min-w-[200px]">
                <div class="w-16 h-16 bg-green-600 rounded-full flex items-center justify-center mb-5 shadow-lg shadow-green-600/30 text-white transition-transform hover:scale-110 duration-300">
                  <Handshake class="w-8 h-8" />
                </div>
                <h4 class="text-4xl font-extrabold text-slate-900 mb-2">200+</h4>
                <p class="text-green-600 font-semibold text-sm uppercase tracking-wide">Petani Terkoneksi</p>
              </div>

              <div class="flex flex-col items-center min-w-[200px]">
                <div class="w-16 h-16 bg-green-600 rounded-full flex items-center justify-center mb-5 shadow-lg shadow-green-600/30 text-white transition-transform hover:scale-110 duration-300">
                  <CalendarDays class="w-8 h-8" />
                </div>
                <h4 class="text-4xl font-extrabold text-slate-900 mb-2">10+</h4>
                <p class="text-green-600 font-semibold text-sm uppercase tracking-wide">Negara Tujuan Ekspor</p>
              </div>

              <div class="flex flex-col items-center min-w-[200px]">
                <div class="w-16 h-16 bg-green-600 rounded-full flex items-center justify-center mb-5 shadow-lg shadow-green-600/30 text-white transition-transform hover:scale-110 duration-300">
                  <Network class="w-8 h-8" />
                </div>
                <h4 class="text-4xl font-extrabold text-slate-900 mb-2">20+</h4>
                <p class="text-green-600 font-semibold text-sm uppercase tracking-wide">Tonase Produksi Tahunan</p>
              </div>

              <div class="flex flex-col items-center min-w-[200px]">
                <div class="w-16 h-16 bg-green-600 rounded-full flex items-center justify-center mb-5 shadow-lg shadow-green-600/30 text-white transition-transform hover:scale-110 duration-300">
                  <Users class="w-8 h-8" />
                </div>
                <h4 class="text-4xl font-extrabold text-slate-900 mb-2">50+</h4>
                <p class="text-green-600 font-semibold text-sm uppercase tracking-wide">Anggota Komunitas Aktif</p>
              </div>

              <div class="flex flex-col items-center min-w-[200px]">
                <div class="w-16 h-16 bg-green-600 rounded-full flex items-center justify-center mb-5 shadow-lg shadow-green-600/30 text-white transition-transform hover:scale-110 duration-300">
                  <Globe class="w-8 h-8" />
                </div>
                <h4 class="text-4xl font-extrabold text-slate-900 mb-2">6+</h4>
                <p class="text-green-600 font-semibold text-sm uppercase tracking-wide">Tahun Pengalaman Kolektif</p>
              </div>

            </div>

            <!-- Duplicate Marquee Content for seamless loop -->
            <div class="flex shrink-0 gap-20 px-10 items-center justify-around min-w-max" aria-hidden="true">
              
              <div class="flex flex-col items-center min-w-[200px]">
                <div class="w-16 h-16 bg-green-600 rounded-full flex items-center justify-center mb-5 shadow-lg shadow-green-600/30 text-white transition-transform hover:scale-110 duration-300">
                  <Handshake class="w-8 h-8" />
                </div>
                <h4 class="text-4xl font-extrabold text-slate-900 mb-2">200+</h4>
                <p class="text-green-600 font-semibold text-sm uppercase tracking-wide">Petani Terkoneksi</p>
              </div>

              <div class="flex flex-col items-center min-w-[200px]">
                <div class="w-16 h-16 bg-green-600 rounded-full flex items-center justify-center mb-5 shadow-lg shadow-green-600/30 text-white transition-transform hover:scale-110 duration-300">
                  <CalendarDays class="w-8 h-8" />
                </div>
                <h4 class="text-4xl font-extrabold text-slate-900 mb-2">10+</h4>
                <p class="text-green-600 font-semibold text-sm uppercase tracking-wide">Negara Tujuan Ekspor</p>
              </div>

              <div class="flex flex-col items-center min-w-[200px]">
                <div class="w-16 h-16 bg-green-600 rounded-full flex items-center justify-center mb-5 shadow-lg shadow-green-600/30 text-white transition-transform hover:scale-110 duration-300">
                  <Network class="w-8 h-8" />
                </div>
                <h4 class="text-4xl font-extrabold text-slate-900 mb-2">20+</h4>
                <p class="text-green-600 font-semibold text-sm uppercase tracking-wide">Tonase Produksi Tahunan</p>
              </div>

              <div class="flex flex-col items-center min-w-[200px]">
                <div class="w-16 h-16 bg-green-600 rounded-full flex items-center justify-center mb-5 shadow-lg shadow-green-600/30 text-white transition-transform hover:scale-110 duration-300">
                  <Users class="w-8 h-8" />
                </div>
                <h4 class="text-4xl font-extrabold text-slate-900 mb-2">50+</h4>
                <p class="text-green-600 font-semibold text-sm uppercase tracking-wide">Anggota Komunitas Aktif</p>
              </div>

              <div class="flex flex-col items-center min-w-[200px]">
                <div class="w-16 h-16 bg-green-600 rounded-full flex items-center justify-center mb-5 shadow-lg shadow-green-600/30 text-white transition-transform hover:scale-110 duration-300">
                  <Globe class="w-8 h-8" />
                </div>
                <h4 class="text-4xl font-extrabold text-slate-900 mb-2">6+</h4>
                <p class="text-green-600 font-semibold text-sm uppercase tracking-wide">Tahun Pengalaman Kolektif</p>
              </div>

            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Coverage Section -->
    <section class="bg-white py-24 border-b border-slate-100">
      <div class="max-w-7xl mx-auto px-6">
        <div class="text-center max-w-4xl mx-auto mb-16">
          <h3 class="text-4xl font-bold text-green-600 mb-6">Wilayah Cakupan Pondok Tani</h3>
          <p class="text-lg text-slate-600 leading-relaxed">
            Dengan pendekatan inovatif dan komitmen Pondok Tani terhadap pertanian berkelanjutan, Pondok Tani siap memimpin menuju masa depan yang lebih cerah bagi petani di seluruh Indonesia.
          </p>
        </div>

        <div class="flex flex-col lg:flex-row items-center gap-12">
          <!-- Map Area -->
          <div class="w-full lg:w-2/3 relative bg-slate-50/50 rounded-3xl p-4 lg:p-8 border border-slate-100 shadow-inner">
            <div class="relative w-full flex items-center justify-center">
              <img src="../../assets/Map.webp" alt="Peta Indonesia" class="w-full h-auto opacity-50 filter grayscale" />
              
              <!-- Map Markers -->
              <!-- Map Markers -->
              <!-- KEPRI -->
              <div class="absolute w-3 h-3 bg-yellow-400 rounded-full border-2 border-white shadow-md" style="top: 35%; left: 21%; transform: translate(-50%, -50%);">
              </div>

              <!-- RIAU -->
              <div class="absolute w-3 h-3 bg-yellow-400 rounded-full border-2 border-white shadow-md" style="top: 32%; left: 15%; transform: translate(-50%, -50%);">
              </div>

              <!-- SUMSEL -->
              <div class="absolute w-3 h-3 bg-yellow-400 rounded-full border-2 border-white shadow-md" style="top: 52%; left: 20%; transform: translate(-50%, -50%);">
              </div>

              <!-- BANTEN -->
              <div class="absolute w-3 h-3 bg-yellow-400 rounded-full border-2 border-white shadow-md" style="top: 71%; left: 24%; transform: translate(-50%, -50%);">
              </div>

              <!-- KALBAR -->
              <div class="absolute w-3 h-3 bg-yellow-400 rounded-full border-2 border-white shadow-md" style="top: 37%; left: 35%; transform: translate(-50%, -50%);">
              </div>

              <!-- KALTIM -->
              <div class="absolute w-3 h-3 bg-yellow-400 rounded-full border-2 border-white shadow-md" style="top: 35%; left: 47%; transform: translate(-50%, -50%);">
              </div>

              <!-- SULUT -->
              <div class="absolute w-3 h-3 bg-yellow-400 rounded-full border-2 border-white shadow-md" style="top: 30%; left: 63%; transform: translate(-50%, -50%);">
              </div>

              <!-- SORONG -->
              <div class="absolute w-3 h-3 bg-yellow-400 rounded-full border-2 border-white shadow-md" style="top: 41%; left: 82%; transform: translate(-50%, -50%);">
              </div>

              <div class="absolute w-3 h-3 bg-yellow-400 rounded-full border-2 border-white shadow-md" style="top: 74%; left: 28%; transform: translate(-50%, -50%);">
              </div>

              <div class="absolute w-3 h-3 bg-yellow-400 rounded-full border-2 border-white shadow-md" style="top: 79%; left: 33.6%; transform: translate(-50%, -50%);">
              </div>

              <div class="absolute w-3 h-3 bg-yellow-400 rounded-full border-2 border-white shadow-md" style="top: 76%; left: 32%; transform: translate(-50%, -50%);">
              </div>

            </div>
          </div>

          <!-- Stats Right -->
          <div class="w-full lg:w-1/3 flex flex-col justify-center">
            
            <div class="flex flex-col gap-3">
              <div class="flex items-center justify-between border-b border-slate-100 pb-2">
                <span class="text-slate-600 font-semibold text-lg uppercase tracking-wider">Jawa Barat</span>
                <span class="text-2xl font-bold text-green-600">208 <span class="text-sm text-slate-400 font-medium lowercase">ha</span></span>
              </div>
              <div class="flex items-center justify-between border-b border-slate-100 pb-2">
                <span class="text-slate-600 font-semibold text-lg uppercase tracking-wider">Kalimantan Timur</span>
                <span class="text-2xl font-bold text-green-600">100 <span class="text-sm text-slate-400 font-medium lowercase">ha</span></span>
              </div>
              <div class="flex items-center justify-between border-b border-slate-100 pb-2">
                <span class="text-slate-600 font-semibold text-lg uppercase tracking-wider">Kepulauan Riau</span>
                <span class="text-2xl font-bold text-green-600">50 <span class="text-sm text-slate-400 font-medium lowercase">ha</span></span>
              </div>
              <div class="flex items-center justify-between border-b border-slate-100 pb-2">
                <span class="text-slate-600 font-semibold text-lg uppercase tracking-wider">Banten</span>
                <span class="text-2xl font-bold text-green-600">35 <span class="text-sm text-slate-400 font-medium lowercase">ha</span></span>
              </div>
              <div class="flex items-center justify-between border-b border-slate-100 pb-2">
                <span class="text-slate-600 font-semibold text-lg uppercase tracking-wider">Jawa Tengah</span>
                <span class="text-2xl font-bold text-green-600">33 <span class="text-sm text-slate-400 font-medium lowercase">ha</span></span>
              </div>  
              <div class="flex items-center justify-between border-b border-slate-100 pb-2">
                <span class="text-slate-600 font-semibold text-lg uppercase tracking-wider">Sorong</span>
                <span class="text-2xl font-bold text-green-600">25 <span class="text-sm text-slate-400 font-medium lowercase">ha</span></span>
              </div>
              <div class="flex items-center justify-between border-b border-slate-100 pb-2">
                <span class="text-slate-600 font-semibold text-lg uppercase tracking-wider">Sumatra Selatan</span>
                <span class="text-2xl font-bold text-green-600">20 <span class="text-sm text-slate-400 font-medium lowercase">ha</span></span>
              </div>
              <div class="flex items-center justify-between border-b border-slate-100 pb-2">
                <span class="text-slate-600 font-semibold text-lg uppercase tracking-wider">Riau</span>
                <span class="text-2xl font-bold text-green-600">11 <span class="text-sm text-slate-400 font-medium lowercase">ha</span></span>
              </div>
              <div class="flex items-center justify-between border-b border-slate-100 pb-2">
                <span class="text-slate-600 font-semibold text-lg uppercase tracking-wider">Kalimantan Barat</span>
                <span class="text-2xl font-bold text-green-600">10 <span class="text-sm text-slate-400 font-medium lowercase">ha</span></span>
              </div>
              <div class="flex items-center justify-between border-b border-slate-100 pb-2">
                <span class="text-slate-600 font-semibold text-lg uppercase tracking-wider">Sulawesi Utara</span>
                <span class="text-2xl font-bold text-green-600">3 <span class="text-sm text-slate-400 font-medium lowercase">ha</span></span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-slate-600 font-semibold text-lg uppercase tracking-wider">Jogja</span>
                <span class="text-2xl font-bold text-green-600">0.5 <span class="text-sm text-slate-400 font-medium lowercase">ha</span></span>
              </div>              
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="bg-white py-20">
      <div class="max-w-7xl mx-auto px-6">
        <div class="text-center mb-16">
          <h3 class="text-4xl font-bold text-gray-900 mb-4">Mengapa Memilih Pondok Tani?</h3>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-6 gap-6">
          <div class="group rounded-3xl border border-slate-200 bg-white p-8 text-center shadow-lg shadow-slate-200/30 transition duration-300 hover:scale-[0.99] hover:border-green-400 hover:shadow-2xl hover:bg-slate-50 xl:col-span-2">
            <div class="flex h-14 w-14 items-center justify-center rounded-2xl border border-green-600 bg-white text-green-600 mx-auto mb-5 transition duration-300 group-hover:bg-green-600 group-hover:text-white group-hover:scale-90">
              <TrendingUp class="w-7 h-7 transition-transform duration-300" />
            </div>
            <h4 class="text-xl font-semibold text-slate-900 mb-3">Ekspor Vanili ke Afrika</h4>
            <p class="text-slate-600 leading-relaxed">
              Pengalaman ekspor vanili sampai ke Benua Afrika, membuka peluang pasar global untuk hasil pertanian kamu.
            </p>
          </div>

          <div class="group rounded-3xl border border-slate-200 bg-white p-8 text-center shadow-lg shadow-slate-200/30 transition duration-300 hover:scale-[0.99] hover:border-green-400 hover:shadow-2xl hover:bg-slate-50 xl:col-span-2">
            <div class="flex h-14 w-14 items-center justify-center rounded-2xl border border-green-600 bg-white text-green-600 mx-auto mb-5 transition duration-300 group-hover:bg-green-600 group-hover:text-white group-hover:scale-90">
              <Sprout class="w-7 h-7 transition-transform duration-300" />
            </div>
            <h4 class="text-xl font-semibold text-slate-900 mb-3">Hasil Sebelum Panen 3-5 Tahun</h4>
            <p class="text-slate-600 leading-relaxed">
              Menciptakan guna hasil sebelum masa panen 3 sampai 5 tahun, memberikan nilai tambah untuk setiap investasi.
            </p>
          </div>

          <div class="group rounded-3xl border border-slate-200 bg-white p-8 text-center shadow-lg shadow-slate-200/30 transition duration-300 hover:scale-[0.99] hover:border-green-400 hover:shadow-2xl hover:bg-slate-50 xl:col-span-2">
            <div class="flex h-14 w-14 items-center justify-center rounded-2xl border border-green-600 bg-white text-green-600 mx-auto mb-5 transition duration-300 group-hover:bg-green-600 group-hover:text-white group-hover:scale-90">
              <Users class="w-7 h-7 transition-transform duration-300" />
            </div>
            <h4 class="text-xl font-semibold text-slate-900 mb-3">Support System Budi Daya</h4>
            <p class="text-slate-600 leading-relaxed">
              Mempunyai support system dalam budi daya vanili, mulai dari pembibitan sampai panen dan pascapanen.
            </p>
          </div>

          <div class="group rounded-3xl border border-slate-200 bg-white p-8 text-center shadow-lg shadow-slate-200/30 transition duration-300 hover:scale-[0.99] hover:border-green-400 hover:shadow-2xl hover:bg-slate-50 md:col-span-2 xl:col-span-2 xl:col-start-2">
            <div class="flex h-14 w-14 items-center justify-center rounded-2xl border border-green-600 bg-white text-green-600 mx-auto mb-5 transition duration-300 group-hover:bg-green-600 group-hover:text-white group-hover:scale-90">
              <Shield class="w-7 h-7 transition-transform duration-300" />
            </div>
            <h4 class="text-xl font-semibold text-slate-900 mb-3">Vanili planifolia Premium</h4>
            <p class="text-slate-600 leading-relaxed">
              Vanili planifolia, spesies terkenal untuk kualitas buah terbaik, potensi buah banyak, dan kadar vanili tinggi.
            </p>
          </div>
          <div class="group rounded-3xl border border-slate-200 bg-white p-8 text-center shadow-lg shadow-slate-200/30 transition duration-300 hover:scale-[0.99] hover:border-green-400 hover:shadow-2xl hover:bg-slate-50 md:col-span-2 xl:col-span-2 xl:col-start-4">
            <div class="flex h-14 w-14 items-center justify-center rounded-2xl border border-green-600 bg-white text-green-600 mx-auto mb-5 transition duration-300 group-hover:bg-green-600 group-hover:text-white group-hover:scale-90">
              <BarChart3 class="w-7 h-7 transition-transform duration-300" />
            </div>
            <h4 class="text-xl font-semibold text-slate-900 mb-3">Garansi Hasil Panen</h4>
            <p class="text-slate-600 leading-relaxed">
              Garansi hasil panen pasti dibeli, memberikan keamanan dan kepastian nilai untuk investasi petani vanili.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- Program Section -->
    <section class="py-20 bg-gradient-to-b from-green-50 to-white">
      <div class="max-w-7xl mx-auto px-6">
        <div class="mb-12 text-center">
          <h3 class="mt-4 text-4xl font-bold text-slate-900">Program Investasi Atau Kemitraan Pondok Tani</h3>
          <p class="mt-3 text-lg text-slate-600">Pilih jalur yang sesuai dengan tujuan Anda: investor atau mitra aktif bersama tim Pondok Tani.</p>
        </div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <article class="rounded-3xl border border-green-200 bg-white p-8 text-center shadow-lg shadow-green-100/60 transition hover:-translate-y-1 hover:shadow-xl">
              <h4 class="text-2xl font-bold text-slate-900">Program Investasi</h4>
              <p class="mt-3 text-slate-600">Skema untuk investor yang ingin menanam modal dan memantau perkembangan kebun dari dashboard secara transparan.</p>
              <button class="mt-6 inline-flex items-center justify-center gap-2 rounded-xl bg-green-600 px-5 py-3 text-white font-semibold hover:bg-green-700" @click="router.push('/paket-investasi')">
                Lihat Program Investasi
                <ArrowRight class="h-4 w-4" />
              </button>
            </article>
          <article class="rounded-3xl border border-green-200 bg-white p-8 text-center shadow-lg shadow-green-100/60 transition hover:-translate-y-1 hover:shadow-xl">
            <h4 class="text-2xl font-bold text-slate-900">Program Kemitraan</h4>
            <p class="mt-3 text-slate-600">Skema untuk calon mitra yang ingin terlibat dalam operasional, budi daya, dan penguatan rantai pasok vanili.</p>
            <button class="mt-6 inline-flex items-center justify-center gap-2 rounded-xl bg-green-600 px-5 py-3 text-white font-semibold hover:bg-green-700" @click="router.push('/kemitraan')">
              Lihat Program Kemitraan
              <ArrowRight class="h-4 w-4" />
            </button>
          </article>
        </div>
      </div>
    </section>

    <!-- Artikel Section -->
    <section class="pb-20 bg-white">
      <div class="max-w-7xl mx-auto px-6">
        <div class="flex items-end justify-between flex-wrap gap-4 mb-10">
          <div>
            <h3 class="mt-4 text-4xl font-bold text-slate-900">Berita Vanili Terbaru</h3>
          </div>
          <a class="rounded-xl border border-green-200 px-5 py-3 font-semibold text-green-600 hover:bg-green-50" href="https://news.google.com/search?q=vanili&hl=id&gl=ID&ceid=ID:id" target="_blank" rel="noopener noreferrer">Lihat Semua Berita</a>
        </div>
        <div v-if="articles.length" class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <router-link
            v-for="item in articles"
            :key="item.id"
            class="rounded-2xl border border-slate-200 overflow-hidden hover:shadow-lg transition"
            :to="{ name: 'article-detail', params: { id: item.id } }"
          >
            <img v-if="item.image_url" class="h-48 w-full object-cover" :src="item.image_url" :alt="item.title" />
            <div class="h-48 w-full bg-gradient-to-br from-green-100 to-green-50 flex items-center justify-center text-green-700 font-semibold" v-else>Artikel Pondok Tani</div>
            <div class="p-6">
              <p class="text-xs font-semibold text-green-600">{{ item.category || 'Artikel' }}</p>
              <h4 class="mt-2 text-xl font-bold text-slate-900">{{ item.title }}</h4>
              <p class="mt-3 text-slate-600">{{ item.excerpt || item.content?.slice(0, 120) }}</p>
            </div>
          </router-link>
        </div>
        <div v-else class="rounded-2xl border border-dashed border-green-200 bg-green-50 p-8 text-center text-green-700">Belum ada artikel publik yang tersedia.</div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="bg-green-600 py-24 relative overflow-hidden">
      <!-- Decorative background elements -->
      <div class="absolute inset-0 opacity-20 bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] from-white via-transparent to-transparent"></div>
      <div class="absolute -top-32 -right-32 w-96 h-96 bg-green-500 rounded-full mix-blend-overlay filter blur-[100px] opacity-40"></div>
      <div class="absolute -bottom-32 -left-32 w-96 h-96 bg-green-300 rounded-full mix-blend-overlay filter blur-[100px] opacity-20"></div>

      <div class="max-w-4xl mx-auto px-6 text-center relative z-10">
        <h3 class="text-4xl md:text-5xl font-bold text-white mb-6 leading-tight drop-shadow-sm">Siap Memulai Perjalanan<br/>Bersama Pondok Tani Land?</h3>
        <p class="text-xl text-green-100 mb-12 leading-relaxed max-w-2xl mx-auto font-medium">
          Bergabunglah dengan ekosistem pertanian modern kami. Pantau lahan Anda dengan transparan dan raih hasil maksimal secara berkelanjutan.
        </p>
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <button
            type="button"
            class="px-8 py-4 bg-white text-green-600 rounded-xl hover:bg-green-50 transition-all font-bold text-lg shadow-xl shadow-green-600/20 hover:shadow-green-600/40 hover:-translate-y-1 flex items-center justify-center gap-2"
            @click="router.push({ path: '/login', query: { mode: 'register' } })"
          >
            Daftar Sekarang
            <ArrowRight class="w-5 h-5" />
          </button>
          <button
            type="button"
            class="px-8 py-4 border-2 border-green-400/50 text-white rounded-xl hover:bg-white/10 hover:border-white transition-all font-semibold text-lg flex items-center justify-center gap-2"
            @click="router.push('/login')"
          >
            Masuk ke Akun
          </button>
        </div>
      </div>
    </section>

  </LandingLayout>
</template>

<style scoped>
@keyframes marquee {
  0% {
    transform: translateX(0);
  }
  100% {
    transform: translateX(-50%);
  }
}

.animate-marquee {
  width: max-content;
  animation: marquee 25s linear infinite;
}

</style>
