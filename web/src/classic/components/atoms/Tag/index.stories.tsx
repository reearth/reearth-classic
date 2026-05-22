import { Meta, StoryObj } from "@storybook/react-vite";

import Tag from ".";

const meta: Meta<typeof Tag> = {
  title: "classic/atoms/Tag",
  component: Tag,
};
export default meta;
type Story = StoryObj<typeof Tag>;

export const Default: Story = {
  args: {
    icon: "cancel",
    onRemove: () => console.log("detatch!"),
    id: "tag",
    label: "tag",
  },
};

export const Remove: Story = {
  args: {
    icon: "bin",
    onRemove: () => console.log("remove!"),
    id: "tag",
    label: "tag",
  },
};
