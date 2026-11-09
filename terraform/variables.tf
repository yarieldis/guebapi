variable "aws_region" {
  description = "AWS region to deploy resources into"
  type        = string
  default     = "us-east-1"
}

variable "service_name" {
  description = "Base name used for AWS resources"
  type        = string
  default     = "guebapi"
}

variable "lambda_zip_path" {
  description = "Path to the zipped Lambda deployment package"
  type        = string
  default     = "lambda.zip"
}

variable "dynamodb_table_name" {
  description = "DynamoDB table name for user storage"
  type        = string
  default     = "users"
}

variable "jwt_secret_key" {
  description = "Secret key used to sign JWTs"
  type        = string
  sensitive   = true
}

variable "jwt_token_duration" {
  description = "JWT token validity duration"
  type        = string
  default     = "72h"
}
