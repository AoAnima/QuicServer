package main

import (
	"context"

	. "aoanima.ru/Logger"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
)

func main() {

	history := memory.NewChatMessageHistory()

	llm, err := ollama.NewChat(
		ollama.WithLLMOptions(ollama.WithModel("orca-mini")),
	)
	if err != nil {
		Ошибка(" %+v ", err)
	}
	ctx := context.Background()
	completion, err := llm.Call(ctx, []schema.ChatMessage{
		schema.SystemChatMessage{Content: "Give a precise answer to the question based on the context. Don't be verbose."},
		schema.HumanChatMessage{Content: "What would be a good company name a company that makes colorful socks? Give me 3 examples."},
	}, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		Инфо(" %+s ", string(chunk))
		return nil
	}),
	)
	if err != nil {
		Ошибка(" %+v ", err)
	}

	Инфо(" %+s ", completion)

}
