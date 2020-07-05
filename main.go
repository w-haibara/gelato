package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

func statusAPIHandler(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"status": map[string]interface{}{
			"arch": runtime.GOARCH,
			"os":   runtime.GOOS,
		},
	}

	j, err := json.Marshal(status)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func consoleAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(nil)
}

func main() {
	fmt.Println("gelato")

	/*
	 * Web Server
	 */
	const appDir = "./page/dist"
	const testPageDir = "./test-page"
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(appDir))))
	http.Handle("/desktop-test", http.StripPrefix("/desktop-test", http.FileServer(http.Dir(testPageDiri+"/desktop"))))

	/*
	 * API Server
	 */
	http.HandleFunc("/api/status", statusAPIHandler)
	http.HandleFunc("/api/console", consoleAPIHandler)

	/*
	 * WebSocket
	 */
	http.HandleFunc("/desktop", captureHandler)


	log.Fatal(http.ListenAndServe(":8080", nil))
}
