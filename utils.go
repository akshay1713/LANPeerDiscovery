package LANPeerDiscovery

import "fmt"


func padLeft(str, pad string, length int) string {
	if len(str) >= length {
		return str
	}
	for {
		str = pad + str
		if len(str) > length {
			return str[0:length]
		}
	}
}

func handleErr(err error, prefix string) {
	if err != nil {
		fmt.Println(prefix, ": ", err)
	}
}

func panicErr(err error){
	if err != nil {
		panic(err)
	}
}

func pos(slice []string, value string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}
