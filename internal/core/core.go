package core

type ConfigRabbitMQ struct {
	User 		string
	Password 	string
	Port		string
	QueueName	string
	TimeDeleyQueue int
}

type Message struct {
	ID		string	`json:id`
	Key		string	`json:"key"`
	Origin	string	`json:"origin"`
	Person	Person	`json:"person,omitempty"`
}

//Person Constructor
func NewMessage(id string, key string, origin string ,person Person) *Message{
	return &Message{
		ID:	id,
		Key: key,
		Origin: origin,
		Person: person,
	}
}

type Person struct {
	ID		string	`json:"id,omitempty"`
	SK		string	`json:"sk,omitempty"`
	Name	string	`json:"name,omitempty"`
	Gender	string	`json:"gender,omitempty"`
}

//Person Constructor
func NewPerson(id string, sk string,name string, gender string) *Person{
	return &Person{
		ID:	id,
		SK: sk,
		Name: name,
		Gender: gender,
	}
}

type DatabaseRDS struct {
    Host 				string `json:"host"`
    Port  				string `json:"port"`
	Schema				string `json:"schema"`
	DatabaseName		string `json:"databaseName"`
	User				string `json:"user"`
	Password			string `json:"password"`
	Db_timeout			int	`json:"db_timeout"`
	Postgres_Driver		string `json:"postgres_driver"`
}