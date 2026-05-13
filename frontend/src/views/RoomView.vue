<script setup lang="ts">
import { computed, toRef } from 'vue';
import { useRouter } from 'vue-router';
import type { Seat } from '@/types';
import { useSeats } from '@/composables/useSeats';
import { useReservationStore } from '@/stores/useReservationStore';

const props = defineProps<{ id: string }>();

const router = useRouter();
const store = useReservationStore();
const { seats, isLoading, error, refresh } = useSeats(toRef(props, 'id'));

const selectedIds = computed<Set<number>>(
    () => new Set(store.selectedSeats.map(s => s.id))
);

const seatsByRow = computed<[number, Seat[]][]>(() => {
  const rows = new Map<number, Seat[]>();
  seats.value.forEach(seat => {
    if (!rows.has(seat.row)) rows.set(seat.row, []);
    rows.get(seat.row)!.push(seat);
  });
  return Array.from(rows.entries()).sort((a, b) => a[0] - b[0]);
});

const hasScreeningInStore = computed<boolean>(() => store.screening !== null);
const canProceed = computed<boolean>(
    () => store.selectedSeats.length > 0 && hasScreeningInStore.value
);

const onSeatClick = (seat: Seat): void => {
  store.toggleSeat(seat);
};

const getSeatClass = (seat: Seat): string => {
  if (seat.is_taken) return 'seat-taken';
  if (selectedIds.value.has(seat.id)) return 'seat-selected';
  return 'seat-available';
};

const getSeatAriaLabel = (seat: Seat): string => {
  const status = seat.is_taken
      ? 'Zajęte'
      : selectedIds.value.has(seat.id) ? 'Wybrane' : 'Wolne';
  return `Rząd ${seat.row}, miejsce ${seat.col}. Status: ${status}`;
};

const goToCheckout = (): void => {
  if (!canProceed.value) return;
  void router.push({ name: 'checkout' });
};
</script>

<template>
  <section class="seat-picker" aria-labelledby="seat-picker-heading">
    <h2 id="seat-picker-heading">Wybierz miejsca</h2>

    <div class="screen-indicator" aria-hidden="true">EKRAN</div>

    <div v-if="isLoading" class="skeleton-seats" aria-busy="true" aria-live="polite">
      <div v-for="i in 12" :key="i" class="seat-skeleton"></div>
    </div>

    <div v-else-if="error" class="error-state" role="alert">
      <p class="error-message">{{ error }}</p>
      <button type="button" class="retry-btn" @click="refresh">
        Spróbuj ponownie
      </button>
    </div>

    <p v-else-if="seats.length === 0" class="empty-state">
      Brak danych o miejscach dla tego seansu.
    </p>

    <div v-else class="cinema-room">
      <div
          v-for="[rowIndex, rowSeats] in seatsByRow"
          :key="rowIndex"
          class="seat-row"
      >
        <span class="row-label" aria-hidden="true">{{ rowIndex }}</span>

        <button
            v-for="seat in rowSeats"
            :key="seat.id"
            type="button"
            class="seat-button"
            :class="getSeatClass(seat)"
            :disabled="seat.is_taken"
            :aria-label="getSeatAriaLabel(seat)"
            :aria-pressed="selectedIds.has(seat.id)"
            @click="onSeatClick(seat)"
        >
          <span class="sr-only">{{ seat.col }}</span>
        </button>
      </div>
    </div>

    <div class="legend" aria-hidden="true">
      <div class="legend-item"><div class="seat-box seat-available"></div> Wolne</div>
      <div class="legend-item"><div class="seat-box seat-selected"></div> Wybrane</div>
      <div class="legend-item"><div class="seat-box seat-taken"></div> Zajęte</div>
    </div>

    <div class="selection-summary">
      <p
          v-if="store.selectedSeats.length === 0"
          class="validation-message"
      >
        Proszę wybrać co najmniej jedno miejsce.
      </p>
      <div v-else class="summary-actions">
        <p>
          Wybrano: <strong>{{ store.selectedSeats.length }}</strong> miejsc(a)
        </p>
        <p
            v-if="!hasScreeningInStore"
            class="validation-message"
        >
          Wróć do repertuaru i wybierz seans, aby kontynuować.
        </p>
        <button
            type="button"
            class="primary-btn"
            :disabled="!canProceed"
            @click="goToCheckout"
        >
          Przejdź do płatności
        </button>
      </div>
    </div>
  </section>
</template>

<style scoped>
.seat-picker {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 2rem;
}

.screen-indicator {
  width: 80%;
  max-width: 400px;
  text-align: center;
  background: linear-gradient(to bottom, #d1d5db, #f3f4f6);
  padding: 0.5rem;
  border-radius: 12px 12px 0 0;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  margin-bottom: 3rem;
  font-weight: bold;
  letter-spacing: 2px;
  color: #4b5563;
}

.cinema-room {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  width: 100%;
  overflow-x: auto;
  padding-bottom: 1rem;
}

.seat-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  min-width: max-content;
}

.row-label {
  width: 20px;
  text-align: right;
  margin-right: 10px;
  font-weight: bold;
  color: #6b7280;
}

.seat-button, .seat-box {
  width: 40px;
  height: 40px;
  border-radius: 8px 8px 4px 4px;
  border: 2px solid transparent;
  cursor: pointer;
  transition: transform 0.1s ease, border-color 0.2s ease;
}

.seat-button:focus-visible {
  outline: 3px solid #3b82f6;
  outline-offset: 2px;
}

.seat-available {
  background-color: #e5e7eb;
  border-color: #d1d5db;
}

.seat-available:hover {
  background-color: #d1d5db;
}

.seat-taken {
  background-color: #ef4444;
  border-color: #dc2626;
  cursor: not-allowed;
  opacity: 0.7;
}

.seat-selected {
  background-color: #f97316;
  border-color: #ea580c;
  transform: scale(1.1);
}

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  border: 0;
}

.legend {
  display: flex;
  gap: 1.5rem;
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid #e5e7eb;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
}

.seat-box {
  width: 20px;
  height: 20px;
  cursor: default;
}

.selection-summary {
  margin-top: 2rem;
  min-height: 60px;
}

.validation-message {
  color: #6b7280;
  font-style: italic;
}

.summary-actions {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.primary-btn {
  background-color: #f97316;
  color: white;
  padding: 0.75rem 2rem;
  border: none;
  border-radius: 8px;
  font-weight: bold;
  cursor: pointer;
  font-size: 1.1rem;
}

.primary-btn:disabled {
  background-color: #fdba74;
  cursor: not-allowed;
}

/* Skeleton seats */
.skeleton-seats {
  display: grid;
  grid-template-columns: repeat(6, 40px);
  gap: 0.5rem;
}

.seat-skeleton {
  width: 40px;
  height: 40px;
  background-color: #e5e7eb;
  border-radius: 8px 8px 4px 4px;
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* Error & empty */
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  border: 1px solid #fecaca;
  border-radius: 8px;
  background-color: #fef2f2;
}

.error-message {
  margin: 0;
  color: #991b1b;
}

.retry-btn {
  padding: 0.5rem 1rem;
  cursor: pointer;
  border: 1px solid #991b1b;
  border-radius: 4px;
  background-color: #ffffff;
  color: #991b1b;
  font: inherit;
}

.retry-btn:hover {
  background-color: #991b1b;
  color: #ffffff;
}

.empty-state {
  color: #6b7280;
  font-style: italic;
}
</style>
