# Comment

locals {
  cfg = { for k, v in local.foo :
    k => {
      v.test[*].test[0].name = "bar"
    }
  }
}
