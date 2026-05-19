import { Meta, StoryObj } from "@storybook/react-vite";

import Video from ".";

const meta: Meta<typeof Video> = {
  title: "classic/molecules/Visualizer/Block/Video",
  component: Video,
  argTypes: {
    onClick: { action: "onClick" },
    onChange: { action: "onChange" },
  },
};
export default meta;
type Story = StoryObj<typeof Video>;

export const Default: Story = {
  args: {
    block: { id: "", property: { default: { url: "https://www.youtube.com/watch?v=oUFJJNQGwhk" } } },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
  },
};

export const Title: Story = {
  args: {
    block: {
      id: "",
      property: { default: { url: "https://www.youtube.com/watch?v=oUFJJNQGwhk", title: "Video" } },
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
