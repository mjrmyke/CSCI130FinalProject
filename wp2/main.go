package finalproject

import (
	"html/template"
	"log"
	"net/http"

	"github.com/nu7hatch/gouuid"
)

func init() {
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))
	http.HandleFunc("/", serve)
	http.HandleFunc("/login/", loginpage)
	http.HandleFunc("/register/", registerpage)
	http.HandleFunc("/state/", statepage)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serve(res http.ResponseWriter, req *http.Request) {
	//create template from index.html file
	temp, err := template.ParseFiles("public/index.html")
	if err != nil {
		log.Fatalln(err, "Failed Parsing File")
	}

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

	temp.Execute(res, cookie)
}

func loginpage(res http.ResponseWriter, req *http.Request) {
	//create template from login.html file
	temp, err := template.ParseFiles("public/login.html")
	if err != nil {
		log.Fatalln(err, "Failed Parsing File")
	}
	temp.Execute(res, nil)
}

func statepage(res http.ResponseWriter, req *http.Request) {
	//create template from login.html file
	temp, err := template.ParseFiles("public/state.html")
	if err != nil {
		log.Fatalln(err, "Failed Parsing File")
	}

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

	temp.Execute(res, Val)
}

func registerpage(res http.ResponseWriter, req *http.Request) {
	//create template from login.html file
	temp, err := template.ParseFiles("public/register.html")
	if err != nil {
		log.Fatalln(err, "Failed Parsing File")
	}
	// var mysess Session
	// var myuser User
	// ctx := appengine.NewContext(req)

	// if req.Method == "POST" {
	// 	myemail := req.FormValue("username")
	// 	mypass1 := req.FormValue("password1")
	// 	mypass2 := req.FormValue("password2")
	// 	myuser.Email = myemail
	//
	// 	key := datastore.NewKey(ctx, "Users", user.Email, 0, nil)
	// 	err := datastore.Get(ctx, key, &user)
	//
	// 	if err == nil {
	// 		log.Infof(ctx, "!!problem in register, User already created!!")
	//
	// 		session.Message = "Email already exists \n "
	// 		temp.ExecuteTemplate(response, "register.html", session)
	// 		return
	// 	}
	// }

	temp.Execute(res, nil)
}
