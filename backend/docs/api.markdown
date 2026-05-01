# GCMDB API 文档

## 概述

GCMDB 提供两组 API：

| 分组 | 路径前缀 | 用途 | 中间件 |
|------|---------|------|--------|
| CMDB 管理 API | `/v1/cmdb/` | 内部管理操作（全量 CRUD） | CORS + Bearer Token 认证 |
| Open API | `/v1/openapi/` | 对外开放接口（简化操作） | Bearer Token 认证 |

### 认证

所有接口均支持 Bearer Token 认证。若服务端未配置 `api_key`，则跳过认证。

```
Authorization: Bearer <your-api-key>
```

### 响应格式

```json
// 成功 (HTTP 200)
{ "code": 1, "msg": "操作成功", "data": ... }

// 失败 (HTTP 400)
{ "code": 0, "msg": "错误信息", "data": null }
```

### 分页参数

大部分列表接口支持以下查询参数：

| 参数 | 类型 | 说明 |
|------|------|------|
| search | string | 模糊搜索关键词 |
| limit | int | 每页数量，默认 10 |
| offset | int | 偏移量 |

---

## 一、模型分组 `/v1/cmdb/models-group`

### 创建模型分组

```
POST /v1/cmdb/models-group
```

```json
{
  "alias": "datacenter",
  "name": "数据中心",
  "description": "数据中心分组",
  "order": 0
}
```

### 查询模型分组列表

```
GET /v1/cmdb/models-group?search=&limit=10&offset=0
```

响应 `data`：

```json
{
  "count": 5,
  "results": [
    { "id": 1, "alias": "datacenter", "name": "数据中心", "description": "...", "order": 0, "created_at": "2026-04-30 10:00:00" }
  ]
}
```

### 查询模型分组详情

```
GET /v1/cmdb/models-group/:id
```

响应 `data`（含分组下的模型列表）：

```json
{
  "id": 1, "alias": "datacenter", "name": "数据中心",
  "models": [
    { "id": 1, "alias": "host", "name": "主机", "group_id": 1, "is_usable": true }
  ]
}
```

### 修改模型分组

```
PATCH /v1/cmdb/models-group/:id
```

```json
{ "name": "新名称", "description": "新描述" }
```

### 删除模型分组

```
DELETE /v1/cmdb/models-group/:id
```

> 若分组下存在模型，将拒绝删除。

---

## 二、模型 `/v1/cmdb/models`

### 创建模型

```
POST /v1/cmdb/models
```

```json
{
  "group_id": 1,
  "alias": "host",
  "name": "主机",
  "description": "物理主机",
  "is_usable": true,
  "icon": "server",
  "order": 0
}
```

### 查询模型列表

```
GET /v1/cmdb/models?search=&limit=10&offset=0
```

### 查询模型详情

```
GET /v1/cmdb/models/:id
```

响应 `data`（含字段分组及字段定义）：

```json
{
  "id": 1, "group_id": 1, "alias": "host", "name": "主机", "is_usable": true,
  "groups": [
    {
      "groups": { "id": 10, "name": "基本信息" },
      "fields": [
        { "id": 100, "alias": "hostname", "name": "主机名", "type": "string", "is_required": true }
      ]
    }
  ]
}
```

### 修改模型

```
PUT /v1/cmdb/models/:id
```

```json
{ "name": "主机", "description": "...", "is_usable": true, "icon": "server", "order": 0 }
```

### 修改模型所属分组

```
PATCH /v1/cmdb/models/:id
```

```json
{ "group_id": 2, "model_id": 1 }
```

### 删除模型

```
DELETE /v1/cmdb/models/:id
```

> 若模型存在关联关系、字段关联、实例关联或实例，将拒绝删除。

---

## 三、模型字段 `/v1/cmdb/models-field`

### 创建字段

```
POST /v1/cmdb/models-field
```

```json
{
  "model_id": 1,
  "field_group_id": 10,
  "alias": "hostname",
  "name": "主机名",
  "type": "string",
  "is_required": true,
  "order": 0
}
```

支持的字段类型：`string`、`number`、`bool`、`date`、`datetime`、`json`

> 创建字段后会异步为该模型的所有已有实例注入默认值。

### 查询模型字段

```
GET /v1/cmdb/models-field/:model_id
```

### 修改字段

```
PUT /v1/cmdb/models-field/:id
```

```json
{ "field_group_id": 10, "name": "新名称", "is_required": false, "order": 1 }
```

> `alias` 和 `type` 创建后不可修改。

### 删除字段

```
DELETE /v1/cmdb/models-field/:id
```

> 若字段被字段关联或唯一约束引用，将拒绝删除。删除后异步从所有实例中移除该字段。

---

## 四、字段分组 `/v1/cmdb/models-field-group`

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/v1/cmdb/models-field-group` | 创建字段分组 |
| GET | `/v1/cmdb/models-field-group/:id` | 查询字段分组详情（含字段列表） |
| PUT | `/v1/cmdb/models-field-group/:id` | 修改字段分组 |
| DELETE | `/v1/cmdb/models-field-group/:id` | 删除字段分组（需无字段） |

### 创建字段分组

```json
{ "model_id": 1, "name": "基本信息", "order": 0 }
```

---

## 五、唯一约束 `/v1/cmdb/models-field-unique`

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/v1/cmdb/models-field-unique` | 创建唯一约束 |
| GET | `/v1/cmdb/models-field-unique/:model_id` | 查询模型的唯一约束 |
| DELETE | `/v1/cmdb/models-field-unique/:id` | 删除唯一约束 |

### 创建唯一约束

```json
{ "model_id": 1, "fields": "hostname,ip", "description": "主机名+IP 唯一" }
```

> `fields` 为逗号分隔的字段别名。若模型已有实例，将拒绝创建。

---

## 六、模型关系 `/v1/cmdb/models-relation`

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/v1/cmdb/models-relation` | 创建模型关系 |
| GET | `/v1/cmdb/models-relation?source_id=` | 查询模型关系 |
| DELETE | `/v1/cmdb/models-relation/:id` | 删除模型关系 |

### 创建模型关系

```json
{ "source_id": 1, "target_id": 2, "type_id": 1, "description": "主机运行应用" }
```

### 查询响应

```json
{
  "count": 2,
  "results": [
    {
      "id": 1, "source_id": 1, "target_id": 2, "type_id": 1,
      "source_display": "主机(host)",
      "target_display": "应用(app)",
      "type_display": "运行"
    }
  ]
}
```

---

## 七、关系类型 `/v1/cmdb/models-relation-type`

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/v1/cmdb/models-relation-type` | 创建关系类型 |
| GET | `/v1/cmdb/models-relation-type?search=&limit=&offset=` | 查询关系类型 |
| PUT | `/v1/cmdb/models-relation-type/:id` | 修改关系类型 |
| DELETE | `/v1/cmdb/models-relation-type/:id` | 删除关系类型（需无引用） |

### 创建关系类型

```json
{ "name": "运行", "s2t": "运行在", "t2s": "被运行于" }
```

---

## 八、字段关联 `/v1/cmdb/models-field-relation`

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/v1/cmdb/models-field-relation` | 创建字段关联 |
| GET | `/v1/cmdb/models-field-relation/:source_model_id` | 查询字段关联 |
| DELETE | `/v1/cmdb/models-field-relation/:id` | 删除字段关联 |

### 创建字段关联

```json
{ "source_model_id": 1, "target_model_id": 2, "source_field_id": 100, "target_field_id": 200 }
```

### 查询响应（含展示名）

```json
[
  {
    "id": 1, "source_model_id": 1, "target_model_id": 2,
    "source_field_id": 100, "target_field_id": 200,
    "source_model_display": "主机(host)",
    "target_model_display": "应用(app)",
    "source_field_display": "主机名(hostname)",
    "target_field_display": "部署主机(host_field)"
  }
]
```

---

## 九、实例 `/v1/cmdb/instance`

### 创建实例

```
POST /v1/cmdb/instance
```

```json
{
  "model_id": 1,
  "model_alias": "host",
  "data": {
    "hostname": "web-01",
    "ip": "10.0.0.1",
    "is_active": true
  }
}
```

> `data` 的 key 必须与模型字段的 `alias` 对应，创建时会校验字段类型和必填。

### 查询实例列表

```
GET /v1/cmdb/instances/:model_id?search=&field=&value=&compare=&limit=10&offset=0
```

支持 4 种搜索模式：

| 参数组合 | 说明 |
|---------|------|
| 无参数 | 返回全部实例，分页 |
| `search` | 模糊搜索所有字段 |
| `field` + `value` + `compare` | 指定字段精确筛选（`=`/`like`/`!=`） |
| `search` + `field` + `value` + `compare` | 组合搜索 |

响应：

```json
{
  "count": 100,
  "results": [
    { "id": 1, "model_id": 1, "model_alias": "host", "data": { "hostname": "web-01" } }
  ]
}
```

### 查询实例详情

```
GET /v1/cmdb/instance/:id
```

### 更新实例

```
PUT /v1/cmdb/instance/:id
```

```json
{ "data": { "hostname": "web-02", "ip": "10.0.0.2" } }
```

### 删除实例

```
DELETE /v1/cmdb/instance/:id
```

> 同时删除该实例的所有关联关系。

---

## 十、实例关系 `/v1/cmdb/instance-relation`

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/v1/cmdb/instance-relation` | 创建实例关系 |
| GET | `/v1/cmdb/source-target-relation/:source_id` | 查询源实例的目标关系 |
| GET | `/v1/cmdb/target-source-relation/:target_id` | 查询目标实例的源关系 |
| DELETE | `/v1/cmdb/instance-relation/:id` | 删除实例关系 |

### 创建实例关系

```json
{
  "source_model_id": 1,
  "target_model_id": 2,
  "source_instance_id": 10,
  "target_instance_id": 20
}
```

### 查询响应（按模型分组）

```json
[
  { "modelId": 2, "modelDisplay": "应用", "count": 3, "instances": [20, 21, 22] }
]
```

---

## 十一、SQL 查询 `/v1/cmdb/search-direct-sql`

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/v1/cmdb/search-direct-sql` | 保存查询 |
| GET | `/v1/cmdb/search-direct-sql` | 查询列表 |
| GET | `/v1/cmdb/search-direct-sql/:id` | 执行查询 |
| PUT | `/v1/cmdb/search-direct-sql/:id` | 修改查询 |
| DELETE | `/v1/cmdb/search-direct-sql/:id` | 删除查询 |

### 保存查询

```json
{ "name": "所有主机", "sql": "SELECT * FROM instance WHERE model_id = 1" }
```

> 仅允许 SELECT 语句。执行结果上限 1000 行。

---

## 十二、Open API `/v1/openapi/`

### 查询模型

```
GET /v1/openapi/model/all
GET /v1/openapi/model/single?id=1
```

### 实例操作

```
POST /v1/openapi/instance/:action
```

支持的 action：`create`、`update`、`delete`、`mul_delete`、`search`、`fulltext`、`direct`

#### search 请求体

```json
{
  "model": "host",
  "fields": ["hostname", "ip"],
  "__condition": {
    "limit": 10,
    "offset": 0,
    "order": ["-created_at"],
    "where": [
      { "search": "web" },
      { "eq": { "is_active": true } },
      { "contains": { "hostname": "prod" } },
      { "in": { "ip": ["10.0.0.1", "10.0.0.2"] } },
      { "or": [
        { "eq": { "system_type": "linux" } },
        { "eq": { "system_type": "windows" } }
      ]}
    ]
  }
}
```

where 支持的操作符：

| 操作符 | 说明 | 值类型 |
|--------|------|--------|
| search | 全局模糊搜索 | string |
| eq | 等于 | object |
| ne / not | 不等于 | object |
| gt | 大于 | object |
| ge | 大于等于 | object |
| lt | 小于 | object |
| le | 小于等于 | object |
| in | 范围查询 | object (值为数组) |
| contains | 包含 | object |
| startswith | 开头匹配 | object |
| endswith | 结尾匹配 | object |
| or | 或条件 | array |

排序：`["-created_at", "id"]`，`-` 前缀表示降序。

### 实例关系操作

```
POST /v1/openapi/instance-relation/create
POST /v1/openapi/instance-relation/delete
```

---

## 十三、任务

```
POST /v1/cmdb/sync-instance-relation
```

同步实例关系（后台任务）。
