terraform {
  required_providers {
    statefile_rce = {
      source  = "offensive-actions/statefile-rce"
      version = "0.1.0"
    }
  }
}

provider "statefile_rce" {}
