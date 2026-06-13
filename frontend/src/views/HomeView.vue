<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { useReservationStore } from '@/stores/useReservationStore';
import { useMovies } from '@/composables/useMovies';
import type { Screening } from '@/types';

const router = useRouter();
const store = useReservationStore();
const { movies, isLoading} = useMovies();

const activeIndex = ref(0);
let interval: any = null;

const availableMovies = computed(() =>
  movies.value.filter(m => m.screenings.some(s => new Date(s.start_time) > new Date()))
);

const featuredMovies = computed(() => availableMovies.value);

const goToMovie = (id: number) => {
  router.push({ name: 'movie', params: { id: String(id) } });
};

const getCarouselItemStyle = (index: number) => {
  const n = featuredMovies.value.length;
  if (n === 0) return {};

  let offset = index - activeIndex.value;

  const half = Math.floor(n / 2);
  if (offset < -half) offset += n;
  if (offset > half) offset -= n;

  const isActive = offset === 0;
  const isVisible = Math.abs(offset) <= 2;

  return {
    transform: `translateX(${offset * 220}px) scale(${isActive ? 1.1 : 0.85})`,
    opacity: isActive ? 1 : (isVisible ? 0.4 : 0),
    zIndex: isActive ? 10 : 1,
    pointerEvents: (isVisible ? 'auto' : 'none') as 'auto' | 'none'
  };
};

onMounted(() => {
  interval = setInterval(() => {
    if (featuredMovies.value.length > 0) {
      activeIndex.value = (activeIndex.value + 1) % featuredMovies.value.length;
    }
  }, 5000);
});

onUnmounted(() => clearInterval(interval));

const selectScreening = (screening: Screening, movie: any) => {
  store.setScreening(screening, movie);
  router.push({ name: 'room', params: { id: screening.id } });
};

const getPosterImage = (poster: string) => poster || '';
const formatTime = (ds: string) => new Date(ds).toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' });
const formatDate = (ds: string) => new Date(ds).toLocaleDateString('pl-PL', { weekday: 'short', day: '2-digit', month: '2-digit' });
</script>

<template>
  <main class="home-page">
    <div v-if="isLoading" class="loader">Ładowanie repertuaru...</div>

    <template v-else-if="availableMovies.length > 0">
      <section class="hero-carousel">
        <div class="carousel-viewport">
          <div class="carousel-track">
            <div
              v-for="(movie, index) in featuredMovies"
              :key="movie.id"
              class="carousel-item"
              :class="{ active: index === activeIndex }"
              :style="getCarouselItemStyle(index)"
              @click="activeIndex === index ? goToMovie(movie.id) : activeIndex = index"
            >
              <div class="poster-wrapper">
                <img :src="getPosterImage(movie.poster)" :alt="movie.title" />
              </div>

              <div class="item-overlay" v-if="index === activeIndex">
                <h2>{{ movie.title }}</h2>
                <p class="carousel-desc">{{ movie.description }}</p>
              </div>
            </div>
          </div>
        </div>

        <div class="carousel-dots">
          <span
            v-for="(_, index) in featuredMovies"
            :key="index"
            class="dot"
            :class="{ active: index === activeIndex }"
            @click="activeIndex = index"
          ></span>
        </div>
      </section>

      <section class="repertuar-section">
        <h2 class="title">Repertuar</h2>
        <div class="movies-grid">
          <article
            v-for="movie in availableMovies"
            :key="movie.id"
            class="movie-card surface"
            @click="goToMovie(movie.id)"
            style="cursor: pointer;"
          >
            <img :src="getPosterImage(movie.poster)" class="card-img" />
            <div class="card-content">
              <h3>{{ movie.title }}</h3>
              <p class="card-desc">{{ movie.description }}</p> <div class="screenings-wrapper">
                <span class="times-label">Dostępne seanse:</span>
                <ul class="screening-list">
                  <li v-for="s in movie.screenings.slice(0, 4)" :key="s.id">
                    <button
                      class="screening-btn"
                      @click.stop="selectScreening(s, movie)"
                    >
                      <span class="date">{{ formatDate(s.start_time) }}</span>
                      <span class="time">{{ formatTime(s.start_time) }}</span>
                      <span class="room">{{ s.room.name }}</span>
                    </button>
                  </li>
                </ul>
              </div>
            </div>
          </article>
        </div>
      </section>
    </template>
  </main>
</template>

<style scoped>
.home-page { padding-bottom: 4rem; overflow-x: hidden; }

.hero-carousel {
  position: relative;
  width: 100%;
  padding: 30px 0;
  background: linear-gradient(to bottom, #f8fafc, #e2e8f0);
  margin-bottom: 3rem;
  overflow: hidden;
}

.carousel-viewport {
  position: relative;
  width: 100%;
  height: 420px;
}

.carousel-track {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
}

.carousel-item {
  position: absolute;
  left: 50%;
  margin-left: -100px;
  width: 200px;
  box-sizing: border-box;
  transition: transform 0.6s ease, opacity 0.6s ease;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.carousel-item.active {
  transform: scale(1.1);
  opacity: 1;
  z-index: 10;
  cursor: default;
}

.poster-wrapper {
  width: 200px;
    height: 300px;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 10px 20px rgba(0,0,0,0.6);
    border: 2px solid transparent;
    transition: border-color 0.3s ease;
    box-sizing: border-box;
}

.carousel-item.active .poster-wrapper {
  border-color: var(--accent);
}

.poster-wrapper img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.item-overlay {
  margin-top: 15px;
  text-align: center;
  width: 220px;
  animation: fadeIn 0.8s ease forwards;
}

.item-overlay h2 {
  font-size: 1.2rem;
  margin: 0 0 6px;
  color: #0f172a;
  text-shadow: none;
}

.carousel-desc {
  font-size: 0.8rem;
  color: #475569;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.4;
}

.carousel-dots {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin-top: 10px;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--border);
  cursor: pointer;
  transition: 0.3s;
}

.dot.active {
  background: var(--accent);
  width: 24px;
  border-radius: 4px;
}

.repertuar-section {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 1.5rem;
}

.title { font-size: 2.2rem; margin-bottom: 2rem; color: var(--text-main); border-bottom: 1px solid var(--border); padding-bottom: 0.5rem; }

.movies-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 2rem;
}

@media (min-width: 900px) {
  .movies-grid { grid-template-columns: repeat(2, 1fr); }
}

.movie-card {
  display: flex;
  gap: 1.2rem;
  padding: 1.2rem;
  transition: transform 0.3s ease;
}

.movie-card:hover { transform: translateY(-3px); }

.card-img {
  width: 130px;
  height: 195px;
  object-fit: cover;
  border-radius: 8px;
  flex-shrink: 0;
}

.card-content {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.card-content h3 { margin: 0 0 8px; font-size: 1.4rem; }

.card-desc {
  font-size: 0.85rem;
  color: var(--text-muted);
  line-height: 1.4;
  margin-bottom: 12px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.times-label {
  display: block;
  font-size: 0.85rem;
  font-weight: bold;
  margin-bottom: 8px;
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
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.2rem;
  padding: 0.5rem 0.8rem;
  cursor: pointer;
  border: 2px solid var(--accent, #f97316);
  border-radius: 8px;
  background-color: var(--bg-color);
  color: var(--text-main);
  transition: all 0.2s ease;
}

.screening-btn:hover {
  background-color: var(--border);
  border-color: var(--accent);
}

.screening-btn .date { font-size: 0.75rem; color: var(--text-muted); text-transform: capitalize; }
.screening-btn .time { font-weight: 700; font-size: 1rem; color: var(--accent); }
.screening-btn .room { font-size: 0.75rem; color: var(--text-muted); }

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.loader { text-align: center; padding: 100px; color: var(--accent); font-size: 1.2rem; }
</style>