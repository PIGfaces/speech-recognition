package recognition

import (
	"sync"
)

type (
	TupuSpeechHdlerIntf interface {
		Recognition(url string, concurrent *sync.WaitGroup)
	}

	TupuSpeechClientIntf interface {
		GetInterface()
		GetAPIClient(privkeyPath string)
	}

	CommSpeechEngin struct {
		SecretId     string
		ResultDataCh chan<- map[string]string
	}
)

var (
	once          sync.Once
	sySpchClients *SySpchClient
)
