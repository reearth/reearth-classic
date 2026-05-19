import { Meta, StoryObj } from "@storybook/react-vite";

import Component from ".";

const meta: Meta<typeof Component> = {
  title: "classic/atoms/SearchBar",
  component: Component,
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Default: Story = {
  args: {
    iconPos: "right",
    placeHolder: "search plugins",
  },
};
