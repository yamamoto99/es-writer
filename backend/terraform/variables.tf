variable "aws_access_key" {
	description = "AWS Access Key"
}

variable "aws_secret_key" {
	description = "AWS Secret Access Key"
}

variable "aws_region" {
	description = "AWS Region"
}

variable "availability_zone_1" {
	description = "The availability zone"
}

variable "availability_zone_2" {
	description = "The availability zone"
}

variable "rds_pass" {
	description = "RDS instance password"
}

variable "rds_username" {
	description = "RDS instance username"
}

variable "key_name" {
	description = "key_name"
}

variable "pub_key" {
	description = "pub_key"
}

variable "https_cert_arn" {
	description = "https_cert_arn"
}

variable "domain_name" {
	description = "The domain name for Route 53"
}

variable "subdomain_name" {
	description = "The subdomain name for Route 53"
}
