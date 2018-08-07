package dcosclient

import (
	"encoding/json"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
)

// Need to see if this can be removed later
//DcosBasicAuth struct
type DcosBasicAuth struct {
	UID      string `json:"uid"`
	Password string `json:"password"`
}

// Need to see if this can be removed later
//DcosAuthResponse struct
type DcosAuthResponse struct {
	Token string `json:"token"`
}

//AsSecret represents the structure of the secret created by the service account script
type AsSecret struct {
	LoginEndpoint string `json:"login_endpoint,omitempty"`
	PrivateKey    string `json:"private_key,omitempty"`
	Scheme        string `json:"scheme,omitempty"`
	UID           string `json:"uid,omitempty"`
}

//AuthToken represents the format expected by the auth API
type AuthToken struct {
	UID   string `json:"uid,omitempty"`
	Token string `json:"token,omitempty"`
}

//TokenClaims blaster
type TokenClaims struct {
	UID string `json:"uid,omitempty"`
	jwt.StandardClaims
}

//Authenticate via a JWT token
func (c *Client) authSecret(asSecStr string) {
	// Get the CA
	//c.downloadFile("dcos-ca.crt", "/ca/dcos-ca.crt")

	asSec := AsSecret{}
	json.Unmarshal([]byte(asSecStr), &asSec)
	logrus.Infof("AS_SECRET read for uid %s", asSec.UID)

	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(asSec.PrivateKey))
	if err != nil {
		logrus.Panicln(err)
	}

	// Only validation serverside is for the 'uid' field
	claims := TokenClaims{
		asSec.UID,
		jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedString, err := token.SignedString(signingKey)
	if err != nil {
		logrus.Panicln(err)
	}
	authToken := AuthToken{
		UID:   asSec.UID,
		Token: signedString,
	}
	c.doAuth(authToken)
}

// Authenticate via User/Password
func (c *Client) authUserPassword(user, pass string) {
	usrPass := DcosBasicAuth{user, pass}
	c.doAuth(usrPass)
}

// Make an auth request, store it in the token
func (c *Client) doAuth(authData interface{}) {
	req, err := c.newRequest("POST", 443, "/acs/api/v1/auth/login", authData)
	if err != nil {
		logrus.Errorln(err)
		logrus.Panicf("Error trying to authenticate with %s", authData)
	}

	body, _ := c.do(req)
	var result DcosAuthResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		logrus.Errorln(body)
		logrus.Errorln(err)
		logrus.Panicln("Couldn't convert to dcosAuthResponse")
	}

	logrus.Infoln("Token obtained successfully")
	c.Token = result.Token
}

func (c *Client) auth() {
	asSecStr := os.Getenv("DCOS_SECRET")
	user := os.Getenv("DCOS_USERID")
	pass := os.Getenv("DCOS_PASSWORD")
	// Did we get a service account with a secret?
	if len(asSecStr) > 0 {
		c.authSecret(asSecStr)
	} else {
		// Did we get username/password?
		if len(user) == 0 || len(pass) == 0 {
			logrus.Panicln("Missing DCOS_SECRET or (DCOS_USERID and DCOS_PASSWORD) environment variables")
		} else {
			c.authUserPassword(user, pass)
		}
	}
}
