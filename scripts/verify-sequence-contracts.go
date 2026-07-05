package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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

	sourceRoot := cleanAbs(*sourceDir)
	docsRoot := cleanAbs(*outDir)

	englishDocs := readCorpus("English sequence docs", filepath.Join(docsRoot, "docs/features/sequence.md"))
	chineseDocs := readCorpus("Chinese sequence docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/features/sequence.md"))
	publicIndex := readCorpus("English public API index", filepath.Join(docsRoot, "docs/reference/public-api-index.md"))
	chinesePublicIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"))

	publicSymbols := exportedTopLevelNames(filepath.Join(sourceRoot, "sequence"))
	expectedSymbols := []string{
		"DBStore", "DBStoreTableName", "ErrInvalidCount",
		"ErrRuleNotFound", "ErrSequenceOverflow",
		"FormatDate", "Generator", "MemoryStore", "NewMemoryStore",
		"NewDBStore", "NewRedisStore", "OverflowError", "OverflowExtend",
		"OverflowReset", "OverflowStrategy", "RedisStore", "ResetCycle",
		"ResetDaily", "ResetMonthly", "ResetNone", "ResetQuarterly",
		"ResetWeekly", "ResetYearly", "Rule", "RuleModel", "Store",
	}

	var failures []string
	failures = append(failures, compareSymbols(publicSymbols, expectedSymbols)...)

	for _, doc := range []corpus{englishDocs, chineseDocs, publicIndex, chinesePublicIndex} {
		for _, symbol := range publicSymbols {
			failures = append(failures, missingTerm(doc, symbol)...)
		}
	}

	sourceChecks := []struct {
		path  string
		terms []string
	}{
		{
			path: "sequence/rule.go",
			terms: []string{
				"ResetNone      ResetCycle = \"N\"",
				"ResetDaily     ResetCycle = \"D\"",
				"ResetWeekly    ResetCycle = \"W\"",
				"ResetMonthly   ResetCycle = \"M\"",
				"ResetQuarterly ResetCycle = \"Q\"",
				"ResetYearly    ResetCycle = \"Y\"",
				"OverflowError OverflowStrategy = \"error\"",
				"OverflowReset OverflowStrategy = \"reset\"",
				"OverflowExtend OverflowStrategy = \"extend\"",
				"CurrentValue     int",
				"LastResetAt      *timex.DateTime",
				"resetAt := *r.LastResetAt",
			},
		},
		{
			path: "sequence/policy.go",
			terms: []string{
				"count < 1",
				"ErrInvalidCount",
				"rule.ResetCycle == \"\" || rule.ResetCycle == ResetNone",
				"rule.LastResetAt == nil",
				"base = rule.StartValue",
				"case OverflowReset:",
				"if resetNeeded",
				"case OverflowExtend:",
				"default:",
				"ErrSequenceOverflow",
			},
		},
		{
			path: "sequence/format.go",
			terms: []string{
				"\"yyyy\", \"2006\"",
				"\"yy\", \"06\"",
				"\"MM\", \"01\"",
				"\"dd\", \"02\"",
				"\"HH\", \"15\"",
				"\"mm\", \"04\"",
				"\"ss\", \"05\"",
				"return \"\"",
			},
		},
		{
			path: "sequence/memory_store.go",
			terms: []string{
				"rules collections.ConcurrentMap[string, *Rule]",
				"locks collections.ConcurrentMap[string, *sync.Mutex]",
				"s.rules.Put(rule.Key, rule.Clone())",
				"mu.Lock()",
				"!ok || !rule.IsActive",
				"rule.CurrentValue = rule.StartValue",
				"rule.CurrentValue += rule.SeqStep * count",
				"return rule.Clone(), rule.CurrentValue, nil",
			},
		},
		{
			path: "sequence/db_store.go",
			terms: []string{
				"const DBStoreTableName = \"sys_sequence_rule\"",
				"type RuleModel struct",
				"`bun:\"table:sys_sequence_rule,alias:ssr\"`",
				"func NewDBStore(db orm.DB) *DBStore",
				"Model((*RuleModel)(nil))",
				"IfNotExists()",
				"s.db.RunInTx(ctx, func(ctx context.Context, tx orm.DB) error",
				"cb.Equals(\"key\", key).",
				"IsTrue(\"is_active\")",
				"ForUpdate()",
				"return ErrRuleNotFound",
				"rule.CurrentValue += rule.SeqStep * count",
				"Set(\"current_value\", newValue)",
				"query.Set(\"last_reset_at\", rule.LastResetAt)",
			},
		},
		{
			path: "sequence/redis_store.go",
			terms: []string{
				"const redisSequencePrefix = \"vef:sequence:\"",
				"func NewRedisStore(client *redis.Client) *RedisStore",
				"rKey := redisSequencePrefix + key",
				"s.client.Watch(ctx, func(tx *redis.Tx) error",
				"if len(fields) == 0",
				"return ErrRuleNotFound",
				"errors.Is(err, redis.TxFailedErr)",
				"func (s *RedisStore) RegisterRule(ctx context.Context, rule *Rule) error",
				"\"current_value\":     rule.CurrentValue",
				"fields[\"last_reset_at\"] = rule.LastResetAt.String()",
				"parseRedisRule",
			},
		},
		{
			path: "sequence/sequence.go",
			terms: []string{
				"Generate(ctx context.Context, key string) (string, error)",
				"GenerateN(ctx context.Context, key string, count int) ([]string, error)",
				"Reserve(ctx context.Context, key string, count int, now timex.DateTime)",
				"read-modify-write must be serialized per key",
			},
		},
		{
			path: "internal/sequence/engine.go",
			terms: []string{
				"GenerateN(ctx, key, 1)",
				"count < 1",
				"sequence.ErrInvalidCount",
				"timex.Now()",
				"seqValue := newValue - (count-1-i)*rule.SeqStep",
				"rule.Prefix",
				"rule.DateFormat",
				"rule.SeqLength",
				"rule.Suffix",
			},
		},
		{
			path: "internal/sequence/module.go",
			terms: []string{
				"fx.As(fx.Self())",
				"fx.As(new(sequence.Store))",
				"NewGenerator",
			},
		},
	}

	for _, check := range sourceChecks {
		source := readCorpus(check.path, filepath.Join(sourceRoot, check.path))
		failures = append(failures, missingTerms(source, check.terms)...)
	}

	docTerms := []string{
		"CurrentValue", "LastResetAt", "Rule.Clone()",
		"ResetCycle", "ResetNone", "ResetDaily", "ResetWeekly",
		"ResetMonthly", "ResetQuarterly", "ResetYearly",
		"LastResetAt == nil", "OverflowError", "OverflowReset",
		"OverflowExtend", "ErrSequenceOverflow", "SeqLength",
		"Store.Reserve", "NewMemoryStore", "MemoryStore.Register",
		"NewDBStore", "DBStore", "DBStoreTableName", "DBStore.Init",
		"RuleModel", "NewRedisStore", "RedisStore",
		"RedisStore.RegisterRule", "vef:sequence:<key>", "WATCH",
		"GenerateN", "SeqStep", "ErrInvalidCount", "ErrRuleNotFound",
		"FormatDate", "yyyy", "yy", "MM", "dd", "HH", "mm", "ss",
	}
	for _, doc := range []corpus{englishDocs, chineseDocs} {
		failures = append(failures, missingTerms(doc, docTerms)...)
		failures = append(failures, forbids(doc, []string{"NewGenerator", "buildSerialNumbers"})...)
	}

	failures = append(failures, missingTerms(englishDocs, []string{
		"non-durable", "lost on process restart", "overwrites any existing rule",
		"deep copy", "read-modify-write path per\nrule key", "not a maximum",
		"being truncated", "unknown", "fall back to `OverflowError`",
		"spaced by that step",
	})...)
	failures = append(failures, missingTerms(chineseDocs, []string{
		"不持久化", "进程重启后计数器和已注册\nrule 都会丢失", "覆盖相同 `Key`",
		"深拷贝", "原子操作", "read-modify-write 路径", "不是最大长度",
		"不会被截断", "防御性\nfallback", "回退为\n`OverflowError`",
		"按该步长间隔递增",
	})...)

	sort.Strings(failures)
	if len(failures) > 0 {
		panic(fmt.Errorf("sequence contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	fmt.Printf("Sequence contract docs verified: %d public symbols, %d source files, 2 doc mirrors\n", len(publicSymbols), len(sourceChecks))
}

func exportedTopLevelNames(dir string) []string {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
	}, 0)
	if err != nil {
		panic(fmt.Errorf("failed to parse sequence package: %w", err))
	}

	names := make(map[string]bool)
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				switch d := decl.(type) {
				case *ast.FuncDecl:
					if d.Recv == nil && d.Name.IsExported() {
						names[d.Name.Name] = true
					}
				case *ast.GenDecl:
					for _, spec := range d.Specs {
						switch s := spec.(type) {
						case *ast.TypeSpec:
							if s.Name.IsExported() {
								names[s.Name.Name] = true
							}
						case *ast.ValueSpec:
							for _, name := range s.Names {
								if name.IsExported() {
									names[name.Name] = true
								}
							}
						}
					}
				}
			}
		}
	}

	result := make([]string, 0, len(names))
	for name := range names {
		result = append(result, name)
	}
	sort.Strings(result)

	return result
}

func compareSymbols(current, expected []string) []string {
	currentSet := set(current)
	expectedSet := set(expected)
	var failures []string
	for _, name := range current {
		if !expectedSet[name] {
			failures = append(failures, "unexpected public sequence symbol: "+name)
		}
	}
	for _, name := range expected {
		if !currentSet[name] {
			failures = append(failures, "missing expected public sequence symbol: "+name)
		}
	}

	return failures
}

func set(values []string) map[string]bool {
	result := make(map[string]bool, len(values))
	for _, value := range values {
		result[value] = true
	}

	return result
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
		failures = append(failures, missingTerm(c, term)...)
	}

	return failures
}

func missingTerm(c corpus, term string) []string {
	if !strings.Contains(c.content, term) {
		return []string{fmt.Sprintf("%s missing term: %s", c.label, term)}
	}

	return nil
}

func forbids(c corpus, terms []string) []string {
	var failures []string
	for _, term := range terms {
		if strings.Contains(c.content, term) {
			failures = append(failures, fmt.Sprintf("%s must not present internal term: %s", c.label, term))
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
