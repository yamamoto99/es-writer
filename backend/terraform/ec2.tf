# ====================
# EC2 Instance
# ====================
resource "aws_instance" "es-writer-app" {
	ami                         = "ami-0f226ae5ce4b11922"
	instance_type               = "t2.micro"
	subnet_id                   = "${aws_subnet.public_app.id}"
	associate_public_ip_address = true
	vpc_security_group_ids      = ["${aws_security_group.es-writer-app-sg.id}"]
	key_name                    = "${var.key_name}"
	tags = {
		Name = "es-writer-app"
	}
	user_data = <<-EOF
		#!/bin/bash
		sudo yum update -y
		sudo amazon-linux-extras install docker -y
		sudo amazon-linux-extras install -y postgresql14
		sudo systemctl start docker
		sudo systemctl enable docker
		sudo usermod -aG docker ec2-user
	EOF
}
