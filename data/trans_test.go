package data

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func xTestSyncMerchant(t *testing.T) {
	url := "http://127.0.0.1:6800/import?key=cilxl123$&type=merchant"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	w := httptest.NewRecorder()

	Import(w, req)

	if w.Code != 200 {
		t.Error(w.Body)
	}
}
