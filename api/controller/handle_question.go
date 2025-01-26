package controller

import (
	"github.com/gin-gonic/gin"
	"walrus_llm_project/common/response"
	"walrus_llm_project/internal/handle_by_llm"
)

type HandleController struct {
	svc *handle_by_llm.HandleService
}

func NewHandleController() *HandleController {
	return &HandleController{
		svc: handle_by_llm.NewHandleService(),
	}
}
func (hc *HandleController) HandleQuestion(c *gin.Context) {
	var req HandleQuestionRequest
	if e := c.ShouldBind(&req); e != nil {
		response.Fail(c, response.ParamParserError)
		return
	}
	resp, err := hc.svc.HandleByLLM(c, &req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Success(c, resp)
}

type HandleQuestionRequest struct {
	Address  string `json:"address"`
	Question string `json:"question"`
}
type HandleQuestionResponse struct {
	Address string `json:"address"`
	Answer  string `json:"answer"`
}

//=============================================== walrus ========================================================

// 上传正式的健康相关的文字
func (hc *HandleController) UploadText(c *gin.Context) {
	var req UploadTextRequest
	if e := c.ShouldBind(&req); e != nil {
		response.Fail(c, response.ParamParserError)
		return
	}
	err := hc.svc.UploadText(c, &req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Success(c, nil)
}

type UploadTextRequest struct {
	Address string `json:"address"`
	Text    string `json:"text"`
}

// 上传正式的文件列表
func (hc *HandleController) UploadFiles(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file"]
	if len(files) == 0 {
		response.Fail(c, response.ParamParserError)
		return
	}
	address := c.Request.FormValue("address")
	err := hc.svc.UploadFiles(c, address, files)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Success(c, nil)

}
