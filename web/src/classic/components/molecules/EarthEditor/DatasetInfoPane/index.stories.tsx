import { Meta, StoryObj } from "@storybook/react-vite";

import DatasetInfoPane from ".";

const meta: Meta<typeof DatasetInfoPane> = {
  title: "classic/molecules/EarthEditor/DatasetInfoPane",
  component: DatasetInfoPane,
};
export default meta;
type Story = StoryObj<typeof DatasetInfoPane>;

export const Default: Story = {
  args: {},
};
