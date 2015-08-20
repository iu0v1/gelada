package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"local/gelada"
)

func main() {

	// create session keys
	sessionKeys := [][]byte{
		[]byte("261AD9502C583BD7D8AA03083598653B"),
		[]byte("E9F6FDFAC2772D33FC5C7B3D6E4DDAFF"),
	}

	// create exception for "no auth" zone
	exceptions := []string{"/noauth/.*"}

	// set options
	options := gelada.Options{
		Path:     "/",
		MaxAge:   60, // 60 seconds
		HTTPOnly: true,

		SessionName:     "test-session",
		SessionLifeTime: 60, // 60 seconds
		SessionKeys:     sessionKeys,

		BindUserAgent: true,
		BindUserHost:  true,

		LoginUserFieldName:     "login",
		LoginPasswordFieldName: "password",
		LoginRoute:             "/login",
		LogoutRoute:            "/logout",

		AuthProvider: checkAuth,

		Exceptions: exceptions,
	}

	// get Gelada
	g, err := gelada.New(options)
	if err != nil {
		panic(err)
	}

	// create handler manager
	hm := &HandlerManager{Gelada: g}

	// create mux router
	router := mux.NewRouter()

	// main page
	router.HandleFunc("/", hm.HandleMainPage)
	// page to view which does not need authorization
	router.HandleFunc("/noauth/page", hm.HandleLoginFreePage)
	// login page
	router.HandleFunc("/login", hm.HandleLoginPage).Methods("GET")
	// function for processing a request for authorization (via POST method)
	router.HandleFunc("/login", g.AuthHandler).Methods("POST")
	// function for processing a request for logout (via POST method)
	router.HandleFunc("/logout", g.LogoutHandler).Methods("POST")

	// wrap around our router
	http.Handle("/", g.GlobalAuth(router))

	// start net/http server at 8082 port
	fmt.Println("start server at 127.0.0.1:8082")
	if err := http.ListenAndServe("127.0.0.1:8082", nil); err != nil {
		panic(err)
	}
}

// HandlerManager need for manage handlers and share some staff beetween them.
type HandlerManager struct {
	Gelada *gelada.Gelada
}

// HandleMainPage - main page.
func (hm *HandlerManager) HandleMainPage(res http.ResponseWriter, req *http.Request) {
	// get session client
	user, err := hm.Gelada.GetClient(req)
	if err != nil {
		fmt.Fprintf(res, "server side error: %v\n", err)
	}

	// create struct for our main page with some additional data
	pageData := struct {
		User         *gelada.Client // client
		ToSessionEnd int            // seconds to end of session
		LogoutRoute  string         // route for logout button
	}{
		User:         user,
		ToSessionEnd: user.TimeToEndOfSession(),
		LogoutRoute:  "/logout",
	}

	mainPage := template.Must(template.New("").Parse(`
		<html><head></head><body>
		<script>
		var sessionTimer = document.getElementById("sessionTimer");
		function startTimer(duration, display) {
		    var timer = duration, minutes, seconds;
		    var tick = setInterval(function() {
		        minutes = parseInt(timer / 60, 10);
		        seconds = parseInt(timer % 60, 10);

		        minutes = minutes < 10 ? "0" + minutes : minutes;
		        seconds = seconds < 10 ? "0" + seconds : seconds;

		        display.textContent = minutes + ":" + seconds;

		        if (--timer < 0) {
		            //timer = duration;
					clearInterval(tick);
		        }
		    }, 1000);
		}
		window.onload = function () {
		    var display = document.querySelector('#time');
		    startTimer("{{.ToSessionEnd}}", display);
		};
		</script>
		<center>
		<h1 style="padding-top:15%;">HELLO {{.User.Username}}!</h1><br>
		<div><span id="time">00:00</span> minutes to end of this session</div><br>
		<form action="{{.LogoutRoute}}" method="post">
			<button type="submit">Logout</button>
		</form>
		</center></body>
		</html>`),
	)
	mainPage.Execute(res, pageData)
}

// HandleLoginPage - login page.
func (hm *HandlerManager) HandleLoginPage(res http.ResponseWriter, req *http.Request) {
	var loginPage = template.Must(template.New("").Parse(`
		<html><head></head><body>
		<center>
		<form id="login_form" action="/login" method="POST" style="padding-top:15%;">
		<p>Login: user | Password: qwerty</p>
		<p>Or go to <a href="/noauth/page">login-free zone</a>.</p>
		<input type="text" name="login" placeholder="Login" autofocus><br>
		<input type="password" placeholder="Password" name="password"><br>
		<input type="submit" value="LOGIN">
		</form></center></body>
		</html>`),
	)
	loginPage.Execute(res, nil)
}

// HandleLoginFreePage - auth-free page.
func (hm *HandlerManager) HandleLoginFreePage(res http.ResponseWriter, req *http.Request) {
	var freePage = template.Must(template.New("").Parse(`
		<html><head></head><body>
		<center>
		<h2 style="padding-top:15%;">Free zone :)</h2><br>
		Auth has no power here!
		</html>`),
	)
	freePage.Execute(res, nil)
}

// auth provider function
func checkAuth(u, p string) bool {
	if u == "user" && p == "qwerty" {
		return true
	}
	return false
}
