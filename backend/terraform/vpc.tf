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
resource "aws_subnet" "public_app_1" {
	vpc_id                  = "${aws_vpc.es-writer.id}"
	cidr_block              = "10.0.1.0/24"
	map_public_ip_on_launch = true
	availability_zone       = var.availability_zone_1
	tags = {
		Name = "public-app-1"
	}
}

# パブリックサブネット
resource "aws_subnet" "public_app_2" {
	vpc_id                  = "${aws_vpc.es-writer.id}"
	cidr_block              = "10.0.5.0/24"
	map_public_ip_on_launch = true
	availability_zone       = var.availability_zone_2
	tags = {
		Name = "public-app-2"
	}
}

# プライベートサブネット
resource "aws_subnet" "private_db_1" {
	vpc_id                  = "${aws_vpc.es-writer.id}"
	cidr_block              = "10.0.3.0/24"
	map_public_ip_on_launch = false
	availability_zone       = var.availability_zone_1
	tags = {
		Name = "private-db-1"
	}
}

# プライベートサブネット
resource "aws_subnet" "private_db_2" {
	vpc_id                  = "${aws_vpc.es-writer.id}"
	cidr_block              = "10.0.4.0/24"
	map_public_ip_on_launch = false
	availability_zone       = var.availability_zone_2
	tags = {
		Name = "private-db-2"
	}
}

# ====================
# Security Group
# ====================
# WEBサーバーのセキュリティグループ
resource "aws_security_group" "es-writer-web-sg" {
	vpc_id      = aws_vpc.es-writer.id
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
}

# アプリケーションサーバーのセキュリティグループ
resource "aws_security_group" "es-writer-app-sg" {
	vpc_id      = aws_vpc.es-writer.id
	name        = "es-writer-app-sg"
	description = "es-writer-app-sg"
	tags = {
		Name = "es-writer-app-sg"
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
		cidr_blocks = ["0.0.0.0/0"]  # セキュリティ上、特定のIPアドレス範囲に制限することを推奨します
	}
}

# WEBサーバーからアプリケーションサーバーへのアクセスを許可するルール
resource "aws_security_group_rule" "web_to_app" {
	type                     = "egress"
	from_port                = 8080
	to_port                  = 8080
	protocol                 = "tcp"
	security_group_id        = aws_security_group.es-writer-web-sg.id
	source_security_group_id = aws_security_group.es-writer-app-sg.id
}

# アプリケーションサーバーがWEBサーバーからのアクセスを受け入れるルール
resource "aws_security_group_rule" "app_from_web" {
	type                     = "ingress"
	from_port                = 80
	to_port                  = 80
	protocol                 = "tcp"
	security_group_id        = aws_security_group.es-writer-app-sg.id
	source_security_group_id = aws_security_group.es-writer-web-sg.id
}

# アプリケーションサーバーからDBサーバーへのアクセスをDB側で許可するルール
resource "aws_security_group_rule" "es-writer-ingress-db" {
	type                     = "ingress"
	from_port                = 5432
	to_port                  = 5432
	protocol                 = "tcp"
	security_group_id        = aws_security_group.es-writer-db-sg.id
	source_security_group_id = aws_security_group.es-writer-app-sg.id
}

# アプリケーションサーバーからDBサーバーへのアクセスをEC2側で許可するルール
resource "aws_security_group_rule" "es-writer-egress-ec2" {
	type                     = "egress"
	from_port                = 5432
	to_port                  = 5432
	protocol                 = "tcp"
	security_group_id        = aws_security_group.es-writer-app-sg.id
	source_security_group_id = aws_security_group.es-writer-db-sg.id
}

# ====================
# Internet Gateway
# ====================
resource "aws_internet_gateway" "es-writer-igw" {
	vpc_id = aws_vpc.es-writer.id
	tags = {
		Name = "es-writer-igw"
	}
}

# ====================
# Route Table
# ====================
# パブリックルートテーブル
resource "aws_route_table" "public_app_rt" {
	vpc_id = aws_vpc.es-writer.id
	tags = {
		Name = "public-app-rt"
	}
}

# インターネットゲートウェイへのルート
resource "aws_route" "public_route" {
	gateway_id             = aws_internet_gateway.es-writer-igw.id
	route_table_id         = aws_route_table.public_app_rt.id
	destination_cidr_block = "0.0.0.0/0"
}

# ルートテーブルとサブネットを関連付け
resource "aws_route_table_association" "public_app_rt_assoc_1" {
	subnet_id      = aws_subnet.public_app_1.id
	route_table_id = aws_route_table.public_app_rt.id
}
