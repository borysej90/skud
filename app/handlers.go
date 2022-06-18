package app

import (
	"context"
	"encoding/json"
	"net/http"
)

type accessChecker interface {
	CheckAccess(ctx context.Context, readerID int64, passcodeCard string) (string, bool, error)
}

func HandleCheckAccess(svc accessChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req checkAccessReq
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			// TODO: log error
			http.Error(w, "", http.StatusInternalServerError)
		}
		defer func() {
			_ = r.Body.Close()
		}()
		var resp checkAccessResp
		resp.Message, resp.Access, err = svc.CheckAccess(r.Context(), req.ReaderID, req.PassCard)
		if err != nil {
			// TODO: log error
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		body, _ := json.Marshal(resp)
		if _, err = w.Write(body); err != nil {
			// TODO: log error
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
