package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"

	"github.com/google/uuid"
	pkceutils "github.com/jimlambrt/go-oauth-pkce-code-verifier"
	"github.com/pkg/browser"
	"github.com/spf13/viper"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const (
	Filename = "exportify-data"
	Port     = 8080
)

var (
	redirectURL = fmt.Sprintf("http://localhost:%d/callback", Port)

	ch   = make(chan *spotify.Client, 1)
	auth *spotifyauth.Authenticator

	state         string
	verifier      *pkceutils.CodeVerifier
	codeVerifier  string
	codeChallenge string
)

func init() {
	state = uuid.NewString()

	var err error
	verifier, err = pkceutils.CreateCodeVerifier()
	if err != nil {
		FatalX(err)
	}

	codeVerifier = verifier.String()
	codeChallenge = verifier.CodeChallengeS256()
}

func main() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		FatalX(err)
	}
	os.Setenv("SPOTIFY_ID", viper.GetString("SPOTIFY_ID"))

	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURL),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopePlaylistReadPrivate),
	)

	http.HandleFunc("/callback", handleOAuth)
	go http.ListenAndServe(fmt.Sprintf(":%d", Port), nil)

	url := auth.AuthURL(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	)

	fmt.Println("🚀 Redirecting to the log in page, check your browser:", url)
	if err := browser.OpenURL(url); err != nil {
		FatalX(err)
	}

	client := <-ch
	ctx := context.Background()

	user, err := client.CurrentUser(ctx)
	if err != nil {
		FatalX(err)
	}
	fmt.Println("🎉 You are logged in as:", user.ID, user.DisplayName)

	playlists, err := client.CurrentUsersPlaylists(ctx)
	if err != nil {
		FatalX(err)
	}

	fmt.Println("🤗 Exporting your playlists, please be patient...")
	var exportifyData []ExportifyPlaylist
	for _, playlist := range playlists.Playlists {
		pt, _ := client.GetPlaylistTracks(ctx, playlist.ID)

		var tracklist []ExportifyTrack
		for _, t := range pt.Tracks {
			var track ExportifyTrack
			track.ID = t.Track.ID
			track.Name = t.Track.Name
			track.Duration = t.Track.Duration
			track.Artists = t.Track.Artists
			track.Endpoint = t.Track.Endpoint
			tracklist = append(tracklist, track)
		}

		data := ExportifyPlaylist{
			PlaylistAttributes: playlist,
			Tracks:             tracklist,
		}

		exportifyData = append(exportifyData, data)
		time.Sleep(1 * time.Second)
	}

	marshaled, err := json.MarshalIndent(exportifyData, "", "  ")
	if err != nil {
		FatalX(err)
	}
	if err := ioutil.WriteFile(fmt.Sprintf("%s.json", Filename), marshaled, 644); err != nil {
		FatalX(err)
	}
	fmt.Println("🥳 It's done, check the contents of file", fmt.Sprintf("'%s.json' !", Filename))
}

func handleOAuth(w http.ResponseWriter, r *http.Request) {
	token, err := auth.Token(r.Context(), state, r, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		FatalX(err)
	}

	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		FatalX(errors.New(fmt.Sprintf("State mismatch: %s != %s\n", st, state)))
	}

	ch <- spotify.New(auth.Client(r.Context(), token))
	fmt.Fprintf(w, "Login attempt successful ✔")
}

type ExportifyPlaylist struct {
	PlaylistAttributes spotify.SimplePlaylist `json:"playlist_attributes"`
	Tracks             []ExportifyTrack       `json:"playlist_tracks"`
}

type ExportifyTrack struct {
	ID       spotify.ID             `json:"id"`
	Artists  []spotify.SimpleArtist `json:"artists"`
	Duration int                    `json:"duration_ms"`
	Endpoint string                 `json:"link"`
	Name     string                 `json:"name"`
}

func FatalX(err error) {
	fmt.Fprintln(os.Stderr, "☹ Error:", err.Error())
	os.Exit(1)
}
