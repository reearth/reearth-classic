import { Meta, StoryObj } from "@storybook/react-vite";
import { useRef } from "react";

import { Map, engines, Engine, InteractionModeType, INTERACTION_MODES } from "@reearth/core";

import { MapRef } from "./types";

import Component, { Props } from ".";

const meta: Meta<Props> = {
  component: Component,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;

const CesiumRenderer = (args: Props & { engine: Engine; interactionMode: InteractionModeType }) => {
  const ref = useRef<MapRef>(null);
  return (
    <>
      <Map
        engine="a"
        engines={{ a: engines.cesium }}
        ref={ref}
        featureFlags={INTERACTION_MODES[args.interactionMode]}
      />
      <Component {...args} mapRef={ref} />
    </>
  );
};

export const Cesium: StoryObj<Props & { engine: Engine; interactionMode: InteractionModeType }> = {
  render: args => <CesiumRenderer {...args} />,
  args: {
    engine: engines.cesium,
  },
};
