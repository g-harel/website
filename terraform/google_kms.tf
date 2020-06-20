data "google_kms_key_ring" "master" {
  name     = "master-222818"
  location = "northamerica-northeast1"
}

data "google_kms_crypto_key" "master" {
  name     = "master-222818"
  key_ring = "${data.google_kms_key_ring.master.id}"
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
  crypto_key = "${data.google_kms_crypto_key.master.self_link}"
  ciphertext = "CiQAfiTmjjKuiVMOWmquTAA4NxJcmpWYiLtBaZoxvs2BBs+WqlgSUQDNbPzaRyz3TpPBhZoH0APDJZSPpeogk4dWg377d13civeUOv+2vqANY/vDIp4eXoCBdQ7TBysD70gF4bo7gPnuaXZV9nc1C5gevTFy2sBruw=="
}

data "google_kms_secret" "cloudflare_key" {
  crypto_key = "${data.google_kms_crypto_key.master.self_link}"
  ciphertext = "CiQAfiTmjh82wX+zrkaJbdfRRmWcGM0NBCqoduJZG7DSbu+o0FkSTgDNbPzalo6xAgjobovm/EWqc7RIK27q0aRbY2MfvpygFvduYeQeMUVoGM5yvR+hp355YW/P0/mKIfvouwHi3MjVp1JhKOrA972wB/gFmQ=="
}
