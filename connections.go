package LANPeerDiscovery

import (
	"net"
	"fmt"
	"strings"
	"io"
	"strconv"
	"encoding/binary"
)


func connectToPeer(ip net.IP, port int) (*net.TCPConn, error) {
	tcpAddr := net.TCPAddr{IP: ip, Port: port}
	chatConn, err := net.DialTCP("tcp", nil, &tcpAddr)
	return chatConn, err
}

type ConnAndType struct{
	Connection *net.TCPConn
	Type       string
}

func waitForTCP(peerManager IPeerManager, listener net.Listener, initiatorConn chan ConnAndType) {
	defer listener.Close()
	for {
		genericConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error while listening in waitforTCP", err)
		}
		conn := genericConn.(*net.TCPConn)
		senderIPString := strings.Split(conn.RemoteAddr().String(), ":")[0]
		senderIPOctets := strings.Split(senderIPString, ".")
		fmt.Println(senderIPString)
		handleErr(err, "Parsing IP")
		var senderIP []byte
		for i := 0; i < len(senderIPOctets); i++ {
			octetInt, _ := strconv.Atoi(senderIPOctets[i])
			senderIP = append(senderIP, byte(octetInt))
		}
		msgType := make([]byte, 1)
		_, err = io.ReadFull(conn, msgType)
		if !peerManager.IsConnected(senderIPString) {
			if msgType[0] == 0 {
				//This msg is a list of IPs & ports
				msgLength := make([]byte, 4)
				_, err = io.ReadFull(conn, msgLength)
				handleErr(err, "Error while reading message ")
				peerInfoLength := binary.BigEndian.Uint32(msgLength)
				peerInfo := make([]byte, peerInfoLength)
				_, err = io.ReadFull(conn, peerInfo)
				handleErr(err, "Error while reading message ")
				senderPort := binary.BigEndian.Uint16([]byte{peerInfo[0], peerInfo[1]})
				newConn, err := connectToPeer(senderIP, int(senderPort))
				handleErr(err, "Error while connecting to sender")
				if newConn == nil {
					fmt.Println("Nil conn")
					continue
				}
				connAndType := ConnAndType{Connection:newConn, Type:"sender"}
				initiatorConn <- connAndType
				for k := 2; k < len(peerInfo); k += 6 {
					peerIP := net.IPv4(peerInfo[k+2], peerInfo[k+3], peerInfo[k+4], peerInfo[k+5])
					peerPort := binary.BigEndian.Uint16([]byte{peerInfo[k], peerInfo[k+1]})
					if !peerManager.IsConnected(peerIP.String()) {
						newConn, err = connectToPeer(peerIP, int(peerPort))
						handleErr(err, "Error while connecting to peer")
						initiatorConn <- ConnAndType{Connection:newConn, Type:"sender"}
					}
				}
			} else {
				initiatorConn <- ConnAndType{Connection:conn, Type:"receiver"}
			}
		} else if msgType[0] == 1 {
			fmt.Println("Checking existing peer")
			initiatorConn <- ConnAndType{Connection:conn, Type:"duplicate_receiver"}
		}
	}
}
