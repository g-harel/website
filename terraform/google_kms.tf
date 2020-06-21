resource "google_kms_key_ring" "secrets" {
  name     = "secrets"
  location = "northamerica-northeast1"
}

resource "google_kms_crypto_key" "secrets" {
  name     = "secrets"
  key_ring = google_kms_key_ring.secrets.self_link

  lifecycle {
    prevent_destroy = true
  }
}

/* Use `gcloud kms encrypt` to regenerate ciphertext

echo -n "my_secret_value" | gcloud kms encrypt \
    --project=website-222818 \
    --location=northamerica-northeast1 \
    --keyring=secrets \
    --key=secrets \
    --plaintext-file - \
    --ciphertext-file - \
    | base64 -w 0

*/

data "google_kms_secret" "github_api_token" {
  crypto_key = google_kms_crypto_key.secrets.self_link
  ciphertext = "CiQA3sn5Bw/1CZ2A0GFd3pKAKoZSIC1zjk1+leIjOc9sLOduWRoSUgA21SmU5cO7cd2VGcXUkaGptOalm+ILXOcFiur6loOk9cTZzzjE6T0K4EvlkXWPgY3E8oSaGs+R+mzx0mE8T8XqXx0qNZpTnKpWAFFABbvumhs="
}

data "google_kms_secret" "cloudflare_token" {
  crypto_key = google_kms_crypto_key.secrets.self_link
  ciphertext = "CiQA3sn5Bwm7o6Y9VeaMAbZMnNH8SwCLXZ95/IeRvVA7oOj6JAoSUAA21SmUgg30UfIThOJu+vK93zRcz/5tUP7jEwPv5BUpDp1a9u/dJwJ6zW5aE2ekQfPoZ++rkAnj0ErvNHpYUfKoTcdWPJglmN6L2eCVpKx3"
}
