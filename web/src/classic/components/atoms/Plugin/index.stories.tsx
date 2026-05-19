import { Meta, StoryObj } from "@storybook/react-vite";
import { useRef } from "react";
import { action } from "storybook/actions";

import Component, { Ref } from ".";

const meta: Meta<typeof Component> = {
  title: "classic/atoms/Plugin(classic)",
  component: Component,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Component>;

let cb: (message: any) => void | undefined;

export const Default: Story = {
  args: {
    src: `/plugins/plugin.js`,
    uiVisible: true,
    iFrameProps: {
      style: {
        width: "300px",
        height: "300px",
        backgroundColor: "#fff",
      },
    },
    exposed: ({ main: { render, postMessage } }) => ({
      console: {
        log: action("console.log"),
      },
      reearth: {
        on(type: string, value: (message: any) => void | undefined) {
          if (type === "message") {
            cb = value;
          }
        },
        ui: {
          show: render,
          postMessage,
        },
      },
    }),
    onMessage: (message: any) => {
      action("onMessage")(message);
      return cb?.(message);
    },
  },
};

export const HiddenIFrame: Story = {
  args: {
    src: `/plugins/hidden.js`,
    uiVisible: true,
    iFrameProps: {
      style: {
        width: "300px",
        height: "300px",
        backgroundColor: "#fff",
      },
    },
    exposed: ({ main: { render, postMessage } }) => ({
      console: {
        log: action("console.log"),
      },
      reearth: {
        on(type: string, value: (message: any) => void | undefined) {
          if (type === "message") {
            cb = value;
          }
        },
        ui: {
          show: render,
          postMessage,
        },
      },
    }),
    onMessage: (message: any) => {
      action("onMessage")(message);
      return cb?.(message);
    },
  },
};

export const SourceCode: Story = {
  args: {
    sourceCode: `console.log("Hello")`,
    exposed: {
      console: {
        log: action("console.log"),
      },
    },
  },
};

const AutoResizeRenderer = (args: Story["args"]) => {
  const ref = useRef<Ref>(null);
  return (
    <Component
      {...args}
      onMessage={msg => {
        ref.current
          ?.arena()
          ?.evalCode(`"onmessage" in globalThis && globalThis.onmessage(${JSON.stringify(msg)})`);
      }}
      ref={ref}
    />
  );
};

export const AutoResize: Story = {
  render: args => <AutoResizeRenderer {...args} />,
  args: {
    sourceCode: `
    render(\`
      <style>body{width: 100px; height: 50px; background:yellow}</style>
      <h1 id="s"></h1>
      <script>
        let a = false;
        const s = document.getElementById("s");

        setInterval(() => {
          a = !a;
          const w = a ? "300px" : "100px";
          document.body.style.width = w;
          parent.postMessage({ width: w }, "*");
        }, 1000);

        const resize = () => {
          s.textContent = window.innerWidth + "x" + window.innerHeight;
        };
        window.onresize = resize;
        resize();
      </script>
    \`);
    onmessage = msg => { resize(msg.width, msg.height); };
  `,
    autoResize: "both",
    uiVisible: true,
    exposed: ({ main: { render, resize } }) => ({ render, resize }),
  },
};
