package agent

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/assistant"
	"nano_cc/internal/config"
	"nano_cc/internal/logger"
	"nano_cc/internal/sessions"
)

func RunREPL() {
	config.Init()

	logger.Info("nano_cc (Go) — starting up...")

	if err := logger.InitFileLogger(); err != nil {
		logger.Warn(fmt.Sprintf("failed to init file logger: %v", err))
	}

	if err := sessions.Skills.LoadAll(); err != nil {
		logger.Warn(fmt.Sprintf("failed to load skills: %v", err))
	}

	cwd, _ := os.Getwd()

	skillDescriptions := sessions.Skills.GetDescriptions()
	skillsSection := ""
	if skillDescriptions != "" {
		skillsSection = fmt.Sprintf("\nSkills available:\n%s\n", skillDescriptions)
	}

	// Load persistent memory
	projectMemory := sessions.LoadProjectMemory()
	userMemory := sessions.LoadUserMemory()
	memorySection := ""
	if ms := sessions.BuildPersistentMemorySection(projectMemory, userMemory); ms != "" {
		memorySection = "\n\n" + ms
	}

	systemPrompt := assistant.BuildSystemPrompt(assistant.PromptOptions{
		Surface:           fmt.Sprintf("the legacy terminal interface at %s", cwd),
		SkillDescriptions: strings.TrimSpace(skillsSection),
		MemorySection:     strings.TrimSpace(memorySection),
	})

	scanner := bufio.NewScanner(os.Stdin)

	messages := []openai.ChatCompletionMessage{
		{Role: "system", Content: systemPrompt},
	}

	for {
		fmt.Print("\nYou: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		switch strings.ToLower(input) {
		case "exit", "quit":
			logger.Info("Goodbye!")
			return
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    "user",
			Content: input,
		})

		reply, updatedMessages, err := AgentLoop(systemPrompt, messages)
		if err != nil {
			logger.Error(fmt.Sprintf("Error: %v", err))
			continue
		}

		messages = updatedMessages
		logger.Assistant(reply)
	}
}
