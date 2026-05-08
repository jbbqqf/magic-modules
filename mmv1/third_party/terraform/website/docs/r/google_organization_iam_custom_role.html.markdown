---
subcategory: "Cloud Platform"
description: |-
 Allows management of a customized Cloud IAM organization role.
---

# google_organization_iam_custom_role

Allows management of a customized Cloud IAM organization role. For more information see
[the official documentation](https://cloud.google.com/iam/docs/understanding-custom-roles)
and
[API](https://cloud.google.com/iam/reference/rest/v1/organizations.roles).

~> **Warning:** Note that custom roles in GCP have the concept of a soft-delete. There are two issues that may arise
 from this and how roles are propagated. 1) creating a role may involve undeleting and then updating a role with the
 same name, possibly causing confusing behavior between undelete and update. 2) A deleted role is permanently deleted
 after 7 days, but it can take up to 30 more days (i.e. between 7 and 37 days after deletion) before the role name is
 made available again. This means a deleted role that has been deleted for more than 7 days cannot be changed at all
 by Terraform, and new roles cannot share that name.
 
## Example Usage

This snippet creates a customized IAM organization role.

```hcl
resource "google_organization_iam_custom_role" "my-custom-role" {
  role_id     = "myCustomRole"
  org_id      = "123456789"
  title       = "My Custom Role"
  description = "A description"
  permissions = ["iam.roles.list", "iam.roles.create", "iam.roles.delete"]
}
```

## Argument Reference

The following arguments are supported:

* `role_id` - (Optional) The role id to use for this role. Either `role_id` or `role_id_prefix` must be set, but not both. If neither is set, a unique role_id beginning with `tf_` is generated.

* `role_id_prefix` - (Optional) Creates a unique `role_id` beginning with the specified prefix. Conflicts with `role_id`. Must satisfy the `role_id` regex (alphanumeric, underscores and dots, 3-64 chars total including the generated suffix). Useful when several roles need to share a common prefix without collisions.

* `org_id` - (Required) The numeric ID of the organization in which you want to create a custom role.

* `title` - (Required) A human-readable title for the role.

* `permissions` (Required) The names of the permissions this role grants when bound in an IAM policy. At least one permission must be specified.

* `stage` - (Optional) The current launch stage of the role.
    Defaults to `GA`.
    List of possible stages is [here](https://cloud.google.com/iam/reference/rest/v1/organizations.roles#Role.RoleLaunchStage).

* `description` - (Optional) A human-readable description for the role.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `deleted` - (Optional) The current deleted state of the role.

* `id` - an identifier for the resource with the format `organizations/{{org_id}}/roles/{{role_id}}`

* `name` - The name of the role in the format `organizations/{{org_id}}/roles/{{role_id}}`. Like `id`, this field can be used as a reference in other resources such as IAM role bindings.

## Import

Customized IAM organization role can be imported using their URI, e.g.

```
$ terraform import google_organization_iam_custom_role.my-custom-role organizations/123456789/roles/myCustomRole
```
