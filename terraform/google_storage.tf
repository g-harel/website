resource "google_storage_bucket" "public_website" {
    name = "g.harel.page"
    location = "northamerica-northeast1"
    storage_class = "REGIONAL"

    website = {
        main_page_suffix = "index.html"
    }
}

resource "google_storage_default_object_acl" "image-store-default-acl" {
    bucket = "${google_storage_bucket.public_website.name}"

    role_entity = [
        "READER:allUsers",
    ]
}