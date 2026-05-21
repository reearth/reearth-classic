import { Meta, StoryObj } from "@storybook/react-vite";
import { fn } from "storybook/test";

import Component from ".";

const meta: Meta<typeof Component> = {
  component: Component,
  args: {
    onMaskClick: fn(),
    onClick: fn(),
    onClickAway: fn(),
    onEnter: fn(),
    onEntered: fn(),
    onExit: fn(),
    onExited: fn(),
    onClose: fn(),
  },
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Default: Story = {
  args: {
    title: "hello",
    visible: true,
  },
};

export const LongTitle: Story = {
  args: {
    title:
      "hellohellohellohellohellohellohellohellohellohhellohellohellohellohellohellohellohellohellohellohhellohellohellohellohellohellohellohellohellohellohhello",
    visible: true,
  },
};

export const LongCJKTitle: Story = {
  args: {
    title:
      "こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。こんにちは。",
    visible: true,
  },
};
