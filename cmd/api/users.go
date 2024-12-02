package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/szymczykkrzysztof/social/internal/store"
	"net/http"
	"strconv"
)

type userKey string

const userCtxKey userKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	followerUser := getUserFromContext(ctx)

	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
	}

	app.store.Users.Follow(ctx, followerUser.ID, payload.UserID)
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	followerUser := getUserFromContext(ctx)
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
	}
	app.store.Users.Unfollow(ctx, followerUser.ID, payload.UserID)
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
		ctx := r.Context()
		user, err := app.store.Users.GetById(ctx, userId)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notFoundResponse(w, r, err)
				return
			default:
				app.internalServerError(w, r, err)
				return
			}
		}
		if err = app.jsonResponse(w, http.StatusOK, user); err != nil {
			app.internalServerError(w, r, err)
		}
		ctx = context.WithValue(ctx, userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
func getUserFromContext(ctx context.Context) *store.User {
	user, ok := ctx.Value(userCtxKey).(*store.User)
	if !ok {
		return nil
	}
	return user
}
