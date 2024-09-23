package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	db "github.com/duyanhitbe/go-boilerplate/internal/database/generated"
	mockdb "github.com/duyanhitbe/go-boilerplate/internal/database/mock"
	"github.com/duyanhitbe/go-boilerplate/internal/hash"
	"github.com/duyanhitbe/go-boilerplate/internal/token"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testCase struct {
		name      string
		buildBody func() *strings.Reader
		buildStub func(store *mockdb.MockStore)
		check     func(w *httptest.ResponseRecorder)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	h := hash.NewBcrypt()
	tk := token.NewJWT("secret_test")
	s := &Server{
		store: store,
		h:     h,
		t:     tk,
	}

	id := uuid.New()
	testCases := []testCase{
		{
			name: "Success",
			buildBody: func() *strings.Reader {
				req := registerRequest{
					Username: "test_username",
					Password: "test_password",
				}
				reqBody, err := json.Marshal(req)
				require.NoError(t, err)
				return strings.NewReader(string(reqBody))
			},
			buildStub: func(store *mockdb.MockStore) {
				user := db.User{
					ID:        id,
					Username:  "test_username",
					Password:  "hashed_password",
					CreatedAt: sql.NullTime{},
					UpdatedAt: sql.NullTime{},
				}
				store.EXPECT().
					FindOneUserByUsername(gomock.Any(), "test_username").
					Times(1).
					Return(nil, sql.ErrNoRows)

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&user, nil)
			},
			check: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, w.Code)

				data, err := io.ReadAll(w.Body)
				require.NoError(t, err)

				var rsp response
				err = json.Unmarshal(data, &rsp)
				require.NoError(t, err)

				rspData, ok := rsp.Data.(map[string]interface{})
				require.True(t, ok)

				require.Equal(t, http.StatusCreated, rsp.Status)
				require.Equal(t, http.StatusText(http.StatusCreated), rsp.Message)
				require.Equal(t, id.String(), rspData["id"])
				require.Equal(t, "test_username", rspData["username"])
			},
		},
		{
			name: "Conflict",
			buildBody: func() *strings.Reader {
				req := registerRequest{
					Username: "test_username",
					Password: "test_password",
				}
				reqBody, err := json.Marshal(req)
				require.NoError(t, err)
				return strings.NewReader(string(reqBody))
			},
			buildStub: func(store *mockdb.MockStore) {
				user := db.User{
					ID:        uuid.New(),
					Username:  "test_username",
					Password:  "hashed_password",
					CreatedAt: sql.NullTime{},
					UpdatedAt: sql.NullTime{},
				}
				store.EXPECT().
					FindOneUserByUsername(gomock.Any(), "test_username").
					Times(1).
					Return(&user, nil)
			},
			check: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, w.Code)

				data, err := io.ReadAll(w.Body)
				require.NoError(t, err)

				var rsp errorResponse
				err = json.Unmarshal(data, &rsp)
				require.NoError(t, err)

				require.Equal(t, http.StatusConflict, rsp.Status)
				require.Equal(t, http.StatusText(http.StatusConflict), rsp.Message)
				require.Equal(t, "username already taken", rsp.Error)
			},
		},
		{
			name: "Bad Request",
			buildBody: func() *strings.Reader {
				req := registerRequest{
					Username: "test_username",
				}
				reqBody, err := json.Marshal(req)
				require.NoError(t, err)
				return strings.NewReader(string(reqBody))
			},
			buildStub: func(store *mockdb.MockStore) {},
			check: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)

				data, err := io.ReadAll(w.Body)
				require.NoError(t, err)

				var rsp errorResponse
				err = json.Unmarshal(data, &rsp)
				require.NoError(t, err)

				require.Equal(t, http.StatusBadRequest, rsp.Status)
				require.Equal(t, http.StatusText(http.StatusBadRequest), rsp.Message)
				require.Equal(t, []interface{}([]interface{}{map[string]interface{}{"field": "password", "message": "This field is required"}}), rsp.Errors)
			},
		},
		{
			name: "FindOneUserByUsername Fail",
			buildBody: func() *strings.Reader {
				req := registerRequest{
					Username: "test_username",
					Password: "test_password",
				}
				reqBody, err := json.Marshal(req)
				require.NoError(t, err)
				return strings.NewReader(string(reqBody))
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					FindOneUserByUsername(gomock.Any(), "test_username").
					Times(1).
					Return(nil, errors.New("something went wrong"))
			},
			check: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)

				data, err := io.ReadAll(w.Body)
				require.NoError(t, err)

				var rsp errorResponse
				err = json.Unmarshal(data, &rsp)
				require.NoError(t, err)

				require.Equal(t, http.StatusInternalServerError, rsp.Status)
				require.Equal(t, http.StatusText(http.StatusInternalServerError), rsp.Message)
				require.Equal(t, "something went wrong", rsp.Error)
			},
		},
		{
			name: "CreateUser Fail",
			buildBody: func() *strings.Reader {
				req := registerRequest{
					Username: "test_username",
					Password: "test_password",
				}
				reqBody, err := json.Marshal(req)
				require.NoError(t, err)
				return strings.NewReader(string(reqBody))
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					FindOneUserByUsername(gomock.Any(), "test_username").
					Times(1).
					Return(nil, sql.ErrNoRows)

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("something went wrong"))
			},
			check: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)

				data, err := io.ReadAll(w.Body)
				require.NoError(t, err)

				var rsp errorResponse
				err = json.Unmarshal(data, &rsp)
				require.NoError(t, err)

				require.Equal(t, http.StatusInternalServerError, rsp.Status)
				require.Equal(t, http.StatusText(http.StatusInternalServerError), rsp.Message)
				require.Equal(t, "something went wrong", rsp.Error)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStub(store)

			router := s.router()
			w := httptest.NewRecorder()

			body := tc.buildBody()
			request, err := http.NewRequest(http.MethodPost, "/api/v1/auth/register", body)
			require.NoError(t, err)
			router.ServeHTTP(w, request)

			tc.check(w)
		})
	}

}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type testCase struct {
		name      string
		buildBody func() *strings.Reader
		buildStub func(store *mockdb.MockStore)
		check     func(w *httptest.ResponseRecorder)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	h := hash.NewBcrypt()
	tk := token.NewJWT("secret_test")
	s := &Server{
		store: store,
		h:     h,
		t:     tk,
	}

	pwd := "test_password"

	testCases := []testCase{
		{
			name: "Success",
			buildBody: func() *strings.Reader {
				req := registerRequest{
					Username: "test_username",
					Password: pwd,
				}
				reqBody, err := json.Marshal(req)
				require.NoError(t, err)
				return strings.NewReader(string(reqBody))
			},
			buildStub: func(store *mockdb.MockStore) {
				hashPwd, err := s.h.Create(pwd)
				require.NoError(t, err)
				user := db.User{
					ID:        uuid.New(),
					Username:  "test_username",
					Password:  hashPwd,
					CreatedAt: sql.NullTime{},
					UpdatedAt: sql.NullTime{},
				}

				store.EXPECT().
					FindOneUserByUsername(gomock.Any(), "test_username").
					Times(1).
					Return(&user, nil)
			},
			check: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)

				data, err := io.ReadAll(w.Body)
				require.NoError(t, err)

				var rsp response
				err = json.Unmarshal(data, &rsp)
				require.NoError(t, err)

				rspData, ok := rsp.Data.(map[string]interface{})
				require.True(t, ok)

				require.Equal(t, http.StatusOK, rsp.Status)
				require.Equal(t, http.StatusText(http.StatusOK), rsp.Message)
				require.Equal(t, "bearer", rspData["type"])
				require.Equal(t, "ms", rspData["unit"])
				require.Equal(t, float64(24*time.Hour.Milliseconds()), rspData["expires_in"])
				require.NotEmpty(t, rspData["access_token"])
				require.NotEmpty(t, rspData["refresh_token"])
			},
		},
		{
			name: "Bad Request",
			buildBody: func() *strings.Reader {
				req := registerRequest{
					Username: "test_username",
				}
				reqBody, err := json.Marshal(req)
				require.NoError(t, err)
				return strings.NewReader(string(reqBody))
			},
			buildStub: func(store *mockdb.MockStore) {

			},
			check: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)

				data, err := io.ReadAll(w.Body)
				require.NoError(t, err)

				var rsp errorResponse
				err = json.Unmarshal(data, &rsp)
				require.NoError(t, err)

				require.Equal(t, http.StatusBadRequest, rsp.Status)
				require.Equal(t, http.StatusText(http.StatusBadRequest), rsp.Message)
				require.Equal(t, []interface{}([]interface{}{map[string]interface{}{"field": "password", "message": "This field is required"}}), rsp.Errors)
			},
		},
		{
			name: "User not found",
			buildBody: func() *strings.Reader {
				req := registerRequest{
					Username: "test_username",
					Password: pwd,
				}
				reqBody, err := json.Marshal(req)
				require.NoError(t, err)
				return strings.NewReader(string(reqBody))
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					FindOneUserByUsername(gomock.Any(), "test_username").
					Times(1).
					Return(nil, sql.ErrNoRows)
			},
			check: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, w.Code)

				data, err := io.ReadAll(w.Body)
				require.NoError(t, err)

				var rsp errorResponse
				err = json.Unmarshal(data, &rsp)
				require.NoError(t, err)

				require.Equal(t, http.StatusForbidden, rsp.Status)
				require.Equal(t, http.StatusText(http.StatusForbidden), rsp.Message)
				require.Equal(t, "user not found", rsp.Error)
			},
		},
		{
			name: "FindOneUserByUsername Fail",
			buildBody: func() *strings.Reader {
				req := registerRequest{
					Username: "test_username",
					Password: pwd,
				}
				reqBody, err := json.Marshal(req)
				require.NoError(t, err)
				return strings.NewReader(string(reqBody))
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					FindOneUserByUsername(gomock.Any(), "test_username").
					Times(1).
					Return(nil, errors.New("something went wrong"))
			},
			check: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)

				data, err := io.ReadAll(w.Body)
				require.NoError(t, err)

				var rsp errorResponse
				err = json.Unmarshal(data, &rsp)
				require.NoError(t, err)

				require.Equal(t, http.StatusInternalServerError, rsp.Status)
				require.Equal(t, http.StatusText(http.StatusInternalServerError), rsp.Message)
				require.Equal(t, "something went wrong", rsp.Error)
			},
		},
		{
			name: "Wrong password",
			buildBody: func() *strings.Reader {
				req := registerRequest{
					Username: "test_username",
					Password: pwd,
				}
				reqBody, err := json.Marshal(req)
				require.NoError(t, err)
				return strings.NewReader(string(reqBody))
			},
			buildStub: func(store *mockdb.MockStore) {
				hashPwd, err := s.h.Create("wrong_password")
				require.NoError(t, err)
				user := db.User{
					ID:        uuid.New(),
					Username:  "test_username",
					Password:  hashPwd,
					CreatedAt: sql.NullTime{},
					UpdatedAt: sql.NullTime{},
				}

				store.EXPECT().
					FindOneUserByUsername(gomock.Any(), "test_username").
					Times(1).
					Return(&user, nil)
			},
			check: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)

				data, err := io.ReadAll(w.Body)
				require.NoError(t, err)

				var rsp errorResponse
				err = json.Unmarshal(data, &rsp)
				require.NoError(t, err)

				require.Equal(t, http.StatusBadRequest, rsp.Status)
				require.Equal(t, http.StatusText(http.StatusBadRequest), rsp.Message)
				require.Equal(t, "invalid password", rsp.Error)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStub(store)

			router := s.router()
			w := httptest.NewRecorder()

			body := tc.buildBody()
			request, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
			require.NoError(t, err)
			router.ServeHTTP(w, request)

			tc.check(w)
		})
	}

}

func TestMe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	h := hash.NewBcrypt()
	tk := token.NewJWT("secret_test")
	s := &Server{
		store: store,
		h:     h,
		t:     tk,
	}

	sub := uuid.New()
	accessToken, err := tk.Create(sub.String(), time.Minute)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	require.NoError(t, err)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken.Token))

	store.EXPECT().
		FindOneUserById(gomock.Any(), sub).
		Return(&db.User{
			ID:        sub,
			Username:  "username_test",
			CreatedAt: sql.NullTime{},
			UpdatedAt: sql.NullTime{},
		}, nil)

	router := s.router()
	router.ServeHTTP(w, request)

	data, err := io.ReadAll(w.Body)
	require.NoError(t, err)

	var rsp response
	err = json.Unmarshal(data, &rsp)
	require.NoError(t, err)

	dataJson, ok := rsp.Data.(map[string]interface{})
	require.True(t, ok)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, http.StatusOK, rsp.Status)
	require.Equal(t, http.StatusText(http.StatusOK), rsp.Message)
	require.Equal(t, sub.String(), dataJson["id"])
	require.Equal(t, "username_test", dataJson["username"])
}
