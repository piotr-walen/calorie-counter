package handlers

import (
	"app/service/middleware"
	"app/service/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func CreateEntry(db *sql.DB, logger *logrus.Logger) http.Handler {
	type RequestObject struct {
		Entry *models.Entry `json:"entry"`
	}
	type ResponseObject struct {
		Error string        `json:"error,omitempty"`
		Entry *models.Entry `json:"entry,omitempty"`
	}
	sendError := func(w http.ResponseWriter, status int, err error) {
		err = errors.Wrap(err, "While creating entry")
		logger.Error(err)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
		out := ResponseObject{
			Error: err.Error(),
		}
		json.NewEncoder(w).Encode(out)
	}
	sendData := func(w http.ResponseWriter, status int, entry *models.Entry) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
		out := ResponseObject{
			Entry: entry,
		}
		json.NewEncoder(w).Encode(out)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in := &RequestObject{}
		err := json.NewDecoder(r.Body).Decode(in)
		if err != nil {
			err = errors.Wrap(err, "While decoding request body")
			sendError(w, http.StatusBadRequest, err)
			return
		}
		if in.Entry == nil {
			err = errors.Wrap(err, "No entry provided")
			sendError(w, http.StatusBadRequest, err)
			return
		}
		entry := in.Entry
		userID, ok := r.Context().Value(middleware.UserID).(int)
		logger.Debugf("USER ID: %d", userID)
		if !ok {
			err = errors.New("While getting UserID from request context")
			sendError(w, http.StatusBadRequest, err)
			return
		}
		entry.UserID = userID
		dbEntry, err := models.CreateEntry(db, entry)
		if err != nil {
			err = errors.Wrap(err, "While creating db entry")
			sendError(w, http.StatusBadRequest, err)
			return
		}
		sendData(w, http.StatusOK, dbEntry)
		return
	})
}

func GetUsersEntries(db *sql.DB, logger *logrus.Logger) http.Handler {
	type Product struct {
		models.Product
		Portions []models.Portion `json:"portions,omitempty"`
	}
	type Entry struct {
		models.Entry
		Product Product `json:"product,omitempty"`
	}
	type RequestObject struct {
		Date time.Time `json:"date,omitempty"`
	}
	type ResponseObject struct {
		Error   string   `json:"error,omitempty"`
		Entries *[]Entry `json:"entries,omitempty"`
	}
	sendError := func(w http.ResponseWriter, status int, err error) {
		err = errors.Wrap(err, "While creating entry")
		logger.Error(err)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
		out := ResponseObject{
			Error: err.Error(),
		}
		json.NewEncoder(w).Encode(out)
	}
	sendData := func(w http.ResponseWriter, status int, entries *[]Entry) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
		out := ResponseObject{
			Entries: entries,
		}
		json.NewEncoder(w).Encode(out)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in := &RequestObject{}
		err := json.NewDecoder(r.Body).Decode(in)
		if err != nil {
			err = errors.Wrap(err, "While decoding request body")
			sendError(w, http.StatusBadRequest, err)
			return

		}
		userID, ok := r.Context().Value(middleware.UserID).(int)
		if !ok {
			err = errors.New("While getting UserID from request context")
			sendError(w, http.StatusBadRequest, err)
			return

		}
		entries, err := models.GetUsersEntries(db, userID, in.Date)
		if err != nil {
			err = errors.Wrap(err, "While getting db users entries")
			sendError(w, http.StatusBadRequest, err)
			return
		}
		popEntries := []Entry{}
		for _, entry := range *entries {
			productID := entry.ProductID
			product, err := models.GetProductById(db, productID)
			if err != nil {
				err = errors.Wrap(err, "While getting db product")
				sendError(w, http.StatusBadRequest, err)
				return
			}
			portions, err := models.GetProductsPortions(db, productID)
			if err != nil {
				err = errors.Wrap(err, "While getting db product portion")
				sendError(w, http.StatusBadRequest, err)
				return
			}
			populated := Product{
				Product:  *product,
				Portions: portions,
			}
			popEntry := Entry{
				Entry:   entry,
				Product: populated,
			}
			popEntries = append(popEntries, popEntry)
		}
		sendData(w, http.StatusOK, &popEntries)
		return
	})
}

func DeleteEntry(db *sql.DB, logger *logrus.Logger) http.Handler {
	type RequestObject struct {
		ID int `json:"id"`
	}
	type ResponseObject struct {
		Error string `json:"error,omitempty"`
	}
	sendError := func(w http.ResponseWriter, status int, err error) {
		err = errors.Wrap(err, "While creating entry")
		logger.Error(err)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
		out := ResponseObject{
			Error: err.Error(),
		}
		json.NewEncoder(w).Encode(out)
	}
	sendData := func(w http.ResponseWriter, status int) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
		out := ResponseObject{}
		json.NewEncoder(w).Encode(out)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in := &RequestObject{}
		err := json.NewDecoder(r.Body).Decode(in)
		if err != nil {
			err = errors.Wrap(err, "While decoding request body")
			sendError(w, http.StatusBadRequest, err)
			return
		}
		userID, ok := r.Context().Value(middleware.UserID).(int)
		if !ok {
			err = errors.New("While getting UserID from request context")
			sendError(w, http.StatusBadRequest, err)
			return

		}
		entry, err := models.GetEntry(db, in.ID)
		if entry.UserID != userID {
			err = errors.New("Permission denied, user id do not match")
			sendError(w, http.StatusUnauthorized, err)
			return
		}
		err = models.DeleteEntry(db, in.ID)
		if err != nil {
			err = errors.Wrap(err, "While deleting db users entry")
			sendError(w, http.StatusBadRequest, err)
			return
		}
		sendData(w, http.StatusOK)
		return
	})
}

func UpdateEntry(db *sql.DB, logger *logrus.Logger) http.Handler {
	type RequestObject struct {
		ID    int           `json:"id"`
		Entry *models.Entry `json:"entry"`
	}
	type ResponseObject struct {
		Error string `json:"error,omitempty"`
	}
	sendError := func(w http.ResponseWriter, status int, err error) {
		err = errors.Wrap(err, "While updating entry")
		logger.Error(err)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
		out := ResponseObject{
			Error: err.Error(),
		}
		json.NewEncoder(w).Encode(out)
	}
	sendData := func(w http.ResponseWriter, status int) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
		out := ResponseObject{}
		json.NewEncoder(w).Encode(out)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in := &RequestObject{}
		err := json.NewDecoder(r.Body).Decode(in)
		if err != nil {
			err = errors.Wrap(err, "While decoding request body")
			sendError(w, http.StatusBadRequest, err)
			return
		}
		userID, ok := r.Context().Value(middleware.UserID).(int)
		if !ok {
			err = errors.New("While getting UserID from request context")
			sendError(w, http.StatusBadRequest, err)
			return

		}
		entry, err := models.GetEntry(db, in.ID)
		if entry.UserID != userID {
			err = errors.New("Permission denied, user id do not match")
			sendError(w, http.StatusUnauthorized, err)
			return
		}
		err = models.UpdateEntry(db, in.ID, in.Entry)
		if err != nil {
			err = errors.Wrap(err, "While db update entry")
			sendError(w, http.StatusBadRequest, err)
			return
		}
		sendData(w, http.StatusOK)
		return
	})
}
