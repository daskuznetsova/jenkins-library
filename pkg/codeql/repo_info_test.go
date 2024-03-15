package codeql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGitRepoInfo(t *testing.T) {
	t.Run("Valid https URL1", func(t *testing.T) {
		var repoInfo RepoInfo
		err := getGitRepoInfo("https://github.hello.test/Testing/fortify.git", &repoInfo)
		assert.NoError(t, err)
		assert.Equal(t, "https://github.hello.test", repoInfo.ServerUrl)
		assert.Equal(t, "fortify", repoInfo.Repo)
		assert.Equal(t, "Testing", repoInfo.Owner)
	})

	t.Run("Valid https URL2", func(t *testing.T) {
		var repoInfo RepoInfo
		err := getGitRepoInfo("https://github.hello.test/Testing/fortify", &repoInfo)
		assert.NoError(t, err)
		assert.Equal(t, "https://github.hello.test", repoInfo.ServerUrl)
		assert.Equal(t, "fortify", repoInfo.Repo)
		assert.Equal(t, "Testing", repoInfo.Owner)
	})
	t.Run("Valid https URL1 with dots", func(t *testing.T) {
		var repoInfo RepoInfo
		err := getGitRepoInfo("https://github.hello.test/Testing/com.sap.fortify.git", &repoInfo)
		assert.NoError(t, err)
		assert.Equal(t, "https://github.hello.test", repoInfo.ServerUrl)
		assert.Equal(t, "com.sap.fortify", repoInfo.Repo)
		assert.Equal(t, "Testing", repoInfo.Owner)
	})

	t.Run("Valid https URL2 with dots", func(t *testing.T) {
		var repoInfo RepoInfo
		err := getGitRepoInfo("https://github.hello.test/Testing/com.sap.fortify", &repoInfo)
		assert.NoError(t, err)
		assert.Equal(t, "https://github.hello.test", repoInfo.ServerUrl)
		assert.Equal(t, "com.sap.fortify", repoInfo.Repo)
		assert.Equal(t, "Testing", repoInfo.Owner)
	})
	t.Run("Valid https URL1 with username and token", func(t *testing.T) {
		var repoInfo RepoInfo
		err := getGitRepoInfo("https://username:token@github.hello.test/Testing/fortify.git", &repoInfo)
		assert.NoError(t, err)
		assert.Equal(t, "https://github.hello.test", repoInfo.ServerUrl)
		assert.Equal(t, "fortify", repoInfo.Repo)
		assert.Equal(t, "Testing", repoInfo.Owner)
	})

	t.Run("Valid https URL2 with username and token", func(t *testing.T) {
		var repoInfo RepoInfo
		err := getGitRepoInfo("https://username:token@github.hello.test/Testing/fortify", &repoInfo)
		assert.NoError(t, err)
		assert.Equal(t, "https://github.hello.test", repoInfo.ServerUrl)
		assert.Equal(t, "fortify", repoInfo.Repo)
		assert.Equal(t, "Testing", repoInfo.Owner)
	})

	t.Run("Invalid https URL as no org/Owner passed", func(t *testing.T) {
		var repoInfo RepoInfo
		assert.Error(t, getGitRepoInfo("https://github.com/fortify", &repoInfo))
	})

	t.Run("Invalid URL as no protocol passed", func(t *testing.T) {
		var repoInfo RepoInfo
		assert.Error(t, getGitRepoInfo("github.hello.test/Testing/fortify", &repoInfo))
	})

	t.Run("Valid ssh URL1", func(t *testing.T) {
		var repoInfo RepoInfo
		err := getGitRepoInfo("git@github.hello.test/Testing/fortify.git", &repoInfo)
		assert.NoError(t, err)
		assert.Equal(t, "https://github.hello.test", repoInfo.ServerUrl)
		assert.Equal(t, "fortify", repoInfo.Repo)
		assert.Equal(t, "Testing", repoInfo.Owner)
	})

	t.Run("Valid ssh URL2", func(t *testing.T) {
		var repoInfo RepoInfo
		err := getGitRepoInfo("git@github.hello.test/Testing/fortify", &repoInfo)
		assert.NoError(t, err)
		assert.Equal(t, "https://github.hello.test", repoInfo.ServerUrl)
		assert.Equal(t, "fortify", repoInfo.Repo)
		assert.Equal(t, "Testing", repoInfo.Owner)
	})
	t.Run("Valid ssh URL1 with dots", func(t *testing.T) {
		var repoInfo RepoInfo
		err := getGitRepoInfo("git@github.hello.test/Testing/com.sap.fortify.git", &repoInfo)
		assert.NoError(t, err)
		assert.Equal(t, "https://github.hello.test", repoInfo.ServerUrl)
		assert.Equal(t, "com.sap.fortify", repoInfo.Repo)
		assert.Equal(t, "Testing", repoInfo.Owner)
	})

	t.Run("Valid ssh URL2 with dots", func(t *testing.T) {
		var repoInfo RepoInfo
		err := getGitRepoInfo("git@github.hello.test/Testing/com.sap.fortify", &repoInfo)
		assert.NoError(t, err)
		assert.Equal(t, "https://github.hello.test", repoInfo.ServerUrl)
		assert.Equal(t, "com.sap.fortify", repoInfo.Repo)
		assert.Equal(t, "Testing", repoInfo.Owner)
	})

	t.Run("Invalid ssh URL as no org/Owner passed", func(t *testing.T) {
		var repoInfo RepoInfo
		assert.Error(t, getGitRepoInfo("git@github.com/fortify", &repoInfo))
	})
}

//func TestInitGitInfo(t *testing.T) {
//	t.Run("Valid URL1", func(t *testing.T) {
//		config := codeqlExecuteScanOptions{Repository: "https://github.hello.test/Testing/codeql.git", AnalyzedRef: "refs/head/branch", CommitID: "abcd1234"}
//		repoInfo, err := getRepoInfo(&config)
//		assert.NoError(t, err)
//		assert.Equal(t, "abcd1234", repoInfo.CommitId)
//		assert.Equal(t, "Testing", repoInfo.Owner)
//		assert.Equal(t, "codeql", repoInfo.Repo)
//		assert.Equal(t, "refs/head/branch", repoInfo.Ref)
//		assert.Equal(t, "https://github.hello.test", repoInfo.ServerUrl)
//	})
//
//	t.Run("Valid URL2", func(t *testing.T) {
//		config := codeqlExecuteScanOptions{Repository: "https://github.hello.test/Testing/codeql", AnalyzedRef: "refs/head/branch", CommitID: "abcd1234"}
//		repoInfo, err := getRepoInfo(&config)
//		assert.NoError(t, err)
//		assert.Equal(t, "abcd1234", repoInfo.CommitId)
//		assert.Equal(t, "Testing", repoInfo.Owner)
//		assert.Equal(t, "codeql", repoInfo.Repo)
//		assert.Equal(t, "refs/head/branch", repoInfo.Ref)
//		assert.Equal(t, "https://github.hello.test", repoInfo.ServerUrl)
//	})
//
//	t.Run("Valid url with dots URL1", func(t *testing.T) {
//		config := codeqlExecuteScanOptions{Repository: "https://github.hello.test/Testing/com.sap.codeql.git", AnalyzedRef: "refs/head/branch", CommitID: "abcd1234"}
//		repoInfo, err := getRepoInfo(&config)
//		assert.NoError(t, err)
//		assert.Equal(t, "abcd1234", repoInfo.CommitId)
//		assert.Equal(t, "Testing", repoInfo.Owner)
//		assert.Equal(t, "com.sap.codeql", repoInfo.Repo)
//		assert.Equal(t, "refs/head/branch", repoInfo.Ref)
//		assert.Equal(t, "https://github.hello.test", repoInfo.ServerUrl)
//	})
//
//	t.Run("Valid url with dots URL2", func(t *testing.T) {
//		config := codeqlExecuteScanOptions{Repository: "https://github.hello.test/Testing/com.sap.codeql", AnalyzedRef: "refs/head/branch", CommitID: "abcd1234"}
//		repoInfo, err := getRepoInfo(&config)
//		assert.NoError(t, err)
//		assert.Equal(t, "abcd1234", repoInfo.CommitId)
//		assert.Equal(t, "Testing", repoInfo.Owner)
//		assert.Equal(t, "com.sap.codeql", repoInfo.Repo)
//		assert.Equal(t, "refs/head/branch", repoInfo.Ref)
//		assert.Equal(t, "https://github.hello.test", repoInfo.ServerUrl)
//	})
//
//	t.Run("Valid url with username and token URL1", func(t *testing.T) {
//		config := codeqlExecuteScanOptions{Repository: "https://username:token@github.hello.test/Testing/codeql.git", AnalyzedRef: "refs/head/branch", CommitID: "abcd1234"}
//		repoInfo, err := getRepoInfo(&config)
//		assert.NoError(t, err)
//		assert.Equal(t, "abcd1234", repoInfo.CommitId)
//		assert.Equal(t, "Testing", repoInfo.Owner)
//		assert.Equal(t, "codeql", repoInfo.Repo)
//		assert.Equal(t, "refs/head/branch", repoInfo.Ref)
//		assert.Equal(t, "https://github.hello.test", repoInfo.ServerUrl)
//	})
//
//	t.Run("Valid url with username and token URL2", func(t *testing.T) {
//		config := codeqlExecuteScanOptions{Repository: "https://username:token@github.hello.test/Testing/codeql", AnalyzedRef: "refs/head/branch", CommitID: "abcd1234"}
//		repoInfo, err := getRepoInfo(&config)
//		assert.NoError(t, err)
//		assert.Equal(t, "abcd1234", repoInfo.CommitId)
//		assert.Equal(t, "Testing", repoInfo.Owner)
//		assert.Equal(t, "codeql", repoInfo.Repo)
//		assert.Equal(t, "refs/head/branch", repoInfo.Ref)
//		assert.Equal(t, "https://github.hello.test", repoInfo.ServerUrl)
//	})
//
//	t.Run("Invalid URL with no org/reponame", func(t *testing.T) {
//		config := codeqlExecuteScanOptions{Repository: "https://github.hello.test", AnalyzedRef: "refs/head/branch", CommitID: "abcd1234"}
//		repoInfo, err := getRepoInfo(&config)
//		assert.NoError(t, err)
//		_, err = orchestrator.GetOrchestratorConfigProvider(nil)
//		assert.Equal(t, "abcd1234", repoInfo.CommitId)
//		assert.Equal(t, "refs/head/branch", repoInfo.Ref)
//		if err != nil {
//			assert.Equal(t, "", repoInfo.Owner)
//			assert.Equal(t, "", repoInfo.Repo)
//			assert.Equal(t, "", repoInfo.ServerUrl)
//		}
//	})
//}
