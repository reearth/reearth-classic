import { Meta, StoryObj } from "@storybook/react-vite";

import Component, { PluginAccordionProps } from ".";

const meta: Meta<PluginAccordionProps> = {
  title: "classic/molecules/Settings/Project/PluginAccordion",
  component: Component,
};
export default meta;
type Story = StoryObj<PluginAccordionProps>;

export const Default: Story = {
  args: {
    plugins: [
      {
        thumbnailUrl: `/sample.svg`,
        title: "Sample",
        isInstalled: true,
        bodyMarkdown: "# Hoge ## Hoge",
        author: "reearth",
        pluginId: "id1",
      },
      {
        thumbnailUrl: `/sample.svg`,
        title: "Sample2",
        isInstalled: false,
        bodyMarkdown: "# Fuga ## Fuga",
        author: "reearth",
        pluginId: "id2",
      },
    ],
  },
};
