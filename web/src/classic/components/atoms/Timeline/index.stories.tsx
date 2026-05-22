import { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";

import Timeline from ".";

const meta: Meta<typeof Timeline> = {
  title: "classic/atoms/Timeline/Timeline",
  component: Timeline,
};
export default meta;
type Story = StoryObj<typeof Timeline>;

export const Normal: Story = {
  render: () => (
    <Timeline
      currentTime={new Date("2022-06-30T12:20:00.000").getTime()}
      range={{
        start: new Date("2022-06-30T21:00:00.000").getTime(),
        end: new Date("2022-07-03T12:21:21.221").getTime(),
      }}
      isOpened={true}
    />
  ),
};

export const DefaultRange: Story = {
  render: () => (
    <Timeline
      // Forward a hour
      currentTime={Date.now() + 3600000}
      isOpened={true}
    />
  ),
};

const MovableRenderer = () => {
  // Forward a hour
  const [currentTime, setCurrentTime] = useState(() => Date.now() + 3600000);
  const [isOpened, setIsOpened] = useState(false);
  return (
    <Timeline
      currentTime={currentTime}
      onClick={setCurrentTime}
      onDrag={setCurrentTime}
      isOpened={isOpened}
      onOpen={() => setIsOpened(true)}
      onClose={() => setIsOpened(false)}
    />
  );
};

export const Movable: Story = {
  render: () => <MovableRenderer />,
};
