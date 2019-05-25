package digitalblasphemy

import (
	"fmt"
	"regexp"
)

const (
	urlMembers  = "https://secure.digitalblasphemy.com/content/"
	urlFreebies = "http://digitalblasphemy.com/seeall.shtml?y=todos&t=0&w=&h=&r=1&f=1"

	freebieURLPattern = "https://secure.digitalblasphemy.com/graphics/HDfree/%sHDfree.jpg"

	indexURLPrefix = "https://secure.digitalblasphemy.com/content/jpgs/"
	indexURLSuffix = "/"
	indexURLSort   = "?C=M;O=D"
)

func makeIndexURL(suffix string) string {
	return indexURLPrefix + suffix + indexURLSuffix
}

var indexURLForResolution = map[string]string{
	"1024x768":  makeIndexURL("1024st"),
	"1152x864":  makeIndexURL("db"),
	"1600x1200": makeIndexURL("1600"),
	"1280x1024": makeIndexURL("1280"),
	"1280x800":  makeIndexURL("1280w"),
	"1366x768":  makeIndexURL("1366"),
	"1600x900":  makeIndexURL("1600x900"),
	"1920x1080": makeIndexURL("1080p"),
	"2560x1440": makeIndexURL("1440p"),
	"1440x900":  makeIndexURL("1440"),
	"1680x1050": makeIndexURL("1680"),
	"1920x1200": makeIndexURL("widescreen"),
	"2560x1600": makeIndexURL("widescreen"),
	"3440x1440": makeIndexURL("21x9"),
	"2880x1800": makeIndexURL("widescreen"),
	"3840x2160": makeIndexURL("4k"),
	"4096x2304": makeIndexURL("4k"),
	"5120x2880": makeIndexURL("5k"),
}

const (
	regexpURLPrefix = "^"
	regexpURLSuffix = `\.jpg$`
	regexpID        = `([\w-]+)`
)

func makeURLRegexp(base string) *regexp.Regexp {
	return regexp.MustCompile(regexpURLPrefix + fmt.Sprintf(base, regexpID) + regexpURLSuffix)
}

var urlRegexpForResolution = map[string]*regexp.Regexp{
	"1024x768":  makeURLRegexp("%s1024st"),
	"1152x864":  makeURLRegexp("%s"),
	"1600x1200": makeURLRegexp("%s1600"),
	"1280x1024": makeURLRegexp("%s1280"),
	"1280x800":  makeURLRegexp("%s1280w"),
	"1366x768":  makeURLRegexp("%s1366"),
	"1600x900":  makeURLRegexp("%s1600x900"),
	"1920x1080": makeURLRegexp("%s1080p"),
	"2560x1440": makeURLRegexp("%s1440p"),
	"1440x900":  makeURLRegexp("%s1440"),
	"1680x1050": makeURLRegexp("%s1680"),
	"1920x1200": makeURLRegexp("%s1920"),
	"2560x1600": makeURLRegexp("%s2560"),
	"3440x1440": makeURLRegexp("%s3440"),
	"2880x1800": makeURLRegexp("%s2880"),
	"3840x2160": makeURLRegexp("%suhd"),
	"4096x2304": makeURLRegexp("%s4ktv"),
	"5120x2880": makeURLRegexp("%s5ktv"),
}

func makeFreebieURL(id string) string {
	return fmt.Sprintf(freebieURLPattern, id)
}
