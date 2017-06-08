package LANPeerDiscovery

import (
	"net"
	"log"
	"time"
)


type UDPBroadcaster struct {
	ports   []string
	appName string
}

func (udpBroadcaster UDPBroadcaster) init(manager IPeerManager) *net.UDPConn  {
	var serverConn *net.UDPConn
	for i := range udpBroadcaster.ports {
		serverAddr, err := net.ResolveUDPAddr("udp", ":"+udpBroadcaster.ports[i])
		if err != nil {
			log.Println("Error while resolving address ", err)
		}
		if serverConn == nil {
			serverConn, err = net.ListenUDP("udp", serverAddr)
			if err != nil {
				log.Println("Error while listening to address", err)
			}
		}
	}
	if serverConn == nil {
		panic("Unable to listen for UDP on any of the ports")
	}
	return serverConn
}

func (udpBroadcaster UDPBroadcaster) broadcastOnAllPorts(tcpListenerPort string) []string {
	ports := udpBroadcaster.ports
	var localAddrs []string
	for i := range ports {
		serverAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:"+ports[i])
		if err != nil {
			panic(err)
		}
		udpConn, err := net.DialUDP("udp", nil, serverAddr)
		if err != nil {
			panic(err)
		}
		go udpBroadcaster.broadcastOnSinglePort(udpConn, tcpListenerPort)
		localAddrs = append(localAddrs, udpConn.LocalAddr().String())
	}
	return localAddrs
}

func (udpBroadcaster UDPBroadcaster) broadcastOnSinglePort(conn *net.UDPConn, port string) {
	defer conn.Close()
	var msg []byte
	msg = append(msg, udpBroadcaster.appName...)
	port = padLeft(port, "0", 5)
	msg = append(msg, port...)
	buf := []byte(msg)
	for i := 0; i < 5; i++ {
		_, err := conn.Write(buf)
		if err != nil {
			log.Println("Error while broadcasting:", err)
			time.Sleep(time.Second * 1)
		}
	}
}
