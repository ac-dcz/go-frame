package geeweb

import (
	"testing"
)

func TestTire(t *testing.T) {
	url1 := "/a/b/c"
	url2 := "/a/:b/c"
	url3 := "a/*b"

	part1, part2, part3 := parsePattern(url1), parsePattern(url2), parsePattern(url3)
	tire := newTire()
	node1, node2, node3 := tire.insertRoute(url1, part1), tire.insertRoute(url2, part2), tire.insertRoute(url3, part3)
	n1, n2, n3 := tire.searchRoute(parsePattern(url1)), tire.searchRoute(parsePattern("/a/dcz/c")), tire.searchRoute(parsePattern("/a/src/file.txt"))

	t.Log(node1, n1)
	t.Log(node2, n2)
	t.Log(node3, n3)

	param2 := parseURL(node2.pattern, "/a/dcz/c")
	param3 := parseURL(node3.pattern, "/a/src/file.txt")

	t.Log(param2, param3)
}
