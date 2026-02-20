# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.2] - 2025-02-15

### Changed

- **ToolCall.Index** – Type changed from `int` to `*int` to distinguish absent from explicitly zero when accumulating streaming deltas

## [1.2.1] - 2025-02-15

### Fixed

- **FileData** – Correct JSON tag from `fileData` to `file_data` for OpenRouter API compatibility

## [1.2.0] - 2025-02-15

### Added

- **ContentPart** – `VideoURL`, `InputAudio`, and `File` fields for multimodal content (video, audio, documents)
- **InputAudio** – Type for audio input in speech-capable models
- **FileData** – Type for document files (PDF, etc.) in OpenRouter requests

## [1.1.0] - 2025-02-15

### Added

- **ToolCall.Index** – `Index` field on `ToolCall` for accumulating streaming deltas (OpenAI format)

## [1.0.0] - 2025-02-14

### Added

- **Chat completions** – Non-streaming and streaming chat via `client.Chat.Create` and `client.Chat.CreateStream`
- **Embeddings** – Text embeddings via `client.Embeddings.Create`
- **Models** – List, filter, and discover models with `client.Models.List`, `ByProvider`, `ByContextLength`, `Cheapest`, `SupportsVision`, and `Get`
- **Cost estimation** – Estimate token costs with `client.EstimateCost`
- **Batch chat** – Run multiple chat requests concurrently with `client.BatchChat`
- **Provider preferences** – Configure provider routing (order, fallbacks, max price, etc.)
- **Debug support** – Request/response logging and `DebugCurl` for exporting requests as curl commands
- **Retry with backoff** – Automatic retries with exponential backoff and jitter for 429, 503, and 408 responses
- **App attribution** – `WithReferer`, `WithTitle`, and `WithForwardedFor` for OpenRouter app tracking
- **Error types** – `OpenRouterError` with `Retryable()`, `Is()`, and sentinel errors for common codes (rate limit, auth, model not found, etc.)

### Dependencies

- Go 1.23.2+
- github.com/MetaDiv-AI/http_caller v1.0.0
- github.com/MetaDiv-AI/logger v1.0.0
