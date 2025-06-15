package pkg

import "net/http"

func HandleAdapter(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, 400, err)
			return
		}
	}
}
