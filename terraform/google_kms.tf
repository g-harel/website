resource "google_kms_key_ring" "master" {
  name     = "master-222818"
  location = "northamerica-northeast1"
}

resource "google_kms_crypto_key" "master" {
  name     = "master-222818"
  key_ring = google_kms_key_ring.master.self_link

  lifecycle {
    prevent_destroy = true
  }
}

/* Use `gcloud kms encrypt` to regenerate ciphertext

echo -n "my_secret_value" | gcloud kms encrypt \
    --project=website-222818 \
    --location=northamerica-northeast1 \
    --keyring=master-222818 \
    --key=master-222818 \
    --plaintext-file - \
    --ciphertext-file - \
    | base64 -w 0

*/

data "google_kms_secret" "github_api_token" {
  crypto_key = google_kms_crypto_key.master.self_link
  ciphertext = "CiQAfiTmjjKuiVMOWmquTAA4NxJcmpWYiLtBaZoxvs2BBs+WqlgSUQDNbPzaRyz3TpPBhZoH0APDJZSPpeogk4dWg377d13civeUOv+2vqANY/vDIp4eXoCBdQ7TBysD70gF4bo7gPnuaXZV9nc1C5gevTFy2sBruw=="
}

data "google_kms_secret" "cloudflare_key" {
  crypto_key = google_kms_crypto_key.master.self_link
  ciphertext = "CiQAfiTmjrSVoMnTwgR3JSi8yBrNDmPUF8gx1YGAGzOt9AwJBYcSUQBhQX3zHT0chIjlxgwUOLCyoCDXUmmfS0Yp1jdugh7JI/kaWQv2GxjxmjrXNmbfYBPRyAs4ozz8GozMBiwDye9bmYEa1v9E39FDHK96Dd8Yfw=="
}
