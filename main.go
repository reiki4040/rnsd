package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

const (
	OPT_REGION       = "region"
	OPT_REGION_SHORT = "r"

	OPT_NAMESPACE_ID       = "namespace-id"
	OPT_NAMESPACE_ID_SHORT = "n"
	OPT_SERVICE_ID         = "service-id"
	OPT_SERVICE_ID_SHORT   = "s"
	OPT_TTL                = "ttl"
	OPT_TTL_SHORT          = "t"
)

var (
	version  string
	revision string

	app = &cli.App{
		Name:    "rnsd",
		Usage:   "control AWS Service Discovery",
		Version: version + "(" + revision + ")",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     OPT_REGION,
				Aliases:  []string{OPT_REGION_SHORT},
				Usage:    "specify AWS region",
				Required: false,
				Value:    "ap-northeast-1",
			},
		},

		Commands: []*cli.Command{
			cmdNamespaces,
			cmdServices,
			cmdModifyTTL,
		},
	}

	cmdNamespaces = &cli.Command{
		Name:    "namespaces",
		Aliases: []string{"ns"},

		Usage:  "show namespaces",
		Action: CmdListNamespaces,
	}

	cmdServices = &cli.Command{
		Name:    "services",
		Aliases: []string{"srv"},
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
	}

	cmdModifyTTL = &cli.Command{
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
				Aliases:  []string{OPT_TTL_SHORT},
				Usage:    "specify set TTL seconds",
				Required: true,
			},
		},
		Usage:  "modify TTL of service",
		Action: CmdModifyTTL,
	}
)

func main() {
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func CmdListNamespaces(c *cli.Context) error {
	r := c.String(OPT_REGION)

	return DoListNamespaces(c.Context, r)
}

func CmdListServices(c *cli.Context) error {
	// option and args
	r := c.String(OPT_REGION)
	nsId := c.String(OPT_NAMESPACE_ID)
	if nsId == "" {
		return fmt.Errorf("namespace id is empty.")
	}

	return DoListServices(c.Context, r, nsId)
}

func CmdModifyTTL(c *cli.Context) error {
	// option and args
	r := c.String(OPT_REGION)
	sId := c.String(OPT_SERVICE_ID)
	if sId == "" {
		return fmt.Errorf("service id is empty.")
	}
	ttl := c.Int64(OPT_TTL)
	if ttl <= 0 {
		return fmt.Errorf("invalid TTL given.")
	}

	return DoModifyTTL(c.Context, r, sId, ttl)
}
