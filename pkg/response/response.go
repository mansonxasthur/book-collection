package response

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
)

type Status int

type Response struct {
	Data interface{} `json:"data"`
	Err  string      `json:"error"`
}

func checkEncoding(data Response) error {
	w := bufio.NewWriter(nil)

	return json.NewEncoder(w).Encode(data)
}

func write(w http.ResponseWriter, data Response, status Status) {
	if err := checkEncoding(data); err != nil {
		log.Printf("error encoding response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(int(status))

	_ = json.NewEncoder(w).Encode(data)
}

func Send(w http.ResponseWriter, data interface{}, status Status) {
	response := Response{
		Data: data,
	}
	write(w, response, status)
}

func Error(w http.ResponseWriter, err error, status Status) {
	response := Response{
		Data: nil,
		Err:  err.Error(),
	}

	write(w, response, status)
}
