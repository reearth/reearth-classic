import { Meta, StoryObj } from "@storybook/react-vite";

import Navigator from ".";

const meta: Meta<typeof Navigator> = {
  component: Navigator,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Navigator>;

export const Normal: Story = {
  render: ({ ...args }) => <Navigator {...args} degree={0} />,
};
