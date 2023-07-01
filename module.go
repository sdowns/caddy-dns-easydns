package easydns

import (
	"os"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	easydns "github.com/libdns/easydns"
)

// Provider lets Caddy read and manipulate DNS records hosted by this DNS provider.
type Provider struct{ *easydns.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.easydns",
		New: func() caddy.Module { return &Provider{new(easydns.Provider)} },
	}
}

// Provision sets up the module. Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()

	// get token value from env variable
	token := os.Getenv("CADDY_EASYDNS_TOKEN")
	key := os.Getenv("CADDY_EASYDNS_KEY")
	url := os.Getenv("CADDY_EASYDNS_URL")
	if url == "" {
		// Use production URL for the EasyDNS API by default.
		// The testing URL is: https://sandbox.rest.easydns.net
		url = "https://rest.easydns.net"
	}

	p.Provider.APIToken = repl.ReplaceAll(p.Provider.APIToken, token)
	p.Provider.APIKey = repl.ReplaceAll(p.Provider.APIKey, key)
	p.Provider.APIUrl = repl.ReplaceAll(p.Provider.APIUrl, url)
	return nil
}

// TODO: This is just an example. Update accordingly.
// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	providername [<api_token>] {
//	    api_token <api_token>
//	}
//
// **THIS IS JUST AN EXAMPLE AND NEEDS TO BE CUSTOMIZED.**
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		// if d.NextArg() {
		// 	// p.Provider.APIToken = d.Val()
		// }
		// if d.NextArg() {
		// 	return d.ArgErr()
		// }
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "api_token":
				if p.Provider.APIToken != "" {
					return d.Err("API token already set")
				}
				if d.NextArg() {
					p.Provider.APIToken = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "api_key":
				if p.Provider.APIKey != "" {
					return d.Err("API key already set")
				}
				if d.NextArg() {
					p.Provider.APIKey = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "api_url":
				if p.Provider.APIUrl != "" {
					return d.Err("API url already set")
				}
				if d.NextArg() {
					p.Provider.APIUrl = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.APIToken == "" {
		return d.Err("missing API token")
	}
	if p.Provider.APIKey == "" {
		return d.Err("missing API key")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
