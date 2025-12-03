package api

import (
	"net/http"

	"github.com/ultrazg/xyz/service"
)

var Handler http.Handler

func init() {
	Handler = service.NewEngine()
}
