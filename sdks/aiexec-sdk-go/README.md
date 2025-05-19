# Aiexec Go SDK
This is the Go SDK for the Aiexec API, which allows you to easily integrate Aiexec into your Go applications.

## Install
```bash
go get github.com/khulnasoft/aiexec
```

## Usage
After installing the SDK, you can use it in your project like this:

```go
package main

import (
	"context"
	"log"
	"strings"

	"github.com/khulnasoft/aiexec"
)

func main() {
	var (
		ctx = context.Background()
		c = aiexec.NewClient("your-aiexec-server-host", "your-api-key-here")

		req = &aiexec.ChatMessageRequest{
			Query: "your-question",
			User: "your-user",
		}

		ch chan aiexec.ChatMessageStreamChannelResponse
		err error
	)

	if ch, err = c.Api().ChatMessagesStream(ctx, req); err != nil {
		return
	}

	var strBuilder strings.Builder

	for {
		select {
		case <-ctx.Done():
			return
		case streamData, isOpen := <-ch:
			if err = streamData.Err; err != nil {
				log.Println(err.Error())
				return
			}
			if !isOpen {
				log.Println(strBuilder.String())
				return
			}

			strBuilder.WriteString(streamData.Answer)
		}
	}
}

```

## License
This SDK is released under the MIT License.