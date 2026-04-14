<template>
  <div class="screen">
    <div class="bg-grid"></div>
    <div class="glow left"></div>
    <div class="glow right"></div>

    <div class="card center-card">
      <div class="center-content">
        <span class="eyebrow">Iniciaste sesión</span>

        <h1>Ver Listado de Reservas</h1>

        <div class="divider">
          <span></span>
          <span class="diamond">◆</span>
          <span></span>
        </div>

        <!-- 📅 NAVEGACIÓN DE MES -->
        <div class="month-nav">
          <button class="month-arrow" @click="cambiarMes(-1)">&#8249;</button>
          <span class="month-label">{{ mesActual }}</span>
          <button class="month-arrow" @click="cambiarMes(1)">&#8250;</button>
        </div>

        <!-- 🗓️ CALENDARIO -->
        <div class="calendar">
          <div class="cal-header">
            <span v-for="dia in diasSemana" :key="dia">{{ dia }}</span>
          </div>
          <div class="cal-grid">
            <!-- Espacios vacíos antes del primer día -->
            <div
              v-for="n in primerDiaSemana"
              :key="'empty-' + n"
              class="day-cell empty"
            ></div>

            <!-- Días del mes -->
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

        <!-- ✅ TABLA (visible solo al seleccionar fecha) -->
        <transition name="fade-slide">
          <div v-if="fechaSeleccionada" class="table-section">
            <p class="date-title">
              Reservas del {{ fechaFormateada }}
            </p>

            <div class="table-container" v-if="reservasFiltradas.length">
              <table>
                <thead>
                  <tr>
                    <th>Nombre</th>
                    <th>WhatsApp</th>
                    <th>Hora</th>
                    <th>Nota</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="r in reservasFiltradas" :key="r.id">
                    <td>{{ r.name }}</td>
                    <td>{{ r.phone }}</td>
                    <td>{{ r.reservation_date.length > 10 ? r.reservation_date.slice(11) : '—' }}</td>
                    <td>{{ r.note || '—' }}</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <p v-else class="empty">No hay reservas para esta fecha.</p>
          </div>
        </transition>

        <!-- 🔘 BOTONES -->
        <div class="btn-row">
          <button class="btn-primary" @click="volver">
            <span>Volver</span>
          </button>
          <button class="btn-primary" @click="cerrarSesion">
            <span>Cerrar sesión</span>
          </button>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const reservas = ref([])

// Estado del calendario
const viewDate = ref(new Date())
const fechaSeleccionada = ref(null) // 'YYYY-MM-DD'

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

// ── Computadas del calendario ─────────────────────────────────────

const mesActual = computed(() => {
  const y = viewDate.value.getFullYear()
  const m = viewDate.value.getMonth()
  return `${MESES[m]} ${y}`
})

const primerDiaSemana = computed(() => {
  const d = new Date(viewDate.value.getFullYear(), viewDate.value.getMonth(), 1)
  return d.getDay() // 0 = Dom
})

const diasEnMes = computed(() => {
  const y = viewDate.value.getFullYear()
  const m = viewDate.value.getMonth()
  return new Date(y, m + 1, 0).getDate()
})

// Extrae solo 'YYYY-MM-DD' de un campo que puede venir como '2026-04-07 16:00'
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

const seleccionarDia = (dia) => {
  const y = viewDate.value.getFullYear()
  const m = viewDate.value.getMonth()
  fechaSeleccionada.value = fmtISO(y, m, dia)
}

const cambiarMes = (delta) => {
  const d = new Date(viewDate.value)
  d.setMonth(d.getMonth() + delta)
  viewDate.value = d
  fechaSeleccionada.value = null
}

// ── Sesión ────────────────────────────────────────────────────────

const cerrarSesion = (mensaje = 'Sesión cerrada') => {
  localStorage.removeItem('token')
  localStorage.removeItem('rol')
  localStorage.removeItem('username')
  alert(mensaje)
  router.push('/api/login')
}

const volver = () => {
  router.push('/dashboard')
}

// ── Fetch de reservas ─────────────────────────────────────────────

const obtenerReservas = async () => {
  try {
    const token = localStorage.getItem('token')

    // Si no hay token o expiró, redirige sin alert ni borrar storage
    if (!token || tokenExpirado(token)) {
      router.push('/api/login')
      return
    }

    const res = await fetch('/api/proyectos', {
      headers: { Authorization: `Bearer ${token}` }
    })

    // Solo limpia sesión si el backend confirma que el token es inválido
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
  padding: 60px 16px;
  font-family: 'Montserrat', sans-serif;
  position: relative;
  overflow: hidden;
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
  margin-top: 6px;
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
  font-size: .78rem;
  letter-spacing: .18em;
  text-transform: uppercase;
  color: #f0e6d0;
  min-width: 170px;
  text-align: center;
}

/* ── Calendario ── */
.calendar { width: 100%; }

.cal-header {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
  margin-bottom: 6px;
}
.cal-header span {
  font-size: .6rem;
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

.day-cell.empty { background: transparent; }

.day-btn {
  background: transparent;
  border: 1px solid rgba(180,145,80,.12);
  color: #666;
  font-size: .78rem;
  padding: 8px 4px 4px;
  cursor: pointer;
  transition: all .2s;
  font-family: 'Montserrat', sans-serif;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  min-height: 42px;
  justify-content: center;
}
.day-btn:hover {
  border-color: rgba(180,145,80,.5);
  color: #f0e6d0;
}

.day-btn.has-data {
  color: #f0e6d0;
  border-color: rgba(180,145,80,.35);
}

.day-btn.today {
  box-shadow: 0 0 0 1px rgba(180,145,80,.45);
  color: #f0e6d0;
}

.day-btn.selected {
  background: #b49150;
  color: #0d0d0d;
  border-color: #b49150;
  font-weight: 600;
}
.day-btn.selected .dot { background: #0d0d0d; }

.dot {
  width: 4px; height: 4px;
  border-radius: 50%;
  background: #b49150;
  flex-shrink: 0;
}

/* ── Tabla ── */
.table-section { width: 100%; }

.date-title {
  font-size: .68rem;
  letter-spacing: .22em;
  text-transform: uppercase;
  color: #b49150;
  margin-bottom: 14px;
}

.table-container {
  width: 100%;
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  color: #f0e6d0;
  font-size: 0.83rem;
}

thead { background: rgba(180,145,80,.1); }

th, td {
  padding: 10px 12px;
  border-bottom: 1px solid rgba(180,145,80,.2);
  text-align: center;
}

th {
  color: #b49150;
  font-weight: 500;
  letter-spacing: .1em;
  font-size: .68rem;
  text-transform: uppercase;
}

tbody tr:hover { background: rgba(180,145,80,.08); }

.empty {
  color: #555;
  font-size: .82rem;
  margin-top: 6px;
}

/* ── Botones ── */
.btn-row {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  justify-content: center;
  margin-top: 8px;
}

.btn-primary {
  padding: 13px 36px;
  border: 1px solid #b49150;
  background: transparent;
  color: #b49150;
  letter-spacing: .3em;
  text-transform: uppercase;
  cursor: pointer;
  transition: all .3s;
  font-family: 'Montserrat', sans-serif;
  font-size: .68rem;
}
.btn-primary:hover {
  background: #b49150;
  color: #0d0d0d;
}

/* ── Transición ── */
.fade-slide-enter-active {
  transition: opacity .3s ease, transform .3s ease;
}
.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(8px);
}
</style>