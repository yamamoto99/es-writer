# ====================
# Create IAM policy
# ====================
resource "aws_iam_policy" "custom_cognito_policy" {
	name_prefix = "CustomCognitoPolicy"
	description = "Custom policy for Cognito User Pools and Clients"
	policy      = jsonencode({
		Version = "2012-10-17",
		Statement = [
			{
				Effect = "Allow",
				Action = [
				"cognito-idp:CreateUserPool",
				"cognito-idp:UpdateUserPool",
				"cognito-idp:DeleteUserPool",
				"cognito-idp:DescribeUserPool",
				"cognito-idp:CreateUserPoolClient",
				"cognito-idp:UpdateUserPoolClient",
				"cognito-idp:DeleteUserPoolClient",
				"cognito-idp:GetUserPoolMfaConfig"
				],
				Resource = [
					"arn:aws:cognito-idp:ap-northeast-1:637423575781:userpool/*"
				],
			},
		],
	})
}

# ====================
# Attach IAM policy to IAM user
# ====================
resource "aws_iam_user_policy_attachment" "attach_custom_cognito_policy" {
	user       = "${var.iam_username}"
	policy_arn = aws_iam_policy.custom_cognito_policy.arn
}
resource "aws_iam_user_policy_attachment" "attach_amazon_ec2_full_access" {
	user       = "${var.iam_username}"
	policy_arn = "arn:aws:iam::aws:policy/AmazonEC2FullAccess"
}

resource "aws_iam_user_policy_attachment" "attach_amazon_vpc_full_access" {
	user       = "${var.iam_username}"
	policy_arn = "arn:aws:iam::aws:policy/AmazonVPCFullAccess"
}

resource "aws_iam_user_policy_attachment" "attach_amazon_rds_full_access" {
	user       = "${var.iam_username}"
	policy_arn = "arn:aws:iam::aws:policy/AmazonRDSFullAccess"
}