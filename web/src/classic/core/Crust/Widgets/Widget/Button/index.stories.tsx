import { Meta, StoryObj } from "@storybook/react-vite";
import { Math as CesiumMath } from "cesium";

import { contextEvents } from "../storybook";

import Button from ".";

const meta: Meta<typeof Button> = {
  component: Button,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Button>;

export const Default: Story = {
  args: {
    widget: {
      id: "",
      property: {
        menu: [
          {
            id: "hoge",
            menuTitle: "Hoge",
            menuType: "camera",
            menuCamera: {
              lat: 0,
              lng: 0,
              height: 1000,
              fov: CesiumMath.toRadians(60),
              heading: 0,
              pitch: 0,
              roll: 0,
            },
          },
          {
            id: "hoge",
            menuType: "border",
          },
          {
            id: "GitHub",
            menuType: "link",
            menuTitle: "GitHub",
            menuLink: "https://github.com",
          },
        ],
      },
    },
    context: { ...contextEvents },
    isBuilt: false,
    isEditable: false,
  },
};
