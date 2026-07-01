package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
	"golang.org/x/net/html"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// читаем файл
	data, err := os.ReadFile("index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("read index.html: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// парсим HTML-форму из файла
	doc, err := FileReadHtml("index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("parse index.html: %v", err), http.StatusInternalServerError)
		return
	}

	// находим форму в HTML-дереве
	formNode := findForm(doc)
	if formNode == nil {
		http.Error(w, "form not found in index.html", http.StatusInternalServerError)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // ограничение
	if err != nil {
		http.Error(w, fmt.Sprintf("parse multipart form: %v", err), http.StatusInternalServerError)
		return
	}

	// получаем файл из поля "myFile"
	file, header, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, fmt.Sprintf("get form file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// читаем данные из файла
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("read file: %v", err), http.StatusInternalServerError)
		return
	}

	result, err := service.Automatic_Detection(string(data))
	if err != nil {
		http.Error(w, fmt.Sprintf("automatic detection: %v", err), http.StatusInternalServerError)
		return
	}

	// сохраняем результат конвертации в локальный файл
	ext := filepath.Ext(header.Filename)
	filename := strings.ReplaceAll(time.Now().UTC().Format(time.RFC3339Nano), ":", "-") + ext
	if err := os.WriteFile(filename, []byte(result), 0644); err != nil {
		http.Error(w, fmt.Sprintf("write result file: %v", err), http.StatusInternalServerError)
		return
	}

	// возвращаем результат
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(result))
}
func FileReadHtml(filename string) (*html.Node, error) {
	// открываем файл
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %w", filename, err)
	}

	defer file.Close()
	// парсим
	doc, err := html.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("parse HTML: %w", err)
	}
	return doc, nil
}

func findForm(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "form" {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findForm(c); result != nil {
			return result
		}
	}
	return nil
}
