import { Meta, StoryObj } from "@storybook/react-vite";

import { Provider } from "../../storybook";

import Component from ".";

const meta: Meta<typeof Component> = {
  title: "classic/molecules/Visualizer/Widget/Storytelling",
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
      property: {
        stories: [
          { layer: "a", title: "a" },
          { layer: "b", title: "b" },
          { layer: "c", title: "c" },
        ],
      },
    },
    isBuilt: false,
    isEditable: false,
  },
};

export const AutoStart: Story = {
  render: args => (
    <Provider>
      <Component {...args} />
    </Provider>
  ),
  args: {
    widget: {
      id: "",
      property: {
        stories: [
          { layer: "a", title: "a" },
          { layer: "b", title: "b" },
          { layer: "c", title: "c" },
        ],
        default: { autoStart: true },
      },
    },
    isBuilt: false,
    isEditable: false,
  },
};
