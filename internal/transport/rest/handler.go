package rest

import (
	"diplom/internal/database/postgre"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

func GetClients(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	sortColumn := r.URL.Query().Get("sortColumn")
	sortOrder := r.URL.Query().Get("sortOrder")
	filterId := r.URL.Query().Get("filterId")
	filterStatus := r.URL.Query().Get("filterStatus")
	filterUsers := r.URL.Query().Get("filterUsers")
	filterEvent := r.URL.Query().Get("filterEvent")
	filterRegion := r.URL.Query().Get("filterRegion")
	if sortColumn == "" {
		sortColumn = "id"
	}
	if sortOrder == "" {
		sortOrder = "asc"
	}
	clients, err := postgre.SortAndFilterClients(db.Db, sortColumn, sortOrder, filterId, filterStatus, filterRegion, filterUsers, filterEvent)
	if err != nil {
		log.Println(loggerString, "Failed to get clients from database", err)
		return
	}
	// Преобразовываем данные в формат JSON и отправляем клиенту
	jsonData, err := json.Marshal(clients)
	if err != nil {
		log.Println(loggerString, "Cannot marshal to json", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

// IndexHandler обрабатывает запросы к главной странице
func IndexHandler(w http.ResponseWriter, r *http.Request, loggerString string) {
	tmpl, err := template.ParseFiles("web/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(loggerString, "Internal Server Error:", err)
		return
	}
	session, _ := store.Get(r, "session-name")
	if session.Values["login"] != nil {
		log.Println(loggerString, session.Values["login"], "visited the server")
	} else {
		log.Println(loggerString, r.RemoteAddr, "visited the server")
	}
	tmpl.Execute(w, nil)
}

func ClientsHandler(w http.ResponseWriter, r *http.Request, loggerString string) {
	session, _ := store.Get(r, "session-name")
	if session.Values["login"] != nil {
		log.Println(loggerString, session.Values["login"], "visited the clients")
	} else {
		log.Println(loggerString, r.RemoteAddr, "visited the clients")
	}
	tmpl, err := template.ParseFiles("web/clients.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(loggerString, "Internal Server Error:", err)
		return
	}
	tmpl.Execute(w, nil)
}

func ClientHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	idStr := r.URL.Path[len("/client/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		log.Println(loggerString, "Invalid client ID:", err)
		return
	}
	client, err := postgre.Select(db.Db, id)
	if err != nil {
		http.Error(w, "Failed to get client from database", http.StatusInternalServerError)
		log.Println(loggerString, "Failed to get client from database:", err)
		return
	}
	htmlContent, err := os.ReadFile("./web/client.html")
	if err != nil {
		http.Error(w, "Failed to read HTML file", http.StatusInternalServerError)
		log.Println(loggerString, "Failed to read HTML file:", err)
		return
	}
	html := string(htmlContent)
	html = strings.Replace(html, "{{ClientID}}", strconv.Itoa(client.ID), -1)
	html = strings.Replace(html, "{{Name}}", client.Name, -1)
	html = strings.Replace(html, "{{Email}}", client.Email, -1)
	html = strings.Replace(html, "{{Phone}}", client.Phone, -1)
	html = strings.Replace(html, "{{Status}}", client.Status, -1)
	html = strings.Replace(html, "{{Region}}", client.Region, -1)
	html = strings.Replace(html, "{{UserName}}", client.UsersName, -1)
	html = strings.Replace(html, "{{Event}}", client.Event, -1)
	html = strings.Replace(html, "{{FCall}}", client.DateFirstCall, -1)
	html = strings.Replace(html, "{{NCall}}", client.DateNextCall, -1)
	html = strings.Replace(html, "{{UID}}", strconv.Itoa(client.UserID), -1)
	html = strings.Replace(html, "{{Total}}", strconv.Itoa(client.Total), -1)
	html = strings.Replace(html, "{{Payday}}", client.Payday, -1)
	html = strings.Replace(html, "{{Info}}", client.Information, -1)
	session, _ := store.Get(r, "session-name")
	if session.Values["login"] != nil {
		log.Println(loggerString, session.Values["login"], "visited the client: ", client.ID)
	} else {
		log.Println(loggerString, r.RemoteAddr, "visited the client: ", client.ID)
	}
	// Отправляем HTML-страницу клиенту
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}
func AccountHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Println(loggerString, "Failed to get session:", err)
	}
	var login string
	if session.Values["login"] != nil {
		var ok bool
		login, ok = session.Values["login"].(string)
		if !ok {
			log.Println(loggerString, "Не смог преобразовать логин в стрингу:", err)
		}
	} else {
		login = ""
	}
	if session.Values["login"] != nil {
		log.Println(loggerString, session.Values["login"], "logged into the server")
	} else {
		log.Println(loggerString, r.RemoteAddr, "visited the account settings")
	}
	user, err := postgre.LoginValidate(db.Db, login)
	if err != nil {
		log.Println(loggerString, "Ошибка проверки пользователя:", err)
		user = &postgre.User{
			ID:       0,
			Login:    "",
			Password: "",
			Name:     "",
			Surname:  "",
			Email:    "",
			Number:   "",
			Position: "",
			Access:   false,
		}
	}
	htmlContent, err := os.ReadFile("./web/account.html")
	if err != nil {
		http.Error(w, "Failed to read HTML file", http.StatusInternalServerError)
		log.Println(loggerString, "Не смог прочитать HTML file:", err)
		return
	}
	html := string(htmlContent)
	html = strings.Replace(html, "{{Name}}", user.Name, -1)
	html = strings.Replace(html, "{{SurName}}", user.Surname, -1)
	html = strings.Replace(html, "{{Position}}", user.Position, -1)
	html = strings.Replace(html, "{{Number}}", user.Number, -1)
	html = strings.Replace(html, "{{Email}}", user.Email, -1)
	var Access string
	if user.Access {
		Access = "Доступ есть"
	} else {
		Access = "Доступа нет, обратитесь к администратору"
	}
	html = strings.Replace(html, "{{Access}}", Access, -1)
	// Отправляем HTML-страницу клиенту
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)

}

func AnalyticsHandler(w http.ResponseWriter, r *http.Request, loggerString string) {
	session, _ := store.Get(r, "session-name")
	if session.Values["login"] != nil {
		log.Println(loggerString, session.Values["login"], "visited the analytics")
	} else {
		log.Println(loggerString, r.RemoteAddr, "visited the analytics")
	}
	tmpl, err := template.ParseFiles("web/analytics.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Internal Server Error:", err)
		return
	}
	tmpl.Execute(w, nil)
}

func NotificationsHandler(w http.ResponseWriter, r *http.Request, loggerString string) {
	session, _ := store.Get(r, "session-name")
	if session.Values["login"] != nil {
		log.Println(loggerString, session.Values["login"], "visited the notif")
	} else {
		log.Println(loggerString, r.RemoteAddr, "visited the notif")
	}
	// Логика обработки запроса для страницы уведомлений
	tmpl, err := template.ParseFiles("web/notifications.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(loggerString, "Internal Server Error:", err)
		return
	}
	tmpl.Execute(w, nil)
}

func Statuses(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	statuses, err := postgre.GetDataBaseHelper("status", db.Db)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(loggerString, "Internal Server Error:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}

func Region(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	region, err := postgre.GetDataBaseHelper("region", db.Db)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(loggerString, "Internal Server Error:", err)
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(region)
}
func Users(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	users, err := postgre.GetDataBaseHelper("users", db.Db)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(loggerString, "Internal Server Error:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
func Event(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	events, err := postgre.GetDataBaseHelper("event", db.Db)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(loggerString, "Internal Server Error:", err)
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	session.Values["user_id"] = nil
	session.Values["login"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "../", http.StatusFound)
}
func LoginHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {

	// Получаем данные из формы
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Проверяем пользователя в базе данных
	user, err := postgre.LoginValidate(db.Db, username)
	fmt.Println(username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":false,"message":"Invalid username"}`))
		log.Println(loggerString, "Not validate login :", err)
		return
	}

	// Сравниваем хэш пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":false,"message":"Invalid password"}`))
		log.Println(loggerString, "Invalid password:", err)
		return
	}

	// Получаем сессию
	session, err := store.Get(r, "session-name")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":false,"message":"Failed to create session"}`))
		log.Println(loggerString, "Failed to create session:", err)
		return
	}
	session.Values["user_id"] = user.ID
	session.Values["login"] = user.Login
	if session.Values["login"] != nil {
		log.Println(loggerString, session.Values["login"], "logged into the server")
	} else {
		log.Println(loggerString, r.RemoteAddr, "visited the loginHandler")
	}
	// Сохраняем сессию
	err = session.Save(r, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":false,"message":"Failed to save session"}`))
		log.Println(loggerString, "Failed to save session:", err)
		return
	}

	// Отправляем успешный ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success":true,"message":"Authentication successful", "redirect":"/"}`))
}

func RegisterHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	nameuser := r.FormValue("nameuser")
	surnameuser := r.FormValue("surnameuser")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	w.Header().Set("Content-Type", "application/json")

	count, err := postgre.LoginCountValidate(db.Db, username)
	if err != nil {
		http.Error(w, `{"success":false,"message":"Cannot validate user"}`, http.StatusBadRequest)
		log.Println(loggerString, "Cannot validate user:", err)
		return
	}

	if *count > 0 {
		http.Error(w, `{"success":false,"message":"Username already exists"}`, http.StatusBadRequest)
		log.Println(loggerString, "Username already exists:", err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"success":false,"message":"Cannot generate hash password"}`, http.StatusInternalServerError)
		log.Println(loggerString, "Cannot generate hash password:", err)
		return
	}

	checkStatusAdd, err := postgre.AddedNewUsers(db.Db, username, nameuser, surnameuser, email, phone, hashedPassword)
	if err != nil {
		http.Error(w, `{"success":false,"message":"Cannot add user"}`, http.StatusInternalServerError)
		log.Println(loggerString, "Cannot add user:", err)
		return
	}

	user, err := postgre.LoginValidate(db.Db, username)
	if err != nil {
		http.Error(w, `{"success":false,"message":"Failed to authorize user"}`, http.StatusUnauthorized)
		log.Println(loggerString, "Failed to authorize user:", err)
		return
	}

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, `{"success":false,"message":"Failed to create session"}`, http.StatusInternalServerError)
		log.Println(loggerString, "Failed to create session:", err)
		return
	}

	session.Values["user_id"] = user.ID
	session.Values["login"] = username

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, `{"success":false,"message":"Failed to save session"}`, http.StatusInternalServerError)
		log.Println(loggerString, "Failed to save session:", err)
		return
	}

	if checkStatusAdd == "ok" {
		w.Write([]byte(`{"success":true,"message":"Registration successful", "redirect":"/"}`))
		log.Println(loggerString, "All is good:", checkStatusAdd)
		return
	}

	w.Write([]byte(`{"success":false,"message":"Unknown error occurred"}`))
}

// SaveClientChangesHandler handles saving changes to the client info
func SaveClientChangesHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["user_id"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println(loggerString, " Unauthorized:", err)
		return
	}

	userID := session.Values["user_id"].(int)
	clientID, err := strconv.Atoi(r.FormValue("client_id"))
	if err != nil {
		log.Println(loggerString, "Я не смог преобразовать тип", err)
	}
	UID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		log.Println(loggerString, "Я не смог преобразовать тип", err)
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	status := r.FormValue("status")
	region := r.FormValue("region")
	event := r.FormValue("event")
	fCall := r.FormValue("fcall")
	nCall := r.FormValue("ncall")
	Uname := r.FormValue("user_name")
	TotalCash, err := strconv.Atoi(r.FormValue("total_cash"))
	if err != nil {
		log.Println(loggerString, "Я не смог преобразовать тип", err)
	}
	Payday := r.FormValue("payday")
	CheckName, _ := postgre.ValidateUsersToChanges(db.Db, session.Values["user_id"].(int))
	if !CheckName.Access {
		http.Error(w, "Нет прав доступа", http.StatusInternalServerError)
		log.Println(loggerString, "Нет прав доступа: ", err)
		return
	}
	ChangeName, err := postgre.GetIdFromName(db.Db, Uname)
	if ChangeName.ID != UID {
		err = postgre.UpdateClientInfo(db.Db, clientID, ChangeName.ID, TotalCash, name, email, phone, status, region, event, fCall, nCall, Uname, Payday)
	} else {
		err = postgre.UpdateClientInfo(db.Db, clientID, UID, TotalCash, name, email, phone, status, region, event, fCall, nCall, Uname, Payday)
	}
	if err != nil {
		http.Error(w, "Failed to save changes", http.StatusInternalServerError)
		log.Println(loggerString, "Failed to save changes:", err)
		return
	}
	UserNameToChangeDescrp, _ := postgre.LoginValidate(db.Db, session.Values["login"].(string))
	changeDescription := fmt.Sprintf("%v добавил(а) изменения ", UserNameToChangeDescrp.Name)
	err = postgre.LogClientChanges(db.Db, userID, clientID, changeDescription)
	if err != nil {
		http.Error(w, "Failed to log changes", http.StatusInternalServerError)
		log.Println(loggerString, "Failed to log changes:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetClientHistoryHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["user_id"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println(loggerString, "Unauthorized: ", err)
		return
	}

	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		http.Error(w, "Client ID is required", http.StatusBadRequest)
		log.Println(loggerString, "Client ID is required: ", err)
		return
	}

	// Fetch client history from the database (implement the function GetClientHistory)
	history, err := postgre.GetClientHistory(db.Db, clientID)
	if err != nil {
		http.Error(w, "Failed to fetch history", http.StatusInternalServerError)
		log.Println(loggerString, "Failed to fetch history: ", err)
		return
	}

	json.NewEncoder(w).Encode(history)
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println(loggerString, "Method not allowed: ")
		return
	}
	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["user_id"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println(loggerString, "Unauthorized: ", err)
		return
	}
	clientID, err := strconv.Atoi(r.FormValue("client_id"))
	if err != nil {
		log.Println(loggerString, "Не смог преобразовать: ", err)
	}
	UID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		log.Println(loggerString, "Не смог преобразовать: ", err)
	}
	description := r.FormValue("description")
	dueDate := r.FormValue("due_date")
	userName := r.FormValue("user_name")
	if description == "" {
		http.Error(w, "Пустой текст задачи: ", http.StatusBadRequest)
		log.Println(loggerString, "Пустой текст задачи: ", err)
		return
	}
	if dueDate == "" {
		http.Error(w, "Пустые время и дата: ", http.StatusBadRequest)
		log.Println(loggerString, "Пустые время и дата: ", err)
		return
	}
	UsersValidate, _ := postgre.ValidateUsersToChanges(db.Db, session.Values["user_id"].(int))
	if !UsersValidate.Access {
		http.Error(w, "Нет доступа: ", http.StatusBadRequest)
		log.Println(loggerString, "Нет доступа: ", err)
		return
	}
	err = postgre.AddTasksUsers(db.Db, UID, clientID, description, dueDate, userName)
	if err != nil {
		http.Error(w, "Cannot add task", http.StatusBadRequest)
		log.Println(loggerString, "Cannot add task: ", err)
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Task added successfully")
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["user_id"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println(loggerString, "Unauthorized: ", err)
		return
	}
	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		log.Println(loggerString, "User ID is required: ", err)
		return
	}

	tasks, err := postgre.GetTasksUsers(db.Db, clientID)
	if err != nil {
		log.Println(loggerString, "Не смог взять задачи: ", err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func GetTasksHandlerToNotif(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["user_id"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println(loggerString, "Unauthorized: ", err)
		return
	}
	userID := session.Values["user_id"].(int)

	tasks, err := postgre.GetTasksUsersToNotif(db.Db, userID)
	if err != nil {
		log.Println(loggerString, "Не смог взять задачи: ", err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func GetCallsDataHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	period := r.URL.Query().Get("period")
	data, err := postgre.GetCallsData(db.Db, period)
	if err != nil {
		http.Error(w, "Невозможно взять данные", http.StatusInternalServerError)
		log.Println(loggerString, "Невозможно взять данные: ", err)
		return
	}
	json.NewEncoder(w).Encode(data)
}
func GetSalesDataHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	period := r.URL.Query().Get("period")
	data, err := postgre.GetSalesData(db.Db, period)
	if err != nil {
		http.Error(w, "Невозможно взять данные", http.StatusInternalServerError)
		log.Println(loggerString, "Невозможно взять данные: ", err)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func ExecuteTaskHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		http.Error(w, "Missing task ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(taskID)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = postgre.DeleteTask(db.Db, id)
	if err != nil {
		log.Println(loggerString, "Failed to delete task:", err)
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task executed and deleted"))
}

func GetClientStatusHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	clientStatuses, err := postgre.GetClientStatus(db.Db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(loggerString, "Неверный запрос:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientStatuses)
}

func GetClientMarketingHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	marketingData, err := postgre.GetClientMarketing(db.Db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(loggerString, "Неверный запрос:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(marketingData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(loggerString, "Json NewEncoder err:", err)
		return
	}
}

func ClientsByRegionHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	clients, err := postgre.GetClientsByRegion(db.Db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(loggerString, "Неверный запрос:", err)
		return
	}

	// Преобразуйте данные в формат JSON
	jsonData, err := json.Marshal(clients)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(loggerString, "JsonMarshal err :", err)
		return
	}

	// Отправьте JSON-ответ обратно клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func GetCallsTodayHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	count, err := postgre.GetCallsToday(db.Db)
	if err != nil {
		log.Println(loggerString, "Неверный запрос:", err)
	}
	response := struct {
		Count int `json:"count"`
	}{
		Count: count,
	}

	// Кодируем структуру в JSON и отправляем клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func GetClientsByUser(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	results, err := postgre.GetClientsByUser(db.Db)
	if err != nil {
		log.Println(loggerString, "Неверный запрос:", err)
	}
	jsonBytes, err := json.Marshal(results)
	if err != nil {
		log.Println(loggerString, "Error encoding JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(loggerString, "Internal Server Error: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func CreateClientHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["user_id"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println(loggerString, "Unauthorized: ", err)
		return
	}

	UID := session.Values["user_id"].(int)
	if err != nil {
		log.Println(loggerString, "Не смог преобразовать: ", err)
	}
	UsersValidate, _ := postgre.ValidateUsersToChanges(db.Db, UID)
	if !UsersValidate.Access {
		http.Error(w, "Нет доступа: ", http.StatusBadRequest)
		log.Println(loggerString, "Нет доступа: ", err)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println(loggerString, "Method not allowed")
		return
	}

	var clientData postgre.Client
	if err := json.NewDecoder(r.Body).Decode(&clientData); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Println(loggerString, "Failed to decode request body: ", err)
		return
	}

	if err := postgre.SaveClientToDB(db.Db, clientData); err != nil {
		log.Println(loggerString, "Failed to save client data to database: ", err)
		http.Error(w, "Failed to save client data to database", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(struct {
		Message    string         `json:"message"`
		ClientData postgre.Client `json:"client_data"`
	}{
		Message:    "Client data saved successfully",
		ClientData: clientData,
	})
	if err != nil {
		log.Println(loggerString, "Failed to encode response: ", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func DeleteClientHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB, loggerString string) {
	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["user_id"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println(loggerString, "Unauthorized: ", err)
		return
	}
	UID := session.Values["user_id"].(int)
	if err != nil {
		log.Println(loggerString, "Не смог преобразовать: ", err)
	}
	UsersValidate, _ := postgre.ValidateUsersToChanges(db.Db, UID)
	if !UsersValidate.Access {
		http.Error(w, "Нет доступа: ", http.StatusBadRequest)
		log.Println(loggerString, "Нет доступа: ", err)
		return
	}
	clientID, err := strconv.Atoi(r.URL.Query().Get("client_id"))
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		log.Println(loggerString, "Invalid client ID: ", err)
		return
	}

	err = postgre.DeleteClient(db.Db, clientID)
	if err != nil {
		log.Println(loggerString, "Failed to delete client: ", err)
		http.Error(w, "Failed to delete client", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Client with ID %d deleted successfully", clientID)
}
