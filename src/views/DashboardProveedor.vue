<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router   = useRouter()
const username = localStorage.getItem('username') || 'Cliente'
const id       = localStorage.getItem('id') || ''

const planes      = ref([])
const loading     = ref(true)
const errorMsg    = ref('')
const viewDate    = ref(new Date())
const verVencidos = ref(false)

const MESES = ['Enero','Febrero','Marzo','Abril','Mayo','Junio','Julio','Agosto','Septiembre','Octubre','Noviembre','Diciembre']
const DIAS  = ['Dom','Lun','Mar','Mié','Jue','Vie','Sáb']

// ✅ ───────────── FIX GLOBAL DE FECHAS ─────────────
const parseFechaLocal = (fechaStr) => {
  const [year, month, day] = fechaStr.split('T')[0].split('-')
  return new Date(year, month - 1, day)
}

const normalizar = (fechaStr) => {
  const d = parseFechaLocal(fechaStr)
  d.setHours(0,0,0,0)
  return d
}
// ────────────────────────────────────────────────


// ── Auth ──────────────────────────────────────────

const tokenExpirado = (token) => {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.exp * 1000 < Date.now()
  } catch { return true }
}

const cerrarSesion = () => {
  localStorage.clear()
  router.push('/login')
}

// ── Fetch ─────────────────────────────────────────

const obtenerPlanes = async () => {
  loading.value = true; errorMsg.value = ''
  try {
    const token = localStorage.getItem('token')
    if (!token || tokenExpirado(token)) { router.push('/login'); return }

    const resC = await fetch(`/api/customers/by-cedula/${username}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!resC.ok) throw new Error('No se encontró el cliente')
    const cliente = await resC.json()

    const resP = await fetch(`/api/plans?idcustomer=${cliente.idcustomer}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!resP.ok) throw new Error('Error al cargar los planes')
    planes.value = await resP.json() || []
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}

// ── Clasificación ─────────────────────────────────

const hoyMidnight = () => {
  const h = new Date()
  h.setHours(0,0,0,0)
  return h
}

// ✅ Planes activos
const planesActivos = computed(() => {
  const hoy = hoyMidnight()
  return planes.value.filter(p => {
    if (p.state?.toLowerCase() !== 'pagado') return false
    const fin = normalizar(p.datefinish)
    return fin >= hoy
  }).sort((a, b) => normalizar(a.datestart) - normalizar(b.datestart))
})

// ✅ Planes vencidos
const planesVencidos = computed(() => {
  const hoy = hoyMidnight()
  return planes.value.filter(p => {
    if (p.state?.toLowerCase() !== 'pagado') return false
    const fin = normalizar(p.datefinish)
    return fin < hoy
  }).sort((a, b) => normalizar(b.datefinish) - normalizar(a.datefinish))
})

// ── Helpers ───────────────────────────────────────

const fmtFecha = (str) => {
  if (!str) return '—'
  const [year, month, day] = str.split('T')[0].split('-')
  return `${day}/${month}/${year}`
}

const fmtFechaLarga = (str) => {
  if (!str) return '—'
  const d = parseFechaLocal(str)
  return d.toLocaleDateString('es-CO', { day: 'numeric', month: 'long', year: 'numeric' })
}

const estadoClass = (s) => {
  if (!s) return 'badge-default'
  if (s.toLowerCase() === 'pagado') return 'badge-pagado'
  if (s.toLowerCase() === 'pendiente') return 'badge-pendiente'
  return 'badge-default'
}

// ✅ Disponibilidad fila
const dispFila = (p) => {
  const hoy = hoyMidnight()
  const ini = normalizar(p.datestart)
  const fin = normalizar(p.datefinish)

  if (hoy < ini) {
    const d = Math.ceil((ini - hoy) / 86400000)
    return { label: `Inicia en ${d}d`, clase: 'disp-proximo' }
  }

  if (hoy > fin) return { label: 'Vencido', clase: 'disp-vencido' }

  const d = Math.ceil((fin - hoy) / 86400000)
  if (d === 0) return { label: 'Vence hoy', clase: 'disp-vencer-hoy' }

  return { label: `${d}d restantes`, clase: 'disp-activo' }
}

// ── Resumen ───────────────────────────────────────

const resumenPlanes = computed(() => {
  const hoy = hoyMidnight()

  return planesActivos.value.map(p => {
    const ini = normalizar(p.datestart)
    const fin = normalizar(p.datefinish)

    const esFuturo = hoy < ini
    const diasRest = Math.max(0, Math.ceil((fin - hoy) / 86400000))

    let estado, mensaje

    if (esFuturo) {
      const d = Math.ceil((ini - hoy) / 86400000)
      estado  = 'proximo'
      mensaje = `Inicia en ${d} día${d !== 1 ? 's' : ''}, el ${fmtFechaLarga(p.datestart)}.`
    } else if (diasRest === 0) {
      estado  = 'vence-hoy'
      mensaje = 'Este plan vence HOY.'
    } else {
      estado  = 'activo'
      mensaje = `Acceso válido hasta el ${fmtFechaLarga(p.datefinish)}. Quedan ${diasRest} días.`
    }

    return { plan: p, estado, mensaje }
  })
})

// ── Calendario ────────────────────────────────────

const mesLabel = computed(() => `${MESES[viewDate.value.getMonth()]} ${viewDate.value.getFullYear()}`)
const primerDia = computed(() => new Date(viewDate.value.getFullYear(), viewDate.value.getMonth(), 1).getDay())
const diasEnMes = computed(() => new Date(viewDate.value.getFullYear(), viewDate.value.getMonth() + 1, 0).getDate())

const cambiarMes = (delta) => {
  const d = new Date(viewDate.value)
  d.setMonth(d.getMonth() + delta)
  viewDate.value = d
}

// ✅ CALENDARIO CORREGIDO
const clasificarDia = (dia) => {
  const y = viewDate.value.getFullYear()
  const m = viewDate.value.getMonth()

  const fecha = new Date(y, m, dia)
  fecha.setHours(0,0,0,0)

  const hoy = hoyMidnight()
  const esHoy = fecha.getTime() === hoy.getTime()

  let disponible = false
  let vencido = false

  for (const p of planes.value) {
    if (p.state?.toLowerCase() !== 'pagado') continue

    const ini = normalizar(p.datestart)
    const fin = normalizar(p.datefinish)

    if (fecha >= ini && fecha <= fin) {
      if (fin < hoy) {
        vencido = true
      } else {
        disponible = true
        break
      }
    }
  }

  return { disponible, vencido, esHoy }
}

onMounted(() => { obtenerPlanes() })
</script>

<template>
  <div class="screen">
    <div class="bg-layer">
      <div class="hex hex-1"></div><div class="hex hex-2"></div>
      <div class="slash slash-1"></div><div class="slash slash-2"></div>
      <div class="noise"></div>
    </div>

    <div class="layout">

      <!-- ── Sidebar ── -->
      <aside class="sidebar">
        <div class="sidebar-logo">
          <svg viewBox="0 0 160 80" xmlns="http://www.w3.org/2000/svg">
            <polygon points="36,4 22,38 32,38 18,72 52,32 38,32 58,4" fill="#f5c500"/>
            <text x="62" y="36" font-family="'Barlow Condensed',sans-serif" font-weight="900" font-size="32" fill="#f5c500" letter-spacing="-1">RAYO</text>
            <text x="62" y="66" font-family="'Barlow Condensed',sans-serif" font-weight="900" font-size="32" fill="#ffffff" letter-spacing="-1">BOX</text>
          </svg>
          <span class="sidebar-sub">CROSS LIFTING</span>
        </div>

        <nav class="sidebar-nav">
          <a class="nav-item active" href="#">
            <svg viewBox="0 0 20 20" fill="none"><rect x="3" y="3" width="14" height="14" rx="1.5" stroke="currentColor" stroke-width="1.5"/><path d="M7 10h6M10 7v6" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/></svg>
            Mis Planes
          </a>
        </nav>

        <div class="sidebar-footer">
          <div class="user-chip">
            <div class="user-avatar">{{ username.charAt(0).toUpperCase() }}</div>
            <div class="user-info">
              <span class="user-name">{{ username }}</span>
              <span class="user-role">ID #{{ id }} · Miembro</span>
            </div>
          </div>
          <button class="logout-btn" @click="cerrarSesion" title="Cerrar sesión">
            <svg viewBox="0 0 20 20" fill="none"><path d="M13 3h4v14h-4M9 14l4-4-4-4M13 10H5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
          </button>
        </div>
      </aside>

      <!-- ── Main ── -->
      <main class="main">

        <!-- Header -->
        <header class="page-header">
          <div class="page-title">
            <span class="page-eyebrow">⚡ Bienvenido de vuelta</span>
            <h1>{{ username }} <span class="id-tag">#{{ id }}</span></h1>
          </div>
          <div class="header-meta">
            <div class="stat-pill">
              <span class="stat-val">{{ planes.length }}</span>
              <span class="stat-lbl">Total</span>
            </div>
            <div class="stat-pill accent">
              <span class="stat-val">{{ planesActivos.length }}</span>
              <span class="stat-lbl">Vigentes</span>
            </div>
            <div class="stat-pill muted" v-if="planesVencidos.length">
              <span class="stat-val">{{ planesVencidos.length }}</span>
              <span class="stat-lbl">Vencidos</span>
            </div>
          </div>
        </header>

        <!-- Loading -->
        <div v-if="loading" class="state-box">
          <div class="bolt-spin">⚡</div>
          <p>Cargando tus planes...</p>
        </div>

        <!-- Error -->
        <div v-else-if="errorMsg" class="state-box error-state">
          <span>⚠</span><p>{{ errorMsg }}</p>
          <button class="retry-btn" @click="obtenerPlanes">Reintentar</button>
        </div>

        <template v-else>

          <!-- ── Sección planes activos/futuros ── -->
          <div class="top-row">

            <!-- Tarjetas de resumen — una por cada plan activo/futuro -->
            <div class="resumenes-col">

              <div v-if="resumenPlanes.length === 0" class="resumen-card resumen-sin-plan">
                <div class="resumen-top">
                  <span class="resumen-icon">🔒</span>
                  <span class="resumen-badge rbadge-sin-plan">SIN PLAN ACTIVO</span>
                </div>
                <p class="resumen-msg">No tienes un plan activo. Habla con tu entrenador.</p>
                <div class="leyenda">
                  <div class="ley-item"><span class="ley-dot dot-activo"></span><span>Día disponible</span></div>
                  <div class="ley-item"><span class="ley-dot dot-vencido-cal"></span><span>Plan vencido</span></div>
                  <div class="ley-item"><span class="ley-dot dot-ninguno"></span><span>Sin acceso</span></div>
                </div>
              </div>

              <div
                v-for="(r, i) in resumenPlanes"
                :key="r.plan.idplan"
                class="resumen-card"
                :class="`resumen-${r.estado}`"
                :style="{ animationDelay: i * 80 + 'ms' }"
              >
                <div class="resumen-top">
                  <span class="resumen-icon">
                    {{ r.estado === 'activo' ? '⚡' : r.estado === 'vence-hoy' ? '⏳' : '🗓' }}
                  </span>
                  <span class="resumen-badge" :class="`rbadge-${r.estado}`">
                    {{ r.estado === 'activo' ? 'ACCESO ACTIVO' : r.estado === 'vence-hoy' ? 'VENCE HOY' : 'PRÓXIMO INICIO' }}
                  </span>
                  <span class="plan-tipo-tag">{{ r.plan.typeplan?.toUpperCase() }}</span>
                </div>

                <p class="resumen-msg">{{ r.mensaje }}</p>

                <div class="resumen-fechas">
                  <div class="rfecha">
                    <span class="rfecha-key">Inicio</span>
                    <span class="rfecha-val">{{ fmtFecha(r.plan.datestart) }}</span>
                  </div>
                  <span class="rfecha-arrow">→</span>
                  <div class="rfecha">
                    <span class="rfecha-key">Fin</span>
                    <span class="rfecha-val">{{ fmtFecha(r.plan.datefinish) }}</span>
                  </div>
                  <div class="rfecha" v-if="r.estado === 'activo' || r.estado === 'vence-hoy'">
                    <span class="rfecha-key">Precio</span>
                    <span class="rfecha-val precio-tag">${{ Number(r.plan.price).toLocaleString('es-CO') }}</span>
                  </div>
                </div>

                <!-- Leyenda solo en la primera tarjeta -->
                <div class="leyenda" v-if="i === 0">
                  <div class="ley-item"><span class="ley-dot dot-activo"></span><span>Día disponible</span></div>
                  <div class="ley-item"><span class="ley-dot dot-vencido-cal"></span><span>Plan vencido</span></div>
                  <div class="ley-item"><span class="ley-dot dot-ninguno"></span><span>Sin acceso</span></div>
                </div>
              </div>

            </div>

            <!-- Calendario visual -->
            <div class="cal-card">
              <div class="cal-nav">
                <button class="cal-arrow" @click="cambiarMes(-1)">&#8249;</button>
                <span class="cal-mes">{{ mesLabel }}</span>
                <button class="cal-arrow" @click="cambiarMes(1)">&#8250;</button>
              </div>

              <div class="cal-header">
                <span v-for="d in DIAS" :key="d">{{ d }}</span>
              </div>

              <div class="cal-grid">
                <div v-for="n in primerDia" :key="'e'+n" class="cal-empty"></div>

                <div
                  v-for="dia in diasEnMes"
                  :key="dia"
                  class="cal-dia"
                  :class="{
                    'dia-activo':  clasificarDia(dia).disponible,
                    'dia-vencido': !clasificarDia(dia).disponible && clasificarDia(dia).vencido,
                    'dia-hoy':     clasificarDia(dia).esHoy,
                    'dia-ninguno': !clasificarDia(dia).disponible && !clasificarDia(dia).vencido
                  }"
                >
                  {{ dia }}
                </div>
              </div>
            </div>

          </div>

          <!-- ── Tabla de planes activos/futuros ── -->
          <div class="section-label">
            <svg viewBox="0 0 20 20" fill="none"><rect x="3" y="3" width="14" height="14" rx="1.5" stroke="currentColor" stroke-width="1.4"/><path d="M7 10h6M10 7v6" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/></svg>
            Planes vigentes
          </div>

          <div class="table-wrapper">
            <div class="table-scroll">
              <table>
                <thead>
                  <tr>
                    <th>#</th>
                    <th>Tipo</th>
                    <th>Precio</th>
                    <th>Fecha Pago</th>
                    <th>Inicio</th>
                    <th>Fin</th>
                    <th>Disponibilidad</th>
                    <th>Estado</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(p, i) in planesActivos" :key="p.idplan"
                      class="table-row" :style="{ animationDelay: i * 45 + 'ms' }">
                    <td><span class="id-badge">{{ p.idplan }}</span></td>
                    <td><span class="tipo-label">{{ p.typeplan }}</span></td>
                    <td><span class="precio-val">${{ Number(p.price).toLocaleString('es-CO') }}</span></td>
                    <td>{{ fmtFecha(p.datepay) }}</td>
                    <td>{{ fmtFecha(p.datestart) }}</td>
                    <td>{{ fmtFecha(p.datefinish) }}</td>
                    <td><span class="disp-badge" :class="dispFila(p).clase">{{ dispFila(p).label }}</span></td>
                    <td><span class="estado-badge" :class="estadoClass(p.state)">{{ p.state }}</span></td>
                  </tr>
                  <tr v-if="!planesActivos.length">
                    <td colspan="8" class="empty-row">⚡ No tienes planes vigentes. Habla con tu entrenador.</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div class="table-footer">
              {{ planesActivos.length }} plan{{ planesActivos.length !== 1 ? 'es' : '' }} vigente{{ planesActivos.length !== 1 ? 's' : '' }}
            </div>
          </div>

          <!-- ── Sección planes vencidos (toggle) ── -->
          <div class="vencidos-header">
            <div class="section-label">
              <svg viewBox="0 0 20 20" fill="none"><path d="M10 3a7 7 0 100 14A7 7 0 0010 3zM10 7v4l2.5 2.5" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round"/></svg>
              Planes vencidos
              <span class="venc-count" v-if="planesVencidos.length">{{ planesVencidos.length }}</span>
            </div>
            <button
              class="btn-ver-vencidos"
              @click="verVencidos = !verVencidos"
              v-if="planesVencidos.length"
            >
              <svg viewBox="0 0 20 20" fill="none">
                <path v-if="!verVencidos" d="M4 8l6 6 6-6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                <path v-else d="M4 12l6-6 6 6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              {{ verVencidos ? 'Ocultar' : 'Ver planes vencidos' }}
            </button>
            <span class="venc-vacio" v-else>Sin historial de planes vencidos</span>
          </div>

          <!-- Panel vencidos desplegable -->
          <transition name="expand">
            <div v-if="verVencidos && planesVencidos.length" class="vencidos-panel">

              <!-- Mini calendario vencidos (mismo mes que el principal, solo días grises) -->
              <div class="venc-top-row">

                <div class="venc-info-card">
                  <div class="section-label" style="margin-bottom:12px">
                    <svg viewBox="0 0 20 20" fill="none"><circle cx="10" cy="10" r="7" stroke="currentColor" stroke-width="1.4"/><path d="M10 7v4l2 2" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/></svg>
                    Historial anterior
                  </div>
                  <div class="venc-resumenes">
                    <div v-for="p in planesVencidos" :key="p.idplan" class="venc-item">
                      <div class="venc-tipo">{{ p.typeplan?.toUpperCase() }}</div>
                      <div class="venc-fechas-row">
                        <span>{{ fmtFecha(p.datestart) }}</span>
                        <span class="venc-arrow">→</span>
                        <span>{{ fmtFecha(p.datefinish) }}</span>
                      </div>
                      <div class="venc-precio">${{ Number(p.price).toLocaleString('es-CO') }}</div>
                    </div>
                  </div>
                </div>

                <!-- Calendario vencidos con días grises -->
                <div class="cal-card cal-vencidos">
                  <div class="cal-nav">
                    <button class="cal-arrow" @click="cambiarMes(-1)">&#8249;</button>
                    <span class="cal-mes">{{ mesLabel }}</span>
                    <button class="cal-arrow" @click="cambiarMes(1)">&#8250;</button>
                  </div>
                  <div class="cal-header">
                    <span v-for="d in DIAS" :key="'v'+d">{{ d }}</span>
                  </div>
                  <div class="cal-grid">
                    <div v-for="n in primerDia" :key="'ve'+n" class="cal-empty"></div>
                    <div
                      v-for="dia in diasEnMes"
                      :key="'vd'+dia"
                      class="cal-dia"
                      :class="{
                        'dia-vencido': clasificarDia(dia).vencido && !clasificarDia(dia).disponible,
                        'dia-ninguno': !clasificarDia(dia).vencido || clasificarDia(dia).disponible
                      }"
                    >
                      {{ dia }}
                    </div>
                  </div>
                </div>

              </div>

              <!-- Tabla vencidos -->
              <div class="table-wrapper table-vencidos">
                <div class="table-scroll">
                  <table>
                    <thead>
                      <tr>
                        <th>#</th>
                        <th>Tipo</th>
                        <th>Precio</th>
                        <th>Fecha Pago</th>
                        <th>Inicio</th>
                        <th>Fin</th>
                        <th>Estado pago</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="(p, i) in planesVencidos" :key="p.idplan"
                          class="table-row row-vencido" :style="{ animationDelay: i * 40 + 'ms' }">
                        <td><span class="id-badge id-badge-venc">{{ p.idplan }}</span></td>
                        <td><span class="tipo-label tipo-venc">{{ p.typeplan }}</span></td>
                        <td><span class="precio-val precio-venc">${{ Number(p.price).toLocaleString('es-CO') }}</span></td>
                        <td>{{ fmtFecha(p.datepay) }}</td>
                        <td>{{ fmtFecha(p.datestart) }}</td>
                        <td>{{ fmtFecha(p.datefinish) }}</td>
                        <td><span class="estado-badge" :class="estadoClass(p.state)">{{ p.state }}</span></td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <div class="table-footer">
                  {{ planesVencidos.length }} plan{{ planesVencidos.length !== 1 ? 'es' : '' }} vencido{{ planesVencidos.length !== 1 ? 's' : '' }} en el historial
                </div>
              </div>

            </div>
          </transition>

        </template>
      </main>
    </div>
  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Barlow+Condensed:wght@400;700;900&family=Barlow:wght@300;400;500&display=swap');
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

/* ── Screen / Layout ── */
.screen { width: 100vw; height: 100vh; overflow: hidden; background: #0a0a0a; font-family: 'Barlow', sans-serif; position: relative; }
.bg-layer { position: fixed; inset: 0; pointer-events: none; z-index: 0; }
.hex { position: absolute; }
.hex-1 { width: 420px; height: 420px; background: conic-gradient(from 30deg,rgba(245,197,0,.08) 0deg,rgba(255,122,0,.06) 80deg,transparent 80deg); clip-path: polygon(50% 0%,100% 25%,100% 75%,50% 100%,0% 75%,0% 25%); top: -100px; right: -80px; }
.hex-2 { width: 260px; height: 260px; background: conic-gradient(from 200deg,rgba(245,197,0,.05) 0deg,transparent 80deg); clip-path: polygon(50% 0%,100% 25%,100% 75%,50% 100%,0% 75%,0% 25%); bottom: -60px; left: 220px; }
.slash { position: absolute; }
.slash-1 { width: 2px; height: 60vh; background: linear-gradient(to bottom,transparent,rgba(245,197,0,.12) 50%,transparent); top: 20%; left: 220px; transform: rotate(12deg); }
.slash-2 { width: 1px; height: 40vh; background: linear-gradient(to bottom,transparent,rgba(255,122,0,.08) 50%,transparent); top: 30%; left: 230px; transform: rotate(12deg); }
.noise { position: absolute; inset: 0; background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='1'/%3E%3C/svg%3E"); opacity: .02; }

.layout { position: relative; z-index: 1; display: flex; width: 100%; height: 100vh; }

/* ── Sidebar ── */
.sidebar { width: 230px; height: 100vh; background: #0d0d0d; border-right: 1px solid rgba(245,197,0,.1); display: flex; flex-direction: column; padding: 32px 0 24px; flex-shrink: 0; }
.sidebar-logo { padding: 0 24px 28px; border-bottom: 1px solid rgba(245,197,0,.08); margin-bottom: 24px; }
.sidebar-logo svg { width: 130px; height: auto; filter: drop-shadow(0 0 18px rgba(245,197,0,.18)); }
.sidebar-sub { display: block; font-family: 'Barlow Condensed',sans-serif; font-size: .55rem; letter-spacing: .35em; text-transform: uppercase; color: #444; margin-top: 4px; }
.sidebar-nav { display: flex; flex-direction: column; gap: 4px; padding: 0 12px; flex: 1; }
.nav-item { display: flex; align-items: center; gap: 10px; padding: 11px 14px; color: #555; font-size: .82rem; font-weight: 500; letter-spacing: .04em; text-decoration: none; border-left: 2px solid transparent; transition: all .2s; cursor: pointer; }
.nav-item svg { width: 16px; height: 16px; flex-shrink: 0; }
.nav-item:hover { color: #f0f0f0; background: rgba(245,197,0,.05); }
.nav-item.active { color: #f5c500; border-left-color: #f5c500; background: rgba(245,197,0,.07); }
.sidebar-footer { padding: 20px 16px 0; border-top: 1px solid rgba(245,197,0,.08); display: flex; align-items: center; gap: 10px; margin-top: 12px; }
.user-chip { display: flex; align-items: center; gap: 10px; flex: 1; min-width: 0; }
.user-avatar { width: 32px; height: 32px; background: linear-gradient(135deg,#f5c500,#ff7a00); color: #0a0a0a; font-family: 'Barlow Condensed',sans-serif; font-weight: 900; font-size: .95rem; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.user-info { display: flex; flex-direction: column; min-width: 0; }
.user-name { font-size: .78rem; color: #e0e0e0; font-weight: 500; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.user-role { font-size: .6rem; color: #555; letter-spacing: .06em; }
.logout-btn { background: none; border: 1px solid #1e1e1e; color: #444; cursor: pointer; width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; transition: all .2s; }
.logout-btn svg { width: 15px; height: 15px; }
.logout-btn:hover { border-color: rgba(245,197,0,.4); color: #f5c500; }

/* ── Main ── */
.main { flex: 1; height: 100vh; overflow-y: auto; padding: 48px 60px; display: flex; flex-direction: column; gap: 22px; min-width: 0; }

/* ── Header ── */
.page-header { display: flex; align-items: flex-end; justify-content: space-between; flex-wrap: wrap; gap: 16px; flex-shrink: 0; }
.page-eyebrow { display: block; font-family: 'Barlow Condensed',sans-serif; font-size: .65rem; letter-spacing: .3em; text-transform: uppercase; color: #f5c500; margin-bottom: 6px; }
.page-title h1 { font-family: 'Barlow Condensed',sans-serif; font-size: 3rem; font-weight: 900; text-transform: uppercase; letter-spacing: .04em; color: #fff; line-height: 1; }
.id-tag { color: #f5c500; font-size: 2.2rem; }
.header-meta { display: flex; gap: 12px; }
.stat-pill { display: flex; flex-direction: column; align-items: center; background: #111; border: 1px solid rgba(245,197,0,.12); padding: 10px 22px; gap: 2px; }
.stat-pill.accent { border-color: rgba(245,197,0,.3); background: rgba(245,197,0,.06); }
.stat-pill.muted  { border-color: rgba(100,100,100,.2); }
.stat-val { font-family: 'Barlow Condensed',sans-serif; font-size: 1.7rem; font-weight: 900; color: #f5c500; line-height: 1; }
.stat-pill.muted .stat-val { color: #555; }
.stat-lbl { font-size: .58rem; letter-spacing: .2em; text-transform: uppercase; color: #555; }

/* ── States ── */
.state-box { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 14px; padding: 80px 20px; color: #555; font-size: .9rem; flex: 1; }
.bolt-spin { font-size: 2.4rem; animation: boltPulse 1.4s ease-in-out infinite; filter: drop-shadow(0 0 10px #f5c500); }
@keyframes boltPulse { 0%,100% { opacity: 1; } 50% { opacity: .3; } }
.error-state { color: #e05a45; }
.retry-btn { background: none; border: 1px solid rgba(220,50,30,.4); color: #e05a45; font-family: 'Barlow Condensed',sans-serif; font-size: .75rem; letter-spacing: .2em; text-transform: uppercase; padding: 8px 20px; cursor: pointer; transition: all .2s; }
.retry-btn:hover { background: rgba(220,50,30,.1); }

/* ── Top row ── */
.top-row { display: grid; grid-template-columns: 1fr 1.2fr; gap: 18px; flex-shrink: 0; }

/* ── Columna de resumenes (apilados) ── */
.resumenes-col { display: flex; flex-direction: column; gap: 12px; }

/* ── Resumen card ── */
.resumen-card {
  background: #111; border: 1px solid rgba(245,197,0,.12);
  padding: 18px 20px; display: flex; flex-direction: column; gap: 12px;
  animation: fadeUp .45s cubic-bezier(.16,1,.3,1) both;
  position: relative; overflow: hidden;
}
.resumen-card::before { content: ''; position: absolute; left: 0; top: 0; bottom: 0; width: 3px; }
.resumen-activo::before    { background: #4ade80; }
.resumen-vence-hoy::before { background: #f59e0b; }
.resumen-proximo::before   { background: #60a5fa; }
.resumen-sin-plan::before  { background: #3a3a3a; }

.resumen-top { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
.resumen-icon { font-size: 1.4rem; line-height: 1; }
.resumen-badge { font-family: 'Barlow Condensed',sans-serif; font-size: .68rem; font-weight: 900; letter-spacing: .25em; text-transform: uppercase; padding: 4px 12px; border: 1px solid; }
.rbadge-activo    { color: #4ade80; border-color: rgba(74,222,128,.4); background: rgba(74,222,128,.07); }
.rbadge-vence-hoy { color: #f59e0b; border-color: rgba(245,158,11,.4); background: rgba(245,158,11,.07); }
.rbadge-proximo   { color: #60a5fa; border-color: rgba(96,165,250,.4); background: rgba(96,165,250,.07); }
.rbadge-sin-plan  { color: #555; border-color: rgba(85,85,85,.3); background: rgba(85,85,85,.05); }

.plan-tipo-tag { font-family: 'Barlow Condensed',sans-serif; font-size: .65rem; font-weight: 700; letter-spacing: .2em; text-transform: uppercase; color: #555; margin-left: auto; }

.resumen-msg { font-size: .85rem; color: #aaa; line-height: 1.5; }

.resumen-fechas { display: flex; align-items: center; gap: 14px; padding: 8px 0; border-top: 1px solid rgba(245,197,0,.07); flex-wrap: wrap; }
.rfecha { display: flex; flex-direction: column; gap: 2px; }
.rfecha-key { font-size: .58rem; letter-spacing: .2em; text-transform: uppercase; color: #555; font-family: 'Barlow Condensed',sans-serif; font-weight: 700; }
.rfecha-val { font-size: .85rem; color: #e0e0e0; font-weight: 500; }
.rfecha-arrow { color: #f5c500; font-weight: 700; font-size: .9rem; }
.precio-tag { color: #f5c500 !important; }

/* ── Leyenda ── */
.leyenda { display: flex; flex-direction: column; gap: 6px; padding-top: 4px; border-top: 1px solid rgba(245,197,0,.07); }
.ley-item { display: flex; align-items: center; gap: 8px; font-size: .7rem; color: #666; letter-spacing: .03em; }
.ley-dot { width: 11px; height: 11px; flex-shrink: 0; border-radius: 2px; }
.dot-activo     { background: rgba(74,222,128,.25); border: 1px solid rgba(74,222,128,.6); }
.dot-vencido-cal { background: rgba(100,100,100,.25); border: 1px solid rgba(100,100,100,.5); }
.dot-ninguno    { background: rgba(255,255,255,.03); border: 1px solid rgba(255,255,255,.06); }

/* ── Calendario ── */
.cal-card { background: #111; border: 1px solid rgba(245,197,0,.12); padding: 20px 22px; display: flex; flex-direction: column; gap: 14px; animation: fadeUp .45s .06s cubic-bezier(.16,1,.3,1) both; }
.cal-vencidos { border-color: rgba(100,100,100,.2); }

.cal-nav { display: flex; align-items: center; justify-content: space-between; }
.cal-mes { font-family: 'Barlow Condensed',sans-serif; font-size: .95rem; font-weight: 700; letter-spacing: .12em; text-transform: uppercase; color: #f0f0f0; }
.cal-arrow { background: transparent; border: 1px solid rgba(245,197,0,.2); color: #f5c500; width: 28px; height: 28px; display: flex; align-items: center; justify-content: center; cursor: pointer; font-size: 1.2rem; line-height: 1; transition: all .2s; }
.cal-arrow:hover { background: rgba(245,197,0,.1); border-color: #f5c500; }

.cal-header { display: grid; grid-template-columns: repeat(7, 1fr); text-align: center; gap: 3px; }
.cal-header span { font-family: 'Barlow Condensed',sans-serif; font-size: .6rem; letter-spacing: .1em; text-transform: uppercase; color: #f5c500; padding: 4px 0; }
.cal-vencidos .cal-header span { color: #555; }

.cal-grid { display: grid; grid-template-columns: repeat(7, 1fr); gap: 3px; }
.cal-empty { aspect-ratio: 1; }

.cal-dia { aspect-ratio: 1; display: flex; align-items: center; justify-content: center; font-family: 'Barlow Condensed',sans-serif; font-size: .8rem; font-weight: 700; cursor: default; border-radius: 2px; transition: all .15s; }

.dia-ninguno { color: #2a2a2a; background: rgba(255,255,255,.02); border: 1px solid rgba(255,255,255,.03); }
.dia-activo  { color: #d1fae5; background: rgba(74,222,128,.18); border: 1px solid rgba(74,222,128,.35); }
.dia-activo:hover { background: rgba(74,222,128,.26); }
.dia-vencido { color: #555; background: rgba(100,100,100,.12); border: 1px solid rgba(100,100,100,.22); }
.dia-hoy     { background: #f5c500 !important; color: #0a0a0a !important; border: none !important; font-weight: 900; box-shadow: 0 0 12px rgba(245,197,0,.4); }

/* ── Section label ── */
.section-label { display: flex; align-items: center; gap: 8px; font-family: 'Barlow Condensed',sans-serif; font-size: .63rem; letter-spacing: .28em; text-transform: uppercase; color: #f5c500; font-weight: 700; flex-shrink: 0; }
.section-label svg { width: 14px; height: 14px; }

/* ── Tabla vigentes ── */
.table-wrapper { display: flex; flex-direction: column; border: 1px solid rgba(245,197,0,.12); background: #111; animation: fadeUp .45s .1s cubic-bezier(.16,1,.3,1) both; flex-shrink: 0; min-height: 80px; max-height: 260px; }
@keyframes fadeUp { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
.table-scroll { overflow-y: auto; flex: 1; }
table { width: 100%; border-collapse: collapse; font-size: .85rem; }
thead { background: rgba(245,197,0,.06); border-bottom: 1px solid rgba(245,197,0,.15); position: sticky; top: 0; z-index: 2; }
th { font-family: 'Barlow Condensed',sans-serif; font-size: .63rem; font-weight: 700; letter-spacing: .22em; text-transform: uppercase; color: #f5c500; padding: 13px 16px; text-align: left; white-space: nowrap; }
td { padding: 12px 16px; border-bottom: 1px solid rgba(255,255,255,.04); color: #ccc; vertical-align: middle; white-space: nowrap; }
.table-row { transition: background .15s; animation: rowIn .3s ease both; }
.table-row:hover { background: rgba(245,197,0,.04); }
.table-row:last-child td { border-bottom: none; }
@keyframes rowIn { from { opacity: 0; transform: translateX(-4px); } to { opacity: 1; transform: translateX(0); } }

.id-badge { display: inline-flex; align-items: center; justify-content: center; width: 26px; height: 26px; background: rgba(245,197,0,.07); border: 1px solid rgba(245,197,0,.2); font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .78rem; color: #f5c500; }
.tipo-label { font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .88rem; color: #f0f0f0; letter-spacing: .04em; text-transform: uppercase; }
.precio-val { font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .88rem; color: #f5c500; }

.estado-badge { display: inline-block; font-family: 'Barlow Condensed',sans-serif; font-size: .66rem; font-weight: 700; letter-spacing: .18em; text-transform: uppercase; padding: 4px 9px; border: 1px solid; }
.badge-pagado   { color: #4ade80; border-color: rgba(74,222,128,.4); background: rgba(74,222,128,.07); }
.badge-pendiente { color: #f59e0b; border-color: rgba(245,158,11,.4); background: rgba(245,158,11,.07); }
.badge-default  { color: #888; border-color: rgba(136,136,136,.3); background: rgba(136,136,136,.06); }

.disp-badge { display: inline-block; font-family: 'Barlow Condensed',sans-serif; font-size: .66rem; font-weight: 700; letter-spacing: .14em; text-transform: uppercase; padding: 4px 9px; border: 1px solid; }
.disp-activo     { color: #4ade80; border-color: rgba(74,222,128,.35); background: rgba(74,222,128,.06); }
.disp-vencer-hoy { color: #f59e0b; border-color: rgba(245,158,11,.35); background: rgba(245,158,11,.06); }
.disp-proximo    { color: #60a5fa; border-color: rgba(96,165,250,.35); background: rgba(96,165,250,.06); }
.disp-vencido    { color: #555; border-color: rgba(85,85,85,.3); background: rgba(85,85,85,.05); }

.table-footer { padding: 10px 16px; border-top: 1px solid rgba(245,197,0,.08); font-size: .68rem; color: #444; letter-spacing: .06em; flex-shrink: 0; }

/* ── Sección vencidos header ── */
.vencidos-header { display: flex; align-items: center; justify-content: space-between; gap: 12px; flex-shrink: 0; }
.venc-count { background: rgba(100,100,100,.15); border: 1px solid rgba(100,100,100,.25); color: #555; font-size: .6rem; font-weight: 700; letter-spacing: .1em; padding: 2px 8px; margin-left: 4px; border-radius: 2px; }
.venc-vacio { font-size: .72rem; color: #333; letter-spacing: .06em; }

.btn-ver-vencidos {
  display: inline-flex; align-items: center; gap: 7px;
  background: transparent; border: 1px solid rgba(100,100,100,.25);
  color: #555; font-family: 'Barlow Condensed',sans-serif;
  font-size: .72rem; font-weight: 700; letter-spacing: .18em; text-transform: uppercase;
  padding: 8px 16px; cursor: pointer; transition: all .2s;
}
.btn-ver-vencidos svg { width: 13px; height: 13px; flex-shrink: 0; }
.btn-ver-vencidos:hover { border-color: rgba(100,100,100,.5); color: #888; background: rgba(100,100,100,.06); }

/* ── Panel vencidos ── */
.vencidos-panel { display: flex; flex-direction: column; gap: 16px; border: 1px solid rgba(100,100,100,.15); background: rgba(100,100,100,.03); padding: 20px; }

.venc-top-row { display: grid; grid-template-columns: 1fr 1.2fr; gap: 18px; }

.venc-info-card { background: #0f0f0f; border: 1px solid rgba(100,100,100,.15); padding: 18px 20px; }

.venc-resumenes { display: flex; flex-direction: column; gap: 10px; }
.venc-item { padding: 10px 0; border-bottom: 1px solid rgba(100,100,100,.1); }
.venc-item:last-child { border-bottom: none; }
.venc-tipo { font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .8rem; letter-spacing: .1em; text-transform: uppercase; color: #555; margin-bottom: 4px; }
.venc-fechas-row { display: flex; align-items: center; gap: 8px; font-size: .8rem; color: #444; }
.venc-arrow { color: #333; }
.venc-precio { font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .8rem; color: #444; margin-top: 4px; }

/* Tabla vencidos */
.table-vencidos { border-color: rgba(100,100,100,.15); max-height: 200px; }
.table-vencidos thead { background: rgba(100,100,100,.06); border-bottom-color: rgba(100,100,100,.15); }
.table-vencidos th { color: #555; }
.row-vencido { opacity: .7; }
.row-vencido:hover { opacity: 1; background: rgba(100,100,100,.05) !important; }
.id-badge-venc { background: rgba(100,100,100,.08) !important; border-color: rgba(100,100,100,.2) !important; color: #555 !important; }
.tipo-venc { color: #555 !important; }
.precio-venc { color: #555 !important; }

/* ── Transición expand ── */
.expand-enter-active { transition: opacity .3s ease, transform .3s ease; }
.expand-leave-active { transition: opacity .2s ease, transform .2s ease; }
.expand-enter-from, .expand-leave-to { opacity: 0; transform: translateY(-8px); }

/* ── Responsive ── */
@media (max-width: 1100px) { .top-row, .venc-top-row { grid-template-columns: 1fr; } }
@media (max-width: 860px) {
  .sidebar { display: none; }
  .main { padding: 28px 20px 40px; }
}
</style>