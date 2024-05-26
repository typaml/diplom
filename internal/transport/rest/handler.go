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

func GetClients(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
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
		fmt.Println("Failed to get clients from database", http.StatusInternalServerError)
		return
	}
	// Преобразовываем данные в формат JSON и отправляем клиенту
	jsonData, err := json.Marshal(clients)
	if err != nil {
		fmt.Println("Cannot marshal to json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

// IndexHandler обрабатывает запросы к главной странице
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr)
	tmpl, err := template.ParseFiles("web/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func ClientsHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {

	tmpl, err := template.ParseFiles("web/clients.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func ClientHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	idStr := r.URL.Path[len("/client/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}
	client, err := postgre.Select(db.Db, id)
	if err != nil {
		http.Error(w, "Failed to get client from database", http.StatusInternalServerError)
		return
	}
	htmlContent, err := os.ReadFile("./web/client.html")
	if err != nil {
		http.Error(w, "Failed to read HTML file", http.StatusInternalServerError)
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

	// Отправляем HTML-страницу клиенту
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}
func AccountHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		fmt.Println("Не смог взять сессию")
	}
	var login string
	if session.Values["login"] != nil {
		var ok bool
		login, ok = session.Values["login"].(string)
		if !ok {
			fmt.Println("Не смог преобразовать логин в стрингу")
		}
	} else {
		login = ""
	}
	user, err := postgre.LoginValidate(db.Db, login)
	if err != nil {
		user = &postgre.User{
			ID:       0,
			Login:    "",
			Password: "",
			Name:     "",
			Surname:  "",
		}
	}
	htmlContent, err := os.ReadFile("./web/account.html")
	if err != nil {
		http.Error(w, "Failed to read HTML file", http.StatusInternalServerError)
		return
	}
	html := string(htmlContent)
	html = strings.Replace(html, "{{Name}}", user.Name, -1)
	html = strings.Replace(html, "{{SurName}}", user.Surname, -1)
	// Отправляем HTML-страницу клиенту
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)

}

func AnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	// Логика обработки запроса для страницы аккаунта
	tmpl, err := template.ParseFiles("web/analytics.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	// Логика обработки запроса для страницы уведомлений
	tmpl, err := template.ParseFiles("web/notifications.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func Statuses(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	statuses, err := postgre.GetDataBaseHelper("status", db.Db)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}

func Region(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	region, err := postgre.GetDataBaseHelper("region", db.Db)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(region)
}
func Users(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	users, err := postgre.GetDataBaseHelper("users", db.Db)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func Event(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	events, err := postgre.GetDataBaseHelper("event", db.Db)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func LoginHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {

	// Парсим данные формы
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form data:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Получаем данные из формы
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Проверяем пользователя в базе данных
	user, err := postgre.LoginValidate(db.Db, username)
	if err != nil {
		http.Error(w, "Not validate login or password", http.StatusUnauthorized)
		return
	}

	// Сравниваем хэш пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Получаем сессию
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}
	session.Values["user_id"] = user.ID
	session.Values["login"] = user.Login

	// Сохраняем сессию
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ клиенту
	http.Redirect(w, r, "/account", http.StatusFound)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	session.Values["user_id"] = nil
	session.Values["login"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "../", http.StatusFound)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	nameuser := r.FormValue("nameuser")
	surnameuser := r.FormValue("surnameuser")
	count, err := postgre.LoginCountValidate(db.Db, username)
	if err != nil {
		http.Error(w, "Cannot validate user", http.StatusBadRequest)
		return
	}

	if *count > 0 {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}
	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Cannot generate hash password", http.StatusInternalServerError)
		return
	}
	checkStatusAdd, err := postgre.AddedNewUsers(db.Db, username, nameuser, surnameuser, hashedPassword)
	if err != nil {
		http.Error(w, "Cannot added users", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	user, err := postgre.LoginValidate(db.Db, username)
	if err != nil {
		http.Error(w, "Failed to authorized users", http.StatusUnauthorized)
	}
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}
	session.Values["user_id"] = user.ID
	session.Values["login"] = username

	// Сохраняем сессию
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}
	if checkStatusAdd == "ok" {
		http.Redirect(w, r, "/account", http.StatusFound)
		return
	}
}

// SaveClientChangesHandler handles saving changes to the client info
func SaveClientChangesHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["user_id"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := session.Values["user_id"].(int)
	clientID, err := strconv.Atoi(r.FormValue("client_id"))
	if err != nil {
		fmt.Println("Я не смог преобразовать тип", clientID)
	}
	UID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		fmt.Println("Я не смог преобразовать тип", UID)
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
		fmt.Println("Я не смог преобразовать тип", TotalCash)
	}
	Payday := r.FormValue("payday")
	ChangeName, _ := postgre.GetIdFromName(db.Db, Uname)
	if ChangeName.ID != UID {
		err = postgre.UpdateClientInfo(db.Db, clientID, ChangeName.ID, TotalCash, name, email, phone, status, region, event, fCall, nCall, Uname, Payday)
	} else {
		err = postgre.UpdateClientInfo(db.Db, clientID, UID, TotalCash, name, email, phone, status, region, event, fCall, nCall, Uname, Payday)
	}
	if err != nil {
		http.Error(w, "Failed to save changes", http.StatusInternalServerError)
		return
	}
	UserNameToChangeDescrp, _ := postgre.LoginValidate(db.Db, session.Values["login"].(string))
	changeDescription := fmt.Sprintf("%v добавил(а) изменения ", UserNameToChangeDescrp.Name)
	err = postgre.LogClientChanges(db.Db, userID, clientID, changeDescription)
	if err != nil {
		http.Error(w, "Failed to log changes", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetClientHistoryHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["user_id"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		http.Error(w, "Client ID is required", http.StatusBadRequest)
		return
	}

	// Fetch client history from the database (implement the function GetClientHistory)
	history, err := postgre.GetClientHistory(db.Db, clientID)
	if err != nil {
		http.Error(w, "Failed to fetch history", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(history)
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["user_id"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	clientID, err := strconv.Atoi(r.FormValue("client_id"))
	if err != nil {
		fmt.Println("Я не смог преобразовать тип", clientID)
	}
	UID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		fmt.Println("Я не смог преобразовать тип", UID)
	}
	description := r.FormValue("description")
	dueDate := r.FormValue("due_date")
	userName := r.FormValue("user_name")
	err = postgre.AddTasksUsers(db.Db, UID, clientID, description, dueDate, userName)
	if err != nil {
		http.Error(w, "Cannot add task", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Task added successfully")
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["user_id"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	tasks, err := postgre.GetTasksUsers(db.Db, clientID)
	if err != nil {
		fmt.Println("Что-то пошло не так")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func GetTasksHandlerToNotif(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["user_id"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := session.Values["user_id"].(int)

	tasks, err := postgre.GetTasksUsersToNotif(db.Db, userID)
	if err != nil {
		fmt.Println("Что-то пошло не так")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func GetCallsDataHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	period := r.URL.Query().Get("period")
	data, err := postgre.GetCallsData(db.Db, period)
	if err != nil {
		http.Error(w, "Невозможно взять данные", http.StatusInternalServerError)
		fmt.Println(err, "ошибка в звонках")
		return
	}
	json.NewEncoder(w).Encode(data)
}
func GetSalesDataHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	period := r.URL.Query().Get("period")
	data, err := postgre.GetSalesData(db.Db, period)
	if err != nil {
		http.Error(w, "Невозможно взять данные", http.StatusInternalServerError)
		fmt.Println(err, "ошибка в продажах")
		return
	}
	json.NewEncoder(w).Encode(data)
}

func ExecuteTaskHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
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
		log.Println("Failed to delete task:", err)
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task executed and deleted"))
}

func GetClientStatusHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	clientStatuses, err := postgre.GetClientStatus(db.Db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Запрос неправильный")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientStatuses)
}

func GetClientMarketingHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	marketingData, err := postgre.GetClientMarketing(db.Db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Запрос неправильный", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(marketingData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ClientsByRegionHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	clients, err := postgre.GetClientsByRegion(db.Db) // Здесь db - ваше подключение к базе данных
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Запрос неправильный")
		return
	}

	// Преобразуйте данные в формат JSON
	jsonData, err := json.Marshal(clients)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправьте JSON-ответ обратно клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func GetCallsTodayHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	count, err := postgre.GetCallsToday(db.Db)
	if err != nil {
		fmt.Println("Запрос неправильный")
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
func GetClientsByUser(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	results, err := postgre.GetClientsByUser(db.Db)
	if err != nil {
		fmt.Println("Запрос неправильный", err)
	}
	jsonBytes, err := json.Marshal(results)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func CreateClientHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var clientData postgre.Client
	if err := json.NewDecoder(r.Body).Decode(&clientData); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		fmt.Println("Я не смог раскодировать json")
		return
	}

	if err := postgre.SaveClientToDB(db.Db, clientData); err != nil {
		fmt.Println("Я не смог сохранить", err)
		http.Error(w, "Failed to save client data to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Client data saved successfully:\n%+v", clientData)
}

func DeleteClientHandler(w http.ResponseWriter, r *http.Request, db *postgre.PostgreClientDB) {
	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["user_id"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	clientID, err := strconv.Atoi(r.URL.Query().Get("client_id"))
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	err = postgre.DeleteClient(db.Db, clientID)
	if err != nil {
		fmt.Println("Error deleting client:", err)
		http.Error(w, "Failed to delete client", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Client with ID %d deleted successfully", clientID)
}
