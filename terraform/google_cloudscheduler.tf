resource "google_cloud_scheduler_job" "trigger" {
    provider = "google-beta"

    name     = "build-trigger"
    schedule = "*/5 * * * *"

    pubsub_target = {
        topic_name = "${google_pubsub_topic.build_triggers.name}"
        data       = "${base64encode("{}")}"
    }
}
