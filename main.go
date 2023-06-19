package main

import (
	"github.com/boltdb/bolt"
	"log"
	"os"
	"rapid-dns/core"
	"rapid-dns/db"
	"rapid-dns/web"
	"strings"
)

func main() {
	dnsPort := os.Getenv("rapid_dns_port")
	webPort := os.Getenv("rapid_web_port")
	domainList := os.Getenv("rapid_domain_list")

	if dnsPort != "" {
		core.Port = dnsPort
	}
	if webPort != "" {
		web.Port = webPort
	}
	if domainList != "" {
		extends := strings.Split(domainList, ",")
		db.Domains = append(db.Domains, extends...)
	}
	log.Println("Rapid-DNS")
	r := db.Init()
	defer func(r *bolt.DB) {
		err := r.Close()
		if err != nil {
			log.Fatal("Db got error while closing")
		}
	}(r)
	db.LoadDefault()
	webSrv := web.Server{}
	go webSrv.StartWeb()
	dnsSrv := core.Server{}
	go dnsSrv.StartDNS()
	log.Printf("Listening and serving on :%s (http) / :%s (dns)", web.Port, core.Port)
	select {}
}
