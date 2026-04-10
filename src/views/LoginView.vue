<script setup>
import { ref, computed, watch } from 'vue'
import { jsPDF } from 'jspdf'
import { useRouter } from 'vue-router'

const router = useRouter()

const nombre   = ref('')
const whatsapp = ref('')
const fecha    = ref('')
const hora     = ref('')
const notas    = ref('')
const enviado  = ref(false)
const mostrarModal = ref(false)
const errorMsg = ref('')
const loading  = ref(false)
const mostrarLogin = ref(false)

const loginEmail    = ref('')
const loginPassword = ref('')
const loginError    = ref('')
const loginLoading  = ref(false)

const loginAdmin = async () => {
  loginLoading.value = true
  loginError.value = ''
  try {
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: loginEmail.value, password: loginPassword.value })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Usuario o contraseña incorrectos')
    localStorage.setItem('token', data.token)
    localStorage.setItem('rol', data.rol)
    localStorage.setItem('username', data.user)
    mostrarLogin.value = false
    if (data.rol === 'admin' || data.rol === 'proveedor') {
      router.push('/api/proyectos')
    } else {
      router.push('/api/login')
    }
  } catch (err) {
    loginError.value = err.message
  } finally {
    loginLoading.value = false
  }
}

const formValido = computed(() =>
  nombre.value && whatsapp.value.length === 12 && fecha.value && hora.value
)

const formatearWhatsApp = (value) => {
  const numeros = value.replace(/\D/g, '').slice(0, 10)
  const partes = []
  if (numeros.length > 0) partes.push(numeros.slice(0, 3))
  if (numeros.length >= 4) partes.push(numeros.slice(3, 6))
  if (numeros.length >= 7) partes.push(numeros.slice(6, 10))
  return partes.join(' ')
}

watch(whatsapp, (newVal) => {
  const formateado = formatearWhatsApp(newVal)
  if (formateado !== newVal) whatsapp.value = formateado
})

const fechasDisponibles = computed(() => {
  const dias = []
  const cursor = new Date()
  cursor.setHours(0, 0, 0, 0)
  while (dias.length < 30) {
    cursor.setDate(cursor.getDate() + 1)
    if (cursor.getDay() !== 0) dias.push(new Date(cursor))
  }
  return dias
})

const formatearFecha = (d) =>
  d.toLocaleDateString('es-CO', { weekday: 'long', day: 'numeric', month: 'long' })
    .replace(/^\w/, c => c.toUpperCase())

const DURACION_MINUTOS = 60

const horasDisponibles = computed(() => {
  const horas = []
  for (let h = 6; h <= 20; h++) {
    for (let m = 0; m < 60; m += DURACION_MINUTOS) {
      if (h === 20 && m > 0) continue
      horas.push(`${String(h).padStart(2,'0')}:${String(m).padStart(2,'0')}`)
    }
  }
  return horas
})

watch(fecha, () => { hora.value = '' })

const generarPDF = () => {
  const doc = new jsPDF()

  const fechaFormateada = new Date(`${fecha.value}T${hora.value}`)
    .toLocaleString('es-CO')

  // 🎨 COLORES
  const dorado = [180, 145, 80]
  const negro = [13, 13, 13]

  // 🧾 FONDO (opcional claro elegante)
  doc.setFillColor(245, 240, 230)
  doc.rect(0, 0, 210, 297, 'F')

  // 🧱 HEADER
  doc.setFont('Helvetica', 'bold')
  doc.setFontSize(22)
  doc.setTextColor(...negro)
  doc.text('Barber Charly', 105, 25, { align: 'center' })

  doc.setFontSize(12)
  doc.setTextColor(120)
  doc.text('Comprobante de Reserva', 105, 32, { align: 'center' })

  // 🔶 Línea decorativa
  doc.setDrawColor(...dorado)
  doc.setLineWidth(1)
  doc.line(40, 36, 170, 36)

  // 📦 CAJA DE DATOS
  doc.setDrawColor(200)
  doc.setLineWidth(0.5)
  doc.roundedRect(20, 50, 170, 80, 5, 5)

  // 📋 CONTENIDO
  let y = 65

  doc.setFontSize(12)
  doc.setTextColor(80)

  doc.text('Nombre:', 30, y)
  doc.setTextColor(...negro)
  doc.text(nombre.value, 80, y)

  y += 12
  doc.setTextColor(80)
  doc.text('Teléfono:', 30, y)
  doc.setTextColor(...negro)
  doc.text(whatsapp.value, 80, y)

  y += 12
  doc.setTextColor(80)
  doc.text('Fecha y hora:', 30, y)
  doc.setTextColor(...negro)
  doc.text(fechaFormateada, 80, y)

  if (notas.value) {
    y += 12
    doc.setTextColor(80)
    doc.text('Notas:', 30, y)
    doc.setTextColor(...negro)

    const splitNotas = doc.splitTextToSize(notas.value, 90)
    doc.text(splitNotas, 80, y)
  }

  // ✂️ MENSAJE FINAL
  doc.setFont('Helvetica', 'italic')
  doc.setFontSize(13)
  doc.setTextColor(...dorado)
  doc.text('✂ Gracias por confiar en nuestro estilo ✂', 105, 150, { align: 'center' })

  // 📅 FOOTER
  doc.setFontSize(9)
  doc.setTextColor(140)
  doc.text('Este documento es un comprobante de tu reserva.', 105, 280, { align: 'center' })

  // 💾 GUARDAR
  doc.save(`reserva-${nombre.value}.pdf`)
}

const enviar = () => {
  if (!formValido.value) { errorMsg.value = 'Completa todos los campos'; return }
  errorMsg.value = ''
  mostrarModal.value = true
}

const confirmarReserva = async () => {
  loading.value = true
  errorMsg.value = ''
  try {
    const payload = {
      name: nombre.value,
      phone: whatsapp.value,
      reservation_date: `${fecha.value} ${hora.value}`,
      note: notas.value
    }
    const token = localStorage.getItem('token')
    const res = await fetch('/api/reservations', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify(payload)
    })
    if (res.status === 401) { localStorage.removeItem('token'); router.push('/api/login'); return }
    const text = await res.text()
    let data
    try { data = JSON.parse(text) } catch { throw new Error('Respuesta inválida: ' + text) }
    if (!res.ok) throw new Error(data.error || 'Error desconocido')
    generarPDF()
    mostrarModal.value = false
    enviado.value = true
    setTimeout(() => router.push('/dashboard'), 2000)
  } catch (error) {
    errorMsg.value = error.message
  } finally {
    loading.value = false
  }
}

const nuevaReserva = () => {
  nombre.value = ''
  whatsapp.value = ''
  fecha.value = ''
  hora.value = ''
  notas.value = ''
  enviado.value = false
  errorMsg.value = ''
}

const scrollToForm = () => {
  document.getElementById('reservas')?.scrollIntoView({ behavior: 'smooth' })
}
</script>

<template>
  <div class="page">
    <div class="bg-grid"></div>
    <div class="glow left"></div>
    <div class="glow right"></div>

    <!-- ── HERO ── -->
    <section class="hero">
      <div class="hero-inner">
        <span class="eyebrow">Barbería &amp; Estilo</span>
        <h1 class="hero-title">
          Barber Charly
          <em>Tu estilo empieza donde termina lo común</em>
        </h1>
        <p class="hero-sub">Agendá tu turno en segundos.</p>
        <div class="hero-actions">
          <button class="btn-hero" @click="scrollToForm">✂ Agendar turno</button>
          <button class="btn-ghost" @click="mostrarLogin = true">Administrador</button>
        </div>
      </div>
    </section>

    <!-- ── FORMULARIO ── -->
    <section class="section-form" id="reservas">
      <div class="section-header">
        <span class="eyebrow">Reservá online</span>
        <h2 class="section-title">Agendar Turno</h2>
        <p class="section-sub">Completá el formulario y confirmamos tu turno al instante.</p>
        <div class="divider"><span></span><span class="diamond">◆</span><span></span></div>
      </div>

      <transition name="fade" mode="out-in">
        <div v-if="enviado" class="success-panel" key="ok">
          <div class="success-icon">✦</div>
          <h3>¡Turno solicitado!</h3>
          <p>Te contactaremos por WhatsApp para confirmar.</p>
          <button class="btn-main btn-outline-gold" @click="nuevaReserva">Hacer otra reserva</button>
        </div>

        <form v-else key="form" class="form" @submit.prevent="enviar">

          <div class="field">
            <label>Nombre completo <span class="req">*</span></label>
            <input v-model="nombre" type="text" required placeholder="Ej: Jhoan Camilo..." autocomplete="name" />
          </div>

          <div class="field">
            <label>WhatsApp <span class="req">*</span></label>
            <input v-model="whatsapp" type="tel" required placeholder="323 517 9341" maxlength="12" autocomplete="tel" />
            <span class="hint">Para confirmar o consultar tu turno</span>
          </div>

          <div class="field">
            <label>Fecha <span class="req">*</span></label>
            <select v-model="fecha" required>
              <option value="" disabled>Seleccioná una fecha</option>
              <option
                v-for="d in fechasDisponibles"
                :key="d.toISOString()"
                :value="d.toISOString().split('T')[0]"
              >{{ formatearFecha(d) }}</option>
            </select>
          </div>

          <div v-if="fecha" class="field">
            <label>Horario <span class="req">*</span></label>
            <div class="horas-grid">
              <button
                v-for="h in horasDisponibles"
                :key="h"
                type="button"
                class="hora-btn"
                :class="{ active: hora === h }"
                @click="hora = h"
              >{{ h }}</button>
            </div>
          </div>

          <div class="field">
            <label>Notas adicionales</label>
            <textarea v-model="notas" rows="3" placeholder="Alguna preferencia o consulta especial..."></textarea>
          </div>

          <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>

          <button type="submit" class="btn-main btn-gold">
            ✂ Confirmar Turno
          </button>
        </form>
      </transition>
    </section>

    <!-- ── FOOTER ── -->
    <footer class="footer">
      <p class="footer-text">© 2026 Barber Charly · Todos los derechos reservados.</p>
    </footer>

    <!-- ── MODAL CONFIRMACIÓN PANTALLA COMPLETA ── -->
    <div v-if="mostrarModal" class="fullscreen-modal">
      <div class="fs-inner">
        <div class="fs-icon">✦</div>
        <h2 class="fs-title">Confirmar reserva</h2>
        <p class="fs-sub">¿Confirmás tu turno?</p>
        <div class="fs-actions">
          <button class="btn-main btn-gold" @click="confirmarReserva" :disabled="loading">
            {{ loading ? 'Guardando...' : 'Sí, confirmar' }}
          </button>
          <button class="btn-main btn-outline-gold" @click="mostrarModal = false">Cancelar</button>
        </div>
        <p v-if="errorMsg" class="error-msg" style="margin-top:16px">{{ errorMsg }}</p>
      </div>
    </div>

    <!-- ── LOGIN PANTALLA COMPLETA ── -->
    <div v-if="mostrarLogin" class="fullscreen-modal">
      <div class="fs-inner">
        <button class="fs-close" @click="mostrarLogin = false">✕</button>
        <div class="fs-icon">🔐</div>
        <h2 class="fs-title">Administrador</h2>
        <p class="fs-sub">Ingresá tus credenciales</p>
        <div class="fs-form">
          
          <div class="field">
            <label>Usuario</label>
            <input type="text" v-model="loginEmail" placeholder="Tu usuario" autocomplete="username" />
          </div>
          <div class="field">
            <label>Contraseña</label>
            <input type="password" v-model="loginPassword" placeholder="••••••••" autocomplete="current-password" />
          </div>
          <p v-if="loginError" class="error-msg">{{ loginError }}</p>
          <button class="btn-main btn-gold" @click="loginAdmin" :disabled="loginLoading">
            {{ loginLoading ? 'Ingresando...' : 'Ingresar' }}
          </button>
          <button 
  class="btn-main btn-outline-gold" 
  @click="mostrarLogin = false"
>
  Cancelar
</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Cormorant+Garamond:ital,wght@0,400;0,600;1,400;1,600&family=Montserrat:wght@300;400;500;600&display=swap');

*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

/* ── BASE ── */
.page {
  min-height: 100vh;
  background: #0d0d0d;
  font-family: 'Montserrat', sans-serif;
  position: relative;
  overflow-x: hidden;
  color: #f0e6d0;
  margin: 0;
  padding: 0;
}



.bg-grid {
  position: fixed; inset: 0;
  background-image:
    linear-gradient(rgba(180,145,80,.025) 1px, transparent 1px),
    linear-gradient(90deg, rgba(180,145,80,.025) 1px, transparent 1px);
  background-size: 60px 60px;
  pointer-events: none; z-index: 0;
}

.glow {
  position: fixed; border-radius: 50%;
  filter: blur(140px); pointer-events: none; z-index: 0;
  width: 500px; height: 500px;
}
.glow.left  { background: radial-gradient(circle, rgba(180,145,80,.12), transparent 70%); top: -100px; left: -150px; }
.glow.right { background: radial-gradient(circle, rgba(180,145,80,.08), transparent 70%); bottom: -100px; right: -150px; }

/* ── HERO ── */
.hero {
  position: relative; z-index: 1;
  min-height: 50vh;
  display: flex; align-items: center; justify-content: center;
  text-align: center;
  padding: 60px 20px 48px;
  border-bottom: 1px solid rgba(180,145,80,.1);
}

.hero-inner { max-width: 680px; width: 100%; }

.eyebrow {
  display: inline-block;
  font-size: 26px; letter-spacing: .35em;
  text-transform: uppercase; color: #b49150;
  margin-bottom: 20px;
}

.hero-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: clamp(2.6rem, 10vw, 7rem);
  font-weight: 600; line-height: 1.1;
  color: #f0e6d0; margin-bottom: 16px;
  display: flex; flex-direction: column; gap: 8px;
}
.hero-title em {
  font-style: italic; color: #b49150;
  font-size: clamp(1.2rem, 4.5vw, 2.6rem);
}

.hero-sub {
  font-size: 26px; font-weight: 300;
  color: #6a5c44; letter-spacing: .04em;
  margin-bottom: 32px;
}

.hero-actions {
  display: flex; flex-direction: column;
  align-items: stretch; gap: 12px;
}

/* ── BOTONES ── */
.btn-main {
  width: 100%;
  /* Altura mínima de 56px — estándar táctil de apps móviles */
  min-height: 80px;
  padding: 0 20px;
  font-family: 'Montserrat', sans-serif;
  /* 16px mínimo iOS, usamos 17px para que se lea cómodo */
  font-size: 22px;
  font-weight: 600;
  letter-spacing: .1em;
  text-transform: uppercase;
  cursor: pointer;
  border: none;
  transition: all .2s;
  display: flex; align-items: center; justify-content: center;
  gap: 10px;
  -webkit-tap-highlight-color: transparent;
  touch-action: manipulation;
}
.btn-main:active { transform: scale(.97); }

.btn-gold { background: #b49150; color: #0d0d0d; }
.btn-gold:hover { background: #c9a862; }
.btn-gold:disabled { opacity: .6; cursor: not-allowed; }

.btn-outline-gold {
  background: transparent;
  color: #b49150;
  border: 2px solid rgba(180,145,80,.4) !important;
}
.btn-outline-gold:hover { border-color: #b49150 !important; }

.btn-hero {
  width: 128%; min-height: 85px;
  padding: 0 20px;
  background: #b49150; color: #0d0d0d;
  border: none;
  font-family: 'Montserrat', sans-serif;
  font-size: 17px; font-weight: 600;
  letter-spacing: .1em; text-transform: uppercase;
  cursor: pointer; transition: all .2s;
  -webkit-tap-highlight-color: transparent;
  touch-action: manipulation;
  display: flex; align-items: center; justify-content: center;
}

.btn-ghost {
  width: 128%; min-height: 85px;
  padding: 0 20px;
  background: transparent; color: #8a7455;
  border: 2px solid rgba(180,145,80,.22);
  font-family: 'Montserrat', sans-serif;
  font-size: 16px; font-weight: 400;
  letter-spacing: .08em; text-transform: uppercase;
  cursor: pointer; transition: all .2s;
  -webkit-tap-highlight-color: transparent;
  display: flex; align-items: center; justify-content: center;
}
.btn-ghost:hover { color: #b49150; border-color: rgba(180,145,80,.5); }

/* ── SECCIÓN FORM — ancho completo en móvil ── */
.section-form {
  position: relative;
  z-index: 1;

  width: 100%;
  max-width: 100%;   /* 🔥 evita que se encoja */
  margin: 0;         /* 🔥 elimina centrado */
  
  padding: 48px 20px 60px;
}

.section-header { text-align: center; margin-bottom: 36px; }

.section-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: clamp(3.5rem, 12vw, 5rem);
  font-weight: 600; color: #f0e6d0;
  letter-spacing: .06em; margin: 30px 0 34px;
}

.section-sub {
  font-size: 22px; font-weight: 340;
  color: #6a5c44; margin-bottom: 44px;
}

.divider {
  display: flex; align-items: center;
  justify-content: center; gap: 16px;
}
.divider span:not(.diamond) {
  width: 80px; height: 10px;
  background: linear-gradient(90deg, transparent, #b49150);
}
.divider span:last-child { background: linear-gradient(90deg, #b49150, transparent); }
.diamond { font-size: .4rem; color: #b49150; }

/* ── FORM ── */
.form { display: flex; flex-direction: column; gap: 44px; }

.field { display: flex; flex-direction: column; gap: 28px; }

.field label {
  /* 14px mínimo legible sin zoom */
  font-size: 24px;
  letter-spacing: .15em;
  text-transform: uppercase; color: #8a7455; font-weight: 540;
}
.req { color: #b49150; }
.hint { font-size: 24px; color: #4a3f30; }

/* ── INPUTS — 16px+ evita zoom automático en iOS ── */
.field input,
.field select,
.field textarea {
  background: #1a1a1a;
  border: 1px solid rgba(180,145,80,.2);
  color: #e8dcc8;
  font-family: 'Montserrat', sans-serif;
  font-weight: 360;
  /* 16px es el umbral exacto donde iOS deja de hacer zoom */
  font-size: 22px;
  /* Altura táctil generosa: mínimo 52px */
  min-height: 90px;
  padding: 20px 22px;
  outline: none;
  width: 100%;
  -webkit-appearance: none;
  border-radius: 0;
  transition: border-color .2s, background .2s;
  resize: none;
}
.field textarea { min-height: 90px; }
.field input::placeholder,
.field textarea::placeholder { color: #3a3028; }
.field input:focus,
.field select:focus,
.field textarea:focus {
  border-color: #b49150;
  background: #181510;
}

.field select {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='10' viewBox='0 0 16 10'%3E%3Cpath d='M1 1l7 7 7-7' stroke='%23b49150' stroke-width='2' fill='none'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 16px center;
  background-color: #1a1a1a;
  padding-right: 64px; cursor: pointer;
}
.field select option { background: #141414; color: #e8dcc8; }

/* ── HORAS — 3 columnas, bien grandes para el dedo ── */
.horas-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 32px;
}

.hora-btn {
  background: #1a1a1a;
  border: 1px solid rgba(180,145,80,.2);
  color: #6a5c44;
  min-height: 74px;
  padding: 0;
  font-size: 24px;
  cursor: pointer; transition: all .15s;
  font-family: 'Montserrat', sans-serif;
  -webkit-tap-highlight-color: transparent;
  touch-action: manipulation;
  border-radius: 0;
  display: flex; align-items: center; justify-content: center;
}
.hora-btn:active { transform: scale(.94); }
.hora-btn.active {
  background: #b49150; color: #0d0d0d;
  border-color: #b49150; font-weight: 640;
}

.error-msg {
  color: #e05555; font-size: 22px;
  text-align: center;
}

/* ── ÉXITO ── */
.success-panel {
  display: flex; flex-direction: column;
  align-items: center; text-align: center;
  padding: 40px 0; gap: 24px;
}
.success-icon { font-size: 3rem; color: #b49150; animation: pulse 2s ease-in-out infinite; }
@keyframes pulse { 0%,100%{opacity:1} 50%{opacity:.3} }
.success-panel h3 {
  font-family: 'Cormorant Garamond', serif;
  font-size: 2.6rem; font-weight: 640; color: #f0e6d0;
}
.success-panel p { color: #7a6a50; font-size: 16px; font-weight: 340; }

/* ── FOOTER ── */
.footer {
  text-align: center; padding: 48px 40px;
  border-top: 1px solid rgba(180,145,80,.08);
  position: relative; z-index: 1;
}
.footer-text { font-size: 13px; letter-spacing: .12em; text-transform: uppercase; color: #2e2820; }

/* ── PANTALLA COMPLETA ── */
.fullscreen-modal {
  position: fixed; inset: 0;
  background: #0d0d0d;
  z-index: 2080;
  display: flex; align-items: center; justify-content: center;
  padding: 60px 40px;
  overflow-y: auto;
}

.fs-inner {
  width: 90%;
  max-width: 90%;   /* 🔥 ocupa toda la pantalla */
  padding: 20px;     /* opcional para respirar */
}

@media (min-width: 660px) {
  .fs-inner { max-width: 900px; } /* 🔥 más ancho en desktop */
}


.fs-close {
  position: absolute; top: -4px; right: 0;
  background: transparent; border: none;
  color: #6a5c44; font-size: 1.6rem;
  cursor: pointer; padding: 10px;
  line-height: 1; transition: color .2s;
  -webkit-tap-highlight-color: transparent;
  min-height: 44px; min-width: 44px;
  display: flex; align-items: center; justify-content: center;
}
.fs-close:hover { color: #b49150; }

.fs-icon { font-size: 2.8rem; color: #b49150; margin-bottom: 16px; }

.fs-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: clamp(3.5rem, 12vw, 5rem);
  font-weight: 600; color: #f0e6d0;
  letter-spacing: .04em; margin-bottom: 10px;
}

.fs-sub {
  font-size: 16px; font-weight: 300;
  color: #6a5c44;
  margin-bottom: 32px;
}

.fs-actions {
  width: 100%;
  display: flex; flex-direction: column;
  gap: 22px;
}

.fs-form {
  width: 100%;
  display: flex; flex-direction: column;
  gap: 60px;
}
.fs-form .field { text-align: left; }

/* ── TRANSICIÓN ── */
.fade-enter-active, .fade-leave-active { transition: opacity .3s, transform .3s; }
.fade-enter-from { opacity: 0; transform: translateY(-8px); }
.fade-leave-to   { opacity: 0; transform: translateY(8px); }

/* ── DESKTOP ── */
@media (min-width: 640px) {
  .hero-actions { flex-direction: row; justify-content: center; }
  .btn-hero, .btn-ghost { width: auto; min-width: 260px; }
  .horas-grid { grid-template-columns: repeat(5, 1fr); }
}
</style>