package params

type CommonQuery struct {
	Search string `form:"search"`
	Offset int    `form:"offset"`
	Limit  int    `form:"limit"`
}

type CommonList struct {
	Count   int64       `json:"count"`
	Results interface{} `json:"results"`
}
