---
subcategory: "Compute Engine"
description: |-
  Get a network within GCE.
---

# google_compute_network

Get a network within GCE from its name or self_link.

## Example Usage

```tf
# Look up a network by name (in the provider's project, or in `project` if set):
data "google_compute_network" "my-network" {
  name = "default-us-east1"
}

# Look up a network by self_link (no need to specify project — it is parsed
# from the link). Useful for cross-project references coming from another
# resource's output:
data "google_compute_network" "shared-vpc" {
  self_link = "https://www.googleapis.com/compute/v1/projects/host-project/global/networks/shared"
}
```

## Argument Reference

The following arguments are supported. Exactly one of `name` or `self_link` must be provided:

* `name` - (Optional) The name of the network. Conflicts with `self_link`.

* `self_link` - (Optional) The full or partial URL of the network. May be one of:
  * a fully-qualified URL (`https://www.googleapis.com/compute/v1/projects/{{project}}/global/networks/{{name}}`),
  * a relative path (`projects/{{project}}/global/networks/{{name}}`),
  * a short relative path (`global/networks/{{name}}`),
  * a name (in which case `project` and `self_link` are equivalent).

  When `self_link` is given, `project` is parsed from the link and the
  data source's `project` argument is ignored. Conflicts with `name`.


- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used. Ignored when `self_link` is set.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - an identifier for the resource with format projects/{{project}}/global/networks/{{name}}

* `description` - Description of this network.

* `network_id` - The numeric unique identifier for the resource.

* `numeric_id` - (Deprecated) The numeric unique identifier for the resource. `numeric_id` is deprecated and will be removed in a future major release. Use `network_id` instead.

* `gateway_ipv4` - The IP address of the gateway.

* `internal_ipv6_range` - The ula internal ipv6 range assigned to this network.

* `network_profile` - A full or partial URL of the network profile to apply to this network.

* `subnetworks_self_links` - the list of subnetworks which belong to the network

* `self_link` - The URI of the resource.
