import { action } from "storybook/actions";
import { Meta, Story } from "@storybook/react-vite";

import ProjectCreationModal, { Props } from ".";

export default {
  title: "classic/molecules/ProjectList/ProjectCreationModal",
  component: ProjectCreationModal,
} as Meta;

export const Default: Story<Props> = args => (
  <ProjectCreationModal {...args} open onClose={action("onClose")} onSubmit={action("onSubmit")} />
);
