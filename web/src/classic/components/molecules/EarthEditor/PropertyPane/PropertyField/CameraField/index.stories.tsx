import { Meta, StoryObj } from "@storybook/react-vite";

import CameraField from ".";

const meta: Meta<typeof CameraField> = {
  title: "classic/molecules/EarthEditor/PropertyPane/PropertyField/CameraField",
  component: CameraField,
};
export default meta;
type Story = StoryObj<typeof CameraField>;

export const HasNoCamera: Story = {
  args: {},
};

export const HasCamera: Story = {
  render: args => (
    <CameraField
      {...args}
      value={{
        lat: 0,
        lng: 0,
        height: 10 ** 8,
        heading: 0,
        pitch: 0,
        roll: 0,
        fov: 1,
      }}
    />
  ),
};

export const OnlyPose: Story = {
  render: args => (
    <CameraField
      {...args}
      onlyPose
      value={{
        lat: 0,
        lng: 0,
        height: 10 ** 8,
        heading: 0,
        pitch: 0,
        roll: 0,
        fov: 1,
      }}
    />
  ),
};
