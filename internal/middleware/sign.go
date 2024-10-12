package middleware

import (
	"bytes"
	"encoding/hex"
	"io"
	"log"
	"net/http"

	"github.com/bbquite/mca-server/internal/utils"
)

func CheckSignMiddleware(shaKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			bodyBytes, _ := io.ReadAll(r.Body)
			r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			log.Printf(`
				MW sign ckeck
				key: %s
				sign: %s
			`, shaKey, r.Header.Get("Hashsha256"))

			shaHeaderSign, err := hex.DecodeString(r.Header.Get("Hashsha256"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Print(err)
				return
			}
			if utils.CheckHMACEqual(shaKey, shaHeaderSign, bodyBytes) {
				log.Print("NORM")
				next.ServeHTTP(w, r)
			} else {
				log.Print("BAD")
				//w.WriteHeader(http.StatusBadRequest)
				//return
			}
			next.ServeHTTP(w, r)
		})
	}
}
