<script setup lang="ts">
import { Sprout } from 'lucide-vue-next';
</script>

<template>
  <footer class="bg-gray-900 text-gray-300 py-12 mt-auto">
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
</template>
