package testdata

import (
	"testing"

	"github.com/garethjevans/jenkins-version/pkg/version"
	"github.com/stretchr/testify/assert"
)

func TestSemVer(t *testing.T) {
	assert.True(t, version.NewVersion("1.0.0").LessThan(version.NewVersion("2.0.0")))
	assert.True(t, version.NewVersion("1.0.0").LessThan(version.NewVersion("1.1.0")))
	assert.True(t, version.NewVersion("1.0.0").LessThan(version.NewVersion("1.0.1")))
	assert.True(t, version.NewVersion("1.0.0.A").LessThan(version.NewVersion("1.0.0.B")))
}
