package edgeosclient

type GetResponse struct {
	Configuration ConfigurationRoot `json:"GET"`
}

type ConfigurationRoot struct {
	Service Service `json:"service"`
}

type Service struct {
	DHCPServer DHCPServer `json:"dhcp-server"`
}

type DHCPServer struct {
	SharedNetworkName map[string]SharedNetwork `json:"shared-network-name"`
}

type SharedNetwork struct {
	Subnet map[string]Subnet `json:"subnet"`
}

type Subnet struct {
	StaticMapping map[string]*StaticMapping `json:"static-mapping"`
}

type StaticMapping struct {
	IPAddress  string `json:"ip-address"`
	MACAddress string `json:"mac-address"`
}
