<template>
  <div class="screen">
    <div class="bg-grid"></div>
    <div class="glow left"></div>
    <div class="glow right"></div>

    <div class="card center-card">

      <div class="center-content">
        <span class="eyebrow">Acción completada</span>

        <h1>Reserva confirmada</h1>

        <p class="subtitle">
          Tu turno ha sido registrado correctamente.
        </p>

        <div class="divider">
          <span></span>
          <span class="diamond">◆</span>
          <span></span>
        </div>

        <!-- BOTÓN -->
        <button class="btn-primary center-btn" @click="volver">
          <span>Volver</span>
        </button>
      </div>

    </div>
  </div>
</template>

<script setup>
const volver = () => {
  console.log('click funcionando')
  window.location.href = 'http://localhost:3000/api/login'
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Cormorant+Garamond:wght@400;600;700&family=Montserrat:wght@300;400;500;600&display=swap');

*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

.screen {
  min-height: 100vh;
  background: #0d0d0d;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px 16px;
  font-family: 'Montserrat', sans-serif;
  position: relative;
  overflow: hidden;
}

/* 🔥 IMPORTANTE: no bloquear clicks */
.bg-grid,
.glow {
  pointer-events: none;
}

.bg-grid {
  position: fixed; inset: 0;
  background-image:
    linear-gradient(rgba(180,145,80,.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(180,145,80,.03) 1px, transparent 1px);
  background-size: 60px 60px;
}

.glow {
  position: fixed;
  width: 500px; height: 500px;
  border-radius: 50%;
  filter: blur(120px);
  opacity: .3;
}
.glow.left  { background: radial-gradient(circle, #b4915025, transparent 70%); top: -100px; left: -150px; }
.glow.right { background: radial-gradient(circle, #b4915015, transparent 70%); bottom: -100px; right: -150px; }

/* 🔥 IMPORTANTE: asegurar que esté encima */
.card {
  position: relative;
  z-index: 10;

  width: 100%;
  max-width: 500px;
  background: #141414;
  border: 1px solid rgba(180,145,80,.2);
  box-shadow: 0 32px 80px rgba(0,0,0,.7);
  padding: 60px 40px;
  text-align: center;
}

.center-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.eyebrow {
  font-size: .65rem;
  letter-spacing: .3em;
  text-transform: uppercase;
  color: #b49150;
}

h1 {
  font-family: 'Cormorant Garamond', serif;
  font-size: 2.3rem;
  color: #f0e6d0;
  letter-spacing: .08em;
}

.subtitle {
  color: #7a6a50;
  font-size: .85rem;
}

.divider {
  display: flex;
  align-items: center;
  gap: 10px;
}

.divider span:not(.diamond) {
  width: 60px;
  height: 1px;
  background: linear-gradient(90deg, transparent, #b49150);
}

.divider span:last-child {
  background: linear-gradient(90deg, #b49150, transparent);
}

.diamond { font-size: .4rem; color: #b49150; }

.btn-primary {
  margin-top: 20px;
  padding: 14px 40px;
  border: 1px solid #b49150;
  background: transparent;
  color: #b49150;
  letter-spacing: .3em;
  text-transform: uppercase;
  cursor: pointer;
  transition: all .3s;
}

.btn-primary:hover {
  background: #b49150;
  color: #0d0d0d;
}
</style>