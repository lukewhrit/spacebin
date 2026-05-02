package server

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/lukewhrit/spacebin/internal/util"
	"golang.org/x/exp/slices"
)

func (s *Server) StaticEditPage(w http.ResponseWriter, r *http.Request) {
	if !s.Config.AccountsEnabled {
		util.RenderError(&resources, w, http.StatusNotFound, errors.New("accounts disabled"))
		return
	}

	username, err := s.authenticatedUsername(r)
	if err != nil || username == "" {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	id := chi.URLParam(r, "document")
	if len(id) != s.Config.IDLength && !slices.Contains(s.Config.Documents, id) {
		util.RenderError(&resources, w, http.StatusBadRequest, fmt.Errorf("invalid document id"))
		return
	}

	doc, err := s.Database.GetDocument(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.RenderError(&resources, w, http.StatusNotFound, err)
			return
		}
		util.RenderError(&resources, w, http.StatusInternalServerError, err)
		return
	}

	if doc.Username != username {
		util.RenderError(&resources, w, http.StatusForbidden, errors.New("you do not own this document"))
		return
	}

	t, err := template.ParseFS(resources, "web/edit.html")
	if err != nil {
		util.RenderError(&resources, w, http.StatusInternalServerError, err)
		return
	}

	err = t.Execute(w, map[string]any{
		"Analytics":       s.Config.Analytics,
		"AccountsEnabled": s.Config.AccountsEnabled,
		"Authenticated":   true,
		"Username":        username,
		"ID":              doc.ID,
		"Content":         doc.Content,
	})
	if err != nil {
		util.RenderError(&resources, w, http.StatusInternalServerError, err)
	}
}

func (s *Server) EditDocument(w http.ResponseWriter, r *http.Request) {
	if !s.Config.AccountsEnabled {
		util.WriteError(w, http.StatusNotFound, errors.New("accounts disabled"))
		return
	}

	username, err := s.authenticatedUsername(r)
	if err != nil || username == "" {
		util.WriteError(w, http.StatusUnauthorized, errors.New("authentication required"))
		return
	}

	id := chi.URLParam(r, "document")
	if len(id) != s.Config.IDLength && !slices.Contains(s.Config.Documents, id) {
		util.RenderError(&resources, w, http.StatusBadRequest, fmt.Errorf("invalid document id"))
		return
	}

	doc, err := s.Database.GetDocument(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.RenderError(&resources, w, http.StatusNotFound, err)
			return
		}
		util.RenderError(&resources, w, http.StatusInternalServerError, err)
		return
	}

	if doc.Username != username {
		util.RenderError(&resources, w, http.StatusForbidden, errors.New("you do not own this document"))
		return
	}

	body, err := util.HandleCreateBody(s.Config.MaxSize, r)
	if err != nil {
		util.RenderError(&resources, w, http.StatusBadRequest, err)
		return
	}

	if err := util.ValidateBody(s.Config.MaxSize, body); err != nil {
		util.RenderError(&resources, w, http.StatusBadRequest, err)
		return
	}

	if err := s.Database.UpdateDocument(r.Context(), id, body.Content); err != nil {
		util.RenderError(&resources, w, http.StatusInternalServerError, err)
		return
	}

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		util.WriteJSON(w, http.StatusOK, map[string]string{"id": id})
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s", doc.ID), http.StatusSeeOther)
}

func (s *Server) RemoveDocument(w http.ResponseWriter, r *http.Request) {
	if !s.Config.AccountsEnabled {
		util.WriteError(w, http.StatusNotFound, errors.New("accounts disabled"))
		return
	}

	username, err := s.authenticatedUsername(r)
	if err != nil || username == "" {
		util.WriteError(w, http.StatusUnauthorized, errors.New("authentication required"))
		return
	}

	id := chi.URLParam(r, "document")
	if len(id) != s.Config.IDLength && !slices.Contains(s.Config.Documents, id) {
		util.RenderError(&resources, w, http.StatusBadRequest, fmt.Errorf("invalid document id"))
		return
	}

	doc, err := s.Database.GetDocument(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.RenderError(&resources, w, http.StatusNotFound, err)
			return
		}
		util.RenderError(&resources, w, http.StatusInternalServerError, err)
		return
	}

	if doc.Username != username {
		util.RenderError(&resources, w, http.StatusForbidden, errors.New("you do not own this document"))
		return
	}

	if err := s.Database.DeleteDocument(r.Context(), id); err != nil {
		util.RenderError(&resources, w, http.StatusInternalServerError, err)
		return
	}

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	http.Redirect(w, r, "/account", http.StatusSeeOther)
}

func (s *Server) RemoveAccount(w http.ResponseWriter, r *http.Request) {
	if !s.Config.AccountsEnabled {
		util.WriteError(w, http.StatusNotFound, errors.New("accounts disabled"))
		return
	}

	username, err := s.authenticatedUsername(r)
	if err != nil || username == "" {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	acc, err := s.Database.GetAccountByUsername(r.Context(), username)
	if err != nil {
		util.RenderError(&resources, w, http.StatusInternalServerError, err)
		return
	}

	if err := s.invalidateSession(r); err != nil {
		util.RenderError(&resources, w, http.StatusInternalServerError, err)
		return
	}

	if err := s.Database.DeleteAccount(r.Context(), fmt.Sprintf("%d", acc.ID)); err != nil {
		util.RenderError(&resources, w, http.StatusInternalServerError, err)
		return
	}

	clearSessionCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
