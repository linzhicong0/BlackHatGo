package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"os"
)

func main() {

	var (
		flDomain      = flag.String("domain", "", "The domain to perform guessing against.")
		flWordList    = flag.String("wordlist", "", "The word list to use for guessing.")
		flWorkerCount = flag.Int("c", 100, "The amount of workers to use.")
		flServerAddr  = flag.String("server", "8.8.8.8:53", "The DNS server to use.")
	)

	flag.Parse()

	if *flDomain == "" || *flWordList == "" {
		fmt.Println("-domain and -wordlist are required")
		os.Exit(1)
	}

	fmt.Println(*flWorkerCount, *flServerAddr)

}

type result struct {
	IPAddress string
	Hostname  string
}

func lookupA(fqdn, serverAddr string) ([]string, error) {

	var msg dns.Msg

	var ips []string

	msg.SetQuestion(dns.Fqdn(fqdn), dns.TypeA)

	in, err := dns.Exchange(&msg, serverAddr)

	if err != nil {
		return ips, err
	}

	if len(in.Answer) < 1 {
		return ips, errors.New("no answer")
	}

	for _, answer := range in.Answer {
		if a, ok := answer.(*dns.A); ok {
			ips = append(ips, a.A.String())
		}
	}

	return ips, nil
}

func lookupCNAME(fqdn, serverAddr string) ([]string, error) {

	var msg dns.Msg

	var fqdns []string

	msg.SetQuestion(dns.Fqdn(fqdn), dns.TypeCNAME)

	in, err := dns.Exchange(&msg, serverAddr)

	if err != nil {
		return fqdns, err
	}

	if len(in.Answer) < 1 {
		return fqdns, errors.New("no answer")
	}

	for _, answer := range in.Answer {
		if a, ok := answer.(*dns.CNAME); ok {
			fqdns = append(fqdns, a.Target)
		}
	}

	return fqdns, nil
}

func lookup(fqdn, serverAddr string) []result {


	var results []result

	var cfqdn = fqdn

	for {
		cnames, err := lookupCNAME(cfqdn, serverAddr)

		if err != nil {
			cfqdn = cnames[0]
			continue
		}

		ips, err := lookupA(cfqdn, serverAddr)

		if err != nil {
			break
		}

		for _, ip := range ips{
			results = append(results, result{
				IPAddress: ip,
				Hostname:  fqdn,
			})
		}
		break
	}

	return results

}