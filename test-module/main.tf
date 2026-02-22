resource "random_pet" "suffix" {
  keepers = {
    bucket_name = var.bucket_name
  }
}

resource "aws_s3_bucket" "this" {
  bucket = "${var.bucket_name}-${random_pet.suffix.id}"
  region = var.region
}