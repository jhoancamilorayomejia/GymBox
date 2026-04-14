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

// Login
const loginEmail    = ref('')
const loginPassword = ref('')
const loginError    = ref('')
const loginLoading  = ref(false)

const loginAdmin = async () => {
  loginLoading.value = true
  loginError.value = ''

  try {
    const payload = {
      username: loginEmail.value,
      password: loginPassword.value
    }

    // ✅ Ruta relativa — funciona en desarrollo y producción
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })

    const data = await res.json()

    if (!res.ok) {
      throw new Error(data.error || 'Usuario o contraseña incorrectos')
    }

    // ✅ Guardar sesión
    localStorage.setItem('token', data.token)
    localStorage.setItem('rol', data.rol)
    localStorage.setItem('username', data.user)

    mostrarLogin.value = false

    // ✅ Redirección por rol
    if (data.rol === 'admin') {
      router.push('/api/proyectos')
    } else if (data.rol === 'proveedor') {
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

// ✅ Validación
const formValido = computed(() => {
  return nombre.value && whatsapp.value.length === 12 && fecha.value && hora.value
})

// ✅ Formatear WhatsApp y limitar a 10 dígitos
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
  if (formateado !== newVal) {
    whatsapp.value = formateado
  }
})

// ✅ Fechas disponibles (próximos 30 días excepto domingos)
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
  d.toLocaleDateString('es-AR', {
    weekday: 'long',
    day: 'numeric',
    month: 'long'
  }).replace(/^\w/, c => c.toUpperCase())

// ✅ Horas disponibles
const horasDisponibles = computed(() => {
  const horas = []
  for (let h = 6; h <= 20; h++) {
    for (let m of [0, 30]) {
      if (h === 20 && m === 30) continue
      const hh = String(h).padStart(2, '0')
      const mm = String(m).padStart(2, '0')
      horas.push(`${hh}:${mm}`)
    }
  }
  return horas
})

// Reset hora si cambia fecha
watch(fecha, () => {
  hora.value = ''
})

// ✅ PDF
const generarPDF = () => {
  const doc = new jsPDF()

  const fechaFormateada = new Date(`${fecha.value}T${hora.value}`)
    .toLocaleString('es-AR')

  doc.setFont("Helvetica", "bold")
  doc.setFontSize(16)
  doc.text("Confirmación de Reserva", 20, 20)

  doc.setFontSize(12)
  doc.setFont("Helvetica", "normal")

  doc.text(`Nombre: ${nombre.value}`, 20, 40)
  doc.text(`Teléfono: ${whatsapp.value}`, 20, 50)
  doc.text(`Fecha y hora: ${fechaFormateada}`, 20, 60)

  if (notas.value) {
    doc.text(`Notas: ${notas.value}`, 20, 70)
  }

  doc.text("Gracias por tu reserva 🙌", 20, 90)

  doc.save(`reserva-${nombre.value}.pdf`)
}

// ✅ Abrir modal
const enviar = () => {
  if (!formValido.value) {
    errorMsg.value = 'Completa todos los campos'
    return
  }
  errorMsg.value = ''
  mostrarModal.value = true
}

// ✅ Confirmar reserva
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
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify(payload)
    })

    if (res.status === 401) {
      localStorage.removeItem('token')
      router.push('/api/login')
      return
    }

    const text = await res.text()

    let data
    try {
      data = JSON.parse(text)
    } catch {
      throw new Error("Respuesta inválida del servidor: " + text)
    }

    if (!res.ok) {
      throw new Error(data.error || 'Error desconocido')
    }

    generarPDF()

    mostrarModal.value = false
    enviado.value = true

    setTimeout(() => {
      router.push('/api/proyectos')
    }, 2000)

  } catch (error) {
    errorMsg.value = error.message
  } finally {
    loading.value = false
  }
}

// ✅ Reset
const nuevaReserva = () => {
  nombre.value = ''
  whatsapp.value = ''
  fecha.value = ''
  hora.value = ''
  notas.value = ''
  enviado.value = false
  errorMsg.value = ''
}
</script>

<template>
  <div class="screen">
    <div class="bg-grid"></div>
    <div class="glow left"></div>
    <div class="glow right"></div>

    <div class="card">
      <header class="card-header">
        <span class="eyebrow">Reservá online</span>
        <h1>Agendar Turno</h1>
        <p class="subtitle">Completá el formulario y confirmamos tu turno al instante.</p>
        <div class="divider"><span></span><span class="diamond">◆</span><span></span></div>
      </header>

      <transition name="fade" mode="out-in">

        <!-- ── ÉXITO ── -->
        <div v-if="enviado" class="success-panel" key="ok">
          <div class="success-icon">✦</div>
          <h2>¡Turno solicitado!</h2>
          <p>Te contactaremos por WhatsApp para confirmar tu reserva.</p>
          <button class="btn-outline" @click="nuevaReserva">Hacer otra reserva</button>
        </div>

        <!-- ── FORMULARIO ── -->
        <form v-else key="form" class="form" @submit.prevent="enviar">

          <div class="field">
            <label>Escriba tu nombre completo<span class="req">*</span></label>
            <input v-model="nombre" type="text" required placeholder="Ej: Jhoan Camilo..." />
          </div>

          <div class="field">
            <label>WhatsApp <span class="req">*</span></label>
            <input
              v-model="whatsapp"
              type="tel"
              required
              placeholder="323 517 9341"
              maxlength="12"
            />
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
              >
                {{ formatearFecha(d) }}
              </option>
            </select>
          </div>

          <!-- Solo aparece si hay fecha -->
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
              >
                {{ h }}
              </button>
            </div>
          </div>

          <div class="field">
            <label>Notas adicionales</label>
            <textarea v-model="notas" rows="3" placeholder="Alguna preferencia o consulta especial..."></textarea>
          </div>

          <button type="submit" class="btn-primary">
            <span>Confirmar Turno</span>
          </button>

          <!-- Enlace para iniciar sesión -->
          <div class="login-link">
            <a href="#" @click.prevent="mostrarLogin = true">¿Eres Administrador? Iniciar sesión</a>
          </div>

        </form>
      </transition>
    </div>

    <!-- Modal confirmación reserva -->
    <div v-if="mostrarModal" class="modal-overlay">
      <div class="modal">
        <h3>Confirmar reserva</h3>
        <p>¿Deseas confirmar tu turno?</p>
        <div class="modal-actions">
          <button class="btn-outline" @click="mostrarModal = false">Cancelar</button>
          <button class="btn-primary" @click="confirmarReserva" :disabled="loading">
            <span>{{ loading ? 'Guardando...' : 'Confirmar' }}</span>
          </button>
        </div>
        <p v-if="errorMsg" style="color:#c0392b; font-size:.75rem; text-align:center;">
          {{ errorMsg }}
        </p>
      </div>
    </div>

    <!-- Modal de Login -->
    <div v-if="mostrarLogin" class="modal-overlay">
      <div class="modal">
        <h3>Iniciar sesión</h3>
        <div class="field">
          <label>Usuario</label>
          <input type="text" v-model="loginEmail" placeholder="Tu usuario" />
        </div>
        <div class="field">
          <label>Contraseña</label>
          <input type="password" v-model="loginPassword" placeholder="********" />
        </div>
        <div class="modal-actions">
          <button class="btn-outline" @click="mostrarLogin = false">Cancelar</button>
          <button class="btn-primary" @click="loginAdmin" :disabled="loginLoading">
            <span>{{ loginLoading ? 'Ingresando...' : 'Ingresar' }}</span>
          </button>
        </div>
        <p v-if="loginError" style="color:#c0392b; font-size:.75rem; text-align:center;">
          {{ loginError }}
        </p>
      </div>
    </div>

  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Cormorant+Garamond:wght@400;600;700&family=Montserrat:wght@300;400;500;600&display=swap');

*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

.screen {
  min-height: 100vh;
  background: #0d0d0d;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px 16px;
  font-family: 'Montserrat', sans-serif;
  position: relative;
  overflow: hidden;
}

.bg-grid {
  position: fixed; inset: 0;
  background-image:
    linear-gradient(rgba(180,145,80,.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(180,145,80,.03) 1px, transparent 1px);
  background-size: 60px 60px;
  pointer-events: none;
}

.glow {
  position: fixed;
  width: 500px; height: 500px;
  border-radius: 50%;
  filter: blur(120px);
  pointer-events: none;
  opacity: .3;
}
.glow.left  { background: radial-gradient(circle, #b4915025, transparent 70%); top: -100px; left: -150px; }
.glow.right { background: radial-gradient(circle, #b4915015, transparent 70%); bottom: -100px; right: -150px; }

/* ── Tarjeta ── */
.card {
  position: relative;
  width: 100%; max-width: 520px;
  background: #141414;
  border: 1px solid rgba(180,145,80,.2);
  box-shadow: 0 32px 80px rgba(0,0,0,.7), inset 0 1px 0 rgba(180,145,80,.08);
  padding: 48px 44px 44px;
  animation: cardIn .65s cubic-bezier(.16,1,.3,1) both;
}
.card::before, .card::after {
  content: ''; position: absolute;
  width: 20px; height: 20px;
  border-color: #b49150; border-style: solid;
}
.card::before { top: 12px; left: 12px; border-width: 1px 0 0 1px; }
.card::after  { bottom: 12px; right: 12px; border-width: 0 1px 1px 0; }

@keyframes cardIn {
  from { opacity: 0; transform: translateY(22px); }
  to   { opacity: 1; transform: translateY(0); }
}

/* ── Header ── */
.card-header { text-align: center; margin-bottom: 36px; }

.eyebrow {
  display: block;
  font-size: .63rem; letter-spacing: .3em;
  text-transform: uppercase; color: #b49150; margin-bottom: 10px;
}
.card-header h1 {
  font-family: 'Cormorant Garamond', serif;
  font-size: 2.6rem; font-weight: 600;
  color: #f0e6d0; letter-spacing: .08em;
  text-transform: uppercase; line-height: 1; margin-bottom: 10px;
}
.subtitle { color: #6a5c44; font-size: .77rem; font-weight: 300; letter-spacing: .04em; margin-bottom: 20px; }

.divider { display: flex; align-items: center; justify-content: center; gap: 10px; color: #b49150; }
.divider span:not(.diamond) { display: block; width: 60px; height: 1px; background: linear-gradient(90deg, transparent, #b49150); }
.divider span:last-child    { background: linear-gradient(90deg, #b49150, transparent); }
.diamond { font-size: .4rem; }

/* ── Campos ── */
.form { display: flex; flex-direction: column; gap: 22px; }

.field { display: flex; flex-direction: column; gap: 8px; }

.field label {
  font-size: .63rem; letter-spacing: .2em;
  text-transform: uppercase; color: #8a7455; font-weight: 500;
}
.req { color: #b49150; }
.hint { font-size: .63rem; color: #4a3f30; letter-spacing: .03em; }

.field input,
.field select,
.field textarea {
  background: #1c1c1c;
  border: 1px solid #2a2218;
  color: #e8dcc8;
  font-family: 'Montserrat', sans-serif;
  font-size: .85rem; font-weight: 300;
  padding: 13px 14px;
  outline: none;
  transition: border-color .25s, background .25s;
  resize: vertical;
}
.field input::placeholder,
.field textarea::placeholder { color: #3a3028; }
.field input:focus,
.field select:focus,
.field textarea:focus { border-color: #b49150; background: #181510; }

.field select {
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='8' viewBox='0 0 12 8'%3E%3Cpath d='M1 1l5 5 5-5' stroke='%23b49150' stroke-width='1.5' fill='none'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 14px center;
  background-color: #1c1c1c;
  padding-right: 38px;
  cursor: pointer;
}
.field select option { background: #1c1c1c; color: #e8dcc8; }

/* ── Botón principal ── */
.btn-primary {
  width: 100%; padding: 15px;
  background: transparent; border: 1px solid #b49150;
  color: #b49150;
  font-family: 'Montserrat', sans-serif;
  font-size: .7rem; font-weight: 600;
  letter-spacing: .3em; text-transform: uppercase;
  cursor: pointer;
  display: flex; align-items: center; justify-content: center;
  position: relative; overflow: hidden;
  transition: color .3s; margin-top: 4px;
}
.btn-primary::before {
  content: ''; position: absolute; inset: 0;
  background: #b49150;
  transform: scaleX(0); transform-origin: left;
  transition: transform .35s cubic-bezier(.4,0,.2,1);
}
.btn-primary:hover::before { transform: scaleX(1); }
.btn-primary:hover { color: #0d0d0d; }
.btn-primary span { position: relative; z-index: 1; }

/* ── Éxito ── */
.success-panel {
  display: flex; flex-direction: column;
  align-items: center; text-align: center;
  padding: 20px 0 10px; gap: 12px;
}
.success-icon { font-size: 2.5rem; color: #b49150; animation: pulse 2s ease-in-out infinite; }
@keyframes pulse { 0%,100% { opacity: 1; } 50% { opacity: .35; } }
.success-panel h2 {
  font-family: 'Cormorant Garamond', serif;
  font-size: 2rem; font-weight: 600; color: #f0e6d0; letter-spacing: .06em;
}
.success-panel p { color: #7a6a50; font-size: .82rem; font-weight: 300; }

.btn-outline {
  padding: 11px 28px; background: transparent;
  border: 1px solid #2e2820; color: #8a7455;
  font-family: 'Montserrat', sans-serif;
  font-size: .68rem; font-weight: 500;
  letter-spacing: .2em; text-transform: uppercase;
  cursor: pointer; margin-top: 8px;
  transition: border-color .2s, color .2s;
}
.btn-outline:hover { border-color: #b49150; color: #b49150; }

/* ── Transición ── */
.fade-enter-active, .fade-leave-active { transition: opacity .35s, transform .35s; }
.fade-enter-from { opacity: 0; transform: translateY(-8px); }
.fade-leave-to   { opacity: 0; transform: translateY(8px); }

/* ── Responsive ── */
@media (max-width: 560px) {
  .card { padding: 36px 22px 32px; }
  .card-header h1 { font-size: 2rem; }
}

/* ── Modal ── */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,.75);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: #141414;
  border: 1px solid rgba(180,145,80,.3);
  padding: 28px;
  width: 90%;
  max-width: 400px;
  text-align: center;
  box-shadow: 0 20px 60px rgba(0,0,0,.8);
}

.modal h3 {
  font-family: 'Cormorant Garamond', serif;
  color: #f0e6d0;
  margin-bottom: 10px;
}

.modal p {
  color: #8a7455;
  margin-bottom: 20px;
  font-size: .8rem;
}

.modal .field {
  text-align: left;
  margin-bottom: 14px;
}

.modal-actions {
  display: flex;
  gap: 10px;
  justify-content: center;
}

/* ── Horas ── */
.horas-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;
}

.hora-btn {
  background: #1c1c1c;
  border: 1px solid #2a2218;
  color: #e8dcc8;
  padding: 12px;
  font-size: .8rem;
  cursor: pointer;
  transition: all .25s;
  font-family: 'Montserrat', sans-serif;
}

.hora-btn:hover {
  border-color: #b49150;
  color: #b49150;
}

.hora-btn.active {
  background: #b49150;
  color: #0d0d0d;
  border-color: #b49150;
}

/* ── Login link ── */
.login-link {
  text-align: center;
  margin-top: 12px;
  font-size: .75rem;
}
.login-link a {
  color: #b49150;
  text-decoration: underline;
  cursor: pointer;
  transition: color .25s;
}
.login-link a:hover {
  color: #e8dcc8;
}
</style>