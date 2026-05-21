import { Meta, StoryObj } from "@storybook/react-vite";
import { fn } from "storybook/test";

import { Provider } from "../../storybook";

import Component from ".";

const meta: Meta<typeof Component> = {
  title: "classic/molecules/Visualizer/Widget/Timeline",
  component: Component,
  args: {
    onExtend: fn(),
    onVisibilityChange: fn(),
    onGetCredits: fn(),
  },
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Default: Story = {
  render: args => (
    <Provider>
      <Component {...args} />
    </Provider>
  ),
  args: {
    widget: {
      id: "",
      extended: {
        horizontally: false,
        vertically: false,
      },
    },
  },
};

export const Extended: Story = {
  render: args => (
    <Provider>
      <Component {...args} />
    </Provider>
  ),
  args: {
    widget: {
      id: "",
      extended: {
        horizontally: true,
        vertically: false,
      },
    },
  },
};
