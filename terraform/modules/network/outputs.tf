output "subnet_ids" {
  description = "IDs of the default VPC subnets"
  value       = data.aws_subnets.default.ids
}

output "security_group_id" {
  description = "Security group ID for ECS"
  value       = aws_security_group.this.id
}

output "vpc_id" {
  description = "The ID of the default VPC"
  value       = data.aws_vpc.default.id
}

output "alb_sg_id" {
  value = aws_security_group.alb_sg.id
}