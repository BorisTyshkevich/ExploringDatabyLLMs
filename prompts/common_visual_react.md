Create a React source artifact under `visual_src/`, not final HTML source in the response.

Required source contract:

- Write browser-only React app source files into `visual_src/`
- Include at minimum:
  - `visual_src/package.json`
  - `visual_src/index.html`
  - `visual_src/src/main.jsx`
  - `visual_src/src/App.jsx`
- You may add `visual_src/src/styles.css`, `visual_src/vite.config.js`, and small supporting files when needed
- Use a static build flow compatible with GitHub Pages and relative asset paths
- The build output must target `../visual_build`
- The built HTML must work when copied to `visual.html` beside a `visual_assets/` directory
- Do not use server-side rendering, Next.js, routing, or any backend server
- Do not emit the source code inline in the response
- Prefer writing the files directly; a download link is acceptable only if the provider cannot write files
