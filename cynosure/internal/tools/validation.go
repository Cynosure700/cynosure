package tools

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
)

// resolvedSchemaCache 以原始 schema 字节为键缓存已编译的 schema，
// 这样重复的工具调用就不必反复解析同一个 schema。
var (
	resolvedSchemaMu    sync.Mutex
	resolvedSchemaCache = map[string]*jsonschema.Resolved{}
)

// ValidateToolArgs 依据工具声明的 JSON Schema（params）校验已解析的工具参数。
// 它通过 github.com/google/jsonschema-go 执行完整的 JSON Schema 校验
// （必填项、类型、枚举、嵌套对象/数组等）。
//
// 行为：
//   - params 为空/缺失，或 params 对象没有任何约束：直接放行。
//   - schema 格式错误（无法解析或解析失败）：直接放行，
//     因为这属于工具定义的编写错误，而非 LLM 的错误。
//   - 其余情况下返回带有清晰、便于 LLM 阅读的错误信息，
//     以便模型修正参数后重试。
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

// resolveSchema 解析并求解 schema 字节，当 schema 为空或无效时返回 false
// （此时应跳过校验）。
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

// compileSchema 反序列化并求解 schema。当 schema 格式错误时返回 nil，
// 以便调用方跳过校验。
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

// trimmedSchemaBytes 将 params 规整为原始 JSON 字节，把 "null"
// 和空值都当作“无 schema”处理。
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

// RawSchemaFromParameters 从 openai.FunctionDefinition.Parameters 的值中提取
// JSON Schema 字节。该字段类型为 any，但在本项目中实际保存的是 json.RawMessage。
// 对其他形态会尽力进行序列化。
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
