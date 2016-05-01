package finalproject

import (
	"encoding/json"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

var temp *template.Template

func init() {
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))
	http.HandleFunc("/", serve)
	http.HandleFunc("/login/", loginpage)
	http.HandleFunc("/register/", registerpage)
	http.HandleFunc("/state/", statepage)
	http.HandleFunc("/home/", homepage)

	temp = template.Must(template.ParseGlob("public/*.html"))

}

func homepage(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	var mysess Session
	var myuser User
	_, session, err := getsess(req)

	//if not logged in
	if err != nil {
		http.Redirect(res, req, `/`, http.StatusSeeOther)
		return
	}

	mysess.id = session

	info, err := memcache.Get(ctx, session)

	json.Unmarshal(info.Value, &myuser)
	mysess.User = myuser

	temp.ExecuteTemplate(res, "homepage.html", mysess)
}

func serve(res http.ResponseWriter, req *http.Request) {
	// ctx := appengine.NewContext(req)
	var mysess Session
	_, session, err := getsess(req)
	mysess.id = session

	if err == nil {
		http.Redirect(res, req, `/index?q=`+mysess.id, http.StatusSeeOther)
	}

	temp.ExecuteTemplate(res, "index.html", nil)
}

func loginpage(res http.ResponseWriter, req *http.Request) {
	var mysess Session
	var myuser User
	ctx := appengine.NewContext(req)

	if req.Method == "POST" {
		//get data from form
		useremail := req.FormValue("username")
		userpass := req.FormValue("password")

		//check the datastore for that info
		key := datastore.NewKey(ctx, "Users", useremail, 0, nil)
		err := datastore.Get(ctx, key, &myuser)

		log.Infof(ctx, "UNAME:", useremail)
		log.Infof(ctx, "PASS:", userpass)

		hiddenpass, err := bcrypt.GenerateFromPassword([]byte(userpass), bcrypt.DefaultCost)

		if err != nil || string(hiddenpass) != myuser.Password {
			log.Infof(ctx, "*** Error Info: Login Failed, given credentials not found in datastore. ***")
			mysess.alerts = "Logged in Failed! \n Email or password incorrect"
		} else {
			//login passed
			mysess.id = makesess(res, req, myuser)
			http.Redirect(res, req, `/dashboard?id=`+mysess.id, http.StatusSeeOther)
		}

		// if err == nil && bcrypt.CompareHashAndPassword([]byte(myuser.Password), []byte(userpass)) == nil {
		// 	mysess.id = makesess(res, req, myuser)
		// 	http.Redirect(res, req, `/index?q=`+mysess.id, http.StatusSeeOther)
		// } else {
		// 	log.Infof(ctx, "User information was not found in datastore, Not Logged in!")
		// 	mysess.alerts = "Login failed!!"
		// }

	}

	temp.ExecuteTemplate(res, "login.html", mysess)
}

func statepage(res http.ResponseWriter, req *http.Request) {

	key := "q"
	Val := req.URL.Query().Get(key)

	temp.ExecuteTemplate(res, "state.html", Val)
}

func registerpage(res http.ResponseWriter, req *http.Request) {

	var mysess Session
	var myuser User
	ctx := appengine.NewContext(req)

	if req.Method == "POST" {
		myemail := req.FormValue("username")
		mypass1 := req.FormValue("password1")
		mypass2 := req.FormValue("password2")
		myuser.Email = myemail

		userkey := datastore.NewKey(ctx, "Users", myuser.Email, 0, nil)
		err := datastore.Get(ctx, userkey, &myuser)

		//Username not unique
		if err == nil {
			log.Infof(ctx, "!!problem in register, User already created!!")

			mysess.alerts = "Email already exists \n "
			temp.ExecuteTemplate(res, "register.html", mysess)
			return
		}

		//passwords dont match
		if mypass1 != mypass2 {
			log.Infof(ctx, "Password does not match confirmed pass")
			mysess.alerts += "passwords do not match!!"
			temp.ExecuteTemplate(res, "register.html", mysess)
		}

		//hash pass
		hiddenpass, err := bcrypt.GenerateFromPassword([]byte(mypass1), bcrypt.DefaultCost)
		if err != nil {
			log.Errorf(ctx, "Could not bcrypt the pass", err)
			http.Error(res, err.Error(), 500)
		}

		//create the registered user
		reggeduser := User{
			Email:    myemail,
			Password: string(hiddenpass),
		}

		//make key for hash table
		userkey = datastore.NewKey(ctx, "User", reggeduser.Email, 0, nil)
		userkey, err = datastore.Put(ctx, userkey, &reggeduser) //save to datastore

		if err != nil {
			log.Errorf(ctx, "While Registering, could not store user to datastore", err)
			http.Error(res, err.Error(), 500)
			return
		}

		newsessionid := makesess(res, req, myuser)
		http.Redirect(res, req, "/home/?q="+newsessionid, http.StatusSeeOther)

	}

	temp.ExecuteTemplate(res, "register.html", mysess)
}
