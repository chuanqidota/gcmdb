<template>
  <div class="inst-relation-panel" v-loading="cache?.loading">
    <template v-if="cache?.groups?.length">
      <div class="inst-relation-header">
        <el-tabs v-model="cache.activeTab" @tab-change="onTabChange" class="inst-relation-tabs">
          <el-tab-pane v-for="g in cache.groups" :key="g.model_id" :name="String(g.model_id)">
            <template #label>
              <span class="tab-label">
                <span :class="['tab-direction', g.direction]">{{ g.direction === 'upstream' ? '←' : '→' }}</span>
                {{ g.model_name }}
                <el-badge :value="g.instances.length" class="tab-badge" />
              </span>
            </template>
          </el-tab-pane>
        </el-tabs>
        <el-button size="small" type="primary" @click="emit('viewTopo', instanceId)" style="margin-left: 12px; flex-shrink: 0;">查看拓扑图</el-button>
      </div>
      <template v-for="g in cache.groups" :key="'c-' + g.model_id">
        <div v-if="String(g.model_id) === cache.activeTab">
          <el-table
            :data="paginated"
            size="small"
            stripe
            highlight-current-row
            v-loading="cache.tabData?.[g.model_id]?.loading"
          >
            <el-table-column prop="instance_id" label="ID" width="60">
              <template #default="{ row: rel }">{{ rel.instance_id }}</template>
            </el-table-column>
            <el-table-column
              v-for="f in (cache.tabData?.[g.model_id]?.fields || [])"
              :key="f.alias"
              :label="f.name"
              show-overflow-tooltip
              min-width="100"
            >
              <template #default="{ row: rel }">{{ rel.data?.[f.alias] ?? '-' }}</template>
            </el-table-column>
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row: rel }">
                <el-button link type="primary" size="small" @click="emit('openDetail', rel.instance_id, rel.model_name, rel.model_alias)">查看</el-button>
                <el-button link type="primary" size="small" @click="handleViewTopo(rel)">拓扑</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="inst-relation-pagination">
            <el-pagination
              v-model:current-page="page"
              v-model:page-size="pageSize"
              :total="total"
              :page-sizes="[10, 20, 50]"
              layout="total, sizes, prev, pager, next"
              small
              background
              @current-change="syncPage"
              @size-change="syncSize"
            />
          </div>
        </div>
      </template>
    </template>
    <template v-else-if="!cache?.loading">
      <div class="inst-relation-empty">暂无关联实例</div>
    </template>
  </div>
</template>

<script>
import { computed, ref, watch } from 'vue'

export default {
  name: 'RelationPanel',
  emits: ['viewTopo', 'openDetail'],
  props: {
    instanceId: { type: Number, required: true },
    cache: { type: Object, default: undefined },
    tabChangeHandler: { type: Function, default: null },
  },
  setup(props, { emit }) {
    const page = ref(1)
    const pageSize = ref(10)

    const currentTab = computed(() => props.cache?.activeTab || '')

    watch(currentTab, (tab) => {
      page.value = 1
      const td = props.cache?.tabData?.[Number(tab)]
      pageSize.value = td?.pageSize || 10
    })

    const tabData = computed(() => {
      const td = props.cache?.tabData?.[currentTab.value]
      return td || null
    })

    const total = computed(() => tabData.value?.instances?.length || 0)

    const paginated = computed(() => {
      const td = tabData.value
      if (!td) return []
      const sz = pageSize.value
      const start = (page.value - 1) * sz
      return td.instances.slice(start, start + sz)
    })

    function syncPage() {
      if (tabData.value) tabData.value.page = page.value
    }

    function syncSize(size) {
      pageSize.value = size
      page.value = 1
      if (tabData.value) {
        tabData.value.pageSize = size
        tabData.value.page = 1
      }
    }

    function onTabChange(tab) {
      const cache = props.cache
      if (!cache) return
      const modelId = Number(tab)
      cache.activeTab = tab
      if (props.tabChangeHandler) {
        props.tabChangeHandler(props.instanceId, modelId)
      }
    }

    function handleViewTopo(rel) {
      emit('viewTopo', rel.instance_id, rel.model_id, rel.model_name, rel.model_alias)
    }

    return { page, pageSize, total, paginated, emit, onTabChange, syncPage, syncSize, handleViewTopo }
  },
}
</script>
