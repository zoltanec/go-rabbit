package web

import (
	"context"
	"encoding/json"
	"net/http"
)

func (a *Api) Status(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	a.log.Info("Status request")

	encoded, _ := json.Marshal(struct {
		Status string
	}{Status: "ok"})
	w.Write(encoded)
}
