# ====================
# Cognito User Pool
# ====================
resource "aws_cognito_user_pool" "pool" {
	# ユーザープールの名前を設定
	name                       = "es-writers-pool"
	
	# ユーザー名の設定を大文字小文字区別に設定
	username_configuration {
		case_sensitive = true
	}

	# 自動検証属性をメールに設定
	auto_verified_attributes   = ["email"]
	# メール検証メッセージを設定
	email_verification_message = "検証コードは {####} です。"
	# メール検証の件名を設定
	email_verification_subject = "検証コード"
	# 多要素認証の設定をオフに設定
	mfa_configuration          = "OFF"
	
	# 管理者によるユーザー作成設定
	admin_create_user_config {
		allow_admin_create_user_only = false
	}

	# メール設定をデフォルトのCognitoアカウントに設定
	email_configuration {
		email_sending_account = "COGNITO_DEFAULT"
	}
	
	# ユーザースキーマを設定
	schema {
		attribute_data_type = "String"
		name                = "email"
	}

	# パスワードポリシーを設定
	password_policy {
		minimum_length    = 8
		require_lowercase = true
		require_numbers   = true
		require_symbols   = true
		require_uppercase = true
	}

	# アカウント回復設定を設定
	account_recovery_setting {
		recovery_mechanism {
			name     = "verified_email"
			priority = 1
		}
	}
}

# ====================
# Cognito User Pool Client
# ====================
resource "aws_cognito_user_pool_client" "client" {
	# ユーザープールクライアントの名前を設定
	name                                 = "es-writers-client"
	# ユーザープールIDを設定
	user_pool_id                         = aws_cognito_user_pool.pool.id

	# OAuthフローを無効に設定
	allowed_oauth_flows_user_pool_client = false

	# 明示的な認証フローを設定
	explicit_auth_flows = [
		"ALLOW_USER_PASSWORD_AUTH",
		"ALLOW_REFRESH_TOKEN_AUTH"
	]

	# クライアントシークレットの生成を無効に設定
	generate_secret              = false
	# サポートされるIDプロバイダーをCognitoに設定
	supported_identity_providers = ["COGNITO"]
}
