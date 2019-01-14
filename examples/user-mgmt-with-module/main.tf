resource "jumpcloud_user_group" "pink_floyd" {
  name = "PinkFloyd"
}

resource "jumpcloud_user_group" "queen" {
  name = "Queen"
}

resource "jumpcloud_user_group" "guitarist" {
  name = "guitarist"
}

resource "jumpcloud_user_group" "drummer" {
  name = "drummer"
}

module "brian_may" {
  source      = "./module"
  firstname   = "Brian"
  lastname    = "May"
  band_id     = "${jumpcloud_user_group.queen.id}"
  position_id = "${jumpcloud_user_group.guitarist.id}"
}

module "nick_mason" {
  source      = "./module"
  firstname   = "Nick"
  lastname    = "Mason"
  band_id     = "${jumpcloud_user_group.pink_floyd.id}"
  position_id = "${jumpcloud_user_group.drummer.id}"
}
