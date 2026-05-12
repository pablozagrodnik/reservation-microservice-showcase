<template>
  <header class="app-header">
    <nav class="app-nav">
      <router-link to="/" class="logo">KinoVue</router-link>
      <router-link to="/admin" class="admin-link">Panel admina</router-link>
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
.app-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}
.logo {
  font-size: 1.5rem;
  font-weight: bold;
  text-decoration: none;
  color: #f97316;
}
.admin-link {
  color: #e5e7eb;
  text-decoration: none;
  font-size: 0.95rem;
  padding: 0.4rem 0.85rem;
  border: 1px solid #4b5563;
  border-radius: 6px;
  transition: background-color 0.15s ease, border-color 0.15s ease;
}
.admin-link:hover {
  background-color: #374151;
  border-color: #6b7280;
}
.admin-link.router-link-active {
  background-color: #f97316;
  border-color: #f97316;
  color: white;
}
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
