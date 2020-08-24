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
	var tag string
	var expectedFilePath string
	var err error

	BeforeEach(func() {
		url = "git@github.com:zapatacomputing/test.git"
		dir = "test"
		branch = "master"
		tag = ""
		curPath, err := os.Getwd()
		Expect(err).ShouldNot(HaveOccurred())
		expectedFilePath = filepath.Join(curPath, dir, "README.md")
	})

	AfterEach(func() {
		curPath, err := os.Getwd()
		Expect(err).ShouldNot(HaveOccurred())
		os.RemoveAll(filepath.Join(curPath, dir))
	})

	Context("with a valid public repo, valid path and filled branch and and empty tag", func() {
		It("should succeed", func() {
			err = Clone(url, dir, branch, tag)
			Expect(err).ShouldNot(HaveOccurred())

			_, err := os.Stat(expectedFilePath)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("with a valid public repo, valid path and filled branch and filled tag", func() {
		BeforeEach(func() {
			tag = "v1.0.0"
		})

		AfterEach(func() {
			tag = ""
		})

		It("should succeed", func() {
			err = Clone(url, dir, branch, tag)
			Expect(err).ShouldNot(HaveOccurred())

			_, err := os.Stat(expectedFilePath)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("with a valid public repo, valid path and empty branch and empty tag", func() {
		BeforeEach(func() {
			branch = ""
		})

		AfterEach(func() {
			branch = "master"
		})

		It("should succeed", func() {
			err = Clone(url, dir, branch, tag)
			Expect(err).ShouldNot(HaveOccurred())

			_, err := os.Stat(expectedFilePath)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("with a valid public repo, valid path and an invalid tag", func() {
		BeforeEach(func() {
			tag = "vnon-existent"
		})

		AfterEach(func() {
			tag = ""
		})

		It("should succeed", func() {
			err = Clone(url, dir, branch, tag)
			Expect(err).Should(HaveOccurred())

			_, err := os.Stat(expectedFilePath)
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("with a valid public repo, valid path and an invalid branch", func() {
		BeforeEach(func() {
			branch = "#$%#"
		})

		AfterEach(func() {
			branch = "master"
		})

		It("should fail", func() {
			err = Clone(url, dir, branch, tag)
			Expect(err).Should(HaveOccurred())

			_, err := os.Stat(expectedFilePath)
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("with a valid public repo, invalid path", func() {
		BeforeEach(func() {
			dir = ""
		})

		AfterEach(func() {
			dir = "test"
		})

		It("should fail", func() {
			err = Clone(url, dir, branch, tag)
			Expect(err).Should(HaveOccurred())

			_, err := os.Stat(expectedFilePath)
			Expect(err).Should(HaveOccurred())
		})
	})
})
