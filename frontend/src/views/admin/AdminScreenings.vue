<script setup lang="ts">
import { ref, onMounted } from 'vue';
import type { Movie, Room } from '@/types';

interface Screening {
  id: number;
  movie: Movie;
  room: Room;
  start_time: string;
}

const API_URL = 'http://localhost:8080';

const screenings = ref<Screening[]>([]);
const movies = ref<Movie[]>([]);
const rooms = ref<Room[]>([]);
const isLoading = ref(false);

const isModalOpen = ref(false);
const isEditing = ref(false);
const isSaving = ref(false);

const formScreening = ref<{ id?: number; movie_id: number | null; room_id: number | null; start_time: string }>({
  movie_id: null,
  room_id: null,
  start_time: ''
});

const formatDateTime = (iso: string) => {
  return new Date(iso).toLocaleString('pl-PL', {
    weekday: 'short', day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit'
  });
};

const formatForInput = (iso: string) => {
  if (!iso) return '';
  return new Date(iso).toISOString().slice(0, 16);
};

const fetchData = async () => {
  isLoading.value = true;
  try {
    const [scrRes, movRes, roomRes] = await Promise.all([
      fetch(`${API_URL}/admin/screenings`),
      fetch(`${API_URL}/movies`),
      fetch(`${API_URL}/admin/rooms`)
    ]);

    if (!scrRes.ok || !movRes.ok || !roomRes.ok) throw new Error();

    screenings.value = await scrRes.json();
    movies.value = await movRes.json();
    rooms.value = await roomRes.json();
  } catch (e) {
    console.error(e);
  } finally {
    isLoading.value = false;
  }
};

const openAddModal = () => {
  isEditing.value = false;
  formScreening.value = { movie_id: null, room_id: null, start_time: '' };
  isModalOpen.value = true;
};

const openEditModal = (scr: Screening) => {
  isEditing.value = true;
  formScreening.value = {
    id: scr.id,
    movie_id: scr.movie.id,
    room_id: scr.room.id,
    start_time: formatForInput(scr.start_time)
  };
  isModalOpen.value = true;
};

const closeModal = () => {
  isModalOpen.value = false;
};

const saveScreening = async () => {
  isSaving.value = true;
  try {
    const url = isEditing.value ? `${API_URL}/admin/screenings/${formScreening.value.id}` : `${API_URL}/admin/screenings`;
    const method = isEditing.value ? 'PATCH' : 'POST';

    const payload = {
      movie_id: Number(formScreening.value.movie_id),
      room_id: Number(formScreening.value.room_id),
      start_time: new Date(formScreening.value.start_time).toISOString()
    };

    const response = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    });

    if (!response.ok) throw new Error();
    await fetchData();
    closeModal();
  } catch (error) {
    alert('Wystąpił błąd podczas zapisywania seansu.');
  } finally {
    isSaving.value = false;
  }
};

onMounted(fetchData);
</script>

<template>
  <div class="admin-view">
    <div class="admin-header">
      <h2>Zarządzanie Seansami</h2>
      <button class="btn-primary" @click="openAddModal">+ Dodaj nowy seans</button>
    </div>

    <div class="table-container surface">
      <div v-if="isLoading" class="loading">Ładowanie danych...</div>

      <table v-else class="admin-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Film</th>
            <th>Sala</th>
            <th>Data rozpoczęcia</th>
            <th class="actions-col">Akcje</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="scr in screenings" :key="scr.id">
            <td class="id-col">#{{ scr.id }}</td>
            <td class="title-col">
              <strong v-if="scr.movie">{{ scr.movie.title }}</strong>
              <span v-else class="text-muted">-</span>
            </td>
            <td>
              <span v-if="scr.room">{{ scr.room.name }}</span>
              <span v-else class="text-muted">Brak sali</span>
            </td>
            <td>{{ formatDateTime(scr.start_time) }}</td>
            <td class="actions-col">
              <button class="btn-action edit" @click="openEditModal(scr)">Edytuj</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content">
        <h3>{{ isEditing ? 'Edytuj seans' : 'Dodaj nowy seans' }}</h3>

        <form @submit.prevent="saveScreening" class="admin-form">
          <div class="form-group">
            <label>Film</label>
            <select v-model="formScreening.movie_id" required>
              <option :value="null" disabled>Wybierz film...</option>
              <option v-for="movie in movies" :key="movie.id" :value="movie.id">{{ movie.title }}</option>
            </select>
          </div>

          <div class="form-group">
            <label>Sala</label>
            <select v-model="formScreening.room_id" required>
              <option :value="null" disabled>Wybierz salę...</option>
              <option v-for="room in rooms" :key="room.id" :value="room.id">{{ room.name }}</option>
            </select>
          </div>

          <div class="form-group">
            <label>Data i godzina</label>
            <input v-model="formScreening.start_time" type="datetime-local" required />
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
.admin-view { padding: 20px; max-width: 1200px; margin: 0 auto; }
.admin-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.admin-header h2 { margin: 0; color: var(--text-main); font-size: 1.8rem; }
.btn-primary { background-color: #f97316 !important; color: #ffffff !important; border: none; padding: 12px 24px; border-radius: 8px; font-weight: 800; cursor: pointer; transition: 0.2s; box-shadow: 0 4px 6px rgba(249, 115, 22, 0.2); }
.btn-primary:hover:not(:disabled) { background-color: #ea580c !important; transform: translateY(-1px); }
.btn-secondary { background-color: #f1f5f9 !important; color: #0f172a !important; border: 1px solid #cbd5e1; padding: 12px 24px; border-radius: 8px; font-weight: 700; cursor: pointer; }
.btn-secondary:hover { background-color: #e2e8f0 !important; }
.table-container { background-color: #ffffff; border-radius: 12px; border: 1px solid var(--border); overflow-x: auto; box-shadow: 0 4px 6px rgba(0,0,0,0.05); }
.admin-table { width: 100%; border-collapse: collapse; text-align: left; }
.admin-table th, .admin-table td { padding: 16px; border-bottom: 1px solid var(--border); color: #0f172a; }
.admin-table th { background-color: #f8fafc; color: #64748b; font-size: 0.85rem; text-transform: uppercase; }
.admin-table tr:hover { background-color: rgba(249, 115, 22, 0.05); }
.id-col { width: 60px; color: var(--text-muted); font-weight: bold; }
.title-col { font-size: 1.1rem; }
.actions-col { width: 160px; text-align: right; }
.btn-action { padding: 6px 12px; margin-left: 8px; border-radius: 6px; font-weight: bold; cursor: pointer; font-size: 0.85rem; }
.btn-action.edit { background: #f1f5f9; border: 1px solid #cbd5e1; color: #0f172a; }
.btn-action.edit:hover { border-color: #f97316; color: #f97316; }
.modal-overlay { position: fixed; inset: 0; background-color: rgba(15, 23, 42, 0.8); backdrop-filter: blur(8px); display: flex; align-items: center; justify-content: center; z-index: 9999; }
.modal-content { width: 90%; max-width: 500px; padding: 2.5rem; background-color: #ffffff !important; border-radius: 16px; box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5); border: 2px solid var(--border); }
.modal-content h3 { margin: 0 0 1.5rem; font-size: 1.6rem; color: #0f172a; border-bottom: 2px solid #f1f5f9; padding-bottom: 0.5rem; }
.form-group { margin-bottom: 1.5rem; }
.form-group label { display: block; font-weight: 700; color: #334155; margin-bottom: 0.5rem; }
.form-group input, .form-group select { width: 100%; padding: 12px; border: 2px solid #e2e8f0; border-radius: 8px; background: #f8fafc; color: #0f172a; font-size: 1rem; }
.form-group input:focus, .form-group select:focus { border-color: #f97316; outline: none; background: #fff; }
.modal-actions { display: flex; justify-content: flex-end; gap: 1rem; margin-top: 2rem; padding-top: 1rem; border-top: 1px solid #f1f5f9; }
.loading { padding: 3rem; text-align: center; color: var(--text-muted); }
</style>