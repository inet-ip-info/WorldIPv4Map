package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// RIR URLs
var urls = []string{
	"https://ftp.arin.net/pub/stats/arin/delegated-arin-extended-latest",
	"https://ftp.ripe.net/pub/stats/ripencc/delegated-ripencc-extended-latest",
	"https://ftp.apnic.net/pub/stats/apnic/delegated-apnic-extended-latest",
	"https://ftp.lacnic.net/pub/stats/lacnic/delegated-lacnic-extended-latest",
	"https://ftp.afrinic.net/pub/stats/afrinic/delegated-afrinic-extended-latest",
}

// FetchData downloads data from the given URL.
func FetchData(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var lines []string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// Record represents an IPv4 allocation record.
type Record struct {
	CC       string
	StartIP  net.IP
	NumHosts int
}

// parseRecord parses a single line of RIR data into a Record.
func parseRecord(line string) (*Record, error) {
	parts := strings.Split(line, "|")
	if parts[2] != "ipv4" {
		return nil, fmt.Errorf("invalid record")
	}

	numHosts, err := strconv.Atoi(parts[4])
	if err != nil {
		return nil, fmt.Errorf("invalid number of hosts: %v", err)
	}

	return &Record{
		CC:       parts[1],
		StartIP:  net.ParseIP(parts[3]),
		NumHosts: numHosts,
	}, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Example usage
	lines := []string{}
	for _, url := range urls {
		var err error
		lines, err = FetchData(url)
		if err != nil {
			fmt.Println("Error fetching data:", err)
			continue
		}
		break // TODO:

	}

	for i, line := range lines {
		parts := strings.Split(line, "|")
		if parts[2] != "ipv4" {
			continue
		}

		startIP := net.ParseIP(parts[3])
		if startIP == nil || startIP.To4() == nil {
			//log.Printf("invalid IPv4 address: %s, line:[%s]", parts[3], line)
			continue
		}

		numHosts, err := strconv.Atoi(parts[4])
		if err != nil {
			log.Printf("line:%d invalid number of hosts: %v", i, err)
			return
		}

		record := Record{
			CC:       parts[1],
			StartIP:  net.ParseIP(parts[3]),
			NumHosts: numHosts,
		}
		fmt.Printf("Parsed record: %s\n", dumpJson(record))

		cidrs, err := convertToCIDR(record)
		if err != nil {
			fmt.Println("Error converting to CIDR:", err)
			continue
		}

		fmt.Printf("%s: %v\n", record.CC, cidrs)
	}
}

func dumpJson(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func highestBitSet(n int) int {
	if n == 0 {
		return -1 // Undefined for 0
	}

	pos := 0
	for n >>= 1; n > 0; n >>= 1 {
		pos++
	}
	return pos
}

func convertToCIDR(record Record) ([]string, error) {
	var cidrs []string
	remainingHosts := record.NumHosts
	currentIP := record.StartIP

	for remainingHosts > 0 {
		// Find the largest power of 2 that fits into the remaining hosts.
		maxSize := 32 - int(math.Log2(float64(remainingHosts)))
		if maxSize < 0 {
			return nil, fmt.Errorf("invalid number of hosts: %d", remainingHosts)
		}

		// Adjust the mask size to not exceed the boundaries of standard CIDR blocks.
		for {
			ipNet := net.IPNet{IP: currentIP, Mask: net.CIDRMask(maxSize, 32)}
			if !ipNet.Contains(nextIP(currentIP, 1<<uint(32-maxSize)-1)) {
				maxSize++
			} else {
				break
			}
		}

		// Add the CIDR block to the list.
		cidrs = append(cidrs, fmt.Sprintf("%s/%d", currentIP, maxSize))

		// Calculate the next starting IP and update the remaining hosts count.
		blockSize := 1 << uint(32-maxSize)
		currentIP = nextIP(currentIP, blockSize)
		remainingHosts -= blockSize
	}

	return cidrs, nil
}

// nextIP calculates the next IP address after a given block size.
func nextIP(ip net.IP, block int) net.IP {
	ipInt := ipToInt(ip)
	return intToIP(ipInt + uint32(block))
}

// Helper functions to convert between net.IP and uint32.
func ipToInt(ip net.IP) uint32 {
	ipv4 := ip.To4()
	if ipv4 == nil {
		panic(fmt.Sprintf("Invalid IPv4 address: %v", ip))
	}
	return uint32(ipv4[0])<<24 + uint32(ipv4[1])<<16 + uint32(ipv4[2])<<8 + uint32(ipv4[3])
}

func intToIP(ipInt uint32) net.IP {
	return net.IPv4(byte(ipInt>>24), byte(ipInt>>16&0xFF), byte(ipInt>>8&0xFF), byte(ipInt&0xFF))
}
