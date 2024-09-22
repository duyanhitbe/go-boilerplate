package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

	user := db.User{
		ID:        uuid.New(),
		Username:  "test_username",
		Password:  "hashed_password",
		CreatedAt: sql.NullTime{},
		UpdatedAt: sql.NullTime{},
	}

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

				require.Equal(t, http.StatusCreated, rsp.Status)
				require.Equal(t, http.StatusText(http.StatusCreated), rsp.Message)
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
				store.EXPECT().
					FindOneUserByUsername(gomock.Any(), "test_username").
					Times(1).
					Return(&user, nil)
			},
			check: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, w.Code)

				data, err := io.ReadAll(w.Body)
				require.NoError(t, err)

				var rsp response
				err = json.Unmarshal(data, &rsp)
				require.NoError(t, err)

				require.Equal(t, http.StatusConflict, rsp.Status)
				require.Equal(t, http.StatusText(http.StatusConflict), rsp.Message)
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

				var rsp response
				err = json.Unmarshal(data, &rsp)
				require.NoError(t, err)

				require.Equal(t, http.StatusBadRequest, rsp.Status)
				require.Equal(t, http.StatusText(http.StatusBadRequest), rsp.Message)
			},
		},
		{
			name: "FindByUsername fail",
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

				var rsp response
				err = json.Unmarshal(data, &rsp)
				require.NoError(t, err)

				require.Equal(t, http.StatusInternalServerError, rsp.Status)
				require.Equal(t, http.StatusText(http.StatusInternalServerError), rsp.Message)
			},
		},
		{
			name: "Create fail",
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

				var rsp response
				err = json.Unmarshal(data, &rsp)
				require.NoError(t, err)

				require.Equal(t, http.StatusInternalServerError, rsp.Status)
				require.Equal(t, http.StatusText(http.StatusInternalServerError), rsp.Message)
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
