package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/reiki4040/rnsd"
)

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
