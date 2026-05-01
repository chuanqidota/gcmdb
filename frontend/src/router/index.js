import { createRouter, createWebHistory } from 'vue-router'
import Layout from '../layout/index.vue'

const routes = [
  {
    path: '/',
    component: Layout,
    redirect: '/home',
    children: [
      { path: 'home', name: 'Home', component: () => import('../views/home/index.vue'), meta: { title: '首页' } },
      { path: 'model-group', name: 'ModelGroup', component: () => import('../views/model-group/index.vue'), meta: { title: '模型分组' } },
      { path: 'model', name: 'Model', component: () => import('../views/model/index.vue'), meta: { title: '模型列表' } },
      { path: 'model/:id', name: 'ModelDetail', component: () => import('../views/model/detail.vue'), meta: { title: '模型详情' } },
      { path: 'model-relation', name: 'ModelRelation', component: () => import('../views/model-relation/index.vue'), meta: { title: '模型关系' } },
      { path: 'model-relation-type', name: 'ModelRelationType', component: () => import('../views/model-relation-type/index.vue'), meta: { title: '关系类型' } },
      { path: 'instance', name: 'Instance', component: () => import('../views/instance/index.vue'), meta: { title: '实例管理' } },
      { path: 'search-direct-sql', name: 'SearchDirectSql', component: () => import('../views/search-direct-sql/index.vue'), meta: { title: 'SQL 查询' } },
      { path: 'search', name: 'Search', component: () => import('../views/search/index.vue'), meta: { title: '综合检索' } },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
