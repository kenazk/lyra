package handler

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/lyraproj/lyra/cmd/goplugin-tf-kubernetes/generated"
	"github.com/lyraproj/puppet-evaluator/eval"
	"github.com/lyraproj/servicesdk/grpc"
	"github.com/lyraproj/servicesdk/service"
	"github.com/terraform-providers/terraform-provider-kubernetes/kubernetes"
)

// Server configures the Terraform provider and creates an instance of the server
func Server(c eval.Context, configuredProvider *schema.Provider) *service.Server {
	sb := service.NewServerBuilder(c, "TerraformKubernetes")
	generated.Initialize(sb, configuredProvider)
	return sb.Server()
}

// Start this server running
func Start() {
	config := &terraform.ResourceConfig{
		Config: map[string]interface{}{
			// "foo": "bar",
		},
	}
	p := kubernetes.Provider().(*schema.Provider)
	err := p.Configure(config)
	if err != nil {
		panic(err)
	}
	eval.Puppet.Do(func(c eval.Context) {
		grpc.Serve(c, Server(c, p))
	})
}
