package main

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
)

// RIR URLs
var urls = []string{
	"https://ftp.arin.net/pub/stats/arin/delegated-arin-extended-latest",
	"https://ftp.ripe.net/pub/stats/ripencc/delegated-ripencc-extended-latest",
	"https://ftp.apnic.net/pub/stats/apnic/delegated-apnic-extended-latest",
	"https://ftp.lacnic.net/pub/stats/lacnic/delegated-lacnic-extended-latest",
	"https://ftp.afrinic.net/pub/stats/afrinic/delegated-afrinic-extended-latest",
}

const (
	CIDRMode       = false
	SubnetMaskMode = true
)

type env struct {
	Mask bool
}

func main() {
	var e env
	envconfig.Process("", &e)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	allCIDRS := map[string][]string{}
	for _, url := range urls {
		f, err := OpenURLFile(url)
		if err != nil {
			log.Printf("Error fetching data: url:%s err:%s", url, err)
			os.Exit(1)
			//continue
		}
		parseIPv4FromRIPsFile(f, allCIDRS)
		f.Close()
		log.Printf("fetch data from %s success", url)
	}
	allCount := 0
	mergeCount := 0
	keys := mapGetSortKeys(allCIDRS)
	printFunc := map[bool]func(string, string){
		CIDRMode: func(cc, ip string) {
			os.Stdout.Write([]byte(cc + "\t" + ip + "\n"))
		},
		SubnetMaskMode: func(cc, ip string) {
			maskip, err := cidrToSubnetMask(ip)
			if err != nil {
				log.Printf("cc:%s convert cidr to subnet mask error:%v", cc, err)
				return

			}
			os.Stdout.Write([]byte(cc + "\t" + maskip + "\n"))
		},
	}[e.Mask]
	for _, cc := range keys {
		cidrs := allCIDRS[cc]
		allCount += len(cidrs)
		ipranges, err := convertCIDRToRanges(cidrs)
		if err != nil {
			log.Printf("cc:%s convert cidr to ranges error:%v", cc, err)
			continue
		}
		for _, iprange := range ipranges {
			mergeCiders, err := convertToCIDR(iprange.StartIP, iprange.NumHosts)
			if err != nil {
				log.Printf("cc:%s convert cidr to ranges error:%v", cc, err)
				continue
			}
			for _, mergeCider := range mergeCiders {
				printFunc(cc, mergeCider)
			}
			mergeCount += len(mergeCiders)
		}
	}
	log.Printf("%d CIDRs", allCount)
	log.Printf("%d CIDRs(merged)", mergeCount)

}
