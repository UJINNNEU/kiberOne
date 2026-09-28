package groups

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kiber/config"
	"net/http"
)

type TeacherRequest struct {
	ID      int `json:"teacher_id"`
	Removed int `json:"removed"`
}

type APIResponse struct {
	Total int     `json:"total"`
	Count int     `json:"count"`
	Page  int     `json:"page"`
	Items []Group `json:"items"` // Массив групп
}

type Group struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// Функция для получения групп учителя по его ID
func getTeacherGroups(id_tutor int, token string) ([]int, error) {

	url := "https://kiberoneomnov.s20.online/v2api/1/group/index"

	// Создаем тело запроса
	requestBody := TeacherRequest{
		ID:      id_tutor,
		Removed: 0,
	}

	// Сериализуем структуру в JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации JSON: %w", err)
	}

	// Создаем POST-запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	// Токен в заголовке
	req.Header.Set("X-ALFACRM-TOKEN", token)

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка сервера: статус %d", resp.StatusCode)
	}

	// Читаем тело ответа
	var bodyBytes []byte
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}
	bodyBytes = buf.Bytes()

	var apiResp APIResponse

	// 1. Парсим весь JSON во временную структуру
	err = json.Unmarshal(bodyBytes, &apiResp)
	if err != nil {
		return nil, err
	}

	var ids_grops []int

	for _, v := range apiResp.Items {
		ids_grops = append(ids_grops, v.Id)
	}
	return ids_grops, nil
}

func GetAllGroups(config *config.Config) error {

	for _, v := range config.Id_tutors {

		data, err := getTeacherGroups(v, config.Token)

		if err == nil || data != nil {
			fmt.Println("найдено аж ", len(data), " групп на всех туторов")
			config.Id_groups = append(config.Id_groups, data...)
		} else {
			fmt.Println("error")
		}
	}
	return nil
}
