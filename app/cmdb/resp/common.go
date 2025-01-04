package resp

type CommonList struct {
	Count   int64       `json:"count"`
	Results interface{} `json:"results"`
}
