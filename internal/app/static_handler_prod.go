//go:build !dev

package app

import (
	"log"
	"net/http"

	"github.com/CP-Payne/wonderpicai/web"
)

func StaticFSHandler() http.Handler {

	log.Println("PROD MODE: Serving static files from embedded FS (via web.StaticFS).")
	return http.StripPrefix("/static/", http.FileServer(http.FS(web.StaticFS)))
}
