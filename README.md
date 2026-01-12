# terraform-provider-sweet

Terraform provider for interacting with Sweet’s API

---

## Usage

```terraform
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

resource "sweet_aws_organization" "example" {
  account_id = "123456789012"
  role_arn = "arn:aws:iam::123456789012:role/example"
  role_name_parameter_arn = "arn:aws:iam::123456789012:role/example"
}


resource "sweet_aws_account" "example" {
  account_id = "123456789012"
  role_arn = "arn:aws:iam::123456789012:role/example"
}
```
Obtain tenant api key & secret from [Sweet platform](https://app.sweet.security) inside settings -> API Tokens
