package main

import (
	"encoding/json"
	"fmt"
)


func main() {
	// Sample JSON data
	jsonData := `{
            "user": {
                    "name": "John",
                    "list": [1, 2, "three", 4],
                    "escaped": "{\"hello\": \"world\"}",
                    "addresses": [
                            {
                                    "street": "123 Main St",
                                    "zipcode": "12345"
                            },
                            {
                                    "street": "456 Oak Ave",
                                    "zipcode": "67890"
                            }
                    ]
            }
    }`

	// Parse JSON
	var data interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		panic(err)
	}

	// Build tree
	tree := BuildTree(data, "", nil)

	// Query by path
	fmt.Println("Query user.addresses:")
	fmt.Println(tree.PrintAsJSON("user.addresses", 0))
	// Output: 12345

	// Print full tree
	fmt.Println("Full tree:")
	fmt.Print(tree.PrintFromRoot())

	// Collapse addresses and print
	tree.Collapse("user.addresses")
	fmt.Println("\nWith addresses collapsed:")
	fmt.Print(tree.PrintFromRoot())

	// Expand addresses and print
	tree.Expand("user.addresses")
	fmt.Println("\nWith addresses expanded:")
	fmt.Print(tree.PrintFromRoot())

	tree.Collapse("user.list")
	tree.Collapse("user.addresses")
	fmt.Println("\nWith addresses and list collapsed:")
	fmt.Print(tree.PrintAsJSONFromRoot())

	// Expand addresses and print
	//tree.Expand("user.addresses")
	tree.Expand("user.list")
	fmt.Println("\nWith addresses expanded:")
	fmt.Print(tree.PrintAsJSONFromRoot())

    // Collapse tree
    tree.Collapse("user")
	fmt.Println("\nWith user collapsed:")
	fmt.Print(tree.PrintAsJSONFromRoot())
}
