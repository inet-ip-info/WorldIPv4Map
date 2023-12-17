# Country-Specific IPv4 Address List

## Overview
This project periodically collects IP address data managed by various Regional Internet Registries (RIRs), converts it into CIDR notation and subnet mask notation, and compiles it into a single file.

## Details
RIR statistical files are data files provided by Regional Internet Registries ([Regional Internet Registry, RIR](https://en.wikipedia.org/wiki/Regional_Internet_registry)) and are published for each region. Specifically, they can be obtained from the URLs listed [here](https://github.com/inet-ip-info/WorldIPv4Map/blob/8981e2c07987fc15be3f005c008b4ec1b960a72b/main.go#L12-L16). These files contain the range and number of IP addresses allocated to each country, allowing you to check the internet resource usage of each country. However, the original data is in the format of 'number of addresses from the start address', which is not suitable for direct use. This project converts these data into a more user-friendly format like CIDR notation (e.g., 192.168.0.0/24) or subnet mask notation (e.g., 192.168.0.0/255.255.255.0) and saves it.

## Usage
Data can be downloaded from the following URLs.

- CIDR Notation: [Download(all-ipv4cidr.tsv.gz)](https://github.com/inet-ip-info/WorldIPv4Map/releases/latest/download/all-ipv4cidr.tsv.gz)
- Subnet Mask Notation: [Download(all-ipv4mask.tsv.gz)](https://github.com/inet-ip-info/WorldIPv4Map/releases/latest/download/all-ipv4mask.tsv.gz)

## Forks Welcome
This repository welcomes your contributions and cooperation. Forking and participating is encouraged to improve the project and to ensure the continuous update and maintenance of this important information resource. We believe that by having many people independently manage and update this list, the accuracy of the data and the speed of updates will be enhanced, benefiting the entire community. We hope that this project will continue with the help of not just myself, but all of you.

## References
In writing this code, the [World's Country-Specific IPv4 Address Allocation List](http://nami.jp/ipv4bycc/) was a great reference. This site provides detailed specifications for converting the IP address list, which was a big hint for coding.

We are deeply grateful to the author of this site for providing us with the idea.
