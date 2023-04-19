<div align="center">

# Go p2p library for cess distributed storage system

[![GitHub license](https://img.shields.io/badge/license-Apache2-blue)](#LICENSE) <a href=""><img src="https://img.shields.io/badge/golang-%3E%3D1.19-blue.svg" /></a> [![Go Reference](https://pkg.go.dev/badge/github.com/CESSProject/p2p-go.svg)](https://pkg.go.dev/github.com/CESSProject/p2p-go) [![Go Report Card](https://goreportcard.com/badge/github.com/CESSProject/p2p-go)](https://goreportcard.com/report/github.com/CESSProject/p2p-go)

</div>

The go p2p library implementation of the CESS distributed storage network, which is customized based on go-libp2p, has many advantages, and at the same time abandons the bandwidth and space redundancy caused by multiple file backups, making nodes more focused on storage allocation Give yourself data to make it more in line with the needs of the CESS network.

## Reporting a Vulnerability
If you find out any vulnerability, Please send an email to frode@cess.one, we are happy to communicate with you.

## Installation
To get the package use the standard:
```
go get -u "github.com/CESSProject/p2p-go"
```
Using Go modules is recommended.

## Documentation 
Please refer to https://pkg.go.dev/github.com/CESSProject/p2p-go

## Examples
Please refer to https://github.com/CESSProject/p2p-go/tree/main/examples

## Usage

The following is an example of creating an p2p host:
```
p2phost, err := p2pgo.New(
    privatekeyPath,
    p2pgo.ListenAddrStrings(addr, port), // regular tcp connections
    p2pgo.Workspace(workspace),
)
```
Creating a p2p host requires you to configure some information, you can refer to `options.go`.

## License
Licensed under [Apache 2.0](https://github.com/CESSProject/p2p-go/blob/main/LICENSE)
