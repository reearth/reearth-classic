import { action } from "storybook/actions";
import { Meta } from "@storybook/react-vite";

import WorkSpaceCreationModal from ".";

export default {
  title: "classic/molecules/Common/WorkSpaceCreationModal",
  component: WorkSpaceCreationModal,
} as Meta;

export const Default = () => (
  <WorkSpaceCreationModal open onClose={action("onClose")} onSubmit={action("onSubmit")} />
);
