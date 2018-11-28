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

resource "cloudflare_zone_settings_override" "primary" {
    name = "${cloudflare_zone.primary.zone}"
    settings = {
        always_use_https = "on"
        browser_cache_ttl = 30
    }
}
