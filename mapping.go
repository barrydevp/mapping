package mapping

import (
  "context"
  "net"
  "strconv"
  "fmt"

  "github.com/coredns/coredns/request"
  "github.com/coredns/coredns/plugin"

  "github.com/miekg/dns"
)

const name = "mapping"

type Mapping struct {
  Next plugin.Handler
  FalconInstance *Falcon
}

// Name implements the Handler interface.
func (m Mapping) Name() string { return name }

// ServeDNS implements the plugin.Handler interface.
func (m Mapping) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
  fmt.Println("Router has been called")

  state := request.Request{W: w, Req: r}

  a := new(dns.Msg)
  a.SetReply(r)
  a.Authoritative = true

  ip := state.IP()
  var rr dns.RR

  fmt.Println("ip: ", ip)

  switch state.Family() {
  case 1:
    rr = new(dns.A)
    rr.(*dns.A).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass()}
    rr.(*dns.A).A = net.ParseIP(ip).To4()
  case 2:
    rr = new(dns.AAAA)
    rr.(*dns.AAAA).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeAAAA, Class: state.QClass()}
    rr.(*dns.AAAA).AAAA = net.ParseIP(ip)
  }

  srv := new(dns.SRV)
  srv.Hdr = dns.RR_Header{Name: "_" + state.Proto() + "." + state.QName(), Rrtype: dns.TypeSRV, Class: state.QClass()}
  if state.QName() == "." {
    srv.Hdr.Name = "_" + state.Proto() + state.QName()
  }
  port, _ := strconv.Atoi(state.Port())
  srv.Port = uint16(port)
  srv.Target = "."

  a.Extra = []dns.RR{rr, srv}
  a.Answer = []dns.RR{rr, srv}

  w.WriteMsg(a)

  return 0, nil
}

