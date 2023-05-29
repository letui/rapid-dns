package main

import (
	"log"
	"rapid-dns/core"
	"rapid-dns/db"
	"rapid-dns/web"
)

func main() {
	log.Println("Rapid-DNS")
	r := db.Init()
	defer r.Close()
	db.LoadDefault()
	webSrv := web.Server{}
	go webSrv.StartWeb()
	dnsSrv := core.Server{}
	log.Println(" Listening and serving on :8053 (http) / :1053 (dns)")
	dnsSrv.StartDNS()

}
