package recognition

import (
	"fmt"
	"sync"

	spch "github.com/tuputech/tupu-go-sdk/recognition/speech/speechsync"
)

type (
	SyncSpeechHdler struct {
		Client *spch.SyncHandler
		CommSpeechEngin
	}

	SySpchClient struct {
		clientMap *sync.Map
	}
)

func (signleton *SySpchClient) GetInterface() *SySpchClient {
	once.Do(func() {
		sySpchClients = &SySpchClient{clientMap: &sync.Map{}}
	})
	return sySpchClients
}

func (signleton *SySpchClient) GetAPIClient(privkeyPath string) *spch.SyncHandler {
	value, ok := signleton.clientMap.Load(privkeyPath)
	if ok {
		return value.(*spch.SyncHandler)
	}
	cli, err := spch.NewSyncHandler(privkeyPath)
	if err != nil {
		panic("create sync handler failed!" + err.Error())
	}
	signleton.clientMap.Store(privkeyPath, cli)
	return cli
}

func (syp *SyncSpeechHdler) Recognition(url string, concurrent *sync.WaitGroup) {
	defer concurrent.Done()
	if syp.Client == nil {
		panic("get sync speech handler failed")
	}
	result, statusCode, err := syp.Client.PerformWithURL(syp.SecretId, []string{url})
	// result, statusCode, err := "test", 200, errors.New("text")
	if err != nil {
		result = fmt.Sprintf("recognition failed: %s", err.Error())
	} else if statusCode != 200 {
		result = fmt.Sprintf("recognition failed statusCode: %d", statusCode)
	}
	// content, _ := json.Marshal(map[string]string{url: result})
	syp.ResultDataCh <- map[string]string{url: result}
}
