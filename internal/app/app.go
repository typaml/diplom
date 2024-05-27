package app

import (
	"diplom/internal/database/postgre"
	"diplom/internal/transport/rest"
	"fmt"
	"log"
	"net/http"
	"os"
)

func Run() {
	db, err := postgre.NewDB()
	fmt.Println(db)
	if err != nil {
		fmt.Println("Database error")
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rest.IndexHandler(w, r, "IndexHandler: ")
	})
	http.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		rest.ClientsHandler(w, r, "ClientsHandler: ")
	})
	http.HandleFunc("/getclients", func(w http.ResponseWriter, r *http.Request) {
		rest.GetClients(w, r, db, "GetClients: ")
	})
	http.HandleFunc("/analytics", func(w http.ResponseWriter, r *http.Request) {
		rest.AnalyticsHandler(w, r, "AnalyticsHandler: ")
	})
	http.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		rest.AccountHandler(w, r, db, "AccountHandler: ")
	})
	http.HandleFunc("/notifications", func(w http.ResponseWriter, r *http.Request) {
		rest.NotificationsHandler(w, r, "NotifHandler: ")
	})
	http.HandleFunc("/client/", func(w http.ResponseWriter, r *http.Request) {
		rest.ClientHandler(w, r, db, "ClientHandler: ")
	})
	http.HandleFunc("/statuses", func(w http.ResponseWriter, r *http.Request) {
		rest.Statuses(w, r, db, "Statuses: ")
	})
	http.HandleFunc("/regions", func(w http.ResponseWriter, r *http.Request) {
		rest.Region(w, r, db, "Region: ")
	})
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		rest.Users(w, r, db, "Users: ")
	})
	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		rest.Event(w, r, db, "Event: ")
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		rest.LoginHandler(w, r, db, "LoginHandler: ")
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		rest.RegisterHandler(w, r, db, "RegisterHandler: ")
	})
	http.HandleFunc("/save-client-changes", func(w http.ResponseWriter, r *http.Request) {
		rest.SaveClientChangesHandler(w, r, db, "SaveClientChangesHanlder")
	})
	http.HandleFunc("/get-client-history", func(w http.ResponseWriter, r *http.Request) {
		rest.GetClientHistoryHandler(w, r, db, "GetClientHistoryHandler: ")
	})
	http.HandleFunc("/add-task", func(w http.ResponseWriter, r *http.Request) {
		rest.AddTaskHandler(w, r, db, "AddTaskHandler: ")
	})
	http.HandleFunc("/get-tasks", func(w http.ResponseWriter, r *http.Request) {
		rest.GetTasksHandler(w, r, db, "GetTaskHandler: ")
	})
	http.HandleFunc("/get-tasks-notif", func(w http.ResponseWriter, r *http.Request) {
		rest.GetTasksHandlerToNotif(w, r, db, "GetTaskHandlerToNotif: ")
	})
	http.HandleFunc("/calls", func(w http.ResponseWriter, r *http.Request) {
		rest.GetCallsDataHandler(w, r, db, "GetCallsDataHandler: ")
	})
	http.HandleFunc("/sales", func(w http.ResponseWriter, r *http.Request) {
		rest.GetSalesDataHandler(w, r, db, "GetSalesDataHandler: ")
	})
	http.HandleFunc("/execute-task", func(w http.ResponseWriter, r *http.Request) {
		rest.ExecuteTaskHandler(w, r, db, "ExecuteTaskHandler: ")
	})
	http.HandleFunc("/getclientsstatus", func(w http.ResponseWriter, r *http.Request) {
		rest.GetClientStatusHandler(w, r, db, "GetClientStatusHandler: ")
	})
	http.HandleFunc("/getclientsevents", func(w http.ResponseWriter, r *http.Request) {
		rest.GetClientMarketingHandler(w, r, db, "GetClientMarketingHandler: ")
	})
	http.HandleFunc("/getclientsregion", func(w http.ResponseWriter, r *http.Request) {
		rest.ClientsByRegionHandler(w, r, db, "ClientsByRegionHandler: ")
	})

	http.HandleFunc("/getcallstoday", func(w http.ResponseWriter, r *http.Request) {
		rest.GetCallsTodayHandler(w, r, db, "GetCallsTodayHandler: ")
	})
	http.HandleFunc("/getclientsbyuser", func(w http.ResponseWriter, r *http.Request) {
		rest.GetClientsByUser(w, r, db, "GetClientsByUser: ")
	})
	http.HandleFunc("/createclient", func(w http.ResponseWriter, r *http.Request) {
		rest.CreateClientHandler(w, r, db, "CreateClientHandler: ")
	})
	http.HandleFunc("/delete-client", func(w http.ResponseWriter, r *http.Request) {
		rest.DeleteClientHandler(w, r, db, "DeleteClientHandler: ")
	})
	http.HandleFunc("/logout", rest.LogoutHandler)
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	port := os.Getenv("PORT")
	log.Println("Сервер заработал на порту: ", port)
	http.ListenAndServe("0.0.0.0:"+port, nil)
}
