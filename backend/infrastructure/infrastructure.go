package infrastructure

import (
	"es-app/model"
	"os"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type IIinfrastructure interface {
	CreateCognitoClient(c echo.Context) (*cognitoidentityprovider.Client, error)
	ValidateToken(c echo.Context, accessToken string) (model.LoginUser, error)
	SignUp(c echo.Context, signUpUser model.SignUpUser) (model.User, error)
	CheckEmail(c echo.Context, checkEmail model.CheckEmail) (bool, error)
	ResendEmail(c echo.Context, resendEmail model.ResendEmail) (bool, error)
	LogIn(c echo.Context, logInUser model.LoginUser) (model.LoginResponse, error)
	RefreshToken(c echo.Context, refreshToken string) (model.LoginResponse, model.LoginUser, error)
	GetUserID(c echo.Context, accessToken string) (string, error)
}

type infrastructure struct {
	svc       *cognitoidentityprovider.Client
	jwtKeyURL string
	clientID  string
}

func NewInfrastructure() IIinfrastructure {
	return &infrastructure{
		svc:       nil,
		jwtKeyURL: os.Getenv("TOKEN_KEY_URL"),
		clientID:  os.Getenv("COGNITO_CLIENT_ID"),
	}
}

func (i *infrastructure) CreateCognitoClient(c echo.Context) (*cognitoidentityprovider.Client, error) {
	cfg, err := config.LoadDefaultConfig(c.Request().Context(), config.WithRegion(os.Getenv("COGNITO_REGION")))
	if err != nil {
		return nil, err
	}
	i.svc = cognitoidentityprovider.NewFromConfig(cfg)
	return i.svc, nil
}

func (i *infrastructure) SignUp(c echo.Context, signUpUser model.SignUpUser) (model.User, error) {
	svc, err := i.CreateCognitoClient(c)
	if err != nil {
		return model.User{}, err
	}

	signUpInput := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(i.clientID),
		Username: aws.String(signUpUser.Username),
		Password: aws.String(signUpUser.Password),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(signUpUser.Email),
			},
		},
	}

	signUpOutput, err := svc.SignUp(c.Request().Context(), signUpInput)
	if err != nil {
		return model.User{}, err
	}

	userID := *signUpOutput.UserSub
	createdAt := time.Now()

	user := model.User{
		UserID:    userID,
		Username:  signUpUser.Username,
		Email:     signUpUser.Email,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	return user, nil
}

func (i *infrastructure) CheckEmail(c echo.Context, checkEmail model.CheckEmail) (bool, error) {
	svc, err := i.CreateCognitoClient(c)
	if err != nil {
		return false, err
	}

	confirmSignUpInput := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(i.clientID),
		Username:         aws.String(checkEmail.Username),
		ConfirmationCode: aws.String(checkEmail.VerificationCode),
	}

	_, err = svc.ConfirmSignUp(c.Request().Context(), confirmSignUpInput)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (i *infrastructure) ResendEmail(c echo.Context, resendEmail model.ResendEmail) (bool, error) {
	svc, err := i.CreateCognitoClient(c)
	if err != nil {
		return false, err
	}

	resendInput := &cognitoidentityprovider.ResendConfirmationCodeInput{
		ClientId: aws.String(i.clientID),
		Username: aws.String(resendEmail.Username),
	}

	_, err = svc.ResendConfirmationCode(c.Request().Context(), resendInput)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (i *infrastructure) LogIn(c echo.Context, logInUser model.LoginUser) (model.LoginResponse, error) {
	svc, err := i.CreateCognitoClient(c)
	if err != nil {
		return model.LoginResponse{}, err
	}

	authInput := &cognitoidentityprovider.InitiateAuthInput{
		ClientId: aws.String(i.clientID),
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME": logInUser.Username,
			"PASSWORD": logInUser.Password,
		},
	}

	authOutput, err := svc.InitiateAuth(c.Request().Context(), authInput)
	if err != nil {
		return model.LoginResponse{}, err
	}

	res := model.LoginResponse{
		AccessToken:  *authOutput.AuthenticationResult.AccessToken,
		IDToken:      *authOutput.AuthenticationResult.IdToken,
		RefreshToken: *authOutput.AuthenticationResult.RefreshToken,
	}

	return res, nil
}

func (i *infrastructure) ValidateToken(c echo.Context, accessToken string) (model.LoginUser, error) {
	svc, err := i.CreateCognitoClient(c)
	if err != nil {
		return model.LoginUser{}, err
	}

	jwtKey, err := keyfunc.NewDefault([]string{i.jwtKeyURL})
	if err != nil {
		return model.LoginUser{}, err
	}

	token, err := jwt.Parse(accessToken, jwtKey.Keyfunc)
	if err != nil {
		return model.LoginUser{}, err
	}

	if !token.Valid {
		return model.LoginUser{}, nil
	}

	input := &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(accessToken),
	}
	loginRes, err := svc.GetUser(c.Request().Context(), input)
	if err != nil {
		return model.LoginUser{}, err
	}

	res := model.LoginUser{
		Username: *loginRes.Username,
	}

	return res, nil
}

func (i *infrastructure) RefreshToken(c echo.Context, refreshToken string) (model.LoginResponse, model.LoginUser, error) {
	svc, err := i.CreateCognitoClient(c)
	if err != nil {
		return model.LoginResponse{}, model.LoginUser{}, err
	}

	refreshInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeRefreshToken,
		ClientId: aws.String(i.clientID),
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
		},
	}

	refreshOutput, err := svc.InitiateAuth(c.Request().Context(), refreshInput)
	if err != nil {
		return model.LoginResponse{}, model.LoginUser{}, err
	}

	newToken := model.LoginResponse{
		AccessToken:  *refreshOutput.AuthenticationResult.AccessToken,
		IDToken:      *refreshOutput.AuthenticationResult.IdToken,
		RefreshToken: refreshToken,
	}

	input := &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(newToken.AccessToken),
	}
	loginRes, err := svc.GetUser(c.Request().Context(), input)
	if err != nil {
		return model.LoginResponse{}, model.LoginUser{}, err
	}

	res := model.LoginUser{
		Username: *loginRes.Username,
	}

	return newToken, res, nil
}

func (i *infrastructure) GetUserID(c echo.Context, accessToken string) (string, error) {
	jwtKey, err := keyfunc.NewDefault([]string{i.jwtKeyURL})
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse(accessToken, jwtKey.Keyfunc)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", err
	}

	return sub, nil
}
