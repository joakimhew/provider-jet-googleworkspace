package members

import (
	"github.com/crossplane/terrajet/pkg/config"
)

// Configure the googleworkspace_group_members resource
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("googleworkspace_group_members", func(r *config.Resource) {
		r.ShortGroup = "group"
		r.ExternalName = config.IdentifierFromProvider
		r.UseAsync = true
		r.References["group_id"] = config.Reference{
			Type: "Group",
		}
	})
}
