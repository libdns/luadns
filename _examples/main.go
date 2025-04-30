// Use the following command line to run this example:
//
//	go run _examples/main.go -email your_email -key your_api_key -zone your_zone.com
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/libdns/libdns"
	"github.com/libdns/luadns"
)

var email string
var key string
var url string
var zone string

func main() {
	flag.StringVar(&email, "email", "joe@example.com", "your email address")
	flag.StringVar(&key, "key", "", "your API key")
	flag.StringVar(&zone, "zone", "example.org.", "zone name")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	provider := luadns.Provider{
		Email:  email,
		APIKey: key,
	}

	ctx := context.Background()
	fmt.Println("===> Running ListZones ...")
	zones, err := provider.ListZones(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	for _, z := range zones {
		fmt.Println(z)
	}

	fmt.Println("===> Running GetRecords ...")
	records, err := provider.GetRecords(ctx, zone)
	if err != nil {
		log.Fatalln(err)
	}
	for _, r := range records {
		fmt.Println(r)
	}

	fmt.Println("===> Running AppendRecords ...")
	created, err := provider.AppendRecords(ctx, zone, []libdns.Record{
		libdns.RR{
			Name: "foo",
			Type: "TXT",
			Data: "foo",
			TTL:  3600 * time.Second,
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	for _, r := range created {
		fmt.Println("created", r)
	}

	fmt.Println("===> Running SetRecords ...")
	updated, err := provider.SetRecords(ctx, zone, []libdns.Record{
		libdns.RR{
			Name: "foo",
			Type: "TXT",
			Data: "bar",
			TTL:  3600 * time.Second,
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	for _, r := range updated {
		fmt.Println("updated", r)
	}

	fmt.Println("===> Running DeleteRecords ...")
	deleted, err := provider.DeleteRecords(ctx, zone, []libdns.Record{
		libdns.RR{
			Name: "foo",
			Type: "TXT",
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	for _, r := range deleted {
		fmt.Println("deleted", r)
	}
}
