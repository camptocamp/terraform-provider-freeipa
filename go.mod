module github.com/camptocamp/terraform-provider-freeipa

go 1.12

replace github.com/tehwalris/go-freeipa => github.com/camptocamp/go-freeipa v0.0.0-20200218165340-4d2a1083fe5f

require (
	github.com/hashicorp/terraform-plugin-sdk v1.7.0
	github.com/tehwalris/go-freeipa v0.0.0-20191205131509-4409391e6e02
)
