import { Meta, StoryObj } from "@storybook/react-vite";

import Component from ".";

const meta: Meta<typeof Component> = {
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
      property: { text: "hogehoge" },
    },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
  },
};
