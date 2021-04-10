package main

import (
	"C"
	"bufio"
	"fmt"
	"os"
	"strings"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/net/context"
)

var (
	verifier    *oidc.IDTokenVerifier
	ctx         context.Context
	clientID    = os.Getenv("OAUTH2_CLIENT_ID")
	providerURL = os.Getenv("OAUTH2_PROVIDER_URL")
)

func setOAuthConfig(clientID string, providerURL string) (bool, string) {
	ctx = context.Background()

	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		return false, err.Error()
	}
	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier = provider.Verifier(oidcConfig)
	return true, ""
}

func verifyToken(token string, nonce string) (bool, string, string) {

	idToken, err := verifier.Verify(ctx, token)
	if err != nil {
		return false, "", err.Error()
	}

	if idToken.Nonce != nonce {
		fmt.Printf("token nonce=[%s], given nonce=[%s]\n", idToken.Nonce, nonce)
		return false, "", "Nonce does not match"
	}

	var claims struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return false, "", "Error extracting claims from id token"
	}

	fmt.Printf("Sub = %s\n", claims.Sub)
	fmt.Printf("Email = %s\n", claims.Email)

	return true, claims.Sub, ""
}

//export setOAuthConfigC
func setOAuthConfigC(client_id *C.char, provider_url *C.char, emsg **C.char) C.int {
	retVal, errMsg := setOAuthConfig(C.GoString(client_id), C.GoString(provider_url))
	if retVal == false {
		*emsg = C.CString(errMsg)
		return 1
	}
	return 0
}

//export verifyTokenC
func verifyTokenC(token *C.char, nonce *C.char, uid **C.char, emsg **C.char) C.int {
	retVal, userId, errMsg := verifyToken(C.GoString(token), C.GoString(nonce))
	if retVal == false {
		*emsg = C.CString(errMsg)
		return 1
	}
	*uid = C.CString(userId)
	return 0
}

func main() {
	setOAuthConfig(clientID, providerURL)

	fmt.Print("Enter ID token: ")
	reader := bufio.NewReader(os.Stdin)
	idtoken, _ := reader.ReadString('\n')
	idtoken = strings.TrimSuffix(idtoken, "\n")

	fmt.Print("Enter nonce: ")
	nonce, _ := reader.ReadString('\n')
	nonce = strings.TrimSuffix(nonce, "\n")

	retVal, userid, emsg := verifyToken(idtoken, nonce)
	if retVal == false {
		fmt.Printf("Failed in verification: error-code=%s\n", emsg)
		return
	}

	fmt.Printf("Token verification success, userid=%s\n", userid)
}
