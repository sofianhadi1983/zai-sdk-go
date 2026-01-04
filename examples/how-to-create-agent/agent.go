package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/sofianhadi1983/zai-sdk-go/api/types/chat"
	"github.com/sofianhadi1983/zai-sdk-go/pkg/zai"
)

type Agent struct {
	client       *zai.Client
	tools        ToolRegistry
	conversation []chat.Message
	config       *Config
}

func NewAgent(client *zai.Client, tools ToolRegistry, config *Config) *Agent {
	return &Agent{
		client:       client,
		tools:        tools,
		conversation: make([]chat.Message, 0),
		config:       config,
	}
}

func (a *Agent) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n\nShutting down gracefully...")

		if a.config.SaveConversation {
			if err := a.saveConversation(); err != nil {
				fmt.Printf("Warning: Failed to save conversation: %v\n", err)
			} else {
				fmt.Println("Conversation saved to conversation_history.json")
			}
		}

		cancel()
	}()

	systemMsg := chat.NewSystemMessage(
		"You are a helpful assistant with access to file operations. " +
			"You can read files, list directories, and write files. " +
			"Use these tools to help users manage their files effectively.",
	)
	a.conversation = append(a.conversation, systemMsg)

	a.colorPrintln(colorHeader, "=== File Agent with glm-4.7 ===")
	a.colorPrintln(colorInfo, "Available tools:")
	fmt.Println("  - read_file: Read file contents")
	fmt.Println("  - list_directory: List directory contents")
	fmt.Println("  - write_file: Write content to file")
	fmt.Println()
	a.colorPrintln(colorInfo, "Example commands:")
	fmt.Println("  \"What files are in the current directory?\"")
	fmt.Println("  \"Read the contents of README.md\"")
	fmt.Println("  \"Create a new file called hello.txt with 'Hello, World!'\"")
	fmt.Println()
	a.colorPrintln(colorInfo, "Type your request (Ctrl+C to exit):")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		a.colorPrint(colorUser, "You: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		if err := a.processUserInput(ctx, input); err != nil {
			if ctx.Err() != nil {
				return nil
			}
			fmt.Printf("Error: %v\n\n", err)
			continue
		}

		a.trimConversation()
	}
}

func (a *Agent) processUserInput(ctx context.Context, input string) error {
	userMsg := chat.NewUserMessage(input)
	a.conversation = append(a.conversation, userMsg)

	for {
		resp, toolsUsed, err := a.performInference(ctx)
		if err != nil {
			return fmt.Errorf("inference failed: %w", err)
		}

		if !toolsUsed {
			if !a.config.EnableStreaming {
				fmt.Println()
				a.colorPrint(colorAssistant, "Assistant: ")
				fmt.Printf("%s", resp.GetContent())
				fmt.Println()
				fmt.Println()
			}
			break
		}
	}

	return nil
}

func (a *Agent) performInferenceStreaming(ctx context.Context) (*chat.ChatCompletionResponse, bool, error) {
	a.logDebug("Starting streaming inference with %d messages in conversation", len(a.conversation))

	req := &chat.ChatCompletionRequest{
		Model:    "glm-4.7",
		Messages: a.conversation,
	}

	sdkTools := a.tools.ToSDKTools()
	if len(sdkTools) > 0 {
		req.Tools = sdkTools
		a.logDebug("Sending %d tools to LLM", len(sdkTools))
	}

	a.logDebug("Calling LLM streaming API...")
	spinnerDone := a.showSpinner("Thinking...")

	stream, err := a.client.Chat.CreateStream(ctx, req)
	a.stopSpinner(spinnerDone)

	if err != nil {
		a.logError("Streaming API call failed: %v", err)
		return nil, false, fmt.Errorf("API call failed: %w", err)
	}
	defer stream.Close()

	a.logDebug("Streaming API call successful")

	fmt.Println()
	a.colorPrint(colorAssistant, "Assistant: ")

	var fullContent strings.Builder
	var toolCalls []chat.ToolCall
	var hasToolCalls bool
	var hasContent bool

	for stream.Next() {
		chunk := stream.Current()
		if chunk == nil {
			continue
		}

		if len(chunk.Choices) == 0 {
			continue
		}

		delta := chunk.Choices[0].Delta

		if delta.Content != "" {
			a.typewriterEffect(delta.Content)
			fullContent.WriteString(delta.Content)
			hasContent = true
		}

		if len(delta.ToolCalls) > 0 {
			hasToolCalls = true
			toolCalls = append(toolCalls, delta.ToolCalls...)
		}
	}

	if err := stream.Err(); err != nil {
		a.logError("Stream error: %v", err)
		return nil, false, fmt.Errorf("stream error: %w", err)
	}

	if hasContent {
		fmt.Println()
		fmt.Println()
	}

	resp := &chat.ChatCompletionResponse{
		Choices: []chat.Choice{
			{
				Message: chat.Message{
					Role:      chat.RoleAssistant,
					Content:   fullContent.String(),
					ToolCalls: toolCalls,
				},
			},
		},
	}

	if hasToolCalls {
		assistantMsg := chat.Message{
			Role:      chat.RoleAssistant,
			Content:   fullContent.String(),
			ToolCalls: toolCalls,
		}
		a.conversation = append(a.conversation, assistantMsg)

		if err := a.executeToolCalls(toolCalls); err != nil {
			return nil, true, fmt.Errorf("tool execution failed: %w", err)
		}

		return resp, true, nil
	}

	assistantMsg := chat.NewAssistantMessage(fullContent.String())
	a.conversation = append(a.conversation, assistantMsg)

	return resp, false, nil
}

func (a *Agent) performInference(ctx context.Context) (*chat.ChatCompletionResponse, bool, error) {
	if a.config.EnableStreaming {
		return a.performInferenceStreaming(ctx)
	}

	a.logDebug("Starting inference with %d messages in conversation", len(a.conversation))

	req := &chat.ChatCompletionRequest{
		Model:    "glm-4.7",
		Messages: a.conversation,
	}

	sdkTools := a.tools.ToSDKTools()
	if len(sdkTools) > 0 {
		req.Tools = sdkTools
		a.logDebug("Sending %d tools to LLM", len(sdkTools))
	}

	a.logDebug("Calling LLM API...")
	spinnerDone := a.showSpinner("Thinking...")
	resp, err := a.client.Chat.Create(ctx, req)
	a.stopSpinner(spinnerDone)

	if err != nil {
		a.logError("API call failed: %v", err)
		return nil, false, fmt.Errorf("API call failed: %w", err)
	}
	a.logDebug("API call successful")

	choice := resp.GetFirstChoice()
	if choice != nil && len(choice.Message.ToolCalls) > 0 {
		toolCalls := choice.Message.ToolCalls

		assistantMsg := chat.Message{
			Role:      chat.RoleAssistant,
			Content:   resp.GetContent(),
			ToolCalls: toolCalls,
		}
		a.conversation = append(a.conversation, assistantMsg)

		if err := a.executeToolCalls(toolCalls); err != nil {
			return nil, true, fmt.Errorf("tool execution failed: %w", err)
		}

		return resp, true, nil
	}

	assistantMsg := chat.NewAssistantMessage(resp.GetContent())
	a.conversation = append(a.conversation, assistantMsg)

	return resp, false, nil
}

func (a *Agent) executeToolCalls(toolCalls []chat.ToolCall) error {
	a.logInfo("Executing %d tool call(s)", len(toolCalls))

	for _, tc := range toolCalls {
		a.colorPrint(colorTool, "[Tool Call] ")
		fmt.Printf("%s with args: %s\n", tc.Function.Name, tc.Function.Arguments)
		a.logDebug("Tool: %s, Args: %s", tc.Function.Name, tc.Function.Arguments)

		var args map[string]interface{}
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
			errMsg := fmt.Sprintf("Failed to parse tool arguments: %v", err)
			toolResult := chat.NewToolMessage(tc.ID, errMsg)
			a.conversation = append(a.conversation, toolResult)
			a.colorPrint(colorError, "[Tool Error] ")
			fmt.Printf("%s\n", errMsg)
			continue
		}

		tool, exists := a.tools[tc.Function.Name]
		if !exists {
			errMsg := fmt.Sprintf("Unknown tool: %s", tc.Function.Name)
			toolResult := chat.NewToolMessage(tc.ID, errMsg)
			a.conversation = append(a.conversation, toolResult)
			a.colorPrint(colorError, "[Tool Error] ")
			fmt.Printf("%s\n", errMsg)
			continue
		}

		result, err := tool.Handler(args)
		if err != nil {
			errMsg := fmt.Sprintf("Tool execution error: %v", err)
			toolResult := chat.NewToolMessage(tc.ID, errMsg)
			a.conversation = append(a.conversation, toolResult)
			a.colorPrint(colorError, "[Tool Error] ")
			fmt.Printf("%s\n", errMsg)
			continue
		}

		toolResult := chat.NewToolMessage(tc.ID, result)
		a.conversation = append(a.conversation, toolResult)

		resultPreview := result
		if len(resultPreview) > 100 {
			resultPreview = resultPreview[:100] + "..."
		}
		a.colorPrint(colorSuccess, "[Tool Result] ")
		fmt.Printf("%s\n", resultPreview)
	}

	return nil
}

func (a *Agent) trimConversation() {
	if a.config.MaxMessages <= 0 {
		return
	}

	if len(a.conversation) <= 1 {
		return
	}

	if len(a.conversation) > a.config.MaxMessages {
		systemMsg := a.conversation[0]
		recentMessages := a.conversation[len(a.conversation)-(a.config.MaxMessages-1):]

		a.conversation = make([]chat.Message, 0, a.config.MaxMessages)
		a.conversation = append(a.conversation, systemMsg)
		a.conversation = append(a.conversation, recentMessages...)

		a.logInfo("Trimmed conversation to %d messages (keeping system + recent)", len(a.conversation))
	}
}

func (a *Agent) saveConversation() error {
	filename := "conversation_history.json"

	saveData := struct {
		Timestamp    time.Time      `json:"timestamp"`
		MessageCount int            `json:"message_count"`
		Messages     []chat.Message `json:"messages"`
	}{
		Timestamp:    time.Now(),
		MessageCount: len(a.conversation),
		Messages:     a.conversation,
	}

	data, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal conversation: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (a *Agent) logInfo(format string, args ...interface{}) {
	if !a.config.EnableLogging {
		return
	}
	if a.config.LogLevel == "info" || a.config.LogLevel == "debug" {
		log.Printf("[INFO] "+format, args...)
	}
}

func (a *Agent) logDebug(format string, args ...interface{}) {
	if !a.config.EnableLogging {
		return
	}
	if a.config.LogLevel == "debug" {
		log.Printf("[DEBUG] "+format, args...)
	}
}

func (a *Agent) logError(format string, args ...interface{}) {
	if !a.config.EnableLogging {
		return
	}
	log.Printf("[ERROR] "+format, args...)
}

func (a *Agent) colorPrint(c *color.Color, format string, args ...interface{}) {
	if a.config.EnableColors {
		c.Printf(format, args...)
	} else {
		fmt.Printf(format, args...)
	}
}

func (a *Agent) colorPrintln(c *color.Color, args ...interface{}) {
	if a.config.EnableColors {
		c.Println(args...)
	} else {
		fmt.Println(args...)
	}
}

var (
	colorHeader    = color.New(color.FgCyan, color.Bold)
	colorUser      = color.New(color.FgGreen)
	colorAssistant = color.New(color.FgBlue)
	colorTool      = color.New(color.FgYellow)
	colorError     = color.New(color.FgRed)
	colorSuccess   = color.New(color.FgGreen)
	colorInfo      = color.New(color.FgCyan)
)

func (a *Agent) showSpinner(message string) chan bool {
	if !a.config.ShowProgress {
		return nil
	}

	done := make(chan bool)
	go func() {
		spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-done:
				fmt.Printf("\r%s\r", strings.Repeat(" ", len(message)+5))
				return
			default:
				fmt.Printf("\r%s %s ", spinner[i], message)
				i = (i + 1) % len(spinner)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	return done
}

func (a *Agent) stopSpinner(done chan bool) {
	if done != nil {
		done <- true
		close(done)
	}
}

func (a *Agent) typewriterEffect(text string) {
	if !a.config.EnableStreaming || a.config.TypingSpeedMs <= 0 {
		fmt.Print(text)
		return
	}

	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(time.Duration(a.config.TypingSpeedMs) * time.Millisecond)
	}
}
