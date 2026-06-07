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

	englishDocs := readCorpus("English cron docs", filepath.Join(docsRoot, "docs/features/cron.md"))
	chineseDocs := readCorpus("Chinese cron docs", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/features/cron.md"))
	publicIndex := readCorpus("English public API index", filepath.Join(docsRoot, "docs/reference/public-api-index.md"))
	chinesePublicIndex := readCorpus("Chinese public API index", filepath.Join(docsRoot, "i18n/zh-Hans/docusaurus-plugin-content-docs/current/reference/public-api-index.md"))

	cronDir := filepath.Join(sourceRoot, "cron")
	publicSymbols := exportedTopLevelNames(cronDir)
	expectedSymbols := []string{
		"CronJobDefinition", "DurationJobDefinition", "DurationRandomJobDefinition",
		"ErrJobNameRequired", "ErrJobTaskHandlerMustFunc", "ErrJobTaskHandlerRequired",
		"Job", "JobDefinition", "JobDescriptorOption", "NewCronJob",
		"NewDurationJob", "NewDurationRandomJob", "NewOneTimeJob",
		"NewScheduler", "OneTimeJobDefinition", "Scheduler", "WithConcurrent",
		"WithContext", "WithLimitedRuns", "WithName", "WithStartAt",
		"WithStartImmediately", "WithStopAt", "WithTags", "WithTask",
	}

	publicMethods := exportedMethodNames(cronDir)
	expectedMethods := []string{
		"Job.ID", "Job.LastRun", "Job.Name", "Job.NextRun", "Job.NextRuns",
		"Job.RunNow", "Job.Tags", "Scheduler.Jobs", "Scheduler.JobsWaitingInQueue",
		"Scheduler.NewJob", "Scheduler.RemoveByTags", "Scheduler.RemoveJob",
		"Scheduler.Start", "Scheduler.StopJobs", "Scheduler.Update",
	}

	var failures []string
	failures = append(failures, compareNames("cron symbol", publicSymbols, expectedSymbols)...)
	failures = append(failures, compareNames("cron method", publicMethods, expectedMethods)...)

	exportedFields := exportedFieldNames(cronDir)
	if len(exportedFields) > 0 {
		failures = append(failures, "cron package should not expose struct fields, found: "+strings.Join(exportedFields, ", "))
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
			"schedulerAdapter", "jobAdapter", "jobInfo", "jobDescriptor",
			"jobTask", "newScheduler", "cronLogger", "jobMonitor",
		})...)
	}

	sourceChecks := []struct {
		path  string
		terms []string
	}{
		{
			path: "cron/scheduler.go",
			terms: []string{
				"type Scheduler interface", "Jobs() []Job",
				"NewJob(definition JobDefinition) (Job, error)",
				"RemoveByTags(tags ...string)", "RemoveJob(id string) error",
				"Start()", "StopJobs() error",
				"Update(id string, definition JobDefinition) (Job, error)",
				"JobsWaitingInQueue() int", "definition.build()",
				"uuid.Parse(id)", "s.scheduler.RemoveJob(uuid)",
				"s.scheduler.Update(uuid, def, task, options...)",
			},
		},
		{
			path: "cron/job.go",
			terms: []string{
				"type Job interface", "ID() string",
				"LastRun() (time.Time, error)", "Name() string",
				"NextRun() (time.Time, error)", "NextRuns(count int) ([]time.Time, error)",
				"RunNow() error", "Tags() []string",
				"j.job.LastRunStartedAt()", "j.job.RunNow()",
			},
		},
		{
			path: "cron/job_descriptor.go",
			terms: []string{
				"if i.name == \"\"", "return nil, ErrJobNameRequired",
				"uuid.NewRandom()", "gocron.WithIdentifier(id)",
				"gocron.WithName(i.name)", "gocron.WithTags(i.tags...)",
				"if !i.allowConcurrent",
				"gocron.WithSingletonMode(gocron.LimitModeWait)",
				"if !i.startAt.IsZero()", "else if i.startImmediately",
				"if !i.stopAt.IsZero()", "if i.limitedRuns > 0",
				"if i.ctx != nil", "if t.handler == nil",
				"return nil, ErrJobTaskHandlerRequired",
				"reflect.ValueOf(t.handler).Kind() != reflect.Func",
				"return nil, ErrJobTaskHandlerMustFunc",
				"gocron.NewTask(t.handler, t.params...)",
				"task, err := d.buildTask()", "options, err := d.buildJobOptions()",
				"failed to build job task: %w", "failed to build job options: %w",
			},
		},
		{
			path: "cron/job_definitions.go",
			terms: []string{
				"type OneTimeJobDefinition struct", "type DurationJobDefinition struct",
				"type DurationRandomJobDefinition struct", "type CronJobDefinition struct",
				"case 0:", "gocron.OneTimeJobStartImmediately()",
				"case 1:", "gocron.OneTimeJobStartDateTime(d.times[0])",
				"gocron.OneTimeJobStartDateTimes(d.times...)",
				"gocron.DurationJob(d.interval)",
				"gocron.DurationRandomJob(d.minInterval, d.maxInterval)",
				"gocron.CronJob(d.expression, d.withSeconds)",
				"func NewOneTimeJob(times []time.Time, options ...JobDescriptorOption) *OneTimeJobDefinition",
				"func NewDurationJob(interval time.Duration, options ...JobDescriptorOption) *DurationJobDefinition",
				"func NewDurationRandomJob(minInterval, maxInterval time.Duration, options ...JobDescriptorOption) *DurationRandomJobDefinition",
				"func NewCronJob(expression string, withSeconds bool, options ...JobDescriptorOption) *CronJobDefinition",
			},
		},
		{
			path: "cron/options.go",
			terms: []string{
				"type JobDescriptorOption func(*jobDescriptor)",
				"func WithName(name string) JobDescriptorOption",
				"func WithTags(tags ...string) JobDescriptorOption",
				"func WithConcurrent() JobDescriptorOption",
				"func WithStartAt(startAt time.Time) JobDescriptorOption",
				"func WithStartImmediately() JobDescriptorOption",
				"func WithStopAt(stopAt time.Time) JobDescriptorOption",
				"func WithLimitedRuns(limitedRuns uint) JobDescriptorOption",
				"func WithContext(ctx context.Context) JobDescriptorOption",
				"func WithTask(handler any, params ...any) JobDescriptorOption",
				"d.name = name", "d.tags = tags", "d.allowConcurrent = true",
				"d.startAt = startAt", "d.startImmediately = true",
				"d.stopAt = stopAt", "d.limitedRuns = limitedRuns",
				"d.ctx = ctx", "d.handler = handler", "d.params = params",
			},
		},
		{
			path: "cron/errors.go",
			terms: []string{
				"ErrJobNameRequired = errors.New(\"job name is required\")",
				"ErrJobTaskHandlerRequired = errors.New(\"job task handler is required\")",
				"ErrJobTaskHandlerMustFunc = errors.New(\"job task handler must be a function\")",
			},
		},
		{
			path: "internal/cron/scheduler.go",
			terms: []string{
				"gocron.WithLocation(time.Local)",
				"gocron.WithStopTimeout(30*time.Second)",
				"gocron.WithLogger(newCronLogger())",
				"gocron.WithMonitorStatus(newJobMonitor())",
				"gocron.WithLimitConcurrentJobs(1000, gocron.LimitModeWait)",
				"scheduler.Start()", "scheduler.Shutdown()",
				"Cron scheduler started", "Cron scheduler stopped",
			},
		},
		{
			path: "internal/cron/module.go",
			terms: []string{
				"fx.Module(", "\"vef:cron\"",
				"fx.Provide(newScheduler, fx.Private)",
				"fx.Provide(cron.NewScheduler)",
			},
		},
		{
			path: "internal/cron/logger.go",
			terms: []string{
				"func (*cronLogger) Debug", "func (*cronLogger) Error",
				"func (*cronLogger) Info", "func (*cronLogger) Warn",
			},
		},
		{
			path: "internal/cron/monitor.go",
			terms: []string{
				"RecordJobTimingWithStatus", "gocron.Success",
				"gocron.Fail", "IncrementJob", "RecordJobTiming",
			},
		},
	}

	for _, check := range sourceChecks {
		source := readCorpus(check.path, filepath.Join(sourceRoot, check.path))
		failures = append(failures, missingTerms(source, check.terms)...)
	}

	failures = append(failures, missingTerms(englishDocs, []string{
		"no exported fields", "caller-provided `gocron.Scheduler`",
		"caller must pass a usable scheduler", "with seconds",
		"standard 5-field form", "singleton wait mode",
		"takes precedence over `WithStartImmediately`",
		"`limitedRuns > 0`", "forwards `params` to it through `gocron.NewTask`",
		"invalid non-UUID strings return a parse error",
		"Use IDs returned by `Job.ID()`", "last run start time",
		"respecting job/scheduler limits and run limits",
		"task handler takes precedence over a missing name",
		"wrapped by the build path", "Use `errors.Is",
		"local time zone", "stop timeout `30s`",
		"scheduler logger and monitor", "concurrent job limit `1000` with wait mode",
		"app stop shuts the scheduler down",
	})...)
	failures = append(failures, missingTerms(chineseDocs, []string{
		"没有 exported fields", "调用方提供的 `gocron.Scheduler`",
		"调用方必须传入可用 scheduler", "带 seconds 字段",
		"标准 5-field 格式", "singleton wait mode",
		"优先级高于 `WithStartImmediately`",
		"`limitedRuns > 0`", "通过 `gocron.NewTask` 转发 `params`",
		"非 UUID 字符串会在 delegation 前返回 parse error",
		"应使用 `Job.ID()` 返回的 ID", "run 的开始时间",
		"遵守 job/scheduler 限制和运行次数限制",
		"先看到 task handler 相关错误", "会被 build path 包裹",
		"使用 `errors.Is", "本地时区", "停止超时 `30s`",
		"scheduler logger 和 monitor", "并发 job 上限 `1000` 且使用 wait mode",
		"应用停止时 shutdown scheduler",
	})...)

	failures = append(failures, runRuntimeChecks(sourceRoot)...)

	sort.Strings(failures)
	if len(failures) > 0 {
		panic(fmt.Errorf("cron contract verification failed:\n%s", strings.Join(failures, "\n")))
	}

	fmt.Printf("Cron contract docs verified: %d public symbols, %d public methods, %d source files, 2 doc mirrors\n", len(publicSymbols), len(publicMethods), len(sourceChecks))
}

func publicDocSurfaceTerms() []string {
	return []string{
		"`cron.Scheduler`",
		"`cron.NewScheduler(scheduler gocron.Scheduler) cron.Scheduler`",
		"`cron.Job`",
		"`cron.JobDefinition`",
		"`cron.JobDescriptorOption`",
		"`cron.OneTimeJobDefinition`",
		"`cron.DurationJobDefinition`",
		"`cron.DurationRandomJobDefinition`",
		"`cron.CronJobDefinition`",
		"`cron.NewOneTimeJob(times []time.Time, options ...cron.JobDescriptorOption) *cron.OneTimeJobDefinition`",
		"`cron.NewDurationJob(interval time.Duration, options ...cron.JobDescriptorOption) *cron.DurationJobDefinition`",
		"`cron.NewDurationRandomJob(minInterval time.Duration, maxInterval time.Duration, options ...cron.JobDescriptorOption) *cron.DurationRandomJobDefinition`",
		"`cron.NewCronJob(expression string, withSeconds bool, options ...cron.JobDescriptorOption) *cron.CronJobDefinition`",
		"`cron.WithName(name string)`",
		"`cron.WithTags(tags ...string)`",
		"`cron.WithConcurrent()`",
		"`cron.WithStartAt(startAt time.Time)`",
		"`cron.WithStartImmediately()`",
		"`cron.WithStopAt(stopAt time.Time)`",
		"`cron.WithLimitedRuns(limitedRuns uint)`",
		"`cron.WithContext(ctx context.Context)`",
		"`cron.WithTask(handler any, params ...any)`",
		"`cron.ErrJobNameRequired`",
		"`cron.ErrJobTaskHandlerRequired`",
		"`cron.ErrJobTaskHandlerMustFunc`",
		"`Jobs` | `Jobs() []cron.Job`",
		"`NewJob` | `NewJob(definition cron.JobDefinition) (cron.Job, error)`",
		"`RemoveByTags` | `RemoveByTags(tags ...string)`",
		"`RemoveJob` | `RemoveJob(id string) error`",
		"`Start` | `Start()`",
		"`StopJobs` | `StopJobs() error`",
		"`Update` | `Update(id string, definition cron.JobDefinition) (cron.Job, error)`",
		"`JobsWaitingInQueue` | `JobsWaitingInQueue() int`",
		"`ID` | `ID() string`",
		"`LastRun` | `LastRun() (time.Time, error)`",
		"`Name` | `Name() string`",
		"`NextRun` | `NextRun() (time.Time, error)`",
		"`NextRuns` | `NextRuns(count int) ([]time.Time, error)`",
		"`RunNow` | `RunNow() error`",
		"`Tags` | `Tags() []string`",
	}
}

func publicIndexTerms() []string {
	return []string{
		"TYPE CronJobDefinition : github.com/coldsmirk/vef-framework-go/cron.CronJobDefinition",
		"TYPE DurationJobDefinition : github.com/coldsmirk/vef-framework-go/cron.DurationJobDefinition",
		"TYPE DurationRandomJobDefinition : github.com/coldsmirk/vef-framework-go/cron.DurationRandomJobDefinition",
		"VAR ErrJobNameRequired : error",
		"VAR ErrJobTaskHandlerMustFunc : error",
		"VAR ErrJobTaskHandlerRequired : error",
		"TYPE Job : github.com/coldsmirk/vef-framework-go/cron.Job",
		"  METHOD ID : func() string",
		"  METHOD LastRun : func() (time.Time, error)",
		"  METHOD Name : func() string",
		"  METHOD NextRun : func() (time.Time, error)",
		"  METHOD NextRuns : func(count int) ([]time.Time, error)",
		"  METHOD RunNow : func() error",
		"  METHOD Tags : func() []string",
		"TYPE JobDefinition : github.com/coldsmirk/vef-framework-go/cron.JobDefinition",
		"TYPE JobDescriptorOption : github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption",
		"FUNC NewCronJob : func(expression string, withSeconds bool, options ...github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption) *github.com/coldsmirk/vef-framework-go/cron.CronJobDefinition",
		"FUNC NewDurationJob : func(interval time.Duration, options ...github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption) *github.com/coldsmirk/vef-framework-go/cron.DurationJobDefinition",
		"FUNC NewDurationRandomJob : func(minInterval time.Duration, maxInterval time.Duration, options ...github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption) *github.com/coldsmirk/vef-framework-go/cron.DurationRandomJobDefinition",
		"FUNC NewOneTimeJob : func(times []time.Time, options ...github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption) *github.com/coldsmirk/vef-framework-go/cron.OneTimeJobDefinition",
		"FUNC NewScheduler : func(scheduler github.com/go-co-op/gocron/v2.Scheduler) github.com/coldsmirk/vef-framework-go/cron.Scheduler",
		"TYPE OneTimeJobDefinition : github.com/coldsmirk/vef-framework-go/cron.OneTimeJobDefinition",
		"TYPE Scheduler : github.com/coldsmirk/vef-framework-go/cron.Scheduler",
		"  METHOD Jobs : func() []github.com/coldsmirk/vef-framework-go/cron.Job",
		"  METHOD JobsWaitingInQueue : func() int",
		"  METHOD NewJob : func(definition github.com/coldsmirk/vef-framework-go/cron.JobDefinition) (github.com/coldsmirk/vef-framework-go/cron.Job, error)",
		"  METHOD RemoveByTags : func(tags ...string)",
		"  METHOD RemoveJob : func(id string) error",
		"  METHOD Start : func()",
		"  METHOD StopJobs : func() error",
		"  METHOD Update : func(id string, definition github.com/coldsmirk/vef-framework-go/cron.JobDefinition) (github.com/coldsmirk/vef-framework-go/cron.Job, error)",
		"FUNC WithConcurrent : func() github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption",
		"FUNC WithContext : func(ctx context.Context) github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption",
		"FUNC WithLimitedRuns : func(limitedRuns uint) github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption",
		"FUNC WithName : func(name string) github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption",
		"FUNC WithStartAt : func(startAt time.Time) github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption",
		"FUNC WithStartImmediately : func() github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption",
		"FUNC WithStopAt : func(stopAt time.Time) github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption",
		"FUNC WithTags : func(tags ...string) github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption",
		"FUNC WithTask : func(handler any, params ...any) github.com/coldsmirk/vef-framework-go/cron.JobDescriptorOption",
	}
}

func runRuntimeChecks(sourceRoot string) []string {
	tmpDir, err := os.MkdirTemp("", "verify-cron-contracts-*")
	if err != nil {
		return []string{fmt.Sprintf("failed to create temp module: %v", err)}
	}
	defer os.RemoveAll(tmpDir)

	goMod := fmt.Sprintf(`module verifycroncontracts

go 1.25.0

require github.com/coldsmirk/vef-framework-go v0.0.0

replace github.com/coldsmirk/vef-framework-go => %s
`, sourceRoot)
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o600); err != nil {
		return []string{fmt.Sprintf("failed to write temp go.mod: %v", err)}
	}

	mainGo := `package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"

	"github.com/coldsmirk/vef-framework-go/cron"
)

func main() {
	raw, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = raw.Shutdown()
	}()

	scheduler := cron.NewScheduler(raw)

	if _, err := scheduler.NewJob(cron.NewDurationJob(time.Minute)); !errors.Is(err, cron.ErrJobTaskHandlerRequired) {
		panic(fmt.Sprintf("missing task should wrap ErrJobTaskHandlerRequired, got %v", err))
	}

	if _, err := scheduler.NewJob(cron.NewDurationJob(time.Minute, cron.WithName("bad"), cron.WithTask("not-a-func"))); !errors.Is(err, cron.ErrJobTaskHandlerMustFunc) {
		panic(fmt.Sprintf("non-function task should wrap ErrJobTaskHandlerMustFunc, got %v", err))
	}

	if _, err := scheduler.NewJob(cron.NewDurationJob(time.Minute, cron.WithTask(func() {}))); !errors.Is(err, cron.ErrJobNameRequired) {
		panic(fmt.Sprintf("missing name should wrap ErrJobNameRequired after task build succeeds, got %v", err))
	}

	if _, err := scheduler.NewJob(cron.NewDurationJob(time.Minute, cron.WithTask("not-a-func"))); !errors.Is(err, cron.ErrJobTaskHandlerMustFunc) || errors.Is(err, cron.ErrJobNameRequired) {
		panic(fmt.Sprintf("task validation must precede missing-name validation, got %v", err))
	}

	if err := scheduler.RemoveJob("not-a-uuid"); err == nil {
		panic("RemoveJob should reject invalid UUID strings")
	}

	if _, err := scheduler.Update("not-a-uuid", cron.NewDurationJob(time.Minute, cron.WithName("valid"), cron.WithTask(func() {}))); err == nil {
		panic("Update should reject invalid UUID strings")
	}

	if _, err := scheduler.NewJob(cron.NewDurationJob(time.Minute, cron.WithName("params"), cron.WithTask(func(value string) {}, "ok"))); err != nil {
		panic(fmt.Sprintf("WithTask should forward params through gocron.NewTask, got %v", err))
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(mainGo), 0o600); err != nil {
		return []string{fmt.Sprintf("failed to write temp main.go: %v", err)}
	}

	cmd := exec.Command("go", "run", "-mod=mod", ".")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{fmt.Sprintf("cron runtime contract check failed: %v\n%s", err, strings.TrimSpace(string(output)))}
	}

	return nil
}

func exportedTopLevelNames(dir string) []string {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
	}, 0)
	if err != nil {
		panic(fmt.Errorf("failed to parse cron package: %w", err))
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
		panic(fmt.Errorf("failed to parse cron package: %w", err))
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
		panic(fmt.Errorf("failed to parse cron package: %w", err))
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
