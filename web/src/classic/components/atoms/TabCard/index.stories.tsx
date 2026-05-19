import { Meta, Story } from "@storybook/react-vite";

import TabCard, { Props } from ".";

export default {
  title: "classic/atoms/TabCard",
  component: TabCard,
} as Meta;

const SampleBody = () => <div>Hello</div>;

export const Default: Story<Props> = args => (
  <TabCard {...args}>
    <SampleBody />
  </TabCard>
);
Default.args = {
  name: "Property",
};
