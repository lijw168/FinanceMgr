package options

type CommResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type BaseOptions struct {
	ID int
}

type ListOptions struct {
	Filter map[string]interface{}
	Offset int
	Limit  int
	Orders map[string]int
}
