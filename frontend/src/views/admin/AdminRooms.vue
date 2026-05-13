<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { BASE_URL } from '@/api/client';

interface Room {
  id: number;
  name: string;
  rows?: number;
  cols?: number;
}

const API_URL = BASE_URL;

const rooms = ref<Room[]>([]);
const isLoading = ref(false);

const isModalOpen = ref(false);
const isEditing = ref(false);
const isSaving = ref(false);

const formRoom = ref<Partial<Room>>({
  name: '',
  rows: 10,
  cols: 10
});

const fetchRooms = async () => {
  isLoading.value = true;
  try {
    const res = await fetch(`${API_URL}/admin/rooms`);
    if (!res.ok) throw new Error();
    rooms.value = await res.json();
  } catch (e) {
    console.error(e);
  } finally {
    isLoading.value = false;
  }
};

const openAddModal = () => {
  isEditing.value = false;
  formRoom.value = { name: '', rows: 10, cols: 10 };
  isModalOpen.value = true;
};

const openEditModal = (room: Room) => {
  isEditing.value = true;
  formRoom.value = { ...room, rows: 10, cols: 10 };
  isModalOpen.value = true;
};

const closeModal = () => {
  isModalOpen.value = false;
};

const saveRoom = async () => {
  isSaving.value = true;
  try {
    const url = isEditing.value ? `${API_URL}/admin/rooms/${formRoom.value.id}` : `${API_URL}/admin/rooms`;
    const method = isEditing.value ? 'PATCH' : 'POST';

    const payload = {
      name: formRoom.value.name,
      rows: Number(formRoom.value.rows),
      cols: Number(formRoom.value.cols)
    };

    const response = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    });

    if (!response.ok) throw new Error();
    await fetchRooms();
    closeModal();
  } catch (error) {
    alert('Wystąpił błąd podczas zapisywania sali.');
  } finally {
    isSaving.value = false;
  }
};

onMounted(fetchRooms);
</script>

<template>
  <div class="admin-view">
    <div class="admin-header">
      <h2>Zarządzanie Salami</h2>
      <button class="btn-primary" @click="openAddModal">+ Dodaj nową salę</button>
    </div>

    <div class="table-container surface">
      <div v-if="isLoading" class="loading">Ładowanie danych...</div>

      <table v-else class="admin-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Nazwa sali</th>
            <th class="actions-col">Akcje</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="room in rooms" :key="room.id">
            <td class="id-col">#{{ room.id }}</td>
            <td class="title-col"><strong>{{ room.name }}</strong></td>
            <td class="actions-col">
              <button class="btn-action edit" @click="openEditModal(room)">Edytuj</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content">
        <h3>{{ isEditing ? 'Edytuj salę' : 'Dodaj nową salę' }}</h3>

        <form @submit.prevent="saveRoom" class="admin-form">
          <div class="form-group">
            <label>Nazwa sali</label>
            <input v-model="formRoom.name" type="text" required />
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>Liczba rzędów</label>
              <input v-model="formRoom.rows" type="number" min="1" max="50" required />
            </div>

            <div class="form-group">
              <label>Liczba kolumn</label>
              <input v-model="formRoom.cols" type="number" min="1" max="50" required />
            </div>
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
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
.form-group { margin-bottom: 1.5rem; }
.form-group label { display: block; font-weight: 700; color: #334155; margin-bottom: 0.5rem; }
.form-group input { width: 100%; padding: 12px; border: 2px solid #e2e8f0; border-radius: 8px; background: #f8fafc; color: #0f172a; font-size: 1rem; }
.form-group input:focus { border-color: #f97316; outline: none; background: #fff; }
.modal-actions { display: flex; justify-content: flex-end; gap: 1rem; margin-top: 2rem; padding-top: 1rem; border-top: 1px solid #f1f5f9; }
.loading { padding: 3rem; text-align: center; color: var(--text-muted); }
</style>