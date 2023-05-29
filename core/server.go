package core

import (
	"github.com/miekg/dns"
	"log"
	"net"
	"rapid-dns/db"
)

type Server struct {
}

func ServeDNS(w dns.ResponseWriter, req *dns.Msg) {
	resp := new(dns.Msg)
	resp.SetReply(req)

	for _, q := range req.Question {
		if q.Qtype == dns.TypeA {
			log.Println(q.Name)
			ip := db.Query(q.Name)
			aResp := dns.A{
				Hdr: dns.RR_Header{
					Name: q.Name, Ttl: 0, Rrtype: dns.TypeA, Class: dns.ClassINET,
				},
				A: net.ParseIP(string(ip)).To4(),
			}
			resp.Answer = append(resp.Answer, &aResp)
		}
	}
	w.WriteMsg(resp)
}

func (srv *Server) StartDNS() {
	dns.HandleFunc(".", ServeDNS)
	dns.ListenAndServe(":53", "udp", nil)
}
