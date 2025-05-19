package real

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"aiexec-plugin/internal/core/aiexec_invocation"
)

type NewAiexecInvocationDaemonPayload struct {
	BaseUrl      string
	CallingKey   string
	WriteTimeout int64
	ReadTimeout  int64
}

func NewAiexecInvocationDaemon(payload NewAiexecInvocationDaemonPayload) (aiexec_invocation.BackwardsInvocation, error) {
	var err error
	invocation := &RealBackwardsInvocation{}
	baseurl, err := url.Parse(payload.BaseUrl)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 120 * time.Second,
			}).Dial,
			IdleConnTimeout: 120 * time.Second,
		},
	}

	invocation.aiexecInnerApiBaseurl = baseurl
	invocation.client = client
	invocation.aiexecInnerApiKey = payload.CallingKey
	invocation.writeTimeout = payload.WriteTimeout
	invocation.readTimeout = payload.ReadTimeout

	return invocation, nil
}
