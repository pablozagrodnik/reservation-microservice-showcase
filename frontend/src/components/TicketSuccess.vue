<template>
  <section class="ticket-wrapper" aria-label="Twoje bilety">
    <p class="sent-info">Płatność zakończona sukcesem!</p>

    <div class="ticket">
      <header class="ticket-header">
        <div class="movie-title">{{ details.movieTitle }}</div>
        <div class="meta">
          <div class="room">Sala: {{ details.roomName }}</div>
          <div class="date">{{ formattedDate }}</div>
        </div>
      </header>

      <div class="ticket-body">
        <div v-if="props.details.poster" class="poster-container">
          <img :src="props.details.poster" :alt="`Plakat: ${props.details.movieTitle}`" class="poster-image" />
        </div>

        <ul class="seats-list">
          <li v-for="seat in details.seats" :key="seat.id">Rząd {{ seat.row }}, miejsce {{ seat.col }}</li>
        </ul>

        <div class="qr">
          <QRCodeVue :value="qrValue" :size="120" render-as="svg" :margin="0" level="H" />
          <div class="qr-caption">Pokaż kod przy wejściu</div>
        </div>
      </div>

      <div class="ticket-footer">
        <div class="price">Suma: <strong>{{ details.total.toFixed(2) }} zł</strong></div>
        <button class="btn-primary" @click="printTicket">Drukuj / Pobierz</button>
      </div>
    </div>

    <p class="sent-info">Bilet został wysłany na adres: <strong>{{ details.email }}</strong></p>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import QRCodeVue from 'qrcode.vue';

interface Seat {
  id: number;
  row: number;
  col: number;
}

interface Details {
  reservationIds: number[];
  movieTitle: string;
  roomName: string;
  startTime?: string;
  seats: Seat[];
  email?: string;
  total: number;
  poster?: string;
}

const props = defineProps<{ details: Details }>();

const formattedDate = computed(() => {
  if (!props.details.startTime) return '';
  const d = new Date(props.details.startTime);
  return d.toLocaleString('pl-PL', { weekday: 'long', day: '2-digit', month: 'long', hour: '2-digit', minute: '2-digit' });
});

const qrValue = computed(() => {
  const ids = props.details.reservationIds?.join(',') ?? '';
  return `Rezerwacja:${ids};film:${props.details.movieTitle};sala:${props.details.roomName}`;
});

const printTicket = (): void => {
  window.print();
};
</script>

<style scoped>
.ticket-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  justify-content: center;
  padding: 1rem;
}

.ticket {
  width: 100%;
  max-width: 700px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(2,6,23,0.12);
  border: 1px solid #e5e7eb;
  overflow: hidden;
  padding: 1rem;
  position: relative;
  background-image: linear-gradient(180deg, rgba(249,115,22,0.04), transparent 30%);
}

.ticket::after {
  content: '';
  position: absolute;
  left: 50%;
  top: 0;
  bottom: 0;
  width: 1px;
  transform: translateX(-50%);
  border-left: 1px dashed rgba(0,0,0,0.08);
}

.ticket-header {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  align-items: flex-start;
}

.movie-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: #111827;
}

.meta {
  text-align: right;
  color: #6b7280;
}

.ticket-body {
  display: flex;
  gap: 1rem;
  padding: 1rem 0;
  align-items: center;
  border-top: 1px dashed #e5e7eb;
  border-bottom: 1px dashed #e5e7eb;
}

.poster-container {
  flex: 0 0 120px;
  width: 120px;
  height: 180px;
  background: #f3f4f6;
  border-radius: 6px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.poster-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.seats-list {
  list-style: none;
  padding: 0;
  margin: 0;
  flex: 1 1 50%;
}

.seats-list li {
  padding: 0.25rem 0;
  color: #374151;
  font-weight: 600;
}

.qr {
  flex: 0 0 140px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.qr-caption {
  margin-top: 0.5rem;
  font-size: 0.85rem;
  color: #6b7280;
}

.ticket-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 0.75rem;
}

.btn-primary {
  background: #f97316;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
}

@media (max-width: 640px) {
  .ticket {
    padding: 0.75rem;
  }
  .ticket-body {
    flex-direction: column;
    gap: 0.75rem;
  }
  .qr { flex: none; }
}

@media print {
  @page {
    margin: 8mm;
  }

  :global(html),
  :global(body) {
    margin: 0 !important;
    padding: 0 !important;
    background: #fff !important;
  }

  :global(#app *) {
    visibility: hidden !important;
  }

  :global(#app .ticket-wrapper),
  :global(#app .ticket-wrapper *) {
    visibility: visible !important;
  }

  :global(#app .ticket-wrapper) {
    position: absolute !important;
    left: 0 !important;
    top: 0 !important;
    width: 100% !important;
    padding: 0 !important;
    gap: 0 !important;
  }

  :global(#app .ticket) {
    width: 150mm !important;
    max-width: 150mm !important;
    margin: 0 auto !important;
    box-shadow: none !important;
    border: 1px solid #e5e7eb !important;
    border-radius: 12px !important;
    background: #fff !important;
  }

  :global(#app .btn-primary),
  :global(#app .sent-info) {
    display: none !important;
  }
}

.sent-info {
  text-align: center;
  margin: 0.75rem auto 0;
  color: #6b7280;
  max-width: 700px;
  width: 100%;
}
</style>

