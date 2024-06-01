package main

import (
	"fmt"
	"net/http"
)

func enableCors(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*r).Method == "OPTIONS" {
		return
	}
}

func main() {
	fh := &FileHandler{}

	http.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w, r)
		if r.Method != "OPTIONS" {
			fh.GetFiles(w, r)
		}
	})
	http.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w, r)
		if r.Method != "OPTIONS" {
			fh.ReadFile(w, r)
		}
	})
	http.HandleFunc("/copy", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w, r)
		if r.Method != "OPTIONS" {
			fh.CopyFile(w, r)
		}
	})
	fmt.Println("Server is running on port 8088...")
	http.ListenAndServe(":8088", nil)
}
