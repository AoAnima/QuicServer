// package main

// import (
// 	"bytes"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// )

// func main() {

// 	client := &http.Client{}
// 	body := []byte(`{"model":"orca-mini", "message":"ты тут ?"}`)

// 	// Создаем новый HTTP клиент

// 	// Создаем новый HTTP запрос с методом POST и URL
// 	request, err := http.NewRequest("POST", "http://localhost:11434/api/generate", bytes.NewBuffer(body))
// 	if err != nil {
// 		panic(err)
// 	}
// 	request.Header.Set("Content-Type", "application/json")
// 	request.Header.Set("Accept", "application/x-ndjson")
// 	// request.Header.Set("User-Agent", fmt.Sprintf("ollama/%s (%s %s) Go/%s", version.Version, runtime.GOARCH, runtime.GOOS, runtime.Version()))

// 	// resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewBuffer(body))

// 	// if err != nil {
// 	// 	fmt.Print(err.Error())
// 	// 	os.Exit(1)
// 	// }

// 	resp, err := client.Do(request)
// 	if err != nil {
// 		panic(err)
// 	}

// 	responseData, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("%+s   resp %+s", responseData, resp)
// }

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
