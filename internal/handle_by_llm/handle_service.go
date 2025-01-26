package handle_by_llm

import (
	"context"
	"mime/multipart"
	"walrus_llm_project/api/controller"
	"walrus_llm_project/common/response"
)

type HandleService struct {
}

func NewHandleService() *HandleService {

	return &HandleService{}
}

func (svc *HandleService) HandleByLLM(ctx context.Context, req *controller.HandleQuestionRequest) (*controller.HandleQuestionResponse, *response.ErrCode) {
	//获取sui上用户的摘要信息
	//返回
	return nil, nil
}

func (svc *HandleService) UploadText(ctx context.Context, req *controller.UploadTextRequest) *response.ErrCode {
	return nil
}

func (svc *HandleService) UploadFiles(ctx context.Context, address string, files []*multipart.FileHeader) *response.ErrCode {
	return nil
}

var HandleByLLMTemplate = `
请你依次执行以下步骤：
① 使用以下上下文来回答最后的问题。如果你不知道答案，就说你不知道，不要试图编造答案。
你应该使答案尽可能详细具体，但不要偏题。如果答案比较长，请酌情进行分段，以提高答案的阅读体验。
如果答案有几点，你应该分点标号回答，让答案清晰具体。
上下文：
{context}
问题: 
{question}
有用的回答:
② 基于提供的上下文，反思回答中有没有不正确或不是基于上下文得到的内容，如果有，回答你不知道
确保你执行了每一个步骤，不要跳过任意一个步骤。
`
var UploadTemplate = `
请你依次执行一下步骤：
1、使用以下上下文 以及 提供的补充条件来生成最新的总结。
2、总结的结果应该只是在原来的上下文的基础，按照 补充条件提供的信息来更新或者补充原来的上下文，以此来生成一份最新的总结。
上下文:
{context}
补充条件:
{more solution}
`
