package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	"github.com/pquerna/otp/totp"
)

var authenticatedUser = false

type User struct {
	Username string
	Password string
	Secret   string
}

// Simple in-memory "database" for demonstration purposes
var users = map[string]*User{
	"john": {Username: "john", Password: "password", Secret: ""},
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello!")
}

func main() {

	dir := http.Dir("./static")
	fs := http.FileServer(dir)
	mux := http.NewServeMux()
	mux.Handle("/", fs)
	fmt.Printf("Heads up: linux files need a !n at the end \n")

	fmt.Print("No. Just this one \n")
	mux.HandleFunc("GET /path/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, "Got path \n!")
	})

	mux.HandleFunc("GET /task/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "python -c import open3d as o3d; print(o3d.__version__) Got task \n %s \n", id)
	})

	//mux.HandleFunc("/hello", helloHandler)

	mux.HandleFunc("POST /recentfile", func(w http.ResponseWriter, r *http.Request) {

		if !authenticatedUser {
			//kickOutUnauthenticatedUser(w, r)
			//return
		}

		// Create or overwrite the file "userchangethis.csv"
		os.Remove("userchangethis.csv") // Remove if the file already exists
		userfile, err := os.OpenFile("userchangethis.csv", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatalf("Error creating/opening file: %s \n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		defer userfile.Close()

		s := bufio.NewReader(r.Body)

		// Skip potential response headers if they're structured as separate lines
		// (assuming headers might be the first 2-4 lines)
		for i := 0; i < 4; i++ {
			line, err := s.ReadString('\n')
			if err != nil {
				break // Exit loop if we hit EOF or another error
			}
			if line == "\n" || line == "\r\n" {
				break // Stop if we reach an empty line, assuming headers are done
			}
		}
		// Read the rest of the content and remove unnecessary patterns (e.g., "-+[0-9]+--")
		var b bytes.Buffer
		b.ReadFrom(s) // Reads remaining content from the body

		// Clean up specific unwanted patterns
		finalstring := b.String()
		m1 := regexp.MustCompile("-+[0-9]+--")
		finalstring = m1.ReplaceAllString(finalstring, "")

		// Write the cleaned content to the CSV file
		_, err = userfile.WriteString(finalstring)
		if err != nil {
			fmt.Println("Error writing to file: \n", err)
			http.Error(w, "Error saving file content \n", http.StatusInternalServerError)
			return
		}

		fmt.Println("File saved as userchangethis.csv without response headers \n")

	})

	mux.HandleFunc("GET /recentfile", func(w http.ResponseWriter, r *http.Request) {
		thisfile, err := os.ReadFile("userchangethis.csv")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		//fmt.Fprintf(w, string(thisfile), "%v\n")
		finalstring := string(thisfile)

		//TODO: TRIM WHITESPACEdocs
		fmt.Fprintf(w, finalstring)
	})

	mux.HandleFunc("GET /runPythonTest", func(w http.ResponseWriter, r *http.Request) {
		runPytonOpen3dScript("Hey")

		fmt.Fprintf(w, "I think we are able to get the python3 filename. \n")
	})

	mux.HandleFunc("PUT /recentfile", func(w http.ResponseWriter, r *http.Request) {
		//SelectRandomRowsToSize("userchangethis.csv", int64(2*1024*1024), "tempfile.csv")
		//	SelectRandomRowsToSize("userchangethis.csv", int64(1*1024*512), "userchangethis.csv")
		print("Processing finished")
	})

	mux.HandleFunc("POST /authenticator", func(w http.ResponseWriter, r *http.Request) {
		authenticatedUser = true
		fmt.Fprintf(w, "You're lucky I dont have the authenticator set up! \n")

	})

	mux.HandleFunc("GET /Home", func(w http.ResponseWriter, r *http.Request) {
		// Execute the index.html template
		t, err := template.ParseFiles("templates/index.html")

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		t.Execute(w, nil)
	})

	//this is the GET section
	//we have another function for PUT requests
	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		// Render the login.html template for GET requests
		t, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
	})

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Form parse error. Contact admin.", http.StatusInternalServerError)
			return
		}

		username := r.Form.Get("username")
		password := r.Form.Get("password")

		user, ok := users[username]

		if !ok || user.Password != password {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		//	if user.Secret == "" {
		http.Redirect(w, r, "/genOTP?username="+username, http.StatusFound)
		//		return
		//	}

		//http.Redirect(w, r, "/dashboard", http.StatusFound)
		t, err := template.ParseFiles("templates/validate.html")
		if err != nil {
			http.Error(w, fmt.Sprint("%v", err), http.StatusInternalServerError)
			return
		}
		//fmt.Printf("Okay, we got the username %d", username)
		t.Execute(w, struct{ Username string }{Username: username})
	})

	mux.HandleFunc("GET /genOTP", func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")

		user, ok := users[username]
		fmt.Printf("USER DETECTED: USERNAME: " + user.Username)

		//user, ok := users[username]
		if !ok {
			http.Redirect(w, r, "/dashboard", http.StatusFound)
			return
		}

		//if (user.Secret) == "" {
		//	secret, err := totp.Generate(totp.GenerateOpts{
		//		Issuer:      "HAMMER_ERC",
		//		AccountName: username,
		//	})
		//	if err != nil {
		//		http.Error(w, "Failed to generate TOTP secret.", http.StatusInternalServerError)
		//		return
		//	}
		//	user.Secret = secret.Secret()
		//	fmt.Println(user.Secret)
		//}

		secret, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "HAMMER_ERC",
			AccountName: username,
		})

		if err != nil {
			http.Error(w, "Failed to generate TOTP secret.", http.StatusInternalServerError)
			return
		}

		user.Secret = secret.Secret()
		//users[0].Secret = secret.Secret()
		fmt.Println("USER SECRET IS: " + user.Secret)

		otpURL := fmt.Sprintf("otpauth://totp/HAMMER_ERC:%s?secret=%s&issuer=HAMMER_ERC", username, user.Secret)

		// Prepare data to pass to the template
		data := struct {
			OTPURL   string
			Username string
		}{
			OTPURL:   otpURL,
			Username: username,
		}

		t, err := template.ParseFiles("templates/qrcodes.html")
		if err != nil {
			//http.Error(w, fmt.Sprint("%v", err), http.StatusInternalServerError)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		t.Execute(w, data)

	})

	mux.HandleFunc("/valOTP", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			username := r.URL.Query().Get("username")
			t, err := template.ParseFiles("templates/validate.html")
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			t.Execute(w, struct{ Username string }{Username: username})
			return

		case "POST":
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Form parse error. Contact admin.", http.StatusInternalServerError)
				return
			}
			currusername := r.FormValue("username")
			otpcode := r.FormValue("otpCode")

			user, ok := users[currusername]

			//user, ok := users[username]
			if !ok {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			isValid := totp.Validate(otpcode, user.Secret)
			if !isValid {
				// If validation fails, redirect back to the validation page
				fmt.Printf("VALIDATION FAILED! SENT SECRET: " + otpcode + "USER SECRET: " + user.Secret)
				http.Redirect(w, r, fmt.Sprintf("/valOTP?username=%s", currusername), http.StatusTemporaryRedirect)
				return
			} else {
				// If OTP is valid, set a session cookie (simplified for this example) and redirect to dashboard
				http.SetCookie(w, &http.Cookie{
					Name:   "authenticatedUser",
					Value:  "true",
					Path:   "/",
					MaxAge: 3600, // 1 hour for example
				})

				http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

			}

		}
	})

	mux.HandleFunc("/whatsmysecret", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintf(w, "secret is: %d", user.Secret)

	})

	mux.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the authenticated user's username from the session cookie
		username, err := r.Cookie("authenticatedUser")
		if err != nil || username.Value == "" {
			// If user is not authenticated, redirect to the homepage
			http.Redirect(w, r, "template/ErrorPage", http.StatusFound)
			return
		}

		// Render the dashboard.html template
		t, err := template.ParseFiles("templates/dashboard.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
	})

	port := os.Getenv("HTTP_PLATFORM_PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":8080", mux)
}

//func dewnewMethod(w http.ResponseWriter, r *http.Request) {
//	var b bytes.Buffer
//	b.ReadFrom(r.Body)
//	fmt.Fprint(w, b.String())
//
//}

func kickOutUnauthenticatedUser(w http.ResponseWriter, r *http.Request) {
	if authenticatedUser == false {
		fmt.Fprintf(w, "Hey. You're unauthenticated! Get out!")
		//return
	}
}

func runPytonOpen3dScript(filename string) bool {
	cmd := exec.Command("./pythonReaction/stlToObj.py", filename)

	cmd.Stdout = os.Stdout
	return true

}
