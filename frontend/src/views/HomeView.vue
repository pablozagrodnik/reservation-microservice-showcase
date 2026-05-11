<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useReservationStore } from '@/stores/useReservationStore';
import type { Movie } from '@/types';

const router = useRouter();
const store = useReservationStore();

const movies = ref<Movie[]>([]);
const isLoading = ref(true);

onMounted(async () => {
  try {
    movies.value = [
      {
        id: 1,
        title: "Diuna: Część Druga",
        description: "Paul Atryda jednoczy się z Chani i Fremenami...",
        poster: "https://fwcdn.pl/...",
        screenings: [
          {
            id: 1,
            start_time: "2026-05-12T19:00:00Z",
            room_id: 1,
            room: { id: 1, name: "Sala 1" }
          }
        ]
      }
    ];
  } finally {
    isLoading.value = false;
  }
});

const formatTime = (dateString: string) => {
  return new Date(dateString).toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' });
};

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('pl-PL', { weekday: 'short', day: '2-digit', month: '2-digit' });
};

const selectScreening = (screening: any) => {
  store.setScreening(screening);
  router.push({ name: 'room', params: { id: screening.id } });
};
</script>

<template>
  <main class="home-page">
    <h1 class="page-title">Repertuar</h1>

    <div v-if="isLoading" class="movies-list" aria-busy="true">
      <div v-for="i in 3" :key="i" class="movie-skeleton">Ładowanie...</div>
    </div>

    <div v-else class="movies-list">
      <article v-for="movie in movies" :key="movie.id" class="movie-card">
        <div class="movie-info">
          <h2>{{ movie.title }}</h2>
          <p class="description">{{ movie.description }}</p>

          <div class="screenings" aria-label="Wybierz seans">
            <h3>Dostępne seanse:</h3>
            <ul class="screening-list">
              <li v-for="screening in movie.screenings" :key="screening.id">
                <button
                    class="screening-btn"
                    @click="selectScreening(screening)"
                    :aria-label="`Wybierz seans na godzinę ${formatTime(screening.start_time)}`"
                >
                  <span class="time">{{ formatTime(screening.start_time) }}</span>
                  <span class="room">{{ screening.room.name }}</span>
                </button>
              </li>
            </ul>
          </div>
        </div>
      </article>
    </div>
  </main>
</template>

<style scoped>
.home-page { padding: 1rem; }
.movies-list { display: flex; flex-direction: column; gap: 1.5rem; }
.movie-card { border: 1px solid #ccc; border-radius: 8px; padding: 1rem; }
.screening-list { display: flex; gap: 0.5rem; list-style: none; padding: 0; }
.screening-btn { padding: 0.5rem 1rem; cursor: pointer; border-radius: 4px; }
.screening-btn:hover { background-color: #f3f4f6; }
</style>