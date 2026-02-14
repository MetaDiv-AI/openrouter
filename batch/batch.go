package batch

import (
	"context"
	"sync"

	"github.com/MetaDiv-AI/openrouter/chat"
)

// ChatBatchProcessor runs multiple chat requests concurrently.
type ChatBatchProcessor struct {
	client      *chat.Service
	concurrency int
}

// NewChatBatchProcessor creates a batch processor for chat requests.
func NewChatBatchProcessor(chatService *chat.Service, concurrency int) *ChatBatchProcessor {
	if concurrency <= 0 {
		concurrency = 5
	}
	return &ChatBatchProcessor{
		client:      chatService,
		concurrency: concurrency,
	}
}

// Run executes the requests concurrently and returns responses in the same order.
func (b *ChatBatchProcessor) Run(ctx context.Context, requests []*chat.ChatRequest) ([]*chat.ChatResponse, []error) {
	if len(requests) == 0 {
		return nil, nil
	}

	results := make([]*chat.ChatResponse, len(requests))
	errs := make([]error, len(requests))
	sem := make(chan struct{}, b.concurrency)
	var wg sync.WaitGroup

	for i, req := range requests {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, r *chat.ChatRequest) {
			defer wg.Done()
			defer func() { <-sem }()
			resp, err := b.client.Create(ctx, r)
			results[idx] = resp
			errs[idx] = err
		}(i, req)
	}

	wg.Wait()
	return results, errs
}
