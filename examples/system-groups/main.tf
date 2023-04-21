# Setup: User groups
resource "jumpcloud_user" "test_user1" {
  username   = "testuser1"
  firstname  = "test"
  lastname   = "user1"
  email      = "testuser1@testorg.org"
  enable_mfa = true
}

resource "jumpcloud_user_group_membership" "test_membership" {
  userid  = jumpcloud_user.test_user1.id
  groupid = jumpcloud_user_group.test_group.id
}

# Setup: System groups
resource "jumpcloud_user_group" "test_group" {
  name = "test_group"
  attributes = {
    posix_groups = "32:testerino"
  }
}

resource "jumpcloud_system_group" "test_group" {
  name = "Jumpcloud Provider Group"
}

resource "jumpcloud_system_group_user_group_membership" "test_group" {
  users_group_id   = jumpcloud_user_group.test_group.id
  systems_group_id = jumpcloud_system_group.test_group.jc_id
}
