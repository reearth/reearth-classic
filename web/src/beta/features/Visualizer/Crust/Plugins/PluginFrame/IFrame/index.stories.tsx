import { Meta, StoryObj } from "@storybook/react-vite";
import { useRef } from "react";

import Component, { Ref } from ".";

const meta: Meta<typeof Component> = {
  component: Component,
  argTypes: {
    onLoad: { action: "onLoad" },
    onMessage: { action: "onMessage" },
  },
  // parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Component>;

const DefaultRenderer = (args: Story["args"]) => {
  const ref = useRef<Ref>(null);
  const postMessage = () => {
    ref.current?.postMessage({ foo: new Date().toISOString() });
  };
  return (
    <div style={{ background: "#fff" }}>
      <Component {...args} ref={ref} />
      <p>
        <button onClick={postMessage}>postMessage</button>
      </p>
    </div>
  );
};

export const Default: Story = {
  render: args => <DefaultRenderer {...args} />,
  args: {
    visible: true,
    iFrameProps: {
      style: {
        width: "400px",
        height: "300px",
      },
    },
    html: `<h1>iframe</h1><script>
  window.addEventListener("message", ev => {
    if (ev.source !== parent) return;
    const p = document.createElement("p");
    p.textContent = JSON.stringify(ev.data);
    document.body.appendChild(p);
    parent.postMessage(ev.data, "*");
  });
  parent.postMessage("loaded", "*");
</script>`,
  },
};
