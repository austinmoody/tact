## Why

The current web UI uses Pico CSS with a modern dark theme, which works well but lacks personality. Applying a Windows 3.1 aesthetic would give TACT a unique, memorable visual identity with nostalgic retro charm. The classic Windows 3.1 look - with its gray backgrounds, beveled 3D borders, system fonts, and blue title bars - is instantly recognizable and would make the time tracking experience more fun.

## What Changes

- Replace Pico CSS with custom Windows 3.1-style CSS
- Implement classic Windows 3.1 visual elements:
  - Gray (#C0C0C0) window backgrounds with white/dark gray beveled borders
  - Blue (#000080) title bars with white text
  - System font styling (MS Sans Serif look using Arial/Helvetica fallback)
  - Classic 3D button effects with inset/outset borders
  - Black window frames and outlines
- Preserve all existing HTMX functionality and page structure
- Maintain responsive behavior for modern devices

## Capabilities

### New Capabilities

None - this is a visual-only change.

### Modified Capabilities

- `web-ui-core`: Update styling specification to use Windows 3.1 theme instead of Pico CSS dark mode. The base layout structure and HTMX integration remain unchanged.

## Impact

- **Modified Code**: `webui/static/css/` - replace `pico.min.css` with custom Win31 CSS or add theme override
- **Modified Templates**: Minor adjustments to template classes/structure if needed for theme compatibility
- **No Backend Changes**: This is purely a frontend visual change
- **No New Dependencies**: Uses standard CSS
- **Build**: No changes to build process

## Visual Design Reference

Key Windows 3.1 design patterns to implement:

1. **Color Palette**
   - Window background: #C0C0C0 (classic Windows gray)
   - Title bar: #000080 (navy blue)
   - Title text: #FFFFFF (white)
   - Button face: #C0C0C0
   - Button highlight: #FFFFFF
   - Button shadow: #808080
   - Window frame: #000000

2. **Border Effects**
   - Raised elements: white top/left, dark gray bottom/right
   - Sunken elements: dark gray top/left, white bottom/right
   - Window frames: thin black outline

3. **Typography**
   - Primary font: Arial, Helvetica, sans-serif (approximating MS Sans Serif)
   - Title bar font: bold
   - No anti-aliasing effects where possible for authentic look

4. **UI Elements**
   - Buttons with classic 3D beveled appearance
   - Text inputs with sunken inset border
   - Title bars with system menu box on left
   - Active/inactive window state differentiation
