package resolver

import (
	"fmt"
	"regexp"
)

var golangRe = regexp.MustCompile(`^golang\.org/x/([^/]+)$`)
var k8sRe = regexp.MustCompile(`^k8s\.io/([^/]+)$`)
var k8sSigsRe = regexp.MustCompile(`^sigs\.k8s\.io/([^/]+)$`)
var gopkgRe = regexp.MustCompile(`(?i)^gopkg\.in/(?:([a-zA-Z0-9][-a-zA-Z0-9]+)/)?([a-zA-Z][-.a-zA-Z0-9]*)\.((?:v0|v[1-9][0-9]*)(?:\.0|\.[1-9][0-9]*){0,2}(?:-unstable)?)(?:\.git)?((?:/[a-zA-Z0-9][-.a-zA-Z0-9]*)*)$`)

func Resolve(name string) string {
	if matches := golangRe.FindStringSubmatch(name); matches != nil {
		return fmt.Sprintf("github.com/golang/%s", matches[1])
	}

	if matches := k8sRe.FindStringSubmatch(name); matches != nil {
		return fmt.Sprintf("github.com/kubernetes/%s", matches[1])
	}

	if matches := k8sSigsRe.FindStringSubmatch(name); matches != nil {
		return fmt.Sprintf("github.com/kubernetes-sigs/%s", matches[1])
	}

	if matches := gopkgRe.FindStringSubmatch(name); matches != nil {
		// URL case 1 with no user means it is go-<pkg>
		if matches[1] == "" {
			matches[1] = "go-" + matches[2]
		}
		return fmt.Sprintf("github.com/%s/%s", matches[1], matches[2])
	}

	return name
}
