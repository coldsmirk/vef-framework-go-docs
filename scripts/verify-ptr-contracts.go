package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type corpus struct {
	label   string
	content string
}

func main() {
	sourceDir := flag.String("source", "../vef-framework-go", "path to the VEF Framework Go source repository")
	outDir := flag.String("out", ".", "path to the VEF Framework Go docs repository")
	flag.Parse()

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)

	englishDocs := readCorpus("English small utilities docs", filepath.Join(docsRoot, "docs/utilities/small-helpers.md"))
	chineseDocs := readCorpus("Chinese small utilities docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/utilities/small-helpers.md"))
	goMod := readCorpus("go.mod", filepath.Join(sourceRoot, "go.mod"))

	var failures []string

	// The ptr package was deliberately removed (refactor!: drop the ptr
	// package in favor of samber/lo and builtin new); the docs must not
	// silently regress into describing it as a live package again.
	ptrDir := filepath.Join(sourceRoot, "ptr")
	if _, err := os.Stat(ptrDir); !os.IsNotExist(err) {
		failures = append(failures, "ptr package directory unexpectedly exists at "+ptrDir+" but docs describe it as removed")
	}

	failures = append(failures, missingTerms(goMod, []string{"github.com/samber/lo"})...)

	failures = append(failures, missingTerms(englishDocs, []string{
		"`ptr` package has been removed",
		"builtin `new`",
		"`ptr.Of(v)`",
		"`new(v)` (builtin)",
		"`lo.EmptyableToPtr(v)`",
		"`ptr.Zero[T]()`",
		"`lo.Empty[T]()`",
		"`ptr.Value(p)`",
		"`lo.FromPtr(p)`",
		"`ptr.Value(p, fallback)` / `ptr.ValueOrElse(p, fn)`",
		"`lo.FromPtrOr(p, fallback)`",
		"`ptr.Coalesce(p1, p2, ...)`",
		"import \"github.com/samber/lo\"",
	})...)
	failures = append(failures, missingTerms(chineseDocs, []string{
		"`ptr` 包已从框架中移除",
		"内置的 `new`",
		"`ptr.Of(v)`",
		"`new(v)`（builtin）",
		"`lo.EmptyableToPtr(v)`",
		"`ptr.Zero[T]()`",
		"`lo.Empty[T]()`",
		"`ptr.Value(p)`",
		"`lo.FromPtr(p)`",
		"`ptr.Value(p, fallback)` / `ptr.ValueOrElse(p, fn)`",
		"`lo.FromPtrOr(p, fallback)`",
		"`ptr.Coalesce(p1, p2, ...)`",
		"import \"github.com/samber/lo\"",
	})...)

	sort.Strings(failures)
	if len(failures) > 0 {
		panic(fmt.Errorf("ptr contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	fmt.Println("PTR contract docs verified: package removal documented, migration table present, samber/lo dependency confirmed, 2 doc mirrors")
}

func readCorpus(label, path string) corpus {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read %s at %s: %w", label, path, err))
	}

	return corpus{label: label, content: string(content)}
}

func missingTerms(c corpus, terms []string) []string {
	var failures []string
	for _, term := range terms {
		failures = append(failures, missingTerm(c, term)...)
	}

	return failures
}

func missingTerm(c corpus, term string) []string {
	if !strings.Contains(c.content, term) {
		return []string{fmt.Sprintf("%s missing term: %s", c.label, term)}
	}

	return nil
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}
