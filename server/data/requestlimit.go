package data

import (
	"net/http"
	"time"
)

var (
	register []string
)

func isInRegister(ip string) bool {
	for i := 0; i < len(register); i++ {
		if register[i] == ip {
			return true
		}
	}
	return false
}

func setLimit(ip string) {
	register = append(register, ip)
	time.AfterFunc(2*time.Second, func() {
		if isInRegister(ip) {
			for i := 0; i < len(register); i++ {
				if register[i] == ip {
					register[i] = ""
				}
			}
		}
	})
}

func CheckRequestLimit(ip string, w http.ResponseWriter) bool {
	if isInRegister(ip) {
		http.Error(w, "to many requests", http.StatusTooManyRequests)
		return false
	} else {
		setLimit(ip)
		return true
	}
}
