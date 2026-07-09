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

	englishDocs := readCorpus("English mapx docs", filepath.Join(docsRoot, "docs/utilities/mapx.md"))
	chineseDocs := readCorpus("Chinese mapx docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/utilities/mapx.md"))
	publicIndex := readCorpus("English public API index", filepath.Join(docsRoot, "docs/reference/public-api-index.md"))
	chinesePublicIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"))

	mapxDir := filepath.Join(sourceRoot, "mapx")
	expectedSurface := packageSurface{
		funcs: []string{
			"FromMap",
			"NewDecoder",
			"ToMap",
			"WithAllowUnsetPointer",
			"WithDecodeHook",
			"WithDecodeNil",
			"WithErrorUnset",
			"WithErrorUnused",
			"WithIgnoreUntaggedFields",
			"WithMatchName",
			"WithMetadata",
			"WithTagName",
			"WithWeaklyTypedInput",
			"WithZeroFields",
		},
		types: []string{
			"DecoderOption",
			"Metadata",
		},
		vars: []string{
			"DecoderHook",
			"ErrCollectionSetIncompatibleKind",
			"ErrCollectionSetNegative",
			"ErrCollectionSetNilElement",
			"ErrCollectionSetNonInteger",
			"ErrCollectionSetNotFinite",
			"ErrCollectionSetOverflow",
			"ErrCollectionSetUnsupportedTarget",
			"ErrInvalidFromMapType",
			"ErrInvalidToMapValue",
		},
	}

	var failures []string
	exported := exportedPackageSurface(mapxDir)
	failures = append(failures, compareNames("mapx const", exported.consts, expectedSurface.consts)...)
	failures = append(failures, compareNames("mapx func", exported.funcs, expectedSurface.funcs)...)
	failures = append(failures, compareNames("mapx type", exported.types, expectedSurface.types)...)
	failures = append(failures, compareNames("mapx var", exported.vars, expectedSurface.vars)...)
	failures = append(failures, compareNames("mapx method", exported.methods, expectedSurface.methods)...)
	failures = append(failures, compareNames("mapx exported field", exported.fields, expectedSurface.fields)...)

	for _, doc := range []corpus{publicIndex, chinesePublicIndex} {
		failures = append(failures, missingTerms(doc, publicIndexTerms())...)
	}
	for _, doc := range []corpus{englishDocs, chineseDocs} {
		failures = append(failures, missingTerms(doc, publicDocSurfaceTerms())...)
	}

	failures = append(failures, missingTerms(englishDocs, []string{
		"Creates a `mapstructure.Decoder` with VEF defaults, then applies options in order",
		"non-struct input returns `ErrInvalidToMapValue`",
		"non-struct `T` returns `ErrInvalidFromMapType`",
		"Sets `DecoderConfig.IgnoreUntaggedFields` to the supplied boolean",
		"Replaces `DecoderConfig.DecodeHook`",
		"compose with `mapx.DecoderHook` yourself to preserve defaults",
		"default is `mapKey == lo.CamelCase(fieldName)`",
		"`WithDecodeHook(hook)` | Replace the default decode hook.",
		"`WithMatchName(fn)` | Custom field-name matcher (default: exact match against `lo.CamelCase(fieldName)`).",
		"`json.RawMessage` — marshals the source value to JSON bytes",
		"`*multipart.FileHeader` — picks the only entry when the source is `[]*multipart.FileHeader` with length 1",
		"`collections.Set` / `SortedSet` / `ConcurrentSet` / `ConcurrentSortedSet` — turns a slice or array into the corresponding set type",
		"Collection-set decoding is registered for `string`, signed integers, unsigned\nintegers, `float32`, and `float64`",
		"It rejects nil elements, string/numeric\nfamily mismatches, numeric overflow, fractional floats targeting integer sets,\nNaN or infinity targeting integer sets, and negative values targeting unsigned\nsets",
		"`WithDecodeHook(myHook)` replaces the default composed hook",
		"compose your hook with `mapx.DecoderHook` before passing it to\n`WithDecodeHook`",
	})...)
	failures = append(failures, missingTerms(chineseDocs, []string{
		"使用 VEF 默认值创建 `mapstructure.Decoder`，然后按顺序应用 options",
		"非 struct 输入返回 `ErrInvalidToMapValue`",
		"非 struct `T` 返回 `ErrInvalidFromMapType`",
		"把 `DecoderConfig.IgnoreUntaggedFields` 设为传入的 boolean",
		"替换 `DecoderConfig.DecodeHook`",
		"需要自行和 `mapx.DecoderHook` compose",
		"默认是 `mapKey == lo.CamelCase(fieldName)`",
		"`WithDecodeHook(hook)` | 替换默认 decode hook。",
		"`WithMatchName(fn)` | 自定义字段名匹配函数（默认与 `lo.CamelCase(fieldName)` 精确比较）。",
		"`json.RawMessage` —— 将源值 marshal 成 JSON bytes",
		"`*multipart.FileHeader` —— 源是长度为 1 的 `[]*multipart.FileHeader` 时取唯一一项",
		"`collections.Set` / `SortedSet` / `ConcurrentSet` / `ConcurrentSortedSet` —— 把 slice 或 array 转为对应的集合类型",
		"collection-set 解码为 `string`、有符号整数、无符号整数、`float32` 和\n`float64` 注册",
		"它会拒绝 nil element、string/numeric family mismatch、\nnumeric overflow、fractional float 转 integer set、NaN 或 infinity 转\ninteger set，以及负数转 unsigned set",
		"`WithDecodeHook(myHook)` 会替换默认 composed hook",
		"自定义 hook 和 `mapx.DecoderHook` compose，再传给 `WithDecodeHook`",
	})...)

	sourceChecks := []struct {
		path  string
		terms []string
	}{
		{
			path: "mapx/decoder.go",
			terms: []string{
				"defaultDecoderTagName = \"json\"",
				"DecoderHook = mapstructure.ComposeDecodeHookFunc(",
				"convertJSONRawMessage",
				"convertFileHeader",
				"convertSliceToCollectionSet",
				"mapstructure.TextUnmarshallerHookFunc()",
				"mapstructure.StringToTimeHookFunc(time.DateTime)",
				"mapstructure.StringToTimeLocationHookFunc()",
				"mapstructure.StringToTimeDurationHookFunc()",
				"mapstructure.StringToURLHookFunc()",
				"mapstructure.StringToIPHookFunc()",
				"mapstructure.StringToIPNetHookFunc()",
				"mapstructure.StringToNetIPPrefixHookFunc()",
				"mapstructure.StringToNetIPAddrHookFunc()",
				"mapstructure.StringToNetIPAddrPortHookFunc()",
				"mapstructure.StringToBasicTypeHookFunc()",
				"DecoderOption func(c *mapstructure.DecoderConfig)",
				"Metadata = mapstructure.Metadata",
				"func NewDecoder(result any, options ...DecoderOption) (*mapstructure.Decoder, error)",
				"TagName:              defaultDecoderTagName",
				"IgnoreUntaggedFields: false",
				"DecodeHook:           DecoderHook",
				"Squash:               true",
				"SquashTagOption:      \"inline\"",
				"return mapKey == lo.CamelCase(fieldName)",
				"ErrorUnused:       false",
				"ErrorUnset:        false",
				"ZeroFields:        false",
				"AllowUnsetPointer: false",
				"Metadata:          nil",
				"WeaklyTypedInput:  false",
				"DecodeNil:         false",
				"Result:            result",
				"for _, option := range options",
				"option(config)",
				"return mapstructure.NewDecoder(config)",
				"func WithTagName(tagName string) DecoderOption",
				"c.TagName = tagName",
				"func WithIgnoreUntaggedFields(ignoreUntaggedFields bool) DecoderOption",
				"c.IgnoreUntaggedFields = ignoreUntaggedFields",
				"func WithDecodeHook(decodeHook mapstructure.DecodeHookFunc) DecoderOption",
				"c.DecodeHook = decodeHook",
				"func WithMatchName(matchName func(mapKey, fieldName string) bool) DecoderOption",
				"c.MatchName = matchName",
				"c.ErrorUnused = true",
				"c.ErrorUnset = true",
				"c.ZeroFields = true",
				"c.AllowUnsetPointer = true",
				"c.Metadata = metadata",
				"c.WeaklyTypedInput = true",
				"c.DecodeNil = true",
				"func ToMap(value any, options ...DecoderOption) (map[string]any, error)",
				"return nil, ErrInvalidToMapValue",
				"func FromMap[T any](value map[string]any, options ...DecoderOption) (*T, error)",
				"reflect.TypeFor[T]().Kind() != reflect.Struct",
				"return nil, ErrInvalidFromMapType",
			},
		},
		{
			path: "mapx/decode_hook.go",
			terms: []string{
				"func convertJSONRawMessage(_, to reflect.Type, value any) (any, error)",
				"json.Marshal(value)",
				"return json.RawMessage(data), nil",
				"func convertFileHeader(from, to reflect.Type, value any) (any, error)",
				"if files := value.([]*multipart.FileHeader); len(files) == 1",
				"return files[0], nil",
				"func convertSliceToCollectionSet(from, to reflect.Type, value any) (any, error)",
				"from.Kind() != reflect.Slice && from.Kind() != reflect.Array",
				"builder, ok := collectionSetBuilders[to]",
				"return builder(reflect.ValueOf(value))",
				"registerCollectionSet[string](registry)",
				"registerCollectionSet[int](registry)",
				"registerCollectionSet[int8](registry)",
				"registerCollectionSet[int16](registry)",
				"registerCollectionSet[int32](registry)",
				"registerCollectionSet[int64](registry)",
				"registerCollectionSet[uint](registry)",
				"registerCollectionSet[uint8](registry)",
				"registerCollectionSet[uint16](registry)",
				"registerCollectionSet[uint32](registry)",
				"registerCollectionSet[uint64](registry)",
				"registerCollectionSet[float32](registry)",
				"registerCollectionSet[float64](registry)",
				"ErrCollectionSetNilElement",
				"ErrCollectionSetIncompatibleKind",
				"ErrCollectionSetOverflow",
				"ErrCollectionSetNotFinite",
				"ErrCollectionSetNonInteger",
				"ErrCollectionSetNegative",
				"ErrCollectionSetUnsupportedTarget",
				"math.IsNaN(f) || math.IsInf(f, 0)",
				"f != math.Trunc(f)",
				"f < float64MinInt64 || f >= float64Pow63",
				"f < 0 || f >= float64Pow64",
			},
		},
		{
			path: "mapx/errors.go",
			terms: []string{
				"ErrInvalidToMapValue = errors.New(\"the value of ToMap function must be a struct\")",
				"ErrInvalidFromMapType = errors.New(\"the type parameter of FromMap function must be a struct\")",
				"ErrCollectionSetNilElement = errors.New(\"nil element cannot be added to collections set\")",
				"ErrCollectionSetIncompatibleKind = errors.New(\"incompatible source kind for collections set element\")",
				"ErrCollectionSetOverflow = errors.New(\"value overflows collections set element type\")",
				"ErrCollectionSetNonInteger = errors.New(\"non-integer value cannot be converted to integer set element\")",
				"ErrCollectionSetNotFinite = errors.New(\"non-finite value cannot be converted to integer set element\")",
				"ErrCollectionSetNegative = errors.New(\"negative value cannot be converted to unsigned set element\")",
				"ErrCollectionSetUnsupportedTarget = errors.New(\"unsupported target kind for collections set\")",
			},
		},
		{
			path: "mapx/decoder_test.go",
			terms: []string{
				"TestNewDecoder",
				"TestToMap",
				"NonStructValue",
				"PointerToStruct",
				"StructWithEmbedding",
				"CustomTagName",
				"TestFromMap",
				"NonStructTypeParameter",
				"TestDecoderOptions",
				"WithIgnoreUntaggedFields(true)",
				"WithWeaklyTypedInput()",
				"WithMetadata(&metadata)",
				"TestComplexTypeConversions",
				"TimeConversion",
				"DurationConversion",
				"URLConversion",
				"IPConversion",
				"TestFileHeaderConversion",
				"SliceWithSingleFileToSinglePointer",
				"EmptySliceToSinglePointer",
				"NilSliceToSinglePointer",
				"TestDecodeOrderDirection",
			},
		},
		{
			path: "mapx/decode_hook_test.go",
			terms: []string{
				"TestConvertSliceToCollectionSetHappyPath",
				"SetString",
				"SortedSetString",
				"ConcurrentSetString",
				"ConcurrentSortedSetString",
				"TestConvertSliceToCollectionSetNumericTypes",
				"SetInt",
				"SetUint16WithinBounds",
				"SetFloat32",
				"TestConvertSliceToCollectionSetRejections",
				"FractionalFloatToInt",
				"OverflowFloatToInt8",
				"NegativeFloatToUint",
				"NaNToInt",
				"InfinityToInt",
				"StringElementToIntSet",
				"NumericElementToStringSet",
				"NilElement",
				"Float64Pow63ToInt64Overflow",
				"Float64Pow64ToUint64Overflow",
				"Float64MinInt64Accepted",
				"TestRegistryCoverage",
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
		panic(fmt.Errorf("mapx contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	topLevelPublic := len(exported.consts) + len(exported.funcs) + len(exported.types) + len(exported.vars)
	fmt.Printf("Mapx contract docs verified: %d top-level public symbols, %d public methods, %d public fields, %d source files, 2 doc mirrors\n",
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
		panic(fmt.Errorf("failed to parse mapx package: %w", err))
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
		"`mapx.DecoderHook`",
		"`mapx.DecoderOption`",
		"`mapx.Metadata`",
		"`mapx.NewDecoder(result, options...)`",
		"`mapx.ToMap(value, options...)`",
		"`mapx.FromMap[T](value, options...)`",
		"`mapx.WithTagName(tagName)`",
		"`mapx.WithIgnoreUntaggedFields(ignore)`",
		"`mapx.WithDecodeHook(decodeHook)`",
		"`mapx.WithMatchName(matchName)`",
		"`mapx.WithErrorUnused()`",
		"`mapx.WithErrorUnset()`",
		"`mapx.WithZeroFields()`",
		"`mapx.WithAllowUnsetPointer()`",
		"`mapx.WithMetadata(metadata)`",
		"`mapx.WithWeaklyTypedInput()`",
		"`mapx.WithDecodeNil()`",
		"`mapx.ErrInvalidToMapValue`",
		"`mapx.ErrInvalidFromMapType`",
		"`mapx.ErrCollectionSetNilElement`",
		"`mapx.ErrCollectionSetIncompatibleKind`",
		"`mapx.ErrCollectionSetOverflow`",
		"`mapx.ErrCollectionSetNonInteger`",
		"`mapx.ErrCollectionSetNotFinite`",
		"`mapx.ErrCollectionSetNegative`",
		"`mapx.ErrCollectionSetUnsupportedTarget`",
	}
}

func publicIndexTerms() []string {
	return []string{
		"## github.com/coldsmirk/vef-framework-go/mapx",
		"VAR DecoderHook : github.com/go-viper/mapstructure/v2.DecodeHookFunc",
		"TYPE DecoderOption : github.com/coldsmirk/vef-framework-go/mapx.DecoderOption",
		"VAR ErrCollectionSetIncompatibleKind : error",
		"VAR ErrCollectionSetNegative : error",
		"VAR ErrCollectionSetNilElement : error",
		"VAR ErrCollectionSetNonInteger : error",
		"VAR ErrCollectionSetNotFinite : error",
		"VAR ErrCollectionSetOverflow : error",
		"VAR ErrCollectionSetUnsupportedTarget : error",
		"VAR ErrInvalidFromMapType : error",
		"VAR ErrInvalidToMapValue : error",
		"FUNC FromMap : func[T any](value map[string]any, options ...github.com/coldsmirk/vef-framework-go/mapx.DecoderOption) (*T, error)",
		"TYPE Metadata : github.com/coldsmirk/vef-framework-go/mapx.Metadata",
		"FUNC NewDecoder : func(result any, options ...github.com/coldsmirk/vef-framework-go/mapx.DecoderOption) (*github.com/go-viper/mapstructure/v2.Decoder, error)",
		"FUNC ToMap : func(value any, options ...github.com/coldsmirk/vef-framework-go/mapx.DecoderOption) (map[string]any, error)",
		"FUNC WithAllowUnsetPointer : func() github.com/coldsmirk/vef-framework-go/mapx.DecoderOption",
		"FUNC WithDecodeHook : func(decodeHook github.com/go-viper/mapstructure/v2.DecodeHookFunc) github.com/coldsmirk/vef-framework-go/mapx.DecoderOption",
		"FUNC WithDecodeNil : func() github.com/coldsmirk/vef-framework-go/mapx.DecoderOption",
		"FUNC WithErrorUnset : func() github.com/coldsmirk/vef-framework-go/mapx.DecoderOption",
		"FUNC WithErrorUnused : func() github.com/coldsmirk/vef-framework-go/mapx.DecoderOption",
		"FUNC WithIgnoreUntaggedFields : func(ignoreUntaggedFields bool) github.com/coldsmirk/vef-framework-go/mapx.DecoderOption",
		"FUNC WithMatchName : func(matchName func(mapKey string, fieldName string) bool) github.com/coldsmirk/vef-framework-go/mapx.DecoderOption",
		"FUNC WithMetadata : func(metadata *github.com/coldsmirk/vef-framework-go/mapx.Metadata) github.com/coldsmirk/vef-framework-go/mapx.DecoderOption",
		"FUNC WithTagName : func(tagName string) github.com/coldsmirk/vef-framework-go/mapx.DecoderOption",
		"FUNC WithWeaklyTypedInput : func() github.com/coldsmirk/vef-framework-go/mapx.DecoderOption",
		"FUNC WithZeroFields : func() github.com/coldsmirk/vef-framework-go/mapx.DecoderOption",
	}
}

func runPackageTests(sourceRoot string) []string {
	cmd := exec.Command("go", "test", "./mapx")
	cmd.Dir = sourceRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("go test ./mapx failed: %v\n%s", err, strings.TrimSpace(string(output)))}
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
