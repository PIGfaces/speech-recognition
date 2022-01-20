package recognition

import (
	"fmt"
	"sync"

	spch "github.com/tuputech/tupu-go-sdk/recognition/speech/speechasync"
)

type (
	AsyncSpeechHdler struct {
		Client *spch.AsyncHandler
		CommSpeechEngin
	}

	AsySpchClient struct {
		clientMap *sync.Map
	}
)

var (
	asySpchClients *AsySpchClient
)

func (signleton *AsySpchClient) GetInterface() *AsySpchClient {
	once.Do(func() {
		asySpchClients = &AsySpchClient{clientMap: &sync.Map{}}
	})
	return asySpchClients
}

func (signleton *AsySpchClient) GetAPIClient(privkeyPath string) *spch.AsyncHandler {
	value, ok := signleton.clientMap.Load(privkeyPath)
	if ok {
		return value.(*spch.AsyncHandler)
	}
	cli, err := spch.NewSpeechHandler(privkeyPath)
	cli.SetServerURL("http://api.speech.tuputech.com/v3/recognition/speech/recording/async/")
	if err != nil {
		panic("create sync handler failed!" + err.Error())
	}
	signleton.clientMap.Store(privkeyPath, cli)
	return cli
}

func (syp *AsyncSpeechHdler) Recognition(url string, concurrent *sync.WaitGroup) {
	defer concurrent.Done()
	if syp.Client == nil {
		panic("get sync speech handler failed")
	}
	result, statusCode, err := syp.Client.Perform(syp.SecretId, url)
	// result, statusCode, err := "test", 200, errors.New("text")
	if err != nil {
		result = fmt.Sprintf("recognition failed: %s", err.Error())
	} else if statusCode != 200 {
		result = fmt.Sprintf("recognition failed statusCode: %d", statusCode)
	}
	// content, _ := json.Marshal(map[string]string{url: result})
	syp.ResultDataCh <- map[string]string{url: result}
}
