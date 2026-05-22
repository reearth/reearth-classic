import { Meta, StoryObj } from "@storybook/react-vite";

import Component from ".";

const meta: Meta<typeof Component> = {
  title: "classic/molecules/Visualizer/Infobox/Frame",
  component: Component,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Default: Story = {
  args: {
    title: "hello",
  },
};

export const LongTitle: Story = {
  args: {
    title:
      "hellohellohellohellohellohellohellohellohellohhellohellohellohellohellohellohellohellohellohellohhellohellohellohellohellohellohellohellohellohellohhello",
  },
};

export const LongCJKTitle: Story = {
  args: {
    title:
      "こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。",
  },
};
