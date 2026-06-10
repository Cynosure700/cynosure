package tools

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
)

// resolvedSchemaCache caches compiled schemas keyed by the raw schema bytes so
// that repeated tool calls do not re-resolve the same schema.
var (
	resolvedSchemaMu    sync.Mutex
	resolvedSchemaCache = map[string]*jsonschema.Resolved{}
)

// ValidateToolArgs validates the parsed tool arguments against the tool's
// declared JSON Schema (params). It performs a full JSON Schema validation
// (required, type, enum, nested objects/arrays, etc.) via
// github.com/google/jsonschema-go.
//
// Behaviour:
//   - Empty/absent params or a params object with no constraints: pass through.
//   - A malformed schema (one that cannot be parsed or resolved): pass through,
//     since that is an authoring bug in the tool definition, not an LLM error.
//   - Otherwise validation errors are returned with a clear, LLM-readable
//     message so the model can correct its arguments and retry.
func ValidateToolArgs(name string, params json.RawMessage, args map[string]any) error {
	resolved, ok := resolveSchema(params)
	if !ok {
		return nil
	}
	if err := resolved.Validate(args); err != nil {
		return fmt.Errorf("invalid arguments for tool %q: %v", name, err)
	}
	return nil
}

// resolveSchema parses and resolves the schema bytes, returning false when the
// schema is empty or invalid (in which case validation should be skipped).
func resolveSchema(params json.RawMessage) (*jsonschema.Resolved, bool) {
	raw := trimmedSchemaBytes(params)
	if len(raw) == 0 {
		return nil, false
	}

	key := string(raw)
	resolvedSchemaMu.Lock()
	defer resolvedSchemaMu.Unlock()
	if cached, ok := resolvedSchemaCache[key]; ok {
		return cached, cached != nil
	}

	resolved := compileSchema(raw)
	resolvedSchemaCache[key] = resolved
	return resolved, resolved != nil
}

// compileSchema unmarshals and resolves the schema. It returns nil when the
// schema is malformed so the caller can skip validation.
func compileSchema(raw []byte) *jsonschema.Resolved {
	var schema jsonschema.Schema
	if err := json.Unmarshal(raw, &schema); err != nil {
		return nil
	}
	resolved, err := schema.Resolve(nil)
	if err != nil {
		return nil
	}
	return resolved
}

// trimmedSchemaBytes normalises the params into raw JSON bytes, treating "null"
// and empty values as "no schema".
func trimmedSchemaBytes(params json.RawMessage) []byte {
	raw := []byte(params)
	if len(raw) == 0 {
		return nil
	}
	if string(raw) == "null" {
		return nil
	}
	return raw
}

// RawSchemaFromParameters extracts the JSON Schema bytes from an
// openai.FunctionDefinition.Parameters value, which is typed as any but holds
// json.RawMessage in this project. Other shapes are marshalled best-effort.
func RawSchemaFromParameters(params any) json.RawMessage {
	switch v := params.(type) {
	case nil:
		return nil
	case json.RawMessage:
		return v
	case []byte:
		return v
	case string:
		return json.RawMessage(v)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil
		}
		return b
	}
}
