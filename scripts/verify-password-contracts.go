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
	"strconv"
	"strings"
)

const (
	passwordPackage = "github.com/coldsmirk/vef-framework-go/password"

	passwordFingerprint = "7759597dd9992c50971c71c33d276b122e4d38cf2f54be381ae6e8db38c03087"
	passwordTopLevel    = 46
	passwordFields      = 0
	passwordMethods     = 3
	passwordEntries     = 49

	passwordGroupedEntries              = 3
	passwordGroupedFields               = 0
	passwordGroupedMethods              = 3
	passwordGroupedReceivers            = 1
	passwordGroupedSignatureFingerprint = "c3f9451be207b1da459663cea4fbf57365324e6ffb224b7453da7e2d2a697d20"
	passwordGroupedReceiverFingerprint  = "bbd49629bbd1ac2431849426c0fdabee3fa50db8fef9da37373c027d692dff98"

	englishPasswordPath = "docs/features/password.md"
	chinesePasswordPath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/features/password.md"
	englishIndexPath    = "docs/reference/public-api-index.md"
	chineseIndexPath    = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"
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
	Disposition    string   `json:"disposition"`
	Coverage       []string `json:"coverage"`
	SourceEvidence []string `json:"source_evidence"`
	Terms          []string `json:"terms"`
}

type liveInventoryEntry struct {
	Package     string   `json:"package"`
	Coverage    []string `json:"coverage"`
	TopLevel    int      `json:"top_level"`
	Fields      int      `json:"fields"`
	Methods     int      `json:"methods"`
	Fingerprint string   `json:"fingerprint"`
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

	englishDocs := readCorpus("English password docs", filepath.Join(docsRoot, englishPasswordPath))
	chineseDocs := readCorpus("Chinese password docs", filepath.Join(docsRoot, chinesePasswordPath))
	englishIndex := readCorpus("English public API index", filepath.Join(docsRoot, englishIndexPath))
	chineseIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, chineseIndexPath))

	audit := loadJSON[auditLedger](auditLedgerPath)
	manifestData := loadJSON[manifest](manifestPath)
	contracts := loadJSON[contractLedger](contractLedgerPath)
	liveManifestEntry := loadLiveInventory(sourceRoot, docsRoot, manifestPath, auditLedgerPath, contractLedgerPath)[passwordPackage]
	liveAuditEntries := passwordEntriesFromAudit(loadLiveAuditLedger(sourceRoot, docsRoot, manifestPath))
	passwordEntries := passwordEntriesFromAudit(audit)

	var failures []string
	failures = append(failures, verifySurfaceEntry("live public API inventory", liveManifestEntry)...)
	failures = append(failures, verifyManifest(manifestData)...)
	failures = append(failures, verifyContractLedger(contracts, sourceRoot)...)
	failures = append(failures, verifyAuditEntries(passwordEntries)...)
	failures = append(failures, verifyLiveAuditEntries(passwordEntries, liveAuditEntries)...)
	failures = append(failures, verifyGroupedPasswordSurface(passwordEntries, []corpus{englishDocs, chineseDocs})...)
	failures = append(failures, verifyGeneratedIndexSection(englishIndex, passwordEntries)...)
	failures = append(failures, verifyGeneratedIndexSection(chineseIndex, passwordEntries)...)
	failures = append(failures, verifyPasswordDocs(passwordEntries, []corpus{englishDocs, chineseDocs})...)
	failures = append(failures, verifySourceTerms(sourceRoot)...)
	failures = append(failures, runGoTests(sourceRoot)...)

	if len(failures) > 0 {
		sort.Strings(failures)
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Println("Password contract docs verified: 49 public entries, 3 grouped encoder methods, prefix/upgrade/cipher contracts")
}

func verifySurfaceEntry(label string, entry manifestEntry) []string {
	var failures []string
	if entry.Package != passwordPackage {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q", label, entry.Package))
	}
	if entry.TopLevel != passwordTopLevel ||
		entry.Fields != passwordFields ||
		entry.Methods != passwordMethods ||
		entry.Fingerprint != passwordFingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s password surface mismatch: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			passwordTopLevel, passwordFields, passwordMethods, passwordFingerprint,
		))
	}

	return failures
}

func verifyManifest(m manifest) []string {
	for _, entry := range m.Packages {
		if entry.Package != passwordPackage {
			continue
		}
		var failures []string
		failures = append(failures, verifySurfaceEntry("API audit manifest", entry)...)
		if !sameSet(entry.Coverage, passwordCoverage()) {
			failures = append(failures, fmt.Sprintf("password manifest coverage mismatch: got %v want %v", entry.Coverage, passwordCoverage()))
		}

		return failures
	}

	return []string{"API audit manifest missing password package"}
}

func verifyContractLedger(contracts contractLedger, sourceRoot string) []string {
	contractID := passwordPackage + "#runtime-contract:password-encoding-prefix-and-upgrade"
	expectedTerms := []string{
		"{bcrypt}",
		"EncoderID",
		"NewCompositeEncoder",
		"UpgradeEncoding",
		"bcrypt",
	}

	var failures []string
	var foundReview bool
	for _, review := range contracts.PackageReviews {
		if review.Package != passwordPackage {
			continue
		}
		foundReview = true
		if review.Disposition != "has-semantic-contracts" {
			failures = append(failures, "password contract review disposition mismatch: "+review.Disposition)
		}
		if review.ReviewedSurface.TopLevel != passwordTopLevel ||
			review.ReviewedSurface.Fields != passwordFields ||
			review.ReviewedSurface.Methods != passwordMethods ||
			review.ReviewedSurface.EntryCount != passwordEntries ||
			review.ReviewedSurface.Fingerprint != passwordFingerprint {
			failures = append(failures, fmt.Sprintf(
				"password contract review surface mismatch: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
				review.ReviewedSurface.TopLevel,
				review.ReviewedSurface.Fields,
				review.ReviewedSurface.Methods,
				review.ReviewedSurface.EntryCount,
				review.ReviewedSurface.Fingerprint,
			))
		}
		if !sameSet(review.Coverage, passwordCoverage()) {
			failures = append(failures, fmt.Sprintf("password contract review coverage mismatch: got %v want %v", review.Coverage, passwordCoverage()))
		}
		if !sameSet(review.ContractIDs, []string{contractID}) {
			failures = append(failures, fmt.Sprintf("password contract ids mismatch: got %v want %v", review.ContractIDs, []string{contractID}))
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, review.SourceEvidence)...)
	}
	if !foundReview {
		failures = append(failures, "contract ledger missing password package review")
	}

	var foundEntry bool
	for _, entry := range contracts.Entries {
		if entry.ID != contractID {
			continue
		}
		foundEntry = true
		if entry.Package != passwordPackage {
			failures = append(failures, "password contract entry package mismatch: "+entry.Package)
		}
		if entry.Disposition != "documented:semantic-contract" {
			failures = append(failures, "password contract entry disposition mismatch: "+entry.Disposition)
		}
		if !sameSet(entry.Coverage, passwordCoverage()) {
			failures = append(failures, fmt.Sprintf("password contract coverage mismatch: got %v want %v", entry.Coverage, passwordCoverage()))
		}
		for _, term := range expectedTerms {
			if !contains(entry.Terms, term) {
				failures = append(failures, "password contract entry missing term "+term)
			}
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, entry.SourceEvidence)...)
	}
	if !foundEntry {
		failures = append(failures, "contract ledger missing password contract entry "+contractID)
	}

	return failures
}

func verifyAuditEntries(entries []auditEntry) []string {
	var failures []string
	if len(entries) != passwordEntries {
		failures = append(failures, fmt.Sprintf("password audit entry count mismatch: got %d want %d", len(entries), passwordEntries))
	}
	counts := map[string]int{}
	dispositionCounts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != passwordPackage {
			failures = append(failures, "non-password audit entry passed into password verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate password audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		dispositionCounts[entry.Disposition]++
		if entry.Symbol == "" || entry.Signature == "" || entry.Disposition == "" {
			failures = append(failures, "password audit entry missing required metadata "+entry.ID)
		}
		if !sameSet(entry.Coverage, passwordCoverage()) {
			failures = append(failures, fmt.Sprintf("password audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage, passwordCoverage()))
		}
	}
	if counts["top"] != passwordTopLevel || counts["field"] != passwordFields || counts["method"] != passwordMethods {
		failures = append(failures, fmt.Sprintf("password audit kind counts mismatch: top/field/method=%d/%d/%d", counts["top"], counts["field"], counts["method"]))
	}
	if dispositionCounts["documented:top-level"] != passwordTopLevel ||
		dispositionCounts["grouped:type-member-family"] != passwordGroupedEntries {
		failures = append(failures, fmt.Sprintf(
			"password audit disposition counts mismatch: top-level/grouped=%d/%d want=%d/%d",
			dispositionCounts["documented:top-level"],
			dispositionCounts["grouped:type-member-family"],
			passwordTopLevel,
			passwordGroupedEntries,
		))
	}

	return failures
}

func verifyLiveAuditEntries(ledgerEntries, liveEntries []auditEntry) []string {
	ledgerByID := entriesByID(ledgerEntries)
	liveByID := entriesByID(liveEntries)
	var failures []string
	for id, live := range liveByID {
		ledger, ok := ledgerByID[id]
		if !ok {
			failures = append(failures, fmt.Sprintf("password missing_in_ledger: %s %s %s", id, live.Symbol, live.Signature))
			continue
		}
		if ledger.Kind != live.Kind || ledger.Symbol != live.Symbol || ledger.Signature != live.Signature {
			failures = append(failures, fmt.Sprintf(
				"password live/ledger signature drift for %s: ledger=%s/%s/%s live=%s/%s/%s",
				id,
				ledger.Kind,
				ledger.Symbol,
				ledger.Signature,
				live.Kind,
				live.Symbol,
				live.Signature,
			))
		}
	}
	for id, ledger := range ledgerByID {
		if _, ok := liveByID[id]; !ok {
			failures = append(failures, fmt.Sprintf("password extra_in_ledger: %s %s %s", id, ledger.Symbol, ledger.Signature))
		}
	}

	return failures
}

func verifyGroupedPasswordSurface(entries []auditEntry, docs []corpus) []string {
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
			failures = append(failures, fmt.Sprintf("password grouped entry has non receiver-qualified symbol %q", entry.Symbol))
			continue
		}
		rows = append(rows, strings.Join([]string{entry.Symbol, entry.Kind, entry.Disposition, entry.Signature}, "\t"))
		receiverCounts[receiver]++
		kindCounts[entry.Kind]++
	}

	failures = append(failures, verifyGroupedFingerprint("password grouped type-member surface", rows, passwordGroupedEntries, passwordGroupedSignatureFingerprint)...)
	if kindCounts["field"] != passwordGroupedFields || kindCounts["method"] != passwordGroupedMethods {
		failures = append(failures, fmt.Sprintf(
			"password grouped kind counts mismatch: got fields/methods=%d/%d want=%d/%d",
			kindCounts["field"],
			kindCounts["method"],
			passwordGroupedFields,
			passwordGroupedMethods,
		))
	}

	var receiverRows []string
	for receiver, count := range receiverCounts {
		receiverRows = append(receiverRows, fmt.Sprintf("%d %s", count, receiver))
	}
	failures = append(failures, verifyGroupedFingerprint("password grouped receiver/type families", receiverRows, passwordGroupedReceivers, passwordGroupedReceiverFingerprint)...)

	for _, doc := range docs {
		for _, term := range []string{
			"49 public password entries",
			"3 grouped password method entries",
			"1 password receiver/type family",
			"0 exported password field entries",
			"3 exported password method entries",
		} {
			if !containsNormalized(doc.content, term) {
				failures = append(failures, doc.label+" missing grouped password audit term "+term)
			}
		}
	}

	return failures
}

func verifyGeneratedIndexSection(index corpus, entries []auditEntry) []string {
	section := packageSection(index.content, passwordPackage)
	if section == "" {
		return []string{index.label + " missing password package section"}
	}

	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s password index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
		}
	}

	return failures
}

func verifyPasswordDocs(entries []auditEntry, docs []corpus) []string {
	var topSymbols []string
	for _, entry := range entries {
		if entry.Kind == "top" {
			topSymbols = append(topSymbols, entry.Symbol)
		}
	}
	sort.Strings(topSymbols)

	terms := []string{
		"`Encoder`",
		"`EncoderID`",
		"`EncoderBcrypt`",
		"`EncoderArgon2`",
		"`EncoderScrypt`",
		"`EncoderPbkdf2`",
		"`EncoderMd5`",
		"`EncoderSha256`",
		"`EncoderPlaintext`",
		"`NewCompositeEncoder`",
		"`NewCipherEncoder`",
		"`UpgradeEncoding`",
		"`ErrDefaultEncoderNotFound`",
		"`ErrCipherRequired`",
		"`ErrEncoderRequired`",
		"`ErrInvalidHashFormat`",
		"`WithBcryptCost`",
		"`WithArgon2Memory`",
		"`WithScryptN`",
		"`WithPbkdf2HashFunction`",
		"`WithSha256SaltPosition`",
		"`WithMd5SaltPosition`",
		"`{bcrypt}`",
		"`{argon2}`",
		"`{sha256}`",
		"`{algorithm}encoded_value`",
		"`{algorithm}$salt$hash`",
		"`bcrypt.DefaultCost`",
		"`64 * 1024`",
		"`310000`",
		"`sha512`",
		"`prefix`",
		"`suffix`",
	}

	var failures []string
	for _, doc := range docs {
		for _, symbol := range topSymbols {
			if !strings.Contains(doc.content, "`"+symbol+"`") &&
				!strings.Contains(doc.content, "`password."+symbol+"`") &&
				!strings.Contains(doc.content, "password."+symbol) &&
				!strings.Contains(doc.content, symbol+"(") {
				failures = append(failures, doc.label+" missing top-level password symbol `"+symbol+"`")
			}
		}
		for _, term := range terms {
			if !containsNormalized(doc.content, term) {
				failures = append(failures, doc.label+" missing password semantic term "+term)
			}
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
			path: "password/password.go",
			terms: []string{
				"type EncoderID string",
				"EncoderBcrypt    EncoderID = \"bcrypt\"",
				"EncoderArgon2    EncoderID = \"argon2\"",
				"EncoderScrypt    EncoderID = \"scrypt\"",
				"EncoderPbkdf2    EncoderID = \"pbkdf2\"",
				"EncoderMd5       EncoderID = \"md5\"",
				"EncoderSha256    EncoderID = \"sha256\"",
				"EncoderPlaintext EncoderID = \"plaintext\"",
				"type Encoder interface",
				"Encode(password string) (string, error)",
				"Matches(password, encodedPassword string) bool",
				"UpgradeEncoding(encodedPassword string) bool",
			},
		},
		{
			path: "password/composite.go",
			terms: []string{
				"func NewCompositeEncoder(defaultEncoderID EncoderID, encoders map[EncoderID]Encoder) Encoder",
				"return fmt.Sprintf(\"{%s}%s\", c.defaultEncoderID, encoded), nil",
				"if encoderID == EncoderID(\"\")",
				"encoderID = c.defaultEncoderID",
				"if encoderID != EncoderID(\"\") && encoderID != c.defaultEncoderID",
				"return true",
				"return encoder.UpgradeEncoding(rawEncoded)",
				"if !strings.HasPrefix(encodedPassword, \"{\")",
				"strings.Index(encodedPassword, \"}\")",
			},
		},
		{
			path: "password/cipher.go",
			terms: []string{
				"func NewCipherEncoder(cipher cryptox.Cipher, encoder Encoder) Encoder",
				"if e.cipher == nil",
				"return \"\", ErrCipherRequired",
				"if e.encoder == nil",
				"return \"\", ErrEncoderRequired",
				"plainPassword, err := e.cipher.Decrypt(password)",
				"return e.encoder.Encode(plainPassword)",
				"return e.encoder.Matches(plainPassword, encodedPassword)",
				"return e.encoder.UpgradeEncoding(encodedPassword)",
			},
		},
		{
			path: "password/bcrypt.go",
			terms: []string{
				"bcryptMinCost = 4",
				"bcryptMaxCost = 31",
				"type BcryptOption func(*bcryptEncoder)",
				"func WithBcryptCost(cost int) BcryptOption",
				"cost: bcrypt.DefaultCost",
				"return cost < e.cost",
			},
		},
		{
			path: "password/argon2.go",
			terms: []string{
				"type Argon2Option func(*argon2Encoder)",
				"func WithArgon2Memory(memory uint32) Argon2Option",
				"func WithArgon2Iterations(iterations uint32) Argon2Option",
				"func WithArgon2Parallelism(parallelism uint8) Argon2Option",
				"memory:      64 * 1024",
				"iterations:  3",
				"parallelism: 4",
				"return fmt.Sprintf(\"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s\"",
				"return params.memory < e.memory || params.iterations < e.iterations || params.parallelism < e.parallelism",
			},
		},
		{
			path: "password/scrypt.go",
			terms: []string{
				"type ScryptOption func(*scryptEncoder)",
				"func WithScryptN(n int) ScryptOption",
				"func WithScryptR(r int) ScryptOption",
				"func WithScryptP(p int) ScryptOption",
				"n: 32768",
				"r: 8",
				"p: 1",
				"return fmt.Sprintf(\"$scrypt$n=%d,r=%d,p=%d$%s$%s\"",
				"return params.n < e.n || params.r < e.r || params.p < e.p",
			},
		},
		{
			path: "password/pbkdf2.go",
			terms: []string{
				"type Pbkdf2Option func(*pbkdf2Encoder)",
				"func WithPbkdf2Iterations(iterations int) Pbkdf2Option",
				"func WithPbkdf2HashFunction(hashFunction string) Pbkdf2Option",
				"iterations:   310000",
				"hashFunction: \"sha256\"",
				"case \"sha256\":",
				"case \"sha512\":",
				"return fmt.Sprintf(\"$pbkdf2-%s$i=%d$%s$%s\"",
				"return params.iterations < e.iterations || params.hashFunction != e.hashFunction",
			},
		},
		{
			path: "password/hash_encoder.go",
			terms: []string{
				"const saltPositionPrefix = \"prefix\"",
				"if e.saltPosition == saltPositionPrefix",
				"return salt + password",
				"return password + salt",
				"return fmt.Sprintf(\"{%s}$%s$%s\", e.algorithm, e.salt, hexHash), nil",
				"prefix := \"{\" + e.algorithm + \"}$\"",
				"func (*hashEncoder) UpgradeEncoding(string) bool",
				"return true",
			},
		},
		{
			path: "password/sha256.go",
			terms: []string{
				"type Sha256Option func(*hashEncoder)",
				"func WithSha256Salt(salt string) Sha256Option",
				"func WithSha256SaltPosition(position string) Sha256Option",
				"saltPosition: \"suffix\"",
				"algorithm:    \"sha256\"",
			},
		},
		{
			path: "password/md5.go",
			terms: []string{
				"type Md5Option func(*hashEncoder)",
				"func WithMd5Salt(salt string) Md5Option",
				"func WithMd5SaltPosition(position string) Md5Option",
				"saltPosition: \"suffix\"",
				"algorithm:    \"md5\"",
			},
		},
		{
			path: "password/plaintext.go",
			terms: []string{
				"func NewPlaintextEncoder() Encoder",
				"return password, nil",
				"return password == encodedPassword",
				"func (*plaintextEncoder) UpgradeEncoding(string) bool",
				"return true",
			},
		},
		{
			path: "password/errors.go",
			terms: []string{
				"ErrInvalidCost",
				"ErrInvalidMemory",
				"ErrInvalidIterations",
				"ErrInvalidParallelism",
				"ErrInvalidEncoderID",
				"ErrInvalidHashFormat",
				"ErrDefaultEncoderNotFound",
				"ErrCipherRequired",
				"ErrEncoderRequired",
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
	cmd := exec.Command("go", "test", "./password")
	cmd.Dir = sourceRoot
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return []string{fmt.Sprintf("go test ./password failed: %v\n%s", err, strings.TrimSpace(output.String()))}
	}

	return nil
}

func passwordEntriesFromAudit(audit auditLedger) []auditEntry {
	var entries []auditEntry
	for _, entry := range audit.Entries {
		if entry.Package == passwordPackage {
			entries = append(entries, entry)
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].ID < entries[j].ID })

	return entries
}

func loadLiveInventory(sourceRoot, docsRoot, manifestPath, auditLedgerPath, contractLedgerPath string) map[string]manifestEntry {
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

	payload := "[" + strings.TrimSpace(string(output)) + "]"
	var entries []liveInventoryEntry
	if err := json.Unmarshal([]byte(payload), &entries); err != nil {
		panic(fmt.Errorf("parse live inventory: %w", err))
	}

	result := map[string]manifestEntry{}
	for _, entry := range entries {
		result[entry.Package] = manifestEntry{
			Package:     entry.Package,
			Coverage:    entry.Coverage,
			TopLevel:    entry.TopLevel,
			Fields:      entry.Fields,
			Methods:     entry.Methods,
			Fingerprint: entry.Fingerprint,
		}
	}

	return result
}

func loadLiveAuditLedger(sourceRoot, docsRoot, manifestPath string) auditLedger {
	cmd := exec.Command(
		"go", "run", filepath.Join(docsRoot, "scripts/verify-api-audit.go"),
		"-source", sourceRoot,
		"-manifest", manifestPath,
		"-print-ledger",
	)
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Errorf("verify-api-audit -print-ledger failed: %w\n%s", err, strings.TrimSpace(string(output))))
	}

	var ledger auditLedger
	if err := json.Unmarshal(output, &ledger); err != nil {
		panic(fmt.Errorf("parse live audit ledger: %w", err))
	}

	return ledger
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

func entriesByID(entries []auditEntry) map[string]auditEntry {
	result := map[string]auditEntry{}
	for _, entry := range entries {
		result[entry.ID] = entry
	}

	return result
}

func receiverForSymbol(symbol string) (string, bool) {
	receiver, _, ok := strings.Cut(symbol, ".")
	if !ok || receiver == "" {
		return "", false
	}

	return receiver, true
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

func sameSet(got, want []string) bool {
	got = sortedUnique(got)
	want = sortedUnique(want)
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}

	return true
}

func sortedUnique(values []string) []string {
	set := map[string]bool{}
	for _, value := range values {
		set[value] = true
	}
	result := make([]string, 0, len(set))
	for value := range set {
		result = append(result, value)
	}
	sort.Strings(result)

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

func passwordCoverage() []string {
	return []string{englishPasswordPath}
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}
