package teacher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kiber/config"
	"net/http"
)

type Teacher_request struct {
	Name string `json:"name"`
}
type Teacher_response struct {
	Items []Items_response `json:"items"`
}

type Items_response struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func getIdByName(name string, token string) (int, error) {

	url := "https://kiberoneomnov.s20.online/v2api/1/teacher/index"

	requestBody := Teacher_request{
		Name: name,
	}

	// Сериализуем структуру в JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return 0, fmt.Errorf("ошибка сериализации JSON: %w", err)
	}

	// Создаем POST-запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	// Токен в заголовке
	req.Header.Set("X-ALFACRM-TOKEN", token)

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("ошибка сервера: статус %d", resp.StatusCode)
	}

	// Читаем тело ответа
	var bodyBytes []byte
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("ошибка чтения ответа: %w", err)
	}
	bodyBytes = buf.Bytes()

	var response Teacher_response

	// Парсим JSON
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return 0, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	// Проверка на пустой список
	if len(response.Items) == 0 {
		return 0, fmt.Errorf("список учителей пуст")
	}

	// Ищем учителя по имени (так как в items может быть несколько человек,
	// лучше искать циклом, а не брать [0] наугад)
	for _, item := range response.Items {
		if item.Name == name {
			return item.Id, nil
		}
	}

	return 0, fmt.Errorf("учитель с именем %s не найден", name)
}

func GetAllIdsTeachers(
	config *config.Config) error {

	var idsTutor []int

	for _, v := range config.Names_tutors {

		id, err := getIdByName(v, config.Token)
		if err == nil {
			idsTutor = append(idsTutor, id)
		} else {
			fmt.Println(err.Error())
		}

	}

	config.Id_tutors = idsTutor
	return nil
}
