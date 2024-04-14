package models

// CarInfo is a car description
type CarInfo struct {
	RegNum string `json:"reg_num,omitempty" db:"reg_num"`
	Mark   string `json:"mark,omitempty" db:"mark"`
	Model  string `json:"model,omitempty" db:"model"`
	Year   int    `json:"year,omitempty" db:"year"`
	Owner  Person `json:"owner,omitempty"`
}

// Person is a description of the car owner
type Person struct {
	Name       string `json:"name,omitempty" db:"owner_name"`
	Surname    string `json:"surname,omitempty" db:"owner_surname"`
	Patronymic string `json:"patronymic,omitempty" db:"owner_patronymic"`
}
