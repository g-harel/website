resource "google_cloud_scheduler_job" "trigger" {
  name     = "build-trigger"
  schedule = "*/5 * * * *"
  region   = "us-central1"

  pubsub_target {
    topic_name = "projects/${google_pubsub_topic.build_triggers.project}/topics/${google_pubsub_topic.build_triggers.name}"
    data       = base64encode("{}")
  }
}
