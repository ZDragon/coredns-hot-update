// Package hotupdate is a CoreDNS plugin that prints "example" to stdout on every packet received.
//
// It serves as an example CoreDNS plugin with numerous code comments.
package hotupdate

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"strings"

	//"github.com/coredns/coredns/plugin/metrics"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

const (
	port = ":50051"
)

// Define log to be a logger with the plugin name in it. This way we can just use log.Info and
// friends to log.
var log = clog.NewWithPlugin("hotupdate")

// server is used to implement helloworld.GreeterServer.
type server struct {
	UnimplementedDNSUpdaterServer
	ctx *HotUpdate
}

// HotUpdate Example is an example plugin to show how to write a plugin.
type HotUpdate struct {
	Next    plugin.Handler
	origins []string // for easy matching, these strings are the index in the map m.
	m       map[string][]dns.RR
}

// ServeDNS implements the plugin.Handler interface. This method gets called when example is used
// in a Server.
func (re HotUpdate) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	qname := state.Name()
	log.Infof("ServeDNS qname: %v", qname)
	zone := plugin.Zones(re.origins).Matches(qname)
	log.Infof("ServeDNS zone: %v", zone)
	if zone == "" {
		return plugin.NextOrFailure(re.Name(), re.Next, ctx, w, r)
	}

	// New we should have some data for this zone, as we just have a list of RR, iterate through them, find the qname
	// and see if the qtype exists. If so reply, if not do the normal DNS thing and return either NXDOMAIN or NODATA.
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	nxdomain := true
	var soa dns.RR
	for _, r := range re.m[zone] {
		if r.Header().Rrtype == dns.TypeSOA && soa == nil {
			soa = r
		}
		if r.Header().Name == qname {
			nxdomain = false
			if r.Header().Rrtype == state.QType() {
				m.Answer = append(m.Answer, r)
			}
		}
	}

	// handle NXDOMAIN, NODATA and normal response here.
	if nxdomain {
		m.Rcode = dns.RcodeNameError
		if soa != nil {
			m.Ns = []dns.RR{soa}
		}
		w.WriteMsg(m)
		return dns.RcodeSuccess, nil
	}

	if len(m.Answer) == 0 {
		if soa != nil {
			m.Ns = []dns.RR{soa}
		}
	}

	w.WriteMsg(m)
	return dns.RcodeSuccess, nil
}

// Name implements the Handler interface.
func (re HotUpdate) Name() string { return "hotupdate" }

// New returns a pointer to a new and intialized Records.
func New() *HotUpdate {
	re := new(HotUpdate)
	re.origins = make([]string, 1)
	re.m = make(map[string][]dns.RR)
	return re
}

func (s *server) Add(ctx context.Context, in *RequestDNSAdd) (*ResponseStatus, error) {
	log.Infof("Received: %v %v", in.Host, in.Ip)
	qname := plugin.Host(in.Host).Normalize()
	log.Infof("Origins len: %v", len(s.ctx.origins))
	zone := plugin.Zones(s.ctx.origins).Matches(qname)
	if zone == "" {
		rr, err := dns.NewRR("$ORIGIN " + qname + "\n" + in.Ip + "\n")
		if err != nil {
			return nil, err
		}
		rr.Header().Name = strings.ToLower(rr.Header().Name)
		s.ctx.origins = append(s.ctx.origins, qname)
		s.ctx.m[qname] = append(s.ctx.m[qname], rr)
	}

	return &ResponseStatus{Message: "Received " + in.Host + " IP " + in.Ip, Status: true}, nil
}
