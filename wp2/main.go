package finalproject

import (
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/nu7hatch/gouuid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

var temp *template.Template

func init() {
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))
	http.HandleFunc("/", serve)
	http.HandleFunc("/login/", loginpage)
	http.HandleFunc("/register/", registerpage)
	http.HandleFunc("/state/", statepage)

	temp = template.Must(template.ParseGlob("public/*.html"))

}

func serve(res http.ResponseWriter, req *http.Request) {

	cookie, err := req.Cookie("session-id")

	if err != nil {
		id, _ := uuid.NewV4()
		cookie = &http.Cookie{
			Name:  "session-id",
			Value: id.String(),
			// Secure: true,
			HttpOnly: true,
		}
	}

	http.SetCookie(res, cookie)

	temp.ExecuteTemplate(res, "index.html", cookie)
}

func loginpage(res http.ResponseWriter, req *http.Request) {

	temp.ExecuteTemplate(res, "login.html", nil)
}

func statepage(res http.ResponseWriter, req *http.Request) {

	// cookie, err := req.Cookie("session-id")
	//
	// if err != nil {
	// 	id, _ := uuid.NewV4()
	// 	cookie = &http.Cookie{
	// 		Name:  "session-id",
	// 		Value: id.String(),
	// 		// Secure: true,
	// 		HttpOnly: true,
	// 	}
	// }

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
		http.Redirect(res, req, "/index?q="+newsessionid, http.StatusSeeOther)

	}

	temp.ExecuteTemplate(res, "register.html", nil)
}
