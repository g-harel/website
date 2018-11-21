resource "google_storage_bucket" "public_website" {
    name = "g.harel.page"
    storage_class = "MULTI_REGIONAL"

    website = {
        main_page_suffix = "index.html"
    }
}

resource "google_storage_default_object_acl" "public_website" {
    bucket = "${google_storage_bucket.public_website.name}"

    role_entity = [
        "READER:allUsers",
    ]
}

resource "google_pubsub_topic" "build_triggers" {
  name = "build-triggers"
}

data "archive_file" "build_function" {
    type = "zip"
    source_dir = "${path.module}/../function"
    output_path = "${path.module}/.functions/build.zip"
}

resource "google_storage_bucket" "functions" {
    name = "functions-222818"
}

resource "google_storage_bucket_object" "build_function" {
  name   = "build.zip"
  source = "${data.archive_file.build_function.output_path}"
  bucket = "${google_storage_bucket.functions.name}"
}

resource "google_cloudfunctions_function" "build" {
    name = "build"
    available_memory_mb = 128
    region = "us-east1"

    source_archive_bucket = "${google_storage_bucket.functions.name}"
    source_archive_object = "${google_storage_bucket_object.build_function.name}"

    event_trigger = {
        event_type = "providers/cloud.pubsub/eventTypes/topic.publish"
        resource = "${google_pubsub_topic.build_triggers.name}"
    }
}
