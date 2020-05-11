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
		responses.ErrorResponse(w, http.StatusNotFound, userNotFoundError.Error())
	}

	responses.CustomResponse(w, resp, 200)
}

func (api *APIHandler) getSelf(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := jwt.GetUserFromContext(ctx)

	resp, err := api.DB.GetSelf(user.Username)
	if err != nil {
		responses.ErrorResponse(w, http.StatusNotFound, userNotFoundError.Error())
	}

	responses.CustomResponse(w, resp, 200)
}

type updateUserRequest struct {
	username, password, newPassword string
}

func (api *APIHandler) updateSelf(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := jwt.GetUserFromContext(ctx)

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, decodingError.Error())
		return
	}

	if req.newPassword != "" && req.password != "" {
		if err := api.DB.ChangePassword(req.username, req.password, req.newPassword); err != nil {
			responses.ErrorResponse(w, http.StatusInternalServerError, passwordError.Error())
			return
		}
	}

	if req.username != "" {
		if err := api.DB.RenameUser(user.UserID, user.Username, req.username); err != nil {
			responses.ErrorResponse(w, http.StatusInternalServerError, usernameError.Error())
			return
		}
		user.Username = req.username
	}

	token, err := jwt.CreateToken(api.Config.TokenSecret, user.Username, user.UserID)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, tokenError.Error())
		return
	}

	responses.CustomResponse(w, map[string]interface{}{
		"user": map[string]string{
			"userID":   user.UserID.String(),
			"username": user.Username,
		},
		"token": token,
	}, 200)
}

func (api *APIHandler) createUser(w http.ResponseWriter, r *http.Request) {
	user, err := api.DB.CreateGuest()

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, accountError.Error())
		return
	}

	token, err := jwt.CreateToken(api.Config.TokenSecret, user.Username, user.UserID)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, tokenError.Error())
		return
	}

	responses.CustomResponse(w, map[string]interface{}{
		"user": map[string]string{
			"userID":   user.UserID.String(),
			"username": user.Username,
		},
		"token": token,
	}, 200)

}

type loginUserRequest struct {
	username, password string
}

func (api *APIHandler) loginUser(w http.ResponseWriter, r *http.Request) {

	var req loginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, decodingError.Error())
		return
	}

	user, err := api.DB.LoginUser(req.username, req.password)
	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, loginError.Error())
	}

	responses.CustomResponse(w, user, 200)
}

func (api *APIHandler) getPortfolio(w http.ResponseWriter, r *http.Request) {
	// TODO
}
