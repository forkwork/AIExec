package tester

import (
	"time"

	"aiexec-plugin/internal/core/aiexec_invocation"
	jsonschema "aiexec-plugin/internal/utils/json_schema"
	"aiexec-plugin/internal/utils/parser"
	"aiexec-plugin/internal/utils/routine"
	"aiexec-plugin/internal/utils/stream"
	"aiexec-plugin/pkg/entities/model_entities"
	"aiexec-plugin/pkg/entities/tool_entities"
)

type MockedAiexecInvocation struct{}

func NewMockedAiexecInvocation() aiexec_invocation.BackwardsInvocation {
	return &MockedAiexecInvocation{}
}

func (m *MockedAiexecInvocation) InvokeLLM(payload *aiexec_invocation.InvokeLLMRequest) (*stream.Stream[model_entities.LLMResultChunk], error) {
	stream := stream.NewStream[model_entities.LLMResultChunk](5)
	routine.Submit(nil, func() {
		if len(payload.Tools) > 0 {
			tool := payload.Tools[0]
			// generate a valid arguments json string
			arguments, err := jsonschema.GenerateValidateJson(tool.Parameters)
			if err != nil {
				stream.WriteError(err)
				return
			}

			stream.WriteBlocking(model_entities.LLMResultChunk{
				Model:             model_entities.LLMModel(payload.Model),
				SystemFingerprint: "test",
				Delta: model_entities.LLMResultChunkDelta{
					Index: parser.ToPtr(0),
					Message: model_entities.PromptMessage{
						Role: model_entities.PROMPT_MESSAGE_ROLE_ASSISTANT,
						ToolCalls: []model_entities.PromptMessageToolCall{
							{
								ID:   "oshiawaseni", // お幸せに
								Type: "function",
								Function: struct {
									Name      string "json:\"name\""
									Arguments string "json:\"arguments\""
								}{
									Name:      tool.Name,
									Arguments: parser.MarshalJson(arguments),
								},
							},
						},
					},
				},
			})
		}

		stream.Write(model_entities.LLMResultChunk{
			Model:             model_entities.LLMModel(payload.Model),
			SystemFingerprint: "test",
			Delta: model_entities.LLMResultChunkDelta{
				Index: &[]int{1}[0],
				Message: model_entities.PromptMessage{
					Role:      model_entities.PROMPT_MESSAGE_ROLE_ASSISTANT,
					Content:   "hello",
					Name:      "test",
					ToolCalls: []model_entities.PromptMessageToolCall{},
				},
			},
		})
		time.Sleep(100 * time.Millisecond)
		stream.Write(model_entities.LLMResultChunk{
			Model:             model_entities.LLMModel(payload.Model),
			SystemFingerprint: "test",
			Delta: model_entities.LLMResultChunkDelta{
				Index: &[]int{1}[0],
				Message: model_entities.PromptMessage{
					Role:      model_entities.PROMPT_MESSAGE_ROLE_ASSISTANT,
					Content:   " world",
					Name:      "test",
					ToolCalls: []model_entities.PromptMessageToolCall{},
				},
			},
		})
		time.Sleep(100 * time.Millisecond)
		stream.Write(model_entities.LLMResultChunk{
			Model:             model_entities.LLMModel(payload.Model),
			SystemFingerprint: "test",
			Delta: model_entities.LLMResultChunkDelta{
				Index: &[]int{2}[0],
				Message: model_entities.PromptMessage{
					Role:      model_entities.PROMPT_MESSAGE_ROLE_ASSISTANT,
					Content:   " world",
					Name:      "test",
					ToolCalls: []model_entities.PromptMessageToolCall{},
				},
			},
		})
		time.Sleep(100 * time.Millisecond)
		stream.Write(model_entities.LLMResultChunk{
			Model:             model_entities.LLMModel(payload.Model),
			SystemFingerprint: "test",
			Delta: model_entities.LLMResultChunkDelta{
				Index: &[]int{3}[0],
				Message: model_entities.PromptMessage{
					Role:      model_entities.PROMPT_MESSAGE_ROLE_ASSISTANT,
					Content:   " !",
					Name:      "test",
					ToolCalls: []model_entities.PromptMessageToolCall{},
				},
			},
		})
		time.Sleep(100 * time.Millisecond)
		stream.Write(model_entities.LLMResultChunk{
			Model:             model_entities.LLMModel(payload.Model),
			SystemFingerprint: "test",
			Delta: model_entities.LLMResultChunkDelta{
				Index: &[]int{3}[0],
				Message: model_entities.PromptMessage{
					Role:      model_entities.PROMPT_MESSAGE_ROLE_ASSISTANT,
					Content:   " !",
					Name:      "test",
					ToolCalls: []model_entities.PromptMessageToolCall{},
				},
				Usage: &model_entities.LLMUsage{
					PromptTokens:     &[]int{100}[0],
					CompletionTokens: &[]int{100}[0],
					TotalTokens:      &[]int{200}[0],
					Latency:          &[]float64{0.4}[0],
					Currency:         &[]string{"USD"}[0],
				},
			},
		})
		stream.Close()
	})
	return stream, nil
}

func (m *MockedAiexecInvocation) InvokeTextEmbedding(payload *aiexec_invocation.InvokeTextEmbeddingRequest) (*model_entities.TextEmbeddingResult, error) {
	result := model_entities.TextEmbeddingResult{
		Model: payload.Model,
		Usage: model_entities.EmbeddingUsage{
			Tokens:      &[]int{100}[0],
			TotalTokens: &[]int{100}[0],
			Latency:     &[]float64{0.1}[0],
			Currency:    &[]string{"USD"}[0],
		},
	}

	for range payload.Texts {
		result.Embeddings = append(result.Embeddings, []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0})
	}

	return &result, nil
}

func (m *MockedAiexecInvocation) InvokeRerank(payload *aiexec_invocation.InvokeRerankRequest) (*model_entities.RerankResult, error) {
	result := model_entities.RerankResult{
		Model: payload.Model,
	}
	for i, doc := range payload.Docs {
		result.Docs = append(result.Docs, model_entities.RerankDocument{
			Index: &[]int{i}[0],
			Score: &[]float64{0.1}[0],
			Text:  &doc,
		})
	}
	return &result, nil
}

func (m *MockedAiexecInvocation) InvokeTTS(payload *aiexec_invocation.InvokeTTSRequest) (*stream.Stream[model_entities.TTSResult], error) {
	stream := stream.NewStream[model_entities.TTSResult](5)
	routine.Submit(nil, func() {
		for i := 0; i < 10; i++ {
			stream.Write(model_entities.TTSResult{
				Result: "a1b2c3d4",
			})
			time.Sleep(100 * time.Millisecond)
		}
		stream.Close()
	})
	return stream, nil
}

func (m *MockedAiexecInvocation) InvokeSpeech2Text(payload *aiexec_invocation.InvokeSpeech2TextRequest) (*model_entities.Speech2TextResult, error) {
	result := model_entities.Speech2TextResult{
		Result: "hello world",
	}
	return &result, nil
}

func (m *MockedAiexecInvocation) InvokeModeration(payload *aiexec_invocation.InvokeModerationRequest) (*model_entities.ModerationResult, error) {
	result := model_entities.ModerationResult{
		Result: true,
	}
	return &result, nil
}

func (m *MockedAiexecInvocation) InvokeTool(payload *aiexec_invocation.InvokeToolRequest) (*stream.Stream[tool_entities.ToolResponseChunk], error) {
	stream := stream.NewStream[tool_entities.ToolResponseChunk](5)
	routine.Submit(nil, func() {
		for i := 0; i < 10; i++ {
			stream.Write(tool_entities.ToolResponseChunk{
				Type: tool_entities.ToolResponseChunkTypeText,
				Message: map[string]any{
					"text": "hello world",
				},
			})
			time.Sleep(100 * time.Millisecond)
		}
		stream.Close()
	})

	return stream, nil
}

func (m *MockedAiexecInvocation) InvokeApp(payload *aiexec_invocation.InvokeAppRequest) (*stream.Stream[map[string]any], error) {
	stream := stream.NewStream[map[string]any](5)
	routine.Submit(nil, func() {
		stream.Write(map[string]any{
			"event":           "agent_message",
			"message_id":      "5ad4cb98-f0c7-4085-b384-88c403be6290",
			"conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2",
			"answer":          "なんで",
			"created_at":      time.Now().Unix(),
		})
		time.Sleep(100 * time.Millisecond)
		stream.Write(map[string]any{
			"event":           "agent_message",
			"message_id":      "5ad4cb98-f0c7-4085-b384-88c403be6290",
			"conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2",
			"answer":          "春日影",
			"created_at":      time.Now().Unix(),
		})
		time.Sleep(100 * time.Millisecond)
		stream.Write(map[string]any{
			"event":           "agent_message",
			"message_id":      "5ad4cb98-f0c7-4085-b384-88c403be6290",
			"conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2",
			"answer":          "やったの",
			"created_at":      time.Now().Unix(),
		})
		time.Sleep(100 * time.Millisecond)
		stream.Write(map[string]any{
			"event":           "message_end",
			"id":              "5e52ce04-874b-4d27-9045-b3bc80def685",
			"conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2",
			"created_at":      time.Now().Unix(),
			"metadata": map[string]any{
				"retriever_resources": []map[string]any{
					{
						"position":      1,
						"dataset_id":    "101b4c97-fc2e-463c-90b1-5261a4cdcafb",
						"dataset_name":  "あなた",
						"document_id":   "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00",
						"document_name": "ご自分のことばかりですのね",
						"score":         0.98457545,
						"content":       "CRYCHICは壊れてしまいましたわ",
					},
				},
				"usage": map[string]any{
					"prompt_tokens":         1033,
					"prompt_unit_price":     "0.001",
					"prompt_price_unit":     "0.001",
					"prompt_price":          "0.0010330",
					"completion_tokens":     135,
					"completion_unit_price": "0.002",
					"completion_price_unit": "0.001",
					"completion_price":      "0.0002700",
					"total_tokens":          1168,
					"total_price":           "0.0013030",
					"currency":              "USD",
					"latency":               1.381760165997548,
				},
			},
		})
		time.Sleep(100 * time.Millisecond)
		stream.Write(map[string]any{
			"event":           "message_file",
			"id":              "5e52ce04-874b-4d27-9045-b3bc80def685",
			"conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2",
			"belongs_to":      "assistant",
			"url":             "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			"created_at":      time.Now().Unix(),
		})
		time.Sleep(100 * time.Millisecond)
		stream.Close()
	})
	return stream, nil
}

func (m *MockedAiexecInvocation) InvokeEncrypt(payload *aiexec_invocation.InvokeEncryptRequest) (map[string]any, error) {
	return payload.Data, nil
}

func (m *MockedAiexecInvocation) InvokeParameterExtractor(payload *aiexec_invocation.InvokeParameterExtractorRequest) (*aiexec_invocation.InvokeNodeResponse, error) {
	resp := &aiexec_invocation.InvokeNodeResponse{
		ProcessData: map[string]any{},
		Outputs:     map[string]any{},
		Inputs: map[string]any{
			"query": payload.Query,
		},
	}

	for _, parameter := range payload.Parameters {
		typ := parameter.Type
		if typ == "string" {
			resp.Outputs[parameter.Name] = "Never gonna give you up ~"
		} else if typ == "number" {
			resp.Outputs[parameter.Name] = 1234567890
		} else if typ == "bool" {
			resp.Outputs[parameter.Name] = true
		} else if typ == "select" {
			options := parameter.Options
			if len(options) == 0 {
				resp.Outputs[parameter.Name] = "Never gonna let you down ~"
			} else {
				resp.Outputs[parameter.Name] = options[0]
			}
		} else if typ == "array[string]" {
			resp.Outputs[parameter.Name] = []string{
				"Never gonna run around and desert you ~",
				"Never gonna make you cry ~",
				"Never gonna say goodbye ~",
				"Never gonna tell a lie and hurt you ~",
			}
		} else if typ == "array[number]" {
			resp.Outputs[parameter.Name] = []int{114, 514, 1919, 810}
		} else if typ == "array[bool]" {
			resp.Outputs[parameter.Name] = []bool{true, false, true, false, true, false, true, false, true, false}
		} else if typ == "array[object]" {
			resp.Outputs[parameter.Name] = []map[string]any{
				{
					"name": "お願い",
					"age":  55555,
				},
				{
					"name": "何でもするがら",
					"age":  99999,
				},
			}
		}
	}

	return resp, nil
}

func (m *MockedAiexecInvocation) InvokeQuestionClassifier(payload *aiexec_invocation.InvokeQuestionClassifierRequest) (*aiexec_invocation.InvokeNodeResponse, error) {
	return &aiexec_invocation.InvokeNodeResponse{
		ProcessData: map[string]any{},
		Outputs: map[string]any{
			"class_name": payload.Classes[0].Name,
		},
		Inputs: map[string]any{},
	}, nil
}

func (m *MockedAiexecInvocation) InvokeSummary(payload *aiexec_invocation.InvokeSummaryRequest) (*aiexec_invocation.InvokeSummaryResponse, error) {
	return &aiexec_invocation.InvokeSummaryResponse{
		Summary: payload.Text,
	}, nil
}

func (m *MockedAiexecInvocation) UploadFile(payload *aiexec_invocation.UploadFileRequest) (*aiexec_invocation.UploadFileResponse, error) {
	return &aiexec_invocation.UploadFileResponse{
		URL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
	}, nil
}

func (m *MockedAiexecInvocation) FetchApp(payload *aiexec_invocation.FetchAppRequest) (map[string]any, error) {
	return map[string]any{
		"name": "test",
	}, nil
}
