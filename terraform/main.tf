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
    token = "${data.google_kms_secret.cloudflare_key.plaintext}"
}

provider "google" {
    credentials = ".secret/terraform-service-account.json"
    project = "website-222818"
    region = "northamerica-northeast1"
}
