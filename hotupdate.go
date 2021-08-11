// Package hotupdate is a CoreDNS plugin that prints "example" to stdout on every packet received.
//
// It serves as an example CoreDNS plugin with numerous code comments.
package hotupdate

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/file"
	"strings"

	//"github.com/coredns/coredns/plugin/metrics"
	clog "github.com/coredns/coredns/plugin/pkg/log"
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
	Next plugin.Handler
	file file.File
}

// ServeDNS implements the plugin.Handler interface. This method gets called when example is used
// in a Server.
func (re *HotUpdate) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	return re.file.ServeDNS(ctx, w, r)
}

// Name implements the Handler interface.
func (re *HotUpdate) Name() string { return "hotupdate" }

// New returns a pointer to a new and intialized Records.
func New() *HotUpdate {
	re := new(HotUpdate)
	re.file = file.File{Zones: file.Zones{Z: make(map[string]*file.Zone), Names: []string{}}}
	return re
}

func (s *server) Add(ctx context.Context, in *RequestDNSAdd) (*ResponseStatus, error) {
	log.Infof("Received: %v %v", in.Host, in.Ip)
	qname := plugin.Host(in.Host).Normalize()
	log.Infof("Origins len: %v", len(s.ctx.file.Zones.Names))
	zone := plugin.Zones(s.ctx.file.Zones.Names).Matches(qname)
	if zone == "" {
		log.Infof("Zone %v empty, try add qname %v", zone, qname)
		rr, err := dns.NewRR("$ORIGIN " + qname + "\n" + in.Ip + "\n")
		if err != nil {
			return nil, err
		}
		rr.Header().Name = strings.ToLower(rr.Header().Name)
		z := file.NewZone(qname, "")
		if err := z.Insert(rr); err != nil {
			return nil, err
		}

		s.ctx.file.Zones.Z["."] = z
		s.ctx.file.Zones.Names = append(s.ctx.file.Zones.Names, ".")
	} else {
		log.Infof("Zone %v found, try add qname %v", zone, qname)
		rr, err := dns.NewRR("$ORIGIN " + qname + "\n" + in.Ip + "\n")
		if err != nil {
			return nil, err
		}
		rr.Header().Name = strings.ToLower(rr.Header().Name)
		z := s.ctx.file.Zones.Z[zone]
		if err := z.Insert(rr); err != nil {
			return nil, err
		}
		s.ctx.file.Zones.Z["."] = z
	}

	return &ResponseStatus{Message: "Received " + in.Host + " IP " + in.Ip, Status: true}, nil
}
