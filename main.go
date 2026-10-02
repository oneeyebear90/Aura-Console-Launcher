package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	TargetAddress = "bedrock.oneeyebear.net:19132"
	LocalListen   = "0.0.0.0:19132"
	DiscordLink   = "https://discord.gg/GA3qsy8vtE"
	AppVersion    = "v1.0.1"
	WindowTitle   = "Aura Minecraft Console Launcher"
	MaxLogs       = 5
	ServerName    = "Aura Minecraft"
)

// Fallback MOTD if Geyser is unreachable on initial boot
var defaultMOTD = fmt.Sprintf("MCPE;%s;0;Latest;0;100;1234567890;Aura Console Proxy;Survival;1;19132;19133;", ServerName)

var raknetMagic = []byte{0x00, 0xff, 0xff, 0x00, 0xfe, 0xfe, 0xfe, 0xfe, 0xfd, 0xfd, 0xfd, 0xfd, 0x12, 0x34, 0x56, 0x78}

var (
	resolvedIP string
	cachedMOTD string
	motdMutex  sync.RWMutex
	logBuffer  []string
	logMutex   sync.Mutex
)

func main() {
	setConsoleTitle(WindowTitle)

	targetAddr, err := net.ResolveUDPAddr("udp4", TargetAddress)
	if err != nil {
		fmt.Printf("[X] Error resolving target server (%s): %v\n", TargetAddress, err)
		waitForExit()
		return
	}
	resolvedIP = targetAddr.String()

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

	// Set initial fallback MOTD
	setCachedMOTD(defaultMOTD)

	// Fetch live MOTD from Geyser in background
	go fetchGeyserMOTD(targetAddr)

	addLog("Proxy initialized. Fetching live MOTD from Geyser...")

	// Packet Forwarding Routine
	go func() {
		for {
			buf := make([]byte, 2048)
			n, clientAddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				continue
			}

			// Intercept RakNet pings for instant response using cached Geyser MOTD
			if n > 1 && (buf[0] == 0x01 || buf[0] == 0x02) {
				go sendInstantPong(conn, clientAddr, buf[:n])
			} else {
				go forwardPacket(conn, targetAddr, clientAddr, buf[:n])
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n[!] Shutting down Aura Console Launcher. See you in-game!")
}

func fetchGeyserMOTD(targetAddr *net.UDPAddr) {
	outConn, err := net.DialUDP("udp4", nil, targetAddr)
	if err != nil {
		return
	}
	defer outConn.Close()

	var ping bytes.Buffer
	ping.WriteByte(0x01)
	_ = binary.Write(&ping, binary.BigEndian, uint64(time.Now().UnixMilli()))
	ping.Write(raknetMagic)
	_ = binary.Write(&ping, binary.BigEndian, uint64(0))

	_, err = outConn.Write(ping.Bytes())
	if err != nil {
		return
	}

	_ = outConn.SetReadDeadline(time.Now().Add(3 * time.Second))
	respBuf := make([]byte, 2048)
	n, _, err := outConn.ReadFromUDP(respBuf)
	if err != nil || n < 35 {
		return
	}

	if respBuf[0] == 0x1c {
		strLen := binary.BigEndian.Uint16(respBuf[33:35])
		if int(35+strLen) <= n {
			geyserMOTD := string(respBuf[35 : 35+strLen])
			setCachedMOTD(geyserMOTD)
			addLog("[+] Synchronized live MOTD header from Geyser")
		}
	}
}

func setCachedMOTD(motd string) {
	motdMutex.Lock()
	cachedMOTD = motd
	motdMutex.Unlock()
}

func getCachedMOTD() string {
	motdMutex.RLock()
	defer motdMutex.RUnlock()
	return cachedMOTD
}

func sendInstantPong(conn *net.UDPConn, clientAddr *net.UDPAddr, packetData []byte) {
	if len(packetData) < 19 {
		return
	}

	clientTime := binary.BigEndian.Uint64(packetData[1:9])
	motd := getCachedMOTD()

	var resp bytes.Buffer
	resp.WriteByte(0x1c)

	_ = binary.Write(&resp, binary.BigEndian, clientTime)
	_ = binary.Write(&resp, binary.BigEndian, uint64(9876543210))
	resp.Write(raknetMagic)

	motdBytes := []byte(motd)
	_ = binary.Write(&resp, binary.BigEndian, uint16(len(motdBytes)))
	resp.Write(motdBytes)

	_, _ = conn.WriteToUDP(resp.Bytes(), clientAddr)
	addLog(fmt.Sprintf("[+] Handled instant LAN ping from [%s]", clientAddr.IP.String()))
}

func addLog(message string) {
	logMutex.Lock()
	timestamp := time.Now().Format("15:04:05")
	entry := fmt.Sprintf("[%s] %s", timestamp, message)

	logBuffer = append(logBuffer, entry)
	if len(logBuffer) > MaxLogs {
		logBuffer = logBuffer[len(logBuffer)-MaxLogs:]
	}
	logMutex.Unlock()

	renderUI()
}

func renderUI() {
	fmt.Print("\033[H\033[2J")

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
	fmt.Printf("[+] Resolved target server: %s -> %s\n", TargetAddress, resolvedIP)
	fmt.Println("[+] Proxy active! Listening for console connections on port 19132...")
	fmt.Println("[+] Open Minecraft on your Console -> PLAY Button -> WORLDS section.")
	fmt.Println("[+] Look for Aura Minecraft with a LAN world tag. Whitelist on Discord.")
	fmt.Println("----------------------------------------------------------------------")

	logMutex.Lock()
	for _, log := range logBuffer {
		fmt.Println(log)
	}
	logMutex.Unlock()
}

func forwardPacket(conn *net.UDPConn, targetAddr, clientAddr *net.UDPAddr, packetData []byte) {
	outConn, err := net.DialUDP("udp4", nil, targetAddr)
	if err != nil {
		return
	}
	defer outConn.Close()

	_, err = outConn.Write(packetData)
	if err != nil {
		return
	}

	_ = outConn.SetReadDeadline(time.Now().Add(2 * time.Second))

	respBuf := make([]byte, 2048)
	respN, _, err := outConn.ReadFromUDP(respBuf)
	if err != nil {
		return
	}

	_, _ = conn.WriteToUDP(respBuf[:respN], clientAddr)
}

func waitForExit() {
	fmt.Println("\nPress Enter to exit...")
	var b [1]byte
	_, _ = os.Stdin.Read(b[:])
}