package util

import "testing"
import "fmt"

func TestMapSort1(t *testing.T) {
	m := make([]map[string]string, 8)
	m = []map[string]string{
		{"action": "a", "file": "a.php"},
		{"action": "a", "file": "/a/b/c/c.jpg"},
		{"action": "a", "file": "/a/b/d/b.jpg"},
		{"action": "a", "file": "/a/b"},
		{"action": "a", "file": "/a"},
		{"action": "a", "file": "/a/b/c"},
		{"action": "a", "file": "/a/b/d"},
		{"action": "a", "file": "/a/b/c/a.jpg"},
	}
	MapSort(m, func(p1, p2 map[string]string) bool {
		return p1["file"] < p2["file"]
	})
	for _, v := range m {
		fmt.Println(v["file"])
	}
}
