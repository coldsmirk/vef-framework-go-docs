package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	datasourcePackage     = "github.com/coldsmirk/vef-framework-go/datasource"
	datasourceFingerprint = "a8d1f60b94e7300151d3df0025eec3b3e387d732829ecfff0ecaf7a660ba3cc3"
	datasourceTopLevel    = 17
	datasourceFields      = 9
	datasourceMethods     = 13
	datasourceEntries     = 39

	englishExtensionPath = "docs/reference/extension-points.md"
	chineseExtensionPath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/extension-points.md"
	englishIndexPath     = "docs/reference/public-api-index.md"
	chineseIndexPath     = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"
)

type corpus struct {
	label   string
	path    string
	content string
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

	englishExtension := readCorpus("English extension-points datasource docs", docsRoot, englishExtensionPath)
	chineseExtension := readCorpus("Chinese extension-points datasource docs", docsRoot, chineseExtensionPath)
	englishIndex := readCorpus("English public API index", docsRoot, englishIndexPath)
	chineseIndex := readCorpus("Chinese public API index", docsRoot, chineseIndexPath)

	entries := loadDatasourceAuditEntries(filepath.Join(docsRoot, "scripts/api-audit-ledger.json"))
	manifestEntry := loadDatasourceManifestEntry(filepath.Join(docsRoot, "scripts/api-audit-manifest.json"))
	review, contract := loadDatasourceContract(filepath.Join(docsRoot, "scripts/api-contract-ledger.json"))
	liveEntry := loadLiveDatasourceEntry(sourceRoot, docsRoot)

	var failures []string
	failures = append(failures, verifySurfaceEntry("live public API inventory", liveEntry)...)
	failures = append(failures, verifySurfaceEntry("API audit manifest", manifestEntry)...)
	failures = append(failures, verifyReviewSurface(review)...)
	failures = append(failures, verifyAuditEntries(entries)...)
	failures = append(failures, verifyCoverage(entries, manifestEntry, review, contract)...)

	for _, index := range []corpus{englishIndex, chineseIndex} {
		failures = append(failures, verifyGeneratedIndexSection(index, entries)...)
	}
	for _, doc := range []corpus{englishExtension, chineseExtension} {
		failures = append(failures, verifyDocumentedSurface(doc, entries)...)
		failures = append(failures, verifyNoPhantomDatasourceRefs(doc, entries)...)
		failures = append(failures, missingTerms(doc, extensionTerms(doc.label))...)
	}

	failures = append(failures, verifyContractLedger(review, contract, sourceRoot)...)
	failures = append(failures, verifySourceContracts(sourceRoot)...)
	failures = append(failures, runSourceTests(sourceRoot)...)

	if len(failures) > 0 {
		sort.Strings(failures)
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Println("datasource contracts verified")
}

func verifySurfaceEntry(label string, entry manifestEntry) []string {
	var failures []string
	if entry.Package != datasourcePackage {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q", label, entry.Package))
	}
	if entry.TopLevel != datasourceTopLevel || entry.Fields != datasourceFields ||
		entry.Methods != datasourceMethods || entry.Fingerprint != datasourceFingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s surface mismatch: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			datasourceTopLevel, datasourceFields, datasourceMethods, datasourceFingerprint,
		))
	}

	return failures
}

func verifyReviewSurface(review contractPackageReview) []string {
	var failures []string
	if review.Package != datasourcePackage {
		failures = append(failures, fmt.Sprintf("contract package review package mismatch: got %q", review.Package))
	}
	if review.Disposition != "has-semantic-contracts" {
		failures = append(failures, fmt.Sprintf("contract package review disposition mismatch: got %q", review.Disposition))
	}
	if review.ReviewedSurface.TopLevel != datasourceTopLevel ||
		review.ReviewedSurface.Fields != datasourceFields ||
		review.ReviewedSurface.Methods != datasourceMethods ||
		review.ReviewedSurface.EntryCount != datasourceEntries ||
		review.ReviewedSurface.Fingerprint != datasourceFingerprint {
		failures = append(failures, fmt.Sprintf(
			"contract package review surface mismatch: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
			review.ReviewedSurface.TopLevel,
			review.ReviewedSurface.Fields,
			review.ReviewedSurface.Methods,
			review.ReviewedSurface.EntryCount,
			review.ReviewedSurface.Fingerprint,
		))
	}
	if !contains(review.ContractIDs, datasourceContractID()) {
		failures = append(failures, "contract package review missing datasource contract id")
	}

	return failures
}

func verifyAuditEntries(entries []auditEntry) []string {
	var failures []string
	if len(entries) != datasourceEntries {
		failures = append(failures, fmt.Sprintf("datasource audit entry count mismatch: got %d want %d", len(entries), datasourceEntries))
	}

	counts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != datasourcePackage {
			failures = append(failures, "non-datasource audit entry passed into datasource verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate datasource audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		if entry.Symbol == "" || entry.Signature == "" {
			failures = append(failures, "datasource audit entry missing symbol/signature "+entry.ID)
		}
	}
	if counts["top"] != datasourceTopLevel || counts["field"] != datasourceFields || counts["method"] != datasourceMethods {
		failures = append(failures, fmt.Sprintf("datasource audit kind counts mismatch: top/field/method=%d/%d/%d", counts["top"], counts["field"], counts["method"]))
	}

	return failures
}

func verifyCoverage(
	entries []auditEntry,
	manifestEntry manifestEntry,
	review contractPackageReview,
	contract contractEntry,
) []string {
	var failures []string
	expected := []string{englishExtensionPath}
	if !sameSet(manifestEntry.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("manifest datasource coverage mismatch: got %v want %v", manifestEntry.Coverage, expected))
	}
	if !sameSet(review.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract package review datasource coverage mismatch: got %v want %v", review.Coverage, expected))
	}
	if !sameSet(contract.Coverage, expected) {
		failures = append(failures, fmt.Sprintf("contract entry datasource coverage mismatch: got %v want %v", contract.Coverage, expected))
	}
	for _, entry := range entries {
		if !sameSet(entry.Coverage, expected) {
			failures = append(failures, fmt.Sprintf("audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage, expected))
		}
	}

	return failures
}

func verifyGeneratedIndexSection(index corpus, entries []auditEntry) []string {
	section := packageSection(index.content, datasourcePackage)
	if section == "" {
		return []string{index.label + " missing datasource package section"}
	}

	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s datasource index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
		}
	}

	return failures
}

func verifyDocumentedSurface(doc corpus, entries []auditEntry) []string {
	var failures []string
	for _, entry := range entries {
		ref := "datasource." + entry.Symbol
		if !hasCodeReference(doc.content, ref) {
			failures = append(failures, fmt.Sprintf("%s missing audited datasource entry `%s`", doc.label, ref))
		}
	}

	return failures
}

func verifyNoPhantomDatasourceRefs(doc corpus, entries []auditEntry) []string {
	valid := map[string]bool{}
	for _, entry := range entries {
		valid["datasource."+entry.Symbol] = true
	}

	var failures []string
	refs := datasourceReferencesFromMarkdownCode(doc.content)
	for _, ref := range refs {
		if valid[ref] {
			continue
		}
		failures = append(failures, fmt.Sprintf("%s references unknown datasource public API: %s", doc.label, ref))
	}

	return failures
}

func hasCodeReference(content, ref string) bool {
	return strings.Contains(content, "`"+ref+"`") ||
		strings.Contains(content, "`"+ref+"(") ||
		strings.Contains(content, "`"+ref+"()")
}

func extensionTerms(label string) []string {
	common := []string{
		datasourceFingerprint,
		"`vef.data_sources.primary`",
		"`orm.DB`",
		"datasource.Provider.Name",
		"datasource.Provider.Load",
		"datasource.ConnectionInfo.Version",
		"datasource.ReconcileOptions.DryRun",
		"datasource.ReconcileReport.Added",
		"datasource.ReconcileReport.Updated",
		"datasource.ReconcileReport.Removed",
		"datasource.ReconcileReport.Errors",
		"datasource.RegisterOptions.CloseGrace",
		"datasource.Spec.Name",
		"datasource.Spec.Config",
		"datasource.Registry.Primary",
		"datasource.Registry.Get",
		"datasource.Registry.Has",
		"datasource.Registry.Names",
		"datasource.Registry.Kind",
		"datasource.Registry.Register",
		"datasource.Registry.Update",
		"datasource.Registry.Unregister",
		"datasource.Registry.Reconcile",
		"datasource.Registry.TestConnection",
		"datasource.Registry.HealthCheck",
		"lexical order",
		"partial failure",
		"5s timeout",
	}
	if strings.HasPrefix(label, "Chinese") {
		return append(common,
			"17 个 top-level symbols",
			"9 个 exported struct fields",
			"13 个 exported methods",
			"provider 顺序未定义",
		)
	}

	return append(common,
		"17 top-level symbols",
		"9 exported struct fields",
		"13 exported methods",
		"Provider order is undefined",
	)
}

func verifyContractLedger(review contractPackageReview, contract contractEntry, sourceRoot string) []string {
	var failures []string
	if contract.ID != datasourceContractID() {
		failures = append(failures, fmt.Sprintf("datasource contract entry id mismatch: got %q", contract.ID))
	}
	if contract.Package != datasourcePackage {
		failures = append(failures, fmt.Sprintf("datasource contract entry package mismatch: got %q", contract.Package))
	}
	if contract.Kind != "runtime-contract" {
		failures = append(failures, fmt.Sprintf("datasource contract entry kind mismatch: got %q", contract.Kind))
	}
	if contract.Contract != "dynamic datasource provider and registry mutation contract" {
		failures = append(failures, fmt.Sprintf("datasource contract label mismatch: got %q", contract.Contract))
	}

	expectedEvidence := datasourceSourceEvidence(sourceRoot)
	if !sameSet(review.SourceEvidence, expectedEvidence) {
		failures = append(failures, fmt.Sprintf("datasource package review source evidence mismatch: got %v want %v", review.SourceEvidence, expectedEvidence))
	}
	if !sameSet(contract.SourceEvidence, expectedEvidence) {
		failures = append(failures, fmt.Sprintf("datasource contract source evidence mismatch: got %v want %v", contract.SourceEvidence, expectedEvidence))
	}
	expectedTestEvidence := []string{
		"internal/datasource/registry_test.go:22",
		"internal/datasource/module_test.go:25",
	}
	if !sameSet(contract.TestEvidence, expectedTestEvidence) {
		failures = append(failures, fmt.Sprintf("datasource contract test evidence mismatch: got %v want %v", contract.TestEvidence, expectedTestEvidence))
	}
	for _, term := range []string{
		"datasource.PrimaryName",
		"datasource.Provider",
		"datasource.Provider.Load",
		"datasource.Registry.Register",
		"datasource.Registry.Update",
		"datasource.Registry.Unregister",
		"datasource.Registry.Reconcile",
		"datasource.Registry.TestConnection",
		"datasource.Registry.HealthCheck",
		"datasource.WithCloseGrace(d)",
		"datasource.WithReconcileDryRun()",
		"datasource.ErrNotFound",
		"datasource.ErrExists",
		"datasource.ErrPrimaryReserved",
		"datasource.ErrNameInvalid",
		"datasource.ErrClosed",
		"5s timeout",
		"lexical order",
		"partial failure",
	} {
		if !contains(contract.Terms, term) {
			failures = append(failures, "datasource contract missing term "+term)
		}
	}

	return failures
}

func verifySourceContracts(sourceRoot string) []string {
	checks := []struct {
		path  string
		terms []string
	}{
		{
			path: "datasource/registry.go",
			terms: []string{
				"const PrimaryName = config.PrimaryDataSourceName",
				"type Registry interface",
				"Primary() orm.DB",
				"Get(name string) (orm.DB, error)",
				"Has(name string) bool",
				"Names() []string",
				"Kind(name string) (config.DBKind, error)",
				"Register(ctx context.Context, name string, cfg config.DataSourceConfig) (orm.DB, error)",
				"Update(ctx context.Context, name string, cfg config.DataSourceConfig, opts ...RegisterOption) (orm.DB, error)",
				"Unregister(ctx context.Context, name string, opts ...RegisterOption) error",
				"Reconcile(ctx context.Context, specs []Spec, opts ...ReconcileOption) (ReconcileReport, error)",
				"TestConnection(ctx context.Context, cfg config.DataSourceConfig) (ConnectionInfo, error)",
				"HealthCheck(ctx context.Context) map[string]error",
				"type ConnectionInfo struct",
				"Version string",
			},
		},
		{
			path: "datasource/provider.go",
			terms: []string{
				"type Provider interface",
				"Name() string",
				"Load(ctx context.Context) ([]Spec, error)",
				"type Spec struct",
				"Name string",
				"Config config.DataSourceConfig",
			},
		},
		{
			path: "datasource/options.go",
			terms: []string{
				"type RegisterOption func(*RegisterOptions)",
				"type RegisterOptions struct",
				"CloseGrace time.Duration",
				"func WithCloseGrace(d time.Duration) RegisterOption",
				"if d > 0",
				"type ReconcileOption func(*ReconcileOptions)",
				"type ReconcileOptions struct",
				"DryRun bool",
				"func WithReconcileDryRun() ReconcileOption",
				"type ReconcileReport struct",
				"Added   []string",
				"Updated []string",
				"Removed []string",
				"Errors  map[string]error",
			},
		},
		{
			path: "datasource/errors.go",
			terms: []string{
				"ErrNotFound = errors.New(\"datasource: data source not found\")",
				"ErrExists = errors.New(\"datasource: data source already registered\")",
				"ErrPrimaryReserved = errors.New(\"datasource: primary data source is reserved\")",
				"ErrNameInvalid = errors.New(\"datasource: data source name invalid\")",
				"ErrClosed = errors.New(\"datasource: registry is closed\")",
			},
		},
		{
			path: "internal/datasource/module.go",
			terms: []string{
				"ProviderParams",
				"Providers []datasource.Provider `group:\"vef:datasource:providers\"`",
				"reg.PrimaryRawDB().PingContext(ctx)",
				"seedStatic(ctx, reg, cfg)",
				"runProviders(ctx, reg, params.Providers)",
				"if name == datasource.PrimaryName",
				"provider.Load(ctx)",
				"fmt.Errorf(\"data source provider %q: %w\", provider.Name(), err)",
				"registerSource(ctx, reg, spec.Name, spec.Config, provider.Name())",
			},
		},
		{
			path: "internal/datasource/registry.go",
			terms: []string{
				"func (r *registry) Names() []string",
				"slices.Sort(out)",
				"func (r *registry) Register(ctx context.Context, name string, cfg config.DataSourceConfig) (orm.DB, error)",
				"PutIfAbsent(name, e)",
				"return nil, datasource.ErrExists",
				"func (r *registry) Update(ctx context.Context, name string, cfg config.DataSourceConfig, opts ...datasource.RegisterOption) (orm.DB, error)",
				"return nil, datasource.ErrNotFound",
				"r.asyncClose(name, oldEntry.rawDB, applyOptions(opts))",
				"func (r *registry) Unregister(_ context.Context, name string, opts ...datasource.RegisterOption) error",
				"return datasource.ErrPrimaryReserved",
				"func (r *registry) Reconcile(ctx context.Context, specs []datasource.Spec, opts ...datasource.ReconcileOption) (datasource.ReconcileReport, error)",
				"r.reconcileMu.Lock()",
				"if s.Name == \"\" || s.Name == datasource.PrimaryName",
				"if ro.DryRun",
				"report.Errors = errs",
				"const defaultTestConnectionTimeout = 5 * time.Second",
				"context.WithTimeout(ctx, defaultTestConnectionTimeout)",
				"database.Version(ctx, cfg.Kind, db)",
				"func (r *registry) HealthCheck(ctx context.Context) map[string]error",
				"r.primary.rawDB.PingContext(ctx)",
				"func validateName(name string) error",
				"unicode.IsSpace(c) || unicode.IsControl(c)",
				"func diffReconcile(current, desired map[string]config.DataSourceConfig) (adds, updates, removes []string)",
				"slices.Sort(adds)",
				"slices.Sort(updates)",
				"slices.Sort(removes)",
			},
		},
		{
			path: "internal/datasource/registry_test.go",
			terms: []string{
				"TestRegistryTestConnection",
				"TestRegistryHealthCheck",
				"TestRegistryReconcileDryRun",
				"TestRegistryReconcilePartialFailureAggregatesErrors",
				"TestRegistryMutationsRejectedAfterShutdown",
				"TestRegistryUpdateWithCloseGrace",
				"TestRegistryUnregisterDrainsInFlight",
			},
		},
		{
			path: "internal/datasource/module_test.go",
			terms: []string{
				"TestProvideRegistryClosesPrimaryOnSeedFailure",
				"TestProvideRegistryDefersPrimaryPingToStart",
				"TestRunProviders",
				"Register duplicate spec name",
			},
		},
	}

	var failures []string
	for _, check := range checks {
		source := readSourceFile(sourceRoot, check.path)
		failures = append(failures, missingTerms(source, check.terms)...)
	}

	return failures
}

func runSourceTests(sourceRoot string) []string {
	return runCommand(sourceRoot, "go", "test", "./datasource", "./internal/datasource")
}

func loadDatasourceAuditEntries(path string) []auditEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var ledger auditLedger
	if err := json.Unmarshal(data, &ledger); err != nil {
		panic(err)
	}

	var entries []auditEntry
	for _, entry := range ledger.Entries {
		if entry.Package == datasourcePackage {
			entries = append(entries, entry)
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Kind != entries[j].Kind {
			return entries[i].Kind < entries[j].Kind
		}

		return entries[i].Symbol < entries[j].Symbol
	})

	return entries
}

func loadDatasourceManifestEntry(path string) manifestEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var m manifest
	if err := json.Unmarshal(data, &m); err != nil {
		panic(err)
	}
	for _, entry := range m.Packages {
		if entry.Package == datasourcePackage {
			return entry
		}
	}

	panic("API audit manifest missing datasource package")
}

func loadDatasourceContract(path string) (contractPackageReview, contractEntry) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var ledger contractLedger
	if err := json.Unmarshal(data, &ledger); err != nil {
		panic(err)
	}

	var review contractPackageReview
	reviewFound := false
	for _, item := range ledger.PackageReviews {
		if item.Package == datasourcePackage {
			review = item
			reviewFound = true

			break
		}
	}
	if !reviewFound {
		panic("API contract ledger missing datasource package review")
	}

	var contract contractEntry
	contractFound := false
	for _, item := range ledger.Entries {
		if item.ID == datasourceContractID() {
			contract = item
			contractFound = true

			break
		}
	}
	if !contractFound {
		panic("API contract ledger missing datasource contract entry")
	}

	return review, contract
}

func loadLiveDatasourceEntry(sourceRoot, docsRoot string) manifestEntry {
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
	for _, entry := range entries {
		if entry.Package == datasourcePackage {
			return entry
		}
	}

	panic("live API inventory missing datasource package")
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

func datasourceReferencesFromMarkdownCode(content string) []string {
	codeParts := markdownCodeParts(content)
	re := regexp.MustCompile(`datasource\.[A-Z][A-Za-z0-9_]*(?:\.[A-Z][A-Za-z0-9_]*)?`)
	seen := map[string]bool{}
	for _, part := range codeParts {
		for _, ref := range re.FindAllString(part, -1) {
			seen[ref] = true
		}
	}

	refs := make([]string, 0, len(seen))
	for ref := range seen {
		refs = append(refs, ref)
	}
	sort.Strings(refs)

	return refs
}

func markdownCodeParts(content string) []string {
	var parts []string
	lines := strings.Split(content, "\n")
	inFence := false
	var fence strings.Builder

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			if inFence {
				parts = append(parts, fence.String())
				fence.Reset()
				inFence = false
			} else {
				inFence = true
			}

			continue
		}

		if inFence {
			fence.WriteString(line)
			fence.WriteByte('\n')

			continue
		}

		parts = append(parts, inlineCodeParts(line)...)
	}

	return parts
}

func inlineCodeParts(line string) []string {
	var parts []string
	for {
		start := strings.IndexByte(line, '`')
		if start < 0 {
			return parts
		}
		rest := line[start+1:]
		end := strings.IndexByte(rest, '`')
		if end < 0 {
			return parts
		}
		parts = append(parts, rest[:end])
		line = rest[end+1:]
	}
}

func datasourceSourceEvidence(sourceRoot string) []string {
	evidence := []string{
		"datasource/registry.go:14",
		"datasource/provider.go:7",
		"datasource/options.go:8",
		"datasource/errors.go:5",
		"internal/datasource/registry.go:53",
		"internal/datasource/module.go:38",
	}
	for _, item := range evidence {
		path, _, _ := strings.Cut(item, ":")
		if _, err := os.Stat(filepath.Join(sourceRoot, path)); err != nil {
			panic(err)
		}
	}

	return evidence
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

func datasourceContractID() string {
	return datasourcePackage + "#runtime-contract:dynamic-datasource-registry"
}

func contains(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}

	return false
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
