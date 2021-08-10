package hotupdate

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"google.golang.org/grpc"
	"net"
)

// init registers this plugin.
func init() { plugin.Register("hotupdate", setup) }

// setup is the function that gets called when the config parser see the token "example". Setup is responsible
// for parsing any extra options the example plugin may have. The first token this function sees is "example".
func setup(c *caddy.Controller) error {
	c.Next() // Ignore "example" and give us the next token.
	if c.NextArg() {
		// If there was another token, return an error, because we don't have any configuration.
		// Any errors returned from this setup function should be wrapped with plugin.Error, so we
		// can present a slightly nicer error message to the user.
		return plugin.Error("hotupdate", c.ArgErr())
	}

	re := New()

	go func() {
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer()
		RegisterDNSUpdaterServer(s, &server{ctx: re})
		log.Infof("server listening at %v", lis.Addr())

		err1 := s.Serve(lis)
		if err1 != nil {
			log.Fatalf("failed to serve: %v", err1)
		}
	}()

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		re.Next = next
		return re
	})

	// All OK, return a nil error.
	return nil
}
