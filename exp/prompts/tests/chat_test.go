package prompts_test

import (
	"reflect"
	"testing"

	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
)

func TestChatTemplate(t *testing.T) {
	systemPrompt, err := prompts.NewPromptTemplate("Here's some context: {context}", []string{"context"})
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	userPrompt, err := prompts.NewPromptTemplate("Hello AI. Give me a long response. {question}", []string{"question"})
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	aiPrompt, err := prompts.NewPromptTemplate("Very good question. My answer to {question} is {answer}", []string{"answer", "question"})
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	messages := []prompts.Message{
		prompts.NewSystemMessage(systemPrompt),
		prompts.NewHumanMessage(userPrompt),
		prompts.NewAiMessage(aiPrompt),
	}

	_, err = prompts.NewChatTemplate(messages, []string{"answer", "context"})
	if err == nil {
		t.Errorf("Expected error creating chat template with too few variables")
	}

	_, err = prompts.NewChatTemplate(messages, []string{"answer", "context", "question", "foo"})
	if err == nil {
		t.Errorf("Expected error creating chat template with too many variables")
	}

	chatTemplate, err := prompts.NewChatTemplate(messages, []string{"answer", "context", "question"})
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	chatMessages, err := chatTemplate.FormatPromptValue(map[string]any{"context": "foo", "question": "bar", "answer": "foobar"})
	expectedChatMessages := []schema.ChatMessage{
		schema.SystemChatMessage{Text: "Here's some context: foo"},
		schema.HumanChatMessage{Text: "Hello AI. Give me a long response. bar"},
		schema.AiChatMessage{Text: "Very good question. My answer to bar is foobar"},
	}
	expectedString := `[{"text":"Here's some context: foo"},{"text":"Hello AI. Give me a long response. bar"},{"text":"Very good question. My answer to bar is foobar"}]`

	if !reflect.DeepEqual(chatMessages.ToChatMessages(), expectedChatMessages) {
		t.Errorf("Chat template format prompt value chat messages not equal to expected. Got: %v. Expect: %v", chatMessages.ToChatMessages(), expectedChatMessages)
	}

	if !(chatMessages.String() == expectedString) {
		t.Errorf("Chat template format prompt value string not equal to expected.\n Got:\n %v\n Expect:\n %v", chatMessages.String(), expectedString)
	}
}
