package web

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/mux"
)

func (a *Api) LogStream(w http.ResponseWriter, r *http.Request) {
	auth_header := r.Header.Get("Authorization")
	if auth_header != a.cfg.LogStreamKey {
		SendError(w, 401, err_unauthorized)
		return
	}

	vars := mux.Vars(r)
	server_id := vars["id"]

	if a.cfg.LogStreamKey == "" || a.cfg.LogStreamDir == "" || server_id == "" {
		SendError(w, 500, errors.New("not configured"))
		return
	}

	logs := []json.RawMessage{}
	err := json.NewDecoder(r.Body).Decode(&logs)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	base_dir := path.Join(a.cfg.LogStreamDir, server_id)
	err = os.MkdirAll(base_dir, 0777)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	time_format := "2006-01"
	increment := time.Hour * 24

	now := time.Now()
	current_logs_file := path.Join(base_dir, fmt.Sprintf("%s.txt", now.Format(time_format)))

	yesterday := now.Add(increment * -1)
	previous_logs_file := path.Join(base_dir, fmt.Sprintf("%s.txt", yesterday.Format(time_format)))

	// check if we have to rotate the logs
	if previous_logs_file != current_logs_file {
		fi, _ := os.Stat(previous_logs_file)
		if fi != nil {
			// file changed and previous not gzipped yet
			gzip_file, err := os.OpenFile(previous_logs_file+".gz", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
			if err != nil {
				SendError(w, 500, err)
				return
			}
			defer gzip_file.Close()

			previous, err := os.Open(previous_logs_file)
			if err != nil {
				SendError(w, 500, err)
				return
			}
			defer previous.Close()

			gzw := gzip.NewWriter(gzip_file)
			_, err = io.Copy(gzw, previous)
			if err != nil {
				SendError(w, 500, err)
				return
			}
			gzw.Close()

			os.Remove(previous_logs_file)
		}
	}

	// append to current file
	f, err := os.OpenFile(current_logs_file, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		SendError(w, 500, fmt.Errorf("could not open current log file '%s': %v", current_logs_file, err))
		return
	}

	for _, log := range logs {
		_, err = f.Write(log)
		if err != nil {
			SendError(w, 500, fmt.Errorf("error while writing log: %v", err))
			return
		}

		_, err = f.WriteString("\n")
		if err != nil {
			SendError(w, 500, err)
			return
		}
	}

	f.Close()
}
