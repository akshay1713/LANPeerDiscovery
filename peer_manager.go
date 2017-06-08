package LANPeerDiscovery

type IPeerManager interface {
	GetAllIPs() []string
	IsConnected(IP string) bool
}
