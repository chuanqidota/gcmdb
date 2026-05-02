import { createRouter, createWebHistory } from 'vue-router'
import Layout from '../layout/index.vue'

const routes = [
  {
    path: '/',
    component: Layout,
    redirect: '/home',
    children: [
      { path: 'home', name: 'Home', component: () => import('../views/home/index.vue'), meta: { title: '首页' } },
      { path: 'model-manage', name: 'ModelManage', component: () => import('../views/model-manage/index.vue'), meta: { title: '模型管理' } },
      { path: 'instance', name: 'Instance', component: () => import('../views/instance/index.vue'), meta: { title: '实例管理' } },
      { path: 'search-direct-sql', name: 'SearchDirectSql', component: () => import('../views/search-direct-sql/index.vue'), meta: { title: 'SQL 查询' } },
      { path: 'search', name: 'Search', component: () => import('../views/search/index.vue'), meta: { title: '综合检索' } },
      { path: 'model-topology', name: 'ModelTopology', component: () => import('../views/model-topology/index.vue'), meta: { title: '模型拓扑' } },
      { path: 'audit', name: 'Audit', component: () => import('../views/audit/index.vue'), meta: { title: '审计日志' } },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
