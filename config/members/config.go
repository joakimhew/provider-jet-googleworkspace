package members

import (
	"github.com/crossplane/terrajet/pkg/config"
)

func Configure(p *config.Provider) {
	p.AddResourceConfigurator("googleworkspace_group_members", func(r *config.Resource) {
		r.ShortGroup = "group"
		r.ExternalName = config.IdentifierFromProvider
	})
}
