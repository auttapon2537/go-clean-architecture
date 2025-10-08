package main

import (
	"bytes"
	"fmt"
	"log"

	projectdocs "github.com/example/go-clean-architecture/docs"
	"github.com/gofiber/fiber/v2"
)

const (
	openAPIJSONFile = "openapi.json"
	openAPIYAMLFile = "openapi.yaml"
)

// OpenAPISpecHandler serves the OpenAPI specification file in the requested format.
func OpenAPISpecHandler(specPath, format string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data, err := projectdocs.OpenAPIFS.ReadFile(specPath)
		if err != nil {
			log.Printf("ERROR: Unable to read OpenAPI spec at %s: %v", specPath, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "OpenAPI specification not available",
			})
		}

		switch format {
		case "json":
			c.Type("json")
		case "yaml":
			c.Type("yaml")
		default:
			c.Type("octet-stream")
		}

		return c.SendStream(bytes.NewReader(data))
	}
}

// OpenAPIDocsHandler serves an offline HTML viewer for the OpenAPI specification.
func OpenAPIDocsHandler(specURL string) fiber.Handler {
	page := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>API Documentation</title>
  <style>
    :root {
      color-scheme: light dark;
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
      --bg: #f4f6f8;
      --card-bg: #ffffff;
      --text: #1f2933;
      --muted: #6b7280;
      --border: #d1d5db;
      --accent: #2563eb;
      --method-get: #047857;
      --method-post: #2563eb;
      --method-put: #d97706;
      --method-delete: #dc2626;
    }
    @media (prefers-color-scheme: dark) {
      :root {
        --bg: #0f172a;
        --card-bg: #111827;
        --text: #e5e7eb;
        --muted: #9ca3af;
        --border: #1f2937;
      }
    }
    body {
      margin: 0;
      background: var(--bg);
      color: var(--text);
    }
    header {
      padding: 32px 24px 16px;
      background: var(--card-bg);
      border-bottom: 1px solid var(--border);
    }
    header h1 {
      margin: 0 0 8px;
      font-size: 28px;
    }
    header p {
      margin: 0;
      color: var(--muted);
    }
    main {
      padding: 24px;
      display: grid;
      gap: 16px;
    }
    .card {
      background: var(--card-bg);
      border: 1px solid var(--border);
      border-radius: 12px;
      padding: 24px;
      box-shadow: 0 1px 2px rgba(15, 23, 42, 0.08);
    }
    .tag {
      display: inline-block;
      font-weight: 600;
      letter-spacing: 0.05em;
      padding: 2px 10px;
      border-radius: 999px;
      font-size: 12px;
      text-transform: uppercase;
      color: #fff;
    }
    .tag.get { background: var(--method-get); }
    .tag.post { background: var(--method-post); }
    .tag.put { background: var(--method-put); }
    .tag.delete { background: var(--method-delete); }
    details {
      border: 1px solid var(--border);
      border-radius: 10px;
      background: var(--card-bg);
      box-shadow: 0 1px 2px rgba(15, 23, 42, 0.04);
    }
    details summary {
      cursor: pointer;
      display: flex;
      align-items: center;
      gap: 16px;
      padding: 16px 20px;
      list-style: none;
    }
    details summary::-webkit-details-marker {
      display: none;
    }
    summary .path {
      font-family: "SFMono-Regular", Menlo, Consolas, monospace;
      font-size: 15px;
    }
    .description {
      margin: 12px 0 0;
      color: var(--muted);
      line-height: 1.5;
    }
    .section {
      padding: 0 20px 20px;
      border-top: 1px solid var(--border);
    }
    .section h4 {
      margin: 20px 0 8px;
      font-size: 15px;
      letter-spacing: 0.03em;
      text-transform: uppercase;
      color: var(--muted);
    }
    table {
      width: 100%%;
      border-collapse: collapse;
      margin-top: 8px;
      font-size: 14px;
    }
    th, td {
      text-align: left;
      padding: 8px 12px;
      border: 1px solid var(--border);
    }
    th {
      background: rgba(148, 163, 184, 0.1);
      font-weight: 600;
      text-transform: uppercase;
      font-size: 12px;
      color: var(--muted);
    }
    pre {
      background: rgba(15, 23, 42, 0.04);
      padding: 12px;
      border-radius: 8px;
      overflow: auto;
      font-size: 13px;
      font-family: "SFMono-Regular", Menlo, Consolas, monospace;
    }
    a.spec-link {
      color: var(--accent);
      text-decoration: none;
      font-weight: 600;
    }
    a.spec-link:hover {
      text-decoration: underline;
    }
    .loading, .error {
      text-align: center;
      padding: 40px;
      color: var(--muted);
    }
    footer {
      text-align: center;
      padding: 16px;
      color: var(--muted);
      font-size: 13px;
    }
  </style>
</head>
<body>
  <header>
    <h1 id="api-title">API Documentation</h1>
    <p id="api-description"></p>
    <p class="spec">
      View raw specs:
      <a class="spec-link" href="%[1]s" target="_blank" rel="noopener">JSON</a>
      ·
      <a class="spec-link" href="/openapi.yaml" target="_blank" rel="noopener">YAML</a>
    </p>
  </header>
  <main>
    <div id="info-card" class="card"></div>
    <div id="paths"></div>
    <div id="schemas" class="card"></div>
  </main>
  <footer>Powered by Go Fiber · Offline-friendly documentation</footer>
  <script>
    const METHOD_COLORS = {
      get: 'get',
      post: 'post',
      put: 'put',
      delete: 'delete',
      patch: 'post',
      options: 'get',
      head: 'get'
    };

    function createTag(method) {
      const tag = document.createElement('span');
      tag.className = 'tag ' + (METHOD_COLORS[method] || 'get');
      tag.textContent = method.toUpperCase();
      return tag;
    }

    function createSummary(method, path, summary) {
      const summaryEl = document.createElement('summary');
      summaryEl.appendChild(createTag(method));
      const pathEl = document.createElement('span');
      pathEl.className = 'path';
      pathEl.textContent = path;
      summaryEl.appendChild(pathEl);
      if (summary) {
        const summaryText = document.createElement('span');
        summaryText.textContent = summary;
        summaryEl.appendChild(summaryText);
      }
      return summaryEl;
    }

    function renderParameters(parameters) {
      if (!parameters || parameters.length === 0) return null;
      const section = document.createElement('div');
      section.className = 'section';
      const title = document.createElement('h4');
      title.textContent = 'Parameters';
      section.appendChild(title);
      const table = document.createElement('table');
      table.innerHTML = '<thead><tr><th>Name</th><th>In</th><th>Type</th><th>Required</th><th>Description</th></tr></thead>';
      const tbody = document.createElement('tbody');
      parameters.forEach(param => {
        const row = document.createElement('tr');
        row.innerHTML = '<td>' + param.name + '</td>' +
          '<td>' + param.in + '</td>' +
          '<td>' + (param.schema ? param.schema.type || '' : '') + '</td>' +
          '<td>' + (param.required ? 'Yes' : 'No') + '</td>' +
          '<td>' + (param.description || '') + '</td>';
        tbody.appendChild(row);
      });
      table.appendChild(tbody);
      section.appendChild(table);
      return section;
    }

    function renderRequestBody(body) {
      if (!body) return null;
      const section = document.createElement('div');
      section.className = 'section';
      const title = document.createElement('h4');
      title.textContent = 'Request Body';
      section.appendChild(title);

      if (body.description) {
        const description = document.createElement('p');
        description.className = 'description';
        description.textContent = body.description;
        section.appendChild(description);
      }

      const content = body.content || {};
      Object.entries(content).forEach(([mime, cfg]) => {
        const mimeTitle = document.createElement('strong');
        mimeTitle.textContent = mime;
        section.appendChild(mimeTitle);
        if (cfg.schema) {
          const pre = document.createElement('pre');
          pre.textContent = JSON.stringify(cfg.schema, null, 2);
          section.appendChild(pre);
        }
      });
      return section;
    }

    function renderResponses(responses) {
      if (!responses) return null;
      const section = document.createElement('div');
      section.className = 'section';
      const title = document.createElement('h4');
      title.textContent = 'Responses';
      section.appendChild(title);

      Object.entries(responses).forEach(([status, resp]) => {
        const wrapper = document.createElement('div');
        const heading = document.createElement('strong');
        heading.textContent = status + ' · ' + (resp.description || 'Response');
        wrapper.appendChild(heading);
        if (resp.content) {
          Object.entries(resp.content).forEach(([mime, cfg]) => {
            const mimeLabel = document.createElement('div');
            mimeLabel.style.marginTop = '6px';
            mimeLabel.style.fontSize = '13px';
            mimeLabel.style.color = 'var(--muted)';
            mimeLabel.textContent = mime;
            wrapper.appendChild(mimeLabel);
            if (cfg.schema) {
              const pre = document.createElement('pre');
              pre.textContent = JSON.stringify(cfg.schema, null, 2);
              wrapper.appendChild(pre);
            }
          });
        }
        section.appendChild(wrapper);
      });

      return section;
    }

    function renderOperation(method, path, operation) {
      const details = document.createElement('details');
      details.appendChild(createSummary(method, path, operation.summary || ''));

      const body = document.createElement('div');
      body.className = 'section';
      if (operation.description) {
        const desc = document.createElement('p');
        desc.className = 'description';
        desc.textContent = operation.description;
        body.appendChild(desc);
      }
      details.appendChild(body);

      const params = renderParameters(operation.parameters);
      if (params) details.appendChild(params);

      const request = renderRequestBody(operation.requestBody);
      if (request) details.appendChild(request);

      const responses = renderResponses(operation.responses);
      if (responses) details.appendChild(responses);

      return details;
    }

    function renderPaths(paths) {
      const container = document.getElementById('paths');
      container.innerHTML = '';
      if (!paths || Object.keys(paths).length === 0) {
        container.innerHTML = '<div class="card">No paths defined.</div>';
        return;
      }

      Object.keys(paths).sort().forEach(path => {
        const operations = paths[path];
        Object.keys(operations).forEach(method => {
          const op = operations[method];
          container.appendChild(renderOperation(method, path, op));
        });
      });
    }

    function renderSchemas(schemas) {
      const container = document.getElementById('schemas');
      container.innerHTML = '';
      if (!schemas) {
        container.style.display = 'none';
        return;
      }
      container.style.display = 'block';
      const title = document.createElement('h2');
      title.textContent = 'Schemas';
      container.appendChild(title);

      Object.entries(schemas).forEach(([name, schema]) => {
        const section = document.createElement('section');
        const heading = document.createElement('h3');
        heading.textContent = name;
        heading.style.marginBottom = '8px';
        section.appendChild(heading);
        const pre = document.createElement('pre');
        pre.textContent = JSON.stringify(schema, null, 2);
        section.appendChild(pre);
        container.appendChild(section);
      });
    }

    function renderInfo(info, servers) {
      document.getElementById('api-title').textContent = info.title || 'API Documentation';
      document.getElementById('api-description').textContent = info.description || '';
      const card = document.getElementById('info-card');
      card.innerHTML = '';
      const version = document.createElement('p');
      version.innerHTML = '<strong>Version:</strong> ' + (info.version || 'N/A');
      card.appendChild(version);
      if (servers && servers.length) {
        const list = document.createElement('ul');
        list.style.paddingLeft = '20px';
        list.style.marginTop = '8px';
        servers.forEach(server => {
          const li = document.createElement('li');
          li.textContent = server.url;
          list.appendChild(li);
        });
        const serversTitle = document.createElement('strong');
        serversTitle.textContent = 'Servers:';
        card.appendChild(serversTitle);
        card.appendChild(list);
      }
    }

    async function init() {
      const container = document.getElementById('paths');
      container.innerHTML = '<div class="card loading">Loading specification...</div>';
      try {
        const response = await fetch('%[1]s');
        if (!response.ok) throw new Error('Unable to load specification');
        const spec = await response.json();
        renderInfo(spec.info || {}, spec.servers || []);
        renderPaths(spec.paths || {});
        renderSchemas(spec.components ? spec.components.schemas : null);
      } catch (error) {
        container.innerHTML = '<div class="card error">Failed to load API specification. ' + error.message + '</div>';
      }
    }

    init();
  </script>
</body>
</html>`, specURL)

	return func(c *fiber.Ctx) error {
		c.Type("html", "utf-8")
		return c.SendString(page)
	}
}
