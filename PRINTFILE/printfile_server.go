package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// サーバーを起動 --- (*1)
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":8888", nil)
}

// アクセスがあった時、アップロードフォームを返す --- (*2)
func indexHandler(w http.ResponseWriter, r *http.Request) {
	s := "<html><body>" +
		"<h1>ファイルを指定してください</h1>" +
		"<form action='/upload' method='post' " +
		" enctype='multipart/form-data'>" +
		"<input type='file' name='upfile'>" +
		"<input type='submit' value='アップロード'>" +
		"</form></body></html>"
	w.Write([]byte(s))
}

// ファイルを投稿した時 --- (*3)
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("upfile") // ファイルを取得
	if err != nil {
		w.Write([]byte("アップロードエラー"))
		return
	}
	data, err := ioutil.ReadAll(file) // ファイルを読み出す
	if err != nil {
		w.Write([]byte("アップロードエラー"))
		return
	}
	s := getBinStr(data) // データをバイナリ表示
	w.Write([]byte("<html><body>" + s +
		"</body></html>"))
}

// バイナリ表示 --- (*4)
func getBinStr(bytes []byte) string {
	// 繰り返し表示する
	result := "<style>" +
		"th { background-color: #f0f0f0; } " +
		".c { background-color: #fff0f0; } " +
		"td { border-bottom: 1px solid silver } " +
		"</style><table>"
	line := "<tr>"
	aline := ""
	cnt := 2
	for i := 0; i < len(bytes); i++ {
		b := bytes[i]
		// アスキー文字の表示用
		c := string(b)
		if b < 32 || b > 126 {
			c = "_"
		}
		if c == ">" {
			c = "&gt;"
		}
		if c == "<" {
			c = "&lt;"
		}
		aline += c
		// 画面に表示する
		m := i % 16
		if m == 0 { // アドレスを追加
			line += fmt.Sprintf("<th>%04d:</th><td>", i)
		}
		line += fmt.Sprintf("%02x", b)
		switch m {
		case 3, 7, 11: // 見やすく区切り線
			line += "</td><td>"
			cnt -= 1
		case 15: // 区切り線とアスキー文字
			result += line + "</td>"
			result += "<td class='c'>" + aline + "</td></tr>\n"
			line = "<tr>"
			aline = ""
			cnt = 2
		default:
			line += " "
		}
	}
	// 表示残しを出力
	if line != "" {
		result += line
		for j := 0; j < cnt; j++ {
			result += "</td><td>"
		}
		result += "</td><td class='c'>" + aline + "</td></tr>"
	}
	result += "</table>"
	return result
}
