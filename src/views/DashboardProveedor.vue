<template>
  <div class="screen">
    <div class="bg-grid"></div>
    <div class="glow left"></div>
    <div class="glow right"></div>

    <div class="card center-card">
      <div class="center-content">
        <span class="eyebrow">Iniciaste sesión - {{ user }}</span>

        <h1>Listado de Reservas</h1>

        <div class="divider">
          <span></span>
          <span class="diamond">◆</span>
          <span></span>
        </div>

        <!-- 📅 CALENDARIO POR MES -->
        <div class="month-nav">
          <button class="month-arrow" @click="cambiarMes(-1)">&#8249;</button>
          <span class="month-label">{{ mesActual }}</span>
          <button class="month-arrow" @click="cambiarMes(1)">&#8250;</button>
        </div>

        <!-- 🗓️ CALENDARIO -->
        <div class="calendar-wrapper">
          <div class="cal-header">
            <span v-for="dia in diasSemana" :key="dia">{{ dia }}</span>
          </div>
          <div class="cal-grid">
            <div
              v-for="n in primerDiaSemana"
              :key="'empty-' + n"
              class="day-cell empty"
            ></div>

            <button
              v-for="dia in diasEnMes"
              :key="dia"
              class="day-btn"
              :class="{
                'has-data': tieneReservas(dia),
                'selected': esSeleccionado(dia),
                'today': esHoy(dia)
              }"
              @click="seleccionarDia(dia)"
            >
              {{ dia }}
              <span v-if="tieneReservas(dia)" class="dot"></span>
            </button>
          </div>
        </div>

        <!-- ✅ PANEL DE RESERVAS -->
        <transition name="fade-slide">
          <div v-if="fechaSeleccionada" class="panel-section" ref="panelRef">

            <div class="panel-header">
              <span class="panel-title">Reservas del {{ fechaFormateada }}</span>
              <span class="badge-count">
                {{ reservasFiltradas.length }} reserva{{ reservasFiltradas.length !== 1 ? 's' : '' }}
              </span>
            </div>

            <div v-if="reservasFiltradas.length">
              <div
                v-for="r in reservasFiltradas"
                :key="r.id"
                class="reserva-card"
              >
                <!-- IZQUIERDA: nombre y teléfono -->
                <div class="rc-info">
                  <span class="rc-name">{{ r.name }}</span>
                  <span class="rc-phone">+57 {{ r.phone }}</span>
                </div>

                <!-- CENTRO: nota -->
                <div class="rc-note-wrap">
                  <span class="rc-note">{{ r.note || '—' }}</span>
                </div>

                <!-- DERECHA: hora + botón -->
                <div class="rc-side">
                  <span class="rc-time">
                    Hora: {{ r.reservation_date.length > 10 ? r.reservation_date.slice(11) : '—' }}
                  </span>
                  <button class="btn-delete" @click="abrirEliminar(r)">
                    Eliminar
                  </button>
                </div>
              </div>
            </div>

            <p v-else class="empty-state">No hay reservas para esta fecha.</p>
          </div>
        </transition>

        <!-- 🔑 BOTÓN CAMBIAR CLAVE -->
        <div class="btn-row" style="margin-top: 0">
          <button class="btn-outline-gold" @click="mostrarCambioPassword = true">
            Cambiar Clave
          </button>
        </div>

        <!-- 🔘 BOTONES -->
        <div class="btn-row">
          <button class="btn-primary" @click="cerrarSesion">
            <span>Cerrar sesión</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 🔴 MODAL ELIMINAR -->
    <div v-if="mostrarEliminar" class="modal-overlay">
      <div class="modal">
        <h3>Eliminar reserva</h3>
        <p>
          ¿Seguro que deseas eliminar la reserva de
          <strong>{{ reservaAEliminar?.name }}</strong>?
        </p>
        <div class="modal-actions">
          <button class="btn-outline" @click="mostrarEliminar = false">Cancelar</button>
          <button class="btn-delete" @click="eliminarReserva" :disabled="eliminando">
            {{ eliminando ? 'Eliminando...' : 'Eliminar' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 🔑 MODAL CAMBIAR CONTRASEÑA -->
    <div v-if="mostrarCambioPassword" class="modal-overlay">
      <div class="modal">
        <h3>Cambiar Contraseña</h3>

        <div class="modal-field">
          <label>Contraseña actual</label>
          <input
            v-model="passwordForm.oldPassword"
            type="password"
            placeholder="Ingresa tu contraseña actual"
            class="modal-input"
          />
        </div>

        <div class="modal-field">
          <label>Nueva contraseña</label>
          <input
            v-model="passwordForm.newPassword"
            type="password"
            placeholder="Mínimo 6 caracteres"
            class="modal-input"
          />
        </div>

        <div class="modal-field">
          <label>Confirmar nueva contraseña</label>
          <input
            v-model="passwordForm.confirmPassword"
            type="password"
            placeholder="Repite la nueva contraseña"
            class="modal-input"
          />
          <span v-if="passwordMismatch" class="error-text">Las contraseñas no coinciden</span>
        </div>

        <p v-if="passwordError" class="error-text" style="margin-bottom: 10px">{{ passwordError }}</p>
        <p v-if="passwordSuccess" class="success-text" style="margin-bottom: 10px">{{ passwordSuccess }}</p>

        <div class="modal-actions">
          <button class="btn-outline" @click="cerrarModalPassword">Cancelar</button>
          <button
            class="btn-primary-sm"
            :disabled="!canSubmitPassword || cambiandoPassword"
            @click="handleChangePassword"
          >
            {{ cambiandoPassword ? 'Guardando...' : 'Actualizar' }}
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const reservas = ref([])
const mostrarEliminar = ref(false)
const reservaAEliminar = ref(null)
const eliminando = ref(false)
const panelRef = ref(null)

const viewDate = ref(new Date())
const fechaSeleccionada = ref(null)
const user = ref(localStorage.getItem('username'))

const MESES = [
  'Enero', 'Febrero', 'Marzo', 'Abril', 'Mayo', 'Junio',
  'Julio', 'Agosto', 'Septiembre', 'Octubre', 'Noviembre', 'Diciembre'
]

const diasSemana = ['Dom', 'Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb']

// ── Helpers ──────────────────────────────────────────────────────

const fmtISO = (y, m, d) =>
  `${y}-${String(m + 1).padStart(2, '0')}-${String(d).padStart(2, '0')}`

const tokenExpirado = (token) => {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.exp * 1000 < Date.now()
  } catch {
    return true
  }
}

// ── Cambio de contraseña ──────────────────────────────────────────

const mostrarCambioPassword = ref(false)
const cambiandoPassword = ref(false)
const passwordError = ref('')
const passwordSuccess = ref('')
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const passwordMismatch = computed(() =>
  passwordForm.value.confirmPassword.length > 0 &&
  passwordForm.value.newPassword !== passwordForm.value.confirmPassword
)

const canSubmitPassword = computed(() =>
  passwordForm.value.oldPassword.length >= 1 &&
  passwordForm.value.newPassword.length >= 6 &&
  passwordForm.value.newPassword === passwordForm.value.confirmPassword
)

const cerrarModalPassword = () => {
  mostrarCambioPassword.value = false
  passwordForm.value = { oldPassword: '', newPassword: '', confirmPassword: '' }
  passwordError.value = ''
  passwordSuccess.value = ''
}

const handleChangePassword = async () => {
  cambiandoPassword.value = true
  passwordError.value = ''
  passwordSuccess.value = ''

  const username = localStorage.getItem('username')
  const token = localStorage.getItem('token')

  try {
    const res = await fetch('/api/users/change-password', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({
        username,
        old_password: passwordForm.value.oldPassword,
        new_password: passwordForm.value.newPassword,
        confirm_password: passwordForm.value.confirmPassword
      })
    })

    const data = await res.json()

    if (!res.ok) {
      passwordError.value = data.error || 'Error al cambiar la contraseña'
    } else {
      passwordSuccess.value = data.message
      passwordForm.value = { oldPassword: '', newPassword: '', confirmPassword: '' }
      setTimeout(() => cerrarModalPassword(), 1500)
    }
  } catch {
    passwordError.value = 'Error de conexión con el servidor'
  } finally {
    cambiandoPassword.value = false
  }
}

// ── Computadas del calendario ─────────────────────────────────────

const mesActual = computed(() => {
  const y = viewDate.value.getFullYear()
  const m = viewDate.value.getMonth()
  return `${MESES[m]} ${y}`
})

const primerDiaSemana = computed(() => {
  const d = new Date(viewDate.value.getFullYear(), viewDate.value.getMonth(), 1)
  return d.getDay()
})

const diasEnMes = computed(() => {
  const y = viewDate.value.getFullYear()
  const m = viewDate.value.getMonth()
  return new Date(y, m + 1, 0).getDate()
})

const soloFecha = (str) => str ? str.slice(0, 10) : ''

const fechasConReservas = computed(() =>
  new Set(reservas.value.map(r => soloFecha(r.reservation_date)))
)

const reservasFiltradas = computed(() => {
  if (!fechaSeleccionada.value) return []
  return reservas.value.filter(r => soloFecha(r.reservation_date) === fechaSeleccionada.value)
})

const fechaFormateada = computed(() => {
  if (!fechaSeleccionada.value) return ''
  const [y, m, d] = fechaSeleccionada.value.split('-')
  return `${parseInt(d)} de ${MESES[parseInt(m) - 1]} de ${y}`
})

// ── Métodos del calendario ────────────────────────────────────────

const tieneReservas = (dia) => {
  const y = viewDate.value.getFullYear()
  const m = viewDate.value.getMonth()
  return fechasConReservas.value.has(fmtISO(y, m, dia))
}

const esSeleccionado = (dia) => {
  const y = viewDate.value.getFullYear()
  const m = viewDate.value.getMonth()
  return fechaSeleccionada.value === fmtISO(y, m, dia)
}

const esHoy = (dia) => {
  const hoy = new Date()
  const y = viewDate.value.getFullYear()
  const m = viewDate.value.getMonth()
  return (
    hoy.getFullYear() === y &&
    hoy.getMonth() === m &&
    hoy.getDate() === dia
  )
}

const seleccionarDia = async (dia) => {
  const y = viewDate.value.getFullYear()
  const m = viewDate.value.getMonth()
  const iso = fmtISO(y, m, dia)

  if (fechaSeleccionada.value === iso) {
    fechaSeleccionada.value = null
    return
  }

  fechaSeleccionada.value = iso

  await nextTick()
  if (panelRef.value) {
    panelRef.value.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
  }
}

const cambiarMes = (delta) => {
  const d = new Date(viewDate.value)
  d.setMonth(d.getMonth() + delta)
  viewDate.value = d
  fechaSeleccionada.value = null
}

// ── Sesión ────────────────────────────────────────────────────────

const cerrarSesion = () => {
  const confirmar = confirm('¿Seguro que deseas cerrar sesión?')
  if (!confirmar) return

  localStorage.removeItem('token')
  localStorage.removeItem('rol')
  localStorage.removeItem('username')
  router.push('/api/login')
}

// ── Fetch de reservas ─────────────────────────────────────────────

const obtenerReservas = async () => {
  try {
    const token = localStorage.getItem('token')

    if (!token || tokenExpirado(token)) {
      router.push('/api/login')
      return
    }

    const res = await fetch('/api/proyectos', {
      headers: { Authorization: `Bearer ${token}` }
    })

    if (res.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('rol')
      localStorage.removeItem('username')
      router.push('/api/login')
      return
    }

    reservas.value = await res.json()
  } catch (error) {
    console.error('Error cargando reservas:', error)
  }
}

const abrirEliminar = (reserva) => {
  reservaAEliminar.value = reserva
  mostrarEliminar.value = true
}

const eliminarReserva = async () => {
  if (!reservaAEliminar.value) return

  eliminando.value = true

  try {
    const token = localStorage.getItem('token')

    const res = await fetch(`/api/reservations/${reservaAEliminar.value.id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` }
    })

    const data = await res.json()
    if (!res.ok) throw new Error(data.error)

    await obtenerReservas()
    mostrarEliminar.value = false
    reservaAEliminar.value = null
  } catch (err) {
    alert(err.message)
  } finally {
    eliminando.value = false
  }
}

onMounted(() => {
  obtenerReservas()
})
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Cormorant+Garamond:wght@400;600;700&family=Montserrat:wght@300;400;500;600&display=swap');

*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

.screen {
  min-height: 100vh;
  background: #0d0d0d;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 16px;
  font-family: 'Montserrat', sans-serif;
  position: relative;
  overflow: auto;
}

.bg-grid,
.glow { pointer-events: none; }

.bg-grid {
  position: fixed; inset: 0;
  background-image:
    linear-gradient(rgba(180,145,80,.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(180,145,80,.03) 1px, transparent 1px);
  background-size: 60px 60px;
}

.glow {
  position: fixed;
  width: 500px; height: 500px;
  border-radius: 50%;
  filter: blur(120px);
  opacity: .3;
}
.glow.left  { background: radial-gradient(circle, #b4915025, transparent 70%); top: -100px; left: -150px; }
.glow.right { background: radial-gradient(circle, #b4915015, transparent 70%); bottom: -100px; right: -150px; }

/* ── Card principal ── */
.card {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 700px;
  background: #141414;
  border: 1px solid rgba(180,145,80,.2);
  box-shadow: 0 32px 80px rgba(0,0,0,.7);
  padding: 60px 40px;
  text-align: center;
}

.center-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.eyebrow {
  font-size: .65rem;
  letter-spacing: .3em;
  text-transform: uppercase;
  color: #b49150;
}

h1 {
  font-family: 'Cormorant Garamond', serif;
  font-size: 2.3rem;
  color: #f0e6d0;
  letter-spacing: .08em;
}

.divider {
  display: flex;
  align-items: center;
  gap: 10px;
}
.divider span:not(.diamond) {
  width: 60px; height: 1px;
  background: linear-gradient(90deg, transparent, #b49150);
}
.divider span:last-child {
  background: linear-gradient(90deg, #b49150, transparent);
}
.diamond { font-size: .4rem; color: #b49150; }

/* ── Navegación de mes ── */
.month-nav {
  display: flex;
  align-items: center;
  gap: 14px;
}

.month-arrow {
  background: transparent;
  border: 1px solid rgba(180,145,80,.3);
  color: #b49150;
  width: 34px; height: 34px;
  cursor: pointer;
  font-size: 1.2rem;
  transition: all .25s;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}
.month-arrow:hover {
  background: #b49150;
  color: #0d0d0d;
}

.month-label {
  font-size: 1rem;
  letter-spacing: .18em;
  text-transform: uppercase;
  color: #f0e6d0;
  min-width: 170px;
  text-align: center;
}

/* ── Calendario compacto ── */
.calendar-wrapper {
  width: 100%;
  background: rgba(255,255,255,.03);
  border: 1px solid rgba(180,145,80,.12);
  padding: 16px;
}

.cal-header {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
  margin-bottom: 8px;
}
.cal-header span {
  font-size: .72rem;
  letter-spacing: .12em;
  text-transform: uppercase;
  color: #b49150;
  text-align: center;
  padding: 4px 0;
}

.cal-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
}

.day-cell.empty { aspect-ratio: 1; }

.day-btn {
  aspect-ratio: 1;
  background: transparent;
  border: 1px solid transparent;
  color: #555;
  font-size: .85rem;
  cursor: pointer;
  transition: all .2s;
  font-family: 'Montserrat', sans-serif;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  width: 100%;
}
.day-btn:hover {
  border-color: rgba(180,145,80,.4);
  color: #f0e6d0;
}

.day-btn.has-data {
  color: #f0e6d0;
  border-color: rgba(180,145,80,.2);
}

.day-btn.today {
  box-shadow: 0 0 0 1px rgba(180,145,80,.5);
  color: #f0e6d0;
}

.day-btn.selected {
  background: #b49150;
  color: #0d0d0d;
  border-color: #b49150;
  font-weight: 600;
}
.day-btn.selected .dot { background: rgba(0,0,0,.4); }

.dot {
  width: 4px; height: 4px;
  border-radius: 50%;
  background: #b49150;
  flex-shrink: 0;
}

/* ── Panel de reservas ── */
.panel-section {
  width: 100%;
  border-top: 1px solid rgba(180,145,80,.2);
  padding-top: 20px;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14px;
}

.panel-title {
  font-size: .68rem;
  letter-spacing: .2em;
  text-transform: uppercase;
  color: #b49150;
}

.badge-count {
  font-size: .65rem;
  letter-spacing: .1em;
  padding: 4px 12px;
  border: 1px solid rgba(180,145,80,.3);
  color: #b49150;
  background: rgba(180,145,80,.08);
}

/* ── Tarjeta de reserva ── */
.reserva-card {
  display: grid;
  grid-template-columns: 1fr 2fr auto;
  gap: 12px;
  align-items: center;
  background: rgba(255,255,255,.03);
  border: 1px solid rgba(180,145,80,.12);
  padding: 14px 16px;
  margin-bottom: 8px;
  text-align: left;
  transition: background .2s;
}
.reserva-card:hover {
  background: rgba(180,145,80,.06);
}
.reserva-card:last-child { margin-bottom: 0; }

/* Columna izquierda */
.rc-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}
.rc-name {
  font-size: .88rem;
  color: #f0e6d0;
  font-weight: 500;
  line-height: 1.3;
  word-break: break-word;
}
.rc-phone {
  font-size: .75rem;
  color: #666;
}

/* Columna central: nota */
.rc-note-wrap {
  min-width: 0;
}
.rc-note {
  font-size: .78rem;
  color: #777;
  line-height: 1.5;
  word-break: break-word;
}

/* Columna derecha: hora + botón */
.rc-side {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 8px;
  flex-shrink: 0;
}
.rc-time {
  font-size: .95rem;
  color: #b49150;
  font-weight: 600;
  white-space: nowrap;
}

.empty-state {
  text-align: center;
  padding: 2rem 1rem;
  color: #444;
  font-size: .82rem;
  letter-spacing: .1em;
}

/* ── Botones ── */
.btn-row {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  justify-content: center;
  width: 100%;
}

.btn-primary {
  padding: 13px 36px;
  border: 1px solid #b49150;
  background: #b49150;
  color: #0d0d0d;
  letter-spacing: .3em;
  text-transform: uppercase;
  cursor: pointer;
  transition: all .3s;
  font-family: 'Montserrat', sans-serif;
  font-size: .9rem;
  width: 100%;
  min-height: 60px;
}
.btn-primary:hover {
  background: #c9a96e;
}

.btn-outline-gold {
  padding: 13px 36px;
  border: 1px solid #b49150;
  background: transparent;
  color: #b49150;
  letter-spacing: .3em;
  text-transform: uppercase;
  cursor: pointer;
  transition: all .3s;
  font-family: 'Montserrat', sans-serif;
  font-size: .9rem;
  width: 100%;
  min-height: 60px;
}
.btn-outline-gold:hover {
  background: rgba(180,145,80,.1);
}

.btn-delete {
  padding: 8px 14px;
  border: 1px solid rgba(192,57,43,.5);
  background: transparent;
  color: #c0392b;
  font-size: .75rem;
  font-family: 'Montserrat', sans-serif;
  letter-spacing: .08em;
  cursor: pointer;
  transition: all .2s;
  white-space: nowrap;
}
.btn-delete:hover {
  background: #c0392b;
  color: #fff;
}
.btn-delete:disabled {
  opacity: .5;
  cursor: not-allowed;
}

/* ── Transición ── */
.fade-slide-enter-active {
  transition: opacity .3s ease, transform .3s ease;
}
.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(8px);
}

/* ── Modal ── */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,.85);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 999;
}

.modal {
  background: #141414;
  border: 1px solid rgba(180,145,80,.25);
  padding: 36px 30px;
  width: 100%;
  max-width: 420px;
  text-align: center;
}

.modal h3 {
  font-family: 'Cormorant Garamond', serif;
  font-size: 1.5rem;
  color: #f0e6d0;
  letter-spacing: .06em;
  margin-bottom: 12px;
}

.modal p {
  color: #8a7455;
  font-size: .88rem;
  margin-bottom: 24px;
  line-height: 1.6;
}

.modal p strong {
  color: #f0e6d0;
}

.modal-actions {
  display: flex;
  gap: 10px;
  justify-content: center;
}

.modal-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 16px;
  text-align: left;
}

.modal-field label {
  font-size: .65rem;
  letter-spacing: .15em;
  text-transform: uppercase;
  color: #b49150;
}

.modal-input {
  background: #1a1a1a;
  border: 1px solid rgba(180,145,80,.25);
  color: #f0e6d0;
  padding: 10px 14px;
  font-family: 'Montserrat', sans-serif;
  font-size: .9rem;
  outline: none;
  transition: border-color .2s;
}
.modal-input:focus {
  border-color: #b49150;
}

.btn-outline {
  padding: 10px 22px;
  border: 1px solid rgba(180,145,80,.3);
  color: #b49150;
  background: transparent;
  font-family: 'Montserrat', sans-serif;
  font-size: .78rem;
  letter-spacing: .1em;
  text-transform: uppercase;
  cursor: pointer;
  transition: all .2s;
}
.btn-outline:hover {
  border-color: #b49150;
  background: rgba(180,145,80,.08);
}

.btn-primary-sm {
  padding: 10px 24px;
  border: 1px solid #b49150;
  background: #b49150;
  color: #0d0d0d;
  font-family: 'Montserrat', sans-serif;
  font-size: .78rem;
  letter-spacing: .15em;
  text-transform: uppercase;
  cursor: pointer;
  transition: all .2s;
}
.btn-primary-sm:disabled {
  opacity: .4;
  cursor: not-allowed;
}
.btn-primary-sm:not(:disabled):hover {
  background: #c9a96e;
}

.error-text {
  color: #c0392b;
  font-size: .78rem;
}
.success-text {
  color: #27ae60;
  font-size: .78rem;
}

/* ── Responsive ── */
@media (max-width: 468px) {
  .card {
    max-width: 95vw;
    padding: 40px 16px;
  }

  h1 {
    font-size: 1.5rem;
  }

  .btn-primary,
  .btn-outline-gold {
    width: 100%;
    min-height: 56px;
  }

  .reserva-card {
    grid-template-columns: 1fr;
    gap: 8px;
  }

  .rc-note-wrap {
    border-top: 1px solid rgba(180,145,80,.08);
    padding-top: 6px;
  }

  .rc-side {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    border-top: 1px solid rgba(180,145,80,.08);
    padding-top: 6px;
  }
}
</style>