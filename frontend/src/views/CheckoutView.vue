<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useReservationStore } from '@/stores/useReservationStore';
import { createReservation } from '@/api/reservations';
import { ApiError } from '@/api/client';
import type { Seat } from '@/types';

const store = useReservationStore();
const router = useRouter();

const email = ref<string>('');
const isProcessing = ref<boolean>(false);
const paymentSuccess = ref<boolean>(false);
const errorMessage = ref<string | null>(null);
const conflictedSeats = ref<Seat[]>([]);

const canSubmit = computed<boolean>(
    () => !!email.value && store.isReadyForCheckout && !isProcessing.value
);

const formatSeatLabel = (seat: Seat): string => `Rząd ${seat.row}, miejsce ${seat.col}`;

const handlePayment = async (): Promise<void> => {
  const screening = store.screening;
  if (!email.value || !screening || store.selectedSeats.length === 0) return;

  isProcessing.value = true;
  errorMessage.value = null;
  conflictedSeats.value = [];

  const seats: Seat[] = [...store.selectedSeats];
  const trimmedEmail = email.value.trim();

  // Promise.allSettled zamiast Promise.all — chcemy wiedzieć,
  // które konkretnie miejsca się nie udały (np. 409), nawet jeśli inne się zapisały.
  const results = await Promise.allSettled(
      seats.map(seat =>
          createReservation({
            screening_id: screening.id,
            seat_id: seat.id,
            user_email: trimmedEmail,
          })
      )
  );

  const failures: { seat: Seat; reason: unknown }[] = [];
  results.forEach((res, idx) => {
    if (res.status === 'rejected') {
      failures.push({ seat: seats[idx], reason: res.reason });
    }
  });

  isProcessing.value = false;

  if (failures.length === 0) {
    paymentSuccess.value = true;
    store.clearReservation();
    return;
  }

  const conflicts = failures.filter(
      f => f.reason instanceof ApiError && f.reason.status === 409
  );

  if (conflicts.length === failures.length) {
    conflictedSeats.value = conflicts.map(c => c.seat);
    errorMessage.value =
        conflicts.length === seats.length
            ? 'Niestety, wybrane miejsca zostały właśnie zajęte przez kogoś innego. Wróć do wyboru miejsc i spróbuj ponownie.'
            : 'Niektóre z wybranych miejsc zostały właśnie zajęte przez kogoś innego. Wróć do wyboru miejsc i spróbuj ponownie.';
  } else {
    errorMessage.value =
        'Wystąpił nieoczekiwany błąd podczas rezerwacji. Spróbuj ponownie za chwilę.';
  }
};

const goHome = (): void => {
  void router.push({ name: 'home' });
};

const goBackToSeats = (): void => {
  const screeningId = store.screening?.id;
  if (screeningId !== undefined) {
    void router.push({ name: 'room', params: { id: String(screeningId) } });
  } else {
    void router.push({ name: 'home' });
  }
};

const formatDate = (dateString?: string): string => {
  if (!dateString) return '';
  return new Date(dateString).toLocaleString('pl-PL', {
    weekday: 'long', day: '2-digit', month: 'long', hour: '2-digit', minute: '2-digit'
  });
};
</script>

<template>
  <main class="checkout-page">
    <div v-if="paymentSuccess" class="success-message" role="alert">
      <h2>🎉 Płatność zakończona sukcesem!</h2>
      <p>Bilety zostały wysłane na adres: <strong>{{ email }}</strong></p>
      <button @click="goHome" class="btn-secondary">Wróć do strony głównej</button>
    </div>

    <div v-else-if="!store.isReadyForCheckout" class="empty-state">
      <h2>Brak wybranych biletów</h2>
      <p>Wybierz seans i miejsca, aby przejść do płatności.</p>
      <button @click="goHome" class="btn-secondary">Przejdź do repertuaru</button>
    </div>

    <div v-else class="checkout-container">
      <section class="summary" aria-labelledby="summary-heading">
        <h2 id="summary-heading">Twoja rezerwacja</h2>

        <div class="movie-details">
          <h3>Sala: {{ store.screening?.room.name }}</h3>
          <p class="date">{{ formatDate(store.screening?.start_time) }}</p>
        </div>

        <ul class="tickets-list">
          <li v-for="seat in store.selectedSeats" :key="seat.id" class="ticket-item">
            <span>Bilet normalny ({{ formatSeatLabel(seat) }})</span>
            <span>25,00 zł</span>
          </li>
        </ul>

        <div class="total">
          <span>Suma do zapłaty:</span>
          <strong>{{ store.totalAmount.toFixed(2) }} zł</strong>
        </div>
      </section>

      <section class="payment-form" aria-labelledby="payment-heading">
        <h2 id="payment-heading">Dane płatności</h2>

        <form @submit.prevent="handlePayment" novalidate>
          <div class="form-group">
            <label for="email">Adres e-mail (do wysyłki biletów)</label>
            <input
                id="email"
                type="email"
                v-model="email"
                required
                placeholder="jan.kowalski@example.com"
                :disabled="isProcessing"
                autocomplete="email"
            />
          </div>

          <div
              v-if="errorMessage"
              class="error-message"
              role="alert"
              aria-live="assertive"
          >
            <p>{{ errorMessage }}</p>
            <ul v-if="conflictedSeats.length > 0" class="conflict-list">
              <li v-for="seat in conflictedSeats" :key="seat.id">
                {{ formatSeatLabel(seat) }}
              </li>
            </ul>
            <button type="button" class="btn-secondary" @click="goBackToSeats">
              Wróć do wyboru miejsc
            </button>
          </div>

          <button
              type="submit"
              class="pay-btn"
              :disabled="!canSubmit"
              :aria-busy="isProcessing"
          >
            <span v-if="isProcessing">Przetwarzanie...</span>
            <span v-else>Zapłać {{ store.totalAmount.toFixed(2) }} zł</span>
          </button>
        </form>
      </section>
    </div>
  </main>
</template>

<style scoped>
.checkout-page {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.checkout-container {
  display: grid;
  grid-template-columns: 1fr;
  gap: 2rem;
}

@media (min-width: 768px) {
  .checkout-container {
    grid-template-columns: 1fr 1fr;
  }
}

.summary, .payment-form, .success-message, .empty-state {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.05);
  border: 1px solid #e5e7eb;
}

.tickets-list {
  list-style: none;
  padding: 0;
  margin: 1.5rem 0;
  border-top: 1px solid #e5e7eb;
  border-bottom: 1px solid #e5e7eb;
}

.ticket-item {
  display: flex;
  justify-content: space-between;
  padding: 0.75rem 0;
}

.total {
  display: flex;
  justify-content: space-between;
  font-size: 1.25rem;
  margin-top: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  margin-bottom: 1.5rem;
}

.form-group label {
  font-weight: 500;
  margin-bottom: 0.5rem;
}

.form-group input {
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 1rem;
}

.form-group input:focus {
  outline: 2px solid #f97316;
  border-color: #f97316;
}

.form-group input:disabled {
  background-color: #f3f4f6;
  cursor: not-allowed;
}

.pay-btn {
  width: 100%;
  padding: 1rem;
  background-color: #f97316;
  color: white;
  font-weight: bold;
  font-size: 1.1rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.pay-btn:disabled {
  background-color: #fdba74;
  cursor: not-allowed;
}

.btn-secondary {
  margin-top: 1rem;
  padding: 0.75rem 1.5rem;
  background-color: #f3f4f6;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  cursor: pointer;
}

.btn-secondary:hover {
  background-color: #e5e7eb;
}

.success-message {
  text-align: center;
  padding: 3rem;
  background-color: #ecfdf5;
  border-color: #a7f3d0;
  color: #065f46;
}

.empty-state {
  text-align: center;
  padding: 3rem;
  color: #4b5563;
}

.error-message {
  margin-bottom: 1.5rem;
  padding: 1rem;
  background-color: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 6px;
  color: #991b1b;
}

.error-message p {
  margin: 0 0 0.5rem;
  font-weight: 500;
}

.conflict-list {
  margin: 0.5rem 0 0.75rem 1.25rem;
  padding: 0;
}

.conflict-list li {
  margin-bottom: 0.25rem;
}

.error-message .btn-secondary {
  margin-top: 0.5rem;
  background-color: white;
  border-color: #fecaca;
  color: #991b1b;
}

.error-message .btn-secondary:hover {
  background-color: #fee2e2;
}
</style>
