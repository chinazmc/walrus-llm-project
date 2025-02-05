package walrus_publisher

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func PublishWalrus(ctx context.Context, reqStr string) (*PublishWalrusResp, error) {
	publishPoint := "http://localhost:31415"
	url := "/v1/blobs"
	method := "PUT"

	payload := strings.NewReader(reqStr)

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, publishPoint+url, payload)

	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var resp PublishWalrusResp
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, err
}

type PublishWalrusResp struct {
	NewlyCreated NewlyCreated `json:"newlyCreated"`
}
type Storage struct {
	ID          string `json:"id"`
	StartEpoch  int    `json:"startEpoch"`
	EndEpoch    int    `json:"endEpoch"`
	StorageSize int    `json:"storageSize"`
}
type BlobObject struct {
	ID              string  `json:"id"`
	RegisteredEpoch int     `json:"registeredEpoch"`
	BlobID          string  `json:"blobId"`
	Size            int     `json:"size"`
	EncodingType    string  `json:"encodingType"`
	CertifiedEpoch  int     `json:"certifiedEpoch"`
	Storage         Storage `json:"storage"`
	Deletable       bool    `json:"deletable"`
}
type RegisterFromScratch struct {
	EncodedLength int `json:"encodedLength"`
	EpochsAhead   int `json:"epochsAhead"`
}
type ResourceOperation struct {
	RegisterFromScratch RegisterFromScratch `json:"registerFromScratch"`
}
type NewlyCreated struct {
	BlobObject        BlobObject        `json:"blobObject"`
	ResourceOperation ResourceOperation `json:"resourceOperation"`
	Cost              int               `json:"cost"`
}
