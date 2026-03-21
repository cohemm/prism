package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// filterEnv returns os.Environ() with the specified keys removed.
func filterEnv(keys ...string) []string {
	env := os.Environ()
	filtered := make([]string, 0, len(env))
	for _, e := range env {
		skip := false
		for _, key := range keys {
			if strings.HasPrefix(e, key+"=") {
				skip = true
				break
			}
		}
		if !skip {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// ---------------------------------------------------------------------------
// Legacy wrappers — delegate to QuerySync from claude_sdk.go
// ---------------------------------------------------------------------------

// queryLLM calls the claude CLI as a subprocess.
func queryLLM(ctx context.Context, prompt string) (string, error) {
	return QuerySync(ctx, prompt, ClaudeOptions{
		MaxTurns: 1,
	})
}

// queryLLMWithSystemPrompt calls the claude CLI with a separate system prompt.
func queryLLMWithSystemPrompt(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	return QuerySync(ctx, userPrompt, ClaudeOptions{
		SystemPrompt: systemPrompt,
		MaxTurns:     1,
	})
}

// queryLLMScoped calls the claude CLI with task-scoped isolation.
func queryLLMScoped(ctx context.Context, stateDir, model, prompt string) (string, error) {
	return QuerySync(ctx, prompt, ClaudeOptions{
		Model:    model,
		Cwd:      stateDir,
		MaxTurns: 1,
	})
}

// queryLLMScopedWithSystemPrompt calls the claude CLI with task-scoped isolation
// and a separate system prompt.
func queryLLMScopedWithSystemPrompt(ctx context.Context, stateDir, model, systemPrompt, userPrompt string) (string, error) {
	return QuerySync(ctx, userPrompt, ClaudeOptions{
		Model:        model,
		SystemPrompt: systemPrompt,
		Cwd:          stateDir,
		MaxTurns:     1,
	})
}

// queryLLMScopedWithSchema calls the claude CLI with --json-schema for structured
// output. Single-turn, no tool access.
func queryLLMScopedWithSchema(ctx context.Context, stateDir, model, jsonSchema, prompt string) (string, error) {
	return QuerySync(ctx, prompt, ClaudeOptions{
		Model:      model,
		JSONSchema: jsonSchema,
		Cwd:        stateDir,
		MaxTurns:   1,
	})
}

// queryLLMScopedWithToolsAndSchema calls the claude CLI with tool access and
// --json-schema for structured output. Multi-turn mode — timeout controls duration.
func queryLLMScopedWithToolsAndSchema(ctx context.Context, stateDir, model, jsonSchema, systemPrompt, userPrompt string, _ int) (string, error) {
	return QuerySync(ctx, userPrompt, ClaudeOptions{
		Model:          model,
		SystemPrompt:   systemPrompt,
		JSONSchema:     jsonSchema,
		PermissionMode: "bypassPermissions",
		Cwd:            stateDir,
		// No MaxTurns — context timeout is the safety net
	})
}

// extractJSON extracts a valid JSON object from output that may contain
// surrounding text (markdown fences, explanatory text, etc.).
// With --json-schema + stream-json, structured_output should be clean JSON,
// but this provides robustness for edge cases.
func extractJSON(s string) (string, error) {
	s = strings.TrimSpace(s)

	// Fast path: output is already valid JSON
	if json.Valid([]byte(s)) {
		return s, nil
	}

	// Strip markdown code fences if present (```json\n...\n```)
	if strings.HasPrefix(s, "```") {
		lines := strings.Split(s, "\n")
		if len(lines) > 2 {
			inner := strings.Join(lines[1:len(lines)-1], "\n")
			inner = strings.TrimSpace(inner)
			if json.Valid([]byte(inner)) {
				return inner, nil
			}
		}
	}

	// Find first '{' or '[' and its matching closer with proper nesting
	objStart := strings.Index(s, "{")
	arrStart := strings.Index(s, "[")

	// Pick whichever comes first
	start := objStart
	if start < 0 || (arrStart >= 0 && arrStart < start) {
		start = arrStart
	}
	if start < 0 {
		return "", fmt.Errorf("no JSON object or array found in output (len=%d)", len(s))
	}

	openCh := s[start]
	var closeCh byte
	if openCh == '{' {
		closeCh = '}'
	} else {
		closeCh = ']'
	}

	depth := 0
	inString := false
	escaped := false
	for i := start; i < len(s); i++ {
		ch := s[i]
		if escaped {
			escaped = false
			continue
		}
		if ch == '\\' && inString {
			escaped = true
			continue
		}
		if ch == '"' {
			inString = !inString
			continue
		}
		if inString {
			continue
		}
		switch ch {
		case openCh:
			depth++
		case closeCh:
			depth--
			if depth == 0 {
				candidate := s[start : i+1]
				if json.Valid([]byte(candidate)) {
					return candidate, nil
				}
			}
		}
	}

	return "", fmt.Errorf("no valid JSON found in output (len=%d, first 200 chars: %.200s)", len(s), s)
}
