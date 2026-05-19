import { Meta, StoryObj } from "@storybook/react-vite";

import { V, location, SceneProperty } from "../storybook";

import Component, { Props } from ".";

const meta: Meta<typeof Component> = {
  title: "classic/molecules/Visualizer/Engine/Cesium/Model",
  component: Component,
};
export default meta;

export const Default: StoryObj<Props & { sceneProperty?: SceneProperty }> = {
  render: args => (
    <V property={args.sceneProperty}>
      <Component {...args} />
    </V>
  ),
  args: {
    sceneProperty: {
      timeline: {
        animation: true,
      },
    },
    layer: {
      id: "",
      isVisible: true,
      property: {
        default: {
          location,
          scale: 1000,
          heading: 130,
          model: `/BoxAnimated.glb`,
        },
      },
    },
  },
};

export const Appearance: StoryObj<Props & { sceneProperty?: SceneProperty }> = {
  render: args => (
    <V property={args.sceneProperty}>
      <Component {...args} />
    </V>
  ),
  args: {
    sceneProperty: {
      atmosphere: {
        enable_shadows: true,
      },
    },
    layer: {
      id: "",
      isVisible: true,
      property: {
        default: {
          location,
          scale: 1000,
          model: `/BoxAnimated.glb`,
          animation: false,
        },
        appearance: {
          shadows: "enabled",
          color: "red",
          colorBlend: "mix",
          colorBlendAmount: 0.5,
          silhouette: true,
          silhouetteColor: "yellow",
          silhouetteSize: 10,
        },
      },
    },
  },
};
