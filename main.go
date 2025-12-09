package main

import (
        "flag"
        "fmt"
        "net"
        "os"
        "strconv"
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

        // iterate through all provided ports and launch a goroutine for each
        for _, port := range ports {
                go runUDPServer(port)
        }

        // Block main thread indefinitely
        select {}
}

func runUDPServer(port int) {
        addr := net.UDPAddr{
                Port: port,
                IP:   net.ParseIP("0.0.0.0"),
        }

        conn, err := net.ListenUDP("udp", &addr)
        if err != nil {
                fmt.Printf("Failed to bind to port %d: %v\n", port, err)
                return
        }
        defer conn.Close()

        fmt.Printf("Server listening on %s\n", conn.LocalAddr().String())

        buffer := make([]byte, 1024)
        for {
                n, remoteAddr, err := conn.ReadFromUDP(buffer)
                if err != nil {
                        fmt.Printf("[Port %d] Error receiving data: %v\n", port, err)
                        continue
                }

                if string(buffer[:n]) == "whoami" {
                        response := fmt.Sprintf("addr:%s", remoteAddr.String())
                        _, err := conn.WriteToUDP([]byte(response), remoteAddr)
                        if err != nil {
                                fmt.Printf("[Port %d] Error sending response: %v\n", port, err)
                                continue
                        }
                        fmt.Printf("[Port %d] Responded to %s\n", port, remoteAddr.String())
                }
        }
}
