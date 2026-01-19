package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ollama/ollama/api"
)

func main() {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	parse_template := `You are a professional {SOURCE_LANG} ({SOURCE_CODE}) to {TARGET_LANG} ({TARGET_CODE}) translator. Your goal is to accurately convey the meaning and nuances of the original {SOURCE_LANG} text while adhering to {TARGET_LANG} grammar, vocabulary, and cultural sensitivities.
Produce only the {TARGET_LANG} translation, without any additional explanations or commentary. Please translate the following {SOURCE_LANG} text into {TARGET_LANG}:


{TEXT}`

	fmt.Printf("parse_template: %v\n", parse_template)

	text := `You are a professional Chinese (zh-Hans) to English (en) translator. Your goal is to accurately convey the meaning and nuances of the original Chinese text while adhering to English grammar, vocabulary, and cultural sensitivities.
Produce only the English translation, without any additional explanations or commentary. Please translate the following Chinese text into English:


我能吞下玻璃而不伤身体。`
	req := &api.GenerateRequest{
		Model:  "translategemma:latest",
		Prompt: text,

		// set streaming to false
		Stream: new(bool),
	}

	ctx := context.Background()
	respFunc := func(resp api.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.
		fmt.Println(resp.Response)
		return nil
	}

	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}

}
