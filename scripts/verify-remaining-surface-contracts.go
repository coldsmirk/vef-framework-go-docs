package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

const (
	sourceModule     = "github.com/coldsmirk/vef-framework-go"
	englishIndexPath = "docs/reference/public-api-index.md"
	chineseIndexPath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"
)

var auditedPackages = []string{
	"github.com/coldsmirk/vef-framework-go",
	"github.com/coldsmirk/vef-framework-go/cqrs",
	"github.com/coldsmirk/vef-framework-go/cryptox",
	"github.com/coldsmirk/vef-framework-go/js",
	"github.com/coldsmirk/vef-framework-go/mold",
	"github.com/coldsmirk/vef-framework-go/monitor",
	"github.com/coldsmirk/vef-framework-go/orm",
	"github.com/coldsmirk/vef-framework-go/password",
	"github.com/coldsmirk/vef-framework-go/schema",
	"github.com/coldsmirk/vef-framework-go/search",
	"github.com/coldsmirk/vef-framework-go/tabular",
	"github.com/coldsmirk/vef-framework-go/timex",
}

var passThroughDocTerms = map[string][]string{
	"github.com/coldsmirk/vef-framework-go/js": {
		"goja pass-through surface",
		"github.com/dop251/goja",
		"js.AstProgram",
		"public API index",
	},
	"github.com/coldsmirk/vef-framework-go/orm": {
		"Bun pass-through surface",
		"github.com/uptrace/bun",
		"public API index",
	},
}

var familyIndexedMethodPackages = map[string][]string{
	"github.com/coldsmirk/vef-framework-go/orm": {
		"VEF-owned ORM method families",
		"receiver/category",
		"public API index",
	},
}

var standardPolicyTerms = map[string][]string{
	"github.com/coldsmirk/vef-framework-go/tabular": {
		"errors.Unwrap",
	},
	"github.com/coldsmirk/vef-framework-go/timex": {
		"MarshalJSON",
		"UnmarshalJSON",
		"MarshalText",
		"UnmarshalText",
		"Scan",
		"Value",
	},
}

type auditLedger struct {
	Entries []auditEntry `json:"entries"`
}

type auditEntry struct {
	ID        string   `json:"id"`
	Package   string   `json:"package"`
	Kind      string   `json:"kind"`
	Symbol    string   `json:"symbol"`
	Signature string   `json:"signature"`
	Coverage  []string `json:"coverage"`
}

type manifest struct {
	Packages []manifestEntry `json:"packages"`
}

type manifestEntry struct {
	Package     string   `json:"package"`
	Tier        string   `json:"tier"`
	Coverage    []string `json:"coverage"`
	TopLevel    int      `json:"top_level"`
	Fields      int      `json:"fields"`
	Methods     int      `json:"methods"`
	Fingerprint string   `json:"fingerprint"`
}

type corpus struct {
	label   string
	content string
}

type aliasSurface struct {
	typeAliases  map[string]string
	valueAliases map[string]string
}

type tierCounts struct {
	A int
	B int
	C int
}

func main() {
	sourceDir := flag.String("source", "../vef-framework-go", "path to the VEF Framework Go source repository")
	outDir := flag.String("out", ".", "path to the VEF Framework Go docs repository")
	flag.Parse()

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)

	ledger := loadJSON[auditLedger](filepath.Join(docsRoot, "scripts/api-audit-ledger.json"))
	manifestByPackage := loadManifestByPackage(filepath.Join(docsRoot, "scripts/api-audit-manifest.json"))
	entriesByPackage := auditEntriesByPackage(ledger.Entries)
	englishIndex := readCorpus("English public API index", filepath.Join(docsRoot, englishIndexPath))
	chineseIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, chineseIndexPath))
	aliases := loadAliasSurfaces(sourceRoot)

	var failures []string
	summaries := make(map[string]tierCounts)
	for _, pkg := range auditedPackages {
		manifestEntry, ok := manifestByPackage[pkg]
		if !ok {
			failures = append(failures, "audited package missing from manifest: "+pkg)
			continue
		}

		docs := packageCorpora(docsRoot, manifestEntry)
		entries := entriesByPackage[pkg]
		failures = append(failures, verifyPackageSurface(pkg, manifestEntry, entries)...)
		failures = append(failures, verifyPackageDocs(pkg, entries, docs, []corpus{englishIndex, chineseIndex}, aliases[pkg], &summaries)...)
	}

	sort.Strings(failures)
	if len(failures) > 0 {
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Println("Remaining public API surface contracts verified:")
	for _, pkg := range auditedPackages {
		counts := summaries[pkg]
		fmt.Printf("- %s: tier A=%d, tier B=%d, tier C=%d\n", pkg, counts.A, counts.B, counts.C)
	}
}

func verifyPackageSurface(pkg string, manifestEntry manifestEntry, entries []auditEntry) []string {
	var failures []string
	if len(entries) != manifestEntry.TopLevel+manifestEntry.Fields+manifestEntry.Methods {
		failures = append(failures, fmt.Sprintf(
			"%s entry count mismatch: got %d want %d",
			pkg,
			len(entries),
			manifestEntry.TopLevel+manifestEntry.Fields+manifestEntry.Methods,
		))
	}
	counts := map[string]int{}
	for _, entry := range entries {
		counts[entry.Kind]++
		if entry.Symbol == "" || entry.Signature == "" {
			failures = append(failures, pkg+" entry missing symbol/signature: "+entry.ID)
		}
	}
	if counts["top"] != manifestEntry.TopLevel || counts["field"] != manifestEntry.Fields || counts["method"] != manifestEntry.Methods {
		failures = append(failures, fmt.Sprintf(
			"%s kind count mismatch: got top/field/method=%d/%d/%d want=%d/%d/%d",
			pkg,
			counts["top"], counts["field"], counts["method"],
			manifestEntry.TopLevel, manifestEntry.Fields, manifestEntry.Methods,
		))
	}

	return failures
}

func verifyPackageDocs(
	pkg string,
	entries []auditEntry,
	docs []corpus,
	indexes []corpus,
	aliases aliasSurface,
	summaries *map[string]tierCounts,
) []string {
	var failures []string
	counts := tierCounts{}
	hasTierB := false
	hasTierC := false
	hasFamilyIndexedA := false

	for _, entry := range entries {
		for _, index := range indexes {
			if !indexMentionsSignature(index.content, entry.Package, entry.Signature) {
				failures = append(failures, fmt.Sprintf("%s missing index signature for %s: %s", index.label, entry.ID, entry.Signature))
			}
		}

		tier := classify(entry, aliases)
		switch tier {
		case "A":
			counts.A++
			if entry.Kind == "method" && isFamilyIndexedMethodPackage(pkg) {
				hasFamilyIndexedA = true
				failures = append(failures, verifyReceiverDocumented(entry, docs)...)
				continue
			}
			for _, doc := range docs {
				if !documentMentionsEntry(doc.content, entry) {
					failures = append(failures, fmt.Sprintf("%s missing tier A surface term for %s", doc.label, entry.Symbol))
				}
			}
		case "B":
			counts.B++
			hasTierB = true
		case "C":
			counts.C++
			hasTierC = true
			failures = append(failures, verifyPassThroughEvidence(entry, aliases)...)
		default:
			failures = append(failures, "unclassified public API entry: "+entry.ID)
		}
	}

	if hasTierB {
		failures = append(failures, verifyPolicyTerms(pkg, docs, standardPolicyTerms[pkg], "tier B standard-interface policy")...)
	}
	if hasTierC {
		failures = append(failures, verifyPolicyTerms(pkg, docs, passThroughDocTerms[pkg], "tier C pass-through policy")...)
	}
	if hasFamilyIndexedA {
		failures = append(failures, verifyPolicyTerms(pkg, docs, familyIndexedMethodPackages[pkg], "tier A method-family policy")...)
	}

	(*summaries)[pkg] = counts

	return failures
}

func classify(entry auditEntry, aliases aliasSurface) string {
	switch entry.Kind {
	case "top":
		if _, ok := aliases.typeAliases[entry.Symbol]; ok {
			return "C"
		}
		if _, ok := aliases.valueAliases[entry.Symbol]; ok {
			return "C"
		}
	case "method":
		receiver := receiverName(entry.Symbol)
		if _, ok := aliases.typeAliases[receiver]; ok {
			return "C"
		}
		if isFamilyIndexedMethodPackage(entry.Package) {
			return "A"
		}
		if isStandardInterfaceMethod(methodName(entry.Symbol)) {
			return "B"
		}
	}

	return "A"
}

func verifyPassThroughEvidence(entry auditEntry, aliases aliasSurface) []string {
	if entry.Kind == "top" {
		if _, ok := aliases.typeAliases[entry.Symbol]; ok {
			return nil
		}
		if _, ok := aliases.valueAliases[entry.Symbol]; ok {
			return nil
		}
		return []string{"tier C top-level entry lacks alias evidence: " + entry.ID}
	}
	if entry.Kind == "method" {
		if _, ok := aliases.typeAliases[receiverName(entry.Symbol)]; ok {
			return nil
		}
		return []string{"tier C method lacks external receiver alias evidence: " + entry.ID}
	}

	return []string{"tier C entry has unsupported kind: " + entry.ID}
}

func verifyReceiverDocumented(entry auditEntry, docs []corpus) []string {
	receiver := receiverName(entry.Symbol)
	if receiver == "" {
		return []string{"method entry missing receiver: " + entry.ID}
	}

	var failures []string
	for _, doc := range docs {
		if !strings.Contains(doc.content, "`"+receiver+"`") && !strings.Contains(doc.content, receiver) {
			failures = append(failures, fmt.Sprintf("%s missing receiver family for %s", doc.label, entry.Symbol))
		}
	}

	return failures
}

func verifyPolicyTerms(pkg string, docs []corpus, terms []string, label string) []string {
	if len(terms) == 0 {
		return []string{fmt.Sprintf("%s missing required %s terms", pkg, label)}
	}

	var failures []string
	for _, doc := range docs {
		for _, term := range terms {
			if !containsTerm(doc.content, term) {
				failures = append(failures, fmt.Sprintf("%s missing %s term: %s", doc.label, label, term))
			}
		}
	}

	return failures
}

func containsTerm(content, term string) bool {
	if strings.Contains(content, term) {
		return true
	}

	return strings.Contains(normalizeWhitespace(content), normalizeWhitespace(term))
}

func normalizeWhitespace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func documentMentionsEntry(content string, entry auditEntry) bool {
	symbol := entry.Symbol
	short := methodName(symbol)
	if entry.Kind == "field" {
		short = fieldName(symbol)
	}

	candidates := []string{
		"`" + symbol + "`",
		"`" + symbol + "(",
		symbol + "(",
		"`" + short + "`",
		"`" + short + "(",
		"." + short + "(",
	}
	if entry.Kind == "field" {
		if jsonField := jsonFieldFromSignature(entry.Signature); jsonField != "" && jsonField != "-" {
			candidates = append(candidates, "`"+jsonField+"`", `"`+jsonField+`"`)
		}
	}

	for _, candidate := range candidates {
		if strings.Contains(content, candidate) {
			return true
		}
	}

	return containsWord(content, short)
}

func indexMentionsSignature(content, pkg, signature string) bool {
	section := packageSection(content, pkg)
	return section != "" && strings.Contains(section, signature)
}

func packageSection(content, pkg string) string {
	marker := "## " + pkg
	start := strings.Index(content, marker)
	if start < 0 {
		return ""
	}

	rest := content[start:]
	next := strings.Index(rest[len(marker):], "\n## ")
	if next < 0 {
		return rest
	}

	return rest[:len(marker)+next]
}

func containsWord(content, word string) bool {
	if word == "" {
		return false
	}

	start := 0
	for {
		idx := strings.Index(content[start:], word)
		if idx < 0 {
			return false
		}
		pos := start + idx
		beforeOK := pos == 0 || !isWordRune(rune(content[pos-1]))
		after := pos + len(word)
		afterOK := after >= len(content) || !isWordRune(rune(content[after]))
		if beforeOK && afterOK {
			return true
		}
		start = after
	}
}

func isWordRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func jsonFieldFromSignature(signature string) string {
	marker := `json:\"`
	start := strings.Index(signature, marker)
	if start < 0 {
		return ""
	}

	value := signature[start+len(marker):]
	end := len(value)
	for _, sep := range []string{",", `\"`} {
		if idx := strings.Index(value, sep); idx >= 0 && idx < end {
			end = idx
		}
	}

	return value[:end]
}

func receiverName(symbol string) string {
	receiver, _, ok := strings.Cut(symbol, ".")
	if !ok {
		return ""
	}

	return receiver
}

func methodName(symbol string) string {
	_, method, ok := strings.Cut(symbol, ".")
	if !ok {
		return symbol
	}

	return method
}

func fieldName(symbol string) string {
	return methodName(symbol)
}

func isStandardInterfaceMethod(name string) bool {
	switch name {
	case "Error", "Unwrap", "MarshalJSON", "UnmarshalJSON", "MarshalText", "UnmarshalText", "Scan", "Value":
		return true
	default:
		return false
	}
}

func isFamilyIndexedMethodPackage(pkg string) bool {
	_, ok := familyIndexedMethodPackages[pkg]
	return ok
}

func packageCorpora(docsRoot string, entry manifestEntry) []corpus {
	var english []string
	var chinese []string
	var englishLabels []string
	var chineseLabels []string
	for _, path := range entry.Coverage {
		englishCorpus := readCorpus(path, filepath.Join(docsRoot, path))
		english = append(english, englishCorpus.content)
		englishLabels = append(englishLabels, path)
		if zhPath, ok := chineseMirrorPath(path); ok {
			chineseCorpus := readCorpus(zhPath, filepath.Join(docsRoot, zhPath))
			chinese = append(chinese, chineseCorpus.content)
			chineseLabels = append(chineseLabels, zhPath)
		}
	}

	result := []corpus{
		{label: strings.Join(englishLabels, " + "), content: strings.Join(english, "\n\n")},
	}
	if len(chinese) > 0 {
		result = append(result, corpus{label: strings.Join(chineseLabels, " + "), content: strings.Join(chinese, "\n\n")})
	}

	return result
}

func chineseMirrorPath(path string) (string, bool) {
	if !strings.HasPrefix(path, "docs/") {
		return "", false
	}

	return filepath.ToSlash(filepath.Join("i18n/zh-Hans/docusaurus-plugin-content-docs/current", strings.TrimPrefix(path, "docs/"))), true
}

func loadAliasSurfaces(sourceRoot string) map[string]aliasSurface {
	return map[string]aliasSurface{
		"github.com/coldsmirk/vef-framework-go/js":  parseExternalAliases(filepath.Join(sourceRoot, "js")),
		"github.com/coldsmirk/vef-framework-go/orm": parseExternalAliases(filepath.Join(sourceRoot, "orm")),
	}
}

func parseExternalAliases(dir string) aliasSurface {
	result := aliasSurface{
		typeAliases:  map[string]string{},
		valueAliases: map[string]string{},
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
	}, 0)
	if err != nil {
		panic(fmt.Errorf("failed to parse package aliases in %s: %w", dir, err))
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			imports := importAliases(file)
			for _, decl := range file.Decls {
				gen, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}
				for _, spec := range gen.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						if !s.Name.IsExported() || s.Assign == token.NoPos {
							continue
						}
						target, ok := externalSelectorTarget(s.Type, imports)
						if ok {
							result.typeAliases[s.Name.Name] = target
						}
					case *ast.ValueSpec:
						for i, name := range s.Names {
							if !name.IsExported() || i >= len(s.Values) {
								continue
							}
							target, ok := externalSelectorTarget(s.Values[i], imports)
							if ok {
								result.valueAliases[name.Name] = target
							}
						}
					}
				}
			}
		}
	}

	return result
}

func importAliases(file *ast.File) map[string]string {
	imports := map[string]string{}
	for _, spec := range file.Imports {
		path := strings.Trim(spec.Path.Value, `"`)
		name := filepath.Base(path)
		if spec.Name != nil {
			name = spec.Name.Name
		}
		imports[name] = path
	}

	return imports
}

func externalSelectorTarget(expr ast.Expr, imports map[string]string) (string, bool) {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return "", false
	}
	importPath := imports[ident.Name]
	if importPath == "" || strings.HasPrefix(importPath, sourceModule+"/") {
		return "", false
	}

	return importPath + "." + sel.Sel.Name, true
}

func auditEntriesByPackage(entries []auditEntry) map[string][]auditEntry {
	result := map[string][]auditEntry{}
	for _, entry := range entries {
		result[entry.Package] = append(result[entry.Package], entry)
	}

	for pkg := range result {
		sort.Slice(result[pkg], func(i, j int) bool {
			if result[pkg][i].Kind != result[pkg][j].Kind {
				return result[pkg][i].Kind < result[pkg][j].Kind
			}
			return result[pkg][i].Symbol < result[pkg][j].Symbol
		})
	}

	return result
}

func loadManifestByPackage(path string) map[string]manifestEntry {
	m := loadJSON[manifest](path)
	result := map[string]manifestEntry{}
	for _, entry := range m.Packages {
		result[entry.Package] = entry
	}

	return result
}

func readCorpus(label, path string) corpus {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read %s at %s: %w", label, path, err))
	}

	return corpus{label: label, content: string(content)}
}

func loadJSON[T any](path string) T {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		panic(fmt.Errorf("failed to parse %s: %w", path, err))
	}

	return result
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}
