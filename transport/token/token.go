package token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"kiber/config"
	"net/http"
)

// Вспомогательная структура для тела запроса авторизации
type authRequest struct {
	Email  string `json:"email"`
	APIKey string `json:"api_key"`
}

// Вспомогательная структура для ответа сервера (предполагаем, что токен приходит здесь)
type authResponse struct {
	Token string `json:"token"`
}

// GetToken получает токен авторизации через POST-запрос и записывает его в конфиг.
func GetToken(cfg *config.Config) error {
	// URL эндпоинта авторизации (замени на реальный адрес твоего API)
	url := "https://kiberoneomnov.s20.online/v2api/auth/login"

	// Формируем тело запроса из конфига
	reqBody := authRequest{
		Email:  cfg.Email,
		APIKey: cfg.Api_key,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных для токена: %w", err)
	}

	// Создаем POST-запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса на токен: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса на токен: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("сервер вернул ошибку %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Читаем и парсим ответ
	var authResp authResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("ошибка парсинга ответа с токеном: %w", err)
	}

	// Записываем полученный токен в конфиг
	cfg.Token = authResp.Token

	return nil
}
