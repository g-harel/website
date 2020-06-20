resource "cloudflare_zone" "primary" {
  zone = "harel.page"
}

resource "cloudflare_record" "g" {
  zone_id = cloudflare_zone.primary.id
  name    = "g"
  value   = "c.storage.googleapis.com"
  type    = "CNAME"
  proxied = true
}

resource "cloudflare_zone_settings_override" "primary" {
  zone_id = cloudflare_zone.primary.id
  settings {
    always_use_https  = "on"
    browser_cache_ttl = 30
  }
}
