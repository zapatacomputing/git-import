package main_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/zapatacomputing/git-import"
)

var _ = Describe("Main", func() {
	var url string
	var dir string
	var branch string
	var expectedFilePath string
	var err error

	BeforeEach(func() {
		url = "git@github.com:zapatacomputing/test.git"
		dir = "test"
		branch = "master"
		curPath, err := os.Getwd()
		Expect(err).ShouldNot(HaveOccurred())
		expectedFilePath = filepath.Join(curPath, dir, "README.md")
	})

	AfterEach(func() {
		curPath, err := os.Getwd()
		Expect(err).ShouldNot(HaveOccurred())
		os.RemoveAll(filepath.Join(curPath, dir))
	})

	Context("with a valid public repo to clone and inputs", func() {
		It("should succeed", func() {
			err = Clone(url, dir, branch)
			Expect(err).ShouldNot(HaveOccurred())

			_, err := os.Stat(expectedFilePath)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when we expect the clone to fail", func() {
		It("should return an error", func() {
			err = Clone("", dir, branch)
			Expect(err).Should(HaveOccurred())
		})
	})
})
