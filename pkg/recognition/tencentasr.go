package recognition

import (
	"fmt"
	"strings"
	"sync"

	asr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/asr/v20190614"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

type (
	TencentAsrHdler struct {
		Client *asr.Client
		CommSpeechEngin
	}

	TencentAsrClient struct {
		SySpchClient
	}
)

var tcentAsrClients *TencentAsrClient

func (signleton *TencentAsrClient) GetInterface() *TencentAsrClient {
	once.Do(func() {
		tcentAsrClients = &TencentAsrClient{SySpchClient: SySpchClient{clientMap: &sync.Map{}}}
	})
	return tcentAsrClients
}

func (signleton *TencentAsrClient) GetAPIClient(secretKey string) *asr.Client {
	value, ok := signleton.clientMap.Load(secretKey)
	if ok {
		return value.(*asr.Client)
	}
	credential := common.NewCredential(
		// secretId
		"AKIDVpsfBQO6XicMq2xbz6wVoqb4HDk0zjDQ",
		// secretKey
		secretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "asr.tencentcloudapi.com"
	client, err := asr.NewClient(credential, "", cpf)

	if err != nil {
		panic("create tencent asr handler failed!" + err.Error())
	}
	signleton.clientMap.Store(secretKey, client)
	return client
}

func (tcentAsr *TencentAsrHdler) Recognition(url string, concurrent *sync.WaitGroup) {
	defer concurrent.Done()
	if tcentAsr.Client == nil {
		panic("get sync speech handler failed")
	}

	request := createAsrRequest(url)
	response, err := tcentAsr.Client.SentenceRecognition(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}
	// content, _ := json.Marshal(map[string]string{url: result})
	tcentAsr.ResultDataCh <- map[string]string{url: response.ToJsonString()}
}

func createAsrRequest(url string) *asr.SentenceRecognitionRequest {

	fileType := getUrlFileType(url)

	request := asr.NewSentenceRecognitionRequest()

	request.ProjectId = common.Uint64Ptr(0)
	request.SubServiceType = common.Uint64Ptr(2)
	request.EngSerViceType = common.StringPtr("16k_zh")
	request.SourceType = common.Uint64Ptr(0)

	// request.Url = common.StringPtr("https://static.tuputech.com/api/image/original/cloud-api/storage-0831/2022-01-18/18-7/6049b867400a3700813f4cef/16425005777310.29285377740674967.wav")
	request.Url = common.StringPtr(url)
	request.VoiceFormat = common.StringPtr(fileType)
	request.UsrAudioKey = common.StringPtr("test")
	return request
}

func getUrlFileType(url string) string {
	fileType := url[strings.LastIndexByte(url, '.')+1:]
	if len(fileType) == 0 {
		return "wav"
	}
	return fileType
}
