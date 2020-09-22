resource "jumpcloud_system_group" "test_group" {
  name = "Jumpcloud Provider Group"
}

resource "jumpcloud_system_group_user_group_membership" "test_group" {
    users_group_id = jumpcloud_user_group.test_group.id
    systems_group_id = jumpcloud_system_group.test_group.jc_id
}
