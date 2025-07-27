package middleware

import (
	"net/http"
)

func AuthMiddleware(next http.Handler, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("devshare_auth")
		if err == nil && cookie.Value == password {
			next.ServeHTTP(w, r)
			return
		}

		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			if r.Form.Get("password") == password {
				http.SetCookie(w, &http.Cookie{
					Name:     "devshare_auth",
					Value:    password,
					Path:     "/",
					HttpOnly: true,
					MaxAge:   3600, // 1 hour
				})

				redirectPath := r.URL.Path
				if redirectPath == "" {
					redirectPath = "/"
				}
				if r.URL.RawQuery != "" {
					redirectPath += "?" + r.URL.RawQuery
				}
				http.Redirect(w, r, redirectPath, http.StatusFound)
				return
			}
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		formAction := r.URL.Path
		if formAction == "" {
			formAction = "/"
		}
		if r.URL.RawQuery != "" {
			formAction += "?" + r.URL.RawQuery
		}

		if formAction == "" || formAction == "/" {
			formAction = r.RequestURI
			if formAction == "" {
				formAction = "/"
			}
		}

		html := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Authentication Required</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			background-color: #f4f4f4;
			color: #333;
			display: flex;
			justify-content: center;
			align-items: center;
			height: 100vh;
			margin: 0;
			padding: 20px;
		}
		.auth-container {
			background: white;
			padding: 1.5rem;
			border-radius: 8px;
			box-shadow: 0 2px 10px rgba(0,0,0,0.1);
			text-align: center;
			max-width: 430px;
			width: 100%;
		}
		.auth-container h2 {
			margin-bottom: 1rem;
			color: #2c3e50;
		}
		.auth-container p {
			margin-bottom: 1.5rem;
			color: #7f8c8d;
		}
		.auth-form {
			display: flex;
			flex-direction: column;
			gap: 1rem;
		}
		.auth-form input[type="password"] {
			padding: 0.75rem;
			border: 1px solid #ddd;
			border-radius: 4px;
			font-size: 1rem;
		}
		.auth-form button {
			padding: 0.75rem;
			background-color: #3498db;
			color: white;
			border: none;
			border-radius: 4px;
			font-size: 1rem;
			cursor: pointer;
			transition: background-color 0.2s;
		}
		.auth-form button:hover {
			background-color: #2980b9;
		}
		.error-message {
			color: #e74c3c;
			margin-top: 1rem;
			font-size: 0.9rem;
		}
	</style>
</head>
<body>
	<div class="auth-container">
		<h1 style="color:#0d92ee">DevShare</h2>
		<a href="https://anophel.com"><img src="https://anophel.com/Anophel-logo.svg" width="200" height="60" alt="Anophel logo"/></a>
		<h2>🔐 Authentication Required</h1>
		<p>Please enter the password to access this development environment.</p>
		<form method="POST" action="` + formAction + `" class="auth-form">
			<input type="password" name="password" placeholder="Enter password" required autofocus>
			<button type="submit">Access Development Environment</button>
		</form>
		<p style="margin-top: 20px; font-size: 0.9rem; color: #95a5a6;">
			If you don't have a password, please contact the developer.
		</p>
		<p style="margin-top: 10px; font-size: 0.8rem; color: #bdc3c7;">
			Built with 💙 by <a href="https://anophel.com">Anophel</a>
		</p>
	</div>
</body>
</html>`

		w.Write([]byte(html))
	})
}
