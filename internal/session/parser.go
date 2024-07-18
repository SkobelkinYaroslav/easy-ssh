package session

import (
	"net"
	"regexp"
	"strconv"
)

func IsValidUsername(username string) bool {
	const usernamePattern = "^[a-zA-Z0-9._-]+$"
	matched, _ := regexp.MatchString(usernamePattern, username)
	return matched
}

func IsValidHostname(hostname string) bool {
	const hostnamePattern = "^(([a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)\\.)*[a-zA-Z]([a-zA-Z0-9-]*[a-zA-Z0-9])?$"
	matched, _ := regexp.MatchString(hostnamePattern, hostname)
	return matched || net.ParseIP(hostname) != nil
}

func IsValidPort(port string) (int, bool) {
	p, err := strconv.Atoi(port)
	if err != nil {
		return 0, false
	}
	if p < 1 || p > 65535 {
		return 0, false
	}
	return p, true
}
