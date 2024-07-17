# ====================
# VPC
# ====================
resource "aws_vpc" "es-writer" {
	cidr_block           = "10.0.0.0/16"
	instance_tenancy     = "default"
	enable_dns_support   = "true"
	enable_dns_hostnames = "true"
	tags = {
		Name = "es-writer-VPC"
	}
}

# ====================
# Subnet
# ====================
# パブリックサブネット
resource "aws_subnet" "public_app" {
	vpc_id                  = "${aws_vpc.es-writer.id}"
	cidr_block              = "10.0.1.0/24"
	map_public_ip_on_launch = true
	availability_zone       = "us-west-2a"
	tags = {
		Name = "public-web"
	}
}

# プライベートサブネット
resource "aws_subnet" "private_db_1" {
	vpc_id                  = "${aws_vpc.es-writer.id}"
	cidr_block              = "10.0.2.0/24"
	map_public_ip_on_launch = false
	availability_zone       = "us-west-2a"
	tags = {
		Name = "private-db-1"
	}
}

# プライベートサブネット
resource "aws_subnet" "private_db_2" {
	vpc_id                  = "${aws_vpc.es-writer.id}"
	cidr_block              = "10.0.3.0/24"
	map_public_ip_on_launch = false
	availability_zone       = "us-west-2c"
	tags = {
		Name = "private-db-2"
	}
}

# ====================
# Security Group
# ====================
# WEBサーバーのセキュリティグループ
resource "aws_security_group" "es-writer-web-sg" {
	vpc_id      = "${aws_vpc.es-writer.id}"
	name        = "es-writer-web-sg"
	description = "es-writer-web-sg"
	tags = {
		Name = "es-writer-web-sg"
	}

	ingress {
		from_port   = 80
		to_port     = 80
		protocol    = "tcp"
		cidr_blocks = ["0.0.0.0/0"]
	}

	ingress {
		from_port   = 443
		to_port     = 443
		protocol    = "tcp"
		cidr_blocks = ["0.0.0.0/0"]
	}

	egress {
		from_port   = 8080
		to_port     = 8080
		protocol    = "tcp"
		security_groups = [aws_security_group.es-writer-app-sg.id]
	}
}

# アプリケーションサーバーのセキュリティグループ
resource "aws_security_group" "es-writer-app-sg" {
	vpc_id      = "${aws_vpc.es-writer.id}"
	name        = "es-writer-app-sg"
	description = "es-writer-app-sg"
	tags = {
	   Name = "es-writer-app-sg"
	}

	ingress {
		from_port   = 80
		to_port     = 80
		protocol    = "tcp"
		cidr_blocks = ["0.0.0.0/0"]
	}

	egress {
		from_port   = 5432
		to_port     = 5432
		protocol    = "tcp"
		security_groups = [aws_security_group.es-writer-db-sg.id]
	}
}

# DBサーバーのセキュリティグループ
resource "aws_security_group" "es-writer-db-sg" {
	name        = "es-writer-db-sg"
	description = "es-writer-db-sg"
	vpc_id      = aws_vpc.es-writer.id
	tags = {
	   Name = "es-writer-db-sg"
	}

	ingress {
		from_port       = 5432
		to_port         = 5432
		protocol        = "tcp"
		security_groups = [aws_security_group.es-writer-app-sg.id]
	}
}

# 管理者用のセキュリティグループ
resource "aws_security_group" "es-writer-admin-sg" {
	name        = "es-writer-admin-sg"
	description = "es-writer-admin-sg"
	vpc_id      = aws_vpc.es-writer.id
	tags = {
	   Name = "es-writer-admin-sg"
	}

	ingress {
		from_port   = 22
		to_port     = 22
		protocol    = "tcp"
		cidr_blocks = ["0.0.0.0/0"]
	}
}
