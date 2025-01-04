package models

// SearchDirectSql
// @Description: 直接搜索sql
type SearchDirectSql struct {
	BaseModel
	Uuid string `gorm:"column:uuid;type:string;size:255;not null;comment:uuid" json:"uuid"`
	Name string `gorm:"column:name;type:string;size:255;not null;comment:名称" json:"name"`
	Sql  string `gorm:"column:sql;type:text;comment:sql语句" json:"sql"`
}

func (m *SearchDirectSql) TableName() string {
	return "search_direct_sql"
}
