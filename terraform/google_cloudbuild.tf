resource "google_cloudbuild_trigger" "master_branch" {
  trigger_template {
    project_id  = "${data.google_project.project.project_id}"
    repo_name   = "${google_sourcerepo_repository.website.name}"
    branch_name = "master"
    dir         = "terraform"
  }

  build {
    step {
      name = "${data.google_container_registry_image.terraform_build_step.image_url}"
      args = ["init"]
    }
    step {
      name = "${data.google_container_registry_image.terraform_build_step.image_url}"
      args = ["apply", "-auto-approve"]
    }
  }
}
