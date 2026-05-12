<script setup lang="ts">
import { computed, onMounted, onUnmounted, reactive, ref } from 'vue';
import type { Movie, Room } from '@/types';
import { fetchMovies } from '@/api/movies';
import {
  createScreening,
  fetchRooms,
  type CreateScreeningPayload,
} from '@/api/admin';
import { ApiError } from '@/api/client';

interface ScreeningForm {
  movieId: number | null;
  roomId: number | null;
  startTime: string;
}

const emptyForm = (): ScreeningForm => ({
  movieId: null,
  roomId: null,
  startTime: '',
});

const movies = ref<Movie[]>([]);
const rooms = ref<Room[]>([]);
const isLoading = ref<boolean>(true);
const loadError = ref<string | null>(null);

const form = reactive<ScreeningForm>(emptyForm());
const isSubmitting = ref<boolean>(false);
const successMessage = ref<string | null>(null);
const errorMessage = ref<string | null>(null);

let initialController: AbortController | null = null;
let submitController: AbortController | null = null;

const canSubmit = computed<boolean>(() =>
    form.movieId !== null &&
    form.roomId !== null &&
    form.startTime.trim().length > 0 &&
    !isSubmitting.value &&
    !isLoading.value
);

const formatDateTime = (iso: string): string => {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return iso;
  return d.toLocaleString('pl-PL', {
    weekday: 'short',
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
};

const loadOptions = async (): Promise<void> => {
  initialController?.abort();
  const controller = new AbortController();
  initialController = controller;

  isLoading.value = true;
  loadError.value = null;

  try {
    const [moviesData, roomsData] = await Promise.all([
      fetchMovies(controller.signal),
      fetchRooms(controller.signal),
    ]);
    if (initialController !== controller) return;
    movies.value = moviesData;
    rooms.value = roomsData;
  } catch (err) {
    if (initialController !== controller) return;
    if (err instanceof DOMException && err.name === 'AbortError') return;

    loadError.value = err instanceof ApiError
        ? `Nie udało się pobrać danych formularza (kod ${err.status}).`
        : 'Nie udało się połączyć z serwerem. Sprawdź połączenie i spróbuj ponownie.';
  } finally {
    if (initialController === controller) {
      isLoading.value = false;
    }
  }
};

const resetForm = (): void => {
  form.movieId = null;
  form.roomId = null;
  form.startTime = '';
};

const handleSubmit = async (): Promise<void> => {
  if (!canSubmit.value || form.movieId === null || form.roomId === null) return;

  const localDate = new Date(form.startTime);
  if (Number.isNaN(localDate.getTime())) {
    errorMessage.value = 'Nieprawidłowa data lub godzina.';
    return;
  }

  submitController?.abort();
  const controller = new AbortController();
  submitController = controller;

  successMessage.value = null;
  errorMessage.value = null;
  isSubmitting.value = true;

  // datetime-local nie ma strefy czasowej — interpretujemy jako czas lokalny
  // i konwertujemy do RFC3339 (UTC), bo tego oczekuje Go time.Time przy decodingu JSON.
  const payload: CreateScreeningPayload = {
    movie_id: form.movieId,
    room_id: form.roomId,
    start_time: localDate.toISOString(),
  };

  const movieTitle = movies.value.find(m => m.id === payload.movie_id)?.title ?? '?';
  const roomName = rooms.value.find(r => r.id === payload.room_id)?.name ?? '?';

  try {
    const created = await createScreening(payload, controller.signal);
    if (submitController !== controller) return;

    successMessage.value =
        `Dodano seans (ID ${created.id}): "${movieTitle}" w sali "${roomName}" na ${formatDateTime(payload.start_time)}.`;
    resetForm();
  } catch (err) {
    if (submitController !== controller) return;
    if (err instanceof DOMException && err.name === 'AbortError') return;

    errorMessage.value = err instanceof ApiError
        ? `Nie udało się dodać seansu (kod ${err.status}).`
        : 'Nie udało się połączyć z serwerem. Sprawdź połączenie i spróbuj ponownie.';
  } finally {
    if (submitController === controller) {
      isSubmitting.value = false;
    }
  }
};

onMounted(() => {
  void loadOptions();
});

onUnmounted(() => {
  initialController?.abort();
  submitController?.abort();
});
</script>

<template>
  <section aria-labelledby="admin-screenings-heading">
    <h2 id="admin-screenings-heading" class="view-title">Dodaj nowy seans</h2>

    <div v-if="isLoading" class="loading-state" aria-busy="true" aria-live="polite">
      Ładowanie danych formularza...
    </div>

    <div v-else-if="loadError" class="error-banner" role="alert">
      <p>{{ loadError }}</p>
      <button type="button" class="retry-btn" @click="loadOptions">
        Spróbuj ponownie
      </button>
    </div>

    <p v-else-if="movies.length === 0 || rooms.length === 0" class="empty-state">
      Brak {{ movies.length === 0 ? 'filmów' : 'sal' }} w systemie. Najpierw dodaj je w odpowiednich zakładkach.
    </p>

    <form
        v-else
        class="admin-form"
        @submit.prevent="handleSubmit"
        novalidate
    >
      <div class="form-group">
        <label for="screening-movie">Film <span aria-hidden="true">*</span></label>
        <select
            id="screening-movie"
            v-model="form.movieId"
            required
            :disabled="isSubmitting"
        >
          <option :value="null" disabled>Wybierz film...</option>
          <option
              v-for="movie in movies"
              :key="movie.id"
              :value="movie.id"
          >
            {{ movie.title }}
          </option>
        </select>
      </div>

      <div class="form-group">
        <label for="screening-room">Sala <span aria-hidden="true">*</span></label>
        <select
            id="screening-room"
            v-model="form.roomId"
            required
            :disabled="isSubmitting"
        >
          <option :value="null" disabled>Wybierz salę...</option>
          <option
              v-for="room in rooms"
              :key="room.id"
              :value="room.id"
          >
            {{ room.name }}
          </option>
        </select>
      </div>

      <div class="form-group">
        <label for="screening-start">Data i godzina rozpoczęcia <span aria-hidden="true">*</span></label>
        <input
            id="screening-start"
            v-model="form.startTime"
            type="datetime-local"
            required
            :disabled="isSubmitting"
        />
      </div>

      <div
          v-if="successMessage"
          class="success-banner"
          role="status"
          aria-live="polite"
      >
        {{ successMessage }}
      </div>

      <div
          v-if="errorMessage"
          class="error-banner"
          role="alert"
          aria-live="assertive"
      >
        {{ errorMessage }}
      </div>

      <button
          type="submit"
          class="submit-btn"
          :disabled="!canSubmit"
          :aria-busy="isSubmitting"
      >
        <span v-if="isSubmitting">Zapisywanie...</span>
        <span v-else>Dodaj seans</span>
      </button>
    </form>
  </section>
</template>

<style scoped>
.view-title {
  margin: 0 0 1.25rem;
  font-size: 1.25rem;
  color: #111827;
}

.admin-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  max-width: 520px;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.form-group label {
  font-weight: 500;
  margin-bottom: 0.4rem;
  color: #374151;
}

.form-group label span {
  color: #dc2626;
}

.form-group input,
.form-group select {
  padding: 0.6rem 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font: inherit;
  background: white;
}

.form-group input:focus,
.form-group select:focus {
  outline: 2px solid #f97316;
  border-color: #f97316;
}

.form-group input:disabled,
.form-group select:disabled {
  background-color: #f3f4f6;
  cursor: not-allowed;
}

.submit-btn {
  align-self: flex-start;
  padding: 0.65rem 1.25rem;
  background-color: #f97316;
  color: white;
  font-weight: 600;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.15s ease;
}

.submit-btn:hover:not(:disabled) {
  background-color: #ea580c;
}

.submit-btn:disabled {
  background-color: #fdba74;
  cursor: not-allowed;
}

.success-banner {
  padding: 0.75rem 1rem;
  border: 1px solid #a7f3d0;
  background-color: #ecfdf5;
  color: #065f46;
  border-radius: 6px;
}

.error-banner {
  padding: 0.75rem 1rem;
  border: 1px solid #fecaca;
  background-color: #fef2f2;
  color: #991b1b;
  border-radius: 6px;
}

.error-banner p {
  margin: 0 0 0.5rem;
}

.retry-btn {
  padding: 0.4rem 0.85rem;
  background-color: white;
  border: 1px solid #fecaca;
  border-radius: 6px;
  color: #991b1b;
  cursor: pointer;
  font: inherit;
}

.retry-btn:hover {
  background-color: #fee2e2;
}

.loading-state {
  padding: 1rem;
  color: #6b7280;
  font-style: italic;
}

.empty-state {
  margin: 0;
  padding: 1rem;
  color: #6b7280;
  font-style: italic;
  background-color: #f9fafb;
  border: 1px dashed #d1d5db;
  border-radius: 6px;
}
</style>
