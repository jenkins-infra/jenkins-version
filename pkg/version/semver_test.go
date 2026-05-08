package version_test

import (
	"fmt"
	"testing"

	"github.com/jenkins-infra/jenkins-version/pkg/version"

	"github.com/stretchr/testify/assert"
)

func TestSemVer(t *testing.T) {
	type test struct {
		v1       string
		v2       string
		lessThan bool
	}

	const alpha1 = "2.0-alpha-1"
	const alpha2 = "2.0-alpha-2"
	const alpha3 = "2.0-alpha-3"
	const alpha4 = "2.0-alpha-4"
	const beta1 = "2.0-beta-1"
	const beta2 = "2.0-beta-2"
	const release100 = "1.0.0"
	const release20 = "2.0"
	const release21 = "2.1"
	const release22 = "2.2"
	const release23 = "2.3"
	const release24 = "2.4"
	const release25 = "2.5"
	const releaseCandidate = "2.0-rc-1"
	tests := []test{
		{v1: release100, v2: "2.0.0", lessThan: true},
		{v1: release100, v2: "1.1.0", lessThan: true},
		{v1: release100, v2: "1.0.1", lessThan: true},
		{v1: "1.0.0-A", v2: "1.0.0-B", lessThan: true},
		{v1: alpha1, v2: release25, lessThan: true},
		{v1: release25, v2: alpha1, lessThan: false},
		{v1: release25, v2: alpha2, lessThan: false},
		{v1: release25, v2: alpha3, lessThan: false},
		{v1: release25, v2: alpha4, lessThan: false},
		{v1: release25, v2: beta1, lessThan: false},
		{v1: release25, v2: beta2, lessThan: false},
		{v1: release24, v2: releaseCandidate, lessThan: false},
		{v1: release25, v2: releaseCandidate, lessThan: false},
		{v1: releaseCandidate, v2: beta1, lessThan: false},
		{v1: alpha1, v2: beta1, lessThan: true},
		{v1: beta1, v2: alpha1, lessThan: false},
		{v1: beta1, v2: release24, lessThan: true},
		{v1: beta1, v2: alpha3, lessThan: false},
		{v1: beta1, v2: alpha4, lessThan: false},
		{v1: beta1, v2: beta2, lessThan: true},
		{v1: beta1, v2: alpha2, lessThan: false},
		{v1: beta1, v2: release20, lessThan: true},
		{v1: beta1, v2: release21, lessThan: true},
		{v1: beta1, v2: release22, lessThan: true},
		{v1: beta1, v2: release23, lessThan: true},
		{v1: release23, v2: beta1, lessThan: false},
		{v1: release22, v2: beta1, lessThan: false},
		{v1: release21, v2: beta1, lessThan: false},
		{v1: release20, v2: beta1, lessThan: false},
		{v1: alpha2, v2: beta1, lessThan: true},
		{v1: beta2, v2: beta1, lessThan: false},
		{v1: alpha1, v2: beta1, lessThan: true},
		{v1: alpha2, v2: beta1, lessThan: true},
		{v1: alpha3, v2: beta1, lessThan: true},
		{v1: alpha4, v2: beta1, lessThan: true},
		{v1: release24, v2: beta1, lessThan: false},
		{v1: "1.518.JENKINS-14362-jzlib", v2: "1.518", lessThan: true},
		{v1: "1.518", v2: "1.518.JENKINS-14362-jzlib", lessThan: false},
		{v1: "1.513.JENKINS-14362-jzlib", v2: "1.513", lessThan: true},
		{v1: "1.516.JENKINS-14362-jzlib", v2: "1.516", lessThan: true},
		{v1: release24, v2: release24, lessThan: false},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("compare %s < %s", tc.v1, tc.v2), func(t *testing.T) {
			version1 := version.NewVersion(tc.v1)
			version2 := version.NewVersion(tc.v2)
			assert.Equal(t, tc.lessThan, version1.LessThan(version2))
		})
	}
}
