package rnsd

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
)

func NewClient(region string) (*Client, error) {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("ap-northeast-1")}))
	client := servicediscovery.New(sess)

	return &Client{
		client: client,
	}, nil
}

type Client struct {
	client *servicediscovery.ServiceDiscovery
}

// not support paging
func (c *Client) ListNamespaces(ctx context.Context) ([]*servicediscovery.NamespaceSummary, error) {
	input := &servicediscovery.ListNamespacesInput{}

	resp, err := c.client.ListNamespacesWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	return resp.Namespaces, nil
}

func (c *Client) GetNamespace(ctx context.Context, namespaceId string) (*servicediscovery.Namespace, error) {
	input := &servicediscovery.GetNamespaceInput{
		Id: aws.String(namespaceId),
	}

	resp, err := c.client.GetNamespaceWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	return resp.Namespace, nil
}

func (c *Client) ListServices(ctx context.Context, nsIds ...string) ([]*servicediscovery.ServiceSummary, error) {
	f := &servicediscovery.ServiceFilter{
		Name: aws.String("NAMESPACE_ID"),
	}

	switch len(nsIds) {
	case 0:
		// required namespace id
		return nil, fmt.Errorf("required namespace id")
	case 1:
		f.Condition = aws.String("EQ")
		f.Values = []*string{aws.String(nsIds[0])}
	default:
		f.Condition = aws.String("IN")
		values := make([]*string, len(nsIds))
		for i, nsid := range nsIds {
			values[i] = aws.String(nsid)
		}
		f.Values = values
	}

	input := &servicediscovery.ListServicesInput{
		Filters: []*servicediscovery.ServiceFilter{f},
	}

	resp, err := c.client.ListServicesWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	return resp.Services, nil
}

func (c *Client) GetService(ctx context.Context, serviceId string) (*servicediscovery.Service, error) {
	input := &servicediscovery.GetServiceInput{
		Id: aws.String(serviceId),
	}

	resp, err := c.client.GetServiceWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	return resp.Service, nil
}

func (c *Client) UpdateTTL(ctx context.Context, serviceId, recordType string, ttl int64) error {
	input := &servicediscovery.UpdateServiceInput{
		Id: aws.String(serviceId),
		Service: &servicediscovery.ServiceChange{
			DnsConfig: &servicediscovery.DnsConfigChange{
				DnsRecords: []*servicediscovery.DnsRecord{
					{
						Type: aws.String(recordType),
						TTL:  aws.Int64(ttl),
					},
				},
			},
		},
	}

	// response has OperationId however does not return
	_, err := c.client.UpdateServiceWithContext(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
