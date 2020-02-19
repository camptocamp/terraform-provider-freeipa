module github.com/camptocamp/terraform-provider-freeipa

go 1.12

replace github.com/tehwalris/go-freeipa => github.com/camptocamp/go-freeipa v0.0.0-20200218165340-4d2a1083fe5f

require (
	github.com/hashicorp/terraform v0.12.7
	github.com/tehwalris/go-freeipa v0.0.0-20191205131509-4409391e6e02
	github.com/tustvold/kerby v0.0.0-20171212165806-ca8245542ec8 // indirect
	gopkg.in/ldap.v2 v2.5.1 // indirect
)
