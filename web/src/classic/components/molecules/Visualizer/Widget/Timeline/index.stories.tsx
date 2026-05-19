import { Meta, StoryObj } from "@storybook/react-vite";

import { Provider } from "../../storybook";

import Component from ".";

const meta: Meta<typeof Component> = {
  title: "classic/molecules/Visualizer/Widget/Timeline",
  component: Component,
  parameters: { actions: { argTypesRegex: "^on.*" } },
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
