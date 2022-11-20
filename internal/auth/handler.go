package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type (
	clientService interface {
		ValidateUser(user, password string) bool
	}

	handler struct {
		logger        *zap.Logger
		clientService clientService
		key           []byte
	}

	Credential struct {
		User     string
		Password string
	}

	Claims struct {
		User string
		jwt.RegisteredClaims
	}
)

const cookieName = "token"

func NewHandler(
	logger *zap.Logger,
	clientService clientService,
	config Config,
) *handler {
	return &handler{
		logger:        logger,
		clientService: clientService,
		key:           config.Key,
	}
}

func (h *handler) RegisterRoutes(r *mux.Router) {
	r.Methods(http.MethodGet).Path("/").HandlerFunc(h.auth)
}

func (h *handler) auth(writer http.ResponseWriter, request *http.Request) {
	var creds Credential
	err := json.NewDecoder(request.Body).Decode(&creds)
	if err != nil {
		h.logger.Warn("decode", zap.Error(err))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if !h.clientService.ValidateUser(creds.User, creds.Password) {
		h.logger.Warn("validate",
			zap.String("user", creds.User),
			zap.String("pass", creds.Password),
			zap.Error(err))
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		User: creds.User,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.key)
	if err != nil {
		h.logger.Warn("sign", zap.Error(err))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(writer, &http.Cookie{
		Name:    cookieName,
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func (h *handler) AuthMiddleware(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		if !h.check(w, r) {
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}

func (h *handler) check(w http.ResponseWriter, r *http.Request) bool {
	c, err := r.Cookie(cookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	tknStr := c.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return h.key, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	return true
}
