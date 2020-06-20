resource "google_cloudfunctions_function" "build" {
  # name used to recreate resource when source changes
  name                  = "build-${substr(base64encode(google_storage_bucket_object.build_function.md5hash), 0, 8)}"
  source_archive_bucket = google_storage_bucket.functions.name
  source_archive_object = google_storage_bucket_object.build_function.name

  region              = "us-east1"
  runtime             = "go111"
  available_memory_mb = 128
  entry_point         = "Build"
  timeout             = 30

  event_trigger {
    event_type = "providers/cloud.pubsub/eventTypes/topic.publish"
    resource   = google_pubsub_topic.build_triggers.name
  }

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
