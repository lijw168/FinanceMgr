package options

type CommResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type BaseOptions struct {
	ID int
}

// type ModifyAttributeOptions struct {
// 	ID          string
// 	Name        *string
// 	Description *string
// }

type ListOptions struct {
	Filter map[string]interface{}
	Offset int
	Limit  int
	Orders map[string]int
}

// type UpdateStatusOptions struct {
// 	ID     string
// 	Status int
// }
