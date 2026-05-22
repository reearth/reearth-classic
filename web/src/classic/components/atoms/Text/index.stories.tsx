import { Meta, StoryObj } from "@storybook/react-vite";

import Text from ".";

const meta: Meta<typeof Text> = {
  title: "classic/atoms/Text",
  component: Text,
};
export default meta;
type Story = StoryObj<typeof Text>;

export const LBold: Story = {
  render: () => (
    <Text size="l" weight="bold">
      LBold
    </Text>
  ),
};

export const MParagraph: Story = {
  render: () => (
    <Text size="m" isParagraph>
      MParagraph
    </Text>
  ),
};

export const MRegularRed: Story = {
  render: () => (
    <Text size="m" color="red">
      MRegular red
    </Text>
  ),
};

export const MRegularRedBgBlue: Story = {
  render: () => (
    <Text size="m" color="red" otherProperties={{ background: "blue", border: "solid 2px red" }}>
      MRegular red blue
    </Text>
  ),
};
