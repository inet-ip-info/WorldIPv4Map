# IPv4 address list by country

## overview
This project regularly collects IP address data managed by multiple Regional Internet Registries (RIRs), converts it to CIDR notation and subnet mask notation, and compiles it into a single file.

## detail
RIR statistical files are data files provided by [Regional Internet Registries (RIR)](https://en.wikipedia.org/wiki/Regional_Internet_registry) and are published for each region. Specifically, it can be obtained from the URL described here [here](https://github.com/inet-ip-info/WorldIPv4Map/blob/8981e2c07987fc15be3f005c008b4ec1b960a72b/main.go#L12-L16).
The file lists the range and number of IP addresses assigned to each country, which allows you to see the Internet resource usage in each country.However, the original data is in the form "Number of addresses from starting address" and is not suitable for direct use.In this project, we will use CIDR notation ("Example: 192.168.0.0/24") and subnet mask notation ("Example: 192.168. 0.0/255.255.255.0") and saved.

## How to Use
You can download the data from the URL below.

### CIDR notation: 
- [Download(all-ipv4cidr.tsv.gz)](https://github.com/inet-ip-info/WorldIPv4Map/releases/latest/download/all-ipv4cidr.tsv.gz)
### Subnet mask notation: 
- [Download(all-ipv4mask.tsv.gz)](https://github.com/inet-ip-info/WorldIPv4Map/releases/latest/download/all-ipv4mask.tsv.gz)

## Welcome Fork
This repository welcomes your contributions and cooperation. We encourage you to fork and participate to make the project better and to ensure continued updating and maintenance of this important information resource. We believe that having more people manage and update this list independently will improve the accuracy of the data, speed of updates, and benefit the entire community. I hope that this project will continue, not just by me, but by everyone's hands.


## Reference site
When writing this code, I made great use of the [IPv4 address allocation list by country in the world] (http://nami.jp/ipv4bycc/). This site has detailed specifications for converting IP address lists, which gave me a great hint for coding.

I would like to express my sincere gratitude to the author of this site for giving me the idea.