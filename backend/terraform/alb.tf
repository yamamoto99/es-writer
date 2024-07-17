# ====================
# ALB for ES Writer
# ====================
resource "aws_lb" "es-writer-alb" {
	name               = "es-writer-alb"
	internal           = false
	load_balancer_type = "application"
	security_groups    = [aws_security_group.es-writer-web-sg.id]
	subnets            = [aws_subnet.public_app_1.id, aws_subnet.public_app_2.id]
	ip_address_type    = "ipv4"

	tags = {
		Name = "tf_test_alb"
	}
}

# ====================
# ALB target group
# ====================
resource "aws_lb_target_group" "es-writer-target-group" {
	name             = "es-writer-target-group"
	target_type      = "instance"
	protocol_version = "HTTP1"
	port             = 8080
	protocol         = "HTTP"

	vpc_id = aws_vpc.es-writer.id

	tags = {
		Name = "es-writer-target-group"
	}

	health_check {
		interval            = 30
		path                = "/"
		port                = "traffic-port"
		protocol            = "HTTP"
		timeout             = 5
		healthy_threshold   = 5
		unhealthy_threshold = 2
		matcher             = "200,301"
	}
}

# ====================
# ALB target group instance
# ====================
resource "aws_lb_target_group_attachment" "es-writer-alb-tga" {
	target_group_arn = aws_lb_target_group.es-writer-target-group.arn
	target_id        = aws_instance.es-writer-app.id
}

# ====================
# ALB listener
# ====================
resource "aws_lb_listener" "es-writer-listener" {
	load_balancer_arn = aws_lb.es-writer-alb.arn
	port              = "80"
	protocol          = "HTTP"

	default_action {
		type             = "forward"
		target_group_arn = aws_lb_target_group.es-writer-target-group.arn
	}
}

resource "aws_lb_listener" "es-writer-listener-https" {
	load_balancer_arn = aws_lb.es-writer-alb.arn
	port              = "443"
	protocol          = "HTTPS"

	certificate_arn   = var.https_cert_arn

	default_action {
		type             = "forward"
		target_group_arn = aws_lb_target_group.es-writer-target-group.arn
	}
}
