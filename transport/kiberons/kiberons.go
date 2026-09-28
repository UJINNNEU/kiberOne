package kiberons

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kiber/config"
	"net/http"
)

type bonusRequest struct {
	Amount int    `json:"amount"`
	Note   string `json:"note"`
}

func addBonusToUser(cfg *config.Config, customerID int) error {
	// Формируем URL с query-параметром customer_id
	url := fmt.Sprintf("https://kiberoneomnov.s20.online/v2api/1/bonus/bonus-add?customer_id=%d", customerID)

	// Готовим тело запроса из конфига
	reqBody := bonusRequest{
		Amount: cfg.Amount,
		Note:   cfg.Note,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("ошибка сериализации бонуса: %w", err)
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса бонуса: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ALFACRM-TOKEN", cfg.Token)

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)

	if resp.StatusCode == http.StatusOK {
		fmt.Println(cfg.Amount, "киберон`ов выставлены ученику с айди", customerID)
	}
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса бонуса для ID %d: %w", customerID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("сервер вернул ошибку %d при начислении бонуса пользователю %d", resp.StatusCode, customerID)
	}

	return nil
}

func AddBonusesToUsers(cfg *config.Config, userIDs []int) error {
	for _, id := range userIDs {
		err := addBonusToUser(cfg, id)
		if err != nil {
			fmt.Println("не удалось начислить бонус пользователю %d: %w", id, err)
			continue
		}
	}
	return nil
}
