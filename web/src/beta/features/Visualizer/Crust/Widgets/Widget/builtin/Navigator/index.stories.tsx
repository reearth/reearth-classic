import { Meta, StoryObj } from "@storybook/react-vite";

import { contextEvents } from "../../storybook";

import Component from ".";

const meta: Meta<typeof Component> = {
  component: Component,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Default: Story = {
  args: {
    context: { ...contextEvents },
  },
};
