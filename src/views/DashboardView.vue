<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const users     = ref([])
const customers = ref([])
const vista     = ref('usuarios')
const loading   = ref(true)
const errorMsg  = ref('')
const busqueda  = ref('')
const filtroRol = ref('todos')

const username = localStorage.getItem('username') || 'Admin'

// ── Sidebar móvil ─────────────────────────────────────────────────
const sidebarOpen = ref(false)
const toggleSidebar = () => { sidebarOpen.value = !sidebarOpen.value }
const closeSidebar  = () => { sidebarOpen.value = false }

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

const getToken = () => {
  const token = localStorage.getItem('token')
  if (!token || tokenExpirado(token)) { router.push('/login'); return null }
  return token
}

// ── Fetch ─────────────────────────────────────────────────────────
const obtenerUsuarios = async () => {
  loading.value = true; errorMsg.value = ''
  try {
    const token = getToken(); if (!token) return
    const res = await fetch('/api/users', { headers: { Authorization: `Bearer ${token}` } })
    if (res.status === 401) { localStorage.removeItem('token'); router.push('/login'); return }
    if (!res.ok) throw new Error('Error al cargar usuarios')
    users.value = await res.json()
  } catch (err) { errorMsg.value = err.message }
  finally { loading.value = false }
}

const obtenerClientes = async () => {
  loading.value = true; errorMsg.value = ''
  try {
    const token = getToken(); if (!token) return
    const res = await fetch('/api/customers', { headers: { Authorization: `Bearer ${token}` } })
    if (!res.ok) throw new Error('Error al cargar clientes')
    const data = await res.json()
    customers.value = Array.isArray(data) ? data : []
  } catch (err) {
    errorMsg.value = err.message
    customers.value = []
  }
  finally { loading.value = false }
}

const cambiarVista = (tipo) => {
  vista.value = tipo; busqueda.value = ''
  closeSidebar()
  if (tipo === 'usuarios') obtenerUsuarios()
  else obtenerClientes()
}

const verPlanes = (cliente) => {
  router.push({
    name: 'dashboard-plan',
    params: { idcustomer: cliente.idcustomer },
    state: { nombre: `${cliente.name} ${cliente.lastname}` }
  })
}

// ── CRUD Usuarios ─────────────────────────────────────────────────
const showModal      = ref(false)
const showEditModal  = ref(false)
const modalError     = ref('')
const guardandoUser  = ref(false)

const newUser  = ref({ username: '', password: '', confirmPassword: '' })
const editUser = ref({ id: null, username: '', password: '', confirmPassword: '' })

const resetNewUser  = () => { newUser.value  = { username: '', password: '', confirmPassword: '' }; modalError.value = '' }
const resetEditUser = () => { editUser.value = { id: null, username: '', password: '', confirmPassword: '' }; modalError.value = '' }

const abrirEditar = (user) => {
  editUser.value = { id: user.id, username: user.username, password: '', confirmPassword: '' }
  modalError.value = ''
  showEditModal.value = true
}

const crearAdmin = async () => {
  modalError.value = ''
  const { username: u, password: p, confirmPassword: c } = newUser.value
  if (!u || !p || !c) { modalError.value = 'Todos los campos son obligatorios'; return }
  if (p !== c)        { modalError.value = 'Las contraseñas no coinciden'; return }
  guardandoUser.value = true
  try {
    const token = getToken(); if (!token) return
    const res = await fetch('/api/new-users', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify({ username: u, password: p })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Error al crear admin')
    showModal.value = false; resetNewUser(); obtenerUsuarios()
  } catch (err) { modalError.value = err.message }
  finally { guardandoUser.value = false }
}

const actualizarUsuario = async () => {
  modalError.value = ''
  const { id: uid, username: u, password: p, confirmPassword: c } = editUser.value
  if (!u) { modalError.value = 'El usuario es obligatorio'; return }
  if (p && p !== c) { modalError.value = 'Las contraseñas no coinciden'; return }
  guardandoUser.value = true
  try {
    const token = getToken(); if (!token) return
    const res = await fetch(`/api/update-users/${uid}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify({ username: u, password: p })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Error al actualizar')
    showEditModal.value = false; resetEditUser(); obtenerUsuarios()
  } catch (err) { modalError.value = err.message }
  finally { guardandoUser.value = false }
}

const eliminarUsuario = async (id) => {
  if (!confirm('¿Eliminar este usuario?')) return
  try {
    const token = getToken(); if (!token) return
    const res = await fetch(`/api/delete-users/${id}`, {
      method: 'DELETE', headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Error al eliminar')
    obtenerUsuarios()
  } catch (err) { alert(err.message) }
}

// ── CRUD Clientes ─────────────────────────────────────────────────
const showCreateCustomerModal = ref(false)
const showEditCustomerModal   = ref(false)
const modalCustomerError      = ref('')
const guardandoCustomer       = ref(false)

const newCustomer  = ref({ cedula: '', name: '', lastname: '', phone: '', password: '', confirmPassword: '' })
const editCustomer = ref({ idcustomer: null, cedula: '', name: '', lastname: '', phone: '', password: '', confirmPassword: '' })

const resetNewCustomer  = () => { newCustomer.value  = { cedula: '', name: '', lastname: '', phone: '', password: '', confirmPassword: '' }; modalCustomerError.value = '' }
const resetEditCustomer = () => { editCustomer.value = { idcustomer: null, cedula: '', name: '', lastname: '', phone: '', password: '', confirmPassword: '' }; modalCustomerError.value = '' }

const abrirEditarCustomer = (c) => {
  editCustomer.value = { idcustomer: c.idcustomer, cedula: c.cedula, name: c.name, lastname: c.lastname, phone: c.phone, password: '', confirmPassword: '' }
  modalCustomerError.value = ''
  showEditCustomerModal.value = true
}

const crearCustomer = async () => {
  modalCustomerError.value = ''
  const nc = newCustomer.value
  if (!nc.cedula || !nc.name || !nc.lastname || !nc.phone || !nc.password || !nc.confirmPassword) {
    modalCustomerError.value = 'Todos los campos son obligatorios'; return
  }
  if (!/^\d{10}$/.test(nc.phone)) {
    modalCustomerError.value = 'El teléfono debe tener exactamente 10 dígitos'; return
  }
  if (nc.password !== nc.confirmPassword) {
    modalCustomerError.value = 'Las contraseñas no coinciden'; return
  }
  guardandoCustomer.value = true
  try {
    const token = getToken(); if (!token) return
    const res = await fetch('/api/new-customers', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify({ cedula: nc.cedula, name: nc.name, lastname: nc.lastname, phone: nc.phone, password: nc.password })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Error al crear cliente')
    showCreateCustomerModal.value = false; resetNewCustomer(); obtenerClientes()
  } catch (err) { modalCustomerError.value = err.message }
  finally { guardandoCustomer.value = false }
}

const actualizarCustomer = async () => {
  modalCustomerError.value = ''
  const ec = editCustomer.value
  if (!ec.cedula || !ec.name) { modalCustomerError.value = 'Cédula y nombre son obligatorios'; return }
  if (ec.phone && !/^\d{10}$/.test(ec.phone)) { modalCustomerError.value = 'El teléfono debe tener exactamente 10 dígitos'; return }
  if (ec.password && ec.password !== ec.confirmPassword) { modalCustomerError.value = 'Las contraseñas no coinciden'; return }
  guardandoCustomer.value = true
  try {
    const token = getToken(); if (!token) return
    const res = await fetch(`/api/update-customers/${ec.idcustomer}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify({ cedula: ec.cedula, name: ec.name, lastname: ec.lastname, phone: ec.phone, password: ec.password })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Error al actualizar')
    showEditCustomerModal.value = false; resetEditCustomer(); obtenerClientes()
  } catch (err) { modalCustomerError.value = err.message }
  finally { guardandoCustomer.value = false }
}

const eliminarCustomer = async (id) => {
  if (!confirm('¿Eliminar este cliente?')) return
  try {
    const token = getToken(); if (!token) return
    const res = await fetch(`/api/delete-customers/${id}`, {
      method: 'DELETE', headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Error al eliminar')
    obtenerClientes()
  } catch (err) { alert(err.message) }
}

// ── Price Plans ───────────────────────────────────────────────────
const PLAN_META = {
  dia:    { label: 'Día',    icon: '☀', color: '#ff7a00', bg: 'rgba(255,122,0,.08)',   border: 'rgba(255,122,0,.35)'   },
  semana: { label: 'Semana', icon: '📅', color: '#f5c500', bg: 'rgba(245,197,0,.08)',   border: 'rgba(245,197,0,.35)'   },
  mes:    { label: 'Mes',    icon: '🗓', color: '#4db8ff', bg: 'rgba(77,184,255,.08)',  border: 'rgba(77,184,255,.35)'  },
  año:    { label: 'Año',    icon: '⚡', color: '#a78bfa', bg: 'rgba(167,139,250,.08)', border: 'rgba(167,139,250,.35)' },
}
const PLAN_SORT = ['dia', 'semana', 'mes', 'año']

const normalizePlanKey = (typeplan) => {
  const aliases = { day: 'dia', week: 'semana', month: 'mes', year: 'año', anio: 'año' }
  const k = typeplan?.toLowerCase().trim()
  return aliases[k] || k
}

const getMeta = (typeplan) => {
  const k = normalizePlanKey(typeplan)
  return PLAN_META[k] || {
    label: typeplan?.charAt(0).toUpperCase() + typeplan?.slice(1) || '?',
    icon: '💠', color: '#888', bg: 'rgba(136,136,136,.06)', border: 'rgba(136,136,136,.25)'
  }
}

const showPricePlanModal = ref(false)
const pricePlans         = ref([])
const loadingPrices      = ref(false)
const priceError         = ref('')
const editingPriceId     = ref(null)
const editingPrice       = ref('')
const savingPrice        = ref(false)
const savePriceError     = ref('')

const sortedPricePlans = computed(() =>
  [...pricePlans.value].sort((a, b) => {
    const ai = PLAN_SORT.indexOf(normalizePlanKey(a.typeplan))
    const bi = PLAN_SORT.indexOf(normalizePlanKey(b.typeplan))
    return (ai === -1 ? 99 : ai) - (bi === -1 ? 99 : bi)
  })
)

const obtenerPricePlans = async () => {
  loadingPrices.value = true; priceError.value = ''
  try {
    const token = getToken(); if (!token) return
    const res = await fetch('/api/priceplan', { headers: { Authorization: `Bearer ${token}` } })
    if (!res.ok) throw new Error('Error al cargar precios')
    pricePlans.value = await res.json()
  } catch (err) { priceError.value = err.message }
  finally { loadingPrices.value = false }
}

const abrirPricePlanModal = () => {
  showPricePlanModal.value = true
  editingPriceId.value = null
  editingPrice.value   = ''
  savePriceError.value = ''
  obtenerPricePlans()
  closeSidebar()
}

const iniciarEditarPrecio = (plan) => {
  editingPriceId.value = plan.idpriceplan
  editingPrice.value   = String(plan.price)
  savePriceError.value = ''
}

const cancelarEditarPrecio = () => {
  editingPriceId.value = null
  editingPrice.value   = ''
  savePriceError.value = ''
}

const guardarPrecio = async (plan) => {
  savePriceError.value = ''
  const val = parseFloat(editingPrice.value)
  if (isNaN(val) || val <= 0) { savePriceError.value = 'El precio debe ser mayor a 0'; return }
  savingPrice.value = true
  try {
    const token = getToken(); if (!token) return
    const res = await fetch(`/api/update-priceplan/${plan.idpriceplan}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify({ price: val })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Error al actualizar')
    plan.price = val
    cancelarEditarPrecio()
  } catch (err) { savePriceError.value = err.message }
  finally { savingPrice.value = false }
}

const formatPrice = (price) =>
  new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', maximumFractionDigits: 0 }).format(price)

// ── Filtrado ──────────────────────────────────────────────────────
const usuariosFiltrados = computed(() =>
  users.value.filter(u => {
    const q = busqueda.value.toLowerCase()
    return (u.username.toLowerCase().includes(q) || u.rol.toLowerCase().includes(q) || String(u.id).includes(q))
      && (filtroRol.value === 'todos' || u.rol === filtroRol.value)
  })
)

const clientesFiltrados = computed(() => {
  const q = busqueda.value.toLowerCase()
  return customers.value.filter(c =>
    c.name?.toLowerCase().includes(q) ||
    c.lastname?.toLowerCase().includes(q) ||
    String(c.cedula).includes(q) ||
    c.phone?.includes(q)
  )
})

const rolesUnicos  = computed(() => [...new Set(users.value.map(u => u.rol))])
const rolColor     = (rol) => rol === 'admin' ? 'badge-admin' : rol === 'customer' ? 'badge-customer' : 'badge-default'
const maskPassword = (pw) => '•'.repeat(Math.min(pw?.length || 8, 12))

onMounted(() => { obtenerUsuarios() })
</script>

<template>
  <div class="screen">

    <!-- Background decorativo -->
    <div class="bg-layer" aria-hidden="true">
      <div class="hex hex-1"></div>
      <div class="hex hex-2"></div>
      <div class="slash slash-1"></div>
      <div class="slash slash-2"></div>
      <div class="noise"></div>
    </div>

    <!-- ── Overlay móvil sidebar ── -->
    <transition name="overlay-fade">
      <div v-if="sidebarOpen" class="sidebar-overlay" @click="closeSidebar"></div>
    </transition>

    <!-- ── Topbar móvil ── -->
    <header class="topbar">
      <button class="topbar-menu" @click="toggleSidebar" :class="{ open: sidebarOpen }" aria-label="Menú">
        <span></span><span></span><span></span>
      </button>
      <div class="topbar-logo">
        <svg viewBox="0 0 160 52" xmlns="http://www.w3.org/2000/svg">
          <polygon points="20,2 12,22 18,22 10,44 32,20 22,20 34,2" fill="#f5c500"/>
          <text x="38" y="22" font-family="'Barlow Condensed',sans-serif" font-weight="900" font-size="20" fill="#f5c500" letter-spacing="-0.5">RAYO</text>
          <text x="38" y="42" font-family="'Barlow Condensed',sans-serif" font-weight="900" font-size="20" fill="#ffffff" letter-spacing="-0.5">BOX</text>
        </svg>
      </div>
      <div class="topbar-actions">
        <div class="topbar-stat">
          <span class="topbar-stat-val">{{ vista === 'usuarios' ? users.length : customers.length }}</span>
          <span class="topbar-stat-lbl">{{ vista === 'usuarios' ? 'Usuarios' : 'Clientes' }}</span>
        </div>
      </div>
    </header>

    <div class="layout">

      <!-- ── Sidebar ── -->
      <aside class="sidebar" :class="{ 'sidebar--open': sidebarOpen }">
        <div class="sidebar-inner">

          <div class="sidebar-logo">
            <svg viewBox="0 0 160 80" xmlns="http://www.w3.org/2000/svg">
              <polygon points="36,4 22,38 32,38 18,72 52,32 38,32 58,4" fill="#f5c500"/>
              <text x="62" y="36" font-family="'Barlow Condensed',sans-serif" font-weight="900" font-size="32" fill="#f5c500" letter-spacing="-1">RAYO</text>
              <text x="62" y="66" font-family="'Barlow Condensed',sans-serif" font-weight="900" font-size="32" fill="#ffffff" letter-spacing="-1">BOX</text>
            </svg>
            <span class="sidebar-sub">CROSS LIFTING</span>
          </div>

          <nav class="sidebar-nav">
            <a class="nav-item" :class="{ active: vista === 'usuarios' }" @click="cambiarVista('usuarios')">
              <svg viewBox="0 0 20 20" fill="none"><path d="M4 6h12M4 10h12M4 14h8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/></svg>
              Usuarios
            </a>
            <a class="nav-item" :class="{ active: vista === 'clientes' }" @click="cambiarVista('clientes')">
              <svg viewBox="0 0 20 20" fill="none"><circle cx="10" cy="8" r="3.5" stroke="currentColor" stroke-width="1.5"/><path d="M3 18c0-3.314 3.134-6 7-6s7 2.686 7 6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/></svg>
              Clientes
            </a>

            <div class="nav-divider"></div>

            <button class="btn-add" v-if="vista === 'usuarios'" @click="showModal = true; resetNewUser(); closeSidebar()">
              <svg viewBox="0 0 20 20" fill="none"><circle cx="10" cy="10" r="7.5" stroke="currentColor" stroke-width="1.4"/><path d="M10 7v6M7 10h6" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/></svg>
              Nuevo Admin
            </button>

            <button class="btn-add" v-if="vista === 'clientes'" @click="showCreateCustomerModal = true; resetNewCustomer(); closeSidebar()">
              <svg viewBox="0 0 20 20" fill="none"><circle cx="10" cy="10" r="7.5" stroke="currentColor" stroke-width="1.4"/><path d="M10 7v6M7 10h6" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/></svg>
              Nuevo Cliente
            </button>

            <button class="btn-add btn-prices" @click="abrirPricePlanModal">
              <svg viewBox="0 0 20 20" fill="none">
                <rect x="3" y="5" width="14" height="11" rx="1" stroke="currentColor" stroke-width="1.4"/>
                <path d="M3 8h14" stroke="currentColor" stroke-width="1.4"/>
                <path d="M7 12h2M12 12h2" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/>
              </svg>
              Ver Precios
            </button>
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

        </div>
      </aside>

      <!-- ── Main ── -->
      <main class="main">

        <!-- Header desktop -->
        <header class="page-header">
          <div class="page-title">
            <span class="page-eyebrow">⚡ Gestión del sistema</span>
            <h1>{{ username }}</h1>
          </div>
          <div class="header-meta">
            <div class="stat-pill">
              <span class="stat-val">{{ vista === 'usuarios' ? users.length : customers.length }}</span>
              <span class="stat-lbl">Total</span>
            </div>
            <div class="stat-pill accent" v-if="vista === 'usuarios'">
              <span class="stat-val">{{ users.filter(u => u.rol === 'admin').length }}</span>
              <span class="stat-lbl">Admins</span>
            </div>
            <div class="stat-pill accent" v-if="vista === 'clientes'">
              <span class="stat-val">{{ customers.length }}</span>
              <span class="stat-lbl">Clientes</span>
            </div>
          </div>
        </header>

        <!-- Tab nav móvil (dentro del main, debajo del topbar) -->
        <div class="mobile-tabs">
          <button class="mobile-tab" :class="{ active: vista === 'usuarios' }" @click="cambiarVista('usuarios')">
            <svg viewBox="0 0 20 20" fill="none"><path d="M4 6h12M4 10h12M4 14h8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/></svg>
            Usuarios
          </button>
          <button class="mobile-tab" :class="{ active: vista === 'clientes' }" @click="cambiarVista('clientes')">
            <svg viewBox="0 0 20 20" fill="none"><circle cx="10" cy="8" r="3.5" stroke="currentColor" stroke-width="1.5"/><path d="M3 18c0-3.314 3.134-6 7-6s7 2.686 7 6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/></svg>
            Clientes
          </button>
        </div>

        <!-- Filtros + acciones -->
        <div class="filters-bar">
          <div class="search-wrap">
            <svg class="search-icon" viewBox="0 0 20 20" fill="none">
              <circle cx="8.5" cy="8.5" r="5" stroke="currentColor" stroke-width="1.4"/>
              <path d="M13 13l3.5 3.5" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/>
            </svg>
            <input v-model="busqueda" type="text"
              :placeholder="vista === 'usuarios' ? 'Buscar usuario, rol o ID...' : 'Buscar nombre, cédula...'"
              class="search-input" />
          </div>

          <!-- Acciones inline móvil -->
          <div class="mobile-actions">
            <button v-if="vista === 'usuarios'" class="fab-inline" @click="showModal = true; resetNewUser()">
              <svg viewBox="0 0 20 20" fill="none"><path d="M10 4v12M4 10h12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            </button>
            <button v-if="vista === 'clientes'" class="fab-inline" @click="showCreateCustomerModal = true; resetNewCustomer()">
              <svg viewBox="0 0 20 20" fill="none"><path d="M10 4v12M4 10h12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            </button>
            <button class="fab-inline fab-prices" @click="abrirPricePlanModal" title="Precios">
              <svg viewBox="0 0 20 20" fill="none"><rect x="3" y="5" width="14" height="11" rx="1" stroke="currentColor" stroke-width="1.5"/><path d="M3 8h14" stroke="currentColor" stroke-width="1.5"/></svg>
            </button>
          </div>

          <div class="rol-filters" v-if="vista === 'usuarios'">
            <button class="rol-chip" :class="{ active: filtroRol === 'todos' }" @click="filtroRol = 'todos'">Todos</button>
            <button v-for="rol in rolesUnicos" :key="rol" class="rol-chip" :class="{ active: filtroRol === rol }" @click="filtroRol = rol">{{ rol }}</button>
          </div>
        </div>

        <!-- Loading -->
        <div v-if="loading" class="state-box">
          <div class="bolt-spin">⚡</div>
          <p>Cargando datos...</p>
        </div>

        <!-- Error usuarios -->
        <div v-else-if="errorMsg && vista === 'usuarios'" class="state-box error-state">
          <span>⚠</span><p>{{ errorMsg }}</p>
          <button class="retry-btn" @click="obtenerUsuarios()">Reintentar</button>
        </div>

        <!-- Contenido principal -->
        <div v-else class="content-area">

          <!-- ── USUARIOS: Cards móvil / Tabla desktop ── -->
          <template v-if="vista === 'usuarios'">

            <!-- Cards en móvil -->
            <div class="cards-list mobile-only">
              <div v-for="(user, i) in usuariosFiltrados" :key="user.id"
                   class="user-card" :style="{ animationDelay: i * 40 + 'ms' }">
                <div class="card-top">
                  <div class="card-avatar">{{ user.username.charAt(0).toUpperCase() }}</div>
                  <div class="card-info">
                    <span class="card-name">{{ user.username }}</span>
                    <span class="role-badge" :class="rolColor(user.rol)">{{ user.rol }}</span>
                  </div>
                  <div class="card-actions">
                    <button v-if="user.rol === 'admin'" class="action-btn edit-btn" @click="abrirEditar(user)">
                      <svg viewBox="0 0 20 20" fill="none"><path d="M14.5 3.5l2 2L6 16l-3 1 1-3 11.5-10.5z" stroke="currentColor" stroke-width="1.3" stroke-linejoin="round"/></svg>
                    </button>
                    <button v-if="user.rol === 'admin' || user.rol === 'customer'" class="action-btn delete-btn" @click="eliminarUsuario(user.id)">
                      <svg viewBox="0 0 20 20" fill="none"><path d="M5 7h10l-1 9H6L5 7zM3 7h14M8 7V5h4v2" stroke="currentColor" stroke-width="1.3" stroke-linecap="round" stroke-linejoin="round"/></svg>
                    </button>
                  </div>
                </div>
                <div class="card-pass">
                  <span class="card-pass-lbl">Contraseña</span>
                  <span class="pass-dots">{{ maskPassword(user.password) }}</span>
                </div>
              </div>
              <div v-if="!usuariosFiltrados.length" class="empty-card">
                ⚡ Sin resultados para "{{ busqueda }}"
              </div>
            </div>

            <!-- Tabla en desktop -->
            <div class="table-wrapper desktop-only">
              <div class="table-scroll">
                <table>
                  <thead>
                    <tr>
                      <th>Usuario</th>
                      <th>Contraseña</th>
                      <th>Rol</th>
                      <th class="th-center">Acciones</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(user, i) in usuariosFiltrados" :key="user.id"
                        class="table-row" :style="{ animationDelay: i * 40 + 'ms' }">
                      <td>
                        <div class="user-cell">
                          <div class="avatar">{{ user.username.charAt(0).toUpperCase() }}</div>
                          <span>{{ user.username }}</span>
                        </div>
                      </td>
                      <td><span class="pass-dots">{{ maskPassword(user.password) }}</span></td>
                      <td><span class="role-badge" :class="rolColor(user.rol)">{{ user.rol }}</span></td>
                      <td class="td-center">
                        <button v-if="user.rol === 'admin'" class="action-btn edit-btn" @click="abrirEditar(user)" title="Editar">
                          <svg viewBox="0 0 20 20" fill="none"><path d="M14.5 3.5l2 2L6 16l-3 1 1-3 11.5-10.5z" stroke="currentColor" stroke-width="1.3" stroke-linejoin="round"/></svg>
                        </button>
                        <button v-if="user.rol === 'admin' || user.rol === 'customer'" class="action-btn delete-btn" @click="eliminarUsuario(user.id)" title="Eliminar">
                          <svg viewBox="0 0 20 20" fill="none"><path d="M5 7h10l-1 9H6L5 7zM3 7h14M8 7V5h4v2" stroke="currentColor" stroke-width="1.3" stroke-linecap="round" stroke-linejoin="round"/></svg>
                        </button>
                        <span v-if="user.rol !== 'admin' && user.rol !== 'customer'" class="no-action">—</span>
                      </td>
                    </tr>
                    <tr v-if="!usuariosFiltrados.length">
                      <td colspan="4" class="empty-row">⚡ Sin resultados para "{{ busqueda }}"</td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <div class="table-footer">
                Mostrando <strong>{{ usuariosFiltrados.length }}</strong> de <strong>{{ users.length }}</strong> usuarios
              </div>
            </div>

          </template>

          <!-- ── CLIENTES: Cards móvil / Tabla desktop ── -->
          <template v-else>

            <!-- Cards en móvil -->
            <div class="cards-list mobile-only">
              <div v-for="(c, i) in clientesFiltrados" :key="c.idcustomer"
                   class="user-card customer-card" :style="{ animationDelay: i * 40 + 'ms' }">
                <div class="card-top">
                  <div class="card-avatar card-avatar--orange">{{ c.name?.charAt(0).toUpperCase() }}</div>
                  <div class="card-info">
                    <span class="card-name">{{ c.name }} {{ c.lastname }}</span>
                    <span class="card-sub">{{ c.cedula }} · +57 {{ c.phone }}</span>
                  </div>
                  <div class="card-actions">
                    <button class="action-btn edit-btn" @click="abrirEditarCustomer(c)">
                      <svg viewBox="0 0 20 20" fill="none"><path d="M14.5 3.5l2 2L6 16l-3 1 1-3 11.5-10.5z" stroke="currentColor" stroke-width="1.3" stroke-linejoin="round"/></svg>
                    </button>
                    <button class="action-btn delete-btn" @click="eliminarCustomer(c.idcustomer)">
                      <svg viewBox="0 0 20 20" fill="none"><path d="M5 7h10l-1 9H6L5 7zM3 7h14M8 7V5h4v2" stroke="currentColor" stroke-width="1.3" stroke-linecap="round" stroke-linejoin="round"/></svg>
                    </button>
                  </div>
                </div>
                <div class="card-bottom">
                  <div class="card-pass">
                    <span class="card-pass-lbl">Contraseña</span>
                    <span class="pass-dots">{{ maskPassword(c.password) }}</span>
                  </div>
                  <button class="btn-plan" @click="verPlanes(c)">
                    <svg viewBox="0 0 20 20" fill="none">
                      <rect x="3" y="3" width="14" height="14" rx="1.5" stroke="currentColor" stroke-width="1.4"/>
                      <path d="M7 10h6M10 7v6" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/>
                    </svg>
                    Ver Planes
                  </button>
                </div>
              </div>
              <div v-if="!clientesFiltrados.length && !busqueda" class="empty-card empty-card--full">
                <span class="empty-icon">👥</span>
                <p class="empty-titulo">Sin clientes registrados</p>
                <p class="empty-sub">Usa el botón + para crear el primero.</p>
              </div>
              <div v-else-if="!clientesFiltrados.length" class="empty-card">
                ⚡ Sin resultados para "{{ busqueda }}"
              </div>
            </div>

            <!-- Tabla desktop -->
            <div class="table-wrapper desktop-only">
              <div class="table-scroll">
                <table>
                  <thead>
                    <tr>
                      <th>Cédula</th>
                      <th>Nombre</th>
                      <th>Apellido</th>
                      <th>Teléfono</th>
                      <th>Contraseña</th>
                      <th class="th-center">Planes</th>
                      <th class="th-center">Acciones</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(c, i) in clientesFiltrados" :key="c.idcustomer"
                        class="table-row" :style="{ animationDelay: i * 40 + 'ms' }">
                      <td>{{ c.cedula }}</td>
                      <td>
                        <div class="user-cell">
                          <div class="avatar">{{ c.name?.charAt(0).toUpperCase() }}</div>
                          <span>{{ c.name }}</span>
                        </div>
                      </td>
                      <td>{{ c.lastname }}</td>
                      <td>+57 {{ c.phone }}</td>
                      <td><span class="pass-dots">{{ maskPassword(c.password) }}</span></td>
                      <td class="td-center">
                        <button class="btn-plan" @click="verPlanes(c)">
                          <svg viewBox="0 0 20 20" fill="none">
                            <rect x="3" y="3" width="14" height="14" rx="1.5" stroke="currentColor" stroke-width="1.4"/>
                            <path d="M7 10h6M10 7v6" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/>
                          </svg>
                          Ver Planes
                        </button>
                      </td>
                      <td class="td-center">
                        <button class="action-btn edit-btn" @click="abrirEditarCustomer(c)" title="Editar">
                          <svg viewBox="0 0 20 20" fill="none"><path d="M14.5 3.5l2 2L6 16l-3 1 1-3 11.5-10.5z" stroke="currentColor" stroke-width="1.3" stroke-linejoin="round"/></svg>
                        </button>
                        <button class="action-btn delete-btn" @click="eliminarCustomer(c.idcustomer)" title="Eliminar">
                          <svg viewBox="0 0 20 20" fill="none"><path d="M5 7h10l-1 9H6L5 7zM3 7h14M8 7V5h4v2" stroke="currentColor" stroke-width="1.3" stroke-linecap="round" stroke-linejoin="round"/></svg>
                        </button>
                      </td>
                    </tr>
                    <tr v-if="!clientesFiltrados.length && !busqueda">
                      <td colspan="7" class="empty-row">
                        <div class="empty-clientes">
                          <span class="empty-icon">👥</span>
                          <p class="empty-titulo">No hay clientes registrados</p>
                          <p class="empty-sub">Crea el primer cliente usando el botón en el menú lateral.</p>
                        </div>
                      </td>
                    </tr>
                    <tr v-else-if="!clientesFiltrados.length">
                      <td colspan="7" class="empty-row">⚡ Sin resultados para "{{ busqueda }}"</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <div v-if="errorMsg && vista === 'clientes'" class="clientes-error-banner">
                <span>⚠</span> {{ errorMsg }}
                <button class="retry-inline" @click="obtenerClientes()">Reintentar</button>
              </div>

              <div class="table-footer">
                Mostrando <strong>{{ clientesFiltrados.length }}</strong> de <strong>{{ customers.length }}</strong> clientes
              </div>
            </div>

            <!-- Error clientes móvil -->
            <div v-if="errorMsg && vista === 'clientes'" class="clientes-error-banner mobile-only">
              <span>⚠</span> {{ errorMsg }}
              <button class="retry-inline" @click="obtenerClientes()">Reintentar</button>
            </div>

          </template>

          <!-- Conteo móvil -->
          <div class="mobile-count mobile-only">
            <span v-if="vista === 'usuarios'">{{ usuariosFiltrados.length }} / {{ users.length }} usuarios</span>
            <span v-else>{{ clientesFiltrados.length }} / {{ customers.length }} clientes</span>
          </div>

        </div>

      </main>
    </div>

    <!-- ══════════════════════════════════════════════════════════ -->
    <!-- MODALES                                                    -->
    <!-- ══════════════════════════════════════════════════════════ -->

    <!-- Modal: Crear Admin -->
    <transition name="modal-fade">
      <div v-if="showModal" class="modal-overlay" @click.self="showModal = false; resetNewUser()">
        <div class="modal-box">
          <div class="modal-header">
            <div>
              <span class="modal-eyebrow">⚡ Nuevo acceso</span>
              <h2>Crear Administrador</h2>
            </div>
            <button class="modal-close" @click="showModal = false; resetNewUser()">✕</button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label>Usuario <span class="req">*</span></label>
              <input v-model="newUser.username" type="text" placeholder="Nombre de usuario" />
            </div>
            <div class="form-group">
              <label>Contraseña <span class="req">*</span></label>
              <input v-model="newUser.password" type="password" placeholder="••••••••" />
            </div>
            <div class="form-group">
              <label>Confirmar Contraseña <span class="req">*</span></label>
              <input v-model="newUser.confirmPassword" type="password" placeholder="••••••••" />
            </div>
            <p v-if="modalError" class="modal-error">⚠ {{ modalError }}</p>
          </div>
          <div class="modal-footer">
            <button class="btn-cancel" @click="showModal = false; resetNewUser()">Cancelar</button>
            <button class="btn-save" @click="crearAdmin" :disabled="guardandoUser">
              {{ guardandoUser ? 'Guardando...' : 'Crear Admin' }}
            </button>
          </div>
        </div>
      </div>
    </transition>

    <!-- Modal: Editar Usuario -->
    <transition name="modal-fade">
      <div v-if="showEditModal" class="modal-overlay" @click.self="showEditModal = false; resetEditUser()">
        <div class="modal-box">
          <div class="modal-header">
            <div>
              <span class="modal-eyebrow">⚡ Modificar acceso</span>
              <h2>Editar Usuario</h2>
            </div>
            <button class="modal-close" @click="showEditModal = false; resetEditUser()">✕</button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label>Usuario <span class="req">*</span></label>
              <input v-model="editUser.username" type="text" />
            </div>
            <div class="form-group">
              <label>Nueva Contraseña <span class="hint-lbl">(opcional)</span></label>
              <input v-model="editUser.password" type="password" placeholder="Dejar vacío para no cambiar" />
            </div>
            <div class="form-group">
              <label>Confirmar Contraseña</label>
              <input v-model="editUser.confirmPassword" type="password" placeholder="••••••••" />
            </div>
            <p v-if="modalError" class="modal-error">⚠ {{ modalError }}</p>
          </div>
          <div class="modal-footer">
            <button class="btn-cancel" @click="showEditModal = false; resetEditUser()">Cancelar</button>
            <button class="btn-save" @click="actualizarUsuario" :disabled="guardandoUser">
              {{ guardandoUser ? 'Guardando...' : 'Actualizar' }}
            </button>
          </div>
        </div>
      </div>
    </transition>

    <!-- Modal: Crear Cliente -->
    <transition name="modal-fade">
      <div v-if="showCreateCustomerModal" class="modal-overlay" @click.self="showCreateCustomerModal = false; resetNewCustomer()">
        <div class="modal-box modal-wide">
          <div class="modal-header">
            <div>
              <span class="modal-eyebrow">⚡ Nuevo miembro</span>
              <h2>Crear Cliente</h2>
            </div>
            <button class="modal-close" @click="showCreateCustomerModal = false; resetNewCustomer()">✕</button>
          </div>
          <div class="modal-body">
            <div class="modal-grid">
              <div class="form-group">
                <label>Cédula <span class="req">*</span></label>
                <input v-model="newCustomer.cedula" type="text" placeholder="Número de cédula" />
              </div>
              <div class="form-group">
                <label>Teléfono <span class="req">*</span></label>
                <input v-model="newCustomer.phone" type="text" placeholder="(300) 000 0000" maxlength="10"
                  @input="newCustomer.phone = newCustomer.phone.replace(/\D/g, '').slice(0,10)" />
              </div>
              <div class="form-group">
                <label>Nombre <span class="req">*</span></label>
                <input v-model="newCustomer.name" type="text" placeholder="Nombre" />
              </div>
              <div class="form-group">
                <label>Apellido <span class="req">*</span></label>
                <input v-model="newCustomer.lastname" type="text" placeholder="Apellido" />
              </div>
              <div class="form-group">
                <label>Contraseña <span class="req">*</span></label>
                <input v-model="newCustomer.password" type="password" placeholder="••••••••" />
              </div>
              <div class="form-group">
                <label>Confirmar Contraseña <span class="req">*</span></label>
                <input v-model="newCustomer.confirmPassword" type="password" placeholder="••••••••" />
              </div>
            </div>
            <p v-if="modalCustomerError" class="modal-error">⚠ {{ modalCustomerError }}</p>
          </div>
          <div class="modal-footer">
            <button class="btn-cancel" @click="showCreateCustomerModal = false; resetNewCustomer()">Cancelar</button>
            <button class="btn-save" @click="crearCustomer" :disabled="guardandoCustomer">
              {{ guardandoCustomer ? 'Guardando...' : 'Crear Cliente' }}
            </button>
          </div>
        </div>
      </div>
    </transition>

    <!-- Modal: Editar Cliente -->
    <transition name="modal-fade">
      <div v-if="showEditCustomerModal" class="modal-overlay" @click.self="showEditCustomerModal = false; resetEditCustomer()">
        <div class="modal-box modal-wide">
          <div class="modal-header">
            <div>
              <span class="modal-eyebrow">⚡ Modificar miembro</span>
              <h2>Editar Cliente</h2>
            </div>
            <button class="modal-close" @click="showEditCustomerModal = false; resetEditCustomer()">✕</button>
          </div>
          <div class="modal-body">
            <div class="modal-grid">
              <div class="form-group">
                <label>Cédula <span class="req">*</span></label>
                <input v-model="editCustomer.cedula" type="text" />
              </div>
              <div class="form-group">
                <label>Teléfono</label>
                <input v-model="editCustomer.phone" type="text" maxlength="10"
                  @input="editCustomer.phone = editCustomer.phone.replace(/\D/g, '').slice(0,10)" />
              </div>
              <div class="form-group">
                <label>Nombre <span class="req">*</span></label>
                <input v-model="editCustomer.name" type="text" />
              </div>
              <div class="form-group">
                <label>Apellido</label>
                <input v-model="editCustomer.lastname" type="text" />
              </div>
              <div class="form-group">
                <label>Nueva Contraseña <span class="hint-lbl">(opcional)</span></label>
                <input v-model="editCustomer.password" type="password" placeholder="Dejar vacío para no cambiar" />
              </div>
              <div class="form-group">
                <label>Confirmar Contraseña</label>
                <input v-model="editCustomer.confirmPassword" type="password" placeholder="••••••••" />
              </div>
            </div>
            <p v-if="modalCustomerError" class="modal-error">⚠ {{ modalCustomerError }}</p>
          </div>
          <div class="modal-footer">
            <button class="btn-cancel" @click="showEditCustomerModal = false; resetEditCustomer()">Cancelar</button>
            <button class="btn-save" @click="actualizarCustomer" :disabled="guardandoCustomer">
              {{ guardandoCustomer ? 'Guardando...' : 'Actualizar' }}
            </button>
          </div>
        </div>
      </div>
    </transition>

    <!-- Modal: Precios del Plan -->
    <transition name="modal-fade">
      <div v-if="showPricePlanModal" class="modal-overlay" @click.self="showPricePlanModal = false; cancelarEditarPrecio()">
        <div class="modal-box pp-modal">
          <div class="modal-header pp-header">
            <div>
              <span class="modal-eyebrow">⚡ Tarifas activas</span>
              <h2>Precios del Plan</h2>
            </div>
            <button class="modal-close" @click="showPricePlanModal = false; cancelarEditarPrecio()">✕</button>
          </div>
          <div class="modal-body pp-body">
            <div v-if="loadingPrices" class="pp-state">
              <div class="bolt-spin">⚡</div>
              <p>Cargando tarifas...</p>
            </div>
            <div v-else-if="priceError" class="pp-state pp-state--error">
              <span class="pp-state-icon">⚠</span>
              <p>{{ priceError }}</p>
              <button class="retry-btn" @click="obtenerPricePlans">Reintentar</button>
            </div>
            <div v-else class="pp-grid">
              <div
                v-for="(plan, idx) in sortedPricePlans"
                :key="plan.idpriceplan"
                class="pp-card"
                :class="{ 'pp-card--active': editingPriceId === plan.idpriceplan }"
                :style="{
                  '--c':  getMeta(plan.typeplan).color,
                  '--bg': getMeta(plan.typeplan).bg,
                  '--br': getMeta(plan.typeplan).border,
                  animationDelay: idx * 70 + 'ms'
                }"
              >
                <div class="pp-card-bar"></div>
                <div class="pp-card-head">
                  <span class="pp-icon">{{ getMeta(plan.typeplan).icon }}</span>
                  <div class="pp-card-head-text">
                    <span class="pp-type">{{ getMeta(plan.typeplan).label }}</span>
                    <span class="pp-id">#{{ plan.idpriceplan }}</span>
                  </div>
                </div>
                <div class="pp-price-area">
                  <template v-if="editingPriceId !== plan.idpriceplan">
                    <p class="pp-price">{{ formatPrice(plan.price) }}</p>
                    <p class="pp-price-desc">por {{ getMeta(plan.typeplan).label.toLowerCase() }}</p>
                  </template>
                  <template v-else>
                    <div class="pp-edit-box">
                      <span class="pp-currency">$</span>
                      <input v-model="editingPrice" type="number" min="0" step="100" class="pp-input"
                        @keyup.enter="guardarPrecio(plan)" @keyup.esc="cancelarEditarPrecio" />
                    </div>
                    <p class="pp-edit-hint">↵ guardar · Esc cancelar</p>
                  </template>
                </div>
                <div class="pp-actions">
                  <template v-if="editingPriceId !== plan.idpriceplan">
                    <button class="pp-btn pp-btn--edit" @click="iniciarEditarPrecio(plan)">
                      <svg viewBox="0 0 20 20" fill="none"><path d="M14.5 3.5l2 2L6 16l-3 1 1-3 11.5-10.5z" stroke="currentColor" stroke-width="1.3" stroke-linejoin="round"/></svg>
                      Cambiar valor
                    </button>
                  </template>
                  <template v-else>
                    <button class="pp-btn pp-btn--save" @click="guardarPrecio(plan)" :disabled="savingPrice">
                      <svg viewBox="0 0 20 20" fill="none"><path d="M4 10l4.5 4.5L16 6" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/></svg>
                      {{ savingPrice ? 'Guardando...' : 'Guardar' }}
                    </button>
                    <button class="pp-btn pp-btn--cancel" @click="cancelarEditarPrecio">
                      <svg viewBox="0 0 20 20" fill="none"><path d="M6 6l8 8M14 6l-8 8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/></svg>
                    </button>
                  </template>
                </div>
              </div>
              <div v-if="!sortedPricePlans.length" class="pp-empty">⚡ No hay planes registrados</div>
            </div>
            <p v-if="savePriceError" class="modal-error pp-save-err">⚠ {{ savePriceError }}</p>
          </div>
          <div class="modal-footer">
            <button class="btn-cancel" @click="showPricePlanModal = false; cancelarEditarPrecio()">Cerrar</button>
          </div>
        </div>
      </div>
    </transition>

  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Barlow+Condensed:wght@400;700;900&family=Barlow:wght@300;400;500&display=swap');

/* ── Reset ── */
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

/* ── Variables ── */
:root {
  --topbar-h: 52px;
  --sidebar-w: 230px;
  --gold: #f5c500;
  --orange: #ff7a00;
  --bg: #0a0a0a;
  --surface: #111;
  --border: rgba(245,197,0,.12);
}

/* ── Screen ── */
.screen {
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  background: #0a0a0a;
  color: #ccc;
  font-family: 'Barlow', sans-serif;
  position: relative;
}

/* ── Background decorativo ── */
.bg-layer { position: fixed; inset: 0; pointer-events: none; z-index: 0; }
.hex { position: absolute; }
.hex-1 { width: 420px; height: 420px; background: conic-gradient(from 30deg,rgba(245,197,0,.08) 0deg,rgba(255,122,0,.06) 80deg,transparent 80deg); clip-path: polygon(50% 0%,100% 25%,100% 75%,50% 100%,0% 75%,0% 25%); top: -100px; right: -80px; }
.hex-2 { width: 260px; height: 260px; background: conic-gradient(from 200deg,rgba(245,197,0,.05) 0deg,transparent 80deg); clip-path: polygon(50% 0%,100% 25%,100% 75%,50% 100%,0% 75%,0% 25%); bottom: -60px; left: 220px; }
.slash { position: absolute; }
.slash-1 { width: 2px; height: 60vh; background: linear-gradient(to bottom,transparent,rgba(245,197,0,.12) 50%,transparent); top: 20%; left: 220px; transform: rotate(12deg); }
.slash-2 { width: 1px; height: 40vh; background: linear-gradient(to bottom,transparent,rgba(255,122,0,.08) 50%,transparent); top: 30%; left: 230px; transform: rotate(12deg); }
.noise { position: absolute; inset: 0; background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='1'/%3E%3C/svg%3E"); opacity: .02; }

/* ── Topbar (solo móvil) ── */
.topbar {
  display: none;
  position: fixed;
  top: 0; left: 0; right: 0;
  height: var(--topbar-h);
  background: #0d0d0d;
  border-bottom: 1px solid var(--border);
  z-index: 200;
  align-items: center;
  padding: 0 16px;
  gap: 12px;
}

.topbar-menu {
  width: 36px; height: 36px;
  background: none; border: 1px solid rgba(245,197,0,.2);
  cursor: pointer;
  display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 5px;
  flex-shrink: 0;
  transition: border-color .2s;
}
.topbar-menu:hover { border-color: var(--gold); }
.topbar-menu span {
  display: block; width: 16px; height: 1.5px; background: #aaa;
  transition: all .25s;
  transform-origin: center;
}
.topbar-menu.open span:nth-child(1) { transform: translateY(6.5px) rotate(45deg); background: var(--gold); }
.topbar-menu.open span:nth-child(2) { opacity: 0; transform: scaleX(0); }
.topbar-menu.open span:nth-child(3) { transform: translateY(-6.5px) rotate(-45deg); background: var(--gold); }

.topbar-logo { flex: 1; }
.topbar-logo svg { width: 90px; height: auto; filter: drop-shadow(0 0 10px rgba(245,197,0,.15)); }

.topbar-actions { display: flex; align-items: center; gap: 8px; }
.topbar-stat { display: flex; flex-direction: column; align-items: flex-end; }
.topbar-stat-val { font-family: 'Barlow Condensed', sans-serif; font-size: 1.1rem; font-weight: 900; color: var(--gold); line-height: 1; }
.topbar-stat-lbl { font-size: .55rem; color: #555; letter-spacing: .15em; text-transform: uppercase; }

/* ── Sidebar overlay (móvil) ── */
.sidebar-overlay {
  position: fixed; inset: 0; background: rgba(0,0,0,.7);
  z-index: 299; backdrop-filter: blur(2px);
}
.overlay-fade-enter-active, .overlay-fade-leave-active { transition: opacity .25s; }
.overlay-fade-enter-from, .overlay-fade-leave-to { opacity: 0; }

/* ── Layout ── */
.layout {
  position: relative; z-index: 1;
  display: flex; width: 100%; height: 100vh;
  background: #0a0a0a;
}

/* ── Sidebar ── */
.sidebar {
  width: var(--sidebar-w);
  height: 100vh;
  background: #0d0d0d;
  border-right: 1px solid var(--border);
  display: flex; flex-direction: column;
  flex-shrink: 0;
  z-index: 300;
}

.sidebar-inner {
  display: flex; flex-direction: column;
  height: 100%;
  padding: 32px 0 24px;
}

.sidebar-logo { padding: 0 24px 28px; border-bottom: 1px solid rgba(245,197,0,.08); margin-bottom: 24px; }
.sidebar-logo svg { width: 130px; height: auto; filter: drop-shadow(0 0 18px rgba(245,197,0,.18)); }
.sidebar-sub { display: block; font-family: 'Barlow Condensed',sans-serif; font-size: .55rem; letter-spacing: .35em; text-transform: uppercase; color: #444; margin-top: 4px; }

.sidebar-nav { display: flex; flex-direction: column; gap: 4px; padding: 0 12px; flex: 1; }
.nav-divider { height: 1px; background: rgba(245,197,0,.06); margin: 8px 2px; }

.nav-item { display: flex; align-items: center; gap: 10px; padding: 11px 14px; color: #555; font-size: .82rem; font-weight: 500; letter-spacing: .04em; text-decoration: none; border-left: 2px solid transparent; transition: all .2s; cursor: pointer; user-select: none; }
.nav-item svg { width: 16px; height: 16px; flex-shrink: 0; }
.nav-item:hover { color: #f0f0f0; background: rgba(245,197,0,.05); }
.nav-item.active { color: var(--gold); border-left-color: var(--gold); background: rgba(245,197,0,.07); }

.btn-add { display: flex; align-items: center; gap: 8px; background: transparent; border: 1px solid rgba(245,197,0,.3); color: var(--gold); font-family: 'Barlow Condensed',sans-serif; font-size: .72rem; font-weight: 700; letter-spacing: .18em; text-transform: uppercase; padding: 10px 14px; cursor: pointer; transition: all .2s; margin-top: 8px; width: 100%; }
.btn-add svg { width: 14px; height: 14px; flex-shrink: 0; }
.btn-add:hover { background: rgba(245,197,0,.08); border-color: var(--gold); }
.btn-prices { border-color: rgba(245,197,0,.2); color: #c9a200; }
.btn-prices:hover { background: rgba(245,197,0,.07); border-color: rgba(245,197,0,.5); color: var(--gold); }

.sidebar-footer { padding: 20px 16px 0; border-top: 1px solid rgba(245,197,0,.08); display: flex; align-items: center; gap: 10px; margin-top: 12px; }
.user-chip { display: flex; align-items: center; gap: 10px; flex: 1; min-width: 0; }
.user-avatar { width: 32px; height: 32px; background: linear-gradient(135deg,var(--gold),var(--orange)); color: #0a0a0a; font-family: 'Barlow Condensed',sans-serif; font-weight: 900; font-size: .95rem; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.user-info { display: flex; flex-direction: column; min-width: 0; }
.user-name { font-size: .78rem; color: #e0e0e0; font-weight: 500; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.user-role { font-size: .6rem; color: #555; letter-spacing: .08em; }
.logout-btn { background: none; border: 1px solid #1e1e1e; color: #444; cursor: pointer; width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; transition: all .2s; }
.logout-btn svg { width: 15px; height: 15px; }
.logout-btn:hover { border-color: rgba(245,197,0,.4); color: var(--gold); }

/* ── Main ── */
.main {
  flex: 1; height: 100vh;
  overflow-y: auto;
  padding: 48px 60px;
  display: flex; flex-direction: column; gap: 26px;
  min-width: 0;
  background: #0a0a0a;
}

/* ── Page header (desktop) ── */
.page-header { display: flex; align-items: flex-end; justify-content: space-between; flex-wrap: wrap; gap: 16px; flex-shrink: 0; }
.page-eyebrow { display: block; font-family: 'Barlow Condensed',sans-serif; font-size: .65rem; letter-spacing: .3em; text-transform: uppercase; color: var(--gold); margin-bottom: 6px; }
.page-title h1 { font-family: 'Barlow Condensed',sans-serif; font-size: 3.2rem; font-weight: 900; text-transform: uppercase; letter-spacing: .04em; color: #fff; line-height: 1; }
.header-meta { display: flex; gap: 12px; }
.stat-pill { display: flex; flex-direction: column; align-items: center; background: #111; border: 1px solid rgba(245,197,0,.12); padding: 10px 22px; gap: 2px; }
.stat-pill.accent { border-color: rgba(245,197,0,.3); background: rgba(245,197,0,.06); }
.stat-val { font-family: 'Barlow Condensed',sans-serif; font-size: 1.7rem; font-weight: 900; color: var(--gold); line-height: 1; }
.stat-lbl { font-size: .58rem; letter-spacing: .2em; text-transform: uppercase; color: #555; }

/* ── Mobile tabs ── */
.mobile-tabs { display: none; }

/* ── Filtros ── */
.filters-bar { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; flex-shrink: 0; }
.search-wrap { position: relative; flex: 1; min-width: 200px; max-width: 500px; }
.search-icon { position: absolute; left: 13px; top: 50%; transform: translateY(-50%); width: 15px; height: 15px; color: #444; pointer-events: none; }
.search-input { width: 100%; background: #111; border: 1px solid rgba(245,197,0,.12); color: #f0f0f0; font-family: 'Barlow',sans-serif; font-size: .85rem; padding: 11px 14px 11px 38px; outline: none; transition: border-color .2s, background .2s; }
.search-input::placeholder { color: #333; }
.search-input:focus { border-color: rgba(245,197,0,.4); background: #0f0e09; }
.rol-filters { display: flex; gap: 8px; flex-wrap: wrap; }
.rol-chip { background: transparent; border: 1px solid rgba(245,197,0,.15); color: #555; font-family: 'Barlow Condensed',sans-serif; font-size: .75rem; font-weight: 700; letter-spacing: .15em; text-transform: uppercase; padding: 8px 18px; cursor: pointer; transition: all .2s; }
.rol-chip:hover { border-color: rgba(245,197,0,.4); color: #ccc; }
.rol-chip.active { background: rgba(245,197,0,.1); border-color: var(--gold); color: var(--gold); }

/* Botones de acción inline en móvil */
.mobile-actions { display: none; gap: 8px; }
.fab-inline { width: 38px; height: 38px; background: rgba(245,197,0,.08); border: 1px solid rgba(245,197,0,.3); color: var(--gold); cursor: pointer; display: flex; align-items: center; justify-content: center; transition: all .2s; flex-shrink: 0; }
.fab-inline svg { width: 16px; height: 16px; }
.fab-inline:hover { background: rgba(245,197,0,.15); border-color: var(--gold); }
.fab-prices { background: rgba(245,197,0,.04); border-color: rgba(245,197,0,.18); color: #c9a200; }
.fab-prices:hover { background: rgba(245,197,0,.1); border-color: rgba(245,197,0,.4); color: var(--gold); }

/* ── Estado loading/error ── */
.state-box { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 14px; padding: 80px 20px; color: #555; font-size: .9rem; flex: 1; }
.bolt-spin { font-size: 2.4rem; animation: boltPulse 1.4s ease-in-out infinite; filter: drop-shadow(0 0 10px var(--gold)); }
@keyframes boltPulse { 0%,100% { opacity: 1; } 50% { opacity: .3; } }
.error-state { color: #e05a45; }
.retry-btn { background: none; border: 1px solid rgba(220,50,30,.4); color: #e05a45; font-family: 'Barlow Condensed',sans-serif; font-size: .75rem; letter-spacing: .2em; text-transform: uppercase; padding: 8px 20px; cursor: pointer; transition: all .2s; }
.retry-btn:hover { background: rgba(220,50,30,.1); }

/* ── Content area ── */
.content-area { display: flex; flex-direction: column; gap: 0; flex: 1; min-height: 0; }

/* ── Cards (móvil) ── */
.cards-list { display: flex; flex-direction: column; gap: 10px; padding-bottom: 80px; }
.user-card {
  background: #111; border: 1px solid rgba(245,197,0,.1);
  padding: 14px 16px;
  animation: rowIn .35s ease both;
  transition: border-color .2s;
}
.user-card:hover { border-color: rgba(245,197,0,.2); }
.customer-card { border-color: rgba(255,122,0,.1); }
.customer-card:hover { border-color: rgba(255,122,0,.25); }

.card-top { display: flex; align-items: center; gap: 10px; }
.card-avatar { width: 38px; height: 38px; flex-shrink: 0; background: linear-gradient(135deg,rgba(245,197,0,.2),rgba(255,122,0,.15)); border: 1px solid rgba(245,197,0,.2); color: var(--gold); font-family: 'Barlow Condensed',sans-serif; font-weight: 900; font-size: 1rem; display: flex; align-items: center; justify-content: center; }
.card-avatar--orange { background: linear-gradient(135deg,rgba(255,122,0,.2),rgba(245,197,0,.1)); border-color: rgba(255,122,0,.25); color: var(--orange); }
.card-info { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 4px; }
.card-name { font-size: .9rem; color: #e0e0e0; font-weight: 500; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.card-sub { font-size: .72rem; color: #555; letter-spacing: .03em; }
.card-actions { display: flex; gap: 4px; flex-shrink: 0; }
.card-pass { display: flex; align-items: center; gap: 8px; margin-top: 10px; padding-top: 10px; border-top: 1px solid rgba(255,255,255,.04); }
.card-pass-lbl { font-family: 'Barlow Condensed',sans-serif; font-size: .6rem; letter-spacing: .2em; text-transform: uppercase; color: #444; }

.card-bottom { display: flex; align-items: center; justify-content: space-between; margin-top: 10px; padding-top: 10px; border-top: 1px solid rgba(255,255,255,.04); gap: 10px; }

.empty-card { text-align: center; color: #444; padding: 40px 20px; font-size: .85rem; letter-spacing: .05em; background: #111; border: 1px dashed #1a1a1a; }
.empty-card--full { padding: 56px 20px; display: flex; flex-direction: column; align-items: center; gap: 10px; }

.mobile-count { text-align: center; font-size: .68rem; color: #333; letter-spacing: .1em; text-transform: uppercase; padding: 8px; }

/* ── Tabla (desktop) ── */
.table-wrapper { display: flex; flex-direction: column; border: 1px solid rgba(245,197,0,.12); background: #111; animation: fadeUp .45s cubic-bezier(.16,1,.3,1) both; flex: 1; min-height: 0; }
@keyframes fadeUp { from { opacity: 0; transform: translateY(12px); } to { opacity: 1; transform: translateY(0); } }
.table-scroll { overflow-y: auto; flex: 1; min-height: 0; }
table { width: 100%; border-collapse: collapse; font-size: .88rem; }
thead { background: rgba(245,197,0,.06); border-bottom: 1px solid rgba(245,197,0,.15); position: sticky; top: 0; z-index: 2; }
th { font-family: 'Barlow Condensed',sans-serif; font-size: .68rem; font-weight: 700; letter-spacing: .25em; text-transform: uppercase; color: var(--gold); padding: 16px 20px; text-align: left; white-space: nowrap; }
.th-center { text-align: center; }
td { padding: 14px 20px; border-bottom: 1px solid rgba(255,255,255,.04); color: #ccc; vertical-align: middle; }
.table-row { transition: background .15s; animation: rowIn .35s ease both; }
.table-row:hover { background: rgba(245,197,0,.04); }
.table-row:last-child td { border-bottom: none; }
@keyframes rowIn { from { opacity: 0; transform: translateX(-6px); } to { opacity: 1; transform: translateX(0); } }
.user-cell { display: flex; align-items: center; gap: 10px; }
.avatar { width: 32px; height: 32px; background: linear-gradient(135deg,rgba(245,197,0,.2),rgba(255,122,0,.15)); border: 1px solid rgba(245,197,0,.2); color: var(--gold); font-family: 'Barlow Condensed',sans-serif; font-weight: 900; font-size: .9rem; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.pass-dots { color: #2a2a2a; letter-spacing: .15em; font-size: .85rem; }
.role-badge { display: inline-block; font-family: 'Barlow Condensed',sans-serif; font-size: .7rem; font-weight: 700; letter-spacing: .2em; text-transform: uppercase; padding: 4px 12px; border: 1px solid; }
.badge-admin    { color: var(--gold); border-color: rgba(245,197,0,.4); background: rgba(245,197,0,.07); }
.badge-customer { color: var(--orange); border-color: rgba(255,122,0,.4); background: rgba(255,122,0,.07); }
.badge-default  { color: #888; border-color: rgba(136,136,136,.3); background: rgba(136,136,136,.06); }
.td-center { text-align: center; }
.no-action { color: #333; font-size: .8rem; }
.action-btn { background: none; border: 1px solid transparent; width: 30px; height: 30px; display: inline-flex; align-items: center; justify-content: center; cursor: pointer; transition: all .2s; }
.action-btn svg { width: 13px; height: 13px; }
.edit-btn { color: #555; }
.edit-btn:hover { color: var(--gold); border-color: rgba(245,197,0,.4); background: rgba(245,197,0,.07); }
.delete-btn { color: #555; }
.delete-btn:hover { color: #e05a45; border-color: rgba(220,50,30,.4); background: rgba(220,50,30,.07); }
.btn-plan { display: inline-flex; align-items: center; gap: 5px; background: transparent; border: 1px solid rgba(245,197,0,.3); color: var(--gold); font-family: 'Barlow Condensed',sans-serif; font-size: .7rem; font-weight: 700; letter-spacing: .15em; text-transform: uppercase; padding: 6px 12px; cursor: pointer; transition: all .2s; white-space: nowrap; }
.btn-plan svg { width: 12px; height: 12px; flex-shrink: 0; }
.btn-plan:hover { background: rgba(245,197,0,.1); border-color: var(--gold); }
.empty-row { text-align: center; color: #444; padding: 56px 20px !important; font-size: .88rem; letter-spacing: .05em; }
.empty-clientes { display: flex; flex-direction: column; align-items: center; gap: 10px; }
.empty-icon { font-size: 2.2rem; opacity: .3; filter: grayscale(1); }
.empty-titulo { font-family: 'Barlow Condensed',sans-serif; font-size: 1.1rem; font-weight: 700; letter-spacing: .08em; text-transform: uppercase; color: #333; }
.empty-sub { font-size: .75rem; color: #2a2a2a; max-width: 320px; line-height: 1.5; }
.clientes-error-banner { display: flex; align-items: center; gap: 10px; padding: 10px 20px; background: rgba(220,50,30,.06); border-top: 1px solid rgba(220,50,30,.2); color: #e05a45; font-size: .78rem; flex-shrink: 0; }
.retry-inline { background: none; border: 1px solid rgba(220,50,30,.35); color: #e05a45; font-family: 'Barlow Condensed',sans-serif; font-size: .65rem; letter-spacing: .15em; text-transform: uppercase; padding: 4px 12px; cursor: pointer; transition: all .2s; margin-left: auto; }
.retry-inline:hover { background: rgba(220,50,30,.1); }
.table-footer { padding: 13px 20px; border-top: 1px solid rgba(245,197,0,.08); font-size: .7rem; color: #444; letter-spacing: .06em; flex-shrink: 0; }
.table-footer strong { color: #888; }

/* ── Modales ── */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,.85); display: flex; align-items: center; justify-content: center; z-index: 1000; padding: 16px; }
.modal-box { background: #111; border: 1px solid rgba(245,197,0,.2); box-shadow: 0 40px 100px rgba(0,0,0,.9); width: 100%; max-width: 380px; max-height: 90vh; overflow-y: auto; display: flex; flex-direction: column; animation: modalIn .3s cubic-bezier(.16,1,.3,1) both; }
.modal-wide { max-width: 560px; }
@keyframes modalIn { from { opacity: 0; transform: scale(.96) translateY(12px); } to { opacity: 1; transform: scale(1) translateY(0); } }
.modal-header { display: flex; align-items: flex-start; justify-content: space-between; padding: 22px 26px 18px; border-bottom: 1px solid rgba(245,197,0,.1); position: sticky; top: 0; background: #111; z-index: 1; }
.modal-eyebrow { display: block; font-family: 'Barlow Condensed',sans-serif; font-size: .6rem; letter-spacing: .3em; text-transform: uppercase; color: var(--gold); margin-bottom: 5px; }
.modal-header h2 { font-family: 'Barlow Condensed',sans-serif; font-size: 1.5rem; font-weight: 900; text-transform: uppercase; color: #fff; line-height: 1; }
.modal-close { background: none; border: none; color: #444; font-size: 1rem; cursor: pointer; padding: 2px; transition: color .2s; line-height: 1; flex-shrink: 0; }
.modal-close:hover { color: var(--gold); }
.modal-body { padding: 22px 26px; }
.modal-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
.form-group { display: flex; flex-direction: column; gap: 6px; margin-bottom: 14px; }
.form-group:last-child { margin-bottom: 0; }
.modal-grid .form-group { margin-bottom: 0; }
.form-group label { font-family: 'Barlow Condensed',sans-serif; font-size: .62rem; letter-spacing: .22em; text-transform: uppercase; color: #666; font-weight: 700; }
.req { color: var(--gold); }
.hint-lbl { color: #444; font-size: .58rem; letter-spacing: .1em; text-transform: none; font-weight: 400; }
.form-group input { background: #0d0d0d; border: 1px solid rgba(245,197,0,.15); color: #f0f0f0; font-family: 'Barlow',sans-serif; font-size: .85rem; padding: 10px 12px; outline: none; transition: border-color .2s, background .2s; }
.form-group input::placeholder { color: #333; }
.form-group input:focus { border-color: rgba(245,197,0,.5); background: #0f0e09; }
.modal-error { color: #e05a45; font-size: .78rem; padding: 9px 13px; background: rgba(220,50,30,.06); border: 1px solid rgba(220,50,30,.25); }
.modal-footer { display: flex; align-items: center; justify-content: flex-end; gap: 10px; padding: 18px 26px; border-top: 1px solid rgba(245,197,0,.1); position: sticky; bottom: 0; background: #111; }
.btn-cancel { background: transparent; border: 1px solid #222; color: #555; font-family: 'Barlow Condensed',sans-serif; font-size: .72rem; font-weight: 700; letter-spacing: .2em; text-transform: uppercase; padding: 9px 20px; cursor: pointer; transition: all .2s; }
.btn-cancel:hover { border-color: rgba(245,197,0,.3); color: #aaa; }
.btn-save { background: var(--gold); border: none; color: #0a0a0a; font-family: 'Barlow Condensed',sans-serif; font-size: .78rem; font-weight: 900; letter-spacing: .2em; text-transform: uppercase; padding: 10px 24px; cursor: pointer; transition: all .2s; }
.btn-save:hover:not(:disabled) { background: #ffd700; box-shadow: 0 0 18px rgba(245,197,0,.3); }
.btn-save:disabled { opacity: .5; cursor: not-allowed; }
.modal-fade-enter-active, .modal-fade-leave-active { transition: opacity .22s; }
.modal-fade-enter-from, .modal-fade-leave-to { opacity: 0; }

/* ── Price Plan Modal ── */
.pp-modal  { max-width: 620px; }
.pp-header { background: linear-gradient(160deg, #0e0e0e 0%, #110f06 100%); }
.pp-body   { padding: 20px 20px 10px; }
.pp-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 12px; }
.pp-card { position: relative; background: #0d0d0d; border: 1px solid var(--br); display: flex; flex-direction: column; overflow: hidden; transition: box-shadow .25s, border-color .25s, transform .2s; animation: ppCardIn .4s cubic-bezier(.16,1,.3,1) both; }
@keyframes ppCardIn { from { opacity: 0; transform: translateY(14px) scale(.96); } to { opacity: 1; transform: translateY(0) scale(1); } }
.pp-card:hover { transform: translateY(-2px); box-shadow: 0 8px 32px color-mix(in srgb, var(--c) 14%, transparent), 0 0 0 1px var(--br); }
.pp-card--active { border-color: var(--c) !important; box-shadow: 0 0 0 1px var(--c), 0 10px 40px color-mix(in srgb, var(--c) 20%, transparent) !important; }
.pp-card-bar { height: 3px; background: linear-gradient(90deg, var(--c), color-mix(in srgb, var(--c) 40%, transparent)); }
.pp-card-head { display: flex; align-items: center; gap: 10px; padding: 14px 16px 6px; }
.pp-icon { font-size: 1.5rem; line-height: 1; filter: drop-shadow(0 0 8px color-mix(in srgb, var(--c) 55%, transparent)); }
.pp-card-head-text { display: flex; flex-direction: column; gap: 2px; flex: 1; }
.pp-type { font-family: 'Barlow Condensed', sans-serif; font-size: 1.1rem; font-weight: 900; text-transform: uppercase; letter-spacing: .14em; color: var(--c); line-height: 1; }
.pp-id { font-family: 'Barlow Condensed', sans-serif; font-size: .58rem; letter-spacing: .18em; color: #333; text-transform: uppercase; }
.pp-price-area { padding: 8px 16px 14px; flex: 1; display: flex; flex-direction: column; gap: 3px; border-top: 1px solid rgba(255,255,255,.04); margin: 0 12px; }
.pp-price { font-family: 'Barlow Condensed', sans-serif; font-size: 1.8rem; font-weight: 900; color: #fff; letter-spacing: .01em; line-height: 1; }
.pp-price-desc { font-size: .65rem; color: #3a3a3a; letter-spacing: .12em; text-transform: uppercase; }
.pp-edit-box { display: flex; align-items: center; gap: 6px; background: #070707; border: 1px solid var(--c); padding: 6px 10px; }
.pp-currency { font-family: 'Barlow Condensed', sans-serif; font-weight: 900; font-size: 1.1rem; color: var(--c); line-height: 1; }
.pp-input { background: transparent; border: none; color: #f0f0f0; font-family: 'Barlow Condensed', sans-serif; font-size: 1.4rem; font-weight: 700; outline: none; width: 100%; min-width: 0; line-height: 1; }
.pp-input::-webkit-inner-spin-button, .pp-input::-webkit-outer-spin-button { opacity: .2; }
.pp-edit-hint { font-size: .58rem; color: #2e2e2e; letter-spacing: .06em; margin-top: 4px; }
.pp-actions { display: flex; gap: 6px; padding: 0 12px 12px; }
.pp-btn { display: inline-flex; align-items: center; justify-content: center; gap: 6px; font-family: 'Barlow Condensed', sans-serif; font-size: .7rem; font-weight: 700; letter-spacing: .18em; text-transform: uppercase; cursor: pointer; transition: all .2s; padding: 9px 14px; border: 1px solid; }
.pp-btn svg { width: 11px; height: 11px; flex-shrink: 0; }
.pp-btn--edit { width: 100%; background: var(--bg); border-color: var(--br); color: var(--c); }
.pp-btn--edit:hover { background: color-mix(in srgb, var(--c) 16%, transparent); border-color: var(--c); box-shadow: 0 0 10px color-mix(in srgb, var(--c) 20%, transparent); }
.pp-btn--save { flex: 1; background: var(--c); border-color: var(--c); color: #0a0a0a; font-weight: 900; }
.pp-btn--save:hover:not(:disabled) { filter: brightness(1.12); box-shadow: 0 0 14px color-mix(in srgb, var(--c) 35%, transparent); }
.pp-btn--save:disabled { opacity: .45; cursor: not-allowed; }
.pp-btn--cancel { background: transparent; border-color: #1e1e1e; color: #555; padding: 9px 11px; }
.pp-btn--cancel:hover { border-color: rgba(220,50,30,.4); color: #e05a45; }
.pp-state { display: flex; flex-direction: column; align-items: center; gap: 12px; padding: 44px 20px; color: #555; font-size: .88rem; }
.pp-state--error { color: #e05a45; }
.pp-state-icon { font-size: 1.8rem; }
.pp-empty { grid-column: span 2; text-align: center; color: #333; padding: 40px 20px; font-size: .88rem; letter-spacing: .05em; }
.pp-save-err { margin-top: 14px; }

/* Helpers */
.mobile-only  { display: none !important; }
.desktop-only { display: flex !important; }

/* ══════════════════════════════════════════════
   RESPONSIVE
══════════════════════════════════════════════ */

/* ── Tablet intermedio ── */
@media (max-width: 1024px) {
  .main { padding: 36px 32px; }
  .page-title h1 { font-size: 2.4rem; }
}

/* ── Móvil / tablet ── */
@media (max-width: 768px) {

  /* Fondo oscuro en TODO */
  html, body, .screen, .layout, .main,
  .content-area, .cards-list { background: #0a0a0a !important; }

  /* Mostrar topbar */
  .topbar { display: flex; }

  /* Sidebar: oculto por defecto, desliza */
  .sidebar {
    position: fixed;
    top: 0; left: 0; bottom: 0;
    width: 260px;
    transform: translateX(-100%);
    transition: transform .28s cubic-bezier(.16,1,.3,1);
  }
  .sidebar--open {
    transform: translateX(0);
    box-shadow: 8px 0 40px rgba(0,0,0,.8);
  }

  /* Layout: block, main cubre toda la pantalla */
  .layout { display: block; position: relative; }
  .main {
    position: fixed;
    top: 0; left: 0; right: 0; bottom: 0;
    width: 100vw;
    height: 100vh;
    overflow-y: auto;
    -webkit-overflow-scrolling: touch;
    padding: calc(var(--topbar-h) + 14px) 14px 20px;
    gap: 12px;
  }

  /* Ocultar page-header de escritorio */
  .page-header { display: none; }

  /* Tabs de navegación */
  .mobile-tabs {
    display: flex; gap: 0;
    border: 1px solid rgba(245,197,0,.18);
    overflow: hidden;
    flex-shrink: 0;
  }
  .mobile-tab {
    flex: 1; display: flex; align-items: center; justify-content: center; gap: 7px;
    padding: 12px 8px;
    background: #0d0d0d; border: none;
    color: #555; font-family: 'Barlow Condensed',sans-serif;
    font-size: .82rem; font-weight: 700; letter-spacing: .12em; text-transform: uppercase;
    cursor: pointer; transition: all .2s;
    border-right: 1px solid rgba(245,197,0,.1);
  }
  .mobile-tab:last-child { border-right: none; }
  .mobile-tab svg { width: 14px; height: 14px; }
  .mobile-tab:hover { color: #ccc; background: rgba(245,197,0,.04); }
  .mobile-tab.active {
    color: var(--gold);
    background: rgba(245,197,0,.08);
    border-bottom: 2px solid var(--gold);
  }

  /* Barra filtros */
  .filters-bar { gap: 8px; flex-wrap: wrap; }
  .search-wrap { max-width: 100%; flex: 1 1 0; min-width: 0; }
  .search-input { background: #0d0d0d; }

  /* Botones de acción inline */
  .mobile-actions { display: flex; }

  /* Cards visibles, tabla oculta */
  .mobile-only  { display: flex !important; }
  .desktop-only { display: none !important; }

  /* Cards */
  .user-card, .empty-card { background: #0f0f0f; }
  .cards-list { flex: 1; }

  /* Rol chips horizontal */
  .rol-filters { flex-wrap: nowrap; overflow-x: auto; padding-bottom: 2px; width: 100%; }
  .rol-chip { white-space: nowrap; padding: 8px 14px; font-size: .72rem; }

  /* Modales desde abajo */
  .modal-overlay { padding: 0; align-items: flex-end; }
  .modal-box {
    max-width: 100%; width: 100%;
    border: 1px solid rgba(245,197,0,.2);
    border-bottom: none;
    border-radius: 0;
    max-height: 92vh;
    background: #111;
    animation: modalInUp .3s cubic-bezier(.16,1,.3,1) both;
  }
  @keyframes modalInUp {
    from { opacity: 0; transform: translateY(30px); }
    to   { opacity: 1; transform: translateY(0); }
  }
  .modal-grid { grid-template-columns: 1fr; }
  .pp-modal { max-width: 100%; }
  .pp-grid { grid-template-columns: 1fr; }
  .pp-price { font-size: 1.5rem; }
}

/* ── Móvil pequeño ── */
@media (max-width: 400px) {
  .main { padding: calc(var(--topbar-h) + 10px) 10px 16px; gap: 10px; }
  .topbar { padding: 0 10px; }
  .search-input { font-size: .8rem; padding: 10px 12px 10px 34px; }
  .card-name { font-size: .82rem; }
  .mobile-tab { font-size: .72rem; padding: 10px 4px; gap: 5px; }
  .mobile-tab svg { width: 12px; height: 12px; }
}
</style>