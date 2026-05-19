import { Meta, StoryObj } from "@storybook/react-vite";

import Navigator from ".";

const meta: Meta<typeof Navigator> = {
  title: "classic/atoms/Navigator",
  component: Navigator,
};
export default meta;
type Story = StoryObj<typeof Navigator>;

export const Normal: Story = {
  render: () => <Navigator degree={0} />,
};
