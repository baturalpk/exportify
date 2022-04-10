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
3. Create a file named ".env" in the current directory and open to edit.
4. Copy the following text into the file and replace "Your_Client_ID" part with your previously obtained Spotify Client ID:<br>
    ```SPOTIFY_ID=Your_Client_ID```
5. Save and close the file.
6. Execute the binary from shell.
7. A window will be opened in your browser and gonna ask for authorization. (Required to be able to export your either private or public Spotify playlists)
8. If everything is okay until now, just lay back and keep an eye on the shell :)
9. Upon the successful process, a file named ```exportify-data.json``` will be created at the same directory 
   that contains your public & private playlists with details.
