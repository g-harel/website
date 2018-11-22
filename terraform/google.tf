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
    source_dir = ".temp/build"
    output_path = ".temp/build.zip"
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
    # name used to recreate resource when source changes
    name = "build-${substr(google_storage_bucket_object.build_function.md5hash, 0, 8)}"
    source_archive_bucket = "${google_storage_bucket.functions.name}"
    source_archive_object = "${google_storage_bucket_object.build_function.name}"

    region = "us-east1"
    runtime = "go111"
    available_memory_mb = 128
    entry_point = "Build"
    timeout = 30

    event_trigger = {
        event_type = "providers/cloud.pubsub/eventTypes/topic.publish"
        resource = "${google_pubsub_topic.build_triggers.name}"
    }

    environment_variables = {
        CONFIG_SRC = "https://raw.githubusercontent.com/g-harel/website/master/.config"
        TEMPLATE_DIR = "/srv/files/templates"
        TEMPLATE_ENTRY = "entry.html"
        GRAPHQL_ENDPOINT = "https://api.github.com/graphql"
        GRAPHQL_TOKEN = "${file(".secret/github.token")}"
        UPLOAD_BUCKET = "${google_storage_bucket.public_website.name}"
        UPLOAD_NAME = "index.html"
        GOOGLE_APPLICATION_CREDENTIALS = "/srv/files/function-service-account.json"
    }
}
