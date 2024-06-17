terraform {
	required_providers {
		aws = {
		source  = "hashicorp/aws"
		version = "~> 3.0"
		}
	}
	required_version = ">= 1.0.0"
}
# providerの宣言
provider "aws" {
	access_key = "${var.aws_access_key}"
	secret_key = "${var.aws_secret_key}"
	region     = "${var.aws_region}"
}

# VPCの設定
resource "aws_vpc" "dev_env" {
	cidr_block           = "10.0.0.0/16"
	instance_tenancy     = "default"
	enable_dns_support   = "true"
	enable_dns_hostnames = "true"
	tags = {
		Name = "dev-env"
	}
}

# パブリックサブネットの作成
resource "aws_subnet" "public_web" {
	vpc_id                  = "${aws_vpc.dev_env.id}"
	cidr_block              = "10.0.0.0/24"
	map_public_ip_on_launch = true
	availability_zone       = "ap-northeast-1a"
	tags = {
		Name = "public-web"
	}
}

# プライベートサブネットの作成
resource "aws_subnet" "private_db_1" {
	vpc_id                  = "${aws_vpc.dev_env.id}"
	cidr_block              = "10.0.1.0/24"
	map_public_ip_on_launch = false
	availability_zone       = "ap-northeast-1a"
	tags = {
		Name = "private-db-1"
	}
}

# プライベートサブネットの作成
resource "aws_subnet" "private_db_2" {
	vpc_id                  = "${aws_vpc.dev_env.id}"
	cidr_block              = "10.0.2.0/24"
	map_public_ip_on_launch = false
	availability_zone       = "ap-northeast-1c"
	tags = {
		Name = "private-db-2"
	}
}

# DB Subnet Groupの設定
resource "aws_db_subnet_group" "db_subnet" {
	name        = "db-subnet"
	description = "It is a DB subnet group on dev-env."
	subnet_ids  =["${aws_subnet.private_db_1.id}","${aws_subnet.private_db_2.id}"]
	tags = {
		Name = "db-subnet"
	}
}

# InternetGatewayの設定
resource "aws_internet_gateway" "dev_env_gw" {
	vpc_id     = "${aws_vpc.dev_env.id}"
	depends_on = [aws_vpc.dev_env]
	tags = {
		Name = "dev-env-gw"
	}
}

# パブリックルートテーブルの作成
resource "aws_route_table" "public_route" {
	vpc_id = "${aws_vpc.dev_env.id}"
	tags   = {
		Name = "public-route"
	}
}

# インターネットゲートウェイへのルートを設定
resource "aws_route" "public_route" {
	destination_cidr_block = "0.0.0.0/0"
	gateway_id             = "${aws_internet_gateway.dev_env_gw.id}"
	route_table_id         = "${aws_route_table.public_route.id}"
}

# ルートテーブルとパブリックサブネットの関連付け
resource "aws_route_table_association" "public_a" {
	subnet_id      = "${aws_subnet.public_web.id}"
	route_table_id = "${aws_route_table.public_route.id}"
}

#  WEBサーバーのセキュリティグループの作成
resource "aws_security_group" "web_security_group" {
	vpc_id      = "${aws_vpc.dev_env.id}"
	name        = "web_security_group"
	description = "it is a security group on http of dev-env"

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
		from_port   = 0
		to_port     = 0
		protocol    = "-1"
		cidr_blocks = ["0.0.0.0/0"]
	}
}

# WEBサーバーの設定
resource "aws_instance" "web" {
	ami                         = "ami-0f9fe1d9214628296"
	instance_type               = "t2.micro"
	subnet_id                   = "${aws_subnet.public_web.id}"
	associate_public_ip_address = true
	vpc_security_group_ids      = ["${aws_security_group.web_security_group.id}"]
	key_name                    = "${var.key_name}"
	root_block_device {
		volume_type = "gp2"
		volume_size = "20"
	}
	ebs_block_device {
		device_name = "/dev/sdf"
		volume_type = "gp2"
		volume_size = "10"
	}
	tags = {
		Name = "progate-aws-app"
	}
	user_data = <<-EOF
		#!/bin/bash
		sudo dnf update -y
		sudo dnf install -y postgresql15
		sudo dnf install docker -y
		sudo systemctl start docker
		sudo systemctl enable docker
		sudo usermod -a -G docker ec2-user
	EOF
}

#  DBサーバーのセキュリティグループの作成
resource "aws_security_group" "db-sg" {
	name        = "db-sg"
    description = "It is a security group on db of dev-env"
    vpc_id      = aws_vpc.dev_env.id
    tags = {
       Name = "db-sg"
    }


	ingress {
		from_port                = 3306
		to_port                  = 3306
		protocol                 = "tcp"
		security_groups = ["${aws_security_group.web_security_group.id}"]
	}

	egress {
		from_port   = 0
		to_port     = 0
		protocol    = "-1"
		cidr_blocks = ["0.0.0.0/0"]
	}
}

# RDSサーバーの設定
resource "aws_db_instance" "main" {
	identifier             = "progate-db"
	allocated_storage      = 20
	storage_type           = "gp2"
	engine                 = "postgres"
	engine_version         = "15.7"
	instance_class         = "db.t3.micro"
	password               = "${var.rds_pass}"
	username               = "${var.rds_username}"
	db_subnet_group_name   = "${aws_db_subnet_group.db_subnet.name}"
	vpc_security_group_ids = ["${aws_security_group.db-sg.id}"]
	skip_final_snapshot    = true
	multi_az               = false
	availability_zone      = "ap-northeast-1a"
	publicly_accessible    = true
	tags = {
		Name = "progate-db"
	}
}
