<script setup lang="ts">
import { ref } from 'vue';
import { useMovies } from '@/composables/useMovies';
import type { Movie } from '@/types';
import { BASE_URL } from '@/api/client';

const { movies, isLoading, refresh } = useMovies();
const API_URL = BASE_URL;

const isModalOpen = ref(false);
const isEditing = ref(false);
const isSaving = ref(false);

const formMovie = ref<Partial<Movie>>({ title: '', description: '', poster: '' });

const openAddModal = () => {
  isEditing.value = false;
  formMovie.value = { title: '', description: '', poster: '' };
  isModalOpen.value = true;
};

const openEditModal = (movie: Movie) => {
  isEditing.value = true;
  formMovie.value = { ...movie };
  isModalOpen.value = true;
};

const closeModal = () => { isModalOpen.value = false; };

const saveMovie = async () => {
  isSaving.value = true;
  try {
    const url = isEditing.value ? `${API_URL}/admin/movies/${formMovie.value.id}` : `${API_URL}/admin/movies`;
    const method = isEditing.value ? 'PATCH' : 'POST';

    const response = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(formMovie.value)
    });

    if (!response.ok) throw new Error('Błąd zapisu');
    await refresh();
    closeModal();
  } catch (error) {
    alert('Błąd zapisu filmu.');
  } finally {
    isSaving.value = false;
  }
};
</script>

<template>
  <div class="admin-view">
    <div class="admin-header">
      <h2>Zarządzanie Filmami</h2>
      <button class="btn-primary" @click="openAddModal">+ Dodaj nowy film</button>
    </div>

    <div class="table-container surface">
      <div v-if="isLoading" class="loading">Ładowanie danych...</div>
      <table v-else class="admin-table">
        <thead>
          <tr>
            <th>ID</th><th>Plakat</th><th>Tytuł filmu</th><th>Opis</th><th class="actions-col">Akcje</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="movie in movies" :key="movie.id">
            <td class="id-col">#{{ movie.id }}</td>
            <td class="poster-col">
              <img :src="movie.poster" class="thumb" v-if="movie.poster" />
              <div v-else class="thumb-placeholder">Brak</div>
            </td>
            <td class="title-col"><strong>{{ movie.title }}</strong></td>
            <td class="desc-col">{{ movie.description }}</td>
            <td class="actions-col">
              <button class="btn-action edit" @click="openEditModal(movie)">Edytuj</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content">
        <h3>{{ isEditing ? 'Edytuj film' : 'Dodaj nowy film' }}</h3>
        <form @submit.prevent="saveMovie" class="admin-form">
          <div class="form-group">
            <label>Tytuł filmu</label>
            <input v-model="formMovie.title" type="text" required />
          </div>
          <div class="form-group">
            <label>Plakat (URL)</label>
            <input v-model="formMovie.poster" type="url" />
          </div>
          <div class="form-group">
            <label>Opis</label>
            <textarea v-model="formMovie.description" rows="4" required></textarea>
          </div>
          <div class="modal-actions">
            <button type="button" class="btn-secondary" @click="closeModal">Anuluj</button>
            <button type="submit" class="btn-primary" :disabled="isSaving">Zapisz</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.admin-view {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.admin-header h2 { margin: 0; color: var(--text-main); font-size: 1.8rem; }

.btn-primary {
  background-color: #f97316 !important;
  color: #ffffff !important;
  border: none;
  padding: 12px 24px;
  border-radius: 8px;
  font-weight: 800;
  cursor: pointer;
  transition: 0.2s;
  box-shadow: 0 4px 6px rgba(249, 115, 22, 0.2);
}

.btn-primary:hover:not(:disabled) {
  background-color: #ea580c !important;
  transform: translateY(-1px);
}

.btn-secondary {
  background-color: #f1f5f9 !important;
  color: #0f172a !important;
  border: 1px solid #cbd5e1;
  padding: 12px 24px;
  border-radius: 8px;
  font-weight: 700;
  cursor: pointer;
}

.btn-secondary:hover {
  background-color: #e2e8f0 !important;
}

.table-container {
  background-color: #ffffff;
  border-radius: 12px;
  border: 1px solid var(--border);
  overflow-x: auto;
  box-shadow: 0 4px 6px rgba(0,0,0,0.05);
}

.admin-table { width: 100%; border-collapse: collapse; text-align: left; }
.admin-table th, .admin-table td { padding: 16px; border-bottom: 1px solid var(--border); color: #0f172a; }
.admin-table th { background-color: #f8fafc; color: #64748b; font-size: 0.85rem; text-transform: uppercase; }

.thumb { width: 50px; height: 75px; object-fit: cover; border-radius: 4px; }

.btn-action {
  padding: 6px 12px;
  margin-left: 8px;
  border-radius: 6px;
  font-weight: bold;
  cursor: pointer;
  font-size: 0.85rem;
}
.btn-action.edit { background: #f1f5f9; border: 1px solid #cbd5e1; color: #0f172a; }
.btn-action.delete { background: #fef2f2; border: 1px solid #fecaca; color: #ef4444; }

.modal-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(15, 23, 42, 0.8);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.modal-content {
  width: 90%;
  max-width: 500px;
  padding: 2.5rem;
  background-color: #ffffff !important;
  border-radius: 16px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  border: 2px solid var(--border);
}

.modal-content h3 { margin: 0 0 1.5rem; font-size: 1.6rem; color: #0f172a; border-bottom: 2px solid #f1f5f9; padding-bottom: 0.5rem; }

.form-group { margin-bottom: 1.5rem; }
.form-group label { display: block; font-weight: 700; color: #334155; margin-bottom: 0.5rem; }
.form-group input, .form-group textarea {
  width: 100%; padding: 12px; border: 2px solid #e2e8f0; border-radius: 8px; background: #f8fafc; color: #0f172a; font-size: 1rem;
}
.form-group input:focus { border-color: #f97316; outline: none; background: #fff; }

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid #f1f5f9;
}
</style>