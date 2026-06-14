<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import LandingLayout from '../../layouts/LandingLayout.vue';
import { getArticleById } from '../../services/content/articles';

const route = useRoute();
const article = ref<any>(null);
const loading = ref(false);

async function loadArticle() {
  loading.value = true;
  try {
    const res = await getArticleById(route.params.id as string);
    article.value = res.data || null;
  } catch (error) {
    console.error('Failed to load article detail:', error);
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  loadArticle();
});
</script>

<template>
  <LandingLayout>
    <section class="py-16 bg-white">
      <div class="mx-auto max-w-5xl px-6">
        <div v-if="loading" class="rounded-2xl border border-gray-200 bg-gray-50 p-10 text-center text-gray-500">Memuat artikel...</div>
        <article v-else-if="article" class="rounded-3xl border border-gray-200 bg-white p-8 shadow-sm">
          <img v-if="article.image_url" :src="article.image_url" :alt="article.title" class="h-72 w-full rounded-2xl object-cover" />
          <div class="mt-6 flex flex-wrap items-center gap-3 text-sm text-gray-500">
            <span class="rounded-full bg-green-50 px-3 py-1 font-semibold text-green-700">{{ article.category || 'Artikel' }}</span>
            <span>{{ article.author || 'Admin' }}</span>
            <span>{{ new Date(article.published_at || article.created_at).toLocaleDateString('id-ID') }}</span>
          </div>
          <h1 class="mt-4 text-3xl font-bold text-slate-900">{{ article.title }}</h1>
          <p class="mt-4 text-lg text-slate-600">{{ article.excerpt }}</p>
          <div class="prose max-w-none mt-6 text-slate-700 whitespace-pre-line">{{ article.content }}</div>
        </article>
        <div v-else class="rounded-2xl border border-dashed border-gray-300 p-10 text-center text-gray-500">Artikel tidak ditemukan.</div>
      </div>
    </section>
  </LandingLayout>
</template>
