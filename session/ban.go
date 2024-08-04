package session

import (
	"github.com/TonimatasDEV/ReposiGO/configuration"
	"strconv"
	"time"
)

var (
	retries = make(map[string]int)
	Bans    = make(map[string]int)
)

func checkBan(ip string) (bool, string) {
	if Bans[ip] > 0 {
		return true, "You are banned " + strconv.Itoa(Bans[ip]) + " seconds"
	}

	return false, ""
}

func addTry(ip string) {
	retries[ip] = retries[ip] + 1

	if retries[ip] >= configuration.ServerConfig.Security.Retries {
		Bans[ip] = configuration.ServerConfig.Security.BanTime
		delete(retries, ip)
	}
}

func BanHandler() {
	for {
		for entry, value := range Bans {
			Bans[entry] = value - 1
		}

		time.Sleep(1 * time.Second)
	}
}
