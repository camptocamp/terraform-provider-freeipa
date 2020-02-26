module github.com/camptocamp/terraform-provider-freeipa

go 1.12

replace github.com/tehwalris/go-freeipa => github.com/camptocamp/go-freeipa v0.0.0-20200226102930-f1a87e4d06ba

require (
	github.com/hashicorp/terraform-plugin-sdk v1.7.0
	github.com/tehwalris/go-freeipa v0.0.0-20191205131509-4409391e6e02
)
