terraform {
    backend "gcs" {
        credentials = ".secret/terraform-service-account.json"
        project = "website-222818"
        bucket = "tf-state-222818"
        prefix  = "terraform/state"
    }
}

provider "archive" {}

provider "cloudflare" {
    email = "gabrielj.harel@gmail.com"
    token = "${file(".secret/cloudflare.key")}"
}

provider "google" {
    credentials = ".secret/terraform-service-account.json"
    project = "website-222818"
    region = "northamerica-northeast1"
}
