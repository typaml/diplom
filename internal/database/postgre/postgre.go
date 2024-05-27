package postgre

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
	"time"
)

type PostgreClientDB struct {
	Db *sql.DB
}
type Client struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	Status        string `json:"status"`
	Region        string `json:"region_name"`
	Event         string `json:"marketing_event"`
	DateFirstCall string `json:"date_first_call"`
	DateNextCall  string `json:"date_next_call"`
	UserID        int    `json:"user_id"`
	UsersName     string `json:"users_name"`
	Payday        string `json:"payday"`
	Total         int    `json:"total_cash"`
}
type HistoryEntry struct {
	ChangeDate        string `json:"change_date"`
	ChangeDescription string `json:"change_description"`
}

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	CreatedAt   string `json:"created_at"`
	ClientID    string `json:"client_id"`
	UserID      string `json:"user_id"`
	UserName    string `json:"username"`
}
type User struct {
	ID       int
	Login    string
	Password string
	Name     string
	Surname  string
	Access   bool
}
type DataBaseHelper struct {
	Name string `json:"name"`
}

func NewDB() (*PostgreClientDB, error) {
	psqlInfo := "postgresql://postgres:rfyQdUUnyLmsPiXKEfReoTDhpPWlPChw@viaduct.proxy.rlwy.net:16561/railway"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		//panic(err)
	}

	// Проверяем, что соединение действительно установлено.
	err = db.Ping()
	if err != nil {
		//panic(err)
		fmt.Println("Не сделалось")
	}

	fmt.Println("Successfully connected to the database")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS public.clients (id SERIAL PRIMARY KEY, name TEXT, phone TEXT, email TEXT, status TEXT, region_name TEXT, user_id INTEGER, marketing_event TEXT, date_first_call TEXT, date_next_call TEXT, users_name TEXT, total INTEGER, payday TEXT);"); err != nil {
		return nil, err
	} else {
		fmt.Println("Table clients created successfully")
	}
	return &PostgreClientDB{Db: db}, nil
}

func (db *PostgreClientDB) Close() {
	db.Db.Close()
}

func Select(db *sql.DB, idclients int) (*Client, error) {
	var client Client
	if err := db.QueryRow("SELECT id , name, phone, email, status, region_name, user_id, marketing_event, date_first_call, date_next_call, users_name, payday, total FROM public.clients WHERE id = $1", idclients).Scan(&client.ID, &client.Name, &client.Phone, &client.Email, &client.Status, &client.Region, &client.UserID, &client.Event, &client.DateFirstCall, &client.DateNextCall, &client.UsersName, &client.Payday, &client.Total); err != nil {
		return nil, err
	}

	return &client, nil
}

func SortAndFilterClients(db *sql.DB, sortColumn, sortOrder, filterId, filterStatus, filterRegion, filterUsers, filterEvent string) ([]Client, error) {
	// Формируем SQL-запрос с учетом параметров сортировки и фильтрации
	query := "SELECT id, name, phone, email, status, region_name, user_id, marketing_event, date_first_call, date_next_call, users_name, total, payday  FROM public.clients"

	// Добавляем условия фильтрации
	var filters []string
	if filterId != "" {
		filters = append(filters, fmt.Sprintf("id = %s", filterId))
	}
	if filterStatus != "" {
		filters = append(filters, fmt.Sprintf("status = '%s'", filterStatus))
	}
	if filterRegion != "" {
		filters = append(filters, fmt.Sprintf("region_name = '%v'", filterRegion))
	}
	if filterUsers != "" {
		filters = append(filters, fmt.Sprintf("users_name = '%s'", filterUsers))
	}
	if filterEvent != "" {
		filters = append(filters, fmt.Sprintf("marketing_event = '%v'", filterEvent))
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	// Добавляем условия сортировки
	query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, sortOrder)
	// Выполняем запрос к базе данных
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Считываем результаты запроса и добавляем их в слайс клиентов
	var clients []Client
	for rows.Next() {
		var client Client
		if err := rows.Scan(&client.ID, &client.Name, &client.Phone, &client.Email, &client.Status, &client.Region, &client.UserID, &client.Event, &client.DateFirstCall, &client.DateNextCall, &client.UsersName, &client.Total, &client.Payday); err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}

func GetDataBaseHelper(DataBaseLoad string, db *sql.DB) ([]DataBaseHelper, error) {
	// Выполняем запрос к базе данных для получения статусов
	query := fmt.Sprintf("SELECT name FROM public.%s", DataBaseLoad)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Создаем слайс для хранения значений статусов
	var DataBaseHelpers []DataBaseHelper

	// Итерируем по результатам запроса и добавляем значения в слайс
	for rows.Next() {
		var DataHelper DataBaseHelper
		err := rows.Scan(&DataHelper.Name)
		if err != nil {
			return nil, err
		}
		DataBaseHelpers = append(DataBaseHelpers, DataHelper)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return DataBaseHelpers, nil
}

func LoginValidate(db *sql.DB, login string) (*User, error) {

	// Поиск пользователя в базе данных
	var user User
	err := db.QueryRow("SELECT id, login, password, name, surname FROM public.users WHERE login=$1", login).Scan(&user.ID, &user.Login, &user.Password, &user.Name, &user.Surname)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}
func GetIdFromName(db *sql.DB, name string) (*User, error) {

	// Поиск пользователя в базе данных
	var user User
	err := db.QueryRow("SELECT id, login, password, name, surname, access FROM public.users WHERE name=$1", name).Scan(&user.ID, &user.Login, &user.Password, &user.Name, &user.Surname, &user.Access)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

func ValidateUsersToChanges(db *sql.DB, id int) (*User, error) {

	// Поиск пользователя в базе данных
	var user User
	err := db.QueryRow("SELECT id, login, password, name, surname, access FROM public.users WHERE id=$1", id).Scan(&user.ID, &user.Login, &user.Password, &user.Name, &user.Surname, &user.Access)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

func LoginCountValidate(db *sql.DB, login string) (*int, error) {

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM public.users WHERE login=$1", login).Scan(&count)
	if err != nil {
		return nil, err
	}
	return &count, nil
}

func AddedNewUsers(db *sql.DB, login, name, surname string, password []byte) (string, error) {
	_, err := db.Exec("INSERT INTO public.users (login, password, name, surname) VALUES ($1, $2, $3, $4)", login, string(password), name, surname)
	if err != nil {
		return "", err
	}
	return "ok", nil
}

func UpdateClientInfo(db *sql.DB, clientID, UID, total int, name, email, phone, status, region, event, fCall, nCall, Uname, Payday string) error {
	query := `
		UPDATE public.clients
		SET name = $1, email = $2, phone = $3, status = $4, region_name = $5, marketing_event = $6, date_first_call = $7, date_next_call = $8, user_id = $9, users_name = $11, payday = $12, total = $13
		WHERE id = $10
	`
	_, err := db.Exec(query, name, email, phone, status, region, event, fCall, nCall, UID, clientID, Uname, Payday, total)
	return err
}

func LogClientChanges(db *sql.DB, userID, clientID int, changeDescription string) error {
	query := `
		INSERT INTO client_history (client_id, user_id, change_description)
		VALUES ($1, $2, $3)
	`
	_, err := db.Exec(query, clientID, userID, changeDescription)
	return err
}

func GetClientHistory(db *sql.DB, clientID string) ([]HistoryEntry, error) {
	query := `
		SELECT change_date, change_description
		FROM public.client_history
		WHERE client_id = $1
		ORDER BY change_date DESC
	`
	rows, err := db.Query(query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []HistoryEntry
	for rows.Next() {
		var entry HistoryEntry
		if err := rows.Scan(&entry.ChangeDate, &entry.ChangeDescription); err != nil {
			return nil, err
		}
		history = append(history, entry)
	}
	return history, nil
}

func AddTasksUsers(db *sql.DB, userID, clientID int, description, data, userName string) error {
	_, err := db.Exec("INSERT INTO tasks (user_id, description, due_date, client_id, user_name) VALUES ($1, $2, $3, $4, $5)", userID, description, data, clientID, userName)
	if err != nil {
		return err
	}
	return nil
}

func GetTasksUsers(db *sql.DB, clientID string) ([]Task, error) {
	CID, _ := strconv.Atoi(clientID)
	rows, err := db.Query("SELECT id, description, due_date, created_at, client_id, user_id, user_name FROM public.tasks WHERE client_id = $1 ORDER BY created_at DESC", CID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Description, &task.DueDate, &task.CreatedAt, &task.ClientID, &task.UserID, &task.UserName); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func GetTasksUsersToNotif(db *sql.DB, userID int) ([]Task, error) {
	rows, err := db.Query("SELECT id, description, due_date, created_at, client_id, user_id, user_name FROM public.tasks WHERE user_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Description, &task.DueDate, &task.CreatedAt, &task.ClientID, &task.UserID, &task.UserName); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

type ChartData struct {
	Labels   []string  `json:"labels"`
	Datasets []Dataset `json:"datasets"`
}

type Dataset struct {
	Label           string `json:"label"`
	Data            []int  `json:"data"`
	BackgroundColor string `json:"backgroundColor"`
	BorderColor     string `json:"borderColor"`
	BorderWidth     int    `json:"borderWidth"`
}

func GetCallsData(db *sql.DB, period string) (ChartData, error) {
	var query string
	var rows *sql.Rows
	var err error

	// Определяем запрос в зависимости от выбранного периода
	switch period {
	case "day":
		query = "SELECT date_first_call, COUNT(*) FROM public.clients WHERE date_first_call IS NOT NULL AND date_first_call != 'none' AND date_trunc('day', TO_TIMESTAMP(date_first_call, 'YYYY-MM-DD\"T\"HH24:MI')) = CURRENT_DATE GROUP BY date_first_call ORDER BY date_first_call"
	case "week":
		query = "SELECT date_first_call, COUNT(*) FROM public.clients WHERE date_first_call IS NOT NULL AND date_first_call != 'none' AND date_trunc('day', TO_TIMESTAMP(date_first_call, 'YYYY-MM-DD\"T\"HH24:MI')) BETWEEN CURRENT_DATE - INTERVAL '1 week' AND CURRENT_DATE GROUP BY date_first_call ORDER BY date_first_call"
	case "month":
		query = "SELECT date_first_call, COUNT(*) FROM public.clients WHERE date_first_call IS NOT NULL AND date_first_call != 'none' AND date_trunc('day', TO_TIMESTAMP(date_first_call, 'YYYY-MM-DD\"T\"HH24:MI')) BETWEEN DATE_TRUNC('month', CURRENT_DATE) AND CURRENT_DATE GROUP BY date_first_call ORDER BY date_first_call"
	case "year":
		query = "SELECT date_first_call, COUNT(*) FROM public.clients WHERE date_first_call IS NOT NULL AND date_first_call != 'none' AND date_trunc('day', TO_TIMESTAMP(date_first_call, 'YYYY-MM-DD\"T\"HH24:MI')) BETWEEN DATE_TRUNC('year', CURRENT_DATE) AND CURRENT_DATE GROUP BY date_first_call ORDER BY date_first_call"
	default:
		return ChartData{}, fmt.Errorf("Invalid period")
	}

	// Выполняем запрос к базе данных
	rows, err = db.Query(query)
	if err != nil {
		return ChartData{}, err
	}
	defer rows.Close()

	// Формируем данные для графика
	var labels []string
	var data []int
	for rows.Next() {
		var dateStr string
		var count int
		if err := rows.Scan(&dateStr, &count); err != nil {
			return ChartData{}, err
		}
		labels = append(labels, dateStr)
		data = append(data, count)
	}

	return ChartData{
		Labels: labels,
		Datasets: []Dataset{
			{
				Label:           "Звонки",
				Data:            data,
				BackgroundColor: "rgba(255, 99, 132, 0.2)",
				BorderColor:     "rgba(255, 99, 132, 1)",
				BorderWidth:     1,
			},
		},
	}, nil
}
func GetSalesData(db *sql.DB, period string) (ChartData, error) {
	var query string
	var rows *sql.Rows
	var err error

	// Определяем запрос в зависимости от выбранного периода
	switch period {
	case "day":
		query = "SELECT payday, SUM(total) FROM public.clients WHERE TRIM(payday) != '' AND COALESCE(NULLIF(TRIM(payday), ''), '')::DATE = CURRENT_DATE GROUP BY payday ORDER BY payday"
	case "week":
		query = "SELECT payday, SUM(total) FROM public.clients WHERE TRIM(payday) != '' AND COALESCE(NULLIF(TRIM(payday), ''), '')::DATE BETWEEN CURRENT_DATE - INTERVAL '1 week' AND CURRENT_DATE GROUP BY payday ORDER BY payday"
	case "month":
		query = "SELECT payday, SUM(total) FROM public.clients WHERE TRIM(payday) != '' AND EXTRACT(MONTH FROM COALESCE(NULLIF(TRIM(payday), ''), '')::DATE) = EXTRACT(MONTH FROM CURRENT_DATE) AND EXTRACT(YEAR FROM COALESCE(NULLIF(TRIM(payday), ''), '')::DATE) = EXTRACT(YEAR FROM CURRENT_DATE) GROUP BY payday ORDER BY payday"
	case "year":
		query = "SELECT payday, SUM(total) FROM public.clients WHERE TRIM(payday) != '' AND EXTRACT(YEAR FROM COALESCE(NULLIF(TRIM(payday), ''), '')::DATE) = EXTRACT(YEAR FROM CURRENT_DATE) GROUP BY payday ORDER BY payday"
	default:
		query = "SELECT payday, SUM(total) FROM public.clients WHERE TRIM(payday) != '' GROUP BY payday ORDER BY payday"
	}

	// Выполняем запрос к базе данных
	rows, err = db.Query(query)
	if err != nil {
		return ChartData{}, err
	}
	defer rows.Close()

	// Формируем данные для графика
	var labels []string
	var data []int
	for rows.Next() {
		var dateStr sql.NullString
		var total int
		if err := rows.Scan(&dateStr, &total); err != nil {
			return ChartData{}, err
		}
		if dateStr.Valid {
			labels = append(labels, dateStr.String)
		} else {
			labels = append(labels, "")
		}
		data = append(data, total)
	}

	return ChartData{
		Labels: labels,
		Datasets: []Dataset{
			{
				Label:           "Продажи",
				Data:            data,
				BackgroundColor: "rgba(54, 162, 235, 0.2)",
				BorderColor:     "rgba(54, 162, 235, 1)",
				BorderWidth:     1,
			},
		},
	}, nil
}

func DeleteTask(db *sql.DB, taskID int) error {
	query := "DELETE FROM tasks WHERE id = $1"
	_, err := db.Exec(query, taskID)
	return err
}

type ClientStatus struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

func GetClientStatus(db *sql.DB) ([]ClientStatus, error) {
	rows, err := db.Query("SELECT status, COUNT(*) FROM clients GROUP BY status")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ClientStatus
	for rows.Next() {
		var cs ClientStatus
		if err := rows.Scan(&cs.Status, &cs.Count); err != nil {
			return nil, err
		}
		results = append(results, cs)
	}

	return results, nil
}

type MarketingData struct {
	MarketingEvent string `json:"marketing_event"`
	Count          int    `json:"count"`
}

func GetClientMarketing(db *sql.DB) ([]MarketingData, error) {
	rows, err := db.Query("SELECT marketing_event, COUNT(*) FROM clients GROUP BY marketing_event")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var marketingData []MarketingData
	for rows.Next() {
		var data MarketingData
		if err := rows.Scan(&data.MarketingEvent, &data.Count); err != nil {
			return nil, err
		}
		marketingData = append(marketingData, data)
	}

	return marketingData, nil
}

type RegionClient struct {
	Region string `json:"region_name"`
	Count  int    `json:"count"`
}

func GetClientsByRegion(db *sql.DB) ([]RegionClient, error) {
	rows, err := db.Query("SELECT region_name, COUNT(*) FROM clients GROUP BY region_name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []RegionClient
	for rows.Next() {
		var client RegionClient
		if err := rows.Scan(&client.Region, &client.Count); err != nil {
			return nil, err
		}
		results = append(results, client)
	}

	return results, nil
}

func GetCallsToday(db *sql.DB) (int, error) {
	today := time.Now()
	dateString := today.Format("2006-01-02")

	// Запрос к базе данных для получения количества звонков за сегодня
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM clients WHERE date_first_call::text LIKE $1", dateString+"%").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

type UserClients struct {
	UserName    string
	ClientCount int
}

func GetClientsByUser(db *sql.DB) ([]UserClients, error) {
	rows, err := db.Query("SELECT users_name, COUNT(*) as client_count FROM clients GROUP BY users_name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []UserClients
	for rows.Next() {
		var uc UserClients
		if err := rows.Scan(&uc.UserName, &uc.ClientCount); err != nil {
			return nil, err
		}
		results = append(results, uc)
	}

	return results, nil
}

func SaveClientToDB(db *sql.DB, clientData Client) error {
	_, err := db.Exec("INSERT INTO clients (name, email, phone, status, region_name, marketing_event) VALUES ($1, $2, $3, $4, $5, $6)", clientData.Name, clientData.Email, clientData.Phone, clientData.Status, clientData.Region, clientData.Event)
	return err

}

func DeleteClient(db *sql.DB, clientID int) error {

	_, err := db.Exec("DELETE FROM public.clients WHERE id = $1", clientID)
	if err != nil {
		return nil
	}

	return nil
}
