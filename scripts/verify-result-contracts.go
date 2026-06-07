package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	resultPackage     = "github.com/coldsmirk/vef-framework-go/result"
	resultFingerprint = "f91600ccb5960c2a405fb3ec5b2b84b38676c6488f4bf2dd45c8c22544b96892"
	resultTopLevel    = 48
	resultFields      = 6
	resultMethods     = 4
	resultEntries     = 58
)

type corpus struct {
	label   string
	content string
}

func main() {
	sourceDir := flag.String("source", "../vef-framework-go", "path to the VEF Framework Go source repository")
	outDir := flag.String("out", ".", "path to the VEF Framework Go docs repository")
	flag.Parse()

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)

	englishResultDocs := readCorpus("English result docs", filepath.Join(docsRoot, "docs/guide/result.md"))
	chineseResultDocs := readCorpus("Chinese result docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/guide/result.md"))
	englishErrorDocs := readCorpus("English error-handling docs", filepath.Join(docsRoot, "docs/guide/error-handling.md"))
	chineseErrorDocs := readCorpus("Chinese error-handling docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/guide/error-handling.md"))
	publicIndex := readCorpus("English public API index", filepath.Join(docsRoot, "docs/reference/public-api-index.md"))
	chinesePublicIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"))

	resultDir := filepath.Join(sourceRoot, "result")
	expectedSurface := packageSurface{
		consts: []string{
			"ErrCodeAccessDenied",
			"ErrCodeBadRequest",
			"ErrCodeDangerousSQL",
			"ErrCodeDefault",
			"ErrCodeForeignKeyViolation",
			"ErrCodeNotFound",
			"ErrCodeNotImplemented",
			"ErrCodeRecordAlreadyExists",
			"ErrCodeRecordNotFound",
			"ErrCodeRequestTimeout",
			"ErrCodeTooManyRequests",
			"ErrCodeUnknown",
			"ErrCodeUnsupportedMediaType",
			"ErrMessage",
			"ErrMessageAccessDenied",
			"ErrMessageDangerousSQL",
			"ErrMessageForeignKeyViolation",
			"ErrMessageNotFound",
			"ErrMessageRecordAlreadyExists",
			"ErrMessageRecordNotFound",
			"ErrMessageRequestTimeout",
			"ErrMessageTooManyRequests",
			"ErrMessageUnknown",
			"ErrMessageUnsupportedMediaType",
			"OkCode",
			"OkMessage",
		},
		funcs: []string{
			"AsErr",
			"Err",
			"ErrNotImplemented",
			"Errf",
			"IsRecordNotFound",
			"Ok",
			"WithCode",
			"WithMessage",
			"WithMessagef",
			"WithStatus",
		},
		types: []string{
			"ErrOption",
			"Error",
			"OkOption",
			"Result",
		},
		vars: []string{
			"ErrAccessDenied",
			"ErrDangerousSQL",
			"ErrForeignKeyViolation",
			"ErrRecordAlreadyExists",
			"ErrRecordNotFound",
			"ErrRequestTimeout",
			"ErrTooManyRequests",
			"ErrUnknown",
		},
		methods: []string{
			"Error.Error",
			"Error.Is",
			"Result.IsOk",
			"Result.Response",
		},
		fields: []string{
			"Error.Code",
			"Error.Message",
			"Error.Status",
			"Result.Code",
			"Result.Data",
			"Result.Message",
		},
	}

	var failures []string
	exported := exportedPackageSurface(resultDir)
	failures = append(failures, compareNames("result const", exported.consts, expectedSurface.consts)...)
	failures = append(failures, compareNames("result func", exported.funcs, expectedSurface.funcs)...)
	failures = append(failures, compareNames("result type", exported.types, expectedSurface.types)...)
	failures = append(failures, compareNames("result var", exported.vars, expectedSurface.vars)...)
	failures = append(failures, compareNames("result method", exported.methods, expectedSurface.methods)...)
	failures = append(failures, compareNames("result exported field", exported.fields, expectedSurface.fields)...)
	failures = append(failures, compareSurfaceCounts(exported)...)

	for _, doc := range []corpus{publicIndex, chinesePublicIndex} {
		failures = append(failures, missingTerms(doc, publicIndexTerms())...)
	}
	for _, doc := range []corpus{englishResultDocs, chineseResultDocs} {
		failures = append(failures, missingTerms(doc, publicDocSurfaceTerms())...)
		failures = append(failures, verifyResultReferences(doc, expectedSurface)...)
		failures = append(failures, verifySurfaceMentioned(doc, exported)...)
	}
	for _, doc := range []corpus{englishErrorDocs, chineseErrorDocs} {
		failures = append(failures, verifyResultReferences(doc, expectedSurface)...)
	}

	failures = append(failures, missingTerms(englishResultDocs, []string{
		"48 top-level exported symbols",
		"6 exported fields",
		"4 exported\nmethods",
		"fingerprint is\n`" + resultFingerprint + "`",
		"Only this type has JSON field tags for the public wire shape",
		"`Result.Code` | Business result code, serialized as `code`",
		"`Result.Message` | Human-readable or i18n-resolved message, serialized as `message`",
		"`Result.Data` | Optional response payload, serialized as `data`; `nil` data is preserved as JSON `null`",
		"`Result.Response(ctx, status...)` | Sends the result as JSON; HTTP status defaults to `200 OK`, or the first supplied status value.",
		"`Error.Code` | Business error code used in the response envelope and in `errors.Is` comparisons.",
		"`Error.Status` | HTTP status used by the app error handler when converting the error into a `Result`.",
		"`Error.Is(target)` | Matches another `result.Error` by `Code` only",
		"`result.Error` intentionally has no JSON tags",
		"Do not serialize `result.Error`\ndirectly",
		"Passing more than one data argument, or passing\ndata after an option, panics",
		"There is no error-message option",
		"`Errf` requires at least one format argument",
		"Two\n`result.Error` values match when their `Code` values are equal",
		"Use `result.AsErr(err)` when you need to read `Code`, `Message`, or `Status`",
		"`result.ErrCodeBadRequest` | Bad request; standalone constant",
		"`ErrCodeNotFound`, `ErrCodeUnsupportedMediaType`, and `ErrCodeBadRequest` do not\nhave predefined `result.Error` values",
		"`result.OkMessage` | `\"ok\"` | default `Ok(...)` message",
		"`result.ErrMessageDangerousSQL` | `\"dangerous_sql\"` | `ErrDangerousSQL`",
		"The database\nand SQL-class business failures intentionally keep HTTP `200 OK`",
		"`result.ErrAccessDenied` | `result.ErrCodeAccessDenied` (`1100`) | `403` | `result.ErrMessageAccessDenied`",
		"`result.ErrDangerousSQL` | `result.ErrCodeDangerousSQL` (`1600`) | `200` | `result.ErrMessageDangerousSQL`",
	})...)
	failures = append(failures, missingTerms(chineseResultDocs, []string{
		"48 个\ntop-level exported symbols",
		"6 个 exported fields",
		"4 个 exported methods",
		"fingerprint 是\n`" + resultFingerprint + "`",
		"只有这个类型带有 public wire shape 的 JSON field tags",
		"`Result.Code` | 业务结果码，序列化为 `code`",
		"`Result.Message` | 面向用户或经 i18n 解析后的消息，序列化为 `message`",
		"`Result.Data` | 可选响应载荷，序列化为 `data`；`nil` data 会保留为 JSON `null`",
		"`Result.Response(ctx, status...)` | 以 JSON 发送结果；HTTP status 默认 `200 OK`，如果传入 status 则使用第一个值。",
		"`Error.Code` | 业务错误码，会进入响应信封，也用于 `errors.Is` 比较。",
		"`Error.Status` | 应用错误处理器把错误转换成 `Result` 时使用的 HTTP status。",
		"`Error.Is(target)` | 只按 `Code` 匹配另一个 `result.Error`",
		"`result.Error` 刻意没有 JSON tags",
		"不要直接序列化 `result.Error`",
		"传入多个 data 参数，或把 data 放在 option 之后，都会 panic",
		"错误消息没有 option 形式",
		"`Errf` 至少需要一个 format arg",
		"两个\n`result.Error` 只要 `Code` 相同就会匹配",
		"需要从 error chain 读取 `Code`、`Message` 或 `Status` 时，使用\n`result.AsErr(err)`",
		"`result.ErrCodeBadRequest` | 错误请求；standalone constant",
		"`ErrCodeNotFound`、`ErrCodeUnsupportedMediaType` 和 `ErrCodeBadRequest` 在这个包里没有预置\n`result.Error` 值",
		"`result.OkMessage` | `\"ok\"` | 默认 `Ok(...)` 消息",
		"`result.ErrMessageDangerousSQL` | `\"dangerous_sql\"` | `ErrDangerousSQL`",
		"数据库和 SQL 类业务失败刻意保持\nHTTP `200 OK`",
		"`result.ErrAccessDenied` | `result.ErrCodeAccessDenied`（`1100`） | `403` | `result.ErrMessageAccessDenied`",
		"`result.ErrDangerousSQL` | `result.ErrCodeDangerousSQL`（`1600`） | `200` | `result.ErrMessageDangerousSQL`",
	})...)

	failures = append(failures, missingTerms(englishErrorDocs, []string{
		"`result.Error` is not serialized directly as the client response",
		"`result.Ok(...)` accepts at most one data argument",
		"The optional message string must be the first `Err(...)` argument",
		"There is no message option for `result.Error`",
		"`result.Error` implements `errors.Is` by comparing `Code` only",
		"`result.ErrCodeBadRequest`, `result.ErrCodeNotFound`, and\n`result.ErrCodeUnsupportedMediaType` are exported building-block constants",
	})...)
	failures = append(failures, missingTerms(chineseErrorDocs, []string{
		"`result.Error` 不会被直接序列化为客户端响应",
		"`result.Ok(...)` 最多接受一个 data 参数",
		"可选 message string 必须是 `Err(...)` 的第一个参数",
		"`result.Error` 没有 message option",
		"`result.Error` 通过只比较 `Code` 实现 `errors.Is`",
		"`result.ErrCodeBadRequest`、`result.ErrCodeNotFound` 和\n`result.ErrCodeUnsupportedMediaType` 是 exported building-block constants",
	})...)

	failures = append(failures, verifyAuditArtifacts(sourceRoot, docsRoot)...)
	failures = append(failures, verifyAppErrorConversion(sourceRoot)...)
	failures = append(failures, verifyMessageKeysInLocales(sourceRoot)...)

	sourceChecks := []struct {
		path  string
		terms []string
	}{
		{
			path: "result/result.go",
			terms: []string{
				"type Result struct",
				"Code    int    `json:\"code\"`",
				"Message string `json:\"message\"`",
				"Data    any    `json:\"data\"`",
				"func (r Result) Response(ctx fiber.Ctx, status ...int) error",
				"statusCode := fiber.StatusOK",
				"statusCode = status[0]",
				"return ctx.Status(statusCode).JSON(r)",
				"func (r Result) IsOk() bool",
				"return r.Code == OkCode",
				"func Ok(dataOrOptions ...any) Result",
				"panic(\"result.Ok: data must come before options\")",
				"panic(\"result.Ok: only one data argument is allowed\")",
				"Message: i18n.T(OkMessage)",
				"for _, opt := range options",
				"opt(&r)",
			},
		},
		{
			path: "result/error.go",
			terms: []string{
				"type Error struct",
				"Code    int",
				"Message string",
				"Status  int",
				"func (e Error) Error() string",
				"return e.Message",
				"func (e Error) Is(target error) bool",
				"other, ok := target.(Error)",
				"return e.Code == other.Code",
				"func Err(messageOrOptions ...any) Error",
				"panic(\"result.Err: message string must be the first argument\")",
				"panic(fmt.Sprintf(\"result.Err: invalid argument type %T at position %d\", v, i))",
				"message = i18n.T(ErrMessage)",
				"Status:  fiber.StatusOK",
				"func Errf(format string, args ...any) Error",
				"panic(\"result.Errf: at least one format argument is required\")",
				"panic(\"result.Errf: format arguments must come before options\")",
				"Message: fmt.Sprintf(format, formatArgs...)",
				"func AsErr(err error) (Error, bool)",
				"return errors.AsType[Error](err)",
				"func IsRecordNotFound(err error) bool",
				"return errors.Is(err, ErrRecordNotFound)",
			},
		},
		{
			path: "result/errors.go",
			terms: []string{
				"ErrAccessDenied = Err(",
				"WithCode(ErrCodeAccessDenied)",
				"WithStatus(fiber.StatusForbidden)",
				"ErrTooManyRequests = Err(",
				"WithCode(ErrCodeTooManyRequests)",
				"WithStatus(fiber.StatusTooManyRequests)",
				"ErrRequestTimeout = Err(",
				"WithCode(ErrCodeRequestTimeout)",
				"WithStatus(fiber.StatusRequestTimeout)",
				"ErrUnknown = Err(",
				"WithCode(ErrCodeUnknown)",
				"WithStatus(fiber.StatusInternalServerError)",
				"ErrRecordNotFound = Err(",
				"WithCode(ErrCodeRecordNotFound)",
				"ErrRecordAlreadyExists = Err(",
				"WithCode(ErrCodeRecordAlreadyExists)",
				"ErrForeignKeyViolation = Err(",
				"WithCode(ErrCodeForeignKeyViolation)",
				"ErrDangerousSQL = Err(",
				"WithCode(ErrCodeDangerousSQL)",
				"func ErrNotImplemented(message string) Error",
				"WithCode(ErrCodeNotImplemented)",
				"WithStatus(fiber.StatusNotImplemented)",
			},
		},
		{
			path: "result/constants.go",
			terms: []string{
				"OkMessage  = \"ok\"",
				"ErrMessage = \"error\"",
				"ErrMessageUnknown              = \"unknown_error\"",
				"ErrMessageNotFound             = \"not_found\"",
				"ErrMessageTooManyRequests      = \"too_many_requests\"",
				"ErrMessageAccessDenied         = \"access_denied\"",
				"ErrMessageUnsupportedMediaType = \"unsupported_media_type\"",
				"ErrMessageRequestTimeout       = \"request_timeout\"",
				"ErrMessageRecordNotFound      = \"record_not_found\"",
				"ErrMessageRecordAlreadyExists = \"record_already_exists\"",
				"ErrMessageForeignKeyViolation = \"foreign_key_violation\"",
				"ErrMessageDangerousSQL        = \"dangerous_sql\"",
				"OkCode = 0",
				"ErrCodeAccessDenied = 1100",
				"ErrCodeNotFound = 1200",
				"ErrCodeUnsupportedMediaType = 1300",
				"ErrCodeBadRequest      = 1400",
				"ErrCodeTooManyRequests = 1401",
				"ErrCodeRequestTimeout  = 1402",
				"ErrCodeNotImplemented = 1500",
				"ErrCodeDangerousSQL = 1600",
				"ErrCodeUnknown = 1900",
				"ErrCodeDefault             = 2000",
				"ErrCodeRecordNotFound      = 2001",
				"ErrCodeRecordAlreadyExists = 2002",
				"ErrCodeForeignKeyViolation = 2003",
			},
		},
		{
			path: "result/option.go",
			terms: []string{
				"type ErrOption func(*Error)",
				"func WithCode(code int) ErrOption",
				"return func(e *Error) { e.Code = code }",
				"func WithStatus(status int) ErrOption",
				"return func(e *Error) { e.Status = status }",
				"type OkOption func(*Result)",
				"func WithMessage(message string) OkOption",
				"return func(r *Result) { r.Message = message }",
				"func WithMessagef(format string, args ...any) OkOption",
				"return func(r *Result) { r.Message = fmt.Sprintf(format, args...) }",
			},
		},
		{
			path: "internal/app/error.go",
			terms: []string{
				"fiberErrorMappings = map[int]fiberErrorMapping",
				"fiber.StatusNotFound",
				"code:    result.ErrCodeNotFound",
				"message: result.ErrMessageNotFound",
				"fiber.StatusUnsupportedMediaType",
				"code:    result.ErrCodeUnsupportedMediaType",
				"message: result.ErrMessageUnsupportedMediaType",
				"if resultErr, ok := result.AsErr(err); ok",
				"return responseError(resultErr, ctx)",
				"func responseError(e result.Error, ctx fiber.Ctx) error",
				"return result.Result{",
				"Code:    e.Code",
				"Message: e.Message",
				"}.Response(ctx, e.Status)",
			},
		},
		{
			path: "i18n/locales/en.json",
			terms: []string{
				"\"ok\": \"Success\"",
				"\"error\": \"Error\"",
				"\"record_not_found\": \"Record not found\"",
				"\"record_already_exists\": \"Record already exists\"",
				"\"foreign_key_violation\": \"Cannot delete or update a record with existing references\"",
				"\"unknown_error\": \"An unexpected error occurred\"",
				"\"not_found\": \"Resource not found\"",
				"\"too_many_requests\": \"Too many requests\"",
				"\"access_denied\": \"Access denied\"",
				"\"unsupported_media_type\": \"Unsupported media type. Only JSON and file data are supported\"",
				"\"request_timeout\": \"Request timeout\"",
				"\"dangerous_sql\": \"Dangerous SQL detected, execution blocked\"",
			},
		},
		{
			path: "i18n/locales/zh-CN.json",
			terms: []string{
				"\"ok\": \"成功\"",
				"\"error\": \"失败\"",
				"\"record_not_found\": \"记录不存在\"",
				"\"record_already_exists\": \"记录已存在\"",
				"\"foreign_key_violation\": \"数据存在关联，无法删除或更新\"",
				"\"unknown_error\": \"出小差了\"",
				"\"not_found\": \"迷路了\"",
				"\"too_many_requests\": \"请求过于频繁\"",
				"\"access_denied\": \"无权限\"",
				"\"unsupported_media_type\": \"请求仅支持JSON或文件数据\"",
				"\"request_timeout\": \"请求超时\"",
				"\"dangerous_sql\": \"检测到危险 SQL 操作, 执行已阻止\"",
			},
		},
		{
			path: "result/result_test.go",
			terms: []string{
				"TestOk",
				"WithMessageOptionOnly",
				"TestOkPanicCases",
				"MultipleDataArguments",
				"DataAfterOption",
				"TestResultResponse",
				"DefaultStatus",
				"CustomStatus",
				"TestResultIsOk",
				"TestResultJSONSerialization",
				"JSONFieldNames",
				"TestOkWithVariousDataTypes",
				"NilData",
				"TestResultResponseIntegration",
			},
		},
		{
			path: "result/error_test.go",
			terms: []string{
				"TestErr",
				"MessageWithCodeAndStatus",
				"TestErrPanicCases",
				"MessageNotFirstArgument",
				"InvalidArgumentType",
				"TestErrf",
				"FormattingWithCodeAndStatus",
				"TestErrfPanicCases",
				"NoFormatArguments",
				"OptionBeforeFormatArgs",
				"MixedFormatArgsAndOptions",
				"TestAsErr",
				"WrappedError",
				"TestIsRecordNotFound",
				"WrappedRecordNotFoundError",
				"TestErrorStructure",
				"TestPredefinedErrors",
				"TestErrorIs",
				"MatchingCodeMatchesIgnoringMessage",
				"ErrfMatchesPredefinedSentinelByCode",
			},
		},
	}
	for _, check := range sourceChecks {
		source := readCorpus(check.path, filepath.Join(sourceRoot, check.path))
		failures = append(failures, missingTerms(source, check.terms)...)
	}

	failures = append(failures, runPackageTests(sourceRoot)...)

	sort.Strings(failures)
	if len(failures) > 0 {
		panic(fmt.Errorf("result contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	topLevelPublic := len(exported.consts) + len(exported.funcs) + len(exported.types) + len(exported.vars)
	fmt.Printf("Result contract docs verified: %d top-level public symbols, %d public methods, %d public fields, %d source files, 2 doc mirrors\n",
		topLevelPublic, len(exported.methods), len(exported.fields), len(sourceChecks))
}

type packageSurface struct {
	consts  []string
	funcs   []string
	types   []string
	vars    []string
	methods []string
	fields  []string
}

func exportedPackageSurface(dir string) packageSurface {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		name := info.Name()

		return !info.IsDir() && strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
	}, 0)
	if err != nil {
		panic(fmt.Errorf("failed to parse result package: %w", err))
	}

	consts := make(map[string]bool)
	funcs := make(map[string]bool)
	typesMap := make(map[string]bool)
	vars := make(map[string]bool)
	methods := make(map[string]bool)
	fields := make(map[string]bool)

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				switch d := decl.(type) {
				case *ast.FuncDecl:
					if d.Recv == nil {
						if d.Name.IsExported() {
							funcs[d.Name.Name] = true
						}

						continue
					}
					if d.Name.IsExported() {
						recv := receiverBaseName(d.Recv.List[0].Type)
						if ast.IsExported(recv) {
							methods[recv+"."+d.Name.Name] = true
						}
					}
				case *ast.GenDecl:
					for _, spec := range d.Specs {
						switch s := spec.(type) {
						case *ast.TypeSpec:
							if s.Name.IsExported() {
								typesMap[s.Name.Name] = true
								if structType, ok := s.Type.(*ast.StructType); ok {
									for _, field := range structType.Fields.List {
										for _, name := range field.Names {
											if name.IsExported() {
												fields[s.Name.Name+"."+name.Name] = true
											}
										}
									}
								}
								if iface, ok := s.Type.(*ast.InterfaceType); ok {
									for _, field := range iface.Methods.List {
										for _, name := range field.Names {
											if name.IsExported() {
												methods[s.Name.Name+"."+name.Name] = true
											}
										}
									}
								}
							}
						case *ast.ValueSpec:
							for _, name := range s.Names {
								if !name.IsExported() {
									continue
								}
								if d.Tok == token.CONST {
									consts[name.Name] = true
								} else {
									vars[name.Name] = true
								}
							}
						}
					}
				}
			}
		}
	}

	return packageSurface{
		consts:  sortedKeys(consts),
		funcs:   sortedKeys(funcs),
		types:   sortedKeys(typesMap),
		vars:    sortedKeys(vars),
		methods: sortedKeys(methods),
		fields:  sortedKeys(fields),
	}
}

func compareSurfaceCounts(surface packageSurface) []string {
	topLevel := len(surface.consts) + len(surface.funcs) + len(surface.types) + len(surface.vars)
	var failures []string
	if topLevel != resultTopLevel {
		failures = append(failures, fmt.Sprintf("result top-level exported symbol count mismatch: got %d want %d", topLevel, resultTopLevel))
	}
	if len(surface.fields) != resultFields {
		failures = append(failures, fmt.Sprintf("result exported field count mismatch: got %d want %d", len(surface.fields), resultFields))
	}
	if len(surface.methods) != resultMethods {
		failures = append(failures, fmt.Sprintf("result exported method count mismatch: got %d want %d", len(surface.methods), resultMethods))
	}
	if topLevel+len(surface.fields)+len(surface.methods) != resultEntries {
		failures = append(failures, fmt.Sprintf("result exported entry count mismatch: got %d want %d", topLevel+len(surface.fields)+len(surface.methods), resultEntries))
	}

	return failures
}

func publicDocSurfaceTerms() []string {
	return []string{
		"`result.Result`",
		"`Result.Code`",
		"`Result.Message`",
		"`Result.Data`",
		"`Result.Response(ctx, status...)`",
		"`Result.IsOk()`",
		"`result.Ok(dataOrOptions...)`",
		"`result.OkOption`",
		"`result.WithMessage(message)`",
		"`result.WithMessagef(format, args...)`",
		"`result.Error`",
		"`Error.Code`",
		"`Error.Message`",
		"`Error.Status`",
		"`Error.Error()`",
		"`Error.Is(target)`",
		"`result.Err(messageOrOptions...)`",
		"`result.Errf(format, args...)`",
		"`result.ErrOption`",
		"`result.WithCode(code)`",
		"`result.WithStatus(status)`",
		"`result.AsErr(err)`",
		"`result.IsRecordNotFound(err)`",
		"`result.ErrNotImplemented(message)`",
		"`result.OkCode` / `result.OkMessage`",
		"`ErrCode*` family",
		"`ErrMessage*` family",
		"`result.ErrAccessDenied`, `result.ErrTooManyRequests`, `result.ErrRequestTimeout`, `result.ErrUnknown`, `result.ErrRecordNotFound`, `result.ErrRecordAlreadyExists`, `result.ErrForeignKeyViolation`, `result.ErrDangerousSQL`",
	}
}

func publicIndexTerms() []string {
	return []string{
		"## " + resultPackage,
		"FUNC AsErr : func(err error) (github.com/coldsmirk/vef-framework-go/result.Error, bool)",
		"FUNC Err : func(messageOrOptions ...any) github.com/coldsmirk/vef-framework-go/result.Error",
		"VAR ErrAccessDenied : github.com/coldsmirk/vef-framework-go/result.Error",
		"CONST ErrCodeAccessDenied : untyped int = 1100",
		"CONST ErrCodeBadRequest : untyped int = 1400",
		"CONST ErrCodeDangerousSQL : untyped int = 1600",
		"CONST ErrCodeDefault : untyped int = 2000",
		"CONST ErrCodeForeignKeyViolation : untyped int = 2003",
		"CONST ErrCodeNotFound : untyped int = 1200",
		"CONST ErrCodeNotImplemented : untyped int = 1500",
		"CONST ErrCodeRecordAlreadyExists : untyped int = 2002",
		"CONST ErrCodeRecordNotFound : untyped int = 2001",
		"CONST ErrCodeRequestTimeout : untyped int = 1402",
		"CONST ErrCodeTooManyRequests : untyped int = 1401",
		"CONST ErrCodeUnknown : untyped int = 1900",
		"CONST ErrCodeUnsupportedMediaType : untyped int = 1300",
		"VAR ErrDangerousSQL : github.com/coldsmirk/vef-framework-go/result.Error",
		"VAR ErrForeignKeyViolation : github.com/coldsmirk/vef-framework-go/result.Error",
		"CONST ErrMessage : untyped string = \"error\"",
		"CONST ErrMessageAccessDenied : untyped string = \"access_denied\"",
		"CONST ErrMessageDangerousSQL : untyped string = \"dangerous_sql\"",
		"CONST ErrMessageForeignKeyViolation : untyped string = \"foreign_key_violation\"",
		"CONST ErrMessageNotFound : untyped string = \"not_found\"",
		"CONST ErrMessageRecordAlreadyExists : untyped string = \"record_already_exists\"",
		"CONST ErrMessageRecordNotFound : untyped string = \"record_not_found\"",
		"CONST ErrMessageRequestTimeout : untyped string = \"request_timeout\"",
		"CONST ErrMessageTooManyRequests : untyped string = \"too_many_requests\"",
		"CONST ErrMessageUnknown : untyped string = \"unknown_error\"",
		"CONST ErrMessageUnsupportedMediaType : untyped string = \"unsupported_media_type\"",
		"FUNC ErrNotImplemented : func(message string) github.com/coldsmirk/vef-framework-go/result.Error",
		"TYPE ErrOption : github.com/coldsmirk/vef-framework-go/result.ErrOption",
		"VAR ErrRecordAlreadyExists : github.com/coldsmirk/vef-framework-go/result.Error",
		"VAR ErrRecordNotFound : github.com/coldsmirk/vef-framework-go/result.Error",
		"VAR ErrRequestTimeout : github.com/coldsmirk/vef-framework-go/result.Error",
		"VAR ErrTooManyRequests : github.com/coldsmirk/vef-framework-go/result.Error",
		"VAR ErrUnknown : github.com/coldsmirk/vef-framework-go/result.Error",
		"FUNC Errf : func(format string, args ...any) github.com/coldsmirk/vef-framework-go/result.Error",
		"TYPE Error : github.com/coldsmirk/vef-framework-go/result.Error",
		"FIELD Code : int [field_order=1 tag=\"\"]",
		"FIELD Message : string [field_order=2 tag=\"\"]",
		"FIELD Status : int [field_order=3 tag=\"\"]",
		"METHOD Error : func() string",
		"METHOD Is : func(target error) bool",
		"FUNC IsRecordNotFound : func(err error) bool",
		"FUNC Ok : func(dataOrOptions ...any) github.com/coldsmirk/vef-framework-go/result.Result",
		"CONST OkCode : untyped int = 0",
		"CONST OkMessage : untyped string = \"ok\"",
		"TYPE OkOption : github.com/coldsmirk/vef-framework-go/result.OkOption",
		"TYPE Result : github.com/coldsmirk/vef-framework-go/result.Result",
		"FIELD Code : int [field_order=1 tag=\"json:\\\"code\\\"\"]",
		"FIELD Message : string [field_order=2 tag=\"json:\\\"message\\\"\"]",
		"FIELD Data : any [field_order=3 tag=\"json:\\\"data\\\"\"]",
		"METHOD IsOk : func() bool",
		"METHOD Response : func(ctx github.com/gofiber/fiber/v3.Ctx, status ...int) error",
		"FUNC WithCode : func(code int) github.com/coldsmirk/vef-framework-go/result.ErrOption",
		"FUNC WithMessage : func(message string) github.com/coldsmirk/vef-framework-go/result.OkOption",
		"FUNC WithMessagef : func(format string, args ...any) github.com/coldsmirk/vef-framework-go/result.OkOption",
		"FUNC WithStatus : func(status int) github.com/coldsmirk/vef-framework-go/result.ErrOption",
	}
}

type contractLedger struct {
	PackageReviews []contractPackageReview `json:"package_reviews"`
}

type contractPackageReview struct {
	Package         string              `json:"package"`
	ReviewedSurface publicSurfaceReview `json:"reviewed_surface"`
}

type publicSurfaceReview struct {
	TopLevel    int    `json:"top_level"`
	Fields      int    `json:"fields"`
	Methods     int    `json:"methods"`
	EntryCount  int    `json:"entry_count"`
	Fingerprint string `json:"fingerprint"`
}

type apiManifest struct {
	Packages []apiManifestPackage `json:"packages"`
}

type apiManifestPackage struct {
	Package     string `json:"package"`
	TopLevel    int    `json:"top_level"`
	Fields      int    `json:"fields"`
	Methods     int    `json:"methods"`
	Fingerprint string `json:"fingerprint"`
}

func verifyAuditArtifacts(sourceRoot, docsRoot string) []string {
	var failures []string
	failures = append(failures, verifyContractLedger(filepath.Join(docsRoot, "scripts/api-contract-ledger.json"))...)
	failures = append(failures, verifyManifest(filepath.Join(docsRoot, "scripts/api-audit-manifest.json"))...)
	failures = append(failures, verifyLiveAuditSurface(sourceRoot, docsRoot)...)

	return failures
}

func verifyLiveAuditSurface(sourceRoot, docsRoot string) []string {
	cmd := exec.Command("go", "run", filepath.Join(docsRoot, "scripts/verify-api-audit.go"), "-source", sourceRoot, "-print-current")
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("live API audit surface print failed: %v\n%s", err, strings.TrimSpace(string(output)))}
	}

	var entries []apiManifestPackage
	wrapped := append([]byte("["), output...)
	wrapped = append(wrapped, ']')
	if err := json.Unmarshal(wrapped, &entries); err != nil {
		return []string{fmt.Sprintf("live API audit surface parse failed: %v", err)}
	}

	for _, entry := range entries {
		if entry.Package != resultPackage {
			continue
		}

		return verifyReviewedSurface("live source API inventory", publicSurfaceReview{
			TopLevel:    entry.TopLevel,
			Fields:      entry.Fields,
			Methods:     entry.Methods,
			EntryCount:  entry.TopLevel + entry.Fields + entry.Methods,
			Fingerprint: entry.Fingerprint,
		})
	}

	return []string{"live API audit surface missing result package"}
}

func verifyContractLedger(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{fmt.Sprintf("failed to read api-contract-ledger.json: %v", err)}
	}

	var ledger contractLedger
	if err := json.Unmarshal(data, &ledger); err != nil {
		return []string{fmt.Sprintf("failed to parse api-contract-ledger.json: %v", err)}
	}

	for _, review := range ledger.PackageReviews {
		if review.Package != resultPackage {
			continue
		}

		return verifyReviewedSurface("api-contract-ledger.json", review.ReviewedSurface)
	}

	return []string{"api-contract-ledger.json missing result package review"}
}

func verifyManifest(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{fmt.Sprintf("failed to read api-audit-manifest.json: %v", err)}
	}

	var manifest apiManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return []string{fmt.Sprintf("failed to parse api-audit-manifest.json: %v", err)}
	}

	for _, entry := range manifest.Packages {
		if entry.Package != resultPackage {
			continue
		}

		return verifyReviewedSurface("api-audit-manifest.json", publicSurfaceReview{
			TopLevel:    entry.TopLevel,
			Fields:      entry.Fields,
			Methods:     entry.Methods,
			EntryCount:  entry.TopLevel + entry.Fields + entry.Methods,
			Fingerprint: entry.Fingerprint,
		})
	}

	return []string{"api-audit-manifest.json missing result package entry"}
}

func verifyReviewedSurface(label string, surface publicSurfaceReview) []string {
	var failures []string
	if surface.TopLevel != resultTopLevel {
		failures = append(failures, fmt.Sprintf("%s result top_level mismatch: got %d want %d", label, surface.TopLevel, resultTopLevel))
	}
	if surface.Fields != resultFields {
		failures = append(failures, fmt.Sprintf("%s result fields mismatch: got %d want %d", label, surface.Fields, resultFields))
	}
	if surface.Methods != resultMethods {
		failures = append(failures, fmt.Sprintf("%s result methods mismatch: got %d want %d", label, surface.Methods, resultMethods))
	}
	if surface.EntryCount != resultEntries {
		failures = append(failures, fmt.Sprintf("%s result entry_count mismatch: got %d want %d", label, surface.EntryCount, resultEntries))
	}
	if surface.Fingerprint != resultFingerprint {
		failures = append(failures, fmt.Sprintf("%s result fingerprint mismatch: got %s want %s", label, surface.Fingerprint, resultFingerprint))
	}

	return failures
}

func verifyResultReferences(doc corpus, surface packageSurface) []string {
	allowed := make(map[string]bool)
	for _, group := range [][]string{surface.consts, surface.funcs, surface.types, surface.vars} {
		for _, name := range group {
			allowed[name] = true
		}
	}

	seen := make(map[string]bool)
	re := regexp.MustCompile(`result\.([A-Z][A-Za-z0-9_]*)`)
	for _, match := range re.FindAllStringSubmatch(doc.content, -1) {
		seen[match[1]] = true
	}

	var failures []string
	for name := range seen {
		if !allowed[name] {
			failures = append(failures, fmt.Sprintf("%s references unknown result public symbol result.%s", doc.label, name))
		}
	}

	return failures
}

func verifySurfaceMentioned(doc corpus, surface packageSurface) []string {
	var failures []string
	for _, name := range topLevelNames(surface) {
		if !containsTerm(doc.content, "result."+name) {
			failures = append(failures, fmt.Sprintf("%s missing exported result symbol result.%s", doc.label, name))
		}
	}

	for _, field := range surface.fields {
		if !containsTerm(doc.content, field) {
			failures = append(failures, fmt.Sprintf("%s missing exported result field %s", doc.label, field))
		}
	}

	for _, method := range surface.methods {
		if !containsTerm(doc.content, method) {
			failures = append(failures, fmt.Sprintf("%s missing exported result method %s", doc.label, method))
		}
	}

	return failures
}

func topLevelNames(surface packageSurface) []string {
	names := append([]string{}, surface.consts...)
	names = append(names, surface.funcs...)
	names = append(names, surface.types...)
	names = append(names, surface.vars...)
	sort.Strings(names)

	return names
}

func verifyAppErrorConversion(sourceRoot string) []string {
	source := readCorpus("app error conversion source", filepath.Join(sourceRoot, "internal/app/error.go"))

	return missingTerms(source, []string{
		"if resultErr, ok := result.AsErr(err); ok",
		"return responseError(resultErr, ctx)",
		"func responseError(e result.Error, ctx fiber.Ctx) error",
		"return result.Result{",
		"Code:    e.Code",
		"Message: e.Message",
		"}.Response(ctx, e.Status)",
	})
}

func verifyMessageKeysInLocales(sourceRoot string) []string {
	keys, failures := resultMessageKeys(filepath.Join(sourceRoot, "result/constants.go"))
	if len(failures) > 0 {
		return failures
	}

	var allFailures []string
	for _, path := range []string{
		filepath.Join(sourceRoot, "i18n/locales/en.json"),
		filepath.Join(sourceRoot, "i18n/locales/zh-CN.json"),
	} {
		data, err := os.ReadFile(path)
		if err != nil {
			allFailures = append(allFailures, fmt.Sprintf("failed to read locale %s: %v", path, err))
			continue
		}
		var locale map[string]any
		if err := json.Unmarshal(data, &locale); err != nil {
			allFailures = append(allFailures, fmt.Sprintf("failed to parse locale %s: %v", path, err))
			continue
		}
		for name, key := range keys {
			if _, ok := locale[key]; !ok {
				allFailures = append(allFailures, fmt.Sprintf("%s missing result message key %s=%q", path, name, key))
			}
		}
	}

	return allFailures
}

func resultMessageKeys(path string) (map[string]string, []string) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, []string{fmt.Sprintf("failed to parse result constants: %v", err)}
	}

	keys := make(map[string]string)
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, name := range valueSpec.Names {
				if name.Name != "OkMessage" && name.Name != "ErrMessage" && !strings.HasPrefix(name.Name, "ErrMessage") {
					continue
				}
				if i >= len(valueSpec.Values) {
					return nil, []string{fmt.Sprintf("result message key %s has no explicit value", name.Name)}
				}
				lit, ok := valueSpec.Values[i].(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					return nil, []string{fmt.Sprintf("result message key %s is not a string literal", name.Name)}
				}
				value, err := strconv.Unquote(lit.Value)
				if err != nil {
					return nil, []string{fmt.Sprintf("failed to unquote result message key %s: %v", name.Name, err)}
				}
				keys[name.Name] = value
			}
		}
	}

	if len(keys) != 12 {
		return nil, []string{fmt.Sprintf("result message key count mismatch: got %d want 12", len(keys))}
	}

	return keys, nil
}

func runPackageTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./result")
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("go test ./result failed: %v\n%s", err, strings.TrimSpace(string(output)))}
	}

	return nil
}

func compareNames(label string, actual, expected []string) []string {
	actualSet := toSet(actual)
	expectedSet := toSet(expected)
	var failures []string

	for _, name := range expected {
		if !actualSet[name] {
			failures = append(failures, fmt.Sprintf("missing %s %s", label, name))
		}
	}
	for _, name := range actual {
		if !expectedSet[name] {
			failures = append(failures, fmt.Sprintf("unexpected %s %s", label, name))
		}
	}

	return failures
}

func toSet(values []string) map[string]bool {
	result := make(map[string]bool, len(values))
	for _, value := range values {
		result[value] = true
	}

	return result
}

func missingTerms(doc corpus, terms []string) []string {
	var failures []string
	for _, term := range terms {
		if !containsTerm(doc.content, term) {
			failures = append(failures, fmt.Sprintf("%s missing term %q", doc.label, term))
		}
	}

	return failures
}

func containsTerm(content string, term string) bool {
	if strings.Contains(content, term) {
		return true
	}

	return strings.Contains(normalizeSpace(content), normalizeSpace(term))
}

func normalizeSpace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func readCorpus(label, path string) corpus {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read %s (%s): %w", label, path, err))
	}

	return corpus{label: label, content: string(data)}
}

func sortedKeys(values map[string]bool) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}

func receiverBaseName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return receiverBaseName(t.X)
	case *ast.IndexExpr:
		return receiverBaseName(t.X)
	case *ast.IndexListExpr:
		return receiverBaseName(t.X)
	default:
		return ""
	}
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}
