<script setup lang="ts">
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useReservationStore } from '@/stores/useReservationStore';
import { useMovies } from '@/composables/useMovies';
import type { Screening, Movie } from '@/types';

const route = useRoute();
const router = useRouter();
const store = useReservationStore();
const { movies, isLoading, error } = useMovies();

const movieId = computed(() => Number(route.params.id));

const movie = computed<Movie | undefined>(() =>
  movies.value.find(m => m.id === movieId.value)
);

const availableScreenings = computed(() => {
  if (!movie.value) return [];
  return movie.value.screenings.filter(s => new Date(s.start_time) > new Date());
});

const selectScreening = (screening: Screening, currentMovie: Movie) => {
  store.setScreening(screening, currentMovie);
  router.push({ name: 'room', params: { id: screening.id } });
};

const getPosterImage = (poster: string) => poster || '';
const formatTime = (ds: string) => new Date(ds).toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' });
const formatDate = (ds: string) => new Date(ds).toLocaleDateString('pl-PL', { weekday: 'long', day: '2-digit', month: 'long' });
</script>

<template>
  <main class="movie-view-page">
    <div v-if="isLoading" class="loader">Ładowanie szczegółów filmu...</div>
    <div v-else-if="error" class="error-state">{{ error }}</div>
    <div v-else-if="!movie" class="error-state">Nie znaleziono filmu.</div>

    <template v-else>
      <div class="movie-container">
        <section class="poster-section">
          <img :src="getPosterImage(movie.poster)" :alt="movie.title" class="movie-poster" />
        </section>

        <section class="info-section">
          <h1 class="movie-title">{{ movie.title }}</h1>
          <div class="movie-description">
            <h2>Opis filmu</h2>
            <p>{{ movie.description }}</p>
          </div>

          <div class="screenings-container">
            <h2>Dostępne seanse</h2>
            <p v-if="availableScreenings.length === 0" class="no-screenings">
              Brak zaplanowanych seansów na najbliższy czas.
            </p>
            <ul v-else class="screening-grid">
              <li v-for="s in availableScreenings" :key="s.id">
                <button class="screening-card-btn" @click="selectScreening(s, movie)">
                  <span class="date">{{ formatDate(s.start_time) }}</span>
                  <span class="time">{{ formatTime(s.start_time) }}</span>
                  <span class="room">{{ s.room.name }}</span>
                </button>
              </li>
            </ul>
          </div>
        </section>
      </div>
    </template>
  </main>
</template>

<style scoped>
.movie-view-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1.5rem;
}

.movie-container {
  display: flex;
  flex-direction: column;
  gap: 2.5rem;
}

@media (min-width: 768px) {
  .movie-container {
    flex-direction: row;
    align-items: flex-start;
  }
}

.poster-section {
  flex-shrink: 0;
  display: flex;
  justify-content: center;
}

.movie-poster {
  width: 300px;
  height: 450px;
  object-fit: cover;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.3);
}

.info-section {
  flex: 1;
}

.movie-title {
  font-size: 2.5rem;
  margin-top: 0;
  margin-bottom: 1.5rem;
  color: var(--text-main);
}

h2 {
  font-size: 1.3rem;
  margin-bottom: 0.75rem;
  color: var(--text-main);
  border-bottom: 1px solid var(--border);
  padding-bottom: 0.25rem;
}

.movie-description p {
  font-size: 1rem;
  line-height: 1.6;
  color: var(--text-muted, #475569);
  margin-bottom: 2rem;
}

.screening-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 1rem;
  list-style: none;
  padding: 0;
  margin: 0;
}

.screening-card-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.3rem;
  width: 100%;
  padding: 1rem;
  cursor: pointer;
  border: 2px solid var(--accent, #f97316);
  border-radius: 8px;
  background-color: var(--bg-color, #ffffff);
  color: var(--text-main);
  transition: all 0.2s ease;
}

.screening-card-btn:hover {
  background-color: var(--accent, #f97316);
  color: white;
}

.screening-card-btn:hover .date,
.screening-card-btn:hover .time,
.screening-card-btn:hover .room {
  color: white;
}

.screening-card-btn .date {
  font-size: 0.8rem;
  color: var(--text-muted);
  text-transform: capitalize;
}

.screening-card-btn .time {
  font-weight: bold;
  font-size: 1.3rem;
  color: var(--accent, #f97316);
}

.screening-card-btn .room {
  font-size: 0.8rem;
  color: var(--text-muted);
}

.loader, .error-state, .no-screenings {
  text-align: center;
  padding: 3rem;
  color: var(--text-muted);
  font-style: italic;
}
</style>