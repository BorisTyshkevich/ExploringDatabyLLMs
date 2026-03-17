# OnTime Dashboard Theme

Use a transport-analytics visual language, not the alert-dashboard purple theme.

## Design goals

- Aviation/corporate rather than incident-response
- Calm, high-contrast, information-dense
- Consistent across static, dynamic, map, and non-map dashboards

## Tokens

```css
:root {
    --bg-top: #eaf3f8;
    --bg-bottom: #f6fafc;
    --panel: #ffffff;
    --panel-alt: #f3f7fa;
    --ink: #163244;
    --muted: #5d7485;
    --navy: #0e3a52;
    --slate: #5c7080;
    --sky: #3c88b5;
    --teal: #1f8a70;
    --amber: #d48a1f;
    --red: #c54f36;
    --grid: rgba(22, 50, 68, 0.12);
    --border: rgba(22, 50, 68, 0.10);
    --shadow: 0 18px 45px rgba(14, 58, 82, 0.10);
    --radius-xl: 22px;
    --radius-lg: 16px;
    --radius-md: 12px;
    --radius-sm: 8px;
}
```

## Semantic colors

- Good / improvement / high OTP: `--teal`
- Neutral / structural context: `--sky` or `--slate`
- Warning / elevated delay: `--amber`
- Worst case / severe disruption: `--red`
- Text and framing: `--navy`, `--ink`

## Layout rules

- Page background should use a soft layered gradient, not a flat color
- Main content cards should use white panels with rounded corners and subtle shadows
- Header should be visually distinct with an aviation-toned gradient or strong navy treatment
- KPI strip should use 4-up responsive cards with short labels and large values
- Default content width: `max-width: 1280px`
- Prefer CSS grid layouts with `repeat(auto-fit, minmax(...))`

## Typography

- Prefer a non-default serif/sans pairing already available on the system, for example:
  - headings: Georgia or ui-serif
  - body/UI: `"Segoe UI", "Helvetica Neue", sans-serif`
- Headline should feel editorial, not generic SaaS
- Keep table and chart labels compact and readable

## Table style

- Sticky or strongly styled header row when practical
- Compact numeric columns, wider label columns
- Row hover state
- Borders should use `--border`, not hard black lines

## Legend and annotation rules

- Always explain what red/amber/green mean in domain terms
- Annotate the single most important row, cell, route, or period
- When ranking, highlight top 1 and optionally top 3

## Map alignment

When maps are used, reuse the same colors and panel styling:

- airport dots: `--navy`
- highlighted airport or route: `--red`
- secondary routes: `--sky`
- contextual fills: muted blue/gray

## Query ledger style

- Ledger entries are collapsible rows with expand/collapse toggle
- Toggle icon: ▶ collapsed, ▼ expanded (monospace font)
- Row hover: `var(--panel-alt)` background
- Row click: expands/collapses the SQL block
- SQL block: monospace, `var(--panel-alt)` background, `var(--radius-sm)` corners
- SQL text: `pre` with word-wrap for long lines
- Status colors: OK=`--teal`, Pending=`--amber`, Failed=`--red`
- Grid columns: toggle (1em), label (flex), role (6em), status (5em), rows (4em)
