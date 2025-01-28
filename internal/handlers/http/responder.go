package http_handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// Responder отправляет JSON-ответ клиенту.
func Responder(w http.ResponseWriter, statusCode int, response interface{}) {
	const op = "http.Respond"

	w.Header().Set("Content-Type", "application/json")

	// Кодируем JSON в буфер
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(response); err != nil {
		// Если произошла ошибка, устанавливаем статус-код 500 и возвращаем сообщение об ошибке
		http.Error(w, ErrInternalServer.Error(), http.StatusInternalServerError)
		log.Println(op, ErrInvalidJSON, err)
		return
	}
	defer buf.Reset()

	// Устанавливаем статус-код и записываем данные
	w.WriteHeader(statusCode)
	_, err := w.Write(buf.Bytes())
	if err != nil {
		log.Printf("%s: %v\n", op, err)
		return
	}
}
