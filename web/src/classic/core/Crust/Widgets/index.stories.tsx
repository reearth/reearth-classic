import { Meta, StoryObj } from "@storybook/react-vite";

import { contextEvents } from "./Widget/storybook";

import Component from ".";

const meta: Meta<typeof Component> = {
  component: Component,
  argTypes: {
    context: { control: false },
  },
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Default: Story = {
  render: args => (
    <div style={{ width: "100%", height: "100%" }}>
      <Component {...args} />
    </div>
  ),
  args: {
    alignSystem: {
      inner: {
        center: {
          top: {
            align: "centered",
            widgets: [
              {
                id: "w",
                pluginId: "reearth",
                extensionId: "button",
                property: {
                  default: {
                    buttonTitle: "Button 1",
                    buttonColor: "white",
                    buttonBgcolor: "red",
                  },
                },
              },
              {
                id: "unknown",
                pluginId: "unknown",
                extensionId: "unknown",
              },
            ],
          },
        },
      },
    },
    context: {
      ...contextEvents,
    },
    editing: false,
    isBuilt: false,
    isEditable: false,
    renderWidget: ({ widget }) => <p>{widget.id}</p>,
  },
};
