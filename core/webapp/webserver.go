package webapp

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	indexHTML = `<!doctype html>
				 <html>
				 <head>
					<meta charset="UTF-8">
					<title>DNS Exporter</title>
				 </head>
				 <body>
					<h1>DNS Exporter</h1>
					<p><a href="%s">Metrics</a></p>
				 </body>
				 </html>`
)

func StartServer(listenAddress string, port int, metricsPath string) {
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, indexHTML, metricsPath)
	})
	http.Handle(metricsPath, promhttp.Handler())

	addr := fmt.Sprintf("%s:%d", listenAddress, port)
	server := &http.Server{
		Addr:    addr,
		Handler: logRequest(http.DefaultServeMux),
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()
	log.Printf("Server started on %s, use <Ctrl-C> to stop\n", addr)

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Printf("\nReceived termination signal. Shutting down gracefully...")

	// Optionally, you can add cleanup logic here before exiting.
	os.Exit(0)
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddress := r.RemoteAddr
		log.Printf("Request received from %s: %s %s\n", ipAddress, r.Method, r.URL.Path)
		handler.ServeHTTP(w, r)
	})
}
