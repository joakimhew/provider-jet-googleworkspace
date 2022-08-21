package group

import (
	"github.com/crossplane/terrajet/pkg/config"
)

// Configure the googleworkspace_group resource
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("googleworkspace_group", func(r *config.Resource) {
		r.ShortGroup = "group"
		r.ExternalName = config.IdentifierFromProvider
		r.UseAsync = true
	})
}
