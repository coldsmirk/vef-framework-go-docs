package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	englishRuntimePath = "docs/reference/runtime-api-index.md"
	chineseRuntimePath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/runtime-api-index.md"
	runtimeLedgerPath  = "scripts/runtime-api-ledger.json"
)

type RuntimeLedger struct {
	SourceModule string            `json:"source_module"`
	Scope        string            `json:"scope"`
	EntryCount   int               `json:"entry_count"`
	Fingerprint  string            `json:"fingerprint"`
	Coverage     []RuntimeCoverage `json:"coverage"`
	Entries      []RuntimeEntry    `json:"entries"`
}

type RuntimeCoverage struct {
	Category      string `json:"category"`
	EntryCount    int    `json:"entry_count"`
	Tier          string `json:"tier"`
	Extractor     string `json:"extractor"`
	Method        string `json:"method"`
	KnownResidual string `json:"known_residual"`
}

type RuntimeEntry struct {
	ID             string   `json:"id"`
	Category       string   `json:"category"`
	Name           string   `json:"name"`
	Value          string   `json:"value,omitempty"`
	Details        []string `json:"details,omitempty"`
	SourceEvidence []string `json:"source_evidence"`
	Terms          []string `json:"terms,omitempty"`
}

type sourceFile struct {
	Path string
	File *ast.File
	Fset *token.FileSet
}

type extractor struct {
	sourceDir string
	files     []sourceFile
	entries   map[string]RuntimeEntry

	constValues map[string]string
	constExprs  map[string]string
	configTypes map[string]configStruct
	configSeen  map[string]bool
}

type configStruct struct {
	Name   string
	Fields []configField
}

type configField struct {
	Name     string
	Key      string
	TypeName string
	TypeExpr string
	Source   string
}

type stringOccurrence struct {
	Value  string
	Source string
	Detail string
}

type resourceRef struct {
	Name   string
	Kind   string
	Source string
}

type actionMeta struct {
	Action             string
	RequiredPermission string
	Public             bool
	EnableAudit        bool
	Kind               string
	Source             string
}

func main() {
	sourceDir := flag.String("source", ".", "path to the VEF Framework Go source repository")
	outDir := flag.String("out", "../vef-framework-go-docs", "path to the VEF Framework Go docs repository")
	write := flag.Bool("write", false, "write generated runtime API index and ledger")
	flag.Parse()

	ledger, err := buildRuntimeLedger(*sourceDir)
	if err != nil {
		panic(err)
	}

	english := englishDocument(ledger)
	chinese := chineseDocument(ledger)
	ledgerJSON := mustJSON(ledger)

	if *write {
		if err := writeFile(filepath.Join(*outDir, englishRuntimePath), english); err != nil {
			panic(err)
		}
		if err := writeFile(filepath.Join(*outDir, chineseRuntimePath), chinese); err != nil {
			panic(err)
		}
		if err := writeFile(filepath.Join(*outDir, runtimeLedgerPath), ledgerJSON); err != nil {
			panic(err)
		}
		fmt.Printf("Runtime API audit index written: %d entries\n", ledger.EntryCount)

		return
	}

	checks := []struct {
		path string
		want string
	}{
		{filepath.Join(*outDir, englishRuntimePath), english},
		{filepath.Join(*outDir, chineseRuntimePath), chinese},
		{filepath.Join(*outDir, runtimeLedgerPath), ledgerJSON},
	}
	for _, check := range checks {
		got, err := os.ReadFile(check.path)
		if err != nil {
			panic(err)
		}
		if string(got) != check.want {
			panic(fmt.Errorf("%s is stale; rerun with -write", check.path))
		}
	}

	fmt.Printf("Runtime API audit verified: %d user-facing entries\n", ledger.EntryCount)
}

func buildRuntimeLedger(sourceDir string) (RuntimeLedger, error) {
	sourceDir, err := filepath.Abs(sourceDir)
	if err != nil {
		return RuntimeLedger{}, err
	}

	files, err := parseSource(sourceDir)
	if err != nil {
		return RuntimeLedger{}, err
	}

	x := &extractor{
		sourceDir:   sourceDir,
		files:       files,
		entries:     make(map[string]RuntimeEntry),
		constValues: make(map[string]string),
		constExprs:  make(map[string]string),
		configTypes: make(map[string]configStruct),
		configSeen:  make(map[string]bool),
	}

	x.collectConstantsAndConfigTypes()
	x.extractRuntimeEnumValues()
	x.extractProtocolConstants()
	x.extractAuthTypes()
	x.extractConfigKeys()
	x.extractConfigDefaults()
	x.extractBuiltInResources()
	x.extractCLI()
	x.extractJSONFields()
	x.extractMCP()
	x.extractValidatorRules()
	x.extractStructTagGrammars()
	x.extractMoldGrammar()
	x.extractJSONSchemaTags()
	x.extractI18NMessageKeys()
	x.extractEventTransportContracts()
	x.extractRESTVerbs()
	x.extractErrorCodes()
	if err := x.verifyRuntimeCoverageBoundaries(); err != nil {
		return RuntimeLedger{}, err
	}

	entries := make([]RuntimeEntry, 0, len(x.entries))
	for _, entry := range x.entries {
		sort.Strings(entry.Details)
		sort.Strings(entry.SourceEvidence)
		entry.Terms = buildTerms(entry)
		entries = append(entries, entry)
	}
	sortEntries(entries)

	ledger := RuntimeLedger{
		SourceModule: "github.com/coldsmirk/vef-framework-go",
		Scope:        "Runtime user-facing API surface: HTTP/RPC protocols, built-in resources, CLI, config, events, error codes, wire JSON fields, tag grammars, MCP surface, and runtime enum/string contracts. Test files, internal logs, and implementation-only strings are excluded.",
		EntryCount:   len(entries),
		Coverage:     buildRuntimeCoverage(entries),
		Entries:      entries,
	}
	ledger.Fingerprint = fingerprint(entries)

	return ledger, nil
}

func buildRuntimeCoverage(entries []RuntimeEntry) []RuntimeCoverage {
	descriptors := map[string]RuntimeCoverage{
		"API default": {
			Tier:          "Tier 3 curated source references",
			Extractor:     "extractProtocolConstants",
			Method:        "Curated defaults from API engine call sites and protocol constants.",
			KnownResidual: "None in generated index; semantic behavior remains covered by guide pages.",
		},
		"API version": {
			Tier:          "Tier 1 AST constants",
			Extractor:     "extractProtocolConstants",
			Method:        "AST scan of api/version.go VersionV* string constants.",
			KnownResidual: "None.",
		},
		"auth strategy": {
			Tier:          "Tier 1 AST constants",
			Extractor:     "extractProtocolConstants",
			Method:        "AST scan of api/auth.go AuthStrategy* string constants.",
			KnownResidual: "None.",
		},
		"auth type": {
			Tier:          "Tier 2 scoped AST constants",
			Extractor:     "extractAuthTypes",
			Method:        "AST scan of internal/security AuthType* constants that are sent through Authentication.Type.",
			KnownResidual: "None in known built-in authenticators.",
		},
		"built-in resource": {
			Tier:          "Tier 2 scoped AST resources",
			Extractor:     "extractBuiltInResources",
			Method:        "AST scan of NewRPCResource/NewRESTResource calls in built-in runtime resource packages.",
			KnownResidual: "None in scanned built-in resource directories.",
		},
		"built-in resource action": {
			Tier:          "Tier 2 scoped AST operations",
			Extractor:     "extractBuiltInResources",
			Method:        "AST scan of explicit OperationSpec values and CRUD builder defaults inside built-in runtime resource packages.",
			KnownResidual: "None in scanned built-in resource directories.",
		},
		"CLI command": {
			Tier:          "Tier 2 Cobra AST",
			Extractor:     "extractCLI",
			Method:        "AST scan of cobra.Command composites under cmd/vef-cli/cmd.",
			KnownResidual: "None in scanned CLI package.",
		},
		"CLI flag": {
			Tier:          "Tier 2 Cobra AST",
			Extractor:     "extractCLI",
			Method:        "AST scan of String/Bool/Int flag helper families and MarkFlagRequired calls under cmd/vef-cli/cmd; unsupported flag definition helpers fail boundary verification.",
			KnownResidual: "None for current Cobra flag definition calls.",
		},
		"config key": {
			Tier:          "Tier 2 config-tag AST",
			Extractor:     "extractConfigKeys",
			Method:        "AST walk of config structs rooted at known vef.* config roots plus vef.data_sources.<name>; verifier fails if a config/ struct with config tags is unreachable.",
			KnownResidual: "None for config/ structs with config tags.",
		},
		"config default": {
			Tier:          "Tier 3 mixed static extraction",
			Extractor:     "extractConfigDefaults",
			Method:        "AST extraction of Effective* accessors, ApplyDefaults assignments, monitor DefaultConfig values, and curated source references for defaults outside those named surfaces; boundary verification fails when a supported default surface is not indexed.",
			KnownResidual: "Defaults outside Effective*/ApplyDefaults/DefaultConfig and curated reviewed call sites require explicit review.",
		},
		"config enum": {
			Tier:          "Tier 2 scoped AST constants",
			Extractor:     "extractProtocolConstants",
			Method:        "AST scan of storage and datasource enum constants used in configuration values.",
			KnownResidual: "None in current config enum files.",
		},
		"config reserved name": {
			Tier:          "Tier 1 AST constants",
			Extractor:     "extractProtocolConstants",
			Method:        "AST scan of reserved configuration-name constants.",
			KnownResidual: "None.",
		},
		"environment variable": {
			Tier:          "Tier 1 AST constants",
			Extractor:     "extractProtocolConstants",
			Method:        "AST scan of Env* constants plus boundary checks for os.Getenv/os.LookupEnv call sites.",
			KnownResidual: "None for string-literal or const-backed environment lookups.",
		},
		"event transport contract": {
			Tier:          "Tier 2 scoped AST constants",
			Extractor:     "extractEventTransportContracts",
			Method:        "AST/source-derived extraction of outbox DLQ headers, topic prefix, retry backoff, and persisted-error bounds.",
			KnownResidual: "None for current built-in event transports.",
		},
		"HTTP endpoint": {
			Tier:          "Tier 2 source-derived constants",
			Extractor:     "extractProtocolConstants",
			Method:        "Source-derived REST/RPC/MCP endpoint constants and call-site evidence.",
			KnownResidual: "None for framework-owned default endpoints.",
		},
		"HTTP header": {
			Tier:          "Tier 1 AST constants",
			Extractor:     "extractProtocolConstants",
			Method:        "AST scan of api/header.go Header* constants.",
			KnownResidual: "None.",
		},
		"HTTP wire field": {
			Tier:          "Tier 3 curated protocol fields",
			Extractor:     "extractProtocolConstants",
			Method:        "Curated source references for fundamental request/result fields shared by REST/RPC.",
			KnownResidual: "None in generated index; JSON DTO fields are covered separately.",
		},
		"i18n message key": {
			Tier:          "Tier 2 AST call/tag scan",
			Extractor:     "extractI18NMessageKeys",
			Method:        "AST scan of literal or const-backed i18n.T calls, validator rule message keys, and label_i18n struct tags.",
			KnownResidual: "None for literal or const-backed keys; dynamic sources are tracked as i18n key indirections.",
		},
		"i18n key indirection": {
			Tier:          "Tier 2 AST call scan",
			Extractor:     "extractI18NMessageKeys",
			Method:        "AST scan of dynamic i18n.T call sites whose key source is another audited surface such as label_i18n tags, validator rules, or Fiber error mappings.",
			KnownResidual: "None for current dynamic i18n.T call sites.",
		},
		"JSON wire field": {
			Tier:          "Tier 2 scoped DTO AST with closed-world boundary check",
			Extractor:     "extractJSONFields",
			Method:        "AST scan of json tags on runtime DTO structs plus a boundary check over every non-test json-tagged struct field.",
			KnownResidual: "None for current non-test source; new json-tagged runtime fields must be indexed or explicitly excluded.",
		},
		"MCP endpoint": {
			Tier:          "Tier 1 AST constants",
			Extractor:     "extractProtocolConstants",
			Method:        "AST scan of the MCP Streamable HTTP endpoint constant.",
			KnownResidual: "None.",
		},
		"MCP jsonschema tag": {
			Tier:          "Tier 2 pinned dependency parser catalog",
			Extractor:     "extractJSONSchemaTags",
			Method:        "Catalog of struct-tag keywords accepted by github.com/invopop/jsonschema v0.14.0, with boundary verification that fails on dependency-version drift and uncovered in-source jsonschema tags.",
			KnownResidual: "None for the pinned jsonschema parser version.",
		},
		"MCP prompt": {
			Tier:          "Tier 2 MCP AST",
			Extractor:     "extractMCP",
			Method:        "AST scan of internal/mcp Prompt composites.",
			KnownResidual: "None in scanned MCP package.",
		},
		"MCP resource": {
			Tier:          "Tier 2 MCP AST",
			Extractor:     "extractMCP",
			Method:        "AST scan of internal/mcp Resource composites.",
			KnownResidual: "None in scanned MCP package.",
		},
		"MCP resource template": {
			Tier:          "Tier 2 MCP AST",
			Extractor:     "extractMCP",
			Method:        "AST scan of internal/mcp ResourceTemplate composites.",
			KnownResidual: "None in scanned MCP package.",
		},
		"MCP tool": {
			Tier:          "Tier 2 MCP AST",
			Extractor:     "extractMCP",
			Method:        "AST scan of internal/mcp Tool composites.",
			KnownResidual: "None in scanned MCP package.",
		},
		"CRUD REST action": {
			Tier:          "Tier 1 AST constants",
			Extractor:     "extractProtocolConstants",
			Method:        "AST scan of CRUD REST action constants.",
			KnownResidual: "None.",
		},
		"CRUD RPC action": {
			Tier:          "Tier 1 AST constants",
			Extractor:     "extractProtocolConstants",
			Method:        "AST scan of CRUD RPC action constants.",
			KnownResidual: "None.",
		},
		"REST action verb": {
			Tier:          "Tier 2 validator AST",
			Extractor:     "extractRESTVerbs",
			Method:        "AST scan of the REST action validator's allowed HTTP verb set.",
			KnownResidual: "None in current validator construction.",
		},
		"RPC form key": {
			Tier:          "Tier 1 AST constants",
			Extractor:     "extractProtocolConstants",
			Method:        "AST scan of FormKey* constants.",
			KnownResidual: "None.",
		},
		"event topic": {
			Tier:          "Tier 2 event constant/method scan",
			Extractor:     "extractProtocolConstants, extractMoldGrammar",
			Method:        "AST scan of EventType*/eventType* constants, EventType() return values, and built-in subscription/route-inspection topic call sites.",
			KnownResidual: "None for framework-owned non-test event topics.",
		},
		"mold tag grammar": {
			Tier:          "Tier 2 parser grammar scan",
			Extractor:     "extractMoldGrammar",
			Method:        "AST scan of the default mold tag name and restricted parser token constants, with boundary verification for parser token coverage.",
			KnownResidual: "None for current mold parser token constants.",
		},
		"mold transformer tag": {
			Tier:          "Tier 2 transformer scan",
			Extractor:     "extractMoldGrammar",
			Method:        "AST scan of built-in FieldTransformer Tag() methods.",
			KnownResidual: "None for current built-in mold transformer Tag() methods.",
		},
		"mold translate kind prefix": {
			Tier:          "Tier 2 translator scan",
			Extractor:     "extractMoldGrammar",
			Method:        "AST scan of built-in Translator Supports(kind) prefix checks.",
			KnownResidual: "None for current built-in translate kind prefixes.",
		},
		"meta tag grammar": {
			Tier:          "Tier 2 AST constants",
			Extractor:     "extractStructTagGrammars",
			Method:        "Catalog of storage meta tag name, dive value, file-reference kinds, and attribute grammar delimiters.",
			KnownResidual: "None for the current parser constants and tag parsing rules.",
		},
		"result error code": {
			Tier:          "Tier 1 AST constants",
			Extractor:     "extractErrorCodes",
			Method:        "AST scan of ErrCode* constants in api_errors.go and result/constants.go.",
			KnownResidual: "None for named error-code constants.",
		},
		"result message key": {
			Tier:          "Tier 1 AST constants",
			Extractor:     "extractProtocolConstants",
			Method:        "AST scan of ErrMessage* constants.",
			KnownResidual: "Inline i18n keys are covered by the i18n message key category.",
		},
		"runtime enum value": {
			Tier:          "Tier 2 typed string constants",
			Extractor:     "extractRuntimeEnumValues",
			Method:        "AST scan of typed string constants in public packages plus runtime internal DTO/transport packages.",
			KnownResidual: "Integer/stringer enum renderings are covered by the generated public API index and package contract ledger.",
		},
		"search tag grammar": {
			Tier:          "Tier 1 AST constants",
			Extractor:     "extractStructTagGrammars",
			Method:        "AST scan of search tag name, attributes, params, ignore marker, and operator/type tokens.",
			KnownResidual: "None for constants in search/constants.go.",
		},
		"tabular tag grammar": {
			Tier:          "Tier 1 AST constants",
			Extractor:     "extractStructTagGrammars",
			Method:        "AST scan of tabular tag name, attributes, and ignore marker.",
			KnownResidual: "None for constants in tabular/constants.go.",
		},
		"validator label tag": {
			Tier:          "Tier 2 validator tag scan",
			Extractor:     "extractValidatorRules",
			Method:        "AST scan of validator struct-tag key constants used by Field.Tag.Get.",
			KnownResidual: "None for current validator label tag lookups.",
		},
		"validator tag": {
			Tier:          "Tier 2 validator AST",
			Extractor:     "extractValidatorRules",
			Method:        "AST scan of built-in validator rule registration calls.",
			KnownResidual: "None for current built-in validator registrations and ValidationRule composites.",
		},
	}

	counts := make(map[string]int)
	for _, entry := range entries {
		counts[entry.Category]++
	}

	categories := make([]string, 0, len(counts))
	for category := range counts {
		if _, ok := descriptors[category]; !ok {
			panic(fmt.Errorf("runtime coverage descriptor missing for category %q", category))
		}
		categories = append(categories, category)
	}
	sort.Strings(categories)

	coverage := make([]RuntimeCoverage, 0, len(categories))
	for _, category := range categories {
		item := descriptors[category]
		item.Category = category
		item.EntryCount = counts[category]
		coverage = append(coverage, item)
	}

	return coverage
}

func parseSource(sourceDir string) ([]sourceFile, error) {
	var files []sourceFile
	err := filepath.WalkDir(sourceDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		name := d.Name()
		if d.IsDir() {
			switch name {
			case ".git", "vendor", "node_modules":
				return filepath.SkipDir
			}

			return nil
		}

		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			return nil
		}

		rel, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parse %s: %w", rel, err)
		}

		files = append(files, sourceFile{Path: rel, File: file, Fset: fset})

		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })

	return files, nil
}

func (x *extractor) add(category, name, value string, source []string, details ...string) {
	if category == "" || name == "" {
		return
	}

	idParts := []string{category, name, value}
	idParts = append(idParts, source...)
	id := stableID(idParts...)
	entry := x.entries[id]
	entry.ID = id
	entry.Category = category
	entry.Name = name
	entry.Value = value
	entry.SourceEvidence = appendUnique(entry.SourceEvidence, source...)
	entry.Details = appendUnique(entry.Details, nonEmpty(details)...)
	x.entries[id] = entry
}

func (x *extractor) entryEvidence(category string) map[string]bool {
	result := make(map[string]bool)
	for _, entry := range x.entries {
		if entry.Category != category {
			continue
		}
		for _, evidence := range entry.SourceEvidence {
			result[evidence] = true
		}
	}

	return result
}

func (x *extractor) entryValues(category string) map[string]bool {
	result := make(map[string]bool)
	for _, entry := range x.entries {
		if entry.Category == category && entry.Value != "" {
			result[entry.Value] = true
		}
	}

	return result
}

func (x *extractor) evidence(sf sourceFile, node ast.Node) string {
	pos := sf.Fset.Position(node.Pos())
	if pos.Line == 0 {
		return sf.Path
	}

	return fmt.Sprintf("%s:%d", sf.Path, pos.Line)
}

func (x *extractor) extractConstStringValues(sf sourceFile, visit func(name, value string, evidence []string)) {
	for _, decl := range sf.File.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		var lastValues []ast.Expr
		for _, spec := range gen.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			if len(vs.Values) > 0 {
				lastValues = vs.Values
			}
			for i, name := range vs.Names {
				var expr ast.Expr
				if i < len(vs.Values) {
					expr = vs.Values[i]
				} else if i < len(lastValues) {
					expr = lastValues[i]
				} else if len(lastValues) == 1 {
					expr = lastValues[0]
				}
				if expr == nil {
					continue
				}
				value, ok := x.evalString(expr)
				if !ok {
					continue
				}
				visit(name.Name, value, []string{x.evidence(sf, name)})
			}
		}
	}
}

func (x *extractor) collectConstantsAndConfigTypes() {
	for _, sf := range x.files {
		for _, decl := range sf.File.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			switch gen.Tok {
			case token.CONST:
				x.collectConstSpec(sf, gen)
			case token.TYPE:
				x.collectTypeSpec(sf, gen)
			}
		}
	}
}

func (x *extractor) collectConstSpec(sf sourceFile, gen *ast.GenDecl) {
	var lastValues []ast.Expr
	for _, spec := range gen.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		if len(vs.Values) > 0 {
			lastValues = vs.Values
		}
		for i, name := range vs.Names {
			if name.Name == "_" {
				continue
			}
			var expr ast.Expr
			if i < len(vs.Values) {
				expr = vs.Values[i]
			} else if i < len(lastValues) {
				expr = lastValues[i]
			} else if len(lastValues) == 1 {
				expr = lastValues[0]
			}
			if expr == nil {
				continue
			}

			key := sf.Path + ":" + name.Name
			exprText := exprString(sf.Fset, expr)
			x.constExprs[key] = exprText
			x.constExprs[name.Name] = exprText
			if value, ok := x.evalString(expr); ok {
				x.constValues[key] = value
				x.constValues[name.Name] = value
				x.propagateStringAliases(name.Name, value)
			} else if value, ok := intValue(expr); ok {
				x.constValues[key] = value
				x.constValues[name.Name] = value
			}
		}
	}
}

func (x *extractor) propagateStringAliases(name, value string) {
	changed := true
	for changed {
		changed = false
		for key, expr := range x.constExprs {
			if _, exists := x.constValues[key]; exists {
				continue
			}
			if expr == name || strings.HasSuffix(expr, "."+name) {
				x.constValues[key] = value
				base := key
				if strings.Contains(base, ":") {
					base = strings.TrimPrefix(base[strings.LastIndex(base, ":"):], ":")
				}
				if _, exists := x.constValues[base]; !exists {
					x.constValues[base] = value
				}
				changed = true
			}
		}
	}
}

func (x *extractor) collectTypeSpec(sf sourceFile, gen *ast.GenDecl) {
	for _, spec := range gen.Specs {
		ts, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		st, ok := ts.Type.(*ast.StructType)
		if ok && strings.HasPrefix(sf.Path, "config/") {
			x.configTypes[ts.Name.Name] = configStruct{
				Name:   ts.Name.Name,
				Fields: configFields(sf, st),
			}
		}
	}
}

func configFields(sf sourceFile, st *ast.StructType) []configField {
	fields := make([]configField, 0)
	for _, field := range st.Fields.List {
		if field.Tag == nil {
			continue
		}
		tag, ok := structTag(field.Tag.Value, "config")
		if !ok || tag == "" || tag == "-" {
			continue
		}
		for _, name := range field.Names {
			fields = append(fields, configField{
				Name:     name.Name,
				Key:      tag,
				TypeName: namedType(field.Type),
				TypeExpr: exprString(sf.Fset, field.Type),
				Source:   fmt.Sprintf("%s:%d", sf.Path, sf.Fset.Position(field.Pos()).Line),
			})
		}
	}

	return fields
}

func (x *extractor) extractProtocolConstants() {
	for _, sf := range x.files {
		ast.Inspect(sf.File, func(node ast.Node) bool {
			vs, ok := node.(*ast.ValueSpec)
			if !ok {
				return true
			}
			for i, name := range vs.Names {
				if i >= len(vs.Values) {
					continue
				}
				value, ok := x.evalString(vs.Values[i])
				if !ok {
					continue
				}

				evidence := []string{x.evidence(sf, name)}
				switch {
				case strings.HasPrefix(name.Name, "VersionV") && sf.Path == "api/version.go":
					x.add("API version", name.Name, value, evidence)
				case strings.HasPrefix(name.Name, "AuthStrategy") && sf.Path == "api/auth.go":
					x.add("auth strategy", name.Name, value, evidence)
				case strings.HasPrefix(name.Name, "Header") && sf.Path == "api/header.go":
					x.add("HTTP header", name.Name, value, evidence)
				case name.Name == "DefaultRPCEndpoint":
					x.add("HTTP endpoint", "RPC endpoint", value, evidence, "POST endpoint for RPC requests")
				case name.Name == "DefaultRESTBasePath":
					x.add("HTTP endpoint", "REST base path", value, evidence)
				case strings.HasPrefix(name.Name, "FormKey"):
					x.add("RPC form key", name.Name, value, evidence)
				case name.Name == "mcpPath":
					x.add("MCP endpoint", "MCP Streamable HTTP endpoint", value, evidence, "all HTTP methods")
				case strings.HasPrefix(name.Name, "Env"):
					x.add("environment variable", name.Name, value, evidence)
				case name.Name == "PrimaryDataSourceName":
					x.add("config reserved name", name.Name, value, evidence, "used under vef.data_sources.<name>")
				case strings.HasPrefix(name.Name, "Storage") && strings.HasPrefix(sf.Path, "config/"):
					x.add("config enum", name.Name, value, evidence)
				case isDBKindConst(name.Name, sf.Path):
					x.add("config enum", name.Name, value, evidence)
				case strings.HasPrefix(name.Name, "RPCAction"):
					x.add("CRUD RPC action", name.Name, value, evidence)
				case strings.HasPrefix(name.Name, "RESTAction"):
					x.add("CRUD REST action", name.Name, value, evidence)
				case strings.HasPrefix(name.Name, "EventType") || (strings.HasPrefix(name.Name, "eventType") && strings.Contains(value, ".")):
					x.add("event topic", name.Name, value, evidence)
				case strings.HasPrefix(name.Name, "ErrMessage"):
					x.add("result message key", name.Name, value, evidence)
				}
			}

			return true
		})
	}

	x.add("API default", "default version", "v1", []string{"internal/api/engine.go:101", "api/version.go:4"})
	x.add("API default", "default timeout", "30s", []string{"internal/api/engine.go:100"})
	x.add("API default", "default auth strategy", "bearer", []string{"internal/api/engine.go:102", "api/auth.go:24"})
	x.add("API default", "default rate limit", "100 requests / 5m", []string{"internal/api/engine.go:104"})
	x.add("HTTP wire field", "api.Identifier.resource", "resource", []string{"api/request.go:14"}, `form:"resource"`)
	x.add("HTTP wire field", "api.Identifier.action", "action", []string{"api/request.go:15"}, `form:"action"`)
	x.add("HTTP wire field", "api.Identifier.version", "version", []string{"api/request.go:16"}, `form:"version"`)
	x.add("HTTP wire field", "api.Request.params", "params", []string{"api/request.go:58"})
	x.add("HTTP wire field", "api.Request.meta", "meta", []string{"api/request.go:59"})
	x.add("HTTP wire field", "result.Result.code", "code", []string{"result/result.go:11"})
	x.add("HTTP wire field", "result.Result.message", "message", []string{"result/result.go:12"})
	x.add("HTTP wire field", "result.Result.data", "data", []string{"result/result.go:13"})
}

func (x *extractor) extractAuthTypes() {
	for _, sf := range x.files {
		if !strings.HasPrefix(sf.Path, "internal/security/") {
			continue
		}
		ast.Inspect(sf.File, func(node ast.Node) bool {
			vs, ok := node.(*ast.ValueSpec)
			if !ok {
				return true
			}
			for i, name := range vs.Names {
				if !strings.HasPrefix(name.Name, "AuthType") || i >= len(vs.Values) {
					continue
				}
				value, ok := x.evalString(vs.Values[i])
				if ok {
					x.add("auth type", name.Name, value, []string{x.evidence(sf, name)})
				}
			}

			return true
		})
	}
}

func (x *extractor) extractErrorCodes() {
	for _, sf := range x.files {
		if !strings.HasSuffix(sf.Path, "api_errors.go") && sf.Path != "result/constants.go" {
			continue
		}
		ast.Inspect(sf.File, func(node ast.Node) bool {
			vs, ok := node.(*ast.ValueSpec)
			if !ok {
				return true
			}
			for i, name := range vs.Names {
				if !strings.HasPrefix(name.Name, "ErrCode") || i >= len(vs.Values) {
					continue
				}
				if value, ok := intValue(vs.Values[i]); ok {
					x.add("result error code", name.Name, value, []string{x.evidence(sf, name)})
				}
			}

			return true
		})
	}
}

func isDBKindConst(name, path string) bool {
	if path != "config/data_sources.go" {
		return false
	}
	switch name {
	case "Oracle", "SQLServer", "Postgres", "MySQL", "SQLite":
		return true
	default:
		return false
	}
}

func (x *extractor) extractRESTVerbs() {
	for _, sf := range x.files {
		if sf.Path != "api/resource.go" {
			continue
		}
		ast.Inspect(sf.File, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok || !strings.Contains(exprString(sf.Fset, call.Fun), "NewHashSetFrom") {
				return true
			}
			for _, arg := range call.Args {
				if value, ok := x.evalString(arg); ok && isHTTPVerb(value) {
					x.add("REST action verb", strings.ToUpper(value), value, []string{x.evidence(sf, arg)})
				}
			}

			return true
		})
	}
}

func (x *extractor) extractConfigKeys() {
	roots := map[string]string{
		"AppConfig":      "vef.app",
		"APIConfig":      "vef.api",
		"CORSConfig":     "vef.cors",
		"SecurityConfig": "vef.security",
		"RedisConfig":    "vef.redis",
		"StorageConfig":  "vef.storage",
		"MonitorConfig":  "vef.monitor",
		"MCPConfig":      "vef.mcp",
		"ApprovalConfig": "vef.approval",
		"EventConfig":    "vef.event",
	}
	for typeName, root := range roots {
		x.walkConfig(typeName, root, nil)
	}
	x.walkConfig("DataSourceConfig", "vef.data_sources.<name>", []string{"internal/config/data_sources.go:20", "config/data_sources.go:22"})
}

func (x *extractor) walkConfig(typeName, prefix string, inheritedEvidence []string) {
	cfg, ok := x.configTypes[typeName]
	if !ok {
		return
	}
	x.configSeen[typeName] = true
	for _, field := range cfg.Fields {
		key := prefix + "." + field.Key
		evidence := appendUnique([]string{field.Source}, inheritedEvidence...)
		x.add("config key", key, field.TypeExpr, evidence, "Go field: "+typeName+"."+field.Name)

		if _, ok := x.configTypes[field.TypeName]; ok {
			x.walkConfig(field.TypeName, key, evidence)
		}
	}
}

func (x *extractor) verifyRuntimeCoverageBoundaries() error {
	var failures []string
	for _, cfg := range x.configTypes {
		if x.configSeen[cfg.Name] || len(cfg.Fields) == 0 {
			continue
		}

		if strings.HasSuffix(cfg.Name, "Config") {
			failures = append(failures, fmt.Sprintf("config struct with config tags is unreachable from runtime config roots: %s", cfg.Name))
		}
	}

	failures = append(failures, x.verifyJSONFieldCoverage()...)
	failures = append(failures, x.verifyEnvLookupCoverage()...)
	failures = append(failures, x.verifyI18NKeyCoverage()...)
	failures = append(failures, x.verifyCLIFlagCoverage()...)
	failures = append(failures, x.verifyConfigDefaultCoverage()...)
	failures = append(failures, x.verifyEventTopicCoverage()...)
	failures = append(failures, x.verifyValidatorCoverage()...)
	failures = append(failures, x.verifyMoldCoverage()...)
	failures = append(failures, x.verifyJSONSchemaCoverage()...)

	if len(failures) > 0 {
		sort.Strings(failures)

		return fmt.Errorf("runtime API audit boundary verification failed:\n%s", strings.Join(failures, "\n"))
	}

	return nil
}

func (x *extractor) extractConfigDefaults() {
	fieldToKey := x.configFieldKeyIndex()

	for _, sf := range x.files {
		if !isConfigDefaultFile(sf.Path) {
			continue
		}

		ast.Inspect(sf.File, func(node ast.Node) bool {
			fn, ok := node.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				return true
			}

			if fn.Name.Name == "ApplyDefaults" {
				ast.Inspect(fn.Body, func(child ast.Node) bool {
					assign, ok := child.(*ast.AssignStmt)
					if !ok || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
						return true
					}
					field := selectorField(assign.Lhs[0])
					if field == "" {
						return true
					}
					if key := fieldToKey["ApprovalConfig."+field]; key != "" {
						x.add("config default", key, exprString(sf.Fset, assign.Rhs[0]), []string{x.evidence(sf, assign)})
					}

					return true
				})
			}

			if fn.Name.Name == "DefaultConfig" {
				x.extractDefaultConfigValues(sf, fn, fieldToKey)
			}

			if strings.HasPrefix(fn.Name.Name, "Effective") {
				key, value := x.effectiveDefault(sf, fn, fieldToKey)
				if key != "" && value != "" {
					x.add("config default", key, value, []string{x.evidence(sf, fn)})
				}
			}

			return true
		})
	}

	x.add("config default", "vef.event.transports.outbox.max_retries", "10", []string{"config/event.go:175"}, "EventConfig.Validate fallback when max_retries is unset")
	x.add("config default", "vef.mcp.require_auth", "true when unset", []string{"internal/mcp/handler.go:34"})
	x.add("environment variable", "SystemDrive", "SystemDrive", []string{"internal/monitor/service.go:189"}, "OS-defined Windows variable read to locate the root filesystem for the disk overview; not framework configuration")
}

func (x *extractor) extractDefaultConfigValues(sf sourceFile, fn *ast.FuncDecl, fieldToKey map[string]string) {
	if fn.Body == nil {
		return
	}
	ast.Inspect(fn.Body, func(node ast.Node) bool {
		ret, ok := node.(*ast.ReturnStmt)
		if !ok || len(ret.Results) != 1 {
			return true
		}
		lit, ok := ret.Results[0].(*ast.CompositeLit)
		if !ok {
			return true
		}
		typeName := namedType(lit.Type)
		if typeName == "" {
			return true
		}
		for _, elt := range lit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			field := keyName(kv.Key)
			key := fieldToKey[typeName+"."+field]
			if key == "" {
				continue
			}
			x.add("config default", key, x.resolveConstExpr(exprString(sf.Fset, kv.Value)), []string{x.evidence(sf, kv)})
		}

		return true
	})
}

func (x *extractor) defaultConfigEvidence(sf sourceFile, fn *ast.FuncDecl) []string {
	var evidence []string
	if fn.Body == nil {
		return evidence
	}
	ast.Inspect(fn.Body, func(node ast.Node) bool {
		ret, ok := node.(*ast.ReturnStmt)
		if !ok || len(ret.Results) != 1 {
			return true
		}
		lit, ok := ret.Results[0].(*ast.CompositeLit)
		if !ok {
			return true
		}
		for _, elt := range lit.Elts {
			if kv, ok := elt.(*ast.KeyValueExpr); ok {
				evidence = append(evidence, x.evidence(sf, kv))
			}
		}

		return true
	})

	return evidence
}

func (x *extractor) applyDefaultsEvidence(sf sourceFile, fn *ast.FuncDecl) []string {
	var evidence []string
	if fn.Body == nil {
		return evidence
	}
	ast.Inspect(fn.Body, func(node ast.Node) bool {
		assign, ok := node.(*ast.AssignStmt)
		if !ok || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
			return true
		}
		if selectorField(assign.Lhs[0]) != "" {
			evidence = append(evidence, x.evidence(sf, assign))
		}

		return true
	})

	return evidence
}

func (x *extractor) verifyJSONFieldCoverage() []string {
	covered := x.entryEvidence("JSON wire field")
	var failures []string
	for _, sf := range x.files {
		if strings.Contains(sf.Path, "/testdata/") {
			continue
		}
		for _, decl := range sf.File.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}
				for _, field := range st.Fields.List {
					if field.Tag == nil {
						continue
					}
					jsonName, ok := structTag(field.Tag.Value, "json")
					if !ok {
						continue
					}
					jsonName = strings.Split(jsonName, ",")[0]
					if jsonName == "" || jsonName == "-" {
						continue
					}

					evidence := x.evidence(sf, field)
					if !covered[evidence] && !isJSONFieldCoverageExcluded(sf.Path) {
						failures = append(failures, fmt.Sprintf("json-tagged field is not covered by Runtime API Index: %s", evidence))
					}
				}
			}
		}
	}

	return failures
}

func isJSONFieldCoverageExcluded(path string) bool {
	if path == "cmd/vef-cli/cmd/modelschema/testdata/models/sample.go" {
		return true
	}

	return false
}

func (x *extractor) verifyEnvLookupCoverage() []string {
	covered := x.entryValues("environment variable")
	var failures []string
	for _, sf := range x.files {
		ast.Inspect(sf.File, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			method := selectorMethod(call.Fun)
			if method != "Getenv" && method != "LookupEnv" {
				return true
			}
			if selectorReceiver(call.Fun) != "os" {
				return true
			}
			if len(call.Args) == 0 {
				return true
			}
			value, ok := x.evalString(call.Args[0])
			if !ok {
				failures = append(failures, fmt.Sprintf("environment lookup uses a dynamic key that Runtime API Index cannot verify: %s", x.evidence(sf, call)))

				return true
			}
			if !covered[value] {
				failures = append(failures, fmt.Sprintf("environment lookup key is not covered by Runtime API Index: %s = %q", x.evidence(sf, call), value))
			}

			return true
		})
	}

	return failures
}

func (x *extractor) verifyI18NKeyCoverage() []string {
	covered := x.entryValues("i18n message key")
	var failures []string
	for _, sf := range x.files {
		for _, occurrence := range x.i18nKeyOccurrences(sf) {
			if occurrence.Value == "<dynamic>" {
				continue
			}
			if !covered[occurrence.Value] {
				failures = append(failures, fmt.Sprintf("i18n message key is not covered by Runtime API Index: %s = %q", occurrence.Source, occurrence.Value))
			}
		}
	}

	return failures
}

func (x *extractor) verifyCLIFlagCoverage() []string {
	covered := x.entryEvidence("CLI flag")
	var failures []string
	for _, sf := range x.files {
		if !strings.HasPrefix(sf.Path, "cmd/vef-cli/cmd/") {
			continue
		}
		ast.Inspect(sf.File, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			fun := exprString(sf.Fset, call.Fun)
			if !isCLIFlagDefinitionCall(fun) {
				return true
			}
			if _, ok := cliFlagHelperKind(fun); !ok {
				if isCLIFlagReader(selectorMethod(call.Fun)) {
					return true
				}
				failures = append(failures, fmt.Sprintf("CLI flag definition helper is not supported by Runtime API Index extractor: %s uses %s", x.evidence(sf, call), selectorMethod(call.Fun)))

				return true
			}
			evidence := x.evidence(sf, call)
			if !covered[evidence] {
				failures = append(failures, fmt.Sprintf("CLI flag helper call is not covered by Runtime API Index: %s", evidence))
			}

			return true
		})
	}

	return failures
}

func (x *extractor) verifyConfigDefaultCoverage() []string {
	covered := x.entryEvidence("config default")
	var failures []string
	for _, sf := range x.files {
		if !isConfigDefaultFile(sf.Path) {
			continue
		}
		for _, decl := range sf.File.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			if strings.HasPrefix(fn.Name.Name, "Effective") && !covered[x.evidence(sf, fn)] {
				failures = append(failures, fmt.Sprintf("config Effective* default accessor is not covered by Runtime API Index: %s", x.evidence(sf, fn)))
			}
			if fn.Name.Name == "DefaultConfig" {
				for _, evidence := range x.defaultConfigEvidence(sf, fn) {
					if !covered[evidence] {
						failures = append(failures, fmt.Sprintf("DefaultConfig assignment is not covered by Runtime API Index: %s", evidence))
					}
				}
			}
			if fn.Name.Name == "ApplyDefaults" {
				for _, evidence := range x.applyDefaultsEvidence(sf, fn) {
					if !covered[evidence] {
						failures = append(failures, fmt.Sprintf("ApplyDefaults assignment is not covered by Runtime API Index: %s", evidence))
					}
				}
			}
		}
	}

	return failures
}

func (x *extractor) verifyEventTopicCoverage() []string {
	covered := x.entryValues("event topic")
	var failures []string
	for _, sf := range x.files {
		if strings.HasSuffix(sf.Path, "_test.go") || strings.Contains(sf.Path, "/testdata/") {
			continue
		}
		ast.Inspect(sf.File, func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.FuncDecl:
				if node.Name.Name != "EventType" {
					return true
				}
				value, ok := x.singleReturnString(sf, node)
				if !ok || value == "" {
					if valueExpr := singleReturnExpr(sf, node); valueExpr != "" && !isRawPayloadEventTypeMethod(sf.Path, node) {
						failures = append(failures, fmt.Sprintf("EventType method uses a dynamic topic that Runtime API Index cannot verify: %s", x.evidence(sf, node)))
					}

					return true
				}
				if !covered[value] {
					failures = append(failures, fmt.Sprintf("EventType method topic is not covered by Runtime API Index: %s = %q", x.evidence(sf, node), value))
				}
			case *ast.CallExpr:
				method := selectorMethod(node.Fun)
				if method != "Subscribe" && method != "HasTransactionalRoute" && method != "HasSubscribableTransport" {
					return true
				}
				if len(node.Args) == 0 {
					return true
				}
				value, ok := x.evalString(node.Args[0])
				if !ok || value == "" {
					return true
				}
				if isTestOnlyEventTopic(value) {
					return true
				}
				if !covered[value] {
					failures = append(failures, fmt.Sprintf("event topic call site is not covered by Runtime API Index: %s = %q", x.evidence(sf, node), value))
				}
			}

			return true
		})
	}

	return failures
}

func (x *extractor) verifyValidatorCoverage() []string {
	coveredTags := x.entryValues("validator tag")
	coveredLabels := x.entryValues("validator label tag")
	var failures []string
	for _, sf := range x.files {
		if !strings.HasPrefix(sf.Path, "validator/") {
			continue
		}
		ast.Inspect(sf.File, func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.CallExpr:
				fun := exprString(sf.Fset, node.Fun)
				switch {
				case strings.HasSuffix(fun, "RegisterValidation") && len(node.Args) > 0:
					value, ok := x.evalString(node.Args[0])
					if !ok || value == "" {
						return true
					}
					if !coveredTags[value] {
						failures = append(failures, fmt.Sprintf("validator RegisterValidation tag is not covered by Runtime API Index: %s = %q", x.evidence(sf, node), value))
					}
				case isStructTagGetCall(node) && len(node.Args) > 0:
					value, ok := x.evalString(node.Args[0])
					if !ok || value == "" {
						failures = append(failures, fmt.Sprintf("validator Field.Tag.Get key is dynamic and not covered by Runtime API Index: %s", x.evidence(sf, node)))

						return true
					}
					if !coveredLabels[value] {
						failures = append(failures, fmt.Sprintf("validator label tag is not covered by Runtime API Index: %s = %q", x.evidence(sf, node), value))
					}
				default:
					if tag, ok := x.validatorTagFromCall(sf, node); ok && !coveredTags[tag] {
						failures = append(failures, fmt.Sprintf("validator rule constructor tag is not covered by Runtime API Index: %s = %q", x.evidence(sf, node), tag))
					}
				}
			case *ast.CompositeLit:
				if !strings.HasSuffix(exprString(sf.Fset, node.Type), "ValidationRule") {
					return true
				}
				for _, elt := range node.Elts {
					kv, ok := elt.(*ast.KeyValueExpr)
					if !ok || keyName(kv.Key) != "RuleTag" {
						continue
					}
					value, ok := x.evalString(kv.Value)
					if !ok || value == "" {
						if isValidatorRuleConstructor(sf.Path, node) && exprString(sf.Fset, kv.Value) == "ruleTag" {
							continue
						}
						failures = append(failures, fmt.Sprintf("ValidationRule.RuleTag is dynamic and not covered by Runtime API Index: %s", x.evidence(sf, kv)))

						continue
					}
					if !coveredTags[value] {
						failures = append(failures, fmt.Sprintf("ValidationRule.RuleTag is not covered by Runtime API Index: %s = %q", x.evidence(sf, kv), value))
					}
				}
			}

			return true
		})
	}

	return failures
}

func isValidatorRuleConstructor(path string, lit *ast.CompositeLit) bool {
	if path != "validator/alphanum.go" && path != "validator/decimal.go" {
		return false
	}

	return strings.HasSuffix(exprString(token.NewFileSet(), lit.Type), "ValidationRule")
}

func (x *extractor) verifyMoldCoverage() []string {
	grammar := x.entryValues("mold tag grammar")
	transformers := x.entryValues("mold transformer tag")
	prefixes := x.entryValues("mold translate kind prefix")
	var failures []string
	for _, sf := range x.files {
		switch sf.Path {
		case "internal/mold/restricted.go":
			x.extractConstStringValues(sf, func(name, value string, evidence []string) {
				if !grammar[value] {
					failures = append(failures, fmt.Sprintf("mold parser token is not covered by Runtime API Index: %s = %q", strings.Join(evidence, ","), value))
				}
			})
		case "internal/mold/mold.go":
			ast.Inspect(sf.File, func(node ast.Node) bool {
				kv, ok := node.(*ast.KeyValueExpr)
				if !ok || keyName(kv.Key) != "tagName" {
					return true
				}
				value, ok := x.evalString(kv.Value)
				if ok && !grammar[value] {
					failures = append(failures, fmt.Sprintf("mold default tag name is not covered by Runtime API Index: %s = %q", x.evidence(sf, kv), value))
				}

				return true
			})
		}
		if !isBuiltInMoldExtensionFile(sf.Path) {
			continue
		}
		ast.Inspect(sf.File, func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.FuncDecl:
				if node.Name.Name != "Tag" {
					return true
				}
				value, ok := x.singleReturnString(sf, node)
				if !ok || value == "" {
					failures = append(failures, fmt.Sprintf("mold transformer Tag() is dynamic and not covered by Runtime API Index: %s", x.evidence(sf, node)))

					return true
				}
				if !transformers[value] {
					failures = append(failures, fmt.Sprintf("mold transformer tag is not covered by Runtime API Index: %s = %q", x.evidence(sf, node), value))
				}
			case *ast.CallExpr:
				if !strings.HasSuffix(exprString(sf.Fset, node.Fun), "strings.HasPrefix") || len(node.Args) != 2 {
					return true
				}
				if exprString(sf.Fset, node.Args[0]) != "kind" {
					return true
				}
				value, ok := x.evalString(node.Args[1])
				if !ok || value == "" {
					failures = append(failures, fmt.Sprintf("mold translate prefix check is dynamic and not covered by Runtime API Index: %s", x.evidence(sf, node)))

					return true
				}
				if !prefixes[value] {
					failures = append(failures, fmt.Sprintf("mold translate prefix is not covered by Runtime API Index: %s = %q", x.evidence(sf, node), value))
				}
			}

			return true
		})
	}

	return failures
}

func (x *extractor) verifyJSONSchemaCoverage() []string {
	var failures []string
	if version := x.goModVersion("github.com/invopop/jsonschema"); version != "v0.14.0" {
		failures = append(failures, fmt.Sprintf("github.com/invopop/jsonschema version changed from v0.14.0 to %s; refresh the MCP jsonschema tag catalog", version))
	}

	covered := x.entryValues("MCP jsonschema tag")
	for _, sf := range x.files {
		for _, field := range structFields(sf) {
			if field.Tag == nil {
				continue
			}
			for _, tagKey := range []string{"jsonschema", "jsonschema_description", "jsonschema_extras"} {
				raw, ok := structTag(field.Tag.Value, tagKey)
				if !ok || raw == "" {
					continue
				}
				if tagKey != "jsonschema" {
					if !covered[tagKey] {
						failures = append(failures, fmt.Sprintf("jsonschema struct tag key is not covered by Runtime API Index: %s = %q", x.evidence(sf, field), tagKey))
					}

					continue
				}
				for _, item := range splitTagCSV(raw) {
					name := strings.SplitN(item, "=", 2)[0]
					if name == "" {
						continue
					}
					if !covered[name] {
						failures = append(failures, fmt.Sprintf("jsonschema tag keyword is not covered by Runtime API Index: %s = %q", x.evidence(sf, field), name))
					}
				}
			}
		}
	}

	return failures
}

func (x *extractor) configFieldKeyIndex() map[string]string {
	roots := map[string]string{
		"AppConfig":        "vef.app",
		"APIConfig":        "vef.api",
		"CORSConfig":       "vef.cors",
		"SecurityConfig":   "vef.security",
		"RedisConfig":      "vef.redis",
		"StorageConfig":    "vef.storage",
		"MonitorConfig":    "vef.monitor",
		"MCPConfig":        "vef.mcp",
		"ApprovalConfig":   "vef.approval",
		"EventConfig":      "vef.event",
		"DataSourceConfig": "vef.data_sources.<name>",
	}

	index := make(map[string]string)
	var walk func(typeName, prefix string)
	walk = func(typeName, prefix string) {
		cfg, ok := x.configTypes[typeName]
		if !ok {
			return
		}
		for _, field := range cfg.Fields {
			key := prefix + "." + field.Key
			index[typeName+"."+field.Name] = key
			if _, ok := x.configTypes[field.TypeName]; ok {
				walk(field.TypeName, key)
			}
		}
	}
	for typeName, root := range roots {
		walk(typeName, root)
	}

	return index
}

func (x *extractor) effectiveDefault(sf sourceFile, fn *ast.FuncDecl, fieldToKey map[string]string) (string, string) {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return "", ""
	}
	typeName := receiverTypeName(fn.Recv.List[0].Type)
	if typeName == "" {
		return "", ""
	}

	effectiveField := strings.TrimPrefix(fn.Name.Name, "Effective")
	var key, value string
	var returnedField, defaultExpr string
	ast.Inspect(fn.Body, func(node ast.Node) bool {
		if key != "" {
			return false
		}
		ret, ok := node.(*ast.ReturnStmt)
		if !ok || len(ret.Results) != 1 {
			return true
		}
		switch expr := ret.Results[0].(type) {
		case *ast.CallExpr:
			fun := exprString(sf.Fset, expr.Fun)
			if strings.HasSuffix(fun, "coalescePositive") && len(expr.Args) == 2 {
				field := selectorField(expr.Args[0])
				if field == "" {
					return true
				}
				key = x.configDefaultKey(sf.Path, typeName, field, fieldToKey)
				value = x.resolveConstExpr(exprString(sf.Fset, expr.Args[1]))
			}
			if strings.HasSuffix(fun, "cmp.Or") && len(expr.Args) == 2 {
				field := selectorField(expr.Args[0])
				if field == "" {
					return true
				}
				key = x.configDefaultKey(sf.Path, typeName, field, fieldToKey)
				value = x.resolveConstExpr(exprString(sf.Fset, expr.Args[1]))
			}
		default:
			if field := selectorField(ret.Results[0]); field != "" {
				returnedField = field
			} else {
				defaultExpr = exprString(sf.Fset, ret.Results[0])
			}
		}

		return true
	})
	if key != "" && value != "" {
		return key, value
	}

	field := effectiveField
	if returnedField != "" && strings.EqualFold(returnedField, effectiveField) {
		field = returnedField
	}
	key = x.configDefaultKey(sf.Path, typeName, field, fieldToKey)
	if key == "" || defaultExpr == "" {
		return "", ""
	}

	return key, x.resolveConstExpr(defaultExpr)
}

func (x *extractor) configDefaultKey(path, typeName, field string, fieldToKey map[string]string) string {
	if field == "" {
		return ""
	}
	if key := fieldToKey[typeName+"."+field]; key != "" {
		return key
	}

	switch path {
	case "event/transport/memory/memory.go":
		return mapFieldToConfigKey("vef.event.transports.memory", map[string]string{
			"QueueSize":  "queue_size",
			"FullPolicy": "full_policy",
		}, field)
	case "event/transport/outbox/outbox.go":
		return mapFieldToConfigKey("vef.event.transports.outbox", map[string]string{
			"RelayInterval":   "relay_interval",
			"MaxRetries":      "max_retries",
			"BatchSize":       "batch_size",
			"LeaseMultiplier": "lease_multiplier",
			"MinLease":        "min_lease",
			"SinkName":        "sink",
		}, field)
	case "event/transport/redisstream/redis_stream.go":
		return mapFieldToConfigKey("vef.event.transports.redis_stream", map[string]string{
			"StreamPrefix":      "stream_prefix",
			"BlockTimeout":      "block_timeout",
			"ClaimIdle":         "claim_idle",
			"ClaimInterval":     "claim_interval",
			"ClaimBatchSize":    "claim_batch_size",
			"ReaperConcurrency":      "reaper_concurrency",
			"SetupTimeout":           "setup_timeout",
			"ConsumerID":             "consumer_id",
			"StartID":                "start_id",
			"IdleGroupSweepInterval": "idle_group_sweep_interval",
		}, field)
	}

	return ""
}

func mapFieldToConfigKey(prefix string, fields map[string]string, field string) string {
	if key := fields[field]; key != "" {
		return prefix + "." + key
	}

	return ""
}

func isConfigDefaultFile(path string) bool {
	return strings.HasPrefix(path, "config/") ||
		path == "event/transport/memory/memory.go" ||
		path == "event/transport/outbox/outbox.go" ||
		path == "event/transport/redisstream/redis_stream.go" ||
		path == "internal/monitor/config.go"
}

func (x *extractor) resolveConstExpr(expr string) string {
	expr = strings.TrimSpace(expr)
	if value, ok := x.constValues[expr]; ok {
		if constExpr, ok := x.constExprs[expr]; ok && constExpr != expr {
			return fmt.Sprintf("%s (%s)", constExpr, value)
		}

		return value
	}
	if constExpr, ok := x.constExprs[expr]; ok {
		return constExpr
	}
	if value, err := strconv.Unquote(expr); err == nil {
		return value
	}

	return expr
}

func (x *extractor) extractBuiltInResources() {
	for _, sf := range x.files {
		if !isRuntimeResourceFile(sf.Path) {
			continue
		}

		ast.Inspect(sf.File, func(node ast.Node) bool {
			lit, ok := node.(*ast.CompositeLit)
			if !ok {
				return true
			}

			res, ok := x.resourceFromComposite(sf, lit)
			if !ok {
				res, ok = x.resourceFromCall(sf, lit)
			}
			if !ok {
				return true
			}
			x.add("built-in resource", res.Name, res.Kind, []string{res.Source})

			for _, action := range x.explicitOperations(sf, lit, res) {
				x.addResourceAction(res, action)
			}
			for _, action := range x.crudOperations(sf, lit, res) {
				x.addResourceAction(res, action)
			}

			return true
		})
	}
}

func (x *extractor) resourceFromCall(sf sourceFile, lit *ast.CompositeLit) (resourceRef, bool) {
	name, kind, ok := resourceCallFromElts(sf, lit.Elts)
	if !ok {
		return resourceRef{}, false
	}

	return resourceRef{Name: name, Kind: kind, Source: x.evidence(sf, lit)}, true
}

func (x *extractor) addResourceAction(res resourceRef, action actionMeta) {
	if action.Action == "" {
		return
	}
	details := []string{"resource kind: " + res.Kind}
	if action.RequiredPermission != "" {
		details = append(details, "permission: "+action.RequiredPermission)
	}
	if action.Public {
		details = append(details, "public")
	}
	if action.EnableAudit {
		details = append(details, "audit enabled")
	}
	x.add("built-in resource action", res.Name+"/"+action.Action, action.Action, []string{action.Source}, details...)
}

func (x *extractor) resourceFromComposite(sf sourceFile, lit *ast.CompositeLit) (resourceRef, bool) {
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok || keyName(kv.Key) != "Resource" {
			continue
		}
		name, kind, ok := resourceCall(sf, kv.Value)
		if !ok {
			continue
		}

		return resourceRef{Name: name, Kind: kind, Source: x.evidence(sf, kv.Value)}, true
	}

	return resourceRef{}, false
}

func resourceCall(sf sourceFile, expr ast.Expr) (string, string, bool) {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return "", "", false
	}
	return resourceCallFromCall(sf, call)
}

func resourceCallFromElts(sf sourceFile, elts []ast.Expr) (string, string, bool) {
	for _, elt := range elts {
		call, ok := elt.(*ast.CallExpr)
		if !ok {
			continue
		}
		if name, kind, ok := resourceCallFromCall(sf, call); ok {
			return name, kind, true
		}
	}

	return "", "", false
}

func resourceCallFromCall(sf sourceFile, call *ast.CallExpr) (string, string, bool) {
	fun := exprString(sf.Fset, call.Fun)
	var kind string
	switch {
	case strings.HasSuffix(fun, "NewRPCResource"):
		kind = "rpc"
	case strings.HasSuffix(fun, "NewRESTResource"):
		kind = "rest"
	default:
		return "", "", false
	}
	if len(call.Args) == 0 {
		return "", "", false
	}
	name, ok := stringValue(call.Args[0])

	return name, kind, ok
}

func (x *extractor) explicitOperations(sf sourceFile, lit *ast.CompositeLit, res resourceRef) []actionMeta {
	var actions []actionMeta
	ast.Inspect(lit, func(node ast.Node) bool {
		op, ok := node.(*ast.CompositeLit)
		if !ok || !isOperationSpec(sf, op) {
			return true
		}
		meta := actionMeta{Kind: res.Kind, Source: x.evidence(sf, op)}
		for _, elt := range op.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			switch keyName(kv.Key) {
			case "Action":
				meta.Action, _ = stringValue(kv.Value)
			case "RequiredPermission":
				meta.RequiredPermission, _ = stringValue(kv.Value)
			case "Public":
				meta.Public = boolValue(kv.Value)
			case "EnableAudit":
				meta.EnableAudit = boolValue(kv.Value)
			}
		}
		if meta.Action != "" {
			actions = append(actions, meta)
		}

		return true
	})

	return actions
}

func isOperationSpec(sf sourceFile, lit *ast.CompositeLit) bool {
	typeText := exprString(sf.Fset, lit.Type)
	return strings.HasSuffix(typeText, "OperationSpec") || strings.Contains(typeText, ".OperationSpec")
}

func (x *extractor) crudOperations(sf sourceFile, lit *ast.CompositeLit, res resourceRef) []actionMeta {
	var actions []actionMeta
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		field := keyName(kv.Key)
		defaultAction, ok := crudActionForField(field, res.Kind)
		if !ok {
			continue
		}
		meta := inspectBuilderChain(sf, kv.Value)
		meta.Source = x.evidence(sf, kv)
		meta.Kind = res.Kind
		if meta.Action == "" {
			meta.Action = defaultAction
		}
		actions = append(actions, meta)
	}

	return actions
}

func inspectBuilderChain(sf sourceFile, expr ast.Expr) actionMeta {
	var meta actionMeta
	ast.Inspect(expr, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		switch sel.Sel.Name {
		case "Action":
			if len(call.Args) > 0 {
				meta.Action, _ = stringValue(call.Args[0])
			}
		case "RequiredPermission":
			if len(call.Args) > 0 {
				meta.RequiredPermission, _ = stringValue(call.Args[0])
			}
		case "Public":
			meta.Public = true
		case "EnableAudit":
			meta.EnableAudit = true
		case "ResourceKind":
			if len(call.Args) > 0 {
				kind := exprString(sf.Fset, call.Args[0])
				if strings.HasSuffix(kind, "KindREST") {
					meta.Kind = "rest"
				} else if strings.HasSuffix(kind, "KindRPC") {
					meta.Kind = "rpc"
				}
			}
		}

		return true
	})

	return meta
}

func crudActionForField(field, kind string) (string, bool) {
	rpc := map[string]string{
		"Create": "create", "Update": "update", "Delete": "delete",
		"CreateMany": "create_many", "UpdateMany": "update_many", "DeleteMany": "delete_many",
		"FindOne": "find_one", "FindAll": "find_all", "FindPage": "find_page",
		"FindOptions": "find_options", "FindTree": "find_tree", "FindTreeOptions": "find_tree_options",
		"Import": "import", "Export": "export",
	}
	rest := map[string]string{
		"Create": "post /", "Update": "put /:id", "Delete": "delete /:id",
		"CreateMany": "post /many", "UpdateMany": "put /many", "DeleteMany": "delete /many",
		"FindOne": "get /:id", "FindAll": "get /", "FindPage": "get /page",
		"FindOptions": "get /options", "FindTree": "get /tree", "FindTreeOptions": "get /tree/options",
		"Import": "post /import", "Export": "get /export",
	}
	if kind == "rest" {
		value, ok := rest[field]
		return value, ok
	}
	value, ok := rpc[field]

	return value, ok
}

func isRuntimeResourceFile(path string) bool {
	return strings.HasPrefix(path, "internal/security/") ||
		strings.HasPrefix(path, "internal/storage/") ||
		strings.HasPrefix(path, "internal/schema/") ||
		strings.HasPrefix(path, "internal/monitor/") ||
		strings.HasPrefix(path, "internal/approval/resource/")
}

func (x *extractor) extractCLI() {
	for _, sf := range x.files {
		if !strings.HasPrefix(sf.Path, "cmd/vef-cli/cmd/") {
			continue
		}

		ast.Inspect(sf.File, func(node ast.Node) bool {
			lit, ok := node.(*ast.CompositeLit)
			if !ok {
				return true
			}
			if command := commandFromComposite(sf, lit); command.Use != "" {
				command.Flags = flagsNearCommand(sf, lit)
				command.Required = requiredFlagsNearCommand(sf, lit)
				x.addCLICommand(command)
			}

			return true
		})

		for _, decl := range sf.File.Decls {
			switch decl := decl.(type) {
			case *ast.FuncDecl:
				x.extractCLIFromBlock(sf, decl.Body)
			case *ast.GenDecl:
				x.extractCLIGenDecl(sf, decl)
			}
		}
	}
}

func (x *extractor) extractCLIGenDecl(sf sourceFile, gen *ast.GenDecl) {
	for _, spec := range gen.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		for _, value := range vs.Values {
			if command := commandFromComposite(sf, value); command.Use != "" {
				x.addCLICommand(command)
			}
		}
	}
}

func (x *extractor) extractCLIFromBlock(sf sourceFile, body *ast.BlockStmt) {
	if body == nil {
		return
	}
	command := cliCommand{}
	ast.Inspect(body, func(node ast.Node) bool {
		if command.Use == "" {
			if assign, ok := node.(*ast.AssignStmt); ok {
				for _, rhs := range assign.Rhs {
					if candidate := commandFromComposite(sf, rhs); candidate.Use != "" {
						command = candidate
					}
				}
			}
		}
		if command.Use == "" {
			return true
		}
		if call, ok := node.(*ast.CallExpr); ok {
			if flag := cliFlagFromCall(sf, call); flag.Name != "" {
				command.Flags = append(command.Flags, flag)
			}
			if required := requiredFlagFromCall(call); required != "" {
				command.Required = appendUnique(command.Required, required)
			}
		}

		return true
	})
	if command.Use != "" {
		x.addCLICommand(command)
	}
}

type cliCommand struct {
	Use      string
	Short    string
	Source   string
	Flags    []cliFlag
	Required []string
}

type cliFlag struct {
	Name    string
	Short   string
	Default string
	Usage   string
	Source  string
}

func commandFromComposite(sf sourceFile, expr ast.Expr) cliCommand {
	if unary, ok := expr.(*ast.UnaryExpr); ok {
		expr = unary.X
	}
	lit, ok := expr.(*ast.CompositeLit)
	if !ok || !strings.Contains(exprString(sf.Fset, lit.Type), "cobra.Command") {
		return cliCommand{}
	}
	cmd := cliCommand{Source: fmt.Sprintf("%s:%d", sf.Path, sf.Fset.Position(lit.Pos()).Line)}
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		switch keyName(kv.Key) {
		case "Use":
			cmd.Use, _ = stringValue(kv.Value)
		case "Short":
			cmd.Short, _ = stringValue(kv.Value)
		}
	}

	return cmd
}

func flagsNearCommand(sf sourceFile, root ast.Node) []cliFlag {
	var flags []cliFlag
	ast.Inspect(root, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		if flag := cliFlagFromCall(sf, call); flag.Name != "" {
			flags = append(flags, flag)
		}

		return true
	})

	return flags
}

func requiredFlagsNearCommand(_ sourceFile, root ast.Node) []string {
	var required []string
	ast.Inspect(root, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		if name := requiredFlagFromCall(call); name != "" {
			required = appendUnique(required, name)
		}

		return true
	})

	return required
}

func (x *extractor) addCLICommand(command cliCommand) {
	use := strings.Fields(command.Use)
	name := command.Use
	if len(use) > 0 {
		name = use[0]
	}
	x.add("CLI command", name, command.Use, []string{command.Source}, command.Short)
	for _, flag := range command.Flags {
		details := []string{flag.Usage}
		if flag.Default != "" {
			details = append(details, "default: "+flag.Default)
		}
		if flag.Short != "" {
			details = append(details, "short: -"+flag.Short)
		}
		if contains(command.Required, flag.Name) {
			details = append(details, "required")
		}
		x.add("CLI flag", name+" --"+flag.Name, "--"+flag.Name, []string{flag.Source}, details...)
	}
}

func cliFlagFromCall(sf sourceFile, call *ast.CallExpr) cliFlag {
	fun := exprString(sf.Fset, call.Fun)
	if !isCLIFlagDefinitionCall(fun) {
		return cliFlag{}
	}
	kind, ok := cliFlagHelperKind(fun)
	if !ok {
		return cliFlag{}
	}
	if len(call.Args) < kind.minArgs {
		return cliFlag{}
	}
	name, ok := stringValue(call.Args[kind.nameArg])
	if !ok {
		return cliFlag{}
	}
	short := ""
	if kind.shortArg >= 0 {
		short, _ = stringValue(call.Args[kind.shortArg])
	}
	def := ""
	if kind.defaultArg >= 0 {
		def = exprString(sf.Fset, call.Args[kind.defaultArg])
	}
	if unquoted, err := strconv.Unquote(def); err == nil {
		def = unquoted
	}
	usage, _ := stringValue(call.Args[kind.usageArg])

	return cliFlag{
		Name:    name,
		Short:   short,
		Default: def,
		Usage:   usage,
		Source:  fmt.Sprintf("%s:%d", sf.Path, sf.Fset.Position(call.Pos()).Line),
	}
}

type cliFlagHelper struct {
	minArgs    int
	nameArg    int
	shortArg   int
	defaultArg int
	usageArg   int
}

func cliFlagHelperKind(fun string) (cliFlagHelper, bool) {
	switch {
	case strings.HasSuffix(fun, ".StringP"),
		strings.HasSuffix(fun, ".BoolP"),
		strings.HasSuffix(fun, ".IntP"):
		return cliFlagHelper{minArgs: 4, nameArg: 0, shortArg: 1, defaultArg: 2, usageArg: 3}, true
	case strings.HasSuffix(fun, ".String"),
		strings.HasSuffix(fun, ".Bool"),
		strings.HasSuffix(fun, ".Int"):
		return cliFlagHelper{minArgs: 3, nameArg: 0, shortArg: -1, defaultArg: 1, usageArg: 2}, true
	case strings.HasSuffix(fun, ".StringVarP"),
		strings.HasSuffix(fun, ".BoolVarP"),
		strings.HasSuffix(fun, ".IntVarP"):
		return cliFlagHelper{minArgs: 5, nameArg: 1, shortArg: 2, defaultArg: 3, usageArg: 4}, true
	case strings.HasSuffix(fun, ".StringVar"),
		strings.HasSuffix(fun, ".BoolVar"),
		strings.HasSuffix(fun, ".IntVar"):
		return cliFlagHelper{minArgs: 4, nameArg: 1, shortArg: -1, defaultArg: 2, usageArg: 3}, true
	default:
		return cliFlagHelper{}, false
	}
}

func isCLIFlagDefinitionCall(fun string) bool {
	return strings.Contains(fun, ".Flags().") || strings.Contains(fun, ".PersistentFlags().")
}

func isCLIFlagReader(method string) bool {
	switch method {
	case "GetString", "GetBool", "GetInt", "Changed", "Lookup", "Visit", "VisitAll",
		"PrintDefaults", "FlagUsages", "FlagUsagesWrapped", "SortFlags":
		return true
	default:
		return false
	}
}

func requiredFlagFromCall(call *ast.CallExpr) string {
	fun := exprString(token.NewFileSet(), call.Fun)
	if !strings.HasSuffix(fun, ".MarkFlagRequired") || len(call.Args) == 0 {
		return ""
	}
	name, _ := stringValue(call.Args[0])

	return name
}

func (x *extractor) extractJSONFields() {
	for _, sf := range x.files {
		if !isRuntimeDTOFile(sf.Path) {
			continue
		}
		for _, decl := range sf.File.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}
				x.addJSONFieldsForStruct(sf, ts.Name.Name, st)
			}
		}
	}
}

func (x *extractor) addJSONFieldsForStruct(sf sourceFile, typeName string, st *ast.StructType) {
	for _, field := range st.Fields.List {
		if field.Tag == nil {
			continue
		}
		jsonName, ok := structTag(field.Tag.Value, "json")
		if !ok {
			continue
		}
		jsonName = strings.Split(jsonName, ",")[0]
		if jsonName == "" || jsonName == "-" {
			continue
		}
		for _, name := range field.Names {
			details := []string{"Go field: " + typeName + "." + name.Name, "type: " + exprString(sf.Fset, field.Type)}
			for _, tagName := range []string{"form", "validate", "search", "tabular", "meta", "mold", "jsonschema"} {
				if tag, ok := structTag(field.Tag.Value, tagName); ok && tag != "" {
					details = append(details, tagName+`: "`+tag+`"`)
				}
			}
			x.add("JSON wire field", typeName+"."+name.Name, jsonName, []string{x.evidence(sf, field)}, details...)
		}
	}
}

func isRuntimeDTOFile(path string) bool {
	exactFiles := map[string]bool{
		"internal/orm/model.go": true,
		// The built-in form-editor parser's designer-document contract and the
		// persisted business-projection record-key shape are runtime wire
		// surfaces despite living under internal/.
		"internal/approval/formeditor/schema.go":  true,
		"internal/approval/binding/record_key.go": true,
	}
	if exactFiles[path] {
		return true
	}

	prefixes := []string{
		"api/", "result/", "page/", "crud/",
		"ai/", "approval/", "storage/", "monitor/", "mcp/",
		"security/", "schema/", "mold/",
		"event/", "internal/security/", "internal/storage/", "internal/schema/",
		"internal/monitor/", "internal/approval/resource/", "internal/approval/shared/",
		"internal/mcp/",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}

func (x *extractor) extractRuntimeEnumValues() {
	for _, sf := range x.files {
		if !isRuntimeEnumFile(sf.Path) {
			continue
		}
		for _, decl := range sf.File.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.CONST {
				continue
			}
			var lastValues []ast.Expr
			for _, spec := range gen.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				typeName := namedType(vs.Type)
				if len(vs.Values) > 0 {
					lastValues = vs.Values
				}
				if typeName == "" {
					continue
				}
				for i, name := range vs.Names {
					var valueExpr ast.Expr
					if i < len(vs.Values) {
						valueExpr = vs.Values[i]
					} else if i < len(lastValues) {
						valueExpr = lastValues[i]
					} else if len(lastValues) == 1 {
						valueExpr = lastValues[0]
					}
					if valueExpr == nil {
						continue
					}
					value, ok := x.evalString(valueExpr)
					if !ok {
						continue
					}
					x.add("runtime enum value", name.Name+" ("+typeName+")", value, []string{x.evidence(sf, name)})
				}
			}
		}
	}
}

func isRuntimeEnumFile(path string) bool {
	if strings.Contains(path, "/testdata/") || strings.HasPrefix(path, "cmd/vef-cli/") {
		return false
	}
	if isRuntimeDTOFile(path) || strings.HasPrefix(path, "config/") || strings.HasPrefix(path, "search/") {
		return true
	}
	if strings.HasPrefix(path, "internal/") {
		return false
	}

	return true
}

func (x *extractor) extractMCP() {
	for _, sf := range x.files {
		if !strings.HasPrefix(sf.Path, "internal/mcp/") {
			continue
		}
		ast.Inspect(sf.File, func(node ast.Node) bool {
			lit, ok := node.(*ast.CompositeLit)
			if !ok {
				return true
			}
			typeText := exprString(sf.Fset, lit.Type)
			switch {
			case strings.HasSuffix(typeText, "Tool"):
				x.extractNamedMCPObject(sf, lit, "MCP tool")
			case strings.HasSuffix(typeText, "Prompt"):
				x.extractNamedMCPObject(sf, lit, "MCP prompt")
			case strings.HasSuffix(typeText, "Resource"):
				x.extractURIMCPObject(sf, lit, "MCP resource")
			case strings.HasSuffix(typeText, "ResourceTemplate"):
				x.extractURITemplateMCPObject(sf, lit, "MCP resource template")
			}

			return true
		})
	}
}

func (x *extractor) extractNamedMCPObject(sf sourceFile, lit *ast.CompositeLit, category string) {
	var name, desc string
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		switch keyName(kv.Key) {
		case "Name":
			name, _ = x.evalString(kv.Value)
		case "Description":
			desc, _ = x.evalString(kv.Value)
		}
	}
	if name != "" {
		x.add(category, name, name, []string{x.evidence(sf, lit)}, desc)
	}
}

func (x *extractor) extractURIMCPObject(sf sourceFile, lit *ast.CompositeLit, category string) {
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if ok && keyName(kv.Key) == "URI" {
			if uri, ok := x.evalString(kv.Value); ok {
				x.add(category, uri, uri, []string{x.evidence(sf, kv)})
			}
		}
	}
}

func (x *extractor) extractURITemplateMCPObject(sf sourceFile, lit *ast.CompositeLit, category string) {
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if ok && keyName(kv.Key) == "URITemplate" {
			if uri, ok := x.evalString(kv.Value); ok {
				x.add(category, uri, uri, []string{x.evidence(sf, kv)})
			}
		}
	}
}

func (x *extractor) validatorTagFromCall(sf sourceFile, call *ast.CallExpr) (string, bool) {
	if len(call.Args) == 0 {
		return "", false
	}
	fun := exprString(sf.Fset, call.Fun)
	if !strings.HasSuffix(fun, "newRegexRule") && !strings.HasSuffix(fun, "newDecimalComparisonRule") {
		return "", false
	}

	return x.evalString(call.Args[0])
}

func isStructTagGetCall(call *ast.CallExpr) bool {
	return selectorMethod(call.Fun) == "Get" && strings.HasSuffix(exprString(token.NewFileSet(), call.Fun), ".Tag.Get")
}

func (x *extractor) extractValidatorRules() {
	for _, sf := range x.files {
		if !strings.HasPrefix(sf.Path, "validator/") {
			continue
		}
		ast.Inspect(sf.File, func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.CallExpr:
				if tag, ok := x.validatorTagFromCall(sf, node); ok {
					x.add("validator tag", tag, tag, []string{x.evidence(sf, node)})
				}
				if isStructTagGetCall(node) && len(node.Args) > 0 {
					if tag, ok := x.evalString(node.Args[0]); ok {
						x.add("validator label tag", tag, tag, []string{x.evidence(sf, node)})
					}
				}
			case *ast.CompositeLit:
				if !strings.HasSuffix(exprString(sf.Fset, node.Type), "ValidationRule") {
					return true
				}
				for _, elt := range node.Elts {
					kv, ok := elt.(*ast.KeyValueExpr)
					if !ok || keyName(kv.Key) != "RuleTag" {
						continue
					}
					if tag, ok := x.evalString(kv.Value); ok {
						x.add("validator tag", tag, tag, []string{x.evidence(sf, kv)})
					}
				}
			}

			return true
		})
	}
}

func (x *extractor) extractStructTagGrammars() {
	for _, sf := range x.files {
		switch sf.Path {
		case "search/constants.go":
			x.extractConstStringValues(sf, func(name, value string, evidence []string) {
				switch {
				case name == "TagSearch":
					x.add("search tag grammar", "tag name", value, evidence)
				case strings.HasPrefix(name, "Attr"):
					x.add("search tag grammar", name, value, evidence)
				case strings.HasPrefix(name, "Param"):
					x.add("search tag grammar", name, value, evidence)
				case name == "IgnoreField":
					x.add("search tag grammar", name, value, evidence)
				case strings.HasPrefix(name, "Type"):
					x.add("search tag grammar", name, value, evidence)
				default:
					if operatorConstNames[name] {
						x.add("search tag grammar", "operator "+name, value, evidence)
					}
				}
			})
		case "tabular/constants.go":
			x.extractConstStringValues(sf, func(name, value string, evidence []string) {
				switch {
				case name == "TagTabular":
					x.add("tabular tag grammar", "tag name", value, evidence)
				case strings.HasPrefix(name, "Attr"):
					x.add("tabular tag grammar", name, value, evidence)
				case name == "IgnoreField":
					x.add("tabular tag grammar", name, value, evidence)
				}
			})
		case "storage/file_refs.go":
			x.extractConstStringValues(sf, func(name, value string, evidence []string) {
				switch {
				case name == "tagMeta":
					x.add("meta tag grammar", "tag name", value, evidence)
				case strings.HasPrefix(name, "MetaType"):
					x.add("meta tag grammar", name, value, evidence)
				}
			})
			x.add("meta tag grammar", "dive", "dive", []string{"storage/file_refs.go:239"})
			x.add("meta tag grammar", "attribute pair delimiter", "space", []string{"storage/file_refs.go:212"})
			x.add("meta tag grammar", "attribute key/value delimiter", ":", []string{"storage/file_refs.go:214"})
		}
	}
}

var operatorConstNames = map[string]bool{
	"Equals": true, "NotEquals": true, "GreaterThan": true, "GreaterThanOrEqual": true,
	"LessThan": true, "LessThanOrEqual": true, "Between": true, "NotBetween": true,
	"In": true, "NotIn": true, "IsNull": true, "IsNotNull": true,
	"Contains": true, "NotContains": true, "StartsWith": true, "NotStartsWith": true,
	"EndsWith": true, "NotEndsWith": true, "ContainsIgnoreCase": true, "NotContainsIgnoreCase": true,
	"StartsWithIgnoreCase": true, "NotStartsWithIgnoreCase": true, "EndsWithIgnoreCase": true,
	"NotEndsWithIgnoreCase": true,
}

func (x *extractor) extractMoldGrammar() {
	for _, sf := range x.files {
		switch sf.Path {
		case "internal/mold/restricted.go":
			x.extractConstStringValues(sf, func(name, value string, evidence []string) {
				x.add("mold tag grammar", name, value, evidence)
			})
		case "internal/mold/mold.go":
			ast.Inspect(sf.File, func(node ast.Node) bool {
				kv, ok := node.(*ast.KeyValueExpr)
				if !ok || keyName(kv.Key) != "tagName" {
					return true
				}
				if value, ok := x.evalString(kv.Value); ok {
					x.add("mold tag grammar", "tag name", value, []string{x.evidence(sf, kv)})
				}

				return true
			})
		}

		if !isBuiltInMoldExtensionFile(sf.Path) {
			continue
		}
		ast.Inspect(sf.File, func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.FuncDecl:
				if node.Name.Name == "Tag" {
					if value, ok := x.singleReturnString(sf, node); ok {
						x.add("mold transformer tag", value, value, []string{x.evidence(sf, node)})
					}
				}
			case *ast.CallExpr:
				if !strings.HasSuffix(exprString(sf.Fset, node.Fun), "strings.HasPrefix") || len(node.Args) != 2 {
					return true
				}
				if exprString(sf.Fset, node.Args[0]) != "kind" {
					return true
				}
				if value, ok := x.evalString(node.Args[1]); ok {
					x.add("mold translate kind prefix", value, value, []string{x.evidence(sf, node)})
				}
			}

			return true
		})
	}
}

func isBuiltInMoldExtensionFile(path string) bool {
	return path == "internal/mold/translate.go" ||
		path == "internal/mold/dictionary_translator.go" ||
		path == "internal/expression/transformer.go"
}

func (x *extractor) extractJSONSchemaTags() {
	tags := []string{
		"-", "required", "nullable",
		"title", "description", "type", "anchor",
		"oneof_required", "anyof_required", "oneof_ref", "oneof_type", "anyof_ref", "anyof_type",
		"default", "example", "enum",
		"minLength", "maxLength", "pattern", "format", "readOnly", "writeOnly",
		"minimum", "maximum", "exclusiveMinimum", "exclusiveMaximum", "multipleOf",
		"minItems", "maxItems", "uniqueItems",
		"jsonschema_description", "jsonschema_extras",
	}
	for _, tag := range tags {
		x.add("MCP jsonschema tag", tag, tag, []string{"github.com/invopop/jsonschema@v0.14.0/reflect.go:613"})
	}
}

func (x *extractor) extractI18NMessageKeys() {
	for _, sf := range x.files {
		for _, occurrence := range x.i18nKeyOccurrences(sf) {
			if occurrence.Value == "<dynamic>" {
				x.add("i18n key indirection", occurrence.Source, occurrence.Detail, []string{occurrence.Source})
				continue
			}
			x.add("i18n message key", occurrence.Value, occurrence.Value, []string{occurrence.Source}, occurrence.Detail)
		}
	}
}

func (x *extractor) extractEventTransportContracts() {
	x.add("event transport contract", "outbox DLQ header", "vef.dlq", []string{"internal/event/transport/outbox/relay.go:37"})
	x.add("event transport contract", "outbox DLQ header value", "1", []string{"internal/event/transport/outbox/relay.go:178"})
	x.add("event transport contract", "outbox DLQ topic prefix", "vef-dlq.", []string{"internal/event/transport/outbox/relay.go:225"})
	x.add("event transport contract", "outbox persisted error max bytes", "256", []string{"internal/event/transport/outbox/relay.go:20"})
	x.add("event transport contract", "outbox retry backoff cap", "1h", []string{"internal/event/transport/outbox/relay.go:211"})
	x.add("event transport contract", "outbox retry backoff formula", "2^retryCount seconds capped at 1h", []string{"internal/event/transport/outbox/relay.go:214"})
}

func (x *extractor) i18nKeyOccurrences(sf sourceFile) []stringOccurrence {
	var result []stringOccurrence
	ast.Inspect(sf.File, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.CallExpr:
			if isI18NCall(node.Fun) && len(node.Args) > 0 {
				value, ok := x.evalString(node.Args[0])
				if !ok {
					result = append(result, stringOccurrence{
						Value:  "<dynamic>",
						Source: x.evidence(sf, node),
						Detail: dynamicI18NDetail(sf.Path, x.evidence(sf, node)),
					})

					return true
				}
				result = append(result, stringOccurrence{
					Value:  value,
					Source: x.evidence(sf, node),
					Detail: "i18n.T call",
				})
			}
			if tag, ok := validatorRuleI18NKey(sf, node, x); ok {
				result = append(result, stringOccurrence{
					Value:  tag,
					Source: x.evidence(sf, node),
					Detail: "validator rule message key",
				})
			}
		case *ast.Field:
			if node.Tag == nil {
				return true
			}
			if tag, ok := structTag(node.Tag.Value, "label_i18n"); ok && tag != "" && tag != "-" {
				result = append(result, stringOccurrence{
					Value:  tag,
					Source: x.evidence(sf, node),
					Detail: "label_i18n struct tag",
				})
			}
		}

		return true
	})

	return result
}

func isI18NCall(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "T" {
		return false
	}
	recv, ok := sel.X.(*ast.Ident)

	return ok && recv.Name == "i18n"
}

func validatorRuleI18NKey(sf sourceFile, call *ast.CallExpr, x *extractor) (string, bool) {
	if !strings.HasPrefix(sf.Path, "validator/") || len(call.Args) == 0 {
		return "", false
	}
	fun := exprString(sf.Fset, call.Fun)
	switch {
	case strings.HasSuffix(fun, "newRegexRule"):
		if len(call.Args) < 4 {
			return "", false
		}

		return x.evalString(call.Args[3])
	case strings.HasSuffix(fun, "newDecimalComparisonRule"):
		if len(call.Args) < 3 {
			return "", false
		}

		return x.evalString(call.Args[2])
	default:
		return "", false
	}
}

func dynamicI18NDetail(path, evidence string) string {
	switch path {
	case "validator/validator.go":
		return "dynamic key sourced from label_i18n struct tags; tag values are indexed separately"
	case "validator/rule.go":
		return "dynamic key sourced from ValidationRule.ErrMessageI18nKey; built-in rule keys are indexed separately"
	case "internal/app/error.go":
		return "dynamic key sourced from fiberErrorMappings; mapped result/security message constants are indexed separately"
	case "internal/api/middleware/audit.go":
		return "dynamic key sourced from app.MapFiberError; mapped result/security message constants are indexed separately"
	default:
		return "dynamic i18n key source reviewed: " + evidence
	}
}

func sortEntries(entries []RuntimeEntry) {
	sort.Slice(entries, func(i, j int) bool {
		a, b := entries[i], entries[j]
		if a.Category != b.Category {
			return a.Category < b.Category
		}
		if a.Name != b.Name {
			return a.Name < b.Name
		}
		if a.Value != b.Value {
			return a.Value < b.Value
		}

		return a.ID < b.ID
	})
}

func fingerprint(entries []RuntimeEntry) string {
	h := sha256.New()
	enc := json.NewEncoder(h)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(entries)

	return hex.EncodeToString(h.Sum(nil))
}

func buildTerms(entry RuntimeEntry) []string {
	terms := make([]string, 0)
	for _, s := range []string{entry.Category, entry.Name, entry.Value} {
		terms = append(terms, splitTerms(s)...)
	}
	for _, detail := range entry.Details {
		terms = append(terms, splitTerms(detail)...)
	}
	sort.Strings(terms)

	return unique(terms)
}

func splitTerms(s string) []string {
	re := regexp.MustCompile(`[A-Za-z0-9_./:@-]+`)
	return re.FindAllString(s, -1)
}

func englishDocument(ledger RuntimeLedger) string {
	var b strings.Builder
	b.WriteString("---\nsidebar_position: 91\n---\n\n")
	b.WriteString("# Runtime API Index\n\n")
	b.WriteString("This page is generated from the current VEF Framework Go source tree. It covers runtime contracts users call, configure, send, receive, or match: HTTP/RPC protocol fields, built-in resources and actions, CLI commands and flags, configuration keys and defaults, events, error codes, wire JSON fields, tag grammars, MCP endpoints/tools/prompts, and runtime enum values.\n\n")
	b.WriteString("It intentionally excludes test fixtures, internal log strings, and implementation-only literals. The exported Go import surface is tracked separately in [Public API Index](./public-api-index).\n\n")
	b.WriteString("The complete public API audit is the union of this runtime index, the exported Go API index, and the package reviews in `scripts/api-contract-ledger.json`. A user-facing API change must update all affected audit artifacts before the docs review is complete.\n\n")
	b.WriteString("Regenerate and verify this page whenever the framework runtime surface changes:\n\n")
	b.WriteString("```bash\n")
	b.WriteString("(cd ../vef-framework-go && go run ../vef-framework-go-docs/scripts/verify-runtime-api-audit.go -source . -out ../vef-framework-go-docs -write)\n")
	b.WriteString("(cd ../vef-framework-go && go run ../vef-framework-go-docs/scripts/verify-runtime-api-audit.go -source . -out ../vef-framework-go-docs)\n")
	b.WriteString("```\n\n")
	b.WriteString(fmt.Sprintf("Fingerprint: `%s`\nEntries: `%d`\n\n", ledger.Fingerprint, ledger.EntryCount))
	writeCoverageSection(&b, ledger.Coverage)
	writeEntrySections(&b, ledger.Entries)

	return strings.TrimRight(b.String(), "\n") + "\n"
}

func chineseDocument(ledger RuntimeLedger) string {
	var b strings.Builder
	b.WriteString("---\nsidebar_position: 91\n---\n\n")
	b.WriteString("# Runtime API Index\n\n")
	b.WriteString("本页由当前 VEF Framework Go 源码生成，覆盖用户会直接调用、配置、发送、接收或匹配的运行时 contract：HTTP/RPC 协议字段、内置 resource/action、CLI 命令和 flags、配置键与默认值、事件、错误码、JSON wire 字段、结构体标签语法、MCP endpoint/tool/prompt，以及运行时枚举值。\n\n")
	b.WriteString("测试 fixture、内部日志字符串和纯实现细节字面量不会进入本索引。导入 Go 包时看到的 exported API 单独由 [Public API Index](./public-api-index) 跟踪。\n\n")
	b.WriteString("完整的公开 API 审计由本运行时索引、exported Go API 索引，以及 `scripts/api-contract-ledger.json` 中的 package review 共同组成。任何用户可见 API 变化都必须同步更新受影响的审计产物，文档审查才算完成。\n\n")
	b.WriteString("框架运行时公开面变化后，使用下面的命令重新生成并验证：\n\n")
	b.WriteString("```bash\n")
	b.WriteString("(cd ../vef-framework-go && go run ../vef-framework-go-docs/scripts/verify-runtime-api-audit.go -source . -out ../vef-framework-go-docs -write)\n")
	b.WriteString("(cd ../vef-framework-go && go run ../vef-framework-go-docs/scripts/verify-runtime-api-audit.go -source . -out ../vef-framework-go-docs)\n")
	b.WriteString("```\n\n")
	b.WriteString(fmt.Sprintf("Fingerprint: `%s`\nEntries: `%d`\n\n", ledger.Fingerprint, ledger.EntryCount))
	writeCoverageSection(&b, ledger.Coverage)
	writeEntrySections(&b, ledger.Entries)

	return strings.TrimRight(b.String(), "\n") + "\n"
}

func writeCoverageSection(b *strings.Builder, coverage []RuntimeCoverage) {
	b.WriteString("## Coverage Evidence\n\n")
	b.WriteString("| Category | Entries | Tier | Extractor | Method | Known residual |\n")
	b.WriteString("| --- | ---: | --- | --- | --- | --- |\n")
	for _, item := range coverage {
		b.WriteString("| `")
		b.WriteString(escapeTable(item.Category))
		b.WriteString("` | ")
		b.WriteString(strconv.Itoa(item.EntryCount))
		b.WriteString(" | ")
		b.WriteString(escapeTable(item.Tier))
		b.WriteString(" | `")
		b.WriteString(escapeTable(item.Extractor))
		b.WriteString("` | ")
		b.WriteString(escapeTable(item.Method))
		b.WriteString(" | ")
		b.WriteString(escapeTable(item.KnownResidual))
		b.WriteString(" |\n")
	}
	b.WriteString("\n")
}

func writeEntrySections(b *strings.Builder, entries []RuntimeEntry) {
	byCategory := make(map[string][]RuntimeEntry)
	for _, entry := range entries {
		byCategory[entry.Category] = append(byCategory[entry.Category], entry)
	}
	categories := make([]string, 0, len(byCategory))
	for category := range byCategory {
		categories = append(categories, category)
	}
	sort.Strings(categories)

	for _, category := range categories {
		b.WriteString("## " + category + "\n\n")
		b.WriteString("| Name | Value | Details | Source |\n")
		b.WriteString("| --- | --- | --- | --- |\n")
		for _, entry := range byCategory[category] {
			b.WriteString("| `")
			b.WriteString(escapeTable(entry.Name))
			b.WriteString("` | ")
			if entry.Value == "" {
				b.WriteString("")
			} else {
				b.WriteString("`")
				b.WriteString(escapeTable(entry.Value))
				b.WriteString("`")
			}
			b.WriteString(" | ")
			b.WriteString(escapeTableDetails(entry.Details))
			b.WriteString(" | ")
			b.WriteString(escapeTable("`" + strings.Join(entry.SourceEvidence, "`, `") + "`"))
			b.WriteString(" |\n")
		}
		b.WriteString("\n")
	}
}

func escapeTable(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

func escapeTableDetails(details []string) string {
	escaped := make([]string, 0, len(details))
	for _, detail := range details {
		escaped = append(escaped, escapeTable(detail))
	}

	return strings.Join(escaped, "<br/>")
}

func mustJSON(ledger RuntimeLedger) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(ledger); err != nil {
		panic(err)
	}

	return buf.String()
}

func writeFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0o644)
}

func stableID(parts ...string) string {
	raw := strings.Join(parts, "\x00")
	sum := sha256.Sum256([]byte(raw))

	return hex.EncodeToString(sum[:12])
}

func stringValue(expr ast.Expr) (string, bool) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}
	value, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}

	return value, true
}

func intValue(expr ast.Expr) (string, bool) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.INT {
		return "", false
	}

	return lit.Value, true
}

func boolValue(expr ast.Expr) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == "true"
}

func structTag(raw, key string) (string, bool) {
	value, err := strconv.Unquote(raw)
	if err != nil {
		return "", false
	}
	tag, ok := reflect.StructTag(value).Lookup(key)

	return tag, ok
}

func exprString(fset *token.FileSet, expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	var buf bytes.Buffer
	if fset == nil {
		fset = token.NewFileSet()
	}
	if err := printer.Fprint(&buf, fset, expr); err != nil {
		return ""
	}

	return strings.TrimSpace(buf.String())
}

func (x *extractor) evalString(expr ast.Expr) (string, bool) {
	switch expr := expr.(type) {
	case *ast.BasicLit:
		return stringValue(expr)
	case *ast.Ident:
		value, ok := x.constValues[expr.Name]
		return value, ok
	case *ast.SelectorExpr:
		value, ok := x.constValues[expr.Sel.Name]
		return value, ok
	case *ast.BinaryExpr:
		if expr.Op != token.ADD {
			return "", false
		}
		left, ok := x.evalString(expr.X)
		if !ok {
			return "", false
		}
		right, ok := x.evalString(expr.Y)
		if !ok {
			return "", false
		}
		return left + right, true
	case *ast.ParenExpr:
		return x.evalString(expr.X)
	default:
		return "", false
	}
}

func keyName(expr ast.Expr) string {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.SelectorExpr:
		return expr.Sel.Name
	default:
		return ""
	}
}

func selectorMethod(expr ast.Expr) string {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return ""
	}

	return sel.Sel.Name
}

func selectorReceiver(expr ast.Expr) string {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	id, ok := sel.X.(*ast.Ident)
	if !ok {
		return ""
	}

	return id.Name
}

func namedType(expr ast.Expr) string {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.SelectorExpr:
		return expr.Sel.Name
	case *ast.StarExpr:
		return namedType(expr.X)
	case *ast.ArrayType:
		return namedType(expr.Elt)
	default:
		return ""
	}
}

func receiverTypeName(expr ast.Expr) string {
	return namedType(expr)
}

func selectorField(expr ast.Expr) string {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return ""
	}

	return sel.Sel.Name
}

func (x *extractor) singleReturnString(sf sourceFile, fn *ast.FuncDecl) (string, bool) {
	expr := singleReturnExprNode(fn)
	if expr == nil {
		return "", false
	}

	return x.evalString(expr)
}

func singleReturnExpr(sf sourceFile, fn *ast.FuncDecl) string {
	expr := singleReturnExprNode(fn)
	if expr == nil {
		return ""
	}

	return exprString(sf.Fset, expr)
}

func singleReturnExprNode(fn *ast.FuncDecl) ast.Expr {
	if fn.Body == nil {
		return nil
	}
	for _, stmt := range fn.Body.List {
		ret, ok := stmt.(*ast.ReturnStmt)
		if !ok || len(ret.Results) != 1 {
			continue
		}

		return ret.Results[0]
	}

	return nil
}

func isRawPayloadEventTypeMethod(path string, fn *ast.FuncDecl) bool {
	return path == "event/event.go" && receiverType(fn) == "RawPayload"
}

func receiverType(fn *ast.FuncDecl) string {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return ""
	}

	return receiverTypeName(fn.Recv.List[0].Type)
}

func structFields(sf sourceFile) []*ast.Field {
	var fields []*ast.Field
	for _, decl := range sf.File.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}
			fields = append(fields, st.Fields.List...)
		}
	}

	return fields
}

func splitTagCSV(raw string) []string {
	if raw == "" {
		return nil
	}
	var result []string
	var current strings.Builder
	escaped := false
	for _, r := range raw {
		if escaped {
			current.WriteRune(r)
			escaped = false

			continue
		}
		if r == '\\' {
			current.WriteRune(r)
			escaped = true

			continue
		}
		if r == ',' {
			result = append(result, current.String())
			current.Reset()

			continue
		}
		current.WriteRune(r)
	}
	result = append(result, current.String())

	return result
}

func (x *extractor) goModVersion(module string) string {
	data, err := os.ReadFile(filepath.Join(x.sourceDir, "go.mod"))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == module {
			return fields[1]
		}
	}

	return ""
}

func isTestOnlyEventTopic(value string) bool {
	return strings.HasPrefix(value, "test.") ||
		strings.HasPrefix(value, "bus.") ||
		strings.HasPrefix(value, "contract.") ||
		strings.HasPrefix(value, "memory.") ||
		value == "anything" ||
		value == "any" ||
		value == "unrouted.event" ||
		value == "with.sink.x" ||
		value == "pub.only.x" ||
		strings.Contains(value, "invalid event type") ||
		strings.Contains(value, "bad/event/type")
}

func appendUnique(slice []string, values ...string) []string {
	seen := make(map[string]bool, len(slice)+len(values))
	for _, item := range slice {
		seen[item] = true
	}
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		slice = append(slice, value)
		seen[value] = true
	}

	return slice
}

func nonEmpty(values []string) []string {
	out := values[:0]
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			out = append(out, value)
		}
	}

	return out
}

func unique(values []string) []string {
	var result []string
	seen := make(map[string]bool, len(values))
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}

	return result
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}

func isHTTPVerb(value string) bool {
	switch value {
	case "get", "post", "put", "delete", "patch", "head", "options", "trace", "connect", "all":
		return true
	default:
		return false
	}
}
