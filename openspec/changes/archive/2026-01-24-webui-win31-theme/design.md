# Design: Windows 3.1 Theme for Web UI

## Overview

This design documents the implementation approach for applying a Windows 3.1 visual theme to the TACT web UI, replacing the Pico CSS dark theme.

## Architecture

### CSS Structure

The theme is implemented as a single CSS file (`webui/static/css/win31.css`) that completely replaces Pico CSS:

```
webui/static/css/
├── win31.css        # Complete Windows 3.1 theme (replaces pico.min.css)
```

### CSS Organization

The stylesheet is organized into logical sections:

1. **Color Variables** - CSS custom properties for the Windows 3.1 palette
2. **Reset and Base Styles** - Normalize browser defaults
3. **Typography** - Font families and text styling
4. **Layout** - Page structure and grid
5. **Navigation** - Header styled as title bar
6. **Buttons** - 3D beveled button styles
7. **Forms** - Input fields with sunken borders
8. **Tables** - Grid-style data display
9. **Cards/Articles** - Window-like containers
10. **Dialogs** - Modal windows
11. **Utilities** - Status colors, tooltips
12. **Responsive** - Mobile adaptations

## Key Design Decisions

### 1. Complete CSS Replacement

**Decision**: Replace Pico CSS entirely rather than override it.

**Rationale**: Windows 3.1 aesthetic is fundamentally different from modern CSS frameworks. Overriding would lead to specificity conflicts and bloated CSS. A clean implementation is simpler to maintain.

### 2. CSS Custom Properties for Colors

**Decision**: Use CSS variables for the color palette.

```css
:root {
    --win31-gray: #C0C0C0;
    --win31-dark-gray: #808080;
    --win31-title-active: #000080;
    /* ... */
}
```

**Rationale**: Enables consistent theming and potential future customization.

### 3. 3D Beveled Borders

**Decision**: Use multi-value `border-color` with `box-shadow` for authentic 3D effects.

```css
.raised {
    border: 2px solid;
    border-color: var(--win31-button-highlight) var(--win31-button-dark-shadow)
                  var(--win31-button-dark-shadow) var(--win31-button-highlight);
    box-shadow: inset 1px 1px 0 var(--win31-light-gray);
}
```

**Rationale**: Achieves the authentic Windows 3.1 beveled appearance without images.

### 4. Navigation as Title Bar

**Decision**: Style the navigation header as a Windows 3.1 title bar with navy background.

**Rationale**: Creates immediate visual recognition of the retro theme while maintaining navigation functionality.

### 5. CSS-Based Tooltips

**Decision**: Use `::before` pseudo-elements with `data-tooltip` attributes instead of native `title` tooltips.

```css
.tooltip::before {
    content: "Parsed: " attr(data-tooltip);
    position: absolute;
    /* ... */
}
```

**Rationale**: Native `title` tooltips are unreliable across browsers. CSS tooltips provide consistent, styleable behavior.

### 6. Dialog Close Buttons

**Decision**: Style dialog close buttons as small Windows 3.1 buttons with "X" text.

**Rationale**: Matches the authentic Windows 3.1 window control appearance.

## Component Mapping

| Web UI Component | Windows 3.1 Equivalent |
|------------------|------------------------|
| Navigation header | Title bar |
| Main content area | Window client area |
| Cards/Articles | Child windows |
| Buttons | 3D push buttons |
| Text inputs | Sunken edit controls |
| Select dropdowns | Combo boxes |
| Tables | List views (grid mode) |
| Dialogs | Popup windows |
| Footer | Status bar |

## Template Changes

Minimal template changes were required:

1. **base.templ**: Changed CSS reference from `pico.min.css` to `win31.css`
2. **entry_list.templ**: Added `tooltip` class and `data-tooltip` attribute for parsed descriptions
3. **context_list.templ**: Added `win31-close-btn` class to dialog close button

## Responsive Considerations

The theme maintains responsiveness while preserving retro aesthetics:

- Navigation collapses to hamburger menu on mobile (future enhancement)
- Tables scroll horizontally on narrow screens
- Font sizes adjust for readability
- Touch targets remain accessible despite compact Windows 3.1 sizing

## Browser Compatibility

The CSS uses widely-supported features:
- CSS custom properties (variables)
- Flexbox and Grid layout
- Pseudo-elements (::before, ::after)
- Standard border and box-shadow properties

No vendor prefixes required for modern browser support.
