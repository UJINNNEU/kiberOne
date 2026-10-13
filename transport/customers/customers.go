package customers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kiber/config"
	"net/http"
)

// attendanceRequest описывает тело POST-запроса для получения данных о посещениях
type attendanceRequest struct {
	Status   int    `json:"status"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
	GroupID  int    `json:"group_id"`
	Page     int    `json:"page"`
}

// lessonDetail описывает посещение конкретного ученика на уроке
type lessonDetail struct {
	CustomerID int `json:"customer_id"`
	IsAttend   int `json:"is_attend"` // 1 - присутствовал, 0 - отсутствовал
}

// lesson представляет одно занятие из массива items
type lesson struct {
	ID      int            `json:"id"`
	Details []lessonDetail `json:"details"`
}

// attendanceResponse описывает корневую структуру JSON-ответа по посещениям
type attendanceResponse struct {
	Total int      `json:"total"`
	Items []lesson `json:"items"`
}

// GetAttendanceData выполняет POST-запрос к API для получения данных о посещениях группы
// за указанный период. Возвращает сырые байты ответа или ошибку.
func getAttendanceData(token, dateFrom, dateTo string, groupID int) ([]byte, error) {
	url := "https://kiberoneomnov.s20.online/v2api/1/lesson/index" // Замените на реальный URL вашего API

	// Формируем тело запроса
	reqBody := attendanceRequest{
		Status:   3,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		GroupID:  groupID,
		Page:     0,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации запроса: %w", err)
	}

	// Создаем POST-запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания HTTP запроса: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
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
		return nil, fmt.Errorf("сервер вернул статус %d", resp.StatusCode)
	}

	// Читаем тело ответа
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	return bodyBytes, nil
}

// GetPerfectAttendanceIDs принимает ID группы и сырой JSON с данными о посещениях,
// возвращает список ID учеников, которые посетили ВСЕ занятия в этой выборке.
func getPerfectAttendanceIDs(groupID int, jsonData []byte) ([]int, error) {

	var resp attendanceResponse

	// 1. Парсим JSON
	if err := json.Unmarshal(jsonData, &resp); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON для группы %d: %w", groupID, err)
	}

	// Определяем общее количество уроков
	totalLessons := resp.Total
	if totalLessons == 0 && len(resp.Items) > 0 {
		totalLessons = len(resp.Items)
	}

	// Если уроков нет, возвращаем пустой срез
	if totalLessons == 0 {
		return []int{}, nil
	}

	// 2. Подсчитываем количество посещений для каждого ученика
	attendanceMap := make(map[int]int)

	for _, l := range resp.Items {
		for _, d := range l.Details {
			if d.IsAttend == 1 {
				attendanceMap[d.CustomerID]++
			}
		}
	}

	// 3. Отбираем тех, у кого число посещений равно общему числу уроков
	var perfectIDs []int
	for customerID, count := range attendanceMap {
		if count == totalLessons {
			perfectIDs = append(perfectIDs, customerID)
		}
	}

	return perfectIDs, nil
}

func GetCustomers(config *config.Config) ([]int, error) {

	var id_perf []int

	for _, id := range config.Id_groups {

		data, err := getAttendanceData(config.Token,
			config.Date_from,
			config.Date_to,
			id)

		if err != nil || data == nil {
			fmt.Print("error group с айди ", id)
			continue
		} else {
			ids_customers, err := getPerfectAttendanceIDs(id, data)

			if err != nil {
				continue
			}

			id_perf = append(id_perf, ids_customers...)
		}

	}

	if len(id_perf) != 0 {
		return id_perf, nil
	} else {
		return nil, errors.New("err last methods")
	}

}

