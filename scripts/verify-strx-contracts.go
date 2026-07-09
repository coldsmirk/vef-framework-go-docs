package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
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
	publicIndex := readCorpus("English public API index", filepath.Join(docsRoot, "docs/reference/public-api-index.md"))
	chinesePublicIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"))

	strxDir := filepath.Join(sourceRoot, "strx")
	expectedSurface := packageSurface{
		consts: []string{"BareAsKey", "BareAsValue", "DefaultKey"},
		funcs: []string{
			"ParseTag",
			"WithBareValueMode",
			"WithPairDelimiter",
			"WithPairDelimiterFunc",
			"WithSpacePairDelimiter",
			"WithValueDelimiter",
		},
		types: []string{"BareValueMode", "ParseOption"},
	}

	var failures []string
	exported := exportedPackageSurface(strxDir)
	failures = append(failures, compareNames("strx const", exported.consts, expectedSurface.consts)...)
	failures = append(failures, compareNames("strx func", exported.funcs, expectedSurface.funcs)...)
	failures = append(failures, compareNames("strx type", exported.types, expectedSurface.types)...)
	failures = append(failures, compareNames("strx var", exported.vars, expectedSurface.vars)...)
	failures = append(failures, compareNames("strx method", exported.methods, expectedSurface.methods)...)
	failures = append(failures, compareNames("strx exported field", exported.fields, expectedSurface.fields)...)

	for _, doc := range []corpus{publicIndex, chinesePublicIndex} {
		failures = append(failures, missingTerms(doc, publicIndexTerms(exported))...)
	}
	for _, doc := range []corpus{englishDocs, chineseDocs} {
		failures = append(failures, missingTerms(doc, publicDocSurfaceTerms(exported))...)
	}

	failures = append(failures, missingTerms(englishDocs, []string{
		"default parsing uses comma-separated pairs",
		"`=` as the key/value separator",
		"`BareAsValue`",
		"`ParseTag` always returns a non-nil map",
		"pair tokens are trimmed after splitting",
		"empty tokens are skipped",
		"split only at the first value delimiter",
		"explicit duplicate keys overwrite earlier values",
		"keys and values themselves are not trimmed",
		"empty keys and empty values are accepted",
		"special characters and Unicode are preserved as raw strings",
		"first bare token is stored under `DefaultKey`",
		"later\n  bare tokens are ignored and logged as warnings",
		"every bare token becomes a key with an empty string value",
		"duplicate bare keys collapse through normal map overwrite behavior",
		"`WithPairDelimiter(delimiter)` replaces the pair separator with a single-rune\n  equality check",
		"`WithPairDelimiterFunc(fn)` replaces the pair separator with `fn`",
		"`WithSpacePairDelimiter()` uses `unicode.IsSpace`",
		"spaces, tabs, and\n  newlines separate pairs",
		"`WithValueDelimiter(delimiter)` changes the key/value separator from `=`",
		"options are applied in order",
		"later options can override earlier separator or\n  bare-value settings",
		"nil `ParseOption` or nil delimiter\n  function will panic when reached",
	})...)
	failures = append(failures, missingTerms(chineseDocs, []string{
		"默认解析使用 comma-separated pairs",
		"`=` 作为 key/value separator",
		"`BareAsValue`",
		"`ParseTag` 总是返回非 nil map",
		"pair token 会在分隔后 trim",
		"空 token 会被跳过",
		"只在第一个 value delimiter 处分割",
		"显式重复 key 会以后出现的值覆盖先出现的值",
		"key 和 value 本身不会再 trim",
		"empty keys 和 empty values 都会被接受",
		"特殊字符和 Unicode 会按 raw strings 保留",
		"第一个 bare token 会写入 `DefaultKey`",
		"后续 bare\n  tokens 会被忽略并记录 warnings",
		"每个 bare token 都会成为 value 为空字符串的 key",
		"duplicate bare keys 按普通 map overwrite 行为折叠",
		"`WithPairDelimiter(delimiter)` 会把 pair separator 替换成单个 rune 的相等判断",
		"`WithPairDelimiterFunc(fn)` 会把 pair separator 替换成 `fn`",
		"`WithSpacePairDelimiter()` 使用 `unicode.IsSpace`",
		"spaces、tabs 和\n  newlines 都会分隔 pairs",
		"`WithValueDelimiter(delimiter)` 会把 key/value separator 从 `=` 改成指定 rune",
		"options 按顺序应用",
		"后面的 options 可以覆盖前面的 separator 或 bare-value\n  setting",
		"nil `ParseOption` 或 nil delimiter function\n  在执行到时会 panic",
	})...)

	sourceChecks := []struct {
		path  string
		terms []string
	}{
		{
			path: "strx/parser.go",
			terms: []string{
				"DefaultKey = \"__default\"",
				"type BareValueMode int",
				"BareAsValue BareValueMode = iota",
				"BareAsKey",
				"type ParseOption func(*parseConfig)",
				"pairSeparator  func(rune) bool",
				"valueSeparator rune",
				"bareValueMode  BareValueMode",
				"pairSeparator:  func(r rune) bool { return r == ',' }",
				"valueSeparator: '='",
				"bareValueMode:  BareAsValue",
				"func ParseTag(input string, opts ...ParseOption) map[string]string",
				"opt(cfg)",
				"result := make(map[string]string)",
				"strings.FieldsFuncSeq(input, cfg.pairSeparator)",
				"pair = strings.TrimSpace(pair)",
				"if pair == \"\"",
				"idx := strings.IndexRune(pair, cfg.valueSeparator)",
				"result[pair[:idx]] = pair[idx+1:]",
				"if cfg.bareValueMode == BareAsKey",
				"result[pair] = \"\"",
				"if _, exists := result[DefaultKey]; exists",
				"logger.Warnf(\"Ignoring duplicate default value %q in input: %s\", pair, input)",
				"result[DefaultKey] = pair",
				"func WithPairDelimiter(delimiter rune) ParseOption",
				"return r == delimiter",
				"func WithPairDelimiterFunc(fn func(rune) bool) ParseOption",
				"c.pairSeparator = fn",
				"func WithSpacePairDelimiter() ParseOption",
				"c.pairSeparator = unicode.IsSpace",
				"func WithValueDelimiter(delimiter rune) ParseOption",
				"c.valueSeparator = delimiter",
				"func WithBareValueMode(mode BareValueMode) ParseOption",
				"c.bareValueMode = mode",
			},
		},
		{
			path: "strx/parser_test.go",
			terms: []string{
				"TestParseTagCommaSeparated",
				"SingleAttributeWithoutKey",
				"DuplicateDefaultAttributes",
				"MultipleEqualsInValue",
				"UnicodeCharacters",
				"EmptyKeyName",
				"DuplicateKeys",
				"TestParseTagSpaceSeparated",
				"TabSeparated",
				"NewlineAsWhitespace",
				"TestParseTagCustomDelimiters",
				"WithPairDelimiterFunc",
				"TestParseTagBareValueMode",
				"BareAsValue_MultipleBareValues",
				"BareAsKey_MultipleBareValues",
				"BareAsKey_DuplicateBareValues",
				"TestParseTagEdgeCases",
				"ValueWithDelimiterCharacter",
				"KeyWithSpaces",
				"ValueWithSpaces",
				"OnlyEqualsSign",
				"TestParseTagOptionsOrdering",
				"OverridingOptions",
			},
		},
	}
	for _, check := range sourceChecks {
		source := readCorpus(check.path, filepath.Join(sourceRoot, check.path))
		failures = append(failures, missingTerms(source, check.terms)...)
	}

	failures = append(failures, runPackageTests(sourceRoot)...)

	sort.Strings(failures)
	if len(failures) > 0 {
		panic(fmt.Errorf("strx contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	topLevelPublic := len(exported.consts) + len(exported.funcs) + len(exported.types) + len(exported.vars)
	fmt.Printf("Strx contract docs verified: %d top-level public symbols, %d public methods, %d public fields, %d source files, 2 doc mirrors\n",
		topLevelPublic, len(exported.methods), len(exported.fields), len(sourceChecks))
}

type packageSurface struct {
	consts  []string
	funcs   []string
	types   []string
	vars    []string
	methods []string
	fields  []string
}

func exportedPackageSurface(dir string) packageSurface {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
	}, 0)
	if err != nil {
		panic(fmt.Errorf("failed to parse strx package: %w", err))
	}

	consts := make(map[string]bool)
	funcs := make(map[string]bool)
	typesMap := make(map[string]bool)
	vars := make(map[string]bool)
	methods := make(map[string]bool)
	fields := make(map[string]bool)

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				switch d := decl.(type) {
				case *ast.FuncDecl:
					if d.Recv == nil {
						if d.Name.IsExported() {
							funcs[d.Name.Name] = true
						}
						continue
					}
					if d.Name.IsExported() {
						recv := receiverBaseName(d.Recv.List[0].Type)
						if ast.IsExported(recv) {
							methods[recv+"."+d.Name.Name] = true
						}
					}
				case *ast.GenDecl:
					for _, spec := range d.Specs {
						switch s := spec.(type) {
						case *ast.TypeSpec:
							if s.Name.IsExported() {
								typesMap[s.Name.Name] = true
								if structType, ok := s.Type.(*ast.StructType); ok {
									for _, field := range structType.Fields.List {
										for _, name := range field.Names {
											if name.IsExported() {
												fields[s.Name.Name+"."+name.Name] = true
											}
										}
									}
								}
								if iface, ok := s.Type.(*ast.InterfaceType); ok {
									for _, field := range iface.Methods.List {
										for _, name := range field.Names {
											if name.IsExported() {
												methods[s.Name.Name+"."+name.Name] = true
											}
										}
									}
								}
							}
						case *ast.ValueSpec:
							for _, name := range s.Names {
								if !name.IsExported() {
									continue
								}
								if d.Tok == token.CONST {
									consts[name.Name] = true
								} else {
									vars[name.Name] = true
								}
							}
						}
					}
				}
			}
		}
	}

	return packageSurface{
		consts:  sortedKeys(consts),
		funcs:   sortedKeys(funcs),
		types:   sortedKeys(typesMap),
		vars:    sortedKeys(vars),
		methods: sortedKeys(methods),
		fields:  sortedKeys(fields),
	}
}

func publicDocSurfaceTerms(surface packageSurface) []string {
	var terms []string
	for _, name := range surface.consts {
		switch name {
		case "DefaultKey":
			terms = append(terms, "`DefaultKey`", "untyped string constant `\"__default\"`")
		case "BareAsValue":
			terms = append(terms, "`BareAsValue`", "`strx.BareValueMode = 0`")
		case "BareAsKey":
			terms = append(terms, "`BareAsKey`", "`strx.BareValueMode = 1`")
		default:
			terms = append(terms, "`"+name+"`")
		}
	}
	for _, name := range surface.types {
		switch name {
		case "BareValueMode":
			terms = append(terms, "`BareValueMode`", "`type BareValueMode int`")
		case "ParseOption":
			terms = append(terms, "`ParseOption`")
		default:
			terms = append(terms, "`"+name+"`")
		}
	}
	for _, name := range surface.funcs {
		switch name {
		case "ParseTag":
			terms = append(terms, "`ParseTag`", "`strx.ParseTag(input string, opts ...strx.ParseOption) map[string]string`")
		case "WithPairDelimiter":
			terms = append(terms, "`WithPairDelimiter`", "`strx.WithPairDelimiter(delimiter rune) strx.ParseOption`")
		case "WithPairDelimiterFunc":
			terms = append(terms, "`WithPairDelimiterFunc`", "`strx.WithPairDelimiterFunc(fn func(rune) bool) strx.ParseOption`")
		case "WithSpacePairDelimiter":
			terms = append(terms, "`WithSpacePairDelimiter`", "`strx.WithSpacePairDelimiter() strx.ParseOption`")
		case "WithValueDelimiter":
			terms = append(terms, "`WithValueDelimiter`", "`strx.WithValueDelimiter(delimiter rune) strx.ParseOption`")
		case "WithBareValueMode":
			terms = append(terms, "`WithBareValueMode`", "`strx.WithBareValueMode(mode strx.BareValueMode) strx.ParseOption`")
		default:
			terms = append(terms, "`"+name+"`")
		}
	}

	return terms
}

func publicIndexTerms(surface packageSurface) []string {
	terms := []string{"## github.com/coldsmirk/vef-framework-go/strx"}
	for _, name := range surface.consts {
		switch name {
		case "DefaultKey":
			terms = append(terms, "CONST DefaultKey : untyped string = \"__default\"")
		case "BareAsValue":
			terms = append(terms, "CONST BareAsValue : github.com/coldsmirk/vef-framework-go/strx.BareValueMode = 0")
		case "BareAsKey":
			terms = append(terms, "CONST BareAsKey : github.com/coldsmirk/vef-framework-go/strx.BareValueMode = 1")
		default:
			terms = append(terms, "CONST "+name+" :")
		}
	}
	for _, name := range surface.types {
		switch name {
		case "BareValueMode":
			terms = append(terms, "TYPE BareValueMode : github.com/coldsmirk/vef-framework-go/strx.BareValueMode")
		case "ParseOption":
			terms = append(terms, "TYPE ParseOption : github.com/coldsmirk/vef-framework-go/strx.ParseOption")
		default:
			terms = append(terms, "TYPE "+name+" :")
		}
	}
	for _, name := range surface.funcs {
		switch name {
		case "ParseTag":
			terms = append(terms, "FUNC ParseTag : func(input string, opts ...github.com/coldsmirk/vef-framework-go/strx.ParseOption) map[string]string")
		case "WithBareValueMode":
			terms = append(terms, "FUNC WithBareValueMode : func(mode github.com/coldsmirk/vef-framework-go/strx.BareValueMode) github.com/coldsmirk/vef-framework-go/strx.ParseOption")
		case "WithPairDelimiter":
			terms = append(terms, "FUNC WithPairDelimiter : func(delimiter rune) github.com/coldsmirk/vef-framework-go/strx.ParseOption")
		case "WithPairDelimiterFunc":
			terms = append(terms, "FUNC WithPairDelimiterFunc : func(fn func(rune) bool) github.com/coldsmirk/vef-framework-go/strx.ParseOption")
		case "WithSpacePairDelimiter":
			terms = append(terms, "FUNC WithSpacePairDelimiter : func() github.com/coldsmirk/vef-framework-go/strx.ParseOption")
		case "WithValueDelimiter":
			terms = append(terms, "FUNC WithValueDelimiter : func(delimiter rune) github.com/coldsmirk/vef-framework-go/strx.ParseOption")
		default:
			terms = append(terms, "FUNC "+name+" :")
		}
	}

	return terms
}

func runPackageTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./strx")
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("go test ./strx failed: %v\n%s", err, strings.TrimSpace(string(output)))}
	}

	return nil
}

func receiverBaseName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return receiverBaseName(t.X)
	case *ast.Ident:
		return t.Name
	case *ast.IndexExpr:
		return receiverBaseName(t.X)
	case *ast.IndexListExpr:
		return receiverBaseName(t.X)
	case *ast.SelectorExpr:
		return t.Sel.Name
	default:
		return ""
	}
}

func compareNames(label string, current, expected []string) []string {
	currentSet := set(current)
	expectedSet := set(expected)
	var failures []string
	for _, name := range current {
		if !expectedSet[name] {
			failures = append(failures, "unexpected public "+label+": "+name)
		}
	}
	for _, name := range expected {
		if !currentSet[name] {
			failures = append(failures, "missing expected public "+label+": "+name)
		}
	}

	return failures
}

func sortedKeys(values map[string]bool) []string {
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)

	return result
}

func set(values []string) map[string]bool {
	result := make(map[string]bool, len(values))
	for _, value := range values {
		result[value] = true
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
