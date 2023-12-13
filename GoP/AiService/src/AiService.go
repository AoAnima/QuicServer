package main

import (
	"context"

	. "aoanima.ru/Logger"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
)

func main() {

	llm, err := ollama.NewChat(
		ollama.WithLLMOptions(
			ollama.WithModel("orca-mini"),
			// ollama.WithSystemPrompt("Давай точный ответ, основываясь на контексте, будь краток и выразителен."),
		),
	)
	if err != nil {
		Ошибка(" %+v ", err)
	}

	// embedder, err := embeddings.NewEmbedder(llm)
	// if err != nil {
	// 	Ошибка(" %+v ", err)
	// }
	// docs := []string{"doc 1", "another doc"}
	// _, err = embedder.EmbedDocuments(context.Background(), docs)
	// if err != nil {
	// 	Ошибка(" %+v ", err)
	// }
	// v, err := embedder.EmbedQuery(context.Background(), "doc 1")
	// if err != nil {
	// 	Ошибка("  %+v \n", err)
	// }
	// Инфо("  %+s \n", v)
	// sm := schema.SystemChatMessage{Content: "Give a precise answer to the question based on the context. Don't be verbose."}
	// Инфо(" sm.GetContent() %+v \n", sm.GetContent())

	ctx := context.Background()
	completion, err := llm.Call(ctx, []schema.ChatMessage{
		// schema.SystemChatMessage{Content: "Give a precise answer to the question based on the context. Don't be verbose."},
		schema.HumanChatMessage{Content: "Write 12 word start on leter S"},
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
