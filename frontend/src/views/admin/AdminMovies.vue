<script setup lang="ts">
import { computed, onUnmounted, reactive, ref } from 'vue';
import { createMovie, type AdminMovie, type CreateMoviePayload } from '@/api/admin';
import { ApiError } from '@/api/client';

interface MovieForm {
  title: string;
  description: string;
  poster: string;
}

const emptyForm = (): MovieForm => ({ title: '', description: '', poster: '' });

const form = reactive<MovieForm>(emptyForm());
const isSubmitting = ref<boolean>(false);
const successMessage = ref<string | null>(null);
const errorMessage = ref<string | null>(null);

let activeController: AbortController | null = null;

const canSubmit = computed<boolean>(
    () => form.title.trim().length > 0 && !isSubmitting.value
);

const resetForm = (): void => {
  form.title = '';
  form.description = '';
  form.poster = '';
};

const handleSubmit = async (): Promise<void> => {
  if (!canSubmit.value) return;

  activeController?.abort();
  const controller = new AbortController();
  activeController = controller;

  successMessage.value = null;
  errorMessage.value = null;
  isSubmitting.value = true;

  const payload: CreateMoviePayload = {
    title: form.title.trim(),
    description: form.description.trim(),
    poster: form.poster.trim(),
  };

  try {
    const movie: AdminMovie = await createMovie(payload, controller.signal);
    if (activeController !== controller) return;

    successMessage.value = `Dodano film "${movie.title}" (ID ${movie.id}).`;
    resetForm();
  } catch (err) {
    if (activeController !== controller) return;
    if (err instanceof DOMException && err.name === 'AbortError') return;

    errorMessage.value = err instanceof ApiError
        ? `Nie udało się dodać filmu (kod ${err.status}).`
        : 'Nie udało się połączyć z serwerem. Sprawdź połączenie i spróbuj ponownie.';
  } finally {
    if (activeController === controller) {
      isSubmitting.value = false;
    }
  }
};

onUnmounted(() => {
  activeController?.abort();
});
</script>

<template>
  <section aria-labelledby="admin-movies-heading">
    <h2 id="admin-movies-heading" class="view-title">Dodaj nowy film</h2>

    <form class="admin-form" @submit.prevent="handleSubmit" novalidate>
      <div class="form-group">
        <label for="movie-title">Tytuł <span aria-hidden="true">*</span></label>
        <input
            id="movie-title"
            v-model="form.title"
            type="text"
            required
            :disabled="isSubmitting"
            autocomplete="off"
        />
      </div>

      <div class="form-group">
        <label for="movie-description">Opis</label>
        <textarea
            id="movie-description"
            v-model="form.description"
            rows="4"
            :disabled="isSubmitting"
        />
      </div>

      <div class="form-group">
        <label for="movie-poster">URL plakatu</label>
        <input
            id="movie-poster"
            v-model="form.poster"
            type="url"
            placeholder="https://..."
            :disabled="isSubmitting"
            autocomplete="off"
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
        <span v-else>Dodaj film</span>
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
.form-group textarea {
  padding: 0.6rem 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font: inherit;
  background: white;
  resize: vertical;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: 2px solid #f97316;
  border-color: #f97316;
}

.form-group input:disabled,
.form-group textarea:disabled {
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
</style>
