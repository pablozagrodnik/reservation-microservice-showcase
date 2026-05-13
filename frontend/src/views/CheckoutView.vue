<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useReservationStore } from '@/stores/useReservationStore';
import { createReservation } from '@/api/reservations';
import { ApiError } from '@/api/client';
import type { Seat } from '@/types';
import TicketSuccess from '@/components/TicketSuccess.vue';

const store = useReservationStore();
const router = useRouter();

const email = ref<string>('');
const isProcessing = ref<boolean>(false);
const paymentSuccess = ref<boolean>(false);
const successDetails = ref<any | null>(null);
const errorMessage = ref<string | null>(null);
const conflictedSeats = ref<Seat[]>([]);
const emailError = ref<string | null>(null);

const canSubmit = computed<boolean>(
    () => !!email.value && isValidEmail(email.value) && store.isReadyForCheckout && !isProcessing.value
);

const isValidEmail = (emailStr: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(emailStr);
};

const onEmailInput = (): void => {
  if (email.value && !isValidEmail(email.value)) {
    emailError.value = 'Podaj prawidłowy adres e-mail.';
  } else {
    emailError.value = null;
  }
};

const formatSeatLabel = (seat: Seat): string => `Rząd ${seat.row}, miejsce ${seat.col}`;

const handlePayment = async (): Promise<void> => {
   const screening = store.screening;
   if (!email.value || !isValidEmail(email.value) || !screening || store.selectedSeats.length === 0) return;

  isProcessing.value = true;
  errorMessage.value = null;
  conflictedSeats.value = [];

  const seats: Seat[] = [...store.selectedSeats];
  const trimmedEmail = email.value.trim();

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
    const reservationIds: number[] = results
        .filter(r => r.status === 'fulfilled')
        .map((r: any) => r.value.id);

    const successfulSeats: Seat[] = results
        .map((r, idx) => ({ res: r, seat: seats[idx] }))
        .filter(x => x.res.status === 'fulfilled')
        .map(x => x.seat);

    successDetails.value = {
      reservationIds,
      movieTitle: store.movie?.title ?? '',
      roomName: screening.room.name,
      startTime: screening.start_time,
      seats: successfulSeats,
      email: trimmedEmail,
      total: store.totalAmount,
      poster: store.movie?.poster ?? ''
    };

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
    <div v-if="paymentSuccess && successDetails" class="success-wrapper" role="alert">
      <TicketSuccess :details="successDetails" />
      <div style="text-align:center; margin-top:2rem;">
        <button @click="goHome" class="btn-secondary">Wróć do strony głównej</button>
      </div>
    </div>

    <div v-else-if="!store.isReadyForCheckout" class="empty-state">
      <h2>Brak wybranych biletów</h2>
      <p>Wybierz seans i miejsca, aby przejść do płatności.</p>
      <button @click="goHome" class="btn-secondary">Przejdź do repertuaru</button>
    </div>

    <div v-else class="checkout-container">
       <section class="summary" aria-labelledby="summary-heading">
         <div class="banner-bg" :style="{ backgroundImage: `url(${store.movie?.poster})` }"></div>
         <div class="summary-content">
           <h2 id="summary-heading">Twoja rezerwacja</h2>

           <div class="movie-details">
             <h3 v-if="store.movie">Film: {{ store.movie.title }}</h3>
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
                 @input="onEmailInput"
             />
             <span v-if="emailError" class="field-error">{{ emailError }}</span>
           </div>

          <div
              v-if="errorMessage"
              class="message-box error"
              role="alert"
              aria-live="assertive"
          >
            <p>{{ errorMessage }}</p>
            <ul v-if="conflictedSeats.length > 0" class="conflict-list">
              <li v-for="seat in conflictedSeats" :key="seat.id">
                {{ formatSeatLabel(seat) }}
              </li>
            </ul>
            <button type="button" class="btn-secondary danger" @click="goBackToSeats">
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

.payment-form, .empty-state {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.05);
  border: 1px solid #e5e7eb;
}

.summary {
  position: relative;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.05);
  overflow: hidden;
  color: white;
}

.banner-bg {
  position: absolute;
  top: 0; left: 0; width: 100%; height: 100%;
  background-size: cover;
  background-position: center;
  filter: brightness(0.2) blur(2px);
  z-index: 0;
}

.summary-content {
  position: relative;
  z-index: 1;
}

.summary-content h2, .summary-content h3 {
  color: white;
  margin-top: 0;
}

.tickets-list {
  list-style: none;
  padding: 0;
  margin: 1.5rem 0;
  border-top: 1px solid rgba(255, 255, 255, 0.2);
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
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

.field-error {
  color: #991b1b;
  font-size: 0.875rem;
  margin-top: 0.25rem;
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
  color: #374151;
  font-weight: 500;
}

.btn-secondary:hover {
  background-color: #e5e7eb;
}

.btn-secondary.danger {
  background-color: white;
  border-color: #fecaca;
  color: #991b1b;
}

.btn-secondary.danger:hover {
  background-color: #fee2e2;
}

.success-wrapper {
  margin-top: 1rem;
}

.empty-state {
  text-align: center;
  padding: 3rem;
  color: #4b5563;
}

.message-box {
  margin-bottom: 1.5rem;
  padding: 1rem;
  border-radius: 6px;
}

.message-box.error {
  background-color: #fef2f2;
  border: 1px solid #fecaca;
  color: #991b1b;
}

.message-box p {
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
</style>