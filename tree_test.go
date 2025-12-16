package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildTree_EmptyObject(t *testing.T) {
	data := map[string]interface{}{}
	tree := BuildTree(data, "", nil)

	// Should have root node only
	assert.NotNil(t, tree)
	assert.Equal(t, 1, len(tree.Nodes)) // Just root
}

func TestBuildTree_Simple(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		path     string
		expected interface{}
		nodeType NodeType
	}{
		{
			"simple string",
			map[string]interface{}{"name": "John"},
			"name",
			"John",
			StringType,
		},
		{
			"simple number",
			map[string]interface{}{"age": 25.0},
			"age",
			25.0,
			NumberType,
		},
		{
			"simple bool",
			map[string]interface{}{"active": true},
			"active",
			true,
			BoolType,
		},
		{
			"simple null",
			map[string]interface{}{"value": nil},
			"value",
			nil,
			NullType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := BuildTree(tt.data, "", nil)
			assert.Equal(t, tt.expected, tree.GetValue(tt.path))
			assert.Equal(t, tt.nodeType, tree.Nodes[tt.path].Type)
		})
	}
}

func TestBuildTree_SimpleArray(t *testing.T) {
	data := map[string]interface{}{
		"elements": []interface{}{
			1, 2.5, "three", false, nil,
		},
	}

	tree := BuildTree(data, "", nil)

	assert.Equal(t, 1, tree.GetValue("elements[0]"))
	assert.Equal(t, 2.5, tree.GetValue("elements[1]"))
	assert.Equal(t, "three", tree.GetValue("elements[2]"))
	assert.Equal(t, NumberType, tree.Nodes["elements[0]"].Type)
	assert.Equal(t, false, tree.GetValue("elements[3]"))
	assert.Equal(t, nil, tree.GetValue("elements[4]"))
	assert.Equal(t, NullType, tree.Nodes["elements[4]"].Type)
}

func TestBuildTree_NestedObject(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  30.0,
		},
	}
	tree := BuildTree(data, "", nil)

	assert.Equal(t, "John", tree.GetValue("user.name"))
	assert.Equal(t, 30.0, tree.GetValue("user.age"))
	assert.Equal(t, "user", tree.Nodes["user.name"].Parent)
}

func TestBuildTree(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  30.0,
		},
		"friends": []interface{}{1, 2, 3, 4},
		"identifications": []interface{}{
			map[string]interface{}{
				"type":   "passport",
				"number": "123456789",
			},
			map[string]interface{}{
				"type":   "license",
				"number": "987654321",
			},
		},
		"email":   "john@email.com",
		"escaped": "{\"meta\": \"data\"}",
		"active":  true,
	}
	tree := BuildTree(data, "", nil)

	assert.Equal(t, "John", tree.GetValue("user.name"))
	assert.Equal(t, 30.0, tree.GetValue("user.age"))
	assert.Equal(t, "user", tree.Nodes["user.name"].Parent)

	assert.Equal(t, 1, tree.GetValue("friends[0]"))
	assert.Equal(t, "friends", tree.Nodes["friends[0]"].Parent)

	assert.Equal(t, "passport",
		tree.GetValue("identifications[0].type"))
	assert.Equal(t, "987654321",
		tree.GetValue("identifications[1].number"))

	assert.Equal(t, "john@email.com", tree.GetValue("email"))
	assert.Equal(t, "{\"meta\": \"data\"}",
		tree.GetValue("escaped"))
	assert.Equal(t, BoolType, tree.Nodes["active"].Type)
}

// TODO (isaac): rewrite with the new rendering system
func TestPrintAsJSON_FullTree(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  30.0,
		},
		// Unmarshal defaults to float64, so a cast is needed
		// for the test
		"friends": []interface{}{
			float64(1), 2.5, "three", true, nil,
		},
		"identifications": []interface{}{
			map[string]interface{}{
				"type":   "passport",
				"number": "123456789",
			},
			map[string]interface{}{
				"type":   "license",
				"number": "987654321",
			},
		},
		"email":   "john@email.com",
		"escaped": "{\"meta\": \"data\"}",
		"active":  true,
	}

	tree := BuildTree(data, "", nil)

	t.Run("print full tree", func(t *testing.T) {
        currentTheme = themes["nocolor"]
		result := tree.PrintAsJSONFromRoot()

		// Parse both, expected and actual JSON
		var actual interface{}
		err := json.Unmarshal([]byte(result), &actual)
		assert.NoError(t, err, "Generated JSON should be valid")

		// expected, actual
		assert.Equal(t, data, actual)
	})
}

func TestPrintAsJSON_FullTree2(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  30.0,
		},
		// Unmarshal defaults to float64, so a cast is needed
		// for the test
		"friends": []interface{}{
			float64(1), 2.5, "three", true, nil, "\n\"hello\"",
		},
		"identifications": []interface{}{
			map[string]interface{}{
				"type":   "passport",
				"number": "123456789",
			},
			map[string]interface{}{
				"type":   "license",
				"number": "987654321",
			},
		},
		"email":   "john@email.com",
		"escaped": "{\"meta\": \"data\"}",
		"active":  true,
	}

	tree := BuildTree(data, "", nil)

	t.Run("print full tree", func(t *testing.T) {
        currentTheme = themes["nocolor"]
		lines := tree.PrintAsJSON2()
        result := ""
        for _, line := range lines {
            result += RenderLine(line, false)
            result += "\n"
        }
        result = strings.TrimSuffix(result, "\n")

		// Parse both, expected and actual JSON
		var actual interface{}
		err := json.Unmarshal([]byte(result), &actual)
		assert.NoError(t, err, "Generated JSON should be valid")

		// expected, actual
		assert.Equal(t, data, actual)
	})
}

// TODO (isaac): rewrite with new rendering system
func TestPrintAsJSON_CollapsedObject(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  30.0,
		},
	}

	tree := BuildTree(data, "", nil)
	tree.Collapse("user")
    currentTheme = themes["nocolor"]

	t.Run("print tree with collapsed object", func(t *testing.T) {
		expected := "{\n  \"user\": {...} // 2 properties\n}"
		assert.Equal(t, expected, tree.PrintAsJSONFromRoot())
	})
}

func TestPrintAsJSON_CollapsedObject2(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  30.0,
		},
	}

	tree := BuildTree(data, "", nil)
	tree.Collapse("user")
    currentTheme = themes["nocolor"]

	t.Run("print tree with collapsed object", func(t *testing.T) {
        lines := tree.PrintAsJSON2()
        actual := ""
        for _, line := range lines {
            actual += RenderLine(line, false)
            actual += "\n"
        }
        actual = strings.TrimSuffix(actual, "\n")

		expected := "{\n  \"user\": {...}\n}"
		assert.Equal(t, expected, actual)
	})
}

// TODO (isaac): rewrite with new rendering system
func TestPrintAsJSON_CollapsedArray(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"emails": []interface{}{
				"user@mail.com", "user@mail.org",
			},
		},
	}

	tree := BuildTree(data, "", nil)
	tree.Collapse("user.emails")

	t.Run("print tree with collapsed object", func(t *testing.T) {
		expected := "{\n  \"user\": {\n    \"emails\": [...] // 2 items\n  }\n}"
		assert.Equal(t, expected, tree.PrintAsJSONFromRoot())
	})
}

func TestPrintAsJSON_CollapsedArray2(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"emails": []interface{}{
				"user@mail.com", "user@mail.org",
			},
		},
	}

	tree := BuildTree(data, "", nil)
	tree.Collapse("user.emails")

	t.Run("print tree with collapsed object", func(t *testing.T) {
        lines := tree.PrintAsJSON2()
        actual := ""
        for _, line := range lines {
            actual += RenderLine(line, false)
            actual += "\n"
        }
        actual = strings.TrimSuffix(actual, "\n")

		expected := "{\n  \"user\": {\n    \"emails\": [...]\n  }\n}"
		assert.Equal(t, expected, actual)
	})
}

func TestGetNode(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  30.0,
		},
	}

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			"get object node",
			"user",
			"user",
		},
		{
			"get nested node",
			"user.age",
			"age",
		},
	}

	tree := BuildTree(data, "", nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, exists := tree.GetNode(tt.path)
			assert.True(t, exists, "Node should exist")
			assert.Equal(t, tt.expected, node.Key)
		})
	}
}

func TestGetNodeAtLine(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  30.0,
		},
		"email": "john@example.com",
		"identifications": []interface{}{
			map[string]interface{}{
				"type":   "passport",
				"number": "123456789",
			},
			map[string]interface{}{
				"type":   "license",
				"number": "987654321",
			},
		},
		"active": true,
	}

	tree := BuildTree(data, "", nil)

	// Use the non JSON print version, because
	// the line numbers follow the same pattern
	// as the JSON print
	printedTree := tree.PrintFromRoot()
	lines := strings.Split(printedTree, "\n")

	for i, line := range lines {
		t.Run(fmt.Sprintf("test line %d", i), func(t *testing.T) {
			// Get the node at each line
			node, exists := tree.GetNodeAtLine(i)
			if exists {
				// If there is a node, check if the printed
				// line contains the node key
				assert.True(t, strings.Contains(line, node.Key))
			}
		})
	}
}

// TODO (isaac): rewrite with new rendering system
func TestGetNodeAtLine_CompressedNode(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  30.0,
		},
		"email": "john@example.com",
		"identifications": []interface{}{
			map[string]interface{}{
				"type":   "passport",
				"number": "123456789",
			},
			map[string]interface{}{
				"type":   "license",
				"number": "987654321",
			},
		},
		"active": true,
	}

    tree := BuildTree(data, "", nil)
    tree.Collapse("identifications[0]")

    // Print to update VirtualToRealLines
    _ = tree.PrintAsJSONFromRoot()

    expectedNode, exists := tree.GetNode("identifications[1]")
    assert.True(t, exists, "Expected node should exist")
    found := 0

    if exists {
        for _, realLineNumber := range(tree.VirtualToRealLines) {
            actualNode, exists := tree.GetNodeAtLine(realLineNumber)
            if exists {
                if actualNode.Path == expectedNode.Path {
                    found++
                }
            }
        }
    }
    assert.Equal(t, 1, found,
        "Expected node not properly found using VirtualToRealLines")
}

func TestGetNodeAtLine_CompressedNode2(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  30.0,
		},
		"email": "john@example.com",
		"identifications": []interface{}{
			map[string]interface{}{
				"type":   "passport",
				"number": "123456789",
			},
			map[string]interface{}{
				"type":   "license",
				"number": "987654321",
			},
		},
		"active": true,
	}

    tree := BuildTree(data, "", nil)
    tree.Collapse("identifications[0]")

    // Print to update VirtualToRealLines
    _ = tree.PrintAsJSON2()

    expectedNode, exists := tree.GetNode("identifications[1]")
    assert.True(t, exists, "Expected node should exist")
    found := 0

    if exists {
        for _, realLineNumber := range(tree.VirtualToRealLines) {
            actualNode, exists := tree.GetNodeAtLine(realLineNumber)
            if exists {
                if actualNode.Path == expectedNode.Path {
                    found++
                }
            }
        }
    }
    assert.Equal(t, 1, found,
        "Expected node not properly found using VirtualToRealLines")
}
