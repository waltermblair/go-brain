package brain_test

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/waltermblair/brain/brain"
)

var _ = Describe("database tests", func() {

	var mockDB     DBClient
	var db     	   *sql.DB
	var mock       sqlmock.Sqlmock

	BeforeEach(func() {
		db, mock, _ = sqlmock.New()
		mockDB = &DBClientImpl{
			"test",
			db,
		}
	})
	AfterEach(func() {
		db.Close()
	})

	It("should FetchNumInputs", func() {
		rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(3)
		mock.ExpectQuery("^SELECT COUNT").WillReturnRows(rows)
		count, err := mockDB.FetchNumInputs(0)
		Ω(count).Should(Equal(3))
		Ω(err).Should(BeNil())
		Ω(mock.ExpectationsWereMet()).To(BeNil())
	})

	It("should SetComponentStatus", func() {
		result := sqlmock.NewResult(0, 1)
		mock.ExpectExec("^UPDATE configurations").WillReturnResult(result)
		err := mockDB.SetComponentStatus(0, "up")
		Ω(err).Should(BeNil())
		Ω(mock.ExpectationsWereMet()).To(BeNil())
	})

	It("should FetchConfig", func() {
		rows := sqlmock.NewRows([]string{"function", "next"}).AddRow("buffer", 1)
		result := sqlmock.NewResult(0, 1)
		rows2 := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(3)
		mock.ExpectQuery("^SELECT (.+)").WillReturnRows(rows)
		mock.ExpectExec("^UPDATE configurations").WillReturnResult(result)
		mock.ExpectQuery("^SELECT COUNT").WillReturnRows(rows2)
		numInputs, nextKeys, fn, err := mockDB.FetchConfig(Config{})
		Ω(numInputs).Should(Equal(3))
		Ω(nextKeys).Should(Equal([]int{1}))
		Ω(fn).Should(Equal("buffer"))
		Ω(err).Should(BeNil())
		Ω(mock.ExpectationsWereMet()).To(BeNil())
	})

})