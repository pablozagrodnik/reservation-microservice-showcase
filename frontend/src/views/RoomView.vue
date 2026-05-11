<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import type { Seat } from '@/types';

// const props = defineProps<{
//   screeningId: number;
// }>();

defineProps<{
  id: string;
}>();

const seats = ref<Seat[]>([]);
const selectedSeatIds = ref<Set<number>>(new Set());
const isLoading = ref(true);

onMounted(async () => {
  try {
    seats.value = [
      { id: 1, row: 1, col: 1, is_taken: true },
      { id: 2, row: 1, col: 2, is_taken: false },
      { id: 3, row: 1, col: 3, is_taken: false },
      { id: 4, row: 2, col: 1, is_taken: false },
      { id: 5, row: 2, col: 2, is_taken: true },
      { id: 6, row: 2, col: 3, is_taken: false },
    ];
  } catch (error) {
    console.error("Błąd podczas pobierania miejsc:", error);
  } finally {
    isLoading.value = false;
  }
});

const seatsByRow = computed(() => {
  const rows = new Map<number, Seat[]>();
  seats.value.forEach(seat => {
    if (!rows.has(seat.row)) rows.set(seat.row, []);
    rows.get(seat.row)!.push(seat);
  });
  return Array.from(rows.entries()).sort((a, b) => a[0] - b[0]);
});

const toggleSeat = (seat: Seat) => {
  if (seat.is_taken) return;

  const updatedSet = new Set(selectedSeatIds.value);
  if (updatedSet.has(seat.id)) {
    updatedSet.delete(seat.id);
  } else {
    updatedSet.add(seat.id);
  }
  selectedSeatIds.value = updatedSet;
};

const getSeatClass = (seat: Seat) => {
  if (seat.is_taken) return 'seat-taken';
  if (selectedSeatIds.value.has(seat.id)) return 'seat-selected';
  return 'seat-available';
};

const getSeatAriaLabel = (seat: Seat) => {
  const status = seat.is_taken ? 'Zajęte' : selectedSeatIds.value.has(seat.id) ? 'Wybrane' : 'Wolne';
  return `Rząd ${seat.row}, miejsce ${seat.col}. Status: ${status}`;
};
</script>

<template>
  <section class="seat-picker" aria-labelledby="seat-picker-heading">
    <h2 id="seat-picker-heading">Wybierz miejsca</h2>

    <div class="screen-indicator" aria-hidden="true">EKRAN</div>

    <div v-if="isLoading" class="skeleton-seats" aria-busy="true">
      <div v-for="i in 12" :key="i" class="seat-skeleton"></div>
    </div>

    <div v-else class="cinema-room">
      <div v-for="[rowIndex, rowSeats] in seatsByRow" :key="rowIndex" class="seat-row">
        <span class="row-label" aria-hidden="true">{{ rowIndex }}</span>

        <button
            v-for="seat in rowSeats"
            :key="seat.id"
            class="seat-button"
            :class="getSeatClass(seat)"
            :disabled="seat.is_taken"
            :aria-label="getSeatAriaLabel(seat)"
            :aria-pressed="selectedSeatIds.has(seat.id)"
            @click="toggleSeat(seat)"
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
      <p v-if="selectedSeatIds.size === 0" class="validation-message">Proszę wybrać co najmniej jedno miejsce.</p>
      <div v-else class="summary-actions">
        <p>Wybrano: <strong>{{ selectedSeatIds.size }}</strong> miejsc(a)</p>
        <button class="primary-btn">Przejdź do płatności</button>
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
}

.seat-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.row-label {
  width: 20px;
  text-align: right;
  margin-right: 10px;
  font-weight: bold;
  color: #6b7280;
}

/* Base Seat Styles */
.seat-button, .seat-box {
  width: 40px;
  height: 40px;
  border-radius: 8px 8px 4px 4px; /* Kształt oparcia */
  border: 2px solid transparent;
  cursor: pointer;
  transition: transform 0.1s ease, border-color 0.2s ease;
}

.seat-button:focus-visible {
  outline: 3px solid #3b82f6;
  outline-offset: 2px;
}

/* Colors logic */
.seat-available {
  background-color: #e5e7eb; /* Szary */
  border-color: #d1d5db;
}

.seat-available:hover {
  background-color: #d1d5db;
}

.seat-taken {
  background-color: #ef4444; /* Czerwony */
  border-color: #dc2626;
  cursor: not-allowed;
  opacity: 0.7;
}

.seat-selected {
  background-color: #f97316; /* Pomarańczowy */
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

/* Skeleton Seats */
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
</style>