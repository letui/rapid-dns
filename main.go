package main

import (
	"log"
	"os"
	"rapid-dns/core"
	"rapid-dns/db"
	"rapid-dns/web"
	"strings"
)

func main() {
	dnsPort := os.Getenv("rapid-dns-port")
	webPort := os.Getenv("rapid-web-port")
	domainList := os.Getenv("rapid-domain-list")

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
	defer r.Close()
	db.LoadDefault()
	webSrv := web.Server{}
	go webSrv.StartWeb()
	dnsSrv := core.Server{}
	go dnsSrv.StartDNS()
	log.Printf("Listening and serving on :%s (http) / :%s (dns)", web.Port, core.Port)
	select {}
}
