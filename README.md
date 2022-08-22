# Terrajet GoogleWorkspace Provider

`provider-jet-googleworkspace` is a [Crossplane](https://crossplane.io/) provider that
is built using [Terrajet](https://github.com/crossplane/terrajet) code
generation tools and exposes XRM-conformant managed resources for the
GoogleWorkspace API.

## Getting Started

Install the provider by using the following command after changing the image tag
to the [latest release](https://github.com/joakimhew/provider-jet-googleworkspace/releases):

```
kubectl crossplane install provider joakimhew/provider-jet-googleworkspace:v0.1.0-alpha
```

Alternatively, you can use declarative installation:

```
kubectl apply -f examples/install.yaml
```

Notice that in this example Provider resource is referencing ControllerConfig with debug enabled.

You can see the API reference [here](https://doc.crds.dev/github.com/joakimhew/provider-jet-googleworkspace).

Create a secret with the credentials for your Google Workspace account:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: example-creds
  namespace: crossplane-system
type: Opaque
stringData:
  credentials: |
    {
      "credentials":"<GOOGLE_CREDENTIALS_JSON>",
      "customer_id":"<CUSTOMER_ID>",
      "impersonated_user_email":"<IMPERSONATED_USER_EMAIL>",
      "oauth_scopes": "<LIST_OF_OAUTH_SCOPES>"
    }
```

Create a `ProviderConfig` that uses the credentials:

```yaml
apiVersion: googleworkspace.jet.crossplane.io/v1alpha1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      name: example-creds
      namespace: crossplane-system
      key: credentials
```

## Developing

Run code-generation pipeline:

```console
go run cmd/generator/main.go "$PWD"
```

Run against a Kubernetes cluster:

```console
make run
```

Build, push, and install:

```console
make all
```

Build binary:

```console
make build
```

## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please
open an [issue](https://github.com/joakimhew/provider-jet-googleworkspace/issues).

## Contact

Please use the following to reach members of the community:

- Slack: Join our [slack channel](https://slack.crossplane.io)
- Forums:
  [crossplane-dev](https://groups.google.com/forum/#!forum/crossplane-dev)
- Twitter: [@crossplane_io](https://twitter.com/crossplane_io)
- Email: [info@crossplane.io](mailto:info@crossplane.io)

## Governance and Owners

provider-jet-googleworkspace is run according to the same
[Governance](https://github.com/crossplane/crossplane/blob/master/GOVERNANCE.md)
and [Ownership](https://github.com/crossplane/crossplane/blob/master/OWNERS.md)
structure as the core Crossplane project.

## Code of Conduct

provider-jet-googleworkspace adheres to the same [Code of
Conduct](https://github.com/crossplane/crossplane/blob/master/CODE_OF_CONDUCT.md)
as the core Crossplane project.

## Licensing

provider-jet-googleworkspace is under the Apache 2.0 license.
