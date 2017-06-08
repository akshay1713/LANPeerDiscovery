package LANPeerDiscovery

import (
	"net"
	"fmt"
	"strings"
	"strconv"
	"encoding/binary"
)

type UDPListener struct {
	listenerPort int
	listenerConn *net.UDPConn
	possibleLocalAddrs []string
	appName string
}

func (udpListener UDPListener) isMessageValid(addr *net.UDPAddr, msg []byte) bool{
	possibleLocalAddrs := udpListener.possibleLocalAddrs
	appName := udpListener.appName
	if pos(possibleLocalAddrs, addr.IP.String()+":"+strconv.Itoa(addr.Port)) != -1 {
		return false
	}

	if string(msg[0:len(appName)]) != appName {
		return false
	}
	return true
}

func (udpListener UDPListener) listenForUDPBroadcast(peerManager IPeerManager) {

	ServerConn := udpListener.listenerConn
	port := udpListener.listenerPort
	defer ServerConn.Close()
	appName := udpListener.appName
	portLen := 5
	typeLen := 0
	buf := make([]byte, len(appName)+portLen+typeLen)
	broadcastRecvdIPs := make(map[string]bool)
	portBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes, uint16(port))
	for {
		_, addr, err := ServerConn.ReadFromUDP(buf)
		if !udpListener.isMessageValid(addr, buf) {
			continue
		}
		recvdPort, err := strconv.Atoi(string(buf[len(appName):]))
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		if _, exists := broadcastRecvdIPs[addr.IP.String()]; !exists {
			broadcastRecvdIPs[addr.IP.String()] = true
			peerIPs := peerManager.GetAllIPs()
			all_ips := []byte{}
			msgLengthBytes := make([]byte, 4)
			totalLen := 2 + len(peerIPs)*6
			binary.BigEndian.PutUint32(msgLengthBytes, uint32(totalLen))
			all_ips = append(all_ips, 0)
			all_ips = append(all_ips, msgLengthBytes...)
			all_ips = append(all_ips, portBytes...)
			for i := range peerIPs {
				splitAddress := strings.Split(peerIPs[i], ":")
				peer_portBytes := make([]byte, 2)
				peer_port, _ := strconv.Atoi(splitAddress[1])
				binary.BigEndian.PutUint16(peer_portBytes, uint16(peer_port))
				splitIP := strings.Split(splitAddress[0], ".")
				for j := 0; j < 4; j++ {
					partIP, _ := strconv.Atoi(splitIP[j])
					all_ips = append(all_ips, byte(partIP))
				}
				all_ips = append(all_ips, peer_portBytes...)
			}
			tcpAddr := net.TCPAddr{IP: addr.IP, Port: recvdPort}
			sConn, err := net.DialTCP("tcp", nil, &tcpAddr)
			if err != nil {
				fmt.Println("Err while connecting to the source of broadcase message", err)
				continue
			}
			sConn.Write(all_ips)
		}
	}

}
