package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
)

const saveFile = "memo.txt" // データファイルの保存先
func main() { // メインプログラム
	// サーバーを起動する --- (*1)
	print("memo server - [URL] http://localhost:8888/\n")
	http.HandleFunc("/", readHandler) // ハンドラを登録
	http.HandleFunc("/w", writeHandler)
	http.ListenAndServe(":8888", nil) // 起動
}

// ルートへアクセスした時メモを表示 --- (*2)
func readHandler(w http.ResponseWriter, r *http.Request) {
	// データファイルを開く
	text, err := ioutil.ReadFile(saveFile)
	if err != nil {
		text = []byte("ここにメモを記入してください。")
	}
	// HTMLのフォームを返す
	htmlText := html.EscapeString(string(text))
	s := "<html>" +
		"<style>textarea { width:99%; height:200px; }</style>" +
		"<form method='POST' action='/w'>" +
		"<textarea name='text'>" + htmlText + "</textarea>" +
		"<input type='submit' value='保存' /></form></html>"
	w.Write([]byte(s))
}

// フォーム投稿した時 --- (*3)
func writeHandler(w http.ResponseWriter, r *http.Request) {
	// 投稿されたフォームを解析
	r.ParseForm()
	if len(r.Form["text"]) == 0 { // 値が書き込まれてない時
		w.Write([]byte("フォームから投稿してください。"))
		return
	}
	text := r.Form["text"][0]
	// データファイルへ書き込む
	ioutil.WriteFile(saveFile, []byte(text), 0644)
	fmt.Println("save: " + text)
	// ルートページへリダイレクトして戻る --- (*4)
	http.Redirect(w, r, "/", 301)
}
