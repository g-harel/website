resource "google_storage_bucket" "public_website" {
  name          = "g.harel.page"
  storage_class = "MULTI_REGIONAL"

  website {
    main_page_suffix = "index.html"
  }
}

resource "google_storage_default_object_acl" "public_website" {
  bucket = google_storage_bucket.public_website.name

  role_entity = [
    "READER:allUsers",
  ]
}

//

resource "google_storage_bucket" "functions" {
  name = "functions-222818"
}

resource "google_storage_bucket_object" "build_config" {
  bucket = google_storage_bucket.functions.name
  name   = ".config"
  source = "../.config"
}

resource "google_storage_bucket_object" "build_function" {
  bucket = google_storage_bucket.functions.name
  # Change name to force update function.
  # https://github.com/hashicorp/terraform-provider-google/issues/1938
  name= format("build.zip#%s", data.archive_file.build_function.output_md5)
  source = data.archive_file.build_function.output_path
}

resource "google_storage_bucket_object" "build_function_templates" {
  bucket = google_storage_bucket.functions.name
  name   = "templates.zip"
  source = data.archive_file.build_function_templates.output_path
}
