package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse the multipart form data
		err := r.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Retrieve the file from the form data
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Create a new file on the server
		dst, err := os.Create(handler.Filename)
		if err != nil {
			http.Error(w, "Error creating file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the file content to the new file
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Error copying file content", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "File uploaded successfully")
	} else {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
	fmt.Print("File server")
}

func main() {
	http.HandleFunc("/upload", handleUpload)

	port := 8000
	fmt.Printf("Server listening on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
