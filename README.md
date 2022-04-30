### Getting Started
1. Log in to your Spotify for Developers dashboard.
2. Create a new app.
3. Go into your newly created app and click 'EDIT SETTINGS'
4. Add following URL into the "Redirect URIs", save and close the menu: `http://localhost:8080/callback`
5. Obtain your app's "Client ID" from the same board.

### Usage
1. Download the proper binary executable depending on your operating system from https://github.com/baturalpk/exportify/releases. <br>
   Optionally, you may want to compile "exportify.go" from the source.
2. Open a shell at the directory where you downloaded/created the binary.
3. Execute the binary from shell.
4. Enter your previously obtained Spotify Client ID.
5. A window will be opened in your browser and gonna ask for authorization. (Required to be able to export your either private or public Spotify playlists)
6. If everything is okay until now, just lay back and keep an eye on the shell :)
7. Upon the successful completion of process, a file called as ```exportify-data.json``` will be created at the same directory: 
   It contains your public & private playlists with details.
