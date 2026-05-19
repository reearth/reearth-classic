import { Meta, StoryObj } from "@storybook/react-vite";

import Component from ".";

const meta: Meta<typeof Component> = {
  title: "classic/atoms/Box",
  component: Component,
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Margin: Story = {
  render: () => (
    <Component m="l">
      <div>Margin</div>
    </Component>
  ),
};

export const Padding: Story = {
  render: () => (
    <Component p="l">
      <div>Padding</div>
    </Component>
  ),
};

export const Border: Story = {
  render: () => (
    <Component border="solid 5px red">
      <div>Border</div>
    </Component>
  ),
};

export const Bg: Story = {
  render: () => (
    <Component bg="red">
      <div>Bg</div>
    </Component>
  ),
};
