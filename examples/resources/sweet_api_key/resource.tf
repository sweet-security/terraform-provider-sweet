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

resource "sweet_api_key" "example" {
  description = "key for sensor installation with only ingestion permissions"
}
