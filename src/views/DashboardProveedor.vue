<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const username = localStorage.getItem('username') || 'Cliente'
const id       = localStorage.getItem('id') || ''

const planes   = ref([])
const loading  = ref(true)
const errorMsg = ref('')

// ── Auth ──────────────────────────────────────────────────────────

const tokenExpirado = (token) => {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.exp * 1000 < Date.now()
  } catch { return true }
}

const cerrarSesion = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('rol')
  localStorage.removeItem('username')
  localStorage.removeItem('id')
  router.push('/login')
}

// ── Fetch planes ──────────────────────────────────────────────────

const obtenerPlanes = async () => {
  loading.value  = true
  errorMsg.value = ''
  try {
    const token = localStorage.getItem('token')
    if (!token || tokenExpirado(token)) { router.push('/login'); return }

    const resC = await fetch(`/api/customers/by-cedula/${username}`, {
  headers: { Authorization: `Bearer ${token}` }
})
const cliente = await resC.json()
const resP = await fetch(`/api/plans?idcustomer=${cliente.idcustomer}`, {
  headers: { Authorization: `Bearer ${token}` }
})
const data = await resP.json()
planes.value = data || []
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}

// ── Helpers ───────────────────────────────────────────────────────


const fmtFecha = (str) => {
  if (!str) return '—'
  const d = new Date(str)
  return isNaN(d) ? str : d.toLocaleDateString('es-CO', { day: '2-digit', month: 'short', year: 'numeric' })
}

const estadoClass = (s) => {
  if (!s) return 'badge-default'
  if (s.toLowerCase() === 'pagado')    return 'badge-pagado'
  if (s.toLowerCase() === 'pendiente') return 'badge-pendiente'
  return 'badge-default'
}

// Disponibilidad calculada por fila según datestart y datefinish
const dispFila = (p) => {
  const hoy = new Date(); hoy.setHours(0, 0, 0, 0)
  const ini = new Date(p.datestart)
  const fin = new Date(p.datefinish)

  if (hoy < ini) {
    const d = Math.ceil((ini - hoy) / 86400000)
    return { label: `Inicia en ${d}d`, clase: 'disp-proximo' }
  }
  if (hoy > fin) return { label: 'Vencido', clase: 'disp-vencido' }
  const d = Math.ceil((fin - hoy) / 86400000)
  if (d === 0) return { label: 'Vence hoy', clase: 'disp-por-vencer' }
  if (d <= 5)  return { label: `${d}d restantes`, clase: 'disp-por-vencer' }
  return { label: `${d}d restantes`, clase: 'disp-activo' }
}

const planesActivos = computed(() =>
  planes.value.filter(p => dispFila(p).clase === 'disp-activo' || dispFila(p).clase === 'disp-por-vencer').length
)

onMounted(() => { obtenerPlanes() })
</script>

<template>
  <div class="screen">

    <div class="bg-layer">
      <div class="hex hex-1"></div>
      <div class="hex hex-2"></div>
      <div class="slash slash-1"></div>
      <div class="slash slash-2"></div>
      <div class="noise"></div>
    </div>

    <div class="layout">

      <!-- ── Sidebar ── -->
      <aside class="sidebar">
        <div class="sidebar-logo">
          <svg viewBox="0 0 160 80" xmlns="http://www.w3.org/2000/svg">
            <polygon points="36,4 22,38 32,38 18,72 52,32 38,32 58,4" fill="#f5c500"/>
            <text x="62" y="36" font-family="'Barlow Condensed',sans-serif" font-weight="900"
                  font-size="32" fill="#f5c500" letter-spacing="-1">RAYO</text>
            <text x="62" y="66" font-family="'Barlow Condensed',sans-serif" font-weight="900"
                  font-size="32" fill="#ffffff" letter-spacing="-1">BOX</text>
          </svg>
          <span class="sidebar-sub">CROSS LIFTING</span>
        </div>

        <nav class="sidebar-nav">
          <a class="nav-item active" href="#">
            <svg viewBox="0 0 20 20" fill="none">
              <rect x="3" y="3" width="14" height="14" rx="1.5" stroke="currentColor" stroke-width="1.5"/>
              <path d="M7 10h6M10 7v6" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/>
            </svg>
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
              <span class="stat-lbl">Total planes</span>
            </div>
            <div class="stat-pill accent">
              <span class="stat-val">{{ planesActivos }}</span>
              <span class="stat-lbl">Vigentes</span>
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
          <span>⚠</span>
          <p>{{ errorMsg }}</p>
          <button class="retry-btn" @click="obtenerPlanes">Reintentar</button>
        </div>

        <!-- Tabla -->
        <div v-else class="table-wrapper">
          <div class="table-scroll">
            <table>
              <thead>
                <tr>
                  <th>#</th>
                  <th>Tipo de Plan</th>
                  <th>Precio</th>
                  <th>Fecha Pago</th>
                  <th>Inicio</th>
                  <th>Fin</th>
                  <th>Disponibilidad</th>
                  <th>Estado</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="(p, i) in planes"
                  :key="p.idplan"
                  class="table-row"
                  :style="{ animationDelay: i * 45 + 'ms' }"
                >
                  <td><span class="id-badge">{{ p.idplan }}</span></td>
                  <td><span class="tipo-label">{{ p.typeplan }}</span></td>
                  <td><span class="precio-val">${{ Number(p.price).toLocaleString('es-CO') }}</span></td>
                  <td>{{ fmtFecha(p.datepay) }}</td>
                  <td>{{ fmtFecha(p.datestart) }}</td>
                  <td>{{ fmtFecha(p.datefinish) }}</td>
                  <td>
                    <span class="disp-badge" :class="dispFila(p).clase">
                      {{ dispFila(p).label }}
                    </span>
                  </td>
                  <td>
                    <span class="estado-badge" :class="estadoClass(p.state)">{{ p.state }}</span>
                  </td>
                </tr>

                <tr v-if="!planes.length">
                  <td colspan="8" class="empty-row">
                    ⚡ No tienes planes registrados aún. Habla con tu entrenador para comenzar.
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <div class="table-footer">
            Mostrando <strong>{{ planes.length }}</strong>
            plan{{ planes.length !== 1 ? 'es' : '' }} de la cuenta
            <strong style="color:#f5c500"> {{ username }} #{{ id }}</strong>
          </div>
        </div>

      </main>
    </div>
  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Barlow+Condensed:wght@400;700;900&family=Barlow:wght@300;400;500&display=swap');

*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

/* ── Pantalla completa ── */
.screen { width: 100vw; height: 100vh; overflow: hidden; background: #0a0a0a; font-family: 'Barlow', sans-serif; position: relative; }

/* ── Fondo ── */
.bg-layer { position: fixed; inset: 0; pointer-events: none; z-index: 0; }
.hex { position: absolute; }
.hex-1 { width: 420px; height: 420px; background: conic-gradient(from 30deg,rgba(245,197,0,.08) 0deg,rgba(255,122,0,.06) 80deg,transparent 80deg); clip-path: polygon(50% 0%,100% 25%,100% 75%,50% 100%,0% 75%,0% 25%); top: -100px; right: -80px; }
.hex-2 { width: 260px; height: 260px; background: conic-gradient(from 200deg,rgba(245,197,0,.05) 0deg,transparent 80deg); clip-path: polygon(50% 0%,100% 25%,100% 75%,50% 100%,0% 75%,0% 25%); bottom: -60px; left: 220px; }
.slash { position: absolute; }
.slash-1 { width: 2px; height: 60vh; background: linear-gradient(to bottom,transparent,rgba(245,197,0,.12) 50%,transparent); top: 20%; left: 220px; transform: rotate(12deg); }
.slash-2 { width: 1px; height: 40vh; background: linear-gradient(to bottom,transparent,rgba(255,122,0,.08) 50%,transparent); top: 30%; left: 230px; transform: rotate(12deg); }
.noise { position: absolute; inset: 0; background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='1'/%3E%3C/svg%3E"); opacity: .02; }

/* ── Layout ── */
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
.main { flex: 1; height: 100vh; overflow-y: auto; padding: 48px 60px 48px; display: flex; flex-direction: column; gap: 26px; min-width: 0; }

/* Header */
.page-header { display: flex; align-items: flex-end; justify-content: space-between; flex-wrap: wrap; gap: 16px; flex-shrink: 0; }
.page-eyebrow { display: block; font-family: 'Barlow Condensed',sans-serif; font-size: .65rem; letter-spacing: .3em; text-transform: uppercase; color: #f5c500; margin-bottom: 6px; }
.page-title h1 { font-family: 'Barlow Condensed',sans-serif; font-size: 3.2rem; font-weight: 900; text-transform: uppercase; letter-spacing: .04em; color: #fff; line-height: 1; }
.id-tag { color: #f5c500; font-size: 2.2rem; }

.header-meta { display: flex; gap: 12px; }
.stat-pill { display: flex; flex-direction: column; align-items: center; background: #111; border: 1px solid rgba(245,197,0,.12); padding: 10px 22px; gap: 2px; }
.stat-pill.accent { border-color: rgba(245,197,0,.3); background: rgba(245,197,0,.06); }
.stat-val { font-family: 'Barlow Condensed',sans-serif; font-size: 1.7rem; font-weight: 900; color: #f5c500; line-height: 1; }
.stat-lbl { font-size: .58rem; letter-spacing: .2em; text-transform: uppercase; color: #555; }

/* States */
.state-box { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 14px; padding: 80px 20px; color: #555; font-size: .9rem; flex: 1; }
.bolt-spin { font-size: 2.4rem; animation: boltPulse 1.4s ease-in-out infinite; filter: drop-shadow(0 0 10px #f5c500); }
@keyframes boltPulse { 0%,100% { opacity: 1; } 50% { opacity: .3; } }
.error-state { color: #e05a45; }
.retry-btn { background: none; border: 1px solid rgba(220,50,30,.4); color: #e05a45; font-family: 'Barlow Condensed',sans-serif; font-size: .75rem; letter-spacing: .2em; text-transform: uppercase; padding: 8px 20px; cursor: pointer; transition: all .2s; }
.retry-btn:hover { background: rgba(220,50,30,.1); }

/* Tabla */
.table-wrapper { display: flex; flex-direction: column; border: 1px solid rgba(245,197,0,.12); background: #111; animation: fadeUp .45s cubic-bezier(.16,1,.3,1) both; flex: 1; min-height: 0; }
@keyframes fadeUp { from { opacity: 0; transform: translateY(12px); } to { opacity: 1; transform: translateY(0); } }
.table-scroll { overflow-y: auto; flex: 1; min-height: 0; }

table { width: 100%; border-collapse: collapse; font-size: .88rem; }
thead { background: rgba(245,197,0,.06); border-bottom: 1px solid rgba(245,197,0,.15); position: sticky; top: 0; z-index: 2; }
th { font-family: 'Barlow Condensed',sans-serif; font-size: .68rem; font-weight: 700; letter-spacing: .25em; text-transform: uppercase; color: #f5c500; padding: 16px 20px; text-align: left; white-space: nowrap; }

td { padding: 15px 20px; border-bottom: 1px solid rgba(255,255,255,.04); color: #ccc; vertical-align: middle; white-space: nowrap; }
.table-row { transition: background .15s; animation: rowIn .35s ease both; }
.table-row:hover { background: rgba(245,197,0,.04); }
.table-row:last-child td { border-bottom: none; }
@keyframes rowIn { from { opacity: 0; transform: translateX(-6px); } to { opacity: 1; transform: translateX(0); } }

/* Celdas */
.id-badge { display: inline-flex; align-items: center; justify-content: center; width: 28px; height: 28px; background: rgba(245,197,0,.07); border: 1px solid rgba(245,197,0,.2); font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .8rem; color: #f5c500; }
.tipo-label { font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .92rem; color: #f0f0f0; letter-spacing: .04em; text-transform: uppercase; }
.precio-val { font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .92rem; color: #f5c500; }

/* Badge estado de pago */
.estado-badge { display: inline-block; font-family: 'Barlow Condensed',sans-serif; font-size: .68rem; font-weight: 700; letter-spacing: .2em; text-transform: uppercase; padding: 4px 10px; border: 1px solid; }
.badge-pagado   { color: #4ade80; border-color: rgba(74,222,128,.4); background: rgba(74,222,128,.07); }
.badge-pendiente { color: #f59e0b; border-color: rgba(245,158,11,.4); background: rgba(245,158,11,.07); }
.badge-default  { color: #888; border-color: rgba(136,136,136,.3); background: rgba(136,136,136,.06); }

/* Badge disponibilidad por fecha */
.disp-badge { display: inline-block; font-family: 'Barlow Condensed',sans-serif; font-size: .68rem; font-weight: 700; letter-spacing: .15em; text-transform: uppercase; padding: 4px 10px; border: 1px solid; }
.disp-activo     { color: #4ade80; border-color: rgba(74,222,128,.35); background: rgba(74,222,128,.06); }
.disp-por-vencer { color: #f59e0b; border-color: rgba(245,158,11,.35); background: rgba(245,158,11,.06); }
.disp-proximo    { color: #60a5fa; border-color: rgba(96,165,250,.35); background: rgba(96,165,250,.06); }
.disp-vencido    { color: #555; border-color: rgba(85,85,85,.3); background: rgba(85,85,85,.05); }

.empty-row { text-align: center; color: #444; padding: 60px 20px !important; font-size: .88rem; letter-spacing: .04em; line-height: 1.7; }

.table-footer { padding: 13px 20px; border-top: 1px solid rgba(245,197,0,.08); font-size: .7rem; color: #444; letter-spacing: .06em; flex-shrink: 0; }
.table-footer strong { color: #888; }

@media (max-width: 860px) {
  .sidebar { display: none; }
  .main { padding: 28px 20px 40px; }
}
</style>