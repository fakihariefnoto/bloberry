package httpx

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	Code    string `json:"code"`
	Content string `json:"content,omitempty"`
}

type Envelope struct {
	Data     interface{} `json:"data,omitempty"`
	Messages []Message   `json:"messages,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(data)
}

func OK(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusOK, Envelope{Data: data})
}

func Created(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusCreated, Envelope{Data: data})
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func MessageEnvelope(w http.ResponseWriter, status int, msgs ...Message) {
	WriteJSON(w, status, Envelope{Messages: msgs})
}

func Error(w http.ResponseWriter, status int, code string) {
	MessageEnvelope(w, status, Message{Code: code})
}

func ErrorWithContent(w http.ResponseWriter, status int, code, content string) {
	MessageEnvelope(w, status, Message{Code: code, Content: content})
}
