import { action } from "storybook/actions";
import { Meta, StoryObj } from "@storybook/react-vite";

import ProjectCreationModal from ".";

const meta: Meta<typeof ProjectCreationModal> = {
  title: "classic/molecules/ProjectList/ProjectCreationModal",
  component: ProjectCreationModal,
};
export default meta;
type Story = StoryObj<typeof ProjectCreationModal>;

export const Default: Story = {
  render: args => (
    <ProjectCreationModal {...args} open onClose={action("onClose")} onSubmit={action("onSubmit")} />
  ),
};
