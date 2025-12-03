package api

import (
	"net/http"
	"sync"

	"github.com/ultrazg/xyz/service"
)

var (
	engineOnce sync.Once
	engine     http.Handler
)

func getHandler() http.Handler {
	engineOnce.Do(func() {
		engine = service.NewEngine()
	})
	return engine
}

func Handler(w http.ResponseWriter, r *http.Request) {
	getHandler().ServeHTTP(w, r)
}
