import { action } from "storybook/actions";
import { Meta, StoryObj } from "@storybook/react-vite";

import ProjectTypeSelectionModal from ".";

const meta: Meta<typeof ProjectTypeSelectionModal> = {
  title: "classic/molecules/common/ProjectTypeSelectionModal",
  component: ProjectTypeSelectionModal,
};
export default meta;
type Story = StoryObj<typeof ProjectTypeSelectionModal>;

export const Default: Story = {
  render: args => (
    <ProjectTypeSelectionModal
      {...args}
      open
      onClose={action("onClose")}
      onSubmit={action("onSubmit")}
    />
  ),
};
