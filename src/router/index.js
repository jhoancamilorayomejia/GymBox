import { createRouter, createWebHistory } from 'vue-router'
import LoginView       from '../views/LoginView.vue'
import DashboardView   from '../views/DashboardView.vue'
import DashboardProveedor from '../views/DashboardProveedor.vue'
import DashboardPlan   from '../views/DashboardPlan.vue'

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
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardView,
      meta: { requiresAuth: true }
    },
    {
      // Planes de un cliente específico
      // Ej: /dashboard/planes/7
      path: '/dashboard/planes/:idcustomer',
      name: 'dashboard-plan',
      component: DashboardPlan,
      meta: { requiresAuth: true }
    },
    {
      path: '/proyectos',
      name: 'proyectos',
      component: DashboardProveedor,
      meta: { requiresAuth: true }
    },
    {
      path: '/',
      redirect: '/login'
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/login'
    }
  ]
})

// 🔐 Middleware global
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')

  if (to.meta.requiresAuth) {
    if (!token || tokenExpirado(token)) {
      localStorage.removeItem('token')
      return next('/login')
    }
  }

  next()
})

export default router