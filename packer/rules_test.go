package packer

import "testing"

func TestLangFromFileName(t *testing.T) {
	cases := []struct {
		name string
		lang string
	}{
		{
			name: "a.js",
			lang: "node",
		},
		{
			name: "a.py",
			lang: "python",
		},
		{
			name: "a.go",
			lang: "go",
		},
		{
			name: "a.rb",
			lang: "ruby",
		},
		{
			name: "a.php",
			lang: "php",
		},
		{
			name: "a.jl",
			lang: "julia",
		},
		{
			name: "a.d",
			lang: "d",
		},
		{
			name: "a.rs",
			lang: "rust",
		},
		{
			name: "a.java",
			lang: "java",
		},
		{
			name: "a.pl",
			lang: "perl",
		},
	}

	for _, c := range cases {
		lang, err := langFromFileName(c.name)
		if err != nil {
			t.Fatal(err)
		}
		if lang != c.lang {
			t.Fatalf("should get %s but got %s", c.lang, lang)
		}
	}
}
