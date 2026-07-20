<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { Plus, Pencil, Trash2, Search, Eye, Image as ImageIcon, Loader2 } from 'lucide-vue-next';
import { createArticle, deleteArticle, getPublishedArticles, updateArticle } from '../../services/content/articles';

const articles = ref<any[]>([]);
const loading = ref(false);
const searchQuery = ref('');
const modalOpen = ref(false);
const editingId = ref<number | null>(null);
const submitting = ref(false);
const error = ref('');
const success = ref('');

const form = reactive({
  title: '',
  category: 'Umum',
  excerpt: '',
  content: '',
  status: 'draft',
  author: 'Admin',
  image: null as File | null,
});

const filteredArticles = computed(() => {
  const query = searchQuery.value.toLowerCase();
  if (!query) return articles.value;
  return articles.value.filter((item) =>
    [item.title, item.category, item.author, item.status].join(' ').toLowerCase().includes(query)
  );
});

function openCreateModal() {
  editingId.value = null;
  form.title = '';
  form.category = 'Umum';
  form.excerpt = '';
  form.content = '';
  form.status = 'draft';
  form.author = 'Admin';
  form.image = null;
  error.value = '';
  success.value = '';
  modalOpen.value = true;
}

function openEditModal(item: any) {
  editingId.value = item.id;
  form.title = item.title || '';
  form.category = item.category || 'Umum';
  form.excerpt = item.excerpt || '';
  form.content = item.content || '';
  form.status = item.status || 'draft';
  form.author = item.author || 'Admin';
  form.image = null;
  error.value = '';
  success.value = '';
  modalOpen.value = true;
}

function closeModal() {
  modalOpen.value = false;
}

async function loadArticles() {
  loading.value = true;
  try {
    const res = await getPublishedArticles({ status: '' });
    articles.value = Array.isArray(res.data) ? res.data : [];
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
}

async function submitArticle() {
  submitting.value = true;
  error.value = '';
  success.value = '';
  try {
    const payload = new FormData();
    payload.append('title', form.title);
    payload.append('category', form.category);
    payload.append('excerpt', form.excerpt);
    payload.append('content', form.content);
    payload.append('status', form.status);
    payload.append('author', form.author);
    if (form.image) payload.append('image', form.image);

    if (editingId.value) {
      await updateArticle(editingId.value, payload);
      success.value = 'Artikel berhasil diperbarui.';
    } else {
      await createArticle(payload);
      success.value = 'Artikel berhasil ditambahkan.';
    }
    await loadArticles();
    closeModal();
  } catch (e: any) {
    error.value = e?.response?.data?.error || 'Gagal menyimpan artikel.';
  } finally {
    submitting.value = false;
  }
}

async function removeArticle(item: any) {
  if (!confirm(`Hapus artikel “${item.title}”?`)) return;
  try {
    await deleteArticle(item.id);
    await loadArticles();
  } catch (e) {
    console.error(e);
  }
}

onMounted(() => {
  loadArticles();
});
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-semibold text-gray-900">Manajemen Artikel</h2>
        <p class="text-gray-600 mt-1">Tambah, edit, dan hapus artikel untuk website.</p>
      </div>
      <button class="inline-flex items-center gap-2 rounded-xl bg-green-600 px-4 py-2.5 text-white hover:bg-green-700" @click="openCreateModal">
        <Plus class="w-4 h-4" /> Tambah Artikel
      </button>
    </div>

    <div class="rounded-2xl border border-gray-200 bg-white p-4 shadow-sm">
      <div class="flex items-center gap-2 rounded-xl border border-gray-200 bg-gray-50 px-3 py-2">
        <Search class="w-4 h-4 text-gray-400" />
        <input v-model="searchQuery" type="text" placeholder="Cari artikel, kategori, atau status" class="w-full bg-transparent outline-none text-sm" />
      </div>
    </div>

    <div v-if="loading" class="rounded-2xl border border-gray-200 bg-white p-10 text-center text-gray-500">Memuat artikel...</div>
    <div v-else class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
      <article v-for="item in filteredArticles" :key="item.id" class="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm">
        <img v-if="item.image_url" :src="item.image_url" alt="" class="h-40 w-full rounded-xl object-cover" />
        <div v-else class="h-40 rounded-xl bg-gray-100 flex items-center justify-center text-gray-400"><ImageIcon class="w-8 h-8" /></div>
        <div class="mt-4 space-y-2">
          <div class="flex items-center justify-between gap-2">
            <span class="rounded-full bg-green-50 px-2.5 py-1 text-xs font-semibold text-green-700">{{ item.category || 'Umum' }}</span>
            <span class="text-xs text-gray-500">{{ item.status === 'publish' ? 'Publish' : 'Draft' }}</span>
          </div>
          <h3 class="text-lg font-semibold text-gray-900">{{ item.title }}</h3>
          <p class="text-sm text-gray-600 line-clamp-3">{{ item.excerpt || item.content }}</p>
          <div class="flex items-center justify-between pt-2 text-xs text-gray-500">
            <span>{{ item.author || 'Admin' }}</span>
            <span>{{ new Date(item.created_at).toLocaleDateString('id-ID') }}</span>
          </div>
        </div>
        <div class="mt-4 flex items-center gap-2">
          <button class="inline-flex items-center gap-1 rounded-lg border border-gray-200 px-3 py-2 text-sm text-gray-700 hover:bg-gray-50" @click="openEditModal(item)"><Pencil class="w-4 h-4" /> Edit</button>
          <button class="inline-flex items-center gap-1 rounded-lg border border-red-200 px-3 py-2 text-sm text-red-600 hover:bg-red-50" @click="removeArticle(item)"><Trash2 class="w-4 h-4" /> Hapus</button>
        </div>
      </article>
    </div>

    <div v-if="modalOpen" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
      <div class="w-full max-w-3xl rounded-3xl bg-white p-6 shadow-2xl">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-xl font-semibold text-gray-900">{{ editingId ? 'Edit Artikel' : 'Tambah Artikel' }}</h3>
            <p class="text-sm text-gray-500">Isi data artikel, status publikasi, dan gambar pendukung.</p>
          </div>
          <button class="rounded-xl border border-gray-200 p-2 text-gray-500 hover:bg-gray-50" @click="closeModal"><span class="text-lg">×</span></button>
        </div>
        <div class="mt-6 grid gap-4 md:grid-cols-2">
          <label class="space-y-1 md:col-span-2"><span class="text-sm font-medium text-gray-700">Judul</span><input v-model="form.title" class="w-full rounded-xl border border-gray-300 px-3 py-2.5" /></label>
          <label class="space-y-1"><span class="text-sm font-medium text-gray-700">Kategori</span><input v-model="form.category" class="w-full rounded-xl border border-gray-300 px-3 py-2.5" /></label>
          <label class="space-y-1"><span class="text-sm font-medium text-gray-700">Penulis</span><input v-model="form.author" class="w-full rounded-xl border border-gray-300 px-3 py-2.5" /></label>
          <label class="space-y-1"><span class="text-sm font-medium text-gray-700">Status</span><select v-model="form.status" class="w-full rounded-xl border border-gray-300 px-3 py-2.5"><option value="draft">Draft</option><option value="publish">Publish</option></select></label>
          <label class="space-y-1 md:col-span-2"><span class="text-sm font-medium text-gray-700">Ringkasan</span><textarea v-model="form.excerpt" rows="3" class="w-full rounded-xl border border-gray-300 px-3 py-2.5"></textarea></label>
          <label class="space-y-1 md:col-span-2"><span class="text-sm font-medium text-gray-700">Isi Artikel</span><textarea v-model="form.content" rows="6" class="w-full rounded-xl border border-gray-300 px-3 py-2.5"></textarea></label>
          <label class="space-y-1 md:col-span-2"><span class="text-sm font-medium text-gray-700">Gambar Artikel</span><input type="file" accept="image/*" @change="(e:any) => form.image = e.target.files?.[0] || null" class="w-full rounded-xl border border-gray-300 px-3 py-2.5" /></label>
        </div>
        <div v-if="error" class="mt-4 rounded-xl bg-red-50 px-3 py-2 text-sm text-red-600">{{ error }}</div>
        <div v-if="success" class="mt-4 rounded-xl bg-green-50 px-3 py-2 text-sm text-green-600">{{ success }}</div>
        <div class="mt-6 flex justify-end gap-3">
          <button class="rounded-xl border border-gray-200 px-4 py-2 text-gray-700 hover:bg-gray-50" @click="closeModal">Batal</button>
          <button class="inline-flex items-center gap-2 rounded-xl bg-green-600 px-4 py-2 text-white hover:bg-green-700 disabled:opacity-60" :disabled="submitting" @click="submitArticle"><Loader2 v-if="submitting" class="w-4 h-4 animate-spin" /> Simpan</button>
        </div>
      </div>
    </div>
  </div>
</template>
