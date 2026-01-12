<template>
  <div class="login-container">
    <Particles />
    <img src="@/assets/wasaLogo.png" alt="WASA Logo" class="header-logo" />

    <div class="login-card">
      <!-- Logo WASA -->
      <div class="logo-section">
        <h1 class="neon-logo">WASA</h1>
        <p class="tagline">Your Messaging Universe</p>
      </div>

      <!-- Login Form -->
      <form @submit.prevent="handleLogin" class="login-form">
        <!-- Username Input -->
        <div class="input-group">
          <div class="input-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
              <circle cx="12" cy="7" r="4"></circle>
            </svg>
          </div>
          <input
            v-model="username"
            type="text"
            class="neon-input"
            placeholder="Username"
            required
            minlength="3"
            maxlength="20"
            :disabled="loading"
            @input="clearError"
          />
        </div>

        <!-- Error Message -->
        <div v-if="error" class="error-message">
          {{ error }}
        </div>

        <!-- Login Button -->
        <button type="submit" class="neon-button" :disabled="loading">
          <span v-if="!loading">LOGIN</span>
          <span v-else class="loading-dots">
            <span>L</span><span>O</span><span>A</span><span>D</span>
            <span>I</span><span>N</span><span>G</span>
          </span>
        </button>

        <!-- Info Text -->
        <p class="info-text">
          Enter any username to login or create a new account
        </p>
      </form>
    </div>

    <!-- Footer -->
    <div class="login-footer">
      <p>Powered by A.C. • 2026</p>
    </div>
  </div>
</template>

<script>
import Particles from './particles.vue'

export default {
  name: 'LoginView',
  components: { Particles },
  data() {
    return {
      username: '',
      error: null,
      loading: false
    }
  },
  methods: {
    async handleLogin() {
      const name = this.username.trim()
      if (!name) { this.error = "Username cannot be empty"; return }
      if (name.length < 3) { this.error = "Username must be at least 3 characters"; return }
      if (name.length > 20) { this.error = "Username must be 20 characters or less"; return }

      this.loading = true
      this.error = null

      try {
        const response = await this.$axios.post('/login', { name })
        const data = response.data
        const token = data.identifier || data.token
        const userId = data.id
        const username = data.username
        const photo = data.photo

        if (!token || !userId) throw new Error('Invalid server response')

        localStorage.setItem('token', token)
        localStorage.setItem('userId', userId)
        localStorage.setItem('username', username || name)

        if (photo) {
          localStorage.setItem('userPhoto', photo.startsWith('data:image') ? photo : `data:image/png;base64,${photo}`)
        } else {
          localStorage.setItem('userPhoto', '')
        }

        if (this.$axios.defaults?.headers) {
          this.$axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
        }

        this.$router.push('/home')

      } catch (err) {
        console.error('Login error:', err)
        if (err.response) {
          const status = err.response.status
          const errorData = err.response.data
          this.error = status === 400 ? (errorData?.error || 'Invalid username format') :
                       status === 409 ? 'Username already taken' :
                       status === 401 ? 'Authentication failed' :
                       status === 500 ? 'Server error. Please try again later.' :
                       errorData?.error || `Server error (${status})`
        } else if (err.request) {
          this.error = 'Cannot connect to server. Is it running?'
        } else {
          this.error = err.message || 'Login failed'
        }
      } finally {
        this.loading = false
      }
    },
    clearError() { this.error = null }
  },
  mounted() {
    try { localStorage.clear() } catch (e) { console.error(e) }
    try { delete this.$axios.defaults.headers.common['Authorization'] } catch (e) {}
  }
}
</script>

<style scoped>
  .header-logo {
  position: absolute;
  top: 20px;
  left: 20px;
  width: auto; 
  height: 120px;
  z-index: 20; 
}

/* Container */
.login-container {
  position: relative;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background: #000000;
  overflow: hidden;
  padding: 2rem;
}



@keyframes float-up {
  0% {
    transform: translateY(0) translateX(0);
    opacity: 1;
  }
  100% {
    transform: translateY(-100vh) translateX(50px);
    opacity: 0;
  }
}

/* Login Card */
.login-card {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 480px;
  padding: 3rem 2.5rem;
  background: rgba(10, 10, 30, 0.9);
  border: 2px solid transparent;
  border-radius: 20px;
  backdrop-filter: blur(10px);
  box-shadow: 
    0 0 40px rgba(0, 229, 255, 0.3),
    0 0 80px rgba(138, 43, 226, 0.2),
    inset 0 0 60px rgba(0, 229, 255, 0.05);
  animation: cardGlow 3s ease-in-out infinite alternate;
}

@keyframes cardGlow {
  from {
    border-color: rgba(0, 229, 255, 0.3);
    box-shadow: 
      0 0 40px rgba(0, 229, 255, 0.3),
      0 0 80px rgba(138, 43, 226, 0.2),
      inset 0 0 60px rgba(0, 229, 255, 0.05);
  }
  to {
    border-color: rgba(138, 43, 226, 0.5);
    box-shadow: 
      0 0 50px rgba(138, 43, 226, 0.4),
      0 0 90px rgba(0, 229, 255, 0.3),
      inset 0 0 80px rgba(138, 43, 226, 0.08);
  }
}

/* Logo Section */
.logo-section {
  text-align: center;
  margin-bottom: 3rem;
}

.neon-logo {
  font-size: 5rem;
  font-weight: 900;
  margin: 0;
  letter-spacing: 0.5rem;
  background: linear-gradient(90deg, #00e5ff 0%, #8a2be2 50%, #00e5ff 100%);
  background-size: 200% auto;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: gradientShift 7s ease-in-out infinite;
  filter: drop-shadow(0 0 20px rgba(0, 229, 255, 0.8));
}

@keyframes gradientShift {
  0%, 100% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
}

.tagline {
  color: #00e5ff;
  font-size: 0.95rem;
  margin-top: 0.5rem;
  letter-spacing: 0.2rem;
  text-transform: uppercase;
  opacity: 0.8;
}

/* Form */
.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

/* Input Group */
.input-group {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 1.25rem;
  color: #00e5ff;
  z-index: 2;
  pointer-events: none;
}

.neon-input {
  width: 100%;
  padding: 1rem 1rem 1rem 3.5rem;
  background: rgba(0, 0, 0, 0.5);
  border: 2px solid rgba(0, 229, 255, 0.3);
  border-radius: 12px !important;
  color: #ffffff;
  font-size: 1rem;
  transition: all 0.3s ease;
  outline: none;
}

.neon-input::placeholder {
  color: rgba(255, 255, 255, 0.4);
}

.neon-input:focus {
  border-color: #00e5ff;
  box-shadow: 
    0 0 20px rgba(0, 229, 255, 0.4),
    inset 0 0 20px rgba(0, 229, 255, 0.1);
  background: rgba(0, 20, 40, 0.6);
}

.neon-input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Error Message */
.error-message {
  padding: 0.875rem 1.25rem;
  background: rgba(255, 0, 100, 0.15);
  border: 1px solid rgba(255, 0, 100, 0.4);
  border-radius: 10px;
  color: #ff0064;
  font-size: 0.875rem;
  text-align: center;
  animation: errorShake 0.5s ease;
}

@keyframes errorShake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-10px); }
  75% { transform: translateX(10px); }
}

/* Neon Button */
.neon-button {
  width: 100%;
  padding: 1.25rem;
  background: linear-gradient(90deg, #00e5ff 0%, #8a2be2 100%);
  border: none;
  border-radius: 12px;
  color: #000000;
  font-size: 1.125rem;
  font-weight: 700;
  letter-spacing: 0.3rem;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;
  box-shadow: 0 0 30px rgba(0, 229, 255, 0.5);
}

.neon-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
  transition: left 0.5s ease;
}

.neon-button:hover::before {
  left: 100%;
}

.neon-button:hover {
  transform: translateY(-2px);
  box-shadow: 
    0 0 40px rgba(0, 229, 255, 0.7),
    0 5px 20px rgba(0, 0, 0, 0.3);
}

.neon-button:active {
  transform: translateY(0);
}

.neon-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

/* Loading Animation */
.loading-dots {
  display: inline-flex;
  gap: 0.2rem;
}

.loading-dots span {
  animation: loadingBounce 1.4s infinite ease-in-out both;
}

.loading-dots span:nth-child(1) { animation-delay: -0.32s; }
.loading-dots span:nth-child(2) { animation-delay: -0.28s; }
.loading-dots span:nth-child(3) { animation-delay: -0.24s; }
.loading-dots span:nth-child(4) { animation-delay: -0.20s; }
.loading-dots span:nth-child(5) { animation-delay: -0.16s; }
.loading-dots span:nth-child(6) { animation-delay: -0.12s; }
.loading-dots span:nth-child(7) { animation-delay: -0.08s; }

@keyframes loadingBounce {
  0%, 80%, 100% { 
    transform: scale(1);
    opacity: 1;
  }
  40% { 
    transform: scale(1.3);
    opacity: 0.7;
  }
}

/* Info Text */
.info-text {
  text-align: center;
  color: rgba(255, 255, 255, 0.5);
  font-size: 0.85rem;
  margin: 0;
  line-height: 1.5;
}

/* Footer */
.login-footer {
  position: relative;
  z-index: 10;
  margin-top: 3rem;
  text-align: center;
  color: rgba(255, 255, 255, 0.3);
  font-size: 0.8rem;
  letter-spacing: 0.1rem;
}

/* Responsive */
@media (max-width: 768px) {
  .login-card {
    padding: 2.5rem 2rem;
  }

  .neon-logo {
    font-size: 3.5rem;
    letter-spacing: 0.3rem;
  }

  .tagline {
    font-size: 0.8rem;
  }
}

@media (max-width: 480px) {
  .login-container {
    padding: 1rem;
  }

  .login-card {
    padding: 2rem 1.5rem;
  }

  .neon-logo {
    font-size: 3rem;
    letter-spacing: 0.2rem;
  }

  .neon-button {
    letter-spacing: 0.2rem;
    font-size: 1rem;
  }
}
</style>