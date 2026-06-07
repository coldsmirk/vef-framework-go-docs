package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Manifest struct {
	SourceModule string                 `json:"source_module"`
	Packages     []ManifestPackageEntry `json:"packages"`
}

type AuditLedger struct {
	SourceModule string             `json:"source_module"`
	Scope        string             `json:"scope"`
	EntryCount   int                `json:"entry_count"`
	Fingerprint  string             `json:"fingerprint"`
	Dispositions map[string]string  `json:"dispositions"`
	Entries      []AuditLedgerEntry `json:"entries"`
}

type AuditLedgerEntry struct {
	ID          string   `json:"id"`
	Package     string   `json:"package"`
	Kind        string   `json:"kind"`
	Symbol      string   `json:"symbol"`
	Signature   string   `json:"signature"`
	Disposition string   `json:"disposition"`
	Coverage    []string `json:"coverage"`
}

type ContractLedger struct {
	SourceModule   string                  `json:"source_module"`
	Scope          string                  `json:"scope"`
	PackageReviews []ContractPackageReview `json:"package_reviews"`
	Entries        []ContractLedgerEntry   `json:"entries"`
}

type ContractLedgerEntry struct {
	ID             string   `json:"id"`
	Package        string   `json:"package"`
	Kind           string   `json:"kind"`
	Contract       string   `json:"contract"`
	Summary        string   `json:"summary"`
	Disposition    string   `json:"disposition"`
	Coverage       []string `json:"coverage"`
	SourceEvidence []string `json:"source_evidence"`
	TestEvidence   []string `json:"test_evidence,omitempty"`
	Terms          []string `json:"terms"`
}

type ContractPackageReview struct {
	Package         string              `json:"package"`
	Disposition     string              `json:"disposition"`
	Summary         string              `json:"summary"`
	ReviewedSurface PublicSurfaceReview `json:"reviewed_surface"`
	Coverage        []string            `json:"coverage"`
	SourceEvidence  []string            `json:"source_evidence"`
	ContractIDs     []string            `json:"contract_ids,omitempty"`
}

type PublicSurfaceReview struct {
	TopLevel    int      `json:"top_level"`
	Fields      int      `json:"fields"`
	Methods     int      `json:"methods"`
	EntryCount  int      `json:"entry_count"`
	Fingerprint string   `json:"fingerprint"`
	ReviewScope []string `json:"review_scope"`
}

type ManifestPackageEntry struct {
	Package     string   `json:"package"`
	Tier        string   `json:"tier"`
	Coverage    []string `json:"coverage"`
	Notes       string   `json:"notes,omitempty"`
	TopLevel    int      `json:"top_level"`
	Fields      int      `json:"fields"`
	Methods     int      `json:"methods"`
	Fingerprint string   `json:"fingerprint"`
}

type PackageInventory struct {
	Package       string
	TopLevel      int
	Fields        int
	Methods       int
	Fingerprint   string
	TopNames      []string
	FieldsByType  map[string]map[string]bool
	MethodsByType map[string]map[string]bool
	Entries       []AuditLedgerEntry
}

const requiredReviewScopeCount = 4

var requiredReviewScope = []string{
	"exported top-level symbols",
	"exported struct fields",
	"exported methods",
	"runtime/user-facing contracts",
}

func main() {
	sourceDir := flag.String("source", ".", "path to the VEF Framework Go source repository")
	manifestPath := flag.String("manifest", "../vef-framework-go-docs/scripts/api-audit-manifest.json", "path to the API audit manifest")
	ledgerPath := flag.String("ledger", "../vef-framework-go-docs/scripts/api-audit-ledger.json", "path to the API audit ledger")
	contractLedgerPath := flag.String("contract-ledger", "../vef-framework-go-docs/scripts/api-contract-ledger.json", "path to the semantic API contract ledger")
	printCurrent := flag.Bool("print-current", false, "print current package counts and fingerprints as JSON entries")
	printLedger := flag.Bool("print-ledger", false, "print a fresh API audit ledger as JSON")
	writeManifest := flag.Bool("write-manifest", false, "refresh package counts and fingerprints in -manifest")
	writeLedger := flag.Bool("write-ledger", false, "write a fresh API audit ledger to -ledger")
	flag.Parse()

	inventory, err := buildInventory(*sourceDir)
	if err != nil {
		panic(err)
	}

	if *printCurrent {
		printCurrentInventory(inventory)

		return
	}

	manifest, err := loadManifest(*manifestPath)
	if err != nil {
		panic(err)
	}

	if *printLedger {
		printAuditLedger(manifest, inventory)

		return
	}

	if *writeManifest {
		if err := writeManifestFile(*manifestPath, manifest, inventory); err != nil {
			panic(err)
		}

		return
	}

	if *writeLedger {
		if err := writeAuditLedger(*ledgerPath, manifest, inventory); err != nil {
			panic(err)
		}

		return
	}

	ledger, err := loadAuditLedger(*ledgerPath)
	if err != nil {
		panic(err)
	}

	contractLedger, err := loadContractLedger(*contractLedgerPath)
	if err != nil {
		panic(err)
	}

	if err := verify(manifest, ledger, contractLedger, inventory, filepath.Dir(*manifestPath)); err != nil {
		panic(err)
	}

	fmt.Printf(
		"API audit manifest verified: %d packages, %d public entries, %d semantic package reviews\n",
		len(inventory),
		ledger.EntryCount,
		len(contractLedger.PackageReviews),
	)
}

func buildInventory(sourceDir string) (map[string]PackageInventory, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps,
		Dir:  sourceDir,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, err
	}

	sort.Slice(pkgs, func(i, j int) bool { return pkgs[i].PkgPath < pkgs[j].PkgPath })

	result := make(map[string]PackageInventory)
	for _, pkg := range pkgs {
		if strings.Contains(pkg.PkgPath, "/internal") || pkg.Name == "main" {
			continue
		}
		if len(pkg.Errors) > 0 {
			return nil, fmt.Errorf("package errors in %s: %v", pkg.PkgPath, pkg.Errors)
		}

		result[pkg.PkgPath] = packageInventory(pkg)
	}

	return result, nil
}

func packageInventory(pkg *packages.Package) PackageInventory {
	var ids []string
	counts := PackageInventory{
		Package:       pkg.PkgPath,
		FieldsByType:  make(map[string]map[string]bool),
		MethodsByType: make(map[string]map[string]bool),
	}

	scope := pkg.Types.Scope()
	for _, name := range exportedNames(scope) {
		obj := scope.Lookup(name)
		if obj == nil {
			continue
		}

		signature := objectString(obj)
		counts.TopLevel++
		counts.TopNames = append(counts.TopNames, name)
		ids = append(ids, "TOP "+signature)
		counts.Entries = append(counts.Entries, AuditLedgerEntry{
			ID:        ledgerID(pkg.PkgPath, "top", name),
			Package:   pkg.PkgPath,
			Kind:      "top",
			Symbol:    name,
			Signature: signature,
		})

		typeName, ok := obj.(*types.TypeName)
		if !ok {
			continue
		}

		for _, field := range exportedFields(typeName.Type()) {
			counts.Fields++
			ids = append(ids, "FIELD "+typeName.Name()+"."+field)
			name := memberName(field)
			addMember(counts.FieldsByType, typeName.Name(), name)
			counts.Entries = append(counts.Entries, AuditLedgerEntry{
				ID:        ledgerID(pkg.PkgPath, "field", typeName.Name()+"."+name),
				Package:   pkg.PkgPath,
				Kind:      "field",
				Symbol:    typeName.Name() + "." + name,
				Signature: field,
			})
		}

		for _, method := range exportedMethodSet(typeName.Type()) {
			counts.Methods++
			ids = append(ids, "METHOD "+typeName.Name()+"."+method)
			name := memberName(method)
			addMember(counts.MethodsByType, typeName.Name(), name)
			counts.Entries = append(counts.Entries, AuditLedgerEntry{
				ID:        ledgerID(pkg.PkgPath, "method", typeName.Name()+"."+name),
				Package:   pkg.PkgPath,
				Kind:      "method",
				Symbol:    typeName.Name() + "." + name,
				Signature: method,
			})
		}
	}

	sort.Strings(ids)
	sort.Strings(counts.TopNames)
	sortLedgerEntries(counts.Entries)
	sum := sha256.Sum256([]byte(strings.Join(ids, "\n")))
	counts.Fingerprint = hex.EncodeToString(sum[:])

	return counts
}

func loadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

func loadAuditLedger(path string) (*AuditLedger, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var ledger AuditLedger
	if err := json.Unmarshal(data, &ledger); err != nil {
		return nil, err
	}

	return &ledger, nil
}

func loadContractLedger(path string) (*ContractLedger, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var ledger ContractLedger
	if err := json.Unmarshal(data, &ledger); err != nil {
		return nil, err
	}

	return &ledger, nil
}

func verify(manifest *Manifest, ledger *AuditLedger, contractLedger *ContractLedger, inventory map[string]PackageInventory, baseDir string) error {
	if manifest.SourceModule != "github.com/coldsmirk/vef-framework-go" {
		return fmt.Errorf("unexpected source_module %q", manifest.SourceModule)
	}
	if ledger.SourceModule != manifest.SourceModule {
		return fmt.Errorf("unexpected ledger source_module %q", ledger.SourceModule)
	}
	if contractLedger.SourceModule != manifest.SourceModule {
		return fmt.Errorf("unexpected contract ledger source_module %q", contractLedger.SourceModule)
	}

	seen := map[string]bool{}
	var failures []string
	repoRoot := filepath.Join(baseDir, "..")
	failures = append(failures, verifyGeneratedIndex(inventory, filepath.Join(repoRoot, "docs/reference/public-api-index.md"), "English")...)
	failures = append(failures, verifyGeneratedIndex(inventory, filepath.Join(repoRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"), "Chinese")...)
	failures = append(failures, verifyTopLevelMentions(inventory, repoRoot)...)
	failures = append(failures, verifyPackageCoverageMentions(manifest, inventory, repoRoot)...)
	failures = append(failures, verifyPublicAPIReferences(inventory, repoRoot)...)
	failures = append(failures, verifyAuditLedger(ledger, manifest, inventory, repoRoot)...)
	failures = append(failures, verifyContractLedger(contractLedger, manifest, inventory, repoRoot)...)

	for _, entry := range manifest.Packages {
		if entry.Package == "" {
			failures = append(failures, "manifest contains an entry with empty package")

			continue
		}
		if seen[entry.Package] {
			failures = append(failures, "manifest contains duplicate package "+entry.Package)

			continue
		}
		seen[entry.Package] = true

		current, ok := inventory[entry.Package]
		if !ok {
			failures = append(failures, "manifest package no longer exists: "+entry.Package)

			continue
		}

		if entry.TopLevel != current.TopLevel ||
			entry.Fields != current.Fields ||
			entry.Methods != current.Methods ||
			entry.Fingerprint != current.Fingerprint {
			failures = append(failures, fmt.Sprintf(
				"public surface changed for %s: manifest top/fields/methods/fingerprint=%d/%d/%d/%s current=%d/%d/%d/%s",
				entry.Package,
				entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
				current.TopLevel, current.Fields, current.Methods, current.Fingerprint,
			))
		}

		if entry.Tier == "" {
			failures = append(failures, "missing tier for "+entry.Package)
		}
		if len(entry.Coverage) == 0 {
			failures = append(failures, "missing coverage for "+entry.Package)
		}
		for _, coverage := range entry.Coverage {
			if strings.Contains(coverage, "://") {
				continue
			}

			path := filepath.Clean(filepath.Join(baseDir, "..", coverage))
			if _, err := os.Stat(path); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					failures = append(failures, fmt.Sprintf("coverage file missing for %s: %s", entry.Package, coverage))
				} else {
					failures = append(failures, fmt.Sprintf("coverage file check failed for %s: %s: %v", entry.Package, coverage, err))
				}
			}

			if zhMirror, ok := chineseMirrorPath(coverage); ok {
				path := filepath.Clean(filepath.Join(baseDir, "..", zhMirror))
				if _, err := os.Stat(path); err != nil {
					if errors.Is(err, os.ErrNotExist) {
						failures = append(failures, fmt.Sprintf("Chinese coverage mirror missing for %s: %s", entry.Package, zhMirror))
					} else {
						failures = append(failures, fmt.Sprintf("Chinese coverage mirror check failed for %s: %s: %v", entry.Package, zhMirror, err))
					}
				}
			}
		}
	}

	for pkg := range inventory {
		if !seen[pkg] {
			failures = append(failures, "unclassified public package: "+pkg)
		}
	}

	if len(failures) > 0 {
		sort.Strings(failures)

		return fmt.Errorf("API audit verification failed:\n%s", strings.Join(failures, "\n"))
	}

	return nil
}

func verifyAuditLedger(ledger *AuditLedger, manifest *Manifest, inventory map[string]PackageInventory, repoRoot string) []string {
	expected := auditLedgerEntries(manifest, inventory)
	expectedByID := make(map[string]AuditLedgerEntry, len(expected))
	for _, entry := range expected {
		expectedByID[entry.ID] = entry
	}

	manifestByPackage := manifestPackages(manifest)
	var failures []string
	if ledger.Scope == "" {
		failures = append(failures, "audit ledger missing scope")
	}
	if ledger.EntryCount != len(expected) {
		failures = append(failures, fmt.Sprintf("audit ledger entry_count mismatch: ledger=%d current=%d", ledger.EntryCount, len(expected)))
	}
	if len(ledger.Entries) != ledger.EntryCount {
		failures = append(failures, fmt.Sprintf("audit ledger entries length mismatch: entry_count=%d entries=%d", ledger.EntryCount, len(ledger.Entries)))
	}

	currentFingerprint := auditLedgerFingerprint(expected)
	if ledger.Fingerprint != currentFingerprint {
		failures = append(failures, fmt.Sprintf("audit ledger fingerprint mismatch: ledger=%s current=%s", ledger.Fingerprint, currentFingerprint))
	}

	seen := make(map[string]bool, len(ledger.Entries))
	for _, entry := range ledger.Entries {
		if entry.ID == "" {
			failures = append(failures, "audit ledger contains an entry with empty id")

			continue
		}
		if seen[entry.ID] {
			failures = append(failures, "audit ledger contains duplicate entry "+entry.ID)

			continue
		}
		seen[entry.ID] = true

		current, ok := expectedByID[entry.ID]
		if !ok {
			failures = append(failures, "audit ledger contains stale entry "+entry.ID)

			continue
		}

		if entry.Package != current.Package ||
			entry.Kind != current.Kind ||
			entry.Symbol != current.Symbol ||
			entry.Signature != current.Signature {
			failures = append(failures, fmt.Sprintf("audit ledger stale metadata for %s", entry.ID))
		}

		if entry.Disposition == "" {
			failures = append(failures, "audit ledger missing disposition for "+entry.ID)
		} else if _, ok := ledger.Dispositions[entry.Disposition]; !ok {
			failures = append(failures, fmt.Sprintf("audit ledger entry %s uses undefined disposition %q", entry.ID, entry.Disposition))
		}

		if len(entry.Coverage) == 0 {
			failures = append(failures, "audit ledger missing coverage for "+entry.ID)
		}
		allowedCoverage := map[string]bool{}
		if manifestEntry, ok := manifestByPackage[entry.Package]; ok {
			for _, coverage := range manifestEntry.Coverage {
				allowedCoverage[coverage] = true
			}
		}
		for _, coverage := range entry.Coverage {
			if !strings.Contains(coverage, "://") && !allowedCoverage[coverage] {
				failures = append(failures, fmt.Sprintf("audit ledger entry %s uses coverage outside package manifest: %s", entry.ID, coverage))
			}
			failures = append(failures, verifyCoverageFile(repoRoot, entry.ID, coverage)...)
		}
	}

	for _, entry := range expected {
		if !seen[entry.ID] {
			failures = append(failures, "audit ledger missing current public API entry "+entry.ID)
		}
	}

	return failures
}

func verifyContractLedger(ledger *ContractLedger, manifest *Manifest, inventory map[string]PackageInventory, repoRoot string) []string {
	var failures []string
	if ledger.Scope == "" {
		failures = append(failures, "contract ledger missing scope")
	}

	manifestByPackage := manifestPackages(manifest)
	failures = append(failures, verifyContractPackageReviews(ledger, manifest, inventory, repoRoot, manifestByPackage)...)

	entryByID := make(map[string]ContractLedgerEntry, len(ledger.Entries))
	seen := make(map[string]bool, len(ledger.Entries))
	for _, entry := range ledger.Entries {
		if entry.ID == "" {
			failures = append(failures, "contract ledger contains an entry with empty id")

			continue
		}
		if seen[entry.ID] {
			failures = append(failures, "contract ledger contains duplicate entry "+entry.ID)

			continue
		}
		seen[entry.ID] = true
		entryByID[entry.ID] = entry

		if entry.Package == "" {
			failures = append(failures, "contract ledger entry "+entry.ID+" has empty package")
		} else if _, ok := inventory[entry.Package]; !ok {
			failures = append(failures, "contract ledger entry "+entry.ID+" references unknown public package "+entry.Package)
		}
		if entry.Kind == "" {
			failures = append(failures, "contract ledger missing kind for "+entry.ID)
		}
		if entry.Contract == "" {
			failures = append(failures, "contract ledger missing contract for "+entry.ID)
		}
		if entry.Summary == "" {
			failures = append(failures, "contract ledger missing summary for "+entry.ID)
		}
		if entry.Disposition == "" {
			failures = append(failures, "contract ledger missing disposition for "+entry.ID)
		}
		if len(entry.Coverage) == 0 {
			failures = append(failures, "contract ledger missing coverage for "+entry.ID)
		}
		if len(entry.SourceEvidence) == 0 {
			failures = append(failures, "contract ledger missing source evidence for "+entry.ID)
		}
		if len(entry.Terms) == 0 {
			failures = append(failures, "contract ledger missing terms for "+entry.ID)
		}

		allowedCoverage := map[string]bool{}
		if manifestEntry, ok := manifestByPackage[entry.Package]; ok {
			for _, coverage := range manifestEntry.Coverage {
				allowedCoverage[coverage] = true
			}
		}
		for _, coverage := range entry.Coverage {
			if !strings.Contains(coverage, "://") && !allowedCoverage[coverage] {
				failures = append(failures, fmt.Sprintf("contract ledger entry %s uses coverage outside package manifest: %s", entry.ID, coverage))
			}
			failures = append(failures, verifyCoverageFile(repoRoot, entry.ID, coverage)...)
		}
		for _, evidence := range entry.SourceEvidence {
			failures = append(failures, verifySourceEvidence(repoRoot, entry.ID, evidence)...)
		}
		for _, evidence := range entry.TestEvidence {
			failures = append(failures, verifySourceEvidence(repoRoot, entry.ID, evidence)...)
		}
		failures = append(failures, verifyContractTerms(repoRoot, entry)...)
	}

	reviewedContracts := make(map[string]bool)
	for _, review := range ledger.PackageReviews {
		for _, id := range review.ContractIDs {
			reviewedContracts[id] = true
			entry, ok := entryByID[id]
			if !ok {
				failures = append(failures, fmt.Sprintf("contract package review for %s references unknown contract id %s", review.Package, id))

				continue
			}
			if entry.Package != review.Package {
				failures = append(failures, fmt.Sprintf("contract package review for %s references contract from %s: %s", review.Package, entry.Package, id))
			}
		}
	}
	for _, entry := range ledger.Entries {
		if !reviewedContracts[entry.ID] {
			failures = append(failures, "contract ledger entry not linked from package review: "+entry.ID)
		}
	}

	return failures
}

func verifyContractPackageReviews(
	ledger *ContractLedger,
	manifest *Manifest,
	inventory map[string]PackageInventory,
	repoRoot string,
	manifestByPackage map[string]ManifestPackageEntry,
) []string {
	allowedDispositions := map[string]bool{
		"has-semantic-contracts":  true,
		"symbol-and-member-only":  true,
		"command-surface-only":    true,
		"configuration-contracts": true,
		"dto-wire-shape":          true,
	}

	var failures []string
	if len(ledger.PackageReviews) != len(manifest.Packages) {
		failures = append(failures, fmt.Sprintf("contract package review count mismatch: reviews=%d manifest=%d", len(ledger.PackageReviews), len(manifest.Packages)))
	}

	seen := make(map[string]bool, len(ledger.PackageReviews))
	for _, review := range ledger.PackageReviews {
		if review.Package == "" {
			failures = append(failures, "contract package review has empty package")

			continue
		}
		if seen[review.Package] {
			failures = append(failures, "contract package review duplicated package "+review.Package)

			continue
		}
		seen[review.Package] = true

		if _, ok := inventory[review.Package]; !ok {
			failures = append(failures, "contract package review references unknown public package "+review.Package)
		}
		if _, ok := manifestByPackage[review.Package]; !ok {
			failures = append(failures, "contract package review references package outside manifest "+review.Package)
		}
		if !allowedDispositions[review.Disposition] {
			failures = append(failures, fmt.Sprintf("contract package review for %s uses unknown disposition %q", review.Package, review.Disposition))
		}
		if review.Summary == "" {
			failures = append(failures, "contract package review missing summary for "+review.Package)
		}
		if current, ok := inventory[review.Package]; ok {
			failures = append(failures, verifyReviewedSurface(review.Package, review.ReviewedSurface, current)...)
		}
		if len(review.Coverage) == 0 {
			failures = append(failures, "contract package review missing coverage for "+review.Package)
		}
		if len(review.SourceEvidence) == 0 {
			failures = append(failures, "contract package review missing source evidence for "+review.Package)
		}
		if review.Disposition == "has-semantic-contracts" && len(review.ContractIDs) == 0 {
			failures = append(failures, "semantic contract package review missing contract_ids for "+review.Package)
		}
		if review.Disposition != "has-semantic-contracts" && len(review.ContractIDs) > 0 {
			failures = append(failures, "non-semantic contract package review has contract_ids for "+review.Package)
		}

		allowedCoverage := map[string]bool{}
		if manifestEntry, ok := manifestByPackage[review.Package]; ok {
			for _, coverage := range manifestEntry.Coverage {
				allowedCoverage[coverage] = true
			}
		}
		for _, coverage := range review.Coverage {
			if !strings.Contains(coverage, "://") && !allowedCoverage[coverage] {
				failures = append(failures, fmt.Sprintf("contract package review for %s uses coverage outside package manifest: %s", review.Package, coverage))
			}
			failures = append(failures, verifyCoverageFile(repoRoot, review.Package, coverage)...)
		}
		for _, evidence := range review.SourceEvidence {
			failures = append(failures, verifySourceEvidence(repoRoot, review.Package, evidence)...)
		}
	}

	for _, entry := range manifest.Packages {
		if !seen[entry.Package] {
			failures = append(failures, "contract package review missing manifest package "+entry.Package)
		}
	}

	return failures
}

func verifyReviewedSurface(pkg string, reviewed PublicSurfaceReview, current PackageInventory) []string {
	var failures []string
	currentEntryCount := current.TopLevel + current.Fields + current.Methods
	if reviewed.TopLevel != current.TopLevel ||
		reviewed.Fields != current.Fields ||
		reviewed.Methods != current.Methods ||
		reviewed.EntryCount != currentEntryCount ||
		reviewed.Fingerprint != current.Fingerprint {
		failures = append(failures, fmt.Sprintf(
			"contract package review surface mismatch for %s: review top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s current=%d/%d/%d/%d/%s",
			pkg,
			reviewed.TopLevel, reviewed.Fields, reviewed.Methods, reviewed.EntryCount, reviewed.Fingerprint,
			current.TopLevel, current.Fields, current.Methods, currentEntryCount, current.Fingerprint,
		))
	}
	failures = append(failures, verifyReviewScope(pkg, reviewed.ReviewScope)...)

	return failures
}

func verifyReviewScope(pkg string, reviewScope []string) []string {
	if len(requiredReviewScope) != requiredReviewScopeCount {
		return []string{fmt.Sprintf("internal verifier requiredReviewScope has %d entries, want %d", len(requiredReviewScope), requiredReviewScopeCount)}
	}

	seen := make(map[string]bool, len(reviewScope))
	for _, scope := range reviewScope {
		if scope == "" {
			return []string{"contract package review contains empty review_scope for " + pkg}
		}
		seen[scope] = true
	}

	var missing []string
	for _, scope := range requiredReviewScope {
		if !seen[scope] {
			missing = append(missing, scope)
		}
	}
	if len(missing) > 0 {
		return []string{fmt.Sprintf("contract package review for %s missing required review_scope: %s", pkg, strings.Join(missing, ", "))}
	}

	return nil
}

func verifySourceEvidence(repoRoot, entryID, evidence string) []string {
	if strings.Contains(evidence, "://") {
		return nil
	}

	path, _, _ := strings.Cut(evidence, ":")
	if path == "" {
		return []string{fmt.Sprintf("contract ledger evidence has empty path for %s: %s", entryID, evidence)}
	}

	sourceRoot := filepath.Clean(filepath.Join(repoRoot, "..", "vef-framework-go"))
	fullPath := filepath.Clean(filepath.Join(sourceRoot, path))
	if rel, err := filepath.Rel(sourceRoot, fullPath); err != nil || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return []string{fmt.Sprintf("contract ledger evidence escapes source root for %s: %s", entryID, evidence)}
	}
	if _, err := os.Stat(fullPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{fmt.Sprintf("contract ledger evidence file missing for %s: %s", entryID, evidence)}
		}

		return []string{fmt.Sprintf("contract ledger evidence file check failed for %s: %s: %v", entryID, evidence, err)}
	}

	return nil
}

func verifyContractTerms(repoRoot string, entry ContractLedgerEntry) []string {
	corpora, err := coverageCorpora(repoRoot, entry.Coverage)
	if err != nil {
		return []string{fmt.Sprintf("failed to read contract coverage corpus for %s: %v", entry.ID, err)}
	}

	var failures []string
	for _, term := range entry.Terms {
		if term == "" {
			failures = append(failures, "contract ledger entry "+entry.ID+" contains an empty term")

			continue
		}
		for _, corpus := range corpora {
			if !strings.Contains(corpus.Content, term) {
				failures = append(failures, fmt.Sprintf("%s contract coverage missing term for %s: %s", corpus.Label, entry.ID, term))
			}
		}
	}

	return failures
}

func auditLedgerEntries(manifest *Manifest, inventory map[string]PackageInventory) []AuditLedgerEntry {
	manifestByPackage := manifestPackages(manifest)
	entries := make([]AuditLedgerEntry, 0)
	for _, inv := range inventory {
		coverage := []string{"docs/reference/public-api-index.md"}
		tier := ""
		if manifestEntry, ok := manifestByPackage[inv.Package]; ok {
			coverage = manifestEntry.Coverage
			tier = manifestEntry.Tier
		}
		for _, entry := range inv.Entries {
			entry.Disposition = defaultDisposition(entry.Package, entry.Kind, tier)
			entry.Coverage = append([]string(nil), coverage...)
			entries = append(entries, entry)
		}
	}
	sortLedgerEntries(entries)

	return entries
}

func manifestPackages(manifest *Manifest) map[string]ManifestPackageEntry {
	result := make(map[string]ManifestPackageEntry, len(manifest.Packages))
	for _, entry := range manifest.Packages {
		result[entry.Package] = entry
	}

	return result
}

func defaultDisposition(pkg, kind, tier string) string {
	if strings.Contains(pkg, "/cmd/vef-cli/") {
		return "excluded:cli-implementation-export"
	}
	if kind == "top" {
		return "documented:top-level"
	}

	switch tier {
	case "builder":
		return "grouped:builder-member-family"
	case "configuration":
		return "grouped:configuration-fields"
	case "dto":
		return "grouped:dto-fields"
	default:
		return "grouped:type-member-family"
	}
}

func verifyCoverageFile(repoRoot, entryID, coverage string) []string {
	if strings.Contains(coverage, "://") {
		return nil
	}

	var failures []string
	path := filepath.Clean(filepath.Join(repoRoot, coverage))
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			failures = append(failures, fmt.Sprintf("audit ledger coverage file missing for %s: %s", entryID, coverage))
		} else {
			failures = append(failures, fmt.Sprintf("audit ledger coverage file check failed for %s: %s: %v", entryID, coverage, err))
		}
	}

	if zhMirror, ok := chineseMirrorPath(coverage); ok {
		path := filepath.Clean(filepath.Join(repoRoot, zhMirror))
		if _, err := os.Stat(path); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				failures = append(failures, fmt.Sprintf("audit ledger Chinese coverage mirror missing for %s: %s", entryID, zhMirror))
			} else {
				failures = append(failures, fmt.Sprintf("audit ledger Chinese coverage mirror check failed for %s: %s: %v", entryID, zhMirror, err))
			}
		}
	}

	return failures
}

func verifyGeneratedIndex(inventory map[string]PackageInventory, path, label string) []string {
	indexInventory, err := parseGeneratedIndex(path)
	if err != nil {
		return []string{fmt.Sprintf("%s generated index parse failed: %v", label, err)}
	}

	var failures []string
	for pkg, current := range inventory {
		indexed, ok := indexInventory[pkg]
		if !ok {
			failures = append(failures, fmt.Sprintf("%s generated index missing package: %s", label, pkg))

			continue
		}

		if indexed.TopLevel != current.TopLevel ||
			indexed.Fields != current.Fields ||
			indexed.Methods != current.Methods ||
			indexed.Fingerprint != current.Fingerprint {
			failures = append(failures, fmt.Sprintf(
				"%s generated index stale for %s: index top/fields/methods/fingerprint=%d/%d/%d/%s current=%d/%d/%d/%s",
				label,
				pkg,
				indexed.TopLevel, indexed.Fields, indexed.Methods, indexed.Fingerprint,
				current.TopLevel, current.Fields, current.Methods, current.Fingerprint,
			))
		}
	}
	for pkg := range indexInventory {
		if _, ok := inventory[pkg]; !ok {
			failures = append(failures, fmt.Sprintf("%s generated index has stale package: %s", label, pkg))
		}
	}

	return failures
}

func parseGeneratedIndex(path string) (map[string]PackageInventory, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	result := make(map[string]PackageInventory)
	var (
		currentPkg  string
		currentType string
		ids         []string
	)

	flush := func() {
		if currentPkg == "" {
			return
		}
		inv := result[currentPkg]
		sort.Strings(ids)
		sum := sha256.Sum256([]byte(strings.Join(ids, "\n")))
		inv.Fingerprint = hex.EncodeToString(sum[:])
		result[currentPkg] = inv
		ids = nil
	}

	for _, line := range strings.Split(string(data), "\n") {
		if after, ok := strings.CutPrefix(line, "## "); ok &&
			strings.HasPrefix(after, "github.com/coldsmirk/vef-framework-go") {
			flush()
			currentPkg = strings.TrimSpace(after)
			currentType = ""
			result[currentPkg] = PackageInventory{
				Package:       currentPkg,
				FieldsByType:  make(map[string]map[string]bool),
				MethodsByType: make(map[string]map[string]bool),
			}

			continue
		}

		if currentPkg == "" {
			continue
		}

		inv := result[currentPkg]
		if obj, ok := cutObjectLine(line); ok {
			inv.TopLevel++
			result[currentPkg] = inv
			ids = append(ids, "TOP "+obj)

			if fields := strings.Fields(obj); len(fields) > 0 {
				currentType = fields[0]
			}

			continue
		}

		if field, ok := strings.CutPrefix(line, "  FIELD "); ok {
			field = strings.TrimSpace(field)
			inv.Fields++
			ids = append(ids, "FIELD "+currentType+"."+field)
			addMember(inv.FieldsByType, currentType, memberName(field))
			result[currentPkg] = inv

			continue
		}

		if method, ok := strings.CutPrefix(line, "  METHOD "); ok {
			method = strings.TrimSpace(method)
			inv.Methods++
			ids = append(ids, "METHOD "+currentType+"."+method)
			addMember(inv.MethodsByType, currentType, memberName(method))
			result[currentPkg] = inv
		}
	}
	flush()

	return result, nil
}

func cutObjectLine(line string) (string, bool) {
	for _, prefix := range []string{"TYPE ", "FUNC ", "CONST ", "VAR "} {
		if after, ok := strings.CutPrefix(line, prefix); ok {
			return after, true
		}
	}

	return "", false
}

func memberName(signature string) string {
	fields := strings.Fields(signature)
	if len(fields) == 0 {
		return ""
	}

	return fields[0]
}

func addMember(target map[string]map[string]bool, typeName, member string) {
	if typeName == "" || member == "" {
		return
	}
	if target[typeName] == nil {
		target[typeName] = make(map[string]bool)
	}

	target[typeName][member] = true
}

func verifyTopLevelMentions(inventory map[string]PackageInventory, repoRoot string) []string {
	corpora, err := docsCorpora(repoRoot)
	if err != nil {
		return []string{"failed to read docs corpus: " + err.Error()}
	}

	var failures []string
	for pkg, inv := range inventory {
		if strings.Contains(pkg, "/cmd/vef-cli/") {
			continue
		}
		for _, name := range inv.TopNames {
			for _, corpus := range corpora {
				if !wordMentioned(corpus.Content, name) {
					failures = append(failures, fmt.Sprintf("%s top-level public symbol missing outside generated index: %s.%s", corpus.Label, pkg, name))
				}
			}
		}
	}

	return failures
}

func verifyPackageCoverageMentions(manifest *Manifest, inventory map[string]PackageInventory, repoRoot string) []string {
	var failures []string
	for _, entry := range manifest.Packages {
		if strings.Contains(entry.Package, "/cmd/vef-cli/") {
			continue
		}

		inv, ok := inventory[entry.Package]
		if !ok {
			continue
		}

		corpora, err := coverageCorpora(repoRoot, entry.Coverage)
		if err != nil {
			failures = append(failures, fmt.Sprintf("failed to read coverage corpus for %s: %v", entry.Package, err))

			continue
		}

		for _, name := range inv.TopNames {
			for _, corpus := range corpora {
				if !wordMentioned(corpus.Content, name) {
					failures = append(failures, fmt.Sprintf("%s package coverage missing top-level public symbol: %s.%s", corpus.Label, entry.Package, name))
				}
			}
		}
	}

	return failures
}

func verifyPublicAPIReferences(inventory map[string]PackageInventory, repoRoot string) []string {
	files, err := docsFiles(repoRoot)
	if err != nil {
		return []string{"failed to read docs files: " + err.Error()}
	}

	knownPackages := packageNames(inventory)
	topLevelByPackageName := topLevelNamesByPackageName(inventory)
	memberTypes := memberTypes(inventory)

	var failures []string
	for _, file := range files {
		for _, ref := range codeReferences(file.Content) {
			if pkg, symbol, ok := splitPackageSymbol(ref); ok && knownPackages[pkg] {
				if !topLevelByPackageName[pkg][symbol] {
					failures = append(failures, fmt.Sprintf("%s references unknown public top-level symbol: %s", file.Path, ref))
				}

				continue
			}

			if typ, member, ok := splitTypeMember(ref); ok && memberTypes[typ] != nil {
				if !memberTypes[typ][member] {
					failures = append(failures, fmt.Sprintf("%s references unknown exported type member: %s", file.Path, ref))
				}
			}
		}
	}

	sort.Strings(failures)

	return failures
}

type docsFile struct {
	Path    string
	Content string
}

func docsFiles(repoRoot string) ([]docsFile, error) {
	var files []docsFile
	for _, root := range []string{
		filepath.Join(repoRoot, "docs"),
		filepath.Join(repoRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current"),
	} {
		if err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".md") || strings.HasSuffix(path, "reference/public-api-index.md") {
				return nil
			}

			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			rel, err := filepath.Rel(repoRoot, path)
			if err != nil {
				rel = path
			}
			files = append(files, docsFile{
				Path:    filepath.ToSlash(rel),
				Content: string(data),
			})

			return nil
		}); err != nil {
			return nil, err
		}
	}

	return files, nil
}

func packageNames(inventory map[string]PackageInventory) map[string]bool {
	result := make(map[string]bool)
	for pkg := range inventory {
		if strings.Contains(pkg, "/cmd/vef-cli/") {
			continue
		}

		result[pathBase(pkg)] = true
	}

	return result
}

func topLevelNamesByPackageName(inventory map[string]PackageInventory) map[string]map[string]bool {
	result := make(map[string]map[string]bool)
	for pkg, inv := range inventory {
		if strings.Contains(pkg, "/cmd/vef-cli/") {
			continue
		}

		name := pathBase(pkg)
		if result[name] == nil {
			result[name] = make(map[string]bool)
		}
		for _, top := range inv.TopNames {
			result[name][top] = true
		}
	}

	return result
}

func memberTypes(inventory map[string]PackageInventory) map[string]map[string]bool {
	result := make(map[string]map[string]bool)
	for pkg, inv := range inventory {
		if strings.Contains(pkg, "/cmd/vef-cli/") {
			continue
		}

		addMembers(result, inv.FieldsByType)
		addMembers(result, inv.MethodsByType)
	}

	return result
}

func addMembers(target map[string]map[string]bool, source map[string]map[string]bool) {
	for typ, members := range source {
		if target[typ] == nil {
			target[typ] = make(map[string]bool)
		}
		for member := range members {
			target[typ][member] = true
		}
	}
}

func pathBase(path string) string {
	if idx := strings.LastIndex(path, "/"); idx >= 0 {
		return path[idx+1:]
	}

	return path
}

var codeReferencePattern = regexp.MustCompile("`([A-Za-z_][A-Za-z0-9_]*(?:\\.[A-Za-z_][A-Za-z0-9_]*)+)`")

func codeReferences(content string) []string {
	matches := codeReferencePattern.FindAllStringSubmatch(content, -1)
	refs := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			refs = append(refs, match[1])
		}
	}

	return refs
}

func splitPackageSymbol(ref string) (string, string, bool) {
	parts := strings.Split(ref, ".")
	if len(parts) != 2 {
		return "", "", false
	}

	if !tokenExported(parts[1]) {
		return "", "", false
	}

	return parts[0], parts[1], true
}

func splitTypeMember(ref string) (string, string, bool) {
	parts := strings.Split(ref, ".")
	if len(parts) != 2 {
		return "", "", false
	}

	if !tokenExported(parts[0]) || !tokenExported(parts[1]) {
		return "", "", false
	}

	return parts[0], parts[1], true
}

type docsCorpus struct {
	Label   string
	Content string
}

func docsCorpora(repoRoot string) ([]docsCorpus, error) {
	roots := []struct {
		label string
		path  string
	}{
		{label: "English", path: filepath.Join(repoRoot, "docs")},
		{label: "Chinese", path: filepath.Join(repoRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current")},
	}

	corpora := make([]docsCorpus, 0, len(roots))
	for _, root := range roots {
		content, err := docsCorpusContent(root.path)
		if err != nil {
			return nil, err
		}
		corpora = append(corpora, docsCorpus{Label: root.label, Content: content})
	}

	return corpora, nil
}

func coverageCorpora(repoRoot string, coverage []string) ([]docsCorpus, error) {
	english := make([]string, 0, len(coverage))
	chinese := make([]string, 0, len(coverage))
	for _, path := range coverage {
		if strings.Contains(path, "://") {
			continue
		}
		if strings.HasSuffix(path, "reference/public-api-index.md") {
			continue
		}

		english = append(english, path)
		if zhPath, ok := chineseMirrorPath(path); ok {
			chinese = append(chinese, zhPath)
		}
	}

	corpora := []docsCorpus{
		{Label: "English", Content: ""},
		{Label: "Chinese", Content: ""},
	}
	var err error
	corpora[0].Content, err = readCoverageFiles(repoRoot, english)
	if err != nil {
		return nil, err
	}
	corpora[1].Content, err = readCoverageFiles(repoRoot, chinese)
	if err != nil {
		return nil, err
	}

	return corpora, nil
}

func readCoverageFiles(repoRoot string, paths []string) (string, error) {
	var builder strings.Builder
	for _, path := range paths {
		data, err := os.ReadFile(filepath.Join(repoRoot, path))
		if err != nil {
			return "", err
		}

		builder.Write(data)
		builder.WriteByte('\n')
	}

	return builder.String(), nil
}

func docsCorpusContent(root string) (string, error) {
	var builder strings.Builder
	if err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") || strings.HasSuffix(path, "reference/public-api-index.md") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		builder.Write(data)
		builder.WriteByte('\n')

		return nil
	}); err != nil {
		return "", err
	}

	return builder.String(), nil
}

func wordMentioned(corpus, word string) bool {
	pattern := regexp.MustCompile(`(^|[^A-Za-z0-9_])` + regexp.QuoteMeta(word) + `([^A-Za-z0-9_]|$)`)

	return pattern.MatchString(corpus)
}

func printCurrentInventory(inventory map[string]PackageInventory) {
	packages := make([]PackageInventory, 0, len(inventory))
	for _, inv := range inventory {
		packages = append(packages, inv)
	}
	sort.Slice(packages, func(i, j int) bool { return packages[i].Package < packages[j].Package })

	for i, inv := range packages {
		if i > 0 {
			fmt.Println(",")
		}
		data, err := json.MarshalIndent(ManifestPackageEntry{
			Package:     inv.Package,
			Tier:        "TODO",
			Coverage:    []string{"docs/reference/public-api-index.md"},
			TopLevel:    inv.TopLevel,
			Fields:      inv.Fields,
			Methods:     inv.Methods,
			Fingerprint: inv.Fingerprint,
		}, "  ", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Print("  " + string(data))
	}
	fmt.Println()
}

func printAuditLedger(manifest *Manifest, inventory map[string]PackageInventory) {
	data, err := json.MarshalIndent(freshAuditLedger(manifest, inventory), "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

func writeAuditLedger(path string, manifest *Manifest, inventory map[string]PackageInventory) error {
	data, err := json.MarshalIndent(freshAuditLedger(manifest, inventory), "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	return os.WriteFile(path, data, 0o644)
}

func writeManifestFile(path string, manifest *Manifest, inventory map[string]PackageInventory) error {
	refreshed := *manifest
	refreshed.Packages = make([]ManifestPackageEntry, 0, len(manifest.Packages))
	for _, entry := range manifest.Packages {
		current, ok := inventory[entry.Package]
		if !ok {
			return fmt.Errorf("manifest package no longer exists: %s", entry.Package)
		}
		entry.TopLevel = current.TopLevel
		entry.Fields = current.Fields
		entry.Methods = current.Methods
		entry.Fingerprint = current.Fingerprint
		refreshed.Packages = append(refreshed.Packages, entry)
	}

	seen := manifestPackages(manifest)
	for pkg := range inventory {
		if _, ok := seen[pkg]; !ok {
			return fmt.Errorf("unclassified public package: %s", pkg)
		}
	}

	data, err := json.MarshalIndent(refreshed, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	return os.WriteFile(path, data, 0o644)
}

func freshAuditLedger(manifest *Manifest, inventory map[string]PackageInventory) AuditLedger {
	entries := auditLedgerEntries(manifest, inventory)

	return AuditLedger{
		SourceModule: manifest.SourceModule,
		Scope:        "Exported top-level symbols, exported struct fields, and exported methods from non-internal/non-main packages.",
		EntryCount:   len(entries),
		Fingerprint:  auditLedgerFingerprint(entries),
		Dispositions: map[string]string{
			"documented:top-level":               "Top-level public symbols are mentioned in the package coverage docs and listed in the generated public API index.",
			"excluded:cli-implementation-export": "Exported CLI command implementation symbols are listed for audit completeness but are not supported import APIs.",
			"grouped:builder-member-family":      "Builder, query, and fluent API members are reviewed as documented method families plus the generated member index.",
			"grouped:configuration-fields":       "Configuration struct fields are reviewed through the configuration reference plus the generated field index.",
			"grouped:dto-fields":                 "Request/response DTO fields are reviewed through the module or built-in resource docs plus the generated field index.",
			"grouped:type-member-family":         "Type fields and methods are reviewed as package-specific member families plus the generated member index.",
		},
		Entries: entries,
	}
}

func auditLedgerFingerprint(entries []AuditLedgerEntry) string {
	lines := make([]string, 0, len(entries))
	for _, entry := range entries {
		lines = append(lines, strings.Join([]string{
			entry.ID,
			entry.Package,
			entry.Kind,
			entry.Symbol,
			entry.Signature,
		}, "\t"))
	}
	sort.Strings(lines)
	sum := sha256.Sum256([]byte(strings.Join(lines, "\n")))

	return hex.EncodeToString(sum[:])
}

func sortLedgerEntries(entries []AuditLedgerEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ID < entries[j].ID
	})
}

func ledgerID(pkg, kind, symbol string) string {
	return pkg + "#" + kind + ":" + symbol
}

func chineseMirrorPath(path string) (string, bool) {
	if !strings.HasPrefix(path, "docs/") || !strings.HasSuffix(path, ".md") {
		return "", false
	}

	return filepath.ToSlash(filepath.Join(
		"i18n/zh-Hans/docusaurus-plugin-content-docs/current",
		strings.TrimPrefix(path, "docs/"),
	)), true
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

func objectString(obj types.Object) string {
	var buf strings.Builder
	buf.WriteString(obj.Name())
	if obj.Type() != nil {
		buf.WriteString(" : ")
		buf.WriteString(types.TypeString(obj.Type(), packagePath))
	}
	if obj, ok := obj.(*types.Const); ok {
		buf.WriteString(" = ")
		buf.WriteString(obj.Val().ExactString())
	}

	return buf.String()
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

func packagePath(p *types.Package) string {
	return p.Path()
}
