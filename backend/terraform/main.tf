terraform {
	// 必要なプロバイダーとそのバージョンを定義
	required_providers {
		aws = {
			source  = "hashicorp/aws"
			version = "~> 3.0"
		}
	}
	// 必要なTerraformのバージョンを定義
	required_version = ">= 1.0.0"
}

provider "aws" {
	// AWSの認証情報を設定
	access_key = "${var.aws_access_key}"
	secret_key = "${var.aws_secret_key}"
	region     = "${var.aws_region}"
}

resource "aws_vpc" "main" {
	// VPCを作成
	cidr_block = "10.0.0.0/16"
	tags = {
		Name = "main-vpc"
	}
}

resource "aws_internet_gateway" "igw" {
	// インターネットゲートウェイを作成
	vpc_id = aws_vpc.main.id
	tags = {
		Name = "main-igw"
	}
}

resource "aws_subnet" "public_a" {
	// パブリックサブネットを作成 (アベイラビリティゾーンa)
	vpc_id            = aws_vpc.main.id
	cidr_block        = "10.0.1.0/24"
	availability_zone = "ap-northeast-1a"
	tags = {
		Name = "public-subnet-a"
	}
}

resource "aws_subnet" "public_c" {
	// パブリックサブネットを作成 (アベイラビリティゾーンc)
	vpc_id            = aws_vpc.main.id
	cidr_block        = "10.0.2.0/24"
	availability_zone = "ap-northeast-1c"
	tags = {
		Name = "public-subnet-c"
	}
}

resource "aws_subnet" "private_a" {
	// プライベートサブネットを作成 (アベイラビリティゾーンa)
	vpc_id            = aws_vpc.main.id
	cidr_block        = "10.0.5.0/24"
	availability_zone = "ap-northeast-1a"
	tags = {
		Name = "private-subnet-a"
	}
}

resource "aws_subnet" "private_c" {
	// プライベートサブネットを作成 (アベイラビリティゾーンc)
	vpc_id            = aws_vpc.main.id
	cidr_block        = "10.0.6.0/24"
	availability_zone = "ap-northeast-1c"
	tags = {
		Name = "private-subnet-c"
	}
}

resource "aws_route_table" "public" {
	// パブリックサブネット用のルートテーブルを作成
	vpc_id = aws_vpc.main.id
	route {
		cidr_block = "0.0.0.0/0"
		gateway_id = aws_internet_gateway.igw.id
	}
	tags = {
		Name = "public-route-table"
	}
}

resource "aws_route_table_association" "public_a" {
	// ルートテーブルをパブリックサブネット (アベイラビリティゾーンa) に関連付け
	subnet_id      = aws_subnet.public_a.id
	route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "public_c" {
	// ルートテーブルをパブリックサブネット (アベイラビリティゾーンc) に関連付け
	subnet_id      = aws_subnet.public_c.id
	route_table_id = aws_route_table.public.id
}

resource "aws_security_group" "web" {
	// Webサーバー用のセキュリティグループを定義
	vpc_id = aws_vpc.main.id
	name = "allow_http"
	description = "Allow HTTP traffic"

	ingress {
		from_port   = 80
		to_port     = 80
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

resource "aws_security_group" "rds" {
	// RDS用のセキュリティグループを定義
	vpc_id = aws_vpc.main.id
	name = "allow_postgres"
	description = "Allow PostgreSQL traffic"

	ingress {
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

resource "aws_instance" "web" {
	// EC2インスタンスを定義
	ami               = "ami-0f9fe1d9214628296"
	instance_type     = "t2.micro"
	subnet_id         = aws_subnet.public_a.id
	associate_public_ip_address = true
	vpc_security_group_ids = [aws_security_group.web.id]
	tags = {
		Name = "progate-aws-app"
	}
}

resource "aws_db_instance" "main" {
	// RDSインスタンスを定義
	identifier        = "progate-db"
	allocated_storage = 20
	storage_type      = "gp2"
	engine            = "postgres"
	engine_version    = "16.3"
	instance_class    = "db.t3.micro"
	password = "${var.rds_pass}"
	username = "${var.rds_username}"
	db_subnet_group_name = aws_db_subnet_group.main.name
	vpc_security_group_ids = [aws_security_group.rds.id]
	skip_final_snapshot = true
	tags = {
		Name = "progate-db"
	}
}

resource "aws_db_subnet_group" "main" {
	// RDS用のサブネットグループを定義
	name       = "main-subnet-group"
	subnet_ids = [aws_subnet.private_a.id, aws_subnet.private_c.id]
	tags = {
		Name = "main-subnet-group"
	}
}
