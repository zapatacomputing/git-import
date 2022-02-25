package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	. "github.com/zapatacomputing/git-import"
)

var _ = Describe("Main", func() {
	var (
		url       string
		dir       string
		branch    string
		tag       string
		revision  string
		pathToCMD string
		session   *gexec.Session
		command   *exec.Cmd
		err       error
	)

	BeforeEach(func() {
		pathToCMD, err = gexec.Build("github.com/zapatacomputing/git-import")
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		curPath, err := os.Getwd()
		Expect(err).ShouldNot(HaveOccurred())
		os.RemoveAll(filepath.Join(curPath, dir))
	})

	Context("with all the arguments passed correctly", func() {
		BeforeEach(func() {
			url = "git@github.com:zapatacomputing/test.git"
			dir = "test"
			branch = "master"
			tag = ""
			revision = ""
			command = exec.Command(pathToCMD, "-url", url, "-dir", dir, "-branch", branch, "-tag", tag, "-revision", revision)
		})

		It("should exit successfully", func() {
			session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			time.Sleep(10 * time.Second)
			Expect(session).Should(gexec.Exit())
			Expect(session.Out).Should(gbytes.Say(`.*Enumerating objects*.`))
			Expect(session.Out).Should(gbytes.Say(`.*Counting objects*.`))
			Expect(session.Out).Should(gbytes.Say(`.*Compressing objects*.`))
			Expect(session.Out).Should(gbytes.Say(`.*Total*.`))
		})
	})

	Context("with a missing url", func() {
		BeforeEach(func() {
			url = ""
			dir = "test"
			branch = "master"
			tag = ""
			revision = ""
			command = exec.Command(pathToCMD, "-url", url, "-dir", dir, "-branch", branch, "-tag", tag, "-revision", revision)
		})

		It("should fail", func() {
			session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			time.Sleep(2 * time.Second)
			Expect(session).Should(gexec.Exit(1))
		})
	})
})

var _ = Describe("Clone", func() {
	var url string
	var dir string
	var branch string
	var tag string
	var revision string
	var expectedFilePath string
	var err error

	BeforeEach(func() {
		url = "git@github.com:zapatacomputing/test.git"
		dir = "test"
		branch = ""
		tag = ""
		revision = ""
		curPath, err := os.Getwd()
		Expect(err).ShouldNot(HaveOccurred())
		expectedFilePath = filepath.Join(curPath, dir, "README.md")
	})

	AfterEach(func() {
		curPath, err := os.Getwd()
		Expect(err).ShouldNot(HaveOccurred())
		os.RemoveAll(filepath.Join(curPath, dir))
	})

	Context("with a valid public repo, valid path and filled branch and and empty tag and empty revision", func() {
		BeforeEach(func() {
			branch = "master"
		})

		AfterEach(func() {
			branch = ""
		})

		It("should succeed", func() {
			err = Clone(url, dir, branch, tag, revision)
			Expect(err).ShouldNot(HaveOccurred())

			_, err := os.Stat(expectedFilePath)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("with a valid public repo, valid path and empty branch and and filled tag and empty revision", func() {
		BeforeEach(func() {
			tag = "v1.0.0"
		})

		AfterEach(func() {
			tag = ""
		})

		It("should succeed", func() {
			err = Clone(url, dir, branch, tag, revision)
			Expect(err).ShouldNot(HaveOccurred())

			_, err := os.Stat(expectedFilePath)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("with a valid public repo, valid path and empty branch and and empty tag and filled revision (branch)", func() {
		BeforeEach(func() {
			revision = "master"
		})

		AfterEach(func() {
			revision = ""
		})

		It("should succeed", func() {
			err = Clone(url, dir, branch, tag, revision)
			Expect(err).ShouldNot(HaveOccurred())

			_, err := os.Stat(expectedFilePath)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("with a valid public repo, valid path and empty branch and and empty tag and filled revision (tag)", func() {
		BeforeEach(func() {
			revision = "v1.0.0"
		})

		AfterEach(func() {
			revision = ""
		})

		It("should succeed", func() {
			err = Clone(url, dir, branch, tag, revision)
			Expect(err).ShouldNot(HaveOccurred())

			_, err := os.Stat(expectedFilePath)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("with a valid public repo, valid path and empty branch and and empty tag and filled revision (commit)", func() {
		BeforeEach(func() {
			revision = "c6ec0ddbb1a750e90e422653891ec8c3254e6c78"
		})

		AfterEach(func() {
			revision = ""
		})

		It("should succeed", func() {
			err = Clone(url, dir, branch, tag, revision)
			Expect(err).ShouldNot(HaveOccurred())

			_, err := os.Stat(expectedFilePath)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("with a valid public repo, valid path and filled branch and filled tag and empty revision", func() {
		BeforeEach(func() {
			branch = "master"
			tag = "v1.0.0"
		})

		AfterEach(func() {
			branch = ""
			tag = ""
		})

		It("should fail", func() {
			err = Clone(url, dir, branch, tag, revision)
			Expect(err).Should(HaveOccurred())
			Expect(strings.ToLower(err.Error())).Should(ContainSubstring("please specify only the branch, only the tag, or only the revision"))

			_, err := os.Stat(expectedFilePath)
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("with a valid public repo, valid path and filled branch and filled tag and filled revision", func() {
		BeforeEach(func() {
			branch = "master"
			tag = "v1.0.0"
			revision = "v1.0.0"
		})

		AfterEach(func() {
			branch = ""
			tag = ""
			revision = ""
		})

		It("should fail", func() {
			err = Clone(url, dir, branch, tag, revision)
			Expect(err).Should(HaveOccurred())
			Expect(strings.ToLower(err.Error())).Should(ContainSubstring("please specify only the branch, only the tag, or only the revision"))

			_, err := os.Stat(expectedFilePath)
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("with a valid public repo, valid path and filled branch and empty tag and filled revision", func() {
		BeforeEach(func() {
			branch = "master"
			revision = "v1.0.0"
		})

		AfterEach(func() {
			branch = ""
			tag = ""
			revision = ""
		})

		It("should fail", func() {
			err = Clone(url, dir, branch, tag, revision)
			Expect(err).Should(HaveOccurred())
			Expect(strings.ToLower(err.Error())).Should(ContainSubstring("please specify only the branch, only the tag, or only the revision"))

			_, err := os.Stat(expectedFilePath)
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("with a valid public repo, valid path empty branch and filled tag and filled revision", func() {
		BeforeEach(func() {
			tag = "v1.0.0"
			revision = "v1.0.0"
		})

		AfterEach(func() {
			branch = ""
			tag = ""
			revision = ""
		})

		It("should fail", func() {
			err = Clone(url, dir, branch, tag, revision)
			Expect(err).Should(HaveOccurred())
			Expect(strings.ToLower(err.Error())).Should(ContainSubstring("please specify only the branch, only the tag, or only the revision"))

			_, err := os.Stat(expectedFilePath)
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("with a valid public repo, valid path and empty branch and empty tag and empty revision", func() {
		It("should fail", func() {
			err = Clone(url, dir, branch, tag, revision)
			Expect(err).Should(HaveOccurred())

			_, err := os.Stat(expectedFilePath)
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("with a valid public repo, valid path and an invalid tag", func() {
		BeforeEach(func() {
			tag = "vnon-existent"
		})

		AfterEach(func() {
			tag = ""
		})

		It("should fail", func() {
			err = Clone(url, dir, branch, tag, revision)
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
			branch = ""
		})

		It("should fail", func() {
			err = Clone(url, dir, branch, tag, revision)
			Expect(err).Should(HaveOccurred())

			_, err := os.Stat(expectedFilePath)
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("with a valid public repo, valid path and an invalid revision", func() {
		BeforeEach(func() {
			revision = "#$%#"
		})

		AfterEach(func() {
			revision = ""
		})

		It("should fail", func() {
			err = Clone(url, dir, branch, tag, revision)
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
			err = Clone(url, dir, branch, tag, revision)
			Expect(err).Should(HaveOccurred())

			_, err := os.Stat(expectedFilePath)
			Expect(err).Should(HaveOccurred())
		})
	})
})
