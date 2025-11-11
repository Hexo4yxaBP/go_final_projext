package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

// UserData represents credentials submitted to /api/signin.
type UserData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// signinHandler validates provided credentials and returns a signed JWT token
// on success. The server uses TODO_PASSWORD as the shared secret; if that
// environment variable is unset no authentication is required.
func signinHandler(w http.ResponseWriter, r *http.Request) {

	var user UserData

	//read request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	pwd := os.Getenv("TODO_PASSWORD")

	if user.Password != pwd {
		writeJson(w, map[string]string{"error": "incorrect login or password"}, http.StatusUnauthorized)
		return
	}

	secret := []byte(user.Password)

	jwtToken := jwt.New(jwt.SigningMethodHS256)

	signedToken, err := jwtToken.SignedString(secret)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusUnauthorized)
		return
	}

	writeJson(w, map[string]string{"token": signedToken}, http.StatusOK)

}

// auth is a middleware that enforces authentication when TODO_PASSWORD is set.
// It expects a JWT token in the "token" cookie and validates it using the
// TODO_PASSWORD as the signing secret.
func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtS string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtS = cookie.Value
			}

			// здесь код для валидации и проверки JWT-токена

			jwtToken, err := jwt.Parse(jwtS, func(t *jwt.Token) (interface{}, error) {
				// секретный ключ для всех токенов одинаковый, поэтому просто возвращаем его
				return []byte(pass), nil
			})

			if err != nil || !jwtToken.Valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
