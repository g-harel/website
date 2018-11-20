resource "cloudflare_zone" "primary" {
    zone = "harel.page"
}

resource "cloudflare_record" "g" {
    domain = "${cloudflare_zone.primary.zone}"
    name = "g"
    value = "c.storage.googleapis.com"
    type = "CNAME"
    proxied = true
}
