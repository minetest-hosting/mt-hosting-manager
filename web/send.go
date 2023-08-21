package web

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SendError(w http.ResponseWriter, code int, err error) {
	logrus.WithFields(logrus.Fields{
		"code":  code,
		"error": err,
	}).Error("http error")

	errmsg := &ErrorResponse{
		Success: false,
		Code:    code,
		Message: err.Error(),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	buf := bytes.NewBuffer([]byte{})
	json.NewEncoder(buf).Encode(errmsg)
	w.Write(buf.Bytes())
}

func SendText(w http.ResponseWriter, txt string) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(txt))
}

func SendJson(w http.ResponseWriter, o interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(o)
	if err != nil {
		SendError(w, 500, err)
	}
}

func Send(w http.ResponseWriter, o interface{}, err error) {
	if err != nil {
		SendError(w, 500, err)
	} else {
		SendJson(w, o)
	}
}
