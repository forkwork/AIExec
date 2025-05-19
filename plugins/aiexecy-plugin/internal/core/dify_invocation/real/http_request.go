package real

import (
	"fmt"
	"reflect"

	"aiexec-plugin/internal/core/aiexec_invocation"
	"aiexec-plugin/internal/utils/http_requests"
	"aiexec-plugin/internal/utils/routine"
	"aiexec-plugin/internal/utils/stream"
	"aiexec-plugin/pkg/entities/model_entities"
	"aiexec-plugin/pkg/entities/tool_entities"
	"aiexec-plugin/pkg/validators"
)

// Send a request to aiexec inner api and validate the response
func Request[T any](i *RealBackwardsInvocation, method string, path string, options ...http_requests.HttpOptions) (*T, error) {
	options = append(options,
		http_requests.HttpHeader(map[string]string{
			"X-Inner-Api-Key": i.aiexecInnerApiKey,
		}),
		http_requests.HttpWriteTimeout(i.writeTimeout),
		http_requests.HttpReadTimeout(i.readTimeout),
	)

	req, err := http_requests.RequestAndParse[BaseBackwardsInvocationResponse[T]](i.client, i.aiexecPath(path), method, options...)
	if err != nil {
		return nil, err
	}

	if req.Error != "" {
		return nil, fmt.Errorf("request failed: %s", req.Error)
	}

	if req.Data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	// check if req.Data is a map[string]any
	if reflect.TypeOf(*req.Data).Kind() == reflect.Map {
		return req.Data, nil
	}

	if err := validators.GlobalEntitiesValidator.Struct(req.Data); err != nil {
		return nil, fmt.Errorf("validate request failed: %s", err.Error())
	}

	return req.Data, nil
}

func StreamResponse[T any](i *RealBackwardsInvocation, method string, path string, options ...http_requests.HttpOptions) (
	*stream.Stream[T], error,
) {
	options = append(
		options, http_requests.HttpHeader(map[string]string{
			"X-Inner-Api-Key": i.aiexecInnerApiKey,
		}),
		http_requests.HttpWriteTimeout(i.writeTimeout),
		http_requests.HttpReadTimeout(i.readTimeout),
	)

	response, err := http_requests.RequestAndParseStream[BaseBackwardsInvocationResponse[T]](i.client, i.aiexecPath(path), method, options...)
	if err != nil {
		return nil, err
	}

	newResponse := stream.NewStream[T](1024)
	newResponse.OnClose(func() {
		response.Close()
	})
	routine.Submit(map[string]string{
		"module":   "aiexec_invocation",
		"function": "StreamResponse",
	}, func() {
		defer newResponse.Close()
		for response.Next() {
			t, err := response.Read()
			if err != nil {
				newResponse.WriteError(err)
				break
			}

			if t.Error != "" {
				newResponse.WriteError(fmt.Errorf("request failed: %s", t.Error))
				break
			}

			if t.Data == nil {
				newResponse.WriteError(fmt.Errorf("data is nil"))
				break
			}

			// check if t.Data is a map[string]any, skip validation if it is
			if reflect.TypeOf(*t.Data).Kind() != reflect.Map {
				if err := validators.GlobalEntitiesValidator.Struct(t.Data); err != nil {
					newResponse.WriteError(fmt.Errorf("validate request failed: %s", err.Error()))
					break
				}
			}

			newResponse.Write(*t.Data)
		}
	})

	return newResponse, nil
}

func (i *RealBackwardsInvocation) InvokeLLM(payload *aiexec_invocation.InvokeLLMRequest) (*stream.Stream[model_entities.LLMResultChunk], error) {
	return StreamResponse[model_entities.LLMResultChunk](i, "POST", "invoke/llm", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeTextEmbedding(payload *aiexec_invocation.InvokeTextEmbeddingRequest) (*model_entities.TextEmbeddingResult, error) {
	return Request[model_entities.TextEmbeddingResult](i, "POST", "invoke/text-embedding", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeRerank(payload *aiexec_invocation.InvokeRerankRequest) (*model_entities.RerankResult, error) {
	return Request[model_entities.RerankResult](i, "POST", "invoke/rerank", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeTTS(payload *aiexec_invocation.InvokeTTSRequest) (*stream.Stream[model_entities.TTSResult], error) {
	return StreamResponse[model_entities.TTSResult](i, "POST", "invoke/tts", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeSpeech2Text(payload *aiexec_invocation.InvokeSpeech2TextRequest) (*model_entities.Speech2TextResult, error) {
	return Request[model_entities.Speech2TextResult](i, "POST", "invoke/speech2text", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeModeration(payload *aiexec_invocation.InvokeModerationRequest) (*model_entities.ModerationResult, error) {
	return Request[model_entities.ModerationResult](i, "POST", "invoke/moderation", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeTool(payload *aiexec_invocation.InvokeToolRequest) (*stream.Stream[tool_entities.ToolResponseChunk], error) {
	return StreamResponse[tool_entities.ToolResponseChunk](i, "POST", "invoke/tool", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeApp(payload *aiexec_invocation.InvokeAppRequest) (*stream.Stream[map[string]any], error) {
	return StreamResponse[map[string]any](i, "POST", "invoke/app", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeParameterExtractor(payload *aiexec_invocation.InvokeParameterExtractorRequest) (*aiexec_invocation.InvokeNodeResponse, error) {
	return Request[aiexec_invocation.InvokeNodeResponse](i, "POST", "invoke/parameter-extractor", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeQuestionClassifier(payload *aiexec_invocation.InvokeQuestionClassifierRequest) (*aiexec_invocation.InvokeNodeResponse, error) {
	return Request[aiexec_invocation.InvokeNodeResponse](i, "POST", "invoke/question-classifier", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) InvokeEncrypt(payload *aiexec_invocation.InvokeEncryptRequest) (map[string]any, error) {
	if !payload.EncryptRequired(payload.Data) {
		return payload.Data, nil
	}

	type resp struct {
		Data map[string]any `json:"data,omitempty"`
	}

	data, err := Request[resp](i, "POST", "invoke/encrypt", http_requests.HttpPayloadJson(payload))
	if err != nil {
		return nil, err
	}

	return data.Data, nil
}

func (i *RealBackwardsInvocation) InvokeSummary(payload *aiexec_invocation.InvokeSummaryRequest) (*aiexec_invocation.InvokeSummaryResponse, error) {
	return Request[aiexec_invocation.InvokeSummaryResponse](i, "POST", "invoke/summary", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) UploadFile(payload *aiexec_invocation.UploadFileRequest) (*aiexec_invocation.UploadFileResponse, error) {
	return Request[aiexec_invocation.UploadFileResponse](i, "POST", "upload/file/request", http_requests.HttpPayloadJson(payload))
}

func (i *RealBackwardsInvocation) FetchApp(payload *aiexec_invocation.FetchAppRequest) (map[string]any, error) {
	type resp struct {
		Data map[string]any `json:"data,omitempty"`
	}

	data, err := Request[resp](i, "POST", "fetch/app/info", http_requests.HttpPayloadJson(payload))
	if err != nil {
		return nil, err
	}

	if data.Data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	return data.Data, nil
}
