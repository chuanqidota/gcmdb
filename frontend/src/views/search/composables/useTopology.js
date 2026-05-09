import { ref, computed, nextTick } from 'vue'
import { Graph } from '@antv/g6'
import { getModelDetail, getInstanceTopology } from '../../../api/search'
import { GROUP_COLORS } from '../../../utils/colors'

export function useTopology(allModels, allModelGroups, inst, instFields) {
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

  function buildNodeLabel(instData, fields) {
    const lines = [`#${instData.id}`]
    const displayFields = fields.filter(f => f.alias !== 'id').slice(0, 4)
    for (const f of displayFields) {
      const val = instData[f.alias]
      if (val !== undefined && val !== null && val !== '') {
        const short = String(val).length > 16 ? String(val).slice(0, 16) + '...' : String(val)
        lines.push(`${f.name || f.alias}: ${short}`)
      }
    }
    return lines.join('\n')
  }

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
    const nodes = instances.map(instData => {
      const nodeId = `inst-${instData.id}`
      instanceIdSet.add(instData.id)
      topoNodeMap.set(nodeId, { ...instData, model_id: model.id, model_name: model.name, model_alias: model.alias })
      return {
        id: nodeId,
        data: { label: buildNodeLabel(instData, fields), instanceId: instData.id, modelId: model.id },
        style: { fill: topoModelColorMap[model.id], stroke: '#fff', lineWidth: 2, radius: 12 },
      }
    })

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

  async function openTopoForInstance(instanceId, modelId, modelName, modelAlias) {
    const model = allModels.value.find(m => m.id === modelId)
    if (!model) return
    destroyTopoGraph()
    topoSelectedInstance.value = null

    const groupIdx = allModelGroups.value.findIndex(g => g.id === model.group_id)
    topoModelColorMap[modelId] = GROUP_COLORS[(groupIdx >= 0 ? groupIdx : 0) % GROUP_COLORS.length]

    const nodeId = `inst-${instanceId}`
    topoNodeMap.set(nodeId, { id: instanceId, model_id: modelId, model_name: modelName, model_alias: modelAlias })
    const nodes = [{
      id: nodeId,
      data: { label: `#${instanceId}\n${modelName}`, instanceId, modelId },
      style: { fill: topoModelColorMap[modelId], stroke: '#fff', lineWidth: 2, radius: 12 },
    }]

    try {
      const res = await getInstanceTopology(instanceId)
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

      topoDrawerTitle.value = `实例拓扑 — ${modelName} #${instanceId}`
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

  return {
    topoDrawerVisible, topoDrawerTitle, topoGraphSearchText, topoSelectedInstance,
    topoGraphContainer, topoNodeCount, topoFieldRows, topoRelations,
    openBatchTopo, openSingleTopo, openTopoForInstance, onTopoGraphSearch, focusTopoNode, onTopoDrawerClose,
  }
}
