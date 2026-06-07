package main

import (
	"flag"
	"fmt"
	"os"
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

	checks := []struct {
		sourcePath      string
		sourceTerms     []string
		docTerms        []string
		englishDocTerms []string
		chineseDocTerms []string
	}{
		{
			sourcePath: "expression/engine.go",
			sourceTerms: []string{
				"Engine interface", "Evaluate(ctx context.Context, source string, env any)",
				"Compile(source string, opts ...CompileOption)", "Program interface",
				"Run(ctx context.Context, env any)", "Source() string",
				"EvaluateAs[T any]", "DecodeValue[T](value)", "Match",
				"e.Compile(source, AsPredicate())", "value.Bool()",
			},
			docTerms: []string{
				"type Engine interface", "Evaluate(ctx context.Context, source string, env any)",
				"Compile(source string, opts ...CompileOption)", "type Program interface",
				"Run(ctx context.Context, env any)", "Source() string",
				"EvaluateAs[T]", "DecodeValue[T](value)", "Match",
				"expression.AsPredicate()", "boolean result",
			},
			englishDocTerms: []string{"already-canceled context", "interrupt evaluation already in flight"},
			chineseDocTerms: []string{"已经取消的 context", "无法中断"},
		},
		{
			sourcePath: "expression/value.go",
			sourceTerms: []string{
				"Value struct", "NewValue(raw any)", "Interface() any",
				"IsNil() bool", "Bool() (bool, error)", "ErrUnexpectedType",
				"Decode(target any) error", "json.Marshal(v.raw)",
				"json.Unmarshal(data, target)", "DecodeValue[T any]",
			},
			docTerms: []string{
				"NewValue(raw)", "Value.Interface()", "Value.IsNil()",
				"Value.Bool()", "expression.ErrUnexpectedType",
				"Value.Decode(target)", "DecodeValue[T](value)", "JSON",
				"nil result",
			},
			englishDocTerms: []string{"zero value"},
			chineseDocTerms: []string{"零值"},
		},
		{
			sourcePath: "expression/options.go",
			sourceTerms: []string{
				"CompileOptions struct", "Predicate bool", "CompileOption func(*CompileOptions)",
				"AsPredicate()", "o.Predicate = true",
			},
			docTerms: []string{
				"CompileOption", "CompileOptions", "CompileOptions.Predicate",
				"AsPredicate()", "Predicate",
			},
		},
		{
			sourcePath: "expression/api_errors.go",
			sourceTerms: []string{
				"ErrCodeEvaluationFailed = 2500", "ErrEvaluationFailed",
				"i18n.T(\"expression_evaluation_failed\")", "result.WithCode(ErrCodeEvaluationFailed)",
			},
			docTerms: []string{
				"expression.ErrEvaluationFailed", "ErrCodeEvaluationFailed", "2500",
				"expression_evaluation_failed",
			},
		},
		{
			sourcePath: "expression/errors.go",
			sourceTerms: []string{
				"ErrUnexpectedType", "expression: unexpected result type",
			},
			docTerms: []string{
				"ErrUnexpectedType", "expression: unexpected result type",
			},
		},
		{
			sourcePath: "internal/expression/module.go",
			sourceTerms: []string{
				"zen.New", "NewEngineResolver", "vef:api:handler_param_resolvers",
				"NewFieldTransformer", "vef:mold:field_transformers",
			},
			docTerms: []string{
				"core boot graph", "Zen-backed engine",
			},
			englishDocTerms: []string{"handler parameter resolver", "mold field transformer"},
			chineseDocTerms: []string{"handler parameter\nresolver", "mold field transformer"},
		},
		{
			sourcePath: "internal/expression/resolver.go",
			sourceTerms: []string{
				"NewEngineResolver", "api.HandlerParamResolver",
				"reflect.TypeFor[expression.Engine]()", "Resolve(fiber.Ctx)",
				"reflect.ValueOf(r.engine)",
			},
			docTerms: []string{
				"expression.Engine",
			},
			englishDocTerms: []string{"handler parameter resolver", "API handlers can request it directly"},
			chineseDocTerms: []string{"handler parameter\nresolver", "API handler 可以直接声明这个参数"},
		},
		{
			sourcePath: "internal/expression/transformer.go",
			sourceTerms: []string{
				"fieldTransformerTag = \"expr\"", "NewFieldTransformer",
				"mold.FieldTransformer", "Tag() string", "fl.Param()",
				"ErrEmptyExpression", "ErrFieldNotSettable", "fl.Struct()",
				"env = s.Interface()", "t.engine.Evaluate(ctx, source, env)",
				"value.Decode(field.Addr().Interface())",
			},
			docTerms: []string{
				"mold:\"expr=price * qty\"", "0x2C",
			},
			englishDocTerms: []string{
				"mold field transformer named `expr`", "containing struct",
				"expression environment", "declaration order", "zero value",
				"decoded result",
			},
			chineseDocTerms: []string{
				"名为 `expr` 的 mold field transformer", "当前结构体",
				"表达式环境", "声明顺序", "零值", "解码后的结果",
			},
		},
		{
			sourcePath: "internal/expression/errors.go",
			sourceTerms: []string{
				"ErrEmptyExpression", "expression: empty expression in field tag",
				"ErrFieldNotSettable", "expression: target field is not settable",
			},
			docTerms: []string{
				"expression: empty expression in field tag",
				"expression: target field is not settable",
			},
		},
		{
			sourcePath: "internal/expression/zen/engine.go",
			sourceTerms: []string{
				"zen.EvaluateExpression[any](source, env)", "zen.EvaluateUnaryExpression",
				"ctx.Err()", "program{source: source, predicate: o.Predicate}",
				"ErrEvaluationFailed", "fmt.Errorf(\"%w: %w\"",
			},
			docTerms: []string{
				"github.com/gorules/zen-go", "CGO_ENABLED=1",
				"best-effort", "Program.Run(...)", "ErrEvaluationFailed",
			},
			englishDocTerms: []string{"already-canceled context", "interrupt evaluation already in flight", "malformed expressions"},
			chineseDocTerms: []string{"已经取消的 context", "无法中断", "格式错误的表达式"},
		},
	}

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)
	englishDocs := readCorpus("English expression docs", filepath.Join(docsRoot, "docs/features/expression.md"))
	chineseDocs := readCorpus("Chinese expression docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/features/expression.md"))
	docs := []corpus{englishDocs, chineseDocs}

	var failures []string
	for _, check := range checks {
		source := readCorpus(check.sourcePath, filepath.Join(sourceRoot, check.sourcePath))
		failures = append(failures, missingTerms(source, check.sourceTerms)...)
		for _, doc := range docs {
			failures = append(failures, missingTerms(doc, check.docTerms)...)
		}
		failures = append(failures, missingTerms(englishDocs, check.englishDocTerms)...)
		failures = append(failures, missingTerms(chineseDocs, check.chineseDocTerms)...)
	}

	sort.Strings(failures)
	if len(failures) > 0 {
		panic(fmt.Errorf("expression contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	fmt.Printf("Expression contract docs verified: %d source files, %d doc mirrors\n", len(checks), len(docs))
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

func cleanAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return filepath.Clean(abs)
}
