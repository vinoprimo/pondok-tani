<template>
    <div class="bg-white text-gray-800">
  
      <!-- 🔝 NAVBAR -->
      <header class="flex justify-between items-center px-20 py-5 border-b border-gray-100">
        <div class="flex items-center gap-3">
          <div class="bg-green-600 w-10 h-10 rounded-xl"></div>
          <div>
            <p class="font-semibold text-gray-900">Omah Vanili</p>
            <p class="text-xs text-gray-400">Sistem Investasi Vanili</p>
          </div>
        </div>
  
        <div class="flex items-center gap-6">
          <router-link to="/login" class="text-green-600 font-medium">
            Masuk
          </router-link>
  
          <button class="bg-green-600 text-white px-5 py-2.5 rounded-xl font-medium shadow-sm hover:opacity-90 transition">
            Daftar Sekarang
          </button>
        </div>
      </header>
  
      <!-- 🌿 HERO -->
      <section class="grid grid-cols-2 gap-16 px-20 py-20 items-center">
  
        <!-- LEFT -->
        <div>
          <div class="inline-flex items-center gap-2 bg-green-50 text-green-600 px-4 py-2 rounded-full text-sm mb-6">
            🌱 Platform Investasi Pertanian Terpercaya
          </div>
  
          <h1 class="text-[52px] leading-[60px] font-bold text-gray-900">
            Investasi Vanili
            <br />
            <span class="text-green-600">Menguntungkan & Transparan</span>
          </h1>
  
          <p class="mt-6 text-gray-500 text-[16px] leading-relaxed max-w-xl">
            Bergabunglah dengan ribuan investor yang telah meraih keuntungan dari investasi vanili premium.
            Pantau perkembangan tanaman Anda secara real-time dengan teknologi terkini.
          </p>
  
          <div class="mt-8 flex gap-4">
            <button class="bg-green-600 text-white px-6 py-3 rounded-xl font-medium shadow-md hover:opacity-90 transition">
              Mulai Investasi →
            </button>
  
            <button class="border border-gray-200 px-6 py-3 rounded-xl text-gray-600 hover:bg-gray-50 transition">
              Pelajari Lebih Lanjut
            </button>
          </div>
        </div>
  
        <!-- RIGHT -->
        <div class="relative">
          <img
            src="https://images.unsplash.com/photo-1598514982494-5c5d2f6c0e1c"
            class="rounded-2xl w-full object-cover shadow-sm"
          />
  
          <div class="absolute bottom-6 left-6 bg-white px-5 py-3 rounded-xl shadow-md">
            <p class="text-green-600 font-semibold text-sm">+32.5%</p>
            <p class="text-xs text-gray-500">ROI Tahun Ini</p>
          </div>
        </div>
      </section>
  
      <!-- 💡 FEATURES -->
      <section class="bg-gray-50 py-20 px-20 text-center">
        <h2 class="text-3xl font-bold text-gray-900">
          Mengapa Memilih Omah Vanili?
        </h2>
  
        <p class="text-gray-500 mt-3 mb-14">
          Platform investasi vanili terlengkap dengan fitur canggih
        </p>
  
        <div class="grid grid-cols-3 gap-8 text-left">
  
          <div class="bg-green-50 p-8 rounded-2xl">
            <div class="bg-green-600 w-12 h-12 rounded-xl mb-4"></div>
            <h3 class="font-semibold text-lg text-gray-900">Monitoring Real-time</h3>
            <p class="text-gray-500 text-sm mt-2">
              Pantau pertumbuhan tanaman, hasil panen, dan proyeksi keuntungan.
            </p>
          </div>
  
          <div class="bg-white p-8 rounded-2xl shadow-sm">
            <div class="bg-blue-500 w-12 h-12 rounded-xl mb-4"></div>
            <h3 class="font-semibold text-lg text-gray-900">Transparan & Aman</h3>
            <p class="text-gray-500 text-sm mt-2">
              Semua aktivitas tervalidasi dengan dokumentasi lengkap.
            </p>
          </div>
  
          <div class="bg-orange-50 p-8 rounded-2xl">
            <div class="bg-orange-500 w-12 h-12 rounded-xl mb-4"></div>
            <h3 class="font-semibold text-lg text-gray-900">Dukungan Penuh</h3>
            <p class="text-gray-500 text-sm mt-2">
              Tim ahli siap membantu kapan saja.
            </p>
          </div>
  
        </div>
      </section>
  
      <!-- ⚙️ CARA KERJA -->
      <section class="py-20 px-20 text-center">
        <h2 class="text-3xl font-bold text-gray-900">Cara Kerja</h2>
  
        <p class="text-gray-500 mt-3 mb-14">
          Mulai investasi vanili hanya dalam 3 langkah mudah
        </p>
  
        <div class="grid grid-cols-3 gap-10">
  
          <div>
            <div class="bg-green-600 text-white w-14 h-14 flex items-center justify-center rounded-full mx-auto mb-4 text-lg font-bold">
              1
            </div>
            <h3 class="font-semibold text-gray-900">Daftar & Pilih Paket</h3>
            <p class="text-gray-500 text-sm mt-2">
              Pilih paket sesuai budget Anda
            </p>
          </div>
  
          <div>
            <div class="bg-green-600 text-white w-14 h-14 flex items-center justify-center rounded-full mx-auto mb-4 text-lg font-bold">
              2
            </div>
            <h3 class="font-semibold text-gray-900">Pantau Perkembangan</h3>
            <p class="text-gray-500 text-sm mt-2">
              Akses dashboard monitoring tanaman
            </p>
          </div>
  
          <div>
            <div class="bg-green-600 text-white w-14 h-14 flex items-center justify-center rounded-full mx-auto mb-4 text-lg font-bold">
              3
            </div>
            <h3 class="font-semibold text-gray-900">Terima Keuntungan</h3>
            <p class="text-gray-500 text-sm mt-2">
              ROI transparan dan tepat waktu
            </p>
          </div>
  
        </div>
      </section>
  
      <!-- 🚀 CTA -->
      <section class="bg-green-600 text-white text-center py-20">
        <h2 class="text-3xl font-bold">
          Siap Memulai Investasi Vanili?
        </h2>
  
        <p class="mt-3 mb-8 text-green-100">
          Bergabunglah dan raih keuntungan hingga 30% per tahun
        </p>
  
        <div class="flex justify-center gap-4">
          <button class="bg-white text-green-600 px-6 py-3 rounded-xl font-semibold shadow-md">
            Mulai Sekarang →
          </button>
  
          <button class="border border-white px-6 py-3 rounded-xl">
            Login Admin
          </button>
        </div>
      </section>
  
      <!-- 📌 FOOTER -->
      <footer class="bg-gray-900 text-gray-300 px-20 py-14">
        <div class="grid grid-cols-4 gap-10">
  
          <div>
            <h3 class="text-white font-semibold mb-3">Omah Vanili</h3>
            <p class="text-sm">
              Platform investasi pertanian vanili terpercaya di Indonesia
            </p>
          </div>
  
          <div>
            <h4 class="text-white font-medium mb-3">Tentang</h4>
            <p class="text-sm">Tentang Kami</p>
            <p class="text-sm">Tim</p>
            <p class="text-sm">Karir</p>
          </div>
  
          <div>
            <h4 class="text-white font-medium mb-3">Layanan</h4>
            <p class="text-sm">Investasi</p>
            <p class="text-sm">Monitoring</p>
            <p class="text-sm">Laporan</p>
          </div>
  
          <div>
            <h4 class="text-white font-medium mb-3">Kontak</h4>
            <p class="text-sm">info@omahvanili.com</p>
            <p class="text-sm">+62 812-3456-7890</p>
            <p class="text-sm">Yogyakarta</p>
          </div>
  
        </div>
  
        <p class="text-center text-sm mt-10 text-gray-500">
          © 2024 Omah Vanili Investment System
        </p>
      </footer>
  
    </div>
  </template>