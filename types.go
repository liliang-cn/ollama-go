package ollama

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Message represents a chat message
type Message struct {
	Role      string      `json:"role"`
	Content   string      `json:"content,omitempty"`
	Images    []Image     `json:"images,omitempty"`
	ToolCalls []ToolCall  `json:"tool_calls,omitempty"`
	ToolName  string      `json:"tool_name,omitempty"`
	Thinking  string      `json:"thinking,omitempty"`
}

// Image represents an image input
type Image struct {
	Data string `json:"-"`
}

// MarshalJSON implements json.Marshaler for Image
func (i Image) MarshalJSON() ([]byte, error) {
	// Handle base64 encoded data
	if strings.HasPrefix(i.Data, "data:image/") {
		// Extract base64 part from data URI
		parts := strings.Split(i.Data, ",")
		if len(parts) == 2 {
			return json.Marshal(parts[1])
		}
	}

	// Handle file path
	if _, err := os.Stat(i.Data); err == nil {
		data, err := os.ReadFile(i.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to read image file: %w", err)
		}
		encoded := base64.StdEncoding.EncodeToString(data)
		return json.Marshal(encoded)
	}

	// Assume it's already base64 encoded
	return json.Marshal(i.Data)
}

// ToolCall represents a tool call
type ToolCall struct {
	Function Function `json:"function"`
}

// Function represents a function call
type Function struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// Tool represents a tool definition
type Tool struct {
	Type     string       `json:"type,omitempty"`
	Function *ToolFunction `json:"function,omitempty"`
}

// ToolFunction represents a tool function definition
type ToolFunction struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// Options contains model configuration options
type Options struct {
	// Load time options
	Numa           *bool `json:"numa,omitempty"`
	NumCtx         *int  `json:"num_ctx,omitempty"`
	NumBatch       *int  `json:"num_batch,omitempty"`
	NumGPU         *int  `json:"num_gpu,omitempty"`
	MainGPU        *int  `json:"main_gpu,omitempty"`
	LowVRAM        *bool `json:"low_vram,omitempty"`
	F16KV          *bool `json:"f16_kv,omitempty"`
	LogitsAll      *bool `json:"logits_all,omitempty"`
	VocabOnly      *bool `json:"vocab_only,omitempty"`
	UseMmap        *bool `json:"use_mmap,omitempty"`
	UseMlock       *bool `json:"use_mlock,omitempty"`
	EmbeddingOnly  *bool `json:"embedding_only,omitempty"`
	NumThread      *int  `json:"num_thread,omitempty"`

	// Runtime options
	NumKeep         *int     `json:"num_keep,omitempty"`
	Seed            *int     `json:"seed,omitempty"`
	NumPredict      *int     `json:"num_predict,omitempty"`
	TopK            *int     `json:"top_k,omitempty"`
	TopP            *float64 `json:"top_p,omitempty"`
	MinP            *float64 `json:"min_p,omitempty"`
	TFSZ            *float64 `json:"tfs_z,omitempty"`
	TypicalP        *float64 `json:"typical_p,omitempty"`
	RepeatLastN     *int     `json:"repeat_last_n,omitempty"`
	Temperature     *float64 `json:"temperature,omitempty"`
	RepeatPenalty   *float64 `json:"repeat_penalty,omitempty"`
	PresencePenalty *float64 `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64 `json:"frequency_penalty,omitempty"`
	Mirostat        *int     `json:"mirostat,omitempty"`
	MirostatTau     *float64 `json:"mirostat_tau,omitempty"`
	MirostatEta     *float64 `json:"mirostat_eta,omitempty"`
	PenalizeNewline *bool    `json:"penalize_newline,omitempty"`
	Stop            []string `json:"stop,omitempty"`
}

// GenerateRequest represents a generation request
type GenerateRequest struct {
	Model     string                 `json:"model"`
	Prompt    string                 `json:"prompt,omitempty"`
	Suffix    string                 `json:"suffix,omitempty"`
	System    string                 `json:"system,omitempty"`
	Template  string                 `json:"template,omitempty"`
	Context   []int                  `json:"context,omitempty"`
	Stream    *bool                  `json:"stream,omitempty"`
	Raw       *bool                  `json:"raw,omitempty"`
	Format    interface{}            `json:"format,omitempty"`
	Options   *Options               `json:"options,omitempty"`
	Images    []Image                `json:"images,omitempty"`
	KeepAlive interface{}            `json:"keep_alive,omitempty"`
	Think     *bool                  `json:"think,omitempty"`
}

// GenerateResponse represents a generation response
type GenerateResponse struct {
	Model              string `json:"model,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
	Response           string `json:"response"`
	Done               bool   `json:"done,omitempty"`
	DoneReason         string `json:"done_reason,omitempty"`
	Context            []int  `json:"context,omitempty"`
	TotalDuration      int64  `json:"total_duration,omitempty"`
	LoadDuration       int64  `json:"load_duration,omitempty"`
	PromptEvalCount    int    `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64  `json:"prompt_eval_duration,omitempty"`
	EvalCount          int    `json:"eval_count,omitempty"`
	EvalDuration       int64  `json:"eval_duration,omitempty"`
	Thinking           string `json:"thinking,omitempty"`
}

// ChatRequest represents a chat request
type ChatRequest struct {
	Model     string                 `json:"model"`
	Messages  []Message              `json:"messages,omitempty"`
	Tools     []Tool                 `json:"tools,omitempty"`
	Stream    *bool                  `json:"stream,omitempty"`
	Format    interface{}            `json:"format,omitempty"`
	Options   *Options               `json:"options,omitempty"`
	KeepAlive interface{}            `json:"keep_alive,omitempty"`
	Think     *bool                  `json:"think,omitempty"`
}

// ChatResponse represents a chat response
type ChatResponse struct {
	Model              string  `json:"model,omitempty"`
	CreatedAt          string  `json:"created_at,omitempty"`
	Message            Message `json:"message"`
	Done               bool    `json:"done,omitempty"`
	DoneReason         string  `json:"done_reason,omitempty"`
	TotalDuration      int64   `json:"total_duration,omitempty"`
	LoadDuration       int64   `json:"load_duration,omitempty"`
	PromptEvalCount    int     `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64   `json:"prompt_eval_duration,omitempty"`
	EvalCount          int     `json:"eval_count,omitempty"`
	EvalDuration       int64   `json:"eval_duration,omitempty"`
}

// EmbedRequest represents an embedding request
type EmbedRequest struct {
	Model     string                 `json:"model"`
	Input     interface{}            `json:"input"` // string or []string
	Truncate  *bool                  `json:"truncate,omitempty"`
	Options   *Options               `json:"options,omitempty"`
	KeepAlive interface{}            `json:"keep_alive,omitempty"`
}

// EmbedResponse represents an embedding response
type EmbedResponse struct {
	Model              string      `json:"model,omitempty"`
	Embeddings         [][]float64 `json:"embeddings"`
	TotalDuration      int64       `json:"total_duration,omitempty"`
	LoadDuration       int64       `json:"load_duration,omitempty"`
	PromptEvalCount    int         `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64       `json:"prompt_eval_duration,omitempty"`
}

// EmbeddingsRequest represents an embeddings request (legacy)
type EmbeddingsRequest struct {
	Model     string                 `json:"model"`
	Prompt    string                 `json:"prompt,omitempty"`
	Options   *Options               `json:"options,omitempty"`
	KeepAlive interface{}            `json:"keep_alive,omitempty"`
}

// EmbeddingsResponse represents an embeddings response (legacy)
type EmbeddingsResponse struct {
	Embedding []float64 `json:"embedding"`
}

// ListResponse represents a model list response
type ListResponse struct {
	Models []ModelInfo `json:"models"`
}

// ModelInfo represents information about a model
type ModelInfo struct {
	Model      string         `json:"model,omitempty"`
	ModifiedAt *time.Time     `json:"modified_at,omitempty"`
	Digest     string         `json:"digest,omitempty"`
	Size       int64          `json:"size,omitempty"`
	Details    *ModelDetails  `json:"details,omitempty"`
}

// ModelDetails represents detailed model information
type ModelDetails struct {
	ParentModel       string   `json:"parent_model,omitempty"`
	Format            string   `json:"format,omitempty"`
	Family            string   `json:"family,omitempty"`
	Families          []string `json:"families,omitempty"`
	ParameterSize     string   `json:"parameter_size,omitempty"`
	QuantizationLevel string   `json:"quantization_level,omitempty"`
}

// ShowRequest represents a show model request
type ShowRequest struct {
	Model   string `json:"model"`
	Verbose *bool  `json:"verbose,omitempty"`
}

// ShowResponse represents a show model response
type ShowResponse struct {
	ModifiedAt   *time.Time             `json:"modified_at,omitempty"`
	Template     string                 `json:"template,omitempty"`
	Modelfile    string                 `json:"modelfile,omitempty"`
	License      string                 `json:"license,omitempty"`
	Details      *ModelDetails          `json:"details,omitempty"`
	ModelInfo    map[string]interface{} `json:"model_info,omitempty"`
	Parameters   string                 `json:"parameters,omitempty"`
	Capabilities []string               `json:"capabilities,omitempty"`
}

// PullRequest represents a pull model request
type PullRequest struct {
	Model    string `json:"model"`
	Insecure *bool  `json:"insecure,omitempty"`
	Stream   *bool  `json:"stream,omitempty"`
}

// PushRequest represents a push model request
type PushRequest struct {
	Model    string `json:"model"`
	Insecure *bool  `json:"insecure,omitempty"`
	Stream   *bool  `json:"stream,omitempty"`
}

// CreateRequest represents a create model request
type CreateRequest struct {
	Model      string                 `json:"model"`
	Modelfile  string                 `json:"modelfile,omitempty"`
	Quantize   string                 `json:"quantize,omitempty"`
	From       string                 `json:"from,omitempty"`
	Files      map[string]string      `json:"files,omitempty"`
	Adapters   map[string]string      `json:"adapters,omitempty"`
	Template   string                 `json:"template,omitempty"`
	License    interface{}            `json:"license,omitempty"` // string or []string
	System     string                 `json:"system,omitempty"`
	Parameters *Options               `json:"parameters,omitempty"`
	Messages   []Message              `json:"messages,omitempty"`
	Stream     *bool                  `json:"stream,omitempty"`
	Path       string                 `json:"path,omitempty"`
}

// DeleteRequest represents a delete model request
type DeleteRequest struct {
	Model string `json:"model"`
}

// CopyRequest represents a copy model request
type CopyRequest struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

// ProcessResponse represents running processes
type ProcessResponse struct {
	Models []ProcessModel `json:"models"`
}

// ProcessModel represents a running model process
type ProcessModel struct {
	Model         string        `json:"model,omitempty"`
	Name          string        `json:"name,omitempty"`
	Digest        string        `json:"digest,omitempty"`
	ExpiresAt     *time.Time    `json:"expires_at,omitempty"`
	Size          int64         `json:"size,omitempty"`
	SizeVRAM      int64         `json:"size_vram,omitempty"`
	Details       *ModelDetails `json:"details,omitempty"`
	ContextLength int           `json:"context_length,omitempty"`
}

// ProgressResponse represents a progress response for long-running operations
type ProgressResponse struct {
	Status    string `json:"status,omitempty"`
	Digest    string `json:"digest,omitempty"`
	Total     int64  `json:"total,omitempty"`
	Completed int64  `json:"completed,omitempty"`
}

// StatusResponse represents a simple status response
type StatusResponse struct {
	Status string `json:"status,omitempty"`
}

// VersionResponse represents a version response
type VersionResponse struct {
	Version string `json:"version"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// ResponseError represents a response error
type ResponseError struct {
	StatusCode int
	Message    string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("ollama: request failed with status %d: %s", e.StatusCode, e.Message)
}