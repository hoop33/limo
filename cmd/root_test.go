package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// `values` can be:
// * Full URL (e.g., https://github.com/hoop33/limo)
// * Owner/Repo (e.g., hoop33/limo)
// * Owner Repo (e.g., hoop33 limo)

func TestParseServiceOwnerRepoShouldParseEmptyWhenNil(t *testing.T) {
	s, o, r := parseServiceOwnerRepo([]string{})
	assert.Equal(t, "", s)
	assert.Equal(t, "", o)
	assert.Equal(t, "", r)
}

func TestParseServiceOwnerRepoShouldParseEmptyWhenEmpty(t *testing.T) {
	s, o, r := parseServiceOwnerRepo([]string{})
	assert.Equal(t, "", s)
	assert.Equal(t, "", o)
	assert.Equal(t, "", r)
}

func TestParseServiceOwnerRepoShouldParseEmptyWhenOneValueMalformatted(t *testing.T) {
	s, o, r := parseServiceOwnerRepo([]string{"this is a test"})
	assert.Equal(t, "", s)
	assert.Equal(t, "", o)
	assert.Equal(t, "", r)
}

func TestParseServiceOwnerRepoShouldParseEmptyWhenTooManyValues(t *testing.T) {
	s, o, r := parseServiceOwnerRepo([]string{"this", "is", "a", "test"})
	assert.Equal(t, "", s)
	assert.Equal(t, "", o)
	assert.Equal(t, "", r)
}

func TestParseServiceOwnerRepoShouldParseAllWhenOneValueFullURL(t *testing.T) {
	s, o, r := parseServiceOwnerRepo([]string{"https://github.com/hoop33/limo"})
	assert.Equal(t, "github", s)
	assert.Equal(t, "hoop33", o)
	assert.Equal(t, "limo", r)
}

func TestParseServiceOwnerRepoShouldParseOwnerRepoWhenOneValueNoService(t *testing.T) {
	s, o, r := parseServiceOwnerRepo([]string{"hoop33/limo"})
	assert.Equal(t, "", s)
	assert.Equal(t, "hoop33", o)
	assert.Equal(t, "limo", r)
}

func TestParseServiceOwnerRepoShouldParseOwnerRepoWhenTwoValuesOwnerRepo(t *testing.T) {
	s, o, r := parseServiceOwnerRepo([]string{"hoop33", "limo"})
	assert.Equal(t, "", s)
	assert.Equal(t, "hoop33", o)
	assert.Equal(t, "limo", r)
}
