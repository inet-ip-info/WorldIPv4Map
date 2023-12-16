package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/bits"
	"net"
	"sort"
	"strconv"
	"strings"
)

type IPRange struct {
	StartIP  net.IP
	NumHosts int
}

func parseIPv4FromRIPsFile(r io.Reader, allCIDRS map[string][]string) {
	cidrs := []string{}
	scanner := bufio.NewScanner(r)
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		i++
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 5 {
			log.Printf("line:%d invalid number of parts: %v", i, parts)
			continue
		}
		if parts[2] != "ipv4" {
			continue
		}
		cc := parts[1]
		if cc == "" {
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

		cidrs, err = convertToCIDR(startIP, numHosts)
		if err != nil {
			log.Printf("line:%d invalid number of hosts: %v", i, err)
			continue
		}
		if _, ok := allCIDRS[cc]; !ok {
			allCIDRS[cc] = []string{}
		}
		allCIDRS[cc] = append(allCIDRS[cc], cidrs...)
	}
}

func convertToCIDR(startIP net.IP, numHosts int) ([]string, error) {
	var cidrs []string
	remainingHosts := numHosts
	currentIP := startIP

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

func convertCIDRToRanges(cidrs []string) ([]IPRange, error) {
	var ranges []IPRange

	for _, cidr := range cidrs {
		ip, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}

		numHosts := 1 << (32 - bits.OnesCount32(uint32(ipNet.Mask[0])<<24|uint32(ipNet.Mask[1])<<16|uint32(ipNet.Mask[2])<<8|uint32(ipNet.Mask[3])))
		ranges = append(ranges, IPRange{StartIP: ip, NumHosts: numHosts})
	}

	// Sort the ranges based on the starting IP address.
	sort.Slice(ranges, func(i, j int) bool {
		return ipToInt(ranges[i].StartIP) < ipToInt(ranges[j].StartIP)
	})

	// Merge adjacent ranges.
	var mergedRanges []IPRange
	for _, r := range ranges {
		if len(mergedRanges) == 0 {
			mergedRanges = append(mergedRanges, r)
			continue
		}

		last := &mergedRanges[len(mergedRanges)-1]
		if ipToInt(nextIP(last.StartIP, last.NumHosts)) == ipToInt(r.StartIP) {
			// Extend the last range if it's adjacent to the current one.
			last.NumHosts += r.NumHosts
		} else {
			mergedRanges = append(mergedRanges, r)
		}
	}

	return mergedRanges, nil
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

func mapGetSortKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func dumpJson(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
