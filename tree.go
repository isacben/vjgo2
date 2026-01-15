package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type JSONTree struct {
	Nodes              map[string]*Node `json:"nodes"`
	LineNumbers        map[int]*Node    `json:"lineNumbers"`
	VirtualToRealLines []int
	Children           map[string][]string `json:"children"`
	Collapsed          map[string]bool     `json:"collapsed"`
	lineCounter        int
	currentRealLine    int
}

func NewJSONTree() *JSONTree {
	return &JSONTree{
		Nodes:       make(map[string]*Node),
		LineNumbers: make(map[int]*Node),
		Children:    make(map[string][]string),
		Collapsed:   make(map[string]bool),
		lineCounter: 0,
	}
}

// ========== Core methods ==========

// GetValue returns the value at the given path
func (jt *JSONTree) GetValue(path string) interface{} {
	if node, exists := jt.Nodes[path]; exists {
		return node.Value
	}
	return nil
}

// Collapse marks a path as collapsed
func (jt *JSONTree) Collapse(path string) {
	jt.Collapsed[path] = true
}

// Expand marks a path as expanded
func (jt *JSONTree) Expand(path string) {
	delete(jt.Collapsed, path)
}

// IsCollapsed checks if a path is collapsed
func (jt *JSONTree) IsCollapsed(path string) bool {
	return jt.Collapsed[path]
}

// AddChild adds a child path to a parent
func (jt *JSONTree) AddChild(parent string, child string) {
	if jt.Children[parent] == nil {
		// Create empty slice to list the children nodes
		jt.Children[parent] = make([]string, 0)
	}
	jt.Children[parent] = append(jt.Children[parent], child)
}

// ========== Utility Methods ==========

// GetNode returns the node at the given path
func (jt *JSONTree) GetNode(path string) (*Node, bool) {
	node, exists := jt.Nodes[path]
	return node, exists
}

// GetNodeAtLine returns the node at a given line number
func (jt *JSONTree) GetNodeAtLine(lineNum int) (*Node, bool) {
	node, exists := jt.LineNumbers[lineNum]
	return node, exists
}

// GetChildren returns all child paths for a given path
func (jt *JSONTree) GetChildren(path string) []string {
	return jt.Children[path]
}

// HasChildren checks if a path has children
func (jt *JSONTree) HasChildren(path string) bool {
	return len(jt.Children[path]) > 0
}

// GetAllPaths returns all paths in the tree
func (jt *JSONTree) GetAllPaths() []string {
	paths := make([]string, 0, len(jt.Nodes))
	for path := range jt.Nodes {
		paths = append(paths, path)
	}
	return paths
}

// SetValue updates the value at a given path
func (jt *JSONTree) SetValue(path string, value interface{}) bool {
	if node, exists := jt.Nodes[path]; exists {
		node.Value = value
		node.Type = getNodeType(value)
		return true
	}
	return false
}

// ========== Pretty Printing ==========

// Print returns a formatted string representation
func (jt *JSONTree) Print(startPath string, indent int) string {
	node, exists := jt.Nodes[startPath]
	if !exists {
		return ""
	}

	result := strings.Repeat(" ", indent) + node.Key + ": "

	if jt.IsCollapsed(startPath) {
		childCount := len(jt.Children[startPath])
		result += fmt.Sprintf("{...} // %d items\n", childCount)
		return result
	}

	switch node.Type {
	case ObjectType:
		result += "{\n"
		for _, childPath := range jt.Children[startPath] {
			result += jt.Print(childPath, indent+2)
		}
		result += strings.Repeat(" ", indent) + "}\n"

	case ArrayType:
		result += "[\n"
		for _, childPath := range jt.Children[startPath] {
			result += jt.Print(childPath, indent+2)
		}
		result += strings.Repeat(" ", indent) + "]\n"

	default:
		valueBytes, _ := json.Marshal(node.Value)
		result += string(valueBytes) + "\n"
	}

	return result
}

// PrintFromRoot prints the entire tree
func (jt *JSONTree) PrintFromRoot() string {
	return jt.Print("", 0)
}

// PrintAsJSON returns the tree as properly formatted JSON
func (jt *JSONTree) PrintAsJSON(startPath string, indent int) string {
	node, exists := jt.Nodes[startPath]
	if !exists {
		// Handle root case when no explicit root node exists
		if startPath == "" {
			// Find if root has single object child or multiple children
			children := jt.Children[startPath]
			if len(children) == 1 {
				// Single root object/array
				return jt.PrintAsJSON(children[0], indent)
			} else if len(children) > 1 {
				// Multiple root elements - wrap in object
				result := "{\n"
				for i, childPath := range children {
					if i > 0 {
						result += ",\n"
					}
					childNode := jt.Nodes[childPath]
					result += strings.Repeat("  ", indent+1) + `"` +
						keyStyle.Render(childNode.Key) + `": `
					result += strings.TrimSpace(jt.PrintAsJSON(childPath, indent+1))
				}
				result += "\n" + strings.Repeat("  ", indent) + "}"
				jt.currentRealLine++
				jt.VirtualToRealLines = append(jt.VirtualToRealLines, jt.currentRealLine)
				return result
			}
		}
		return ""
	}

	jt.currentRealLine = node.LineNumber
	jt.VirtualToRealLines = append(jt.VirtualToRealLines, jt.currentRealLine)

	if jt.IsCollapsed(startPath) {
		childCount := len(jt.Children[startPath])
		if jt.Nodes[startPath].Type == ArrayType {
			return fmt.Sprintf("[...] // %d items", childCount)
		}
		return fmt.Sprintf("{...} // %d properties", childCount)
	}

	switch node.Type {
	case ObjectType:
		children := jt.Children[startPath]
		if len(children) == 0 {
			return "{}"
		}

		result := "{\n"
		for i, childPath := range children {
			if i > 0 {
				result += ",\n"
			}
			childNode := jt.Nodes[childPath]
			// Quote the key and add colon
			result += strings.Repeat("  ", indent+1) + `"` +
				keyStyle.Render(childNode.Key) + `": `
			result += strings.TrimSpace(jt.PrintAsJSON(childPath, indent+1))
		}
		result += "\n" + strings.Repeat("  ", indent) + "}"
		jt.currentRealLine++
		jt.VirtualToRealLines = append(jt.VirtualToRealLines, jt.currentRealLine)
		return result

	case ArrayType:
		children := jt.Children[startPath]
		if len(children) == 0 {
			return "[]"
		}

		result := "[\n"
		for i, childPath := range children {
			if i > 0 {
				result += ",\n"
			}
			result += strings.Repeat("  ", indent+1)
			result += strings.TrimSpace(jt.PrintAsJSON(childPath, indent+1))
		}
		result += "\n" + strings.Repeat("  ", indent) + "]"
		jt.currentRealLine++
		jt.VirtualToRealLines = append(jt.VirtualToRealLines, jt.currentRealLine)
		return result

	case StringType:
		return stringStyle.Render(
			`"` + strings.ReplaceAll(node.Value.(string), `"`, `\"`) + `"`)

	case NumberType:
		return numberStyle.Render(fmt.Sprintf("%v", node.Value))

	case BoolType:
		return booleanStyle.Render(fmt.Sprintf("%t", node.Value.(bool)))

	case NullType:
		return nullStyle.Render("null")

	default:
		// Fallback to JSON marshal
		valueBytes, _ := json.Marshal(node.Value)
		return string(valueBytes)
	}
}

// PrintAsJSONFromRoot prints the entire tree as JSON
func (jt *JSONTree) PrintAsJSONFromRoot() string {
	jt.currentRealLine = 0
	jt.VirtualToRealLines = jt.VirtualToRealLines[:0]
	return jt.PrintAsJSON("", 0)
}

func (jt *JSONTree) PrintAsJSON2() []LineMetadata {
	var result []LineMetadata
	jt.currentRealLine = 0
	jt.VirtualToRealLines = jt.VirtualToRealLines[:0]
	jt.collectLines("", 0, &result, true, true)
	return result
}

func (jt *JSONTree) collectLines(startPath string, indent int, result *[]LineMetadata, isRoot bool, isLast bool) {
	node, exists := jt.Nodes[startPath]
	if !exists {
		// Handle root case
		if startPath == "" {
			children := jt.Children[startPath]
			for i, childPath := range children {
				isLastChild := i == len(children)-1
				jt.collectLines(childPath, indent, result, false, isLastChild)
			}
		}
		return
	}

	switch node.Type {
	case ObjectType:
		// Add opening brace
		if isRoot {
			openBrace := LineMetadata{
				LineNumber:  len(*result),
				LineType:    OpenBracket,
				Content:     "{",
				NodePath:    startPath,
				NodeType:    node.Type,
				Indent:      indent,
				BracketChar: "{",
				IsCollapsed: jt.IsCollapsed(startPath),
				HasChildren: jt.HasChildren(startPath),
			}
			*result = append(*result, openBrace)
			jt.VirtualToRealLines = append(jt.VirtualToRealLines, node.LineNumber)
		}

		// Add key line if this isn't root
		if !isRoot && node.Key != "" {
			keyLine := LineMetadata{
				LineNumber:     len(*result),
				LineType:       ContentWithBrace,
				Content:        node.Key,
				NodePath:       startPath,
				NodeType:       node.Type,
				Key:            node.Key,
				Value:          node.Value,
				IsArrayElement: node.IsArrayElement,
				Indent:         indent,
				BracketChar:    "{",
				IsCollapsed:    jt.IsCollapsed(startPath),
				HasChildren:    jt.HasChildren(startPath),
				IsLastChild:    isLast,
			}
			*result = append(*result, keyLine)
			jt.VirtualToRealLines = append(jt.VirtualToRealLines, node.LineNumber)
		}

		// Add children if not collapsed
		if !jt.IsCollapsed(startPath) {
			children := jt.Children[startPath]
			for i, childPath := range children {
				isLastChild := i == len(children)-1
				jt.collectLines(childPath, indent+1, result, false, isLastChild)
			}

			// Add closing brace
			closeBrace := LineMetadata{
				LineNumber:  len(*result),
				LineType:    CloseBracket,
				Content:     "}",
				NodePath:    startPath,
				NodeType:    node.Type,
				Indent:      indent,
				BracketChar: "}",
				IsLastChild: isLast,
			}
			*result = append(*result, closeBrace)
			jt.VirtualToRealLines = append(jt.VirtualToRealLines, node.ClosingLineNumber)
		}

	case ArrayType:
		// Add opening brace if it's root
		if isRoot {
			openBrace := LineMetadata{
				LineNumber:  len(*result),
				LineType:    OpenBracket,
				Content:     "[",
				NodePath:    startPath,
				NodeType:    node.Type,
				Indent:      indent,
				BracketChar: "[",
				IsCollapsed: jt.IsCollapsed(startPath),
				HasChildren: jt.HasChildren(startPath),
			}
			*result = append(*result, openBrace)
			jt.VirtualToRealLines = append(jt.VirtualToRealLines, node.LineNumber)
		}

		// Add key line if this isn't root
		if !isRoot && node.Key != "" {
			keyLine := LineMetadata{
				LineNumber:  len(*result),
				LineType:    ContentWithBrace,
				Content:     node.Key,
				NodePath:    startPath,
				NodeType:    node.Type,
				Key:         node.Key,
				Value:       node.Value,
				Indent:      indent,
				BracketChar: "[",
				IsCollapsed: jt.IsCollapsed(startPath),
				HasChildren: jt.HasChildren(startPath),
				IsLastChild: isLast,
			}
			*result = append(*result, keyLine)
			jt.VirtualToRealLines = append(jt.VirtualToRealLines, node.LineNumber)
		}

		// Add opening bracket
		// openBracket := LineMetadata{
		// 	LineNumber:  len(*result),
		// 	LineType:    OpenBracket,
		// 	Content:     "[",
		// 	NodePath:    startPath,
		// 	NodeType:    node.Type,
		// 	Indent:      indent,
		// 	BracketChar: "[",
		// 	IsCollapsed: jt.IsCollapsed(startPath),
		// 	HasChildren: jt.HasChildren(startPath),
		// }
		// *result = append(*result, openBracket)

		// Add children if not collapsed
		if !jt.IsCollapsed(startPath) {
			children := jt.Children[startPath]
			for i, childPath := range children {
				isLastChild := i == len(children)-1
				jt.collectLines(childPath, indent+1, result, false, isLastChild)
			}

			// Add closing bracket
			closeBracket := LineMetadata{
				LineNumber:  len(*result),
				LineType:    CloseBracket,
				Content:     "]",
				NodePath:    startPath,
				NodeType:    node.Type,
				Indent:      indent,
				BracketChar: "]",
				IsLastChild: isLast,
			}
			*result = append(*result, closeBracket)
			jt.VirtualToRealLines = append(jt.VirtualToRealLines, node.ClosingLineNumber)
		}

	default:
		// Primitive values (string, number, boolean, null)
		valueLine := LineMetadata{
			LineNumber: len(*result),
			LineType:   ContentLine,
			// Content:        fmt.Sprintf("%v", node.Value),
			NodePath:       startPath,
			NodeType:       node.Type,
			Key:            node.Key,
			Value:          node.Value,
			Indent:         indent,
			IsArrayElement: node.IsArrayElement,
			IsLastChild:    isLast,
		}

		if node.Type == StringType {
			escapedBytes, err := json.Marshal(node.Value)
			if err != nil {
				panic(err)
			}

			// Remove outer quotes
			valueLine.Content = string(escapedBytes[1 : len(escapedBytes)-1])
		} else {
			valueLine.Content = fmt.Sprintf("%v", node.Value)
		}

		*result = append(*result, valueLine)
		jt.VirtualToRealLines = append(jt.VirtualToRealLines, node.LineNumber)
	}
}

// ========== Tree Building ==========

// BuildTree constructs the tree from JSON data
func BuildTree(data interface{}, basePath string, tree *JSONTree) *JSONTree {
	if tree == nil {
		tree = NewJSONTree()
	}

	// Create root node if this is the initial call
	if basePath == "" {
		rootNode := &Node{
			Path:           "",
			Type:           getNodeType(data),
			Value:          data,
			Parent:         "",
			Depth:          0,
			Key:            "root", // or you could use ""
			IsArrayElement: false,
			LineNumber:     0,
		}
		tree.Nodes[""] = rootNode
		tree.LineNumbers[rootNode.LineNumber] = rootNode
		tree.lineCounter++
	}

	createNode := func(path string, value interface{},
		key string, isParentArray bool) *Node {

		nodeType := getNodeType(value)
		depth := getDepth(path)
		isArrayElement := regexp.MustCompile(`\[\d+\]$`).MatchString(path) ||
			isParentArray

		node := &Node{
			Path:           path,
			Type:           nodeType,
			Value:          value,
			Parent:         basePath,
			Depth:          depth,
			Key:            key,
			IsArrayElement: isArrayElement,
			LineNumber:     tree.lineCounter,
		}
		tree.lineCounter++
		return node
	}

	// data.(type) syntax is specific to the switch statements
	// it can be used alone, or with variable assignment, like
	// in this case, where "v" gets the actual map
	switch v := data.(type) {
	case map[string]interface{}:
		// map[string]interface{} if for JSON objects
		for key, value := range v {
			childPath := buildChildPath(basePath, key, false)
			node := createNode(childPath, value, key, false)
			tree.Nodes[childPath] = node
			tree.LineNumbers[node.LineNumber] = node
			tree.AddChild(basePath, childPath)

			// Recursively build for nested objects/arrays
			if isNested(value) {
				BuildTree(value, childPath, tree)
			}
		}

		if node, exists := tree.Nodes[basePath]; exists {
			node.ClosingLineNumber = tree.lineCounter
		}

		tree.lineCounter++ // count the "}"

	case []interface{}:
		// []interface{} is for JSON arrays
		for i, value := range v {
			key := strconv.Itoa(i)
			childPath := buildChildPath(basePath, key, true)
			node := createNode(childPath, value, fmt.Sprintf("[%d]", i), true)
			tree.Nodes[childPath] = node
			tree.LineNumbers[node.LineNumber] = node
			tree.AddChild(basePath, childPath)

			// Recursively build for nested objects/arrays
			if isNested(value) {
				BuildTree(value, childPath, tree)
			}
		}

		if node, exists := tree.Nodes[basePath]; exists {
			node.ClosingLineNumber = tree.lineCounter
		}

		tree.lineCounter++ // count the "]"
	}

	return tree
}
