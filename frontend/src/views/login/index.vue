<template>
  <div class="login-page" @mousemove="onMouseMove" @mouseleave="onMouseLeave">
    <!-- Canvas 粒子网络 -->
    <canvas ref="canvasRef" class="login-canvas"></canvas>

    <!-- 鼠标聚光灯 -->
    <div class="spotlight" :style="spotlightStyle"></div>

    <el-card class="login-card" shadow="always" :style="cardStyle">
      <div class="login-header">
        <div class="logo-icon">
          <svg width="36" height="36" viewBox="0 0 28 28" fill="none">
            <rect width="28" height="28" rx="8" fill="url(#logoGrad)"/>
            <path d="M8 10h12M8 14h8M8 18h10" stroke="#fff" stroke-width="2" stroke-linecap="round"/>
            <defs>
              <linearGradient id="logoGrad" x1="0" y1="0" x2="28" y2="28">
                <stop stop-color="#2563EB"/>
                <stop offset="1" stop-color="#7C3AED"/>
              </linearGradient>
            </defs>
          </svg>
        </div>
        <h2 class="login-title">GCMDB</h2>
        <p class="login-subtitle">配置管理数据库</p>
      </div>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="0" @keyup.enter="handleLogin">
        <el-form-item prop="username">
          <el-input v-model="form.username" placeholder="用户名" size="large" prefix-icon="User" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="form.password" type="password" placeholder="密码" size="large" prefix-icon="Lock" show-password />
        </el-form-item>
        <el-form-item>
          <el-button class="login-btn" type="primary" size="large" :loading="loading" style="width: 100%" @click="handleLogin">登 录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { login } from '../../api/auth'

const router = useRouter()
const formRef = ref(null)
const loading = ref(false)
const form = ref({ username: '', password: '' })
const canvasRef = ref(null)

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

const handleLogin = async () => {
  await formRef.value.validate()
  loading.value = true
  try {
    await login(form.value)
    localStorage.setItem('gcmdb_logged_in', '1')
    ElMessage.success('登录成功')
    router.push('/')
  } catch {
    // error already shown by interceptor
  } finally {
    loading.value = false
  }
}

// === 鼠标跟踪 ===
const mouseX = ref(-1000)
const mouseY = ref(-1000)
const tiltX = ref(0)
const tiltY = ref(0)

const prefersReducedMotion = typeof window !== 'undefined' && window.matchMedia('(prefers-reduced-motion: reduce)').matches

function onMouseMove(e) {
  mouseX.value = e.clientX
  mouseY.value = e.clientY
  if (!prefersReducedMotion) {
    const x = (e.clientX / window.innerWidth - 0.5) * 2
    const y = (e.clientY / window.innerHeight - 0.5) * 2
    tiltX.value = y * -10
    tiltY.value = x * 15
  }
}

function onMouseLeave() {
  mouseX.value = -1000
  mouseY.value = -1000
  tiltX.value = 0
  tiltY.value = 0
}

const spotlightStyle = computed(() => ({
  background: `radial-gradient(600px circle at ${mouseX.value}px ${mouseY.value}px, rgba(37,99,235,0.12), transparent 40%)`,
}))

const cardStyle = computed(() => {
  if (prefersReducedMotion) return {}
  return {
    transform: `perspective(1000px) rotateX(${tiltX.value}deg) rotateY(${tiltY.value}deg)`,
  }
})

// === Canvas 粒子网络 ===
let animId = null
let particles = []

onMounted(() => {
  if (prefersReducedMotion) return
  const canvas = canvasRef.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')

  function resize() {
    canvas.width = window.innerWidth
    canvas.height = window.innerHeight
  }
  resize()
  window.addEventListener('resize', resize)

  // 创建粒子
  const count = 70
  particles = []
  for (let i = 0; i < count; i++) {
    particles.push({
      x: Math.random() * canvas.width,
      y: Math.random() * canvas.height,
      vx: (Math.random() - 0.5) * 0.6,
      vy: (Math.random() - 0.5) * 0.6,
      r: Math.random() * 2 + 2,
      baseAlpha: Math.random() * 0.4 + 0.4,
    })
  }

  function draw() {
    ctx.clearRect(0, 0, canvas.width, canvas.height)

    const mx = mouseX.value
    const my = mouseY.value

    for (let i = 0; i < count; i++) {
      const p = particles[i]

      // 鼠标吸引力
      if (mx > 0 && my > 0) {
        const dx = mx - p.x
        const dy = my - p.y
        const dist = Math.sqrt(dx * dx + dy * dy)
        if (dist < 200 && dist > 1) {
          const force = 0.02 * (1 - dist / 200)
          p.vx += (dx / dist) * force
          p.vy += (dy / dist) * force
        }
      }

      // 速度衰减
      p.vx *= 0.99
      p.vy *= 0.99

      p.x += p.vx
      p.y += p.vy

      // 边界反弹
      if (p.x < 0 || p.x > canvas.width) p.vx *= -1
      if (p.y < 0 || p.y > canvas.height) p.vy *= -1

      // 鼠标距离影响亮度
      let alpha = p.baseAlpha
      if (mx > 0 && my > 0) {
        const dx = mx - p.x
        const dy = my - p.y
        const dist = Math.sqrt(dx * dx + dy * dy)
        if (dist < 200) {
          alpha = Math.min(1, p.baseAlpha + 0.5 * (1 - dist / 200))
        }
      }

      // 绘制粒子（带发光）
      ctx.beginPath()
      ctx.arc(p.x, p.y, p.r, 0, Math.PI * 2)
      ctx.fillStyle = `rgba(96, 165, 250, ${alpha})`
      ctx.shadowBlur = 15
      ctx.shadowColor = 'rgba(59, 130, 246, 0.6)'
      ctx.fill()
      ctx.shadowBlur = 0

      // 绘制连线
      for (let j = i + 1; j < count; j++) {
        const q = particles[j]
        const dx = p.x - q.x
        const dy = p.y - q.y
        const dist = Math.sqrt(dx * dx + dy * dy)
        if (dist < 120) {
          const lineAlpha = 0.25 * (1 - dist / 120)
          ctx.beginPath()
          ctx.moveTo(p.x, p.y)
          ctx.lineTo(q.x, q.y)
          ctx.strokeStyle = `rgba(96, 165, 250, ${lineAlpha})`
          ctx.lineWidth = 1
          ctx.stroke()
        }
      }
    }

    animId = requestAnimationFrame(draw)
  }

  draw()

  onBeforeUnmount(() => {
    cancelAnimationFrame(animId)
    window.removeEventListener('resize', resize)
  })
})
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: #0F172A;
}

.login-canvas {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.spotlight {
  position: absolute;
  inset: 0;
  pointer-events: none;
  z-index: 1;
}

/* 毛玻璃卡片 */
.login-card {
  --el-card-bg-color: transparent;
  width: 400px;
  padding: 12px;
  background: rgba(15, 23, 42, 0.75);
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border: 1px solid rgba(96, 165, 250, 0.2);
  box-shadow:
    0 30px 80px rgba(0, 0, 0, 0.5),
    0 0 40px rgba(37, 99, 235, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.05);
  border-radius: 20px;
  animation: cardEnter 0.8s cubic-bezier(0.16, 1, 0.3, 1) both;
  transition: transform 0.15s ease-out;
  position: relative;
  z-index: 2;
}

.login-card::before {
  content: '';
  position: absolute;
  inset: -1px;
  border-radius: 21px;
  background: linear-gradient(135deg, rgba(96, 165, 250, 0.3), rgba(124, 58, 237, 0.1), rgba(96, 165, 250, 0.3));
  z-index: -1;
  opacity: 0.6;
  mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  mask-composite: exclude;
  -webkit-mask-composite: xor;
  padding: 1px;
}

@keyframes cardEnter {
  from {
    opacity: 0;
    transform: translateY(60px) scale(0.9);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.login-header {
  text-align: center;
  margin-bottom: 36px;
}

.logo-icon {
  display: inline-block;
  margin-bottom: 14px;
  animation: logoPulse 3s ease-in-out infinite;
  filter: drop-shadow(0 0 20px rgba(37, 99, 235, 0.5));
}

@keyframes logoPulse {
  0%, 100% {
    transform: scale(1);
    filter: drop-shadow(0 0 20px rgba(37, 99, 235, 0.5));
  }
  50% {
    transform: scale(1.08);
    filter: drop-shadow(0 0 35px rgba(37, 99, 235, 0.7));
  }
}

.login-title {
  font-size: 24px;
  font-weight: 700;
  color: #F1F5F9;
  margin: 0 0 6px 0;
  letter-spacing: 2px;
}

.login-subtitle {
  font-size: 13px;
  color: #64748B;
  margin: 0;
}

/* 输入框样式覆盖 */
.login-card :deep(.el-input__wrapper) {
  background: rgba(30, 41, 59, 0.8);
  border: 1px solid rgba(96, 165, 250, 0.15);
  box-shadow: none;
  border-radius: 10px;
}

.login-card :deep(.el-input__wrapper:hover) {
  border-color: rgba(96, 165, 250, 0.3);
}

.login-card :deep(.el-input__wrapper.is-focus) {
  border-color: rgba(96, 165, 250, 0.5);
  box-shadow: 0 0 15px rgba(37, 99, 235, 0.15);
}

.login-card :deep(.el-input__inner) {
  color: #E2E8F0;
}

.login-card :deep(.el-input__inner::placeholder) {
  color: #475569;
}

.login-card :deep(.el-input__prefix .el-icon) {
  color: #64748B;
}

/* 按钮微交互 */
.login-btn {
  border: none;
  border-radius: 10px;
  height: 44px;
  font-size: 15px;
  font-weight: 600;
  letter-spacing: 4px;
  background: linear-gradient(135deg, #2563EB, #7C3AED);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 30px rgba(37, 99, 235, 0.5);
}

.login-btn:active {
  transform: translateY(0);
}

/* 尊重用户偏好 */
@media (prefers-reduced-motion: reduce) {
  .login-card {
    animation: none;
  }
  .logo-icon {
    animation: none;
  }
  .login-btn:hover {
    transform: none;
  }
}
</style>
