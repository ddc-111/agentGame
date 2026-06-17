package mcp

import (
	"encoding/json"
	"testing"
)

func TestMCPRequest_Struct(t *testing.T) {
	req := MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/list",
	}

	if req.JSONRPC != "2.0" {
		t.Errorf("MCPRequest.JSONRPC = %v, want 2.0", req.JSONRPC)
	}
	if req.ID != 1 {
		t.Errorf("MCPRequest.ID = %v, want 1", req.ID)
	}
	if req.Method != "tools/list" {
		t.Errorf("MCPRequest.Method = %v, want tools/list", req.Method)
	}
}

func TestMCPResponse_Struct(t *testing.T) {
	resp := MCPResponse{
		JSONRPC: "2.0",
		ID:      1,
		Result:  map[string]string{"status": "ok"},
	}

	if resp.JSONRPC != "2.0" {
		t.Errorf("MCPResponse.JSONRPC = %v, want 2.0", resp.JSONRPC)
	}
	if resp.Error != nil {
		t.Errorf("MCPResponse.Error = %v, want nil", resp.Error)
	}
}

func TestMCPError_Struct(t *testing.T) {
	mcpErr := MCPError{
		Code:    -32600,
		Message: "Invalid Request",
	}

	if mcpErr.Code != -32600 {
		t.Errorf("MCPError.Code = %v, want -32600", mcpErr.Code)
	}
	if mcpErr.Message != "Invalid Request" {
		t.Errorf("MCPError.Message = %v, want Invalid Request", mcpErr.Message)
	}
}

func TestTool_Struct(t *testing.T) {
	tool := Tool{
		Name:        "test_tool",
		Description: "A test tool",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"param1": map[string]string{"type": "string"},
			},
		},
	}

	if tool.Name != "test_tool" {
		t.Errorf("Tool.Name = %v, want test_tool", tool.Name)
	}
	if tool.Description != "A test tool" {
		t.Errorf("Tool.Description = %v, want A test tool", tool.Description)
	}
}

func TestToolCall_Struct(t *testing.T) {
	args := json.RawMessage(`{"param1":"value1"}`)
	call := ToolCall{
		Name:      "test_tool",
		Arguments: args,
	}

	if call.Name != "test_tool" {
		t.Errorf("ToolCall.Name = %v, want test_tool", call.Name)
	}

	var parsed map[string]string
	if err := json.Unmarshal(call.Arguments, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal arguments: %v", err)
	}
	if parsed["param1"] != "value1" {
		t.Errorf("ToolCall.Arguments[param1] = %v, want value1", parsed["param1"])
	}
}

func TestToolResult_Struct(t *testing.T) {
	result := ToolResult{
		Content: []Content{
			{Type: "text", Text: "result text"},
		},
		IsError: false,
	}

	if len(result.Content) != 1 {
		t.Errorf("ToolResult.Content length = %v, want 1", len(result.Content))
	}
	if result.Content[0].Type != "text" {
		t.Errorf("ToolResult.Content[0].Type = %v, want text", result.Content[0].Type)
	}
	if result.IsError != false {
		t.Errorf("ToolResult.IsError = %v, want false", result.IsError)
	}
}

func TestContent_Struct(t *testing.T) {
	content := Content{
		Type: "text",
		Text: "test content",
	}

	if content.Type != "text" {
		t.Errorf("Content.Type = %v, want text", content.Type)
	}
	if content.Text != "test content" {
		t.Errorf("Content.Text = %v, want test content", content.Text)
	}
}

func TestResource_Struct(t *testing.T) {
	resource := Resource{
		URI:         "file:///test.txt",
		Name:        "test.txt",
		Description: "A test file",
		MIMEType:    "text/plain",
	}

	if resource.URI != "file:///test.txt" {
		t.Errorf("Resource.URI = %v, want file:///test.txt", resource.URI)
	}
	if resource.Name != "test.txt" {
		t.Errorf("Resource.Name = %v, want test.txt", resource.Name)
	}
}

func TestResourceContent_Struct(t *testing.T) {
	rc := ResourceContent{
		URI:      "file:///test.txt",
		MIMEType: "text/plain",
		Text:     "content",
	}

	if rc.URI != "file:///test.txt" {
		t.Errorf("ResourceContent.URI = %v, want file:///test.txt", rc.URI)
	}
	if rc.Text != "content" {
		t.Errorf("ResourceContent.Text = %v, want content", rc.Text)
	}
}

func TestPrompt_Struct(t *testing.T) {
	prompt := Prompt{
		Name:        "test_prompt",
		Description: "A test prompt",
		Arguments: []PromptArgument{
			{Name: "arg1", Description: "First argument", Required: true},
		},
	}

	if prompt.Name != "test_prompt" {
		t.Errorf("Prompt.Name = %v, want test_prompt", prompt.Name)
	}
	if len(prompt.Arguments) != 1 {
		t.Errorf("Prompt.Arguments length = %v, want 1", len(prompt.Arguments))
	}
}

func TestPromptArgument_Struct(t *testing.T) {
	arg := PromptArgument{
		Name:        "test_arg",
		Description: "A test argument",
		Required:    true,
	}

	if arg.Name != "test_arg" {
		t.Errorf("PromptArgument.Name = %v, want test_arg", arg.Name)
	}
	if arg.Required != true {
		t.Errorf("PromptArgument.Required = %v, want true", arg.Required)
	}
}

func TestPromptMessage_Struct(t *testing.T) {
	msg := PromptMessage{
		Role: "user",
		Content: Content{
			Type: "text",
			Text: "test",
		},
	}

	if msg.Role != "user" {
		t.Errorf("PromptMessage.Role = %v, want user", msg.Role)
	}
	if msg.Content.Text != "test" {
		t.Errorf("PromptMessage.Content.Text = %v, want test", msg.Content.Text)
	}
}

func TestMCPRequest_JSON(t *testing.T) {
	jsonStr := `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`
	var req MCPRequest

	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Failed to unmarshal MCPRequest: %v", err)
	}

	if req.JSONRPC != "2.0" {
		t.Errorf("MCPRequest.JSONRPC = %v, want 2.0", req.JSONRPC)
	}
	if req.Method != "tools/list" {
		t.Errorf("MCPRequest.Method = %v, want tools/list", req.Method)
	}
}

func TestMCPResponse_JSON(t *testing.T) {
	resp := MCPResponse{
		JSONRPC: "2.0",
		ID:      1,
		Result:  map[string]string{"status": "ok"},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal MCPResponse: %v", err)
	}

	var parsed MCPResponse
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal MCPResponse: %v", err)
	}

	if parsed.JSONRPC != "2.0" {
		t.Errorf("MCPResponse.JSONRPC = %v, want 2.0", parsed.JSONRPC)
	}
}

func TestMCPError_JSON(t *testing.T) {
	jsonStr := `{"code":-32600,"message":"Invalid Request"}`
	var mcpErr MCPError

	if err := json.Unmarshal([]byte(jsonStr), &mcpErr); err != nil {
		t.Fatalf("Failed to unmarshal MCPError: %v", err)
	}

	if mcpErr.Code != -32600 {
		t.Errorf("MCPError.Code = %v, want -32600", mcpErr.Code)
	}
	if mcpErr.Message != "Invalid Request" {
		t.Errorf("MCPError.Message = %v, want Invalid Request", mcpErr.Message)
	}
}

func TestToolResult_JSON(t *testing.T) {
	jsonStr := `{"content":[{"type":"text","text":"result"}],"isError":false}`
	var result ToolResult

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Failed to unmarshal ToolResult: %v", err)
	}

	if len(result.Content) != 1 {
		t.Errorf("ToolResult.Content length = %v, want 1", len(result.Content))
	}
	if result.Content[0].Text != "result" {
		t.Errorf("ToolResult.Content[0].Text = %v, want result", result.Content[0].Text)
	}
}
