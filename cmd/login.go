package cmd

import (
	"context"
	urand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"tidbyt.dev/pixlet/cmd/config"

	cowsay "github.com/Code-Hex/Neo-cowsay/v2"
	"github.com/Code-Hex/Neo-cowsay/v2/decoration"
	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2/authhandler"
)

var loginCommandJSON bool

func init() {
	LoginCmd.Flags().BoolVar(&loginCommandJSON, "json", false, "output login information as json")
}

var LoginCmd = &cobra.Command{
	Use:     "login",
	Short:   "Login to your Tidbyt account",
	Example: "login",
	Run:     login,
}

func login(cmd *cobra.Command, args []string) {
	server := &http.Server{
		Addr: config.OAuthCallbackAddr,
	}

	var authCode, authState string
	done := make(chan bool)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authCode = r.URL.Query().Get("code")
		authState = r.URL.Query().Get("state")

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

		return authCode, authState, nil
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
		config.OAuthConf,
		state,
		handler,
		pkce,
	).Token()
	if err != nil {
		fmt.Println("getting auth token:", err)
		os.Exit(1)
	}

	config.PrivateConfig.Set("token", tok)

	if err := config.PrivateConfig.WriteConfig(); err != nil {
		fmt.Println("persisting token:", err)
		os.Exit(1)
	}

	if loginCommandJSON {
		b, err := json.Marshal(tok)
		if err != nil {
			fmt.Printf("could not marshal token: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", b)
		os.Exit(0)
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
