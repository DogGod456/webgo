package handlers

import (
	"net/http"
)

func deleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем сессию
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Очистка всех значений в сессии
	session.Options.MaxAge = -1 // Устанавливаем MaxAge в -1, чтобы удалить cookie
	err = session.Save(r, w)    // Сохраняем изменения
	if err != nil {
		http.Error(w, "Unable to delete session", http.StatusInternalServerError)
		return
	}
}
