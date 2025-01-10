package params

type CommonQuery struct {
	Search string `form:"search" label:"搜索"`
	Offset int    `form:"offset" label:"分页offset"`
	Limit  int    `form:"limit" label:"分页limit"`
}
