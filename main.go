package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// example関数はHTML文字列を返す
func example() string {
	return "<p>example</p>"
}

func main() {

	// "/api/example"に対応
	http.HandleFunc("/api/example", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, example())
	})

	// "/api/john"に対応
	http.HandleFunc("/api/john", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		hello("John").Render(r.Context(), w)
	})
	// "/"および"/index.html"に対応
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			data, err := os.ReadFile("./index.html") // ioutil.ReadFile → os.ReadFile に変更
			if err != nil {
				http.Error(w, "index.html not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
			return
		}
		// その他のパスは404
		http.NotFound(w, r)
	})

	// サーバー起動
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
