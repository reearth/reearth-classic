import { Meta, StoryObj } from "@storybook/react-vite";

import Component from ".";

const meta: Meta<typeof Component> = {
  title: "classic/molecules/EarthEditor/PropertyPane/WidgetAlignSystemToggle",
  component: Component,
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Default: Story = {
  args: {
    checked: false,
  },
};
