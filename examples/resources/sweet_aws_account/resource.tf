terraform {
  required_providers {
    sweet = {
      source = "sweet-security/sweet"
    }
  }
}

provider "sweet" {
  api_key = "00000000-0000-0000-0000-000000000000"
  secret = "00000000-0000-0000-0000-000000000000"
}

resource "sweet_aws_account" "example" {
  account_id = "123456789012"
  role_arn = "arn:aws:iam::123456789012:role/example"
}
