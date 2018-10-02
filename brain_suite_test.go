package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBrain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Brain Suite")
}
