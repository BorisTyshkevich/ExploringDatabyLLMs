# Edit Mode

Every dashboard should support a lightweight native edit mode.

## Activation

- Enable with `?edit=1`
- Default view remains read-only

## Required capabilities

- Move cards between primary and secondary regions
- Hide cards
- Restore hidden cards
- Persist layout in `localStorage`
- Primary region accepts only one main visual

## Storage

Use a page-specific layout key:

`OnTimeAnalystDashboard::<dashboardId>::layout`

## Implementation guidance

Prefer native browser behavior over external drag libraries:

- simple up/down/move-to-primary/move-to-secondary buttons are acceptable
- native HTML drag-and-drop is acceptable if kept minimal

## Required DOM concepts

- primary container
- secondary container
- hidden container or restore tray
- each card has a stable `data-card-id`

## Rerender hooks

After any layout change:

- persist layout
- rerender or resize SVG charts
- if a Leaflet map exists, invalidate size after the card is visible again

## Visibility

- edit controls visible only in edit mode
- hidden cards excluded from tab order in read-only mode
