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
    source_dir = "../functions"
    output_path = ".temp/build.zip"
}

data "archive_file" "build_function_templates" {
    type = "zip"
    source_dir = "../templates"
    output_path = ".temp/templates.zip"
}

resource "google_storage_bucket" "functions" {
    name = "functions-222818"
}

resource "google_storage_bucket_object" "build_function" {
    bucket = "${google_storage_bucket.functions.name}"
    name   = "build.zip"
    source = "${data.archive_file.build_function.output_path}"
}

resource "google_storage_bucket_object" "build_function_templates" {
    bucket = "${google_storage_bucket.functions.name}"
    name   = "templates.zip"
    source = "${data.archive_file.build_function_templates.output_path}"
}

resource "google_cloudfunctions_function" "build" {
    # name used to recreate resource when source changes
    name = "build-${substr(base64encode(google_storage_bucket_object.build_function.md5hash), 0, 8)}"
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
        GRAPHQL_ENDPOINT = "https://api.github.com/graphql"
        GRAPHQL_TOKEN = "${file(".secret/github.token")}"
        TEMPLATE_BUCKET = "${google_storage_bucket.functions.name}"
        TEMPLATE_OBJECT = "${google_storage_bucket_object.build_function_templates.name}"
        TEMPLATE_ENTRY = "entry.html"
        UPLOAD_BUCKET = "${google_storage_bucket.public_website.name}"
        UPLOAD_OBJECT = "index.html"
    }
}
