package uuid

import (
	"crypto/md5"
	"fmt"
)

// This function is named ExampleNewMD5()
// it with the Examples type.
func ExampleNewMD5() {
	u := NewMD5(NameSpaceDNS, []byte("golang.org"))
	fmt.Println(u.String())
	// Output: c4e4c1e8-2e52-3de5-aacb-c9bc7208105d
}

// This function is named ExampleNewSHA1()
// it with the Examples type.
func ExampleNewSHA1() {
	u := NewSHA1(NameSpaceDNS, []byte("golang.org"))
	fmt.Println(u.String())
	// Output: 53447179-a84a-5086-927b-77f5951d9e4e
}

// This function is named ExampleNewHash()
// it with the Examples type.
func ExampleNewHash() {
	u := NewHash(md5.New(), NameSpaceURL, []byte("https://golang.org"), 3)
	fmt.Println(u.String())
	// Output: 282a1487-8b16-3fc9-ab77-be3810dacc14
}
