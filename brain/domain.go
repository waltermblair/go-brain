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
	Function 	string		`json:"function"`
	NumInputs 	int			`json:"num-inputs"`
	NextKeys 	[]int		`json:"next-keys"`
}

type ConfigRecord struct {
	This		int
	Function 	string
	NumInputs	int
	NextKey 	int
}