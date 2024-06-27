# # ====================
# # EC2 Instance
# # ====================
# resource "aws_instance" "progate-web" {
# 	ami                         = "ami-0f226ae5ce4b11922"
# 	instance_type               = "t2.micro"
# 	subnet_id                   = "${aws_subnet.public_web.id}"
# 	associate_public_ip_address = true
# 	vpc_security_group_ids      = ["${aws_security_group.progate-web-sg.id}"]
# 	key_name                    = "${var.key_name}"
# 	tags = {
# 		Name = "progate-app"
# 	}
# 	user_data = <<-EOF
# 		#!/bin/bash
# 		sudo yum update -y
# 		sudo amazon-linux-extras install docker -y
# 		sudo amazon-linux-extras install -y postgresql14
# 		sudo systemctl start docker
# 		sudo systemctl enable docker
# 		sudo usermod -aG docker ec2-user
# 	EOF
# }

# # ====================
# # Elastic IP
# # ====================
# resource "aws_eip" "progate" {
# 	instance = aws_instance.progate-web.id
# }