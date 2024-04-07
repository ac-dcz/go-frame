package geeweb

import "net/http"

type Context struct {
	w http.ResponseWriter
	r *http.Request
}
