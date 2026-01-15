package main

import (
	"fmt"
	"regexp"
)

type NodeType string

const (
	StringType NodeType = "string"
	NumberType NodeType = "number"
	BoolType   NodeType = "boolean"
	ObjectType NodeType = "object"
	ArrayType  NodeType = "array"
	NullType   NodeType = "null"
)

type Node struct {
	Path              string      `json:"path"`
	Type              NodeType    `json:"type"`
	Value             interface{} `json:"value"`
	Parent            string      `json:"parent"`
	Depth             int         `json:"depth"`
	Key               string      `json:"key"`
	IsArrayElement    bool        `json:"isArrayElement"`
	LineNumber        int
	ClosingLineNumber int
}

// Helper functions
func getNodeType(value interface{}) NodeType {
	if value == nil {
		return NullType
	}

	switch value.(type) {
	case string:
		return StringType
	case float64, int, int64:
		return NumberType
	case bool:
		return BoolType
	case map[string]interface{}:
		return ObjectType
	case []interface{}:
		return ArrayType
	default:
		return StringType
	}
}

func getDepth(path string) int {
	if path == "" {
		return 0
	}
	re := regexp.MustCompile(`[.\[]`)
	matches := re.FindAllString(path, -1)
	return len(matches)
}

func buildChildPath(basePath, key string,
	isArray bool) string {
	if basePath == "" {
		return key
	}

	if isArray {
		return fmt.Sprintf("%s[%s]", basePath, key)
	}

	return fmt.Sprintf("%s.%s", basePath, key)
}

func isNested(value interface{}) bool {
	switch value.(type) {
	case map[string]interface{}, []interface{}:
		return true
	default:
		return false
	}
}
