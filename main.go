package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	TargetAddress = "bedrock.oneeyebear.net:19132"
	LocalListen   = "0.0.0.0:19132"
	DiscordLink   = "https://discord.gg/GA3qsy8vtE"
	AppVersion    = "v1.0.0"
)

func main() {
	printBanner()

	// Resolve target server IP
	targetAddr, err := net.ResolveUDPAddr("udp4", TargetAddress)
	if err != nil {
		fmt.Printf("[X] Error resolving target server (%s): %v\n", TargetAddress, err)
		waitForExit()
		return
	}

	// Listen on local UDP port 19132 for console ping requests
	localAddr, err := net.ResolveUDPAddr("udp4", LocalListen)
	if err != nil {
		fmt.Printf("[X] Error binding to local port 19132: %v\n", err)
		waitForExit()
		return
	}

	conn, err := net.ListenUDP("udp4", localAddr)
	if err != nil {
		fmt.Println("[X] Error: Could not start local listener on port 19132.")
		fmt.Println("    Make sure Minecraft PC or another proxy isn't open on this computer.")
		waitForExit()
		return
	}
	defer conn.Close()

	fmt.Println("[+] Proxy active! Listening for console connections on port 19132...")
	fmt.Println("[+] Open Minecraft on your Console -> Friends tab -> LAN Games section.")
	fmt.Println("----------------------------------------------------------------------")

	// Packet Forwarding Routine
	buf := make([]byte, 2048)
	go func() {
		for {
			n, clientAddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				continue
			}

			// Forward incoming console packet to Aura Bedrock server
			outConn, err := net.DialUDP("udp4", nil, targetAddr)
			if err != nil {
				continue
			}

			_, _ = outConn.Write(buf[:n])

			// Relay response from server back to console client
			respBuf := make([]byte, 2048)
			respN, _, err := outConn.ReadFromUDP(respBuf)
			if err == nil {
				_, _ = conn.WriteToUDP(respBuf[:respN], clientAddr)
			}
			outConn.Close()
		}
	}()

	// Wait for user termination signal (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n[!] Shutting down Aura Console Launcher. See you in-game!")
}

func printBanner() {
	banner := `
  █████╗ ██╗  ██╗██████╗  █████╗ 
 ██╔══██╗██║  ██║██╔══██╗██╔══██╗
 ███████║██║  ██║██████╔╝███████║
 ██╔══██║██║  ██║██╔══██╗██╔══██║
 ██║  ██║╚█████╔╝██║  ██║██║  ██║
 ╚═╝  ╚═╝ ╚════╝ ╚═╝  ╚═╝╚═╝  ╚═╝
`
	fmt.Println(banner)
	fmt.Println("======================================================================")
	fmt.Printf("           AURA MINECRAFT CONSOLE LAUNCHER (%s)\n", AppVersion)
	fmt.Println("======================================================================")
	fmt.Printf(" [!] Destination Server: %s\n", TargetAddress)
	fmt.Printf(" [!] Need help or support? Join our Discord:\n     %s\n", DiscordLink)
	fmt.Println("======================================================================")
	fmt.Println(" [!] KEEP THIS WINDOW OPEN WHILE PLAYING ON YOUR CONSOLE!")
	fmt.Println("======================================================================")
	fmt.Println()
}

func waitForExit() {
	fmt.Println("\nPress Enter to exit...")
	var b [1]byte
	_, _ = os.Stdin.Read(b[:])
}
