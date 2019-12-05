package freeipa

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The hostname of the FreeIPA master to use",
			},

			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "admin",
				Description: "The username to use to authenticate to the FreeIPA master",
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The password to use to authenticate to the FreeIPA master",
				Sensitive:   true,
			},

			"base_dn": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The LDAP BaseDN",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"freeipa_user": resourceUser(),
			"freeipa_group": resourceGroup(),
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Host:       d.Get("host").(string),
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
		BaseDN:     d.Get("base_dn").(string),
	}

	return config.NewConnection()
}
