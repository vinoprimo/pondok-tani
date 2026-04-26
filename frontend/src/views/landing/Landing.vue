<script setup lang="ts">
import { ArrowRight, Sprout, TrendingUp, Shield, BarChart3, Users, Handshake, BookOpenText } from 'lucide-vue-next';
import { useRouter } from "vue-router";
import LandingLayout from '../../layouts/LandingLayout.vue';
import heroImage from '../../assets/pondok-tani-hero.jpg';
import heroImageCard from '../../assets/pondok-tani-hero-card.png';
import berita1 from '../../assets/berita-1.jpg';
import berita2 from '../../assets/berita-2.jpg';
import berita3 from '../../assets/berita-3.jpg';

const router = useRouter();
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
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          <div class="rounded-[2rem] p-10 shadow-2xl">
            <h2 class="text-5xl font-bold mb-6 leading-tight text-white">
              Mulai Investasi Vanili sebagai
              <span class="text-green-600">  Aset Emas Hijau</span> Anda di Masa Depan
            </h2>
            <p class="text-xl mb-8 leading-relaxed text-white/85">
              Pondok Tani membantu Anda berinvestasi pada vanili sebagai aset emas hijau melalui pengelolaan yang profesional dan terstruktur, sehingga memberikan potensi return jangka panjang yang stabil dan berkelanjutan.
            </p>
            <div class="flex flex-col gap-4 sm:flex-row">
              <button
                type="button"
                class="px-8 py-4 bg-green-600 text-white rounded-xl hover:bg-green-600 transition-colors font-semibold flex items-center gap-2 text-lg"
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
          <div class="relative rounded-3xl overflow-hidden h-96 w-full ring-5 ring-white/90 bg-white/10 backdrop-blur-sm">
            <div
              class="absolute inset-0 bg-cover bg-center opacity-90"
              :style="{ backgroundImage: `url(${heroImageCard})` }"
            ></div>
            <div class="absolute inset-0 bg-gradient-to-t from-black/90 via-transparent to-transparent"></div>
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
          <a class="rounded-xl border border-green-200 px-5 py-3 font-semibold text-green-700 hover:bg-green-50" href="https://news.google.com/search?q=vanili&hl=id&gl=ID&ceid=ID:id" target="_blank" rel="noopener noreferrer">Lihat Semua Berita</a>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <a class="rounded-2xl border border-slate-200 overflow-hidden hover:shadow-lg transition" href="https://news.google.com/search?q=vanili+indonesia+terbaru&hl=id&gl=ID&ceid=ID:id" target="_blank" rel="noopener noreferrer">
            <img class="h-48 w-full object-cover" :src="berita1" alt="Berita vanili Indonesia" />
            <div class="p-6">
                <p class="text-xs font-semibold text-green-600">Google News</p>
              <h4 class="mt-2 text-xl font-bold text-slate-900">Berita Vanili Indonesia Terbaru</h4>
              <p class="mt-3 text-slate-600">Update perkembangan budidaya, petani, dan industri vanili di Indonesia.</p>
            </div>
          </a>
          <a class="rounded-2xl border border-slate-200 overflow-hidden hover:shadow-lg transition" href="https://news.google.com/search?q=harga+vanili+dunia&hl=id&gl=ID&ceid=ID:id" target="_blank" rel="noopener noreferrer">
            <img class="h-48 w-full object-cover" :src="berita2" alt="Berita harga vanili dunia" />
            <div class="p-6">
                <p class="text-xs font-semibold text-green-600">Google News</p>
              <h4 class="mt-2 text-xl font-bold text-slate-900">Pergerakan Harga Vanili Global</h4>
              <p class="mt-3 text-slate-600">Pantau kabar terbaru terkait harga vanili dunia dan permintaan pasar internasional.</p>
            </div>
          </a>
          <a class="rounded-2xl border border-slate-200 overflow-hidden hover:shadow-lg transition" href="https://news.google.com/search?q=teknik+budidaya+vanili+terkini&hl=id&gl=ID&ceid=ID:id" target="_blank" rel="noopener noreferrer">
            <img class="h-48 w-full object-cover" :src="berita3" alt="Teknik budidaya vanili" />
            <div class="p-6">
                <p class="text-xs font-semibold text-green-600">Google News</p>
              <h4 class="mt-2 text-xl font-bold text-slate-900">Teknik Budidaya Vanili Terkini</h4>
              <p class="mt-3 text-slate-600">Kumpulan berita dan tips terbaru seputar teknik budidaya untuk hasil optimal.</p>
            </div>
          </a>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="bg-gradient-to-r from-green-700 to-green-800 py-20">
      <div class="max-w-4xl mx-auto px-6 text-center">
        <h3 class="text-4xl font-bold text-white mb-6">Siap Memulai Investasi Vanili?</h3>
        <p class="text-xl text-green-50 mb-8">
          Bergabunglah dengan ribuan investor sukses dan raih keuntungan hingga 30% per tahun
        </p>
        <div class="flex gap-4 justify-center">
          <button
            type="button"
            class="px-8 py-4 bg-white text-green-700 rounded-xl hover:bg-green-50 transition-colors font-semibold text-lg flex items-center gap-2"
            @click="router.push({ path: '/login', query: { mode: 'register' } })"
          >
            Mulai Sekarang
            <ArrowRight class="w-5 h-5" />
          </button>
          <button
            type="button"
            class="px-8 py-4 border-2 border-white text-white rounded-xl hover:bg-white/10 transition-colors font-semibold text-lg"
            @click="router.push('/login')"
          >
            Masuk sebagai admin
          </button>
        </div>
      </div>
    </section>

  </LandingLayout>
</template>
