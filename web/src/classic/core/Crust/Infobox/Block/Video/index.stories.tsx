import { Meta, StoryObj } from "@storybook/react-vite";

import Video from ".";

const meta: Meta<typeof Video> = {
  component: Video,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Video>;

export const Default: Story = {
  args: {
    block: { id: "", property: { url: "https://www.youtube.com/watch?v=oUFJJNQGwhk" } },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
  },
};

export const Title: Story = {
  args: {
    block: {
      id: "",
      property: { url: "https://www.youtube.com/watch?v=oUFJJNQGwhk", title: "Video" },
    },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
  },
};

export const NoVideo: Story = {
  args: {
    isSelected: false,
    isBuilt: false,
    isEditable: true,
  },
};
