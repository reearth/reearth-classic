import { Meta, StoryObj } from "@storybook/react-vite";

import { contextEvents } from "../../storybook";

import Component from ".";

const meta: Meta<typeof Component> = {
  component: Component,
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Default: Story = {
  args: {
    widget: {
      id: "",
      extended: {
        horizontally: false,
        vertically: false,
      },
    },
    context: { ...contextEvents },
  },
};

export const Extended: Story = {
  args: {
    widget: {
      id: "",
      extended: {
        horizontally: true,
        vertically: false,
      },
    },
    context: { ...contextEvents },
  },
};
