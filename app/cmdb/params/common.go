package params

type CommonQuery struct {
	Search string `form:"search"`
	Offset int    `form:"offset"`
	Limit  int    `form:"limit"`
}
