import type {ReactNode} from 'react';
import Link from '@docusaurus/Link';
import Translate, {translate} from '@docusaurus/Translate';
import useBaseUrl from '@docusaurus/useBaseUrl';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import Heading from '@theme/Heading';

type StatItem = {
  label: string;
  value: string;
};

type WorkflowItem = {
  index: string;
  title: string;
  description: string;
};

type FeatureItem = {
  badge: string;
  title: string;
  description: string;
};

type PathItem = {
  title: string;
  description: string;
  to: string;
};

const heroCode = [
  'func main() {',
  '  vef.Run(',
  '    app.Module,',
  '    vef.ProvideAPIResource(resources.NewUserResource),',
  '  )',
  '}',
];

const statList: StatItem[] = [
  {
    label: translate({id: 'homepage.stats.endpoint.label', message: 'RPC endpoint'}),
    value: '/api',
  },
  {
    label: translate({id: 'homepage.stats.version.label', message: 'Default version'}),
    value: 'v1',
  },
  {
    label: translate({id: 'homepage.stats.timeout.label', message: 'Default timeout'}),
    value: '30s',
  },
  {
    label: translate({id: 'homepage.stats.rate.label', message: 'Default rate limit'}),
    value: '100 / 5m',
  },
];

const workflowList: WorkflowItem[] = [
  {
    index: '01',
    title: translate({id: 'homepage.workflow.one.title', message: 'Compose modules'}),
    description: translate({
      id: 'homepage.workflow.one.description',
      message:
        'Start with vef.Run(...) and let FX assemble config, database, ORM, middleware, security, storage, MCP, and the application server.',
    }),
  },
  {
    index: '02',
    title: translate({id: 'homepage.workflow.two.title', message: 'Register resources'}),
    description: translate({
      id: 'homepage.workflow.two.description',
      message:
        'Expose APIs with api.NewRPCResource(...) or api.NewRESTResource(...), then register them through vef.ProvideAPIResource(...).',
    }),
  },
  {
    index: '03',
    title: translate({id: 'homepage.workflow.three.title', message: 'Keep handlers lean'}),
    description: translate({
      id: 'homepage.workflow.three.description',
      message:
        'Inject fiber.Ctx, orm.DB, Principal, Logger, Params, Meta, Storage, Event, and Cron directly into handlers instead of wiring glue code by hand.',
    }),
  },
];

const featureList: FeatureItem[] = [
  {
    badge: 'API',
    title: translate({id: 'homepage.feature.api.title', message: 'Unified resource model'}),
    description: translate({
      id: 'homepage.feature.api.description',
      message:
        'RPC and REST are first-class resource types. Each operation carries auth, timeout, rate limit, audit, and versioning behavior.',
    }),
  },
  {
    badge: 'CRUD',
    title: translate({id: 'homepage.feature.crud.title', message: 'Generic CRUD builders'}),
    description: translate({
      id: 'homepage.feature.crud.description',
      message:
        'Create, update, delete, paging, tree queries, import/export, batch operations, and hooks are composed from typed builders.',
    }),
  },
  {
    badge: 'ORM',
    title: translate({id: 'homepage.feature.orm.title', message: 'Model-driven data layer'}),
    description: translate({
      id: 'homepage.feature.orm.description',
      message:
        'Build on Bun with audit models, transactional helpers, search tags, pagination, and source-backed defaults for enterprise data access.',
    }),
  },
  {
    badge: 'AUTH',
    title: translate({id: 'homepage.feature.auth.title', message: 'Security and data scopes'}),
    description: translate({
      id: 'homepage.feature.auth.description',
      message:
        'Bearer, Signature, and public endpoints work with RBAC permission checks and request-scoped data permissions.',
    }),
  },
  {
    badge: 'SYS',
    title: translate({id: 'homepage.feature.builtin.title', message: 'Built-in system resources'}),
    description: translate({
      id: 'homepage.feature.builtin.description',
      message:
        'Authentication, storage, schema inspection, monitoring, and MCP are available as ready-to-use resources and middleware.',
    }),
  },
  {
    badge: 'FX',
    title: translate({id: 'homepage.feature.modular.title', message: 'Modular runtime'}),
    description: translate({
      id: 'homepage.feature.modular.description',
      message:
        'Extension points are grouped through FX and kept explicit: API resources, app middleware, CQRS behaviors, handler resolvers, and MCP providers.',
    }),
  },
];

const pathList: PathItem[] = [
  {
    title: translate({id: 'homepage.paths.install.title', message: 'Install and configure'}),
    description: translate({
      id: 'homepage.paths.install.description',
      message: 'Set up the project, application.toml, and the minimal runtime dependencies.',
    }),
    to: '/docs/getting-started/installation',
  },
  {
    title: translate({id: 'homepage.paths.quick.title', message: 'Build your first resource'}),
    description: translate({
      id: 'homepage.paths.quick.description',
      message: 'Walk through the smallest app that actually serves an API endpoint.',
    }),
    to: '/docs/getting-started/quick-start',
  },
  {
    title: translate({id: 'homepage.paths.routing.title', message: 'Understand routing'}),
    description: translate({
      id: 'homepage.paths.routing.description',
      message: 'Learn how RPC requests, REST routes, params, meta, auth, and handlers fit together.',
    }),
    to: '/docs/guide/routing',
  },
  {
    title: translate({id: 'homepage.paths.crud.title', message: 'Scale with CRUD builders'}),
    description: translate({
      id: 'homepage.paths.crud.description',
      message: 'Move from hand-written handlers to generic operations, hooks, search, and pagination.',
    }),
    to: '/docs/guide/crud',
  },
];

function HeroCodeBlock() {
  return (
    <div className="home-panel hero-code-panel">
      <div className="home-panel-header">
        <span className="home-panel-dot" />
        <span className="home-panel-dot" />
        <span className="home-panel-dot" />
        <span className="home-panel-title">
          <Translate id="homepage.hero.panel.title">Bootstrap shape</Translate>
        </span>
      </div>
      <pre className="home-code-block" aria-hidden="true">
        <code>
          {heroCode.map((line) => `${line}\n`)}
        </code>
      </pre>
      <ul className="home-panel-list">
        <li>
          <Translate id="homepage.hero.panel.pointOne">
            Resources are registered explicitly and resolved by the API engine.
          </Translate>
        </li>
        <li>
          <Translate id="homepage.hero.panel.pointTwo">
            CRUD builders, result helpers, and parameter injection remove repetitive glue code.
          </Translate>
        </li>
        <li>
          <Translate id="homepage.hero.panel.pointThree">
            Defaults come from source: version `v1`, timeout `30s`, and rate limit `100 / 5m`.
          </Translate>
        </li>
      </ul>
    </div>
  );
}

function Home() {
  const {siteConfig} = useDocusaurusContext();
  const pageTitle = translate({
    id: 'homepage.meta.title',
    message: 'Home',
  });
  const pageDescription = translate({
    id: 'homepage.meta.description',
    message:
      'VEF Framework Go is a resource-driven Go web framework built on FX, Fiber, and Bun, with built-in security, storage, monitoring, and MCP support.',
  });

  return (
    <Layout title={pageTitle} description={pageDescription}>
      <main className="homepage-main">
        <div className="home-shell">
          <section className="home-hero">
            <div className="home-hero-copy">
              <p className="home-kicker">
                <Translate id="homepage.hero.kicker">
                  FX assembly, resource-driven APIs, and built-in infrastructure
                </Translate>
              </p>
              <Heading as="h1" className="home-title">
                {siteConfig.title}
              </Heading>
              <p className="home-description">
                <Translate id="homepage.hero.description">
                  Build Go services around explicit resources, typed handlers, and framework defaults that are visible in source instead of hidden behind scaffolding.
                </Translate>
              </p>
              <div className="home-actions">
                <Link className="button button--primary button--lg" to="/docs/getting-started/quick-start">
                  <Translate id="homepage.hero.primary">Read quick start</Translate>
                </Link>
                <Link className="button button--secondary button--lg" to="/docs/guide/routing">
                  <Translate id="homepage.hero.secondary">Explore routing</Translate>
                </Link>
              </div>
              <div className="home-stat-grid home-stat-grid--metrics">
                {statList.map((item) => (
                  <div className="home-stat-card home-stat-card--metric" key={item.label}>
                    <span className="home-stat-label">{item.label}</span>
                    <strong className="home-stat-value">{item.value}</strong>
                  </div>
                ))}
              </div>
            </div>

            <div className="home-hero-side">
              <img
                src={useBaseUrl('/img/logo.svg')}
                alt="VEF Framework Logo"
                className="home-hero-logo"
              />
              <HeroCodeBlock />
            </div>
          </section>

          <section className="home-section">
            <div className="home-section-head">
              <p className="home-eyebrow">
                <Translate id="homepage.workflow.eyebrow">Mental model</Translate>
              </p>
              <Heading as="h2">
                <Translate id="homepage.workflow.title">
                  Start from composition, not from routes
                </Translate>
              </Heading>
              <p>
                <Translate id="homepage.workflow.subtitle">
                  VEF is easiest to understand when you follow the same shape as the runtime: compose modules, register resources, and let the API engine dispatch operations.
                </Translate>
              </p>
            </div>

            <div className="workflow-grid">
              {workflowList.map((item) => (
                <article className="home-panel workflow-card" key={item.index}>
                  <span className="workflow-index">{item.index}</span>
                  <Heading as="h3">{item.title}</Heading>
                  <p>{item.description}</p>
                </article>
              ))}
            </div>
          </section>

          <section className="home-section">
            <div className="home-section-head">
              <p className="home-eyebrow">
                <Translate id="homepage.features.eyebrow">Capabilities</Translate>
              </p>
              <Heading as="h2">
                <Translate id="homepage.features.title">
                  The parts you should reach for first
                </Translate>
              </Heading>
              <p>
                <Translate id="homepage.features.subtitle">
                  These are the framework surfaces that show up early in real projects and are backed directly by the current source tree.
                </Translate>
              </p>
            </div>

            <div className="feature-grid">
              {featureList.map((item) => (
                <article className="home-panel feature-card" key={item.title}>
                  <span className="feature-badge">{item.badge}</span>
                  <Heading as="h3">{item.title}</Heading>
                  <p>{item.description}</p>
                </article>
              ))}
            </div>
          </section>

          <section className="home-section">
            <div className="home-section-head">
              <p className="home-eyebrow">
                <Translate id="homepage.paths.eyebrow">Start here</Translate>
              </p>
              <Heading as="h2">
                <Translate id="homepage.paths.title">
                  Follow the shortest path to useful context
                </Translate>
              </Heading>
              <p>
                <Translate id="homepage.paths.subtitle">
                  These pages are the fastest way to understand how a VEF application is structured and how the main APIs are intended to be used.
                </Translate>
              </p>
            </div>

            <div className="link-grid">
              {pathList.map((item) => (
                <Link className="home-panel link-card" key={item.title} to={item.to}>
                  <span className="link-card-arrow" aria-hidden="true">
                    →
                  </span>
                  <Heading as="h3">{item.title}</Heading>
                  <p>{item.description}</p>
                </Link>
              ))}
            </div>
          </section>
        </div>
      </main>
    </Layout>
  );
}

export default function Homepage(): ReactNode {
  return <Home />;
}
