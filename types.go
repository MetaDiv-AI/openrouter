package openrouter

import (
	"github.com/MetaDiv-AI/openrouter/chat"
	"github.com/MetaDiv-AI/openrouter/embeddings"
)

// Re-export commonly used types for convenience.

type (
	ChatRequest    = chat.ChatRequest
	ChatResponse   = chat.ChatResponse
	Message        = chat.Message
	Choice         = chat.Choice
	Usage          = chat.Usage
	Tool           = chat.Tool
	ToolCall       = chat.ToolCall
	ResponseFormat = chat.ResponseFormat
	CreateRequest  = embeddings.CreateRequest
	CreateResponse = embeddings.CreateResponse
)
