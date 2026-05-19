import { Meta, StoryObj } from "@storybook/react-vite";

import { Provider } from "../storybook";

import Component from ".";

const meta: Meta<typeof Component> = {
  title: "classic/molecules/Visualizer/Block",
  component: Component,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Default: Story = {
  args: {
    block: {
      id: "",
      pluginId: "reearth",
      extensionId: "textblock",
      property: { default: { text: "hogehoge" } },
    },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
  },
};

export const Plugin: Story = {
  render: args => (
    <Provider>
      <div style={{ background: "#fff" }}>
        <Component {...args} />
      </div>
    </Provider>
  ),
  args: {
    block: {
      id: "",
      pluginId: "plugins",
      extensionId: "block",
      property: {
        location: { lat: 0, lng: 100 },
      },
    },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
    pluginBaseUrl: "",
  },
};
