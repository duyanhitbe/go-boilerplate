package server

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	db "github.com/duyanhitbe/go-boilerplate/internal/database/generated"
	mockdb "github.com/duyanhitbe/go-boilerplate/internal/database/mock"
	"github.com/duyanhitbe/go-boilerplate/internal/hash"
	"github.com/duyanhitbe/go-boilerplate/internal/token"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	type testCase struct {
		name string
		req  registerRequest
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	h := hash.NewBcrypt()
	tk := token.NewJWT("secret_test")

	t.Run("Success", func(t *testing.T) {
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
			Return(nil, sql.ErrNoRows)

		store.EXPECT().
			CreateUser(gomock.Any(), gomock.Any()).
			Times(1).
			Return(&user, nil)

		s := &Server{
			store: store,
			h:     h,
			t:     tk,
		}

		router := s.router()
		w := httptest.NewRecorder()

		req := registerRequest{
			Username: "test_username",
			Password: "test_password",
		}

		reqBody, err := json.Marshal(req)
		require.NoError(t, err)
		request, err := http.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(string(reqBody)))
		require.NoError(t, err)
		router.ServeHTTP(w, request)

		require.Equal(t, http.StatusCreated, w.Code)

		data, err := io.ReadAll(w.Body)
		require.NoError(t, err)

		var rsp response
		err = json.Unmarshal(data, &rsp)
		require.NoError(t, err)

		require.Equal(t, http.StatusCreated, rsp.Status)
		require.Equal(t, http.StatusText(http.StatusCreated), rsp.Message)
	})
}
