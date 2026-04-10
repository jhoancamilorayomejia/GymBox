import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import DashboardProveedor from '../views/DashboardProveedor.vue'

// 🔍 Validar expiración del token
const tokenExpirado = (token) => {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.exp * 1000 < Date.now()
  } catch {
    return true
  }
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/api/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardView,
      meta: { requiresAuth: true } // 🔐 PROTEGIDA
    },
    {
      path: '/api/proyectos',
      name: 'dashboard-proveedor',
      component: DashboardProveedor,
      meta: { requiresAuth: true } // 🔐 PROTEGIDA
    },
    {
      path: '/',
      redirect: '/api/login'
    }
  ]
})


// 🔥 MIDDLEWARE GLOBAL (FRONTEND)
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')

  // 🔐 si la ruta requiere auth
  if (to.meta.requiresAuth) {

    // ❌ no hay token
    if (!token) {
      next('/api/login')
      return
    }

    // ❌ token expirado
    if (tokenExpirado(token)) {
      localStorage.removeItem('token')
      next('/api/login')
      return
    }
  }

  // ✅ todo bien
  next()
})

export default router