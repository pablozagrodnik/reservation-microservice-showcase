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
  grid-template-columns: 250px 1fr;
  gap: 2rem;
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
  align-items: start;
}

.admin-sidebar {
  background-color: var(--surface-color);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 1.5rem;
  position: sticky;
  top: 100px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
}

.sidebar-title {
  font-size: 0.9rem;
  margin: 0 0 1rem;
  color: var(--text-muted);
  letter-spacing: 1px;
  text-transform: uppercase;
  font-weight: bold;
}

.sidebar-nav {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.sidebar-link {
  display: block;
  padding: 0.8rem 1rem;
  border-radius: 8px;
  text-decoration: none;
  color: var(--text-main);
  font-weight: 600;
  transition: all 0.2s ease;
}

.sidebar-link:hover {
  background-color: rgba(249, 115, 22, 0.1);
  color: var(--accent);
  transform: translateX(4px);
}

.sidebar-link.sidebar-link-active {
  background-color: #f1f5f9 !important;
  color: #000000 !important;
  font-weight: 800 !important;
  box-shadow: none !important;
  border-left: 4px solid var(--accent) !important;
}

.sidebar-link.sidebar-link-active:hover {
  background-color: #e2e8f0 !important;
  color: #000000 !important;
  transform: none;
}

.sidebar-link:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}

.admin-content {
  background: transparent;
  padding: 0;
  min-height: 500px;
}

@media (max-width: 768px) {
  .admin-shell {
    grid-template-columns: 1fr;
  }

  .admin-sidebar {
    position: static;
    margin-bottom: 1rem;
  }

  .sidebar-nav {
    flex-direction: row;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .sidebar-link {
    padding: 0.6rem 1rem;
  }
}
</style>