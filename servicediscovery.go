package rnsd

import (
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
func (c *Client) ListNamespaces() ([]*servicediscovery.NamespaceSummary, error) {
	input := &servicediscovery.ListNamespacesInput{}

	resp, err := c.client.ListNamespaces(input)
	if err != nil {
		return nil, err
	}

	return resp.Namespaces, nil
}

func (c *Client) ListServices(nsIds ...string) ([]*servicediscovery.ServiceSummary, error) {
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

	resp, err := c.client.ListServices(input)
	if err != nil {
		return nil, err
	}

	return resp.Services, nil
}

func (c *Client) UpdateTTL(serviceId string, ttl int64) error {
	input := &servicediscovery.UpdateServiceInput{
		Id: aws.String(serviceId),
		Service: &servicediscovery.ServiceChange{
			DnsConfig: &servicediscovery.DnsConfigChange{
				DnsRecords: []*servicediscovery.DnsRecord{
					{TTL: aws.Int64(ttl)},
				},
			},
		},
	}

	// response has OperationId however does not return
	_, err := c.client.UpdateService(input)
	if err != nil {
		return err
	}

	return nil
}
