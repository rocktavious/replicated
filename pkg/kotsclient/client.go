package kotsclient

import (
	"fmt"
	"github.com/pkg/errors"
	channels "github.com/replicatedhq/replicated/gen/go/v1"
	"github.com/replicatedhq/replicated/pkg/graphql"
	"github.com/replicatedhq/replicated/pkg/types"
)

type Client interface {
	ListApps() ([]types.AppAndChannels, error)
	GetApp(appID string) (*types.App, error)

	CreateRelease(appID string, multiyaml string) (*types.ReleaseInfo, error)
	ListReleases(appID string) ([]types.ReleaseInfo, error)
	UpdateRelease(appID string, sequence int64, yaml string) error
	PromoteRelease(appID string, sequence int64, label string, notes string, channelIDs ...string) error

	ListChannels(appID string) ([]types.Channel, error)
	CreateChannel(appID string, name string, description string) (string, error)
	GetChannel(appID, channelID string) (*channels.AppChannel, []channels.ChannelRelease, error)

	ListCustomers(appID string) ([]types.Customer, error)
}

type AppOptions struct {
	Name string
}

type ChannelOptions struct {
	Name        string
	Description string
}

// Client communicates with the Replicated Vendor GraphQL API.
type GraphQLClient struct {
	GraphQLClient *graphql.Client
}

func NewGraphQLClient(origin string, apiKey string) *GraphQLClient {
	c := &GraphQLClient{GraphQLClient: graphql.NewClient(origin, apiKey)}

	return c
}

func (c *GraphQLClient) ExecuteRequest(requestObj graphql.Request, deserializeTarget interface{}) error {
	return c.GraphQLClient.ExecuteRequest(requestObj, deserializeTarget)
}

func (c *GraphQLClient) GetCustomerByNameOrID(appID string, customerNameOrID string) (*types.Customer, error) {
	customers, err := c.ListCustomers(appID)
	if err != nil {
		return nil, errors.Wrap(err, "list customers")
	}

	matchingCustomers := make([]*types.Customer, 0)
	for _, customer := range customers {
		if customer.ID == customerNameOrID || customer.Name == customerNameOrID {
			matchingCustomers = append(matchingCustomers, &types.Customer{
				ID:   customer.ID,
				Name: customer.Name,
			})
		}
	}

	if len(matchingCustomers) == 0 {
		return nil, fmt.Errorf("could not find customer %q", customerNameOrID)
	}

	if len(matchingCustomers) > 1 {
		return nil, fmt.Errorf("customer %q is ambiguous, please use customer ID", customerNameOrID)
	}

    return matchingCustomers[0], nil
}

func (c *GraphQLClient) GetCustomerLicense(appID string, customerID string) ([]byte, error) {
	return nil, errors.New("not implemented")
}
