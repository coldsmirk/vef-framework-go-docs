package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	auditLedgerPath = "scripts/api-audit-ledger.json"
	securityPackage = "github.com/coldsmirk/vef-framework-go/security"

	securityGroupedEntryCount           = 174
	securityGroupedFieldCount           = 75
	securityGroupedMethodCount          = 99
	securityGroupedReceiverCount        = 66
	securityGroupedSignatureFingerprint = "929d935fa6c1a0035edeb168ad53f46e62a29c2b83fee32ae322131e53a8cb9b"
	securityGroupedReceiverFingerprint  = "6064db54b11081b0eec6e15754a4fae7d0995dfb3b69bc887dcaf09b570b3637"
)

type corpus struct {
	label   string
	content string
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

type check struct {
	sourcePath      string
	sourceTerms     []string
	docTerms        []string
	englishDocTerms []string
	chineseDocTerms []string
}

func main() {
	sourceDir := flag.String("source", "../vef-framework-go", "path to the VEF Framework Go source repository")
	outDir := flag.String("out", ".", "path to the VEF Framework Go docs repository")
	flag.Parse()

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)
	authDocs := []corpus{
		readCorpus("English authentication docs", filepath.Join(docsRoot, "docs/security/authentication.md")),
		readCorpus("Chinese authentication docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/security/authentication.md")),
	}
	authzDocs := []corpus{
		readCorpus("English authorization docs", filepath.Join(docsRoot, "docs/security/authorization.md")),
		readCorpus("Chinese authorization docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/security/authorization.md")),
	}
	dataPermDocs := []corpus{
		readCorpus("English data permission docs", filepath.Join(docsRoot, "docs/security/data-permissions.md")),
		readCorpus("Chinese data permission docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/security/data-permissions.md")),
	}
	audit := loadJSON[auditLedger](filepath.Join(docsRoot, auditLedgerPath))

	var failures []string
	failures = append(failures, verifyGroupedSecuritySurface(audit, authDocs)...)
	failures = append(failures, runChecks(sourceRoot, authDocs, authChecks())...)
	failures = append(failures, runChecks(sourceRoot, authzDocs, authzChecks())...)
	failures = append(failures, runChecks(sourceRoot, dataPermDocs, dataPermChecks())...)

	sort.Strings(failures)
	if len(failures) > 0 {
		panic(fmt.Errorf("security contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	fmt.Printf("Security contract docs verified: 174 grouped field/method entries, %d auth docs, %d authorization docs, %d data-permission docs\n",
		len(authDocs), len(authzDocs), len(dataPermDocs))
}

func verifyGroupedSecuritySurface(audit auditLedger, docs []corpus) []string {
	var rows []string
	receiverCounts := map[string]int{}
	kindCounts := map[string]int{}
	var failures []string
	for _, entry := range audit.Entries {
		if entry.Package != securityPackage || !strings.HasPrefix(entry.Disposition, "grouped:") {
			continue
		}

		receiver, ok := receiverForSymbol(entry.Symbol)
		if !ok {
			failures = append(failures, fmt.Sprintf("security grouped entry has non receiver-qualified symbol %q", entry.Symbol))
			continue
		}

		rows = append(rows, strings.Join([]string{entry.Symbol, entry.Kind, entry.Signature}, "\t"))
		receiverCounts[receiver]++
		kindCounts[entry.Kind]++
	}

	failures = append(failures, verifyGroupedFingerprint(
		"security grouped field/method surface",
		rows,
		securityGroupedEntryCount,
		securityGroupedSignatureFingerprint,
	)...)
	if kindCounts["field"] != securityGroupedFieldCount || kindCounts["method"] != securityGroupedMethodCount {
		failures = append(failures, fmt.Sprintf(
			"security grouped kind counts mismatch: got fields/methods=%d/%d want=%d/%d",
			kindCounts["field"],
			kindCounts["method"],
			securityGroupedFieldCount,
			securityGroupedMethodCount,
		))
	}

	receiverRows := make([]string, 0, len(receiverCounts))
	for receiver, count := range receiverCounts {
		receiverRows = append(receiverRows, fmt.Sprintf("%d %s", count, receiver))
	}
	failures = append(failures, verifyGroupedFingerprint(
		"security grouped receiver/type families",
		receiverRows,
		securityGroupedReceiverCount,
		securityGroupedReceiverFingerprint,
	)...)

	for _, doc := range docs {
		for _, term := range []string{
			"174 grouped security field/method entries",
			"66 receiver/type families",
			"75 public field entries",
			"99 public method entries",
		} {
			if !containsNormalized(doc.content, term) {
				failures = append(failures, doc.label+" missing grouped security audit term "+term)
			}
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

func receiverForSymbol(symbol string) (string, bool) {
	receiver, _, ok := strings.Cut(symbol, ".")
	if !ok || receiver == "" {
		return "", false
	}

	return receiver, true
}

func authChecks() []check {
	return []check{
		{
			sourcePath: "internal/security/auth_resource.go",
			sourceTerms: []string{
				"api.NewRPCResource(\n\t\t\t\"security/auth\"",
				"Action:    \"login\"",
				"Action:    \"refresh\"",
				"Action: \"logout\"",
				"Action:    \"resolve_challenge\"",
				"Action: \"get_user_info\"",
				"RateLimit: &api.RateLimitConfig{Max: params.SecurityConfig.LoginRateLimit}",
				"RateLimit: &api.RateLimitConfig{Max: params.SecurityConfig.RefreshRateLimit}",
				"validate:\"required\" label_i18n:\"auth_type\"",
				"validate:\"required\" label_i18n:\"auth_refresh_token\"",
				"validate:\"required\" label_i18n:\"auth_challenge_token\"",
				"return result.Ok().Response(ctx)",
				"a.userInfoLoader.LoadUserInfo(ctx.Context(), principal, params)",
			},
			docTerms: []string{
				"security/auth",
				"`login`",
				"`refresh`",
				"`logout`",
				"`resolve_challenge`",
				"`get_user_info`",
				"vef.security.login_rate_limit",
				"vef.security.refresh_rate_limit",
				"validate:\"required\"",
				"`UserInfoLoader.LoadUserInfo(...)`",
			},
			englishDocTerms: []string{
				"`logout` returns an ok result immediately",
				"does not revoke or blacklist a\nserver-side token",
				"arbitrary `params`, forwarded to `UserInfoLoader.LoadUserInfo(...)`",
			},
			chineseDocTerms: []string{
				"`logout` 会立即返回 ok 结果",
				"不会在服务端吊销或拉黑 token",
				"任意 `params`，会转发给 `UserInfoLoader.LoadUserInfo(...)`",
			},
		},
		{
			sourcePath: "security/security.go",
			sourceTerms: []string{
				"`json:\"accessToken\"`",
				"`json:\"refreshToken\"`",
				"`json:\"type\"`",
				"`json:\"principal\"`",
				"`json:\"credentials\"`",
				"`json:\"enabled\"`",
				"`json:\"ipWhitelist\"`",
			},
			docTerms: []string{
				"`accessToken`",
				"`refreshToken`",
				"`type`",
				"`principal`",
				"`credentials`",
				"`ExternalAppConfig.IPWhitelist`",
			},
		},
		{
			sourcePath: "security/jwt.go",
			sourceTerms: []string{
				"JWTIssuer          = \"vef\"",
				"DefaultJWTAudience = \"vef-app\"",
				"DefaultJWTSecret",
				"jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name})",
				"jwt.WithLeeway(10 * time.Second)",
				"jwt.WithIssuedAt()",
				"jwt.WithExpirationRequired()",
				"hex.DecodeString(cmp.Or(config.Secret, DefaultJWTSecret))",
				"config.Audience = cmp.Or(config.Audience, DefaultJWTAudience)",
			},
			docTerms: []string{
				"`JWTIssuer`",
				"`vef`",
				"`DefaultJWTAudience`",
				"`DefaultJWTSecret`",
				"`HS256`",
				"10",
				"leeway",
				"`iat`",
				"`exp`",
			},
		},
		{
			sourcePath: "security/jwt_config.go",
			sourceTerms: []string{
				"claimJWTID     = \"jti\"",
				"claimSubject   = \"sub\"",
				"claimIssuer    = \"iss\"",
				"claimAudience  = \"aud\"",
				"claimIssuedAt  = \"iat\"",
				"claimNotBefore = \"nbf\"",
				"claimExpiresAt = \"exp\"",
				"claimType      = \"typ\"",
				"claimRoles     = \"rls\"",
				"claimDetails   = \"det\"",
			},
			docTerms: []string{
				"`jti`",
				"`sub`",
				"`iss`",
				"`aud`",
				"`iat`",
				"`nbf`",
				"`exp`",
				"`typ`",
				"`rls`",
				"`det`",
			},
		},
		{
			sourcePath: "internal/security/jwt_token_generator.go",
			sourceTerms: []string{
				"AccessTokenExpires  = time.Minute * 30",
				"RefreshTokenExpires = time.Hour * 24 * 7",
				"WithSubject(fmt.Sprintf(\"%s@%s\", principal.ID, principal.Name))",
				"WithType(security.TokenTypeAccess)",
				"WithType(security.TokenTypeRefresh)",
				"return g.jwt.Generate(claimsBuilder, AccessTokenExpires, 0)",
				"return g.jwt.Generate(claimsBuilder, g.refreshExpires, g.refreshNotBefore)",
			},
			docTerms: []string{
				"`30m`",
				"`168h`",
				"`id@name`",
				"`TokenTypeAccess`",
				"`TokenTypeRefresh`",
				"`jti`",
			},
		},
		{
			sourcePath: "security/principal.go",
			sourceTerms: []string{
				"PrincipalTypeUser PrincipalType = \"user\"",
				"PrincipalTypeExternalApp PrincipalType = \"external_app\"",
				"PrincipalTypeSystem PrincipalType = orm.OperatorSystem",
				"Name: \"系统\"",
				"PrincipalAnonymous = NewUser(orm.OperatorAnonymous, \"匿名\")",
				"`json:\"type\"`",
				"`json:\"id\"`",
				"`json:\"name\"`",
				"`json:\"roles\"`",
				"`json:\"details\"`",
				"ErrUserDetailsNotStruct",
				"ErrExternalAppDetailsNotStruct",
				"p.Details = nil",
				"var detailsMap map[string]any",
			},
			docTerms: []string{
				"`user`",
				"`external_app`",
				"`system`",
				"`PrincipalSystem`",
				"`PrincipalAnonymous`",
				"`系统`",
				"`匿名`",
				"`ErrUserDetailsNotStruct`",
				"`ErrExternalAppDetailsNotStruct`",
				"`map[string]any`",
				"`nil`",
			},
		},
		{
			sourcePath: "security/jwt_challenge_token_store.go",
			sourceTerms: []string{
				"ChallengeTokenExpires       = 5 * time.Minute",
				"ClaimChallengePending       = \"pnd\"",
				"ClaimChallengePrincipalType = \"ptp\"",
				"ClaimChallengeResolved      = \"rsd\"",
				"ClaimChallengePrincipalName = \"pnm\"",
				"WithType(TokenTypeChallenge)",
				"WithSubject(principal.ID)",
				"claimsAccessor.Type() != TokenTypeChallenge",
			},
			docTerms: []string{
				"`ChallengeTokenExpires`",
				"`pnd`",
				"`ptp`",
				"`rsd`",
				"`pnm`",
				"`typ: \"challenge\"`",
			},
			englishDocTerms: []string{
				"subject is principal ID only",
			},
			chineseDocTerms: []string{
				"subject 只保存 principal ID",
			},
		},
		{
			sourcePath: "security/redis_challenge_token_store.go",
			sourceTerms: []string{
				"const redisChallengePrefix = \"vef:security:challenge:\"",
				"id.GenerateUUID()",
				"s.client.Set(ctx, key, data, ChallengeTokenExpires)",
			},
			docTerms: []string{
				"`vef:security:challenge:<token>`",
				"UUID token",
				"`ChallengeTokenExpires`",
			},
		},
		{
			sourcePath: "security/otp.go",
			sourceTerms: []string{
				"ChallengeTypeSMS   = \"sms_otp\"",
				"ChallengeTypeEmail = \"email_otp\"",
				"ChallengeType:  challengeType",
				"ChallengeOrder: order",
				"panic(\"security: OTPChallengeProviderConfig.ChallengeType is required\")",
				"panic(\"security: OTPChallengeProviderConfig.Evaluator is required\")",
				"panic(\"security: OTPChallengeProviderConfig.Verifier is required\")",
				"return nil, ErrOTPCodeRequired",
				"return nil, ErrOTPCodeInvalid",
			},
			docTerms: []string{
				"`sms_otp`",
				"`email_otp`",
				"`ChallengeType`",
				"`ChallengeOrder`",
				"panic",
				"`ErrOTPCodeRequired`",
				"`ErrOTPCodeInvalid`",
			},
		},
		{
			sourcePath: "security/totp.go",
			sourceTerms: []string{
				"ChallengeTypeTOTP      = \"totp\"",
				"TOTPDefaultDestination = \"Authenticator App\"",
				"return totp.Validate(code, secret), nil",
				"ChallengeOrder: 100",
			},
			docTerms: []string{
				"`totp`",
				"`Authenticator App`",
				"`100`",
			},
		},
		{
			sourcePath: "security/password_change.go",
			sourceTerms: []string{
				"ChallengeTypePasswordChange = \"password_change\"",
				"PasswordChangeReasonFirstLogin = \"first_login\"",
				"PasswordChangeReasonExpired    = \"expired\"",
				"func (*PasswordChangeChallengeProvider) Order() int   { return 400 }",
				"panic(\"security: PasswordChangeChecker is required\")",
				"panic(\"security: PasswordChanger is required\")",
				"return nil, ErrNewPasswordRequired",
			},
			docTerms: []string{
				"`password_change`",
				"`first_login`",
				"`expired`",
				"`400`",
				"panic",
				"`ErrNewPasswordRequired`",
			},
		},
		{
			sourcePath: "security/department_selection.go",
			sourceTerms: []string{
				"ChallengeTypeDepartmentSelection = \"department_selection\"",
				"`json:\"departments\"`",
				"`json:\"meta,omitempty\"`",
				"func (*DepartmentSelectionChallengeProvider) Order() int   { return 500 }",
				"panic(\"security: DepartmentLoader is required\")",
				"panic(\"security: DepartmentSelector is required\")",
				"return nil, ErrDepartmentRequired",
			},
			docTerms: []string{
				"`department_selection`",
				"`departments`",
				"`meta`",
				"`500`",
				"panic",
				"`ErrDepartmentRequired`",
			},
		},
		{
			sourcePath: "security/signature.go",
			sourceTerms: []string{
				"SignatureAlgHmacSHA256 SignatureAlgorithm = \"HMAC-SHA256\"",
				"SignatureAlgHmacSHA512 SignatureAlgorithm = \"HMAC-SHA512\"",
				"SignatureAlgHmacSM3    SignatureAlgorithm = \"HMAC-SM3\"",
				"defaultSignatureTimestampTolerance = 5 * time.Minute",
				"nonceTTLBuffer                     = 1 * time.Minute",
				"nonceStore:         NewMemoryNonceStore()",
				"return nil, ErrSignatureSecretRequired",
				"return fmt.Appendf(nil, \"app_id=%s&method=%s&nonce=%s&path=%s&timestamp=%d\"",
				"if s.nonceStore == nil",
				"s.timestampTolerance+nonceTTLBuffer",
			},
			docTerms: []string{
				"`HMAC-SHA256`",
				"`HMAC-SHA512`",
				"`HMAC-SM3`",
				"`5m`",
				"tolerance + `1m`",
				"`MemoryNonceStore`",
				"`ErrSignatureSecretRequired`",
				"app_id=<appID>&method=<method>&nonce=<nonce>&path=<path>&timestamp=<timestamp>",
				"`WithNonceStore(nil)`",
				"signature payload",
			},
		},
		{
			sourcePath: "security/ip_whitelist.go",
			sourceTerms: []string{
				"whitelist = strings.TrimSpace(whitelist)",
				"validator.isEmpty = true",
				"validator.invalid = true",
				"if v.isEmpty",
				"if v.invalid",
			},
			docTerms: []string{
				"whitelist",
				"fail-closed",
			},
			englishDocTerms: []string{
				"allows all IPs",
				"denies all requests",
			},
			chineseDocTerms: []string{
				"允许全部 IP",
				"拒绝全部请求",
			},
		},
		{
			sourcePath: "internal/security/signature_authenticator.go",
			sourceTerms: []string{
				"const AuthTypeSignature = \"signature\"",
				"contextx.RequestMethod(ctx)",
				"contextx.RequestPath(ctx)",
				"details, ok := principal.Details.(*security.ExternalAppConfig)",
				"return security.ErrExternalAppDisabled",
				"requestIP := contextx.RequestIP(ctx)",
				"return security.ErrIPNotAllowed",
				"return security.ErrSignatureInvalid",
			},
			docTerms: []string{
				"`signature`",
				"HTTP method",
				"path",
				"`ExternalAppConfig.IPWhitelist`",
				"`ErrExternalAppDisabled`",
				"`ErrIPNotAllowed`",
				"`ErrSignatureInvalid`",
			},
		},
		{
			sourcePath: "security/api_errors.go",
			sourceTerms: []string{
				"ErrCodeUnauthenticated               = 1000",
				"ErrCodeUnsupportedAuthenticationType = 1001",
				"ErrCodeTokenExpired                  = 1002",
				"ErrCodeTokenInvalid                  = 1003",
				"ErrCodeTokenNotValidYet              = 1004",
				"ErrCodeTokenInvalidIssuer            = 1005",
				"ErrCodeTokenInvalidAudience          = 1006",
				"ErrCodePrincipalInvalid              = 1007",
				"ErrCodeCredentialsInvalid            = 1008",
				"ErrCodeAppIDRequired                 = 1009",
				"ErrCodeTimestampRequired             = 1010",
				"ErrCodeSignatureRequired             = 1011",
				"ErrCodeTimestampInvalid              = 1012",
				"ErrCodeSignatureExpired              = 1013",
				"ErrCodeExternalAppNotFound           = 1014",
				"ErrCodeExternalAppDisabled           = 1015",
				"ErrCodeIPNotAllowed                  = 1016",
				"ErrCodeSignatureInvalid              = 1017",
				"ErrCodeNonceRequired                 = 1018",
				"ErrCodeNonceInvalid                  = 1019",
				"ErrCodeNonceAlreadyUsed              = 1020",
				"ErrCodeAuthHeaderMissing             = 1021",
				"ErrCodeAuthHeaderInvalid             = 1022",
				"ErrCodeChallengeTokenInvalid  = 1031",
				"ErrCodeChallengeTypeInvalid   = 1033",
				"ErrCodeChallengeResolveFailed = 1034",
				"ErrCodeOTPCodeRequired        = 1035",
				"ErrCodeOTPCodeInvalid         = 1036",
				"ErrCodeNewPasswordRequired    = 1037",
				"ErrCodeDepartmentRequired     = 1038",
			},
			docTerms: []string{
				"`ErrCodeUnauthenticated`",
				"`ErrCodeUnsupportedAuthenticationType`",
				"`ErrCodeTokenExpired`",
				"`ErrCodeTokenInvalid`",
				"`ErrCodeTokenNotValidYet`",
				"`ErrCodeTokenInvalidIssuer`",
				"`ErrCodeTokenInvalidAudience`",
				"`ErrCodePrincipalInvalid`",
				"`ErrCodeCredentialsInvalid`",
				"`ErrCodeAppIDRequired`",
				"`ErrCodeTimestampRequired`",
				"`ErrCodeSignatureRequired`",
				"`ErrCodeTimestampInvalid`",
				"`ErrCodeSignatureExpired`",
				"`ErrCodeExternalAppNotFound`",
				"`ErrCodeExternalAppDisabled`",
				"`ErrCodeIPNotAllowed`",
				"`ErrCodeSignatureInvalid`",
				"`ErrCodeNonceRequired`",
				"`ErrCodeNonceInvalid`",
				"`ErrCodeNonceAlreadyUsed`",
				"`ErrCodeAuthHeaderMissing`",
				"`ErrCodeAuthHeaderInvalid`",
				"`ErrCodeChallengeTokenInvalid`",
				"`ErrCodeChallengeTypeInvalid`",
				"`ErrCodeChallengeResolveFailed`",
				"`ErrCodeOTPCodeRequired`",
				"`ErrCodeOTPCodeInvalid`",
				"`ErrCodeNewPasswordRequired`",
				"`ErrCodeDepartmentRequired`",
				"`400`",
				"`401`",
			},
		},
	}
}

func authzChecks() []check {
	return []check{
		{
			sourcePath: "security/cached_role_permission_loader.go",
			sourceTerms: []string{
				"eventTypeRolePermissionsChanged = \"vef.security.role_permissions.changed\"",
				"`json:\"roles\"`",
				"return bus.Publish(ctx, &RolePermissionsChangedEvent{Roles: roles})",
				"panic(fmt.Errorf(\"subscribe role_permissions.changed: %w\", err))",
				"return c.cache.Invalidate(ctx, evt.Roles...)",
			},
			docTerms: []string{
				"`vef.security.role_permissions.changed`",
				"`roles`",
				"empty `roles`",
			},
			englishDocTerms: []string{
				"panics",
			},
			chineseDocTerms: []string{
				"empty `roles`",
				"panic",
			},
		},
		{
			sourcePath: "internal/security/rbac_permission_checker.go",
			sourceTerms: []string{
				"if principal == nil || len(principal.Roles) == 0",
				"if c.loader == nil",
				"if _, exists := permissions[permissionToken]; exists",
			},
			docTerms: []string{
				"`RolePermissionsLoader`",
			},
			englishDocTerms: []string{
				"principal is nil",
				"no roles",
				"permission map contains",
			},
			chineseDocTerms: []string{
				"principal is nil",
				"no roles",
				"权限 map",
			},
		},
		{
			sourcePath: "internal/security/rbac_data_permission_resolver.go",
			sourceTerms: []string{
				"selectedScope security.DataScope",
				"maxPriority   = -1",
				"if priority := dataScope.Priority(); priority > maxPriority",
			},
			docTerms: []string{
				"`DataScope.Priority()`",
				"highest",
			},
		},
		{
			sourcePath: "security/login_event.go",
			sourceTerms: []string{
				"const eventTypeLogin = \"vef.security.login\"",
				"`json:\"authType\"`",
				"`json:\"userId\"`",
				"`json:\"username\"`",
				"`json:\"loginIp\"`",
				"`json:\"userAgent\"`",
				"`json:\"traceId\"`",
				"`json:\"isOk\"`",
				"`json:\"failReason\"`",
				"`json:\"errorCode\"`",
				"func SubscribeLoginEvent(",
			},
			docTerms: []string{
				"`vef.security.login`",
				"`authType`",
				"`userId`",
				"`username`",
				"`loginIp`",
				"`userAgent`",
				"`traceId`",
				"`isOk`",
				"`failReason`",
				"`errorCode`",
				"`SubscribeLoginEvent`",
			},
		},
		{
			sourcePath: "security/user_info.go",
			sourceTerms: []string{
				"GenderMale    Gender = \"male\"",
				"GenderFemale  Gender = \"female\"",
				"GenderUnknown Gender = \"unknown\"",
				"UserMenuTypeDirectory UserMenuType = \"directory\"",
				"UserMenuTypeMenu      UserMenuType = \"menu\"",
				"UserMenuTypeView      UserMenuType = \"view\"",
				"UserMenuTypeDashboard UserMenuType = \"dashboard\"",
				"UserMenuTypeReport    UserMenuType = \"report\"",
				"`json:\"permissionTokens\"`",
				"`json:\"metadata,omitempty\"`",
				"`json:\"children,omitempty\"`",
			},
			docTerms: []string{
				"`male`",
				"`female`",
				"`unknown`",
				"`directory`",
				"`menu`",
				"`view`",
				"`dashboard`",
				"`report`",
				"`permissionTokens`",
				"`metadata`",
				"`children`",
			},
		},
	}
}

func dataPermChecks() []check {
	return []check{
		{
			sourcePath: "security/data_scopes.go",
			sourceTerms: []string{
				"PrioritySelf = 10",
				"PriorityDepartment = 20",
				"PriorityDepartmentAndSub = 30",
				"PriorityOrganization = 40",
				"PriorityOrganizationAndSub = 50",
				"PriorityCustom = 60",
				"PriorityAll = 10000",
				"return \"all\"",
				"return \"self\"",
				"createdByColumn: cmp.Or(createdByColumn, orm.ColumnCreatedBy)",
				"cb.Equals(s.createdByColumn, principal.ID)",
			},
			docTerms: []string{
				"`PrioritySelf` (`10`)",
				"`PriorityDepartment` (`20`)",
				"`PriorityDepartmentAndSub` (`30`)",
				"`PriorityOrganization` (`40`)",
				"`PriorityOrganizationAndSub` (`50`)",
				"`PriorityCustom` (`60`)",
				"`PriorityAll` (`10000`)",
				"`all`",
				"`self`",
				"`created_by`",
				"principal ID",
			},
		},
		{
			sourcePath: "security/request_scoped_data_perm_applier.go",
			sourceTerms: []string{
				"if a.dataScope == nil",
				"return nil",
				"queryBuilder, ok := query.(orm.QueryBuilder)",
				"return ErrQueryNotQueryBuilder",
				"return ErrQueryModelNotSet",
				"if !a.dataScope.Supports(a.principal, table)",
				"failed to apply data scope %q",
			},
			docTerms: []string{
				"`RequestScopedDataPermApplier.Apply(...)`",
				"skip without error",
				"`orm.QueryBuilder`",
				"`ErrQueryNotQueryBuilder`",
				"`ErrQueryModelNotSet`",
				"`DataScope.Supports(...)`",
				"scope key",
			},
		},
		{
			sourcePath: "security/permission.go",
			sourceTerms: []string{
				"type PermissionChecker interface",
				"type RolePermissionsLoader interface",
				"type DataScope interface",
				"type DataPermissionResolver interface",
				"type DataPermissionApplier interface",
			},
			docTerms: []string{
				"`RolePermissionsLoader`",
				"`DataScope`",
				"`DataPermissionResolver`",
				"`DataPermissionApplier`",
			},
		},
		{
			sourcePath: "security/department_selection.go",
			sourceTerms: []string{
				"type DepartmentOption struct",
				"`json:\"id\"`",
				"`json:\"name\"`",
				"`json:\"departments\"`",
				"`json:\"meta,omitempty\"`",
				"type DepartmentLoader interface",
				"type DepartmentSelector interface",
				"type DepartmentSelectionChallengeProvider struct",
			},
			docTerms: []string{
				"`DepartmentOption`",
				"`DepartmentLoader`",
				"`DepartmentSelector`",
				"`DepartmentSelectionChallengeProvider`",
			},
		},
	}
}

func runChecks(sourceRoot string, docs []corpus, checks []check) []string {
	var failures []string
	for _, check := range checks {
		source := readCorpus(check.sourcePath, filepath.Join(sourceRoot, check.sourcePath))
		failures = append(failures, missingTerms(source, check.sourceTerms)...)
		for _, doc := range docs {
			failures = append(failures, missingTerms(doc, check.docTerms)...)
		}
		if len(check.englishDocTerms) > 0 && len(docs) > 0 {
			failures = append(failures, missingTerms(docs[0], check.englishDocTerms)...)
		}
		if len(check.chineseDocTerms) > 0 && len(docs) > 1 {
			failures = append(failures, missingTerms(docs[1], check.chineseDocTerms)...)
		}
	}

	return failures
}

func readCorpus(label, path string) corpus {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read %s at %s: %w", label, path, err))
	}

	return corpus{label: label, content: string(content)}
}

func missingTerms(c corpus, terms []string) []string {
	var failures []string
	for _, term := range terms {
		if !strings.Contains(c.content, term) {
			failures = append(failures, fmt.Sprintf("%s missing term: %s", c.label, term))
		}
	}

	return failures
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
