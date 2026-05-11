<script setup lang="ts">
import { computed, ref } from 'vue';
import LandingLayout from '../../layouts/LandingLayout.vue';
import { MapPin, Mail, Instagram, Facebook, Youtube, Send } from 'lucide-vue-next';

const contactEmail = 'info@pondoktaniland.com';
const subject = ref('Konsultasi Investasi Pondok Tani Land');
const name = ref('');
const email = ref('');
const message = ref('');

const canSend = computed(() => name.value.trim() && email.value.trim() && message.value.trim());

function sendConsultationEmail() {
  const body = `Halo Pondok Tani Land,%0D%0A%0D%0ANama: ${encodeURIComponent(name.value)}%0D%0AEmail: ${encodeURIComponent(email.value)}%0D%0A%0D%0AMessage:%0D%0A${encodeURIComponent(message.value)}`;
  window.location.href = `mailto:${contactEmail}?subject=${encodeURIComponent(subject.value)}&body=${body}`;
}
</script>

<template>
  <LandingLayout>
    <section class="mx-auto max-w-7xl px-6 py-14 lg:py-20">
      <div class="mb-12 flex flex-col gap-6 lg:flex-row lg:items-end lg:justify-between">
        <div class="max-w-3xl space-y-4">
          <p class="text-sm font-semibold uppercase tracking-[0.25em] text-green-700">Kontak Kami</p>
          <h1 class="text-4xl font-bold text-slate-900 sm:text-5xl">Kami Siap Membantu Konsultasi Investasi Anda</h1>
          <p class="text-base leading-8 text-slate-600">Isi formulir konsultasi di bawah ini atau hubungi kami melalui sosial media dan lokasi Google Maps.</p>
        </div>
      </div>

      <div class="grid gap-8 lg:grid-cols-[1.4fr_1fr]">
        <div class="rounded-3xl border border-slate-200 bg-white p-8 shadow-sm">
          <div class="space-y-6">
            <div>
              <h2 class="text-2xl font-semibold text-slate-900 mb-2">Form Email Konsultasi</h2>
              <p class="text-sm text-slate-600">Masukkan data Anda dan pertanyaan investasi. Kami akan membalas melalui email.</p>
            </div>

            <div class="grid gap-5">
              <label class="space-y-2 text-sm text-slate-700">
                Nama Lengkap
                <input
                  v-model="name"
                  type="text"
                  placeholder="Contoh: Budi Santoso"
                  class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-slate-900 outline-none transition focus:border-green-500 focus:ring-2 focus:ring-green-100"
                />
              </label>
              <label class="space-y-2 text-sm text-slate-700">
                Email
                <input
                  v-model="email"
                  type="email"
                  placeholder="nama@contoh.com"
                  class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-slate-900 outline-none transition focus:border-green-500 focus:ring-2 focus:ring-green-100"
                />
              </label>
              <label class="space-y-2 text-sm text-slate-700">
                Subjek
                <input
                  v-model="subject"
                  type="text"
                  placeholder="Konsultasi mengenai paket investasi"
                  class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-slate-900 outline-none transition focus:border-green-500 focus:ring-2 focus:ring-green-100"
                />
              </label>
              <label class="space-y-2 text-sm text-slate-700">
                Pesan
                <textarea
                  v-model="message"
                  rows="6"
                  placeholder="Tuliskan pertanyaan, kebutuhan investasi, atau informasi yang ingin Anda dapatkan."
                  class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-slate-900 outline-none transition focus:border-green-500 focus:ring-2 focus:ring-green-100 resize-none"
                ></textarea>
              </label>
            </div>

            <button
              type="button"
              :disabled="!canSend"
              @click="sendConsultationEmail"
              class="inline-flex w-full items-center justify-center gap-3 rounded-2xl bg-green-600 px-6 py-4 text-white shadow-sm transition hover:bg-green-700 disabled:cursor-not-allowed disabled:bg-slate-300"
            >
              <Send class="h-5 w-5" />
              Kirim Konsultasi via Email
            </button>
          </div>
        </div>

        <div class="space-y-8">
          <div class="rounded-3xl border border-slate-200 bg-white p-8 shadow-sm">
            <div class="flex items-center gap-4 mb-6">
              <div class="inline-flex h-12 w-12 items-center justify-center rounded-2xl bg-green-100 text-green-600">
                <Mail class="h-6 w-6" />
              </div>
              <div>
                <h3 class="text-xl font-semibold text-slate-900">Email Konsultasi</h3>
                <p class="text-sm text-slate-600">info@pondoktaniland.com</p>
              </div>
            </div>

            <div class="space-y-5">
              <div class="flex items-start gap-4">
                <div class="inline-flex h-11 w-11 items-center justify-center rounded-2xl bg-green-100 text-green-600">
                  <Instagram class="h-5 w-5" />
                </div>
                <div>
                  <h4 class="font-semibold text-slate-900">Instagram</h4>
                  <p class="text-sm text-slate-600">@pondoktaniland</p>
                </div>
              </div>
              <div class="flex items-start gap-4">
                <div class="inline-flex h-11 w-11 items-center justify-center rounded-2xl bg-green-100 text-green-600">
                  <Facebook class="h-5 w-5" />
                </div>
                <div>
                  <h4 class="font-semibold text-slate-900">Facebook</h4>
                  <p class="text-sm text-slate-600">Pondok Tani Land</p>
                </div>
              </div>
              <div class="flex items-start gap-4">
                <div class="inline-flex h-11 w-11 items-center justify-center rounded-2xl bg-green-100 text-green-600">
                  <Youtube class="h-5 w-5" />
                </div>
                <div>
                  <h4 class="font-semibold text-slate-900">YouTube</h4>
                  <p class="text-sm text-slate-600">Pondok Tani Land</p>
                </div>
              </div>
            </div>
          </div>

          <div class="overflow-hidden rounded-3xl border border-slate-200 bg-white shadow-sm">
            <iframe
              class="h-[360px] w-full border-0"
              loading="lazy"
              allowfullscreen
              src="https://www.google.com/maps?q=-7.5690098,110.1440301&z=17&output=embed"
            ></iframe>
            <div class="p-6">
              <h3 class="text-xl font-semibold text-slate-900">Lokasi Kami</h3>
              <p class="mt-2 text-sm text-slate-600">Pondok Tani Land</p>
              <p class="mt-3 text-sm text-slate-600">Jalan Utama, Yogyakarta, Indonesia</p>
            </div>
          </div>
        </div>
      </div>
    </section>
  </LandingLayout>
</template>
