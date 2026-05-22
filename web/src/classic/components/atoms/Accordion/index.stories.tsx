import { Meta, StoryObj } from "@storybook/react-vite";

import Component from ".";

const meta: Meta<typeof Component> = {
  title: "classic/atoms/accordion",
  component: Component,
};
export default meta;
type Story = StoryObj<typeof Component>;

const SampleHeading = <div style={{ color: "white", background: "black" }}>heading</div>;

const SampleContent = (
  <div style={{ color: "white", background: "black", padding: "20px" }}>hoge</div>
);

export const Default: Story = {
  args: {
    items: [
      {
        id: "hoge",
        heading: SampleHeading,
        content: SampleContent,
      },
      {
        id: "hogefuga",
        heading: SampleHeading,
        content: SampleContent,
      },
    ],
  },
};
