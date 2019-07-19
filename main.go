// filename: burst-dns
package main

import (
	"./dns"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

const nameServer = "nameserver 127.0.0.1"

func addToResolvConf() {
	// TODO: remove on exit
	resolvconfFile := viper.GetString("dns.resolvconf")
	for {
		b, err := ioutil.ReadFile(resolvconfFile)
		if err != nil {
			log.Printf("Fatal error resolv.conf: %s\n", err)
			continue
		}

		nameServers := string(b)
		if !strings.Contains(nameServers, nameServer) {
			nameServers = nameServer + "\n" + nameServers
			err := ioutil.WriteFile(resolvconfFile, []byte(nameServers), 0644)
			if err != nil {
				log.Printf("Fatal error update resolv.conf: %s\n", err)
				continue
			}
			log.Printf(nameServer+" was added to %s", resolvconfFile)
		}
		time.Sleep(5 * time.Second)
	}
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/burst-dns/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s\n", err)
	}

	go addToResolvConf()

	s := dns.BurstDNS{}
	s.Listen()
}
