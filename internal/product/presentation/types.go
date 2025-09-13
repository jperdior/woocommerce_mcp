package presentation

// MCP Protocol types for the presentation layer

// Tool represents an MCP tool definition
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

// InputSchema represents the input schema for an MCP tool
type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

// Property represents a property in the input schema
type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// ToolCall represents a call to an MCP tool
type ToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// ToolResult represents the result of an MCP tool call
type ToolResult struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

// Content represents content in an MCP tool result
type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// NewToolResult creates a new tool result with text content
func NewToolResult(text string) ToolResult {
	return ToolResult{
		Content: []Content{{Type: "text", Text: text}},
		IsError: false,
	}
}

// NewErrorToolResult creates a new error tool result
func NewErrorToolResult(text string) ToolResult {
	return ToolResult{
		Content: []Content{{Type: "text", Text: text}},
		IsError: true,
	}
}

// AddContent adds content to the tool result
func (tr *ToolResult) AddContent(contentType, text string) {
	tr.Content = append(tr.Content, Content{
		Type: contentType,
		Text: text,
	})
}

// SetError sets the error flag on the tool result
func (tr *ToolResult) SetError(isError bool) {
	tr.IsError = isError
}

// IsEmpty checks if the tool result has no content
func (tr *ToolResult) IsEmpty() bool {
	return len(tr.Content) == 0
}

// GetFirstContent returns the first content item, or empty string if none
func (tr *ToolResult) GetFirstContent() string {
	if len(tr.Content) == 0 {
		return ""
	}
	return tr.Content[0].Text
}

// GetAllContent returns all content concatenated
func (tr *ToolResult) GetAllContent() string {
	var result string
	for i, content := range tr.Content {
		if i > 0 {
			result += "\n"
		}
		result += content.Text
	}
	return result
}
