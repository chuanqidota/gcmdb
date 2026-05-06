package utils

import "encoding/json"

// ParseJSONColumns 将 GORM 扫描结果中的 []uint8 和 JSON 字符串反序列化为结构化数据。
// GORM 扫描 MySQL JSON 列到 map[string]interface{} 时，值可能为 []uint8 或 string，
// 需要手动反序列化为 JSON 对象/数组，否则前端收到的是转义字符串而非结构化数据。
func ParseJSONColumns(rows []map[string]interface{}) {
	for _, row := range rows {
		for k, v := range row {
			switch val := v.(type) {
			case []uint8:
				var parsed interface{}
				if json.Unmarshal(val, &parsed) == nil {
					row[k] = parsed
				}
			case string:
				if len(val) > 0 && (val[0] == '{' || val[0] == '[') {
					var parsed interface{}
					if json.Unmarshal([]byte(val), &parsed) == nil {
						row[k] = parsed
					}
				}
			}
		}
	}
}
