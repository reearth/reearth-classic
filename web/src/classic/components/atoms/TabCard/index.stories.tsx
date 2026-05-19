import { Meta, StoryObj } from "@storybook/react-vite";

import TabCard from ".";

const meta: Meta<typeof TabCard> = {
  title: "classic/atoms/TabCard",
  component: TabCard,
};
export default meta;
type Story = StoryObj<typeof TabCard>;

const SampleBody = () => <div>Hello</div>;

export const Default: Story = {
  render: args => (
    <TabCard {...args}>
      <SampleBody />
    </TabCard>
  ),
  args: {
    name: "Property",
  },
};
