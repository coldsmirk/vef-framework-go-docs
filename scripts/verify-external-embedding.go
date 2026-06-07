package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

const sourceModule = "github.com/coldsmirk/vef-framework-go"

type auditLedger struct {
	Entries []auditEntry `json:"entries"`
}

type auditEntry struct {
	ID string `json:"id"`
}

type manifest struct {
	Packages []manifestEntry `json:"packages"`
}

type manifestEntry struct {
	Package  string   `json:"package"`
	Coverage []string `json:"coverage"`
}

type externalMethod struct {
	Package string
	Type    string
	Method  string
	Origin  string
}

type policy struct {
	Terms []string
}

var externalSurfacePolicies = map[string]policy{
	groupKey("github.com/coldsmirk/vef-framework-go", "go.uber.org/fx"): {
		Terms: []string{"go.uber.org/fx", "Lifecycle.Append", "public API index"},
	},
	groupKey("github.com/coldsmirk/vef-framework-go/ai", "io"): {
		Terms: []string{"io.Closer", "MessageStream", "StringStream"},
	},
	groupKey("github.com/coldsmirk/vef-framework-go/ai/stream", "io"): {
		Terms: []string{"io.Writer", "ResponseWriter"},
	},
	groupKey("github.com/coldsmirk/vef-framework-go/cache", "io"): {
		Terms: []string{"io.Closer", "Close()"},
	},
	groupKey("github.com/coldsmirk/vef-framework-go/decimal", "github.com/shopspring/decimal"): {
		Terms: []string{"shopspring/decimal.Decimal", "Decimal.*", "public API index"},
	},
	groupKey("github.com/coldsmirk/vef-framework-go/js", "github.com/dop251/goja"): {
		Terms: []string{"goja pass-through surface", "github.com/dop251/goja", "public API index"},
	},
	groupKey("github.com/coldsmirk/vef-framework-go/js", "github.com/dop251/goja/ast"): {
		Terms: []string{"goja pass-through surface", "js.AstProgram", "public API index"},
	},
	groupKey("github.com/coldsmirk/vef-framework-go/mcp", "github.com/modelcontextprotocol/go-sdk/mcp"): {
		Terms: []string{"MCP SDK pass-through surface", "github.com/modelcontextprotocol/go-sdk/mcp", "public API index"},
	},
	groupKey("github.com/coldsmirk/vef-framework-go/orm", "fmt"): {
		Terms: []string{"fmt.Stringer", "String()"},
	},
	groupKey("github.com/coldsmirk/vef-framework-go/orm", "github.com/uptrace/bun"): {
		Terms: []string{"Bun pass-through surface", "github.com/uptrace/bun", "public API index"},
	},
	groupKey("github.com/coldsmirk/vef-framework-go/orm", "github.com/uptrace/bun/schema"): {
		Terms: []string{"Bun pass-through surface", "Bun/schema aliases", "public API index"},
	},
}

func main() {
	sourceDir := flag.String("source", "../vef-framework-go", "path to the VEF Framework Go source repository")
	outDir := flag.String("out", ".", "path to the VEF Framework Go docs repository")
	flag.Parse()

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)

	methods := externalMethods(sourceRoot)
	ledgerIDs := loadLedgerIDs(filepath.Join(docsRoot, "scripts/api-audit-ledger.json"))
	manifestByPackage := loadManifest(filepath.Join(docsRoot, "scripts/api-audit-manifest.json"))

	var failures []string
	groups := map[string]int{}
	for _, method := range methods {
		groups[groupKey(method.Package, method.Origin)]++
		id := fmt.Sprintf("%s#method:%s.%s", method.Package, method.Type, method.Method)
		if !ledgerIDs[id] {
			failures = append(failures, "external promoted method missing from audit ledger: "+id)
		}
	}

	for key, count := range groups {
		policy, ok := externalSurfacePolicies[key]
		if !ok {
			failures = append(failures, fmt.Sprintf("external method group lacks policy: %s (%d methods)", key, count))
			continue
		}

		pkg, _, _ := strings.Cut(key, " <- ")
		entry, ok := manifestByPackage[pkg]
		if !ok {
			failures = append(failures, "external method package missing from manifest: "+pkg)
			continue
		}
		failures = append(failures, verifyTerms(docsRoot, entry, policy.Terms)...)
	}

	if len(failures) > 0 {
		sort.Strings(failures)
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Printf(
		"External pass-through method contracts verified: %d methods across %d package/origin groups\n",
		len(methods),
		len(groups),
	)
}

func externalMethods(sourceRoot string) []externalMethod {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedDeps,
		Dir:  sourceRoot,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		panic(err)
	}

	var methods []externalMethod
	for _, pkg := range pkgs {
		if strings.Contains(pkg.PkgPath, "/internal") || pkg.Name == "main" {
			continue
		}
		if len(pkg.Errors) > 0 {
			panic(fmt.Errorf("package errors in %s: %v", pkg.PkgPath, pkg.Errors))
		}

		scope := pkg.Types.Scope()
		for _, name := range scope.Names() {
			if !tokenExported(name) {
				continue
			}

			typeName, ok := scope.Lookup(name).(*types.TypeName)
			if !ok {
				continue
			}
			methods = append(methods, externalMethodSet(pkg.PkgPath, typeName)...)
		}
	}

	sort.Slice(methods, func(i, j int) bool {
		if methods[i].Package != methods[j].Package {
			return methods[i].Package < methods[j].Package
		}
		if methods[i].Type != methods[j].Type {
			return methods[i].Type < methods[j].Type
		}
		if methods[i].Method != methods[j].Method {
			return methods[i].Method < methods[j].Method
		}
		return methods[i].Origin < methods[j].Origin
	})

	return methods
}

func externalMethodSet(pkgPath string, typeName *types.TypeName) []externalMethod {
	seen := map[string]bool{}
	var methods []externalMethod
	for _, typ := range []types.Type{typeName.Type(), types.NewPointer(typeName.Type())} {
		set := types.NewMethodSet(typ)
		for i := 0; i < set.Len(); i++ {
			obj := set.At(i).Obj()
			if !tokenExported(obj.Name()) {
				continue
			}
			origin := packagePath(obj.Pkg())
			if !isExternalPackage(origin) {
				continue
			}

			key := obj.Name() + "\x00" + origin
			if seen[key] {
				continue
			}
			seen[key] = true
			methods = append(methods, externalMethod{
				Package: pkgPath,
				Type:    typeName.Name(),
				Method:  obj.Name(),
				Origin:  origin,
			})
		}
	}

	return methods
}

func verifyTerms(docsRoot string, entry manifestEntry, terms []string) []string {
	var failures []string
	for _, corpus := range []struct {
		label   string
		content string
	}{
		{label: entry.Package + " English coverage", content: readCoverage(docsRoot, entry.Coverage, false)},
		{label: entry.Package + " Chinese coverage", content: readCoverage(docsRoot, entry.Coverage, true)},
	} {
		for _, term := range terms {
			if !containsTerm(corpus.content, term) {
				failures = append(failures, fmt.Sprintf("%s missing external surface policy term: %s", corpus.label, term))
			}
		}
	}

	return failures
}

func readCoverage(docsRoot string, paths []string, chinese bool) string {
	var parts []string
	for _, path := range paths {
		if strings.Contains(path, "://") {
			continue
		}
		if chinese {
			zhPath, ok := chineseMirrorPath(path)
			if !ok {
				continue
			}
			path = zhPath
		}
		parts = append(parts, readFile(filepath.Join(docsRoot, path)))
	}

	return strings.Join(parts, "\n\n")
}

func loadLedgerIDs(path string) map[string]bool {
	ledger := loadJSON[auditLedger](path)
	result := map[string]bool{}
	for _, entry := range ledger.Entries {
		result[entry.ID] = true
	}

	return result
}

func loadManifest(path string) map[string]manifestEntry {
	manifest := loadJSON[manifest](path)
	result := map[string]manifestEntry{}
	for _, entry := range manifest.Packages {
		result[entry.Package] = entry
	}

	return result
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

func isExternalPackage(pkgPath string) bool {
	return pkgPath != "" && pkgPath != sourceModule && !strings.HasPrefix(pkgPath, sourceModule+"/")
}

func packagePath(pkg *types.Package) string {
	if pkg == nil {
		return ""
	}

	return pkg.Path()
}

func groupKey(pkg, origin string) string {
	return pkg + " <- " + origin
}

func chineseMirrorPath(path string) (string, bool) {
	if !strings.HasPrefix(path, "docs/") {
		return "", false
	}

	return filepath.ToSlash(filepath.Join("i18n/zh-Hans/docusaurus-plugin-content-docs/current", strings.TrimPrefix(path, "docs/"))), true
}

func loadJSON[T any](path string) T {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		panic(err)
	}

	return result
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(data)
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}

func tokenExported(name string) bool {
	if name == "" || name[0] == '_' {
		return false
	}
	r := rune(name[0])

	return 'A' <= r && r <= 'Z'
}
