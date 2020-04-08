package v1

import (
	"encoding/json"
	"fantasymarket/utils/http/middleware/jwt"
	"fantasymarket/utils/http/responses"
	"net/http"
)

func (api *APIHandler) getUser(w http.ResponseWriter, r *http.Request) {
	// GET   /user/{userID}		(Get info about specific user)

}

func (api *APIHandler) getSelf(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(jwt.UserKey).(jwt.UserClaims)

	// responses.CustomResponse(w, resp, 200)

	responses.CustomResponse(w, user.UserID, 200)
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
	}

	responses.CustomResponse(w, map[string]string{"status": "success"}, 200)
	// TODO: generate token
}

func (api *APIHandler) createUser(w http.ResponseWriter, r *http.Request) {
	user, err := api.DB.CreateGuest()

	if err != nil {
		responses.ErrorResponse(w, http.StatusInternalServerError, "error creating new user account")
		return
	}

	responses.CustomResponse(w, map[string]string{
		"status":   "success",
		"username": user.Username,
		"userID":   user.UserID.String(),
	}, 200)
	// TODO: generate token
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
