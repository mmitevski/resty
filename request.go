package resty

import (
	"encoding/json"
	"net/http"

	"github.com/go-zoo/bone"
)

type paramHandler struct {
	r *http.Request
}

func (ph *paramHandler) Param(key string) string {
	return bone.GetValue(ph.r, key)
}

func (ph *paramHandler) Query(key string) string {
	q := ph.Queries(key)
	if len(q) > 0 {
		return q[0]
	}
	return ""
}

func (ph *paramHandler) Queries(key string) []string {
	return bone.GetQuery(ph.r, key)
}

func (ph *paramHandler) ScanBody(dst interface{}) error {
	return json.NewDecoder(ph.r.Body).Decode(dst)
}
