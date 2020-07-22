rnsd
===

[![GitHub release](http://img.shields.io/github/release/reiki4040/rnsd.svg?style=flat-square)][release]
![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)

[release]: https://github.com/reiki4040/rnsd/releases

control AWS Service Discovery command for modify TTL that ECS service SRV record.

features
- show namespaces.
- show services in the namespace.
- modify TTL of the service.

*currently status is developing. probably command interface(maybe output) will change.*

## install

recommend use homebrew 

```
brew tap reiki4040/tap/rnsd
```

or download binary or build from code.

## usage

```
NAME:
   rnsd - control AWS Service Discovery

USAGE:
   rnsd [global options] command [command options] [arguments...]

VERSION:
   0.1.0(68ab404)

COMMANDS:
   namespaces, ns   show namespaces
   services, srv    show services
   modify-ttl, ttl  modify TTL of service
   help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --region value, -r value  specify AWS region (default: "ap-northeast-1")
   --help, -h                show help (default: false)
   --version, -v             print the version (default: false)
```

show namespaces.
```
rnsd ns

ns-xxxxxxxxxxxxxxx  your.private
```

show services in the namespace.
```
rnsd srv -n ns-xxxxxxxxxxxxxxx

srv-xxxxxxxxxxxxxxxa    rp  rp.your.private SRV 300
srv-xxxxxxxxxxxxxxxb    rp2 rp2.your.private    SRV 300
srv-xxxxxxxxxxxxxxxc    rp3 rp3.your.private    SRV 45
```

modify service TTL.
```
rnsd ttl -s srv-xxxxxxxxxxxxxxxa -t 30
```
