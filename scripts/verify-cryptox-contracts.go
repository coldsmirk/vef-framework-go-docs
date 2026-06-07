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
	cryptoxPackage = "github.com/coldsmirk/vef-framework-go/cryptox"

	cryptoxFingerprint = "10f32dd3a1d0db4bcce8f4dd13998ec92edcb6954a96ef80e3473f38b322718b"
	cryptoxTopLevel    = 77
	cryptoxFields      = 0
	cryptoxMethods     = 8
	cryptoxEntries     = 85

	cryptoxGroupedEntries              = 8
	cryptoxGroupedFields               = 0
	cryptoxGroupedMethods              = 8
	cryptoxGroupedReceivers            = 3
	cryptoxGroupedSignatureFingerprint = "802b292a995b81ebae5e8e5d583c627c149178d9068166ed301117a228ba9f90"
	cryptoxGroupedReceiverFingerprint  = "b0289cc9ebdef776d3592c4bf53ff8caf6467fd34fa1b41680b29511d8f7b021"

	englishCryptoxPath = "docs/features/cryptox.md"
	chineseCryptoxPath = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/features/cryptox.md"
	englishIndexPath   = "docs/reference/public-api-index.md"
	chineseIndexPath   = "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"
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

	englishDocs := readCorpus("English cryptox docs", filepath.Join(docsRoot, englishCryptoxPath))
	chineseDocs := readCorpus("Chinese cryptox docs", filepath.Join(docsRoot, chineseCryptoxPath))
	englishIndex := readCorpus("English public API index", filepath.Join(docsRoot, englishIndexPath))
	chineseIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, chineseIndexPath))

	audit := loadJSON[auditLedger](auditLedgerPath)
	manifestData := loadJSON[manifest](manifestPath)
	contracts := loadJSON[contractLedger](contractLedgerPath)
	liveManifestEntry := loadLiveInventory(sourceRoot, docsRoot, manifestPath, auditLedgerPath, contractLedgerPath)[cryptoxPackage]
	liveAuditEntries := cryptoxEntriesFromAudit(loadLiveAuditLedger(sourceRoot, docsRoot, manifestPath))
	cryptoxEntries := cryptoxEntriesFromAudit(audit)

	var failures []string
	failures = append(failures, verifySurfaceEntry("live public API inventory", liveManifestEntry)...)
	failures = append(failures, verifyManifest(manifestData)...)
	failures = append(failures, verifyContractLedger(contracts, sourceRoot)...)
	failures = append(failures, verifyAuditEntries(cryptoxEntries)...)
	failures = append(failures, verifyLiveAuditEntries(cryptoxEntries, liveAuditEntries)...)
	failures = append(failures, verifyGroupedCryptoxSurface(cryptoxEntries, []corpus{englishDocs, chineseDocs})...)
	failures = append(failures, verifyGeneratedIndexSection(englishIndex, cryptoxEntries)...)
	failures = append(failures, verifyGeneratedIndexSection(chineseIndex, cryptoxEntries)...)
	failures = append(failures, verifyCryptoxDocs(cryptoxEntries, []corpus{englishDocs, chineseDocs})...)
	failures = append(failures, verifySourceTerms(sourceRoot)...)
	failures = append(failures, runGoTests(sourceRoot)...)

	if len(failures) > 0 {
		sort.Strings(failures)
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Println("Cryptox contract docs verified: 85 public entries, 8 grouped cipher/signer methods, algorithm modes and key encodings")
}

func verifySurfaceEntry(label string, entry manifestEntry) []string {
	var failures []string
	if entry.Package != cryptoxPackage {
		failures = append(failures, fmt.Sprintf("%s package mismatch: got %q", label, entry.Package))
	}
	if entry.TopLevel != cryptoxTopLevel ||
		entry.Fields != cryptoxFields ||
		entry.Methods != cryptoxMethods ||
		entry.Fingerprint != cryptoxFingerprint {
		failures = append(failures, fmt.Sprintf(
			"%s cryptox surface mismatch: got top/fields/methods/fingerprint=%d/%d/%d/%s want=%d/%d/%d/%s",
			label,
			entry.TopLevel, entry.Fields, entry.Methods, entry.Fingerprint,
			cryptoxTopLevel, cryptoxFields, cryptoxMethods, cryptoxFingerprint,
		))
	}

	return failures
}

func verifyManifest(m manifest) []string {
	for _, entry := range m.Packages {
		if entry.Package != cryptoxPackage {
			continue
		}
		var failures []string
		failures = append(failures, verifySurfaceEntry("API audit manifest", entry)...)
		if !sameSet(entry.Coverage, cryptoxCoverage()) {
			failures = append(failures, fmt.Sprintf("cryptox manifest coverage mismatch: got %v want %v", entry.Coverage, cryptoxCoverage()))
		}

		return failures
	}

	return []string{"API audit manifest missing cryptox package"}
}

func verifyContractLedger(contracts contractLedger, sourceRoot string) []string {
	contractID := cryptoxPackage + "#runtime-contract:crypto-algorithm-modes-and-encodings"
	expectedTerms := []string{
		"CipherSigner",
		"AesModeGcm",
		"RsaModeOAEP",
		"SM4ModeCBC",
		"EciesCurveP256",
	}

	var failures []string
	var foundReview bool
	for _, review := range contracts.PackageReviews {
		if review.Package != cryptoxPackage {
			continue
		}
		foundReview = true
		if review.Disposition != "has-semantic-contracts" {
			failures = append(failures, "cryptox contract review disposition mismatch: "+review.Disposition)
		}
		if review.ReviewedSurface.TopLevel != cryptoxTopLevel ||
			review.ReviewedSurface.Fields != cryptoxFields ||
			review.ReviewedSurface.Methods != cryptoxMethods ||
			review.ReviewedSurface.EntryCount != cryptoxEntries ||
			review.ReviewedSurface.Fingerprint != cryptoxFingerprint {
			failures = append(failures, fmt.Sprintf(
				"cryptox contract review surface mismatch: got top/fields/methods/entries/fingerprint=%d/%d/%d/%d/%s",
				review.ReviewedSurface.TopLevel,
				review.ReviewedSurface.Fields,
				review.ReviewedSurface.Methods,
				review.ReviewedSurface.EntryCount,
				review.ReviewedSurface.Fingerprint,
			))
		}
		if !sameSet(review.Coverage, cryptoxCoverage()) {
			failures = append(failures, fmt.Sprintf("cryptox contract review coverage mismatch: got %v want %v", review.Coverage, cryptoxCoverage()))
		}
		if !sameSet(review.ContractIDs, []string{contractID}) {
			failures = append(failures, fmt.Sprintf("cryptox contract ids mismatch: got %v want %v", review.ContractIDs, []string{contractID}))
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, review.SourceEvidence)...)
	}
	if !foundReview {
		failures = append(failures, "contract ledger missing cryptox package review")
	}

	var foundEntry bool
	for _, entry := range contracts.Entries {
		if entry.ID != contractID {
			continue
		}
		foundEntry = true
		if entry.Package != cryptoxPackage {
			failures = append(failures, "cryptox contract entry package mismatch: "+entry.Package)
		}
		if entry.Disposition != "documented:semantic-contract" {
			failures = append(failures, "cryptox contract entry disposition mismatch: "+entry.Disposition)
		}
		if !sameSet(entry.Coverage, cryptoxCoverage()) {
			failures = append(failures, fmt.Sprintf("cryptox contract coverage mismatch: got %v want %v", entry.Coverage, cryptoxCoverage()))
		}
		for _, term := range expectedTerms {
			if !contains(entry.Terms, term) {
				failures = append(failures, "cryptox contract entry missing term "+term)
			}
		}
		failures = append(failures, verifySourceEvidence(sourceRoot, entry.SourceEvidence)...)
	}
	if !foundEntry {
		failures = append(failures, "contract ledger missing cryptox contract entry "+contractID)
	}

	return failures
}

func verifyAuditEntries(entries []auditEntry) []string {
	var failures []string
	if len(entries) != cryptoxEntries {
		failures = append(failures, fmt.Sprintf("cryptox audit entry count mismatch: got %d want %d", len(entries), cryptoxEntries))
	}
	counts := map[string]int{}
	dispositionCounts := map[string]int{}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.Package != cryptoxPackage {
			failures = append(failures, "non-cryptox audit entry passed into cryptox verifier: "+entry.ID)
		}
		if seen[entry.ID] {
			failures = append(failures, "duplicate cryptox audit entry "+entry.ID)
		}
		seen[entry.ID] = true
		counts[entry.Kind]++
		dispositionCounts[entry.Disposition]++
		if entry.Symbol == "" || entry.Signature == "" || entry.Disposition == "" {
			failures = append(failures, "cryptox audit entry missing required metadata "+entry.ID)
		}
		if !sameSet(entry.Coverage, cryptoxCoverage()) {
			failures = append(failures, fmt.Sprintf("cryptox audit entry %s coverage mismatch: got %v want %v", entry.ID, entry.Coverage, cryptoxCoverage()))
		}
	}
	if counts["top"] != cryptoxTopLevel || counts["field"] != cryptoxFields || counts["method"] != cryptoxMethods {
		failures = append(failures, fmt.Sprintf("cryptox audit kind counts mismatch: top/field/method=%d/%d/%d", counts["top"], counts["field"], counts["method"]))
	}
	if dispositionCounts["documented:top-level"] != cryptoxTopLevel ||
		dispositionCounts["grouped:type-member-family"] != cryptoxGroupedEntries {
		failures = append(failures, fmt.Sprintf(
			"cryptox audit disposition counts mismatch: top-level/grouped=%d/%d want=%d/%d",
			dispositionCounts["documented:top-level"],
			dispositionCounts["grouped:type-member-family"],
			cryptoxTopLevel,
			cryptoxGroupedEntries,
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
			failures = append(failures, fmt.Sprintf("cryptox missing_in_ledger: %s %s %s", id, live.Symbol, live.Signature))
			continue
		}
		if ledger.Kind != live.Kind || ledger.Symbol != live.Symbol || ledger.Signature != live.Signature {
			failures = append(failures, fmt.Sprintf(
				"cryptox live/ledger signature drift for %s: ledger=%s/%s/%s live=%s/%s/%s",
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
			failures = append(failures, fmt.Sprintf("cryptox extra_in_ledger: %s %s %s", id, ledger.Symbol, ledger.Signature))
		}
	}

	return failures
}

func verifyGroupedCryptoxSurface(entries []auditEntry, docs []corpus) []string {
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
			failures = append(failures, fmt.Sprintf("cryptox grouped entry has non receiver-qualified symbol %q", entry.Symbol))
			continue
		}
		rows = append(rows, strings.Join([]string{entry.Symbol, entry.Kind, entry.Disposition, entry.Signature}, "\t"))
		receiverCounts[receiver]++
		kindCounts[entry.Kind]++
	}

	failures = append(failures, verifyGroupedFingerprint("cryptox grouped type-member surface", rows, cryptoxGroupedEntries, cryptoxGroupedSignatureFingerprint)...)
	if kindCounts["field"] != cryptoxGroupedFields || kindCounts["method"] != cryptoxGroupedMethods {
		failures = append(failures, fmt.Sprintf(
			"cryptox grouped kind counts mismatch: got fields/methods=%d/%d want=%d/%d",
			kindCounts["field"],
			kindCounts["method"],
			cryptoxGroupedFields,
			cryptoxGroupedMethods,
		))
	}

	var receiverRows []string
	for receiver, count := range receiverCounts {
		receiverRows = append(receiverRows, fmt.Sprintf("%d %s", count, receiver))
	}
	failures = append(failures, verifyGroupedFingerprint("cryptox grouped receiver/type families", receiverRows, cryptoxGroupedReceivers, cryptoxGroupedReceiverFingerprint)...)

	for _, doc := range docs {
		for _, term := range []string{
			"85 public cryptox entries",
			"8 grouped cryptox method entries",
			"3 cryptox receiver/type families",
			"0 exported cryptox field entries",
			"8 exported cryptox method entries",
		} {
			if !containsNormalized(doc.content, term) {
				failures = append(failures, doc.label+" missing grouped cryptox audit term "+term)
			}
		}
	}

	return failures
}

func verifyGeneratedIndexSection(index corpus, entries []auditEntry) []string {
	section := packageSection(index.content, cryptoxPackage)
	if section == "" {
		return []string{index.label + " missing cryptox package section"}
	}

	var failures []string
	for _, entry := range entries {
		if !strings.Contains(section, entry.Signature) {
			failures = append(failures, fmt.Sprintf("%s cryptox index missing signature for %s: %s", index.label, entry.ID, entry.Signature))
		}
	}

	return failures
}

func verifyCryptoxDocs(entries []auditEntry, docs []corpus) []string {
	var topSymbols []string
	for _, entry := range entries {
		if entry.Kind == "top" {
			topSymbols = append(topSymbols, entry.Symbol)
		}
	}
	sort.Strings(topSymbols)

	terms := []string{
		"`Cipher`",
		"`Signer`",
		"`CipherSigner`",
		"`AesModeGcm`",
		"`AesModeCbc`",
		"`RsaModeOAEP`",
		"`RsaModePKCS1v15`",
		"`RsaSignModePSS`",
		"`RsaSignModePKCS1v15`",
		"`SM4ModeCBC`",
		"`SM4ModeECB`",
		"`EciesCurveP256`",
		"`EciesCurveX25519`",
		"`EcdsaCurveP256`",
		"NewAESFromHex",
		"NewAESFromBase64",
		"NewRSAFromPEM",
		"NewSM2FromPEM",
		"NewECDSAFromPEM",
		"NewECIESFromBytes",
		"GenerateECDSAKey",
		"GenerateECIESKey",
		"`ErrInvalidAESKeySize`",
		"`ErrInvalidSM4KeySize`",
		"`ErrInvalidSignature`",
		"Base64",
		"Hex",
		"PEM",
		"P-256",
	}

	var failures []string
	for _, doc := range docs {
		for _, symbol := range topSymbols {
			if !strings.Contains(doc.content, "`"+symbol+"`") &&
				!strings.Contains(doc.content, "`cryptox."+symbol+"`") &&
				!strings.Contains(doc.content, "cryptox."+symbol) &&
				!strings.Contains(doc.content, symbol+"(") {
				failures = append(failures, doc.label+" missing top-level cryptox symbol `"+symbol+"`")
			}
		}
		for _, term := range terms {
			if !containsNormalized(doc.content, term) {
				failures = append(failures, doc.label+" missing cryptox semantic term "+term)
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
			path: "cryptox/cipher.go",
			terms: []string{
				"type Cipher interface",
				"Encrypt(plaintext string) (string, error)",
				"Decrypt(ciphertext string) (string, error)",
				"type Signer interface",
				"Sign(data string) (signature string, err error)",
				"Verify(data, signature string) (bool, error)",
				"type CipherSigner interface",
				"Cipher",
				"Signer",
			},
		},
		{
			path: "cryptox/errors.go",
			terms: []string{
				"ErrAtLeastOneKeyRequired",
				"ErrPublicKeyRequiredForEncrypt",
				"ErrPrivateKeyRequiredForDecrypt",
				"ErrPrivateKeyRequiredForSign",
				"ErrPublicKeyRequiredForVerify",
				"ErrInvalidAESKeySize",
				"ErrInvalidSM4KeySize",
				"ErrInvalidSignature",
			},
		},
		{
			path: "cryptox/aes_cipher.go",
			terms: []string{
				"type AESMode string",
				"AesModeCbc AESMode = \"CBC\"",
				"AesModeGcm AESMode = \"GCM\"",
				"type AESOption func(*aesCipher)",
				"func WithAESIv(iv []byte) AESOption",
				"func WithAESMode(mode AESMode) AESOption",
				"func NewAES(key []byte, opts ...AESOption) (Cipher, error)",
				"len(key) != 16 && len(key) != 24 && len(key) != 32",
				"mode: AesModeGcm",
				"if cipher.mode == AesModeCbc",
				"func NewAESFromHex(keyHex string, opts ...AESOption) (Cipher, error)",
				"func NewAESFromBase64(keyBase64 string, opts ...AESOption) (Cipher, error)",
				"base64.StdEncoding.EncodeToString(ciphertext)",
			},
		},
		{
			path: "cryptox/rsa_cipher.go",
			terms: []string{
				"type RSAMode string",
				"RsaModeOAEP     RSAMode = \"OAEP\"",
				"RsaModePKCS1v15 RSAMode = \"PKCS1v15\"",
				"type RSASignMode string",
				"RsaSignModePSS      RSASignMode = \"PSS\"",
				"RsaSignModePKCS1v15 RSASignMode = \"PKCS1v15\"",
				"type RSAOption func(*rsaCipher)",
				"func WithRSAMode(mode RSAMode) RSAOption",
				"func WithRSASignMode(signMode RSASignMode) RSAOption",
				"func NewRSA(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, opts ...RSAOption) (CipherSigner, error)",
				"mode:       RsaModeOAEP",
				"signMode:   RsaSignModePSS",
				"if publicKey == nil && privateKey != nil",
				"func NewRSAFromPEM(privatePEM, publicPEM []byte, opts ...RSAOption) (CipherSigner, error)",
				"func NewRSAFromHex(privateKeyHex, publicKeyHex string, opts ...RSAOption) (CipherSigner, error)",
				"func NewRSAFromBase64(privateKeyBase64, publicKeyBase64 string, opts ...RSAOption) (CipherSigner, error)",
				"rsa.EncryptOAEP",
				"rsa.SignPSS",
				"fmt.Errorf(\"%w: %w\", ErrInvalidSignature, err)",
			},
		},
		{
			path: "cryptox/sm4_cipher.go",
			terms: []string{
				"type SM4Mode string",
				"SM4ModeCBC SM4Mode = \"CBC\"",
				"SM4ModeECB SM4Mode = \"ECB\"",
				"type SM4Option func(*sm4Cipher)",
				"func WithSM4Iv(iv []byte) SM4Option",
				"func WithSM4Mode(mode SM4Mode) SM4Option",
				"func NewSM4(key []byte, opts ...SM4Option) (Cipher, error)",
				"mode: SM4ModeCBC",
				"if cipher.mode == SM4ModeCBC",
				"func NewSM4FromHex(keyHex string, opts ...SM4Option) (Cipher, error)",
				"func NewSM4FromBase64(keyBase64 string, opts ...SM4Option) (Cipher, error)",
				"sm4.Sm4Ecb",
			},
		},
		{
			path: "cryptox/ecdsa_cipher.go",
			terms: []string{
				"type ECDSACurve string",
				"EcdsaCurveP224 ECDSACurve = \"P224\"",
				"EcdsaCurveP256 ECDSACurve = \"P256\"",
				"EcdsaCurveP384 ECDSACurve = \"P384\"",
				"EcdsaCurveP521 ECDSACurve = \"P521\"",
				"func NewECDSA(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (Signer, error)",
				"func NewECDSAFromPEM(privatePEM, publicPEM []byte) (Signer, error)",
				"func NewECDSAFromHex(privateKeyHex, publicKeyHex string) (Signer, error)",
				"func NewECDSAFromBase64(privateKeyBase64, publicKeyBase64 string) (Signer, error)",
				"func GenerateECDSAKey(curve ECDSACurve) (*ecdsa.PrivateKey, error)",
				"ellipticCurve = elliptic.P256()",
				"asn1.Marshal(ecdsaSignature{R: r, S: s})",
				"fmt.Errorf(\"%w: %w\", ErrInvalidSignature, err)",
			},
		},
		{
			path: "cryptox/sm2_cipher.go",
			terms: []string{
				"func NewSM2(privateKey *sm2.PrivateKey, publicKey *sm2.PublicKey) (CipherSigner, error)",
				"func NewSM2FromPEM(privatePEM, publicPEM []byte) (CipherSigner, error)",
				"func NewSM2FromHex(privateKeyHex, publicKeyHex string) (CipherSigner, error)",
				"func NewSM2FromBase64(privateKeyBase64, publicKeyBase64 string) (CipherSigner, error)",
				"sm2.Encrypt(s.publicKey, []byte(plaintext), rand.Reader, sm2.C1C3C2)",
				"sm2.Decrypt(s.privateKey, encryptedData, sm2.C1C3C2)",
				"s.privateKey.Sign(rand.Reader, []byte(data), nil)",
				"s.publicKey.Verify([]byte(data), signatureBytes)",
			},
		},
		{
			path: "cryptox/ecies_cipher.go",
			terms: []string{
				"type ECIESCurve string",
				"EciesCurveP256   ECIESCurve = \"P256\"",
				"EciesCurveP384   ECIESCurve = \"P384\"",
				"EciesCurveP521   ECIESCurve = \"P521\"",
				"EciesCurveX25519 ECIESCurve = \"X25519\"",
				"func NewECIES(privateKey *ecdh.PrivateKey, publicKey *ecdh.PublicKey) (Cipher, error)",
				"func NewECIESFromBytes(privateKeyBytes, publicKeyBytes []byte, curve ECIESCurve) (Cipher, error)",
				"func NewECIESFromHex(privateKeyHex, publicKeyHex string, curve ECIESCurve) (Cipher, error)",
				"func NewECIESFromBase64(privateKeyBase64, publicKeyBase64 string, curve ECIESCurve) (Cipher, error)",
				"func GenerateECIESKey(curve ECIESCurve) (*ecdh.PrivateKey, error)",
				"EciesCurveX25519: ecdh.X25519()",
				"return ecdh.P256()",
				"hkdf.New(sha256.New, sharedSecret, nil, nil)",
				"gcm.Seal(nil, nonce, []byte(plaintext), nil)",
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
	cmd := exec.Command("go", "test", "./cryptox")
	cmd.Dir = sourceRoot
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return []string{fmt.Sprintf("go test ./cryptox failed: %v\n%s", err, strings.TrimSpace(output.String()))}
	}

	return nil
}

func cryptoxEntriesFromAudit(audit auditLedger) []auditEntry {
	var entries []auditEntry
	for _, entry := range audit.Entries {
		if entry.Package == cryptoxPackage {
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

func cryptoxCoverage() []string {
	return []string{englishCryptoxPath}
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}
