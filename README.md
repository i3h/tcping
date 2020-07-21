# Introduction

Tcping is ping probe command line tool, supporting ICMP, TCP and HTTP protocols.

You can also use it to query IP information from [https://ifconfig.is](https://ifconfig.is).

# Features

- Support ICMP/TCP/HTTP protocols
- Query basic IP information

# Installation

1. Download latest [release](https://github.com/i3h/tcping/releases/latest) (recommend)

2. Use go get

```
go get -u github.com/i3h/tcping/cmd/tcping
```

3. Build on your own

```
git clone https://github.com/i3h/tcping.git
cd tcping/cmd/tcping
go build
```

# Usage

```
  -h string
        HTTP Ping
  -i string
        ICMP Ping
  -m string
        MTR Trace
  -q string
        Query ip information
  -t string
        TCP Ping
  -v    Version
```

# Examples

```bash
# Test port
$ tcping google.com 443

TCP    OPEN      [2404:6800:4003:c03::71]:443

# Test with protocol
$ tcping https://google.com

Continent:    North America
Country  :    United States
Latitude :    37.751000
Longitude:    -97.822000
TimeZone :    America/Chicago
IsEU     :    false
ASN      :    15169
ORG      :    GOOGLE

Proxy     :    false
Scheme    :    https
Host      :    google.com
DNS Lookup:    0.85 ms
TCP       :    1.62 ms
TLS       :    3.11 ms
Process   :    32.91 ms
Transfer  :    0.15 ms
Total     :    38.72 ms

# HTTP ping
$ tcping -h https://google.com

Proxy     :    false
Scheme    :    https
Host      :    google.com
DNS Lookup:    0.92 ms
TCP       :    1.71 ms
TLS       :    2.99 ms
Process   :    32.24 ms
Transfer  :    0.14 ms
Total     :    38.10 ms

# ICMP ping
$ tcping -i google.com

ICMP   OPEN      172.217.194.113    2.0 ms

# Query IP info
$ tcping -q google.com

Continent:    North America
Country  :    United States
Latitude :    37.751000
Longitude:    -97.822000
TimeZone :    America/Chicago
IsEU     :    false
ASN      :    15169
ORG      :    GOOGLE

# Test port
$ tcping -t google.com:443

TCP    OPEN      google.com:443
```

# Note

Root permission is required when running ICMP ping, since it needs to open raw socket.

You can either use sudo command, or set setuid bit for tcping.

```bash
# Use sudo for one-time ping
$ sudo tcping -i google.com

# Set setuid bit
$ sudo chown root:root tcping
$ sudo chmod u+s tcping

```

# License

See the [LICENSE](https://github.com/i3h/tcping/blob/master/LICENSE.md) file for license rights and limitations (MIT).

# Acknowledgements

[lmas/icmp_ping.go](https://gist.github.com/lmas/c13d1c9de3b2224f9c26435eb56e6ef3)

[sparrc/go-ping](https://github.com/sparrc/go-ping)

[davecheney/httpstat](https://github.com/davecheney/httpstat)
