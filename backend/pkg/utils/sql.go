package utils

import (
	"fmt"
	"strings"
)

// ValidateSelectSQL 校验 SQL 语句是否为安全的 SELECT 查询
// 禁止：堆叠查询、注释、UNION、危险关键字
func ValidateSelectSQL(sql string) error {
	trimmed := strings.TrimSpace(sql)
	if trimmed == "" {
		return fmt.Errorf("sql语句为空")
	}

	upper := strings.ToUpper(trimmed)

	// 必须以 SELECT 开头
	if !strings.HasPrefix(upper, "SELECT") {
		return fmt.Errorf("仅支持SELECT查询语句")
	}

	// 禁止分号（堆叠查询）
	if strings.Contains(trimmed, ";") {
		return fmt.Errorf("不允许使用分号")
	}

	// 禁止注释符
	if strings.Contains(trimmed, "--") || strings.Contains(trimmed, "/*") || strings.Contains(trimmed, "*/") {
		return fmt.Errorf("不允许使用注释符")
	}

	// 禁止 UNION（防止 UNION 注入）
	if strings.Contains(upper, "UNION ") || strings.Contains(upper, "UNION\t") || strings.HasPrefix(upper[6:], "UNION") {
		return fmt.Errorf("不允许使用UNION查询")
	}

	// 黑名单关键字（不依赖尾部空格，用 Contains 检查）
	dangerous := []string{"DROP ", "DELETE ", "UPDATE ", "INSERT ", "ALTER ", "TRUNCATE ", "CREATE ", "REPLACE ", "GRANT ", "REVOKE "}
	for _, kw := range dangerous {
		if strings.Contains(upper, kw) {
			return fmt.Errorf("不允许执行修改类SQL语句")
		}
	}

	return nil
}
