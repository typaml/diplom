package app

import (
	"diplom/internal/database/postgre"
	"diplom/internal/transport/rest"
	"fmt"
	"net/http"
	"os"
)

func Run() {
	db, err := postgre.NewDB()
	fmt.Println(db)
	if err != nil {
		fmt.Println("Database error")
	}
	http.HandleFunc("/", rest.IndexHandler)
	http.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		rest.ClientsHandler(w, r, db)
	})
	http.HandleFunc("/getclients", func(w http.ResponseWriter, r *http.Request) {
		rest.GetClients(w, r, db)
	})
	http.HandleFunc("/analytics", rest.AnalyticsHandler)
	http.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		rest.AccountHandler(w, r, db)
	})
	http.HandleFunc("/notifications", rest.NotificationsHandler)
	http.HandleFunc("/client/", func(w http.ResponseWriter, r *http.Request) {
		rest.ClientHandler(w, r, db)
	})
	http.HandleFunc("/statuses", func(w http.ResponseWriter, r *http.Request) {
		rest.Statuses(w, r, db)
	})
	http.HandleFunc("/regions", func(w http.ResponseWriter, r *http.Request) {
		rest.Region(w, r, db)
	})
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		rest.Users(w, r, db)
	})
	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		rest.Event(w, r, db)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		rest.LoginHandler(w, r, db)
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		rest.RegisterHandler(w, r, db)
	})
	http.HandleFunc("/save-client-changes", func(w http.ResponseWriter, r *http.Request) {
		rest.SaveClientChangesHandler(w, r, db)
	})
	http.HandleFunc("/get-client-history", func(w http.ResponseWriter, r *http.Request) {
		rest.GetClientHistoryHandler(w, r, db)
	})
	http.HandleFunc("/add-task", func(w http.ResponseWriter, r *http.Request) {
		rest.AddTaskHandler(w, r, db)
	})
	http.HandleFunc("/get-tasks", func(w http.ResponseWriter, r *http.Request) {
		rest.GetTasksHandler(w, r, db)
	})
	http.HandleFunc("/get-tasks-notif", func(w http.ResponseWriter, r *http.Request) {
		rest.GetTasksHandlerToNotif(w, r, db)
	})
	http.HandleFunc("/calls", func(w http.ResponseWriter, r *http.Request) {
		rest.GetCallsDataHandler(w, r, db)
	})
	http.HandleFunc("/sales", func(w http.ResponseWriter, r *http.Request) {
		rest.GetSalesDataHandler(w, r, db)
	})
	http.HandleFunc("/execute-task", func(w http.ResponseWriter, r *http.Request) {
		rest.ExecuteTaskHandler(w, r, db)
	})
	http.HandleFunc("/getclientsstatus", func(w http.ResponseWriter, r *http.Request) {
		rest.GetClientStatusHandler(w, r, db)
	})
	http.HandleFunc("/getclientsevents", func(w http.ResponseWriter, r *http.Request) {
		rest.GetClientMarketingHandler(w, r, db)
	})
	http.HandleFunc("/getclientsregion", func(w http.ResponseWriter, r *http.Request) {
		rest.ClientsByRegionHandler(w, r, db)
	})

	http.HandleFunc("/getcallstoday", func(w http.ResponseWriter, r *http.Request) {
		rest.GetCallsTodayHandler(w, r, db)
	})
	http.HandleFunc("/getclientsbyuser", func(w http.ResponseWriter, r *http.Request) {
		rest.GetClientsByUser(w, r, db)
	})
	http.HandleFunc("/createclient", func(w http.ResponseWriter, r *http.Request) {
		rest.CreateClientHandler(w, r, db)
	})
	http.HandleFunc("/delete-client", func(w http.ResponseWriter, r *http.Request) {
		rest.DeleteClientHandler(w, r, db)
	})
	http.HandleFunc("/logout", rest.LogoutHandler)
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	port := os.Getenv("PORT")
	http.ListenAndServe("0.0.0.0:"+port, nil)
}
