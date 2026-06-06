resource "cloudflare_zone" "primary" {
  name = "harel.page"
  account = {
    id = "187b22511a4652c53ac8facff126ae30"
  }
}

resource "cloudflare_dns_record" "g" {
  zone_id = cloudflare_zone.primary.id
  name    = "g.harel.page"
  content = "c.storage.googleapis.com"
  type    = "CNAME"
  proxied = true
  ttl     = 1
}

resource "cloudflare_zone_setting" "always_use_https" {
  zone_id    = cloudflare_zone.primary.id
  setting_id = "always_use_https"
  value      = "on"
}

resource "cloudflare_zone_setting" "browser_cache_ttl" {
  zone_id    = cloudflare_zone.primary.id
  setting_id = "browser_cache_ttl"
  value      = 30
}
