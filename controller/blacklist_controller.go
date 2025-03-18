package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"product-controller/models"
	"product-controller/repository"
	"strconv"
)

type BlacklistController struct {
	BlacklistRepo *repository.BlacklistRepository
}

func NewBlacklistController(blacklistRepo *repository.BlacklistRepository) *BlacklistController {
	return &BlacklistController{
		BlacklistRepo: blacklistRepo,
	}
}

func (c *BlacklistController) GetAllBlacklistWords(w http.ResponseWriter, r *http.Request) {
	words, err := c.BlacklistRepo.GetAllBlacklistWords()
	if err != nil {
		http.Error(w, "Błąd pobierania blacklisty", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(words)
}

func (c *BlacklistController) AddBlacklistWord(w http.ResponseWriter, r *http.Request) {
	var word models.BlacklistWord
	err := json.NewDecoder(r.Body).Decode(&word)
	if err != nil {
		http.Error(w, "Niepoprawne dane wejściowe", http.StatusBadRequest)
		return
	}

	if word.Word == "" {
		http.Error(w, "Pole 'Word' jest wymagane", http.StatusBadRequest)
		return
	}

	err = c.BlacklistRepo.AddBlacklistWord(&word)
	if err != nil {
		http.Error(w, "Błąd dodawania słowa do blacklisty", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(word)
}

func (c *BlacklistController) DeleteBlacklistWord(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "Nieprawidłowe ID", http.StatusBadRequest)
		return
	}

	err = c.BlacklistRepo.DeleteBlacklistWord(uint(id))
	if err != nil {
		http.Error(w, "Błąd usuwania słowa z blacklisty", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
