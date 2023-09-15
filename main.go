package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53domains"
)

func parseNames(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	names := []string{}
	var s bytes.Buffer
	for {
		b := make([]byte, 1)
		_, err := f.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if string(b[0]) == "\n" {
			names = append(names, s.String())
			s.Reset()
			continue
		}
		s.Write(b)
	}
	f.Close()
	return names, nil
}

func main() {
	names, err := parseNames("./names.csv")
	if err != nil {
		log.Fatalf("error loading names: %s", err.Error())
	}

	prefixes, err := parseNames("./prefix.csv")
	if err != nil {
		log.Fatalf("error loading prefix: %s", err.Error())
	}

	suffixes, err := parseNames("./suffix.csv")
	if err != nil {
		log.Fatalf("error loading suffix: %s", err.Error())
	}

	domains := []string{}
	for _, name := range names {
		domains = append(domains, name+".com")
		for _, prefix := range prefixes {
			domains = append(domains, prefix+name+".com")
		}

		for _, suffix := range suffixes {
			domains = append(domains, name+suffix+".com")
		}
	}
	log.Printf("domains: %d", len(domains))

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))
	svc := route53domains.New(sess)
	for _, domain := range domains {
		out, err := svc.CheckDomainAvailability(&route53domains.CheckDomainAvailabilityInput{
			DomainName: aws.String(domain),
		})
		if err != nil {
			log.Printf("error checking %s: %s", domain, err.Error())
		}
		if aws.StringValue(out.Availability) == "AVAILABLE" {
			fmt.Println(domain)
		}
	}
}
