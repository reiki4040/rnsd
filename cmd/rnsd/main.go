package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

const (
	OPT_NAMESPACE_ID       = "namespace-id"
	OPT_NAMESPACE_ID_SHORT = "ns"
	OPT_SERVICE_ID         = "service-id"
	OPT_SERVICE_ID_SHORT   = "srv"
	OPT_TTL                = "ttl"
)

var (
	app = &cli.App{
		Name:  "rnsd",
		Usage: "control AWS Service Discovery",
		Commands: []*cli.Command{
			cmdNamespaces,
			cmdServices,
			cmdModifyTTL,
		},
	}

	cmdNamespaces = &cli.Command{
		Name:    "namespaces",
		Aliases: []string{"lns"},

		Usage:  "show namespaces",
		Action: CmdListNamespaces,
	}

	cmdServices = &cli.Command{
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
				Aliases:  nil,
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
	}
}

func CmdListNamespaces(c *cli.Context) error {
	// no option and args
	return DoListNamespaces(c.Context, "ap-northeast-1")
}

func CmdListServices(c *cli.Context) error {
	// option and args
	nsId := c.String(OPT_NAMESPACE_ID)
	if nsId == "" {
		return fmt.Errorf("namespace id is empty.")
	}

	return DoListServices(c.Context, "ap-northeast-1", nsId)
}

func CmdModifyTTL(c *cli.Context) error {
	// option and args
	sId := c.String(OPT_SERVICE_ID)
	if sId == "" {
		return fmt.Errorf("service id is empty.")
	}
	ttl := c.Int64(OPT_TTL)
	if ttl <= 0 {
		return fmt.Errorf("invalid TTL given.")
	}

	return DoModifyTTL(c.Context, "ap-northeast-1", sId, ttl)
}
