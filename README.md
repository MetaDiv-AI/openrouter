# OpenRouter Go SDK

Go client for the [OpenRouter](https://openrouter.ai) API, built with [http_caller](https://github.com/MetaDiv-AI/http_caller) and [logger](https://github.com/MetaDiv-AI/logger).

## Installation

```bash
go get github.com/MetaDiv-AI/openrouter
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/MetaDiv-AI/logger"
    "github.com/MetaDiv-AI/openrouter"
)

func main() {
    log := logger.New().Development().Build()
    defer log.Sync()

    client, err := openrouter.NewClient(
        openrouter.WithAPIKey(os.Getenv("OPENROUTER_API_KEY")),
        openrouter.WithReferer("https://myapp.com"),
        openrouter.WithTitle("My App"),
        openrouter.WithLogger(log),
        openrouter.WithMaxRetries(3),
    )
    if err != nil {
        panic(err)
    }

    ctx := context.Background()

    // Chat completion
    resp, err := client.Chat.Create(ctx, &openrouter.ChatRequest{
        Model: "anthropic/claude-sonnet-4",
        Messages: []openrouter.Message{{Role: "user", Content: "Hello"}},
    })
    if err != nil {
        panic(err)
    }
    fmt.Println(resp.Choices[0].Message.Content)
}
```

## Streaming

```go
import "io"

stream, err := client.Chat.CreateStream(ctx, &openrouter.ChatRequest{
    Model:    "anthropic/claude-sonnet-4",
    Messages: []openrouter.Message{{Role: "user", Content: "Hello"}},
})
if err != nil {
    panic(err)
}
defer stream.Close()

for {
    chunk, err := stream.Next()
    if err == io.EOF {
        break
    }
    if err != nil {
        panic(err)
    }
    if chunk != nil && len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
        fmt.Print(chunk.Choices[0].Delta.Content)
    }
}
```

## Models

```go
models, _ := client.Models.List(ctx)
cheapest, _ := client.Models.Cheapest(ctx)
visionModels, _ := client.Models.SupportsVision(ctx)
byProvider, _ := client.Models.ByProvider(ctx, "anthropic")
```

## Cost Estimation

```go
cost, _ := client.EstimateCost(ctx, "anthropic/claude-sonnet-4", 100, 50)
```

## Batch Chat

```go
requests := []*openrouter.ChatRequest{
    {Model: "openai/gpt-4", Messages: []openrouter.Message{{Role: "user", Content: "Hi 1"}}},
    {Model: "openai/gpt-4", Messages: []openrouter.Message{{Role: "user", Content: "Hi 2"}}},
}
responses, errs := client.BatchChat(ctx, requests, 5)
```

## Embeddings

```go
resp, err := client.Embeddings.Create(ctx, &openrouter.CreateRequest{
    Model: "openai/text-embedding-3-small",
    Input: "The quick brown fox",
})
```

## Debug

```go
client, _ := openrouter.NewClient(
    openrouter.WithAPIKey(apiKey),
    openrouter.WithDebug(true),  // enables request/response logging via logger
)

// Export as curl command
curlCmd := client.DebugCurl(&openrouter.ChatRequest{
    Model: "anthropic/claude-sonnet-4",
    Messages: []openrouter.Message{{Role: "user", Content: "Hello"}},
})
fmt.Println(curlCmd)
```

## License

MIT
