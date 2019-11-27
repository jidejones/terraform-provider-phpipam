package phpipam

import (
	"log"

	"github.com/lord-kyron/phpipam-sdk-go/controllers/addresses"
	"github.com/lord-kyron/phpipam-sdk-go/controllers/sections"
	"github.com/lord-kyron/phpipam-sdk-go/controllers/subnets"
	"github.com/lord-kyron/phpipam-sdk-go/controllers/vlans"
	"github.com/lord-kyron/phpipam-sdk-go/phpipam"
	"github.com/lord-kyron/phpipam-sdk-go/phpipam/session"
)

// Config provides the configuration for the PHPIPAM provider.
type Config struct {
	// The application ID required for API requests. This needs to be created
	// in the PHPIPAM console. It can also be supplied via the PHPIPAM_APP_ID
	// environment variable.
	AppID string

	// The API endpoint. This defaults to http://localhost/api, and can also be
	// supplied via the PHPIPAM_ENDPOINT_ADDR environment variable.
	Endpoint string

	// The password for the PHPIPAM account. This can also be supplied via the
	// PHPIPAM_PASSWORD environment variable.
	Password string

	// The user name for the PHPIPAM account. This can also be supplied via the
	// PHPIPAM_USER_NAME environment variable.
	Username string
}

// ProviderPHPIPAMClient is a structure that contains the client connections
// necessary to interface with the PHPIPAM API controllers. Example:
// subnets.Controller, or addresses.Controller.
type ProviderPHPIPAMClient struct {
	// The client for the addresses controller.
	addressesController *addresses.Controller

	// The client for the sections controller.
	sectionsController *sections.Controller

	// The client for the subnets controller.
	subnetsController *subnets.Controller

	// The client for the vlans controller.
	vlansController *vlans.Controller
}

// Client configures and returns a fully initialized PingdomClient.
func (c *Config) Client() (interface{}, error) {
	cfg := phpipam.Config{
		AppID:    c.AppID,
		Endpoint: c.Endpoint,
		Password: c.Password,
		Username: c.Username,
	}
	log.Printf("[DEBUG] Initializing PHPIPAM controllers")
	sess := session.NewSession(cfg)

	// Create the client object and return it
	client := ProviderPHPIPAMClient{
		addressesController: addresses.NewController(sess),
		sectionsController:  sections.NewController(sess),
		subnetsController:   subnets.NewController(sess),
		vlansController:     vlans.NewController(sess),
	}

	// Validate that our conneciton is okay
	if err := c.ValidateConnection(client.sectionsController); err != nil {
		return nil, err
	}

	return &client, nil
}

// ValidateConnection ensures that we can connect to PHPIPAM early, so that we
// do not fail in the middle of a TF run if it can be prevented.
func (c *Config) ValidateConnection(sc *sections.Controller) error {
	_, err := sc.ListSections()
	return err
}