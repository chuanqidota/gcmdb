# openapi 使用说明

## 简介

- cmdb的数据写入和读取

## 写入

### 接口uri示例说明

```
https://***/v1/openapi/动作
如：
https://***/openapi/v1/instance/zeus/host/create
```

- 接口请求方法都为 POST
- 动作有 create、update、delete，对应增、改、删
- create 传参

``` json script 
{
    "data": [
        {
            "inner_ip": "192.168.1.140",
            "internet_ip": "1.1.1.140",
            "hostname": "主机140",
            "system_type": "windows",
            "host_number": 50,
            "appid": "zeus"
        },{
            "inner_ip": "192.168.1.141",
            "internet_ip": "1.1.1.141",
            "hostname": "主机141",
            "system_type": "linux",
            "appid": "zeus"
        }
    ]
}

```

- create带数据关系 传参
``` json script
传参说明
"__related": {"id": [123, 124, 125]}
"__related": {"info": "asset", "condition": [{"ip": "1.1.1.1", "os": "linux"}]}   # 根据条件匹配一条或多条
```

``` json script 
案例：
{
    "data": [
        {
            "inner_ip": "192.168.1.140",
            "internet_ip": "1.1.1.140",
            "hostname": "主机140",
            "system_type": "windows",
            "host_number": 50,
            "appid": "zeus",
            "__related": {"id": [123, 124, 125]}  # 与id为123的数据实例建立关系
        },{
            "inner_ip": "192.168.1.141",
            "internet_ip": "1.1.1.141",
            "hostname": "主机141",
            "system_type": "linux",
            "appid": "zeus",
            "__related": {"info_id": 2, "condition": [{"ip": "1.1.1.1", "os": "linux"}]}
        }
    ]
}

```

- update 传参（__condition为更新条件，其它元素为待更新字段和更新值）

``` json script 
{
    "data": [
        {
            "__condition": {"uuid": "0b377654-143a-41ed-9a74-4190b30da4c4"}, 
            "inner_ip": "192.168.1.144",
            "host_number": 55，
            "__related": {"id": [123, 124, 125]}
        },{
            "__condition": {"hostname": "主机141", "internet_ip": "1.1.1.141"},
            "inner_ip": "192.168.1.212",
            "__related": {"info_id": "asset", "condition": [{"ip": "1.1.1.1", "os": "linux"}]}
        }
    ]
}
```

- delete 传参

``` json script 
{
    "data": [
        {"uuid": "0b377654-143a-41ed-9a74-4190b30da4c4"},
        {"inner_ip": "192.168.1.212", "host_number": 25}
    ]
}
```

## 读取

### 接口uri示例说明

```
https://***/openapi/v1/search/租户/模型
如：
https://***/openapi/v1/search/zbase/product
```

- 接口请求方法为 POST

### 查询条件传参示例

``` json script
{
    "page": 1,
    "pageSize": 5,
    "order": ["name", "-appid"],
    "and_condition": {"dept_name": "业务中台", "owner_name": "张三"},
    "or_condition": {"second_owner_name": "李四"},
    "in_condition": {"name": ["product-config-web", "base-center-schedule"]},
    "like_condition": {"owner_name": "张"},
    "not_condition": {"owner_name": "张三"}
}
```

### 查询条件说明

- page 和 pageSize 用于分页，默认返回第1页，每页15条数据
- order 用于排序，-表示降序
- and_condition 与查询条件
- or_condition 或查询条件
- in_condition 范围查询
- like_condition 模糊查询
- not_condition 不等于查询


## 数据关系

### 接口uri示例说明

```
https://***/openapi/v1/related/租户/方法
如：
https://***/openapi/v1/related/zbase/create
```

- 接口请求方法为 POST

``` json script
create
{
    "source_info": "asset",
    "source": 45,
    "target_info": "production",
    "targets": [
        31,
        30
    ]
}
```

``` json script
delete，每次只能删除一条数据关系
{
    "source_info": "asset",
    "source": 45,
    "target_info": "production",
    "target": 32
}
```
