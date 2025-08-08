package ollama

import "context"

// Default client instance
var defaultClient *Client

func init() {
	var err error
	defaultClient, err = NewClient()
	if err != nil {
		panic("failed to create default ollama client: " + err.Error())
	}
}

// Generate generates a response using the default client.
// It sends a prompt to the specified model and returns the complete response.
//
// Example:
//
//	response, err := ollama.Generate(ctx, "gemma3", "Why is the sky blue?")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(response.Response)
func Generate(ctx context.Context, model, prompt string, options ...func(*GenerateRequest)) (*GenerateResponse, error) {
	req := &GenerateRequest{
		Model:  model,
		Prompt: prompt,
	}
	
	for _, opt := range options {
		opt(req)
	}
	
	return defaultClient.Generate(ctx, req)
}

// GenerateStream generates a streaming response using the default client.
// It returns channels for receiving streaming responses and errors.
// The response channel will be closed when the stream is complete.
//
// Example:
//
//	responseChan, errorChan := ollama.GenerateStream(ctx, "gemma3", "Tell me a story")
//	for {
//		select {
//		case response, ok := <-responseChan:
//			if !ok {
//				return
//			}
//			fmt.Print(response.Response)
//		case err := <-errorChan:
//			if err != nil {
//				log.Fatal(err)
//			}
//		}
//	}
func GenerateStream(ctx context.Context, model, prompt string, options ...func(*GenerateRequest)) (<-chan *GenerateResponse, <-chan error) {
	req := &GenerateRequest{
		Model:  model,
		Prompt: prompt,
	}
	
	for _, opt := range options {
		opt(req)
	}
	
	return defaultClient.GenerateStream(ctx, req)
}

// Chat sends a chat message using the default client.
// It provides a conversational interface where you can maintain message history.
//
// Example:
//
//	messages := []ollama.Message{
//		{Role: "user", Content: "Hello!"},
//	}
//	response, err := ollama.Chat(ctx, "gemma3", messages)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(response.Message.Content)
func Chat(ctx context.Context, model string, messages []Message, options ...func(*ChatRequest)) (*ChatResponse, error) {
	req := &ChatRequest{
		Model:    model,
		Messages: messages,
	}
	
	for _, opt := range options {
		opt(req)
	}
	
	return defaultClient.Chat(ctx, req)
}

// ChatStream sends a chat message with streaming response using the default client.
// Similar to Chat but returns streaming responses for real-time interaction.
//
// Example:
//
//	messages := []ollama.Message{
//		{Role: "user", Content: "Tell me a joke"},
//	}
//	responseChan, errorChan := ollama.ChatStream(ctx, "gemma3", messages)
//	for {
//		select {
//		case response, ok := <-responseChan:
//			if !ok {
//				return
//			}
//			fmt.Print(response.Message.Content)
//		case err := <-errorChan:
//			if err != nil {
//				log.Fatal(err)
//			}
//		}
//	}
func ChatStream(ctx context.Context, model string, messages []Message, options ...func(*ChatRequest)) (<-chan *ChatResponse, <-chan error) {
	req := &ChatRequest{
		Model:    model,
		Messages: messages,
	}
	
	for _, opt := range options {
		opt(req)
	}
	
	return defaultClient.ChatStream(ctx, req)
}

// Embed creates embeddings using the default client.
// It converts text into numerical vectors that can be used for semantic similarity.
//
// The input parameter can be a string or []string for multiple inputs.
//
// Example:
//
//	response, err := ollama.Embed(ctx, "nomic-embed-text", "The quick brown fox")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Generated %d embeddings\n", len(response.Embeddings))
func Embed(ctx context.Context, model string, input interface{}, options ...func(*EmbedRequest)) (*EmbedResponse, error) {
	req := &EmbedRequest{
		Model: model,
		Input: input,
	}
	
	for _, opt := range options {
		opt(req)
	}
	
	return defaultClient.Embed(ctx, req)
}

// Embeddings creates embeddings using the legacy API and default client
func Embeddings(ctx context.Context, model, prompt string, options ...func(*EmbeddingsRequest)) (*EmbeddingsResponse, error) {
	req := &EmbeddingsRequest{
		Model:  model,
		Prompt: prompt,
	}
	
	for _, opt := range options {
		opt(req)
	}
	
	return defaultClient.Embeddings(ctx, req)
}

// List lists all available models using the default client.
// Returns information about each model including name, size, and modification time.
//
// Example:
//
//	models, err := ollama.List(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, model := range models.Models {
//		fmt.Println(model.Name)
//	}
func List(ctx context.Context) (*ListResponse, error) {
	return defaultClient.List(ctx)
}

// Show returns detailed information about a specific model using the default client.
// Includes model parameters, template, and system information.
//
// Example:
//
//	info, err := ollama.Show(ctx, "gemma3")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Model: %s, Size: %d\n", info.Details.Family, info.Details.ParameterSize)
func Show(ctx context.Context, model string) (*ShowResponse, error) {
	return defaultClient.Show(ctx, &ShowRequest{Model: model})
}

// Pull downloads a model using the default client.
// Downloads the specified model from the Ollama registry.
//
// Example:
//
//	err := ollama.Pull(ctx, "gemma3")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Model downloaded successfully")
func Pull(ctx context.Context, model string, options ...func(*PullRequest)) (*StatusResponse, error) {
	req := &PullRequest{Model: model}
	
	for _, opt := range options {
		opt(req)
	}
	
	return defaultClient.Pull(ctx, req)
}

// PullStream downloads a model with progress updates using the default client.
// Returns channels for receiving progress updates and errors.
//
// Example:
//
//	progressChan, errorChan := ollama.PullStream(ctx, "gemma3")
//	for {
//		select {
//		case progress, ok := <-progressChan:
//			if !ok {
//				fmt.Println("Download completed!")
//				return
//			}
//			if progress.Total > 0 {
//				percentage := float64(progress.Completed) / float64(progress.Total) * 100
//				fmt.Printf("Progress: %.1f%%\n", percentage)
//			}
//		case err := <-errorChan:
//			if err != nil {
//				log.Fatal(err)
//			}
//		}
//	}
func PullStream(ctx context.Context, model string, options ...func(*PullRequest)) (<-chan *ProgressResponse, <-chan error) {
	req := &PullRequest{Model: model}
	
	for _, opt := range options {
		opt(req)
	}
	
	return defaultClient.PullStream(ctx, req)
}

// Push uploads a model using the default client
func Push(ctx context.Context, model string, options ...func(*PushRequest)) (*StatusResponse, error) {
	req := &PushRequest{Model: model}
	
	for _, opt := range options {
		opt(req)
	}
	
	return defaultClient.Push(ctx, req)
}

// PushStream uploads a model with progress using the default client
func PushStream(ctx context.Context, model string, options ...func(*PushRequest)) (<-chan *ProgressResponse, <-chan error) {
	req := &PushRequest{Model: model}
	
	for _, opt := range options {
		opt(req)
	}
	
	return defaultClient.PushStream(ctx, req)
}

// Create creates a new model using the default client
func Create(ctx context.Context, model, modelfile string, options ...func(*CreateRequest)) (*StatusResponse, error) {
	req := &CreateRequest{
		Model:     model,
		Modelfile: modelfile,
	}
	
	for _, opt := range options {
		opt(req)
	}
	
	return defaultClient.Create(ctx, req)
}

// CreateStream creates a new model with progress using the default client
func CreateStream(ctx context.Context, model, modelfile string, options ...func(*CreateRequest)) (<-chan *ProgressResponse, <-chan error) {
	req := &CreateRequest{
		Model:     model,
		Modelfile: modelfile,
	}
	
	for _, opt := range options {
		opt(req)
	}
	
	return defaultClient.CreateStream(ctx, req)
}

// Delete deletes a model using the default client
func Delete(ctx context.Context, model string) (*StatusResponse, error) {
	return defaultClient.Delete(ctx, &DeleteRequest{Model: model})
}

// Copy copies a model using the default client
func Copy(ctx context.Context, source, destination string) (*StatusResponse, error) {
	return defaultClient.Copy(ctx, &CopyRequest{Source: source, Destination: destination})
}

// Ps shows running processes using the default client
func Ps(ctx context.Context) (*ProcessResponse, error) {
	return defaultClient.Ps(ctx)
}

// Option functions for common configurations

// WithSystem sets the system prompt
func WithSystem(system string) func(*GenerateRequest) {
	return func(req *GenerateRequest) {
		req.System = system
	}
}

// WithChatSystem sets the system prompt for chat
func WithChatSystem(system string) func(*ChatRequest) {
	return func(req *ChatRequest) {
		// Add system message as the first message
		messages := []Message{{Role: "system", Content: system}}
		messages = append(messages, req.Messages...)
		req.Messages = messages
	}
}

// WithOptions sets the model options
func WithOptions(options *Options) func(interface{}) {
	return func(req interface{}) {
		switch r := req.(type) {
		case *GenerateRequest:
			r.Options = options
		case *ChatRequest:
			r.Options = options
		case *EmbedRequest:
			r.Options = options
		case *EmbeddingsRequest:
			r.Options = options
		}
	}
}

// WithFormat sets the response format
func WithFormat(format interface{}) func(interface{}) {
	return func(req interface{}) {
		switch r := req.(type) {
		case *GenerateRequest:
			r.Format = format
		case *ChatRequest:
			r.Format = format
		}
	}
}

// WithKeepAlive sets the keep alive duration
func WithKeepAlive(keepAlive interface{}) func(interface{}) {
	return func(req interface{}) {
		switch r := req.(type) {
		case *GenerateRequest:
			r.KeepAlive = keepAlive
		case *ChatRequest:
			r.KeepAlive = keepAlive
		case *EmbedRequest:
			r.KeepAlive = keepAlive
		case *EmbeddingsRequest:
			r.KeepAlive = keepAlive
		}
	}
}

// WithImages adds images to the request
func WithImages(images []Image) func(interface{}) {
	return func(req interface{}) {
		switch r := req.(type) {
		case *GenerateRequest:
			r.Images = images
		}
	}
}

// WithTools adds tools to the chat request
func WithTools(tools []Tool) func(*ChatRequest) {
	return func(req *ChatRequest) {
		req.Tools = tools
	}
}

// WithThinking enables thinking mode
func WithThinking() func(interface{}) {
	return func(req interface{}) {
		switch r := req.(type) {
		case *GenerateRequest:
			r.Think = BoolPtr(true)
		case *ChatRequest:
			r.Think = BoolPtr(true)
		}
	}
}

// CreateBlob uploads a file using the default client
func CreateBlob(ctx context.Context, path string) (string, error) {
	return defaultClient.CreateBlob(ctx, path)
}

// WithTemplate sets the template for create requests
func WithTemplate(template string) func(*CreateRequest) {
	return func(req *CreateRequest) {
		req.Template = template
	}
}

// WithLicense sets the license for create requests
func WithLicense(license interface{}) func(*CreateRequest) {
	return func(req *CreateRequest) {
		req.License = license
	}
}

// WithFiles sets the files for create requests
func WithFiles(files map[string]string) func(*CreateRequest) {
	return func(req *CreateRequest) {
		req.Files = files
	}
}

// WithAdapters sets the adapters for create requests
func WithAdapters(adapters map[string]string) func(*CreateRequest) {
	return func(req *CreateRequest) {
		req.Adapters = adapters
	}
}

// WithCreateMessages sets the messages for create requests
func WithCreateMessages(messages []Message) func(*CreateRequest) {
	return func(req *CreateRequest) {
		req.Messages = messages
	}
}

// WithCreateOptions sets the parameters for create requests
func WithCreateOptions(options *Options) func(*CreateRequest) {
	return func(req *CreateRequest) {
		req.Parameters = options
	}
}

// WithCreateSystem sets the system prompt for create requests
func WithCreateSystem(system string) func(*CreateRequest) {
	return func(req *CreateRequest) {
		req.System = system
	}
}

// WithQuantize sets the quantization for create requests
func WithQuantize(quantize string) func(*CreateRequest) {
	return func(req *CreateRequest) {
		req.Quantize = quantize
	}
}

// WithFrom sets the from model for create requests
func WithFrom(from string) func(*CreateRequest) {
	return func(req *CreateRequest) {
		req.From = from
	}
}

// WithRaw sets the raw mode for generate requests
func WithRaw() func(*GenerateRequest) {
	return func(req *GenerateRequest) {
		req.Raw = BoolPtr(true)
	}
}

// WithSuffix sets the suffix for generate requests
func WithSuffix(suffix string) func(*GenerateRequest) {
	return func(req *GenerateRequest) {
		req.Suffix = suffix
	}
}

// WithTemplate sets the template for generate requests
func WithGenerateTemplate(template string) func(*GenerateRequest) {
	return func(req *GenerateRequest) {
		req.Template = template
	}
}

// WithContext sets the context for generate requests
func WithContext(context []int) func(*GenerateRequest) {
	return func(req *GenerateRequest) {
		req.Context = context
	}
}

// WithTruncate sets truncate option for embed requests
func WithTruncate(truncate bool) func(*EmbedRequest) {
	return func(req *EmbedRequest) {
		req.Truncate = BoolPtr(truncate)
	}
}

// WithInsecure sets insecure option for pull/push requests
func WithInsecure(insecure bool) func(interface{}) {
	return func(req interface{}) {
		switch r := req.(type) {
		case *PullRequest:
			r.Insecure = BoolPtr(insecure)
		case *PushRequest:
			r.Insecure = BoolPtr(insecure)
		}
	}
}