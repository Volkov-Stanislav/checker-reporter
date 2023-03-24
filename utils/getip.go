package utils

import (
	"net"
)

func GetLocalIPAdresses() (result []string) {
	ifaces, err := net.Interfaces()

	if err != nil {
		return result
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return result
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			result = append(result, ip.String())
		}
	}

	return result
}
