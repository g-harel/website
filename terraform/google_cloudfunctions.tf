resource "google_cloudfunctions2_function" "build" {
  name        = "website-builder-${substr(base64encode(google_storage_bucket_object.build_function.md5hash), 0, 16)}"
  location    = "us-east1"
  description = "Website builder Cloud Function (Gen 2)"

  build_config {
    runtime     = "go126"
    entry_point = "Build"
    source {
      storage_source {
        bucket = google_storage_bucket.functions.name
        object = google_storage_bucket_object.build_function.name
      }
    }
  }

  service_config {
    max_instance_count = 3000
    available_memory   = "128Mi"
    timeout_seconds    = 30
    environment_variables = {
      CONFIG_BUCKET    = google_storage_bucket.functions.name
      CONFIG_OBJECT    = google_storage_bucket_object.build_config.name
      GRAPHQL_ENDPOINT = "https://api.github.com/graphql"
      GRAPHQL_TOKEN    = data.google_kms_secret.github_api_token.plaintext
      TEMPLATE_BUCKET  = google_storage_bucket.functions.name
      TEMPLATE_OBJECT  = google_storage_bucket_object.build_function_templates.name
      TEMPLATE_ENTRY   = "entry.html"
      UPLOAD_BUCKET    = google_storage_bucket.public_website.name
      UPLOAD_OBJECT    = "index.html"
    }
  }

  event_trigger {
    trigger_region = "us-east1"
    event_type     = "google.cloud.pubsub.topic.v1.messagePublished"
    pubsub_topic   = google_pubsub_topic.build_triggers.id
  }
}
