package LANPeerDiscovery

import (
	"github.com/akshay1713/goUtils"
	"net"
	"strconv"
)

func GetConnectionsChan(candidatePorts []string, peerManager IPeerManager, appName string) chan ConnAndType {
	connChan := make(chan ConnAndType)
	tcpListener, err := net.Listen("tcp", ":0")
	goUtils.PanicErr(err)
	go waitForTCP(peerManager, tcpListener, connChan)
	toSendPort := strconv.Itoa(tcpListener.Addr().(*net.TCPAddr).Port)
	portInt, _ := strconv.Atoi(toSendPort)
	udpBroadcaster := UDPBroadcaster{ports: candidatePorts, appName: appName}
	broadcastListenerConn := udpBroadcaster.init(peerManager)
	possibleLocalPorts := udpBroadcaster.broadcastOnAllPorts(toSendPort)
	udpListener := UDPListener{
		listenerPort:       portInt,
		listenerConn:       broadcastListenerConn,
		possibleLocalAddrs: possibleLocalPorts,
		appName:            appName,
	}
	go udpListener.listenForUDPBroadcast(peerManager)
	return connChan
}
