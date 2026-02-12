<template>
  <div class="app-wrapper">
    <!-- Unified Top Bar (only when NOT on login page) -->
    <header v-if="$route.path !== '/login'" class="app-topbar">
      <div class="header-left">
        <h1 class="neon-title">WASA</h1>
        <span class="subtitle">Chats</span>
      </div>
      
      <div class="header-actions">

  <!-- LEFT -->
  <div class="actions-left">

    <!-- Home -->
    <RouterLink to="/home" class="btn-icon" title="Home">
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"
        fill="none" stroke="currentColor" stroke-width="2">
        <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
        <polyline points="9 22 9 12 15 12 15 22"/>
      </svg>
    </RouterLink>

    <!-- New Chat -->
    <button
      class="btn-icon"
      title="New Chat"
      @click="emitNewChat"
    >
      <!-- chat bubble -->
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"
        fill="none" stroke="currentColor" stroke-width="2">
        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
      </svg>
    </button>

    <!-- New Group -->
    <RouterLink to="/groups/create" class="btn-icon" title="New Group">
      <!-- users + plus -->
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"
        fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="9" cy="7" r="4"/>
        <path d="M3 21v-2a4 4 0 0 1 4-4h4"/>
        <path d="M16 11v6"/>
        <path d="M13 14h6"/>
      </svg>
    </RouterLink>

  </div>

  <!-- RIGHT -->
  <div class="actions-right">
     <!-- Profile (ICONA ORIGINALE TENUTA) -->
    <RouterLink to="/profile" class="btn-icon" title="Profile">
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"
        fill="none" stroke="currentColor" stroke-width="2">
        <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
        <circle cx="12" cy="7" r="4"/>
      </svg>
    </RouterLink>


    <!-- Users (ICONA ORIGINALE TENUTA) -->
    <RouterLink to="/users" class="btn-icon" title="Users">
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"
        fill="none" stroke="currentColor" stroke-width="2">
        <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
        <circle cx="9" cy="7" r="4"/>
        <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
        <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
      </svg>
    </RouterLink>
    <!-- Logout -->
    <button @click="logout" class="btn-logout" title="Logout">
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"
        fill="none" stroke="currentColor" stroke-width="2">
        <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
        <polyline points="16 17 21 12 16 7"/>
        <line x1="21" y1="12" x2="9" y2="12"/>
      </svg>
    </button>

  </div>
</div>


    </header>

    <!-- Main Content -->
    <main :class="{ 'with-topbar': $route.path !== '/login', 'fullscreen': $route.path === '/login' }">
      <RouterView @new-chat="handleNewChat" @refresh="handleRefresh" />
    </main>
  </div>
</template>

<script setup>
import { RouterLink, RouterView, useRouter } from 'vue-router'
import { onMounted } from 'vue'
import axios from './services/axios.js'

const router = useRouter()

onMounted(async () => {
  try {
    await axios.get('/liveness')
  } catch (e) {
    console.error('Liveness check failed', e)
  }
})

// Event handlers for HomeView communication
const emitNewChat = () => {
  window.dispatchEvent(new CustomEvent('app-new-chat'))
}

const emitRefresh = () => {
  window.dispatchEvent(new CustomEvent('app-refresh'))
}

const handleNewChat = () => {
  // Handle if needed
}

const handleRefresh = () => {
  // Handle if needed
}

const logout = () => {
  try {
    const preserve = localStorage.getItem('leftConversations')
    localStorage.clear()
    if (preserve) localStorage.setItem('leftConversations', preserve)
  } catch (e) {
    console.error('Error clearing localStorage:', e)
  }
  
  try { 
    delete axios.defaults.headers.common['Authorization'] 
  } catch (e) {
    console.error('Error clearing axios headers:', e)
  }
  
  router.push('/login')
}
</script>

<style>
/* App Wrapper */
/* App Wrapper */
.app-wrapper {
  width: 100%;
  min-height: 100vh;
  background: #000000;
  display: flex;
  flex-direction: column;
}

/* Unified Top Bar */
.app-topbar {
  position: sticky;
  top: 0;
  z-index: 200; /* AUMENTATO: deve stare sopra tutto */
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 2rem;
  height: 76px; /* AGGIUNTO: altezza fissa */
  background: rgba(10, 10, 30, 0.98); /* AUMENTATO: più opaco */
  border-bottom: 2px solid rgba(0, 229, 255, 0.3);
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
}
.actions-left,
.actions-right {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  margin-left: 2rem;
}

.header-left {
  display: flex;
  align-items: baseline;
  gap: 1rem;
}

.neon-title {
  font-size: 2rem;
  font-weight: 900;
  margin: 0;
  background: linear-gradient(90deg, #00e5ff 0%, #8a2be2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: 0.3rem;
  filter: drop-shadow(0 0 10px rgba(0, 229, 255, 0.6));
}

.subtitle {
  color: rgba(255, 255, 255, 0.6);
  font-size: 1rem;
  font-weight: 400;
}

.header-actions {
  display: flex;
  gap: 0.75rem;
  align-items: center;
}

.btn-icon {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 229, 255, 0.1);
  border: 2px solid rgba(0, 229, 255, 0.3);
  border-radius: 12px;
  color: #00e5ff;
  cursor: pointer;
  transition: all 0.3s ease;
  text-decoration: none;
}

.btn-icon:hover {
  background: rgba(0, 229, 255, 0.2);
  border-color: #00e5ff;
  box-shadow: 0 0 15px rgba(0, 229, 255, 0.4);
  transform: translateY(-2px);
}

.btn-icon.router-link-active {
  background: rgba(0, 229, 255, 0.25);
  border-color: #00e5ff;
  box-shadow: 0 0 20px rgba(0, 229, 255, 0.5);
}

.btn-logout {
  padding: 0.75rem 1.25rem;
  background: transparent;
  border: 2px solid rgba(255, 0, 100, 0.4);
  border-radius: 12px;
  color: #ff0064;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s ease;
  font-weight: 600;
}

.btn-logout:hover {
  background: rgba(255, 0, 100, 0.15);
  border-color: #ff0064;
  box-shadow: 0 0 15px rgba(255, 0, 100, 0.4);
}

/* Main Content - MODIFIED */
main {
  flex: 1;
  width: 100%;
  min-height: 0; /* AGGIUNTO: permette di shrinkare correttamente */
}

main.fullscreen {
  height: 100vh;
}

main.with-topbar {
  height: calc(100vh - 76px); /* CHANGED: da min-height a height */
  /* overflow rimosso - lo scroll è gestito da .messages-neon */
}

/* Responsive */
@media (max-width: 768px) {
  .app-topbar {
    padding: 1rem 1.5rem;
    height: 68px; /* AGGIUNTO */
  }

  main.with-topbar {
    height: calc(100vh - 68px); /* AGGIORNATO */
  }

  .header-left {
    gap: 0.5rem;
  }

  .neon-title {
    font-size: 1.5rem;
    letter-spacing: 0.2rem;
  }

  .subtitle {
    font-size: 0.875rem;
  }

  .header-actions {
    gap: 0.5rem;
  }

  .btn-icon {
    width: 40px;
    height: 40px;
  }

  .btn-logout {
    padding: 0.6rem 1rem;
  }
}

@media (max-width: 480px) {
  .app-topbar {
    height: 64px; /* AGGIUNTO */
  }

  main.with-topbar {
    height: calc(100vh - 64px); /* AGGIORNATO */
  }

  .subtitle {
    display: none;
  }
  
  .header-actions {
    gap: 0.35rem;
  }
  
  .btn-icon {
    width: 36px;
    height: 36px;
  }
  
  .btn-logout {
    padding: 0.6rem 0.8rem;
  }
}
</style>