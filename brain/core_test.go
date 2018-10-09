package brain_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/waltermblair/brain/brain"
	"io/ioutil"
	"os"
	"strconv"
)

type MockRabbitClientImpl struct {
	config 		Config
}

func NewMockRabbitClient(cfg Config) RabbitClient {
	r := MockRabbitClientImpl{
		cfg,
	}
	return &r
}

func (r *MockRabbitClientImpl) RunConsumer() (bool, error) {
	return true, nil
}
func (r *MockRabbitClientImpl) Publish (m MessageBody, s string) error {
	output := strconv.FormatBool(m.Input[0])
	return errors.New("next-key: " + s + " output: " + output)
}
func (r *MockRabbitClientImpl) InitRabbit() {}

var _ = Describe("Core", func() {

	var s 			Service
	var msgRun		Message
	var index 		int
	var cfg    	    Config

	BeforeSuite(func() {
		file, _ := os.Open("../resources/json/run.json")
		bytes, _ := ioutil.ReadAll(file)
		json.Unmarshal(bytes, &msgRun)
		index = 0
	})

	Describe("Core", func() {

		var mockRabbit RabbitClient
		var mockDB     DBClient
		var db     	   *sql.DB
		var mock       sqlmock.Sqlmock
		var fn 		   string
		var nextKeys   []int

		BeforeEach(func() {
			mockRabbit = NewMockRabbitClient(cfg)
			db, mock, _ = sqlmock.New()
			mockDB = &DBClientImpl{
				"test",
				db,
			}
			fn = "not"
			nextKeys = []int{1, 2, 3}
			rows0 := sqlmock.NewRows([]string{"function", "next"}).
				AddRow(fn, nextKeys[0]).
				AddRow(fn, nextKeys[1])
			rows1 := sqlmock.NewRows([]string{"function", "next"}).
				AddRow(fn, nextKeys[2])
			rows3 := sqlmock.NewRows([]string{"function", "next"}).
				AddRow(fn, nextKeys[2])
			mock.ExpectQuery("^SELECT (.+) FROM configurations (.+) WHERE c.this = 0$").WillReturnRows(rows0)
			mock.ExpectQuery("^SELECT (.+) FROM configurations (.+)$").WillReturnRows(rows1)
			mock.ExpectQuery("^SELECT (.+) FROM configurations (.+)$").WillReturnRows(rows3)
			s = NewService(mockDB)
			cfg = msgRun.Body.Configs[index]

		})
		AfterEach(func() {
			db.Close()
		})

		Describe("Fetch Component Config", func() {
			It("should fetch component 1 config", func() {
				result := s.FetchComponentConfig(cfg, mockDB)
				Ω(result.ID).Should(Equal(cfg.ID))
				Ω(result.Status).Should(Equal(cfg.Status))
				Ω(result.Function).Should(Equal(fn))
				Ω(result.NextKeys[0]).Should(Equal(nextKeys[2]))
			})
		})

		Describe("Select Input", func() {
			It("should select input when component is in brain's nextKeys", func() {
				expected := msgRun.Body.Input[index]
				result := s.SelectInput(msgRun.Body.Input, cfg)
				Ω(result).Should(Equal(expected))
			})

		})

		It("should build message", func() {
			expected := Config{
				1,
				"up",
				"not",
				[]int{3},
			}
			result := s.BuildMessage(msgRun.Body.Input, cfg, mockDB)
			Ω(result.Configs[0]).Should(Equal(expected))
		})

		It("should run demo", func() {
			output, err := s.RunDemo(msgRun.Body, mockRabbit, mockDB)
			Ω(output).Should(BeTrue())
			Ω(err).Should(BeNil())
			Ω(mock.ExpectationsWereMet()).To(BeNil())
		})
	})
})
