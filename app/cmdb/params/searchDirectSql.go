package params

type UpdateSearchDirectSql struct {
	Name string `json:"name" label:"名称" binding:"required"`
	Sql  string `json:"sql" label:"sql语句" binding:"required"`
}
