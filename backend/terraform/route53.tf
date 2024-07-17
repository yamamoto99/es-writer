# ====================
# Route53
# ====================
data "aws_route53_zone" "es-writer" {
	name         = var.domain_name
	private_zone = false
}

resource "aws_route53_record" "es-writer-record" {
	zone_id = data.aws_route53_zone.es-writer.zone_id
	name    = var.subdomain_name
	type    = "A"
	
	alias {
		name                   = aws_lb.es-writer-alb.dns_name
		zone_id                = aws_lb.es-writer-alb.zone_id
		evaluate_target_health = true
	}
}
