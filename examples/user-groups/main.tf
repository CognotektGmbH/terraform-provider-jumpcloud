resource "jumpcloud_user_group" "example_resource" {
  name = "Examplegroup"

  attributes {
    posix_groups = "32:test"
  }
}
