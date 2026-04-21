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
const id       = localStorage.getItem('id') || ''

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
    customers.value = await res.json()
  } catch (err) { errorMsg.value = err.message }
  finally { loading.value = false }
}

const cambiarVista = (tipo) => {
  vista.value = tipo; busqueda.value = ''
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

const newUser = ref({ username: '', password: '', confirmPassword: '' })
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

const newCustomer = ref({ cedula: '', name: '', lastname: '', phone: '', password: '', confirmPassword: '' })
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

// Metadatos visuales por tipo de plan
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
          <a class="nav-item" :class="{ active: vista === 'usuarios' }" @click="cambiarVista('usuarios')">
            <svg viewBox="0 0 20 20" fill="none"><path d="M4 6h12M4 10h12M4 14h8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/></svg>
            Usuarios
          </a>
          <a class="nav-item" :class="{ active: vista === 'clientes' }" @click="cambiarVista('clientes')">
            <svg viewBox="0 0 20 20" fill="none"><circle cx="10" cy="8" r="3.5" stroke="currentColor" stroke-width="1.5"/><path d="M3 18c0-3.314 3.134-6 7-6s7 2.686 7 6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/></svg>
            Clientes
          </a>

          <button class="btn-add" v-if="vista === 'usuarios'" @click="showModal = true; resetNewUser()">
            <svg viewBox="0 0 20 20" fill="none"><circle cx="10" cy="10" r="7.5" stroke="currentColor" stroke-width="1.4"/><path d="M10 7v6M7 10h6" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/></svg>
            Nuevo Admin
          </button>

          <button class="btn-add btn-prices" @click="abrirPricePlanModal">
            <svg viewBox="0 0 20 20" fill="none">
              <rect x="3" y="5" width="14" height="11" rx="1" stroke="currentColor" stroke-width="1.4"/>
              <path d="M3 8h14" stroke="currentColor" stroke-width="1.4"/>
              <path d="M7 12h2M12 12h2" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/>
            </svg>
            Ver Precios
          </button>

          <button class="btn-add" v-if="vista === 'clientes'" @click="showCreateCustomerModal = true; resetNewCustomer()">
            <svg viewBox="0 0 20 20" fill="none"><circle cx="10" cy="10" r="7.5" stroke="currentColor" stroke-width="1.4"/><path d="M10 7v6M7 10h6" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/></svg>
            Nuevo Cliente
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
      </aside>

      <!-- ── Main ── -->
      <main class="main">

        <header class="page-header">
          <div class="page-title">
            <span class="page-eyebrow">⚡ Gestión del sistema</span>
            <h1>{{ username }} <span class="id-hint">#{{ id }}</span></h1>
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

        <!-- Filtros -->
        <div class="filters-bar">
          <div class="search-wrap">
            <svg class="search-icon" viewBox="0 0 20 20" fill="none">
              <circle cx="8.5" cy="8.5" r="5" stroke="currentColor" stroke-width="1.4"/>
              <path d="M13 13l3.5 3.5" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/>
            </svg>
            <input v-model="busqueda" type="text"
              :placeholder="vista === 'usuarios' ? 'Buscar por usuario, rol o ID...' : 'Buscar por nombre, cédula o teléfono...'"
              class="search-input" />
          </div>
          <div class="rol-filters" v-if="vista === 'usuarios'">
            <button class="rol-chip" :class="{ active: filtroRol === 'todos' }" @click="filtroRol = 'todos'">Todos</button>
            <button v-for="rol in rolesUnicos" :key="rol" class="rol-chip" :class="{ active: filtroRol === rol }" @click="filtroRol = rol">{{ rol }}</button>
          </div>
        </div>

        <!-- Loading / Error -->
        <div v-if="loading" class="state-box">
          <div class="bolt-spin">⚡</div>
          <p>Cargando datos...</p>
        </div>
        <div v-else-if="errorMsg" class="state-box error-state">
          <span>⚠</span><p>{{ errorMsg }}</p>
          <button class="retry-btn" @click="vista === 'usuarios' ? obtenerUsuarios() : obtenerClientes()">Reintentar</button>
        </div>

        <!-- Tabla -->
        <div v-else class="table-wrapper">
          <div class="table-scroll">

            <!-- ── USUARIOS ── -->
            <table v-if="vista === 'usuarios'">
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

            <!-- ── CLIENTES ── -->
            <table v-else>
              <thead>
                <tr>
                  <th>ID</th>
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
                  <td><span class="id-badge-sm">{{ c.idcustomer }}</span></td>
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
                <tr v-if="!clientesFiltrados.length">
                  <td colspan="8" class="empty-row">⚡ No hay clientes registrados</td>
                </tr>
              </tbody>
            </table>

          </div>

          <div class="table-footer">
            <span v-if="vista === 'usuarios'">
              Mostrando <strong>{{ usuariosFiltrados.length }}</strong> de <strong>{{ users.length }}</strong> usuarios
            </span>
            <span v-else>
              Mostrando <strong>{{ clientesFiltrados.length }}</strong> de <strong>{{ customers.length }}</strong> clientes
            </span>
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

    <!-- ══════════════════════════════════════════════════════════ -->
    <!-- Modal: Precios del Plan — diseño tarjetas                  -->
    <!-- ══════════════════════════════════════════════════════════ -->
    <transition name="modal-fade">
      <div v-if="showPricePlanModal" class="modal-overlay" @click.self="showPricePlanModal = false; cancelarEditarPrecio()">
        <div class="modal-box pp-modal">

          <!-- Header -->
          <div class="modal-header pp-header">
            <div>
              <span class="modal-eyebrow">⚡ Tarifas activas</span>
              <h2>Precios del Plan</h2>
            </div>
            <button class="modal-close" @click="showPricePlanModal = false; cancelarEditarPrecio()">✕</button>
          </div>

          <!-- Body -->
          <div class="modal-body pp-body">

            <!-- Loading -->
            <div v-if="loadingPrices" class="pp-state">
              <div class="bolt-spin">⚡</div>
              <p>Cargando tarifas...</p>
            </div>

            <!-- Error carga -->
            <div v-else-if="priceError" class="pp-state pp-state--error">
              <span class="pp-state-icon">⚠</span>
              <p>{{ priceError }}</p>
              <button class="retry-btn" @click="obtenerPricePlans">Reintentar</button>
            </div>

            <!-- Grid de tarjetas -->
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
                <!-- Franja de color superior -->
                <div class="pp-card-bar"></div>

                <!-- Cabecera de la tarjeta -->
                <div class="pp-card-head">
                  <span class="pp-icon">{{ getMeta(plan.typeplan).icon }}</span>
                  <div class="pp-card-head-text">
                    <span class="pp-type">{{ getMeta(plan.typeplan).label }}</span>
                    <span class="pp-id">#{{ plan.idpriceplan }}</span>
                  </div>
                </div>

                <!-- Área del precio -->
                <div class="pp-price-area">

                  <!-- Modo visualización -->
                  <template v-if="editingPriceId !== plan.idpriceplan">
                    <p class="pp-price">{{ formatPrice(plan.price) }}</p>
                    <p class="pp-price-desc">por {{ getMeta(plan.typeplan).label.toLowerCase() }}</p>
                  </template>

                  <!-- Modo edición -->
                  <template v-else>
                    <div class="pp-edit-box">
                      <span class="pp-currency">$</span>
                      <input
                        v-model="editingPrice"
                        type="number"
                        min="0"
                        step="100"
                        class="pp-input"
                        @keyup.enter="guardarPrecio(plan)"
                        @keyup.esc="cancelarEditarPrecio"
                      />
                    </div>
                    <p class="pp-edit-hint">↵ guardar &nbsp;·&nbsp; Esc cancelar</p>
                  </template>
                </div>

                <!-- Botones de acción -->
                <div class="pp-actions">

                  <!-- Estado normal -->
                  <template v-if="editingPriceId !== plan.idpriceplan">
                    <button class="pp-btn pp-btn--edit" @click="iniciarEditarPrecio(plan)">
                      <svg viewBox="0 0 20 20" fill="none">
                        <path d="M14.5 3.5l2 2L6 16l-3 1 1-3 11.5-10.5z" stroke="currentColor" stroke-width="1.3" stroke-linejoin="round"/>
                      </svg>
                      Cambiar valor
                    </button>
                  </template>

                  <!-- Estado edición -->
                  <template v-else>
                    <button class="pp-btn pp-btn--save" @click="guardarPrecio(plan)" :disabled="savingPrice">
                      <svg viewBox="0 0 20 20" fill="none">
                        <path d="M4 10l4.5 4.5L16 6" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
                      </svg>
                      {{ savingPrice ? 'Guardando...' : 'Guardar' }}
                    </button>
                    <button class="pp-btn pp-btn--cancel" @click="cancelarEditarPrecio">
                      <svg viewBox="0 0 20 20" fill="none">
                        <path d="M6 6l8 8M14 6l-8 8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                      </svg>
                    </button>
                  </template>

                </div>
              </div><!-- /pp-card -->

              <!-- Sin planes -->
              <div v-if="!sortedPricePlans.length" class="pp-empty">
                ⚡ No hay planes registrados
              </div>
            </div><!-- /pp-grid -->

            <!-- Error al guardar -->
            <p v-if="savePriceError" class="modal-error pp-save-err">⚠ {{ savePriceError }}</p>

          </div><!-- /pp-body -->

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

*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

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
.btn-add { display: flex; align-items: center; gap: 8px; background: transparent; border: 1px solid rgba(245,197,0,.3); color: #f5c500; font-family: 'Barlow Condensed',sans-serif; font-size: .72rem; font-weight: 700; letter-spacing: .18em; text-transform: uppercase; padding: 10px 14px; cursor: pointer; transition: all .2s; margin-top: 8px; }
.btn-add svg { width: 14px; height: 14px; flex-shrink: 0; }
.btn-add:hover { background: rgba(245,197,0,.08); border-color: #f5c500; }
.btn-prices { border-color: rgba(245,197,0,.2); color: #c9a200; }
.btn-prices:hover { background: rgba(245,197,0,.07); border-color: rgba(245,197,0,.5); color: #f5c500; }
.sidebar-footer { padding: 20px 16px 0; border-top: 1px solid rgba(245,197,0,.08); display: flex; align-items: center; gap: 10px; margin-top: 12px; }
.user-chip { display: flex; align-items: center; gap: 10px; flex: 1; min-width: 0; }
.user-avatar { width: 32px; height: 32px; background: linear-gradient(135deg,#f5c500,#ff7a00); color: #0a0a0a; font-family: 'Barlow Condensed',sans-serif; font-weight: 900; font-size: .95rem; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.user-info { display: flex; flex-direction: column; min-width: 0; }
.user-name { font-size: .78rem; color: #e0e0e0; font-weight: 500; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.user-role { font-size: .6rem; color: #555; letter-spacing: .08em; }
.logout-btn { background: none; border: 1px solid #1e1e1e; color: #444; cursor: pointer; width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; transition: all .2s; }
.logout-btn svg { width: 15px; height: 15px; }
.logout-btn:hover { border-color: rgba(245,197,0,.4); color: #f5c500; }

/* ── Main ── */
.main { flex: 1; height: 100vh; overflow-y: auto; padding: 48px 60px; display: flex; flex-direction: column; gap: 26px; min-width: 0; }
.page-header { display: flex; align-items: flex-end; justify-content: space-between; flex-wrap: wrap; gap: 16px; flex-shrink: 0; }
.page-eyebrow { display: block; font-family: 'Barlow Condensed',sans-serif; font-size: .65rem; letter-spacing: .3em; text-transform: uppercase; color: #f5c500; margin-bottom: 6px; }
.page-title h1 { font-family: 'Barlow Condensed',sans-serif; font-size: 3.2rem; font-weight: 900; text-transform: uppercase; letter-spacing: .04em; color: #fff; line-height: 1; }
.id-hint { color: #f5c500; font-size: 2.4rem; }
.header-meta { display: flex; gap: 12px; }
.stat-pill { display: flex; flex-direction: column; align-items: center; background: #111; border: 1px solid rgba(245,197,0,.12); padding: 10px 22px; gap: 2px; }
.stat-pill.accent { border-color: rgba(245,197,0,.3); background: rgba(245,197,0,.06); }
.stat-val { font-family: 'Barlow Condensed',sans-serif; font-size: 1.7rem; font-weight: 900; color: #f5c500; line-height: 1; }
.stat-lbl { font-size: .58rem; letter-spacing: .2em; text-transform: uppercase; color: #555; }
.filters-bar { display: flex; align-items: center; gap: 16px; flex-wrap: wrap; flex-shrink: 0; }
.search-wrap { position: relative; flex: 1; min-width: 240px; max-width: 500px; }
.search-icon { position: absolute; left: 13px; top: 50%; transform: translateY(-50%); width: 15px; height: 15px; color: #444; pointer-events: none; }
.search-input { width: 100%; background: #111; border: 1px solid rgba(245,197,0,.12); color: #f0f0f0; font-family: 'Barlow',sans-serif; font-size: .85rem; padding: 11px 14px 11px 38px; outline: none; transition: border-color .2s, background .2s; }
.search-input::placeholder { color: #333; }
.search-input:focus { border-color: rgba(245,197,0,.4); background: #0f0e09; }
.rol-filters { display: flex; gap: 8px; flex-wrap: wrap; }
.rol-chip { background: transparent; border: 1px solid rgba(245,197,0,.15); color: #555; font-family: 'Barlow Condensed',sans-serif; font-size: .75rem; font-weight: 700; letter-spacing: .15em; text-transform: uppercase; padding: 8px 18px; cursor: pointer; transition: all .2s; }
.rol-chip:hover { border-color: rgba(245,197,0,.4); color: #ccc; }
.rol-chip.active { background: rgba(245,197,0,.1); border-color: #f5c500; color: #f5c500; }
.state-box { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 14px; padding: 80px 20px; color: #555; font-size: .9rem; flex: 1; }
.bolt-spin { font-size: 2.4rem; animation: boltPulse 1.4s ease-in-out infinite; filter: drop-shadow(0 0 10px #f5c500); }
@keyframes boltPulse { 0%,100% { opacity: 1; } 50% { opacity: .3; } }
.error-state { color: #e05a45; }
.retry-btn { background: none; border: 1px solid rgba(220,50,30,.4); color: #e05a45; font-family: 'Barlow Condensed',sans-serif; font-size: .75rem; letter-spacing: .2em; text-transform: uppercase; padding: 8px 20px; cursor: pointer; transition: all .2s; }
.retry-btn:hover { background: rgba(220,50,30,.1); }
.table-wrapper { display: flex; flex-direction: column; border: 1px solid rgba(245,197,0,.12); background: #111; animation: fadeUp .45s cubic-bezier(.16,1,.3,1) both; flex: 1; min-height: 0; }
@keyframes fadeUp { from { opacity: 0; transform: translateY(12px); } to { opacity: 1; transform: translateY(0); } }
.table-scroll { overflow-y: auto; flex: 1; min-height: 0; }
table { width: 100%; border-collapse: collapse; font-size: .88rem; }
thead { background: rgba(245,197,0,.06); border-bottom: 1px solid rgba(245,197,0,.15); position: sticky; top: 0; z-index: 2; }
th { font-family: 'Barlow Condensed',sans-serif; font-size: .68rem; font-weight: 700; letter-spacing: .25em; text-transform: uppercase; color: #f5c500; padding: 16px 20px; text-align: left; white-space: nowrap; }
.th-center { text-align: center; }
td { padding: 14px 20px; border-bottom: 1px solid rgba(255,255,255,.04); color: #ccc; vertical-align: middle; }
.table-row { transition: background .15s; animation: rowIn .35s ease both; }
.table-row:hover { background: rgba(245,197,0,.04); }
.table-row:last-child td { border-bottom: none; }
@keyframes rowIn { from { opacity: 0; transform: translateX(-6px); } to { opacity: 1; transform: translateX(0); } }
.user-cell { display: flex; align-items: center; gap: 10px; }
.avatar { width: 32px; height: 32px; background: linear-gradient(135deg,rgba(245,197,0,.2),rgba(255,122,0,.15)); border: 1px solid rgba(245,197,0,.2); color: #f5c500; font-family: 'Barlow Condensed',sans-serif; font-weight: 900; font-size: .9rem; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.pass-dots { color: #2a2a2a; letter-spacing: .15em; font-size: .85rem; }
.id-badge-sm { display: inline-flex; align-items: center; justify-content: center; min-width: 28px; height: 22px; padding: 0 6px; background: rgba(245,197,0,.07); border: 1px solid rgba(245,197,0,.18); font-family: 'Barlow Condensed',sans-serif; font-weight: 700; font-size: .75rem; color: #f5c500; }
.role-badge { display: inline-block; font-family: 'Barlow Condensed',sans-serif; font-size: .7rem; font-weight: 700; letter-spacing: .2em; text-transform: uppercase; padding: 4px 12px; border: 1px solid; }
.badge-admin    { color: #f5c500; border-color: rgba(245,197,0,.4); background: rgba(245,197,0,.07); }
.badge-customer { color: #ff7a00; border-color: rgba(255,122,0,.4); background: rgba(255,122,0,.07); }
.badge-default  { color: #888; border-color: rgba(136,136,136,.3); background: rgba(136,136,136,.06); }
.td-center { text-align: center; }
.no-action { color: #333; font-size: .8rem; }
.action-btn { background: none; border: 1px solid transparent; width: 30px; height: 30px; display: inline-flex; align-items: center; justify-content: center; cursor: pointer; transition: all .2s; }
.action-btn svg { width: 13px; height: 13px; }
.edit-btn { color: #555; }
.edit-btn:hover { color: #f5c500; border-color: rgba(245,197,0,.4); background: rgba(245,197,0,.07); }
.delete-btn { color: #555; }
.delete-btn:hover { color: #e05a45; border-color: rgba(220,50,30,.4); background: rgba(220,50,30,.07); }
.btn-plan { display: inline-flex; align-items: center; gap: 5px; background: transparent; border: 1px solid rgba(245,197,0,.3); color: #f5c500; font-family: 'Barlow Condensed',sans-serif; font-size: .7rem; font-weight: 700; letter-spacing: .15em; text-transform: uppercase; padding: 6px 12px; cursor: pointer; transition: all .2s; white-space: nowrap; }
.btn-plan svg { width: 12px; height: 12px; flex-shrink: 0; }
.btn-plan:hover { background: rgba(245,197,0,.1); border-color: #f5c500; }
.empty-row { text-align: center; color: #444; padding: 56px 20px !important; font-size: .88rem; letter-spacing: .05em; }
.table-footer { padding: 13px 20px; border-top: 1px solid rgba(245,197,0,.08); font-size: .7rem; color: #444; letter-spacing: .06em; flex-shrink: 0; }
.table-footer strong { color: #888; }

/* ══ MODALES BASE ══ */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,.82); display: flex; align-items: center; justify-content: center; z-index: 1000; padding: 20px; }
.modal-box { background: #111; border: 1px solid rgba(245,197,0,.2); box-shadow: 0 40px 100px rgba(0,0,0,.9); width: 100%; max-width: 380px; display: flex; flex-direction: column; animation: modalIn .3s cubic-bezier(.16,1,.3,1) both; }
.modal-wide { max-width: 560px; }
@keyframes modalIn { from { opacity: 0; transform: scale(.96) translateY(12px); } to { opacity: 1; transform: scale(1) translateY(0); } }
.modal-header { display: flex; align-items: flex-start; justify-content: space-between; padding: 22px 26px 18px; border-bottom: 1px solid rgba(245,197,0,.1); }
.modal-eyebrow { display: block; font-family: 'Barlow Condensed',sans-serif; font-size: .6rem; letter-spacing: .3em; text-transform: uppercase; color: #f5c500; margin-bottom: 5px; }
.modal-header h2 { font-family: 'Barlow Condensed',sans-serif; font-size: 1.5rem; font-weight: 900; text-transform: uppercase; color: #fff; line-height: 1; }
.modal-close { background: none; border: none; color: #444; font-size: 1rem; cursor: pointer; padding: 2px; transition: color .2s; line-height: 1; }
.modal-close:hover { color: #f5c500; }
.modal-body { padding: 22px 26px; }
.modal-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
.form-group { display: flex; flex-direction: column; gap: 6px; margin-bottom: 14px; }
.form-group:last-child { margin-bottom: 0; }
.modal-grid .form-group { margin-bottom: 0; }
.form-group label { font-family: 'Barlow Condensed',sans-serif; font-size: .62rem; letter-spacing: .22em; text-transform: uppercase; color: #666; font-weight: 700; }
.req { color: #f5c500; }
.hint-lbl { color: #444; font-size: .58rem; letter-spacing: .1em; text-transform: none; font-weight: 400; }
.form-group input { background: #0d0d0d; border: 1px solid rgba(245,197,0,.15); color: #f0f0f0; font-family: 'Barlow',sans-serif; font-size: .85rem; padding: 10px 12px; outline: none; transition: border-color .2s, background .2s; }
.form-group input::placeholder { color: #333; }
.form-group input:focus { border-color: rgba(245,197,0,.5); background: #0f0e09; }
.modal-error { color: #e05a45; font-size: .78rem; padding: 9px 13px; background: rgba(220,50,30,.06); border: 1px solid rgba(220,50,30,.25); }
.modal-footer { display: flex; align-items: center; justify-content: flex-end; gap: 10px; padding: 18px 26px; border-top: 1px solid rgba(245,197,0,.1); }
.btn-cancel { background: transparent; border: 1px solid #222; color: #555; font-family: 'Barlow Condensed',sans-serif; font-size: .72rem; font-weight: 700; letter-spacing: .2em; text-transform: uppercase; padding: 9px 20px; cursor: pointer; transition: all .2s; }
.btn-cancel:hover { border-color: rgba(245,197,0,.3); color: #aaa; }
.btn-save { background: #f5c500; border: none; color: #0a0a0a; font-family: 'Barlow Condensed',sans-serif; font-size: .78rem; font-weight: 900; letter-spacing: .2em; text-transform: uppercase; padding: 10px 24px; cursor: pointer; transition: all .2s; }
.btn-save:hover:not(:disabled) { background: #ffd700; box-shadow: 0 0 18px rgba(245,197,0,.3); }
.btn-save:disabled { opacity: .5; cursor: not-allowed; }
.modal-fade-enter-active, .modal-fade-leave-active { transition: opacity .22s; }
.modal-fade-enter-from, .modal-fade-leave-to { opacity: 0; }

/* ══════════════════════════════════════════════════════════════════
   PRICE PLAN MODAL — tarjetas rediseñadas
══════════════════════════════════════════════════════════════════ */

.pp-modal  { max-width: 620px; }
.pp-header { background: linear-gradient(160deg, #0e0e0e 0%, #110f06 100%); }
.pp-body   { padding: 20px 20px 10px; }

/* Grid 2 columnas */
.pp-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

/* ── Tarjeta ── */
.pp-card {
  position: relative;
  background: #0d0d0d;
  border: 1px solid var(--br);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  transition: box-shadow .25s, border-color .25s, transform .2s;
  animation: ppCardIn .4s cubic-bezier(.16,1,.3,1) both;
}
@keyframes ppCardIn {
  from { opacity: 0; transform: translateY(14px) scale(.96); }
  to   { opacity: 1; transform: translateY(0) scale(1); }
}

.pp-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 32px color-mix(in srgb, var(--c) 14%, transparent),
              0 0 0 1px var(--br);
}

.pp-card--active {
  border-color: var(--c) !important;
  box-shadow: 0 0 0 1px var(--c),
              0 10px 40px color-mix(in srgb, var(--c) 20%, transparent) !important;
}

/* Franja superior de color */
.pp-card-bar {
  height: 3px;
  background: linear-gradient(90deg, var(--c), color-mix(in srgb, var(--c) 40%, transparent));
}

/* Cabecera de tarjeta */
.pp-card-head {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 16px 6px;
}

.pp-icon {
  font-size: 1.5rem;
  line-height: 1;
  filter: drop-shadow(0 0 8px color-mix(in srgb, var(--c) 55%, transparent));
}

.pp-card-head-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
}

.pp-type {
  font-family: 'Barlow Condensed', sans-serif;
  font-size: 1.1rem;
  font-weight: 900;
  text-transform: uppercase;
  letter-spacing: .14em;
  color: var(--c);
  line-height: 1;
}

.pp-id {
  font-family: 'Barlow Condensed', sans-serif;
  font-size: .58rem;
  letter-spacing: .18em;
  color: #333;
  text-transform: uppercase;
}

/* Área de precio */
.pp-price-area {
  padding: 8px 16px 14px;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 3px;
  /* Línea divisora sutil */
  border-top: 1px solid rgba(255,255,255,.04);
  margin: 0 12px;
}

.pp-price {
  font-family: 'Barlow Condensed', sans-serif;
  font-size: 1.8rem;
  font-weight: 900;
  color: #fff;
  letter-spacing: .01em;
  line-height: 1;
}

.pp-price-desc {
  font-size: .65rem;
  color: #3a3a3a;
  letter-spacing: .12em;
  text-transform: uppercase;
}

/* Modo edición */
.pp-edit-box {
  display: flex;
  align-items: center;
  gap: 6px;
  background: #070707;
  border: 1px solid var(--c);
  padding: 6px 10px;
}

.pp-currency {
  font-family: 'Barlow Condensed', sans-serif;
  font-weight: 900;
  font-size: 1.1rem;
  color: var(--c);
  line-height: 1;
}

.pp-input {
  background: transparent;
  border: none;
  color: #f0f0f0;
  font-family: 'Barlow Condensed', sans-serif;
  font-size: 1.4rem;
  font-weight: 700;
  outline: none;
  width: 100%;
  min-width: 0;
  line-height: 1;
}

.pp-input::-webkit-inner-spin-button,
.pp-input::-webkit-outer-spin-button { opacity: .2; }

.pp-edit-hint {
  font-size: .58rem;
  color: #2e2e2e;
  letter-spacing: .06em;
  margin-top: 4px;
}

/* Botones */
.pp-actions {
  display: flex;
  gap: 6px;
  padding: 0 12px 12px;
}

/* Botón base */
.pp-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-family: 'Barlow Condensed', sans-serif;
  font-size: .7rem;
  font-weight: 700;
  letter-spacing: .18em;
  text-transform: uppercase;
  cursor: pointer;
  transition: all .2s;
  padding: 9px 14px;
  border: 1px solid;
}

.pp-btn svg { width: 11px; height: 11px; flex-shrink: 0; }

/* Cambiar valor */
.pp-btn--edit {
  width: 100%;
  background: var(--bg);
  border-color: var(--br);
  color: var(--c);
}
.pp-btn--edit:hover {
  background: color-mix(in srgb, var(--c) 16%, transparent);
  border-color: var(--c);
  box-shadow: 0 0 10px color-mix(in srgb, var(--c) 20%, transparent);
}

/* Guardar */
.pp-btn--save {
  flex: 1;
  background: var(--c);
  border-color: var(--c);
  color: #0a0a0a;
  font-weight: 900;
}
.pp-btn--save:hover:not(:disabled) { filter: brightness(1.12); box-shadow: 0 0 14px color-mix(in srgb, var(--c) 35%, transparent); }
.pp-btn--save:disabled { opacity: .45; cursor: not-allowed; }

/* Cancelar (solo icono) */
.pp-btn--cancel {
  background: transparent;
  border-color: #1e1e1e;
  color: #555;
  padding: 9px 11px;
}
.pp-btn--cancel:hover { border-color: rgba(220,50,30,.4); color: #e05a45; }

/* Loading / error dentro del modal */
.pp-state { display: flex; flex-direction: column; align-items: center; gap: 12px; padding: 44px 20px; color: #555; font-size: .88rem; }
.pp-state--error { color: #e05a45; }
.pp-state-icon { font-size: 1.8rem; }

/* Sin planes */
.pp-empty { grid-column: span 2; text-align: center; color: #333; padding: 40px 20px; font-size: .88rem; letter-spacing: .05em; }

/* Error al guardar */
.pp-save-err { margin-top: 14px; }

@media (max-width: 860px) {
  .sidebar { display: none; }
  .main { padding: 28px 20px 48px; }
  .modal-grid { grid-template-columns: 1fr; }
  .pp-grid { grid-template-columns: 1fr; }
  .pp-modal { max-width: 100%; }
}
</style>