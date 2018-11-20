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

	var msgRun		Message
	var msgConfig	Message
	var cfg    	    Config

	BeforeSuite(func() {
		file, _ := os.Open("../resources/json/run.json")
		bytes, _ := ioutil.ReadAll(file)
		json.Unmarshal(bytes, &msgRun)
		file, _ = os.Open("../resources/json/msgConfig.json")
		bytes, _ = ioutil.ReadAll(file)
		json.Unmarshal(bytes, &msgConfig)
	})

	Describe("Core", func() {

		var s			Service
		var mockRabbit  RabbitClient
		var mockDB      DBClient
		var db     	    *sql.DB
		var mock        sqlmock.Sqlmock

		BeforeEach(func() {
			mockRabbit = NewMockRabbitClient(cfg)
			db, mock, _ = sqlmock.New()
			mockDB = &DBClientImpl{
				"test",
				db,
			}
			// Expectations each time service is created
			rows := sqlmock.NewRows([]string{"function", "next"}).
				AddRow("buffer", 1).
				AddRow("buffer", 2)
			result := sqlmock.NewResult(0, 1)
			rows2 := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)
			mock.ExpectQuery("^SELECT (.+)").WillReturnRows(rows)
			mock.ExpectExec("^UPDATE configurations").WillReturnResult(result)
			mock.ExpectQuery("^SELECT COUNT").WillReturnRows(rows2)
			s, _ = NewService(mockDB)
		})
		AfterEach(func() {
			db.Close()
		})

		It("should NewService", func() {
			Ω(s.GetConfig().ID).Should(Equal(0))
			Ω(s.GetConfig().Status).Should(Equal(""))
			Ω(s.GetConfig().Function).Should(Equal(""))
			Ω(s.GetConfig().NumInputs).Should(Equal(1))
			Ω(s.GetConfig().NextKeys).Should(Equal([]int{1,2}))
			Ω(mock.ExpectationsWereMet()).To(BeNil())
		})

		It("should FetchComponentConfig", func() {
			rows3 := sqlmock.NewRows([]string{"function", "next"}).
				AddRow("not", 2).
				AddRow("not", 3)
			result := sqlmock.NewResult(1, 1)
			rows4 := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)
			mock.ExpectQuery("^SELECT (.+)").WillReturnRows(rows3)
			mock.ExpectExec("^UPDATE configurations").WillReturnResult(result)
			mock.ExpectQuery("^SELECT COUNT").WillReturnRows(rows4)
			cfg, err := s.FetchComponentConfig(msgConfig.Body.Configs[0], mockDB)
			Ω(cfg.ID).Should(Equal(1))
			Ω(err).Should(BeNil())
			Ω(mock.ExpectationsWereMet()).To(BeNil())
		})

		It("should BuildInputMessage", func() {
			msg := s.BuildInputMessage(true)
			Ω(msg.Input[0]).Should(BeTrue())
		})

		It("should BuildConfigMessage", func() {
			rows5 := sqlmock.NewRows([]string{"function", "next"}).
				AddRow("not", 2).
				AddRow("not", 3)
			result := sqlmock.NewResult(1, 1)
			rows6 := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)
			mock.ExpectQuery("^SELECT (.+)").WillReturnRows(rows5)
			mock.ExpectExec("^UPDATE configurations").WillReturnResult(result)
			mock.ExpectQuery("^SELECT COUNT").WillReturnRows(rows6)
			msg := s.BuildConfigMessage(msgConfig.Body.Configs[0], mockDB)
			Ω(msg.Input).Should(BeNil())
			Ω(msg.Configs[0].ID).Should(Equal(1))
			Ω(msg.Configs[0].Status).Should(Equal("up"))
			Ω(msg.Configs[0].Function).Should(Equal("not"))
			Ω(msg.Configs[0].NumInputs).Should(Equal(1))
			Ω(msg.Configs[0].NextKeys).Should(Equal([]int{2, 3}))
			Ω(mock.ExpectationsWereMet()).To(BeNil())
		})

		// TODO test runDemo
		It("should RunDemo", func() {

		})
	})
})
