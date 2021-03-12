package options

type CommResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type BaseOptions struct {
	Id string
}

type ModifyAttributeOptions struct {
	Id          string
	Name        *string
	Description *string
}

type ListOptions struct {
	Filter map[string]interface{}
	Offset int
	Limit  int
	Orders map[string]int
}

type UpdateStatusOptions struct {
	Id     string
	Status int
}
