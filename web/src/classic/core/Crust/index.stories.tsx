import { Meta, StoryObj } from "@storybook/react-vite";
import { useRef } from "react";

import { engine } from "../engines/Cesium";
import Map, { Engine } from "../Map";

import { MapRef } from "./types";

import Component, { Props } from ".";

const meta: Meta<Props> = {
  component: Component,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;

const CesiumRenderer = (args: Props & { engine: Engine }) => {
  const ref = useRef<MapRef>(null);
  return (
    <>
      <Map engine="a" engines={{ a: engine }} ref={ref} />
      <Component {...args} mapRef={ref} />
    </>
  );
};

export const Cesium: StoryObj<Props & { engine: Engine }> = {
  render: args => <CesiumRenderer {...args} />,
  args: {
    engine: engine,
  },
};
