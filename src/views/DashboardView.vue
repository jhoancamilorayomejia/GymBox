<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const users     = ref([])
const customers = ref([])
const vista = ref('usuarios')
const loading   = ref(true)
const errorMsg  = ref('')
const busqueda  = ref('')
const filtroRol = ref('todos')

// ── Auth helpers ──────────────────────────────────────────────────

const tokenExpirado = (token) => {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.exp * 1000 < Date.now()
  } catch {
    return true
  }
}

const cerrarSesion = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('rol')
  localStorage.removeItem('username')
  router.push('/login')
}

// ── Fetch ─────────────────────────────────────────────────────────

const obtenerUsuarios = async () => {
  loading.value  = true
  errorMsg.value = ''
  try {
    const token = localStorage.getItem('token')
    if (!token || tokenExpirado(token)) { router.push('/login'); return }
    const res = await fetch('/api/users', { headers: { Authorization: `Bearer ${token}` } })
    if (res.status === 401) { localStorage.removeItem('token'); router.push('/login'); return }
    if (!res.ok) throw new Error('Error al cargar usuarios')
    users.value = await res.json()
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}

const obtenerClientes = async () => {
  loading.value  = true
  errorMsg.value = ''
  try {
    const token = localStorage.getItem('token')
    if (!token || tokenExpirado(token)) { router.push('/login'); return }
    const res = await fetch('/api/customers', { headers: { Authorization: `Bearer ${token}` } })
    if (!res.ok) throw new Error('Error al cargar clientes')
    customers.value = await res.json()
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}

const cambiarVista = (tipo) => {
  vista.value = tipo
  busqueda.value = ''
  if (tipo === 'usuarios') obtenerUsuarios()
  else if (tipo === 'clientes') obtenerClientes()
}

// ── Ir a planes de un cliente ──────────────────────────────────────

const verPlanes = (cliente) => {
  router.push({
    name: 'dashboard-plan',
    params: { idcustomer: cliente.idcustomer },
    // pasamos nombre para mostrarlo en el header de la siguiente vista
    state: { nombre: `${cliente.name} ${cliente.lastname}` }
  })
}

// ── Filtrado ──────────────────────────────────────────────────────

const usuariosFiltrados = computed(() => {
  return users.value.filter(u => {
    const matchBusqueda =
      u.username.toLowerCase().includes(busqueda.value.toLowerCase()) ||
      u.rol.toLowerCase().includes(busqueda.value.toLowerCase()) ||
      String(u.id).includes(busqueda.value)
    const matchRol = filtroRol.value === 'todos' || u.rol === filtroRol.value
    return matchBusqueda && matchRol
  })
})

const clientesFiltrados = computed(() => {
  const q = busqueda.value.toLowerCase()
  return customers.value.filter(c =>
    c.name?.toLowerCase().includes(q) ||
    c.lastname?.toLowerCase().includes(q) ||
    String(c.cedula).includes(q) ||
    c.phone?.includes(q)
  )
})

const rolesUnicos = computed(() => [...new Set(users.value.map(u => u.rol))])

const rolColor = (rol) => {
  if (rol === 'admin') return 'badge-admin'
  if (rol === 'customer') return 'badge-customer'
  return 'badge-default'
}

const maskPassword = (pw) => '•'.repeat(Math.min(pw?.length || 8, 12))

const username = localStorage.getItem('username') || 'Admin'

onMounted(() => { obtenerUsuarios() })
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
          <a class="nav-item" :class="{ active: vista === 'usuarios' }" @click="cambiarVista('usuarios')">
            <svg viewBox="0 0 20 20" fill="none"><path d="M4 6h12M4 10h12M4 14h8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/></svg>
            Usuarios
          </a>
          <a class="nav-item" :class="{ active: vista === 'clientes' }" @click="cambiarVista('clientes')">
            <svg viewBox="0 0 20 20" fill="none"><circle cx="10" cy="8" r="3.5" stroke="currentColor" stroke-width="1.5"/><path d="M3 18c0-3.314 3.134-6 7-6s7 2.686 7 6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/></svg>
            Clientes
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

        <header class="page-header">
          <div class="page-title">
            <span class="page-eyebrow">⚡ Gestión del sistema</span>
            <h1>{{ vista === 'usuarios' ? 'Usuarios' : 'Clientes' }}</h1>
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
          </div>
        </header>

        <!-- Filtros -->
        <div class="filters-bar">
          <div class="search-wrap">
            <svg class="search-icon" viewBox="0 0 20 20" fill="none">
              <circle cx="8.5" cy="8.5" r="5" stroke="currentColor" stroke-width="1.4"/>
              <path d="M13 13l3.5 3.5" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/>
            </svg>
            <input
              v-model="busqueda"
              type="text"
              :placeholder="vista === 'usuarios' ? 'Buscar por usuario, rol o ID...' : 'Buscar por nombre, cédula o teléfono...'"
              class="search-input"
            />
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

        <!-- Error -->
        <div v-else-if="errorMsg" class="state-box error-state">
          <span>⚠</span>
          <p>{{ errorMsg }}</p>
          <button class="retry-btn" @click="vista === 'usuarios' ? obtenerUsuarios() : obtenerClientes()">Reintentar</button>
        </div>

        <!-- Tabla -->
        <div v-else class="table-wrapper">
          <div class="table-scroll">

            <!-- USUARIOS -->
            <table v-if="vista === 'usuarios'">
              <thead>
                <tr>
                  <th>Usuario</th>
                  <th>Contraseña</th>
                  <th>Rol</th>
                  <th class="th-actions">Acciones</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(user, i) in usuariosFiltrados" :key="user.id"
                    class="table-row" :style="{ animationDelay: i * 40 + 'ms' }">
                  <td class="td-user">
                    <div class="user-cell">
                      <div class="avatar">{{ user.username.charAt(0).toUpperCase() }}</div>
                      <span>{{ user.username }}</span>
                    </div>
                  </td>
                  <td class="td-pass"><span class="pass-dots">{{ maskPassword(user.password) }}</span></td>
                  <td><span class="role-badge" :class="rolColor(user.rol)">{{ user.rol }}</span></td>
                  <td class="td-actions">
                    <button class="action-btn edit-btn" title="Editar">
                      <svg viewBox="0 0 20 20" fill="none"><path d="M14.5 3.5l2 2L6 16l-3 1 1-3 11.5-10.5z" stroke="currentColor" stroke-width="1.3" stroke-linejoin="round"/></svg>
                    </button>
                    <button class="action-btn delete-btn" title="Eliminar">
                      <svg viewBox="0 0 20 20" fill="none"><path d="M5 7h10l-1 9H6L5 7zM3 7h14M8 7V5h4v2" stroke="currentColor" stroke-width="1.3" stroke-linecap="round" stroke-linejoin="round"/></svg>
                    </button>
                  </td>
                </tr>
                <tr v-if="!usuariosFiltrados.length">
                  <td colspan="4" class="empty-row">⚡ Sin resultados para "{{ busqueda }}"</td>
                </tr>
              </tbody>
            </table>

            <!-- CLIENTES -->
            <table v-else>
              <thead>
                <tr>
                  <th>Cédula</th>
                  <th>Nombre</th>
                  <th>Apellido</th>
                  <th>Teléfono</th>
                  <th>Contraseña</th>
                  <th class="th-actions">Planes</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(c, i) in clientesFiltrados" :key="c.idcustomer"
                    class="table-row" :style="{ animationDelay: i * 40 + 'ms' }">
                  <td>{{ c.cedula }}</td>
                  <td class="td-user">
                    <div class="user-cell">
                      <div class="avatar">{{ c.name?.charAt(0).toUpperCase() }}</div>
                      <span>{{ c.name }}</span>
                    </div>
                  </td>
                  <td>{{ c.lastname }}</td>
                  <td>{{ c.phone }}</td>
                  <td class="td-pass"><span class="pass-dots">{{ maskPassword(c.password) }}</span></td>
                  <td class="td-actions">
                    <!-- ✅ BOTÓN REDIRECCIÓN A PLANES -->
                    <button class="btn-plan" @click="verPlanes(c)" title="Ver y gestionar planes">
                      <svg viewBox="0 0 20 20" fill="none">
                        <rect x="3" y="3" width="14" height="14" rx="1.5" stroke="currentColor" stroke-width="1.4"/>
                        <path d="M7 10h6M10 7v6" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/>
                      </svg>
                      Ver Planes
                    </button>
                  </td>
                </tr>
                <tr v-if="!clientesFiltrados.length">
                  <td colspan="6" class="empty-row">⚡ No hay clientes registrados</td>
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
  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Barlow+Condensed:wght@400;700;900&family=Barlow:wght@300;400;500&display=swap');

*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

.screen {
  min-height: 100vh;
  background: #0a0a0a;
  font-family: 'Barlow', sans-serif;
  position: relative;
  overflow: hidden;
}

.bg-layer { position: fixed; inset: 0; pointer-events: none; z-index: 0; }
.hex { position: absolute; }
.hex-1 {
  width: 420px; height: 420px;
  background: conic-gradient(from 30deg, rgba(245,197,0,.08) 0deg, rgba(255,122,0,.06) 80deg, transparent 80deg);
  clip-path: polygon(50% 0%,100% 25%,100% 75%,50% 100%,0% 75%,0% 25%);
  top: -100px; right: -80px;
}
.hex-2 {
  width: 260px; height: 260px;
  background: conic-gradient(from 200deg, rgba(245,197,0,.05) 0deg, transparent 80deg);
  clip-path: polygon(50% 0%,100% 25%,100% 75%,50% 100%,0% 75%,0% 25%);
  bottom: -60px; left: 220px;
}
.slash { position: absolute; }
.slash-1 { width: 2px; height: 60vh; background: linear-gradient(to bottom,transparent,rgba(245,197,0,.12) 50%,transparent); top: 20%; left: 220px; transform: rotate(12deg); }
.slash-2 { width: 1px; height: 40vh; background: linear-gradient(to bottom,transparent,rgba(255,122,0,.08) 50%,transparent); top: 30%; left: 230px; transform: rotate(12deg); }
.noise { position: absolute; inset: 0; background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='1'/%3E%3C/svg%3E"); opacity: .02; }

.layout { position: relative; z-index: 1; display: flex; min-height: 100vh; }

/* ── Sidebar ── */
.sidebar { width: 230px; min-height: 100vh; background: #0d0d0d; border-right: 1px solid rgba(245,197,0,.1); display: flex; flex-direction: column; padding: 32px 0 24px; flex-shrink: 0; position: sticky; top: 0; height: 100vh; }
.sidebar-logo { padding: 0 24px 28px; border-bottom: 1px solid rgba(245,197,0,.08); margin-bottom: 24px; }
.sidebar-logo svg { width: 130px; height: auto; filter: drop-shadow(0 0 18px rgba(245,197,0,.18)); }
.sidebar-sub { display: block; font-family: 'Barlow Condensed',sans-serif; font-size: .55rem; letter-spacing: .35em; text-transform: uppercase; color: #444; margin-top: 4px; padding-left: 2px; }

.sidebar-nav { display: flex; flex-direction: column; gap: 4px; padding: 0 12px; flex: 1; }
.nav-item { display: flex; align-items: center; gap: 10px; padding: 10px 14px; color: #555; font-size: .8rem; font-weight: 500; letter-spacing: .04em; text-decoration: none; border-left: 2px solid transparent; transition: all .2s; cursor: pointer; }
.nav-item svg { width: 16px; height: 16px; flex-shrink: 0; }
.nav-item:hover { color: #f0f0f0; background: rgba(245,197,0,.05); }
.nav-item.active { color: #f5c500; border-left-color: #f5c500; background: rgba(245,197,0,.07); }
.sidebar-footer {
  padding: 20px 16px 0; border-top: 1px solid rgba(245,197,0,.08);
  display: flex; align-items: center; gap: 10px; margin-top: 12px;
}

.user-chip { display: flex; align-items: center; gap: 10px; flex: 1; min-width: 0; }
.user-avatar {
  width: 32px; height: 32px;
  background: linear-gradient(135deg, #f5c500, #ff7a00);
  color: #0a0a0a; font-family: 'Barlow Condensed',sans-serif;
  font-weight: 900; font-size: .95rem;
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}
.user-info { display: flex; flex-direction: column; min-width: 0; }
.user-name { font-size: .78rem; color: #e0e0e0; font-weight: 500; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.user-role { font-size: .6rem; color: #555; letter-spacing: .08em; }

.logout-btn {
  background: none; border: 1px solid #1e1e1e; color: #444; cursor: pointer;
  width: 32px; height: 32px; display: flex; align-items: center; justify-content: center;
  flex-shrink: 0; transition: all .2s;
}
.logout-btn svg { width: 15px; height: 15px; }
.logout-btn:hover { border-color: rgba(245,197,0,.4); color: #f5c500; }

/* ── Main ── */
.main { flex: 1; padding: 70px 80px 60px; display: flex; flex-direction: column; gap: 24px; min-width: 0; }


.page-header { display: flex; align-items: flex-end; justify-content: space-between; flex-wrap: wrap; gap: 16px; }

.page-eyebrow { display: block; font-family: 'Barlow Condensed',sans-serif; font-size: .65rem; letter-spacing: .3em; text-transform: uppercase; color: #f5c500; margin-bottom: 6px; }

.page-title h1 { font-family: 'Barlow Condensed',sans-serif; font-size: 2.8rem; font-weight: 900; text-transform: uppercase; letter-spacing: .04em; color: #fff; line-height: 1; }

.header-meta { display: flex; gap: 12px; }
.stat-pill { display: flex; flex-direction: column; align-items: center; background: #111; border: 1px solid rgba(245,197,0,.12); padding: 10px 20px; gap: 2px; }
.stat-pill.accent { border-color: rgba(245,197,0,.3); background: rgba(245,197,0,.06); }
.stat-val { font-family: 'Barlow Condensed',sans-serif; font-size: 1.6rem; font-weight: 900; color: #f5c500; line-height: 1; }
.stat-lbl { font-size: .58rem; letter-spacing: .2em; text-transform: uppercase; color: #555; }

/* ── Filtros ── */
.filters-bar { display: flex; align-items: center; gap: 16px; flex-wrap: wrap; }

.search-wrap { position: relative; flex: 1; min-width: 240px; max-width: 400px; }
.search-icon { position: absolute; left: 13px; top: 50%; transform: translateY(-50%); width: 15px; height: 15px; color: #444; pointer-events: none; }
.search-input { width: 100%; background: #111; border: 1px solid rgba(245,197,0,.12); color: #f0f0f0; font-family: 'Barlow',sans-serif; font-size: .85rem; padding: 10px 14px 10px 38px; outline: none; transition: border-color .2s, background .2s; }
.search-input::placeholder { color: #333; }
.search-input:focus { border-color: rgba(245,197,0,.4); background: #0f0e09; }

.rol-filters { display: flex; gap: 8px; flex-wrap: wrap; }
.rol-chip { background: transparent; border: 1px solid rgba(245,197,0,.15); color: #555; font-family: 'Barlow Condensed',sans-serif; font-size: .75rem; font-weight: 700; letter-spacing: .15em; text-transform: uppercase; padding: 7px 16px; cursor: pointer; transition: all .2s; }
.rol-chip:hover { border-color: rgba(245,197,0,.4); color: #ccc; }
.rol-chip.active { background: rgba(245,197,0,.1); border-color: #f5c500; color: #f5c500; }

/* ── States ── */
.state-box { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 14px; padding: 80px 20px; color: #555; font-size: .9rem; letter-spacing: .05em; }
.bolt-spin { font-size: 2.4rem; animation: boltPulse 1.4s ease-in-out infinite; filter: drop-shadow(0 0 10px #f5c500); }
@keyframes boltPulse { 0%,100% { opacity: 1; } 50% { opacity: .3; } }
.error-state { color: #e05a45; }
.retry-btn { background: none; border: 1px solid rgba(220,50,30,.4); color: #e05a45; font-family: 'Barlow Condensed',sans-serif; font-size: .75rem; letter-spacing: .2em; text-transform: uppercase; padding: 8px 20px; cursor: pointer; transition: all .2s; }
.retry-btn:hover { background: rgba(220,50,30,.1); }

/* ── Tabla ── */
.table-wrapper { display: flex; flex-direction: column; border: 1px solid rgba(245,197,0,.12); background: #111; animation: fadeUp .5s cubic-bezier(.16,1,.3,1) both; }
@keyframes fadeUp { from { opacity: 0; transform: translateY(16px); } to { opacity: 1; transform: translateY(0); } }
.table-scroll { overflow-x: auto; }

table { width: 100%; border-collapse: collapse; font-size: .85rem; }
thead { background: rgba(245,197,0,.06); border-bottom: 1px solid rgba(245,197,0,.15); }
th { font-family: 'Barlow Condensed',sans-serif; font-size: .65rem; font-weight: 700; letter-spacing: .25em; text-transform: uppercase; color: #f5c500; padding: 14px 18px; text-align: left; }


.th-actions { text-align: center; }
td { padding: 14px 18px; border-bottom: 1px solid rgba(255,255,255,.04); color: #ccc; vertical-align: middle; }

.table-row { transition: background .15s; animation: rowIn .4s ease both; }
.table-row:hover { background: rgba(245,197,0,.04); }
.table-row:last-child td { border-bottom: none; }
@keyframes rowIn { from { opacity: 0; transform: translateX(-8px); } to { opacity: 1; transform: translateX(0); } }

.user-cell { display: flex; align-items: center; gap: 10px; }
.avatar { width: 32px; height: 32px; background: linear-gradient(135deg,rgba(245,197,0,.2),rgba(255,122,0,.15)); border: 1px solid rgba(245,197,0,.2); color: #f5c500; font-family: 'Barlow Condensed',sans-serif; font-weight: 900; font-size: .9rem; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.pass-dots { color: #333; letter-spacing: .15em; font-size: .8rem; }

.role-badge { display: inline-block; font-family: 'Barlow Condensed',sans-serif; font-size: .7rem; font-weight: 700; letter-spacing: .2em; text-transform: uppercase; padding: 4px 12px; border: 1px solid; }
.badge-admin   { color: #f5c500; border-color: rgba(245,197,0,.4); background: rgba(245,197,0,.07); }
.badge-customer { color: #ff7a00; border-color: rgba(255,122,0,.4); background: rgba(255,122,0,.07); }
.badge-default { color: #888; border-color: rgba(136,136,136,.3); background: rgba(136,136,136,.06); }

.td-actions { text-align: center; }
.action-btn { background: none; border: 1px solid transparent; width: 30px; height: 30px; display: inline-flex; align-items: center; justify-content: center; cursor: pointer; transition: all .2s; }
.action-btn svg { width: 14px; height: 14px; }
.edit-btn { color: #555; }
.edit-btn:hover { color: #f5c500; border-color: rgba(245,197,0,.4); background: rgba(245,197,0,.07); }
.delete-btn { color: #555; }
.delete-btn:hover { color: #e05a45; border-color: rgba(220,50,30,.4); background: rgba(220,50,30,.07); }

/* ── Botón Ver Planes ── */
.btn-plan {
  display: inline-flex; align-items: center; gap: 6px;
  background: transparent;
  border: 1px solid rgba(245,197,0,.3);
  color: #f5c500;
  font-family: 'Barlow Condensed', sans-serif;
  font-size: .72rem; font-weight: 700;
  letter-spacing: .18em; text-transform: uppercase;
  padding: 6px 14px;
  cursor: pointer;
  transition: all .2s;
  white-space: nowrap;
}
.btn-plan svg { width: 13px; height: 13px; flex-shrink: 0; }
.btn-plan:hover { background: rgba(245,197,0,.1); border-color: #f5c500; }

.empty-row { text-align: center; color: #444; padding: 48px 20px !important; font-size: .85rem; letter-spacing: .05em; }

.table-footer { padding: 12px 18px; border-top: 1px solid rgba(245,197,0,.08); font-size: .7rem; color: #444; letter-spacing: .06em; }
.table-footer strong { color: #888; }

@media (max-width: 0px) {
  .sidebar { display: none; }
  .main { padding: 28px 20px 48px; }
  .modal-grid { grid-template-columns: 1fr; }
}

</style>