<script setup lang="ts">
import { useMovies } from '@/composables/useMovies';
import heroImg from '@/assets/hero.png';

const { movies } = useMovies();

const getRandomStyle = () => {
  const top = Math.floor(Math.random() * 80) + 5;
  const left = Math.floor(Math.random() * 80) + 5;
  const rotation = Math.floor(Math.random() * 40) - 20;
  return {
    top: `${top}%`,
    left: `${left}%`,
    transform: `rotate(${rotation}deg)`,
    zIndex: 0
  };
};
</script>

<template>
  <div class="background-overlay">
    <img
      v-for="(movie) in movies"
      :key="movie.id"
      :src="movie.poster"
      class="bg-poster"
      :style="getRandomStyle()"
      alt=""
    />
    <div class="glass-layer"></div>
  </div>

  <header class="app-header">
    <nav class="app-nav">
      <div class="nav-left"></div> <router-link to="/" class="logo">
        <img :src="heroImg" alt="KinoVue Logo" class="logo-img" />
        <span>KinoVue</span>
      </router-link>

      <div class="nav-right">
        <router-link to="/admin" class="admin-link">Panel Admina</router-link>
      </div>
    </nav>
  </header>

  <router-view v-slot="{ Component }">
    <transition name="fade" mode="out-in">
      <component :is="Component" />
    </transition>
  </router-view>
</template>

<style scoped>
.background-overlay { position: fixed; top: 0; left: 0; width: 100vw; height: 100vh; z-index: -1; overflow: hidden; background-color: var(--bg-color); }
.bg-poster { position: absolute; width: 320px; height: auto; opacity: 0.65; filter: grayscale(10%); pointer-events: none; box-shadow: 0 10px 30px rgba(0,0,0,0.3); }
.glass-layer { position: absolute; inset: 0; background: rgba(248, 250, 252, 0.65); backdrop-filter: blur(6px); }

.app-header {
  background-color: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(5px);
  border-bottom: 1px solid var(--border);
  padding: 0.5rem 1rem;
  position: sticky;
  top: 0;
  z-index: 1000;
}

.app-nav {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
}

.logo {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  text-decoration: none;
  color: var(--accent);
  font-weight: 800;
  font-size: 1.4rem;
  justify-self: center;
}

.logo-img { height: 40px; width: auto; }

.nav-right {
  justify-self: end;
}

.admin-link {
  color: var(--text-muted);
  text-decoration: none;
  font-size: 0.85rem;
  font-weight: bold;
  padding: 0.5rem 1rem;
  border: 1px solid var(--border);
  border-radius: 6px;
  transition: all 0.2s ease;
  background: var(--surface-color);
}

.admin-link:hover {
  color: var(--accent);
  border-color: var(--accent);
}

.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

@media (max-width: 640px) {
  .admin-link { display: none; }
}
</style>