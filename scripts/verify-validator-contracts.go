package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
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

	englishDocs := readCorpus("English validation docs", filepath.Join(docsRoot, "docs/guide/validation.md"))
	chineseDocs := readCorpus("Chinese validation docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/guide/validation.md"))
	publicIndex := readCorpus("English public API index", filepath.Join(docsRoot, "docs/reference/public-api-index.md"))
	chinesePublicIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"))

	validatorDir := filepath.Join(sourceRoot, "validator")
	expectedSurface := packageSurface{
		funcs: []string{
			"RegisterTypeFunc",
			"RegisterValidationRules",
			"Validate",
		},
		types: []string{
			"CustomTypeFunc",
			"ValidationRule",
		},
		fields: []string{
			"ValidationRule.CallValidationEvenIfNull",
			"ValidationRule.ErrMessageI18nKey",
			"ValidationRule.ErrMessageTemplate",
			"ValidationRule.ParseParam",
			"ValidationRule.RuleTag",
			"ValidationRule.Validate",
		},
	}

	var failures []string
	exported := exportedPackageSurface(validatorDir)
	failures = append(failures, compareNames("validator const", exported.consts, expectedSurface.consts)...)
	failures = append(failures, compareNames("validator func", exported.funcs, expectedSurface.funcs)...)
	failures = append(failures, compareNames("validator type", exported.types, expectedSurface.types)...)
	failures = append(failures, compareNames("validator var", exported.vars, expectedSurface.vars)...)
	failures = append(failures, compareNames("validator method", exported.methods, expectedSurface.methods)...)
	failures = append(failures, compareNames("validator exported field", exported.fields, expectedSurface.fields)...)

	for _, doc := range []corpus{publicIndex, chinesePublicIndex} {
		failures = append(failures, missingTerms(doc, publicIndexTerms())...)
	}
	for _, doc := range []corpus{englishDocs, chineseDocs} {
		failures = append(failures, missingTerms(doc, publicDocSurfaceTerms())...)
	}

	failures = append(failures, missingTerms(englishDocs, []string{
		"5 top-level exported symbols",
		"6 exported fields",
		"no exported\nmethods",
		"fingerprint is\n`747c46e376e5a77e49d9266c675461dfca729dc5cb2e9b54b231c86aab67ddcf`",
		"success returns `nil`, validation failure returns the first translated error as a bad-request `result.Error`",
		"Registers each `ValidationRule` with the shared validator and both built-in translators",
		"Registers a custom type extractor by delegating to `RegisterCustomTypeFunc`",
		"Alias-compatible function type `func(field reflect.Value) any`",
		"Passed through to go-playground `RegisterValidation`",
		"if the translation is missing, the i18n key may appear in the error message",
		"startup panics instead of leaving the validator partially configured",
		"| accepted values | `^1[3-9]\\d{9}$` |",
		"non-`decimal.Decimal` fields or invalid threshold params fail validation",
		"All three alphanum rules require at least one character",
		"business code | `result.ErrCodeBadRequest`",
		"HTTP status | `400 Bad Request`",
		"read the active\n`i18n.CurrentLanguage()` at validation time",
		"Chinese (`zh-CN`) uses the Chinese\ntranslator; other supported languages use the English translator",
		"`ErrMessageI18nKey` is checked first",
		"When\n`i18n.T(key)` returns a value different from the key",
		"go-playground translator renders `ErrMessageTemplate`",
	})...)
	failures = append(failures, missingTerms(chineseDocs, []string{
		"5 个\ntop-level exported symbols",
		"6 个 exported fields",
		"没有 exported methods",
		"fingerprint 是\n`747c46e376e5a77e49d9266c675461dfca729dc5cb2e9b54b231c86aab67ddcf`",
		"成功返回 `nil`，校验失败时把第一个翻译后的错误包装成 bad-request `result.Error`",
		"注册到共享 validator 和两个内置 translators",
		"通过 `RegisterCustomTypeFunc` 注册自定义类型提取函数",
		"与 `func(field reflect.Value) any` 兼容",
		"透传给 go-playground `RegisterValidation`",
		"如果翻译缺失，错误消息里可能显示这个 i18n key",
		"启动会 panic，不会留下部分配置好的 validator",
		"| 接受的值 | `^1[3-9]\\d{9}$` |",
		"字段不是 `decimal.Decimal` 或阈值参数不是合法 decimal 时校验失败",
		"这三个 alphanum 规则都要求至少 1 个字符",
		"业务码 | `result.ErrCodeBadRequest`",
		"HTTP 状态 | `400 Bad Request`",
		"校验发生时读取当前\n`i18n.CurrentLanguage()`",
		"中文（`zh-CN`）使用中文 translator；其他支持语言使用英文 translator",
		"自定义规则消息会先检查 `ErrMessageI18nKey`",
		"当 `i18n.T(key)` 返回不同于\nkey 的真实消息",
		"由 go-playground translator 渲染\n`ErrMessageTemplate`",
	})...)

	sourceChecks := []struct {
		path  string
		terms []string
	}{
		{
			path: "validator/validator.go",
			terms: []string{
				"validator = v.New(v.WithRequiredStructEnabled())",
				"zhtranslation.RegisterDefaultTranslations(validator, zhTranslator)",
				"entranslation.RegisterDefaultTranslations(validator, enTranslator)",
				"validator.RegisterTagNameFunc(func(field reflect.StructField) string",
				"field.Tag.Get(tagLabel)",
				"field.Tag.Get(tagLabelI18n)",
				"lo.CoalesceOrEmpty(i18n.T(label), field.Name)",
				"func activeTranslator() ut.Translator",
				"if i18n.CurrentLanguage() == i18n.DefaultLanguage",
				"return translators[langZh]",
				"return translators[langEn]",
				"func RegisterValidationRules(rules ...ValidationRule) error",
				"streams.FromSlice(rules).ForEachErr",
				"rule.register(validator, translators)",
				"type CustomTypeFunc = func(field reflect.Value) any",
				"func RegisterTypeFunc(fn CustomTypeFunc, types ...any)",
				"validator.RegisterCustomTypeFunc(fn, types...)",
				"func Validate(value any) error",
				"err := validator.Struct(value)",
				"errors.As(err, &validationErrors)",
				"validationErrors[0].Translate(activeTranslator())",
				"result.WithCode(result.ErrCodeBadRequest)",
				"result.WithStatus(fiber.StatusBadRequest)",
			},
		},
		{
			path: "validator/rule.go",
			terms: []string{
				"presetValidationRules = []ValidationRule",
				"newPhoneNumberRule()",
				"newDecimalMinRule()",
				"newDecimalMaxRule()",
				"newAlphanumUsRule()",
				"newAlphanumUsSlashRule()",
				"newAlphanumUsDotRule()",
				"type ValidationRule struct",
				"RuleTag                  string",
				"ErrMessageTemplate       string",
				"ErrMessageI18nKey        string",
				"Validate                 func(fl v.FieldLevel) bool",
				"ParseParam               func(fe v.FieldError) []string",
				"CallValidationEvenIfNull bool",
				"validator.RegisterValidation(vr.RuleTag, vr.Validate, vr.CallValidationEvenIfNull)",
				"for _, translator := range translators",
				"t.Add(vr.RuleTag, vr.ErrMessageTemplate, false)",
				"i18n.T(vr.ErrMessageI18nKey)",
				"msg != vr.ErrMessageI18nKey",
				"vr.replacePlaceholders(msg, vr.ParseParam(fe))",
				"t.T(vr.RuleTag, vr.ParseParam(fe)...)",
			},
		},
		{
			path: "validator/setup.go",
			terms: []string{
				"func setup()",
				"RegisterValidationRules(presetValidationRules...)",
				"panic(err)",
			},
		},
		{
			path: "validator/phone_number.go",
			terms: []string{
				"phoneNumberRegex = regexp.MustCompile(`^1[3-9]\\d{9}$`)",
				"return newRegexRule(\"phone_number\", phoneNumberRegex, \"{0}格式不正确\", \"validator_phone_number\")",
			},
		},
		{
			path: "validator/decimal.go",
			terms: []string{
				"return newDecimalComparisonRule(\"dec_min\", \"{0}最小只能为{1}\", \"validator_decimal_min\"",
				"dec.GreaterThanOrEqual(threshold)",
				"return newDecimalComparisonRule(\"dec_max\", \"{0}必须小于或等于{1}\", \"validator_decimal_max\"",
				"dec.LessThanOrEqual(threshold)",
				"fl.Field().Interface().(decimal.Decimal)",
				"decimal.NewFromString(fl.Param())",
				"return []string{fe.Field(), fe.Param()}",
			},
		},
		{
			path: "validator/alphanum.go",
			terms: []string{
				"regexp.MustCompile(`^[a-zA-Z0-9_]+$`)",
				"regexp.MustCompile(`^[a-zA-Z0-9_/]+$`)",
				"regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)",
				"newRegexRule(\"alphanum_us\", alphanumUsRegex",
				"newRegexRule(\"alphanum_us_slash\", alphanumUsSlashRegex",
				"newRegexRule(\"alphanum_us_dot\", alphanumUsDotRegex",
				"CallValidationEvenIfNull: false",
				"regex.MatchString(fl.Field().String())",
				"return []string{fe.Field()}",
			},
		},
		{
			path: "i18n/locales/en.json",
			terms: []string{
				"\"validator_phone_number\": \"{0} format is invalid\"",
				"\"validator_decimal_min\": \"{0} must be at least {1}\"",
				"\"validator_decimal_max\": \"{0} must be less than or equal to {1}\"",
				"\"validator_alphanum_us\": \"{0} can only contain letters, numbers and underscores\"",
				"\"validator_alphanum_us_slash\": \"{0} can only contain letters, numbers, underscores and slashes\"",
				"\"validator_alphanum_us_dot\": \"{0} can only contain letters, numbers, underscores and dots\"",
			},
		},
		{
			path: "i18n/locales/zh-CN.json",
			terms: []string{
				"\"validator_phone_number\": \"{0}格式不正确\"",
				"\"validator_decimal_min\": \"{0}最小只能为{1}\"",
				"\"validator_decimal_max\": \"{0}必须小于或等于{1}\"",
				"\"validator_alphanum_us\": \"{0}只能包含字母、数字和下划线\"",
				"\"validator_alphanum_us_slash\": \"{0}只能包含字母、数字、下划线和斜线\"",
				"\"validator_alphanum_us_dot\": \"{0}只能包含字母、数字、下划线和点\"",
			},
		},
		{
			path: "validator/validator_test.go",
			terms: []string{
				"TestValidate",
				"MissingRequiredField",
				"MultipleErrors",
				"PointerStringNil",
				"DecimalWithNonDecimalType",
				"DecimalWithInvalidParam",
				"TestReplacePlaceholders",
			},
		},
		{
			path: "validator/i18n_test.go",
			terms: []string{
				"TestValidatorI18nMessages",
				"ChineseMessages",
				"EnglishMessages",
				"TestValidatorBuiltInRuleLanguage",
				"must be a valid email address",
				"必须是一个有效的邮箱",
				"TestValidatorI18nPhoneNumber",
				"TestValidatorI18nAlphanumRules",
			},
		},
		{
			path: "validator/phone_number_test.go",
			terms: []string{
				"TestPhoneNumberValidation",
				"ValidPhoneStartsWith13",
				"ValidPhoneStartsWith19",
				"InvalidStartsWith12",
				"InvalidTooShort",
				"InvalidContainsLetters",
			},
		},
		{
			path: "validator/decimal_test.go",
			terms: []string{
				"TestDecimalMinValidation",
				"ValidExactMinimum",
				"InvalidBelowMinimum",
				"TestDecimalMaxValidation",
				"ValidExactMaximum",
				"InvalidAboveMaximum",
				"TestDecimalRangeValidation",
			},
		},
		{
			path: "validator/alphanum_test.go",
			terms: []string{
				"TestAlphanumUs",
				"WithUnderscores",
				"WithSlash",
				"EmptyString",
				"TestAlphanumUsSlash",
				"WithSlashes",
				"WithDot",
				"TestAlphanumUsDot",
				"WithDots",
				"WithSpecialChars",
				"TestAlphanumRulesCombined",
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
		panic(fmt.Errorf("validator contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	topLevelPublic := len(exported.consts) + len(exported.funcs) + len(exported.types) + len(exported.vars)
	fmt.Printf("Validator contract docs verified: %d top-level public symbols, %d public methods, %d public fields, %d source files, 2 doc mirrors\n",
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
		panic(fmt.Errorf("failed to parse validator package: %w", err))
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

func publicDocSurfaceTerms() []string {
	return []string{
		"`validator.Validate(value)`",
		"`validator.RegisterValidationRules(rules...)`",
		"`validator.RegisterTypeFunc(fn, types...)`",
		"`validator.CustomTypeFunc`",
		"`validator.ValidationRule`",
		"`ValidationRule.RuleTag`",
		"`ValidationRule.ErrMessageTemplate`",
		"`ValidationRule.ErrMessageI18nKey`",
		"`ValidationRule.Validate`",
		"`ValidationRule.ParseParam`",
		"`ValidationRule.CallValidationEvenIfNull`",
	}
}

func publicIndexTerms() []string {
	return []string{
		"## github.com/coldsmirk/vef-framework-go/validator",
		"TYPE CustomTypeFunc : github.com/coldsmirk/vef-framework-go/validator.CustomTypeFunc",
		"FUNC RegisterTypeFunc : func(fn github.com/coldsmirk/vef-framework-go/validator.CustomTypeFunc, types ...any)",
		"FUNC RegisterValidationRules : func(rules ...github.com/coldsmirk/vef-framework-go/validator.ValidationRule) error",
		"FUNC Validate : func(value any) error",
		"TYPE ValidationRule : github.com/coldsmirk/vef-framework-go/validator.ValidationRule",
		"FIELD RuleTag : string [field_order=1 tag=\"\"]",
		"FIELD ErrMessageTemplate : string [field_order=2 tag=\"\"]",
		"FIELD ErrMessageI18nKey : string [field_order=3 tag=\"\"]",
		"FIELD Validate : func(fl github.com/go-playground/validator/v10.FieldLevel) bool [field_order=4 tag=\"\"]",
		"FIELD ParseParam : func(fe github.com/go-playground/validator/v10.FieldError) []string [field_order=5 tag=\"\"]",
		"FIELD CallValidationEvenIfNull : bool [field_order=6 tag=\"\"]",
	}
}

func runPackageTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./validator")
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("go test ./validator failed: %v\n%s", err, strings.TrimSpace(string(output)))}
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

func readCorpus(label string, path string) corpus {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read %s at %s: %w", label, path, err))
	}

	return corpus{label: label, content: string(data)}
}

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(fmt.Errorf("failed to resolve %s: %w", path, err))
	}

	return filepath.Clean(abs)
}
