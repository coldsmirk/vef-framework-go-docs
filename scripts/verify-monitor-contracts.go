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
	"reflect"
	"sort"
	"strconv"
	"strings"
)

const (
	monitorPackage = "github.com/coldsmirk/vef-framework-go/monitor"

	monitorFingerprint = "2dcd6b85f44ba51532759515b9641c9584332b083cddbcf310ad045b7e1ea15c"
	monitorTopLevel    = 27
	monitorFields      = 181
	monitorMethods     = 9
	monitorEntries     = 217

	monitorGroupedEntries              = 190
	monitorGroupedFields               = 181
	monitorGroupedMethods              = 9
	monitorGroupedReceivers            = 23
	monitorGroupedSignatureFingerprint = "323c1b9612fb0905a31851a622199f0d0c1c67e7f466935109bcec70d3da76f9"
	monitorGroupedReceiverFingerprint  = "777aa0b0f78aad775efc5bf8f890af7fb30ec5749424c7919d11aeed39b23fb2"

	englishMonitorPath = "docs/infrastructure/monitor.md"
	chineseMonitorPath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/infrastructure/monitor.md"
	englishBuiltInPath = "docs/reference/built-in-resources.md"
	chineseBuiltInPath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/built-in-resources.md"
)

type corpus struct {
	label   string
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
	Coverage       []string `json:"coverage"`
	SourceEvidence []string `json:"source_evidence"`
	Terms          []string `json:"terms"`
}

type runtimeLedger struct {
	Entries []runtimeEntry `json:"entries"`
}

type runtimeEntry struct {
	Category string   `json:"category"`
	Name     string   `json:"name"`
	Value    string   `json:"value"`
	Details  []string `json:"details"`
}

func main() {
	sourceDir := flag.String("source", "../vef-framework-go", "path to the VEF Framework Go source repository")
	outDir := flag.String("out", ".", "path to the VEF Framework Go docs repository")
	flag.Parse()

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)

	englishMonitor := readCorpus("English monitor docs", filepath.Join(docsRoot, englishMonitorPath))
	chineseMonitor := readCorpus("Chinese monitor docs", filepath.Join(docsRoot, chineseMonitorPath))
	englishBuiltIn := readCorpus("English built-in resources docs", filepath.Join(docsRoot, englishBuiltInPath))
	chineseBuiltIn := readCorpus("Chinese built-in resources docs", filepath.Join(docsRoot, chineseBuiltInPath))

	audit := loadJSON[auditLedger](filepath.Join(docsRoot, "scripts/api-audit-ledger.json"))
	manifestData := loadJSON[manifest](filepath.Join(docsRoot, "scripts/api-audit-manifest.json"))
	contracts := loadJSON[contractLedger](filepath.Join(docsRoot, "scripts/api-contract-ledger.json"))
	runtime := loadJSON[runtimeLedger](filepath.Join(docsRoot, "scripts/runtime-api-ledger.json"))

	entries := monitorEntriesFromAudit(audit)
	var failures []string
	failures = append(failures, verifyManifest(manifestData)...)
	failures = append(failures, verifyContractLedger(contracts, sourceRoot)...)
	failures = append(failures, verifyAuditEntries(entries)...)
	failures = append(failures, verifyGroupedMonitorSurface(entries)...)
	failures = append(failures, verifyRuntimeActions(runtime, []corpus{englishMonitor, chineseMonitor, englishBuiltIn, chineseBuiltIn})...)
	failures = append(failures, verifyMonitorDocs(entries, []corpus{englishMonitor, chineseMonitor})...)
	failures = append(failures, verifySourceTerms(sourceRoot)...)
	failures = append(failures, runGoTests(sourceRoot)...)

	if len(failures) > 0 {
		sort.Strings(failures)
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Println("Monitor contract docs verified: sys/monitor runtime actions")
}

func verifyManifest(m manifest) []string {
	for _, entry := range m.Packages {
		if entry.Package == monitorPackage {
			return verifySurface("manifest", entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint)
		}
	}

	return []string{"manifest missing monitor package"}
}

func verifyContractLedger(contracts contractLedger, sourceRoot string) []string {
	var failures []string
	var foundReview bool
	for _, review := range contracts.PackageReviews {
		if review.Package != monitorPackage {
			continue
		}
		foundReview = true
		if review.Disposition != "has-semantic-contracts" {
			failures = append(failures, "monitor contract review disposition mismatch: "+review.Disposition)
		}
		failures = append(failures, verifySurface("contract review", review.ReviewedSurface.TopLevel, review.ReviewedSurface.Fields, review.ReviewedSurface.Methods, review.ReviewedSurface.Fingerprint)...)
		if review.ReviewedSurface.EntryCount != monitorEntries {
			failures = append(failures, fmt.Sprintf("monitor contract review entry_count=%d want=%d", review.ReviewedSurface.EntryCount, monitorEntries))
		}
		if !contains(review.ContractIDs, monitorContractID()) {
			failures = append(failures, "monitor contract review missing contract id "+monitorContractID())
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, review.SourceEvidence)...)
	}
	if !foundReview {
		failures = append(failures, "contract ledger missing monitor package review")
	}

	var foundContract bool
	for _, entry := range contracts.Entries {
		if entry.ID != monitorContractID() {
			continue
		}
		foundContract = true
		if entry.Package != monitorPackage || entry.Kind != "dynamic-resource" {
			failures = append(failures, fmt.Sprintf("monitor contract entry shape mismatch: package=%s kind=%s", entry.Package, entry.Kind))
		}
		for _, term := range []string{"sys/monitor", "get_overview", "get_build_info", "60", "ErrNotReady", "ErrCollectionFailed"} {
			if !contains(entry.Terms, term) {
				failures = append(failures, "monitor contract missing term "+term)
			}
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, entry.SourceEvidence)...)
	}
	if !foundContract {
		failures = append(failures, "contract ledger missing monitor contract entry")
	}

	return failures
}

func verifySurface(label string, topLevel, fields, methods int, fingerprint string) []string {
	if topLevel == monitorTopLevel && fields == monitorFields && methods == monitorMethods && fingerprint == monitorFingerprint {
		return nil
	}

	return []string{fmt.Sprintf(
		"%s monitor surface mismatch: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
		label, topLevel, fields, methods, fingerprint, monitorTopLevel, monitorFields, monitorMethods, monitorFingerprint,
	)}
}

func verifyAuditEntries(entries []auditEntry) []string {
	var failures []string
	if len(entries) != monitorEntries {
		failures = append(failures, fmt.Sprintf("monitor audit entry count mismatch: got %d want %d", len(entries), monitorEntries))
	}
	counts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if seen[entry.ID] {
			failures = append(failures, "duplicate monitor audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		if entry.Symbol == "" || entry.Signature == "" || entry.Disposition == "" {
			failures = append(failures, "monitor audit entry missing metadata "+entry.ID)
		}
	}
	if counts["top"] != monitorTopLevel || counts["field"] != monitorFields || counts["method"] != monitorMethods {
		failures = append(failures, fmt.Sprintf("monitor audit kind counts mismatch: top/field/method=%d/%d/%d", counts["top"], counts["field"], counts["method"]))
	}

	return failures
}

func verifyGroupedMonitorSurface(entries []auditEntry) []string {
	var rows []string
	receiverCounts := map[string]int{}
	kindCounts := map[string]int{}
	var failures []string
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Disposition, "grouped:") {
			continue
		}
		receiver, ok := receiverForSymbol(entry.Symbol)
		if !ok {
			failures = append(failures, fmt.Sprintf("monitor grouped entry has non receiver-qualified symbol %q", entry.Symbol))
			continue
		}
		rows = append(rows, strings.Join([]string{entry.Symbol, entry.Kind, entry.Disposition, entry.Signature}, "\t"))
		receiverCounts[receiver]++
		kindCounts[entry.Kind]++
	}

	failures = append(failures, verifyGroupedFingerprint("monitor grouped DTO/service surface", rows, monitorGroupedEntries, monitorGroupedSignatureFingerprint)...)
	if kindCounts["field"] != monitorGroupedFields || kindCounts["method"] != monitorGroupedMethods {
		failures = append(failures, fmt.Sprintf("monitor grouped kind counts mismatch: got fields/methods=%d/%d want=%d/%d", kindCounts["field"], kindCounts["method"], monitorGroupedFields, monitorGroupedMethods))
	}

	var receiverRows []string
	for receiver, count := range receiverCounts {
		receiverRows = append(receiverRows, fmt.Sprintf("%d %s", count, receiver))
	}
	failures = append(failures, verifyGroupedFingerprint("monitor grouped receiver families", receiverRows, monitorGroupedReceivers, monitorGroupedReceiverFingerprint)...)

	return failures
}

func verifyRuntimeActions(runtime runtimeLedger, docs []corpus) []string {
	actions := map[string]bool{}
	var resourceFound bool
	for _, entry := range runtime.Entries {
		if entry.Category == "built-in resource" && entry.Name == "sys/monitor" && entry.Value == "rpc" {
			resourceFound = true
		}
		if entry.Category == "built-in resource action" && strings.HasPrefix(entry.Name, "sys/monitor/") {
			actions[strings.TrimPrefix(entry.Name, "sys/monitor/")] = true
		}
	}

	wantActions := []string{
		"get_overview",
		"get_cpu",
		"get_memory",
		"get_disk",
		"get_network",
		"get_host",
		"get_process",
		"get_load",
		"get_build_info",
	}
	var failures []string
	if !resourceFound {
		failures = append(failures, "runtime ledger missing sys/monitor built-in resource")
	}
	failures = append(failures, compareSets("runtime sys/monitor actions", sortedKeys(actions), wantActions)...)
	for _, doc := range docs {
		if !strings.Contains(doc.content, "`sys/monitor`") {
			failures = append(failures, doc.label+" missing sys/monitor")
		}
		for _, action := range wantActions {
			if !strings.Contains(doc.content, "`"+action+"`") {
				failures = append(failures, doc.label+" missing monitor action "+action)
			}
		}
	}

	return failures
}

func verifyMonitorDocs(entries []auditEntry, docs []corpus) []string {
	receivers := map[string]bool{}
	jsonFields := map[string]bool{}
	for _, entry := range entries {
		receiver, ok := receiverForSymbol(entry.Symbol)
		if ok {
			receivers[receiver] = true
		}
		if entry.Kind == "field" {
			if jsonName := jsonNameFromSignature(entry.Signature); jsonName != "" {
				jsonFields[jsonName] = true
			}
		}
	}

	var failures []string
	for _, doc := range docs {
		for _, receiver := range sortedKeys(receivers) {
			if !strings.Contains(doc.content, "`monitor."+receiver+"`") {
				failures = append(failures, doc.label+" missing monitor receiver/type "+receiver)
			}
		}
		for _, field := range sortedKeys(jsonFields) {
			if !strings.Contains(doc.content, "`"+field+"`") {
				failures = append(failures, doc.label+" missing monitor JSON field "+field)
			}
		}
		for _, term := range []string{
			"`monitor.ErrNotReady`",
			"`ErrCodeNotReady`",
			"`monitor.ErrCollectionFailed`",
			"`ErrCodeCollectionFailed`",
			"`unknown`",
			"`vefVersion`",
			"`10s`",
			"`2s`",
			"60",
		} {
			if !strings.Contains(doc.content, term) {
				failures = append(failures, doc.label+" missing monitor semantic term "+term)
			}
		}
		if strings.Contains(doc.content, "`v0.0.0`") || strings.Contains(doc.content, "`2022-08-08 01:00:00`") {
			failures = append(failures, doc.label+" still documents stale monitor build-info fallback")
		}
	}

	return failures
}

func verifySourceTerms(sourceRoot string) []string {
	checks := []struct {
		path  string
		terms []string
	}{
		{
			path: "monitor/service.go",
			terms: []string{
				"type Service interface",
				"BuildInfo() *BuildInfo",
				"`json:\"vefVersion\"`",
				"`json:\"appVersion\"`",
				"`json:\"gitCommit\"`",
			},
		},
		{
			path: "monitor/api_errors.go",
			terms: []string{
				"ErrCodeNotReady",
				"ErrCodeCollectionFailed",
				"ErrNotReady",
				"ErrCollectionFailed",
			},
		},
		{
			path: "internal/monitor/resource.go",
			terms: []string{
				"api.NewRPCResource(",
				"\"sys/monitor\"",
				"var defaultRateLimit = &api.RateLimitConfig{Max: 60}",
				"Action: \"get_build_info\"",
				"return monitor.ErrNotReady",
				"return monitor.ErrCollectionFailed",
			},
		},
		{
			path: "internal/monitor/config.go",
			terms: []string{
				"DefaultSampleInterval = 10 * time.Second",
				"DefaultSampleDuration = 2 * time.Second",
			},
		},
		{
			path: "internal/monitor/service.go",
			terms: []string{
				"AppVersion: \"unknown\"",
				"BuildTime:  \"unknown\"",
				"GitCommit:  \"unknown\"",
				"buildInfo.VEFVersion = version.VEFVersion",
				"if cfg.SampleInterval > 0",
				"if cfg.SampleDuration > 0",
				"overview.Build = s.BuildInfo()",
				"return nil, ErrCPUInfoNotReady",
				"return nil, ErrProcessInfoNotReady",
			},
		},
		{
			path: "internal/monitor/service_internal_test.go",
			terms: []string{
				"assert.Equal(t, \"unknown\", got.AppVersion",
				"assert.Equal(t, \"unknown\", got.BuildTime",
				"assert.Equal(t, \"unknown\", got.GitCommit",
				"assert.Equal(t, version.VEFVersion, got.VEFVersion",
			},
		},
	}

	var failures []string
	for _, check := range checks {
		source := readCorpus(check.path, filepath.Join(sourceRoot, check.path))
		for _, term := range check.terms {
			if !strings.Contains(source.content, term) {
				failures = append(failures, source.label+" missing source term "+term)
			}
		}
	}

	return failures
}

func runGoTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./monitor", "./internal/monitor")
	cmd.Dir = sourceRoot
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return []string{fmt.Sprintf("go test ./monitor ./internal/monitor failed: %v\n%s", err, strings.TrimSpace(output.String()))}
	}

	return nil
}

func monitorEntriesFromAudit(audit auditLedger) []auditEntry {
	var entries []auditEntry
	for _, entry := range audit.Entries {
		if entry.Package == monitorPackage {
			entries = append(entries, entry)
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].ID < entries[j].ID })

	return entries
}

func verifySourceEvidence(sourceRoot string, evidence []string) []string {
	var failures []string
	for _, item := range evidence {
		path, lineText, ok := strings.Cut(item, ":")
		if !ok {
			failures = append(failures, "bad source evidence format "+item)
			continue
		}
		if _, err := strconv.Atoi(lineText); err != nil {
			failures = append(failures, "bad source evidence line "+item)
		}
		if _, err := os.Stat(filepath.Join(sourceRoot, path)); err != nil {
			failures = append(failures, "source evidence missing "+item)
		}
	}

	return failures
}

func verifyGroupedFingerprint(label string, rows []string, wantCount int, wantFingerprint string) []string {
	gotFingerprint := fingerprintRows(rows)
	var failures []string
	if len(rows) != wantCount {
		failures = append(failures, fmt.Sprintf("%s count mismatch: got %d want %d", label, len(rows), wantCount))
	}
	if gotFingerprint != wantFingerprint {
		failures = append(failures, fmt.Sprintf("%s fingerprint mismatch: got %s want %s", label, gotFingerprint, wantFingerprint))
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

func jsonNameFromSignature(signature string) string {
	tagMarker := " tag=\""
	start := strings.Index(signature, tagMarker)
	if start < 0 {
		return ""
	}
	rest := signature[start+len(tagMarker):]
	end := strings.LastIndex(rest, "\"]")
	if end < 0 {
		return ""
	}
	tag := strings.ReplaceAll(rest[:end], `\"`, `"`)
	jsonTag := reflect.StructTag(tag).Get("json")
	name, _, _ := strings.Cut(jsonTag, ",")
	if name == "" || name == "-" {
		return ""
	}

	return name
}

func receiverForSymbol(symbol string) (string, bool) {
	receiver, _, ok := strings.Cut(symbol, ".")
	if !ok || receiver == "" {
		return "", false
	}

	return receiver, true
}

func monitorContractID() string {
	return monitorPackage + "#dynamic-resource:sys-monitor-resource-dtos"
}

func readCorpus(label, path string) corpus {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read %s at %s: %w", label, path, err))
	}

	return corpus{label: label, content: string(content)}
}

func loadJSON[T any](path string) T {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var result T
	if err := json.Unmarshal(content, &result); err != nil {
		panic(err)
	}

	return result
}

func compareSets(label string, got, want []string) []string {
	got = sortedUnique(got)
	want = sortedUnique(want)
	var failures []string
	gotSet := sliceSet(got)
	wantSet := sliceSet(want)
	for _, item := range want {
		if !gotSet[item] {
			failures = append(failures, label+" missing "+item)
		}
	}
	for _, item := range got {
		if !wantSet[item] {
			failures = append(failures, label+" has phantom "+item)
		}
	}

	return failures
}

func sortedKeys[V any](m map[string]V) []string {
	result := make([]string, 0, len(m))
	for key := range m {
		result = append(result, key)
	}
	sort.Strings(result)

	return result
}

func sortedUnique(values []string) []string {
	set := sliceSet(values)
	result := make([]string, 0, len(set))
	for value := range set {
		result = append(result, value)
	}
	sort.Strings(result)

	return result
}

func sliceSet(values []string) map[string]bool {
	result := map[string]bool{}
	for _, value := range values {
		result[value] = true
	}

	return result
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}

	return false
}

func containsNormalized(content, term string) bool {
	return strings.Contains(content, term) ||
		strings.Contains(strings.Join(strings.Fields(content), " "), strings.Join(strings.Fields(term), " "))
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}
