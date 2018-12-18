resource "jumpcloud_user_group" "language" {
    name = "language"
    attributes {}
}

resource "jumpcloud_user_group" "cars" {
    name = "cars"
}

resource "jumpcloud_user_group" "foo" {
    name = "barz"
    attributes {
      posix_groups = "32:testerino"
    }
}
