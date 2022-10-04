package cmd

import (
	"context"
	urand "crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	cowsay "github.com/Code-Hex/Neo-cowsay/v2"
	"github.com/Code-Hex/Neo-cowsay/v2/decoration"
	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/authhandler"
)

const (
	oauthClientID     = "d8ae7ea0-4a1a-46b0-b556-6d742687223a"
	oauthCallbackAddr = "localhost:8085"

	oauthBaseURL  = "https://login.tidbyt.com/oauth2/auth"
	oauthTokenURL = "https://login.tidbyt.com/oauth2/token"
)

var LoginCmd = &cobra.Command{
	Use:     "login",
	Short:   "Login to your Tidbyt account",
	Example: "login",
	Run:     login,
}

var (
	oauthConf = &oauth2.Config{
		ClientID: "d8ae7ea0-4a1a-46b0-b556-6d742687223a",
		Scopes:   []string{"device", "offline_access"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.tidbyt.com/oauth2/auth",
			TokenURL: "https://login.tidbyt.com/oauth2/token",
		},
		RedirectURL: fmt.Sprintf("http://%s", oauthCallbackAddr),
	}
)

type authResult struct {
	code  string
	state string
	err   error
}

func login(cmd *cobra.Command, args []string) {
	server := &http.Server{
		Addr: oauthCallbackAddr,
	}

	authResult := struct {
		code  string
		state string
	}{}

	done := make(chan bool, 1)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authResult.code = r.URL.Query().Get("code")
		authResult.state = r.URL.Query().Get("state")

		w.Header().Set("content-type", "text/plain")
		io.WriteString(w, "please close this window and return to pixlet")
		done <- true
	})

	handler := func(url string) (string, string, error) {
		if err := browser.OpenURL(url); err != nil {
			fmt.Println("Open this URL in your browser to login to your Tidbyt account:")
		} else {
			fmt.Println("We just opened your browser to visit:")
		}
		fmt.Printf("\n%s\n\n", url)

		go server.ListenAndServe()
		<-done

		return authResult.code, authResult.state, nil
	}

	cv, err := cv.CreateCodeVerifier()
	if err != nil {
		fmt.Println("creating PKCE code verifier:", err)
		os.Exit(1)
	}

	buf := make([]byte, 32)
	if _, err := urand.Read(buf); err != nil {
		fmt.Println("couldn't generate enough random bytes for state")
		os.Exit(1)
	}
	state := base64.RawURLEncoding.EncodeToString(buf)

	pkce := &authhandler.PKCEParams{
		Challenge:       cv.CodeChallengeS256(),
		ChallengeMethod: "S256",
		Verifier:        cv.String(),
	}

	tok, err := authhandler.TokenSourceWithPKCE(
		context.Background(),
		oauthConf,
		state,
		handler,
		pkce,
	).Token()
	if err != nil {
		fmt.Println("getting auth token:", err)
		os.Exit(1)
	}

	privateConfig.Set("token", tok)

	if err := privateConfig.WriteConfig(); err != nil {
		fmt.Println("persisting token:", err)
		os.Exit(1)
	}

	say, err := cowsay.Say(
		"A journey of a thousand API calls begins with a single login.",
		cowsay.Type("turtle"),
	)
	if err == nil {
		rand.Seed(time.Now().Unix())
		w := decoration.NewWriter(os.Stdout, decoration.WithAurora(rand.Intn(256)), decoration.WithBold())
		fmt.Fprintln(w, say)
	} else {
		fmt.Println("successfully logged in")
	}
}
