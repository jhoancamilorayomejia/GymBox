<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route  = useRoute()

const idcustomer    = Number(route.params.idcustomer)
const nombreCliente = history.state?.nombre || `Cliente #${idcustomer}`

const planes    = ref([])
const loading   = ref(true)
const errorMsg  = ref('')
const showModal = ref(false)
const saving    = ref(false)
const saveError = ref('')

// ── Filtros ─────────────────────────────────────────────
const filtro = ref({
  tipo_plan: '',
  estado: '',
  fecha_pago: '',
  fecha_inicio: '',
  fecha_fin: ''
})

const planesFiltrados = computed(() => {
  return planes.value.filter(p => {
    const tipoOK = filtro.value.tipo_plan
      ? p.typeplan?.toLowerCase().includes(filtro.value.tipo_plan.toLowerCase())
      : true
    const estadoOK = filtro.value.estado
      ? p.state?.toLowerCase() === filtro.value.estado.toLowerCase()
      : true
    const fechaPagoOK = filtro.value.fecha_pago
      ? p.datepay?.startsWith(filtro.value.fecha_pago)
      : true
    const fechaInicioOK = filtro.value.fecha_inicio
      ? p.datestart?.startsWith(filtro.value.fecha_inicio)
      : true
    const fechaFinOK = filtro.value.fecha_fin
      ? p.datefinish?.startsWith(filtro.value.fecha_fin)
      : true
    return tipoOK && estadoOK && fechaPagoOK && fechaInicioOK && fechaFinOK
  })
})

const eliminarPlanesFiltrados = async () => {
  if (!planesFiltrados.value.length) {
    alert('No hay planes para eliminar')
    return
  }
  if (!confirm(`Vas a eliminar ${planesFiltrados.value.length} planes filtrados. ¿Continuar?`)) return
  try {
    const token = getToken()
    if (!token) return
    for (const plan of planesFiltrados.value) {
      await fetch(`/api/delete-plans/${plan.idplan}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` }
      })
    }
    await obtenerPlanes()
  } catch (err) {
    alert(err.message)
  }
}

// ── Helper: parsea "YYYY-MM-DD" o "YYYY-MM-DDTHH:..." como fecha LOCAL (sin desfase UTC)
const parseLocalDate = (str) => {
  if (!str) return null
  const [y, m, d] = str.split('T')[0].split('-').map(Number)
  const fecha = new Date(y, m - 1, d)
  fecha.setHours(0, 0, 0, 0)
  return fecha
}

// ── Modal calendario ──────────────────────────────────────────────
const showCalModal   = ref(false)
const planCalendario = ref(null)
const calViewDate    = ref(new Date())

const MESES = ['Enero','Febrero','Marzo','Abril','Mayo','Junio','Julio','Agosto','Septiembre','Octubre','Noviembre','Diciembre']
const DIAS  = ['Dom','Lun','Mar','Mié','Jue','Vie','Sáb']

const abrirCalendario = (plan) => {
  planCalendario.value = plan
  const ini = parseLocalDate(plan.datestart)
  calViewDate.value = new Date(ini.getFullYear(), ini.getMonth(), 1)
  showCalModal.value = true
}

const calMesLabel  = computed(() => `${MESES[calViewDate.value.getMonth()]} ${calViewDate.value.getFullYear()}`)
const calPrimerDia = computed(() => new Date(calViewDate.value.getFullYear(), calViewDate.value.getMonth(), 1).getDay())
const calDiasEnMes = computed(() => new Date(calViewDate.value.getFullYear(), calViewDate.value.getMonth() + 1, 0).getDate())

const cambiarCalMes = (delta) => {
  const d = new Date(calViewDate.value)
  d.setMonth(d.getMonth() + delta)
  calViewDate.value = d
}

// ── Clasificar día del calendario (CORREGIDO: usa parseLocalDate) ──
const clasificarDiaCal = (dia) => {
  if (!planCalendario.value) return { disponible: false, vencido: false, esHoy: false }

  const y     = calViewDate.value.getFullYear()
  const m     = calViewDate.value.getMonth()
  const fecha = new Date(y, m, dia); fecha.setHours(0, 0, 0, 0)
  const hoy   = new Date(); hoy.setHours(0, 0, 0, 0)
  const esHoy = fecha.getTime() === hoy.getTime()

  const ini = parseLocalDate(planCalendario.value.datestart)
  const fin = parseLocalDate(planCalendario.value.datefinish)

  if (fecha < ini || fecha > fin) return { disponible: false, vencido: false, esHoy }

  const esPasado = fin < hoy
  return {
    disponible: !esPasado,
    vencido:    esPasado,
    esHoy
  }
}

// ── Info resumida del plan en el modal (CORREGIDO: usa parseLocalDate) ──
const planCalInfo = computed(() => {
  if (!planCalendario.value) return null
  const p   = planCalendario.value
  const hoy = new Date(); hoy.setHours(0, 0, 0, 0)
  const ini = parseLocalDate(p.datestart)
  const fin = parseLocalDate(p.datefinish)

  let estado, colorClass
  if (fin < hoy)      { estado = 'Vencido';       colorClass = 'info-vencido' }
  else if (ini > hoy) { estado = 'Próximo inicio'; colorClass = 'info-proximo' }
  else                { estado = 'Activo';         colorClass = 'info-activo'  }

  const diasRest = fin >= hoy ? Math.ceil((fin - hoy) / 86400000) : 0

  return { estado, colorClass, diasRest, ini, fin }
})

// ── Formularios ───────────────────────────────────────────────────
const getFechaHoy = () => {
  const h = new Date()
  return `${h.getFullYear()}-${String(h.getMonth()+1).padStart(2,'0')}-${String(h.getDate()).padStart(2,'0')}`
}

const formPlan = ref({ tipo_plan: '', precio: '', fecha_pago: getFechaHoy(), fecha_inicio: '', fecha_fin: '', state: '' })
const resetForm = () => {
  formPlan.value = { tipo_plan: '', precio: '', fecha_pago: getFechaHoy(), fecha_inicio: '', fecha_fin: '', state: '' }
  saveError.value = ''
}

const planEditandoId = ref(null)
const showEditModal  = ref(false)

const abrirEditar = (plan) => {
  planEditandoId.value = plan.idplan
  formPlan.value = {
    tipo_plan:    plan.typeplan,
    precio:       plan.price,
    fecha_pago:   plan.datepay?.split('T')[0]    || plan.datepay,
    fecha_inicio: plan.datestart?.split('T')[0]  || plan.datestart,
    fecha_fin:    plan.datefinish?.split('T')[0] || plan.datefinish,
    state:        plan.state
  }
  saveError.value = ''
  showEditModal.value = true
}

// ── Auth ──────────────────────────────────────────────────────────
const tokenExpirado = (token) => {
  try { return JSON.parse(atob(token.split('.')[1])).exp * 1000 < Date.now() }
  catch { return true }
}

const getToken = () => {
  const t = localStorage.getItem('token')
  if (!t || tokenExpirado(t)) { localStorage.removeItem('token'); router.push('/login'); return null }
  return t
}

const cerrarSesion = () => {
  ['token','rol','username'].forEach(k => localStorage.removeItem(k))
  router.push('/login')
}

// ── Fetch ─────────────────────────────────────────────────────────
const obtenerPlanes = async () => {
  loading.value = true; errorMsg.value = ''
  try {
    const token = getToken(); if (!token) return
    const res = await fetch(`/api/plans?idcustomer=${idcustomer}`, { headers: { Authorization: `Bearer ${token}` } })
    if (res.status === 401) { localStorage.removeItem('token'); router.push('/login'); return }
    if (!res.ok) throw new Error('Error al cargar planes')
    const data = await res.json()
    const arr  = Array.isArray(data) ? data : (data?.data ?? [])
    planes.value = arr.filter(p => Number(p.idcustomer) === Number(idcustomer))
  } catch (err) { errorMsg.value = err.message }
  finally { loading.value = false }
}

const guardarPlan = async () => {
  saveError.value = ''
  const f = formPlan.value
  if (!f.tipo_plan || !f.precio || !f.fecha_pago || !f.fecha_inicio || !f.fecha_fin || !f.state) {
    saveError.value = 'Todos los campos son obligatorios'; return
  }
  saving.value = true
  try {
    const token = getToken(); if (!token) return
    const res = await fetch('/api/plans', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify({ idcustomer, typeplan: f.tipo_plan, price: Number(f.precio), datepay: f.fecha_pago, datestart: f.fecha_inicio, datefinish: f.fecha_fin, state: f.state })
    })
    if (res.status === 401) { localStorage.removeItem('token'); router.push('/login'); return }
    if (!res.ok) { const d = await res.json(); throw new Error(d.error || 'Error al guardar') }
    showModal.value = false; resetForm(); await obtenerPlanes()
  } catch (err) { saveError.value = err.message }
  finally { saving.value = false }
}

const actualizarPlan = async () => {
  saveError.value = ''
  const f = formPlan.value
  if (!f.tipo_plan || !f.precio || !f.fecha_pago || !f.fecha_inicio || !f.fecha_fin || !f.state) {
    saveError.value = 'Todos los campos son obligatorios'; return
  }
  saving.value = true
  try {
    const token = getToken(); if (!token) return
    const res = await fetch(`/api/update-plans/${planEditandoId.value}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify({ idplan: planEditandoId.value, typeplan: f.tipo_plan, price: Number(f.precio), datepay: f.fecha_pago, datestart: f.fecha_inicio, datefinish: f.fecha_fin, state: f.state })
    })
    if (res.status === 401) { localStorage.removeItem('token'); router.push('/login'); return }
    if (!res.ok) { const d = await res.json(); throw new Error(d.error || 'Error al actualizar') }
    showEditModal.value = false; await obtenerPlanes()
  } catch (err) { saveError.value = err.message }
  finally { saving.value = false }
}

const eliminarPlan = async (id) => {
  if (!confirm('¿Seguro que deseas eliminar este plan?')) return
  try {
    const token = getToken(); if (!token) return
    const res = await fetch(`/api/delete-plans/${id}`, { method: 'DELETE', headers: { Authorization: `Bearer ${token}` } })
    if (res.status === 401) { localStorage.removeItem('token'); router.push('/login'); return }
    if (!res.ok) { const d = await res.json(); throw new Error(d.error || 'Error al eliminar') }
    await obtenerPlanes()
  } catch (err) { alert(err.message) }
}

// ── Helpers ───────────────────────────────────────────────────────
const fmtFecha = (str) => {
  if (!str) return '—'
  const [year, month, day] = str.split('T')[0].split('-')
  return `${day}/${month}/${year}`
}

const estadoClass = (s) => {
  if (!s) return 'badge-default'
  if (s.toLowerCase() === 'pagado')    return 'badge-pagado'
  if (s.toLowerCase() === 'pendiente') return 'badge-pendiente'
  return 'badge-default'
}

const username = localStorage.getItem('username') || 'Admin'

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
            Planes
          </a>
          <a class="nav-item" @click="router.push('/dashboard')">
            <svg viewBox="0 0 20 20" fill="none"><path d="M12 16l-5-6 5-6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
            Regresar
          </a>
        </nav>
        <div class="sidebar-footer">
          <div class="user-chip">
            <div class="user-avatar">{{ username.charAt(0).toUpperCase() }}</div>
            <div class="user-info">
              <span class="user-name">{{ username }}</span>
              <span class="user-role">Administrador</span>
            </div>
          </div>
          <button class="logout-btn" @click="cerrarSesion" title="Cerrar sesión">
            <svg viewBox="0 0 20 20" fill="none"><path d="M13 3h4v14h-4M9 14l4-4-4-4M13 10H5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
          </button>
        </div>
      </aside>

      <!-- ── Main ── -->
      <main class="main">

        <div class="breadcrumb" @click="router.push('/dashboard')">
          <svg viewBox="0 0 20 20" fill="none"><path d="M12 16l-5-6 5-6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
          Volver a Clientes
        </div>

        <header class="page-header">
          <div class="page-title">
            <span class="page-eyebrow">⚡ Planes de membresía</span>
            <h1>{{ nombreCliente }}</h1>
          </div>
          <div class="header-meta">
            <div class="stat-pill">
              <span class="stat-val">{{ planes.length }}</span>
              <span class="stat-lbl">Planes</span>
            </div>
            <div class="stat-pill accent">
              <span class="stat-val">{{ planes.filter(p => p.state?.toLowerCase() === 'pagado').length }}</span>
              <span class="stat-lbl">Pagados</span>
            </div>
          </div>
        </header>

        <div v-if="loading" class="state-box">
          <div class="bolt-spin">⚡</div><p>Cargando planes...</p>
        </div>

        <div v-else-if="errorMsg" class="state-box error-state">
          <span>⚠</span><p>{{ errorMsg }}</p>
          <button class="retry-btn" @click="obtenerPlanes">Reintentar</button>
        </div>

        <div v-else class="table-wrapper">

          <div class="filtros-box">
            <div class="filtro-item">
              <label>Tipo de plan</label>
              <input v-model="filtro.tipo_plan" placeholder="Ej: mes" />
            </div>
            <div class="filtro-item">
              <label>Estado</label>
              <select v-model="filtro.estado">
                <option value="">Todos</option>
                <option value="pagado">Pagado</option>
                <option value="pendiente">Pendiente</option>
              </select>
            </div>
            <div class="filtro-item">
              <label>Fecha de pago</label>
              <input type="date" v-model="filtro.fecha_pago" />
            </div>
            <div class="filtro-item">
              <label>Fecha inicio</label>
              <input type="date" v-model="filtro.fecha_inicio" />
            </div>
            <div class="filtro-item">
              <label>Fecha fin</label>
              <input type="date" v-model="filtro.fecha_fin" />
            </div>
            <div class="filtro-item btn-box">
              <button @click="filtro = { tipo_plan:'', estado:'', fecha_pago:'', fecha_inicio:'', fecha_fin:'' }">
                Limpiar
              </button>
              <button @click="eliminarPlanesFiltrados">
                Borrar visibles
              </button>
            </div>
          </div>

          <div class="table-scroll">
            <table>
              <thead>
                <tr>
                  <th>Tipo de Plan</th>
                  <th>Precio</th>
                  <th>Fecha Pago</th>
                  <th>Inicio</th>
                  <th>Fin</th>
                  <th>Estado</th>
                  <th>Calendario</th>
                  <th class="th-add">
                    <button class="btn-agregar" @click="showModal = true; resetForm()">
                      <svg viewBox="0 0 20 20" fill="none">
                        <circle cx="10" cy="10" r="7.5" stroke="currentColor" stroke-width="1.4"/>
                        <path d="M10 7v6M7 10h6" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/>
                      </svg>
                      Agregar Plan
                    </button>
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(p, i) in planesFiltrados" :key="p.idplan"
                    class="table-row" :style="{ animationDelay: i * 40 + 'ms' }">
                  <td><span class="tipo-label">{{ p.typeplan }}</span></td>
                  <td><span class="precio-val">${{ Number(p.price).toLocaleString('es-CO') }}</span></td>
                  <td>{{ fmtFecha(p.datepay) }}</td>
                  <td>{{ fmtFecha(p.datestart) }}</td>
                  <td>{{ fmtFecha(p.datefinish) }}</td>
                  <td><span class="estado-badge" :class="estadoClass(p.state)">{{ p.state }}</span></td>
                  <td class="td-cal">
                    <button class="btn-cal" @click="abrirCalendario(p)">
                      <svg viewBox="0 0 20 20" fill="none">
                        <rect x="3" y="4" width="14" height="13" rx="1.5" stroke="currentColor" stroke-width="1.4"/>
                        <path d="M7 2v4M13 2v4M3 9h14" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/>
                      </svg>
                      Ver calendario
                    </button>
                  </td>
                  <td class="td-center">
                    <button class="action-btn edit-btn" @click="abrirEditar(p)" title="Editar">
                      <svg viewBox="0 0 20 20" fill="none"><path d="M14.5 3.5l2 2L6 16l-3 1 1-3 11.5-10.5z" stroke="currentColor" stroke-width="1.3" stroke-linejoin="round"/></svg>
                    </button>
                    <button class="action-btn delete-btn" @click="eliminarPlan(p.idplan)" title="Eliminar">
                      <svg viewBox="0 0 20 20" fill="none"><path d="M5 7h10l-1 9H6L5 7zM3 7h14M8 7V5h4v2" stroke="currentColor" stroke-width="1.3" stroke-linecap="round" stroke-linejoin="round"/></svg>
                    </button>
                  </td>
                </tr>
                <tr v-if="!planes.length">
                  <td colspan="8" class="empty-row">
                    ⚡ Este cliente no tiene planes aún.
                    <span class="empty-cta" @click="showModal = true; resetForm()">Crear el primero →</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="table-footer">
            Mostrando <strong>{{ planes.length }}</strong> plan{{ planes.length !== 1 ? 'es' : '' }} para <strong>{{ nombreCliente }}</strong>
          </div>
        </div>

      </main>
    </div>

    <!-- ══════════════════════════════════════════════════════════ -->
    <!-- MODAL CALENDARIO VISUAL                                    -->
    <!-- ══════════════════════════════════════════════════════════ -->
    <transition name="modal-fade">
      <div v-if="showCalModal" class="modal-overlay" @click.self="showCalModal = false">
        <div class="modal modal-cal">

          <div class="modal-header">
            <div>
              <span class="modal-eyebrow">⚡ Disponibilidad</span>
              <h2>{{ planCalendario?.typeplan?.toUpperCase() }}</h2>
              <p class="modal-sub">{{ nombreCliente }}</p>
            </div>
            <button class="modal-close" @click="showCalModal = false">✕</button>
          </div>

          <!-- Info rápida del plan -->
          <div class="cal-info-row" v-if="planCalInfo">
            <div class="cal-info-item">
              <span class="ci-key">Estado</span>
              <span class="ci-val" :class="planCalInfo.colorClass">{{ planCalInfo.estado }}</span>
            </div>
            <div class="cal-info-item">
              <span class="ci-key">Inicio</span>
              <span class="ci-val">{{ fmtFecha(planCalendario.datestart) }}</span>
            </div>
            <div class="cal-info-divider">→</div>
            <div class="cal-info-item">
              <span class="ci-key">Fin</span>
              <span class="ci-val">{{ fmtFecha(planCalendario.datefinish) }}</span>
            </div>
            <div class="cal-info-item" v-if="planCalInfo.diasRest > 0">
              <span class="ci-key">Días rest.</span>
              <span class="ci-val ci-dias">{{ planCalInfo.diasRest }}</span>
            </div>
            <div class="cal-info-item">
              <span class="ci-key">Precio</span>
              <span class="ci-val">${{ Number(planCalendario?.price || 0).toLocaleString('es-CO') }}</span>
            </div>
          </div>

          <!-- Calendario -->
          <div class="modal-body cal-body">

            <div class="cal-nav">
              <button class="cal-arrow" @click="cambiarCalMes(-1)">&#8249;</button>
              <span class="cal-mes">{{ calMesLabel }}</span>
              <button class="cal-arrow" @click="cambiarCalMes(1)">&#8250;</button>
            </div>

            <div class="cal-header-row">
              <span v-for="d in DIAS" :key="d">{{ d }}</span>
            </div>

            <div class="cal-grid">
              <div v-for="n in calPrimerDia" :key="'e'+n" class="cal-empty"></div>
              <div
                v-for="dia in calDiasEnMes"
                :key="dia"
                class="cal-dia"
                :class="{
                  'dia-activo':  clasificarDiaCal(dia).disponible && !clasificarDiaCal(dia).esHoy,
                  'dia-vencido': clasificarDiaCal(dia).vencido   && !clasificarDiaCal(dia).esHoy,
                  'dia-hoy':     clasificarDiaCal(dia).esHoy,
                  'dia-ninguno': !clasificarDiaCal(dia).disponible && !clasificarDiaCal(dia).vencido && !clasificarDiaCal(dia).esHoy
                }"
              >
                {{ dia }}
              </div>
            </div>

            <!-- Leyenda -->
            <div class="cal-leyenda">
              <div class="ley-item"><span class="ley-dot dot-activo"></span> Día disponible</div>
              <div class="ley-item"><span class="ley-dot dot-vencido"></span> Día vencido</div>
              <div class="ley-item"><span class="ley-dot dot-hoy"></span> Hoy</div>
              <div class="ley-item"><span class="ley-dot dot-ninguno"></span> Fuera del plan</div>
            </div>

          </div>

          <div class="modal-footer">
            <button class="btn-cancel" @click="showCalModal = false">Cerrar</button>
          </div>

        </div>
      </div>
    </transition>

    <!-- ══════════════════════════════════════════════════════════ -->
    <!-- MODAL CREAR PLAN                                           -->
    <!-- ══════════════════════════════════════════════════════════ -->
    <transition name="modal-fade">
      <div v-if="showModal" class="modal-overlay" @click.self="showModal = false; resetForm()">
        <div class="modal">
          <div class="modal-header">
            <div>
              <span class="modal-eyebrow">⚡ Nuevo plan</span>
              <h2>Asignar membresía</h2>
              <p class="modal-sub">{{ nombreCliente }}</p>
            </div>
            <button class="modal-close" @click="showModal = false; resetForm()">✕</button>
          </div>
          <div class="modal-body">
            <div class="modal-grid">
              <div class="mfield">
                <label>Tipo de plan <span class="req">*</span></label>
                <select v-model="formPlan.tipo_plan">
                  <option disabled value="">Seleccionar...</option>
                  <option value="dia">Día</option>
                  <option value="semana">Semana</option>
                  <option value="mes">Mes</option>
                  <option value="anio">Año</option>
                </select>
              </div>
              <div class="mfield">
                <label>Precio (COP) <span class="req">*</span></label>
                <input v-model="formPlan.precio" type="number" placeholder="Ej: 120000" min="0" />
              </div>
              <div class="mfield">
                <label>Fecha de pago <span class="req">*</span></label>
                <input v-model="formPlan.fecha_pago" type="date" disabled />
              </div>
              <div class="mfield">
                <label>Estado <span class="req">*</span></label>
                <select v-model="formPlan.state">
                  <option disabled value="">Seleccionar...</option>
                  <option value="pagado">Pagado</option>
                  <option value="pendiente">Pendiente</option>
                </select>
              </div>
              <div class="mfield">
                <label>Fecha inicio <span class="req">*</span></label>
                <input v-model="formPlan.fecha_inicio" type="date" />
              </div>
              <div class="mfield">
                <label>Fecha fin <span class="req">*</span></label>
                <input v-model="formPlan.fecha_fin" type="date" />
              </div>
            </div>
            <p v-if="saveError" class="save-error">⚠ {{ saveError }}</p>
          </div>
          <div class="modal-footer">
            <button class="btn-cancel" @click="showModal = false; resetForm()">Cancelar</button>
            <button class="btn-save" @click="guardarPlan" :disabled="saving">
              <svg v-if="saving" class="spinner" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3" stroke-dasharray="31.4" stroke-dashoffset="10"/>
              </svg>
              {{ saving ? 'Guardando...' : 'Guardar Plan' }}
            </button>
          </div>
        </div>
      </div>
    </transition>

    <!-- ══════════════════════════════════════════════════════════ -->
    <!-- MODAL EDITAR PLAN                                          -->
    <!-- ══════════════════════════════════════════════════════════ -->
    <transition name="modal-fade">
      <div v-if="showEditModal" class="modal-overlay" @click.self="showEditModal = false">
        <div class="modal">
          <div class="modal-header">
            <div>
              <span class="modal-eyebrow">✏️ Editar plan</span>
              <h2>Actualizar membresía</h2>
              <p class="modal-sub">ID Plan: {{ planEditandoId }}</p>
            </div>
            <button class="modal-close" @click="showEditModal = false">✕</button>
          </div>
          <div class="modal-body">
            <div class="modal-grid">
              <div class="mfield">
                <label>Tipo de plan *</label>
                <select v-model="formPlan.tipo_plan">
                  <option value="dia">Día</option>
                  <option value="semana">Semana</option>
                  <option value="mes">Mes</option>
                  <option value="anio">Año</option>
                </select>
              </div>
              <div class="mfield">
                <label>Precio *</label>
                <input v-model="formPlan.precio" type="number" />
              </div>
              <div class="mfield">
                <label>Fecha de pago</label>
                <input v-model="formPlan.fecha_pago" type="date" />
              </div>
              <div class="mfield">
                <label>Estado *</label>
                <select v-model="formPlan.state">
                  <option value="pagado">Pagado</option>
                  <option value="pendiente">Pendiente</option>
                </select>
              </div>
              <div class="mfield">
                <label>Fecha inicio *</label>
                <input v-model="formPlan.fecha_inicio" type="date" />
              </div>
              <div class="mfield">
                <label>Fecha fin *</label>
                <input v-model="formPlan.fecha_fin" type="date" />
              </div>
            </div>
            <p v-if="saveError" class="save-error">⚠ {{ saveError }}</p>
          </div>
          <div class="modal-footer">
            <button class="btn-cancel" @click="showEditModal = false">Cancelar</button>
            <button class="btn-save" @click="actualizarPlan" :disabled="saving">
              {{ saving ? 'Actualizando...' : 'Actualizar Plan' }}
            </button>
          </div>
        </div>
      </div>
    </transition>

  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Barlow+Condensed:wght@400;700;900&family=Barlow:wght@300;400;500&display=swap');
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

.screen { width: 100vw; height: 100vh; overflow: hidden; background: #0a0a0a; font-family: 'Barlow', sans-serif; position: relative; }

/* Fondo */
.bg-layer { position: fixed; inset: 0; pointer-events: none; z-index: 0; }
.hex { position: absolute; }
.hex-1 { width: 420px; height: 420px; background: conic-gradient(from 30deg,rgba(245,197,0,.08) 0deg,rgba(255,122,0,.06) 80deg,transparent 80deg); clip-path: polygon(50% 0%,100% 25%,100% 75%,50% 100%,0% 75%,0% 25%); top: -100px; right: -80px; }
.hex-2 { width: 260px; height: 260px; background: conic-gradient(from 200deg,rgba(245,197,0,.05) 0deg,transparent 80deg); clip-path: polygon(50% 0%,100% 25%,100% 75%,50% 100%,0% 75%,0% 25%); bottom: -60px; left: 220px; }
.slash { position: absolute; }
.slash-1 { width: 2px; height: 60vh; background: linear-gradient(to bottom,transparent,rgba(245,197,0,.12) 50%,transparent); top: 20%; left: 220px; transform: rotate(12deg); }
.slash-2 { width: 1px; height: 40vh; background: linear-gradient(to bottom,transparent,rgba(255,122,0,.08) 50%,transparent); top: 30%; left: 230px; transform: rotate(12deg); }
.noise { position: absolute; inset: 0; background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='1'/%3E%3C/svg%3E"); opacity: .02; }

.layout { position: relative; z-index: 1; display: flex; width: 100%; height: 100vh; }

/* Sidebar */
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
.user-role { font-size: .6rem; color: #555; letter-spacing: .08em; }
.logout-btn { background: none; border: 1px solid #1e1e1e; color: #444; cursor: pointer; width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; transition: all .2s; }
.logout-btn svg { width: 15px; height: 15px; }
.logout-btn:hover { border-color: rgba(245,197,0,.4); color: #f5c500; }

/* Main */
.main { flex: 1; height: 100vh; overflow-y: auto; padding: 48px 60px; display: flex; flex-direction: column; gap: 26px; min-width: 0; }

.breadcrumb { display: inline-flex; align-items: center; gap: 6px; font-family: 'Barlow Condensed',sans-serif; font-size: .72rem; letter-spacing: .18em; text-transform: uppercase; color: #555; cursor: pointer; transition: color .2s; width: fit-content; flex-shrink: 0; }
.breadcrumb svg { width: 14px; height: 14px; }
.breadcrumb:hover { color: #f5c500; }

.page-header { display: flex; align-items: flex-end; justify-content: space-between; flex-wrap: wrap; gap: 16px; flex-shrink: 0; }
.page-eyebrow { display: block; font-family: 'Barlow Condensed',sans-serif; font-size: .65rem; letter-spacing: .3em; text-transform: uppercase; color: #f5c500; margin-bottom: 6px; }
.page-title h1 { font-family: 'Barlow Condensed',sans-serif; font-size: 3.2rem; font-weight: 900; text-transform: uppercase; letter-spacing: .04em; color: #fff; line-height: 1; }
.header-meta { display: flex; gap: 12px; }
.stat-pill { display: flex; flex-direction: column; align-items: center; background: #111; border: 1px solid rgba(245,197,0,.12); padding: 10px 22px; gap: 2px; }
.stat-pill.accent { border-color: rgba(245,197,0,.3); background: rgba(245,197,0,.06); }
.stat-val { font-family: 'Barlow Condensed',sans-serif; font-size: 1.7rem; font-weight: 900; color: #f5c500; line-height: 1; }
.stat-lbl { font-size: .58rem; letter-spacing: .2em; text-transform: uppercase; color: #555; }

.state-box { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 14px; padding: 80px 20px; color: #555; font-size: .9rem; }
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
.th-add { text-align: right; padding-right: 16px; }
.btn-agregar { display: inline-flex; align-items: center; gap: 7px; background: #f5c500; border: none; color: #0a0a0a; font-family: 'Barlow Condensed',sans-serif; font-size: .72rem; font-weight: 900; letter-spacing: .2em; text-transform: uppercase; padding: 8px 18px; cursor: pointer; transition: all .2s; white-space: nowrap; }
.btn-agregar svg { width: 14px; height: 14px; flex-shrink: 0; }
.btn-agregar:hover { background: #ffd700; box-shadow: 0 0 16px rgba(245,197,0,.35); }
td { padding: 14px 20px; border-bottom: 1px solid rgba(255,255,255,.04); color: #ccc; vertical-align: middle; white-space: nowrap; }
.table-row { transition: background .15s; animation: rowIn .35s ease both; }
.table-row:hover { background: rgba(245,197,0,.04); }
.table-row:last-child td { border-bottom: none; }
@keyframes rowIn { from { opacity: 0; transform: translateX(-6px); } to { opacity: 1; transform: translateX(0); } }
.tipo-label { font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .92rem; color: #f0f0f0; letter-spacing: .04em; }
.precio-val { font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .95rem; color: #f5c500; }
.estado-badge { display: inline-block; font-family: 'Barlow Condensed',sans-serif; font-size: .7rem; font-weight: 700; letter-spacing: .2em; text-transform: uppercase; padding: 4px 12px; border: 1px solid; }
.badge-pagado    { color: #4ade80; border-color: rgba(74,222,128,.4); background: rgba(74,222,128,.07); }
.badge-pendiente { color: #f59e0b; border-color: rgba(245,158,11,.4); background: rgba(245,158,11,.07); }
.badge-default   { color: #888; border-color: rgba(136,136,136,.3); background: rgba(136,136,136,.06); }

/* Botón calendario */
.td-cal { text-align: center; }
.btn-cal { display: inline-flex; align-items: center; gap: 6px; background: transparent; border: 1px solid rgba(245,197,0,.25); color: #f5c500; font-family: 'Barlow Condensed',sans-serif; font-size: .7rem; font-weight: 700; letter-spacing: .16em; text-transform: uppercase; padding: 6px 14px; cursor: pointer; transition: all .2s; white-space: nowrap; }
.btn-cal svg { width: 13px; height: 13px; flex-shrink: 0; }
.btn-cal:hover { background: rgba(245,197,0,.1); border-color: #f5c500; }

.td-center { text-align: center; }
.action-btn { background: none; border: 1px solid transparent; width: 30px; height: 30px; display: inline-flex; align-items: center; justify-content: center; cursor: pointer; transition: all .2s; }
.action-btn svg { width: 13px; height: 13px; }
.edit-btn { color: #555; }
.edit-btn:hover { color: #f5c500; border-color: rgba(245,197,0,.4); background: rgba(245,197,0,.07); }
.delete-btn { color: #555; }
.delete-btn:hover { color: #e05a45; border-color: rgba(220,50,30,.4); background: rgba(220,50,30,.07); }
.empty-row { text-align: center; color: #444; padding: 56px 20px !important; font-size: .88rem; }
.empty-cta { color: #f5c500; cursor: pointer; margin-left: 8px; text-decoration: underline; }
.table-footer { padding: 13px 20px; border-top: 1px solid rgba(245,197,0,.08); font-size: .7rem; color: #444; letter-spacing: .06em; flex-shrink: 0; }
.table-footer strong { color: #888; }

/* ══ MODALES ══ */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,.78); display: flex; align-items: center; justify-content: center; z-index: 1000; padding: 20px; }
.modal { background: #111; border: 1px solid rgba(245,197,0,.2); box-shadow: 0 40px 100px rgba(0,0,0,.9),0 0 60px rgba(245,197,0,.05); width: 100%; max-width: 560px; display: flex; flex-direction: column; animation: modalIn .3s cubic-bezier(.16,1,.3,1) both; }
.modal-cal { max-width: 480px; }
@keyframes modalIn { from { opacity: 0; transform: scale(.96) translateY(12px); } to { opacity: 1; transform: scale(1) translateY(0); } }

.modal-header { display: flex; align-items: flex-start; justify-content: space-between; padding: 22px 26px 18px; border-bottom: 1px solid rgba(245,197,0,.1); }
.modal-eyebrow { display: block; font-family: 'Barlow Condensed',sans-serif; font-size: .6rem; letter-spacing: .3em; text-transform: uppercase; color: #f5c500; margin-bottom: 5px; }
.modal-header h2 { font-family: 'Barlow Condensed',sans-serif; font-size: 1.5rem; font-weight: 900; text-transform: uppercase; color: #fff; line-height: 1; }
.modal-sub { font-size: .75rem; color: #666; margin-top: 4px; }
.modal-close { background: none; border: none; color: #444; font-size: 1rem; cursor: pointer; padding: 2px; transition: color .2s; line-height: 1; }
.modal-close:hover { color: #f5c500; }

/* Info rápida del plan */
.cal-info-row { display: flex; align-items: center; gap: 16px; flex-wrap: wrap; padding: 14px 26px; border-bottom: 1px solid rgba(245,197,0,.08); background: rgba(245,197,0,.02); }
.cal-info-item { display: flex; flex-direction: column; gap: 2px; }
.ci-key { font-family: 'Barlow Condensed',sans-serif; font-size: .58rem; letter-spacing: .2em; text-transform: uppercase; color: #555; font-weight: 700; }
.ci-val { font-size: .85rem; color: #e0e0e0; font-weight: 500; }
.ci-dias { color: #f5c500; font-family: 'Barlow Condensed',sans-serif; font-weight: 900; font-size: 1rem; }
.cal-info-divider { color: #f5c500; font-weight: 700; font-size: .9rem; padding-top: 12px; }
.info-activo  { color: #4ade80 !important; }
.info-vencido { color: #555 !important; }
.info-proximo { color: #60a5fa !important; }

/* Calendario dentro del modal */
.cal-body { padding: 20px 26px; display: flex; flex-direction: column; gap: 14px; }

.cal-nav { display: flex; align-items: center; justify-content: space-between; }
.cal-mes { font-family: 'Barlow Condensed',sans-serif; font-size: .95rem; font-weight: 700; letter-spacing: .12em; text-transform: uppercase; color: #f0f0f0; }
.cal-arrow { background: transparent; border: 1px solid rgba(245,197,0,.2); color: #f5c500; width: 28px; height: 28px; display: flex; align-items: center; justify-content: center; cursor: pointer; font-size: 1.2rem; line-height: 1; transition: all .2s; }
.cal-arrow:hover { background: rgba(245,197,0,.1); border-color: #f5c500; }

.cal-header-row { display: grid; grid-template-columns: repeat(7, 1fr); text-align: center; gap: 3px; }
.cal-header-row span { font-family: 'Barlow Condensed',sans-serif; font-size: .6rem; letter-spacing: .1em; text-transform: uppercase; color: #f5c500; padding: 4px 0; }

.cal-grid { display: grid; grid-template-columns: repeat(7, 1fr); gap: 3px; }
.cal-empty { aspect-ratio: 1; }
.cal-dia { aspect-ratio: 1; display: flex; align-items: center; justify-content: center; font-family: 'Barlow Condensed',sans-serif; font-size: .8rem; font-weight: 700; border-radius: 2px; transition: all .15s; }

.dia-ninguno { color: #bdb9b9; background: rgba(255,255,255,.02); border: 1px solid rgba(255,255,255,.03); }
.dia-activo  { color: #d1fae5; background: rgba(74,222,128,.18); border: 1px solid rgba(74,222,128,.35); }
.dia-activo:hover { background: rgba(74,222,128,.26); }
.dia-vencido { color: #555; background: rgba(248,0,0,.12); border: 1px solid rgba(100,100,100,.22); }
.dia-hoy     { background: #f5c500 !important; color: #0a0a0a !important; border: none !important; font-weight: 900; box-shadow: 0 0 12px rgba(245,197,0,.4); }

/* Leyenda */
.cal-leyenda { display: flex; gap: 16px; flex-wrap: wrap; padding-top: 4px; border-top: 1px solid rgba(245,197,0,.07); }
.ley-item { display: flex; align-items: center; gap: 6px; font-size: .7rem; color: #666; }
.ley-dot { width: 11px; height: 11px; border-radius: 2px; flex-shrink: 0; }
.dot-activo  { background: rgba(74,222,128,.25); border: 1px solid rgba(74,222,128,.6); }
.dot-vencido { background: rgba(248,0,0,.12); border: 1px solid rgba(100,100,100,.5); }
.dot-hoy     { background: #f5c500; border: 1px solid #f5c500; }
.dot-ninguno { background: rgba(255,255,255,.03); border: 1px solid rgba(255,255,255,.08); }

/* Formulario modal */
.modal-body { padding: 12px 65px; }
.modal-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
.mfield { display: flex; flex-direction: column; gap: 7px; }
.mfield label { font-family: 'Barlow Condensed',sans-serif; font-size: .62rem; letter-spacing: .25em; text-transform: uppercase; color: #666; font-weight: 700; }
.req { color: #f5c500; }
.mfield input, .mfield select { background: #0d0d0d; border: 1px solid rgba(245,197,0,.15); color: #f0f0f0; font-family: 'Barlow',sans-serif; font-size: .85rem; padding: 10px 12px; outline: none; transition: border-color .2s, background .2s; appearance: none; }
.mfield input::placeholder { color: #333; }
.mfield input:focus, .mfield select:focus { border-color: rgba(245,197,0,.5); background: #0f0e09; }
.mfield select { cursor: pointer; background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='10' height='6' viewBox='0 0 10 6'%3E%3Cpath d='M1 1l4 4 4-4' stroke='%23f5c500' stroke-width='1.2' fill='none'/%3E%3C/svg%3E"); background-repeat: no-repeat; background-position: right 12px center; background-color: #0d0d0d; padding-right: 32px; }
.mfield select option { background: #111; color: #f0f0f0; }
input:disabled { opacity: 0.5; cursor: not-allowed; }
.save-error { color: #e05a45; font-size: .78rem; margin-top: 12px; padding: 9px 13px; background: rgba(220,50,30,.06); border: 1px solid rgba(220,50,30,.25); }

.modal-footer { display: flex; align-items: center; justify-content: flex-end; gap: 10px; padding: 18px 26px; border-top: 1px solid rgba(245,197,0,.1); }
.btn-cancel { background: transparent; border: 1px solid #222; color: #555; font-family: 'Barlow Condensed',sans-serif; font-size: .72rem; font-weight: 700; letter-spacing: .2em; text-transform: uppercase; padding: 9px 20px; cursor: pointer; transition: all .2s; }
.btn-cancel:hover { border-color: rgba(245,197,0,.3); color: #aaa; }
.btn-save { display: inline-flex; align-items: center; gap: 8px; background: #f5c500; border: none; color: #0a0a0a; font-family: 'Barlow Condensed',sans-serif; font-size: .78rem; font-weight: 900; letter-spacing: .2em; text-transform: uppercase; padding: 10px 24px; cursor: pointer; transition: all .2s; }
.btn-save:hover:not(:disabled) { background: #ffd700; box-shadow: 0 0 20px rgba(245,197,0,.3); }
.btn-save:disabled { opacity: .5; cursor: not-allowed; }
.spinner { width: 14px; height: 14px; animation: spin .7s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.modal-fade-enter-active, .modal-fade-leave-active { transition: opacity .25s; }
.modal-fade-enter-from, .modal-fade-leave-to { opacity: 0; }

@media (max-width: 860px) {
  .sidebar { display: none; }
  .main { padding: 28px 20px 48px; }
  .modal-grid { grid-template-columns: 1fr; }
}

/* Filtros */
.filtros-box { display: flex; gap: 10px; flex-wrap: wrap; padding: 12px; background: #0d0d0d; border: 1px solid rgba(245,197,0,.1); }
.filtros-box input, .filtros-box select { background: #111; border: 1px solid rgba(245,197,0,.2); color: #fff; padding: 6px 10px; font-size: 0.8rem; }
.filtros-box button { background: #f5c500; border: none; padding: 6px 12px; cursor: pointer; font-weight: bold; }
.filtro-item { display: flex; flex-direction: column; gap: 4px; }
.filtro-item label { font-size: 0.6rem; letter-spacing: 0.15em; text-transform: uppercase; color: #f5c500; font-weight: 700; }
.btn-box { display: flex; gap: 8px; }
.btn-box button { flex: 1; }
</style>