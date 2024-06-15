terraform {
	// 必要なプロバイダとそのバージョンを定義
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

resource "aws_security_group" "web" {
	// Webサーバー用のセキュリティグループを定義
	name = "allow_http"
	description = "Allow HTTP inbound traffic"

	// HTTPのインバウンドトラフィックを許可
	ingress {
		from_port = 80
		to_port = 80
		protocol = "tcp"
		cidr_blocks = ["0.0.0.0/0"]
	}

	// すべてのアウトバウンドトラフィックを許可
	egress {
		from_port = 0
		to_port = 0
		protocol = "-1"
		cidr_blocks = ["0.0.0.0/0"]
	}
}

resource "aws_instance" "web" {
	// EC2インスタンスを定義
	ami = "ami-0f9fe1d9214628296"
	instance_type = "t2.micro"

	// セキュリティグループを設定
	vpc_security_group_ids = [aws_security_group.web.id]

	// タグを設定
	tags = {
		Name = "progate-aws-app"
	}
}

resource "aws_security_group" "rds" {
	// RDS用のセキュリティグループを定義
	name = "allow_postgres"
	description = "Allow PostgreSQL inbound traffic"

	// PostgreSQLのインバウンドトラフィックを許可
	ingress {
		from_port = 5432
		to_port = 5432
		protocol = "tcp"
		cidr_blocks = ["0.0.0.0/0"]
	}

	// すべてのアウトバウンドトラフィックを許可
	egress {
		from_port = 0
		to_port = 0
		protocol = "-1"
		cidr_blocks = ["0.0.0.0/0"]
	}
}

resource "aws_db_instance" "main" {
	// RDSインスタンスを定義
	identifier = "progate-db"
	allocated_storage = 20
	storage_type = "gp2"
	engine = "postgres"
	engine_version = "16.3"
	instance_class = "db.t3.micro"
	password = "${var.rds_pass}"
	username = "${var.rds_username}"

	// セキュリティグループを設定
	vpc_security_group_ids = [aws_security_group.rds.id]
	skip_final_snapshot = true

	// タグを設定
	tags = {
		Name = "progate-db"
	}
}
