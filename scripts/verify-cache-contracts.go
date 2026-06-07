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

	englishDocs := readCorpus("English cache docs", filepath.Join(docsRoot, "docs/features/cache.md"))
	chineseDocs := readCorpus("Chinese cache docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/features/cache.md"))
	publicIndex := readCorpus("English public API index", filepath.Join(docsRoot, "docs/reference/public-api-index.md"))
	chinesePublicIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"))

	cacheDir := filepath.Join(sourceRoot, "cache")
	publicSymbols := exportedTopLevelNames(cacheDir)
	expectedSymbols := []string{
		"Cache", "ErrCacheClosed", "ErrLoaderRequired", "ErrMemoryLimitExceeded",
		"ErrTypeAssertionFailed", "EvictionPolicy", "EvictionPolicyFIFO",
		"EvictionPolicyLFU", "EvictionPolicyLRU", "EvictionPolicyNone",
		"GetFunc", "Invalidating", "Key", "KeyBuilder", "KeyedLoaderFunc",
		"LoaderFunc", "MemoryOption", "NewInvalidating", "NewMemory",
		"NewPrefixKeyBuilder", "NewRedis", "PrefixKeyBuilder", "RedisOption",
		"SetFunc", "SingleflightMixin", "WithMemDefaultTTL",
		"WithMemEvictionPolicy", "WithMemGCInterval", "WithMemMaxSize",
		"WithRdsDefaultTTL",
	}

	publicMethods := exportedMethodNames(cacheDir)
	expectedMethods := []string{
		"Cache.Clear", "Cache.Close", "Cache.Contains", "Cache.Delete",
		"Cache.ForEach", "Cache.Get", "Cache.GetOrLoad", "Cache.Keys",
		"Cache.Set", "Cache.Size", "Invalidating.Get",
		"Invalidating.Invalidate", "KeyBuilder.Build", "PrefixKeyBuilder.Build",
		"SingleflightMixin.GetOrLoad",
	}

	var failures []string
	failures = append(failures, compareNames("cache symbol", publicSymbols, expectedSymbols)...)
	failures = append(failures, compareNames("cache method", publicMethods, expectedMethods)...)

	exportedFields := exportedFieldNames(cacheDir)
	if len(exportedFields) > 0 {
		failures = append(failures, "cache package should not expose struct fields, found: "+strings.Join(exportedFields, ", "))
	}

	for _, doc := range []corpus{publicIndex, chinesePublicIndex} {
		for _, term := range publicIndexTerms() {
			failures = append(failures, missingTerm(doc, term)...)
		}
	}

	for _, doc := range []corpus{englishDocs, chineseDocs} {
		for _, term := range publicDocSurfaceTerms() {
			failures = append(failures, missingTerm(doc, term)...)
		}
		failures = append(failures, forbids(doc, []string{
			"memoryCache", "redisCache", "globEscaper", "defaultKeyBuilder",
			"cacheKeyPrefix", "lruHandler", "lfuHandler", "fifoHandler",
			"noOpEvictionHandler", "jsonSerializer", "newJSONSerializer",
		})...)
	}

	sourceChecks := []struct {
		path  string
		terms []string
	}{
		{
			path: "cache/cache.go",
			terms: []string{
				"type Cache[T any] interface", "io.Closer",
				"GetOrLoad(ctx context.Context, key string, loader LoaderFunc[T], ttl ...time.Duration) (T, error)",
				"const cacheKeyPrefix = \"vef:cache\"",
				"func NewMemory[T any](opts ...MemoryOption) Cache[T]",
				"func NewRedis[T any](client *redis.Client, namespace string, opts ...RedisOption) Cache[T]",
				"panic(\"redis cache requires a non-nil redis client\")",
				"panic(\"cache.NewRedis requires a non-empty namespace\")",
				"prefix := defaultKeyBuilder.Build(cacheKeyPrefix, namespace)",
			},
		},
		{
			path: "cache/options.go",
			terms: []string{
				"type MemoryOption func(*memoryConfig)", "func WithMemMaxSize(size int64) MemoryOption",
				"func WithMemDefaultTTL(ttl time.Duration) MemoryOption",
				"func WithMemEvictionPolicy(policy EvictionPolicy) MemoryOption",
				"func WithMemGCInterval(interval time.Duration) MemoryOption",
				"type RedisOption func(*redisConfig)",
				"func WithRdsDefaultTTL(ttl time.Duration) RedisOption",
				"cfg.maxSize = size", "cfg.defaultTTL = ttl", "cfg.evictionPolicy = policy",
				"cfg.gcInterval = interval",
			},
		},
		{
			path: "cache/errors.go",
			terms: []string{
				"ErrMemoryLimitExceeded = errors.New(\"memory cache size limit exceeded\")",
				"ErrCacheClosed = errors.New(\"cache closed\")",
				"ErrLoaderRequired = errors.New(\"cache loader is required\")",
				"ErrTypeAssertionFailed = errors.New(\"singleflight: type assertion failed\")",
			},
		},
		{
			path: "cache/eviction.go",
			terms: []string{
				"type EvictionPolicy int", "EvictionPolicyNone EvictionPolicy = iota",
				"EvictionPolicyLRU", "EvictionPolicyLFU", "EvictionPolicyFIFO",
				"FIFO doesn't track access, only insertion order",
				"Return oldest (front of list)",
				"Return the first entry (oldest by insertion order due to FIFO within bucket)",
			},
		},
		{
			path: "cache/memory.go",
			terms: []string{
				"maxSize:        0", "defaultTTL:     0",
				"evictionPolicy: EvictionPolicyLRU", "gcInterval:     5 * time.Minute",
				"if cfg.gcInterval <= 0", "cfg.gcInterval = 5 * time.Minute",
				"if cfg.maxSize <= 0", "cfg.evictionPolicy = EvictionPolicyNone",
				"cfg.evictionPolicy = EvictionPolicyLRU",
				"size            atomic.Int64",
				"for m.size.Load() >= m.maxSize",
				"return ErrMemoryLimitExceeded",
				"if len(ttl) > 0 && ttl[0] > 0",
				"} else if m.defaultTTL > 0",
				"return ErrCacheClosed", "return nil, nil",
				"return 0, nil", "close(m.stopGC)",
			},
		},
		{
			path: "cache/redis.go",
			terms: []string{
				"strings.NewReplacer", "`*`, `\\*`", "`?`, `\\?`",
				"`[`, `\\[`", "`]`, `\\]`",
				"basePrefix: keyBuilder.Build()",
				"serializer: newJSONSerializer[T]()",
				"return globEscaper.Replace(c.basePrefix) + \"*\"",
				"return globEscaper.Replace(c.keyBuilder.Build(prefix)) + \"*\"",
				"func (c *redisCache[T]) stripPrefix(cacheKey string) string",
				"logger.Warnf(\"redis cache deserialize failed for key %s, treating as miss: %v\", cacheKey, err)",
				"return fmt.Errorf(\"redis cache foreach deserialize failed for key %s: %w\", cacheKey, err)",
				"if c.basePrefix == \"\" {",
				"return c.client.FlushDB(ctx).Err()",
				"return c.client.DBSize(ctx).Result()",
				"Close marks the cache as closed. The underlying Redis client remains managed externally.",
				"c.closed.Store(true)",
			},
		},
		{
			path: "cache/serializer.go",
			terms: []string{
				"import \"encoding/json\"", "return json.Marshal(value)",
				"err = json.Unmarshal(data, &value)",
			},
		},
		{
			path: "cache/key_builder.go",
			terms: []string{
				"var defaultKeyBuilder = NewPrefixKeyBuilder(\"\")",
				"func Key(keyParts ...string) string",
				"type KeyBuilder interface", "Build(keyParts ...string) string",
				"type PrefixKeyBuilder struct", "separator string",
				"separator: \":\"", "return strings.Join(keyParts, k.separator)",
				"if len(keyParts) == 0", "return k.prefix",
			},
		},
		{
			path: "cache/singleflight_mixin.go",
			terms: []string{
				"type LoaderFunc[T any] func(ctx context.Context) (T, error)",
				"type GetFunc[T any] func(context.Context, string) (T, bool)",
				"type SetFunc[T any] func(context.Context, string, T, ...time.Duration) error",
				"type SingleflightMixin[T any] struct",
				"if loader == nil", "return value, ErrLoaderRequired",
				"if value, found := getFn(ctx, cacheKey); found",
				"result, err, _ := m.group.Do(cacheKey, func() (any, error)",
				"Double-check: Another goroutine might have loaded it while we waited",
				"if setErr := setFn(ctx, cacheKey, value, ttl...); setErr != nil",
				"return value, ErrTypeAssertionFailed",
			},
		},
		{
			path: "cache/invalidating.go",
			terms: []string{
				"type KeyedLoaderFunc[T any] func(ctx context.Context, key string) (T, error)",
				"type Invalidating[T any] struct",
				"func NewInvalidating[T any](loader KeyedLoaderFunc[T], logger logx.Logger) *Invalidating[T]",
				"cache:  NewMemory[T]()", "loader: loader", "logger: logger",
				"return i.loader(ctx, key)", "if len(keys) == 0",
				"i.cache.Clear(ctx)", "i.cache.Delete(ctx, key)",
				"i.logger.Info(\"Cleared all cache entries\")",
				"i.logger.Infof(\"Cleared cache entry %q\", key)",
			},
		},
	}

	for _, check := range sourceChecks {
		source := readCorpus(check.path, filepath.Join(sourceRoot, check.path))
		failures = append(failures, missingTerms(source, check.terms)...)
	}

	failures = append(failures, missingTerms(englishDocs, []string{
		"no exported fields", "entry count", "`size <= 0` means unlimited",
		"`interval <= 0` falls back to `5m`", "non-durable",
		"process-local", "forces `EvictionPolicyNone`",
		"falls back to `EvictionPolicyLRU`", "`ttl <= 0` does not create an expiration",
		"`Get` and `Contains` remove expired entries lazily", "`Keys` and `ForEach` skip them",
		"panics if `client` is nil or `namespace` is empty",
		"does not close the underlying client", "serialized with JSON",
		"deserialization failures as misses", "`ForEach` returns an error",
		"vef:cache:<namespace>", "strip that internal prefix",
		"`Clear` deletes only keys under the cache namespace",
		"`Size` counts only that namespace", "Redis glob metacharacters",
		"`*`, `?`, `[`, `]`, and `\\`", "`cache.Key()` returns `\"\"`",
		"`cache.NewPrefixKeyBuilder(\"app\").Build()`", "`app:user:123`",
		"checks the cache before joining the singleflight call",
		"checks it again inside the singleflight function",
		"propagates loader errors", "propagates write errors",
		"Pass a non-nil `cache.KeyedLoaderFunc[T]` and a non-nil `logx.Logger`",
		"clears the whole cache when `keys` is empty",
		"deletes exactly the named keys",
		"write path is called after `Close()`",
	})...)
	failures = append(failures, missingTerms(chineseDocs, []string{
		"没有 exported fields", "按条目数量", "`size <= 0` 表示无限制",
		"`interval <= 0` 回退到 `5m`", "不持久化",
		"进程内", "强制使用 `EvictionPolicyNone`",
		"回退到 `EvictionPolicyLRU`", "`ttl <= 0` 本身不会创建过期时间",
		"`Get` 和 `Contains` 会懒删除过期条目", "`Keys` 和 `ForEach` 会跳过它们",
		"`client` 为 nil 或 `namespace` 为空时 panic",
		"不会关闭底层 client", "使用 JSON 序列化",
		"反序列化失败都当成 miss", "`ForEach` 在无法读取或反序列化",
		"vef:cache:<namespace>", "剥离这个内部 prefix",
		"`Clear` 只删除该 cache namespace 下的 key",
		"`Size` 也只统计该 namespace", "Redis glob 元字符",
		"`*`、`?`、`[`、`]` 和 `\\`", "`cache.Key()` 返回 `\"\"`",
		"`cache.NewPrefixKeyBuilder(\"app\").Build()`", "`app:user:123`",
		"先在加入 singleflight 前读一次 cache",
		"singleflight 函数内部读一次", "loader error 和写入 error 都会原样返回",
		"非 nil 的 `cache.KeyedLoaderFunc[T]` 和非 nil 的 `logx.Logger`",
		"`keys` 为空时清空整个缓存", "否则只删除指定 key",
		"写入路径在 `Close()` 之后被调用",
	})...)

	sort.Strings(failures)
	if len(failures) > 0 {
		panic(fmt.Errorf("cache contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	fmt.Printf("Cache contract docs verified: %d public symbols, %d public methods, %d source files, 2 doc mirrors\n", len(publicSymbols), len(publicMethods), len(sourceChecks))
}

func publicDocSurfaceTerms() []string {
	return []string{
		"`cache.Cache[T]`",
		"`cache.NewMemory[T](opts ...cache.MemoryOption) cache.Cache[T]`",
		"`cache.NewRedis[T](client *redis.Client, namespace string, opts ...cache.RedisOption) cache.Cache[T]`",
		"`cache.MemoryOption`", "`cache.RedisOption`",
		"`cache.WithMemMaxSize(size int64)`",
		"`cache.WithMemDefaultTTL(ttl time.Duration)`",
		"`cache.WithMemEvictionPolicy(policy cache.EvictionPolicy)`",
		"`cache.WithMemGCInterval(interval time.Duration)`",
		"`cache.WithRdsDefaultTTL(ttl time.Duration)`",
		"`cache.EvictionPolicy`", "`cache.EvictionPolicyNone`",
		"`cache.EvictionPolicyLRU`", "`cache.EvictionPolicyLFU`",
		"`cache.EvictionPolicyFIFO`", "`cache.Key(keyParts ...string) string`",
		"`cache.KeyBuilder`", "`cache.PrefixKeyBuilder`",
		"`cache.NewPrefixKeyBuilder(prefix string) *cache.PrefixKeyBuilder`",
		"`cache.LoaderFunc[T]`", "`func(ctx context.Context) (T, error)`",
		"`cache.KeyedLoaderFunc[T]`", "`func(ctx context.Context, key string) (T, error)`",
		"`cache.GetFunc[T]`", "`func(context.Context, string) (T, bool)`",
		"`cache.SetFunc[T]`", "`func(context.Context, string, T, ...time.Duration) error`",
		"`cache.SingleflightMixin[T]`", "`cache.Invalidating[T]`",
		"`cache.NewInvalidating[T](loader cache.KeyedLoaderFunc[T], logger logx.Logger) *cache.Invalidating[T]`",
		"`cache.ErrMemoryLimitExceeded`", "`cache.ErrCacheClosed`",
		"`cache.ErrLoaderRequired`", "`cache.ErrTypeAssertionFailed`",
		"`Get` | `Get(ctx context.Context, key string) (T, bool)`",
		"`GetOrLoad` | `GetOrLoad(ctx context.Context, key string, loader cache.LoaderFunc[T], ttl ...time.Duration) (T, error)`",
		"`Set` | `Set(ctx context.Context, key string, value T, ttl ...time.Duration) error`",
		"`Contains` | `Contains(ctx context.Context, key string) bool`",
		"`Delete` | `Delete(ctx context.Context, key string) error`",
		"`Clear` | `Clear(ctx context.Context) error`",
		"`Keys` | `Keys(ctx context.Context, prefix ...string) ([]string, error)`",
		"`ForEach` | `ForEach(ctx context.Context, callback func(key string, value T) bool, prefix ...string) error`",
		"`Size` | `Size(ctx context.Context) (int64, error)`",
		"`Close` | `Close() error`",
		"Build(keyParts ...string) string",
		"func (m *SingleflightMixin[T]) GetOrLoad(",
		"`Get` | `Get(ctx context.Context, key string) (T, error)`",
		"`Invalidate` | `Invalidate(ctx context.Context, keys ...string) error`",
	}
}

func publicIndexTerms() []string {
	return []string{
		"TYPE Cache : github.com/coldsmirk/vef-framework-go/cache.Cache[T any]",
		"  METHOD Clear : func(ctx context.Context) error",
		"  METHOD Close : func() error",
		"  METHOD Contains : func(ctx context.Context, key string) bool",
		"  METHOD Delete : func(ctx context.Context, key string) error",
		"  METHOD ForEach : func(ctx context.Context, callback func(key string, value T) bool, prefix ...string) error",
		"  METHOD Get : func(ctx context.Context, key string) (T, bool)",
		"  METHOD GetOrLoad : func(ctx context.Context, key string, loader github.com/coldsmirk/vef-framework-go/cache.LoaderFunc[T], ttl ...time.Duration) (T, error)",
		"  METHOD Keys : func(ctx context.Context, prefix ...string) ([]string, error)",
		"  METHOD Set : func(ctx context.Context, key string, value T, ttl ...time.Duration) error",
		"  METHOD Size : func(ctx context.Context) (int64, error)",
		"VAR ErrCacheClosed : error",
		"VAR ErrLoaderRequired : error",
		"VAR ErrMemoryLimitExceeded : error",
		"VAR ErrTypeAssertionFailed : error",
		"TYPE EvictionPolicy : github.com/coldsmirk/vef-framework-go/cache.EvictionPolicy",
		"CONST EvictionPolicyFIFO : github.com/coldsmirk/vef-framework-go/cache.EvictionPolicy = 3",
		"CONST EvictionPolicyLFU : github.com/coldsmirk/vef-framework-go/cache.EvictionPolicy = 2",
		"CONST EvictionPolicyLRU : github.com/coldsmirk/vef-framework-go/cache.EvictionPolicy = 1",
		"CONST EvictionPolicyNone : github.com/coldsmirk/vef-framework-go/cache.EvictionPolicy = 0",
		"TYPE GetFunc : github.com/coldsmirk/vef-framework-go/cache.GetFunc[T any]",
		"TYPE Invalidating : github.com/coldsmirk/vef-framework-go/cache.Invalidating[T any]",
		"  METHOD Get : func(ctx context.Context, key string) (T, error)",
		"  METHOD Invalidate : func(ctx context.Context, keys ...string) error",
		"FUNC Key : func(keyParts ...string) string",
		"TYPE KeyBuilder : github.com/coldsmirk/vef-framework-go/cache.KeyBuilder",
		"  METHOD Build : func(keyParts ...string) string",
		"TYPE KeyedLoaderFunc : github.com/coldsmirk/vef-framework-go/cache.KeyedLoaderFunc[T any]",
		"TYPE LoaderFunc : github.com/coldsmirk/vef-framework-go/cache.LoaderFunc[T any]",
		"TYPE MemoryOption : github.com/coldsmirk/vef-framework-go/cache.MemoryOption",
		"FUNC NewInvalidating : func[T any](loader github.com/coldsmirk/vef-framework-go/cache.KeyedLoaderFunc[T], logger github.com/coldsmirk/vef-framework-go/logx.Logger) *github.com/coldsmirk/vef-framework-go/cache.Invalidating[T]",
		"FUNC NewMemory : func[T any](opts ...github.com/coldsmirk/vef-framework-go/cache.MemoryOption) github.com/coldsmirk/vef-framework-go/cache.Cache[T]",
		"FUNC NewPrefixKeyBuilder : func(prefix string) *github.com/coldsmirk/vef-framework-go/cache.PrefixKeyBuilder",
		"FUNC NewRedis : func[T any](client *github.com/redis/go-redis/v9.Client, namespace string, opts ...github.com/coldsmirk/vef-framework-go/cache.RedisOption) github.com/coldsmirk/vef-framework-go/cache.Cache[T]",
		"TYPE PrefixKeyBuilder : github.com/coldsmirk/vef-framework-go/cache.PrefixKeyBuilder",
		"TYPE RedisOption : github.com/coldsmirk/vef-framework-go/cache.RedisOption",
		"TYPE SetFunc : github.com/coldsmirk/vef-framework-go/cache.SetFunc[T any]",
		"TYPE SingleflightMixin : github.com/coldsmirk/vef-framework-go/cache.SingleflightMixin[T any]",
		"  METHOD GetOrLoad : func(ctx context.Context, cacheKey string, loader github.com/coldsmirk/vef-framework-go/cache.LoaderFunc[T], ttl []time.Duration, getFn github.com/coldsmirk/vef-framework-go/cache.GetFunc[T], setFn github.com/coldsmirk/vef-framework-go/cache.SetFunc[T]) (value T, _ error)",
		"FUNC WithMemDefaultTTL : func(ttl time.Duration) github.com/coldsmirk/vef-framework-go/cache.MemoryOption",
		"FUNC WithMemEvictionPolicy : func(policy github.com/coldsmirk/vef-framework-go/cache.EvictionPolicy) github.com/coldsmirk/vef-framework-go/cache.MemoryOption",
		"FUNC WithMemGCInterval : func(interval time.Duration) github.com/coldsmirk/vef-framework-go/cache.MemoryOption",
		"FUNC WithMemMaxSize : func(size int64) github.com/coldsmirk/vef-framework-go/cache.MemoryOption",
		"FUNC WithRdsDefaultTTL : func(ttl time.Duration) github.com/coldsmirk/vef-framework-go/cache.RedisOption",
	}
}

func exportedTopLevelNames(dir string) []string {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
	}, 0)
	if err != nil {
		panic(fmt.Errorf("failed to parse cache package: %w", err))
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

	return sortedKeys(names)
}

func exportedMethodNames(dir string) []string {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
	}, 0)
	if err != nil {
		panic(fmt.Errorf("failed to parse cache package: %w", err))
	}

	methods := make(map[string]bool)
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				switch d := decl.(type) {
				case *ast.FuncDecl:
					if d.Recv == nil || !d.Name.IsExported() {
						continue
					}
					recv := receiverBaseName(d.Recv.List[0].Type)
					if ast.IsExported(recv) {
						methods[recv+"."+d.Name.Name] = true
					}
				case *ast.GenDecl:
					for _, spec := range d.Specs {
						typeSpec, ok := spec.(*ast.TypeSpec)
						if !ok || !typeSpec.Name.IsExported() {
							continue
						}
						if iface, ok := typeSpec.Type.(*ast.InterfaceType); ok {
							for _, field := range iface.Methods.List {
								if len(field.Names) == 0 {
									if isEmbeddedIOCloser(field.Type) {
										methods[typeSpec.Name.Name+".Close"] = true
									}
									continue
								}
								for _, name := range field.Names {
									if name.IsExported() {
										methods[typeSpec.Name.Name+"."+name.Name] = true
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return sortedKeys(methods)
}

func exportedFieldNames(dir string) []string {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
	}, 0)
	if err != nil {
		panic(fmt.Errorf("failed to parse cache package: %w", err))
	}

	fields := make(map[string]bool)
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}
				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok || !typeSpec.Name.IsExported() {
						continue
					}
					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}
					for _, field := range structType.Fields.List {
						for _, name := range field.Names {
							if name.IsExported() {
								fields[typeSpec.Name.Name+"."+name.Name] = true
							}
						}
					}
				}
			}
		}
	}

	return sortedKeys(fields)
}

func receiverBaseName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return receiverBaseName(t.X)
	case *ast.Ident:
		return t.Name
	case *ast.IndexExpr:
		return receiverBaseName(t.X)
	case *ast.IndexListExpr:
		return receiverBaseName(t.X)
	case *ast.SelectorExpr:
		return t.Sel.Name
	default:
		return ""
	}
}

func isEmbeddedIOCloser(expr ast.Expr) bool {
	selector, ok := expr.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "Closer" {
		return false
	}

	ident, ok := selector.X.(*ast.Ident)
	return ok && ident.Name == "io"
}

func compareNames(label string, current, expected []string) []string {
	currentSet := set(current)
	expectedSet := set(expected)
	var failures []string
	for _, name := range current {
		if !expectedSet[name] {
			failures = append(failures, "unexpected public "+label+": "+name)
		}
	}
	for _, name := range expected {
		if !currentSet[name] {
			failures = append(failures, "missing expected public "+label+": "+name)
		}
	}

	return failures
}

func sortedKeys(values map[string]bool) []string {
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)

	return result
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
