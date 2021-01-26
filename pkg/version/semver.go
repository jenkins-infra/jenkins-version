package version

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Semver struct {
	Major string
	Minor string
	Patch string
	Other string
}

func NewVersion(in string) Semver {
	parts := strings.Split(in, ".")
	if len(parts) == 1 {
		return Semver{Major: parts[0]}
	} else if len(parts) == 2 {
		return Semver{Major: parts[0], Minor: parts[1]}
	} else if len(parts) == 3 {
		return Semver{Major: parts[0], Minor: parts[1], Patch: parts[2]}
	} else {
		return Semver{Major: parts[0], Minor: parts[1], Patch: parts[2], Other: parts[3]}
	}
}

func toInt(in string) (int, error) {
	s, err := strconv.Atoi(in)
	if err != nil {
		return -1, errors.Wrap(err, "unable to parse as int: "+in)
	}
	return s, nil
}

func (v *Semver) String() string {
	return fmt.Sprintf("%d.%d.%d.%s", v.Major, v.Minor, v.Patch, v.Other)
}

func (v Semver) LessThan(o Semver) bool {
	if v.Major != o.Major {
		return v.Major < o.Major
	}

	if v.Minor != o.Minor {
		return v.Minor < o.Minor
	}

	if v.Patch != o.Patch {
		return v.Patch < o.Patch
	}

	return strings.Compare(v.Other, o.Other) < 0
}

type bySemVer []string

func (s bySemVer) Len() int {
	return len(s)
}
func (s bySemVer) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s bySemVer) Less(i, j int) bool {
	v1 := NewVersion(s[i])
	v2 := NewVersion(s[j])

	return v1.LessThan(v2)
}
