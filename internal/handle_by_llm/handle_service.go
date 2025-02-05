package handle_by_llm

import (
	"context"
	"fmt"
	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"mime/multipart"
	"walrus_llm_project/api/controller"
	"walrus_llm_project/client/walrus_publisher"
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
func (svc *HandleService) GetObjectData(ctx context.Context, objectId string) (map[string]interface{}, *response.ErrCode) {
	var cli = sui.NewSuiClient(constant.BvTestnetEndpoint)

	rsp, err := cli.SuiGetObject(ctx, models.SuiGetObjectRequest{
		ObjectId: "0xeeb964d1e640219c8ddb791cc8548f3242a3392b143ff47484a3753291cad898",
		// only fetch the effects field
		Options: models.SuiObjectDataOptions{
			ShowContent: true,
			ShowDisplay: true,
			ShowType:    true,
			ShowBcs:     true,
			//ShowOwner:               true,
			//ShowPreviousTransaction: true,
			ShowStorageRebate: true,
		},
	})

	if err != nil {
		return nil, response.SystemError.ReplaceMsg(err.Error())
	}
	return rsp.Data.Content.Fields, nil
}
func (svc *HandleService) UploadText(ctx context.Context, req *controller.UploadTextRequest) *response.ErrCode {
	res, err := walrus_publisher.PublishWalrus(ctx, req.Text)
	if err != nil {
		return response.SystemError.ReplaceMsg(err.Error())
	}
	blobId := res.NewlyCreated.BlobObject.BlobID
	address := req.Address
	println(blobId + "::::" + address)
	var cli = sui.NewSuiClient(constant.BvTestnetEndpoint)

	gasObj := "0x58c103930dc52c0ab86319d99218e301596fda6fd80c4efafd7f4c9df1d0b6d0"
	signerAddress := ""
	rsp, err := cli.MoveCall(ctx, models.MoveCallRequest{
		Signer:          signerAddress,
		PackageObjectId: "0x7d584c9a27ca4a546e8203b005b0e9ae746c9bec6c8c3c0bc84611bcf4ceab5f",
		Module:          "auction",
		Function:        "start_an_auction",
		TypeArguments:   []interface{}{},
		Arguments: []interface{}{
			"0x342e959f8d9d1fa9327a05fd54fefd929bbedad47190bdbb58743d8ba3bd3420",
			"0x3fd0fdedb84cf1f59386b6251ba6dd2cb495094da26e0a5a38239acd9d437f96",
			"0xb3de4235cb04167b473de806d00ba351e5860500253cf8e62d711e578e1d92ae",
			"BlockVision",
			"0xc699c6014da947778fe5f740b2e9caf905ca31fb4c81e346f467ae126e3c03f1",
		},
		Gas:       &gasObj,
		GasBudget: "100000000",
	})
	if err != nil {
		fmt.Println(err.Error())
		return response.SystemError.ReplaceMsg(err.Error())
	}
	fmt.Println(rsp)
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
