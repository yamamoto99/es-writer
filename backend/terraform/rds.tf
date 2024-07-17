# ====================
# RDS Parameter Group
# ====================
resource "aws_db_parameter_group" "es-writer-db-pg" {
	name        = "es-writer-db-pg"
	family      = "postgres14"
	description = "esWriter DB parameter group"

	parameter {
		name  = "rds.force_ssl"
		value = "1"
	}
}

# ====================
# RDS Instance
# ====================
# RDSサーバーの設定
resource "aws_db_instance" "main" {
	identifier             = "es-writer-db"
	db_name                = "es_writer_db"
	allocated_storage      = 20
	storage_type           = "gp2"
	engine                 = "postgres"
	engine_version         = "14.12"
	instance_class         = "db.t3.micro"
	password               = "${var.rds_pass}"
	username               = "${var.rds_username}"
	db_subnet_group_name   = "${aws_db_subnet_group.es-writer-db-subnet.name}"
	vpc_security_group_ids = ["${aws_security_group.es-writer-db-sg.id}"]
	parameter_group_name   = "${aws_db_parameter_group.es-writer-db-pg.name}"
	skip_final_snapshot    = true
	multi_az               = false
	availability_zone      = var.availability_zone_1
	publicly_accessible    = false
	tags = {
		Name = "es-writer-db"
	}
}

# ====================
# RDS Subnet Group
# ====================
resource "aws_db_subnet_group" "es-writer-db-subnet" {
	name        = "es-writer-db-subnet"
	description = "es-writer-db-subnet"
	subnet_ids  =["${aws_subnet.private_db_1.id}","${aws_subnet.private_db_2.id}"]
	tags = {
		Name = "es-writer-db-subnet"
	}
}