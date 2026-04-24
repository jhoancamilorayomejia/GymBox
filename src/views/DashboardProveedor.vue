<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'

const router   = useRouter()
const username = localStorage.getItem('username') || 'Cliente'

const planes      = ref([])
const loading     = ref(true)
const errorMsg    = ref('')
const viewDate    = ref(new Date())
const verVencidos = ref(false)

const filtroActivos = ref({ datepay: '', datestart: '', datefinish: '' })
const filtroVencidos = ref({ datepay: '', datestart: '', datefinish: '' })

const matchFiltro = (p, f) => {
  const fmt = (str) => str ? str.split('T')[0] : ''
  if (f.datepay    && !fmt(p.datepay).includes(f.datepay))       return false
  if (f.datestart  && !fmt(p.datestart).includes(f.datestart))   return false
  if (f.datefinish && !fmt(p.datefinish).includes(f.datefinish)) return false
  return true
}

const planesFiltradosActivos = computed(() =>
  planesActivos.value.filter(p => matchFiltro(p, filtroActivos.value))
)
const planesFiltradosVencidos = computed(() =>
  planesVencidos.value.filter(p => matchFiltro(p, filtroVencidos.value))
)

const limpiarFiltroActivos  = () => { filtroActivos.value  = { datepay: '', datestart: '', datefinish: '' } }
const limpiarFiltroVencidos = () => { filtroVencidos.value = { datepay: '', datestart: '', datefinish: '' } }

const filtroActivosActivo  = computed(() => Object.values(filtroActivos.value).some(v => v))
const filtroVencidosActivo = computed(() => Object.values(filtroVencidos.value).some(v => v))

const MESES = ['Enero','Febrero','Marzo','Abril','Mayo','Junio','Julio','Agosto','Septiembre','Octubre','Noviembre','Diciembre']
const DIAS  = ['Dom','Lun','Mar','Mié','Jue','Vie','Sáb']

const ORDEN_PLANES = ['dia', 'semana', 'mes', 'anio']
const LABEL_PLANES = { dia: 'Día', semana: 'Semana', mes: 'Mes', anio: 'Año' }
const ICON_PLANES  = { dia: '☀️', semana: '📅', mes: '🗓', anio: '⚡' }

const parseFechaLocal = (fechaStr) => {
  const [year, month, day] = fechaStr.split('T')[0].split('-')
  return new Date(year, month - 1, day)
}
const normalizar = (fechaStr) => {
  const d = parseFechaLocal(fechaStr)
  d.setHours(0,0,0,0)
  return d
}

const tokenExpirado = (token) => {
  try { return JSON.parse(atob(token.split('.')[1])).exp * 1000 < Date.now() }
  catch { return true }
}
const cerrarSesion = () => { localStorage.clear(); router.push('/login') }

const obtenerPlanes = async () => {
  loading.value = true; errorMsg.value = ''
  try {
    const token = localStorage.getItem('token')
    if (!token || tokenExpirado(token)) { router.push('/login'); return }
    const resC = await fetch(`/api/customers/by-cedula/${username}`, { headers: { Authorization: `Bearer ${token}` } })
    if (!resC.ok) throw new Error('No se encontró el cliente')
    const cliente = await resC.json()
    const resP = await fetch(`/api/plans?idcustomer=${cliente.idcustomer}`, { headers: { Authorization: `Bearer ${token}` } })
    if (!resP.ok) throw new Error('Error al cargar los planes')
    planes.value = await resP.json() || []
  } catch (err) { errorMsg.value = err.message }
  finally { loading.value = false }
}

const hoyMidnight = () => { const h = new Date(); h.setHours(0,0,0,0); return h }

const planesActivos = computed(() => {
  const hoy = hoyMidnight()
  return planes.value.filter(p => p.state?.toLowerCase() === 'pagado' && normalizar(p.datefinish) >= hoy)
    .sort((a,b) => normalizar(a.datestart) - normalizar(b.datestart))
})
const planesVencidos = computed(() => {
  const hoy = hoyMidnight()
  return planes.value.filter(p => p.state?.toLowerCase() === 'pagado' && normalizar(p.datefinish) < hoy)
    .sort((a,b) => normalizar(b.datefinish) - normalizar(a.datefinish))
})

const mostrarModal     = ref(false)
const pricePlans       = ref([])
const planSeleccionado = ref(null)
const cantidadUnidades = ref(1)
const preferenceId     = ref(null)
const mostrarWallet    = ref(false)
const fechaInicio      = ref('')

const fechaMinima = computed(() => new Date().toISOString().split('T')[0])

const fechaFin = computed(() => {
  if (!fechaInicio.value || !planSeleccionado.value) return ''
  const [y, m, d] = fechaInicio.value.split('-').map(Number)
  const fecha = new Date(y, m - 1, d)
  const tipo  = planSeleccionado.value.typeplan
  const cant  = cantidadUnidades.value
  if (tipo === 'dia')    fecha.setDate(fecha.getDate() + cant - 1)
  if (tipo === 'semana') fecha.setDate(fecha.getDate() + cant * 7 - 1)
  if (tipo === 'mes')    fecha.setMonth(fecha.getMonth() + cant, fecha.getDate() - 1)
  if (tipo === 'anio')   fecha.setFullYear(fecha.getFullYear() + cant, fecha.getMonth(), fecha.getDate() - 1)
  const yy = fecha.getFullYear()
  const mm = String(fecha.getMonth() + 1).padStart(2, '0')
  const dd = String(fecha.getDate()).padStart(2, '0')
  return `${yy}-${mm}-${dd}`
})

const fmtDateInput = (str) => {
  if (!str) return '—'
  const [y, m, d] = str.split('-')
  return `${d}/${m}/${y}`
}

const descripcionRango = computed(() => {
  if (!fechaInicio.value || !planSeleccionado.value) return ''
  const tipo = planSeleccionado.value.typeplan
  const cant = cantidadUnidades.value
  const unidad = LABEL_PLANES[tipo] || tipo
  const plural = cant > 1
  const unidadPlural = tipo === 'mes' ? (plural ? 'Meses' : 'Mes')
                     : tipo === 'anio' ? (plural ? 'Años' : 'Año')
                     : unidad + (plural ? 's' : '')
  return `${cant} ${unidadPlural} · ${fmtDateInput(fechaInicio.value)} → ${fmtDateInput(fechaFin.value)}`
})

const obtenerPricePlans = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/customer-priceplans', { headers: { Authorization: `Bearer ${token}` } })
    if (!res.ok) throw new Error('Error al obtener planes')
    const data = await res.json()
    pricePlans.value = [...data].sort((a, b) => {
      const ia = ORDEN_PLANES.indexOf(a.typeplan)
      const ib = ORDEN_PLANES.indexOf(b.typeplan)
      return (ia === -1 ? 99 : ia) - (ib === -1 ? 99 : ib)
    })
  } catch (err) { console.error('Error cargando pricePlans:', err) }
}

const abrirModal = async () => {
  await obtenerPricePlans()
  planSeleccionado.value = null; cantidadUnidades.value = 1
  mostrarWallet.value = false; preferenceId.value = null
  fechaInicio.value = ''; mostrarModal.value = true
}
const cerrarModal = () => {
  mostrarModal.value = false; planSeleccionado.value = null
  cantidadUnidades.value = 1; mostrarWallet.value = false
  preferenceId.value = null; fechaInicio.value = ''
}
const seleccionarPlan = (plan) => {
  planSeleccionado.value = plan; cantidadUnidades.value = 1
  mostrarWallet.value = false; preferenceId.value = null
}

const totalAPagar = computed(() => {
  if (!planSeleccionado.value) return 0
  return planSeleccionado.value.price * cantidadUnidades.value
})

const realizarPago = async () => {
  if (!planSeleccionado.value) { alert('Selecciona un plan'); return }
  if (!fechaInicio.value) { alert('Selecciona una fecha de inicio'); return }
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/payment/preference', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
      body: JSON.stringify({
        amount: totalAPagar.value, username, planType: planSeleccionado.value.typeplan,
        quantity: cantidadUnidades.value, dateStart: fechaInicio.value, dateEnd: fechaFin.value
      })
    })
    if (!res.ok) throw new Error('Error al crear la preferencia')
    const data = await res.json()
    preferenceId.value = data.preferenceId
    mostrarWallet.value = true
    await nextTick()
    inicializarMercadoPago()
  } catch (err) { console.error('Error pago:', err); alert('Error al procesar el pago. Intenta nuevamente.') }
}

const inicializarMercadoPago = () => {
  const mp = new window.MercadoPago(import.meta.env.VITE_MP_PUBLIC_KEY, { locale: 'es-CO' })
  mp.bricks().create('wallet', 'wallet_container', {
    initialization: { preferenceId: preferenceId.value },
    customization: { texts: { valueProp: 'smart_option' } }
  })
}

const mostrarModalPassword = ref(false)
const guardandoPassword    = ref(false)
const passwordData = ref({ password: '', confirmPassword: '' })

const abrirModalPassword = () => {
  passwordData.value.password = ''; passwordData.value.confirmPassword = ''
  mostrarModalPassword.value = true
}
const cerrarModalPassword = () => {
  mostrarModalPassword.value = false
  passwordData.value.password = ''; passwordData.value.confirmPassword = ''
}

const cambiarPassword = async () => {
  if (!passwordData.value.password || !passwordData.value.confirmPassword) {
    alert('Completa todos los campos'); return
  }
  if (passwordData.value.password !== passwordData.value.confirmPassword) {
    alert('Las contraseñas no coinciden'); return
  }
  if (passwordData.value.password.length < 8) {
    alert('La contraseña debe tener al menos 8 caracteres'); return
  }
  const tieneLetras  = /[a-zA-Z]/.test(passwordData.value.password)
  const tieneNumeros = /[0-9]/.test(passwordData.value.password)
  if (!tieneLetras || !tieneNumeros) {
    alert('La contraseña debe contener al menos una letra y un número'); return
  }
  guardandoPassword.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/users/update-password-by-username', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
      body: JSON.stringify({ username, password: passwordData.value.password })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Error al cambiar contraseña')
    alert('✅ Contraseña actualizada correctamente')
    cerrarModalPassword()
  } catch (err) { console.error(err); alert('❌ ' + err.message) }
  finally { guardandoPassword.value = false }
}

const fmtFecha = (str) => {
  if (!str) return '—'
  const [year, month, day] = str.split('T')[0].split('-')
  return `${day}/${month}/${year}`
}
const fmtFechaLarga = (str) => {
  if (!str) return '—'
  return parseFechaLocal(str).toLocaleDateString('es-CO', { day:'numeric', month:'long', year:'numeric' })
}
const estadoClass = (s) => {
  if (!s) return 'badge-default'
  if (s.toLowerCase() === 'pagado') return 'badge-pagado'
  if (s.toLowerCase() === 'pendiente') return 'badge-pendiente'
  return 'badge-default'
}
const dispFila = (p) => {
  const hoy = new Date(); hoy.setHours(0,0,0,0)
  const inicio = new Date(p.datestart); inicio.setHours(0,0,0,0)
  const fin = new Date(p.datefinish); fin.setHours(0,0,0,0)
  if (inicio > hoy) return { label: 'Próximo inicio', clase: 'disp-proximo' }
  if (fin < hoy)    return { label: 'Vencido', clase: 'disp-vencido' }
  const diff = Math.round((fin - hoy) / 86400000)
  if (diff === 0)   return { label: 'Vence hoy', clase: 'disp-vencer-hoy' }
  return { label: `${diff} Días restantes`, clase: 'disp-activo' }
}

const resumenPlanes = computed(() => {
  const hoy = hoyMidnight()
  return planesActivos.value.map(p => {
    const ini = normalizar(p.datestart)
    const fin = normalizar(p.datefinish)
    const esFuturo = hoy < ini
    const diasRest = Math.max(0, Math.ceil((fin-hoy)/86400000))
    let estado, mensaje
    if (esFuturo) { const d=Math.ceil((ini-hoy)/86400000); estado='proximo'; mensaje=`Inicia en ${d} día${d!==1?'s':''}, el ${fmtFechaLarga(p.datestart)}.` }
    else if (diasRest===0) { estado='vence-hoy'; mensaje='Este plan vence HOY.' }
    else { estado='activo'; mensaje=`Acceso válido hasta el ${fmtFechaLarga(p.datefinish)}. Quedan ${diasRest} días.` }
    return { plan:p, estado, mensaje }
  })
})

const mesLabel  = computed(() => `${MESES[viewDate.value.getMonth()]} ${viewDate.value.getFullYear()}`)
const primerDia = computed(() => new Date(viewDate.value.getFullYear(), viewDate.value.getMonth(), 1).getDay())
const diasEnMes = computed(() => new Date(viewDate.value.getFullYear(), viewDate.value.getMonth()+1, 0).getDate())
const cambiarMes = (delta) => { const d=new Date(viewDate.value); d.setMonth(d.getMonth()+delta); viewDate.value=d }

const clasificarDia = (dia) => {
  const y=viewDate.value.getFullYear(), m=viewDate.value.getMonth()
  const fecha=new Date(y,m,dia); fecha.setHours(0,0,0,0)
  const hoy=hoyMidnight()
  const esHoy=fecha.getTime()===hoy.getTime()
  let disponible=false, vencido=false
  for (const p of planes.value) {
    if (p.state?.toLowerCase()!=='pagado') continue
    const ini=normalizar(p.datestart), fin=normalizar(p.datefinish)
    if (fecha>=ini && fecha<=fin) { if(fin<hoy){vencido=true}else{disponible=true;break} }
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
          <button class="change-pass-btn" @click="abrirModalPassword">🔑 Cambiar contraseña</button>
          <div class="suscripcion-wrap">
            <button class="btn-suscripcion" @click="abrirModal">
              <span class="btn-sus-icon"><svg viewBox="0 0 20 20" fill="none"><path d="M10 4v12M4 10h12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg></span>
              <span class="btn-sus-text"><span class="btn-sus-title">Nueva suscripción</span><span class="btn-sus-sub">Renueva o agrega un plan</span></span>
              <span class="btn-sus-bolt">⚡</span>
            </button>
          </div>
        </nav>
        <div class="sidebar-footer">
          <div class="user-chip">
            <div class="user-avatar">{{ username.charAt(0).toUpperCase() }}</div>
            <div class="user-info">
              <span class="user-name">{{ username }}</span>
              <span class="user-role">Cliente</span>
            </div>
          </div>
          <button class="logout-btn" @click="cerrarSesion" title="Cerrar sesión">
            <svg viewBox="0 0 20 20" fill="none"><path d="M13 3h4v14h-4M9 14l4-4-4-4M13 10H5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
          </button>
        </div>
      </aside>

      <!-- ── Main ── -->
      <main class="main">
        <header class="page-header">
          <div class="page-title">
            <span class="page-eyebrow">⚡ Bienvenido de vuelta</span>
            <h1>{{ username }}</h1>
          </div>
          <div class="header-meta">
            <div class="stat-pill"><span class="stat-val">{{ planes.length }}</span><span class="stat-lbl">Total</span></div>
            <div class="stat-pill accent"><span class="stat-val">{{ planesActivos.length }}</span><span class="stat-lbl">Vigentes</span></div>
            <div class="stat-pill muted" v-if="planesVencidos.length"><span class="stat-val">{{ planesVencidos.length }}</span><span class="stat-lbl">Vencidos</span></div>
          </div>
        </header>

        <!-- ── ACCIONES MÓVIL ── -->
        <div class="mobile-actions">
          <button class="mob-btn-pass" @click="abrirModalPassword">
            <span class="mob-btn-icon">🔑</span><span class="mob-btn-text">Cambiar contraseña</span>
          </button>
          <button class="mob-btn-sus" @click="abrirModal">
            <span class="mob-btn-icon">⚡</span><span class="mob-btn-text">Nueva suscripción</span>
          </button>
          <button class="mob-btn-logout" @click="cerrarSesion">
            <svg viewBox="0 0 20 20" fill="none"><path d="M13 3h4v14h-4M9 14l4-4-4-4M13 10H5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
          </button>
        </div>

        <div v-if="loading" class="state-box">
          <div class="bolt-spin">⚡</div><p>Cargando tus planes...</p>
        </div>
        <div v-else-if="errorMsg" class="state-box error-state">
          <span>⚠</span><p>{{ errorMsg }}</p>
          <button class="retry-btn" @click="obtenerPlanes">Reintentar</button>
        </div>

        <template v-else>
          <div class="top-row">
            <div class="resumenes-col">
              <div v-if="resumenPlanes.length === 0" class="resumen-card resumen-sin-plan">
                <div class="resumen-top"><span class="resumen-icon">🔒</span><span class="resumen-badge rbadge-sin-plan">SIN PLAN ACTIVO</span></div>
                <p class="resumen-msg">No tienes un plan activo.</p>
                <div class="leyenda">
                  <div class="ley-item"><span class="ley-dot dot-activo"></span><span>Día disponible</span></div>
                  <div class="ley-item"><span class="ley-dot dot-vencido-cal"></span><span>Plan vencido</span></div>
                  <div class="ley-item"><span class="ley-dot dot-ninguno"></span><span>Sin acceso</span></div>
                </div>
              </div>
              <div v-for="(r,i) in resumenPlanes" :key="r.plan.idplan" class="resumen-card" :class="`resumen-${r.estado}`" :style="{ animationDelay: i*80+'ms' }">
                <div class="resumen-top">
                  <span class="resumen-icon">{{ r.estado==='activo'?'⚡':r.estado==='vence-hoy'?'⏳':'🗓' }}</span>
                  <span class="resumen-badge" :class="`rbadge-${r.estado}`">{{ r.estado==='activo'?'ACCESO ACTIVO':r.estado==='vence-hoy'?'VENCE HOY':'PRÓXIMO INICIO' }}</span>
                  <span class="plan-tipo-tag">{{ r.plan.typeplan?.toUpperCase() }}</span>
                </div>
                <p class="resumen-msg">{{ r.mensaje }}</p>
                <div class="resumen-fechas">
                  <div class="rfecha"><span class="rfecha-key">Inicio</span><span class="rfecha-val">{{ fmtFecha(r.plan.datestart) }}</span></div>
                  <span class="rfecha-arrow">→</span>
                  <div class="rfecha"><span class="rfecha-key">Fin</span><span class="rfecha-val">{{ fmtFecha(r.plan.datefinish) }}</span></div>
                  <div class="rfecha" v-if="r.estado==='activo'||r.estado==='vence-hoy'"><span class="rfecha-key">Precio</span><span class="rfecha-val precio-tag">${{ Number(r.plan.price).toLocaleString('es-CO') }}</span></div>
                </div>
                <div class="leyenda" v-if="i===0">
                  <div class="ley-item"><span class="ley-dot dot-activo"></span><span>Día disponible</span></div>
                  <div class="ley-item"><span class="ley-dot dot-vencido-cal"></span><span>Plan vencido</span></div>
                  <div class="ley-item"><span class="ley-dot dot-ninguno"></span><span>Sin acceso</span></div>
                </div>
              </div>
            </div>

            <div class="cal-card">
              <div class="cal-nav">
                <button class="cal-arrow" @click="cambiarMes(-1)">&#8249;</button>
                <span class="cal-mes">{{ mesLabel }}</span>
                <button class="cal-arrow" @click="cambiarMes(1)">&#8250;</button>
              </div>
              <div class="cal-header"><span v-for="d in DIAS" :key="d">{{ d }}</span></div>
              <div class="cal-grid">
                <div v-for="n in primerDia" :key="'e'+n" class="cal-empty"></div>
                <div v-for="dia in diasEnMes" :key="dia" class="cal-dia"
                  :class="{ 'dia-activo':clasificarDia(dia).disponible, 'dia-vencido':!clasificarDia(dia).disponible&&clasificarDia(dia).vencido, 'dia-hoy':clasificarDia(dia).esHoy, 'dia-ninguno':!clasificarDia(dia).disponible&&!clasificarDia(dia).vencido }">
                  {{ dia }}
                </div>
              </div>
            </div>
          </div>

          <!-- ── Tabla planes vigentes ── -->
          <div class="section-label">
            <svg viewBox="0 0 20 20" fill="none"><rect x="3" y="3" width="14" height="14" rx="1.5" stroke="currentColor" stroke-width="1.4"/><path d="M7 10h6M10 7v6" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/></svg>
            Planes vigentes
          </div>

          <div class="filtro-bar">
            <div class="filtro-fields">
              <div class="filtro-field">
                <label class="filtro-label">FECHA DE PAGO</label>
                <input type="date" class="filtro-date" v-model="filtroActivos.datepay" />
              </div>
              <div class="filtro-field">
                <label class="filtro-label">FECHA INICIO</label>
                <input type="date" class="filtro-date" v-model="filtroActivos.datestart" />
              </div>
              <div class="filtro-field">
                <label class="filtro-label">FECHA FIN</label>
                <input type="date" class="filtro-date" v-model="filtroActivos.datefinish" />
              </div>
            </div>
            <div class="filtro-meta">
              <span class="filtro-count-total">{{ planesFiltradosActivos.length }} / {{ planesActivos.length }} planes</span>
              <button v-if="filtroActivosActivo" class="filtro-reset" @click="limpiarFiltroActivos">✕ Limpiar</button>
            </div>
          </div>

          <div class="table-wrapper">
            <div class="table-scroll">
              <table>
                <thead><tr><th>#</th><th>Tipo</th><th>Precio</th><th>Fecha Pago</th><th>Inicio</th><th>Fin</th><th>Disponibilidad</th><th>Estado</th></tr></thead>
                <tbody>
                  <tr v-for="(p,i) in planesFiltradosActivos" :key="p.idplan" class="table-row" :style="{ animationDelay:i*45+'ms' }">
                    <td><span class="id-badge">{{ p.idplan }}</span></td>
                    <td><span class="tipo-label">{{ p.typeplan }}</span></td>
                    <td><span class="precio-val">${{ Number(p.price).toLocaleString('es-CO') }}</span></td>
                    <td>{{ fmtFecha(p.datepay) }}</td>
                    <td>{{ fmtFecha(p.datestart) }}</td>
                    <td>{{ fmtFecha(p.datefinish) }}</td>
                    <td><span class="disp-badge" :class="dispFila(p).clase">{{ dispFila(p).label }}</span></td>
                    <td><span class="estado-badge" :class="estadoClass(p.state)">{{ p.state }}</span></td>
                  </tr>
                  <tr v-if="!planesFiltradosActivos.length">
                    <td colspan="8" class="empty-row">{{ filtroActivosActivo ? '🔍 Sin resultados para ese filtro.' : '⚡ No tienes planes vigentes. Habla con tu entrenador.' }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div class="table-footer">
              {{ planesFiltradosActivos.length }} plan{{ planesFiltradosActivos.length!==1?'es':'' }} vigente{{ planesFiltradosActivos.length!==1?'s':'' }}
              <span v-if="filtroActivosActivo" class="filtro-footer-hint"> · filtrado de {{ planesActivos.length }} total</span>
            </div>
          </div>

          <!-- ── Planes vencidos ── -->
          <div class="vencidos-header">
            <div class="section-label">
              <svg viewBox="0 0 20 20" fill="none"><path d="M10 3a7 7 0 100 14A7 7 0 0010 3zM10 7v4l2.5 2.5" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round"/></svg>
              Planes vencidos <span class="venc-count" v-if="planesVencidos.length">{{ planesVencidos.length }}</span>
            </div>
            <button class="btn-ver-vencidos" @click="verVencidos=!verVencidos" v-if="planesVencidos.length">
              <svg viewBox="0 0 20 20" fill="none"><path v-if="!verVencidos" d="M4 8l6 6 6-6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/><path v-else d="M4 12l6-6 6 6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
              {{ verVencidos?'Ocultar':'Ver planes vencidos' }}
            </button>
            <span class="venc-vacio" v-else>Sin historial de planes vencidos</span>
          </div>

          <transition name="expand">
            <div v-if="verVencidos && planesVencidos.length" class="vencidos-panel">
              <div class="venc-top-row">
                <div class="venc-info-card">
                  <div class="section-label" style="margin-bottom:12px">
                    <svg viewBox="0 0 20 20" fill="none"><circle cx="10" cy="10" r="7" stroke="currentColor" stroke-width="1.4"/><path d="M10 7v4l2 2" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/></svg>
                    Historial anterior
                  </div>
                  <div class="venc-resumenes">
                    <div v-for="p in planesVencidos" :key="p.idplan" class="venc-item">
                      <div class="venc-tipo">{{ p.typeplan?.toUpperCase() }}</div>
                      <div class="venc-fechas-row"><span>{{ fmtFecha(p.datestart) }}</span><span class="venc-arrow">→</span><span>{{ fmtFecha(p.datefinish) }}</span></div>
                      <div class="venc-precio">${{ Number(p.price).toLocaleString('es-CO') }}</div>
                    </div>
                  </div>
                </div>
                <div class="cal-card cal-vencidos">
                  <div class="cal-nav">
                    <button class="cal-arrow" @click="cambiarMes(-1)">&#8249;</button>
                    <span class="cal-mes">{{ mesLabel }}</span>
                    <button class="cal-arrow" @click="cambiarMes(1)">&#8250;</button>
                  </div>
                  <div class="cal-header"><span v-for="d in DIAS" :key="'v'+d">{{ d }}</span></div>
                  <div class="cal-grid">
                    <div v-for="n in primerDia" :key="'ve'+n" class="cal-empty"></div>
                    <div v-for="dia in diasEnMes" :key="'vd'+dia" class="cal-dia"
                      :class="{ 'dia-vencido':clasificarDia(dia).vencido&&!clasificarDia(dia).disponible, 'dia-ninguno':!clasificarDia(dia).vencido||clasificarDia(dia).disponible }">
                      {{ dia }}
                    </div>
                  </div>
                </div>
              </div>

              <div class="filtro-bar">
                <div class="filtro-fields">
                  <div class="filtro-field">
                    <label class="filtro-label">FECHA DE PAGO</label>
                    <input type="date" class="filtro-date filtro-date-venc" v-model="filtroVencidos.datepay" />
                  </div>
                  <div class="filtro-field">
                    <label class="filtro-label filtro-label-venc">FECHA INICIO</label>
                    <input type="date" class="filtro-date filtro-date-venc" v-model="filtroVencidos.datestart" />
                  </div>
                  <div class="filtro-field">
                    <label class="filtro-label filtro-label-venc">FECHA FIN</label>
                    <input type="date" class="filtro-date filtro-date-venc" v-model="filtroVencidos.datefinish" />
                  </div>
                </div>
                <div class="filtro-meta">
                  <span class="filtro-count-total">{{ planesFiltradosVencidos.length }} / {{ planesVencidos.length }} planes</span>
                  <button v-if="filtroVencidosActivo" class="filtro-reset" @click="limpiarFiltroVencidos">✕ Limpiar</button>
                </div>
              </div>

              <div class="table-wrapper table-vencidos">
                <div class="table-scroll">
                  <table>
                    <thead><tr><th>#</th><th>Tipo</th><th>Precio</th><th>Fecha Pago</th><th>Inicio</th><th>Fin</th><th>Estado pago</th></tr></thead>
                    <tbody>
                      <tr v-for="(p,i) in planesFiltradosVencidos" :key="p.idplan" class="table-row row-vencido" :style="{ animationDelay:i*40+'ms' }">
                        <td><span class="id-badge id-badge-venc">{{ p.idplan }}</span></td>
                        <td><span class="tipo-label tipo-venc">{{ p.typeplan }}</span></td>
                        <td><span class="precio-val precio-venc">${{ Number(p.price).toLocaleString('es-CO') }}</span></td>
                        <td>{{ fmtFecha(p.datepay) }}</td>
                        <td>{{ fmtFecha(p.datestart) }}</td>
                        <td>{{ fmtFecha(p.datefinish) }}</td>
                        <td><span class="estado-badge" :class="estadoClass(p.state)">{{ p.state }}</span></td>
                      </tr>
                      <tr v-if="!planesFiltradosVencidos.length">
                        <td colspan="7" class="empty-row">{{ filtroVencidosActivo ? '🔍 Sin resultados para ese filtro.' : 'Sin planes vencidos.' }}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <div class="table-footer">
                  {{ planesFiltradosVencidos.length }} plan{{ planesFiltradosVencidos.length!==1?'es':'' }} vencido{{ planesFiltradosVencidos.length!==1?'s':'' }} en el historial
                  <span v-if="filtroVencidosActivo" class="filtro-footer-hint"> · filtrado de {{ planesVencidos.length }} total</span>
                </div>
              </div>
            </div>
          </transition>

        </template>
      </main>
    </div>

    <!-- MODAL DE PAGO -->
    <transition name="modal-fade">
      <div v-if="mostrarModal" class="modal-overlay" @click.self="cerrarModal">
        <div class="modal-pay">
          <div class="mpay-header">
            <div>
              <span class="mpay-eyebrow">⚡ RAYOBOX · CROSS LIFTING</span>
              <h2 class="mpay-title">Nueva Suscripción</h2>
              <p class="mpay-sub">Selecciona tu plan, fechas y cantidad</p>
            </div>
            <button class="mpay-close" @click="cerrarModal">✕</button>
          </div>
          <div class="mpay-body">
            <div class="mpay-left">
              <p class="mpay-section-label">ELIGE UN PLAN</p>
              <div class="mpay-planes">
                <button v-for="p in pricePlans" :key="p.idpriceplan" class="mpay-plan-card" :class="{ 'mpay-plan-activo': planSeleccionado?.idpriceplan === p.idpriceplan }" @click="seleccionarPlan(p)">
                  <span class="mpay-plan-icon">{{ ICON_PLANES[p.typeplan] || '📌' }}</span>
                  <div class="mpay-plan-info">
                    <span class="mpay-plan-nombre">{{ LABEL_PLANES[p.typeplan] || p.typeplan }}</span>
                    <span class="mpay-plan-precio">${{ Number(p.price).toLocaleString('es-CO') }}</span>
                  </div>
                  <span class="mpay-plan-check" v-if="planSeleccionado?.idpriceplan === p.idpriceplan">✓</span>
                </button>
              </div>
              <div class="mpay-cantidad-box" v-if="planSeleccionado">
                <p class="mpay-section-label" style="margin-bottom:10px">CANTIDAD DE {{ (LABEL_PLANES[planSeleccionado.typeplan] || planSeleccionado.typeplan).toUpperCase() + (cantidadUnidades > 1 ? 'S' : '') }}</p>
                <div class="mpay-cantidad-ctrl">
                  <button class="mpay-qty-btn" @click="cantidadUnidades = Math.max(1, cantidadUnidades - 1)">−</button>
                  <span class="mpay-qty-val">{{ cantidadUnidades }}</span>
                  <button class="mpay-qty-btn" @click="cantidadUnidades++">+</button>
                </div>
                <p class="mpay-cantidad-hint">{{ cantidadUnidades }} {{ LABEL_PLANES[planSeleccionado.typeplan] || planSeleccionado.typeplan }}{{ cantidadUnidades > 1 ? 's' : '' }} de membresía</p>
              </div>
              <div class="mpay-fecha-box" v-if="planSeleccionado">
                <p class="mpay-section-label" style="margin-bottom:10px">FECHA DE INICIO</p>
                <input type="date" class="mpay-date-input" v-model="fechaInicio" :min="fechaMinima"/>
                <div class="mpay-fecha-fin" v-if="fechaFin">
                  <div class="mpay-fecha-row">
                    <div class="mpay-fecha-item"><span class="mpay-fecha-key">INICIO</span><span class="mpay-fecha-val">{{ fmtDateInput(fechaInicio) }}</span></div>
                    <span class="mpay-fecha-arrow">→</span>
                    <div class="mpay-fecha-item"><span class="mpay-fecha-key">FIN</span><span class="mpay-fecha-val">{{ fmtDateInput(fechaFin) }}</span></div>
                  </div>
                  <p class="mpay-fecha-desc">{{ descripcionRango }}</p>
                </div>
              </div>
            </div>
            <div class="mpay-right" :class="{ 'mpay-right-vacio': !planSeleccionado }">
              <div v-if="!planSeleccionado" class="mpay-vacio">
                <span class="mpay-vacio-icon">⚡</span>
                <p>Selecciona un plan<br>para ver el resumen</p>
              </div>
              <template v-else>
                <p class="mpay-section-label">RESUMEN DE PAGO</p>
                <div class="mpay-resumen">
                  <div class="mpay-resumen-plan">
                    <span class="mpay-resumen-icon">{{ ICON_PLANES[planSeleccionado.typeplan] || '📌' }}</span>
                    <div>
                      <span class="mpay-resumen-tipo">{{ LABEL_PLANES[planSeleccionado.typeplan] || planSeleccionado.typeplan }}</span>
                      <span class="mpay-resumen-unit">× {{ cantidadUnidades }} unidad{{ cantidadUnidades > 1 ? 'es' : '' }}</span>
                    </div>
                  </div>
                  <div class="mpay-lineas">
                    <div class="mpay-linea"><span>Precio unitario</span><span>${{ Number(planSeleccionado.price).toLocaleString('es-CO') }}</span></div>
                    <div class="mpay-linea" v-if="cantidadUnidades > 1"><span>Cantidad</span><span>× {{ cantidadUnidades }}</span></div>
                    <div class="mpay-linea" v-if="fechaInicio"><span>Fecha inicio</span><span>{{ fmtDateInput(fechaInicio) }}</span></div>
                    <div class="mpay-linea" v-if="fechaFin"><span>Fecha fin</span><span>{{ fmtDateInput(fechaFin) }}</span></div>
                    <div class="mpay-divider"></div>
                    <div class="mpay-linea mpay-linea-total"><span>TOTAL</span><span class="mpay-total-val">${{ totalAPagar.toLocaleString('es-CO') }}</span></div>
                  </div>
                  <div class="mpay-desglose" v-if="cantidadUnidades > 1">
                    <p class="mpay-desglose-label">Desglose</p>
                    <div class="mpay-desglose-chips">
                      <span v-for="n in cantidadUnidades" :key="n" class="mpay-chip">{{ LABEL_PLANES[planSeleccionado.typeplan] || planSeleccionado.typeplan }} {{ n }}: ${{ Number(planSeleccionado.price * n).toLocaleString('es-CO') }}</span>
                    </div>
                  </div>
                </div>
                <div class="mpay-alerta" v-if="!fechaInicio"><span>📅</span> Selecciona una fecha de inicio para continuar</div>
                <button class="mpay-btn-pagar" @click="realizarPago" :disabled="mostrarWallet || !fechaInicio">
                  <span class="mpay-pagar-bolt">⚡</span>
                  <span class="mpay-pagar-text">
                    <span class="mpay-pagar-titulo">{{ mostrarWallet ? 'Procesando...' : 'Pagar ahora' }}</span>
                    <span class="mpay-pagar-monto">${{ totalAPagar.toLocaleString('es-CO') }}</span>
                  </span>
                  <svg class="mpay-pagar-arrow" viewBox="0 0 20 20" fill="none"><path d="M4 10h12M10 4l6 6-6 6" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/></svg>
                </button>
                <div v-if="mostrarWallet" id="wallet_container"></div>
                <p class="mpay-nota">Al realizar el pago, tu entrenador activará el plan en el sistema.</p>
              </template>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- MODAL CAMBIAR PASSWORD -->
    <transition name="modal-fade">
      <div v-if="mostrarModalPassword" class="modal-overlay" @click.self="cerrarModalPassword">
        <div class="modal-pay" style="max-width: 420px;">
          <div class="mpay-header">
            <div>
              <span class="mpay-eyebrow">⚡ RAYOBOX · SEGURIDAD</span>
              <h2 class="mpay-title">Cambiar contraseña</h2>
              <p class="mpay-sub">Usuario: <strong style="color:#f5c500">{{ username }}</strong></p>
            </div>
            <button class="mpay-close" @click="cerrarModalPassword">✕</button>
          </div>
          <div class="mpay-body" style="display:flex; flex-direction:column; gap:16px; padding:24px 28px;">
            <div class="pass-field">
              <label class="pass-label">NUEVA CONTRASEÑA</label>
              <input type="password" placeholder="Mínimo 8 caracteres, letras y números" v-model="passwordData.password" class="mpay-date-input"/>
            </div>
            <div class="pass-field">
              <label class="pass-label">CONFIRMAR CONTRASEÑA</label>
              <input type="password" placeholder="Repite la contraseña" v-model="passwordData.confirmPassword" class="mpay-date-input"/>
            </div>
            <div v-if="passwordData.password" class="pass-requisitos">
              <div class="req-item" :class="passwordData.password.length >= 8 ? 'req-ok' : 'req-fail'">
                {{ passwordData.password.length >= 8 ? '✓' : '✕' }} Mínimo 8 caracteres
              </div>
              <div class="req-item" :class="/[a-zA-Z]/.test(passwordData.password) ? 'req-ok' : 'req-fail'">
                {{ /[a-zA-Z]/.test(passwordData.password) ? '✓' : '✕' }} Al menos una letra
              </div>
              <div class="req-item" :class="/[0-9]/.test(passwordData.password) ? 'req-ok' : 'req-fail'">
                {{ /[0-9]/.test(passwordData.password) ? '✓' : '✕' }} Al menos un número
              </div>
            </div>
            <div v-if="passwordData.confirmPassword" class="pass-match" :class="passwordData.password === passwordData.confirmPassword ? 'match-ok' : 'match-fail'">
              {{ passwordData.password === passwordData.confirmPassword ? '✓ Las contraseñas coinciden' : '✕ Las contraseñas no coinciden' }}
            </div>
            <button class="mpay-btn-pagar" @click="cambiarPassword" :disabled="guardandoPassword" style="margin-top:4px">
              <span class="mpay-pagar-bolt">🔑</span>
              <span class="mpay-pagar-text">
                <span class="mpay-pagar-titulo">{{ guardandoPassword ? 'Guardando...' : 'Guardar contraseña' }}</span>
                <span class="mpay-pagar-monto">{{ username }}</span>
              </span>
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

/* ═══════════════════════════════════════
   BASE SCREEN & BACKGROUND
═══════════════════════════════════════ */
.screen {
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  background: #0a0a0a;
  font-family: 'Barlow', sans-serif;
  position: relative;
}
.bg-layer { position: fixed; inset: 0; pointer-events: none; z-index: 0; }
.hex { position: absolute; }
.hex-1 { width: 420px; height: 420px; background: conic-gradient(from 30deg,rgba(245,197,0,.08) 0deg,rgba(255,122,0,.06) 80deg,transparent 80deg); clip-path: polygon(50% 0%,100% 25%,100% 75%,50% 100%,0% 75%,0% 25%); top: -100px; right: -80px; }
.hex-2 { width: 260px; height: 260px; background: conic-gradient(from 200deg,rgba(245,197,0,.05) 0deg,transparent 80deg); clip-path: polygon(50% 0%,100% 25%,100% 75%,50% 100%,0% 75%,0% 25%); bottom: -60px; left: 260px; }
.slash { position: absolute; }
.slash-1 { width: 2px; height: 60vh; background: linear-gradient(to bottom,transparent,rgba(245,197,0,.12) 50%,transparent); top: 20%; left: 260px; transform: rotate(12deg); }
.slash-2 { width: 1px; height: 40vh; background: linear-gradient(to bottom,transparent,rgba(255,122,0,.08) 50%,transparent); top: 30%; left: 270px; transform: rotate(12deg); }
.noise { position: absolute; inset: 0; background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='1'/%3E%3C/svg%3E"); opacity: .02; }

/* ═══════════════════════════════════════
   LAYOUT — CORRECCIÓN PRINCIPAL
   El sidebar es una columna vertical fija a la izquierda
═══════════════════════════════════════ */
.layout {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: row;   /* sidebar izquierda + main derecha */
  width: 100%;
  height: 100vh;
}

/* ═══════════════════════════════════════
   SIDEBAR — DESKTOP
═══════════════════════════════════════ */
.sidebar {
  width: 240px;
  min-width: 240px;
  flex-shrink: 0;
  height: 100vh;
  background: #0d0d0d;
  border-right: 1px solid rgba(245,197,0,.1);
  display: flex;
  flex-direction: column;   /* logo → nav → footer en columna */
  padding-top: 32px;
  overflow-y: auto;
  overflow-x: hidden;
}

.sidebar-logo {
  padding: 0 24px 28px;
  border-bottom: 1px solid rgba(245,197,0,.08);
  margin-bottom: 24px;
  flex-shrink: 0;
}
.sidebar-logo svg { width: 130px; height: auto; filter: drop-shadow(0 0 18px rgba(245,197,0,.18)); }
.sidebar-sub { display: block; font-family: 'Barlow Condensed',sans-serif; font-size: .55rem; letter-spacing: .35em; text-transform: uppercase; color: #444; margin-top: 4px; }

.sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 0 12px;
  flex: 1;
}

.nav-item { display: flex; align-items: center; gap: 10px; padding: 11px 14px; color: #555; font-size: .82rem; font-weight: 500; letter-spacing: .04em; text-decoration: none; border-left: 2px solid transparent; transition: all .2s; cursor: pointer; }
.nav-item svg { width: 16px; height: 16px; flex-shrink: 0; }
.nav-item:hover { color: #f0f0f0; background: rgba(245,197,0,.05); }
.nav-item.active { color: #f5c500; border-left-color: #f5c500; background: rgba(245,197,0,.07); }

.change-pass-btn { display: flex; align-items: center; gap: 8px; padding: 10px 14px; width: 100%; text-align: left; background: transparent; border: none; border-left: 2px solid transparent; color: #555; font-size: .82rem; font-weight: 500; letter-spacing: .04em; cursor: pointer; transition: all .2s; font-family: 'Barlow', sans-serif; }
.change-pass-btn:hover { color: #f0f0f0; background: rgba(245,197,0,.05); border-left-color: rgba(245,197,0,.3); }

.suscripcion-wrap { padding: 4px 2px 0; margin-top: 8px; }
.btn-suscripcion { width: 100%; display: flex; align-items: center; gap: 10px; background: linear-gradient(135deg, rgba(245,197,0,.12) 0%, rgba(255,122,0,.08) 100%); border: 1px solid rgba(245,197,0,.35); padding: 12px 14px; cursor: pointer; transition: all .25s; position: relative; overflow: hidden; text-align: left; }
.btn-suscripcion::before { content: ''; position: absolute; left: 0; top: 0; bottom: 0; width: 3px; background: linear-gradient(to bottom, #f5c500, #ff7a00); }
.btn-suscripcion:hover { background: linear-gradient(135deg, rgba(245,197,0,.2) 0%, rgba(255,122,0,.14) 100%); border-color: rgba(245,197,0,.6); box-shadow: 0 0 20px rgba(245,197,0,.12); }
.btn-sus-icon { width: 28px; height: 28px; background: rgba(245,197,0,.15); border: 1px solid rgba(245,197,0,.3); display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.btn-sus-icon svg { width: 13px; height: 13px; stroke: #f5c500; }
.btn-sus-text { display: flex; flex-direction: column; gap: 1px; flex: 1; min-width: 0; }
.btn-sus-title { font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .82rem; letter-spacing: .08em; text-transform: uppercase; color: #f5c500; }
.btn-sus-sub { font-size: .6rem; color: #666; letter-spacing: .04em; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.btn-sus-bolt { font-size: 1rem; opacity: .4; }

.sidebar-footer {
  padding: 20px 16px;
  border-top: 1px solid rgba(245,197,0,.08);
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
  margin-top: auto;
}
.user-chip { display: flex; align-items: center; gap: 10px; flex: 1; min-width: 0; }
.user-avatar { width: 32px; height: 32px; background: linear-gradient(135deg,#f5c500,#ff7a00); color: #0a0a0a; font-family: 'Barlow Condensed',sans-serif; font-weight: 900; font-size: .95rem; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.user-info { display: flex; flex-direction: column; min-width: 0; }
.user-name { font-size: .78rem; color: #e0e0e0; font-weight: 500; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.user-role { font-size: .6rem; color: #555; letter-spacing: .06em; }
.logout-btn { background: none; border: 1px solid #1e1e1e; color: #444; cursor: pointer; width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; transition: all .2s; }
.logout-btn svg { width: 15px; height: 15px; }
.logout-btn:hover { border-color: rgba(245,197,0,.4); color: #f5c500; }

/* ═══════════════════════════════════════
   MAIN CONTENT
═══════════════════════════════════════ */
.main {
  flex: 1;
  height: 100vh;
  overflow-y: auto;
  padding: 48px 60px;
  display: flex;
  flex-direction: column;
  gap: 22px;
  min-width: 0;
}

.page-header { display: flex; align-items: flex-end; justify-content: space-between; flex-wrap: wrap; gap: 16px; flex-shrink: 0; }
.page-eyebrow { display: block; font-family: 'Barlow Condensed',sans-serif; font-size: .65rem; letter-spacing: .3em; text-transform: uppercase; color: #f5c500; margin-bottom: 6px; }
.page-title h1 { font-family: 'Barlow Condensed',sans-serif; font-size: 3rem; font-weight: 900; text-transform: uppercase; letter-spacing: .04em; color: #fff; line-height: 1; }
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
@keyframes boltPulse { 0%,100%{opacity:1}50%{opacity:.3} }
.error-state { color: #e05a45; }
.retry-btn { background: none; border: 1px solid rgba(220,50,30,.4); color: #e05a45; font-family: 'Barlow Condensed',sans-serif; font-size: .75rem; letter-spacing: .2em; text-transform: uppercase; padding: 8px 20px; cursor: pointer; transition: all .2s; }
.retry-btn:hover { background: rgba(220,50,30,.1); }

/* ── Top row (resumen + calendario) ── */
.top-row { display: grid; grid-template-columns: 1fr 1.2fr; gap: 18px; flex-shrink: 0; }
.resumenes-col { display: flex; flex-direction: column; gap: 12px; }
.resumen-card { background: #111; border: 1px solid rgba(245,197,0,.12); padding: 18px 20px; display: flex; flex-direction: column; gap: 12px; animation: fadeUp .45s cubic-bezier(.16,1,.3,1) both; position: relative; overflow: hidden; }
.resumen-card::before { content:''; position:absolute; left:0; top:0; bottom:0; width:3px; }
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
.leyenda { display: flex; flex-direction: column; gap: 6px; padding-top: 4px; border-top: 1px solid rgba(245,197,0,.07); }
.ley-item { display: flex; align-items: center; gap: 8px; font-size: .7rem; color: #666; }
.ley-dot { width: 11px; height: 11px; flex-shrink: 0; border-radius: 2px; }
.dot-activo      { background: rgba(74,222,128,.25); border: 1px solid rgba(74,222,128,.6); }
.dot-vencido-cal { background: rgba(102,28,28,0.25); border: 1px solid rgba(100,100,100,.5); }
.dot-ninguno     { background: rgba(255,255,255,.03); border: 1px solid rgba(255,255,255,.06); }

/* ── Calendario ── */
.cal-card { background: #111; border: 1px solid rgba(245,197,0,.12); padding: 20px 22px; display: flex; flex-direction: column; gap: 14px; animation: fadeUp .45s .06s cubic-bezier(.16,1,.3,1) both; }
.cal-vencidos { border-color: rgba(100,100,100,.2); }
.cal-nav { display: flex; align-items: center; justify-content: space-between; }
.cal-mes { font-family: 'Barlow Condensed',sans-serif; font-size: .95rem; font-weight: 700; letter-spacing: .12em; text-transform: uppercase; color: #f0f0f0; }
.cal-arrow { background: transparent; border: 1px solid rgba(245,197,0,.2); color: #f5c500; width: 28px; height: 28px; display: flex; align-items: center; justify-content: center; cursor: pointer; font-size: 1.2rem; line-height: 1; transition: all .2s; }
.cal-arrow:hover { background: rgba(245,197,0,.1); border-color: #f5c500; }
.cal-header { display: grid; grid-template-columns: repeat(7,1fr); text-align: center; gap: 3px; }
.cal-header span { font-family: 'Barlow Condensed',sans-serif; font-size: .6rem; letter-spacing: .1em; text-transform: uppercase; color: #f5c500; padding: 4px 0; }
.cal-vencidos .cal-header span { color: #555; }
.cal-grid { display: grid; grid-template-columns: repeat(7,1fr); gap: 3px; }
.cal-empty { aspect-ratio: 1; }
.cal-dia { aspect-ratio: 1; display: flex; align-items: center; justify-content: center; font-family: 'Barlow Condensed',sans-serif; font-size: .8rem; font-weight: 700; cursor: default; border-radius: 2px; transition: all .15s; }
.dia-ninguno { color: #2a2a2a; background: rgba(255,255,255,.02); border: 1px solid rgba(255,255,255,.03); }
.dia-activo  { color: #d1fae5; background: rgba(74,222,128,.18); border: 1px solid rgba(74,222,128,.35); }
.dia-activo:hover { background: rgba(74,222,128,.26); }
.dia-vencido { color: #fffbfb; background: rgba(148,21,21,0.12); border: 1px solid rgba(100,100,100,.22); }
.dia-hoy     { background: #f5c500 !important; color: #0a0a0a !important; border: none !important; font-weight: 900; box-shadow: 0 0 12px rgba(245,197,0,.4); }

/* ── Labels de sección ── */
.section-label { display: flex; align-items: center; gap: 8px; font-family: 'Barlow Condensed',sans-serif; font-size: .63rem; letter-spacing: .28em; text-transform: uppercase; color: #f5c500; font-weight: 700; flex-shrink: 0; }
.section-label svg { width: 14px; height: 14px; }

/* ── Filtros ── */
.filtro-bar { display: flex; flex-direction: column; gap: 10px; flex-shrink: 0; background: #111; border: 1px solid rgba(245,197,0,.12); padding: 14px 16px; }
.filtro-fields { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; }
.filtro-field { display: flex; flex-direction: column; gap: 6px; }
.filtro-label { font-family: 'Barlow Condensed',sans-serif; font-size: .58rem; letter-spacing: .25em; text-transform: uppercase; color: #f5c500; font-weight: 700; opacity: .8; }
.filtro-label-venc { color: #666 !important; }
.filtro-date { background: #0d0d0d; border: 1px solid rgba(245,197,0,.2); color: #e0e0e0; padding: 8px 10px; font-family: 'Barlow',sans-serif; font-size: .8rem; outline: none; transition: border-color .2s; color-scheme: dark; width: 100%; }
.filtro-date:focus { border-color: rgba(245,197,0,.6); }
.filtro-date::-webkit-calendar-picker-indicator { filter: invert(1) sepia(1) saturate(3) hue-rotate(10deg); cursor: pointer; opacity: .5; }
.filtro-date-venc { border-color: rgba(100,100,100,.2); color: #666; }
.filtro-date-venc:focus { border-color: rgba(100,100,100,.5); }
.filtro-date-venc::-webkit-calendar-picker-indicator { filter: none; opacity: .3; }
.filtro-meta { display: flex; align-items: center; justify-content: flex-end; gap: 12px; }
.filtro-count-total { font-family: 'Barlow Condensed',sans-serif; font-size: .65rem; letter-spacing: .15em; text-transform: uppercase; color: #555; }
.filtro-reset { background: none; border: 1px solid rgba(245,197,0,.2); color: #f5c500; font-family: 'Barlow Condensed',sans-serif; font-size: .62rem; letter-spacing: .18em; text-transform: uppercase; padding: 5px 12px; cursor: pointer; transition: all .2s; }
.filtro-reset:hover { background: rgba(245,197,0,.08); border-color: rgba(245,197,0,.5); }
.filtro-footer-hint { color: #3a3a3a; font-style: italic; }

/* ── Tablas ── */
.table-wrapper { display: flex; flex-direction: column; border: 1px solid rgba(245,197,0,.12); background: #111; animation: fadeUp .45s .1s cubic-bezier(.16,1,.3,1) both; flex-shrink: 0; min-height: 80px; max-height: 260px; }
@keyframes fadeUp { from{opacity:0;transform:translateY(10px)}to{opacity:1;transform:translateY(0)} }
.table-scroll { overflow-y: auto; flex: 1; }
table { width: 100%; border-collapse: collapse; font-size: .85rem; }
thead { background: rgba(245,197,0,.06); border-bottom: 1px solid rgba(245,197,0,.15); position: sticky; top: 0; z-index: 2; }
th { font-family: 'Barlow Condensed',sans-serif; font-size: .63rem; font-weight: 700; letter-spacing: .22em; text-transform: uppercase; color: #f5c500; padding: 13px 16px; text-align: left; white-space: nowrap; }
td { padding: 12px 16px; border-bottom: 1px solid rgba(255,255,255,.04); color: #ccc; vertical-align: middle; white-space: nowrap; }
.table-row { transition: background .15s; animation: rowIn .3s ease both; }
.table-row:hover { background: rgba(245,197,0,.04); }
.table-row:last-child td { border-bottom: none; }
@keyframes rowIn { from{opacity:0;transform:translateX(-4px)}to{opacity:1;transform:translateX(0)} }
.id-badge { display: inline-flex; align-items: center; justify-content: center; width: 26px; height: 26px; background: rgba(245,197,0,.07); border: 1px solid rgba(245,197,0,.2); font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .78rem; color: #f5c500; }
.tipo-label { font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .88rem; color: #f0f0f0; letter-spacing: .04em; text-transform: uppercase; }
.precio-val { font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .88rem; color: #f5c500; }
.estado-badge { display: inline-block; font-family: 'Barlow Condensed',sans-serif; font-size: .66rem; font-weight: 700; letter-spacing: .18em; text-transform: uppercase; padding: 4px 9px; border: 1px solid; }
.badge-pagado    { color: #4ade80; border-color: rgba(74,222,128,.4); background: rgba(74,222,128,.07); }
.badge-pendiente { color: #f59e0b; border-color: rgba(245,158,11,.4); background: rgba(245,158,11,.07); }
.badge-default   { color: #888; border-color: rgba(136,136,136,.3); background: rgba(136,136,136,.06); }
.disp-badge { display: inline-block; font-family: 'Barlow Condensed',sans-serif; font-size: .66rem; font-weight: 700; letter-spacing: .14em; text-transform: uppercase; padding: 4px 9px; border: 1px solid; }
.disp-activo     { color: #4ade80; border-color: rgba(74,222,128,.35); background: rgba(74,222,128,.06); }
.disp-vencer-hoy { color: #f59e0b; border-color: rgba(245,158,11,.35); background: rgba(245,158,11,.06); }
.disp-proximo    { color: #60a5fa; border-color: rgba(96,165,250,.35); background: rgba(96,165,250,.06); }
.disp-vencido    { color: #555; border-color: rgba(85,85,85,.3); background: rgba(85,85,85,.05); }
.table-footer { padding: 10px 16px; border-top: 1px solid rgba(245,197,0,.08); font-size: .68rem; color: #444; letter-spacing: .06em; flex-shrink: 0; }
.empty-row { text-align: center; color: #444; padding: 40px 20px !important; font-size: .88rem; }

/* ── Vencidos ── */
.vencidos-header { display: flex; align-items: center; justify-content: space-between; gap: 12px; flex-shrink: 0; }
.venc-count { background: rgba(204,12,12,0.15); border: 1px solid rgba(100,100,100,.25); color: #555; font-size: .6rem; font-weight: 700; letter-spacing: .1em; padding: 2px 8px; margin-left: 4px; border-radius: 2px; }
.venc-vacio { font-size: .72rem; color: #333; letter-spacing: .06em; }
.btn-ver-vencidos { display: inline-flex; align-items: center; gap: 7px; background: transparent; border: 1px solid rgba(100,100,100,.25); color: #555; font-family: 'Barlow Condensed',sans-serif; font-size: .72rem; font-weight: 700; letter-spacing: .18em; text-transform: uppercase; padding: 8px 16px; cursor: pointer; transition: all .2s; }
.btn-ver-vencidos svg { width: 13px; height: 13px; flex-shrink: 0; }
.btn-ver-vencidos:hover { border-color: rgba(100,100,100,.5); color: #888; background: rgba(100,100,100,.06); }
.vencidos-panel { display: flex; flex-direction: column; gap: 16px; border: 1px solid rgba(90,88,88,0.15); background: rgba(100,100,100,.03); padding: 20px; }
.venc-top-row { display: grid; grid-template-columns: 1fr 1.2fr; gap: 18px; }
.venc-info-card { background: #0f0f0f; border: 1px solid rgba(100,100,100,.15); padding: 18px 20px; }
.venc-resumenes { display: flex; flex-direction: column; gap: 10px; }
.venc-item { padding: 10px 0; border-bottom: 1px solid rgba(100,100,100,.1); }
.venc-item:last-child { border-bottom: none; }
.venc-tipo { font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .8rem; letter-spacing: .1em; text-transform: uppercase; color: #555; margin-bottom: 4px; }
.venc-fechas-row { display: flex; align-items: center; gap: 8px; font-size: .8rem; color: #444; }
.venc-arrow { color: #333; }
.venc-precio { font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .8rem; color: #444; margin-top: 4px; }
.table-vencidos { border-color: rgba(100,100,100,.15); max-height: 200px; }
.table-vencidos thead { background: rgba(9,51,128,0.06); border-bottom-color: rgba(100,100,100,.15); }
.table-vencidos th { color: #ffffff; }
.row-vencido { opacity: .7; }
.row-vencido:hover { opacity: 1; background: rgba(100,100,100,.05) !important; }
.id-badge-venc { background: rgba(100,100,100,.08) !important; border-color: rgba(100,100,100,.2) !important; color: #555 !important; }
.tipo-venc { color: #555 !important; }
.precio-venc { color: #555 !important; }

/* ── Transiciones ── */
.expand-enter-active { transition: opacity .3s ease, transform .3s ease; }
.expand-leave-active { transition: opacity .2s ease, transform .2s ease; }
.expand-enter-from, .expand-leave-to { opacity: 0; transform: translateY(-8px); }
.modal-fade-enter-active, .modal-fade-leave-active { transition: opacity .25s; }
.modal-fade-enter-from, .modal-fade-leave-to { opacity: 0; }

/* ═══════════════════════════════════════
   MODAL DE PAGO
═══════════════════════════════════════ */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,.82); display: flex; align-items: center; justify-content: center; z-index: 1000; padding: 20px; }
.modal-pay { background: #0f0f0f; border: 1px solid rgba(245,197,0,.2); box-shadow: 0 40px 100px rgba(0,0,0,.9), 0 0 80px rgba(245,197,0,.04); width: 100%; max-width: 780px; max-height: 90vh; overflow: hidden; display: flex; flex-direction: column; animation: modalIn .32s cubic-bezier(.16,1,.3,1) both; }
@keyframes modalIn { from{opacity:0;transform:scale(.96) translateY(14px)}to{opacity:1;transform:scale(1) translateY(0)} }
.mpay-header { display: flex; align-items: flex-start; justify-content: space-between; padding: 22px 28px 18px; border-bottom: 1px solid rgba(245,197,0,.1); background: linear-gradient(to right, rgba(245,197,0,.04), transparent); }
.mpay-eyebrow { display: block; font-family: 'Barlow Condensed',sans-serif; font-size: .58rem; letter-spacing: .35em; text-transform: uppercase; color: #f5c500; opacity: .6; margin-bottom: 5px; }
.mpay-title { font-family: 'Barlow Condensed',sans-serif; font-size: 1.8rem; font-weight: 900; text-transform: uppercase; color: #fff; line-height: 1; letter-spacing: .04em; }
.mpay-sub { font-size: .75rem; color: #555; margin-top: 5px; }
.mpay-close { background: none; border: 1px solid #222; color: #444; width: 30px; height: 30px; display: flex; align-items: center; justify-content: center; cursor: pointer; transition: all .2s; font-size: .85rem; flex-shrink: 0; }
.mpay-close:hover { border-color: rgba(245,197,0,.4); color: #f5c500; }
.mpay-body { display: grid; grid-template-columns: 1fr 1fr; flex: 1; overflow: hidden; min-height: 0; }
.mpay-left { padding: 22px 24px; border-right: 1px solid rgba(245,197,0,.08); overflow-y: auto; display: flex; flex-direction: column; gap: 18px; }
.mpay-section-label { font-family: 'Barlow Condensed',sans-serif; font-size: .6rem; letter-spacing: .3em; text-transform: uppercase; color: #f5c500; font-weight: 700; opacity: .7; }
.mpay-planes { display: flex; flex-direction: column; gap: 8px; }
.mpay-plan-card { display: flex; align-items: center; gap: 14px; background: #111; border: 1px solid rgba(255,255,255,.06); padding: 14px 16px; cursor: pointer; transition: all .2s; text-align: left; width: 100%; position: relative; }
.mpay-plan-card:hover { border-color: rgba(245,197,0,.3); background: rgba(245,197,0,.04); }
.mpay-plan-activo { border-color: #f5c500 !important; background: rgba(245,197,0,.08) !important; box-shadow: inset 3px 0 0 #f5c500; }
.mpay-plan-icon { font-size: 1.4rem; width: 36px; text-align: center; flex-shrink: 0; }
.mpay-plan-info { display: flex; flex-direction: column; gap: 2px; flex: 1; }
.mpay-plan-nombre { font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: 1rem; letter-spacing: .06em; text-transform: uppercase; color: #e0e0e0; }
.mpay-plan-precio { font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .85rem; color: #f5c500; }
.mpay-plan-check { font-size: 1rem; color: #f5c500; font-weight: 900; margin-left: auto; }
.mpay-cantidad-box { background: #111; border: 1px solid rgba(245,197,0,.1); padding: 16px 18px; }
.mpay-cantidad-ctrl { display: flex; align-items: center; gap: 0; border: 1px solid rgba(245,197,0,.2); width: fit-content; }
.mpay-qty-btn { width: 36px; height: 36px; background: rgba(245,197,0,.08); border: none; color: #f5c500; font-size: 1.2rem; font-weight: 700; cursor: pointer; transition: all .18s; display: flex; align-items: center; justify-content: center; }
.mpay-qty-btn:hover { background: rgba(245,197,0,.18); }
.mpay-qty-val { width: 52px; height: 36px; display: flex; align-items: center; justify-content: center; font-family: 'Barlow Condensed',sans-serif; font-weight: 900; font-size: 1.2rem; color: #fff; background: #0d0d0d; border-left: 1px solid rgba(245,197,0,.15); border-right: 1px solid rgba(245,197,0,.15); }
.mpay-cantidad-hint { font-size: .72rem; color: #555; margin-top: 10px; letter-spacing: .04em; }
.mpay-fecha-box { background: #111; border: 1px solid rgba(245,197,0,.1); padding: 16px 18px; display: flex; flex-direction: column; gap: 12px; }
.mpay-date-input { width: 100%; background: #0d0d0d; border: 1px solid rgba(245,197,0,.25); color: #f0f0f0; padding: 10px 14px; font-family: 'Barlow',sans-serif; font-size: .85rem; cursor: pointer; outline: none; transition: border-color .2s; color-scheme: dark; }
.mpay-date-input:focus { border-color: #f5c500; }
.mpay-date-input::-webkit-calendar-picker-indicator { filter: invert(1) sepia(1) saturate(5) hue-rotate(10deg); cursor: pointer; }
.mpay-fecha-fin { background: rgba(245,197,0,.04); border: 1px solid rgba(245,197,0,.12); padding: 10px 14px; display: flex; flex-direction: column; gap: 8px; }
.mpay-fecha-row { display: flex; align-items: center; gap: 12px; }
.mpay-fecha-item { display: flex; flex-direction: column; gap: 2px; }
.mpay-fecha-key { font-family: 'Barlow Condensed',sans-serif; font-size: .55rem; letter-spacing: .2em; text-transform: uppercase; color: #666; font-weight: 700; }
.mpay-fecha-val { font-family: 'Barlow Condensed',sans-serif; font-size: .9rem; font-weight: 700; color: #f5c500; }
.mpay-fecha-arrow { color: #f5c500; font-weight: 700; font-size: 1rem; }
.mpay-fecha-desc { font-size: .68rem; color: #555; letter-spacing: .04em; }
.mpay-right { padding: 22px 24px; overflow-y: auto; display: flex; flex-direction: column; gap: 18px; background: rgba(245,197,0,.015); }
.mpay-right-vacio { align-items: center; justify-content: center; }
.mpay-vacio { display: flex; flex-direction: column; align-items: center; gap: 12px; text-align: center; color: #333; }
.mpay-vacio-icon { font-size: 2.5rem; filter: grayscale(1); opacity: .3; }
.mpay-vacio p { font-size: .82rem; line-height: 1.6; }
.mpay-resumen { background: #111; border: 1px solid rgba(245,197,0,.12); display: flex; flex-direction: column; gap: 0; overflow: hidden; }
.mpay-resumen-plan { display: flex; align-items: center; gap: 14px; padding: 16px 18px; border-bottom: 1px solid rgba(245,197,0,.08); }
.mpay-resumen-icon { font-size: 1.8rem; }
.mpay-resumen-tipo { display: block; font-family: 'Barlow Condensed',sans-serif; font-weight: 900; font-size: 1.1rem; letter-spacing: .06em; text-transform: uppercase; color: #f0f0f0; }
.mpay-resumen-unit { font-size: .7rem; color: #555; letter-spacing: .06em; }
.mpay-lineas { display: flex; flex-direction: column; gap: 0; padding: 12px 18px; }
.mpay-linea { display: flex; justify-content: space-between; align-items: center; padding: 8px 0; font-size: .82rem; color: #777; border-bottom: 1px solid rgba(255,255,255,.03); }
.mpay-linea:last-child { border-bottom: none; }
.mpay-divider { border-top: 1px solid rgba(245,197,0,.15); margin: 4px 0; }
.mpay-linea-total { color: #fff; font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .9rem; letter-spacing: .1em; padding-top: 10px; }
.mpay-total-val { font-size: 1.3rem; font-weight: 900; color: #f5c500; }
.mpay-desglose { padding: 12px 18px 16px; border-top: 1px solid rgba(245,197,0,.06); }
.mpay-desglose-label { font-family: 'Barlow Condensed',sans-serif; font-size: .58rem; letter-spacing: .25em; text-transform: uppercase; color: #444; font-weight: 700; margin-bottom: 8px; }
.mpay-desglose-chips { display: flex; flex-direction: column; gap: 4px; max-height: 120px; overflow-y: auto; padding-right: 4px; }
.mpay-desglose-chips::-webkit-scrollbar { width: 4px; }
.mpay-desglose-chips::-webkit-scrollbar-track { background: transparent; }
.mpay-desglose-chips::-webkit-scrollbar-thumb { background: rgba(245,197,0,.25); border-radius: 2px; }
.mpay-desglose-chips::-webkit-scrollbar-thumb:hover { background: rgba(245,197,0,.45); }
.mpay-chip { font-size: .72rem; color: #666; padding: 4px 10px; background: rgba(255,255,255,.03); border: 1px solid rgba(255,255,255,.05); border-radius: 1px; font-family: 'Barlow Condensed',sans-serif; letter-spacing: .04em; }
.mpay-alerta { display: flex; align-items: center; gap: 8px; background: rgba(245,158,11,.06); border: 1px solid rgba(245,158,11,.25); padding: 10px 14px; font-size: .75rem; color: #f59e0b; letter-spacing: .03em; }
.mpay-btn-pagar { display: flex; align-items: center; gap: 14px; background: linear-gradient(135deg, #f5c500 0%, #ff9500 100%); border: none; padding: 16px 22px; cursor: pointer; transition: all .25s; width: 100%; text-align: left; position: relative; overflow: hidden; }
.mpay-btn-pagar::before { content: ''; position: absolute; inset: 0; background: linear-gradient(135deg, rgba(255,255,255,.15) 0%, transparent 60%); pointer-events: none; }
.mpay-btn-pagar:hover:not(:disabled) { filter: brightness(1.08); box-shadow: 0 8px 30px rgba(245,197,0,.35); transform: translateY(-1px); }
.mpay-btn-pagar:active:not(:disabled) { transform: translateY(0); filter: brightness(.95); }
.mpay-btn-pagar:disabled { opacity: 0.45; cursor: not-allowed; transform: none !important; filter: none !important; box-shadow: none !important; }
.mpay-pagar-bolt { font-size: 1.4rem; flex-shrink: 0; }
.mpay-pagar-text { display: flex; flex-direction: column; gap: 1px; flex: 1; }
.mpay-pagar-titulo { font-family: 'Barlow Condensed',sans-serif; font-weight: 900; font-size: 1rem; letter-spacing: .15em; text-transform: uppercase; color: #0a0a0a; }
.mpay-pagar-monto { font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .82rem; color: rgba(0,0,0,.6); letter-spacing: .06em; }
.mpay-pagar-arrow { width: 20px; height: 20px; stroke: rgba(0,0,0,.5); flex-shrink: 0; }
#wallet_container { width: 100%; min-height: 48px; }
.mpay-nota { font-size: .68rem; color: #444; line-height: 1.5; letter-spacing: .03em; padding: 0 2px; }

/* ── Modal password ── */
.pass-field { display: flex; flex-direction: column; gap: 6px; }
.pass-label { font-family: 'Barlow Condensed',sans-serif; font-size: .58rem; letter-spacing: .28em; text-transform: uppercase; color: #f5c500; font-weight: 700; opacity: .7; }
.pass-match { font-size: .75rem; padding: 8px 12px; border: 1px solid; letter-spacing: .03em; }
.match-ok   { color: #4ade80; border-color: rgba(74,222,128,.35); background: rgba(74,222,128,.06); }
.match-fail { color: #f87171; border-color: rgba(248,113,113,.35); background: rgba(248,113,113,.06); }
.pass-requisitos { display: flex; flex-direction: column; gap: 5px; padding: 10px 12px; background: rgba(255,255,255,.02); border: 1px solid rgba(255,255,255,.06); }
.req-item { font-size: .72rem; letter-spacing: .03em; transition: color .2s; }
.req-ok   { color: #4ade80; }
.req-fail { color: #555; }

/* ═══════════════════════════════════════
   ACCIONES MÓVIL (ocultas en desktop)
═══════════════════════════════════════ */
.mobile-actions { display: none; }

/* ═══════════════════════════════════════
   RESPONSIVE
═══════════════════════════════════════ */

/* Tablet: sidebar más angosto */
@media (max-width: 1100px) {
  .sidebar { width: 200px; min-width: 200px; }
  .main { padding: 36px 32px; }
  .top-row, .venc-top-row { grid-template-columns: 1fr; }
}

/* Mobile: sidebar oculto, aparecen botones móviles */
@media (max-width: 860px) {
  .sidebar { display: none; }

  .main { padding: 28px 20px 40px; gap: 16px; }

  .mobile-actions {
    display: flex;
    gap: 8px;
    align-items: stretch;
    flex-shrink: 0;
  }

  .mpay-body { grid-template-columns: 1fr; overflow-y: auto; }
  .mpay-left { border-right: none; border-bottom: 1px solid rgba(245,197,0,.08); }

  .filtro-fields { grid-template-columns: 1fr; }

  .mob-btn-pass,
  .mob-btn-sus {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 12px 14px;
    border: 1px solid;
    background: transparent;
    font-family: 'Barlow Condensed', sans-serif;
    font-weight: 700;
    font-size: .78rem;
    letter-spacing: .1em;
    text-transform: uppercase;
    cursor: pointer;
    transition: all .2s;
    flex: 1;
  }
  .mob-btn-pass {
    color: #aaa;
    border-color: rgba(255,255,255,.1);
    background: rgba(255,255,255,.03);
  }
  .mob-btn-pass:hover {
    border-color: rgba(245,197,0,.35);
    color: #f5c500;
    background: rgba(245,197,0,.05);
  }
  .mob-btn-sus {
    color: #0a0a0a;
    border-color: #f5c500;
    background: linear-gradient(135deg, #f5c500, #ff9500);
    position: relative;
    overflow: hidden;
  }
  .mob-btn-sus::before {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(135deg, rgba(255,255,255,.15) 0%, transparent 60%);
    pointer-events: none;
  }
  .mob-btn-sus:hover { filter: brightness(1.08); box-shadow: 0 4px 20px rgba(245,197,0,.35); }
  .mob-btn-icon { font-size: 1rem; flex-shrink: 0; }
  .mob-btn-text { white-space: nowrap; }
  .mob-btn-logout {
    width: 44px;
    height: auto;
    flex-shrink: 0;
    background: transparent;
    border: 1px solid rgba(255,255,255,.08);
    color: #444;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all .2s;
  }
  .mob-btn-logout svg { width: 16px; height: 16px; }
  .mob-btn-logout:hover { border-color: rgba(245,197,0,.4); color: #f5c500; }
}
</style>