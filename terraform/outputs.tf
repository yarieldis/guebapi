output "api_url" {
  description = "Invoke URL of the API Gateway HTTP API"
  value       = aws_apigatewayv2_stage.default.invoke_url
}

output "dynamodb_table_name" {
  description = "Name of the DynamoDB users table"
  value       = aws_dynamodb_table.users.name
}
