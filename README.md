# Country-Specific IPv4 Address List

## Overview
This project periodically collects IP address data managed by various Regional Internet Registries (RIRs), converts it into CIDR notation and subnet mask notation, and compiles it into a single file.

## Details
RIR statistical files are data files provided by Regional Internet Registries ([Regional Internet Registry, RIR](https://en.wikipedia.org/wiki/Regional_Internet_registry)) and are published for each region. Specifically, they can be obtained from the URLs listed [here](https://github.com/inet-ip-info/WorldIPv4Map/blob/8981e2c07987fc15be3f005c008b4ec1b960a72b/main.go#L12-L16). These files contain the range and number of IP addresses allocated to each country, allowing you to check the internet resource usage of each country. However, the original data is in the format of 'number of addresses from the start address', which is not suitable for direct use. This project converts these data into a more user-friendly format like CIDR notation (e.g., 192.168.0.0/24) or subnet mask notation (e.g., 192.168.0.0/255.255.255.0) and saves it.

## Usage
Data can be downloaded from the following URLs.

- CIDR Notation: [Download(all-ipv4cidr.tsv.gz)](https://github.com/inet-ip-info/WorldIPv4Map/releases/latest/download/all-ipv4cidr.tsv.gz)
- Subnet Mask Notation: [Download(all-ipv4mask.tsv.gz)](https://github.com/inet-ip-info/WorldIPv4Map/releases/latest/download/all-ipv4mask.tsv.gz)

## Usage Example
The following example demonstrates how to use ipset and iptables to only allow IPv4 addresses with the country code JP.

### Installing ipset Command
```bash
apt install -y ipset
```

### Downloading and Setting Up all-ipv4cidr.tsv.gz
```bash
URL=https://github.com/inet-ip-info/WorldIPv4Map/releases/latest/download/all-ipv4cidr.tsv.gz
CIDRFILE=/var/lib/ipset/ipset_list
TIMEOUT_DAYS=7
SETNAME=allow_list

find $CIDRFILE -type f -mtime +$TIMEOUT_DAYS -exec rm -f {} \;
[[ -f $CIDRFILE ]] ||
	curl -sL $URL |
	zcat |
	sed -n 's/^JP\t//p' \
		>$CIDRFILE

/usr/sbin/ipset create $SETNAME hash:net
/usr/sbin/ipset flush $SETNAME 2>/tmp/ipset.err.log

while read line; do
	/usr/sbin/ipset add $SETNAME $line 2>>/tmp/ipset.err.log
done <$CIDRFILE
```

### Configuring iptables to Allow Specific Ports
```bash
# UDP（26900-26903)
/sbin/iptables -A INPUT -p udp --dport 26900:26903 -m set --match-set $SETNAME src -j ACCEPT
/sbin/iptables -A INPUT -p udp --dport 26900:26903 -j DROP
```

## About Automatic Updates
This IP list is automatically updated daily at 3 AM Japan time using GitHub Actions. It is recommended to regularly download the data to maintain the most current information.



## Forks Welcome
This repository welcomes your contributions and cooperation. Forking and participating is encouraged to improve the project and to ensure the continuous update and maintenance of this important information resource. We believe that by having many people independently manage and update this list, the accuracy of the data and the speed of updates will be enhanced, benefiting the entire community. We hope that this project will continue with the help of not just myself, but all of you.

## References
In writing this code, the [World's Country-Specific IPv4 Address Allocation List](http://nami.jp/ipv4bycc/) was a great reference. This site provides detailed specifications for converting the IP address list, which was a big hint for coding.

We are deeply grateful to the author of this site for providing us with the idea.
