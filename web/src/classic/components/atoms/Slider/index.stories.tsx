import { action } from "storybook/actions";
import { Meta } from "@storybook/react-vite";

import Slider from ".";

export default {
  title: "classic/atoms/Slider",
  component: Slider,
} as Meta;

export const Default = () => <Slider value={120} min={0} max={100} onChange={action("onchange")} />;
export const Frame = () => (
  <Slider value={120} min={0} max={100} onChange={action("onchange")} frame />
);
