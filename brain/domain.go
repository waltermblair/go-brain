package brain

type Message struct {
	Header 		int 		`json:"header"`
	Body 		MessageBody `json:"body"`
}

type MessageBody struct {
	Configs 	[]Config 	`json:"configuration"`
	Input 		[]bool 		`json:"input"`
}

type Config struct {
	ID 			int 		`json:"id"`
	Status 		string 		`json:"status"`
	Function 	string
	NextKeys 	[]int
}

type ConfigRecord struct {
	ID			int			`json:"this"`
	Function 	string      `json:"function"`
	NextKey 	int 		`json:"next"`
}