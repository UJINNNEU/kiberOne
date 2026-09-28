package config

type Config struct {
	Api_key      string `json:"API_KEY"`
	Email        string `json:"email"`
	Token        string
	Amount       int      `json:"amount"`
	Note         string   `json:"note"`
	Date_from    string   `json:"date_from"`
	Date_to      string   `json:"date_to"`
	Names_tutors []string `json:"names_tutors"`
	Id_tutors    []int
	Id_groups    []int
}

func (c *Config) validConfig() bool {
	//return c.Api_key != "" && c.Email != "" && c.Id_tutor != 0 && c.Name != ""
	return true
}

func LoadConfig() (*Config, error) {
	return &Config{
		"api key 1",
		"mail.ru",
		"4a0c1ece4a5f48cd1b84d5854c5003e7",
		25,
		"акция: не пропусти месяц занятий, выставленно через API",
		"2026-09-01",
		"2026-09-30",
		[]string{"Олейник Данил"},
		[]int{},
		[]int{},
	}, nil
}
