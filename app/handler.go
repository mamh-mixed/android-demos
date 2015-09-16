package app

import (
	"encoding/json"
	"net/http"

	"github.com/omigo/log"
)

func registerHandle(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("username")
	password := r.FormValue("password")
	transtime := r.FormValue("transtime")
	sign := r.FormValue("sign")

	ret := User.register(userName, password, transtime, sign)

	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
		return
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}
