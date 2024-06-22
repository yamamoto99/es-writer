# ====================
# VPC
# ====================
resource "aws_vpc" "progate" {
	cidr_block           = "10.0.0.0/16"
	instance_tenancy     = "default"
	enable_dns_support   = "true"
	enable_dns_hostnames = "true"
	tags = {
		Name = "progate-VPC"
	}
}

# ====================
# Subnet
# ====================
# パブリックサブネット
resource "aws_subnet" "public_web" {
	vpc_id                  = "${aws_vpc.progate.id}"
	cidr_block              = "10.0.1.0/24"
	map_public_ip_on_launch = true
	availability_zone       = "us-west-2a"
	tags = {
		Name = "public-web"
	}
}

# プライベートサブネット
resource "aws_subnet" "private_db_1" {
	vpc_id                  = "${aws_vpc.progate.id}"
	cidr_block              = "10.0.2.0/24"
	map_public_ip_on_launch = false
	availability_zone       = "us-west-2a"
	tags = {
		Name = "private-db-1"
	}
}

# プライベートサブネット
resource "aws_subnet" "private_db_2" {
	vpc_id                  = "${aws_vpc.progate.id}"
	cidr_block              = "10.0.3.0/24"
	map_public_ip_on_launch = false
	availability_zone       = "us-west-2c"
	tags = {
		Name = "private-db-2"
	}
}

# ====================
# Internet Gateway
# ====================
resource "aws_internet_gateway" "progate" {
	vpc_id = "${aws_vpc.progate.id}"
	tags = {
		Name = "progate"
	}
}

# ====================
# Route Table
# ====================
# パブリックルートテーブル
resource "aws_route_table" "progate" {
	vpc_id = "${aws_vpc.progate.id}"
	tags   = {
		Name = "public-route"
	}
}

# インターネットゲートウェイへのルート
resource "aws_route" "public_route" {
	gateway_id             = "${aws_internet_gateway.progate.id}"
	route_table_id         = "${aws_route_table.progate.id}"
	destination_cidr_block = "0.0.0.0/0"
}

# ルートテーブルとパブリックサブネットの関連付け
resource "aws_route_table_association" "public_a" {
	subnet_id      = "${aws_subnet.public_web.id}"
	route_table_id = "${aws_route_table.progate.id}"
}

# ====================
# Security Group
# ====================
# WEBサーバーのセキュリティグループ
resource "aws_security_group" "progate-web-sg" {
	vpc_id      = "${aws_vpc.progate.id}"
	name        = "progate-web-sg"
	description = "progate-web-sg"
	tags = {
	   Name = "progate-web-sg"
	}

	ingress {
		from_port   = 80
		to_port     = 80
		protocol    = "tcp"
		cidr_blocks = ["0.0.0.0/0"]
	}

	ingress {
		from_port   = 22
		to_port     = 22
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
		from_port   = 5432
		to_port     = 5432
		protocol    = "tcp"
		cidr_blocks = ["0.0.0.0/0"]
	}

	egress {
		from_port   = 0
		to_port     = 0
		protocol    = "-1"
		cidr_blocks = ["0.0.0.0/0"]
	}
}

# DBサーバーのセキュリティグループ
resource "aws_security_group" "progate-db-sg" {
	name        = "progate-db-sg"
	description = "progate-db-sg"
	vpc_id      = aws_vpc.progate.id
	tags = {
	   Name = "db-sg"
	}

	ingress {
		from_port       = 5432
		to_port         = 5432
		protocol        = "tcp"
		security_groups = ["${aws_security_group.progate-web-sg.id}"]
	}

	egress {
		from_port   = 0
		to_port     = 0
		protocol    = "-1"
		cidr_blocks = ["0.0.0.0/0"]
	}
}