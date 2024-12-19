package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// HandleNotes обработчик методов
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

// GetAllNotes получение всех данных о заметках пользователя
func GetAllNotes(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	userID := session.Values["user_id"]

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
		notes = append(notes, note)
	}

	jsonResponse, err := json.Marshal(notes)
	if err != nil {
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// CreateNote Обработчик для создания новой заметки
func CreateNote(w http.ResponseWriter, r *http.Request) {
	var note Note
	session, _ := store.Get(r, "session-name")
	note.UserID = session.Values["user_id"].(int)

	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Ошибка декодирования данных", http.StatusBadRequest)
		return
	}

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

// HandleNoteByID обработчик методов для одной заметки (по ID)
func HandleNoteByID(w http.ResponseWriter, r *http.Request) {
	idNoteStr := r.URL.Path[len("/notes/"):]
	idNote, err := strconv.Atoi(idNoteStr)

	if err != nil {
		http.Error(w, "Ошибка конвертации ID ", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		GetNoteByID(w, r, idNote)
	case http.MethodPatch: // Patch - используется для частичного обловления ресурса на сервере
		UpdateNoteById(w, r, idNote)
	case http.MethodDelete:
		DeleteNoteByID(w, r, idNote)
	}

}

// GetNoteByID получение всей информации конкретной заметки
func GetNoteByID(w http.ResponseWriter, r *http.Request, id int) {
	session, _ := store.Get(r, "session-name")
	userID := session.Values["user_id"]

	rows, err := db.Query("SELECT id, user_id, title, content FROM notes WHERE id = $1 and user_id = $2", id, userID)
	if err != nil {
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
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
	if err != nil {
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}

// UpdateNoteById Обработчик для обновления существующей заметки
func UpdateNoteById(w http.ResponseWriter, r *http.Request, id int) {
	var note Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Ошибка декодирования данных ", http.StatusBadRequest)
		return
	}

	note.ID = id

	_, err := db.Exec("UPDATE notes SET content=$1 WHERE id=$2", note.Content, id)
	if err != nil {
		http.Error(w, "Ошибка обновления заметки ", http.StatusInternalServerError)
		return
	}
}

// DeleteNoteByID удаление конкретной заметки
func DeleteNoteByID(w http.ResponseWriter, r *http.Request, id int) {
	session, _ := store.Get(r, "session-name")
	userID := session.Values["user_id"]

	rows, err := db.Query("SELECT id, user_id, title, content FROM notes WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	if !rows.Next() {
		return
	}

	var note Note
	if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content); err != nil {
		http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
		return
	}

	if note.UserID != userID {
		http.Error(w, "Ошибка доступа", http.StatusForbidden)
		return
	}

	_, err = db.Exec("DELETE FROM notes WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Ошибка удаления записи", http.StatusInternalServerError)
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
