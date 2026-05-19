import { Meta, StoryObj } from "@storybook/react-vite";

import { Provider } from "../storybook";

import Component from ".";

const meta: Meta<typeof Component> = {
  component: Component,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Default: Story = {
  render: args => (
    <Provider>
      <div style={{ background: "#fff" }}>
        <Component {...args} />
      </div>
    </Provider>
  ),
  args: {
    pluginId: "plugins",
    extensionId: "plugin",
    pluginBaseUrl: "",
    visible: true,
  },
};

export const Headless: Story = {
  render: args => (
    <Provider>
      <Component {...args} />
    </Provider>
  ),
  args: {
    sourceCode: `console.log("hello world");`,
    visible: true,
  },
};
