1. Install Go and set up oauth2 authroization with google or dax
   - for google head to console.cloud.google.com
   - click api and services on the left side menu
   - select credentials from menu
   - register an oauth2 application
   - set the redirect URI to http://127.0.0.1:5556/auth/google/callback (that's what the gettoken.go program uses)
   - copy the CLIENT ID and CLIENT SECRET generated

2. run ./build.sh to build the go code based archive library, C header file and C example program

3. Edit the set_env file to set proper values to the variables OAUTH2_CLIENT_ID and OAUTH2_CLIENT_SECRET

4. source ./set_env
   > go run gettoken.go
   > Click on the listening link to launch a browser consent screen for OAuth2
   > Copy the ID token and Nonce values from the terminal window

5. source ./set_env
   > run the test C program by calling ./example 
   > Enter the ID token and nonce when prompted from the data copied from the gettoken.go program output 
   > The userid or (subject id) shall be displayed
