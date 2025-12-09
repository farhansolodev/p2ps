package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"

	"golang.org/x/net/ipv4"
)

// portList implements flag.Value to accept multiple -p flags
type portList []int

func (p *portList) String() string {
	return fmt.Sprintf("%v", *p)
}

func (p *portList) Set(value string) error {
	port, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*p = append(*p, port)
	return nil
}

func main() {
	var ports portList
	flag.Var(&ports, "p", "UDP port to listen on (can be repeated)")
	flag.Parse()

	if len(ports) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No ports specified, e.g: -p 5000 -p 5001")
		os.Exit(1)
	}

	for _, port := range ports {
		go runUDPServer(port)
	}

	select {}
}

func runUDPServer(port int) {
	// 1. Create a PacketConn listening on 0.0.0.0 (all interfaces)
	conn, err := net.ListenPacket("udp4", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		fmt.Printf("Failed to bind to port %d: %v\n", port, err)
		return
	}
	defer conn.Close()

	// 2. Wrap it in an ipv4.PacketConn to access control messages
	pconn := ipv4.NewPacketConn(conn)
	if err := pconn.SetControlMessage(ipv4.FlagDst|ipv4.FlagInterface, true); err != nil {
		fmt.Printf("Failed to set control message: %v\n", err)
		return
	}

	fmt.Printf("Server listening on %s\n", conn.LocalAddr().String())

	buffer := make([]byte, 1024)
	for {
		// 3. Use ReadFrom to get the payload AND the control message (cm)
		n, cm, remoteAddr, err := pconn.ReadFrom(buffer)
		if err != nil {
			fmt.Printf("[Port %d] Error receiving data: %v\n", port, err)
			continue
		}

		if string(buffer[:n]) == "whoami" {
			// cm.Dst contains the IP address the client sent the packet TO (our local IP)
			// We must use this as the Source IP for the reply.
			response := fmt.Sprintf("addr:%s", remoteAddr.String())

			// 4. Create a ControlMessage for the reply, setting Src to the original Dst
			replyCM := &ipv4.ControlMessage{
				Src: cm.Dst,
			}

			// 5. WriteTo forces the kernel to use the specified Source IP
			_, err := pconn.WriteTo([]byte(response), replyCM, remoteAddr)
			if err != nil {
				fmt.Printf("[Port %d] Error sending response: %v\n", port, err)
				continue
			}
			fmt.Printf("[Port %d] Responded to %s via %s\n", port, remoteAddr.String(), cm.Dst.String())
		}
	}
}
