# Changelog

## 1.0.0 (2026-01-21)


### ⚠ BREAKING CHANGES

* **api:** API field names changed
    - raw_text → user_input (original user input, immutable)
    - description → parsed_description (LLM-generated description)

### Features

* add backend environment with FastAPI and Docker ([fc01fba](https://github.com/austinmoody/tact/commit/fc01fbaa02712c965924a06c822b72198c22d44f))
* add entries API with CRUD endpoints ([89d0afa](https://github.com/austinmoody/tact/commit/89d0afaf61463815cf242d07dac487b80978bbfc))
* add GitHub Copilot custom agents for OpenSpec workflow ([858e5b2](https://github.com/austinmoody/tact/commit/858e5b23097bd08bafa3b3400d5a6ac9482c9fea))
* add LLM entry parsing with Ollama and Anthropic providers ([fb59086](https://github.com/austinmoody/tact/commit/fb59086b5cbb6a13b320ca8fa59ed93eea0f856d))
* add MCP server for Tact API ([2898900](https://github.com/austinmoody/tact/commit/28989004e3ca6fc537fe5282ec4e208d6ec636dd))
* add RAG-based context parsing for time entries ([4dd4a8a](https://github.com/austinmoody/tact/commit/4dd4a8a07b1214e2e46fa8019af50ce6cf6992fc))
* add SQLite data model with Alembic migrations ([7eeb6d0](https://github.com/austinmoody/tact/commit/7eeb6d0a2eabec773f153970f21e229145ff9fa9))
* add structured API logging for all operations ([71ed0d5](https://github.com/austinmoody/tact/commit/71ed0d5f101d46fbf50f9d0acb9cfcbc913ada5a))
* add Time Codes and Work Types CRUD API ([49f1632](https://github.com/austinmoody/tact/commit/49f16327bdedfdc079960a1425a7dcf75f465d33))
* add TUI dashboard with Go/Bubbletea ([4cfc459](https://github.com/austinmoody/tact/commit/4cfc459abc7ac8869a3f9a645796233c2559311f))
* **api:** add configurable duration rounding ([17e477a](https://github.com/austinmoody/tact/commit/17e477a45723bce510f504532b417303bc4a954f))
* **api:** add learning from manual corrections ([4987b96](https://github.com/austinmoody/tact/commit/4987b9627afe11ac3cc716f68b57838c766ed41b))
* **api:** add learning from manual corrections ([67edbcf](https://github.com/austinmoody/tact/commit/67edbcf018f8b3f1508c1ae4965f5333391e4ff1))
* **api:** rename time entry fields for clarity ([2b727ac](https://github.com/austinmoody/tact/commit/2b727ac83a4828df259fd15511d258c6c158b9b3))
* **ci:** add GitHub Actions workflows for CI/CD ([78f0a5d](https://github.com/austinmoody/tact/commit/78f0a5d0dcd12752bebd0534c4aa41bf8ecf317d))
* implement parse_notes field for time entries ([44f3c20](https://github.com/austinmoody/tact/commit/44f3c2089deb369bb885cbb03903383e8f4fefba))
* **llm:** add configurable timeout for Ollama requests ([1c151ec](https://github.com/austinmoody/tact/commit/1c151ece4319d1026903029211ef5acc1393e1a3))
* **macos:** add completed timers section with quick restart ([87daab8](https://github.com/austinmoody/tact/commit/87daab8c05bad0316f47a3e56539dc1e776a5764))
* **macos:** add native macOS timer app (Tact Timer) ([d73db50](https://github.com/austinmoody/tact/commit/d73db506f4d8c74259d37d71fc2a725ee633b924))
* **mcp:** add project and context tools ([1918c3d](https://github.com/austinmoody/tact/commit/1918c3deafb71ee1fa4a07aa245e70c1495c952e))
* **ollama:** auto-pull models on first use ([732de52](https://github.com/austinmoody/tact/commit/732de5298c6c90e255a4d8b425b8e245e6ef1033))
* redesign TUI with entry-focused home screen ([d99fb7b](https://github.com/austinmoody/tact/commit/d99fb7be022866cd7ade3b0fabc38385dedb24d0))
* remove description field from projects ([9224474](https://github.com/austinmoody/tact/commit/9224474999e05d57b7337931046dc7b21455363b))
* simplify time code model by removing unused fields ([657d5ad](https://github.com/austinmoody/tact/commit/657d5ad513b45c73e36f3dcd713407bea7f8361b))
* **tui:** add project selection to time code add/edit modal ([d7c0db7](https://github.com/austinmoody/tact/commit/d7c0db7e2cbadbede4e704df7da33e8c98d03058))
* **tui:** add Projects and Context management screens ([2f40142](https://github.com/austinmoody/tact/commit/2f401423afa025a8420329c8a909bfc5adf3dbc9))
* **tui:** add timer functionality with floating panel ([6c68c10](https://github.com/austinmoody/tact/commit/6c68c100f0699cc382681bc05b4eed69a1aab7d7))
* **tui:** upgrade to Bubble Tea v2 for improved input handling ([4b61628](https://github.com/austinmoody/tact/commit/4b61628569c369f14ef14ed976b2090d2cdfaae3))


### Bug Fixes

* address Copilot review feedback ([2878b9f](https://github.com/austinmoody/tact/commit/2878b9fd99c0e7422cc21af07e50057269d10f2c))
* address lint issues found by new CI ([6707f60](https://github.com/austinmoody/tact/commit/6707f6093d98dd190b3c66245479e253572e7b87))
* **db:** merge alembic migration heads ([2b47160](https://github.com/austinmoody/tact/commit/2b47160a896e9fcac08a4d1be65f1ff2ac235abf))
* **llm:** convert float duration_minutes to int ([2bd46b6](https://github.com/austinmoody/tact/commit/2bd46b6417920981f91929d78f1361d4b27e649b))
* **llm:** convert float duration_minutes to int ([3d398a0](https://github.com/austinmoody/tact/commit/3d398a0438b1f5e52b646780cd2267a83f831090))
* **mcp:** remove get_summary tool (endpoint doesn't exist) ([57b7ad9](https://github.com/austinmoody/tact/commit/57b7ad9ab38a05bb0359893fe5c9b89215ad24b8))
* **parser:** prevent API blocking during parse operations ([ccdf22e](https://github.com/austinmoody/tact/commit/ccdf22e9e707ccf003441a74ebed92d48d737b00))
* **parser:** require both time_code and duration for parsed status ([3f9036a](https://github.com/austinmoody/tact/commit/3f9036aaf493f24a8bd37101ac10b5b38dbb2656))
* **parser:** require both time_code and duration for parsed status ([c4df57f](https://github.com/austinmoody/tact/commit/c4df57f0976b1972fcf671fd1cbcb929c63b4ca9))
* **parser:** set status based on confidence threshold ([185b03c](https://github.com/austinmoody/tact/commit/185b03c450d97b68cb88127aa41e7e1a6006b5a1))
* **parser:** set status based on confidence threshold ([a239ae9](https://github.com/austinmoody/tact/commit/a239ae9ab6a2ec5ca62e167f9afc967ae1f0dc9b))
* resolve all ruff lint errors ([c4aa461](https://github.com/austinmoody/tact/commit/c4aa461eefd7db2e213e8890ea1f45d89a4627e0))
* resolve remaining lint issues in backend code ([85b1cc4](https://github.com/austinmoody/tact/commit/85b1cc40b1d6f7c1c42b9748fdf63c173241e6e2))
* **tui:** route async messages to modals for context loading ([7f13855](https://github.com/austinmoody/tact/commit/7f1385540dde23953c7a2883af566250bc47febc))
