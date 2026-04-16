<script setup lang="ts">
import { ArrowRight, Sprout, TrendingUp, Shield, BarChart3, Users } from 'lucide-vue-next';
import { useRouter } from "vue-router";

const router = useRouter();
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-green-50 via-white to-green-50">
    <!-- Header -->
    <header class="border-b border-gray-200 bg-white/80 backdrop-blur-sm sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 bg-green-600 rounded-lg flex items-center justify-center">
            <Sprout class="w-6 h-6 text-white" />
          </div>
          <div>
            <h1 class="font-semibold text-gray-900">Omah Vanili</h1>
            <p class="text-xs text-gray-500">Sistem Investasi Vanili</p>
          </div>
        </div>
        <div class="flex items-center gap-4">
          <button
            type="button"
            class="px-4 py-2 text-green-600 hover:bg-green-50 rounded-lg transition-colors font-medium"
            @click="router.push('/login')"
          >
            Masuk
          </button>
          <button
            type="button"
            class="px-6 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-medium"
            @click="router.push({ path: '/login', query: { mode: 'register' } })"
          >
            Daftar Sekarang
          </button>
        </div>
      </div>
    </header>

    <!-- Hero Section -->
    <section class="max-w-7xl mx-auto px-6 py-20">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
        <div>
          <div class="inline-block px-4 py-2 bg-green-100 text-green-700 rounded-full text-sm font-medium mb-6">
            🌱 Platform Investasi Pertanian Terpercaya
          </div>
          <h2 class="text-5xl font-bold text-gray-900 mb-6 leading-tight">
            Investasi Vanili
            <br />
            <span class="text-green-600">Menguntungkan & Transparan</span>
          </h2>
          <p class="text-xl text-gray-600 mb-8 leading-relaxed">
            Bergabunglah dengan ribuan investor yang telah meraih keuntungan dari investasi vanili premium.
            Pantau perkembangan tanaman Anda secara real-time dengan teknologi terkini.
          </p>
          <div class="flex gap-4">
            <button
              type="button"
              class="px-8 py-4 bg-green-600 text-white rounded-xl hover:bg-green-700 transition-colors font-semibold flex items-center gap-2 text-lg"
              @click="router.push('/login')"
            >
              Mulai Investasi
              <ArrowRight class="w-5 h-5" />
            </button>
            <button
              type="button"
              class="px-8 py-4 border-2 border-gray-300 text-gray-700 rounded-xl hover:bg-gray-50 transition-colors font-semibold text-lg"
            >
              Pelajari Lebih Lanjut
            </button>
          </div>
          <div class="flex items-center gap-8 mt-12">
            <div>
              <p class="text-3xl font-bold text-gray-900">1,200+</p>
              <p class="text-gray-600">Investor Aktif</p>
            </div>
            <div>
              <p class="text-3xl font-bold text-gray-900">25%+</p>
              <p class="text-gray-600">ROI Rata-rata</p>
            </div>
            <div>
              <p class="text-3xl font-bold text-gray-900">50 Ha</p>
              <p class="text-gray-600">Lahan Produktif</p>
            </div>
          </div>
        </div>
        <div class="relative">
          <div class="bg-gradient-to-br from-green-100 to-green-50 rounded-3xl p-8 shadow-2xl">
            <img
              src="https://images.unsplash.com/photo-1637922808382-0e5930886159?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx2YW5pbGxhJTIwcGxhbnRhdGlvbiUyMGZhcm0lMjBhZ3JpY3VsdHVyZXxlbnwxfHx8fDE3NzQ4NDIxNzJ8MA&ixlib=rb-4.1.0&q=80&w=1080"
              alt="Perkebunan vanili"
              class="rounded-2xl shadow-lg w-full h-96 object-cover"
            />
            <div class="absolute -bottom-6 -left-6 bg-white rounded-2xl shadow-xl p-6 border border-gray-100">
              <div class="flex items-center gap-3 mb-2">
                <TrendingUp class="w-8 h-8 text-green-600" />
                <div>
                  <p class="text-2xl font-bold text-green-600">+32.5%</p>
                  <p class="text-sm text-gray-600">ROI Tahun Ini</p>
                </div>
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
          <h3 class="text-4xl font-bold text-gray-900 mb-4">Mengapa Memilih Omah Vanili?</h3>
          <p class="text-xl text-gray-600">Platform investasi vanili terlengkap dengan fitur canggih</p>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div
            class="bg-gradient-to-br from-green-50 to-white rounded-2xl p-8 border border-gray-100 hover:shadow-xl transition-shadow"
          >
            <div class="w-14 h-14 bg-green-600 rounded-xl flex items-center justify-center mb-4">
              <BarChart3 class="w-8 h-8 text-white" />
            </div>
            <h4 class="text-xl font-semibold text-gray-900 mb-3">Monitoring Real-time</h4>
            <p class="text-gray-600 leading-relaxed">
              Pantau pertumbuhan tanaman, hasil panen, dan proyeksi keuntungan Anda secara langsung dari dashboard yang mudah digunakan.
            </p>
          </div>

          <div
            class="bg-gradient-to-br from-blue-50 to-white rounded-2xl p-8 border border-gray-100 hover:shadow-xl transition-shadow"
          >
            <div class="w-14 h-14 bg-blue-600 rounded-xl flex items-center justify-center mb-4">
              <Shield class="w-8 h-8 text-white" />
            </div>
            <h4 class="text-xl font-semibold text-gray-900 mb-3">Transparan & Aman</h4>
            <p class="text-gray-600 leading-relaxed">
              Setiap aktivitas perawatan divalidasi dengan foto dan dokumen. Investasi Anda dijamin aman dan terverifikasi.
            </p>
          </div>

          <div
            class="bg-gradient-to-br from-orange-50 to-white rounded-2xl p-8 border border-gray-100 hover:shadow-xl transition-shadow"
          >
            <div class="w-14 h-14 bg-orange-600 rounded-xl flex items-center justify-center mb-4">
              <Users class="w-8 h-8 text-white" />
            </div>
            <h4 class="text-xl font-semibold text-gray-900 mb-3">Dukungan Penuh</h4>
            <p class="text-gray-600 leading-relaxed">
              Tim ahli kami siap membantu Anda kapan saja. Dapatkan laporan lengkap dan konsultasi gratis untuk investasi Anda.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- How It Works Section -->
    <section class="py-20">
      <div class="max-w-7xl mx-auto px-6">
        <div class="text-center mb-16">
          <h3 class="text-4xl font-bold text-gray-900 mb-4">Cara Kerja</h3>
          <p class="text-xl text-gray-600">Mulai investasi vanili hanya dalam 3 langkah mudah</p>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div class="text-center">
            <div
              class="w-20 h-20 bg-green-600 rounded-full flex items-center justify-center mx-auto mb-6 text-white text-3xl font-bold"
            >
              1
            </div>
            <h4 class="text-xl font-semibold text-gray-900 mb-3">Daftar & Pilih Paket</h4>
            <p class="text-gray-600">
              Buat akun dan pilih paket investasi yang sesuai dengan budget Anda. Mulai dari Rp 5 juta.
            </p>
          </div>

          <div class="text-center">
            <div
              class="w-20 h-20 bg-green-600 rounded-full flex items-center justify-center mx-auto mb-6 text-white text-3xl font-bold"
            >
              2
            </div>
            <h4 class="text-xl font-semibold text-gray-900 mb-3">Pantau Perkembangan</h4>
            <p class="text-gray-600">
              Akses dashboard untuk melihat pertumbuhan tanaman, aktivitas perawatan, dan estimasi hasil panen.
            </p>
          </div>

          <div class="text-center">
            <div
              class="w-20 h-20 bg-green-600 rounded-full flex items-center justify-center mx-auto mb-6 text-white text-3xl font-bold"
            >
              3
            </div>
            <h4 class="text-xl font-semibold text-gray-900 mb-3">Terima Keuntungan</h4>
            <p class="text-gray-600">
              Setelah panen, terima bagi hasil sesuai kesepakatan. ROI transparan dan tepat waktu.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="bg-gradient-to-r from-green-600 to-green-700 py-20">
      <div class="max-w-4xl mx-auto px-6 text-center">
        <h3 class="text-4xl font-bold text-white mb-6">Siap Memulai Investasi Vanili?</h3>
        <p class="text-xl text-green-50 mb-8">
          Bergabunglah dengan ribuan investor sukses dan raih keuntungan hingga 30% per tahun
        </p>
        <div class="flex gap-4 justify-center">
          <button
            type="button"
            class="px-8 py-4 bg-white text-green-600 rounded-xl hover:bg-green-50 transition-colors font-semibold text-lg flex items-center gap-2"
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

    <!-- Footer -->
    <footer class="bg-gray-900 text-gray-300 py-12">
      <div class="max-w-7xl mx-auto px-6">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-8 mb-8">
          <div>
            <div class="flex items-center gap-3 mb-4">
              <div class="w-10 h-10 bg-green-600 rounded-lg flex items-center justify-center">
                <Sprout class="w-6 h-6 text-white" />
              </div>
              <div>
                <h4 class="font-semibold text-white">Omah Vanili</h4>
              </div>
            </div>
            <p class="text-sm text-gray-400">Platform investasi pertanian vanili terpercaya di Indonesia</p>
          </div>

          <div>
            <h5 class="font-semibold text-white mb-4">Tentang</h5>
            <ul class="space-y-2 text-sm">
              <li><a href="#" class="hover:text-green-400 transition-colors">Tentang Kami</a></li>
              <li><a href="#" class="hover:text-green-400 transition-colors">Tim</a></li>
              <li><a href="#" class="hover:text-green-400 transition-colors">Karir</a></li>
            </ul>
          </div>

          <div>
            <h5 class="font-semibold text-white mb-4">Layanan</h5>
            <ul class="space-y-2 text-sm">
              <li><a href="#" class="hover:text-green-400 transition-colors">Investasi</a></li>
              <li><a href="#" class="hover:text-green-400 transition-colors">Monitoring</a></li>
              <li><a href="#" class="hover:text-green-400 transition-colors">Laporan</a></li>
            </ul>
          </div>

          <div>
            <h5 class="font-semibold text-white mb-4">Kontak</h5>
            <ul class="space-y-2 text-sm">
              <li class="text-gray-400">Email: info@omahvanili.com</li>
              <li class="text-gray-400">Telp: +62 812-3456-7890</li>
              <li class="text-gray-400">Alamat: Yogyakarta, Indonesia</li>
            </ul>
          </div>
        </div>

        <div class="border-t border-gray-800 pt-8 text-center text-sm text-gray-400">
          <p>&copy; 2024 Sistem Investasi Omah Vanili. Semua hak dilindungi.</p>
        </div>
      </div>
    </footer>
  </div>
</template>
