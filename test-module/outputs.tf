output "bucket_name" {
  value = aws_s3_bucket.this.bucket
  description = "test desc for bucket_name"
}

output "other_test" {
  value = "whatever"
  sensitive = true
}