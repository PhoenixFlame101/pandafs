package main

import (
	"fmt"
	"io"
	"net/http"
)

type Command1 struct {
	Command   string `json:"command"`
	Filename1 string `json:"filename"`
	Filename2 string `json:"filename2"`
}

func main() {
	// Define a handler function
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is POST
		if r.Method == http.MethodPost {
			// Read the request body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error parsing request body", http.StatusInternalServerError)
				return
			}

			// var reqData Command1
			// err = json.Unmarshal(body, &reqData)
			// if err != nil {
			// 	http.Error(w, "Error decoding JSON data", http.StatusBadRequest)
			// 	return
			// }
			// Print the received data
			fmt.Printf("Received data: %s\n ", string(body))
		}

		// Respond to the client
		fmt.Fprintln(w, "Hello, this is a basic Go server!")
	}

	// Register the handler function for the root route "/"
	http.HandleFunc("/", handler)

	// Start the HTTP server on port 9001
	port := 9001
	fmt.Printf("Server listening on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
