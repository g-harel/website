terraform {
  backend "gcs" {
    bucket = "tf-state-222818"
    prefix = "terraform/state"
  }
}

provider "archive" {}

provider "cloudflare" {
  email = "gabrielj.harel@gmail.com"
  api_token = data.google_kms_secret.cloudflare_token.plaintext
}

provider "google" {
  project = "website-222818"
  region  = "northamerica-northeast1"
}

data "google_project" "project" {}
