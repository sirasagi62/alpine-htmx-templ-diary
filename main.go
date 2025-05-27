package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

// example関数はHTML文字列を返す
func example() string {
	return "<p>example</p>"
}

type DiaryEntry struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Emoji   string `json:"emoji"`
	Date    string `json:"date"`
}

// 日記のエントリーのスライスを表す構造体
type DiaryEntries []DiaryEntry

// JSONファイルを読み込み、DiaryEntriesを返す関数
func LoadDiaryEntries(filePath string) (DiaryEntries, error) {
	// JSONファイルを読み込む
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// JSONデータをデコードする
	var entries DiaryEntries
	decoder := json.NewDecoder(bytes.NewReader(data))
	err = decoder.Decode(&entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
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

	http.HandleFunc("/view/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		diaries, err := LoadDiaryEntries("diary.json")
		if err != nil {
			http.Error(w, "Can't load diary.json", http.StatusInternalServerError)
			return
		}
		diaryCardList(diaries).Render(r.Context(), w)
	})

	http.HandleFunc("/view/item", func(w http.ResponseWriter, r *http.Request) {
		idx_str := r.URL.Query().Get("idx")
		if idx_str == "" {
			idx_str = "0"
		}
		idx, _ := strconv.Atoi(idx_str)

		diaries, err := LoadDiaryEntries("diary.json")
		if err != nil {
			http.Error(w, "Can't load diary.json", http.StatusInternalServerError)
			return
		}

		if idx >= len(diaries) {
			http.Error(w, "The entry is not found.", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		diaryDetail(diaries[idx]).Render(r.Context(), w)
	})

	// "/"および"/index.html"に対応
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			data, err := os.ReadFile("./index.html")
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
