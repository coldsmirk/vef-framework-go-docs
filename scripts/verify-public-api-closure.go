package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const (
	englishIndexPath  = "docs/reference/public-api-index.md"
	chineseIndexPath  = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"
	runtimeLedgerPath = "scripts/runtime-api-ledger.json"

	groupedPublicEntryCount           = 4033
	groupedPublicSignatureFingerprint = "4d89b0294d243989871de465647ba3017ce19e417b4987e78d3ea0b7dc9384d5"
)

type auditLedger struct {
	SourceModule string             `json:"source_module"`
	EntryCount   int                `json:"entry_count"`
	Entries      []auditLedgerEntry `json:"entries"`
}

type auditLedgerEntry struct {
	ID          string   `json:"id"`
	Package     string   `json:"package"`
	Kind        string   `json:"kind"`
	Symbol      string   `json:"symbol"`
	Signature   string   `json:"signature"`
	Disposition string   `json:"disposition"`
	Coverage    []string `json:"coverage"`
}

type manifest struct {
	SourceModule string          `json:"source_module"`
	Packages     []manifestEntry `json:"packages"`
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
	SourceModule   string                  `json:"source_module"`
	PackageReviews []contractPackageReview `json:"package_reviews"`
	Entries        []contractLedgerEntry   `json:"entries"`
}

type contractPackageReview struct {
	Package         string              `json:"package"`
	Disposition     string              `json:"disposition"`
	ReviewedSurface publicSurfaceReview `json:"reviewed_surface"`
	Coverage        []string            `json:"coverage"`
	SourceEvidence  []string            `json:"source_evidence"`
	ContractIDs     []string            `json:"contract_ids"`
}

type publicSurfaceReview struct {
	TopLevel    int    `json:"top_level"`
	Fields      int    `json:"fields"`
	Methods     int    `json:"methods"`
	EntryCount  int    `json:"entry_count"`
	Fingerprint string `json:"fingerprint"`
}

type contractLedgerEntry struct {
	ID             string   `json:"id"`
	Package        string   `json:"package"`
	Kind           string   `json:"kind"`
	Coverage       []string `json:"coverage"`
	SourceEvidence []string `json:"source_evidence"`
	TestEvidence   []string `json:"test_evidence"`
	Terms          []string `json:"terms"`
}

type runtimeLedger struct {
	EntryCount int `json:"entry_count"`
}

func main() {
	sourceDir := flag.String("source", "../vef-framework-go", "path to the VEF Framework Go source repository")
	outDir := flag.String("out", ".", "path to the VEF Framework Go docs repository")
	flag.Parse()

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)

	manifestPath := filepath.Join(docsRoot, "scripts/api-audit-manifest.json")
	auditLedgerPath := filepath.Join(docsRoot, "scripts/api-audit-ledger.json")
	contractLedgerPath := filepath.Join(docsRoot, "scripts/api-contract-ledger.json")

	m := loadJSON[manifest](manifestPath)
	audit := loadJSON[auditLedger](auditLedgerPath)
	contracts := loadJSON[contractLedger](contractLedgerPath)
	liveEntries := loadLiveEntries(sourceRoot, docsRoot, manifestPath, auditLedgerPath, contractLedgerPath)

	englishIndex := readFile(filepath.Join(docsRoot, englishIndexPath))
	chineseIndex := readFile(filepath.Join(docsRoot, chineseIndexPath))

	var failures []string
	failures = append(failures, verifyGlobalGate(sourceRoot, docsRoot, manifestPath, auditLedgerPath, contractLedgerPath)...)
	failures = append(failures, verifyExternalEmbeddingGate(sourceRoot, docsRoot)...)
	failures = append(failures, verifyRuntimeGate(sourceRoot, docsRoot)...)
	failures = append(failures, verifyManifestClosure(m, liveEntries, docsRoot)...)
	failures = append(failures, verifyAuditLedgerClosure(audit, m, englishIndex, chineseIndex, docsRoot)...)
	failures = append(failures, verifyGroupedPublicSurface(audit)...)
	failures = append(failures, verifyContractLedgerClosure(contracts, m, docsRoot, sourceRoot)...)

	if len(failures) > 0 {
		sort.Strings(failures)
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	runtimeAudit := loadJSON[runtimeLedger](filepath.Join(docsRoot, runtimeLedgerPath))
	fmt.Printf(
		"Public API closure verified: %d packages, %d public entries, %d grouped entries locked; runtime API closure verified: %d user-facing entries\n",
		len(m.Packages),
		audit.EntryCount,
		groupedPublicEntryCount,
		runtimeAudit.EntryCount,
	)
}

func verifyGlobalGate(sourceRoot, docsRoot, manifestPath, auditLedgerPath, contractLedgerPath string) []string {
	cmd := exec.Command(
		"go", "run", filepath.Join(docsRoot, "scripts/verify-api-audit.go"),
		"-source", sourceRoot,
		"-manifest", manifestPath,
		"-ledger", auditLedgerPath,
		"-contract-ledger", contractLedgerPath,
	)
	cmd.Dir = sourceRoot
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return []string{fmt.Sprintf("verify-api-audit gate failed: %v\n%s", err, strings.TrimSpace(output.String()))}
	}

	return nil
}

func verifyExternalEmbeddingGate(sourceRoot, docsRoot string) []string {
	cmd := exec.Command(
		"go", "run", filepath.Join(docsRoot, "scripts/verify-external-embedding.go"),
		"-source", sourceRoot,
		"-out", docsRoot,
	)
	cmd.Dir = sourceRoot
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return []string{fmt.Sprintf("verify-external-embedding gate failed: %v\n%s", err, strings.TrimSpace(output.String()))}
	}

	return nil
}

func verifyRuntimeGate(sourceRoot, docsRoot string) []string {
	cmd := exec.Command(
		"go", "run", filepath.Join(docsRoot, "scripts/verify-runtime-api-audit.go"),
		"-source", sourceRoot,
		"-out", docsRoot,
	)
	cmd.Dir = docsRoot
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return []string{fmt.Sprintf("verify-runtime-api-audit gate failed: %v\n%s", err, strings.TrimSpace(output.String()))}
	}

	return nil
}

func verifyManifestClosure(m manifest, liveEntries map[string]manifestEntry, docsRoot string) []string {
	var failures []string
	if m.SourceModule != "github.com/coldsmirk/vef-framework-go" {
		failures = append(failures, "manifest source_module mismatch: "+m.SourceModule)
	}

	seen := map[string]bool{}
	for _, entry := range m.Packages {
		if seen[entry.Package] {
			failures = append(failures, "duplicate manifest package "+entry.Package)
			continue
		}
		seen[entry.Package] = true

		live, ok := liveEntries[entry.Package]
		if !ok {
			failures = append(failures, "manifest package missing from live inventory "+entry.Package)
			continue
		}
		if entry.TopLevel != live.TopLevel ||
			entry.Fields != live.Fields ||
			entry.Methods != live.Methods ||
			entry.Fingerprint != live.Fingerprint {
			failures = append(failures, fmt.Sprintf(
				"manifest surface drift for %s: manifest=%d/%d/%d/%s live=%d/%d/%d/%s",
				entry.Package,
				entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
				live.TopLevel, live.Fields, live.Methods, live.Fingerprint,
			))
		}

		failures = append(failures, verifyCoveragePaths(docsRoot, entry.Package, entry.Coverage)...)
	}

	for pkg := range liveEntries {
		if !seen[pkg] {
			failures = append(failures, "live public package missing from manifest "+pkg)
		}
	}

	return failures
}

func verifyAuditLedgerClosure(audit auditLedger, m manifest, englishIndex, chineseIndex string, docsRoot string) []string {
	var failures []string
	if audit.SourceModule != m.SourceModule {
		failures = append(failures, "audit ledger source_module mismatch: "+audit.SourceModule)
	}
	if audit.EntryCount != len(audit.Entries) {
		failures = append(failures, fmt.Sprintf("audit ledger entry_count mismatch: entry_count=%d entries=%d", audit.EntryCount, len(audit.Entries)))
	}

	manifestByPackage := manifestPackages(m)
	entryCounts := map[string]map[string]int{}
	seen := map[string]bool{}
	for _, entry := range audit.Entries {
		if seen[entry.ID] {
			failures = append(failures, "duplicate audit ledger entry "+entry.ID)
			continue
		}
		seen[entry.ID] = true

		manifestEntry, ok := manifestByPackage[entry.Package]
		if !ok {
			failures = append(failures, "audit ledger entry references package outside manifest "+entry.ID)
			continue
		}
		if entry.Kind != "top" && entry.Kind != "field" && entry.Kind != "method" {
			failures = append(failures, fmt.Sprintf("audit ledger entry %s has unknown kind %q", entry.ID, entry.Kind))
		}
		if entry.Symbol == "" || entry.Signature == "" || entry.Disposition == "" {
			failures = append(failures, "audit ledger entry missing required metadata "+entry.ID)
		}
		failures = append(failures, verifyCoverageSubset(entry.ID, entry.Coverage, manifestEntry.Coverage)...)
		failures = append(failures, verifyCoveragePaths(docsRoot, entry.ID, entry.Coverage)...)
		failures = append(failures, verifyPublicIndexSignature(englishIndex, "English public API index", entry)...)
		failures = append(failures, verifyPublicIndexSignature(chineseIndex, "Chinese public API index", entry)...)

		if entryCounts[entry.Package] == nil {
			entryCounts[entry.Package] = map[string]int{}
		}
		entryCounts[entry.Package][entry.Kind]++
	}

	for _, entry := range m.Packages {
		counts := entryCounts[entry.Package]
		gotEntries := counts["top"] + counts["field"] + counts["method"]
		wantEntries := entry.TopLevel + entry.Fields + entry.Methods
		if counts["top"] != entry.TopLevel || counts["field"] != entry.Fields ||
			counts["method"] != entry.Methods || gotEntries != wantEntries {
			failures = append(failures, fmt.Sprintf(
				"audit ledger package count mismatch for %s: got top/field/method/entries=%d/%d/%d/%d want=%d/%d/%d/%d",
				entry.Package,
				counts["top"], counts["field"], counts["method"], gotEntries,
				entry.TopLevel, entry.Fields, entry.Methods, wantEntries,
			))
		}
	}

	return failures
}

func verifyGroupedPublicSurface(audit auditLedger) []string {
	var rows []string
	for _, entry := range audit.Entries {
		if !strings.HasPrefix(entry.Disposition, "grouped:") {
			continue
		}
		rows = append(rows, strings.Join([]string{
			entry.Package,
			entry.Symbol,
			entry.Kind,
			entry.Disposition,
			entry.Signature,
		}, "\t"))
	}

	gotFingerprint := fingerprintRows(rows)
	var failures []string
	if len(rows) != groupedPublicEntryCount {
		failures = append(failures, fmt.Sprintf("grouped public API surface count mismatch: got %d want %d", len(rows), groupedPublicEntryCount))
	}
	if gotFingerprint != groupedPublicSignatureFingerprint {
		failures = append(failures, fmt.Sprintf("grouped public API surface fingerprint mismatch: got %s want %s", gotFingerprint, groupedPublicSignatureFingerprint))
	}

	return failures
}

func fingerprintRows(rows []string) string {
	sorted := append([]string(nil), rows...)
	sort.Strings(sorted)

	hash := sha256.New()
	for _, row := range sorted {
		hash.Write([]byte(row))
		hash.Write([]byte("\n"))
	}

	return hex.EncodeToString(hash.Sum(nil))
}

func verifyContractLedgerClosure(contracts contractLedger, m manifest, docsRoot, sourceRoot string) []string {
	var failures []string
	if contracts.SourceModule != m.SourceModule {
		failures = append(failures, "contract ledger source_module mismatch: "+contracts.SourceModule)
	}

	manifestByPackage := manifestPackages(m)
	contractByID := map[string]contractLedgerEntry{}
	for _, entry := range contracts.Entries {
		if _, ok := contractByID[entry.ID]; ok {
			failures = append(failures, "duplicate contract ledger entry "+entry.ID)
			continue
		}
		contractByID[entry.ID] = entry
		manifestEntry, ok := manifestByPackage[entry.Package]
		if !ok {
			failures = append(failures, "contract ledger entry references package outside manifest "+entry.ID)
			continue
		}
		if entry.Kind == "" || len(entry.Terms) == 0 || len(entry.SourceEvidence) == 0 {
			failures = append(failures, "contract ledger entry missing required closure metadata "+entry.ID)
		}
		failures = append(failures, verifyCoverageSubset(entry.ID, entry.Coverage, manifestEntry.Coverage)...)
		failures = append(failures, verifyCoveragePaths(docsRoot, entry.ID, entry.Coverage)...)
		failures = append(failures, verifySourceEvidence(sourceRoot, entry.ID, entry.SourceEvidence)...)
		failures = append(failures, verifySourceEvidence(sourceRoot, entry.ID, entry.TestEvidence)...)
	}

	reviewedPackages := map[string]bool{}
	reviewedContracts := map[string]bool{}
	for _, review := range contracts.PackageReviews {
		if reviewedPackages[review.Package] {
			failures = append(failures, "duplicate contract package review "+review.Package)
			continue
		}
		reviewedPackages[review.Package] = true

		manifestEntry, ok := manifestByPackage[review.Package]
		if !ok {
			failures = append(failures, "contract package review references package outside manifest "+review.Package)
			continue
		}
		wantEntries := manifestEntry.TopLevel + manifestEntry.Fields + manifestEntry.Methods
		if review.ReviewedSurface.TopLevel != manifestEntry.TopLevel ||
			review.ReviewedSurface.Fields != manifestEntry.Fields ||
			review.ReviewedSurface.Methods != manifestEntry.Methods ||
			review.ReviewedSurface.EntryCount != wantEntries ||
			review.ReviewedSurface.Fingerprint != manifestEntry.Fingerprint {
			failures = append(failures, fmt.Sprintf("contract package review surface mismatch for %s", review.Package))
		}
		failures = append(failures, verifyCoverageSubset(review.Package, review.Coverage, manifestEntry.Coverage)...)
		failures = append(failures, verifyCoveragePaths(docsRoot, review.Package, review.Coverage)...)
		failures = append(failures, verifySourceEvidence(sourceRoot, review.Package, review.SourceEvidence)...)
		if review.Disposition == "has-semantic-contracts" && len(review.ContractIDs) == 0 {
			failures = append(failures, "semantic package review missing contract ids "+review.Package)
		}
		for _, id := range review.ContractIDs {
			reviewedContracts[id] = true
			contract, ok := contractByID[id]
			if !ok {
				failures = append(failures, fmt.Sprintf("package review %s references unknown contract %s", review.Package, id))
				continue
			}
			if contract.Package != review.Package {
				failures = append(failures, fmt.Sprintf("package review %s references contract from %s: %s", review.Package, contract.Package, id))
			}
		}
	}

	for _, entry := range m.Packages {
		if !reviewedPackages[entry.Package] {
			failures = append(failures, "manifest package missing contract package review "+entry.Package)
		}
	}
	for _, contract := range contracts.Entries {
		if !reviewedContracts[contract.ID] {
			failures = append(failures, "contract entry not linked from package review "+contract.ID)
		}
	}

	return failures
}

func verifyPublicIndexSignature(indexContent, label string, entry auditLedgerEntry) []string {
	section := packageSection(indexContent, entry.Package)
	if section == "" {
		return []string{label + " missing package section for " + entry.Package}
	}
	if !strings.Contains(section, entry.Signature) {
		return []string{fmt.Sprintf("%s missing signature for %s: %s", label, entry.ID, entry.Signature)}
	}

	return nil
}

func verifyCoverageSubset(entryID string, got, allowed []string) []string {
	if len(got) == 0 {
		return []string{"missing coverage for " + entryID}
	}

	allowedSet := map[string]bool{}
	for _, item := range allowed {
		allowedSet[item] = true
	}

	var failures []string
	for _, item := range got {
		if !strings.Contains(item, "://") && !allowedSet[item] {
			failures = append(failures, fmt.Sprintf("%s uses coverage outside package manifest: %s", entryID, item))
		}
	}

	return failures
}

func verifyCoveragePaths(docsRoot, entryID string, coverage []string) []string {
	var failures []string
	for _, item := range coverage {
		if strings.Contains(item, "://") {
			continue
		}
		path := filepath.Join(docsRoot, item)
		if _, err := os.Stat(path); err != nil {
			failures = append(failures, fmt.Sprintf("%s coverage file missing or unreadable: %s: %v", entryID, item, err))
		}
		if zhPath, ok := chineseMirrorPath(item); ok {
			path := filepath.Join(docsRoot, zhPath)
			if _, err := os.Stat(path); err != nil {
				failures = append(failures, fmt.Sprintf("%s Chinese coverage file missing or unreadable: %s: %v", entryID, zhPath, err))
			}
		}
	}

	return failures
}

func verifySourceEvidence(sourceRoot, entryID string, evidence []string) []string {
	var failures []string
	for _, item := range evidence {
		if strings.Contains(item, "://") {
			continue
		}
		path, line, ok := strings.Cut(item, ":")
		if !ok || line == "" {
			failures = append(failures, "source evidence missing line number for "+entryID+": "+item)
			continue
		}
		if _, err := os.Stat(filepath.Join(sourceRoot, path)); err != nil {
			failures = append(failures, fmt.Sprintf("source evidence file missing for %s: %s: %v", entryID, item, err))
		}
	}

	return failures
}

func loadLiveEntries(sourceRoot, docsRoot, manifestPath, auditLedgerPath, contractLedgerPath string) map[string]manifestEntry {
	cmd := exec.Command(
		"go", "run", filepath.Join(docsRoot, "scripts/verify-api-audit.go"),
		"-source", sourceRoot,
		"-manifest", manifestPath,
		"-ledger", auditLedgerPath,
		"-contract-ledger", contractLedgerPath,
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
		panic(fmt.Errorf("failed to parse live inventory: %w", err))
	}

	result := map[string]manifestEntry{}
	for _, entry := range entries {
		result[entry.Package] = entry
	}

	return result
}

func manifestPackages(m manifest) map[string]manifestEntry {
	result := make(map[string]manifestEntry, len(m.Packages))
	for _, entry := range m.Packages {
		result[entry.Package] = entry
	}

	return result
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
