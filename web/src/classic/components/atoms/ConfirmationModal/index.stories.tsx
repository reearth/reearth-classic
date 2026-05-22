import { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";

import ConfirmationModal, { Props } from ".";

const meta: Meta<typeof ConfirmationModal> = {
  title: "classic/atoms/Modal/ConfirmationModal",
  component: ConfirmationModal,
};
export default meta;
type Story = StoryObj<typeof ConfirmationModal>;

const DefaultRenderer = (args: Story["args"]) => {
  const [isOpen, setOpen] = useState(false);
  return (
    <>
      <ConfirmationModal
        {...(args as Props)}
        isOpen={isOpen}
        onClose={() => setOpen(false)}
        onCancel={() => setOpen(false)}
      />
      <button onClick={() => setOpen(true)}>click</button>
    </>
  );
};

export const Default: Story = {
  render: args => <DefaultRenderer {...args} />,
  args: {
    body: <div>Are you sure to delete this</div>,
    title: "Delete Sample",
    onProceed: () => console.log("Proceed"),
  },
};
