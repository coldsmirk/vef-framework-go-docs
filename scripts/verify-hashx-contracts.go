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

	englishDocs := readCorpus("English hashx docs", filepath.Join(docsRoot, "docs/utilities/hashx.md"))
	chineseDocs := readCorpus("Chinese hashx docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/utilities/hashx.md"))
	publicIndex := readCorpus("English public API index", filepath.Join(docsRoot, "docs/reference/public-api-index.md"))
	chinesePublicIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"))

	hashxDir := filepath.Join(sourceRoot, "hashx")
	expectedSurface := packageSurface{
		funcs: []string{
			"HmacMD5",
			"HmacSHA1",
			"HmacSHA256",
			"HmacSHA512",
			"HmacSM3",
			"MD5",
			"MD5Bytes",
			"SHA1",
			"SHA1Bytes",
			"SHA256",
			"SHA256Bytes",
			"SHA512",
			"SHA512Bytes",
			"SM3",
			"SM3Bytes",
		},
	}

	var failures []string
	exported := exportedPackageSurface(hashxDir)
	failures = append(failures, compareNames("hashx const", exported.consts, expectedSurface.consts)...)
	failures = append(failures, compareNames("hashx func", exported.funcs, expectedSurface.funcs)...)
	failures = append(failures, compareNames("hashx type", exported.types, expectedSurface.types)...)
	failures = append(failures, compareNames("hashx var", exported.vars, expectedSurface.vars)...)
	failures = append(failures, compareNames("hashx method", exported.methods, expectedSurface.methods)...)
	failures = append(failures, compareNames("hashx exported field", exported.fields, expectedSurface.fields)...)

	for _, doc := range []corpus{publicIndex, chinesePublicIndex} {
		failures = append(failures, missingTerms(doc, publicIndexTerms())...)
	}
	for _, doc := range []corpus{englishDocs, chineseDocs} {
		failures = append(failures, missingTerms(doc, publicDocSurfaceTerms())...)
	}

	failures = append(failures, missingTerms(englishDocs, []string{
		"15 top-level exported symbols",
		"no exported fields",
		"no exported\nmethods",
		"fingerprint is\n`22e7f661d37170d375f54592fa00078a3ea92b1b93459672709422aab54d5a01`",
		"Converts `data` to bytes, calls `MD5Bytes`",
		"Hashes raw bytes with `crypto/md5`",
		"Converts `data` to bytes, calls `SHA1Bytes`",
		"Hashes raw bytes with `crypto/sha1`",
		"Converts `data` to bytes, calls `SHA256Bytes`",
		"Hashes raw bytes with `crypto/sha256`",
		"Converts `data` to bytes, calls `SHA512Bytes`",
		"Hashes raw bytes with `crypto/sha512`",
		"Converts `data` to bytes, calls `SM3Bytes`",
		"Hashes raw bytes with `github.com/tjfoc/gmsm/sm3`",
		"nil slice is accepted and hashes like an empty\nslice",
		"lowercase hex-encoded MACs",
		"`MD5`, `MD5Bytes`, `HmacMD5`",
		"`SHA1`, `SHA1Bytes`, `HmacSHA1`",
		"`SHA256`, `SHA256Bytes`, `HmacSHA256`",
		"`SHA512`, `SHA512Bytes`, `HmacSHA512`",
		"`SM3`, `SM3Bytes`, `HmacSM3`",
		"`MD5` and `SHA1` are kept for compatibility checksums and legacy integrations",
		"Do not use them for password storage",
	})...)
	failures = append(failures, missingTerms(chineseDocs, []string{
		"15 个\ntop-level exported symbols",
		"没有 exported fields",
		"没有 exported methods",
		"fingerprint 是\n`22e7f661d37170d375f54592fa00078a3ea92b1b93459672709422aab54d5a01`",
		"将 `data` 转成 bytes，调用 `MD5Bytes`",
		"使用 `crypto/md5` 对原始 bytes 求 hash",
		"将 `data` 转成 bytes，调用 `SHA1Bytes`",
		"使用 `crypto/sha1` 对原始 bytes 求 hash",
		"将 `data` 转成 bytes，调用 `SHA256Bytes`",
		"使用 `crypto/sha256` 对原始 bytes 求 hash",
		"将 `data` 转成 bytes，调用 `SHA512Bytes`",
		"使用 `crypto/sha512` 对原始 bytes 求 hash",
		"将 `data` 转成 bytes，调用 `SM3Bytes`",
		"使用 `github.com/tjfoc/gmsm/sm3` 对原始 bytes 求 hash",
		"nil slice 也会被接受",
		"并按 empty slice 一样求 hash",
		"lowercase\nhex-encoded MAC",
		"`MD5`, `MD5Bytes`, `HmacMD5`",
		"`SHA1`, `SHA1Bytes`, `HmacSHA1`",
		"`SHA256`, `SHA256Bytes`, `HmacSHA256`",
		"`SHA512`, `SHA512Bytes`, `HmacSHA512`",
		"`SM3`, `SM3Bytes`, `HmacSM3`",
		"`MD5` 和 `SHA1` 只适合 compatibility checksum 和 legacy integration",
		"不要把它们用于密码存储",
	})...)

	sourceChecks := []struct {
		path  string
		terms []string
	}{
		{
			path: "hashx/hashx.go",
			terms: []string{
				"func MD5(data string) string",
				"return MD5Bytes([]byte(data))",
				"func MD5Bytes(data []byte) string",
				"sum := md5.Sum(data)",
				"return hex.EncodeToString(sum[:])",
				"func SHA1(data string) string",
				"return SHA1Bytes([]byte(data))",
				"func SHA1Bytes(data []byte) string",
				"sum := sha1.Sum(data)",
				"func SHA256(data string) string",
				"return SHA256Bytes([]byte(data))",
				"func SHA256Bytes(data []byte) string",
				"sum := sha256.Sum256(data)",
				"func SHA512(data string) string",
				"return SHA512Bytes([]byte(data))",
				"func SHA512Bytes(data []byte) string",
				"sum := sha512.Sum512(data)",
				"func SM3(data string) string",
				"return SM3Bytes([]byte(data))",
				"func SM3Bytes(data []byte) string",
				"sum := sm3.Sm3Sum(data)",
				"return hex.EncodeToString(sum)",
				"func HmacMD5(key, data []byte) string",
				"mac := hmac.New(md5.New, key)",
				"func HmacSHA1(key, data []byte) string",
				"mac := hmac.New(sha1.New, key)",
				"func HmacSHA256(key, data []byte) string",
				"mac := hmac.New(sha256.New, key)",
				"func HmacSHA512(key, data []byte) string",
				"mac := hmac.New(sha512.New, key)",
				"func HmacSM3(key, data []byte) string",
				"mac := hmac.New(sm3.New, key)",
				"_, _ = mac.Write(data)",
				"return hex.EncodeToString(mac.Sum(nil))",
			},
		},
		{
			path: "hashx/hashx_test.go",
			terms: []string{
				"TestMD5",
				"TestSHA1",
				"TestSHA256",
				"TestSHA512",
				"TestSM3",
				"TestHmacMD5",
				"TestHmacSHA1",
				"TestHmacSHA256",
				"TestHmacSHA512",
				"TestHmacSM3",
				"TestHashFunctionsNilInput",
				"TestHashOutputFormat",
				"Should contain only lowercase hex characters",
				"StandardTestVector",
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
		panic(fmt.Errorf("hashx contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	topLevelPublic := len(exported.consts) + len(exported.funcs) + len(exported.types) + len(exported.vars)
	fmt.Printf("Hashx contract docs verified: %d top-level public symbols, %d public methods, %d public fields, %d source files, 2 doc mirrors\n",
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
		panic(fmt.Errorf("failed to parse hashx package: %w", err))
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

func publicDocSurfaceTerms() []string {
	return []string{
		"`hashx.MD5(data string)`",
		"`hashx.MD5Bytes(data []byte)`",
		"`hashx.SHA1(data string)`",
		"`hashx.SHA1Bytes(data []byte)`",
		"`hashx.SHA256(data string)`",
		"`hashx.SHA256Bytes(data []byte)`",
		"`hashx.SHA512(data string)`",
		"`hashx.SHA512Bytes(data []byte)`",
		"`hashx.SM3(data string)`",
		"`hashx.SM3Bytes(data []byte)`",
		"`hashx.HmacMD5(key, data []byte)`",
		"`hashx.HmacSHA1(key, data []byte)`",
		"`hashx.HmacSHA256(key, data []byte)`",
		"`hashx.HmacSHA512(key, data []byte)`",
		"`hashx.HmacSM3(key, data []byte)`",
	}
}

func publicIndexTerms() []string {
	return []string{
		"## github.com/coldsmirk/vef-framework-go/hashx",
		"FUNC HmacMD5 : func(key []byte, data []byte) string",
		"FUNC HmacSHA1 : func(key []byte, data []byte) string",
		"FUNC HmacSHA256 : func(key []byte, data []byte) string",
		"FUNC HmacSHA512 : func(key []byte, data []byte) string",
		"FUNC HmacSM3 : func(key []byte, data []byte) string",
		"FUNC MD5 : func(data string) string",
		"FUNC MD5Bytes : func(data []byte) string",
		"FUNC SHA1 : func(data string) string",
		"FUNC SHA1Bytes : func(data []byte) string",
		"FUNC SHA256 : func(data string) string",
		"FUNC SHA256Bytes : func(data []byte) string",
		"FUNC SHA512 : func(data string) string",
		"FUNC SHA512Bytes : func(data []byte) string",
		"FUNC SM3 : func(data string) string",
		"FUNC SM3Bytes : func(data []byte) string",
	}
}

func runPackageTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./hashx")
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("go test ./hashx failed: %v\n%s", err, strings.TrimSpace(string(output)))}
	}

	return nil
}

func compareNames(label string, actual, expected []string) []string {
	actualSet := toSet(actual)
	expectedSet := toSet(expected)
	var failures []string

	for _, name := range expected {
		if !actualSet[name] {
			failures = append(failures, fmt.Sprintf("missing %s %s", label, name))
		}
	}
	for _, name := range actual {
		if !expectedSet[name] {
			failures = append(failures, fmt.Sprintf("unexpected %s %s", label, name))
		}
	}

	return failures
}

func toSet(values []string) map[string]bool {
	result := make(map[string]bool, len(values))
	for _, value := range values {
		result[value] = true
	}

	return result
}

func missingTerms(doc corpus, terms []string) []string {
	var failures []string
	for _, term := range terms {
		if !containsTerm(doc.content, term) {
			failures = append(failures, fmt.Sprintf("%s missing term %q", doc.label, term))
		}
	}

	return failures
}

func containsTerm(content string, term string) bool {
	if strings.Contains(content, term) {
		return true
	}

	return strings.Contains(normalizeSpace(content), normalizeSpace(term))
}

func normalizeSpace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func sortedKeys(values map[string]bool) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}

func receiverBaseName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return receiverBaseName(t.X)
	case *ast.IndexExpr:
		return receiverBaseName(t.X)
	case *ast.IndexListExpr:
		return receiverBaseName(t.X)
	default:
		return ""
	}
}

func readCorpus(label string, path string) corpus {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read %s at %s: %w", label, path, err))
	}

	return corpus{label: label, content: string(data)}
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(fmt.Errorf("failed to resolve %s: %w", path, err))
	}

	return filepath.Clean(abs)
}
