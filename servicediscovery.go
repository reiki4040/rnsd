package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery/types"
)

func NewClient(region string) (*Client, error) {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	svc := servicediscovery.NewFromConfig(cfg)

	return &Client{
		client: svc,
	}, nil
}

type Client struct {
	client *servicediscovery.Client
}

// not support paging
func (c *Client) ListNamespaces(ctx context.Context) ([]types.NamespaceSummary, error) {
	input := &servicediscovery.ListNamespacesInput{}

	resp, err := c.client.ListNamespaces(ctx, input)
	if err != nil {
		return nil, err
	}

	return resp.Namespaces, nil
}

func (c *Client) GetNamespace(ctx context.Context, namespaceId string) (*types.Namespace, error) {
	input := &servicediscovery.GetNamespaceInput{
		Id: aws.String(namespaceId),
	}

	resp, err := c.client.GetNamespace(ctx, input)
	if err != nil {
		return nil, err
	}

	return resp.Namespace, nil
}

func (c *Client) ListServices(ctx context.Context, nsIds ...string) ([]types.ServiceSummary, error) {
	f := types.ServiceFilter{
		Name: types.ServiceFilterNameNamespaceId,
	}

	switch len(nsIds) {
	case 0:
		// required namespace id
		return nil, fmt.Errorf("required namespace id")
	case 1:
		f.Condition = types.FilterConditionEq
		f.Values = nsIds
	default:
		f.Condition = types.FilterConditionIn
		f.Values = nsIds
	}

	input := &servicediscovery.ListServicesInput{
		Filters: []types.ServiceFilter{f},
	}

	resp, err := c.client.ListServices(ctx, input)
	if err != nil {
		return nil, err
	}

	return resp.Services, nil
}

func (c *Client) GetService(ctx context.Context, serviceId string) (*types.Service, error) {
	input := &servicediscovery.GetServiceInput{
		Id: aws.String(serviceId),
	}

	resp, err := c.client.GetService(ctx, input)
	if err != nil {
		return nil, err
	}

	return resp.Service, nil
}

func (c *Client) UpdateTTL(ctx context.Context, serviceId string, recordType types.RecordType, ttl int64) error {
	input := &servicediscovery.UpdateServiceInput{
		Id: aws.String(serviceId),
		Service: &types.ServiceChange{
			DnsConfig: &types.DnsConfigChange{
				DnsRecords: []types.DnsRecord{
					{
						Type: recordType,
						TTL:  aws.Int64(ttl),
					},
				},
			},
		},
	}

	// response has OperationId however does not return
	_, err := c.client.UpdateService(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
