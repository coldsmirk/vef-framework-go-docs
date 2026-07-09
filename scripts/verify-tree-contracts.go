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
	"regexp"
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

	englishDocs := readCorpus("English tree docs", filepath.Join(docsRoot, "docs/utilities/tree.md"))
	chineseDocs := readCorpus("Chinese tree docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/utilities/tree.md"))
	englishCrudDocs := readCorpus("English CRUD docs", filepath.Join(docsRoot, "docs/data-access/crud.md"))
	chineseCrudDocs := readCorpus("Chinese CRUD docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/data-access/crud.md"))
	publicIndex := readCorpus("English public API index", filepath.Join(docsRoot, "docs/reference/public-api-index.md"))
	chinesePublicIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"))

	treeDir := filepath.Join(sourceRoot, "tree")
	expectedSurface := packageSurface{
		funcs:  []string{"Build", "FindNode", "FindNodePath"},
		types:  []string{"Adapter"},
		fields: []string{"Adapter.GetChildren", "Adapter.GetID", "Adapter.GetParentID", "Adapter.SetChildren"},
	}

	var failures []string
	exported := exportedPackageSurface(treeDir)
	failures = append(failures, compareNames("tree const", exported.consts, expectedSurface.consts)...)
	failures = append(failures, compareNames("tree func", exported.funcs, expectedSurface.funcs)...)
	failures = append(failures, compareNames("tree type", exported.types, expectedSurface.types)...)
	failures = append(failures, compareNames("tree var", exported.vars, expectedSurface.vars)...)
	failures = append(failures, compareNames("tree method", exported.methods, expectedSurface.methods)...)
	failures = append(failures, compareNames("tree exported field", exported.fields, expectedSurface.fields)...)

	for _, doc := range []corpus{publicIndex, chinesePublicIndex} {
		failures = append(failures, missingTerms(doc, publicIndexTerms(exported))...)
	}
	for _, doc := range []corpus{englishDocs, chineseDocs} {
		failures = append(failures, missingTerms(doc, publicDocSurfaceTerms(exported))...)
	}

	failures = append(failures, missingTerms(englishDocs, []string{
		"`Build(nil, adapter)` and `Build([]T{}, adapter)` return a non-nil empty\n  slice (`[]T{}`)",
		"`GetID` values are raw string keys",
		"empty-ID nodes are not\n  indexed for parent lookup and their own children are not populated",
		"empty-ID nodes can still appear in the returned roots or in a parent's\n  children",
		"`GetParentID(node) == nil` makes the node a root",
		"parent ID that does not exist in the indexed node map also makes the\n  node a root",
		"closed cycles whose parent chain never reaches a root are omitted",
		"`Build` uses visited tracking while assigning children",
		"`Build` calls `SetChildren` on elements of the input slice and returns value\n  copies",
		"`GetChildren` is not called by `Build`",
		"missing adapter callbacks panic naturally",
		"empty `targetID` returns the zero value of `T` and `false`",
		"a missing target also returns the zero value of `T` and `false`",
		"duplicate IDs are not de-duplicated; the first traversal match wins",
		"`FindNode` does not add cycle protection around `GetChildren`",
		"an empty `targetID`, a missing target, or an empty tree returns `nil, false`",
		"a found target returns the full root-to-node path and `true`",
		"`FindNodePath` does not add cycle protection around `GetChildren`",
	})...)
	failures = append(failures, missingTerms(chineseDocs, []string{
		"`Build(nil, adapter)` 和 `Build([]T{}, adapter)` 返回非 nil 空切片",
		"`GetID` 的值按原始 string key 使用",
		"空 ID 节点不会进入 parent lookup 的索引",
		"空 ID 节点仍可能因为自身 parent 关系出现在返回的 roots 中",
		"`GetParentID(node) == nil` 时，该节点是 root",
		"parent ID 如果不存在于已索引 node map，也会让该节点成为 root",
		"形成 parent chain 永远到不了 root 的闭环",
		"`Build` 在设置 children 时使用 visited tracking",
		"`Build` 会对输入 slice 的元素调用 `SetChildren`",
		"`Build` 不会调用 `GetChildren`",
		"adapter callback 缺失时",
		"空 `targetID` 返回 `T` 的零值和 `false`",
		"目标不存在时也返回 `T` 的零值和 `false`",
		"duplicate IDs 不会被去重",
		"`FindNode` 不会围绕 `GetChildren` 额外加 cycle protection",
		"空 `targetID`、目标不存在或空树都会返回 `nil, false`",
		"命中目标时返回完整 root-to-node path 和 `true`",
		"`FindNodePath` 不会围绕 `GetChildren` 额外加 cycle protection",
	})...)

	for _, doc := range []corpus{englishCrudDocs, chineseCrudDocs} {
		failures = append(failures, missingTerms(doc, []string{
			"func buildCategoryTree(flat []Category) []Category",
			"adapter := tree.Adapter[Category]{",
			"return tree.Build(flat, adapter)",
			"crud.NewFindTree[Category, CategorySearch](buildCategoryTree)",
		})...)
		failures = append(failures, forbiddenDirectTreeBuild(doc)...)
	}

	sourceChecks := []struct {
		path  string
		terms []string
	}{
		{
			path: "tree/builder.go",
			terms: []string{
				"type Adapter[T any] struct",
				"GetID func(T) string",
				"GetParentID func(T) *string",
				"GetChildren func(T) []T",
				"SetChildren func(*T, []T)",
				"func Build[T any](nodes []T, adapter Adapter[T]) []T",
				"if len(nodes) == 0",
				"return []T{}",
				"if id := adapter.GetID(*node); id != \"\"",
				"childrenMap[*parentID] = append(childrenMap[*parentID], node)",
				"visited := make(map[string]bool)",
				"if id == \"\" || visited[id]",
				"adapter.SetChildren(nodePtr, children)",
				"if parentID == nil || nodeMap[*parentID] == nil",
				"func FindNode[T any](roots []T, targetID string, adapter Adapter[T]) (T, bool)",
				"if targetID == \"\"",
				"return lo.Empty[T](), false",
				"func FindNodePath[T any](roots []T, targetID string, adapter Adapter[T]) ([]T, bool)",
				"return nil, false",
				"adapter.GetChildren(node)",
			},
		},
		{
			path: "tree/builder_test.go",
			terms: []string{
				"HandlesEmptySlice",
				"HandlesNodesWithEmptyIDs",
				"HandlesCircularReferencesGracefully",
				"HandlesPartialCircularReferences",
				"ReturnsFalseForNonExistentNode",
				"ReturnsFalseForEmptyTargetID",
				"ReturnsNilForNonExistentNode",
				"ReturnsNilForEmptyTargetID",
				"AdapterWithNilFunctionsPanics",
				"NodesWithSpecialCharactersInIDs",
				"FindsFirstOccurrenceWithDuplicateIDs",
			},
		},
		{
			path: "crud/crud.go",
			terms: []string{
				"func NewFindTree[TModel, TSearch any](",
				"treeBuilder func(flatModels []TModel) []TModel",
				"RPCActionFindTree",
				"RESTActionFindTree",
			},
		},
		{
			path: "crud/find_tree.go",
			terms: []string{
				"type FindTree[TModel, TSearch any] interface",
				"WithIDColumn(name string) FindTree[TModel, TSearch]",
				"WithParentIDColumn(name string) FindTree[TModel, TSearch]",
				"treeBuilder    func(flatModels []TModel) []TModel",
				"models := a.treeBuilder(flatModels)",
			},
		},
		{
			path: "internal/approval/resource/category.go",
			terms: []string{
				"func buildFlowCategoryTree(flatCategories []approval.FlowCategory) []approval.FlowCategory",
				"adapter := tree.Adapter[approval.FlowCategory]{",
				"GetID: func(c approval.FlowCategory) string",
				"GetParentID: func(c approval.FlowCategory) *string",
				"SetChildren: func(c *approval.FlowCategory, children []approval.FlowCategory)",
				"return tree.Build(flatCategories, adapter)",
				"crud.NewFindTree[approval.FlowCategory, CategorySearch](buildFlowCategoryTree)",
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
		panic(fmt.Errorf("tree contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	topLevelPublic := len(exported.consts) + len(exported.funcs) + len(exported.types) + len(exported.vars)
	fmt.Printf("Tree contract docs verified: %d top-level public symbols, %d public methods, %d public fields, %d source files, 2 doc mirrors\n",
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
		panic(fmt.Errorf("failed to parse tree package: %w", err))
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
	for _, name := range surface.types {
		switch name {
		case "Adapter":
			terms = append(terms, "`Adapter[T]`", "`type Adapter[T any] struct`")
		default:
			terms = append(terms, "`"+name+"`")
		}
	}
	for _, name := range surface.fields {
		switch name {
		case "Adapter.GetID":
			terms = append(terms, "`Adapter.GetID`", "`func(T) string`")
		case "Adapter.GetParentID":
			terms = append(terms, "`Adapter.GetParentID`", "`func(T) *string`")
		case "Adapter.GetChildren":
			terms = append(terms, "`Adapter.GetChildren`", "`func(T) []T`")
		case "Adapter.SetChildren":
			terms = append(terms, "`Adapter.SetChildren`", "`func(*T, []T)`")
		default:
			terms = append(terms, "`"+name+"`")
		}
	}
	for _, name := range surface.funcs {
		switch name {
		case "Build":
			terms = append(terms, "`Build`", "`tree.Build[T any](nodes []T, adapter tree.Adapter[T]) []T`")
		case "FindNode":
			terms = append(terms, "`FindNode`", "`tree.FindNode[T any](roots []T, targetID string, adapter tree.Adapter[T]) (T, bool)`")
		case "FindNodePath":
			terms = append(terms, "`FindNodePath`", "`tree.FindNodePath[T any](roots []T, targetID string, adapter tree.Adapter[T]) ([]T, bool)`")
		default:
			terms = append(terms, "`"+name+"`")
		}
	}

	return terms
}

func publicIndexTerms(surface packageSurface) []string {
	terms := []string{"## github.com/coldsmirk/vef-framework-go/tree"}
	for _, name := range surface.types {
		switch name {
		case "Adapter":
			terms = append(terms, "TYPE Adapter : github.com/coldsmirk/vef-framework-go/tree.Adapter[T any]")
		default:
			terms = append(terms, "TYPE "+name+" :")
		}
	}
	for _, name := range surface.fields {
		switch name {
		case "Adapter.GetID":
			terms = append(terms, "FIELD GetID : func(T) string")
		case "Adapter.GetParentID":
			terms = append(terms, "FIELD GetParentID : func(T) *string")
		case "Adapter.GetChildren":
			terms = append(terms, "FIELD GetChildren : func(T) []T")
		case "Adapter.SetChildren":
			terms = append(terms, "FIELD SetChildren : func(*T, []T)")
		default:
			terms = append(terms, "FIELD "+strings.TrimPrefix(name, "Adapter.")+" :")
		}
	}
	for _, name := range surface.funcs {
		switch name {
		case "Build":
			terms = append(terms, "FUNC Build : func[T any](nodes []T, adapter github.com/coldsmirk/vef-framework-go/tree.Adapter[T]) []T")
		case "FindNode":
			terms = append(terms, "FUNC FindNode : func[T any](roots []T, targetID string, adapter github.com/coldsmirk/vef-framework-go/tree.Adapter[T]) (T, bool)")
		case "FindNodePath":
			terms = append(terms, "FUNC FindNodePath : func[T any](roots []T, targetID string, adapter github.com/coldsmirk/vef-framework-go/tree.Adapter[T]) ([]T, bool)")
		default:
			terms = append(terms, "FUNC "+name+" :")
		}
	}

	return terms
}

func runPackageTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./tree")
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("go test ./tree failed: %v\n%s", err, strings.TrimSpace(string(output)))}
	}

	return nil
}

func forbiddenDirectTreeBuild(c corpus) []string {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`NewFindTree\s*\[[^\]]+\]\s*\(\s*tree\.Build\s*\)`),
		regexp.MustCompile(`NewFindTree\s*\([^)]*tree\.Build`),
	}
	for _, pattern := range patterns {
		if pattern.MatchString(c.content) {
			return []string{fmt.Sprintf("%s contains forbidden direct tree.Build NewFindTree example", c.label)}
		}
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
