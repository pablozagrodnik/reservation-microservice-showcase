<script setup lang="ts">
import { computed, onUnmounted, reactive, ref } from 'vue';
import { createRoom, type CreateRoomPayload, type CreateRoomResponse } from '@/api/admin';
import { ApiError } from '@/api/client';

interface RoomForm {
  name: string;
  rows: number | null;
  cols: number | null;
}

const emptyForm = (): RoomForm => ({ name: '', rows: null, cols: null });

const form = reactive<RoomForm>(emptyForm());
const isSubmitting = ref<boolean>(false);
const successMessage = ref<string | null>(null);
const errorMessage = ref<string | null>(null);

let activeController: AbortController | null = null;

const isPositiveInt = (value: number | null): value is number =>
    value !== null && Number.isInteger(value) && value > 0;

const canSubmit = computed<boolean>(
    () =>
        form.name.trim().length > 0 &&
        isPositiveInt(form.rows) &&
        isPositiveInt(form.cols) &&
        !isSubmitting.value
);

const resetForm = (): void => {
  form.name = '';
  form.rows = null;
  form.cols = null;
};

const handleSubmit = async (): Promise<void> => {
  if (!canSubmit.value) return;
  if (!isPositiveInt(form.rows) || !isPositiveInt(form.cols)) return;

  activeController?.abort();
  const controller = new AbortController();
  activeController = controller;

  successMessage.value = null;
  errorMessage.value = null;
  isSubmitting.value = true;

  const payload: CreateRoomPayload = {
    name: form.name.trim(),
    rows: form.rows,
    cols: form.cols,
  };

  try {
    const result: CreateRoomResponse = await createRoom(payload, controller.signal);
    if (activeController !== controller) return;

    const totalSeats = payload.rows * payload.cols;
    successMessage.value =
        `Utworzono salę "${result.room.name}" (ID ${result.room.id}) z ${totalSeats} miejscami.`;
    resetForm();
  } catch (err) {
    if (activeController !== controller) return;
    if (err instanceof DOMException && err.name === 'AbortError') return;

    errorMessage.value = err instanceof ApiError
        ? `Nie udało się utworzyć sali (kod ${err.status}).`
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
  <section aria-labelledby="admin-rooms-heading">
    <h2 id="admin-rooms-heading" class="view-title">Utwórz nową salę</h2>
    <p class="view-hint">Backend wygeneruje siatkę miejsc na podstawie liczby rzędów i kolumn.</p>

    <form class="admin-form" @submit.prevent="handleSubmit" novalidate>
      <div class="form-group">
        <label for="room-name">Nazwa sali <span aria-hidden="true">*</span></label>
        <input
            id="room-name"
            v-model="form.name"
            type="text"
            required
            :disabled="isSubmitting"
            autocomplete="off"
        />
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="room-rows">Liczba rzędów <span aria-hidden="true">*</span></label>
          <input
              id="room-rows"
              v-model.number="form.rows"
              type="number"
              min="1"
              max="20"
              step="1"
              inputmode="numeric"
              required
              :disabled="isSubmitting"
          />
        </div>

        <div class="form-group">
          <label for="room-cols">Liczba kolumn <span aria-hidden="true">*</span></label>
          <input
              id="room-cols"
              v-model.number="form.cols"
              type="number"
              min="1"
              max="25"
              step="1"
              inputmode="numeric"
              required
              :disabled="isSubmitting"
          />
        </div>
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
        <span v-else>Utwórz salę</span>
      </button>
    </form>
  </section>
</template>

<style scoped>
.view-title {
  margin: 0 0 0.25rem;
  font-size: 1.25rem;
  color: #111827;
}

.view-hint {
  margin: 0 0 1.25rem;
  color: #6b7280;
  font-size: 0.95rem;
}

.admin-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  max-width: 520px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
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

.form-group input {
  padding: 0.6rem 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font: inherit;
  background: white;
}

.form-group input:focus {
  outline: 2px solid #f97316;
  border-color: #f97316;
}

.form-group input:disabled {
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
