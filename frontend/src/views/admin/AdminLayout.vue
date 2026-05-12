<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router';
</script>

<template>
  <div class="admin-shell">
    <aside class="admin-sidebar" aria-label="Nawigacja panelu administratora">
      <h2 class="sidebar-title">Panel admina</h2>
      <nav>
        <ul class="sidebar-nav">
          <li>
            <RouterLink
                :to="{ name: 'admin-movies' }"
                class="sidebar-link"
                active-class="sidebar-link-active"
            >
              Filmy
            </RouterLink>
          </li>
          <li>
            <RouterLink
                :to="{ name: 'admin-rooms' }"
                class="sidebar-link"
                active-class="sidebar-link-active"
            >
              Sale
            </RouterLink>
          </li>
          <li>
            <RouterLink
                :to="{ name: 'admin-screenings' }"
                class="sidebar-link"
                active-class="sidebar-link-active"
            >
              Seanse
            </RouterLink>
          </li>
        </ul>
      </nav>
    </aside>

    <main class="admin-content">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.admin-shell {
  display: grid;
  grid-template-columns: 220px 1fr;
  gap: 1.5rem;
  max-width: 1100px;
  margin: 0 auto;
  padding: 1.5rem 1rem;
  text-align: left;
  align-items: start;
}

.admin-sidebar {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 1rem;
  position: sticky;
  top: 1rem;
}

.sidebar-title {
  font-size: 1rem;
  margin: 0 0 0.75rem;
  color: #374151;
  letter-spacing: 0.02em;
  text-transform: uppercase;
}

.sidebar-nav {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.sidebar-link {
  display: block;
  padding: 0.55rem 0.75rem;
  border-radius: 6px;
  text-decoration: none;
  color: #374151;
  font-weight: 500;
  transition: background-color 0.15s ease, color 0.15s ease;
}

.sidebar-link:hover {
  background-color: #f3f4f6;
}

.sidebar-link:focus-visible {
  outline: 2px solid #2563eb;
  outline-offset: 2px;
}

.sidebar-link-active {
  background-color: #f97316;
  color: white;
}

.sidebar-link-active:hover {
  background-color: #ea580c;
}

.admin-content {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 1.5rem;
  min-height: 320px;
}

@media (max-width: 720px) {
  .admin-shell {
    grid-template-columns: 1fr;
  }

  .admin-sidebar {
    position: static;
  }

  .sidebar-nav {
    flex-direction: row;
    flex-wrap: wrap;
  }
}
</style>
