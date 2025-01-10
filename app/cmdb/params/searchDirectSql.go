package params

type UpdateSearchDirectSql struct {
	Name string `json:"name" label:"名称"`
	Sql  string `json:"sql" label:"sql语句"`
}
