package utils

import (
	"math/rand"
	"net"
	"net/http"
	"strings"
)

var newLine = "\n"

//RandSeq ...
func RandSeq(n int) string {

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	lLength := 62 // len(letters)
	result := make([]rune, n)
	for i := range result {
		result[i] = letters[rand.Intn(lLength)]
	}
	return string(result)
}

//CheckIP ....
func CheckIP(allowIPs []string, r *http.Request) bool {
	clientIP := GetRealAddr(r)
	clientIPSplitted := Split(clientIP, ".")
	for _, ip := range allowIPs {
		k := 0
		splittedIP := Split(ip, ".")
		for i, item := range splittedIP {
			if item == clientIPSplitted[i] || item == "*" {
				k++
			}
		}
		if k == len(splittedIP) {
			return true
		}
	}
	return false
}

//GetRealAddr ...
func GetRealAddr(r *http.Request) string {

	remoteIP := ""
	// the default is the originating ip. but we try to find better options because this is almost
	// never the right IP
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) == 2 {
		remoteIP = parts[0]
	}
	// If we have a forwarded-for header, take the address from there
	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip.String()
		}
		// parse X-Real-Ip header
	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip.String()
		}
	}

	return remoteIP

}

//Split ....
func Split(str string, separator string) []string {

	returnSplittedString := make([]string, 0)
	splittedString := strings.Split(str, separator)
	if len(splittedString) == 0 && len(str) > 0 {
		splittedString = append(splittedString, str)
	}

	for _, element := range splittedString {
		if len(strings.TrimSpace(element)) > 0 {
			returnSplittedString = append(returnSplittedString, element)
		}
	}
	return returnSplittedString
}
