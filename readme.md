# terraform-to-composition
Hack/ POC for generating XRs/XRDs from Terraform modules.

Terraform Module -> JSON Schema -> Go Structs -> go text/template -> XR/XRD YAML using OpenTofu Provider + function-patch-and-transform

### Resources
- https://github.com/HewlettPackard/terraschema
- https://docs.crossplane.io/v2.1/guides/function-patch-and-transform/#fromcompositefieldpath
- https://marketplace.upbound.io/providers/upbound/provider-opentofu/v1.0.3/resources/opentofu.m.upbound.io/Workspace/v1beta1

### Todo
- Handle module outputs
- cleanup/ wire CLI args
- output files to disk instead of stdout