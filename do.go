package main

import (
	"context"
	"fmt"
)

func DoListNamespaces(ctx context.Context, region string) error {
	client, err := NewClient(region)
	if err != nil {
		return err
	}

	nsList, err := client.ListNamespaces(ctx)
	if err != nil {
		return err
	}

	for _, ns := range nsList {
		fmt.Printf("%s\t%s\n", *ns.Id, *ns.Name)
	}

	return nil
}
func DoListServices(ctx context.Context, region, nsId string) error {
	client, err := NewClient(region)
	if err != nil {
		return err
	}

	ns, err := client.GetNamespace(ctx, nsId)
	if err != nil {
		return err
	}

	serviceList, err := client.ListServices(ctx, nsId)
	if err != nil {
		return err
	}

	for _, srv := range serviceList {
		fmt.Printf("%s\t%s\t%s", *srv.Id, *srv.Name, *srv.Name+"."+*ns.Name)
		if len(srv.DnsConfig.DnsRecords) > 0 {
			for _, d := range srv.DnsConfig.DnsRecords {
				fmt.Printf("\t%s\t%d", d.Type, *d.TTL)
			}
		}
		fmt.Println()
	}

	return nil
}

func DoModifyTTL(ctx context.Context, region, sId string, ttl int64) error {
	client, err := NewClient(region)
	if err != nil {
		return err
	}

	s, err := client.GetService(ctx, sId)
	if err != nil {
		return err
	}

	if s.DnsConfig == nil {
		return fmt.Errorf("service does not have DNS")
	}

	if s.DnsConfig.DnsRecords == nil || len(s.DnsConfig.DnsRecords) == 0 {
		return fmt.Errorf("service does not have DNS")
	}

	/*
		var recordType types.RecordType
		switch s.DnsConfig.DnsRecords[0].Type {
		case "A":
			recordType = types.RecordTypeA
		case "CNAME":
			recordType = types.RecordTypeCname
		case "SRV":
			recordType = types.RecordTypeSrv
		default:
			return fmt.Errorf("unsupported recorde type %s: ", s.DnsConfig.DnsRecords[0].Type)
		}
	*/

	err = client.UpdateTTL(ctx, sId, s.DnsConfig.DnsRecords[0].Type, ttl)
	if err != nil {
		return err
	}

	return nil
}
