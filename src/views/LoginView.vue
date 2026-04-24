<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const showModal = ref(false)
const loginEmail    = ref('')
const loginPassword = ref('')
const loginError    = ref('')
const loginLoading  = ref(false)
const showPassword  = ref(false)
const mobileMenuOpen = ref(false)

const loginAdmin = async () => {
  loginLoading.value = true
  loginError.value   = ''
  try {
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: loginEmail.value, password: loginPassword.value })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Usuario o contraseña incorrectos')
    localStorage.setItem('token',    data.token)
    localStorage.setItem('rol',      data.rol)
    localStorage.setItem('username', data.user)
    localStorage.setItem('id',       data.id)
    showModal.value = false
    if (data.rol === 'admin')         router.push('/dashboard')
    else if (data.rol === 'customer') router.push('/proyectos')
    else                              router.push('/login')
  } catch (err) {
    loginError.value = err.message
  } finally {
    loginLoading.value = false
  }
}

/* ── Planes dinámicos ── */
const plansLoading = ref(true)
const plansError   = ref('')
const plans        = ref([])

const PLAN_META = {
  DIA:    { tag: 'Pase único',  highlight: false, weight: 1, perks: ['Acceso completo por 1 día','Todas las zonas de entrenamiento','Sin compromiso'] },
  SEMANA: { tag: 'Más popular', highlight: true,  weight: 2, perks: ['Acceso ilimitado 7 días','Todas las zonas de entrenamiento','Pago digital disponible','Visualiza horarios en plataforma'] },
  MES:    { tag: 'Recomendado', highlight: false, weight: 3, perks: ['Acceso ilimitado 30 días','Todas las zonas de entrenamiento','Pago digital disponible','Visualiza horarios en plataforma'] },
  ANO:    { tag: 'Mejor valor', highlight: false, weight: 4, perks: ['Acceso ilimitado 365 días','Todas las zonas de entrenamiento','Pago digital disponible','Visualiza horarios en plataforma','Ahorrás vs mensual'] },
}

const normalizeKey = (str) => str.toUpperCase().normalize('NFD').replace(/[\u0300-\u036f]/g,'').replace(/\s+/g,'')
const formatPrice  = (p)   => new Intl.NumberFormat('es-CO',{minimumFractionDigits:0}).format(p)

onMounted(async () => {
  try {
    const res  = await fetch('/api/login-price_plans')
    if (!res.ok) throw new Error('No se pudieron cargar los planes')
    const data = await res.json()
    plans.value = data
      .map(p => {
        const key  = normalizeKey(p.typeplan)
        const meta = PLAN_META[key] || { tag:'', highlight:false, weight:99, perks:[] }
        return { ...p, ...meta, name: p.typeplan }
      })
      .sort((a,b) => a.weight - b.weight)
  } catch (e) {
    plansError.value = e.message
  } finally {
    plansLoading.value = false
  }
})

/* ── Disciplinas ── */
const disciplines = [
  {
    name: 'CrossFit',
    desc: 'Entrenamiento funcional de alta intensidad basado en movimientos compuestos y varianza constante.',
    img:  'https://images.unsplash.com/photo-1517963879433-6ad2b056d712?w=700&q=80',
    alt:  'Atleta en plataforma de levantamiento olímpico con barra cargada'
  },
  {
    name: 'Entrenamiento Funcional',
    desc: 'Metodología orientada a patrones de movimiento naturales, mejorando fuerza, coordinación y resistencia.',
    img:  'https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=700&q=80',
    alt:  'Persona realizando entrenamiento funcional'
  },
  {
    name: 'Powerlifting',
    desc: 'Disciplina de fuerza máxima centrada en sentadilla, press de banca y peso muerto con progresión estructurada.',
    img:  'https://images.unsplash.com/photo-1534438327276-14e5300c3a48?w=700&q=80',
    alt:  'Atleta ejecutando press de banca en banco plano'
  },
  {
    name: 'Sala de Máquinas',
    desc: 'Zona equipada con maquinaria de resistencia para hipertrofia, rehabilitación y entrenamiento individualizado.',
    img:  'https://images.unsplash.com/photo-1574680096145-d05b474e2155?w=700&q=80',
    alt:  'Persona entrenando en máquina hack squat'
  }
]

const benefits = [
  { icon: '📅', title: 'Disponibilidad en tiempo real',  desc: 'Consultá desde la plataforma qué equipos y horarios están disponibles antes de ir a entrenar.' },
  { icon: '💳', title: 'Pago 100% digital',              desc: 'Abonás tu plan de forma segura con Mercado Pago desde la app o web. Sin filas, sin efectivo.' },
  { icon: '⚡', title: 'Gestión desde cualquier lugar',  desc: 'Accedé a tu historial, renovaciones y estado de membresía desde donde estés, 24/7.' }
]

const showRegisterModal = ref(false)
const regName      = ref('')
const regLastName  = ref('')
const regCedula    = ref('')
const regPhone     = ref('')
const regPassword  = ref('')
const regPassword2 = ref('')
const regError     = ref('')
const regLoading   = ref(false)
const showRegPassword = ref(false)

const registerCustomer = async () => {
  regError.value = ''
  if (regPassword.value !== regPassword2.value) {
    regError.value = 'Las contraseñas no coinciden'
    return
  }
  regLoading.value = true
  try {
    const res = await fetch('/api/new-public-customers', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: regName.value,
        lastname: regLastName.value,
        cedula: regCedula.value,
        phone: regPhone.value,
        password: regPassword.value
      })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Error al registrarse')
    showRegisterModal.value = false
    showModal.value = true
  } catch (err) {
    regError.value = err.message
  } finally {
    regLoading.value = false
  }
}
</script>

<template>
  <div class="site">

    <!-- NAVBAR -->
    <nav class="navbar">
      <div class="nav-inner">
        <div class="nav-logo">
          <svg viewBox="0 0 180 80" xmlns="http://www.w3.org/2000/svg" class="logo-svg">
            <polygon points="52,4 36,40 48,40 32,76 72,32 56,32 78,4" fill="#f5c500"/>
            <text x="78" y="38" font-family="'Barlow Condensed',sans-serif" font-weight="900" font-size="34" fill="#f5c500" letter-spacing="-1">RAYO</text>
            <text x="78" y="68" font-family="'Barlow Condensed',sans-serif" font-weight="900" font-size="34" fill="#fff" letter-spacing="-1">BOX</text>
            <text x="78" y="78" font-family="'Barlow Condensed',sans-serif" font-weight="400" font-size="9" fill="#888" letter-spacing="3">CROSS LIFTING</text>
          </svg>
        </div>
        <div class="nav-links">
          <a href="#sede">Sede</a>
          <a href="#planes">Planes</a>
          <a href="#disciplinas">Disciplinas</a>
          <a href="#plataforma">Plataforma</a>
          <button class="btn-nav-login" @click="showModal = true">Iniciar Sesión</button>
        </div>
        <!-- Mobile nav -->
        <div class="nav-mobile">
          <button class="btn-nav-login" @click="showModal = true">Ingresar</button>
          <button class="hamburger" @click="mobileMenuOpen = !mobileMenuOpen" :class="{ open: mobileMenuOpen }">
            <span></span><span></span><span></span>
          </button>
        </div>
      </div>
      <!-- Mobile dropdown menu -->
      <transition name="menu-drop">
        <div v-if="mobileMenuOpen" class="mobile-menu">
          <a href="#sede"        @click="mobileMenuOpen = false">Sede</a>
          <a href="#planes"      @click="mobileMenuOpen = false">Planes</a>
          <a href="#disciplinas" @click="mobileMenuOpen = false">Disciplinas</a>
          <a href="#plataforma"  @click="mobileMenuOpen = false">Plataforma</a>
          <button class="mobile-menu-register" @click="mobileMenuOpen = false; showRegisterModal = true">Registrarse</button>
        </div>
      </transition>
    </nav>

    <!-- HERO -->
    <header class="hero">
      <div class="hero-bg">
        <div class="hero-geo geo-1"></div>
        <div class="hero-geo geo-2"></div>
        <div class="hero-geo geo-3"></div>
        <div class="hero-noise"></div>
      </div>
      <div class="hero-content">
        <span class="hero-eyebrow">⚡ Bienvenido a</span>
        <div class="hero-logo-img-wrap">
          <img src="@/assets/rayo.png" alt="Rayo Box Cross Lifting" class="hero-logo-img" />
        </div>
        <p class="hero-sub">CrossFit · Funcional · Powerlifting · Máquinas<br/>Sede La Paila, Valle del Cauca</p>
        <div class="hero-ctas">
          <a href="#planes" class="btn-primary">Ver planes</a>
          <button class="btn-secondary" @click="showModal = true">Iniciar Sesión</button>
          <button class="btn-secondary" @click="showRegisterModal = true">Registrarse</button>
        </div>
      </div>
      <div class="hero-badge">
        <span class="badge-num">+4</span>
        <span class="badge-txt">Disciplinas disponibles</span>
      </div>
    </header>

    <!-- SEDE -->
    <section class="sede-section" id="sede">
      <div class="section-inner">
        <div class="section-label">📍 Ubicación</div>
        <h2 class="section-title">Sede <span class="accent">La Paila</span>, Valle del Cauca</h2>
        <p class="section-sub">Encontrá nuestro box y comenzá tu transformación hoy mismo.</p>
        <div class="sede-card">
          <div class="sede-info">
            <div class="sede-detail"><span class="sd-icon">📍</span><div><strong>Dirección</strong><p>La Paila, Valle del Cauca, Colombia</p></div></div>
            <div class="sede-detail"><span class="sd-icon">⏰</span><div><strong>Horarios</strong><p>Lunes a Sábado: 5:00am – 10:00pm</p></div></div>
            <div class="sede-detail"><span class="sd-icon">📞</span><div><strong>Contacto</strong><p>Escribinos por WhatsApp</p></div></div>
            <a href="#" class="btn-sede">Ver en el mapa</a>
          </div>
          <div class="sede-map-placeholder">
            <div class="map-inner">
              <div class="map-pin">📍</div>
              <span>Rayo Box — La Paila</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- PLANES -->
    <section class="planes-section" id="planes">
      <div class="section-inner">
        <div class="section-label">💪 Membresías</div>
        <h2 class="section-title">Elegí tu <span class="accent">plan</span></h2>
        <p class="section-sub">Accedé a Rayo Box con la modalidad que mejor se adapte a tu ritmo de entrenamiento.</p>

        <div v-if="plansLoading" class="plans-loading">
          <div class="loading-spinner"></div>
          <span>Cargando planes...</span>
        </div>
        <div v-else-if="plansError" class="plans-error"><span>⚠️ {{ plansError }}</span></div>
        <div v-else class="planes-grid">
          <div v-for="plan in plans" :key="plan.idpriceplan" class="plan-card" :class="{ featured: plan.highlight }">
            <div v-if="plan.highlight" class="plan-badge">{{ plan.tag }}</div>
            <div v-else class="plan-tag-small">{{ plan.tag }}</div>
            <div class="plan-name">{{ plan.name }}</div>
            <div class="plan-price">
              <span class="price-currency">$</span>
              <span class="price-num">{{ formatPrice(plan.price) }}</span>
            </div>
            <ul class="plan-perks">
              <li v-for="perk in plan.perks" :key="perk"><span class="perk-check">✓</span> {{ perk }}</li>
            </ul>
            <button class="btn-secondary" @click="showRegisterModal = true">Inscribirme</button>
          </div>
        </div>

        <p class="planes-note">* Los precios pueden variar según disponibilidad y promociones vigentes.</p>
      </div>
    </section>

    <!-- DISCIPLINAS -->
    <section class="disciplinas-section" id="disciplinas">
      <div class="section-inner">
        <div class="section-label">🏋️ Modalidades</div>
        <h2 class="section-title">Metodologías de <span class="accent">alto rendimiento</span></h2>
        <p class="section-sub">En Rayo Box integramos disciplinas complementarias para un desarrollo físico integral y progresivo.</p>
        <div class="disc-grid">
          <div class="disc-card" v-for="d in disciplines" :key="d.name">
            <div class="disc-img-wrap">
              <img :src="d.img" :alt="d.alt" class="disc-img" loading="lazy" />
              <div class="disc-img-overlay"></div>
            </div>
            <div class="disc-body">
              <div class="disc-line"></div>
              <h3 class="disc-name">{{ d.name }}</h3>
              <p class="disc-desc">{{ d.desc }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- PLATAFORMA -->
    <section class="plataforma-section" id="plataforma">
      <div class="section-inner plat-inner">
        <div class="plat-text">
          <div class="section-label">⚡ Plataforma digital</div>
          <h2 class="section-title">Registrate y gestioná<br/><span class="accent">tu membresía online</span></h2>
          <p class="section-sub">Más que un box, una plataforma completa para que tengas el control total de tu entrenamiento.</p>
          <div class="benefits-list">
            <div class="benefit-item" v-for="b in benefits" :key="b.title">
              <div class="benefit-icon">{{ b.icon }}</div>
              <div>
                <strong>{{ b.title }}</strong>
                <p>{{ b.desc }}</p>
              </div>
            </div>
          </div>
          <button class="btn-primary" @click="showModal = true">Acceder a la plataforma</button>
        </div>

        <div class="plat-visual">
          <div class="screen-mock">
            <div class="mock-bar"><span></span><span></span><span></span></div>
            <div class="mock-content">
              <div class="mock-layout">

                <div class="mock-sidebar">
                  <div class="mock-logo-small">⚡ RAYO<br/><span>BOX</span></div>
                  <div class="mock-nav-item active"><span class="nav-icon">☰</span> Mis Planes</div>
                  <div class="mock-plan-btn">
                    <span class="nav-icon">＋</span>
                    <div class="mock-plan-btn-text">
                      <strong>NUEVA<br/>SUSCRIPCIÓN</strong>
                      <small>Renueva o agrega</small>
                    </div>
                  </div>
                  <div class="mock-user-pill">
                    <div class="mock-user-num">1</div>
                    <div class="mock-user-info">
                      <span>Usuario</span>
                      <small>ID #1 · Miembro</small>
                    </div>
                  </div>
                </div>

                <div class="mock-main">
                  <div class="mock-welcome">⚡ BIENVENIDO DE VUELTA</div>
                  <div class="mock-uid">Usuario <span class="mock-uid-tag">#1</span></div>

                  <div class="mock-counters">
                    <div class="mock-counter-box">
                      <div class="mock-counter-num">1</div>
                      <div class="mock-counter-lbl">TOTAL</div>
                    </div>
                    <div class="mock-counter-box active-box">
                      <div class="mock-counter-num">1</div>
                      <div class="mock-counter-lbl">VIGENTES</div>
                    </div>
                  </div>

                  <div class="mock-plan-active">
                    <div class="mock-plan-header">
                      <span class="mock-bolt">⚡</span>
                      <span class="mock-badge-active">ACCESO ACTIVO</span>
                      <span class="mock-plan-type">SEMANA</span>
                    </div>
                    <p class="mock-plan-desc">Acceso válido hasta el 28 de abril de 2026. Quedan 6 días.</p>
                    <div class="mock-plan-dates">
                      <div><small>INICIO</small><strong>20/04/2026</strong></div>
                      <div class="mock-arrow">→</div>
                      <div><small>FIN</small><strong>28/04/2026</strong></div>
                      <div class="mock-plan-price-tag"><small>PRECIO</small><strong class="accent">$20.000</strong></div>
                    </div>
                  </div>

                  <div class="mock-calendar">
                    <div class="mock-cal-header">
                      <span class="mock-cal-nav">‹</span>
                      <strong>ABRIL 2026</strong>
                      <span class="mock-cal-nav">›</span>
                    </div>
                    <div class="mock-cal-days">
                      <span v-for="d in ['D','L','M','M','J','V','S']" :key="d" class="mock-cal-dow">{{d}}</span>
                    </div>
                    <div class="mock-cal-grid">
                      <span v-for="i in 2" :key="'e'+i" class="mock-cal-cell empty"></span>
                      <span v-for="n in 19" :key="n" class="mock-cal-cell muted">{{ n }}</span>
                      <span class="mock-cal-cell avail">20</span>
                      <span class="mock-cal-cell avail">21</span>
                      <span class="mock-cal-cell today">22</span>
                      <span class="mock-cal-cell avail">23</span>
                      <span class="mock-cal-cell avail">24</span>
                      <span class="mock-cal-cell avail">25</span>
                      <span class="mock-cal-cell avail">26</span>
                      <span class="mock-cal-cell avail">27</span>
                      <span class="mock-cal-cell avail">28</span>
                      <span class="mock-cal-cell muted">29</span>
                      <span class="mock-cal-cell muted">30</span>
                    </div>
                    <div class="mock-cal-legend">
                      <span><span class="legend-dot avail"></span>Disponible</span>
                      <span><span class="legend-dot venc"></span>Vencido</span>
                      <span><span class="legend-dot noacc"></span>Sin acceso</span>
                    </div>
                  </div>

                  <div class="mock-table-section">
                    <div class="mock-table-title">⊞ PLANES VIGENTES</div>
                    <div class="mock-table-row header-row"><span>#</span><span>TIPO</span><span>PRECIO</span><span>ESTADO</span></div>
                    <div class="mock-table-row">
                      <span class="mock-id-badge">18</span>
                      <span>SEMANA</span>
                      <span class="accent">$20.000</span>
                      <span class="mock-badge-estado active">ACTIVO</span>
                    </div>
                  </div>

                  <div class="mock-mp-btn">
                    <svg viewBox="0 0 60 24" class="mp-logo" xmlns="http://www.w3.org/2000/svg">
                      <ellipse cx="12" cy="12" rx="10" ry="10" fill="#009ee3"/>
                      <path d="M8 12 a4 4 0 0 1 8 0" stroke="#fff" stroke-width="2" fill="none"/>
                      <circle cx="12" cy="8" r="1.5" fill="#fff"/>
                    </svg>
                    Pagar con Mercado Pago
                  </div>
                </div>

              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- FOOTER -->
    <footer class="footer">
      <div class="footer-inner">
        <div class="footer-logo">
          <svg viewBox="0 0 160 70" xmlns="http://www.w3.org/2000/svg" style="width:120px">
            <polygon points="46,3 30,36 42,36 26,67 64,28 48,28 70,3" fill="#f5c500"/>
            <text x="70" y="33" font-family="'Barlow Condensed',sans-serif" font-weight="900" font-size="30" fill="#f5c500">RAYO</text>
            <text x="70" y="60" font-family="'Barlow Condensed',sans-serif" font-weight="900" font-size="30" fill="#fff">BOX</text>
          </svg>
          <p>Sistema de Gestión v1.0<br/>La Paila, Valle del Cauca</p>
        </div>
        <div class="footer-links">
          <a href="#sede">Sede</a>
          <a href="#planes">Planes</a>
          <a href="#disciplinas">Disciplinas</a>
          <a @click="showModal = true" style="cursor:pointer">Iniciar Sesión</a>
        </div>
      </div>
      <div class="footer-copy">© 2025 Rayo Box Cross Lifting · Todos los derechos reservados</div>
    </footer>

    <!-- MODAL LOGIN -->
    <transition name="modal-fade">
      <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
        <div class="modal-card">
          <button class="modal-close" @click="showModal = false">✕</button>

          <div class="modal-brand">
            <svg viewBox="0 0 180 80" xmlns="http://www.w3.org/2000/svg" style="width:140px">
              <polygon points="52,4 36,40 48,40 32,76 72,32 56,32 78,4" fill="#f5c500"/>
              <text x="78" y="38" font-family="'Barlow Condensed',sans-serif" font-weight="900" font-size="34" fill="#f5c500" letter-spacing="-1">RAYO</text>
              <text x="78" y="68" font-family="'Barlow Condensed',sans-serif" font-weight="900" font-size="34" fill="#fff" letter-spacing="-1">BOX</text>
              <text x="78" y="78" font-family="'Barlow Condensed',sans-serif" font-weight="400" font-size="9" fill="#555" letter-spacing="3">CROSS LIFTING</text>
            </svg>
          </div>

          <div class="modal-header">
            <div class="modal-header-top">
              <div class="bolt-badge">⚡</div>
              <h1>Acceso al Sistema</h1>
            </div>
            <p>Ingresá tus credenciales para continuar</p>
          </div>
          <div class="modal-divider"></div>

          <form class="login-form" @submit.prevent="loginAdmin">
            <div class="field" :class="{ filled: loginEmail }">
              <label>Usuario</label>
              <div class="input-wrap">
                <svg class="input-icon" viewBox="0 0 20 20" fill="none">
                  <circle cx="10" cy="7" r="3.5" stroke="currentColor" stroke-width="1.5"/>
                  <path d="M3 17c0-3.314 3.134-6 7-6s7 2.686 7 6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                </svg>
                <input v-model="loginEmail" type="text" placeholder="Tu nombre de usuario" autocomplete="username" required/>
              </div>
            </div>
            <div class="field" :class="{ filled: loginPassword }">
              <label>Contraseña</label>
              <div class="input-wrap">
                <svg class="input-icon" viewBox="0 0 20 20" fill="none">
                  <rect x="4" y="9" width="12" height="9" rx="1.5" stroke="currentColor" stroke-width="1.5"/>
                  <path d="M7 9V6.5a3 3 0 016 0V9" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                </svg>
                <input v-model="loginPassword" :type="showPassword ? 'text' : 'password'" placeholder="••••••••" autocomplete="current-password" required/>
                <button type="button" class="eye-btn" @click="showPassword = !showPassword">
                  <svg v-if="!showPassword" viewBox="0 0 20 20" fill="none">
                    <path d="M2 10s3-6 8-6 8 6 8 6-3 6-8 6-8-6-8-6z" stroke="currentColor" stroke-width="1.4"/>
                    <circle cx="10" cy="10" r="2.5" stroke="currentColor" stroke-width="1.4"/>
                  </svg>
                  <svg v-else viewBox="0 0 20 20" fill="none">
                    <path d="M3 3l14 14M8.5 8.7A2.5 2.5 0 0013 12.3M6 6.1C3.9 7.4 2 10 2 10s3 6 8 6c1.6 0 3-.5 4.2-1.3M10 4c4.2.5 6.6 3.7 8 6" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/>
                  </svg>
                </button>
              </div>
            </div>
            <transition name="shake">
              <div v-if="loginError" class="error-box"><span>⚡</span> {{ loginError }}</div>
            </transition>
            <button class="btn-login" type="submit" :disabled="loginLoading">
              <span class="btn-label">
                <svg v-if="loginLoading" class="spinner" viewBox="0 0 24 24">
                  <circle cx="12" cy="12" r="10" stroke="rgba(0,0,0,.25)" stroke-width="3" fill="none"/>
                  <path d="M12 2a10 10 0 0 1 10 10" stroke="#0a0a0a" stroke-width="3" stroke-linecap="round" fill="none"/>
                </svg>
                {{ loginLoading ? 'Ingresando...' : 'Ingresar al sistema' }}
              </span>
            </button>
          </form>

          <p class="modal-switch-link">
            ¿No tenés cuenta?
            <button @click="showModal = false; showRegisterModal = true">Registrate gratis</button>
          </p>
          <p class="modal-footer-note">Rayo Box · Sistema de Gestión v1.0</p>
        </div>
      </div>
    </transition>

    <!-- MODAL REGISTRO -->
    <transition name="modal-fade">
      <div v-if="showRegisterModal" class="modal-overlay" @click.self="showRegisterModal = false">
        <div class="modal-card">
          <button class="modal-close" @click="showRegisterModal = false">✕</button>

          <div class="modal-brand">
            <svg viewBox="0 0 180 80" xmlns="http://www.w3.org/2000/svg" style="width:140px">
              <polygon points="52,4 36,40 48,40 32,76 72,32 56,32 78,4" fill="#f5c500"/>
              <text x="78" y="38" font-family="'Barlow Condensed',sans-serif" font-weight="900" font-size="34" fill="#f5c500" letter-spacing="-1">RAYO</text>
              <text x="78" y="68" font-family="'Barlow Condensed',sans-serif" font-weight="900" font-size="34" fill="#fff" letter-spacing="-1">BOX</text>
              <text x="78" y="78" font-family="'Barlow Condensed',sans-serif" font-weight="400" font-size="9" fill="#555" letter-spacing="3">CROSS LIFTING</text>
            </svg>
          </div>

          <div class="modal-header">
            <div class="modal-header-top">
              <div class="bolt-badge">⚡</div>
              <h1>Crear cuenta</h1>
            </div>
            <p>Completá tus datos para unirte a Rayo Box</p>
          </div>
          <div class="modal-divider"></div>

          <form class="login-form" @submit.prevent="registerCustomer">
            <div class="reg-grid">
              <div class="field" :class="{ filled: regName }">
                <label>Nombre</label>
                <div class="input-wrap">
                  <svg class="input-icon" viewBox="0 0 20 20" fill="none">
                    <circle cx="10" cy="7" r="3.5" stroke="currentColor" stroke-width="1.5"/>
                    <path d="M3 17c0-3.314 3.134-6 7-6s7 2.686 7 6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                  </svg>
                  <input v-model="regName" type="text" placeholder="Juan" required />
                </div>
              </div>
              <div class="field" :class="{ filled: regLastName }">
                <label>Apellido</label>
                <div class="input-wrap">
                  <svg class="input-icon" viewBox="0 0 20 20" fill="none">
                    <circle cx="10" cy="7" r="3.5" stroke="currentColor" stroke-width="1.5"/>
                    <path d="M3 17c0-3.314 3.134-6 7-6s7 2.686 7 6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                  </svg>
                  <input v-model="regLastName" type="text" placeholder="Pérez" required />
                </div>
              </div>
            </div>

            <div class="reg-grid">
              <div class="field" :class="{ filled: regCedula }">
                <label>Cédula</label>
                <div class="input-wrap">
                  <svg class="input-icon" viewBox="0 0 20 20" fill="none">
                    <rect x="3" y="4" width="14" height="12" rx="1.5" stroke="currentColor" stroke-width="1.5"/>
                    <path d="M6 9h8M6 12h5" stroke="currentColor" stroke-width="1.3" stroke-linecap="round"/>
                  </svg>
                  <input v-model="regCedula" type="text" placeholder="1234567890" required />
                </div>
              </div>
              <div class="field" :class="{ filled: regPhone }">
                <label>Teléfono</label>
                <div class="input-wrap">
                  <svg class="input-icon" viewBox="0 0 20 20" fill="none">
                    <path d="M4 4a1 1 0 011-1h2.5a1 1 0 011 .8l.5 2.5a1 1 0 01-.3.9l-1 1a8 8 0 004.1 4.1l1-1a1 1 0 01.9-.3l2.5.5a1 1 0 01.8 1V16a1 1 0 01-1 1C7.2 17 3 12.8 3 5a1 1 0 011-1z" stroke="currentColor" stroke-width="1.4"/>
                  </svg>
                  <input v-model="regPhone" type="tel" placeholder="3001234567" required />
                </div>
              </div>
            </div>

            <div class="field" :class="{ filled: regPassword }">
              <label>Contraseña</label>
              <div class="input-wrap">
                <svg class="input-icon" viewBox="0 0 20 20" fill="none">
                  <rect x="4" y="9" width="12" height="9" rx="1.5" stroke="currentColor" stroke-width="1.5"/>
                  <path d="M7 9V6.5a3 3 0 016 0V9" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                </svg>
                <input v-model="regPassword" :type="showRegPassword ? 'text' : 'password'" placeholder="Mínimo 8 caracteres" required />
                <button type="button" class="eye-btn" @click="showRegPassword = !showRegPassword">
                  <svg v-if="!showRegPassword" viewBox="0 0 20 20" fill="none">
                    <path d="M2 10s3-6 8-6 8 6 8 6-3 6-8 6-8-6-8-6z" stroke="currentColor" stroke-width="1.4"/>
                    <circle cx="10" cy="10" r="2.5" stroke="currentColor" stroke-width="1.4"/>
                  </svg>
                  <svg v-else viewBox="0 0 20 20" fill="none">
                    <path d="M3 3l14 14M8.5 8.7A2.5 2.5 0 0013 12.3M6 6.1C3.9 7.4 2 10 2 10s3 6 8 6c1.6 0 3-.5 4.2-1.3M10 4c4.2.5 6.6 3.7 8 6" stroke="currentColor" stroke-width="1.4" stroke-linecap="round"/>
                  </svg>
                </button>
              </div>
            </div>

            <div class="field" :class="{ filled: regPassword2 }">
              <label>Confirmar contraseña</label>
              <div class="input-wrap">
                <svg class="input-icon" viewBox="0 0 20 20" fill="none">
                  <rect x="4" y="9" width="12" height="9" rx="1.5" stroke="currentColor" stroke-width="1.5"/>
                  <path d="M7 9V6.5a3 3 0 016 0V9" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                </svg>
                <input v-model="regPassword2" :type="showRegPassword ? 'text' : 'password'" placeholder="Repetí tu contraseña" required />
              </div>
            </div>

            <transition name="shake">
              <div v-if="regError" class="error-box"><span>⚡</span> {{ regError }}</div>
            </transition>

            <button class="btn-login" type="submit" :disabled="regLoading">
              <span class="btn-label">
                <svg v-if="regLoading" class="spinner" viewBox="0 0 24 24">
                  <circle cx="12" cy="12" r="10" stroke="rgba(0,0,0,.25)" stroke-width="3" fill="none"/>
                  <path d="M12 2a10 10 0 0 1 10 10" stroke="#0a0a0a" stroke-width="3" stroke-linecap="round" fill="none"/>
                </svg>
                {{ regLoading ? 'Creando cuenta...' : 'Crear mi cuenta' }}
              </span>
            </button>
          </form>

          <p class="modal-switch-link">
            ¿Ya tenés cuenta?
            <button @click="showRegisterModal = false; showModal = true">Iniciá sesión</button>
          </p>
          <p class="modal-footer-note">Rayo Box · Sistema de Gestión v1.0</p>
        </div>
      </div>
    </transition>

  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Barlow+Condensed:wght@400;700;900&family=Barlow:wght@300;400;500;600&display=swap');
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
.site{font-family:'Barlow',sans-serif;background:#0a0a0a;color:#e8e8e8;overflow-x:hidden}

/* ── NAVBAR ── */
.navbar{position:fixed;top:0;left:0;right:0;z-index:100;background:rgba(10,10,10,.92);backdrop-filter:blur(12px);border-bottom:1px solid rgba(245,197,0,.12)}
.nav-inner{max-width:1200px;margin:0 auto;padding:0 32px;height:64px;display:flex;align-items:center;justify-content:space-between}
.logo-svg{width:130px;height:auto}
.nav-links{display:flex;align-items:center;gap:28px}
.nav-links a{color:#aaa;font-size:.82rem;letter-spacing:.1em;text-transform:uppercase;text-decoration:none;transition:color .2s}
.nav-links a:hover{color:#f5c500}
.btn-nav-login{background:#f5c500;border:none;color:#0a0a0a;font-family:'Barlow Condensed',sans-serif;font-size:.85rem;font-weight:900;letter-spacing:.2em;text-transform:uppercase;padding:9px 22px;cursor:pointer;transition:background .2s,transform .15s}
.btn-nav-login:hover{background:#ffd633;transform:translateY(-1px)}

/* Mobile nav elements — hidden on desktop */
.nav-mobile{display:none;align-items:center;gap:12px}
.hamburger{background:none;border:none;cursor:pointer;padding:6px;display:flex;flex-direction:column;gap:5px;width:32px}
.hamburger span{display:block;height:2px;background:#f5c500;border-radius:2px;transition:transform .25s, opacity .25s}
.hamburger.open span:nth-child(1){transform:translateY(7px) rotate(45deg)}
.hamburger.open span:nth-child(2){opacity:0}
.hamburger.open span:nth-child(3){transform:translateY(-7px) rotate(-45deg)}

/* Mobile dropdown */
.mobile-menu{background:rgba(10,10,10,.97);border-bottom:1px solid rgba(245,197,0,.12);display:flex;flex-direction:column;padding:12px 20px 20px}
.mobile-menu a{color:#aaa;font-size:.9rem;letter-spacing:.1em;text-transform:uppercase;text-decoration:none;padding:13px 0;border-bottom:1px solid #1a1a1a;transition:color .2s}
.mobile-menu a:hover{color:#f5c500}
.mobile-menu-register{margin-top:16px;background:#f5c500;border:none;color:#0a0a0a;font-family:'Barlow Condensed',sans-serif;font-size:.9rem;font-weight:900;letter-spacing:.2em;text-transform:uppercase;padding:13px;cursor:pointer;width:100%}
.menu-drop-enter-active,.menu-drop-leave-active{transition:opacity .2s, transform .2s}
.menu-drop-enter-from,.menu-drop-leave-to{opacity:0;transform:translateY(-8px)}

/* ── HERO ── */
.hero{min-height:100vh;display:flex;align-items:center;padding:100px 32px 80px;position:relative;overflow:hidden}
.hero-bg{position:absolute;inset:0;pointer-events:none}
.hero-geo{position:absolute}
.geo-1{width:520px;height:520px;background:conic-gradient(from 30deg,rgba(245,197,0,.18) 0deg,rgba(255,122,0,.12) 90deg,transparent 90deg);clip-path:polygon(50% 0%,100% 25%,100% 75%,50% 100%,0% 75%,0% 25%);top:-160px;right:-120px;animation:rotateSlow 40s linear infinite}
.geo-2{width:340px;height:340px;background:conic-gradient(from 200deg,rgba(245,197,0,.1) 0deg,transparent 70deg);clip-path:polygon(50% 0%,100% 25%,100% 75%,50% 100%,0% 75%,0% 25%);bottom:-80px;left:-60px;animation:rotateSlow 28s linear infinite reverse}
.geo-3{width:2px;height:80vh;background:linear-gradient(to bottom,transparent,rgba(245,197,0,.25) 50%,transparent);top:10%;left:55%;transform:rotate(12deg)}
.hero-noise{position:absolute;inset:0;background-image:url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='1'/%3E%3C/svg%3E");opacity:.025}
@keyframes rotateSlow{to{transform:rotate(360deg)}}
.hero-content{max-width:1200px;margin:0 auto;position:relative;z-index:1;animation:slideUp .8s cubic-bezier(.16,1,.3,1) both}
@keyframes slideUp{from{opacity:0;transform:translateY(32px)}to{opacity:1;transform:translateY(0)}}
.hero-eyebrow{display:inline-block;font-size:.7rem;letter-spacing:.35em;text-transform:uppercase;color:#f5c500;margin-bottom:10px}
.hero-logo-img-wrap{margin-bottom:2px;max-width:850px}
.hero-logo-img{width:100%;height:auto;display:block;border-radius:4px;filter:drop-shadow(0 0 32px rgba(245,197,0,.3))}
.accent{color:#f5c500}
.hero-sub{font-size:1rem;color:#666;line-height:1.7;margin-bottom:40px;letter-spacing:.04em}
.hero-ctas{display:flex;gap:16px;flex-wrap:wrap}
.btn-primary{background:#f5c500;color:#0a0a0a;font-family:'Barlow Condensed',sans-serif;font-size:.9rem;font-weight:900;letter-spacing:.25em;text-transform:uppercase;padding:15px 36px;border:none;cursor:pointer;text-decoration:none;display:inline-block;transition:background .2s,transform .15s}
.btn-primary:hover{background:#ffd633;transform:translateY(-2px)}
.btn-secondary{background:transparent;border:1.5px solid rgba(245,197,0,.5);color:#f5c500;font-family:'Barlow Condensed',sans-serif;font-size:.9rem;font-weight:900;letter-spacing:.25em;text-transform:uppercase;padding:15px 36px;cursor:pointer;transition:border-color .2s,background .2s,color .2s}
.btn-secondary:hover{background:rgba(245,197,0,.08);border-color:#f5c500}
.hero-badge{position:absolute;bottom:40px;right:32px;display:flex;flex-direction:column;align-items:flex-end;z-index:1}
.badge-num{font-family:'Barlow Condensed',sans-serif;font-size:3.5rem;font-weight:900;color:#f5c500;line-height:1}
.badge-txt{font-size:.7rem;color:#555;letter-spacing:.1em;text-transform:uppercase}

/* ── SECTIONS ── */
.section-inner{max-width:1200px;margin:0 auto;padding:0 32px}
.section-label{font-size:.65rem;letter-spacing:.35em;text-transform:uppercase;color:#f5c500;margin-bottom:12px}
.section-title{font-family:'Barlow Condensed',sans-serif;font-size:clamp(2rem,4vw,3rem);font-weight:900;text-transform:uppercase;color:#fff;margin-bottom:12px;line-height:1.1}
.section-sub{font-size:.9rem;color:#666;max-width:560px;line-height:1.7;margin-bottom:48px}

/* ── SEDE ── */
.sede-section{padding:100px 0;background:#0d0d0d;border-top:1px solid rgba(245,197,0,.08);border-bottom:1px solid rgba(245,197,0,.08)}
.sede-card{display:grid;grid-template-columns:1fr 1fr;gap:40px;background:#111;border:1px solid rgba(245,197,0,.1);padding:40px}
.sede-info{display:flex;flex-direction:column;gap:24px}
.sede-detail{display:flex;gap:14px;align-items:flex-start}
.sd-icon{font-size:1.3rem;margin-top:2px}
.sede-detail strong{display:block;font-size:.8rem;letter-spacing:.1em;text-transform:uppercase;color:#f5c500;margin-bottom:4px}
.sede-detail p{font-size:.85rem;color:#888}
.btn-sede{display:inline-block;margin-top:8px;background:#f5c500;color:#0a0a0a;font-family:'Barlow Condensed',sans-serif;font-size:.8rem;font-weight:900;letter-spacing:.2em;text-transform:uppercase;padding:12px 28px;text-decoration:none;cursor:pointer;border:none;width:fit-content;transition:background .2s}
.btn-sede:hover{background:#ffd633}
.sede-map-placeholder{background:#0d0d0d;border:1px solid #1e1e1e;display:flex;align-items:center;justify-content:center;min-height:220px}
.map-inner{text-align:center}
.map-pin{font-size:2.5rem;margin-bottom:12px;display:block}
.map-inner span{font-size:.75rem;letter-spacing:.15em;text-transform:uppercase;color:#555}

/* ── PLANES ── */
.planes-section{padding:100px 0}
.plans-loading{display:flex;align-items:center;gap:16px;padding:48px 0;color:#555;font-size:.85rem;letter-spacing:.1em}
.loading-spinner{width:20px;height:20px;border:2px solid #222;border-top-color:#f5c500;border-radius:50%;animation:spin .7s linear infinite}
.plans-error{padding:20px;background:rgba(220,50,30,.08);border:1px solid rgba(220,50,30,.3);color:#e05a45;font-size:.82rem;margin-bottom:24px}
.planes-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:20px}
.plan-card{background:#111;border:1px solid #1e1e1e;padding:32px 28px;display:flex;flex-direction:column;gap:16px;position:relative;transition:border-color .25s,transform .2s}
.plan-card:hover{border-color:rgba(245,197,0,.3);transform:translateY(-4px)}
.plan-card.featured{border-color:#f5c500;background:linear-gradient(145deg,#141200,#111)}
.plan-badge{position:absolute;top:-12px;left:24px;background:#f5c500;color:#0a0a0a;font-family:'Barlow Condensed',sans-serif;font-size:.65rem;font-weight:900;letter-spacing:.2em;text-transform:uppercase;padding:4px 12px}
.plan-tag-small{font-size:.65rem;letter-spacing:.2em;text-transform:uppercase;color:#555}
.plan-name{font-family:'Barlow Condensed',sans-serif;font-size:1.8rem;font-weight:900;text-transform:uppercase;color:#fff}
.plan-price{display:flex;align-items:baseline;gap:4px;color:#f5c500}
.price-currency{font-size:.9rem;font-weight:600}
.price-num{font-family:'Barlow Condensed',sans-serif;font-size:2.4rem;font-weight:900;line-height:1}
.plan-perks{list-style:none;display:flex;flex-direction:column;gap:8px;flex:1}
.plan-perks li{font-size:.8rem;color:#888;display:flex;gap:8px;align-items:flex-start}
.perk-check{color:#f5c500;font-weight:700;flex-shrink:0}
.planes-note{font-size:.7rem;color:#444;margin-top:24px;text-align:center}

/* ── DISCIPLINAS ── */
.disciplinas-section{padding:100px 0;background:#0d0d0d;border-top:1px solid rgba(245,197,0,.08);border-bottom:1px solid rgba(245,197,0,.08)}
.disc-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(260px,1fr));gap:2px;background:#111;border:1px solid #1a1a1a}
.disc-card{background:#0d0d0d;display:flex;flex-direction:column;overflow:hidden;transition:background .2s}
.disc-card:hover{background:#111}
.disc-card:hover .disc-img{transform:scale(1.05)}
.disc-img-wrap{position:relative;height:200px;overflow:hidden}
.disc-img{width:100%;height:100%;object-fit:cover;display:block;transition:transform .5s ease;filter:brightness(.65) saturate(.8)}
.disc-img-overlay{position:absolute;inset:0;background:linear-gradient(to bottom,transparent 40%,#0d0d0d 100%);pointer-events:none}
.disc-body{padding:24px 28px 28px;flex:1}
.disc-line{width:32px;height:2px;background:#f5c500;margin-bottom:16px}
.disc-name{font-family:'Barlow Condensed',sans-serif;font-size:1.3rem;font-weight:900;text-transform:uppercase;color:#fff;margin-bottom:10px}
.disc-desc{font-size:.82rem;color:#666;line-height:1.65}

/* ── PLATAFORMA ── */
.plataforma-section{padding:100px 0}
.plat-inner{display:grid;grid-template-columns:1fr 1fr;gap:64px;align-items:start}
.benefits-list{display:flex;flex-direction:column;gap:24px;margin-bottom:36px}
.benefit-item{display:flex;gap:16px;align-items:flex-start}
.benefit-icon{font-size:1.4rem;flex-shrink:0;margin-top:2px}
.benefit-item strong{display:block;font-size:.82rem;letter-spacing:.06em;text-transform:uppercase;color:#f5c500;margin-bottom:4px}
.benefit-item p{font-size:.82rem;color:#666;line-height:1.6}

/* ── MOCK ── */
.screen-mock{background:#111;border:1px solid rgba(245,197,0,.15);border-radius:4px;box-shadow:0 40px 80px rgba(0,0,0,.6),0 0 40px rgba(245,197,0,.06);overflow:hidden}
.mock-bar{background:#0d0d0d;padding:8px 12px;display:flex;gap:6px;align-items:center;border-bottom:1px solid #1a1a1a}
.mock-bar span{width:8px;height:8px;border-radius:50%;background:#222}
.mock-bar span:first-child{background:rgba(245,197,0,.4)}
.mock-content{padding:0}
.mock-layout{display:flex;min-height:480px}
.mock-sidebar{width:130px;flex-shrink:0;background:#0a0a0a;border-right:1px solid #1a1a1a;padding:16px 12px;display:flex;flex-direction:column;gap:12px}
.mock-logo-small{font-family:'Barlow Condensed',sans-serif;font-weight:900;color:#f5c500;font-size:.9rem;line-height:1.1;margin-bottom:4px}
.mock-logo-small span{color:#fff}
.mock-nav-item{display:flex;align-items:center;gap:6px;font-size:.65rem;letter-spacing:.1em;text-transform:uppercase;color:#555;padding:6px 8px;cursor:pointer}
.mock-nav-item.active{background:rgba(245,197,0,.08);color:#f5c500;border-left:2px solid #f5c500}
.nav-icon{font-size:.8rem;flex-shrink:0}
.mock-plan-btn{display:flex;align-items:center;gap:8px;background:rgba(245,197,0,.1);border:1px solid rgba(245,197,0,.25);padding:8px;cursor:pointer}
.mock-plan-btn-text strong{display:block;font-family:'Barlow Condensed',sans-serif;font-size:.7rem;font-weight:900;color:#f5c500;text-transform:uppercase;line-height:1.1}
.mock-plan-btn-text small{font-size:.55rem;color:#555}
.mock-user-pill{margin-top:auto;display:flex;align-items:center;gap:6px;background:#111;border:1px solid #1e1e1e;padding:6px 8px}
.mock-user-num{width:20px;height:20px;background:#f5c500;color:#0a0a0a;font-family:'Barlow Condensed',sans-serif;font-weight:900;font-size:.75rem;display:flex;align-items:center;justify-content:center;flex-shrink:0}
.mock-user-info span{display:block;font-size:.58rem;color:#ccc}
.mock-user-info small{font-size:.52rem;color:#444}
.mock-main{flex:1;padding:16px;overflow:hidden;display:flex;flex-direction:column;gap:12px;position:relative}
.mock-welcome{font-size:.5rem;letter-spacing:.25em;color:#f5c500;text-transform:uppercase}
.mock-uid{font-family:'Barlow Condensed',sans-serif;font-size:1.3rem;font-weight:900;color:#fff;line-height:1}
.mock-uid-tag{color:#f5c500;font-size:1rem}
.mock-counters{display:flex;gap:8px;position:absolute;top:44px;right:16px}
.mock-counter-box{border:1px solid #1e1e1e;padding:6px 12px;text-align:center}
.mock-counter-box.active-box{border-color:#f5c500;background:rgba(245,197,0,.06)}
.mock-counter-num{font-family:'Barlow Condensed',sans-serif;font-size:1.1rem;font-weight:900;color:#fff}
.mock-counter-lbl{font-size:.45rem;letter-spacing:.15em;color:#555;text-transform:uppercase}
.mock-plan-active{background:#0d0d0d;border-left:3px solid #f5c500;padding:10px 12px;display:flex;flex-direction:column;gap:6px}
.mock-plan-header{display:flex;align-items:center;gap:8px;flex-wrap:wrap}
.mock-bolt{font-size:.8rem}
.mock-badge-active{background:rgba(245,197,0,.15);color:#f5c500;font-size:.5rem;font-weight:700;letter-spacing:.2em;padding:2px 8px;text-transform:uppercase}
.mock-plan-type{font-size:.55rem;color:#555;letter-spacing:.15em;text-transform:uppercase;margin-left:auto}
.mock-plan-desc{font-size:.6rem;color:#888}
.mock-plan-dates{display:flex;align-items:center;gap:10px;flex-wrap:wrap}
.mock-plan-dates>div{display:flex;flex-direction:column;gap:2px}
.mock-plan-dates small{font-size:.48rem;letter-spacing:.15em;color:#555;text-transform:uppercase}
.mock-plan-dates strong{font-size:.65rem;color:#ccc}
.mock-plan-price-tag strong{font-size:.7rem!important}
.mock-arrow{color:#333;font-size:.7rem}
.mock-calendar{background:#0a0a0a;border:1px solid #1a1a1a;padding:10px}
.mock-cal-header{display:flex;justify-content:space-between;align-items:center;margin-bottom:8px}
.mock-cal-header strong{font-family:'Barlow Condensed',sans-serif;font-size:.7rem;font-weight:900;letter-spacing:.15em;color:#ccc}
.mock-cal-nav{font-size:.8rem;color:#444;cursor:pointer;width:16px;height:16px;display:flex;align-items:center;justify-content:center;border:1px solid #1e1e1e}
.mock-cal-days{display:grid;grid-template-columns:repeat(7,1fr);gap:2px;margin-bottom:4px}
.mock-cal-dow{text-align:center;font-size:.45rem;letter-spacing:.1em;color:#444;text-transform:uppercase;padding:2px 0}
.mock-cal-grid{display:grid;grid-template-columns:repeat(7,1fr);gap:2px}
.mock-cal-cell{aspect-ratio:1;display:flex;align-items:center;justify-content:center;font-size:.55rem;color:#555}
.mock-cal-cell.avail{background:rgba(30,80,40,.6);color:#6dbf7e;font-weight:600}
.mock-cal-cell.today{background:#f5c500;color:#0a0a0a;font-weight:900}
.mock-cal-cell.muted{color:#2e2e2e}
.mock-cal-legend{display:flex;gap:10px;margin-top:6px;flex-wrap:wrap}
.mock-cal-legend span{display:flex;align-items:center;gap:4px;font-size:.48rem;color:#444;letter-spacing:.08em}
.legend-dot{width:7px;height:7px;display:inline-block}
.legend-dot.avail{background:rgba(30,80,40,.6)}
.legend-dot.venc{background:rgba(80,30,20,.6)}
.legend-dot.noacc{background:#1a1a1a;border:1px solid #2e2e2e}
.mock-table-title{font-size:.52rem;letter-spacing:.2em;color:#f5c500;text-transform:uppercase;margin-bottom:6px}
.mock-table-row{display:grid;grid-template-columns:28px 1fr 60px 60px;gap:6px;padding:5px 6px;font-size:.58rem;color:#888;align-items:center;border-bottom:1px solid #111}
.mock-table-row.header-row{font-size:.48rem;letter-spacing:.15em;color:#333;text-transform:uppercase;border-bottom:1px solid #1a1a1a}
.mock-id-badge{background:#111;border:1px solid #2a2a2a;color:#555;font-size:.5rem;font-weight:700;padding:2px 4px;text-align:center}
.mock-badge-estado{font-size:.48rem;font-weight:700;letter-spacing:.1em;padding:2px 6px;text-transform:uppercase}
.mock-badge-estado.active{color:#4caf50;border:1px solid rgba(76,175,80,.3)}
.mock-mp-btn{background:#f5c500;color:#0a0a0a;font-family:'Barlow Condensed',sans-serif;font-size:.65rem;font-weight:900;letter-spacing:.15em;text-align:center;padding:8px;cursor:pointer;display:flex;align-items:center;justify-content:center;gap:6px;text-transform:uppercase}
.mp-logo{width:18px;height:18px;flex-shrink:0}

/* ── FOOTER ── */
.footer{background:#060606;border-top:1px solid rgba(245,197,0,.1);padding:48px 0 24px}
.footer-inner{max-width:1200px;margin:0 auto;padding:0 32px 32px;display:flex;justify-content:space-between;align-items:flex-start;flex-wrap:wrap;gap:32px;border-bottom:1px solid #111}
.footer-logo p{font-size:.75rem;color:#444;margin-top:10px;line-height:1.6}
.footer-links{display:flex;flex-direction:column;gap:12px}
.footer-links a{color:#555;font-size:.78rem;letter-spacing:.1em;text-transform:uppercase;text-decoration:none;transition:color .2s}
.footer-links a:hover{color:#f5c500}
.footer-copy{max-width:1200px;margin:0 auto;padding:20px 32px 0;font-size:.68rem;color:#2a2a2a;letter-spacing:.1em;text-transform:uppercase}

/* ── MODAL ── */
.modal-overlay{position:fixed;inset:0;z-index:200;background:rgba(0,0,0,.88);backdrop-filter:blur(10px);display:flex;align-items:center;justify-content:center;padding:24px;overflow-y:auto}
.modal-fade-enter-active,.modal-fade-leave-active{transition:opacity .3s}
.modal-fade-enter-from,.modal-fade-leave-to{opacity:0}

/* Card principal — más ancha y con acento lateral dorado */
.modal-card{
  position:relative;width:100%;max-width:500px;
  background:#0f0f0f;
  border:1px solid rgba(245,197,0,.2);
  border-left:4px solid #f5c500;
  box-shadow:0 48px 120px rgba(0,0,0,.95),0 0 80px rgba(245,197,0,.08);
  padding:52px 48px 44px;
  animation:cardIn .4s cubic-bezier(.16,1,.3,1) both;
  margin:auto;
}
@keyframes cardIn{from{opacity:0;transform:scale(.96) translateY(20px)}to{opacity:1;transform:scale(1) translateY(0)}}

/* Botón cerrar */
.modal-close{
  position:absolute;top:18px;right:20px;
  background:rgba(255,255,255,.04);border:1px solid #1e1e1e;
  color:#555;font-size:.85rem;width:32px;height:32px;
  display:flex;align-items:center;justify-content:center;
  cursor:pointer;transition:background .2s,color .2s,border-color .2s;
  border-radius:2px;
}
.modal-close:hover{background:rgba(245,197,0,.1);border-color:rgba(245,197,0,.4);color:#f5c500}

/* Brand logo dentro del modal */
.modal-brand{margin-bottom:20px}

/* Header */
.modal-header{margin-bottom:36px}
.modal-header-top{display:flex;align-items:center;gap:14px;margin-bottom:8px}
.bolt-badge{
  width:44px;height:44px;flex-shrink:0;
  background:rgba(245,197,0,.1);border:1px solid rgba(245,197,0,.3);
  display:flex;align-items:center;justify-content:center;
  font-size:1.3rem;
  animation:boltPulse 2.5s ease-in-out infinite;
}
@keyframes boltPulse{0%,100%{border-color:rgba(245,197,0,.3)}50%{border-color:rgba(245,197,0,.7);background:rgba(245,197,0,.15)}}
.modal-header h1{
  font-family:'Barlow Condensed',sans-serif;
  font-size:2rem;font-weight:900;text-transform:uppercase;
  letter-spacing:.04em;color:#fff;line-height:1;margin:0;
}
.modal-header p{font-size:.82rem;color:#555;margin-top:6px;letter-spacing:.03em;padding-left:58px}

/* Divisor */
.modal-divider{height:1px;background:linear-gradient(to right,rgba(245,197,0,.3),transparent);margin-bottom:32px}

/* Formulario */
.login-form{display:flex;flex-direction:column;gap:22px}
.field{display:flex;flex-direction:column;gap:8px}
.field label{
  font-family:'Barlow Condensed',sans-serif;
  font-size:.7rem;letter-spacing:.28em;text-transform:uppercase;
  color:#444;font-weight:700;transition:color .2s;
}
.field.filled label,.field:focus-within label{color:#f5c500}

/* Input wrap con borde inferior animado */
.input-wrap{position:relative;display:flex;align-items:center}
.input-icon{position:absolute;left:14px;width:16px;height:16px;color:#333;pointer-events:none;transition:color .2s;flex-shrink:0}
.field:focus-within .input-icon{color:#f5c500}

.input-wrap input{
  width:100%;
  background:#080808;
  border:1px solid #1c1c1c;
  border-bottom:2px solid #1c1c1c;
  color:#f0f0f0;
  font-family:'Barlow',sans-serif;font-size:.95rem;
  padding:15px 16px 15px 44px;
  outline:none;
  transition:border-color .25s,border-bottom-color .25s,background .25s;
  border-radius:0;
}
.input-wrap input::placeholder{color:#252525;font-size:.88rem}
.input-wrap input:focus{
  border-color:#252525;
  border-bottom-color:#f5c500;
  background:#0b0a07;
}
.field.filled .input-wrap input{border-bottom-color:rgba(245,197,0,.4)}

/* Input sin icono izquierdo (registro) */
.input-wrap input.no-icon{padding-left:16px}

/* Eye button */
.eye-btn{position:absolute;right:12px;background:none;border:none;cursor:pointer;width:24px;height:24px;color:#2e2e2e;padding:0;display:flex;align-items:center;justify-content:center;transition:color .2s}
.eye-btn:hover{color:#f5c500}
.eye-btn svg{width:16px;height:16px}

/* Error */
.error-box{
  background:rgba(220,50,30,.06);
  border-left:3px solid #e05a45;
  color:#e05a45;font-size:.8rem;
  padding:11px 14px;display:flex;align-items:center;gap:10px;
}
.shake-enter-active{animation:shake .4s cubic-bezier(.36,.07,.19,.97)}
@keyframes shake{10%,90%{transform:translateX(-2px)}20%,80%{transform:translateX(3px)}30%,50%,70%{transform:translateX(-4px)}40%,60%{transform:translateX(4px)}}

/* Botón submit — relleno dorado desde el inicio */
.btn-login{
  position:relative;width:100%;padding:17px;
  background:#f5c500;
  border:none;
  color:#0a0a0a;
  font-family:'Barlow Condensed',sans-serif;font-size:1rem;
  font-weight:900;letter-spacing:.3em;text-transform:uppercase;
  cursor:pointer;overflow:hidden;margin-top:6px;
  display:flex;align-items:center;justify-content:center;
  transition:background .25s,transform .15s;
}
.btn-login::after{
  content:'';position:absolute;inset:0;
  background:rgba(255,255,255,0);
  transition:background .25s;
}
.btn-login:hover{background:#ffd633;transform:translateY(-1px)}
.btn-login:active{transform:translateY(0)}
.btn-login:disabled{opacity:.45;cursor:not-allowed;transform:none}
/* Eliminamos btn-bg ya que ahora usamos fondo sólido */
.btn-bg{display:none}
.btn-label{position:relative;z-index:1;display:flex;align-items:center;gap:8px}
.spinner{width:15px;height:15px;animation:spin .7s linear infinite;border:2px solid rgba(0,0,0,.2);border-top-color:#0a0a0a;border-radius:50%}
@keyframes spin{to{transform:rotate(360deg)}}

/* Nota footer */
.modal-footer-note{text-align:center;font-size:.65rem;color:#222;letter-spacing:.12em;text-transform:uppercase;margin-top:28px}

/* Link secundario (¿Ya tenés cuenta?) */
.modal-switch-link{text-align:center;font-size:.78rem;color:#444;margin-top:18px}
.modal-switch-link button{background:none;border:none;color:#f5c500;cursor:pointer;font-size:.78rem;font-weight:600;text-decoration:underline;padding:0;transition:color .2s}
.modal-switch-link button:hover{color:#ffd633}

/* ── REGISTRO GRID ── */
.reg-grid{display:grid;grid-template-columns:1fr 1fr;gap:14px}

/* ═══════════════════════════════════════
   RESPONSIVE — TABLET  (≤1000px)
═══════════════════════════════════════ */
@media(max-width:1000px){
  .plat-inner{grid-template-columns:1fr}
  .plat-visual{max-width:600px;margin:0 auto;width:100%}
}

/* ═══════════════════════════════════════
   RESPONSIVE — TABLET  (≤900px)
═══════════════════════════════════════ */
@media(max-width:900px){
  /* Navbar */
  .nav-links{display:none}
  .nav-mobile{display:flex}
  .nav-inner{padding:0 20px}

  /* Sede */
  .sede-card{grid-template-columns:1fr}

  /* Hero badge */
  .hero-badge{display:none}

  /* Sections */
  .sede-section,.planes-section,.disciplinas-section,.plataforma-section{padding:72px 0}
}

/* ═══════════════════════════════════════
   RESPONSIVE — MÓVIL  (≤600px)
═══════════════════════════════════════ */
@media(max-width:1600px){
  /* Navbar */
  .nav-inner{height:56px;padding:0 16px}
  .logo-svg{width:100px}
  .btn-nav-login{font-size:.75rem;padding:7px 14px;letter-spacing:.1em}

  /* Hero */
  .hero{padding:80px 16px 60px;min-height:auto}
  .hero-logo-img-wrap{max-width:100%}
  .hero-sub{font-size:.85rem;margin-bottom:28px}
  .hero-ctas{flex-direction:column;gap:12px}
  .btn-primary,.btn-secondary{width:100%;text-align:center;padding:14px 20px;font-size:.85rem}

  /* Sections padding & lateral */
  .section-inner{padding:0 16px}
  .sede-section,.planes-section,.disciplinas-section,.plataforma-section{padding:56px 0}
  .section-sub{margin-bottom:32px}

  /* Sede */
  .sede-card{padding:24px 16px;gap:20px}
  .sede-map-placeholder{min-height:160px}
  .btn-sede{width:100%;text-align:center}

  /* Planes — 1 columna en móvil */
  .planes-grid{grid-template-columns:1fr}
  .plan-card{padding:28px 20px}

  /* Disciplinas */
  .disc-grid{grid-template-columns:1fr}
  .disc-img-wrap{height:180px}

  /* Plataforma mock — ocultar sidebar para no aplastar */
  .mock-sidebar{display:none}
  .mock-layout{min-height:auto}
  .mock-counters{position:static;margin-top:4px}
  .mock-plan-dates{gap:6px}
  .screen-mock{font-size:90%}

  /* Footer */
  .footer-inner{flex-direction:column;gap:20px;padding:0 16px 24px}
  .footer-links{flex-direction:row;flex-wrap:wrap;gap:14px}
  .footer-copy{padding:16px 16px 0;text-align:center}

  /* Modal */
  .modal-card{padding:32px 20px 28px;border-left-width:3px}
  .modal-header h1{font-size:1.6rem}
  .modal-overlay{padding:12px;align-items:flex-start}
  .bolt-badge{width:36px;height:36px;font-size:1.1rem}
  .modal-header p{padding-left:0;margin-top:10px}

  /* Registro grid → 1 columna */
  .reg-grid{grid-template-columns:1fr}
}

/* ═══════════════════════════════════════
   RESPONSIVE — MÓVIL PEQUEÑO  (≤380px)
═══════════════════════════════════════ */
@media(max-width:1380px){
  .hero{padding:72px 12px 48px}
  .section-inner{padding:0 12px}
  .modal-card{padding:24px 14px 20px;border-left-width:3px}
  .plan-card{padding:22px 14px}
  .mock-cal-cell{font-size:.48rem}
}
</style>