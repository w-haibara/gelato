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

func main() {
	/*
	 * Title View (genelated by http://patorjk.com/software/taag/#p=display&h=0&v=0&f=Epic&t=Gelato)
	 */
	fmt.Println(`
(  ____ \(  ____ \( \      (  ___  )\__   __/(  ___  )
| (    \/| (    \/| (      | (   ) |   ) (   | (   ) |
| |      | (__    | |      | (___) |   | |   | |   | |
| | ____ |  __)   | |      |  ___  |   | |   | |   | |
| | \_  )| (      | |      | (   ) |   | |   | |   | |
| (___) || (____/\| (____/\| )   ( |   | |   | (___) |
(_______)(_______/(_______/|/     \|   )_(   (_______)
`)

	/*
	 * Web Server
	 */
	const appDir = "./page/dist"
	const testPageDir = "./test-page"
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(appDir))))
	http.Handle("/console-test", http.StripPrefix("/console-test", http.FileServer(http.Dir(testPageDir+"/console"))))
	http.Handle("/desktop-test", http.StripPrefix("/desktop-test", http.FileServer(http.Dir(testPageDir+"/desktop"))))

	/*
	 * API Server
	 */
	http.HandleFunc("/api/status", statusAPIHandler)

	/*
	 * WebSocket
	 */
	http.HandleFunc("/console", consoleHandler)
	http.HandleFunc("/desktop", desktopHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
