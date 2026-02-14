package internal

import "bytes"

// ParseSSELine extracts the JSON payload from an SSE "data:" line.
// Returns (payload, true) if the line contains "[DONE]", (payload, false) otherwise.
// Skips comment lines (starting with ":") and empty lines.
func ParseSSELine(line []byte) ([]byte, bool) {
	line = bytes.TrimSpace(line)
	if len(line) == 0 {
		return nil, false
	}
	if line[0] == ':' {
		return nil, false
	}
	if !bytes.HasPrefix(line, []byte("data: ")) {
		return nil, false
	}
	payload := bytes.TrimSpace(line[6:])
	if bytes.Equal(payload, []byte("[DONE]")) {
		return nil, true
	}
	return payload, false
}
