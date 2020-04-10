package v1

import (
	"encoding/json"
	"fantasymarket/utils/http/middleware/jwt"
	"fantasymarket/utils/http/responses"
	"net/http"

	"github.com/go-chi/chi"
)

func (api *APIHandler) getUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	resp, err := api.DB.GetSelf(username)
	if err != nil {
		responses.ErrorResponse(w, http.StatusNotFound, err.Error())
	}

	responses.CustomResponse(w, resp, 200)
}

func (api *APIHandler) getSelf(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(jwt.UserKey).(jwt.UserClaims)

	resp, err := api.DB.GetSelf(user.Username)
	if err != nil {
		responses.ErrorResponse(w, http.StatusNotFound, err.Error())
	}

	responses.CustomResponse(w, resp, 200)
}

type updateUserRequest struct {
	username, password, newPassword string
}

func (api *APIHandler) updateSelf(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(jwt.UserKey).(jwt.UserClaims)

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "error parsing request")
		return
	}

	if req.newPassword != "" && req.password != "" {
		if err := api.DB.ChangePassword(req.username, req.password, req.newPassword); err != nil {
			responses.ErrorResponse(w, http.StatusInternalServerError, "error updating password")
			return
		}
	}

	if req.username != "" {
		if err := api.DB.RenameUser(user.Username, req.username); err != nil {
			responses.ErrorResponse(w, http.StatusInternalServerError, "error updating username")
			return
		}
		user.Username = req.username
	}

	token, err := jwt.CreateToken("secret", user.Username, user.UserID)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "error generating user token")
		return
	}

	responses.CustomResponse(w, map[string]string{
		"username": user.Username,
		"userID":   user.UserID,
		"token":    token,
	}, 200)
}

func (api *APIHandler) createUser(w http.ResponseWriter, r *http.Request) {
	user, err := api.DB.CreateGuest()

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "error creating new user account")
		return
	}

	token, err := jwt.CreateToken("secret", user.Username, user.UserID.String())
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "error generating user token")
		return
	}

	responses.CustomResponse(w, map[string]string{
		"username": user.Username,
		"userID":   user.UserID.String(),
		"token":    token,
	}, 200)

}

type loginUserRequest struct {
	username, password string
}

func (api *APIHandler) loginUser(w http.ResponseWriter, r *http.Request) {

	var req loginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "error parsing request")
		return
	}

	user, err := api.DB.LoginUser(req.username, req.password)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, err.Error())
	}

	responses.CustomResponse(w, user, 200)
}

func (api *APIHandler) getPortfolio(w http.ResponseWriter, r *http.Request) {
	// TODO
}
