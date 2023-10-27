/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

// IsIPv4 is used to determine whether ipAddr is an ipv4 address
func IsIPv4(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	return ip != nil && strings.Contains(ipAddr, ".")
}

// IsIPv6 is used to determine whether ipAddr is an ipv6 address
func IsIPv6(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	return ip != nil && strings.Contains(ipAddr, ":")
}

// CalcPathSHA256 is used to calculate the sha256 value
// of a file with a given path.
func CalcPathSHA256(fpath string) (string, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return CalcFileSHA256(f)
}

// CalcFileSHA256 is used to calculate the sha256 value
// of the file type.
func CalcFileSHA256(f *os.File) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// CalcSHA256 is used to calculate the sha256 value
// of the data.
func CalcSHA256(data []byte) (string, error) {
	if len(data) <= 0 {
		return "", errors.New("data is nil")
	}
	h := sha256.New()
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Get local ip address
func GetLocalIp() ([]string, error) {
	var result = make([]string, 0)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return result, err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				if ip[len(ip)-1] == byte(49) && ip[len(ip)-3] == byte(48) {
					continue
				}
				result = append(result, ipnet.IP.String())
			}
		}
	}
	return result, nil
}

// Get external ip address
func GetExternalIp() (string, error) {
	var (
		err        error
		externalIp string
	)

	var errstr string
	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
	}

	resp, err := client.Get("http://myexternalip.com/raw")
	if err != nil {
		errstr += err.Error()
	}
	if err == nil {
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		externalIp = fmt.Sprintf("%s", string(b))
		if IsIPv4(externalIp) {
			return externalIp, nil
		}
	}

	ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel1()
	output, err := exec.CommandContext(ctx1, "bash", "-c", "curl ifconfig.co").Output()
	if err != nil {
		errstr += err.Error()
	}
	if err == nil {
		externalIp = strings.ReplaceAll(string(output), "\n", "")
		externalIp = strings.ReplaceAll(externalIp, " ", "")
		if IsIPv4(externalIp) {
			return externalIp, nil
		}
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()
	output, err = exec.CommandContext(ctx2, "bash", "-c", "curl cip.cc | grep  IP | awk '{print $3;}'").Output()
	if err != nil {
		errstr += err.Error()
	}
	if err == nil {
		externalIp = strings.ReplaceAll(string(output), "\n", "")
		externalIp = strings.ReplaceAll(externalIp, " ", "")
		if IsIPv4(externalIp) {
			return externalIp, nil
		}
	}

	ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel3()
	output, err = exec.CommandContext(ctx3, "bash", "-c", `curl ipinfo.io | grep \"ip\" | awk '{print $2;}'`).Output()
	if err != nil {
		errstr += err.Error()
	}
	if err == nil {
		externalIp = strings.ReplaceAll(string(output), "\"", "")
		externalIp = strings.ReplaceAll(externalIp, ",", "")
		externalIp = strings.ReplaceAll(externalIp, "\n", "")
		if IsIPv4(externalIp) {
			return externalIp, nil
		}
	}

	ctx4, cancel4 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel4()
	output, err = exec.CommandContext(ctx4, "bash", "-c", `curl ifcfg.me`).Output()
	if err != nil {
		errstr += err.Error()
	}
	if err == nil {
		externalIp = strings.ReplaceAll(string(output), "\"", "")
		externalIp = strings.ReplaceAll(externalIp, ",", "")
		externalIp = strings.ReplaceAll(externalIp, "\n", "")
		if IsIPv4(externalIp) {
			return externalIp, nil
		}
	}

	return "", errors.New(errstr) //errors.New("please configure the public ip")
}

// ParseMultiaddrs
func ParseMultiaddrs(domain string) ([]string, error) {
	var result = make([]string, 0)
	var realDns = make([]string, 0)

	addr, err := ma.NewMultiaddr(domain)
	if err == nil {
		_, err = peer.AddrInfoFromP2pAddr(addr)
		if err == nil {
			result = append(result, domain)
			return result, nil
		}
	}

	dnsnames, err := net.LookupTXT(domain)
	if err != nil {
		return result, err
	}

	for _, v := range dnsnames {
		if strings.Contains(v, "ip4") && strings.Contains(v, "tcp") && strings.Count(v, "=") == 1 {
			result = append(result, strings.TrimPrefix(v, "dnsaddr="))
		}
	}

	trims := strings.Split(domain, ".")
	domainname := fmt.Sprintf("%s.%s", trims[len(trims)-2], trims[len(trims)-1])

	for _, v := range dnsnames {
		trims = strings.Split(v, "/")
		for _, vv := range trims {
			if strings.Contains(vv, domainname) {
				realDns = append(realDns, vv)
				break
			}
		}
	}

	for _, v := range realDns {
		dnses, err := net.LookupTXT("_dnsaddr." + v)
		if err != nil {
			continue
		}
		for i := 0; i < len(dnses); i++ {
			if strings.Contains(dnses[i], "ip4") && strings.Contains(dnses[i], "tcp") && strings.Count(dnses[i], "=") == 1 {
				var multiaddr = strings.TrimPrefix(dnses[i], "dnsaddr=")
				result = append(result, multiaddr)
			}
		}
	}

	return result, nil
}

func FreeLocalPort(port uint32) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), time.Second*3)
	if err != nil {
		return true
	}
	conn.Close()
	return false
}

func FindFile(dir, name string) string {
	var result string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == name {
			result = path
			return nil
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
	return result
}
