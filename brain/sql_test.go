package brain_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/waltermblair/brain/brain"
)

// TODO - mock this out like core-test
var _ = Describe("database tests", func() {

	var client DBClient
	var err    error

	dbURL := "root:root@tcp(localhost:3306)/store"

	BeforeEach(func() {
		client, err = NewDBClient(dbURL)
	})

	PIt("should connect to db", func() {
		Ω(client.HealthCheck()).Should(BeNil())
	})

})