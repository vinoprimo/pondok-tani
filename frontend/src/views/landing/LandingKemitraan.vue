<script setup lang="ts">
import { useRouter } from 'vue-router';
import { 
  CheckCircle, Users, Leaf, TrendingUp, Award, Target, BookOpen, 
  Handshake, MapPin, Sprout, FileCheck, ChevronRight, MessageCircle, 
  DollarSign, Shield 
} from 'lucide-vue-next';
import LandingLayout from '../../layouts/LandingLayout.vue';

const router = useRouter();
</script>

<template>
  <LandingLayout>
    <!-- Hero Section -->
    <section class="bg-gradient-to-r from-green-600 to-green-600 text-white py-20 relative overflow-hidden">
      <!-- Decorative background elements -->
      <div class="absolute inset-0 opacity-10 bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] from-white via-transparent to-transparent"></div>
      
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        <div class="max-w-3xl">
          <h1 class="text-4xl md:text-5xl font-bold mb-6 leading-tight drop-shadow-sm">
            Bangun Kemitraan Vanili yang Tumbuh Bersama
          </h1>
          <p class="text-xl text-green-100 leading-relaxed mb-8">
            Bergabunglah dengan program kemitraan Pondok Tani Land dan dapatkan dukungan penuh dari
            pelatihan, pendampingan teknis, akses pasar, hingga jaminan pembelian hasil panen.
          </p>
        </div>
      </div>
    </section>

    <!-- Benefits Overview -->
    <section class="py-20 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-16">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
            Keuntungan Menjadi Mitra
          </h2>
          <p class="text-xl text-gray-600">
            Berbagai benefit yang mendukung kesuksesan budidaya vanili Anda
          </p>
        </div>
        <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
          <div class="bg-gradient-to-br from-green-50 to-green-100 p-8 rounded-xl text-center border border-green-100 shadow-sm hover:shadow-md transition">
            <div class="w-16 h-16 bg-green-600 rounded-full flex items-center justify-center mx-auto mb-6 shadow-md">
              <Target class="w-8 h-8 text-white" />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">Akses Pasar Terjamin</h3>
            <p class="text-gray-700 leading-relaxed">
              Kami membeli seluruh hasil panen berkualitas dengan harga kompetitif dan transparan
            </p>
          </div>
          <div class="bg-gradient-to-br from-green-50 to-emerald-100 p-8 rounded-xl text-center border border-emerald-100 shadow-sm hover:shadow-md transition">
            <div class="w-16 h-16 bg-emerald-600 rounded-full flex items-center justify-center mx-auto mb-6 shadow-md">
              <BookOpen class="w-8 h-8 text-white" />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">Pelatihan Gratis</h3>
            <p class="text-gray-700 leading-relaxed">
              Dapatkan pelatihan budidaya vanili dari ahli agrikultur berpengalaman
            </p>
          </div>
          <div class="bg-gradient-to-br from-teal-50 to-teal-100 p-8 rounded-xl text-center border border-teal-100 shadow-sm hover:shadow-md transition">
            <div class="w-16 h-16 bg-teal-600 rounded-full flex items-center justify-center mx-auto mb-6 shadow-md">
              <Users class="w-8 h-8 text-white" />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">Pendampingan Rutin</h3>
            <p class="text-gray-700 leading-relaxed">
              Tim lapangan siap membantu dengan kunjungan berkala dan konsultasi teknis
            </p>
          </div>
          <div class="bg-gradient-to-br from-lime-50 to-lime-100 p-8 rounded-xl text-center border border-lime-100 shadow-sm hover:shadow-md transition">
            <div class="w-16 h-16 bg-lime-600 rounded-full flex items-center justify-center mx-auto mb-6 shadow-md">
              <Sprout class="w-8 h-8 text-white" />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">Bibit Berkualitas</h3>
            <p class="text-gray-700 leading-relaxed">
              Akses ke bibit vanili unggul dengan harga khusus untuk kemitraan
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- Partnership Requirements -->
    <section class="py-20 bg-slate-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-16">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
            Syarat Menjadi Mitra
          </h2>
          <p class="text-xl text-gray-600">
            Persyaratan yang harus dipenuhi untuk bergabung bersama kami
          </p>
        </div>
        <div class="grid md:grid-cols-2 gap-8 max-w-5xl mx-auto">
          <div class="bg-white p-8 rounded-xl shadow-sm border border-slate-200">
            <div class="flex items-start gap-4">
              <div class="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center flex-shrink-0">
                <MapPin class="w-6 h-6 text-green-600" />
              </div>
              <div>
                <h3 class="text-xl font-bold text-gray-900 mb-2">Kepemilikan Lahan</h3>
                <p class="text-gray-600 mb-4 leading-relaxed">
                  Memiliki atau menguasai lahan minimal 0.5 hektar yang sesuai untuk budidaya vanili
                </p>
                <ul class="space-y-2 text-sm text-gray-600">
                  <li class="flex items-start gap-2">
                    <CheckCircle class="w-4 h-4 text-green-600 flex-shrink-0 mt-0.5" />
                    <span>Lahan berkontur atau dataran rendah</span>
                  </li>
                  <li class="flex items-start gap-2">
                    <CheckCircle class="w-4 h-4 text-green-600 flex-shrink-0 mt-0.5" />
                    <span>Akses air memadai</span>
                  </li>
                  <li class="flex items-start gap-2">
                    <CheckCircle class="w-4 h-4 text-green-600 flex-shrink-0 mt-0.5" />
                    <span>Naungan alami atau buatan tersedia</span>
                  </li>
                </ul>
              </div>
            </div>
          </div>

          <div class="bg-white p-8 rounded-xl shadow-sm border border-slate-200">
            <div class="flex items-start gap-4">
              <div class="w-12 h-12 bg-emerald-100 rounded-lg flex items-center justify-center flex-shrink-0">
                <FileCheck class="w-6 h-6 text-emerald-700" />
              </div>
              <div>
                <h3 class="text-xl font-bold text-gray-900 mb-2">Dokumen Administrasi</h3>
                <p class="text-gray-600 mb-4 leading-relaxed">
                  Melengkapi dokumen-dokumen yang diperlukan untuk verifikasi kemitraan
                </p>
                <ul class="space-y-2 text-sm text-gray-600">
                  <li class="flex items-start gap-2">
                    <CheckCircle class="w-4 h-4 text-emerald-600 flex-shrink-0 mt-0.5" />
                    <span>KTP dan Kartu Keluarga</span>
                  </li>
                  <li class="flex items-start gap-2">
                    <CheckCircle class="w-4 h-4 text-emerald-600 flex-shrink-0 mt-0.5" />
                    <span>Bukti kepemilikan/sewa lahan</span>
                  </li>
                  <li class="flex items-start gap-2">
                    <CheckCircle class="w-4 h-4 text-emerald-600 flex-shrink-0 mt-0.5" />
                    <span>Foto lahan dari berbagai sudut</span>
                  </li>
                </ul>
              </div>
            </div>
          </div>

          <div class="bg-white p-8 rounded-xl shadow-sm border border-slate-200">
            <div class="flex items-start gap-4">
              <div class="w-12 h-12 bg-teal-100 rounded-lg flex items-center justify-center flex-shrink-0">
                <Handshake class="w-6 h-6 text-teal-700" />
              </div>
              <div>
                <h3 class="text-xl font-bold text-gray-900 mb-2">Komitmen Kerja Sama</h3>
                <p class="text-gray-600 mb-4 leading-relaxed">
                  Bersedia mengikuti program kemitraan sesuai standar yang ditetapkan
                </p>
                <ul class="space-y-2 text-sm text-gray-600">
                  <li class="flex items-start gap-2">
                    <CheckCircle class="w-4 h-4 text-teal-600 flex-shrink-0 mt-0.5" />
                    <span>Mengikuti pelatihan wajib</span>
                  </li>
                  <li class="flex items-start gap-2">
                    <CheckCircle class="w-4 h-4 text-teal-600 flex-shrink-0 mt-0.5" />
                    <span>Menerapkan SOP budidaya</span>
                  </li>
                  <li class="flex items-start gap-2">
                    <CheckCircle class="w-4 h-4 text-teal-600 flex-shrink-0 mt-0.5" />
                    <span>Kontrak kemitraan minimal 3 tahun</span>
                  </li>
                </ul>
              </div>
            </div>
          </div>

          <div class="bg-white p-8 rounded-xl shadow-sm border border-slate-200">
            <div class="flex items-start gap-4">
              <div class="w-12 h-12 bg-lime-100 rounded-lg flex items-center justify-center flex-shrink-0">
                <MessageCircle class="w-6 h-6 text-lime-700" />
              </div>
              <div>
                <h3 class="text-xl font-bold text-gray-900 mb-2">Komunikasi Aktif</h3>
                <p class="text-gray-600 mb-4 leading-relaxed">
                  Mampu berkomunikasi dan melaporkan perkembangan secara berkala
                </p>
                <ul class="space-y-2 text-sm text-gray-600">
                  <li class="flex items-start gap-2">
                    <CheckCircle class="w-4 h-4 text-lime-600 flex-shrink-0 mt-0.5" />
                    <span>Memiliki smartphone atau alat komunikasi</span>
                  </li>
                  <li class="flex items-start gap-2">
                    <CheckCircle class="w-4 h-4 text-lime-600 flex-shrink-0 mt-0.5" />
                    <span>Dapat mengakses WhatsApp/email</span>
                  </li>
                  <li class="flex items-start gap-2">
                    <CheckCircle class="w-4 h-4 text-lime-600 flex-shrink-0 mt-0.5" />
                    <span>Responsif terhadap update dari tim</span>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Partnership Process -->
    <section class="py-20 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-16">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
            Alur Bergabung Kemitraan
          </h2>
          <p class="text-xl text-gray-600">
            Proses mudah dan terstruktur untuk memulai kemitraan bersama kami
          </p>
        </div>
        <div class="relative">
          <!-- Timeline Line -->
          <div class="hidden md:block absolute top-1/2 left-0 right-0 h-1 bg-gradient-to-r from-green-200 via-emerald-200 to-teal-200 transform -translate-y-1/2"></div>

          <div class="grid md:grid-cols-5 gap-8 relative">
            <div class="text-center">
              <div class="relative mb-6">
                <div class="w-20 h-20 bg-gradient-to-br from-green-500 to-green-600 text-white rounded-full flex items-center justify-center mx-auto font-bold text-2xl shadow-lg relative z-10 border-4 border-white">
                  1
                </div>
              </div>
              <h3 class="text-lg font-bold text-gray-900 mb-2">Pendaftaran Online</h3>
              <p class="text-sm text-gray-600">
                Isi formulir pendaftaran dan upload dokumen persyaratan
              </p>
            </div>

            <div class="text-center">
              <div class="relative mb-6">
                <div class="w-20 h-20 bg-gradient-to-br from-emerald-500 to-emerald-600 text-white rounded-full flex items-center justify-center mx-auto font-bold text-2xl shadow-lg relative z-10 border-4 border-white">
                  2
                </div>
              </div>
              <h3 class="text-lg font-bold text-gray-900 mb-2">Verifikasi Data</h3>
              <p class="text-sm text-gray-600">
                Tim akan memverifikasi kelengkapan dokumen (1-3 hari kerja)
              </p>
            </div>

            <div class="text-center">
              <div class="relative mb-6">
                <div class="w-20 h-20 bg-gradient-to-br from-teal-500 to-teal-600 text-white rounded-full flex items-center justify-center mx-auto font-bold text-2xl shadow-lg relative z-10 border-4 border-white">
                  3
                </div>
              </div>
              <h3 class="text-lg font-bold text-gray-900 mb-2">Survei Lapangan</h3>
              <p class="text-sm text-gray-600">
                Kunjungan tim untuk menilai kelayakan lahan dan kondisi
              </p>
            </div>

            <div class="text-center">
              <div class="relative mb-6">
                <div class="w-20 h-20 bg-gradient-to-br from-lime-500 to-lime-600 text-white rounded-full flex items-center justify-center mx-auto font-bold text-2xl shadow-lg relative z-10 border-4 border-white">
                  4
                </div>
              </div>
              <h3 class="text-lg font-bold text-gray-900 mb-2">Pelatihan Awal</h3>
              <p class="text-sm text-gray-600">
                Ikuti pelatihan intensif budidaya vanili (2-3 hari)
              </p>
            </div>

            <div class="text-center">
              <div class="relative mb-6">
                <div class="w-20 h-20 bg-gradient-to-br from-green-600 to-green-600 text-white rounded-full flex items-center justify-center mx-auto font-bold text-2xl shadow-lg relative z-10 border-4 border-white">
                  5
                </div>
              </div>
              <h3 class="text-lg font-bold text-gray-900 mb-2">Aktivasi Kemitraan</h3>
              <p class="text-sm text-gray-600">
                Penandatanganan kontrak dan mulai beroperasi
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Support Programs -->
    <section class="py-20 bg-slate-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-16">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
            Program Dukungan untuk Mitra
          </h2>
          <p class="text-xl text-gray-600">
            Berbagai program yang kami sediakan untuk kesuksesan budidaya Anda
          </p>
        </div>
        <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
          <div class="bg-white p-8 rounded-xl shadow-sm hover:shadow-md transition-shadow border border-slate-100">
            <div class="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center mb-6">
              <BookOpen class="w-6 h-6 text-green-600" />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">Pelatihan Berkala</h3>
            <p class="text-gray-600 mb-4 leading-relaxed">
              Pelatihan rutin setiap 3 bulan untuk update teknik budidaya terbaru dan best practices
            </p>
            <ul class="space-y-2 text-sm text-gray-600">
              <li class="flex items-center gap-2"><CheckCircle class="w-4 h-4 text-green-600" /><span>Workshop praktek lapangan</span></li>
              <li class="flex items-center gap-2"><CheckCircle class="w-4 h-4 text-green-600" /><span>Webinar dengan expert</span></li>
              <li class="flex items-center gap-2"><CheckCircle class="w-4 h-4 text-green-600" /><span>Materi digital lengkap</span></li>
            </ul>
          </div>

          <div class="bg-white p-8 rounded-xl shadow-sm hover:shadow-md transition-shadow border border-slate-100">
            <div class="w-12 h-12 bg-emerald-100 rounded-lg flex items-center justify-center mb-6">
              <Users class="w-6 h-6 text-emerald-700" />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">Pendampingan Teknis</h3>
            <p class="text-gray-600 mb-4 leading-relaxed">
              Tim field officer yang siap membantu menyelesaikan masalah budidaya di lapangan
            </p>
            <ul class="space-y-2 text-sm text-gray-600">
              <li class="flex items-center gap-2"><CheckCircle class="w-4 h-4 text-emerald-600" /><span>Kunjungan rutin minimal 1x/bulan</span></li>
              <li class="flex items-center gap-2"><CheckCircle class="w-4 h-4 text-emerald-600" /><span>Konsultasi via WhatsApp 24/7</span></li>
              <li class="flex items-center gap-2"><CheckCircle class="w-4 h-4 text-emerald-600" /><span>Emergency support saat masalah</span></li>
            </ul>
          </div>

          <div class="bg-white p-8 rounded-xl shadow-sm hover:shadow-md transition-shadow border border-slate-100">
            <div class="w-12 h-12 bg-teal-100 rounded-lg flex items-center justify-center mb-6">
              <Sprout class="w-6 h-6 text-teal-700" />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">Bantuan Input Produksi</h3>
            <p class="text-gray-600 mb-4 leading-relaxed">
              Akses mudah ke bibit, pupuk, dan perlengkapan budidaya dengan harga khusus mitra
            </p>
            <ul class="space-y-2 text-sm text-gray-600">
              <li class="flex items-center gap-2"><CheckCircle class="w-4 h-4 text-teal-600" /><span>Bibit vanili unggul bersertifikat</span></li>
              <li class="flex items-center gap-2"><CheckCircle class="w-4 h-4 text-teal-600" /><span>Pupuk organik berkualitas</span></li>
              <li class="flex items-center gap-2"><CheckCircle class="w-4 h-4 text-teal-600" /><span>Sistem pembayaran fleksibel</span></li>
            </ul>
          </div>

          <div class="bg-white p-8 rounded-xl shadow-sm hover:shadow-md transition-shadow border border-slate-100">
            <div class="w-12 h-12 bg-lime-100 rounded-lg flex items-center justify-center mb-6">
              <DollarSign class="w-6 h-6 text-lime-700" />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">Jaminan Pembelian</h3>
            <p class="text-gray-600 mb-4 leading-relaxed">
              Kami berkomitmen membeli seluruh hasil panen yang memenuhi standar kualitas
            </p>
            <ul class="space-y-2 text-sm text-gray-600">
              <li class="flex items-center gap-2"><CheckCircle class="w-4 h-4 text-lime-600" /><span>Harga kompetitif dan transparan</span></li>
              <li class="flex items-center gap-2"><CheckCircle class="w-4 h-4 text-lime-600" /><span>Pembayaran cepat maksimal 7 hari</span></li>
              <li class="flex items-center gap-2"><CheckCircle class="w-4 h-4 text-lime-600" /><span>Bonus untuk kualitas premium</span></li>
            </ul>
          </div>

          <div class="bg-white p-8 rounded-xl shadow-sm hover:shadow-md transition-shadow border border-slate-100">
            <div class="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center mb-6">
              <Shield class="w-6 h-6 text-green-600" />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">Asuransi & Mitigasi</h3>
            <p class="text-gray-600 mb-4 leading-relaxed">
              Dukungan mitigasi jika terjadi masalah alam atau penyakit pada tanaman
            </p>
            <ul class="space-y-2 text-sm text-gray-600">
              <li class="flex items-center gap-2"><CheckCircle class="w-4 h-4 text-green-600" /><span>Pendampingan recovery tanaman</span></li>
              <li class="flex items-center gap-2"><CheckCircle class="w-4 h-4 text-green-600" /><span>Bantuan bibit pengganti</span></li>
              <li class="flex items-center gap-2"><CheckCircle class="w-4 h-4 text-green-600" /><span>Analisa masalah pasca panen</span></li>
            </ul>
          </div>

          <div class="bg-white p-8 rounded-xl shadow-sm hover:shadow-md transition-shadow border border-slate-100">
            <div class="w-12 h-12 bg-emerald-100 rounded-lg flex items-center justify-center mb-6">
              <Award class="w-6 h-6 text-emerald-700" />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">Program Reward</h3>
            <p class="text-gray-600 mb-4 leading-relaxed">
              Apresiasi untuk mitra berprestasi dengan berbagai bentuk penghargaan
            </p>
            <ul class="space-y-2 text-sm text-gray-600">
              <li class="flex items-center gap-2"><CheckCircle class="w-4 h-4 text-emerald-600" /><span>Bonus produktivitas tahunan</span></li>
              <li class="flex items-center gap-2"><CheckCircle class="w-4 h-4 text-emerald-600" /><span>Apresiasi mitra terbaik</span></li>
              <li class="flex items-center gap-2"><CheckCircle class="w-4 h-4 text-emerald-600" /><span>Kesempatan studi banding</span></li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <!-- FAQ -->
    <!-- <section class="py-20 bg-slate-50">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <h2 class="text-3xl font-bold text-gray-900 mb-12 text-center">
          Pertanyaan Umum Seputar Kemitraan
        </h2>
        <div class="space-y-6">
          <div class="bg-white p-6 rounded-xl shadow-sm border border-slate-100">
            <h3 class="text-lg font-bold text-gray-900 mb-2">
              Berapa minimal luas lahan untuk menjadi mitra?
            </h3>
            <p class="text-gray-600 leading-relaxed">
              Minimal 0.5 hektar (5000 m²). Namun kami juga membuka kesempatan untuk lahan yang lebih kecil dengan sistem kelompok tani.
            </p>
          </div>
          <div class="bg-white p-6 rounded-xl shadow-sm border border-slate-100">
            <h3 class="text-lg font-bold text-gray-900 mb-2">
              Apakah ada biaya untuk bergabung sebagai mitra?
            </h3>
            <p class="text-gray-600 leading-relaxed">
              Tidak ada biaya pendaftaran. Anda hanya perlu menyiapkan lahan dan komitmen untuk mengikuti program pelatihan serta menerapkan SOP budidaya.
            </p>
          </div>
          <div class="bg-white p-6 rounded-xl shadow-sm border border-slate-100">
            <h3 class="text-lg font-bold text-gray-900 mb-2">
              Berapa lama waktu yang dibutuhkan hingga panen pertama?
            </h3>
            <p class="text-gray-600 leading-relaxed">
              Vanili umumnya mulai berbunga di tahun ke-2 sampai ke-3. Panen optimal biasanya dimulai pada tahun ke-3 hingga ke-4 setelah tanam.
            </p>
          </div>
          <div class="bg-white p-6 rounded-xl shadow-sm border border-slate-100">
            <h3 class="text-lg font-bold text-gray-900 mb-2">
              Bagaimana sistem pembayaran hasil panen?
            </h3>
            <p class="text-gray-600 leading-relaxed">
              Kami membeli hasil panen sesuai grade kualitas dengan harga kompetitif. Pembayaran dilakukan maksimal 7 hari kerja setelah penerimaan dan verifikasi hasil panen.
            </p>
          </div>
          <div class="bg-white p-6 rounded-xl shadow-sm border border-slate-100">
            <h3 class="text-lg font-bold text-gray-900 mb-2">
              Apakah saya boleh menjual hasil panen ke pihak lain?
            </h3>
            <p class="text-gray-600 leading-relaxed">
              Sesuai kontrak kemitraan, hasil panen yang memenuhi standar kualitas diprioritaskan untuk dibeli oleh Pondok Tani Land dengan harga terbaik. Namun untuk hasil di luar standar, mitra bebas menjual ke pihak lain.
            </p>
          </div>
        </div>
      </div>
    </section> -->

    <!-- CTA Section -->
    <section class="py-24 bg-green-600 text-white relative overflow-hidden">
      <!-- Decorative background elements -->
      <div class="absolute inset-0 opacity-20 bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] from-white via-transparent to-transparent"></div>
      
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center relative z-10">
        <h2 class="text-3xl md:text-4xl font-bold mb-6 drop-shadow-sm">
          Siap Bergabung Menjadi Mitra Pondok Tani Land?
        </h2>
        <p class="text-xl text-green-100 mb-10 leading-relaxed max-w-3xl mx-auto">
          Bergabunglah dengan ratusan petani yang telah merasakan manfaat kemitraan bersama kami.
          Dapatkan dukungan penuh dari pelatihan hingga jaminan pembelian hasil panen.
        </p>
        <div class="flex flex-wrap justify-center gap-4">
          <button 
            @click="router.push('/login?mode=register&role=mitra')"
            class="bg-white text-green-600 px-10 py-4 rounded-xl font-bold text-lg hover:bg-green-50 transition-all shadow-xl flex items-center gap-2 hover:-translate-y-0.5"
          >
            Daftar Jadi Mitra Sekarang
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
