<script setup lang="ts">
import { ref } from 'vue';
import { useReservationStore } from '@/stores/useReservationStore';

const store = useReservationStore();
const email = ref('');
const isProcessing = ref(false);
const paymentSuccess = ref(false);

const handlePayment = async () => {
  if (!email.value || !store.isReadyForCheckout) return;

  isProcessing.value = true;

  await new Promise(resolve => setTimeout(resolve, 1500));

  isProcessing.value = false;
  paymentSuccess.value = true;

};

const formatDate = (dateString?: string) => {
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
      <button @click="store.clearReservation()" class="btn-secondary">Wróć do strony głównej</button>
    </div>

    <div v-else-if="!store.isReadyForCheckout" class="empty-state">
      <h2>Brak wybranch biletów</h2>
      <p>Wybierz seans i miejsca, aby przejść do płatności.</p>
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
            <span>Bilet normalny (Rząd {{ seat.row }}, Miejsce {{ seat.col }})</span>
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

        <form @submit.prevent="handlePayment">
          <div class="form-group">
            <label for="email">Adres e-mail (do wysyłki biletów)</label>
            <input
                id="email"
                type="email"
                v-model="email"
                required
                placeholder="jan.kowalski@example.com"
            />
          </div>

          <button
              type="submit"
              class="pay-btn"
              :disabled="isProcessing || !email"
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

.summary, .payment-form, .success-message {
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

.success-message {
  text-align: center;
  padding: 3rem;
  background-color: #ecfdf5;
  border-color: #a7f3d0;
  color: #065f46;
}
</style>