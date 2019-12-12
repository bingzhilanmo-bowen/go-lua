package controller

import (
	"bingzhilanmo/go-lua/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RunParams struct {
	Keyword string `json:"keyword"`
	Params map[string]string `json:"params"`
}

// @Summary 创建一个新的脚本
// @Tags Keyword
// @Accept  json
// @Produce json
// @Param create body service.KeywordDto true "service.KeywordDto"
// @Success 200 {object} Response
// @Failure 500 {string} string "Internal Error"
// @Router  /api/v1/keyword/created [post]
func CreateNewKeyword(c *gin.Context) {
	resp := &Response{}
	var dto *service.KeywordDto
	err := c.BindJSON(&dto)
	if err != nil {
		resp.Code = BAD_PARAMS
		resp.Message = "params error !!!"
	} else {

		ns, err := service.CreateNewKeyword(dto)
		if err == nil {
			resp.Code = SUCCESS
			resp.Message = "success"
			resp.Data = ns
		} else {
			resp.Code = CREATE_SC_ERROR
			resp.Message = err.Error()
		}
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary 修改脚本
// @Tags Keyword
// @Accept  json
// @Produce json
// @Param create body service.KeywordDto true "service.KeywordDto"
// @Success 200 {object} Response
// @Failure 500 {string} string "Internal Error"
// @Router  /api/v1/keyword/updated [put]
func UpdateKeyword(c *gin.Context) {
	resp := &Response{}
	var dto *service.KeywordDto
	err := c.BindJSON(&dto)
	if err != nil {
		resp.Code = BAD_PARAMS
		resp.Message = "params error !!!"
	} else {

		ns, err := service.UpdateKeyword(dto)
		if err == nil {
			resp.Code = SUCCESS
			resp.Message = "success"
			resp.Data = ns
		} else {
			resp.Code = CREATE_SC_ERROR
			resp.Message = err.Error()
		}
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary 查询列表
// @Tags Keyword
// @Accept  json
// @Produce json
// @Param page_no query int false "page_size"
// @Param page_size query int false "page_size"
// @Param keyword query string false "keyword"
// @Success 200 {object} Response
// @Failure 500 {string} string "Internal Error"
// @Router  /api/v1/keyword/page [get]
func QueryKeywordPage(c *gin.Context) {
	resp := &Response{}

	pageNo, _ := strconv.Atoi(c.Query("page_no"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	keyword := c.Query("keyword")

	dto := &service.KeywordQueryDto{}

	dto.PageNo = pageNo
	dto.PageSize = pageSize
	dto.Keyword = keyword

	list, count := service.QueryKeywordPage(dto)

	data := &PageQueryResponse{
		List:     list,
		Count:    count,
		PageNo:   dto.PageNo,
		PageSize: dto.PageSize,
	}
	resp.Code = SUCCESS
	resp.Message = "success"
	resp.Data = data

	c.JSON(http.StatusOK, resp)
}

// @Summary 查询关键字详情
// @Tags Keyword
// @Accept  json
// @Produce json
// @Param id query string true "id"
// @Success 200 {object} Response
// @Failure 500 {string} string "Internal Error"
// @Router  /api/v1/keyword/detail [get]
func QueryKeywordDetail(c *gin.Context) {
	resp := &Response{}

	id := c.Query("id")
	detail, err := service.QueryKeywordDetail(id)

	if err == nil {
		resp.Code = SUCCESS
		resp.Message = "success"
		resp.Data = detail
	} else {
		resp.Code = CREATE_SC_ERROR
		resp.Message = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary 运行一个脚本
// @Tags Keyword
// @Accept  json
// @Produce json
// @Param create body controller.RunParams true "controller.RunParams"
// @Success 200 {object} Response
// @Failure 500 {string} string "Internal Error"
// @Router  /api/v1/keyword/run [post]
func RunKeywordScript(c *gin.Context) {
	resp := &Response{}
	var dto *RunParams
	err := c.BindJSON(&dto)
	if err != nil {
		resp.Code = BAD_PARAMS
		resp.Message = "params error !!!"
	} else {

		ns, err := service.RunKeywordScript(dto.Keyword, dto.Params)
		if err == nil {
			resp.Code = SUCCESS
			resp.Message = "success"
			resp.Data = ns
		} else {
			resp.Code = CREATE_SC_ERROR
			resp.Message = err.Error()
		}
	}

	c.JSON(http.StatusOK, resp)
}