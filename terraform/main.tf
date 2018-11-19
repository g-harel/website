terraform {
    backend "gcs" {
        credentials = "service-account.json"
        project = "website-222818"
        bucket = "tf-state-222818"
        prefix  = "terraform/state"
    }
}

provider "google" {
    credentials = "service-account.json"
    project = "website-222818"
    region = "northamerica-northeast1"
}
