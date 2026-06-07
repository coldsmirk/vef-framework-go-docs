package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

const (
	englishPath = "docs/reference/public-api-index.md"
	chinesePath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"
)

func main() {
	sourceDir := flag.String("source", ".", "path to the VEF Framework Go source repository")
	outDir := flag.String("out", "../vef-framework-go-docs", "path to the VEF Framework Go docs repository")
	flag.Parse()

	inventory, err := buildInventory(*sourceDir)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(filepath.Join(*outDir, englishPath), []byte(englishDocument(inventory)), 0o644); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(*outDir, chinesePath), []byte(chineseDocument(inventory)), 0o644); err != nil {
		panic(err)
	}
}

func buildInventory(sourceDir string) (string, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps,
		Dir:  sourceDir,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return "", err
	}

	sort.Slice(pkgs, func(i, j int) bool { return pkgs[i].PkgPath < pkgs[j].PkgPath })

	var buf strings.Builder
	for _, pkg := range pkgs {
		if strings.Contains(pkg.PkgPath, "/internal") || pkg.Name == "main" {
			continue
		}
		if len(pkg.Errors) > 0 {
			return "", fmt.Errorf("package errors in %s: %v", pkg.PkgPath, pkg.Errors)
		}

		fmt.Fprintf(&buf, "\n## %s\n", pkg.PkgPath)
		scope := pkg.Types.Scope()
		for _, name := range exportedNames(scope) {
			obj := scope.Lookup(name)
			if obj == nil {
				continue
			}
			writeObject(&buf, obj)
		}
	}

	return strings.TrimRight(buf.String(), "\n"), nil
}

func writeObject(buf *strings.Builder, obj types.Object) {
	switch obj := obj.(type) {
	case *types.TypeName:
		fmt.Fprintf(buf, "TYPE %s\n", typeString(obj))
		for _, field := range exportedFields(obj.Type()) {
			fmt.Fprintf(buf, "  FIELD %s\n", field)
		}
		for _, method := range exportedMethodSet(obj.Type()) {
			fmt.Fprintf(buf, "  METHOD %s\n", method)
		}
	case *types.Func:
		fmt.Fprintf(buf, "FUNC %s\n", typeString(obj))
	case *types.Const:
		fmt.Fprintf(buf, "CONST %s = %s\n", typeString(obj), obj.Val().ExactString())
	case *types.Var:
		fmt.Fprintf(buf, "VAR %s\n", typeString(obj))
	default:
		fmt.Fprintf(buf, "OTHER %s\n", typeString(obj))
	}
}

func exportedNames(scope *types.Scope) []string {
	names := make([]string, 0)
	for _, name := range scope.Names() {
		if tokenExported(name) {
			names = append(names, name)
		}
	}
	sort.Strings(names)

	return names
}

func tokenExported(name string) bool {
	if name == "" || name[0] == '_' {
		return false
	}
	r := rune(name[0])

	return 'A' <= r && r <= 'Z'
}

func typeString(obj types.Object) string {
	var buf bytes.Buffer
	buf.WriteString(obj.Name())
	if obj.Type() != nil {
		buf.WriteString(" : ")
		buf.WriteString(types.TypeString(obj.Type(), packagePath))
	}

	return buf.String()
}

func exportedMethodSet(t types.Type) []string {
	methods := map[string]bool{}
	for _, m := range exportedMethods(t) {
		methods[m] = true
	}
	for _, m := range exportedMethods(types.NewPointer(t)) {
		methods[m] = true
	}

	merged := make([]string, 0, len(methods))
	for m := range methods {
		merged = append(merged, m)
	}
	sort.Strings(merged)

	return merged
}

func exportedMethods(t types.Type) []string {
	set := types.NewMethodSet(t)
	names := make([]string, 0)
	for i := 0; i < set.Len(); i++ {
		sel := set.At(i)
		if tokenExported(sel.Obj().Name()) {
			names = append(names, fmt.Sprintf("%s : %s", sel.Obj().Name(), types.TypeString(sel.Obj().Type(), packagePath)))
		}
	}
	sort.Strings(names)

	return names
}

func exportedFields(t types.Type) []string {
	st := structType(t)
	if st == nil {
		return nil
	}

	fields := make([]string, 0)
	hidden := make(map[string]bool)
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		hidden[f.Name()] = true
		if tokenExported(f.Name()) {
			fields = append(fields, fieldSignature(f, i+1, st.Tag(i)))
		}
	}

	rootSeen := make(map[string]bool)
	if key := embeddedTypeKey(t); key != "" {
		rootSeen[key] = true
	}
	frontier := embeddedFields(st, "", rootSeen)
	for depth := 1; len(frontier) > 0; depth++ {
		level := make(map[string][]string)
		var next []embeddedField
		for _, embedded := range frontier {
			embeddedStruct := structType(embedded.Type)
			if embeddedStruct == nil {
				continue
			}
			for i := 0; i < embeddedStruct.NumFields(); i++ {
				f := embeddedStruct.Field(i)
				if tokenExported(f.Name()) && !hidden[f.Name()] {
					level[f.Name()] = append(level[f.Name()], promotedFieldSignature(f, embedded.Path, depth, i+1, embeddedStruct.Tag(i)))
				}
			}
			next = append(next, embeddedFields(embeddedStruct, embedded.Path, embedded.Seen)...)
		}

		names := make([]string, 0, len(level))
		for name := range level {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			hidden[name] = true
			if len(level[name]) == 1 {
				fields = append(fields, level[name][0])
			}
		}
		frontier = next
	}

	return fields
}

type embeddedField struct {
	Type types.Type
	Path string
	Seen map[string]bool
}

func fieldSignature(f *types.Var, order int, tag string) string {
	return fmt.Sprintf(
		"%s : %s [field_order=%d tag=%q]",
		f.Name(),
		types.TypeString(f.Type(), packagePath),
		order,
		tag,
	)
}

func promotedFieldSignature(f *types.Var, path string, depth, order int, tag string) string {
	return fmt.Sprintf(
		"%s : %s [promoted_from=%s depth=%d field_order=%d tag=%q]",
		f.Name(),
		types.TypeString(f.Type(), packagePath),
		path,
		depth,
		order,
		tag,
	)
}

func structType(t types.Type) *types.Struct {
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}
	if named, ok := t.(*types.Named); ok {
		t = named.Underlying()
	}
	st, ok := t.(*types.Struct)
	if !ok {
		return nil
	}

	return st
}

func embeddedFields(st *types.Struct, prefix string, seen map[string]bool) []embeddedField {
	fields := make([]embeddedField, 0)
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		if !f.Embedded() {
			continue
		}

		key := embeddedTypeKey(f.Type())
		if key == "" || seen[key] {
			continue
		}

		nextSeen := copySeenTypes(seen)
		nextSeen[key] = true
		path := f.Name()
		if prefix != "" {
			path = prefix + "." + path
		}
		fields = append(fields, embeddedField{
			Type: f.Type(),
			Path: path,
			Seen: nextSeen,
		})
	}

	return fields
}

func embeddedTypeKey(t types.Type) string {
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}
	if named, ok := t.(*types.Named); ok {
		obj := named.Obj()
		if obj.Pkg() != nil {
			return obj.Pkg().Path() + "." + obj.Name()
		}

		return obj.Name()
	}

	return types.TypeString(t, packagePath)
}

func copySeenTypes(seen map[string]bool) map[string]bool {
	result := make(map[string]bool, len(seen)+1)
	for key, value := range seen {
		result[key] = value
	}

	return result
}

func packagePath(p *types.Package) string {
	return p.Path()
}

func englishDocument(inventory string) string {
	return fmt.Sprintf(`---
sidebar_position: 90
---

# Public API Index

This page is an audit index generated from the current VEF Framework Go source tree. It lists exported symbols from non-%[1]sinternal/%[1]s packages, plus exported fields and methods on exported types, so documentation reviews have a concrete no-omissions checklist.

Topic guides remain the source for supported usage patterns and examples. Inclusion in this index is not a stability guarantee; treat the topic guides as the supported API contract.

Coverage rule used by the docs audit:

- every exported non-%[1]sinternal/%[1]s symbol is listed here, including exported fields and methods
- exported constant values, exported struct field order, exported struct field tags, and unambiguous promoted exported fields are part of the generated signatures
- every exported top-level symbol outside %[1]scmd/vef-cli/**%[1]s is mentioned in a topical guide or reference page
- every exported top-level symbol, field, and method has an entry in %[1]sscripts/api-audit-ledger.json%[1]s with a documented, grouped, or excluded audit disposition
- every package review in %[1]sscripts/api-contract-ledger.json%[1]s pins the exact top-level/field/method counts and fingerprint reviewed from the current source tree; drift is a failed audit
- user-facing runtime contracts that are not ordinary exported Go symbols are tracked separately in [Runtime API Index](./runtime-api-index)
- entries under %[1]scmd/vef-cli/**%[1]s are command implementation symbols kept for export-audit completeness; they are not a supported import API, and users should consume the CLI through the command surface documented in [CLI Tools](../advanced/cli-tools)

Regenerate the index and audit ledger whenever the framework source public surface changes:

%[2]sbash
(cd ../vef-framework-go && go run ../vef-framework-go-docs/scripts/gen-api-index.go -source . -out ../vef-framework-go-docs)
(cd ../vef-framework-go && go run ../vef-framework-go-docs/scripts/verify-api-audit.go -source . -manifest ../vef-framework-go-docs/scripts/api-audit-manifest.json -ledger ../vef-framework-go-docs/scripts/api-audit-ledger.json -write-ledger)
(cd ../vef-framework-go && go run ../vef-framework-go-docs/scripts/verify-api-audit.go -source . -manifest ../vef-framework-go-docs/scripts/api-audit-manifest.json -ledger ../vef-framework-go-docs/scripts/api-audit-ledger.json -contract-ledger ../vef-framework-go-docs/scripts/api-contract-ledger.json)
%[2]s

%[2]stext
%s
%[2]s
`, "`", "```", inventory)
}

func chineseDocument(inventory string) string {
	return fmt.Sprintf(`---
sidebar_position: 90
---

# 公开 API 索引

这一页是从当前 VEF Framework Go 源码生成的审计索引。它列出所有非 %[1]sinternal/%[1]s 包中的 exported symbols，以及 exported 类型上的公开字段和方法，用来给文档审查提供一个可核对的无遗漏清单。

各专题文档仍然负责说明推荐用法和示例。出现在这个索引中并不代表稳定性承诺；支持的 API contract 以专题文档为准。

本次文档审计采用的覆盖规则：

- 所有非 %[1]sinternal/%[1]s exported symbol 都列在本页，包括 exported 字段和方法
- exported 常量值、exported struct 字段顺序、exported struct 字段 tag，以及无歧义的 promoted exported fields 都会进入生成签名
- 除 %[1]scmd/vef-cli/**%[1]s 外，所有 exported top-level symbol 都必须在专题或参考文档中出现
- 每一个 exported top-level symbol、field 和 method 都必须在 %[1]sscripts/api-audit-ledger.json%[1]s 中有 documented、grouped 或 excluded 审计处置
- %[1]sscripts/api-contract-ledger.json%[1]s 中的每个 package review 都会固定当前源码已审查的 top-level/field/method 数量和 fingerprint；一旦漂移，审计失败
- 不属于普通 exported Go symbol 的用户可见运行时 contract 单独由 [Runtime API Index](./runtime-api-index) 跟踪
- %[1]scmd/vef-cli/**%[1]s 条目只为导出审计完整性保留；它们不是受支持的 import API，用户应按 [CLI 工具](../advanced/cli-tools) 记录的命令面使用

当框架源码的公开面发生变化时，请同步重新生成索引和审计账本：

%[2]sbash
(cd ../vef-framework-go && go run ../vef-framework-go-docs/scripts/gen-api-index.go -source . -out ../vef-framework-go-docs)
(cd ../vef-framework-go && go run ../vef-framework-go-docs/scripts/verify-api-audit.go -source . -manifest ../vef-framework-go-docs/scripts/api-audit-manifest.json -ledger ../vef-framework-go-docs/scripts/api-audit-ledger.json -write-ledger)
(cd ../vef-framework-go && go run ../vef-framework-go-docs/scripts/verify-api-audit.go -source . -manifest ../vef-framework-go-docs/scripts/api-audit-manifest.json -ledger ../vef-framework-go-docs/scripts/api-audit-ledger.json -contract-ledger ../vef-framework-go-docs/scripts/api-contract-ledger.json)
%[2]s

%[2]stext
%s
%[2]s
`, "`", "```", inventory)
}
