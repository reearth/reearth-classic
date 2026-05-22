import { Meta, StoryObj } from "@storybook/react-vite";

import Component from ".";

const meta: Meta<typeof Component> = {
  title: "classic/molecules/Visualizer/Infobox",
  component: Component,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Component>;

const baseArgs = {
  blocks: [
    {
      id: "a",
      pluginId: "reearth",
      extensionId: "textblock",
      property: {
        default: {
          text: "# Hello",
          markdown: true,
        },
      },
    },
    {
      id: "b",
      pluginId: "reearth",
      extensionId: "textblock",
      property: {
        default: {
          text: "# World!",
          markdown: true,
        },
      },
    },
  ],
  layer: {
    id: "z",
    property: {
      default: {
        title: "Hoge",
        size: "small",
      },
    },
  },
  selectedBlockId: undefined,
  title: "Name",
  infoboxKey: "",
  isBuilt: false,
  isEditable: false,
  visible: true,
};

export const Default: Story = {
  args: baseArgs,
};

export const Large: Story = {
  args: {
    ...baseArgs,
    layer: {
      id: "z",
      property: {
        ...baseArgs.layer?.property,
        default: {
          ...baseArgs.layer?.property?.default,
          size: "large",
        },
      },
    },
  },
};
