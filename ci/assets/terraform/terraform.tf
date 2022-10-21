terraform {
  backend "oss" {}
  required_providers {
    alicloud = {
      source  = "xdops/alicloud"
    }
  }
}