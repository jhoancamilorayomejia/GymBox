<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const loginEmail    = ref('')
const loginPassword = ref('')
const loginError    = ref('')
const loginLoading  = ref(false)
const showPassword  = ref(false)

const loginAdmin = async () => {
  loginLoading.value = true
  loginError.value   = ''

  try {
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: loginEmail.value,
        password: loginPassword.value
      })
    })

    const data = await res.json()

    if (!res.ok) throw new Error(data.error || 'Usuario o contraseña incorrectos')

    localStorage.setItem('token',    data.token)
    localStorage.setItem('rol',      data.rol)
    localStorage.setItem('username', data.user)

    if (data.rol === 'admin') {
  router.push('/dashboard')
} else if (data.rol === 'customer') {
  router.push('/proyectos')
} else {
  router.push('/login')
}

  } catch (err) {
    loginError.value = err.message
  } finally {
    loginLoading.value = false
  }
}
</script>

<template>
  <div class="screen">
    <!-- Fondo geométrico animado -->
    <div class="bg-layer">
      <div class="hex hex-1"></div>
      <div class="hex hex-2"></div>
      <div class="hex hex-3"></div>
      <div class="slash slash-1"></div>
      <div class="slash slash-2"></div>
      <div class="slash slash-3"></div>
      <div class="noise"></div>
    </div>

    <!-- Panel izquierdo — branding -->
    <aside class="brand-panel">
      <div class="brand-inner">
        <!-- Logo SVG recreado desde la imagen -->
        <div class="logo-wrap">
          <svg class="logo-svg" viewBox="0 0 220 110" xmlns="http://www.w3.org/2000/svg">
            <!-- Rayo -->
            <polygon points="88,8 68,54 82,54 62,102 108,46 90,46 115,8" fill="#f5c500"/>
            <!-- RAYO -->
            <text x="110" y="50" font-family="'Barlow Condensed', sans-serif" font-weight="900"
                  font-size="44" fill="#f5c500" letter-spacing="-1">RAYO</text>
            <!-- BOX -->
            <text x="110" y="90" font-family="'Barlow Condensed', sans-serif" font-weight="900"
                  font-size="44" fill="#ffffff" letter-spacing="-1">BOX</text>
            <!-- CROSS LIFTING -->
            <text x="110" y="104" font-family="'Barlow Condensed', sans-serif" font-weight="400"
                  font-size="11" fill="#888" letter-spacing="4">CROSS LIFTING</text>
          </svg>
        </div>

        <div class="brand-tagline">
          <span class="tag-line">SISTEMA DE GESTIÓN</span>
          <h2>Entrená.<br/>Gestioná.<br/>Dominá.</h2>
        </div>

        <div class="stats-row">
          <div class="stat">
            <span class="stat-num">100%</span>
            <span class="stat-label">Digital</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat">
            <span class="stat-num">24/7</span>
            <span class="stat-label">Acceso</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat">
            <span class="stat-num">∞</span>
            <span class="stat-label">Poder</span>
          </div>
        </div>
      </div>
    </aside>

    <!-- Panel derecho — login -->
    <main class="login-panel">
      <div class="login-card">
        <!-- Esquinas decorativas -->
        <span class="corner tl"></span>
        <span class="corner br"></span>

        <div class="login-header">
          <div class="bolt-icon">⚡</div>
          <h1>Acceso al Sistema</h1>
          <p>Ingresá tus credenciales para continuar</p>
        </div>

        <form class="login-form" @submit.prevent="loginAdmin">

          <div class="field" :class="{ filled: loginEmail }">
            <label>Usuario</label>
            <div class="input-wrap">
              <svg class="input-icon" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                <circle cx="10" cy="7" r="3.5" stroke="currentColor" stroke-width="1.4"/>
                <path d="M3 17c0-3.314 3.134-6 7-6s7 2.686 7 6" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/>
              </svg>
              <input
                v-model="loginEmail"
                type="text"
                placeholder="Tu usuario"
                autocomplete="username"
                required
              />
            </div>
          </div>

          <div class="field" :class="{ filled: loginPassword }">
            <label>Contraseña</label>
            <div class="input-wrap">
              <svg class="input-icon" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                <rect x="4" y="9" width="12" height="9" rx="1.5" stroke="currentColor" stroke-width="1.4"/>
                <path d="M7 9V6.5a3 3 0 016 0V9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/>
              </svg>
              <input
                v-model="loginPassword"
                :type="showPassword ? 'text' : 'password'"
                placeholder="••••••••"
                autocomplete="current-password"
                required
              />
              <button
                type="button"
                class="eye-btn"
                @click="showPassword = !showPassword"
                :title="showPassword ? 'Ocultar' : 'Mostrar'"
              >
                <svg v-if="!showPassword" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M2 10s3-6 8-6 8 6 8 6-3 6-8 6-8-6-8-6z" stroke="currentColor" stroke-width="1.4"/>
                  <circle cx="10" cy="10" r="2.5" stroke="currentColor" stroke-width="1.4"/>
                </svg>
                <svg v-else viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M3 3l14 14M8.5 8.7A2.5 2.5 0 0013 12.3M6 6.1C3.9 7.4 2 10 2 10s3 6 8 6c1.6 0 3-.5 4.2-1.3M10 4c4.2.5 6.6 3.7 8 6" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/>
                </svg>
              </button>
            </div>
          </div>

          <transition name="shake">
            <div v-if="loginError" class="error-box">
              <span>⚡</span> {{ loginError }}
            </div>
          </transition>

          <button class="btn-login" type="submit" :disabled="loginLoading">
            <span class="btn-bg"></span>
            <span class="btn-label">
              <svg v-if="loginLoading" class="spinner" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3" stroke-dasharray="31.4" stroke-dashoffset="10"/>
              </svg>
              {{ loginLoading ? 'Ingresando...' : 'INGRESAR' }}
            </span>
          </button>

        </form>

        <p class="footer-note">
          Rayo Box · Sistema de Gestión v1.0
        </p>
      </div>
    </main>
  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Barlow+Condensed:wght@400;700;900&family=Barlow:wght@300;400;500&display=swap');

*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

/* ── Layout ── */
.screen {
  min-height: 100vh;
  display: flex;
  font-family: 'Barlow', sans-serif;
  background: #0a0a0a;
  position: relative;
  overflow: hidden;
}

/* ── Fondo geométrico ── */
.bg-layer { position: fixed; inset: 0; pointer-events: none; }

.hex, .slash { position: absolute; }

.hex-1 {
  width: 380px; height: 380px;
  background: conic-gradient(from 30deg, #f5c500 0deg, #ff7a00 80deg, transparent 80deg);
  clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);
  top: -120px; left: -100px; opacity: .12;
  animation: rotateSlow 30s linear infinite;
}
.hex-2 {
  width: 280px; height: 280px;
  background: conic-gradient(from 200deg, #ff7a00 0deg, #f5c500 90deg, transparent 90deg);
  clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);
  bottom: -80px; right: 38%; opacity: .10;
  animation: rotateSlow 22s linear infinite reverse;
}
.hex-3 {
  width: 180px; height: 180px;
  background: #f5c500;
  clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);
  top: 40%; right: 8%; opacity: .06;
  animation: rotateSlow 18s linear infinite;
}

.slash-1 {
  width: 3px; height: 60vh;
  background: linear-gradient(to bottom, transparent, #f5c500 50%, transparent);
  top: 20%; left: 37.5%; opacity: .18;
  transform: rotate(15deg);
}
.slash-2 {
  width: 2px; height: 40vh;
  background: linear-gradient(to bottom, transparent, #ff7a00 50%, transparent);
  top: 10%; left: 39%; opacity: .12;
  transform: rotate(15deg);
}
.slash-3 {
  width: 2px; height: 30vh;
  background: linear-gradient(to bottom, transparent, #f5c500 50%, transparent);
  bottom: 5%; left: 35%; opacity: .10;
  transform: rotate(15deg);
}

.noise {
  position: absolute; inset: 0;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='1'/%3E%3C/svg%3E");
  opacity: .025;
}

@keyframes rotateSlow { to { transform: rotate(360deg); } }

/* ── Brand Panel ── */
.brand-panel {
  width: 42%;
  min-height: 100vh;
  background: linear-gradient(145deg, #111 0%, #0a0a0a 60%);
  border-right: 1px solid rgba(245,197,0,.15);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px 48px;
  position: relative;
}

.brand-panel::after {
  content: '';
  position: absolute;
  top: 0; right: -1px;
  width: 1px; height: 100%;
  background: linear-gradient(to bottom, transparent, #f5c500 40%, #ff7a00 60%, transparent);
}

.brand-inner {
  display: flex;
  flex-direction: column;
  gap: 44px;
  max-width: 340px;
}

.logo-wrap {
  animation: logoIn .8s cubic-bezier(.16,1,.3,1) both;
}

.logo-svg {
  width: 240px;
  height: auto;
  filter: drop-shadow(0 0 28px rgba(245,197,0,.25));
}

@keyframes logoIn {
  from { opacity: 0; transform: translateX(-30px); }
  to   { opacity: 1; transform: translateX(0); }
}

.brand-tagline {
  animation: logoIn .8s .15s cubic-bezier(.16,1,.3,1) both;
}

.tag-line {
  display: block;
  font-family: 'Barlow Condensed', sans-serif;
  font-size: .65rem; letter-spacing: .35em;
  text-transform: uppercase; color: #f5c500;
  margin-bottom: 10px;
}

.brand-tagline h2 {
  font-family: 'Barlow Condensed', sans-serif;
  font-size: 3.4rem; font-weight: 900;
  line-height: 1; letter-spacing: -.01em;
  color: #fff;
  text-transform: uppercase;
}

.stats-row {
  display: flex;
  align-items: center;
  gap: 20px;
  animation: logoIn .8s .3s cubic-bezier(.16,1,.3,1) both;
}

.stat { display: flex; flex-direction: column; gap: 2px; }
.stat-num {
  font-family: 'Barlow Condensed', sans-serif;
  font-size: 1.7rem; font-weight: 900;
  color: #f5c500; line-height: 1;
}
.stat-label {
  font-size: .65rem; letter-spacing: .15em;
  text-transform: uppercase; color: #555;
}
.stat-divider {
  width: 1px; height: 36px;
  background: linear-gradient(to bottom, transparent, #333, transparent);
}

/* ── Login Panel ── */
.login-panel {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
}

.login-card {
  position: relative;
  width: 100%; max-width: 420px;
  background: #111;
  border: 1px solid rgba(245,197,0,.12);
  box-shadow: 0 40px 100px rgba(0,0,0,.8), 0 0 60px rgba(245,197,0,.04);
  padding: 48px 44px 40px;
  animation: cardIn .7s .1s cubic-bezier(.16,1,.3,1) both;
}

@keyframes cardIn {
  from { opacity: 0; transform: translateY(24px); }
  to   { opacity: 1; transform: translateY(0); }
}

/* Esquinas decorativas */
.corner {
  position: absolute;
  width: 18px; height: 18px;
  border-color: #f5c500; border-style: solid;
}
.corner.tl { top: 10px; left: 10px; border-width: 1.5px 0 0 1.5px; }
.corner.br { bottom: 10px; right: 10px; border-width: 0 1.5px 1.5px 0; }

/* ── Header login ── */
.login-header {
  text-align: center;
  margin-bottom: 36px;
}

.bolt-icon {
  font-size: 2rem; margin-bottom: 10px;
  filter: drop-shadow(0 0 12px #f5c500);
  display: block;
  animation: boltPulse 2.5s ease-in-out infinite;
}

@keyframes boltPulse {
  0%,100% { filter: drop-shadow(0 0 8px #f5c500); }
  50%      { filter: drop-shadow(0 0 22px #f5c500) drop-shadow(0 0 40px rgba(245,197,0,.4)); }
}

.login-header h1 {
  font-family: 'Barlow Condensed', sans-serif;
  font-size: 1.8rem; font-weight: 900;
  text-transform: uppercase;
  letter-spacing: .06em;
  color: #fff; margin-bottom: 6px;
}

.login-header p {
  font-size: .77rem; color: #555;
  letter-spacing: .05em;
}

/* ── Formulario ── */
.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 7px;
}

.field label {
  font-family: 'Barlow Condensed', sans-serif;
  font-size: .65rem; letter-spacing: .3em;
  text-transform: uppercase;
  color: #666; font-weight: 700;
  transition: color .2s;
}

.field.filled label,
.field:focus-within label {
  color: #f5c500;
}

.input-wrap {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 13px;
  width: 16px; height: 16px;
  color: #444;
  pointer-events: none;
  transition: color .2s;
}

.field:focus-within .input-icon { color: #f5c500; }

.input-wrap input {
  width: 100%;
  background: #0d0d0d;
  border: 1px solid #222;
  color: #f0f0f0;
  font-family: 'Barlow', sans-serif;
  font-size: .9rem; font-weight: 400;
  padding: 13px 14px 13px 40px;
  outline: none;
  transition: border-color .25s, background .25s, box-shadow .25s;
}

.input-wrap input::placeholder { color: #333; }

.input-wrap input:focus {
  border-color: #f5c500;
  background: #0f0e09;
  box-shadow: 0 0 0 3px rgba(245,197,0,.07), inset 0 1px 0 rgba(245,197,0,.06);
}

.eye-btn {
  position: absolute; right: 12px;
  background: none; border: none; cursor: pointer;
  width: 20px; height: 20px;
  color: #444; padding: 0;
  display: flex; align-items: center; justify-content: center;
  transition: color .2s;
}
.eye-btn:hover { color: #f5c500; }
.eye-btn svg { width: 16px; height: 16px; }

/* ── Error ── */
.error-box {
  background: rgba(220, 50, 30, .08);
  border: 1px solid rgba(220, 50, 30, .3);
  color: #e05a45;
  font-size: .78rem;
  padding: 10px 14px;
  display: flex; align-items: center; gap: 8px;
  letter-spacing: .03em;
}

.shake-enter-active { animation: shake .4s cubic-bezier(.36,.07,.19,.97); }
@keyframes shake {
  10%, 90% { transform: translateX(-2px); }
  20%, 80% { transform: translateX(3px); }
  30%, 50%, 70% { transform: translateX(-4px); }
  40%, 60% { transform: translateX(4px); }
}

/* ── Botón login ── */
.btn-login {
  position: relative;
  width: 100%;
  padding: 15px;
  background: transparent;
  border: 1.5px solid #f5c500;
  color: #f5c500;
  font-family: 'Barlow Condensed', sans-serif;
  font-size: .95rem; font-weight: 900;
  letter-spacing: .35em; text-transform: uppercase;
  cursor: pointer;
  overflow: hidden;
  margin-top: 4px;
  transition: color .3s;
  display: flex; align-items: center; justify-content: center;
}

.btn-bg {
  position: absolute; inset: 0;
  background: #f5c500;
  transform: scaleX(0); transform-origin: left;
  transition: transform .4s cubic-bezier(.4,0,.2,1);
}

.btn-login:hover .btn-bg,
.btn-login:focus .btn-bg { transform: scaleX(1); }
.btn-login:hover,
.btn-login:focus { color: #0a0a0a; outline: none; }
.btn-login:disabled { opacity: .5; cursor: not-allowed; }

.btn-label {
  position: relative; z-index: 1;
  display: flex; align-items: center; gap: 10px;
}

.spinner {
  width: 16px; height: 16px;
  animation: spin .7s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

/* ── Footer ── */
.footer-note {
  text-align: center;
  font-size: .65rem;
  color: #2e2e2e;
  letter-spacing: .1em;
  text-transform: uppercase;
  margin-top: 28px;
}

/* ── Responsive ── */
@media (max-width: 820px) {
  .screen { flex-direction: column; }
  .brand-panel {
    width: 100%; min-height: auto;
    padding: 48px 32px 36px;
    border-right: none;
    border-bottom: 1px solid rgba(245,197,0,.15);
  }
  .brand-inner { align-items: center; text-align: center; gap: 28px; }
  .brand-panel::after { display: none; }
  .stats-row { justify-content: center; }
  .brand-tagline h2 { font-size: 2.6rem; }
  .login-card { padding: 36px 26px 30px; }
}
</style>