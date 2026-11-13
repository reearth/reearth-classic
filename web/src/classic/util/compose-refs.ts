import composeRefsModule from "@seznam/compose-react-refs";

// The package is CJS with `exports.default`. Since Vite 8 the default import of such a module
// resolves to `module.exports` (`{ default: fn }`) rather than the function itself.
const composeRefs: typeof composeRefsModule =
  (composeRefsModule as unknown as { default?: typeof composeRefsModule }).default ??
  composeRefsModule;

export default composeRefs;
