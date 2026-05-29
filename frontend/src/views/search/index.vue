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
    <div v-show="activeTab !== 'api' && activeTab !== 'llm'" class="search-hero" :class="{ 'hero-column': activeTab === 'instance' }">
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

      <!-- 实例搜索 -->
      <template v-if="activeTab === 'instance'">
        <div class="search-box">
          <el-select
            v-model="inst.modelId"
            placeholder="选择模型..."
            filterable
            class="search-model-select"
            @change="onModelChange"
          >
            <el-option
              v-for="m in allModels"
              :key="m.id"
              :label="`${m.name} (${m.alias})`"
              :value="m.id"
            />
          </el-select>
          <el-input v-model="inst.globalSearch" placeholder="全局搜索所有字段..." clearable @keyup.enter="doInstanceSearch" class="search-input">
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
          <el-button type="primary" class="search-btn" @click="doInstanceSearch" :loading="inst.loading" :disabled="!inst.modelId">搜索</el-button>
          <el-button class="search-btn" @click="resetInstance">重置</el-button>
        </div>
        <div class="cond-action-row">
          <el-button size="small" @click="inst.conditions.push({ group: 1, field: '', op: 'eq', val: '' })">+ 添加条件</el-button>
          <el-button size="small" @click="addOrGroup">+ 添加 OR 组</el-button>
        </div>
      </template>

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
    <div v-show="activeTab !== 'api' && activeTab !== 'llm'" class="results-area">

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
        <!-- 条件构建器（高级条件） -->
        <div class="condition-builder" v-if="inst.modelId">
          <!-- 条件行 -->
            <template v-for="(idx, di) in sortedIndices" :key="idx">
              <div v-if="di > 0 && (inst.conditions[idx].group || 1) > (inst.conditions[sortedIndices[di - 1]].group || 1)" class="or-divider">
                <span>OR</span>
              </div>
              <div class="condition-row" :class="{ 'or-condition': (inst.conditions[idx].group || 1) > 1 }">
                <el-tag v-if="(inst.conditions[idx].group || 1) === 1" size="small" type="info">AND</el-tag>
                <el-tag v-else size="small" type="warning">OR</el-tag>
                <el-select v-model="inst.conditions[idx].field" placeholder="字段" class="cond-field" :disabled="inst.conditions[idx].op === 'search'">
                  <el-option v-for="f in inst.fields" :key="f.alias" :value="f.alias">
                    <span>{{ f.name }}</span>
                    <el-tag v-if="f.is_indexed" size="small" type="primary" style="margin-left: 6px">索引</el-tag>
                  </el-option>
                </el-select>
                <el-select v-model="inst.conditions[idx].op" class="cond-op" @change="v => { if (v === 'search') inst.conditions[idx].field = '' }">
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
                  <el-option label="模糊搜索" value="search" />
                </el-select>
                <el-input v-model="inst.conditions[idx].val" :placeholder="inst.conditions[idx].op === 'search' ? '全文关键词' : '值'" class="cond-val" />
                <el-button :icon="Delete" circle @click="inst.conditions.splice(idx, 1)" />
                <template v-if="di === sortedIndices.length - 1">
                  <span class="condition-spacer"></span>
                  <el-button size="small" @click="inst.conditions.push({ group: 1, field: '', op: 'eq', val: '' })">+ 添加条件</el-button>
                  <el-button size="small" @click="addOrGroup">+ 添加 OR 组</el-button>
                </template>
              </div>
            </template>
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
                  <RelationPanel
                    :instance-id="row.id"
                    :cache="instRelationCache[row.id]"
                    :tab-change-handler="handleTabChange"
                    @view-topo="(instId, modelId, modelName, modelAlias) => openTopoForInstance(instId, modelId, modelName, modelAlias)"
                    @open-detail="(instId, name, alias) => openRelDrawer(instId, name, alias)"
                  />
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
    <div v-if="activeTab === 'api'" class="api-test-layout" :class="{ resizing: apiResizing }">
      <!-- 左栏：接口列表 -->
      <div class="api-col-endpoints">
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
      </div>

      <!-- 中栏：参数表单 -->
      <div class="api-col-params" :style="{ width: apiParamsWidth + 'px' }">
        <div v-if="apiSelected" class="api-param-section">
          <div class="api-panel-title">参数 — {{ apiSelected.name }}</div>
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
                <template v-for="(idx, di) in sortedApiIndices" :key="idx">
                  <div v-if="di > 0 && (apiConditions[idx].group || 1) > (apiConditions[sortedApiIndices[di - 1]].group || 1)" class="or-divider-sm">
                    <span>OR</span>
                  </div>
                  <div class="api-cond-row" :class="{ 'or-row': (apiConditions[idx].group || 1) > 1 }">
                    <el-tag v-if="(apiConditions[idx].group || 1) === 1" size="small" type="info">AND</el-tag>
                    <el-tag v-else size="small" type="warning">OR</el-tag>
                    <el-select v-model="apiConditions[idx].field" placeholder="字段" size="small" style="width: 100px" filterable :disabled="apiConditions[idx].op === 'search'">
                      <el-option v-for="f in apiModelFields" :key="f.alias" :label="f.name" :value="f.alias" />
                    </el-select>
                    <el-select v-model="apiConditions[idx].op" placeholder="操作符" size="small" style="width: 110px" @change="v => { if (v === 'search') apiConditions[idx].field = '' }">
                      <el-option v-for="op in apiOperators" :key="op.value" :label="op.label" :value="op.value" />
                    </el-select>
                    <el-input v-model="apiConditions[idx].value" :placeholder="apiConditions[idx].op === 'search' ? '全文关键词' : '值'" size="small" style="flex: 1" />
                    <el-button link type="danger" size="small" @click="apiConditions.splice(idx, 1)">
                      <el-icon><Delete /></el-icon>
                    </el-button>
                  </div>
                </template>
                <div style="display: flex; gap: 6px; margin-top: 4px;">
                  <el-button size="small" @click="apiConditions.push({ group: 1, field: '', op: 'eq', value: '' })">
                    <el-icon><Plus /></el-icon> 添加条件
                  </el-button>
                  <el-button size="small" @click="addApiOrGroup">+ 添加 OR 组</el-button>
                </div>
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
        <div v-else class="api-param-empty">
          <el-icon size="40" color="var(--color-text-muted)"><ArrowLeft /></el-icon>
          <span>请在左侧选择接口</span>
        </div>
      </div>

      <!-- 拖拽调宽手柄 -->
      <div
        class="api-resize-handle"
        :class="{ active: apiResizing }"
        @mousedown="onApiResizeStart"
      />

      <!-- 右栏：请求预览 + 响应 -->
      <div class="api-col-response">
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

    <!-- 大模型检索 -->
    <div v-if="activeTab === 'llm'" class="llm-search-layout">
      <!-- 搜索输入区 -->
      <div class="llm-search-box">
        <el-input
          v-model="llm.message"
          placeholder="用自然语言描述你想查询的数据... 例如：找出所有交换机的管理IP"
          clearable
          @keyup.enter="doLlmSearch"
          class="llm-search-input"
          :disabled="llm.loading"
        >
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <el-button type="primary" class="llm-search-btn" @click="doLlmSearch" :loading="llm.loading">
          发送
        </el-button>
        <el-button class="llm-search-btn" @click="doClearSession" :disabled="llm.loading">
          新对话
        </el-button>
      </div>

      <!-- 模型范围选择 -->
      <div class="llm-model-filter">
        <span class="llm-filter-label">查询范围：</span>
        <el-select
          v-model="llm.modelFilter"
          placeholder="全部模型"
          multiple
          filterable
          collapse-tags
          collapse-tags-tooltip
          clearable
          size="small"
          class="llm-model-select"
        >
          <el-option v-for="m in allModels" :key="m.id" :label="m.name" :value="m.alias" />
        </el-select>
        <span v-if="llm.tokenUsed > 0" class="llm-token-info">
          Token: {{ llm.tokenUsed }}
        </span>
      </div>

      <!-- 加载状态 -->
      <div v-if="llm.loading" class="llm-loading">
        <el-skeleton :rows="3" animated />
      </div>

      <!-- 思考过程 -->
      <div v-if="llm.thinking && !llm.loading" class="llm-thinking">
        <div class="llm-section-title">思考过程</div>
        <div class="llm-thinking-content">{{ llm.thinking }}</div>
      </div>

      <!-- AI 分析 -->
      <div v-if="llm.aiAnalysis && !llm.loading" class="llm-analysis">
        <div class="llm-section-title">AI 分析</div>
        <div class="llm-analysis-content">{{ llm.aiAnalysis }}</div>
      </div>

      <!-- 查询结果 -->
      <div v-if="llm.results.length > 0 && !llm.loading" class="llm-results">
        <template v-for="(result, idx) in llm.results" :key="idx">
          <div class="llm-result-section">
            <div class="llm-result-header">
              <span class="llm-result-model">{{ result.model_name }}</span>
              <el-tag size="small" type="info">{{ result.model_alias }}</el-tag>
              <span class="llm-result-count">共 {{ result.count }} 条</span>
            </div>
            <el-table
              v-if="result.data && result.data.length > 0"
              :data="result.data"
              stripe
              size="small"
              highlight-current-row
              show-overflow-tooltip
              max-height="400"
            >
              <el-table-column
                v-for="(val, key) in result.data[0]"
                :key="key"
                :prop="String(key)"
                :label="String(key)"
                :min-width="120"
              />
            </el-table>
            <el-empty v-else description="无数据" :image-size="60" />
          </div>
        </template>
      </div>

      <!-- 无结果提示 -->
      <div v-if="!hasResults && !llm.loading && llm.history.length === 0" class="llm-empty">
        <el-empty description="用自然语言描述你想查询的数据，例如：找出所有交换机的管理IP" />
      </div>

      <!-- 对话历史 -->
      <div v-if="llm.history.length > 0" class="llm-history">
        <div class="llm-history-toggle" @click="llm.showHistory = !llm.showHistory">
          <el-icon><ArrowDown v-if="!llm.showHistory" /><ArrowUp v-else /></el-icon>
          对话历史 ({{ llm.history.length }})
        </div>
        <transition name="llm-slide">
          <div v-show="llm.showHistory" class="llm-history-list">
            <div v-for="(item, idx) in llm.history" :key="idx" class="llm-history-item">
              <div class="llm-history-question">
                <span class="llm-history-role">Q:</span>
                {{ item.message }}
                <span class="llm-history-time">{{ item.timestamp }}</span>
              </div>
              <div v-if="item.aiAnalysis" class="llm-history-answer">
                <span class="llm-history-role">A:</span>
                {{ item.aiAnalysis }}
              </div>
            </div>
          </div>
        </transition>
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

  <!-- 关联实例抽屉 -->
  <el-drawer v-model="relDrawerVisible" :title="relDrawerCurrent ? `${relDrawerCurrent.modelName} #${relDrawerCurrent.instanceId}` : '关联实例'" size="85%" direction="rtl" @close="closeRelDrawer">
    <div class="rel-drawer-content">
      <!-- 面包屑导航 -->
      <div v-if="relDrawerStack.length > 1" class="rel-breadcrumb">
        <span
          v-for="(item, idx) in relDrawerStack"
          :key="idx"
          class="rel-breadcrumb-item"
          :class="{ active: idx === relDrawerStack.length - 1 }"
          @click="idx < relDrawerStack.length - 1 && popToRelDrawer(idx)"
        >
          {{ item.modelName }} #{{ item.instanceId }}
          <el-icon v-if="idx < relDrawerStack.length - 1" class="rel-breadcrumb-sep"><ArrowRight /></el-icon>
        </span>
      </div>

      <!-- 当前实例字段值 -->
      <div v-if="relDrawerCurrent" class="rel-detail-section">
        <div class="rel-section-title">字段值</div>
        <div v-loading="relDrawerCurrent.loading">
          <el-table v-if="relDrawerCurrent.fields.length" :data="relDrawerFieldRows" size="small" stripe highlight-current-row max-height="300">
            <el-table-column prop="name" label="字段" width="150" />
            <el-table-column prop="alias" label="别名" width="120" />
            <el-table-column prop="value" label="值" show-overflow-tooltip />
          </el-table>
          <div v-else-if="!relDrawerCurrent.loading" class="rel-empty">暂无字段</div>
        </div>
      </div>

      <!-- 关联列表 -->
      <div v-if="relDrawerCurrent" class="rel-detail-section">
        <div class="rel-section-title">关联实例</div>
        <RelationPanel
          :instance-id="relDrawerCurrent.instanceId"
          :cache="relDrawerRelationCache[relDrawerCurrent.instanceId]"
          :tab-change-handler="handleRelDrawerTabChange"
          @view-topo="(instId, modelId, modelName, modelAlias) => openTopoForInstance(instId, modelId, modelName, modelAlias)"
          @open-detail="(instId, name, alias) => pushRelDrawer(instId, name, alias)"
        />
      </div>
    </div>
  </el-drawer>

</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { Search, Delete, CircleCheckFilled, ArrowDown, ArrowRight, ArrowLeft, Close, CopyDocument, ArrowUp } from '@element-plus/icons-vue'
import { getAllModels, getModelGroups, getModelRelationTypes } from '../../api/search'
import { useFulltext } from './composables/useFulltext'
import { useInstance } from './composables/useInstance'
import { useModelBrowse } from './composables/useModelBrowse'
import { useApiDebug } from './composables/useApiDebug'
import { useTopology } from './composables/useTopology'
import { useLlmSearch } from './composables/useLlmSearch'
import RelationPanel from './components/RelationPanel.vue'

const typeTag = { string: '', number: 'success', bool: 'warning', date: 'info', datetime: 'info', json: 'danger' }

const tabs = [
  { name: 'fulltext', label: '全文检索' },
  { name: 'instance', label: '实例搜索' },
  { name: 'model', label: '模型浏览' },
  { name: 'api', label: '接口调试' },
  { name: 'llm', label: '大模型检索' },
]

const activeTab = ref('fulltext')
const allModels = ref([])
const allModelGroups = ref([])
const allRelationTypes = ref([])

function switchTab(name) {
  activeTab.value = name
}

// Composables
const { ft, doFulltextSearch } = useFulltext(allModels)
const {
  inst, instRelationCache, instRelations,
  onModelChange, doInstanceSearch, calcColWidth, resetInstance, addOrGroup,
  onInstExpandChange, onInstTabChange, handleTabChange,
  relDrawerVisible, relDrawerStack, relDrawerCurrent, relDrawerRelationCache,
  openRelDrawer, pushRelDrawer, popRelDrawer, popToRelDrawer, closeRelDrawer,
  handleRelDrawerTabChange,
} = useInstance(allModels)

// 按 group 排序后的条件索引数组（保留原数组引用，v-model 不断裂）
const sortedIndices = computed(() => {
  const arr = inst.value.conditions
  return arr.map((_, i) => i).sort((a, b) => (arr[a].group || 1) - (arr[b].group || 1))
})
const { modelSearch, expandedModelId, modelFieldsCache, groupedModels, toggleFields, resolveRelation } = useModelBrowse(allModels, allModelGroups, allRelationTypes)
const {
  apiEndpointGroups, apiSelectedKey, apiSelected, apiParamValues, apiConditions,
  apiSending, apiRequestUrl, apiRequestBody, apiResponseData, apiResponseStatus, apiResponseTime,
  apiOperators, apiModelFields, apiSavedSqls, addApiOrGroup, sortedApiIndices,
  getApiOptions, selectApiEndpoint, sendApiRequest, copyApiRequest, copyApiResponse,
  isApiReadOnly, loadApiSavedSqls, formatJson,
} = useApiDebug(allModels)
const {
  topoDrawerVisible, topoDrawerTitle, topoGraphSearchText, topoSelectedInstance,
  topoGraphContainer, topoNodeCount, topoFieldRows, topoRelations,
  openBatchTopo, openSingleTopo, openTopoForInstance, onTopoGraphSearch, focusTopoNode, onTopoDrawerClose,
} = useTopology(allModels, allModelGroups, inst)

const { llm, hasResults, doLlmSearch, doClearSession } = useLlmSearch(allModels)

// 接口调试 — 中栏拖拽调宽
const apiParamsWidth = ref(300)
const apiResizing = ref(false)
let startX = 0
let startWidth = 0

function onApiResizeStart(e) {
  e.preventDefault()
  apiResizing.value = true
  startX = e.clientX
  startWidth = apiParamsWidth.value
  document.addEventListener('mousemove', onApiResizeMove)
  document.addEventListener('mouseup', onApiResizeEnd)
}

function onApiResizeMove(e) {
  const delta = e.clientX - startX
  apiParamsWidth.value = Math.min(600, Math.max(200, startWidth + delta))
}

function onApiResizeEnd() {
  apiResizing.value = false
  document.removeEventListener('mousemove', onApiResizeMove)
  document.removeEventListener('mouseup', onApiResizeEnd)
}

onBeforeUnmount(() => {
  document.removeEventListener('mousemove', onApiResizeMove)
  document.removeEventListener('mouseup', onApiResizeEnd)
})

// 关联抽屉字段行
const relDrawerFieldRows = computed(() => {
  if (!relDrawerCurrent.value) return []
  const fields = relDrawerCurrent.value.fields || []
  const data = relDrawerCurrent.value.instanceData || {}
  return fields.map(f => ({
    name: f.name,
    alias: f.alias,
    value: data[f.alias] !== undefined && data[f.alias] !== null ? String(data[f.alias]) : '-',
  }))
})

// ===== 初始化 =====
onMounted(async () => {
  try {
    const [modelsRes, groupsRes, typesRes] = await Promise.all([
      getAllModels(),
      getModelGroups(),
      getModelRelationTypes(),
    ])
    allModels.value = (modelsRes.data || []).filter(m => m.is_usable !== false)
    allModelGroups.value = groupsRes.data || []
    allRelationTypes.value = typesRes.data || []
  } catch {}
  loadApiSavedSqls()
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
  margin-top: 24px;
  margin-bottom: 24px;
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
  margin-bottom: 16px;
}

.search-hero.hero-column {
  flex-direction: column;
  align-items: center;
}

.cond-action-row {
  display: flex;
  gap: 8px;
  margin-top: 8px;
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
  width: 200px;
  flex-shrink: 0;
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

/* 条件构建器 */
.condition-builder {
  margin-bottom: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--color-border);
  max-width: 960px;
  width: 100%;
}


.condition-row {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 6px;
}

.cond-field {
  width: 180px;
}

.cond-op {
  width: 120px;
}

.cond-val {
  flex: 2;
  min-width: 120px;
}

.condition-spacer {
  flex: 1;
}

.or-divider {
  text-align: center;
  padding: 6px 0;
  color: #e6a23c;
  font-weight: bold;
  font-size: 13px;
  letter-spacing: 2px;
}

.or-divider span {
  background: var(--color-muted);
  padding: 2px 16px;
  border-radius: 4px;
  border: 1px dashed #e6a23c;
}

.or-condition {
  background: #fdf6ec;
  border-radius: 6px;
  padding: 4px 8px;
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

/* ===== 接口调试（三栏布局） ===== */
.api-test-layout {
  display: flex;
  gap: 12px;
  width: 100%;
  height: var(--page-height);
  align-self: stretch;
}

.api-test-layout.resizing {
  user-select: none;
  cursor: col-resize;
}

/* 左栏：接口列表 */
.api-col-endpoints {
  width: 180px;
  flex-shrink: 0;
  overflow-y: auto;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
}

/* 中栏：参数表单 */
.api-col-params {
  width: 300px;
  min-width: 200px;
  max-width: 600px;
  flex-shrink: 0;
  overflow-y: auto;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
}

/* 拖拽调宽手柄 */
.api-resize-handle {
  width: 6px;
  cursor: col-resize;
  border-radius: 3px;
  background: transparent;
  transition: background 0.15s;
  flex-shrink: 0;
  position: relative;
}

.api-resize-handle:hover,
.api-resize-handle.active {
  background: var(--color-primary-light-5, #c0c4cc);
}

.api-resize-handle::after {
  content: '';
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  width: 2px;
  height: 24px;
  border-radius: 1px;
  background: var(--color-border);
  opacity: 0;
  transition: opacity 0.15s;
}

.api-resize-handle:hover::after,
.api-resize-handle.active::after {
  opacity: 1;
}

.api-param-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 12px;
  color: var(--color-text-muted);
  font-size: 13px;
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
  gap: 6px;
  padding: 5px 8px 5px 16px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 12px;
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
  margin-top: 0;
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

.or-divider-sm {
  text-align: center;
  padding: 4px 0;
  color: #e6a23c;
  font-weight: bold;
  font-size: 11px;
  letter-spacing: 2px;
}

.or-divider-sm span {
  padding: 1px 10px;
  border-radius: 3px;
  border: 1px dashed #e6a23c;
}

.or-row {
  background: #fdf6ec;
  border-radius: 4px;
  padding: 2px 4px;
}

.api-json-textarea :deep(textarea) {
  font-family: 'Fira Code', 'Cascadia Code', 'JetBrains Mono', monospace;
  font-size: 12px;
  line-height: 1.5;
}

/* 右栏：请求+响应 */
.api-col-response {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-width: 0;
  overflow-y: auto;
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

/* ===== 关联实例抽屉 ===== */
.rel-drawer-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
  height: 100%;
}

.rel-breadcrumb {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
  padding: 8px 12px;
  background: var(--color-muted);
  border-radius: var(--radius-md);
  font-size: 13px;
}

.rel-breadcrumb-item {
  cursor: pointer;
  color: var(--color-primary);
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.rel-breadcrumb-item.active {
  color: var(--color-text-primary);
  font-weight: 600;
  cursor: default;
}

.rel-breadcrumb-sep {
  color: var(--color-text-muted);
  font-size: 12px;
}

.rel-detail-section {
  flex-shrink: 0;
}

.rel-section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 10px;
  padding-bottom: 6px;
  border-bottom: 1px solid var(--color-border);
}

.rel-empty {
  text-align: center;
  padding: 24px;
  font-size: 13px;
  color: var(--color-text-muted);
}

/* ===== 大模型检索 ===== */
.llm-search-layout {
  width: 100%;
  max-width: 960px;
  align-self: center;
}

.llm-search-box {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.llm-search-input {
  flex: 1;
}

.llm-search-input :deep(.el-input__wrapper) {
  border-radius: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  padding: 4px 20px;
  height: 44px;
  transition: box-shadow 0.2s ease;
}

.llm-search-input :deep(.el-input__wrapper:focus-within) {
  box-shadow: 0 2px 16px rgba(37, 99, 235, 0.15);
}

.llm-search-btn {
  border-radius: 24px;
  padding: 0 24px;
  height: 44px;
  font-weight: 600;
  flex-shrink: 0;
}

.llm-model-filter {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  padding: 8px 12px;
  background: var(--color-muted);
  border-radius: var(--radius-md);
}

.llm-filter-label {
  font-size: 13px;
  color: var(--color-text-secondary);
  flex-shrink: 0;
}

.llm-model-select {
  flex: 1;
  max-width: 500px;
}

.llm-token-info {
  font-size: 12px;
  color: var(--color-text-muted);
  margin-left: auto;
}

.llm-loading {
  padding: 16px 0;
}

.llm-thinking,
.llm-analysis {
  margin-bottom: 16px;
  padding: 16px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
}

.llm-section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 10px;
  padding-bottom: 6px;
  border-bottom: 1px solid var(--color-border);
}

.llm-thinking-content {
  font-size: 13px;
  color: var(--color-text-secondary);
  white-space: pre-wrap;
  line-height: 1.6;
}

.llm-analysis-content {
  font-size: 14px;
  color: var(--color-text-primary);
  white-space: pre-wrap;
  line-height: 1.8;
}

.llm-results {
  margin-bottom: 16px;
}

.llm-result-section {
  margin-bottom: 16px;
  padding: 16px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.llm-result-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.llm-result-model {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.llm-result-count {
  font-size: 13px;
  color: var(--color-text-muted);
  margin-left: auto;
}

.llm-empty {
  padding: 60px 0;
}

.llm-history {
  margin-top: 24px;
  border-top: 1px solid var(--color-border);
  padding-top: 16px;
}

.llm-history-toggle {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--color-text-secondary);
  cursor: pointer;
  padding: 8px 0;
  user-select: none;
}

.llm-history-toggle:hover {
  color: var(--color-primary);
}

.llm-history-list {
  margin-top: 8px;
}

.llm-history-item {
  padding: 12px;
  background: var(--color-muted);
  border-radius: var(--radius-md);
  margin-bottom: 8px;
}

.llm-history-question {
  font-size: 14px;
  color: var(--color-text-primary);
  margin-bottom: 6px;
}

.llm-history-answer {
  font-size: 13px;
  color: var(--color-text-secondary);
  padding-left: 20px;
}

.llm-history-role {
  font-weight: 600;
  color: var(--color-primary);
  margin-right: 6px;
}

.llm-history-time {
  font-size: 12px;
  color: var(--color-text-muted);
  margin-left: 8px;
}

.llm-slide-enter-active,
.llm-slide-leave-active {
  transition: all 0.2s ease;
}

.llm-slide-enter-from,
.llm-slide-leave-to {
  opacity: 0;
  max-height: 0;
  overflow: hidden;
}

.llm-slide-enter-to,
.llm-slide-leave-from {
  max-height: 500px;
}
</style>
