import { Meta, StoryObj } from "@storybook/react-vite";
import { Math as CesiumMath } from "cesium";

import { Provider } from "../../storybook";

import Button from ".";

const meta: Meta<typeof Button> = {
  title: "classic/molecules/Visualizer/Widget/Button",
  component: Button,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Button>;

export const Default: Story = {
  render: args => (
    <Provider>
      <Button {...args} />
    </Provider>
  ),
  args: {
    widget: {
      id: "",
      property: {
        buttons: [
          {
            id: "menu",
            buttonType: "menu",
            buttonTitle: "Menu",
            buttonPosition: "topleft",
            buttonStyle: "text",
          },
          {
            id: "google",
            buttonType: "link",
            buttonLink: "https://google.com",
            buttonTitle: "Google",
            buttonPosition: "topleft",
            buttonStyle: "text",
            buttonBgcolor: "red",
          },
          {
            id: "twitter",
            buttonType: "link",
            buttonLink: "https://twitter.com",
            buttonTitle: "Twitter",
            buttonPosition: "bottomright",
            buttonStyle: "text",
            buttonBgcolor: "#0ff",
          },
          {
            id: "hoge",
            buttonType: "camera",
            buttonPosition: "bottomright",
            buttonCamera: {
              lat: 0,
              lng: 0,
              height: 1000,
              fov: CesiumMath.toRadians(60),
              heading: 0,
              pitch: 0,
              roll: 0,
            },
            buttonTitle: "hoge",
          },
          {
            id: "menu2",
            buttonType: "menu",
            buttonIcon: "/sample.svg",
            buttonPosition: "bottomleft",
            buttonStyle: "icon",
          },
          {
            id: "menu3",
            buttonType: "menu",
            buttonTitle: "Menu",
            buttonIcon: "/sample.svg",
            buttonPosition: "bottomleft",
            buttonStyle: "texticon",
          },
        ],
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
    isBuilt: false,
    isEditable: false,
  },
};
