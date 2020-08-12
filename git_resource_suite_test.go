package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGitImport(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GitImport Suite")
}
