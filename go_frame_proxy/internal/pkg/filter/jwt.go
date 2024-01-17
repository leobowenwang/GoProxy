package filter

import (
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/project-sesame/sesame-gateway/internal/pkg/database"
	"github.com/project-sesame/sesame-gateway/internal/pkg/util"
	"go.uber.org/zap"
)

type jwtToken struct {
	Name  string `json:"username"`
	Roles string `json:"roles"`
	jwt.StandardClaims
}

type signedToken struct {
	Token string `json:"token"`
}

var secret string = "your_secret_key"

func getCertificate(certFilePath string) (*x509.Certificate, error) {
	certBytes, err := ioutil.ReadFile(certFilePath)
	if err != nil {
		return nil, fmt.Errorf("error while reading the content of certificate.: %w", err)
	}

	certBlock, _ := pem.Decode(certBytes)
	if certBlock == nil {
		return nil, fmt.Errorf("error while decoding Certificate")
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error while parsing Certificate: %w", err)
	}

	return cert, nil
}

func GetPrivateKey(keyFilePath string) (privateKey *rsa.PrivateKey, err error) {

	keyBytes, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("error while reading the content of Keyfile.: %w", err)
	}

	keyBlock, _ := pem.Decode(keyBytes)
	if keyBlock == nil {
		return nil, fmt.Errorf("error while decoding Keyfile")
	}

	key, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return
	}

	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("failed to assert private key type")
	}

	return
}

func CheckCertFile(certFile, keyFile string) (cert *x509.Certificate, privateKey *rsa.PrivateKey, err error) {

	// Load certificate and private key from files

	cert, err = getCertificate(certFile)

	if err != nil {
		return nil, nil, err
	}

	privateKey, err = GetPrivateKey(keyFile)
	if err != nil {
		return nil, nil, err
	}

	return
}

func BasicToJWT(certFile, keyFile string) func(next echo.HandlerFunc) echo.HandlerFunc {

	cert, privateKey, _ := CheckCertFile(certFile, keyFile)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			username, password, ok := c.Request().BasicAuth()
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing credentials")
			}

			token, err := parseToJWT(username, password, secret, cert, privateKey)
			if err != nil {
				// Return appropriate error status code and message
				switch err {
				case sql.ErrNoRows, echo.ErrUnauthorized:
					return echo.NewHTTPError(http.StatusUnauthorized, "Incorrect username or password")
				case echo.ErrForbidden:
					return echo.NewHTTPError(http.StatusForbidden, "Invalid Password")
				default:
					return err
				}
			}

			// Set token in response header
			util.Logger.Debug("Token: " + token)

			//c.Response().Header().Set("Authorization", "Bearer "+token)
			c.Request().Header.Set("Authorization", "Bearer "+token)

			return next(c)
		}
	}
}

func parseToJWT(username, password string, secret string, cert *x509.Certificate, privateKey *rsa.PrivateKey) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			util.Logger.Error("Panic in parseJWT", zap.Any("error", r))
		}
	}()

	database := &database.Database{}
	authenticated, roles, err := BCryptAuthenticate(database, username, password)
	if err != nil {
		// Log error
		util.Logger.Error("Error authenticating user:", zap.Error(err))

		return "", echo.ErrForbidden
	}

	if !authenticated {
		return "", echo.ErrUnauthorized
	}

	tokenString, err := GenerateJWT(username, roles, privateKey)
	if err != nil {
		util.Logger.Error("Error generating JWT token:", zap.Error(err))

		return "", fmt.Errorf("error generating JWT token")
	}

	return tokenString, nil
}

func GenerateJWT(username string, roles string, privateKey *rsa.PrivateKey) (tokenString string, err error) {
	// Set custom claims
	claims := &jwtToken{
		Name:  username,
		Roles: roles,
		StandardClaims: jwt.StandardClaims{
			// Expires after 5 minutes
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Generate encoded token
	tokenString, err = token.SignedString(privateKey)

	return
}
