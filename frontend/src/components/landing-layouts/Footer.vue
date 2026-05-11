<script setup lang="ts">
import { MapPin, Mail, Phone } from 'lucide-vue-next';
</script>

<template>
  <footer class="bg-slate-900 text-slate-300 pt-16 pb-8 mt-auto border-t-[6px] border-green-600">
    <div class="max-w-7xl mx-auto px-6">
      <div class="grid grid-cols-1 md:grid-cols-12 gap-12 mb-16">
        <div class="md:col-span-4">
          <div class="flex items-center gap-3 mb-6">
            <img src="../../assets/logo-pondok-tani.png" alt="Pondok Tani Land" class="w-12 h-12 object-contain rounded-xl bg-white p-1.5 shadow-md" />
            <div>
              <h4 class="text-2xl font-bold text-white tracking-wide">Pondok Tani Land</h4>
            </div>
          </div>
          <p class="text-sm text-slate-400 leading-relaxed mb-6 pr-4">
            Platform manajemen lahan pertanian modern dan terpercaya di Indonesia. Kami menjembatani pengelolaan lahan dengan transparansi penuh dan efisiensi maksimal.
          </p>
        </div>

        <div class="md:col-span-2">
          <h5 class="text-white font-semibold mb-6 uppercase tracking-wider text-sm">Perusahaan</h5>
          <ul class="space-y-4 text-sm">
            <li><a href="#" class="text-slate-400 hover:text-green-400 transition-colors">Tentang Kami</a></li>
            <li><a href="#" class="text-slate-400 hover:text-green-400 transition-colors">Tim & Pakar</a></li>
            <li><a href="#" class="text-slate-400 hover:text-green-400 transition-colors">Karir</a></li>
            <li><a href="#" class="text-slate-400 hover:text-green-400 transition-colors">Blog & Edukasi</a></li>
          </ul>
        </div>

        <div class="md:col-span-2">
          <h5 class="text-white font-semibold mb-6 uppercase tracking-wider text-sm">Layanan</h5>
          <ul class="space-y-4 text-sm">
            <li><router-link to="/paket-investasi" class="text-slate-400 hover:text-green-400 transition-colors">Paket Lahan</router-link></li>
            <li><router-link to="/kemitraan" class="text-slate-400 hover:text-green-400 transition-colors">Program Mitra</router-link></li>
            <li><a href="#" class="text-slate-400 hover:text-green-400 transition-colors">Sistem Monitoring</a></li>
            <li><a href="#" class="text-slate-400 hover:text-green-400 transition-colors">Laporan Hasil</a></li>
          </ul>
        </div>

        <div class="md:col-span-4">
          <h5 class="text-white font-semibold mb-6 uppercase tracking-wider text-sm">Hubungi Kami</h5>
          <ul class="space-y-4 text-sm">
            <li class="flex items-start gap-3">
              <span class="text-green-500 mt-0.5"><MapPin class="w-4 h-4" /></span>
              <span class="text-slate-400 leading-relaxed">Pondok Tani HQ<br>Kawasan Pertanian Terpadu<br>Yogyakarta, Indonesia 55123</span>
            </li>
            <li class="flex items-center gap-3">
              <span class="text-green-500"><Mail class="w-4 h-4" /></span>
              <span class="text-slate-400">halo@pondoktaniland.com</span>
            </li>
            <li class="flex items-center gap-3">
              <span class="text-green-500"><Phone class="w-4 h-4" /></span>
              <span class="text-slate-400">+62 811-2233-4455</span>
            </li>
          </ul>
        </div>
      </div>

      <div class="border-t border-slate-800 pt-8 flex flex-col md:flex-row items-center justify-between gap-4 text-sm text-slate-500">
        <p>&copy; {{ new Date().getFullYear() }} Pondok Tani Land. Seluruh hak cipta dilindungi.</p>
        <div class="flex gap-6">
          <a href="#" class="hover:text-white transition-colors">Syarat & Ketentuan</a>
          <a href="#" class="hover:text-white transition-colors">Kebijakan Privasi</a>
        </div>
      </div>
    </div>
  </footer>
</template>
