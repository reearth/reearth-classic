import { Meta, StoryObj } from "@storybook/react-vite";

import { contextEvents } from "../storybook";

import Component from ".";

const meta: Meta<typeof Component> = {
  component: Component,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Default: Story = {
  args: {
    widget: {
      id: "",
      property: {
        stories: [
          { layer: "a", title: "a" },
          { layer: "b", title: "b" },
          { layer: "c", title: "c" },
        ],
      },
    },
    context: { ...contextEvents },
    isBuilt: false,
    isEditable: false,
  },
};

export const AutoStart: Story = {
  args: {
    widget: {
      id: "",
      property: {
        stories: [
          { layer: "a", title: "a" },
          { layer: "b", title: "b" },
          { layer: "c", title: "c" },
        ],
        default: { autoStart: true },
      },
    },
    context: { ...contextEvents },
    isBuilt: false,
    isEditable: false,
  },
};
