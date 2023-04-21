locals {
  username_lowercase = "${lower(var.firstname)}.${lower(var.lastname)}"
}

resource "jumpcloud_user" "user_from_module" {
  username  = local.username_lowercase
  firstname = var.firstname
  lastname  = var.lastname
  email     = "${local.username_lowercase}@rockstar.org"

  # Even rockstars need MFA
  enable_mfa = true
}

resource "jumpcloud_user_group_membership" "band" {
  userid  = jumpcloud_user.user_from_module.id
  groupid = var.band_id
}

resource "jumpcloud_user_group_membership" "position" {
  userid  = jumpcloud_user.user_from_module.id
  groupid = var.position_id
}
