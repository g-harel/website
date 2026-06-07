resource "google_cloudbuild_trigger" "master_branch" {
  name     = "master-branch-trigger"
  location = "northamerica-northeast1"

  github {
    owner = "g-harel"
    name  = "website"
    push {
      branch = "^master$"
    }
  }

  build {
    step {
      name = "hashicorp/terraform:1.15.5"
      dir  = "terraform"
      args = ["init"]
    }
    step {
      name = "hashicorp/terraform:1.15.5"
      dir  = "terraform"
      args = ["apply", "-auto-approve", "-lock-timeout=10m"]
    }
  }
}
