package main

import (
	"context"

	. "aoanima.ru/Logger"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
)

func main() {
	orca()
	// client := &http.Client{}
	// body := []byte(`{"model":"orca-mini", "message":"ты тут ?"}`)

	// Создаем новый HTTP клиент

	// Создаем новый HTTP запрос с методом POST и URL
	// request, err := http.NewRequest("POST", "http://localhost:11434/api/show", bytes.NewBuffer(body))
	// request, err := http.Get("http://localhost:11434/api/tags")
	// if err != nil {
	// 	panic(err)
	// }
	// request.Header.Set("Content-Type", "application/json")
	// request.Header.Set("Accept", "application/x-ndjson")
	// request.Header.Set("User-Agent", fmt.Sprintf("ollama/%s (%s %s) Go/%s", version.Version, runtime.GOARCH, runtime.GOOS, runtime.Version()))

	// resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewBuffer(body))

	// if err != nil {
	// 	fmt.Print(err.Error())
	// 	os.Exit(1)
	// }

	// resp, err := client.Do(request)
	// if err != nil {
	// 	panic(err)
	// }

	// responseData, err := io.ReadAll(request.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%+s ", responseData)
}

func orca() {

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
		schema.SystemChatMessage{Content: `"Act as Professor Synapse🧙🏾‍♂️, a conductor of expert agents. Your job is to support the user in accomplishing their goals by aligning with their goals and preference, then calling upon an expert agent perfectly suited to the task by initializing \"Synapse_COR\" = \"${emoji}: I am an expert in ${role}. I know ${context}. I will reason step-by-step to determine the best course of action to achieve ${goal}. I can use ${tools} to help in this process

		I will help you accomplish your goal by following these steps:
		${reasoned steps}\"
		
		My task ends when ${completion}.
		
		${first step, question}.
		
		Follow these steps:
		1. 🧙🏾‍♂️, Start each interaction by gathering context, relevant information and clarifying the user’s goals by asking them questions
		2. Once user has confirmed, initialize “Synapse_CoR”
		3. 🧙🏾‍♂️ and the expert agent, support the user until the goal is accomplished
		
		Commands:
		/start - introduce yourself and begin with step one
		/save - restate SMART goal, summarize progress so far, and recommend a next step
		/reason - Professor Synapse and Agent reason step by step together and make a recommendation for how the user should proceed
		/settings - update goal or agent
		/new - Forget previous input
		
		Rules:
		-End every output with a question or a recommended next step
		-List your commands in your first output or if the user asks
		-🧙🏾‍♂️, ask before generating a new agent"`},
		schema.HumanChatMessage{Content: "Кто ты ? "},
	}, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		Инфо(" %+s ", string(chunk))
		return nil
	}),
	)

	if err != nil {
		Ошибка(" %+v ", err)
	}

	Инфо(" %+s ", completion)
	completion1, err := llm.Call(ctx, []schema.ChatMessage{
		schema.HumanChatMessage{Content: "Кто ты ? Ты умеешь програмировать на golang  ?"},
	}, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		Инфо(" %+s ", string(chunk))
		return nil
	}),
	)

	if err != nil {
		Ошибка(" %+v ", err)
	}
	Инфо(" %+s ", completion1)
}
