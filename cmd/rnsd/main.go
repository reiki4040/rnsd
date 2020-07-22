package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"log"
	"os"

	"github.com/reiki4040/rnsd"
	"github.com/urfave/cli/v2"
)

const (
	OPT_NAMESPACE_ID       = "namespace-id"
	OPT_NAMESPACE_ID_SHORT = "ns"
	OPT_SERVICE_ID         = "service-id"
	OPT_SERVICE_ID_SHORT   = "srv"
	OPT_TTL                = "ttl"
)

func CmdListNamespaces(c *cli.Context) error {
	// option and args

	return DoListNamespaces(c.Context, "ap-northeast-1")
}

func DoListNamespaces(ctx context.Context, region string) error {
	client, err := rnsd.NewClient(region)
	if err != nil {
		return err
	}

	nsList, err := client.ListNamespaces()
	if err != nil {
		return err
	}

	for _, ns := range nsList {
		fmt.Printf("%s\t%s\n", *ns.Id, *ns.Name)
	}

	return nil
}

func CmdListServices(c *cli.Context) error {
	// option and args
	nsId := c.String("namespace-id")

	return DoListServices(c.Context, "ap-northeast-1", nsId)
}

func DoListServices(ctx context.Context, region, nsId string) error {
	client, err := rnsd.NewClient(region)
	if err != nil {
		return err
	}

	ns, err := client.GetNamespace(nsId)
	if err != nil {
		return err
	}

	serviceList, err := client.ListServices(nsId)
	if err != nil {
		return err
	}

	for _, srv := range serviceList {
		fmt.Printf("%s\t%s\t%s", aws.StringValue(srv.Id), aws.StringValue(srv.Name), aws.StringValue(srv.Name)+"."+aws.StringValue(ns.Name))
		if len(srv.DnsConfig.DnsRecords) > 0 {
			for _, d := range srv.DnsConfig.DnsRecords {
				fmt.Printf("\t%s\t%d", aws.StringValue(d.Type), aws.Int64Value(d.TTL))
			}
		}
		fmt.Println()
	}

	return nil
}

func CmdModifyTTL(c *cli.Context) error {
	// option and args
	sId := c.String("service-id")
	ttl := c.Int64("ttl")

	return DoModifyTTL(c.Context, "ap-northeast-1", sId, ttl)
}

func DoModifyTTL(ctx context.Context, region, sId string, ttl int64) error {
	client, err := rnsd.NewClient(region)
	if err != nil {
		return err
	}

	err = client.UpdateTTL(sId, ttl)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var app = &cli.App{
		Name:  "rnsd",
		Usage: "control AWS Service Discovery",
		Commands: []*cli.Command{
			{
				Name:    "namespaces",
				Aliases: []string{"lns"},

				Usage:  "show namespaces",
				Action: CmdListNamespaces,
			},
			{
				Name:    "services",
				Aliases: []string{"lsrv"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     OPT_NAMESPACE_ID,
						Aliases:  []string{OPT_NAMESPACE_ID_SHORT},
						Usage:    "specify namespace id (like: ns-XXXXXXXXX)",
						Required: true,
					},
				},
				Usage:  "show services",
				Action: CmdListServices,
			},
			{
				Name:    "modify-ttl",
				Aliases: []string{"ttl"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     OPT_SERVICE_ID,
						Aliases:  []string{OPT_SERVICE_ID_SHORT},
						Usage:    "specify service id (like: srv-XXXXXXXXX)",
						Required: true,
					},
					&cli.Int64Flag{
						Name:     OPT_TTL,
						Aliases:  nil,
						Usage:    "specify set TTL seconds",
						Required: true,
					},
				},
				Usage:  "modify TTL of service",
				Action: CmdModifyTTL,
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
