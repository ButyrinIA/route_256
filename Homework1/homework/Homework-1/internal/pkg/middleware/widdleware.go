package middleware

import (
	"log"
	"net/http"
)

// логгирования полей POST, PUT, DELETE запросов

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Method: %s, URL: %s\n", req.Method, req.URL.Path)

		if req.Method == http.MethodPost || req.Method == http.MethodPut || req.Method == http.MethodDelete {
			// Логгирование полей POST, PUT, DELETE запросов
			log.Println("Request Body:")
			if err := req.ParseForm(); err != nil {
				log.Println("Error parsing form:", err)
			}
			for key, values := range req.PostForm {
				for _, value := range values {
					log.Printf("%s: %s\n", key, value)
				}
			}
		}

		next.ServeHTTP(w, req)
	})
}
