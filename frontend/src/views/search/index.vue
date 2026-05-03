<template>
  <div class="search-page">
    <!-- Tab 栏 -->
    <div class="tab-bar">
      <span
        v-for="tab in tabs"
        :key="tab.name"
        class="tab-item"
        :class="{ active: activeTab === tab.name }"
        @click="switchTab(tab.name)"
      >{{ tab.label }}</span>
    </div>

    <!-- 搜索区 (接口调试 tab 时隐藏) -->
    <div v-show="activeTab !== 'api'" class="search-hero">
      <!-- 全文检索 -->
      <div v-if="activeTab === 'fulltext'" class="search-box">
        <el-input
          v-model="ft.keyword"
          placeholder="输入关键词搜索全部模型数据..."
          clearable
          @keyup.enter="doFulltextSearch"
          class="search-input"
        />
        <el-select v-model="ft.modelAlias" placeholder="全部模型" clearable filterable class="search-model-select">
          <el-option label="全部模型" value="" />
          <el-option v-for="m in allModels" :key="m.id" :label="m.name" :value="m.alias" />
        </el-select>
        <el-button type="primary" class="search-btn" @click="doFulltextSearch" :loading="ft.loading">搜索</el-button>
      </div>

      <!-- 实例搜索：模型选择 + 操作按钮（与全文检索同款居中布局） -->
      <div v-if="activeTab === 'instance'" class="instance-search-box">
        <el-select
          v-model="inst.modelId"
          placeholder="选择模型开始搜索..."
          filterable
          class="instance-model-select"
          @change="onModelChange"
        >
          <el-option
            v-for="m in allModels"
            :key="m.id"
            :label="`${m.name} (${m.alias})`"
            :value="m.id"
          />
        </el-select>
        <el-button type="primary" class="search-btn" @click="doInstanceSearch" :loading="inst.loading" :disabled="!inst.modelId">搜索</el-button>
        <el-button class="search-btn" @click="resetInstance">重置</el-button>
      </div>

      <!-- 模型浏览 -->
      <div v-if="activeTab === 'model'" class="search-box">
        <el-input
          v-model="modelSearch"
          placeholder="搜索模型名称、别名..."
          clearable
          class="search-input"
        >
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
      </div>
    </div>

    <!-- 结果区 (接口调试 tab 时隐藏) -->
    <div v-show="activeTab !== 'api'" class="results-area">

      <!-- 全文检索结果 -->
      <template v-if="activeTab === 'fulltext' && (ft.results.length > 0 || ft.loading)">
        <div class="result-count" v-if="ft.count > 0">找到 <strong>{{ ft.count }}</strong> 条结果</div>
        <el-skeleton :loading="ft.loading" animated :rows="5">
          <template #template>
            <div style="padding: 8px 0;">
              <el-skeleton-item v-for="i in 5" :key="i" variant="text" style="display: block; width: 100%; height: 40px; margin-bottom: 6px;" />
            </div>
          </template>
          <template #default>
            <el-table :data="ft.results" stripe size="small" highlight-current-row>
              <el-table-column prop="id" label="ID" width="80" />
              <el-table-column prop="model_name" label="模型" width="120">
                <template #default="{ row }"><el-tag size="small">{{ row.model_name }}</el-tag></template>
              </el-table-column>
              <el-table-column prop="model_alias" label="别名" width="100" />
              <el-table-column label="数据" min-width="300">
                <template #default="{ row }">
                  <span v-for="(v, k) in (row.data || {})" :key="k" class="data-tag">
                    <strong>{{ k }}:</strong> {{ v }}
                  </span>
                </template>
              </el-table-column>
              <template #empty>
                <el-empty description="输入关键词或选择模型开始搜索" />
              </template>
            </el-table>
          </template>
        </el-skeleton>
        <div v-if="ft.count > 0" class="pagination-wrap">
          <el-pagination
            v-model:current-page="ft.page"
            v-model:page-size="ft.pageSize"
            :total="ft.count"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            @current-change="doFulltextSearch"
            @size-change="doFulltextSearch"
          />
        </div>
      </template>

      <!-- 实例搜索：条件构建器 + 结果 -->
      <template v-if="activeTab === 'instance'">
        <!-- 条件构建器 -->
        <div v-if="inst.modelId" class="condition-builder">
          <div v-for="(c, i) in inst.conditions" :key="i" class="condition-row">
            <el-select v-model="c.field" placeholder="字段" class="cond-field">
              <el-option v-for="f in inst.fields" :key="f.alias" :value="f.alias">
                <span>{{ f.name }}</span>
                <el-tag v-if="f.is_indexed" size="small" type="primary" style="margin-left: 6px">索引</el-tag>
              </el-option>
            </el-select>
            <el-select v-model="c.op" class="cond-op">
              <el-option label="=" value="eq" />
              <el-option label="!=" value="ne" />
              <el-option label=">" value="gt" />
              <el-option label=">=" value="ge" />
              <el-option label="<" value="lt" />
              <el-option label="<=" value="le" />
              <el-option label="包含" value="contains" />
              <el-option label="开头是" value="startswith" />
              <el-option label="结尾是" value="endswith" />
              <el-option label="在列表" value="in" />
            </el-select>
            <el-input v-model="c.val" placeholder="值" class="cond-val" />
            <el-button :icon="Delete" circle @click="inst.conditions.splice(i, 1)" />
          </div>
          <div class="condition-actions">
            <el-button @click="inst.conditions.push({ field: '', op: 'eq', val: '' })">+ 添加条件</el-button>
          </div>
        </div>
        <div class="result-count" v-if="inst.count > 0">
          找到 <strong>{{ inst.count }}</strong> 条实例
          <el-button type="primary" size="small" style="margin-left: 12px" @click="openBatchTopo" :loading="inst.topoLoading">查看拓扑</el-button>
        </div>
        <el-skeleton :loading="inst.loading" animated :rows="5">
          <template #template>
            <div style="padding: 8px 0;">
              <el-skeleton-item v-for="i in 5" :key="i" variant="text" style="display: block; width: 100%; height: 40px; margin-bottom: 6px;" />
            </div>
          </template>
          <template #default>
            <el-table :data="inst.results" stripe size="small" highlight-current-row @expand-change="onInstExpandChange" row-key="id">
              <el-table-column type="expand" width="45">
                <template #default="{ row }">
                  <div class="inst-relation-panel" v-loading="instRelationCache[row.id]?.loading">
                    <template v-if="instRelationCache[row.id]?.groups?.length">
                      <div class="inst-relation-header">
                        <el-tabs v-model="instRelationCache[row.id].activeTab" @tab-change="(id) => onInstTabChange(row, Number(id))" class="inst-relation-tabs">
                          <el-tab-pane v-for="g in instRelationCache[row.id].groups" :key="g.model_id" :name="String(g.model_id)">
                            <template #label>
                              <span class="tab-label">
                                <span :class="['tab-direction', g.direction]">{{ g.direction === 'upstream' ? '←' : '→' }}</span>
                                {{ g.model_name }}
                                <el-badge :value="g.instances.length" class="tab-badge" />
                              </span>
                            </template>
                          </el-tab-pane>
                        </el-tabs>
                      </div>
                      <template v-for="g in instRelationCache[row.id].groups" :key="'content-' + g.model_id">
                        <div v-if="String(g.model_id) === instRelationCache[row.id].activeTab">
                          <el-table
                            :data="getInstPaginatedInstances(row.id, g.model_id)"
                            size="small"
                            stripe
                            highlight-current-row
                            v-loading="instRelationCache[row.id]?.tabData?.[g.model_id]?.loading"
                          >
                            <el-table-column prop="instance_id" label="ID" width="60">
                              <template #default="{ row: rel }">{{ rel.instance_id }}</template>
                            </el-table-column>
                            <el-table-column prop="relation_type" label="关系类型" width="120">
                              <template #default="{ row: rel }">
                                <el-tag size="small" type="primary">{{ rel.relation_type || '-' }}</el-tag>
                              </template>
                            </el-table-column>
                            <el-table-column
                              v-for="f in (instRelationCache[row.id]?.tabData?.[g.model_id]?.fields || [])"
                              :key="f.alias"
                              :label="f.name"
                              show-overflow-tooltip
                              min-width="100"
                            >
                              <template #default="{ row: rel }">{{ rel.data?.[f.alias] ?? '-' }}</template>
                            </el-table-column>
                          </el-table>
                          <div class="inst-relation-pagination">
                            <el-pagination
                              v-model:current-page="instRelationCache[row.id].tabData[g.model_id].page"
                              v-model:page-size="instRelationCache[row.id].tabData[g.model_id].pageSize"
                              :total="instRelationCache[row.id].tabData[g.model_id].instances.length"
                              :page-sizes="[10, 20, 50]"
                              layout="total, sizes, prev, pager, next"
                              small
                              background
                            />
                          </div>
                        </div>
                      </template>
                    </template>
                    <template v-else-if="!instRelationCache[row.id]?.loading">
                      <div class="inst-relation-empty">暂无关联实例</div>
                    </template>
                  </div>
                </template>
              </el-table-column>
              <el-table-column
                v-for="col in inst.columns"
                :key="col"
                :prop="col"
                :label="inst.columnLabelMap[col] || col"
                :min-width="calcColWidth(inst.columnLabelMap[col] || col)"
                show-overflow-tooltip
              />
              <el-table-column label="操作" width="100" fixed="right">
                <template #default="{ row }">
                  <el-button link type="primary" @click="openSingleTopo(row)">拓扑</el-button>
                </template>
              </el-table-column>
              <template #empty>
                <el-empty description="选择模型并设置条件后点击搜索" />
              </template>
            </el-table>
          </template>
        </el-skeleton>
        <div v-if="inst.count > 0" class="pagination-wrap">
          <el-pagination
            v-model:current-page="inst.page"
            v-model:page-size="inst.pageSize"
            :total="inst.count"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            @current-change="doInstanceSearch"
            @size-change="doInstanceSearch"
          />
        </div>
        <!-- 实例关系表格 -->
        <div v-if="instRelations.length" class="inst-relations-section">
          <div class="section-title">实例关系（{{ instRelations.length }} 条）</div>
          <el-table :data="instRelations" stripe size="small" highlight-current-row>
            <el-table-column prop="sourceLabel" label="源实例" min-width="150" />
            <el-table-column prop="relationType" label="关系类型" width="150">
              <template #default="{ row }">
                <el-tag size="small" type="primary">{{ row.relationType }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="targetLabel" label="目标实例" min-width="150" />
          </el-table>
        </div>
      </template>

      <!-- 模型浏览结果 -->
      <template v-if="activeTab === 'model'">
        <div v-for="group in groupedModels" :key="group.id" class="model-group-section">
          <div class="model-group-header">{{ group.name }}</div>
          <div class="model-list">
            <template v-for="m in group.models" :key="m.id">
              <div
                class="model-card"
                :class="{ expanded: expandedModelId === m.id }"
                @click="toggleFields(m)"
              >
                <div class="model-card-header">
                  <div class="model-card-info">
                    <span class="model-card-name">{{ m.name }}</span>
                    <el-tag size="small" type="info">{{ m.alias }}</el-tag>
                  </div>
                  <el-icon class="expand-icon" :class="{ rotated: expandedModelId === m.id }"><ArrowDown /></el-icon>
                </div>
                <div v-if="m.description" class="model-card-desc">{{ m.description }}</div>
              </div>
              <!-- 字段表内联展开（按分组） -->
              <div v-if="expandedModelId === m.id" class="fields-inline">
                <template v-if="modelFieldsCache[m.id]?.loading">
                  <el-skeleton :rows="5" animated>
                    <template #template>
                      <div style="padding: 8px 0;">
                        <el-skeleton-item v-for="i in 5" :key="i" variant="text" style="display: block; width: 100%; height: 40px; margin-bottom: 6px;" />
                      </div>
                    </template>
                  </el-skeleton>
                </template>
                <template v-else-if="modelFieldsCache[m.id]?.fieldGroups?.length">
                  <div v-for="group in modelFieldsCache[m.id].fieldGroups" :key="group.name" class="field-group">
                    <div class="field-group-title">{{ group.name || '未分组' }}</div>
                    <el-table :data="group.fields" stripe size="small" :show-header="group._showHeader" highlight-current-row>
                      <el-table-column prop="alias" label="别名" width="150">
                        <template #default="{ row }">
                          <span>{{ row.alias }}</span>
                          <el-tag v-if="modelFieldsCache[m.id]?.uniqueAliases?.includes(row.alias)" size="small" type="warning" class="unique-tag">唯一</el-tag>
                          <el-tag v-if="row.is_indexed" size="small" type="primary" class="unique-tag">索引</el-tag>
                        </template>
                      </el-table-column>
                      <el-table-column prop="name" label="名称" width="150" />
                      <el-table-column prop="type" label="类型" width="100">
                        <template #default="{ row }">
                          <el-tag size="small" :type="typeTag[row.type] || 'info'">{{ row.type }}</el-tag>
                        </template>
                      </el-table-column>
                      <el-table-column prop="is_required" label="必填" width="80" align="center">
                        <template #default="{ row }">
                          <el-icon v-if="row.is_required" color="var(--color-destructive)"><CircleCheckFilled /></el-icon>
                          <span v-else>-</span>
                        </template>
                      </el-table-column>
                      <el-table-column prop="order" label="排序" width="80" align="center" />
                      <template #empty>
                        <el-empty description="暂无字段数据" />
                      </template>
                    </el-table>
                  </div>
                </template>
                <template v-else>
                  <div class="fields-loading">暂无字段</div>
                </template>
                <!-- 模型关系 -->
                <div v-if="modelFieldsCache[m.id]?.relations?.length" class="model-relations">
                  <div class="relations-title">模型关系</div>
                  <div v-for="rel in modelFieldsCache[m.id].relations" :key="rel.id" class="relation-item">
                    <span class="relation-type">{{ resolveRelation(rel).typeLabel }}</span>
                    <span class="relation-target">{{ resolveRelation(rel).targetName }}</span>
                    <span v-if="resolveRelation(rel).description" class="relation-desc">{{ resolveRelation(rel).description }}</span>
                  </div>
                </div>
              </div>
            </template>
          </div>
        </div>
      </template>
    </div>

    <!-- 接口调试 -->
    <div v-if="activeTab === 'api'" class="api-test-layout">
      <!-- 左侧：接口选择 + 参数面板 -->
      <div class="api-left-panel">
        <div class="api-panel-title">OpenAPI 接口</div>
        <div class="api-endpoint-groups">
          <div v-for="group in apiEndpointGroups" :key="group.name" class="api-endpoint-group">
            <div class="api-group-title" @click="group.collapsed = !group.collapsed">
              <el-icon><ArrowRight v-if="group.collapsed" /><ArrowDown v-else /></el-icon>
              {{ group.name }}
            </div>
            <div v-show="!group.collapsed" class="api-group-items">
              <div
                v-for="ep in group.endpoints"
                :key="ep.key"
                :class="['api-endpoint-item', { active: apiSelectedKey === ep.key }]"
                @click="selectApiEndpoint(ep)"
              >
                <span class="api-method-tag" :class="ep.method.toLowerCase()">{{ ep.method }}</span>
                <span class="api-ep-name">{{ ep.name }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 参数表单 -->
        <div v-if="apiSelected" class="api-param-section">
          <div class="api-panel-title">参数</div>
          <div class="api-param-form">
            <div v-for="p in apiSelected.params" :key="p.key" class="api-param-row">
              <label class="api-param-label">
                {{ p.label }}
                <span v-if="p.required" class="api-required">*</span>
              </label>

              <el-input v-if="p.type === 'input'" v-model="apiParamValues[p.key]" :placeholder="p.placeholder || ''" size="small" />
              <el-input-number v-else-if="p.type === 'number'" v-model="apiParamValues[p.key]" :min="0" :placeholder="p.label" size="small" style="width: 100%" />
              <el-select v-else-if="p.type === 'select'" v-model="apiParamValues[p.key]" :placeholder="p.placeholder || '请选择'" size="small" style="width: 100%" filterable>
                <el-option v-for="opt in getApiOptions(p)" :key="opt.value" :label="opt.label" :value="opt.value" />
              </el-select>
              <el-select v-else-if="p.type === 'multi-select'" v-model="apiParamValues[p.key]" :placeholder="p.placeholder || '可多选'" size="small" style="width: 100%" multiple filterable collapse-tags collapse-tags-tooltip>
                <el-option v-for="opt in getApiOptions(p)" :key="opt.value" :label="opt.label" :value="opt.value" />
              </el-select>

              <!-- 多行文本 -->
              <el-input v-else-if="p.type === 'textarea'" v-model="apiParamValues[p.key]" type="textarea" :rows="4" :placeholder="p.placeholder || ''" size="small" class="api-json-textarea" />

              <!-- 条件构建器 -->
              <div v-else-if="p.type === 'condition-builder'" class="api-condition-builder">
                <div v-for="(cond, idx) in apiConditions" :key="idx" class="api-cond-row">
                  <el-select v-model="cond.field" placeholder="字段" size="small" style="width: 100px" filterable>
                    <el-option v-for="f in apiModelFields" :key="f.alias" :label="f.name" :value="f.alias" />
                  </el-select>
                  <el-select v-model="cond.op" placeholder="操作符" size="small" style="width: 110px">
                    <el-option v-for="op in apiOperators" :key="op.value" :label="op.label" :value="op.value" />
                  </el-select>
                  <el-input v-model="cond.value" placeholder="值" size="small" style="flex: 1" />
                  <el-button link type="danger" size="small" @click="apiConditions.splice(idx, 1)">
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </div>
                <el-button size="small" @click="apiConditions.push({ field: '', op: 'eq', value: '' })">
                  <el-icon><Plus /></el-icon> 添加条件
                </el-button>
              </div>
            </div>

            <template v-if="isApiReadOnly">
              <el-alert type="warning" :closable="false" show-icon style="margin-top: 12px">
                <template #title>该接口为写入操作，仅支持构造请求参数，不允许直接执行</template>
              </el-alert>
            </template>
            <el-button v-else type="primary" style="width: 100%; margin-top: 12px" @click="sendApiRequest" :loading="apiSending">
              发送请求
            </el-button>
          </div>
        </div>
      </div>

      <!-- 右侧：请求预览 + 响应 -->
      <div class="api-right-panel">
        <!-- 请求预览 -->
        <div class="api-request-preview">
          <div class="api-preview-header">
            <span class="api-panel-title">请求</span>
            <el-button size="small" link @click="copyApiRequest">
              <el-icon><CopyDocument /></el-icon> 复制
            </el-button>
          </div>
          <div class="api-preview-url">{{ apiRequestUrl }}</div>
          <pre v-if="apiRequestBody" class="api-preview-body">{{ apiRequestBody }}</pre>
        </div>

        <!-- 响应结果 -->
        <div class="api-response-section">
          <div class="api-response-header">
            <span class="api-panel-title">响应</span>
            <div class="api-response-actions">
              <span v-if="apiResponseStatus !== null" class="api-response-meta">
                <el-tag :type="apiResponseStatus === 200 ? 'success' : 'danger'" size="small">{{ apiResponseStatus }}</el-tag>
                <span class="api-response-time">{{ apiResponseTime }}ms</span>
              </span>
              <el-button v-if="apiResponseData !== null" size="small" link @click="copyApiResponse">
                <el-icon><CopyDocument /></el-icon> 复制
              </el-button>
            </div>
          </div>
          <div v-if="apiResponseData === null" class="api-response-empty">
            点击「发送请求」查看响应结果
          </div>
          <pre v-else class="api-response-body"><code>{{ formatJson(apiResponseData) }}</code></pre>
        </div>
      </div>
    </div>
  </div>

  <!-- 拓扑图 Drawer -->
  <el-drawer v-model="topoDrawerVisible" :title="topoDrawerTitle" size="85%" direction="rtl" @close="onTopoDrawerClose">
    <div class="topo-drawer-content">
      <div class="topo-drawer-toolbar">
        <el-input v-model="topoGraphSearchText" placeholder="搜索节点..." clearable size="small" style="width: 240px" @input="onTopoGraphSearch">
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <span class="topo-node-count">已加载 {{ topoNodeCount }} 个节点</span>
      </div>
      <div class="topo-drawer-body">
        <div ref="topoGraphContainer" class="topo-graph" />
        <transition name="topo-slide">
          <div v-if="topoSelectedInstance" class="topo-detail-panel">
            <div class="panel-header">
              <div class="panel-title-row">
                <span class="panel-model-name">{{ topoSelectedInstance.model_name }}</span>
                <el-tag size="small" type="info">ID: {{ topoSelectedInstance.id }}</el-tag>
              </div>
              <el-icon class="panel-close" @click="topoSelectedInstance = null"><Close /></el-icon>
            </div>
            <div class="panel-section">
              <div class="panel-section-title">字段值</div>
              <el-table :data="topoFieldRows" size="small" stripe highlight-current-row max-height="360">
                <el-table-column prop="name" label="字段" width="120" />
                <el-table-column prop="alias" label="别名" width="100" />
                <el-table-column prop="value" label="值" show-overflow-tooltip />
              </el-table>
            </div>
            <div class="panel-section">
              <div class="panel-section-title">关系</div>
              <template v-if="topoRelations.length">
                <div v-for="rel in topoRelations" :key="rel.key" class="topo-relation-row" @click="focusTopoNode(rel.instanceId)">
                  <span class="topo-rel-direction" :class="rel.direction">{{ rel.direction === 'upstream' ? '上游' : '下游' }}</span>
                  <span class="topo-rel-type">{{ rel.typeLabel }}</span>
                  <span class="topo-rel-target">{{ rel.modelName }} #{{ rel.instanceId }}</span>
                </div>
              </template>
              <template v-else>
                <div class="panel-loading">暂无关系</div>
              </template>
            </div>
          </div>
        </transition>
      </div>
    </div>
  </el-drawer>
</template>

<script setup>
import { ref, computed, reactive, onMounted, nextTick, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Graph } from '@antv/g6'
import { Search, Delete, CircleCheckFilled, ArrowDown, ArrowRight, Close, CopyDocument } from '@element-plus/icons-vue'
import { getAllModels, getModelGroups, getModelRelationTypes, getModelDetail, fulltextSearch, searchInstance, getInstanceTopology } from '../../api/search'
import { listSearchDirectSql } from '../../api/searchDirectSql'
import openapiRequest from '../../api/openapiRequest'

const typeTag = { string: '', number: 'success', bool: 'warning', date: 'info', datetime: 'info', json: 'danger' }

const tabs = [
  { name: 'fulltext', label: '全文检索' },
  { name: 'instance', label: '实例搜索' },
  { name: 'model', label: '模型浏览' },
  { name: 'api', label: '接口调试' },
]

const activeTab = ref('fulltext')
const allModels = ref([])

function switchTab(name) {
  activeTab.value = name
}

// ===== 全文检索 =====
const ft = ref({
  keyword: '',
  modelAlias: '',
  results: [],
  count: 0,
  page: 1,
  pageSize: 10,
  loading: false,
})

async function doFulltextSearch() {
  const hasKeyword = ft.value.keyword.trim().length > 0
  const hasModel = !!ft.value.modelAlias
  // 无关键词且无模型选择时不搜索
  if (!hasKeyword && !hasModel) return

  ft.value.loading = true
  try {
    // 有关键词 → 全文检索
    if (hasKeyword) {
      const params = {
        search: ft.value.keyword.trim(),
        limit: ft.value.pageSize,
        offset: (ft.value.page - 1) * ft.value.pageSize,
      }
      if (hasModel) params.model_alias = ft.value.modelAlias
      const res = await fulltextSearch(params)
      ft.value.results = res.data?.result || []
      ft.value.count = res.data?.count || 0
    } else {
      // 无关键词但选了模型 → 用实例搜索列出全部
      const body = {
        model: ft.value.modelAlias,
        __condition: {
          limit: ft.value.pageSize,
          offset: (ft.value.page - 1) * ft.value.pageSize,
          order: ['-id'],
        },
      }
      const res = await searchInstance(body)
      const results = res.data?.results || []
      ft.value.results = results.map(r => ({
        id: r.id,
        model_name: ft.value.modelAlias,
        model_alias: ft.value.modelAlias,
        data: Object.fromEntries(Object.entries(r).filter(([k]) => !['id', 'model_id', 'model_alias'].includes(k))),
      }))
      ft.value.count = res.data?.count || 0
    }
  } catch {
    ft.value.results = []
    ft.value.count = 0
  } finally {
    ft.value.loading = false
  }
}

// ===== 实例搜索 =====
const inst = ref({
  modelId: null,
  fields: [],
  conditions: [],
  results: [],
  columns: [],
  columnLabelMap: {},
  count: 0,
  page: 1,
  pageSize: 10,
  loading: false,
  topoLoading: false,
})

async function onModelChange(id) {
  inst.value.fields = []
  inst.value.conditions = []
  inst.value.results = []
  inst.value.columns = []
  inst.value.count = 0
  inst.value.page = 1
  instRelationCache.value = {}
  if (!id) return
  try {
    const m = allModels.value.find(x => x.id === id)
    const res = await getModelDetail(m?.alias || id)
    inst.value.fields = res.data?.model_fields || []
  } catch {
    inst.value.fields = []
  }
}

async function doInstanceSearch() {
  const m = allModels.value.find(x => x.id === inst.value.modelId)
  if (!m) return
  inst.value.loading = true
  try {
    const where = []
    for (const c of inst.value.conditions) {
      if (!c.field || !c.op) continue
      if (c.op === 'in') {
        const vals = c.val.split(',').map(v => v.trim()).filter(Boolean)
        if (vals.length) where.push({ in: { [c.field]: vals } })
      } else if (c.val !== '') {
        where.push({ [c.op]: { [c.field]: c.val } })
      }
    }
    const body = {
      model: m.alias,
      fields: inst.value.fields.map(f => f.alias),
      __condition: {
        limit: inst.value.pageSize,
        offset: (inst.value.page - 1) * inst.value.pageSize,
        order: ['-id'],
      },
    }
    if (where.length) body.__condition.where = where
    const res = await searchInstance(body)
    inst.value.results = res.data?.results || []
    inst.value.count = res.data?.count || 0
    if (inst.value.results.length) {
      inst.value.columns = Object.keys(inst.value.results[0]).filter(k => !['model_id', 'model_alias'].includes(k))
      // 构建列标签映射: alias → "name(alias)"
      const labelMap = { id: 'ID' }
      for (const f of inst.value.fields) {
        labelMap[f.alias] = f.name ? `${f.name}(${f.alias})` : f.alias
      }
      inst.value.columnLabelMap = labelMap
    }
  } catch {
    inst.value.results = []
    inst.value.count = 0
  } finally {
    inst.value.loading = false
  }
  // 加载实例关系
  loadInstanceRelations()
}

function calcColWidth(label) {
  // 中文字符按 14px，英文/数字按 8px，加上 padding
  let w = 0
  for (const ch of label) {
    w += /[一-鿿]/.test(ch) ? 14 : 8
  }
  return Math.max(w + 24, 80)
}

function resetInstance() {
  inst.value.conditions = []
  inst.value.results = []
  inst.value.columns = []
  inst.value.columnLabelMap = {}
  inst.value.count = 0
  inst.value.page = 1
  instRelationCache.value = {}
}

// ===== 模型浏览 =====
const modelSearch = ref('')
const expandedModelId = ref(null)
const allModelGroups = ref([])
const allRelationTypes = ref([])
// 缓存: { [modelId]: { fieldGroups, uniqueAliases, relations, loading } }
const modelFieldsCache = ref({})

const groupedModels = computed(() => {
  const kw = modelSearch.value.toLowerCase()
  const filtered = kw
    ? allModels.value.filter(m =>
        m.name.toLowerCase().includes(kw) ||
        m.alias.toLowerCase().includes(kw) ||
        (m.description || '').toLowerCase().includes(kw))
    : allModels.value

  const groupMap = new Map()
  for (const g of allModelGroups.value) {
    groupMap.set(g.id, { ...g, models: [] })
  }
  const ungrouped = { id: 0, name: '未分组', alias: '', models: [] }

  for (const m of filtered) {
    if (m.group_id && groupMap.has(m.group_id)) {
      groupMap.get(m.group_id).models.push(m)
    } else {
      ungrouped.models.push(m)
    }
  }

  const result = []
  for (const g of groupMap.values()) {
    if (g.models.length) result.push(g)
  }
  if (ungrouped.models.length) result.push(ungrouped)
  return result
})

async function toggleFields(m) {
  if (expandedModelId.value === m.id) {
    expandedModelId.value = null
    return
  }
  expandedModelId.value = m.id
  if (modelFieldsCache.value[m.id]) return
  modelFieldsCache.value[m.id] = { fieldGroups: [], uniqueAliases: [], loading: true }
  try {
    const res = await getModelDetail(m.alias)
    const data = res.data || {}
    const fields = data.model_fields || []
    const groups = data.model_field_groups || []
    const uniques = data.model_field_uniques || []

    // 唯一字段别名
    const uniqueSet = new Set()
    for (const u of uniques) {
      if (u.fields) {
        u.fields.split(',').map(f => f.trim()).filter(Boolean).forEach(f => uniqueSet.add(f))
      }
    }

    // 按分组归类字段
    const groupMap = new Map()
    for (const g of groups) {
      groupMap.set(g.id, { name: g.name, fields: [], _showHeader: true })
    }
    const ungrouped = { name: '', fields: [], _showHeader: true }

    for (const f of fields) {
      if (f.field_group_id && groupMap.has(f.field_group_id)) {
        groupMap.get(f.field_group_id).fields.push(f)
      } else {
        ungrouped.fields.push(f)
      }
    }

    const fieldGroups = []
    for (const g of groupMap.values()) {
      if (g.fields.length) fieldGroups.push(g)
    }
    if (ungrouped.fields.length) fieldGroups.push(ungrouped)

    // 只有第一组显示表头，其余组隐藏避免重复
    fieldGroups.forEach((g, i) => { g._showHeader = i === 0 })

    modelFieldsCache.value[m.id] = {
      fieldGroups,
      uniqueAliases: [...uniqueSet],
      relations: data.model_relations || [],
      loading: false,
    }
  } catch {
    modelFieldsCache.value[m.id] = { fieldGroups: [], uniqueAliases: [], relations: [], loading: false }
  }
}

function resolveRelation(rel) {
  const isSource = rel.source_id === expandedModelId.value
  const targetId = isSource ? rel.target_id : rel.source_id
  const targetModel = allModels.value.find(m => m.id === targetId)
  const relType = allRelationTypes.value.find(t => t.id === rel.type_id)
  const direction = isSource ? 's2t' : 't2s'
  return {
    targetName: targetModel ? `${targetModel.name}(${targetModel.alias})` : `#${targetId}`,
    typeLabel: relType ? relType[direction] : '',
    description: rel.description || '',
  }
}

// ===== 拓扑（Drawer 模式） =====
const GROUP_COLORS = [
  '#4F8EF7', '#F59E0B', '#10B981', '#EF4444', '#8B5CF6',
  '#EC4899', '#06B6D4', '#F97316', '#6366F1', '#14B8A6',
  '#E11D48', '#A855F7', '#0EA5E9', '#D946EF', '#22C55E',
]

const topoDrawerVisible = ref(false)
const topoDrawerTitle = ref('实例拓扑')
const topoGraphSearchText = ref('')
const topoSelectedInstance = ref(null)
const topoGraphContainer = ref(null)
let topoGraph = null
const topoNodeMap = new Map()
const topoEdgeSet = new Set()
let topoModelColorMap = {}
const topoFieldMaps = {}
const topoNodeCount = computed(() => topoNodeMap.size)

// 实例关系表格数据
const instRelations = ref([])

// ===== 实例展开（关联实例） =====
// instRelationCache[instanceId] = { groups, loading, activeTab, tabData: { [modelId]: { instances, fields, loading, page, pageSize } } }
const instRelationCache = ref({})

const onInstExpandChange = (row, expanded) => {
  if (expanded.length && !instRelationCache.value[row.id]) {
    loadInstRelations(row)
  }
}

const loadInstRelations = async (row) => {
  instRelationCache.value[row.id] = { groups: [], loading: true, activeTab: '', tabData: {} }
  try {
    const res = await getInstanceTopology(row.id)
    const topoData = res.data || {}
    const upstream = (topoData.upstream || []).map(r => ({ ...r, direction: 'upstream' }))
    const downstream = (topoData.downstream || []).map(r => ({ ...r, direction: 'downstream' }))
    const allRels = [...upstream, ...downstream]

    // Group by model
    const groupMap = new Map()
    for (const rel of allRels) {
      const key = rel.model_id
      if (!groupMap.has(key)) {
        groupMap.set(key, { model_id: rel.model_id, model_name: rel.model_name, model_alias: rel.model_alias, direction: rel.direction, instances: [] })
      }
      groupMap.get(key).instances.push(rel)
    }

    const groups = [...groupMap.values()]
    instRelationCache.value[row.id].groups = groups
    instRelationCache.value[row.id].loading = false

    if (groups.length) {
      const tabToLoad = String(groups[0].model_id)
      instRelationCache.value[row.id].activeTab = tabToLoad
      onInstTabChange(row, Number(tabToLoad))
    }
  } catch {
    instRelationCache.value[row.id] = { groups: [], loading: false, activeTab: '', tabData: {} }
  }
}

const onInstTabChange = async (row, modelId) => {
  const cache = instRelationCache.value[row.id]
  if (!cache || cache.tabData[modelId]) return

  const group = cache.groups.find(g => g.model_id === modelId)
  if (!group) return

  cache.tabData[modelId] = { instances: group.instances, fields: [], loading: true, page: 1, pageSize: 10 }
  instRelationCache.value = { ...instRelationCache.value }

  try {
    const detailRes = await getModelDetail(group.model_alias)
    const fields = (detailRes.data?.model_fields || []).slice(0, 8)
    cache.tabData[modelId].fields = fields
  } catch {
    cache.tabData[modelId].fields = []
  } finally {
    cache.tabData[modelId].loading = false
    instRelationCache.value = { ...instRelationCache.value }
  }
}

const getInstPaginatedInstances = (instanceId, modelId) => {
  const tabData = instRelationCache.value[instanceId]?.tabData?.[modelId]
  if (!tabData) return []
  const size = tabData.pageSize || 10
  const start = (tabData.page - 1) * size
  return tabData.instances.slice(start, start + size)
}

const topoFieldRows = computed(() => {
  if (!topoSelectedInstance.value) return []
  const data = topoSelectedInstance.value.data || {}
  const modelAlias = topoSelectedInstance.value.model_alias || ''
  const fieldMap = topoFieldMaps[modelAlias] || {}
  return Object.entries(data).map(([alias, value]) => {
    const field = fieldMap[alias]
    return {
      alias,
      name: field?.name || alias,
      value: typeof value === 'object' ? JSON.stringify(value) : String(value ?? ''),
    }
  })
})

const topoRelations = computed(() => {
  if (!topoSelectedInstance.value) return []
  return (topoSelectedInstance.value._relations || []).map(r => ({
    key: `${r.direction}-${r.instance_id}`,
    instanceId: r.instance_id,
    direction: r.direction,
    typeLabel: r.relation_type,
    modelName: r.model_name,
  }))
})

function destroyTopoGraph() {
  if (topoGraph) {
    topoGraph.destroy()
    topoGraph = null
  }
  topoNodeMap.clear()
  topoEdgeSet.clear()
  topoModelColorMap = {}
}

function buildNodeLabel(inst, fields) {
  const lines = [`#${inst.id}`]
  const displayFields = fields.filter(f => f.alias !== 'id').slice(0, 4)
  for (const f of displayFields) {
    const val = inst[f.alias]
    if (val !== undefined && val !== null && val !== '') {
      const short = String(val).length > 16 ? String(val).slice(0, 16) + '...' : String(val)
      lines.push(`${f.name || f.alias}: ${short}`)
    }
  }
  return lines.join('\n')
}

// 批量加载拓扑图（当前页实例）
async function openBatchTopo() {
  const model = allModels.value.find(m => m.id === inst.value.modelId)
  if (!model || !inst.value.results.length) return
  destroyTopoGraph()
  topoSelectedInstance.value = null

  const fields = inst.value.fields
  const groupIdx = allModelGroups.value.findIndex(g => g.id === model.group_id)
  topoModelColorMap[model.id] = GROUP_COLORS[(groupIdx >= 0 ? groupIdx : 0) % GROUP_COLORS.length]

  const instances = inst.value.results
  const instanceIdSet = new Set()
  const nodes = instances.map(inst => {
    const nodeId = `inst-${inst.id}`
    instanceIdSet.add(inst.id)
    topoNodeMap.set(nodeId, { ...inst, model_id: model.id, model_name: model.name, model_alias: model.alias })
    return {
      id: nodeId,
      data: { label: buildNodeLabel(inst, fields), instanceId: inst.id, modelId: model.id },
      style: { fill: topoModelColorMap[model.id], stroke: '#fff', lineWidth: 2, radius: 12 },
    }
  })

  // Batch-fetch topology to build edges
  const topoPromises = instances.map(i => getInstanceTopology(i.id).catch(() => null))
  const topoResults = await Promise.all(topoPromises)
  const edges = []
  for (const res of topoResults) {
    if (!res?.data) continue
    const { instance: srcInst, upstream, downstream } = res.data
    if (!srcInst) continue
    const srcId = srcInst.id
    for (const rel of (downstream || [])) {
      const edgeId = `edge-inst-${srcId}-inst-${rel.instance_id}`
      if (topoEdgeSet.has(edgeId)) continue
      const relNodeId = `inst-${rel.instance_id}`
      if (!topoNodeMap.has(relNodeId)) {
        topoModelColorMap[rel.model_id] = topoModelColorMap[rel.model_id] || GROUP_COLORS[Object.keys(topoModelColorMap).length % GROUP_COLORS.length]
        topoNodeMap.set(relNodeId, { id: rel.instance_id, model_id: rel.model_id, model_name: rel.model_name, model_alias: rel.model_alias, ...(rel.data || {}) })
        nodes.push({ id: relNodeId, data: { label: `#${rel.instance_id}\n${rel.model_name}`, instanceId: rel.instance_id, modelId: rel.model_id }, style: { fill: topoModelColorMap[rel.model_id], stroke: '#fff', lineWidth: 1.5, radius: 10 } })
      }
      topoEdgeSet.add(edgeId)
      edges.push({ id: edgeId, source: `inst-${srcId}`, target: relNodeId, data: { relType: rel.relation_type || '' } })
    }
    for (const rel of (upstream || [])) {
      const edgeId = `edge-inst-${rel.instance_id}-inst-${srcId}`
      if (topoEdgeSet.has(edgeId)) continue
      const relNodeId = `inst-${rel.instance_id}`
      if (!topoNodeMap.has(relNodeId)) {
        topoModelColorMap[rel.model_id] = topoModelColorMap[rel.model_id] || GROUP_COLORS[Object.keys(topoModelColorMap).length % GROUP_COLORS.length]
        topoNodeMap.set(relNodeId, { id: rel.instance_id, model_id: rel.model_id, model_name: rel.model_name, model_alias: rel.model_alias, ...(rel.data || {}) })
        nodes.push({ id: relNodeId, data: { label: `#${rel.instance_id}\n${rel.model_name}`, instanceId: rel.instance_id, modelId: rel.model_id }, style: { fill: topoModelColorMap[rel.model_id], stroke: '#fff', lineWidth: 1.5, radius: 10 } })
      }
      topoEdgeSet.add(edgeId)
      edges.push({ id: edgeId, source: relNodeId, target: `inst-${srcId}`, data: { relType: rel.relation_type || '' } })
    }
  }

  topoDrawerTitle.value = `实例拓扑 — ${model.name}（${instances.length} 条）`
  topoDrawerVisible.value = true
  await nextTick()
  initTopoGraph(nodes, edges)
}

// 单实例拓扑
async function openSingleTopo(row) {
  const model = allModels.value.find(m => m.id === inst.value.modelId)
  if (!model) return
  destroyTopoGraph()
  topoSelectedInstance.value = null

  const groupIdx = allModelGroups.value.findIndex(g => g.id === model.group_id)
  topoModelColorMap[model.id] = GROUP_COLORS[(groupIdx >= 0 ? groupIdx : 0) % GROUP_COLORS.length]

  const nodeId = `inst-${row.id}`
  topoNodeMap.set(nodeId, { ...row, model_id: model.id, model_name: model.name, model_alias: model.alias })
  const nodes = [{
    id: nodeId,
    data: { label: buildNodeLabel(row, inst.value.fields), instanceId: row.id, modelId: model.id },
    style: { fill: topoModelColorMap[model.id], stroke: '#fff', lineWidth: 2, radius: 12 },
  }]

  // Fetch topology for this instance
  try {
    const res = await getInstanceTopology(row.id)
    const topoData = res.data || {}
    const edges = []
    for (const rel of (topoData.downstream || [])) {
      const relNodeId = `inst-${rel.instance_id}`
      if (!topoNodeMap.has(relNodeId)) {
        topoModelColorMap[rel.model_id] = topoModelColorMap[rel.model_id] || GROUP_COLORS[Object.keys(topoModelColorMap).length % GROUP_COLORS.length]
        topoNodeMap.set(relNodeId, { id: rel.instance_id, model_id: rel.model_id, model_name: rel.model_name, model_alias: rel.model_alias, ...(rel.data || {}) })
        nodes.push({ id: relNodeId, data: { label: `#${rel.instance_id}\n${rel.model_name}`, instanceId: rel.instance_id, modelId: rel.model_id }, style: { fill: topoModelColorMap[rel.model_id], stroke: '#fff', lineWidth: 1.5, radius: 10 } })
      }
      const edgeId = `edge-${nodeId}-${relNodeId}`
      if (!topoEdgeSet.has(edgeId)) {
        topoEdgeSet.add(edgeId)
        edges.push({ id: edgeId, source: nodeId, target: relNodeId, data: { relType: rel.relation_type || '' } })
      }
    }
    for (const rel of (topoData.upstream || [])) {
      const relNodeId = `inst-${rel.instance_id}`
      if (!topoNodeMap.has(relNodeId)) {
        topoModelColorMap[rel.model_id] = topoModelColorMap[rel.model_id] || GROUP_COLORS[Object.keys(topoModelColorMap).length % GROUP_COLORS.length]
        topoNodeMap.set(relNodeId, { id: rel.instance_id, model_id: rel.model_id, model_name: rel.model_name, model_alias: rel.model_alias, ...(rel.data || {}) })
        nodes.push({ id: relNodeId, data: { label: `#${rel.instance_id}\n${rel.model_name}`, instanceId: rel.instance_id, modelId: rel.model_id }, style: { fill: topoModelColorMap[rel.model_id], stroke: '#fff', lineWidth: 1.5, radius: 10 } })
      }
      const edgeId = `edge-${relNodeId}-${nodeId}`
      if (!topoEdgeSet.has(edgeId)) {
        topoEdgeSet.add(edgeId)
        edges.push({ id: edgeId, source: relNodeId, target: nodeId, data: { relType: rel.relation_type || '' } })
      }
    }

    topoDrawerTitle.value = `实例拓扑 — ${model.name} #${row.id}`
    topoDrawerVisible.value = true
    await nextTick()
    initTopoGraph(nodes, edges)
  } catch { /* ignore */ }
}

function initTopoGraph(nodes, edges) {
  if (!topoGraphContainer.value) return
  const rect = topoGraphContainer.value.getBoundingClientRect()

  topoGraph = new Graph({
    container: topoGraphContainer.value,
    width: rect.width,
    height: rect.height,
    data: { nodes, edges },
    node: {
      style: {
        labelText: (d) => d.data?.label || '',
        labelFontSize: 10,
        labelFontWeight: 400,
        labelFill: '#1E293B',
        labelPlacement: 'center',
        labelMaxWidth: 160,
        wordWrap: true,
        cursor: 'pointer',
      },
    },
    edge: {
      style: {
        labelText: (d) => d.data?.relType || '',
        labelFontSize: 9,
        labelFill: '#64748B',
        endArrow: true,
        endArrowSize: 5,
        stroke: '#C8D0DE',
        lineWidth: 1.5,
      },
    },
    layout: {
      type: 'force',
      preventOverlap: true,
      nodeSize: 80,
      linkDistance: 200,
      strength: { charge: -500 },
    },
    behaviors: ['drag-canvas', 'zoom-canvas', 'drag-element'],
    autoFit: 'view',
  })

  topoGraph.render().then(() => {
    topoGraph.on('node:click', async (e) => {
      const nodeId = e.target.id
      const instData = topoNodeMap.get(nodeId)
      if (!instData) return
      await loadTopoDetail(instData.id || instData.instance_id, instData)
    })
    topoGraph.on('canvas:click', () => {
      topoSelectedInstance.value = null
    })
  })
}

async function loadTopoDetail(instId, instData) {
  const data = {}
  for (const [k, v] of Object.entries(instData)) {
    if (!['id', 'model_id', 'model_alias', 'model_name', 'instance_id'].includes(k)) {
      data[k] = v
    }
  }
  const modelAlias = instData.model_alias || ''
  topoSelectedInstance.value = { id: instId, model_name: instData.model_name || '', model_alias: modelAlias, data, _relations: [] }

  if (modelAlias && !topoFieldMaps[modelAlias]) {
    try {
      const detailRes = await getModelDetail(modelAlias)
      const fields = detailRes.data?.model_fields || []
      const fm = {}
      for (const f of fields) fm[f.alias] = f
      topoFieldMaps[modelAlias] = fm
    } catch { /* ignore */ }
  }

  try {
    const res = await getInstanceTopology(instId)
    const topoData = res.data || {}
    const upstream = (topoData.upstream || []).map(r => ({ ...r, direction: 'upstream' }))
    const downstream = (topoData.downstream || []).map(r => ({ ...r, direction: 'downstream' }))
    topoSelectedInstance.value._relations = [...upstream, ...downstream]

    if (topoData.instance?.data) {
      topoSelectedInstance.value.data = topoData.instance.data
    }

    // Expand: add new nodes and edges
    const newNodes = []
    const newEdges = []
    for (const rel of [...upstream, ...downstream]) {
      const relNodeId = `inst-${rel.instance_id}`
      if (!topoNodeMap.has(relNodeId)) {
        topoModelColorMap[rel.model_id] = topoModelColorMap[rel.model_id] || GROUP_COLORS[Object.keys(topoModelColorMap).length % GROUP_COLORS.length]
        topoNodeMap.set(relNodeId, { id: rel.instance_id, model_id: rel.model_id, model_name: rel.model_name, model_alias: rel.model_alias, ...(rel.data || {}) })
        const relFields = topoFieldMaps[rel.model_alias]
        const label = relFields ? buildNodeLabel({ id: rel.instance_id, ...(rel.data || {}) }, Object.values(relFields)) : `#${rel.instance_id}\n${rel.model_name}`
        newNodes.push({ id: relNodeId, data: { label, instanceId: rel.instance_id, modelId: rel.model_id }, style: { fill: topoModelColorMap[rel.model_id], stroke: '#fff', lineWidth: 1.5, radius: 10 } })
      }
      const sourceId = rel.direction === 'downstream' ? `inst-${instId}` : relNodeId
      const targetId = rel.direction === 'downstream' ? relNodeId : `inst-${instId}`
      const edgeId = `edge-${sourceId}-${targetId}`
      if (!topoEdgeSet.has(edgeId)) {
        topoEdgeSet.add(edgeId)
        newEdges.push({ id: edgeId, source: sourceId, target: targetId, data: { relType: rel.relation_type || '' } })
      }
    }
    if (newNodes.length || newEdges.length) {
      for (const n of newNodes) topoGraph.addNodeData(n)
      for (const e of newEdges) topoGraph.addEdgeData(e)
      topoGraph.layout()
    }
  } catch { /* ignore */ }
}

function onTopoGraphSearch(text) {
  if (!topoGraph) return
  const allNodeIds = topoGraph.getNodeData().map(n => n.id)
  if (!text || !text.trim()) {
    topoGraph.updateNodeData(allNodeIds.map(id => ({ id, style: { opacity: 1, lineWidth: 2 } })))
    return
  }
  const kw = text.trim().toLowerCase()
  const matchedIds = []
  const updates = allNodeIds.map(nid => {
    const instData = topoNodeMap.get(nid)
    if (!instData) return { id: nid, style: { opacity: 0.15, lineWidth: 1 } }
    const dataStr = JSON.stringify(instData.data || instData).toLowerCase()
    const idStr = String(instData.id || instData.instance_id)
    const matched = idStr.includes(kw) || dataStr.includes(kw)
    if (matched) matchedIds.push(nid)
    return { id: nid, style: { opacity: matched ? 1 : 0.15, lineWidth: matched ? 3 : 1 } }
  })
  topoGraph.updateNodeData(updates)
  if (matchedIds.length === 1) topoGraph.focusElement(matchedIds[0])
}

function focusTopoNode(instId) {
  const nodeId = `inst-${instId}`
  if (!topoGraph) return
  if (topoGraph.getNodeData().find(n => n.id === nodeId)) {
    topoGraph.focusElement(nodeId)
    const instData = topoNodeMap.get(nodeId)
    if (instData) loadTopoDetail(instData.id || instData.instance_id, instData)
  }
}

function onTopoDrawerClose() {
  destroyTopoGraph()
  topoSelectedInstance.value = null
  topoGraphSearchText.value = ''
}

// 加载实例关系表格
async function loadInstanceRelations() {
  if (!inst.value.results.length) { instRelations.value = []; return }
  const model = allModels.value.find(m => m.id === inst.value.modelId)
  if (!model) return
  const idSet = new Set(inst.value.results.map(r => r.id))
  const topoPromises = inst.value.results.map(i => getInstanceTopology(i.id).catch(() => null))
  const topoResults = await Promise.all(topoPromises)
  const relations = []
  for (const res of topoResults) {
    if (!res?.data) continue
    const { instance: srcInst, upstream, downstream } = res.data
    if (!srcInst) continue
    for (const rel of (downstream || [])) {
      if (idSet.has(rel.instance_id)) {
        relations.push({ sourceId: srcInst.id, sourceLabel: `${model.name} #${srcInst.id}`, relationType: rel.relation_type || '', targetId: rel.instance_id, targetLabel: `${rel.model_name} #${rel.instance_id}` })
      }
    }
  }
  instRelations.value = relations
}

// ===== 接口调试 =====
const apiSavedSqls = ref([])
const apiModelFields = ref([])
const apiEndpointGroups = reactive([
  {
    name: '模型查询',
    collapsed: false,
    endpoints: [
      { key: 'model-all', name: '所有模型', method: 'GET', path: '/model/all', params: [] },
      { key: 'model-group', name: '模型分组', method: 'GET', path: '/model/group', params: [] },
      { key: 'model-single', name: '模型详情', method: 'GET', path: '/model/single', params: [
        { key: 'alias', label: '模型', type: 'select', required: true, source: 'models', placeholder: '选择模型' }
      ]},
      { key: 'model-relation-type', name: '关系类型', method: 'GET', path: '/model/relation-type', params: [] },
      { key: 'model-relation', name: '模型关系', method: 'GET', path: '/model/relation', params: [] },
    ]
  },
  {
    name: '实例查询',
    collapsed: false,
    endpoints: [
      { key: 'instance-fulltext', name: '全文搜索', method: 'GET', path: '/instance/fulltext', params: [
        { key: 'search', label: '搜索词', type: 'input', required: true },
        { key: 'model_alias', label: '模型过滤', type: 'select', source: 'models', placeholder: '可选' },
        { key: 'limit', label: '数量', type: 'number' },
        { key: 'offset', label: '偏移', type: 'number' },
      ]},
      { key: 'instance-search', name: '结构化搜索', method: 'GET', path: '/instance/search', params: [
        { key: 'model', label: '模型', type: 'select', required: true, source: 'models', placeholder: '选择模型' },
        { key: 'fields', label: '返回字段', type: 'multi-select', source: 'fields', placeholder: '不选则返回全部' },
        { key: 'where', label: '查询条件', type: 'condition-builder' },
        { key: 'order', label: '排序', type: 'select', source: 'orderFields', placeholder: '可选' },
        { key: 'limit', label: '数量', type: 'number' },
        { key: 'offset', label: '偏移', type: 'number' },
      ]},
      { key: 'instance-detail', name: '实例详情', method: 'GET', path: '/instance/detail', params: [
        { key: 'id', label: '实例ID', type: 'input', required: true, placeholder: '输入实例ID' },
      ]},
      { key: 'instance-topology', name: '拓扑关系', method: 'GET', path: '/instance/topology', params: [
        { key: 'id', label: '实例ID', type: 'input', required: true, placeholder: '输入实例ID' },
        { key: 'model', label: '模型过滤', type: 'select', source: 'models', placeholder: '可选' },
      ]},
      { key: 'instance-source', name: '查询上游实例', method: 'GET', path: '/instance/source', params: [
        { key: 'id', label: '目标实例ID', type: 'input', required: true },
        { key: 'model', label: '源模型', type: 'select', required: true, source: 'models', placeholder: '选择源模型' },
      ]},
      { key: 'instance-target', name: '查询下游实例', method: 'GET', path: '/instance/target', params: [
        { key: 'id', label: '源实例ID', type: 'input', required: true },
        { key: 'model', label: '目标模型', type: 'select', required: true, source: 'models', placeholder: '选择目标模型' },
      ]},
      { key: 'instance-direct', name: '执行SQL查询', method: 'GET', path: '/instance/direct', params: [
        { key: 'uuid', label: '查询', type: 'select', source: 'sqls', required: true, placeholder: '选择已保存的查询' },
      ]},
    ]
  },
  {
    name: '实例写入',
    collapsed: false,
    endpoints: [
      { key: 'instance-create', name: '创建实例', method: 'POST', path: '/instance/create', params: [
        { key: 'model_alias', label: '模型', type: 'select', required: true, source: 'models', placeholder: '选择模型' },
        { key: 'data', label: '数据 (JSON)', type: 'textarea', required: true, placeholder: '{"field1": "value1", ...}' },
      ]},
      { key: 'instance-update', name: '更新实例', method: 'POST', path: '/instance/update', params: [
        { key: 'id', label: '实例ID', type: 'input', required: true },
        { key: 'data', label: '数据 (JSON)', type: 'textarea', required: true, placeholder: '{"field1": "new_value", ...}' },
      ]},
      { key: 'instance-delete', name: '删除实例', method: 'POST', path: '/instance/delete', params: [
        { key: 'model_alias', label: '模型', type: 'select', required: true, source: 'models', placeholder: '选择模型' },
        { key: 'id', label: '实例ID', type: 'input', required: true },
      ]},
    ]
  },
  {
    name: '实例关系',
    collapsed: false,
    endpoints: [
      { key: 'relation-create', name: '创建关系', method: 'POST', path: '/instance-relation/create', params: [
        { key: 'source_model', label: '源模型', type: 'select', required: true, source: 'models', placeholder: '选择源模型' },
        { key: 'target_model', label: '目标模型', type: 'select', required: true, source: 'models', placeholder: '选择目标模型' },
        { key: 'source_id', label: '源实例ID', type: 'input', required: true },
        { key: 'target_ids', label: '目标实例ID', type: 'input', required: true, placeholder: '多个用逗号分隔' },
      ]},
      { key: 'relation-delete', name: '删除关系', method: 'POST', path: '/instance-relation/delete', params: [
        { key: 'source_model', label: '源模型', type: 'select', required: true, source: 'models', placeholder: '选择源模型' },
        { key: 'target_model', label: '目标模型', type: 'select', required: true, source: 'models', placeholder: '选择目标模型' },
        { key: 'source_id', label: '源实例ID', type: 'input', required: true },
        { key: 'target_id', label: '目标实例ID', type: 'input', required: true },
      ]},
    ]
  },
])

const apiSelectedKey = ref('')
const apiSelected = ref(null)
const apiParamValues = reactive({})
const apiConditions = ref([])
const apiSending = ref(false)
const apiRequestUrl = ref('')
const apiRequestBody = ref('')
const apiResponseData = ref(null)
const apiResponseStatus = ref(null)
const apiResponseTime = ref(0)

const apiOperators = [
  { label: '等于 (=)', value: 'eq' },
  { label: '不等于 (!=)', value: 'ne' },
  { label: '大于 (>)', value: 'gt' },
  { label: '大于等于 (>=)', value: 'ge' },
  { label: '小于 (<)', value: 'lt' },
  { label: '小于等于 (<=)', value: 'le' },
  { label: '包含 (in)', value: 'in' },
  { label: '模糊匹配', value: 'contains' },
  { label: '前缀匹配', value: 'startswith' },
  { label: '后缀匹配', value: 'endswith' },
]

const apiOrderFields = [
  { label: 'ID 升序', value: 'id' },
  { label: 'ID 降序', value: '-id' },
  { label: '创建时间 升序', value: 'created_at' },
  { label: '创建时间 降序', value: '-created_at' },
  { label: '更新时间 升序', value: 'updated_at' },
  { label: '更新时间 降序', value: '-updated_at' },
]

const getApiOptions = (p) => {
  switch (p.source) {
    case 'models':
      return allModels.value.map(m => ({ label: `${m.name} (${m.alias})`, value: m.alias }))
    case 'fields':
      return apiModelFields.value.map(f => ({ label: `${f.name} (${f.alias})`, value: f.alias }))
    case 'sqls':
      return apiSavedSqls.value.map(s => ({ label: `${s.name} (${s.uuid?.slice(0, 8)}...)`, value: s.uuid }))
    case 'orderFields':
      return apiOrderFields
    default:
      return []
  }
}

const selectApiEndpoint = (ep) => {
  apiSelectedKey.value = ep.key
  apiSelected.value = ep
  // 重置参数：先清空已有键，再设置新键（避免 delete 破坏 reactive 追踪）
  for (const key of Object.keys(apiParamValues)) {
    apiParamValues[key] = undefined
  }
  ep.params.forEach(p => {
    if (p.type === 'number') apiParamValues[p.key] = undefined
    else if (p.type === 'multi-select') apiParamValues[p.key] = []
    else if (p.type === 'textarea') apiParamValues[p.key] = ''
    else apiParamValues[p.key] = ''
  })
  apiConditions.value = []
  apiModelFields.value = []
  updateApiPreview()
}

// 监听模型选择变化（加载字段列表）
watch(() => apiParamValues.model, async (alias) => {
  if (!alias) { apiModelFields.value = []; return }
  try {
    const res = await getModelDetail(alias)
    apiModelFields.value = (res.data?.model_fields || [])
  } catch {
    apiModelFields.value = []
  }
})

// 通用 GET 参数构建
const buildGetParams = (params, paramValues, conditions) => {
  const qs = new URLSearchParams()
  for (const p of params) {
    const val = paramValues[p.key]

    if (p.type === 'condition-builder') {
      const where = conditions.filter(c => c.field && c.op).map(c => {
        let v = c.value
        if (c.op === 'in') v = v.split(',').map(s => s.trim())
        return { [c.op]: { [c.field]: v } }
      })
      if (where.length) qs.set('where', JSON.stringify(where))
    } else if (p.type === 'multi-select') {
      if (val?.length) qs.set(p.key, val.join(','))
    } else if (p.type === 'number') {
      if (val !== undefined && val !== null && val !== '') qs.set(p.key, String(val))
    } else {
      if (val) qs.set(p.key, val)
    }
  }
  const str = qs.toString()
  return str ? `?${str}` : ''
}

// 通用 POST body 构建
const buildPostBody = (ep, paramValues, conditions) => {
  const body = {}
  for (const p of ep.params) {
    const val = paramValues[p.key]
    if (p.type === 'textarea') {
      // JSON 字段：尝试解析
      if (val) {
        try { body[p.key] = JSON.parse(val) } catch { body[p.key] = val }
      }
    } else if (p.type === 'number') {
      if (val !== undefined && val !== null && val !== '') body[p.key] = val
    } else if (p.type === 'multi-select') {
      if (val?.length) body[p.key] = val
    } else {
      if (val) body[p.key] = val
    }
  }
  // target_ids 逗号分隔转数组
  if (body.target_ids && typeof body.target_ids === 'string') {
    body.target_ids = body.target_ids.split(',').map(s => parseInt(s.trim(), 10)).filter(n => !isNaN(n))
  }
  // source_id / target_id / id 转数字
  if (body.source_id) body.source_id = parseInt(body.source_id, 10) || body.source_id
  if (body.target_id) body.target_id = parseInt(body.target_id, 10) || body.target_id
  if (body.id && ep.key !== 'instance-update') body.id = parseInt(body.id, 10) || body.id
  return body
}

const buildApiRequest = () => {
  if (!apiSelected.value) return { url: '', body: null, method: '' }

  const ep = apiSelected.value
  const method = ep.method
  let url = ep.path
  let body = null

  if (method === 'GET') {
    url = ep.path + buildGetParams(ep.params, apiParamValues, apiConditions.value)
  } else {
    body = buildPostBody(ep, apiParamValues, apiConditions.value)
  }

  return { url, body, method }
}

const updateApiPreview = () => {
  const { url, body, method } = buildApiRequest()
  apiRequestUrl.value = method ? `${method} ${url}` : ''
  apiRequestBody.value = body ? JSON.stringify(body, null, 2) : ''
}

watch(apiParamValues, updateApiPreview, { deep: true })
watch(apiConditions, updateApiPreview, { deep: true })

// 检查当前选中接口是否为只读（危险操作）
const isApiReadOnly = computed(() => {
  const key = apiSelectedKey.value
  return key === 'instance-create' || key === 'instance-update' || key === 'instance-delete' ||
         key === 'relation-create' || key === 'relation-delete'
})

const sendApiRequest = async () => {
  if (!apiSelected.value || apiSending.value || isApiReadOnly.value) return

  const req = buildApiRequest()
  if (!req.url) return
  const { url, body, method } = req

  apiSending.value = true
  apiResponseData.value = null
  apiResponseStatus.value = null
  apiResponseTime.value = 0

  const startTime = Date.now()
  try {
    let res
    if (method === 'GET') {
      res = await openapiRequest.get(url)
    } else {
      res = await openapiRequest.post(url, body)
    }
    apiResponseTime.value = Date.now() - startTime
    apiResponseStatus.value = 200
    apiResponseData.value = res
  } catch (err) {
    apiResponseTime.value = Date.now() - startTime
    // openapiRequest interceptor may reject with Error (no .response) or axios error
    apiResponseStatus.value = err.response?.status || err.status || 0
    apiResponseData.value = err.response?.data || { error: err.message || '请求失败' }
  } finally {
    apiSending.value = false
  }
}

const copyApiRequest = () => {
  const { url, body, method } = buildApiRequest()
  let text = `${method} ${url}`
  if (body) text += `\n\n${JSON.stringify(body, null, 2)}`
  navigator.clipboard.writeText(text)
  ElMessage.success('已复制到剪贴板')
}

const formatJson = (data) => {
  try {
    return JSON.stringify(data, null, 2)
  } catch {
    return String(data)
  }
}

const copyApiResponse = () => {
  if (apiResponseData.value === null) return
  navigator.clipboard.writeText(formatJson(apiResponseData.value))
  ElMessage.success('已复制到剪贴板')
}

// ===== 初始化 =====
onMounted(async () => {
  try {
    const [modelsRes, groupsRes, typesRes, sqlRes] = await Promise.all([
      getAllModels(),
      getModelGroups(),
      getModelRelationTypes(),
      listSearchDirectSql({ limit: 1000 }).catch(() => null),
    ])
    allModels.value = (modelsRes.data || []).filter(m => m.is_usable !== false)
    allModelGroups.value = groupsRes.data || []
    allRelationTypes.value = typesRes.data || []
    apiSavedSqls.value = sqlRes?.data?.results || []
  } catch {}
})
</script>

<style scoped>
.search-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-height: var(--page-height);
}

/* ===== Tab 栏 ===== */
.tab-bar {
  display: flex;
  gap: 32px;
  margin-top: 40px;
  margin-bottom: 32px;
}

.tab-item {
  font-size: 16px;
  color: var(--color-text-secondary);
  cursor: pointer;
  padding-bottom: 8px;
  border-bottom: 3px solid transparent;
  transition: all 0.2s ease;
  user-select: none;
}

.tab-item:hover {
  color: var(--color-text-primary);
}

.tab-item.active {
  color: var(--color-primary);
  border-bottom-color: var(--color-primary);
  font-weight: 600;
}

/* ===== 搜索区 ===== */
.search-hero {
  width: 100%;
  display: flex;
  justify-content: center;
  margin-bottom: 40px;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  max-width: 640px;
}

.search-input {
  flex: 1;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  padding: 4px 20px;
  height: 40px;
  transition: box-shadow 0.2s ease;
}

.search-input :deep(.el-input__wrapper:focus-within) {
  box-shadow: 0 2px 16px rgba(37, 99, 235, 0.15);
}

.search-model-select {
  width: 160px;
}

.search-model-select :deep(.el-input__wrapper) {
  border-radius: 24px;
  height: 40px;
}

.search-btn {
  border-radius: 24px;
  padding: 0 28px;
  height: 40px;
  font-weight: 600;
  flex-shrink: 0;
}

/* 实例搜索：复用全文检索的 hero 居中布局 */
.instance-search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  max-width: 640px;
}

.instance-search-box .instance-model-select {
  flex: 1;
}

.instance-search-box .instance-model-select :deep(.el-input__wrapper) {
  border-radius: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  padding: 4px 20px;
  height: 40px;
  transition: box-shadow 0.2s ease;
}

.instance-search-box .instance-model-select :deep(.el-input__wrapper:focus-within) {
  box-shadow: 0 2px 16px rgba(37, 99, 235, 0.15);
}

.instance-search-box .search-btn {
  border-radius: 24px;
  padding: 0 28px;
  height: 40px;
  font-weight: 600;
  flex-shrink: 0;
}

/* ===== 结果区 ===== */
.results-area {
  width: 100%;
  max-width: 960px;
}

.result-count {
  font-size: 13px;
  color: var(--color-text-muted);
  margin-bottom: 12px;
}

.result-count strong {
  color: var(--color-primary);
}

/* 数据标签 */
.data-tag {
  display: inline-block;
  font-size: 12px;
  padding: 1px 6px;
  margin: 2px 4px 2px 0;
  background: var(--color-muted);
  border-radius: 4px;
  color: var(--color-text-secondary);
}

.data-tag strong {
  color: var(--color-text-primary);
}

/* 分页 */
.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

/* 条件构建器外层间距 */

/* 条件构建器 */
.condition-builder {
  margin-bottom: 16px;
  padding: 16px;
  background: var(--color-muted);
  border-radius: var(--radius-md);
}

.condition-row {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
}

.cond-field {
  width: 180px;
}

.cond-op {
  width: 120px;
}

.cond-val {
  flex: 1;
}

.condition-actions {
  display: flex;
  gap: 8px;
  margin-top: 4px;
}

/* 模型分组 */
.model-group-section {
  margin-bottom: 24px;
}

.model-group-header {
  font-size: 16px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin-bottom: 10px;
  padding-bottom: 8px;
  border-bottom: 2px solid var(--color-primary);
}

/* 模型列表 */
.model-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.model-card {
  padding: 12px 16px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.15s ease;
}

.model-card:hover {
  border-color: var(--color-primary);
  background: var(--color-muted);
}

.model-card.expanded {
  border-color: var(--color-primary);
  background: rgba(37, 99, 235, 0.03);
}

.model-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.model-card-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.model-card-name {
  font-weight: 600;
  color: var(--color-text-primary);
}

.model-card-desc {
  margin-top: 4px;
  font-size: 12px;
  color: var(--color-text-muted);
}

.expand-icon {
  transition: transform 0.2s ease;
  color: var(--color-text-muted);
}

.expand-icon.rotated {
  transform: rotate(180deg);
}

/* 字段内联展开 */
.fields-inline {
  margin: -2px 0 6px 0;
  padding: 12px 16px;
  border: 1px solid var(--color-primary);
  border-top: none;
  border-radius: 0 0 var(--radius-md) var(--radius-md);
  background: rgba(37, 99, 235, 0.02);
}

.fields-loading {
  text-align: center;
  padding: 16px;
  font-size: 13px;
  color: var(--color-text-muted);
}

/* 字段分组 */
.field-group {
  margin-bottom: 12px;
}

.field-group:last-child {
  margin-bottom: 0;
}

.field-group-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-bottom: 6px;
  padding-left: 2px;
}

.unique-tag {
  margin-left: 6px;
  vertical-align: middle;
}

/* 模型关系 */
.model-relations {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px dashed var(--color-border);
}

.relations-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-bottom: 8px;
}

.relation-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 0;
  font-size: 13px;
}

.relation-type {
  color: var(--color-primary);
  font-weight: 500;
}

.relation-target {
  font-weight: 600;
  color: var(--color-text-primary);
}

.relation-desc {
  color: var(--color-text-muted);
  font-size: 12px;
}

/* ===== 实例展开（关联实例） ===== */
.inst-relation-panel {
  padding: 12px 20px;
  background: var(--color-muted);
  border-radius: var(--radius-md);
  margin: 4px 0;
}

.inst-relation-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.inst-relation-header :deep(.el-tabs__header) {
  margin: 0;
}

.inst-relation-tabs {
  flex: 1;
  min-width: 0;
}

.tab-label {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.tab-direction {
  font-weight: 700;
  font-size: 13px;
}

.tab-direction.upstream {
  color: var(--color-primary);
}

.tab-direction.downstream {
  color: var(--color-warning);
}

.tab-badge {
  margin-left: 4px;
}

.tab-badge :deep(.el-badge__content) {
  font-size: 10px;
}

.inst-relation-pagination {
  display: flex;
  justify-content: center;
  margin-top: 8px;
}

.inst-relation-empty {
  text-align: center;
  padding: 24px;
  font-size: 13px;
  color: var(--color-text-muted);
}

/* ===== 实例关系 ===== */
.inst-relations-section {
  margin-top: 16px;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 10px;
  padding-bottom: 6px;
  border-bottom: 1px solid var(--color-border);
}

/* ===== 拓扑 Drawer ===== */
.topo-drawer-content {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.topo-drawer-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid var(--color-border);
  flex-shrink: 0;
}

.topo-node-count {
  font-size: 13px;
  color: var(--color-text-muted);
  margin-left: auto;
}

.topo-drawer-body {
  display: flex;
  flex: 1;
  min-height: 0;
  position: relative;
}

.topo-graph {
  flex: 1;
  background: var(--color-background);
  border-radius: var(--radius-md);
  min-width: 0;
}

.topo-detail-panel {
  width: 340px;
  flex-shrink: 0;
  background: var(--color-surface);
  border-left: 1px solid var(--color-border);
  padding: 16px;
  overflow-y: auto;
  box-shadow: -4px 0 16px rgba(0, 0, 0, 0.06);
}

.panel-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 12px;
}

.panel-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.panel-model-name {
  font-size: 15px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.panel-close {
  cursor: pointer;
  font-size: 18px;
  color: var(--color-text-muted);
  padding: 4px;
  border-radius: 4px;
  flex-shrink: 0;
}

.panel-close:hover {
  background: var(--color-muted);
  color: var(--color-text-primary);
}

.panel-section {
  margin-bottom: 14px;
}

.panel-section-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 8px;
  padding-bottom: 6px;
  border-bottom: 1px solid var(--color-border);
}

.panel-loading {
  text-align: center;
  padding: 12px;
  font-size: 13px;
  color: var(--color-text-muted);
}

.topo-relation-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  cursor: pointer;
  border-radius: 4px;
  padding: 4px 6px;
  transition: background 0.15s ease;
}

.topo-relation-row:hover {
  background: var(--color-muted);
}

.topo-rel-direction {
  font-size: 11px;
  font-weight: 600;
  padding: 1px 6px;
  border-radius: 3px;
  flex-shrink: 0;
}

.topo-rel-direction.upstream {
  background: #DBEAFE;
  color: #2563EB;
}

.topo-rel-direction.downstream {
  background: #D1FAE5;
  color: #059669;
}

.topo-rel-type {
  color: var(--color-primary);
  font-weight: 500;
}

.topo-rel-target {
  font-weight: 500;
  color: var(--color-text-primary);
}

.topo-slide-enter-active,
.topo-slide-leave-active {
  transition: all 0.25s ease;
}

.topo-slide-enter-from,
.topo-slide-leave-to {
  transform: translateX(100%);
  opacity: 0;
}

/* ===== 接口调试 ===== */
.api-test-layout {
  display: flex;
  gap: 16px;
  width: 100%;
  height: var(--page-height);
  align-self: stretch;
}

.api-left-panel {
  width: 320px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 0;
  overflow-y: auto;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
}

.api-panel-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-primary);
  padding: 12px 16px 8px;
}

.api-endpoint-groups {
  padding: 0 8px;
}

.api-group-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-secondary);
  padding: 8px;
  cursor: pointer;
  border-radius: 6px;
}

.api-group-title:hover {
  background: var(--color-muted);
}

.api-group-items {
  padding: 0 0 4px 0;
}

.api-endpoint-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px 6px 20px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  color: var(--color-text-secondary);
  transition: all 0.15s;
}

.api-endpoint-item:hover {
  background: var(--color-muted);
  color: var(--color-text-primary);
}

.api-endpoint-item.active {
  background: var(--color-primary);
  color: #fff;
}

.api-method-tag {
  font-size: 10px;
  font-weight: 700;
  padding: 1px 5px;
  border-radius: 3px;
  flex-shrink: 0;
}

.api-method-tag.get {
  background: #DBEAFE;
  color: #1D4ED8;
}

.api-method-tag.post {
  background: #FEF3C7;
  color: #92400E;
}

.api-endpoint-item.active .api-method-tag {
  background: rgba(255,255,255,0.25);
  color: #fff;
}

.api-ep-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.api-param-section {
  border-top: 1px solid var(--color-border);
  margin-top: 4px;
}

.api-param-form {
  padding: 0 16px 16px;
}

.api-param-row {
  margin-bottom: 12px;
}

.api-param-label {
  display: block;
  font-size: 12px;
  font-weight: 500;
  color: var(--color-text-secondary);
  margin-bottom: 4px;
}

.api-required {
  color: #EF4444;
}

.api-condition-builder {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.api-cond-row {
  display: flex;
  align-items: center;
  gap: 4px;
}

.api-json-textarea :deep(textarea) {
  font-family: 'Fira Code', 'Cascadia Code', 'JetBrains Mono', monospace;
  font-size: 12px;
  line-height: 1.5;
}

/* 接口调试 右侧 */
.api-right-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-width: 0;
}

.api-request-preview {
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
  overflow: hidden;
}

.api-preview-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
}

.api-preview-header .api-panel-title {
  padding: 12px 0 8px;
}

.api-preview-url {
  padding: 0 16px 8px;
  font-family: 'Fira Code', 'Cascadia Code', 'JetBrains Mono', monospace;
  font-size: 13px;
  color: var(--color-primary);
  word-break: break-all;
}

.api-preview-body {
  margin: 0;
  padding: 12px 16px;
  background: #1E293B;
  color: #E2E8F0;
  font-family: 'Fira Code', 'Cascadia Code', 'JetBrains Mono', monospace;
  font-size: 12px;
  line-height: 1.6;
  overflow-x: auto;
  max-height: 200px;
}

.api-response-section {
  flex: 1;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.api-response-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
}

.api-response-header .api-panel-title {
  padding: 12px 0 8px;
}

.api-response-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.api-response-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}

.api-response-time {
  font-size: 12px;
  color: var(--color-text-muted);
}

.api-response-empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-muted);
  font-size: 13px;
}

.api-response-body {
  flex: 1;
  margin: 0;
  padding: 12px 16px;
  background: #1E293B;
  color: #E2E8F0;
  font-family: 'Fira Code', 'Cascadia Code', 'JetBrains Mono', monospace;
  font-size: 12px;
  line-height: 1.6;
  overflow: auto;
}
</style>
