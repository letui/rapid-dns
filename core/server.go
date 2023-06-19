package core

import (
	"github.com/miekg/dns"
	"log"
	"net"
	"rapid-dns/db"
)

var Port = "53"

type Server struct {
}

func ServeDNS(w dns.ResponseWriter, req *dns.Msg) {
	resp := new(dns.Msg)
	resp.SetReply(req)

	for _, q := range req.Question {
		if q.Qtype == dns.TypeA {
			ip := db.Query(q.Name)
			if ip != nil {
				aResp := dns.A{
					Hdr: dns.RR_Header{
						Name: q.Name, Ttl: 0, Rrtype: dns.TypeA, Class: dns.ClassINET,
					},
					A: net.ParseIP(string(ip)).To4(),
				}
				resp.Answer = append(resp.Answer, &aResp)
			} else {
				m := new(dns.Msg)
				m.SetQuestion(dns.Fqdn(q.Name), dns.TypeA)
				c := new(dns.Client)
				in, rtt, err := c.Exchange(m, "114.114.114.114:53")
				if err != nil {
					log.Println(err, rtt.String())
				}
				resp.Answer = append(resp.Answer, in.Answer[0])
			}
		}
	}
	w.WriteMsg(resp)
}

func (srv *Server) StartDNS() {
	dns.HandleFunc(".", ServeDNS)
	dns.ListenAndServe("0.0.0.0:"+Port, "udp", nil)
}
