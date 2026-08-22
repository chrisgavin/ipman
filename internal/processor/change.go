package processor

type Change struct {
	Kind      string `json:"kind"`
	Operation string `json:"operation"`
	Provider  string `json:"provider"`
	Network   string `json:"network"`
	Site      string `json:"site"`
	Pool      string `json:"pool"`
	Name      string `json:"name"`
	Type      string `json:"type,omitempty"`
	Summary   string `json:"summary"`
}

const (
	dnsKind  = "dns"
	dhcpKind = "dhcp"
)
