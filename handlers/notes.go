package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func HandleNotes(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		CreateNote(w, r)
	case http.MethodGet:
		GetAllNotes(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

}

func GetAllNotes(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	userID := session.Values["user_id"] // Извлекаем ID пользователя из сессии

	rows, err := db.Query("SELECT id, user_id, title, content FROM notes WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notes []Note

	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content); err != nil {
			http.Error(w, "Ошибка чтения данных заметки", http.StatusInternalServerError)
			return
		}
		notes = append(notes, note) // Добавляем заметку в список заметок пользователя
	}

	jsonResponse, err := json.Marshal(notes)
	log.Println(jsonResponse)
	if err != nil {
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// Обработчик для создания новой заметки (AJAX)
func CreateNote(w http.ResponseWriter, r *http.Request) {
	var note Note
	session, _ := store.Get(r, "session-name")
	userID := session.Values["user_id"]
	note.UserID = userID.(int)
	log.Println(userID, store, session)
	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Ошибка декодирования данных", http.StatusBadRequest)
		return
	}

	//session, _ := store.Get(r, "session-name")
	//userID := session.Values["user_id"]
	//note.UserID = session.Values["user_id"].(int) // Извлекаем ID пользователя из сессии

	err := db.QueryRow("INSERT INTO notes (user_id,title ,content ) VALUES ($1,$2,$3) RETURNING id",
		note.UserID, note.Title, note.Content).Scan(&note.ID)

	if err != nil {
		http.Error(w, "Ошибка создания заметки", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(note)
	if err != nil {
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func NoteByID(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Path[len("/notes/"):]
	log.Println(idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Ошибка конвертации ID ", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		GetNoteByID(w, r, id)
	case http.MethodPatch: // Patch - используется для частичного обловления ресурса на сервере

	}

}

func GetNoteByID(w http.ResponseWriter, r *http.Request, id int) {
	session, _ := store.Get(r, "session-name")
	userID := session.Values["user_id"]

	rows, err := db.Query("SELECT id, user_id, title, content FROM notes WHERE id = $1 and user_id = $2", id, userID)
	if err != nil {
		http.Error(w, "Ошибка получения данных2", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	if !rows.Next() {
		http.Error(w, "Заметка отсутствует", http.StatusNotFound)
		return
	}

	var note Note
	if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content); err != nil {
		http.Error(w, "Ошибка чтения данных заметки", http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(note)
	log.Println(jsonResponse)
	if err != nil {
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}

// Обработчик для обновления существующей заметки (AJAX)
func UpdateNote(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/notes/"):]
	log.Println(idStr)
	id, err := strconv.Atoi(idStr) // Конвертация строки в целое число
	if err != nil {
		http.Error(w, "Ошибка конвертации ID ", http.StatusBadRequest)
		return
	}
	log.Println(id)
	var note Note

	//if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
	//	http.Error(w, "Ошибка декодирования данных ", http.StatusBadRequest)
	//	return
	//}

	note.ID = id // Устанавливаем ID для обновления
	log.Println("qwe", note.ID)
	if _, err := db.Exec("UPDATE notes SET content=$1 WHERE id=$2", note.Content, id); err != nil {
		http.Error(w, "Ошибка обновления заметки ", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(note)
	log.Println(jsonResponse, err)
	if err != nil {
		http.Error(w, "Ошибка формирования ответа ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json ")
	w.Write(jsonResponse)
}

/*func GetAllNotes(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	var user User

	err := db.QueryRow("SELECT id FROM users WHERE user_name = $1", username).Scan(&user.UserID)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	rows, err := db.Query("SELECT id, user_id, title, content FROM notes WHERE user_id = $1", user.UserID)
	if err != nil {
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notes []Note

	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content); err != nil {
			http.Error(w, "Ошибка чтения данных заметки", http.StatusInternalServerError)
			return
		}
		note.UserID = user.UserID   // Устанавливаем UserID для заметки
		notes = append(notes, note) // Добавляем заметку в список заметок пользователя
	}

	jsonResponse, err := json.Marshal(notes)
	if err != nil {
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// Обработчик для создания новой заметки (AJAX)
func CreateNote(w http.ResponseWriter, r *http.Request) {
	var note Note

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Ошибка декодирования данных", http.StatusBadRequest)
		return
	}

	err := db.QueryRow("INSERT INTO notes (user_id, title, content) VALUES ($1, $2, $3) RETURNING id",
		note.UserID, note.Title, note.Content).Scan(&note.ID)

	if err != nil {
		http.Error(w, "Ошибка создания заметки", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(note)
	if err != nil {
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// Обработчик для обновления существующей заметки (AJAX)
func UpdateNote(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/notes/"):]

	var note Note

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Ошибка декодирования данных", http.StatusBadRequest)
		return
	}

	note.ID, _ = strconv.Atoi(id) // Устанавливаем ID для обновления

	if _, err := db.Exec("UPDATE notes SET title=$1, content=$2 WHERE id=$3",
		note.Title, note.Content, id); err != nil {
		http.Error(w, "Ошибка обновления заметки", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(note)
	if err != nil {
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func Notes(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	hashedPassword := r.URL.Query().Get("&")

	title := r.FormValue("title")
	color := r.FormValue("color")
	content := r.FormValue("content")

	// Сохраняем заметку в базе данных (добавьте свою логику сохранения)
	_, err := db.Exec("INSERT INTO notes (title, content, color, user_id) VALUES ($1, $2, $3, (SELECT id FROM users WHERE user_name = $4 and hashed_password = $5))", title, content, color, username, hashedPassword)
	if err != nil {
		http.Error(w, "Ошибка заметки или доступа", http.StatusInternalServerError)
		return
	}
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/notes.html")
		if err != nil {
			log.Println("Error loading template:", err)
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil) // Отображаем страницу заметок

	} else if r.Method == http.MethodPost {

		var note Note
		var user User
		user, _ = UserInformationFromDB(username)

		err := json.NewDecoder(r.Body).Decode(&note)
		if err != nil {
			http.Error(w, "Ошибка декодирования JSON", http.StatusBadRequest)
			return
		}

		note.UserID = user.UserID // Замените на ID текущего пользователя (например, из сессии)

		query := "INSERT INTO notes (user_id, title, content, color) VALUES ($1, $2, $3, $4)"
		_, err = db.Exec(query, note.UserID, note.Title, note.Content, note.Color)
		if err != nil {
			http.Error(w, "Ошибка сохранения заметки", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated) // Успешное создание заметки
		json.NewEncoder(w).Encode(note)   // Возвращаем созданную заметку в формате JSON

	}
} */
