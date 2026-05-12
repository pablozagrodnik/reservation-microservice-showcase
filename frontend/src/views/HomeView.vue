<script setup lang="ts">
import { useRouter } from 'vue-router';
import { useReservationStore } from '@/stores/useReservationStore';
import { useMovies } from '@/composables/useMovies';
import type { Screening } from '@/types';

const router = useRouter();
const store = useReservationStore();
const { movies, isLoading, error, refresh } = useMovies();

const formatTime = (dateString: string): string => {
  return new Date(dateString).toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' });
};

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString('pl-PL', { weekday: 'short', day: '2-digit', month: '2-digit' });
};

const selectScreening = (screening: Screening): void => {
  store.setScreening(screening);
  router.push({ name: 'room', params: { id: screening.id } });
};
</script>

<template>
  <main class="home-page">
    <h1 class="page-title">Repertuar</h1>

    <div
        v-if="isLoading"
        class="movies-list"
        aria-busy="true"
        aria-live="polite"
        aria-label="Ładowanie repertuaru"
    >
      <article v-for="i in 3" :key="i" class="movie-skeleton" aria-hidden="true">
        <div class="skeleton-title" />
        <div class="skeleton-line" />
        <div class="skeleton-line short" />
        <div class="skeleton-screenings">
          <span class="skeleton-pill" />
          <span class="skeleton-pill" />
          <span class="skeleton-pill" />
        </div>
      </article>
    </div>

    <div
        v-else-if="error"
        class="error-state"
        role="alert"
    >
      <p class="error-message">{{ error }}</p>
      <button type="button" class="retry-btn" @click="refresh">
        Spróbuj ponownie
      </button>
    </div>

    <p v-else-if="movies.length === 0" class="empty-state">
      Brak dostępnych seansów w tej chwili.
    </p>

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
                    type="button"
                    class="screening-btn"
                    @click="selectScreening(screening)"
                    :aria-label="`Wybierz seans ${formatDate(screening.start_time)} o godzinie ${formatTime(screening.start_time)} w sali ${screening.room.name}`"
                >
                  <span class="date">{{ formatDate(screening.start_time) }}</span>
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
.home-page {
  padding: 1rem;
  max-width: 960px;
  margin: 0 auto;
}

.page-title {
  font-size: 1.5rem;
  margin: 0 0 1rem;
}

.movies-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.movie-card {
  border: 1px solid #ccc;
  border-radius: 8px;
  padding: 1rem;
}

.movie-card h2 {
  font-size: 1.25rem;
  margin: 0 0 0.5rem;
}

.description {
  margin: 0 0 1rem;
  color: #4b5563;
}

.screenings h3 {
  font-size: 1rem;
  margin: 0 0 0.5rem;
}

.screening-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  list-style: none;
  padding: 0;
  margin: 0;
}

.screening-btn {
  display: inline-flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.125rem;
  padding: 0.5rem 0.875rem;
  cursor: pointer;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  background-color: #ffffff;
  font: inherit;
  color: inherit;
  transition: background-color 0.15s ease, border-color 0.15s ease;
}

.screening-btn:hover {
  background-color: #f3f4f6;
  border-color: #9ca3af;
}

.screening-btn:focus-visible {
  outline: 2px solid #2563eb;
  outline-offset: 2px;
}

.screening-btn .date {
  font-size: 0.75rem;
  color: #6b7280;
  text-transform: capitalize;
}

.screening-btn .time {
  font-weight: 600;
}

.screening-btn .room {
  font-size: 0.75rem;
  color: #6b7280;
}

/* Skeleton */
.movie-skeleton {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 1rem;
  background-color: #ffffff;
}

.skeleton-title,
.skeleton-line,
.skeleton-pill {
  background: linear-gradient(90deg, #e5e7eb 0%, #f3f4f6 50%, #e5e7eb 100%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.4s ease-in-out infinite;
  border-radius: 4px;
}

.skeleton-title {
  height: 1.25rem;
  width: 60%;
  margin-bottom: 0.75rem;
}

.skeleton-line {
  height: 0.75rem;
  width: 100%;
  margin-bottom: 0.5rem;
}

.skeleton-line.short {
  width: 75%;
  margin-bottom: 1rem;
}

.skeleton-screenings {
  display: flex;
  gap: 0.5rem;
}

.skeleton-pill {
  height: 2.25rem;
  width: 5rem;
}

@keyframes skeleton-shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

@media (prefers-reduced-motion: reduce) {
  .skeleton-title,
  .skeleton-line,
  .skeleton-pill {
    animation: none;
  }
}

/* Error & empty */
.error-state {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
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

.retry-btn:focus-visible {
  outline: 2px solid #991b1b;
  outline-offset: 2px;
}

.empty-state {
  color: #6b7280;
  font-style: italic;
}

@media (min-width: 640px) {
  .home-page {
    padding: 2rem;
  }

  .page-title {
    font-size: 2rem;
  }
}
</style>
