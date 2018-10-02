package brain_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/waltermblair/brain/brain"
)

var _ = Describe("database tests", func() {

	var client DBClient
	var err    error

	dbURL := "root:root@tcp(localhost:3306)/store"

	BeforeEach(func() {
		client, err = NewDBClient(dbURL)
	})

	It("should connect to db", func() {
		Ω(client.HealthCheck()).Should(BeNil())
	})

})