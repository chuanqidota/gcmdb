package resp

type CommonList struct {
	Count   int64       `json:"count" label:"数量"`
	Results interface{} `json:"results" label:"结果"`
}
