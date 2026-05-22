import { Meta, StoryObj } from "@storybook/react-vite";

import Component from ".";

const meta: Meta<typeof Component> = {
  component: Component,
};
export default meta;
type Story = StoryObj<typeof Component>;

const ExampleDiv = () => (
  <>
    <div>hoge</div>
    <div>fuga</div>
  </>
);

export const SpaceBetween: Story = {
  render: args => (
    <Component {...args}>
      <ExampleDiv />
    </Component>
  ),
  args: {
    justify: "space-between",
  },
};

export const GapChildren: Story = {
  render: args => (
    <Component {...args}>
      <div>hoge</div>
      <div>fuga</div>
    </Component>
  ),
  args: {
    gap: "20px",
  },
};

export const DirectionVertical: Story = {
  render: args => (
    <Component {...args}>
      <ExampleDiv />
    </Component>
  ),
  args: {
    direction: "column",
  },
};
