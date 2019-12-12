package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

const (
	SUCCESS = "200"
	BAD_PARAMS  = "40400001"
	CREATE_SC_ERROR = "50000002"


)

type Response struct {
	Code 		string 		`json:"code"`
	Message 	string 		`json:"message"`
	Data 		interface{}	`json:"data"`
}

type PageList struct {
	Size 		int 		`json:"size"`
	Offset 		int 		`json:"offset"`
	Total 		int 		`json:"total"`
	DataList 	interface{}	`json:"list"`
}
type PageQueryResponse struct {
	PageNo   int `json:"page_no"`
	PageSize int `json:"page_size"`
    Count int `json:"count"`
	List interface{}  `json:"list"`
}

//var MsgFlags map[int]string = {
//
//}


func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := uuid.NewV4()
		c.Writer.Header().Set("X-Request-Id", uuid.String())
		c.Next()
	}
}
