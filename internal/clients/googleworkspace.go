/*
Copyright 2021 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package clients

import (
	"context"
	"encoding/json"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/terrajet/pkg/terraform"

	"github.com/joakimhew/provider-jet-googleworkspace/apis/v1alpha1"
)

const (
	// error messages
	errNoProviderConfig     = "no providerConfigRef provided"
	errGetProviderConfig    = "cannot get referenced ProviderConfig"
	errTrackUsage           = "cannot track ProviderConfig usage"
	errExtractCredentials   = "cannot extract credentials"
	errUnmarshalCredentials = "cannot unmarshal googleworkspace credentials as JSON"

	keyCredentials           = "credentials"
	keyCustomerId            = "customer_id"
	keyOauthScopes           = "oauth_scopes"
	keyImpersonatedUserEmail = "impersonated_user_email"
)

// TerraformSetupBuilder builds Terraform a terraform.SetupFn function which
// returns Terraform provider setup configuration
func TerraformSetupBuilder(version, providerSource, providerVersion string) terraform.SetupFn {
	return func(ctx context.Context, client client.Client, mg resource.Managed) (terraform.Setup, error) {
		ps := terraform.Setup{
			Version: version,
			Requirement: terraform.ProviderRequirement{
				Source:  providerSource,
				Version: providerVersion,
			},
		}

		configRef := mg.GetProviderConfigReference()
		if configRef == nil {
			return ps, errors.New(errNoProviderConfig)
		}
		pc := &v1alpha1.ProviderConfig{}
		if err := client.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc); err != nil {
			return ps, errors.Wrap(err, errGetProviderConfig)
		}

		t := resource.NewProviderConfigUsageTracker(client, &v1alpha1.ProviderConfigUsage{})
		if err := t.Track(ctx, mg); err != nil {
			return ps, errors.Wrap(err, errTrackUsage)
		}

		data, err := resource.CommonCredentialExtractor(ctx, pc.Spec.Credentials.Source, client, pc.Spec.Credentials.CommonCredentialSelectors)
		if err != nil {
			return ps, errors.Wrap(err, errExtractCredentials)
		}
		googleworkspaceCreds := map[string]string{}
		if err := json.Unmarshal(data, &googleworkspaceCreds); err != nil {
			return ps, errors.Wrap(err, errUnmarshalCredentials)
		}

		// set environment variables for sensitive provider configuration
		// Deprecated: In shared gRPC mode we do not support injecting
		// credentials via the environment variables. You should specify
		// credentials via the Terraform main.tf.json instead.
		/*ps.Env = []string{
			fmt.Sprintf("%s=%s", "HASHICUPS_USERNAME", googleworkspaceCreds["username"]),
			fmt.Sprintf("%s=%s", "HASHICUPS_PASSWORD", googleworkspaceCreds["password"]),
		}*/
		// set credentials in Terraform provider configuration
		/*ps.Configuration = map[string]interface{}{
			"username": googleworkspaceCreds["username"],
			"password": googleworkspaceCreds["password"],
		}*/

		ps.Configuration = map[string]interface{}{}
		if v, ok := googleworkspaceCreds[keyCredentials]; ok {
			ps.Configuration[keyCredentials] = v
		}
		if v, ok := googleworkspaceCreds[keyCustomerId]; ok {
			ps.Configuration[keyCustomerId] = v
		}
		if v, ok := googleworkspaceCreds[keyOauthScopes]; ok {
			list := []string{}
			if err := json.Unmarshal([]byte(v), &list); err != nil {
				return ps, errors.Wrap(err, errUnmarshalCredentials)
			}
			ps.Configuration[keyOauthScopes] = list
		}
		if v, ok := googleworkspaceCreds[keyImpersonatedUserEmail]; ok {
			ps.Configuration[keyImpersonatedUserEmail] = v
		}
		return ps, nil
	}
}
