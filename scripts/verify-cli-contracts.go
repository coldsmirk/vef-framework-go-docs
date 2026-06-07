package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	cliDocsPath         = "docs/advanced/cli-tools.md"
	chineseCliDocsPath  = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/advanced/cli-tools.md"
	englishIndexPath    = "docs/reference/public-api-index.md"
	chineseIndexPath    = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"
	cliEntryDisposition = "excluded:cli-implementation-export"
)

type cliPackage struct {
	pkg         string
	topLevel    int
	fields      int
	methods     int
	entries     int
	fingerprint string
	contractID  string
}

var cliPackages = []cliPackage{
	{
		pkg:         "github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd",
		topLevel:    6,
		fields:      3,
		methods:     1,
		entries:     10,
		fingerprint: "6a01b8fdcb43f6842164be353432a6dbc7849601835c454228aab6cb5ef046ef",
		contractID:  "github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd#command-contract:root-command-and-version",
	},
	{
		pkg:         "github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd/buildinfo",
		topLevel:    2,
		fields:      0,
		methods:     0,
		entries:     2,
		fingerprint: "a9f40a22aaf4f4e6313cea5a7fcd439a5dcde2d0b13f977e954753c1317ab33e",
		contractID:  "github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd/buildinfo#command-contract:build-info-generation",
	},
	{
		pkg:         "github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd/create",
		topLevel:    2,
		fields:      0,
		methods:     0,
		entries:     2,
		fingerprint: "26171a8454bd55208efc47d3ba16ce5744a971956d17bfac4972c1468619cd3b",
		contractID:  "github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd/create#command-contract:create-command-placeholder",
	},
	{
		pkg:         "github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd/modelschema",
		topLevel:    9,
		fields:      12,
		methods:     0,
		entries:     21,
		fingerprint: "19164973da27a846f72a4df3b55d320998b55e57ee2b3dc40dd7abc4868e8735",
		contractID:  "github.com/coldsmirk/vef-framework-go/cmd/vef-cli/cmd/modelschema#command-contract:model-schema-generation",
	},
}

type corpus struct {
	label   string
	path    string
	content string
}

type auditLedger struct {
	Entries []auditEntry `json:"entries"`
}

type auditEntry struct {
	ID          string   `json:"id"`
	Package     string   `json:"package"`
	Kind        string   `json:"kind"`
	Symbol      string   `json:"symbol"`
	Signature   string   `json:"signature"`
	Disposition string   `json:"disposition"`
	Coverage    []string `json:"coverage"`
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

type contractLedger struct {
	PackageReviews []contractPackageReview `json:"package_reviews"`
	Entries        []contractEntry         `json:"entries"`
}

type contractPackageReview struct {
	Package         string        `json:"package"`
	Disposition     string        `json:"disposition"`
	ReviewedSurface reviewSurface `json:"reviewed_surface"`
	Coverage        []string      `json:"coverage"`
	SourceEvidence  []string      `json:"source_evidence"`
	ContractIDs     []string      `json:"contract_ids"`
}

type reviewSurface struct {
	TopLevel    int    `json:"top_level"`
	Fields      int    `json:"fields"`
	Methods     int    `json:"methods"`
	EntryCount  int    `json:"entry_count"`
	Fingerprint string `json:"fingerprint"`
}

type contractEntry struct {
	ID             string   `json:"id"`
	Package        string   `json:"package"`
	Kind           string   `json:"kind"`
	Contract       string   `json:"contract"`
	Disposition    string   `json:"disposition"`
	Coverage       []string `json:"coverage"`
	SourceEvidence []string `json:"source_evidence"`
	TestEvidence   []string `json:"test_evidence"`
	Terms          []string `json:"terms"`
}

func main() {
	sourceDir := flag.String("source", "../vef-framework-go", "path to the VEF Framework Go source repository")
	outDir := flag.String("out", ".", "path to the VEF Framework Go docs repository")
	flag.Parse()

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)

	englishDocs := readCorpus("English CLI docs", docsRoot, cliDocsPath)
	chineseDocs := readCorpus("Chinese CLI docs", docsRoot, chineseCliDocsPath)
	englishIndex := readCorpus("English public API index", docsRoot, englishIndexPath)
	chineseIndex := readCorpus("Chinese public API index", docsRoot, chineseIndexPath)

	entries := loadCLIAuditEntries(filepath.Join(docsRoot, "scripts/api-audit-ledger.json"))
	manifestEntries := loadCLIManifestEntries(filepath.Join(docsRoot, "scripts/api-audit-manifest.json"))
	reviews, contracts := loadCLIContracts(filepath.Join(docsRoot, "scripts/api-contract-ledger.json"))
	liveEntries := loadLiveCLIEntries(sourceRoot, docsRoot)

	var failures []string
	for _, pkg := range cliPackages {
		failures = append(failures, verifySurfaceEntry("live public API inventory", pkg, liveEntries[pkg.pkg])...)
		failures = append(failures, verifySurfaceEntry("API audit manifest", pkg, manifestEntries[pkg.pkg])...)
		failures = append(failures, verifyReviewSurface(pkg, reviews[pkg.pkg])...)
		failures = append(failures, verifyAuditEntries(pkg, entries[pkg.pkg])...)
		failures = append(failures, verifyCoverage(pkg, entries[pkg.pkg], manifestEntries[pkg.pkg], reviews[pkg.pkg], contracts[pkg.contractID])...)
		for _, index := range []corpus{englishIndex, chineseIndex} {
			failures = append(failures, verifyGeneratedIndexSection(index, pkg, entries[pkg.pkg])...)
		}
	}

	for _, doc := range []corpus{englishDocs, chineseDocs} {
		failures = append(failures, verifyCLIDocumentedSurface(doc)...)
		failures = append(failures, verifyCLIContractTerms(doc, contracts)...)
	}

	failures = append(failures, verifyContractLedger(reviews, contracts, sourceRoot)...)
	failures = append(failures, verifySourceContracts(sourceRoot)...)
	failures = append(failures, runSourceTests(sourceRoot)...)

	if len(failures) > 0 {
		sort.Strings(failures)
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Println("CLI contracts verified")
}

func verifySurfaceEntry(label string, pkg cliPackage, entry manifestEntry) []string {
	var failures []string
	if entry.Package != pkg.pkg {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q want %q", label, entry.Package, pkg.pkg))
	}
	if entry.TopLevel != pkg.topLevel || entry.Fields != pkg.fields ||
		entry.Methods != pkg.methods || entry.Fingerprint != pkg.fingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s surface mismatch for %s: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			pkg.pkg,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			pkg.topLevel, pkg.fields, pkg.methods, pkg.fingerprint,
		))
	}

	return failures
}

func verifyReviewSurface(pkg cliPackage, review contractPackageReview) []string {
	var failures []string
	if review.Package != pkg.pkg {
		failures = append(failures, fmt.Sprintf("contract package review package mismatch: got %q want %q", review.Package, pkg.pkg))
	}
	if review.Disposition != "has-semantic-contracts" {
		failures = append(failures, fmt.Sprintf("contract package review disposition mismatch for %s: got %q", pkg.pkg, review.Disposition))
	}
	if review.ReviewedSurface.TopLevel != pkg.topLevel ||
		review.ReviewedSurface.Fields != pkg.fields ||
		review.ReviewedSurface.Methods != pkg.methods ||
		review.ReviewedSurface.EntryCount != pkg.entries ||
		review.ReviewedSurface.Fingerprint != pkg.fingerprint {
		failures = append(failures, fmt.Sprintf(
			"contract package review surface mismatch for %s: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
			pkg.pkg,
			review.ReviewedSurface.TopLevel,
			review.ReviewedSurface.Fields,
			review.ReviewedSurface.Methods,
			review.ReviewedSurface.EntryCount,
			review.ReviewedSurface.Fingerprint,
		))
	}
	if !sameSet(review.ContractIDs, []string{pkg.contractID}) {
		failures = append(failures, fmt.Sprintf("contract package review contract ids mismatch for %s: got %v want %v", pkg.pkg, review.ContractIDs, []string{pkg.contractID}))
	}

	return failures
}

func verifyAuditEntries(pkg cliPackage, entries []auditEntry) []string {
	var failures []string
	if len(entries) != pkg.entries {
		failures = append(failures, fmt.Sprintf("CLI audit entry count mismatch for %s: got %d want %d", pkg.pkg, len(entries), pkg.entries))
	}

	counts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != pkg.pkg {
			failures = append(failures, "non-CLI package entry passed into CLI verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate CLI audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		if entry.Symbol == "" || entry.Signature == "" {
			failures = append(failures, "CLI audit entry missing symbol/signature "+entry.ID)
		}
		if entry.Disposition != cliEntryDisposition {
			failures = append(failures, fmt.Sprintf("CLI audit entry %s disposition mismatch: got %q want %q", entry.ID, entry.Disposition, cliEntryDisposition))
		}
	}
	if counts["top"] != pkg.topLevel || counts["field"] != pkg.fields || counts["method"] != pkg.methods {
		failures = append(failures, fmt.Sprintf("CLI audit kind counts mismatch for %s: top/field/method=%d/%d/%d", pkg.pkg, counts["top"], counts["field"], counts["method"]))
	}

	return failures
}

func verifyCoverage(
	pkg cliPackage,
	entries []auditEntry,
	manifestEntry manifestEntry,
	review contractPackageReview,
	contract contractEntry,
) []string {
	var failures []string
	expected := []string{cliDocsPath}
	if manifestEntry.Tier != "cli" {
		failures = append(failures, fmt.Sprintf("manifest tier mismatch for %s: got %q want cli", pkg.pkg, manifestEntry.Tier))
	}
	if !sameSet(manifestEntry.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("manifest CLI coverage mismatch for %s: got %v want %v", pkg.pkg, manifestEntry.Coverage, expected))
	}
	if !sameSet(review.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract package review CLI coverage mismatch for %s: got %v want %v", pkg.pkg, review.Coverage, expected))
	}
	if !sameSet(contract.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract entry CLI coverage mismatch for %s: got %v want %v", pkg.pkg, contract.Coverage, expected))
	}
	for _, entry := range entries {
		if !sameSet(entry.Coverage, expected) {
			failures = append(failures, fmt.Sprintf("audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage, expected))
		}
	}

	return failures
}

func verifyGeneratedIndexSection(index corpus, pkg cliPackage, entries []auditEntry) []string {
	section := packageSection(index.content, pkg.pkg)
	if section == "" {
		return []string{index.label + " missing CLI package section for " + pkg.pkg}
	}

	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s CLI index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
		}
	}

	return failures
}

func verifyCLIDocumentedSurface(doc corpus) []string {
	var failures []string
	failures = append(failures, missingTerms(doc, []string{
		"cmd/vef-cli/**",
		"vef-cli",
		"generate-build-info",
		"generate-model-schema",
		"create",
		"--version",
		"--output",
		"--package",
		"--input",
		"--name",
		"--path",
		"--module",
		"BuildInfo = &monitor.BuildInfo",
		"AppVersion",
		"BuildTime",
		"GitCommit",
		"git describe --tags --always --dirty",
		"none",
		"orm.BaseModel",
		"bun:\"embed:prefix_\"",
		"bun:\",scanonly\"",
		"Columns()",
		"raw=true",
		"ColTable",
		"__",
	})...)

	for _, pkg := range cliPackages {
		failures = append(failures, missingTerms(doc, []string{
			pkg.pkg,
			strconv.Itoa(pkg.entries),
			pkg.fingerprint,
		})...)
	}

	alternatives := map[string][]string{
		"CLI import boundary": {
			"instead of importing command implementation packages",
			"不应\nimport 命令实现包",
			"不是受支持的 import API",
		},
		"create placeholder": {
			"not implemented",
			"还没有实现",
			"not-implemented error",
		},
		"version dirty suffix": {
			"-dirty",
		},
	}
	for label, terms := range alternatives {
		if !containsAnyTerm(doc.content, terms) {
			failures = append(failures, fmt.Sprintf("%s missing %s term: one of %v", doc.label, label, terms))
		}
	}

	return failures
}

func verifyCLIContractTerms(doc corpus, contracts map[string]contractEntry) []string {
	var failures []string
	for _, pkg := range cliPackages {
		contract := contracts[pkg.contractID]
		for _, term := range contract.Terms {
			if !strings.Contains(doc.content, term) {
				failures = append(failures, fmt.Sprintf("%s missing CLI contract term for %s: %s", doc.label, contract.ID, term))
			}
		}
	}

	return failures
}

func verifyContractLedger(reviews map[string]contractPackageReview, contracts map[string]contractEntry, sourceRoot string) []string {
	var failures []string
	for _, pkg := range cliPackages {
		review := reviews[pkg.pkg]
		contract := contracts[pkg.contractID]
		if contract.ID != pkg.contractID {
			failures = append(failures, fmt.Sprintf("CLI contract id mismatch for %s: got %q", pkg.pkg, contract.ID))
		}
		if contract.Package != pkg.pkg {
			failures = append(failures, fmt.Sprintf("CLI contract package mismatch for %s: got %q", pkg.pkg, contract.Package))
		}
		if contract.Kind != "command-contract" {
			failures = append(failures, fmt.Sprintf("CLI contract kind mismatch for %s: got %q", pkg.pkg, contract.Kind))
		}
		if contract.Disposition != "documented:semantic-contract" {
			failures = append(failures, fmt.Sprintf("CLI contract disposition mismatch for %s: got %q", pkg.pkg, contract.Disposition))
		}
		if len(contract.Terms) == 0 {
			failures = append(failures, "CLI contract terms empty for "+contract.ID)
		}

		allEvidence := append([]string{}, review.SourceEvidence...)
		allEvidence = append(allEvidence, contract.SourceEvidence...)
		allEvidence = append(allEvidence, contract.TestEvidence...)
		for _, item := range allEvidence {
			path, lineText, ok := strings.Cut(item, ":")
			if !ok || lineText == "" {
				failures = append(failures, "CLI contract evidence missing line number: "+item)
				continue
			}
			if _, err := os.Stat(filepath.Join(sourceRoot, path)); err != nil {
				failures = append(failures, "CLI contract evidence missing file: "+item)
			}
		}
	}

	return failures
}

func verifySourceContracts(sourceRoot string) []string {
	files := map[string]string{
		"cmd/vef-cli/cmd/root.go":                       readSourceFile(sourceRoot, "cmd/vef-cli/cmd/root.go").content,
		"cmd/vef-cli/cmd/version.go":                    readSourceFile(sourceRoot, "cmd/vef-cli/cmd/version.go").content,
		"cmd/vef-cli/cmd/version_test.go":               readSourceFile(sourceRoot, "cmd/vef-cli/cmd/version_test.go").content,
		"cmd/vef-cli/cmd/buildinfo/command.go":          readSourceFile(sourceRoot, "cmd/vef-cli/cmd/buildinfo/command.go").content,
		"cmd/vef-cli/cmd/buildinfo/generator.go":        readSourceFile(sourceRoot, "cmd/vef-cli/cmd/buildinfo/generator.go").content,
		"cmd/vef-cli/cmd/buildinfo/generator_test.go":   readSourceFile(sourceRoot, "cmd/vef-cli/cmd/buildinfo/generator_test.go").content,
		"cmd/vef-cli/cmd/create/command.go":             readSourceFile(sourceRoot, "cmd/vef-cli/cmd/create/command.go").content,
		"cmd/vef-cli/cmd/modelschema/command.go":        readSourceFile(sourceRoot, "cmd/vef-cli/cmd/modelschema/command.go").content,
		"cmd/vef-cli/cmd/modelschema/generator.go":      readSourceFile(sourceRoot, "cmd/vef-cli/cmd/modelschema/generator.go").content,
		"cmd/vef-cli/cmd/modelschema/generator_test.go": readSourceFile(sourceRoot, "cmd/vef-cli/cmd/modelschema/generator_test.go").content,
	}

	checks := []struct {
		file string
		term string
	}{
		{"cmd/vef-cli/cmd/root.go", `Use:   "vef-cli"`},
		{"cmd/vef-cli/cmd/root.go", "rootCmd.SetVersionTemplate(Banner"},
		{"cmd/vef-cli/cmd/root.go", "create.Command()"},
		{"cmd/vef-cli/cmd/root.go", "buildinfo.Command()"},
		{"cmd/vef-cli/cmd/root.go", "modelschema.Command()"},
		{"cmd/vef-cli/cmd/version.go", "type VersionInfo struct"},
		{"cmd/vef-cli/cmd/version.go", "Version string"},
		{"cmd/vef-cli/cmd/version.go", "Date    string"},
		{"cmd/vef-cli/cmd/version.go", "Dirty   bool"},
		{"cmd/vef-cli/cmd/version.go", `version += "-dirty"`},
		{"cmd/vef-cli/cmd/version.go", `return fmt.Sprintf("Version: %s | Built: %s", version, v.Date)`},
		{"cmd/vef-cli/cmd/version_test.go", "TestVersionInfoString"},
		{"cmd/vef-cli/cmd/version_test.go", "DirtyAppendsSuffix"},
		{"cmd/vef-cli/cmd/version_test.go", "TestMergeBuildInfo"},
		{"cmd/vef-cli/cmd/buildinfo/command.go", `Use:   "generate-build-info"`},
		{"cmd/vef-cli/cmd/buildinfo/command.go", `cmd.Flags().StringP("output", "o", "build_info.go"`},
		{"cmd/vef-cli/cmd/buildinfo/command.go", `cmd.Flags().StringP("package", "p", "main"`},
		{"cmd/vef-cli/cmd/buildinfo/generator.go", "var BuildInfo = &monitor.BuildInfo"},
		{"cmd/vef-cli/cmd/buildinfo/generator.go", "AppVersion: getVersion(ctx)"},
		{"cmd/vef-cli/cmd/buildinfo/generator.go", "BuildTime:  timex.Now().String()"},
		{"cmd/vef-cli/cmd/buildinfo/generator.go", "GitCommit:  getCommit(ctx)"},
		{"cmd/vef-cli/cmd/buildinfo/generator.go", "os.MkdirAll(dir, 0o755)"},
		{"cmd/vef-cli/cmd/buildinfo/generator.go", `"git", "describe", "--tags", "--always", "--dirty"`},
		{"cmd/vef-cli/cmd/buildinfo/generator.go", `return "dev"`},
		{"cmd/vef-cli/cmd/buildinfo/generator.go", `return "none"`},
		{"cmd/vef-cli/cmd/buildinfo/generator_test.go", "WritesParseableFileInNestedDir"},
		{"cmd/vef-cli/cmd/buildinfo/generator_test.go", "monitor.BuildInfo"},
		{"cmd/vef-cli/cmd/create/command.go", `Use:   "create"`},
		{"cmd/vef-cli/cmd/create/command.go", "ErrNotImplemented"},
		{"cmd/vef-cli/cmd/create/command.go", "vef-cli create is not implemented yet, please generate the project manually"},
		{"cmd/vef-cli/cmd/create/command.go", `cmd.Flags().StringP("name", "n", "", "Project name (required)")`},
		{"cmd/vef-cli/cmd/create/command.go", `cmd.Flags().StringP("path", "p", ".", "Directory path to create the project")`},
		{"cmd/vef-cli/cmd/create/command.go", `cmd.Flags().StringP("module", "m", "", "Go module name`},
		{"cmd/vef-cli/cmd/create/command.go", `cmd.MarkFlagRequired("name")`},
		{"cmd/vef-cli/cmd/modelschema/command.go", `Use:   "generate-model-schema"`},
		{"cmd/vef-cli/cmd/modelschema/command.go", `cmd.Flags().StringP("input", "i", "", "Input model file or directory path")`},
		{"cmd/vef-cli/cmd/modelschema/command.go", `cmd.Flags().StringP("output", "o", "", "Output schema file or directory path")`},
		{"cmd/vef-cli/cmd/modelschema/command.go", `cmd.Flags().StringP("package", "p", "schemas"`},
		{"cmd/vef-cli/cmd/modelschema/command.go", "errInputOutputMismatch"},
		{"cmd/vef-cli/cmd/modelschema/generator.go", "ErrNoGoFilesFound"},
		{"cmd/vef-cli/cmd/modelschema/generator.go", "ErrNoPackagesFound"},
		{"cmd/vef-cli/cmd/modelschema/generator.go", "ErrMultiplePackages"},
		{"cmd/vef-cli/cmd/modelschema/generator.go", "ErrFileNotFoundInPackage"},
		{"cmd/vef-cli/cmd/modelschema/generator.go", "func GenerateFile(inputFile, outputFile, packageName string) error"},
		{"cmd/vef-cli/cmd/modelschema/generator.go", "func GenerateDirectory(inputDir, outputDir, packageName string) error"},
		{"cmd/vef-cli/cmd/modelschema/generator.go", `filepath.Glob(filepath.Join(inputDir, "*.go"))`},
		{"cmd/vef-cli/cmd/modelschema/generator.go", "isOrmBaseModel"},
		{"cmd/vef-cli/cmd/modelschema/generator.go", "parseBunTag"},
		{"cmd/vef-cli/cmd/modelschema/generator.go", "extractEmbedPrefixFromTag"},
		{"cmd/vef-cli/cmd/modelschema/generator.go", "hasScanonlyTagFromTag"},
		{"cmd/vef-cli/cmd/modelschema/generator.go", "reservedMethodNames"},
		{"cmd/vef-cli/cmd/modelschema/generator.go", `strings.HasPrefix(strings.TrimSpace(part), "rel:")`},
		{"cmd/vef-cli/cmd/modelschema/generator_test.go", "SchemaTypeIsUnexportedAndOnlyVarIsExported"},
		{"cmd/vef-cli/cmd/modelschema/generator_test.go", "ScanonlyExcludedFromColumns"},
		{"cmd/vef-cli/cmd/modelschema/generator_test.go", "EmbedPrefixApplied"},
		{"cmd/vef-cli/cmd/modelschema/generator_test.go", "ReservedMethodNameEscaped"},
	}

	var failures []string
	for _, check := range checks {
		if !strings.Contains(files[check.file], check.term) {
			failures = append(failures, fmt.Sprintf("%s missing source contract term: %s", check.file, check.term))
		}
	}

	return failures
}

func runSourceTests(sourceRoot string) []string {
	return runCommand(sourceRoot, "go", "test", "./cmd/vef-cli/cmd/...")
}

func loadCLIAuditEntries(path string) map[string][]auditEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var ledger auditLedger
	if err := json.Unmarshal(data, &ledger); err != nil {
		panic(err)
	}

	result := map[string][]auditEntry{}
	for _, entry := range ledger.Entries {
		if _, ok := cliPackageByName(entry.Package); ok {
			result[entry.Package] = append(result[entry.Package], entry)
		}
	}
	for pkg := range result {
		sort.Slice(result[pkg], func(i, j int) bool {
			return result[pkg][i].ID < result[pkg][j].ID
		})
	}

	return result
}

func loadCLIManifestEntries(path string) map[string]manifestEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var m manifest
	if err := json.Unmarshal(data, &m); err != nil {
		panic(err)
	}

	result := map[string]manifestEntry{}
	for _, entry := range m.Packages {
		if _, ok := cliPackageByName(entry.Package); ok {
			result[entry.Package] = entry
		}
	}

	return result
}

func loadCLIContracts(path string) (map[string]contractPackageReview, map[string]contractEntry) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var ledger contractLedger
	if err := json.Unmarshal(data, &ledger); err != nil {
		panic(err)
	}

	reviews := map[string]contractPackageReview{}
	for _, review := range ledger.PackageReviews {
		if _, ok := cliPackageByName(review.Package); ok {
			reviews[review.Package] = review
		}
	}

	contracts := map[string]contractEntry{}
	for _, item := range ledger.Entries {
		for _, pkg := range cliPackages {
			if item.ID == pkg.contractID {
				contracts[item.ID] = item
			}
		}
	}

	return reviews, contracts
}

func loadLiveCLIEntries(sourceRoot, docsRoot string) map[string]manifestEntry {
	cmd := exec.Command(
		"go", "run", filepath.Join(docsRoot, "scripts/verify-api-audit.go"),
		"-source", sourceRoot,
		"-manifest", filepath.Join(docsRoot, "scripts/api-audit-manifest.json"),
		"-ledger", filepath.Join(docsRoot, "scripts/api-audit-ledger.json"),
		"-contract-ledger", filepath.Join(docsRoot, "scripts/api-contract-ledger.json"),
		"-print-current",
	)
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Errorf("verify-api-audit -print-current failed: %w\n%s", err, strings.TrimSpace(string(output))))
	}

	var entries []manifestEntry
	payload := "[" + strings.TrimSpace(string(output)) + "]"
	if err := json.Unmarshal([]byte(payload), &entries); err != nil {
		panic(fmt.Errorf("failed to parse verify-api-audit -print-current output: %w", err))
	}

	result := map[string]manifestEntry{}
	for _, entry := range entries {
		if _, ok := cliPackageByName(entry.Package); ok {
			result[entry.Package] = entry
		}
	}

	return result
}

func cliPackageByName(name string) (cliPackage, bool) {
	for _, pkg := range cliPackages {
		if pkg.pkg == name {
			return pkg, true
		}
	}

	return cliPackage{}, false
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

func missingTerms(doc corpus, terms []string) []string {
	var failures []string
	for _, term := range terms {
		if !strings.Contains(doc.content, term) {
			failures = append(failures, fmt.Sprintf("%s missing term: %s", doc.label, term))
		}
	}

	return failures
}

func containsAnyTerm(content string, terms []string) bool {
	for _, term := range terms {
		if strings.Contains(content, term) {
			return true
		}
	}

	return false
}

func readCorpus(label, root, relPath string) corpus {
	content, err := os.ReadFile(filepath.Join(root, relPath))
	if err != nil {
		panic(err)
	}

	return corpus{label: label, path: relPath, content: string(content)}
}

func readSourceFile(sourceRoot, relPath string) corpus {
	content, err := os.ReadFile(filepath.Join(sourceRoot, relPath))
	if err != nil {
		panic(err)
	}

	return corpus{label: relPath, path: relPath, content: string(content)}
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}

func sameSet(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	gotCopy := append([]string(nil), got...)
	wantCopy := append([]string(nil), want...)
	sort.Strings(gotCopy)
	sort.Strings(wantCopy)
	for i := range gotCopy {
		if gotCopy[i] != wantCopy[i] {
			return false
		}
	}

	return true
}

func runCommand(dir, name string, args ...string) []string {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return []string{fmt.Sprintf("%s failed: %v\n%s", strings.Join(append([]string{name}, args...), " "), err, strings.TrimSpace(output.String()))}
	}

	return nil
}
