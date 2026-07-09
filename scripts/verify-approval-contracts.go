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
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	approvalDocsPath        = "docs/approval.md"
	chineseApprovalDocsPath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/approval.md"
	runtimeLedgerPath       = "scripts/runtime-api-ledger.json"
	manifestPath            = "scripts/api-audit-manifest.json"
	contractLedgerPath      = "scripts/api-contract-ledger.json"
	auditLedgerPath         = "scripts/api-audit-ledger.json"
)

type corpus struct {
	label   string
	content string
}

type manifest struct {
	SourceModule string          `json:"source_module"`
	Packages     []manifestEntry `json:"packages"`
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
	SourceModule   string                  `json:"source_module"`
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
	SourceModule string         `json:"source_module"`
	Entries      []runtimeEntry `json:"entries"`
}

type auditLedger struct {
	Entries []auditEntry `json:"entries"`
}

type auditEntry struct {
	Package     string `json:"package"`
	Kind        string `json:"kind"`
	Symbol      string `json:"symbol"`
	Signature   string `json:"signature"`
	Disposition string `json:"disposition"`
}

type runtimeEntry struct {
	Category string   `json:"category"`
	Name     string   `json:"name"`
	Value    string   `json:"value"`
	Details  []string `json:"details"`
}

type liveInventoryEntry struct {
	Package     string   `json:"package"`
	Coverage    []string `json:"coverage"`
	TopLevel    int      `json:"top_level"`
	Fields      int      `json:"fields"`
	Methods     int      `json:"methods"`
	Fingerprint string   `json:"fingerprint"`
}

type expectedPackage struct {
	pkg         string
	topLevel    int
	fields      int
	methods     int
	entries     int
	fingerprint string
	contracts   []string
}

type groupedSurface struct {
	pkg                  string
	entryCount           int
	fieldCount           int
	methodCount          int
	receiverCount        int
	signatureFingerprint string
	receiverFingerprint  string
}

type stringConst struct {
	name     string
	typeName string
	value    string
}

type transition struct {
	kind string
	from string
	to   string
}

type eventContract struct {
	constName   string
	topic       string
	structName  string
	constructor string
	fields      []string
}

type errorContract struct {
	errName   string
	codeName  string
	codeValue string
	message   string
}

var approvalPackages = []expectedPackage{
	{
		pkg:         "github.com/coldsmirk/vef-framework-go/approval",
		topLevel:    378,
		fields:      823,
		methods:     132,
		entries:     1333,
		fingerprint: "d9c2c0bd0a05bfc880bd85055e097223a5570839220d00d9cf2b9e7bc5648519",
		contracts: []string{
			"github.com/coldsmirk/vef-framework-go/approval#dynamic-resource:approval-built-in-resources",
			"github.com/coldsmirk/vef-framework-go/approval#event-contract:approval-domain-events",
			"github.com/coldsmirk/vef-framework-go/approval#result-contract:approval-error-envelope",
			"github.com/coldsmirk/vef-framework-go/approval#runtime-contract:tenant-and-business-binding",
			"github.com/coldsmirk/vef-framework-go/approval#wire-contract:flow-definition-and-form-json",
		},
	},
	{
		pkg:         "github.com/coldsmirk/vef-framework-go/approval/admin",
		topLevel:    6,
		fields:      60,
		methods:     0,
		entries:     66,
		fingerprint: "1c843a8f8400f22c84a8fe843ad7638ca4ad20c5cf48237cf5b835f1628d9f3e",
		contracts: []string{
			"github.com/coldsmirk/vef-framework-go/approval/admin#dto-wire-shape:approval-admin-dtos",
		},
	},
	{
		pkg:         "github.com/coldsmirk/vef-framework-go/approval/my",
		topLevel:    8,
		fields:      68,
		methods:     0,
		entries:     76,
		fingerprint: "e6ad65e45fb72650108d7289633637687116fb63a06bf797d73073d34205d681",
		contracts: []string{
			"github.com/coldsmirk/vef-framework-go/approval/my#dto-wire-shape:approval-my-dtos",
		},
	},
}

var approvalGroupedSurfaces = []groupedSurface{
	{
		pkg:                  "github.com/coldsmirk/vef-framework-go/approval",
		entryCount:           955,
		fieldCount:           823,
		methodCount:          132,
		receiverCount:        123,
		signatureFingerprint: "34b0824567611923b21fcb74850f57977ca8620deeba5c896af7377315dbebd0",
		receiverFingerprint:  "8e1d8ecef9a98f7e87c581cabed6222b0120ec964693097bac7640f613cc585a",
	},
	{
		pkg:                  "github.com/coldsmirk/vef-framework-go/approval/admin",
		entryCount:           60,
		fieldCount:           60,
		methodCount:          0,
		receiverCount:        6,
		signatureFingerprint: "fcb429a9dc59a49d9a371acb008118c580825d4df0b157242dbca33c3e95c644",
		receiverFingerprint:  "0695f298d8c00d3acb00ad14407dab7f29b4cdfaf1f659006f8ca54e4d6bb39d",
	},
	{
		pkg:                  "github.com/coldsmirk/vef-framework-go/approval/my",
		entryCount:           68,
		fieldCount:           68,
		methodCount:          0,
		receiverCount:        8,
		signatureFingerprint: "1427ec9ba629b3fa01f49cdefac1d4243d881283c6b15ce7377e1687023c97cd",
		receiverFingerprint:  "720d08f689e8c5f2fbe676ffb985a37b1b723dd0f5a62773ee646385fe9d28c4",
	},
}

func main() {
	sourceDir := flag.String("source", "../vef-framework-go", "path to the VEF Framework Go source repository")
	outDir := flag.String("out", ".", "path to the VEF Framework Go docs repository")
	flag.Parse()

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)

	english := readCorpus("English approval docs", filepath.Join(docsRoot, approvalDocsPath))
	chinese := readCorpus("Chinese approval docs", filepath.Join(docsRoot, chineseApprovalDocsPath))
	docs := []corpus{english, chinese}

	m := loadJSON[manifest](filepath.Join(docsRoot, manifestPath))
	contracts := loadJSON[contractLedger](filepath.Join(docsRoot, contractLedgerPath))
	runtime := loadJSON[runtimeLedger](filepath.Join(docsRoot, runtimeLedgerPath))
	audit := loadJSON[auditLedger](filepath.Join(docsRoot, auditLedgerPath))
	liveInventory := loadLiveInventory(sourceRoot, docsRoot)

	var failures []string
	failures = append(failures, verifyPackageSurfaces(m, contracts, liveInventory)...)
	failures = append(failures, verifyGroupedApprovalSurfaces(audit, docs)...)
	failures = append(failures, verifyRuntimeResources(runtime, docs)...)

	constsByName, constsByType := extractStringConstants(sourceRoot)
	failures = append(failures, verifyApprovalEnums(constsByType, docs)...)
	failures = append(failures, verifyStateMachines(sourceRoot, constsByName, docs)...)
	failures = append(failures, verifyRequestDTOs(sourceRoot, docs)...)
	failures = append(failures, verifyResponseDTOs(sourceRoot, docs)...)
	failures = append(failures, verifyFlowAndFormWireShapes(sourceRoot, docs)...)
	failures = append(failures, verifyEvents(sourceRoot, constsByName, docs)...)
	failures = append(failures, verifyErrors(sourceRoot, docs)...)
	failures = append(failures, verifyResubmitAndAvailableActions(sourceRoot, docs)...)
	failures = append(failures, verifyTenantBindingAndExtensionTerms(docs)...)
	failures = append(failures, runSourceTests(sourceRoot)...)

	if len(failures) > 0 {
		sort.Strings(failures)
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Printf("Approval contract docs verified: 3 public packages, 1059 grouped field/method entries, %d runtime approval resources/actions, source-derived state/event/error/DTO contracts, 2 doc mirrors\n", approvalRuntimeSurfaceCount(runtime))
}

func verifyPackageSurfaces(m manifest, contracts contractLedger, live map[string]liveInventoryEntry) []string {
	var failures []string
	if m.SourceModule != "github.com/coldsmirk/vef-framework-go" {
		failures = append(failures, "manifest source_module mismatch: "+m.SourceModule)
	}
	if contracts.SourceModule != m.SourceModule {
		failures = append(failures, "contract ledger source_module mismatch: "+contracts.SourceModule)
	}

	manifestByPackage := map[string]manifestEntry{}
	for _, entry := range m.Packages {
		manifestByPackage[entry.Package] = entry
	}
	reviewByPackage := map[string]contractPackageReview{}
	for _, review := range contracts.PackageReviews {
		reviewByPackage[review.Package] = review
	}
	contractByID := map[string]contractEntry{}
	for _, entry := range contracts.Entries {
		contractByID[entry.ID] = entry
	}

	for _, expected := range approvalPackages {
		liveEntry, ok := live[expected.pkg]
		if !ok {
			failures = append(failures, "live public inventory missing "+expected.pkg)
			continue
		}
		failures = append(failures, verifySurface("live inventory", expected, liveEntry.TopLevel, liveEntry.Fields, liveEntry.Methods, liveEntry.Fingerprint)...)

		manifestEntry, ok := manifestByPackage[expected.pkg]
		if !ok {
			failures = append(failures, "manifest missing "+expected.pkg)
			continue
		}
		failures = append(failures, verifySurface("manifest", expected, manifestEntry.TopLevel, manifestEntry.Fields, manifestEntry.Methods, manifestEntry.Fingerprint)...)
		if !contains(manifestEntry.Coverage, approvalDocsPath) {
			failures = append(failures, expected.pkg+" manifest coverage must include "+approvalDocsPath)
		}

		review, ok := reviewByPackage[expected.pkg]
		if !ok {
			failures = append(failures, "contract package review missing "+expected.pkg)
			continue
		}
		if review.Disposition != "has-semantic-contracts" {
			failures = append(failures, expected.pkg+" review disposition = "+review.Disposition)
		}
		failures = append(failures, verifySurface(
			"contract review",
			expected,
			review.ReviewedSurface.TopLevel,
			review.ReviewedSurface.Fields,
			review.ReviewedSurface.Methods,
			review.ReviewedSurface.Fingerprint,
		)...)
		if review.ReviewedSurface.EntryCount != expected.entries {
			failures = append(failures, fmt.Sprintf("%s contract review entry_count=%d want=%d", expected.pkg, review.ReviewedSurface.EntryCount, expected.entries))
		}
		for _, id := range expected.contracts {
			if !contains(review.ContractIDs, id) {
				failures = append(failures, expected.pkg+" review missing contract id "+id)
			}
			entry, ok := contractByID[id]
			if !ok {
				failures = append(failures, "contract ledger missing "+id)
				continue
			}
			if entry.Package != expected.pkg {
				failures = append(failures, id+" package mismatch: "+entry.Package)
			}
			if !contains(entry.Coverage, approvalDocsPath) {
				failures = append(failures, id+" coverage must include "+approvalDocsPath)
			}
			if len(entry.SourceEvidence) == 0 || len(entry.Terms) == 0 {
				failures = append(failures, id+" missing source evidence or terms")
			}
		}
	}

	return failures
}

func verifySurface(label string, expected expectedPackage, topLevel, fields, methods int, fingerprint string) []string {
	if topLevel == expected.topLevel && fields == expected.fields && methods == expected.methods && fingerprint == expected.fingerprint {
		return nil
	}

	return []string{fmt.Sprintf(
		"%s surface drift for %s: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
		label,
		expected.pkg,
		topLevel, fields, methods, fingerprint,
		expected.topLevel, expected.fields, expected.methods, expected.fingerprint,
	)}
}

func verifyGroupedApprovalSurfaces(audit auditLedger, docs []corpus) []string {
	surfacesByPackage := map[string]groupedSurface{}
	for _, surface := range approvalGroupedSurfaces {
		surfacesByPackage[surface.pkg] = surface
	}

	rowsByPackage := map[string][]string{}
	receiverCountsByPackage := map[string]map[string]int{}
	kindCountsByPackage := map[string]map[string]int{}
	var failures []string
	for _, entry := range audit.Entries {
		if _, ok := surfacesByPackage[entry.Package]; !ok || !strings.HasPrefix(entry.Disposition, "grouped:") {
			continue
		}

		receiver, ok := receiverForSymbol(entry.Symbol)
		if !ok {
			failures = append(failures, fmt.Sprintf("%s grouped entry has non receiver-qualified symbol %q", entry.Package, entry.Symbol))
			continue
		}

		row := strings.Join([]string{entry.Symbol, entry.Kind, entry.Signature}, "\t")
		rowsByPackage[entry.Package] = append(rowsByPackage[entry.Package], row)
		if receiverCountsByPackage[entry.Package] == nil {
			receiverCountsByPackage[entry.Package] = map[string]int{}
		}
		if kindCountsByPackage[entry.Package] == nil {
			kindCountsByPackage[entry.Package] = map[string]int{}
		}
		receiverCountsByPackage[entry.Package][receiver]++
		kindCountsByPackage[entry.Package][entry.Kind]++
	}

	totalGrouped := 0
	for _, surface := range approvalGroupedSurfaces {
		rows := rowsByPackage[surface.pkg]
		totalGrouped += len(rows)
		failures = append(failures, verifyGroupedSurfaceFingerprint(surface.pkg, rows, surface.entryCount, surface.signatureFingerprint)...)

		kindCounts := kindCountsByPackage[surface.pkg]
		if kindCounts["field"] != surface.fieldCount || kindCounts["method"] != surface.methodCount {
			failures = append(failures, fmt.Sprintf(
				"%s grouped kind counts mismatch: got fields/methods=%d/%d want=%d/%d",
				surface.pkg,
				kindCounts["field"],
				kindCounts["method"],
				surface.fieldCount,
				surface.methodCount,
			))
		}

		receiverRows := receiverRows(receiverCountsByPackage[surface.pkg])
		failures = append(failures, verifyGroupedSurfaceFingerprint(surface.pkg+" receiver families", receiverRows, surface.receiverCount, surface.receiverFingerprint)...)
	}
	if totalGrouped != 1083 {
		failures = append(failures, fmt.Sprintf("approval grouped surface total mismatch: got %d want 1083", totalGrouped))
	}

	return failures
}

func verifyGroupedSurfaceFingerprint(label string, rows []string, wantCount int, wantFingerprint string) []string {
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

func receiverRows(counts map[string]int) []string {
	rows := make([]string, 0, len(counts))
	for receiver, count := range counts {
		rows = append(rows, fmt.Sprintf("%d %s", count, receiver))
	}

	return rows
}

func receiverForSymbol(symbol string) (string, bool) {
	receiver, _, ok := strings.Cut(symbol, ".")
	if !ok || receiver == "" {
		return "", false
	}

	return receiver, true
}

func containsNormalized(content, term string) bool {
	return strings.Contains(content, term) || strings.Contains(strings.Join(strings.Fields(content), " "), strings.Join(strings.Fields(term), " "))
}

func verifyRuntimeResources(runtime runtimeLedger, docs []corpus) []string {
	var failures []string
	if runtime.SourceModule != "github.com/coldsmirk/vef-framework-go" {
		failures = append(failures, "runtime ledger source_module mismatch: "+runtime.SourceModule)
	}

	resources := map[string]bool{}
	actionsByResource := map[string]map[string]runtimeEntry{}
	for _, entry := range runtime.Entries {
		if entry.Category == "built-in resource" && strings.HasPrefix(entry.Name, "approval/") {
			resources[entry.Name] = true
		}
		if entry.Category == "built-in resource action" && strings.HasPrefix(entry.Name, "approval/") {
			resource, action, ok := strings.Cut(entry.Name, "/")
			if !ok {
				continue
			}
			// approval/resource/action splits into approval + resource/action; join first two segments.
			parts := strings.Split(entry.Name, "/")
			if len(parts) < 3 {
				continue
			}
			resource = parts[0] + "/" + parts[1]
			action = strings.Join(parts[2:], "/")
			if actionsByResource[resource] == nil {
				actionsByResource[resource] = map[string]runtimeEntry{}
			}
			actionsByResource[resource][action] = entry
		}
	}

	expectedResources := []string{
		"approval/category",
		"approval/delegation",
		"approval/flow",
		"approval/instance",
		"approval/my",
		"approval/admin",
	}
	failures = append(failures, compareSets("approval built-in resources", sortedKeys(resources), expectedResources)...)

	for _, doc := range docs {
		for _, resource := range expectedResources {
			if !strings.Contains(doc.content, "`"+resource+"`") {
				failures = append(failures, doc.label+" missing resource "+resource)
			}
			gotDocActions := docActionsForResource(doc.content, resource)
			gotRuntimeActions := sortedKeys(actionsByResource[resource])
			failures = append(failures, compareSets(doc.label+" "+resource+" action rows", gotDocActions, gotRuntimeActions)...)

			for action, entry := range actionsByResource[resource] {
				row := rowForFirstColumn(doc.content, action, resource)
				if row == "" {
					failures = append(failures, doc.label+" missing action row "+resource+"/"+action)
					continue
				}
				for _, detail := range entry.Details {
					if perm, ok := strings.CutPrefix(detail, "permission: "); ok && !strings.Contains(row, "`"+perm+"`") {
						failures = append(failures, doc.label+" action "+resource+"/"+action+" missing permission "+perm)
					}
					if detail == "audit enabled" && !containsAny(row, "Audited", "开启审计") {
						failures = append(failures, doc.label+" action "+resource+"/"+action+" missing audit note")
					}
				}
			}
		}
		if !strings.Contains(doc.content, "RequiredPermission") ||
			!containsAny(doc.content, "authenticated principal", "已认证") {
			failures = append(failures, doc.label+" missing approval/my no-RequiredPermission authenticated-principal contract")
		}
		if !containsAny(doc.content, "max `10`", "最多 `10`") || !strings.Contains(doc.content, "`1m`") {
			failures = append(failures, doc.label+" missing urge_task 10 per 1m rate limit")
		}
	}

	return failures
}

func verifyApprovalEnums(constsByType map[string][]stringConst, docs []corpus) []string {
	var failures []string
	enumTypes := []string{
		"BindingMode",
		"VersionStatus",
		"InitiatorKind",
		"StorageMode",
		"NodeKind",
		"ExecutionType",
		"ApprovalMethod",
		"PassRule",
		"EmptyAssigneeAction",
		"SameApplicantAction",
		"RollbackType",
		"RollbackDataStrategy",
		"AddAssigneeType",
		"ConsecutiveApproverAction",
		"AssigneeKind",
		"InstanceStatus",
		"TaskStatus",
		"ConditionKind",
		"ActionType",
		"CCKind",
		"CCTiming",
		"FieldKind",
		"TimeoutAction",
		"Permission",
	}
	for _, doc := range docs {
		for _, typeName := range enumTypes {
			values := constsByType[typeName]
			if len(values) == 0 {
				failures = append(failures, "source enum extraction found no constants for "+typeName)
				continue
			}
			if !strings.Contains(doc.content, "`"+typeName+"`") && typeName != "BindingMode" && typeName != "NodeKind" {
				failures = append(failures, doc.label+" missing enum type "+typeName)
			}
			for _, c := range values {
				if !strings.Contains(doc.content, "`"+c.name+"`") {
					failures = append(failures, doc.label+" missing enum constant "+c.name)
				}
				if !strings.Contains(doc.content, "`"+c.value+"`") {
					failures = append(failures, doc.label+" missing enum wire value "+typeName+"."+c.name+"="+c.value)
				}
			}
		}
	}

	return failures
}

func verifyStateMachines(sourceRoot string, constsByName map[string]stringConst, docs []corpus) []string {
	transitions := extractStateTransitions(sourceRoot, constsByName)
	var instanceRows []string
	var taskRows []string
	for _, t := range transitions {
		row := "`" + t.from + "` | `" + t.to + "`"
		if t.kind == "instance" {
			instanceRows = append(instanceRows, row)
		}
		if t.kind == "task" {
			taskRows = append(taskRows, row)
		}
	}
	sort.Strings(instanceRows)
	sort.Strings(taskRows)

	wantInstance := []string{
		"`returned` | `running`",
		"`running` | `approved`",
		"`running` | `rejected`",
		"`running` | `returned`",
		"`running` | `terminated`",
		"`running` | `withdrawn`",
		"`returned` | `terminated`",
		"`returned` | `withdrawn`",
		"`withdrawn` | `running`",
		"`withdrawn` | `terminated`",
	}
	wantTask := []string{
		"`pending` | `approved`",
		"`pending` | `canceled`",
		"`pending` | `handled`",
		"`pending` | `rejected`",
		"`pending` | `removed`",
		"`pending` | `rolled_back`",
		"`pending` | `transferred`",
		"`pending` | `waiting`",
		"`waiting` | `canceled`",
		"`waiting` | `pending`",
		"`waiting` | `removed`",
		"`waiting` | `skipped`",
	}
	failures := compareSets("source instance transitions", instanceRows, wantInstance)
	failures = append(failures, compareSets("source task transitions", taskRows, wantTask)...)
	for _, doc := range docs {
		for _, row := range append(wantInstance, wantTask...) {
			if !strings.Contains(doc.content, row) {
				failures = append(failures, doc.label+" missing state transition row "+row)
			}
		}
		if !strings.Contains(doc.content, "Withdrawn/Returned") && !strings.Contains(doc.content, "已撤回/已退回") {
			failures = append(failures, doc.label+" missing withdrawn/returned resubmit lifecycle text")
		}
	}

	return failures
}

func verifyRequestDTOs(sourceRoot string, docs []corpus) []string {
	typeSpec := map[string]struct {
		path   string
		prefix string
	}{
		"CategoryParams":               {"internal/approval/resource/category.go", "params"},
		"CategorySearch":               {"internal/approval/resource/category.go", "meta"},
		"DelegationParams":             {"internal/approval/resource/delegation.go", "params"},
		"DelegationSearch":             {"internal/approval/resource/delegation.go", "meta"},
		"CreateFlowParams":             {"internal/approval/resource/flow.go", "params"},
		"CreateInitiatorParams":        {"internal/approval/resource/flow.go", ""},
		"DeployFlowParams":             {"internal/approval/resource/flow.go", "params"},
		"PublishVersionParams":         {"internal/approval/resource/flow.go", "params"},
		"GetGraphParams":               {"internal/approval/resource/flow.go", "params"},
		"FindFlowsParams":              {"internal/approval/resource/flow.go", "params"},
		"UpdateParams":                 {"internal/approval/resource/flow.go", "params"},
		"ToggleActiveParams":           {"internal/approval/resource/flow.go", "params"},
		"FindVersionsParams":           {"internal/approval/resource/flow.go", "params"},
		"StartParams":                  {"internal/approval/resource/instance.go", "params"},
		"ProcessTaskParams":            {"internal/approval/resource/instance.go", "params"},
		"WithdrawParams":               {"internal/approval/resource/instance.go", "params"},
		"ResubmitParams":               {"internal/approval/resource/instance.go", "params"},
		"AddCCParams":                  {"internal/approval/resource/instance.go", "params"},
		"MarkCCReadParams":             {"internal/approval/resource/instance.go", "params"},
		"AddAssigneeParams":            {"internal/approval/resource/instance.go", "params"},
		"RemoveAssigneeParams":         {"internal/approval/resource/instance.go", "params"},
		"UrgeTaskParams":               {"internal/approval/resource/instance.go", "params"},
		"FindAvailableFlowsParams":     {"internal/approval/resource/my.go", "params"},
		"FindInitiatedParams":          {"internal/approval/resource/my.go", "params"},
		"FindPendingTasksParams":       {"internal/approval/resource/my.go", "params"},
		"FindCompletedTasksParams":     {"internal/approval/resource/my.go", "params"},
		"FindCCRecordsParams":          {"internal/approval/resource/my.go", "params"},
		"GetPendingCountsParams":       {"internal/approval/resource/my.go", "params"},
		"GetInstanceDetailParams":      {"internal/approval/resource/my.go", "params"},
		"AdminFindInstancesParams":     {"internal/approval/resource/admin.go", "params"},
		"AdminFindTasksParams":         {"internal/approval/resource/admin.go", "params"},
		"AdminGetInstanceDetailParams": {"internal/approval/resource/admin.go", "params"},
		"AdminFindActionLogsParams":    {"internal/approval/resource/admin.go", "params"},
		"AdminTerminateInstanceParams": {"internal/approval/resource/admin.go", "params"},
		"AdminReassignTaskParams":      {"internal/approval/resource/admin.go", "params"},
		"AdminGetMetricsParams":        {"internal/approval/resource/admin.go", "params"},
	}

	fieldsByPath := map[string]map[string][]string{}
	for typeName, spec := range typeSpec {
		if fieldsByPath[spec.path] == nil {
			fieldsByPath[spec.path] = extractStructJSONFields(sourceRoot, spec.path)
		}
		fields := fieldsByPath[spec.path][typeName]
		if len(fields) == 0 {
			return []string{"source request DTO extraction found no json fields for " + typeName}
		}
	}

	var failures []string
	for _, doc := range docs {
		for typeName, spec := range typeSpec {
			fields := fieldsByPath[spec.path][typeName]
			if !strings.Contains(doc.content, "`"+typeName+"`") && !strings.HasPrefix(typeName, "CreateInitiator") {
				failures = append(failures, doc.label+" missing request DTO type "+typeName)
			}
			for _, field := range fields {
				term := "`" + field + "`"
				if spec.prefix != "" {
					term = "`" + spec.prefix + "." + field + "`"
				}
				if !strings.Contains(doc.content, term) {
					failures = append(failures, doc.label+" missing request field "+typeName+"."+field+" as "+term)
				}
			}
		}
		if !strings.Contains(doc.content, "`approve`") || !strings.Contains(doc.content, "`handle`") ||
			!strings.Contains(doc.content, "`transfer`") || !strings.Contains(doc.content, "`rollback`") {
			failures = append(failures, doc.label+" missing process_task action oneof values")
		}
		instanceSection := markdownSection(doc.content, "### `approval/instance`")
		startRow := rowForMarker(instanceSection, "| `start` | `params.tenantId`")
		if !containsAny(startRow, "max 512 chars", "最多 512 字符") {
			failures = append(failures, doc.label+" approval/instance start row missing businessRef max 512 constraint")
		}
		processTaskRow := rowForMarker(instanceSection, "| `process_task` | `params.taskId`")
		if !containsAny(processTaskRow, "max 20 entries", "最多 20 项") ||
			!containsAny(processTaskRow, "each max 512 chars", "每项最多 512 字符") {
			failures = append(failures, doc.label+" approval/instance process_task row missing attachments max 20/max 512 constraint")
		}
	}

	return failures
}

func verifyResponseDTOs(sourceRoot string, docs []corpus) []string {
	adminFields := map[string][]string{}
	for _, path := range []string{
		"approval/admin/instance.go",
		"approval/admin/instance_detail.go",
		"approval/admin/task.go",
		"approval/admin/metrics.go",
	} {
		for typeName, fields := range extractStructJSONFields(sourceRoot, path) {
			adminFields[typeName] = fields
		}
	}
	myFields := map[string][]string{}
	for _, path := range []string{
		"approval/my/pending_tasks.go",
		"approval/my/completed_tasks.go",
		"approval/my/cc_records.go",
		"approval/my/initiated_instances.go",
		"approval/my/available_flows.go",
		"approval/my/instance_detail.go",
		"approval/my/pending_counts.go",
	} {
		for typeName, fields := range extractStructJSONFields(sourceRoot, path) {
			myFields[typeName] = fields
		}
	}

	var failures []string
	for _, doc := range docs {
		dtoSection := responseDTOSection(doc.content)
		for typeName, fields := range adminFields {
			row := rowForMarker(dtoSection, "`admin."+typeName+"`")
			if row == "" {
				failures = append(failures, doc.label+" missing admin DTO row "+typeName)
				continue
			}
			for _, field := range fields {
				if !strings.Contains(row, "`"+field+"`") {
					failures = append(failures, doc.label+" admin DTO "+typeName+" row missing field "+field)
				}
			}
		}
		for typeName, fields := range myFields {
			row := rowForMarker(dtoSection, "`my."+typeName+"`")
			if row == "" {
				failures = append(failures, doc.label+" missing my DTO row "+typeName)
				continue
			}
			for _, field := range fields {
				if !strings.Contains(row, "`"+field+"`") {
					failures = append(failures, doc.label+" my DTO "+typeName+" row missing field "+field)
				}
			}
		}
	}

	return failures
}

func verifyFlowAndFormWireShapes(sourceRoot string, docs []corpus) []string {
	typeSources := map[string][]string{
		"approval/flow_definition.go": {"FlowDefinition", "NodeDefinition", "Position", "EdgeDefinition"},
		"approval/node_data.go":       {"BaseNodeData", "TaskNodeData", "ApprovalNodeData", "CCNodeData", "ConditionNodeData"},
		"approval/form_field.go":      {"FormFieldDefinition", "FieldOption", "ValidationRule"},
		"approval/assignee.go":        {"AssigneeDefinition", "CCDefinition"},
		"approval/condition.go":       {"Condition", "ConditionGroup", "ConditionBranch"},
	}
	fieldsByType := map[string][]string{}
	for path, typeNames := range typeSources {
		extracted := extractStructJSONFields(sourceRoot, path)
		for _, typeName := range typeNames {
			fieldsByType[typeName] = extracted[typeName]
		}
	}

	var failures []string
	for _, doc := range docs {
		for typeName, fields := range fieldsByType {
			if len(fields) == 0 {
				failures = append(failures, "source flow/form extraction found no json fields for "+typeName)
				continue
			}
			if !strings.Contains(doc.content, "`"+typeName+"`") {
				failures = append(failures, doc.label+" missing flow/form type "+typeName)
			}
			for _, field := range fields {
				if !strings.Contains(doc.content, "`"+field+"`") {
					failures = append(failures, doc.label+" missing flow/form JSON field "+typeName+"."+field)
				}
			}
		}
		for _, term := range []string{
			"`sourceHandle`",
			"`rollbackTargetKeys`",
			"`editable`",
			"`required`",
			"`visible`",
			"`hidden`",
			"`minLength`",
			"`maxLength`",
			"`pattern`",
			"64 KiB",
		} {
			if !strings.Contains(doc.content, term) {
				failures = append(failures, doc.label+" missing flow/form behavior term "+term)
			}
		}
	}

	return failures
}

func verifyEvents(sourceRoot string, constsByName map[string]stringConst, docs []corpus) []string {
	events := extractEvents(sourceRoot, constsByName)
	var failures []string
	eventTypes := eventTypeConstants(constsByName)
	if len(events) != len(eventTypes) {
		failures = append(failures, fmt.Sprintf("source event extraction got %d events, want %d public EventType constants", len(events), len(eventTypes)))
	}
	eventByConst := map[string]eventContract{}
	for _, event := range events {
		eventByConst[event.constName] = event
	}
	for _, eventType := range eventTypes {
		if _, ok := eventByConst[eventType.name]; !ok {
			failures = append(failures, "source event extraction missing struct/EventType method for "+eventType.name)
		}
	}
	for _, doc := range docs {
		for _, event := range events {
			row := rowForMarker(doc.content, "`"+event.constName+"`")
			if row == "" {
				failures = append(failures, doc.label+" missing event row "+event.constName)
				continue
			}
			for _, term := range []string{event.topic, event.structName, event.constructor} {
				if !strings.Contains(row, "`"+term+"`") {
					failures = append(failures, doc.label+" event "+event.constName+" row missing "+term)
				}
			}
			for _, field := range event.fields {
				if field == "tenantId" || field == "occurredTime" {
					continue
				}
				if !strings.Contains(row, "`"+field+"`") {
					failures = append(failures, doc.label+" event "+event.constName+" row missing payload field "+field)
				}
			}
		}
		if !strings.Contains(doc.content, "`tenantId`") || !strings.Contains(doc.content, "`occurredTime`") || !strings.Contains(doc.content, "event.WithTx") {
			failures = append(failures, doc.label+" missing common approval event fields or event.WithTx")
		}
	}

	return failures
}

func eventTypeConstants(constsByName map[string]stringConst) []stringConst {
	var eventTypes []stringConst
	for _, c := range constsByName {
		if strings.HasPrefix(c.name, "EventType") {
			eventTypes = append(eventTypes, c)
		}
	}
	sort.Slice(eventTypes, func(i, j int) bool { return eventTypes[i].name < eventTypes[j].name })

	return eventTypes
}

func verifyErrors(sourceRoot string, docs []corpus) []string {
	codeValues := extractIntConstants(sourceRoot, "internal/approval/shared/errors.go")
	errorContracts := extractErrorContracts(sourceRoot, codeValues)
	messageKeys := extractMessageKeys(sourceRoot)

	var failures []string
	for _, doc := range docs {
		for _, errorContract := range errorContracts {
			row := rowForMarker(doc.content, "`"+errorContract.errName+"`")
			if row == "" {
				failures = append(failures, doc.label+" missing approval error row "+errorContract.errName)
				continue
			}
			for _, term := range []string{errorContract.codeValue, errorContract.codeName, errorContract.errName, errorContract.message} {
				if !strings.Contains(row, "`"+term+"`") {
					failures = append(failures, doc.label+" error row "+errorContract.errName+" missing "+term)
				}
			}
		}
		for _, codeName := range sortedKeys(codeValues) {
			code := strconv.Itoa(codeValues[codeName])
			if !strings.Contains(doc.content, "`"+codeName+"`") || !strings.Contains(doc.content, "`"+code+"`") {
				failures = append(failures, doc.label+" missing approval error code "+codeName+"="+code)
			}
		}
		for _, message := range messageKeys {
			if !strings.Contains(doc.content, message) {
				failures = append(failures, doc.label+" missing approval i18n message key "+message)
			}
		}
		for _, term := range []string{
			"approval.ErrCrossTenantAccess",
			"approval.ErrInvalidBusinessIdentifier",
			"approval.ErrUnknownNodeKind",
			"approval.ErrNodeDataUnmarshal",
		} {
			if !strings.Contains(doc.content, term) {
				failures = append(failures, doc.label+" missing public approval sentinel term "+term)
			}
		}
		if !containsAny(doc.content, "not `result.Error`", "不是 `result.Error`") {
			failures = append(failures, doc.label+" missing public approval sentinel term not result.Error")
		}
	}

	return failures
}

func verifyResubmitAndAvailableActions(sourceRoot string, docs []corpus) []string {
	resubmitSource := readFile(filepath.Join(sourceRoot, "internal/approval/command/resubmit_instance.go"))
	availableActionsSource := readFile(filepath.Join(sourceRoot, "internal/approval/query/get_my_instance_detail.go"))
	var failures []string
	for _, term := range []string{
		"engine.InstanceStateMachine.CanTransition(instance.Status, approval.InstanceRunning)",
		"shared.ErrResubmitNotAllowed",
		"approval.ActionResubmit",
		"NewInstanceResubmittedEvent",
	} {
		if !strings.Contains(resubmitSource, term) {
			failures = append(failures, "resubmit handler missing expected source term "+term)
		}
	}
	for _, term := range []string{
		"engine.InstanceStateMachine.CanTransition(instance.Status, approval.InstanceWithdrawn)",
		"actions.Add(\"withdraw\")",
		"engine.InstanceStateMachine.CanTransition(instance.Status, approval.InstanceRunning)",
		"actions.Add(\"resubmit\")",
		"actions.Add(\"handle\")",
		"actions.Add(\"approve\")",
		"actions.Add(\"reject\")",
		"actions.Add(\"transfer\")",
		"actions.Add(\"rollback\")",
		"actions.Add(\"add_assignee\")",
		"actions.Add(\"add_cc\")",
		"actions.Add(\"urge\")",
	} {
		if !strings.Contains(availableActionsSource, term) {
			failures = append(failures, "availableActions source missing expected term "+term)
		}
	}
	for _, doc := range docs {
		if !containsAny(doc.content, "returned or withdrawn", "已退回或已撤回") {
			failures = append(failures, doc.label+" missing command resubmit returned-or-withdrawn contract")
		}
		for _, term := range []string{"`availableActions`", "`withdraw`", "`resubmit`", "`rejected`", "`returned`", "`handle`", "`approve`", "`reject`", "`transfer`", "`rollback`", "`add_assignee`", "`add_cc`", "`urge`"} {
			if !strings.Contains(doc.content, term) {
				failures = append(failures, doc.label+" missing availableActions term "+term)
			}
		}
		if !containsAny(doc.content, "Command handlers still perform", "命令 handler 仍会独立") {
			failures = append(failures, doc.label+" missing availableActions query-hint vs command-validation distinction")
		}
	}

	return failures
}

func verifyTenantBindingAndExtensionTerms(docs []corpus) []string {
	var failures []string
	terms := []string{
		"approval:super_admin",
		"CallerContext",
		"SystemCaller",
		"PrincipalTenantResolver",
		"PrincipalDepartmentResolver",
		"BusinessRefProvider",
		"BusinessRefResolver",
		"InstanceLifecycleHook",
		"vef:approval:lifecycle_hooks",
		"BusinessTable",
		"BusinessPKField",
		"BusinessStatusField",
		"ValidateBusinessIdentifier",
		"^[A-Za-z_][A-Za-z0-9_]{0,62}$",
		"InstanceBindingFailedEvent",
		"InstanceCompletedEvent",
	}
	staleTerms := []string{"BusinessTitleField", "businessTitleField"}
	for _, doc := range docs {
		for _, term := range terms {
			if !strings.Contains(doc.content, term) {
				failures = append(failures, doc.label+" missing tenant/binding/extension term "+term)
			}
		}
		for _, term := range staleTerms {
			if strings.Contains(doc.content, term) {
				failures = append(failures, doc.label+" contains stale tenant/binding/extension term "+term)
			}
		}
	}

	return failures
}

func extractStringConstants(sourceRoot string) (map[string]stringConst, map[string][]stringConst) {
	result := map[string]stringConst{}
	byType := map[string][]stringConst{}
	for _, relPath := range []string{"approval/enums.go", "approval/events.go"} {
		file := parseGoFile(filepath.Join(sourceRoot, relPath))
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.CONST {
				continue
			}
			currentType := ""
			for _, spec := range gen.Specs {
				valueSpec := spec.(*ast.ValueSpec)
				if ident, ok := valueSpec.Type.(*ast.Ident); ok {
					currentType = ident.Name
				}
				for i, name := range valueSpec.Names {
					if i >= len(valueSpec.Values) {
						continue
					}
					lit, ok := valueSpec.Values[i].(*ast.BasicLit)
					if !ok || lit.Kind != token.STRING {
						continue
					}
					value, err := strconv.Unquote(lit.Value)
					if err != nil {
						panic(err)
					}
					c := stringConst{name: name.Name, typeName: currentType, value: value}
					result[name.Name] = c
					byType[currentType] = append(byType[currentType], c)
				}
			}
		}
	}
	for typeName := range byType {
		sort.Slice(byType[typeName], func(i, j int) bool { return byType[typeName][i].name < byType[typeName][j].name })
	}

	return result, byType
}

func responseDTOSection(content string) string {
	section := markdownSection(content, "### Response DTO Fields")
	if section != "" {
		return section
	}

	return markdownSection(content, "### 响应 DTO 字段")
}

func extractStateTransitions(sourceRoot string, constsByName map[string]stringConst) []transition {
	file := parseGoFile(filepath.Join(sourceRoot, "internal/approval/engine/state_machine.go"))
	var transitions []transition
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "AddTransition" || len(call.Args) != 2 {
			return true
		}
		fromName, okFrom := selectorName(call.Args[0])
		toName, okTo := selectorName(call.Args[1])
		if !okFrom || !okTo {
			return true
		}
		from := constsByName[fromName]
		to := constsByName[toName]
		kind := ""
		if from.typeName == "InstanceStatus" && to.typeName == "InstanceStatus" {
			kind = "instance"
		}
		if from.typeName == "TaskStatus" && to.typeName == "TaskStatus" {
			kind = "task"
		}
		if kind != "" {
			transitions = append(transitions, transition{kind: kind, from: from.value, to: to.value})
		}
		return true
	})
	sort.Slice(transitions, func(i, j int) bool {
		if transitions[i].kind != transitions[j].kind {
			return transitions[i].kind < transitions[j].kind
		}
		if transitions[i].from != transitions[j].from {
			return transitions[i].from < transitions[j].from
		}
		return transitions[i].to < transitions[j].to
	})

	return transitions
}

func extractStructJSONFields(sourceRoot, relPath string) map[string][]string {
	file := parseGoFile(filepath.Join(sourceRoot, relPath))
	result := map[string][]string{}
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			typeSpec := spec.(*ast.TypeSpec)
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			var fields []string
			for _, field := range structType.Fields.List {
				if field.Tag == nil {
					continue
				}
				name := jsonName(field.Tag.Value)
				if name == "" {
					continue
				}
				fields = append(fields, name)
			}
			sort.Strings(fields)
			result[typeSpec.Name.Name] = fields
		}
	}

	return result
}

func extractEvents(sourceRoot string, constsByName map[string]stringConst) []eventContract {
	eventFiles := []string{
		"approval/events_instance.go",
		"approval/events_task.go",
		"approval/events_timeout.go",
		"approval/events_node.go",
		"approval/events_cc.go",
		"approval/events_flow.go",
	}
	fieldsByStruct := map[string][]string{}
	constByStruct := map[string]string{}
	for _, relPath := range eventFiles {
		file := parseGoFile(filepath.Join(sourceRoot, relPath))
		for typeName, fields := range extractStructJSONFields(sourceRoot, relPath) {
			if strings.HasSuffix(typeName, "Event") {
				fieldsByStruct[typeName] = fields
			}
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Name.Name != "EventType" || fn.Recv == nil || len(fn.Recv.List) != 1 || fn.Body == nil {
				continue
			}
			receiver := receiverTypeName(fn.Recv.List[0].Type)
			if receiver == "" {
				continue
			}
			for _, stmt := range fn.Body.List {
				ret, ok := stmt.(*ast.ReturnStmt)
				if !ok || len(ret.Results) != 1 {
					continue
				}
				ident, ok := ret.Results[0].(*ast.Ident)
				if ok {
					constByStruct[receiver] = ident.Name
				}
			}
		}
	}

	var events []eventContract
	for structName, fields := range fieldsByStruct {
		constName := constByStruct[structName]
		c := constsByName[constName]
		if constName == "" || c.value == "" {
			continue
		}
		events = append(events, eventContract{
			constName:   constName,
			topic:       c.value,
			structName:  structName,
			constructor: "New" + structName,
			fields:      fields,
		})
	}
	sort.Slice(events, func(i, j int) bool { return events[i].constName < events[j].constName })

	return events
}

func extractIntConstants(sourceRoot, relPath string) map[string]int {
	file := parseGoFile(filepath.Join(sourceRoot, relPath))
	result := map[string]int{}
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			valueSpec := spec.(*ast.ValueSpec)
			for i, name := range valueSpec.Names {
				if i >= len(valueSpec.Values) {
					continue
				}
				lit, ok := valueSpec.Values[i].(*ast.BasicLit)
				if !ok || lit.Kind != token.INT {
					continue
				}
				value, err := strconv.Atoi(lit.Value)
				if err != nil {
					panic(err)
				}
				result[name.Name] = value
			}
		}
	}

	return result
}

func extractErrorContracts(sourceRoot string, codeValues map[string]int) []errorContract {
	file := parseGoFile(filepath.Join(sourceRoot, "internal/approval/shared/api_errors.go"))
	var result []errorContract
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.VAR {
			continue
		}
		for _, spec := range gen.Specs {
			valueSpec := spec.(*ast.ValueSpec)
			for i, name := range valueSpec.Names {
				if i >= len(valueSpec.Values) {
					continue
				}
				message, codeName := findErrorMessageAndCode(valueSpec.Values[i])
				if message == "" || codeName == "" {
					continue
				}
				result = append(result, errorContract{
					errName:   name.Name,
					codeName:  codeName,
					codeValue: strconv.Itoa(codeValues[codeName]),
					message:   message,
				})
			}
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].errName < result[j].errName })

	return result
}

func findErrorMessageAndCode(expr ast.Expr) (string, string) {
	message := ""
	codeName := ""
	ast.Inspect(expr, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "T" && len(call.Args) == 1 {
			if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
				message, _ = strconv.Unquote(lit.Value)
			}
		}
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "WithCode" && len(call.Args) == 1 {
			if ident, ok := call.Args[0].(*ast.Ident); ok {
				codeName = ident.Name
			}
		}
		return true
	})

	return message, codeName
}

func extractMessageKeys(sourceRoot string) []string {
	file := parseGoFile(filepath.Join(sourceRoot, "internal/approval/shared/messages.go"))
	var result []string
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			valueSpec := spec.(*ast.ValueSpec)
			for _, value := range valueSpec.Values {
				lit, ok := value.(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					continue
				}
				s, err := strconv.Unquote(lit.Value)
				if err != nil {
					panic(err)
				}
				result = append(result, s)
			}
		}
	}
	sort.Strings(result)

	return result
}

func runSourceTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./approval/...", "./internal/approval/...", "-count=1")
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("go test ./approval/... ./internal/approval/... failed: %v\n%s", err, strings.TrimSpace(string(output)))}
	}

	return nil
}

func loadLiveInventory(sourceRoot, docsRoot string) map[string]liveInventoryEntry {
	cmd := exec.Command(
		"go", "run", filepath.Join(docsRoot, "scripts/verify-api-audit.go"),
		"-source", sourceRoot,
		"-print-current",
	)
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Errorf("verify-api-audit -print-current failed: %w\n%s", err, strings.TrimSpace(string(output))))
	}
	payload := "[" + strings.TrimSpace(string(output)) + "]"
	var entries []liveInventoryEntry
	if err := json.Unmarshal([]byte(payload), &entries); err != nil {
		panic(fmt.Errorf("parse live inventory: %w", err))
	}
	result := map[string]liveInventoryEntry{}
	for _, entry := range entries {
		result[entry.Package] = entry
	}

	return result
}

func approvalRuntimeSurfaceCount(runtime runtimeLedger) int {
	count := 0
	for _, entry := range runtime.Entries {
		if (entry.Category == "built-in resource" || entry.Category == "built-in resource action") && strings.HasPrefix(entry.Name, "approval/") {
			count++
		}
	}

	return count
}

func docActionsForResource(content, resource string) []string {
	section := markdownSection(content, "### `"+resource+"`")
	matches := regexp.MustCompile("(?m)^\\| `([^`]+)` \\|").FindAllStringSubmatch(section, -1)
	set := map[string]bool{}
	for _, match := range matches {
		action := match[1]
		if action == "Action" {
			continue
		}
		set[action] = true
	}

	return sortedKeys(set)
}

func rowForFirstColumn(content, firstColumn, sectionMarker string) string {
	section := markdownSection(content, "### `"+sectionMarker+"`")
	return rowForMarker(section, "| `"+firstColumn+"` |")
}

func rowForMarker(content, marker string) string {
	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(line, marker) {
			return line
		}
	}

	return ""
}

func markdownSection(content, marker string) string {
	start := strings.Index(content, marker)
	if start < 0 {
		return ""
	}
	rest := content[start:]
	nextHeading := regexp.MustCompile(`(?m)^#{2,3} `).FindAllStringIndex(rest[len(marker):], -1)
	if len(nextHeading) == 0 {
		return rest
	}

	return rest[:len(marker)+nextHeading[0][0]]
}

func parseGoFile(path string) *ast.File {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	return file
}

func selectorName(expr ast.Expr) (string, bool) {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}

	return sel.Sel.Name, true
}

func receiverTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name
		}
	}

	return ""
}

func jsonName(tagLiteral string) string {
	tag, err := strconv.Unquote(tagLiteral)
	if err != nil {
		panic(err)
	}
	value := reflect.StructTag(tag).Get("json")
	if value == "" || value == "-" {
		return ""
	}
	name, _, _ := strings.Cut(value, ",")
	if name == "" || name == "-" {
		return ""
	}

	return name
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

func containsAny(content string, terms ...string) bool {
	for _, term := range terms {
		if strings.Contains(content, term) {
			return true
		}
	}

	return false
}

func readCorpus(label, path string) corpus {
	return corpus{label: label, content: readFile(path)}
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(data)
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

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}

func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	err := cmd.Run()

	return output.String(), err
}
