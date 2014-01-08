package debian

import (
	"strconv"
	"strings"
	"unicode"
)

// Using documentation from: http://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-Version

// CompareVersions compares two package versions
func CompareVersions(ver1, ver2 string) int {
	e1, u1, d1 := parseVersion(ver1)
	e2, u2, d2 := parseVersion(ver2)

	r := compareVersionPart(e1, e2)
	if r != 0 {
		return r
	}

	r = compareVersionPart(u1, u2)
	if r != 0 {
		return r
	}

	return compareVersionPart(d1, d2)
}

// parseVersions breaks down full version to components (possibly empty)
func parseVersion(ver string) (epoch, upstream, debian string) {
	i := strings.LastIndex(ver, "-")
	if i != -1 {
		debian, ver = ver[i+1:], ver[:i]
	}

	i = strings.Index(ver, ":")
	if i != -1 {
		epoch, ver = ver[:i], ver[i+1:]
	}

	upstream = ver

	return
}

// compareLexicographic compares in "Debian lexicographic" way, see below compareVersionPart for details
func compareLexicographic(s1, s2 string) int {
	i := 0
	l1, l2 := len(s1), len(s2)

	for {
		if i == l1 && i == l2 {
			// s1 equal to s2
			break
		}

		if i == l2 {
			// s1 is longer than s2
			if s1[i] == '~' {
				return -1 // s1 < s2
			}
			return 1 // s1 > s2
		}

		if i == l1 {
			// s2 is longer than s1
			if s2[i] == '~' {
				return 1 // s1 > s2
			}
			return -1 // s1 < s2
		}

		if s1[i] == s2[i] {
			i++
			continue
		}

		if s1[i] == '~' {
			return -1
		}

		if s2[i] == '~' {
			return 1
		}

		c1, c2 := unicode.IsLetter(rune(s1[i])), unicode.IsLetter(rune(s2[i]))
		if c1 && !c2 {
			return -1
		}
		if !c1 && c2 {
			return 1
		}

		if s1[i] < s2[i] {
			return -1
		}
		return 1
	}
	return 0
}

// compareVersionPart compares parts of full version
//
// From Debian Policy Manual:
//
// "The strings are compared from left to right.
//
// First the initial part of each string consisting entirely of non-digit characters is
// determined. These two parts (one of which may be empty) are compared lexically. If a
// difference is found it is returned. The lexical comparison is a comparison of ASCII values
// modified so that all the letters sort earlier than all the non-letters and so that a tilde
// sorts before anything, even the end of a part. For example, the following parts are in sorted
// order from earliest to latest: ~~, ~~a, ~, the empty part.
//
// Then the initial part of the remainder of each string which consists entirely of digit
// characters is determined. The numerical values of these two parts are compared, and any difference
// found is returned as the result of the comparison. For these purposes an empty string (which can only occur at
// the end of one or both version strings being compared) counts as zero.

// These two steps (comparing and removing initial non-digit strings and initial digit strings) are
// repeated until a difference is found or both strings are exhausted."
func compareVersionPart(part1, part2 string) int {
	i1, i2 := 0, 0
	l1, l2 := len(part1), len(part2)

	for {
		j1, j2 := i1, i2
		for j1 < l1 && !unicode.IsDigit(rune(part1[j1])) {
			j1++
		}

		for j2 < l2 && !unicode.IsDigit(rune(part2[j2])) {
			j2++
		}

		s1, s2 := part1[i1:j1], part2[i2:j2]
		r := compareLexicographic(s1, s2)
		if r != 0 {
			return r
		}

		i1, i2 = j1, j2

		for j1 < l1 && unicode.IsDigit(rune(part1[j1])) {
			j1++
		}

		for j2 < l2 && unicode.IsDigit(rune(part2[j2])) {
			j2++
		}

		s1, s2 = part1[i1:j1], part2[i2:j2]
		n1, _ := strconv.Atoi(s1)
		n2, _ := strconv.Atoi(s2)

		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}

		i1, i2 = j1, j2

		if i1 == l1 && i2 == l2 {
			break
		}
	}
	return 0
}