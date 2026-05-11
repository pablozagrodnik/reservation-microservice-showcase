<template>
  <header class="app-header">
    <nav>
      <router-link to="/" class="logo">🎬 KinoVue</router-link>
    </nav>
  </header>

  <router-view v-slot="{ Component }">
    <transition name="fade" mode="out-in">
      <component :is="Component" />
    </transition>
  </router-view>
</template>

<style>
.app-header {
  padding: 1rem 2rem;
  background: #1f2937;
  color: white;
}
.logo {
  font-size: 1.5rem;
  font-weight: bold;
  text-decoration: none;
  color: #f97316;
}
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>