// This file has been automatically migrated to valid ESM format by Storybook.
import { fileURLToPath } from "node:url";
import { resolve, dirname } from "path";

import type { StorybookConfig } from "@storybook/react-vite";
import { mergeConfig } from "vite";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const config: StorybookConfig = {
  stories: ["../src/**/*.stories.@(js|ts|tsx|mdx)"],
  addons: [
    "@storybook/addon-onboarding",
    "@storybook/addon-a11y",
    "@storybook/addon-themes",
    "@chromatic-com/storybook",
    "@storybook/addon-docs",
  ],

  framework: {
    name: "@storybook/react-vite",
    options: {},
  },

  core: {
    disableTelemetry: true,
  },

  staticDirs: ["./public"],

  viteFinal(config, { configType }) {
    return mergeConfig(config, {
      define: {
        "process.env.QTS_DEBUG": "false", // quickjs-emscripten
      },

      build:
        configType === "PRODUCTION"
          ? {
              // https://github.com/storybookjs/builder-vite/issues/409
              minify: false,
              sourcemap: false,
            }
          : {},
      resolve: {
        alias: [
          {
            find: "crypto",
            replacement: "crypto-js",
          },
          {
            find: "@reearth/cesium-mvt-imagery-provider",
            replacement: resolve(
              __dirname,
              "..",
              "node_modules/@reearth/cesium-mvt-imagery-provider",
            ),
          },
          {
            find: "@reearth/core",
            replacement: resolve(__dirname, "..", "node_modules/@reearth/core"),
          },
          // cesium-mvt-imagery-provider: force resolution to the JS bundle, not the .d.ts types file
          {
            find: "cesium-mvt-imagery-provider",
            replacement: resolve(
              __dirname,
              "..",
              "node_modules/cesium-mvt-imagery-provider/dist/cesium-mvt-imagery-provider.mjs",
            ),
          },
          // quickjs-emscripten: force resolution to the JS bundle, not the .d.ts types file
          {
            find: "quickjs-emscripten-sync",
            replacement: resolve(
              __dirname,
              "..",
              "node_modules/quickjs-emscripten-sync/dist/quickjs-emscripten-sync.mjs",
            ),
          },
          // react-align: force resolution to the JS bundle, not the .d.ts types file
          {
            find: "react-align",
            replacement: resolve(__dirname, "..", "node_modules/react-align/dist/react-align.mjs"),
          },
          {
            find: "@reearth",
            replacement: resolve(__dirname, "..", "src"),
          },
          {
            find: "csv-parse",
            replacement: "csv-parse/browser/esm",
          },
        ],
      },
      optimizeDeps: {
        include: ["quickjs-emscripten-sync", "prop-types", "hoist-non-react-statics", "react-is"],
        exclude: ["react-align"],
      },
      server: {
        watch: {
          // https://github.com/storybookjs/storybook/issues/22253#issuecomment-1673229400
          ignored: ["**/.env"],
        },
      },
    });
  },

  docs: {},

  typescript: {
    reactDocgen: "react-docgen-typescript",
  },
};
export default config;
